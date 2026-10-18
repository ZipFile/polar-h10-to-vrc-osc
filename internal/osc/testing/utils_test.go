package osc_testing

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/hypebeast/go-osc/osc"
)

func TestNewClientServerPair(t *testing.T) {
	messages := make(chan *osc.Message)
	bundles := make(chan *osc.Bundle)
	client, server := NewClientServerPair(t, messages, bundles)

	require.NotNil(t, client)
	require.NotNil(t, server)

	t.Run("message", func(t *testing.T) {
		expected := osc.NewMessage("/test", "test", int32(123))
		err := client.Send(expected)

		require.NoError(t, err)

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no message received")
		case msg := <-messages:
			require.Equal(t, expected, msg)
		}
	})

	t.Run("bundle", func(t *testing.T) {
		expected := osc.NewBundle(time.Date(2026, 10, 18, 19, 0, 0, 0, time.Local))

		require.NoError(t, expected.Append(osc.NewMessage("/test", "test", int32(123))))
		require.NoError(t, expected.Append(osc.NewMessage("/tset", int32(321), "tset")))

		require.NoError(t, client.Send(expected))

		select {
		case <-time.After(100 * time.Millisecond):
			t.Error("no bundle received")
		case bundle := <-bundles:
			require.Equal(t, expected, bundle)
		}
	})
}
