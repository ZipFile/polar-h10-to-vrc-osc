package common

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

var _ Sender = (*NoopSender)(nil)
var _ Sender = &NoopSender{}

func TestNoopSender(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var s *NoopSender
		require.NoError(t, s.e())
	})

	t.Run("ok", func(t *testing.T) {
		s := &NoopSender{}
		require.NoError(t, s.e())
	})

	t.Run("err", func(t *testing.T) {
		s := &NoopSender{Err: require.ErrTest}

		require.ErrorIs(t, s.Zero(), require.ErrTest)
		require.ErrorIs(t, s.Value(120), require.ErrTest)
		require.ErrorIs(t, s.ActiveStatus(true), require.ErrTest)
		require.ErrorIs(t, s.ConnectedStatus(true), require.ErrTest)
		require.ErrorIs(t, s.BeatStatus(true), require.ErrTest)
	})
}
