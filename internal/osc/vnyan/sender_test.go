package vnyan_osc

import (
	"testing"
	"time"

	osc_testing "github.com/ZipFile/polar-h10-to-osc/internal/osc/testing"
	"github.com/go-openapi/testify/v2/require"
	"github.com/hypebeast/go-osc/osc"
)

func assertMsg(t *testing.T, msg *osc.Message, expected int) {
	require.Len(t, msg.Arguments, 2)
	require.Equal(t, msg.Address, "/VNyan/Param/Float")
	require.Equal(t, "heartrate", msg.Arguments[0])
	require.Equal(t, float32(expected), msg.Arguments[1])
}

func TestSender_Value(t *testing.T) {
	messages := make(chan *osc.Message)
	client, _ := osc_testing.NewClientServerPair(t, messages, nil)
	sender := Sender{client}

	t.Run("zero", func(t *testing.T) {
		require.NoError(t, sender.Zero())

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no message received")
		case msg := <-messages:
			assertMsg(t, msg, 0)
		}
	})

	t.Run("value", func(t *testing.T) {
		require.NoError(t, sender.Value(72))

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no message received")
		case msg := <-messages:
			assertMsg(t, msg, 72)
		}
	})
}

func TestSender_BoolNoops(t *testing.T) {
	var s Sender
	table := []struct {
		name string
		fn   func(bool) error
	}{
		{"ActiveStatus", s.ActiveStatus},
		{"ConnectedStatus", s.ConnectedStatus},
		{"BeatStatus", s.BeatStatus},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(false); err != nil {
				t.Errorf("%s(false) = %v", tt.name, err)
			}
		})
	}
}
