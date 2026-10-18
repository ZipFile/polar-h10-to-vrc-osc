package vrc_osc

import (
	"testing"
	"time"

	osc_testing "github.com/ZipFile/polar-h10-to-osc/internal/osc/testing"
	"github.com/go-openapi/testify/v2/require"
	"github.com/hypebeast/go-osc/osc"
)

func assertMsg(t *testing.T, msg *osc.Message, name string, expected any) {
	require.Len(t, msg.Arguments, 1)
	require.Equal(t, msg.Address, "/avatar/parameters/"+name)
	require.Equal(t, expected, msg.Arguments[0])
}

func TestSender_SendEmpty(t *testing.T) {
	var sender Sender

	require.NoError(t, sender.send())
	require.NoError(t, sender.send(nil))
	require.NoError(t, sender.send(nil, nil))
}

func TestSender(t *testing.T) {
	messages := make(chan *osc.Message)
	bundles := make(chan *osc.Bundle)
	client, _ := osc_testing.NewClientServerPair(t, messages, bundles)
	params := PresetDefault
	sender := Sender{*params, client}

	t.Run("one", func(t *testing.T) {
		require.NoError(t, sender.send(fmt("test", int32(123))))

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no message received")
		case msg := <-messages:
			assertMsg(t, msg, "test", int32(123))
		}
	})

	t.Run("bundle", func(t *testing.T) {
		err := sender.send(
			nil,
			fmt("test", int32(123)),
			fmt("", "nope"),
			fmt("tset", int32(321)),
			nil,
		)

		require.NoError(t, err)

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no bundles received")
		case bundle := <-bundles:
			require.Equal(t, uint64(1), bundle.Timetag.TimeTag())
			require.Len(t, bundle.Messages, 2)
			assertMsg(t, bundle.Messages[0], "test", int32(123))
			assertMsg(t, bundle.Messages[1], "tset", int32(321))
		}
	})

	t.Run("zero", func(t *testing.T) {
		err := sender.Zero()

		require.NoError(t, err)

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no bundles received")
		case bundle := <-bundles:
			require.Len(t, bundle.Messages, 6)
			assertMsg(t, bundle.Messages[0], params.StatusConnected, false)
			assertMsg(t, bundle.Messages[1], params.StatusActive, false)
			assertMsg(t, bundle.Messages[2], params.StatusBeat, false)
			assertMsg(t, bundle.Messages[3], params.ValueAbsolute, float32(0))
			assertMsg(t, bundle.Messages[4], params.ValuePercent, float32(0))
			assertMsg(t, bundle.Messages[5], params.ValuePercentFull, float32(0))
		}
	})

	t.Run("value", func(t *testing.T) {
		require.NoError(t, sender.Value(75))

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no bundle received")
		case bundle := <-bundles:
			require.Len(t, bundle.Messages, 3)
			assertMsg(t, bundle.Messages[0], params.ValueAbsolute, int32(75))
			assertMsg(t, bundle.Messages[1], params.ValuePercent, float32(0.5))
			assertMsg(t, bundle.Messages[2], params.ValuePercentFull, float32(0))
		}
	})

	var btable = []struct {
		name  string
		fn    func(bool) error
		param string
	}{
		{"active", sender.ActiveStatus, params.StatusActive},
		{"connected", sender.ConnectedStatus, params.StatusConnected},
		{"beat", sender.BeatStatus, params.StatusBeat},
	}

	for _, bt := range btable {
		t.Run(bt.name, func(t *testing.T) {
			require.NoError(t, bt.fn(true))

			select {
			case <-time.After(100 * time.Millisecond):
				t.Error("no message received")
			case msg := <-messages:
				assertMsg(t, msg, bt.param, true)
			}
		})
	}
}
