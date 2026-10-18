package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/ZipFile/polar-h10-to-osc/internal/common"
	"github.com/ZipFile/polar-h10-to-osc/internal/osc/vnyan"
	"github.com/ZipFile/polar-h10-to-osc/internal/osc/vrc"
	"github.com/hypebeast/go-osc/osc"
	"github.com/siiimooon/go-polar/pkg/h10"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		if !connected {
			cancel()
		}
	})

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
	params := *vrc_osc.PresetDefault
	params.Min = 32
	params.Max = 192
	vrcSender := &vrc_osc.Sender{
		Client: osc.NewClient(vrc_osc.DefaultHost, vrc_osc.DefaultPort),
		Params: params,
	}
	vnyanSender := &vnyan_osc.Sender{
		Client: osc.NewClient(vnyan_osc.DefaultHost, vnyan_osc.DefaultPort),
	}
	relay := NewOSCRelay(
		common.CombinedSender{vrcSender, vnyanSender},
		true,
	)
	httpAddr := defaultAddr

	println("Sending VRChat OSC data to", fmt.Sprintf("%s:%d", vrc_osc.DefaultHost, vrc_osc.DefaultPort))
	println("Sending VNyan OSC data to", fmt.Sprintf("%s:%d", vnyan_osc.DefaultHost, vnyan_osc.DefaultPort))
	println("URL: http://" + httpAddr)
	println("Press Ctrl+C to exit")

	go reader.StreamHeartRate(ctx, hr)
	go relay.Do(ctx, hr)
	go serveHttp(hr, httpAddr, ctx)

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
