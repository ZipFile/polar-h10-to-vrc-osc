package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

const DefaultScanTimeout = 10 * time.Second
const DeviceNamePrefix = "Polar H10 "

var ErrNoAdapterProvided = errors.New("no bluetooth adapter provided")

type Scanner struct {
	Adapter     *bluetooth.Adapter
	MACAddress  bluetooth.MAC
	DeviceName  string
	ScanTimeout time.Duration
}

func (s *Scanner) Match(device bluetooth.ScanResult) bool {
	if s == nil {
		goto end
	}

	if device.Address.MAC == s.MACAddress {
		return true
	}

	if s.DeviceName != "" {
		return device.LocalName() == s.DeviceName
	}

end:
	return strings.HasPrefix(device.LocalName(), DeviceNamePrefix)
}

func (s *Scanner) GetScanTimeout() time.Duration {
	if s == nil || s.ScanTimeout == 0 {
		return DefaultScanTimeout
	}
	return s.ScanTimeout
}

func (s *Scanner) Scan(ctx context.Context) (device *bluetooth.ScanResult, err error) {
	if s == nil || s.Adapter == nil {
		return nil, ErrNoAdapterProvided
	}

	ctx, cancel := context.WithTimeout(ctx, s.GetScanTimeout())
	defer cancel()

	go func() {
		<-ctx.Done()
		s.Adapter.StopScan()
	}()

	err = s.Adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if s.Match(result) {
			device = &result
			adapter.StopScan()
		}
	})

	return
}
