package common

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

type testSender struct {
	zeroCalls            int
	valueCalls           int
	activeStatusCalls    int
	connectedStatusCalls int
	beatStatusCalls      int

	lastValue  int
	lastStatus bool
	err        error
}

func (m *testSender) Zero() error {
	m.zeroCalls++
	return m.err
}

func (m *testSender) Value(v int) error {
	m.valueCalls++
	m.lastValue = v
	return m.err
}

func (m *testSender) ActiveStatus(v bool) error {
	m.activeStatusCalls++
	m.lastStatus = v
	return m.err
}

func (m *testSender) ConnectedStatus(v bool) error {
	m.connectedStatusCalls++
	m.lastStatus = v
	return m.err
}

func (m *testSender) BeatStatus(v bool) error {
	m.beatStatusCalls++
	m.lastStatus = v
	return m.err
}

func TestCombinedSender(t *testing.T) {
	m1 := &testSender{}
	m2 := &testSender{}
	combined := CombinedSender{m1, m2}

	t.Run("Zero", func(t *testing.T) {
		require.NoError(t, combined.Zero())
		require.Equal(t, 1, m1.zeroCalls)
		require.Equal(t, 1, m2.zeroCalls)
	})

	t.Run("Value", func(t *testing.T) {
		require.NoError(t, combined.Value(85))
		require.Equal(t, 1, m1.valueCalls)
		require.Equal(t, 85, m1.lastValue)
		require.Equal(t, 1, m2.valueCalls)
		require.Equal(t, 85, m2.lastValue)
	})

	t.Run("ActiveStatus", func(t *testing.T) {
		require.NoError(t, combined.ActiveStatus(true))
		require.Equal(t, 1, m1.activeStatusCalls)
		require.True(t, m1.lastStatus)
		require.Equal(t, 1, m2.activeStatusCalls)
		require.True(t, m2.lastStatus)
	})

	t.Run("ConnectedStatus", func(t *testing.T) {
		require.NoError(t, combined.ConnectedStatus(true))
		require.Equal(t, 1, m1.connectedStatusCalls)
		require.True(t, m1.lastStatus)
		require.Equal(t, 1, m2.connectedStatusCalls)
		require.True(t, m2.lastStatus)
	})

	t.Run("BeatStatus", func(t *testing.T) {
		require.NoError(t, combined.BeatStatus(true))
		require.Equal(t, 1, m1.beatStatusCalls)
		require.True(t, m1.lastStatus)
		require.Equal(t, 1, m2.beatStatusCalls)
		require.True(t, m2.lastStatus)
	})

	t.Run("Errors", func(t *testing.T) {
		errMock := &testSender{err: require.ErrTest}
		c := CombinedSender{m1, errMock}
		err := c.Value(90)
		require.ErrorIs(t, err, require.ErrTest)
	})
}
