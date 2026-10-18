package vrc_osc

import (
	"net"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/hashicorp/mdns"
)

func TestInfo_MDNSParse(t *testing.T) {
	info := Info{}

	const name = "VRChat-Client-TEST"

	t.Run("junk", func(t *testing.T) {
		ch := make(chan *mdns.ServiceEntry)

		go func() {
			ch <- nil
			ch <- &mdns.ServiceEntry{}
			ch <- &mdns.ServiceEntry{Name: "test"}
			ch <- &mdns.ServiceEntry{Name: name}
			close(ch)
		}()

		info.MDNSParseChan(ch)

		require.False(t, false, info.Complete())
	})

	t.Run("ipv4", func(t *testing.T) {
		ch := make(chan *mdns.ServiceEntry)

		go func() {
			ch <- &mdns.ServiceEntry{
				Name:   name,
				Port:   9000,
				AddrV4: net.IPv4(127, 0, 0, 1),
			}
			close(ch)
		}()

		info.MDNSParseChan(ch)

		require.EqualT(t, "127.0.0.1", info.Host)
		require.EqualT(t, 9000, info.Port)
	})

	t.Run("ipv6", func(t *testing.T) {
		ch := make(chan *mdns.ServiceEntry)

		go func() {
			ch <- &mdns.ServiceEntry{
				Name:         name,
				Port:         9999,
				AddrV6IPAddr: &net.IPAddr{IP: net.IPv6loopback},
			}
			close(ch)
		}()

		info.MDNSParseChan(ch)

		require.EqualT(t, "::1", info.Host)
		require.EqualT(t, 9999, info.Port)
	})
}
