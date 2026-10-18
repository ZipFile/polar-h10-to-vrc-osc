package osc_testing

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/hypebeast/go-osc/osc"
)

type TestDispatcher struct {
	Messages chan<- *osc.Message
	Bundles  chan<- *osc.Bundle
}

var _ osc.Dispatcher = (*TestDispatcher)(nil)

func (d *TestDispatcher) Dispatch(packet osc.Packet) {
	switch p := packet.(type) {
	case *osc.Message:
		d.Messages <- p
	case *osc.Bundle:
		d.Bundles <- p
	}
}

func NewClientServerPair(
	t *testing.T,
	messages chan<- *osc.Message,
	bundles chan<- *osc.Bundle,
) (*osc.Client, *osc.Server) {
	port := rand.IntN(16384) + 49152
	client := osc.NewClient("127.0.0.1", port)
	server := &osc.Server{
		Addr:       fmt.Sprintf("127.0.0.1:%d", port),
		Dispatcher: &TestDispatcher{messages, bundles},
	}

	t.Cleanup(func() {
		_ = server.CloseConnection()
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			require.ErrorContains(t, err, "use of closed network connection")
		}
	}()

	return client, server
}
