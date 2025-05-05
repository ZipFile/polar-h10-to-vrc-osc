package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/siiimooon/go-polar/pkg/h10"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var oscHost = "127.0.0.1"
var oscPort = 9000

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	scanner := Scanner{Adapter: adapter}
	scanResult, err := scanner.Scan(ctx)

	must("scan for device", err)

	if scanResult == nil {
		println("no devices found")
		os.Exit(1)
	}

	println("Device:", scanResult.LocalName())
	println("Address:", scanResult.Address.String())

	device, err := adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{
		MinInterval: bluetooth.NewDuration(500 * time.Millisecond),
		MaxInterval: bluetooth.NewDuration(4 * time.Second),
		Timeout:     bluetooth.NewDuration(30 * time.Second),
	})

	must("connect to device", err)
	println("Connected!")

	defer device.Disconnect()

	reader := h10.New(&device)
	hr := make(chan h10.HeartRateMeasurement, 10)
	relay := OSCRelay{
		Client:      osc.NewClient(oscHost, oscPort),
		MinHR:       DefaultMinHR,
		MaxHR:       DefaultMaxHR,
		IsConnected: true,
	}

	println("Sending OSC data to", fmt.Sprintf("%s:%d", oscHost, oscPort))
	println("Min HR:", relay.MinHR)
	println("Max HR:", relay.MaxHR)
	println("Press Ctrl+C to exit")

	go reader.StreamHeartRate(ctx, hr)
	go relay.Do(ctx, hr)

	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
}

func must(action string, err error) {
	if err != nil {
		panic("Failed to " + action + ": " + err.Error())
	}
}
