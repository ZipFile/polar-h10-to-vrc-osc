package vrc_osc

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

func (i *Info) MDNSParse(entry *mdns.ServiceEntry) {
	if entry == nil || !strings.HasPrefix(entry.Name, "VRChat-Client-") {
		return
	}

	if entry.Port != 0 {
		i.Port = entry.Port
	}

	if entry.AddrV4 != nil {
		i.Host = entry.AddrV4.String()
	} else if entry.AddrV6IPAddr != nil {
		i.Host = entry.AddrV6IPAddr.String()
	}
}

func (i *Info) MDNSParseChan(entries <-chan *mdns.ServiceEntry) {
	for entry := range entries {
		i.MDNSParse(entry)
	}
}

func (i *Info) MDNSQuery(ctx context.Context) error {
	entries := make(chan *mdns.ServiceEntry, 4)

	defer close(entries)

	params := mdns.QueryParam{
		Service:             "_osc._udp",
		Domain:              "local",
		Timeout:             500 * time.Millisecond,
		Entries:             entries,
		WantUnicastResponse: true,
		DisableIPv4:         false,
		DisableIPv6:         false,
	}

	go i.MDNSParseChan(entries)

	return mdns.QueryContext(ctx, &params)
}
