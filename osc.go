package main

import (
	"context"
	"time"

	"github.com/ZipFile/polar-h10-to-osc/internal/common"
	"github.com/siiimooon/go-polar/pkg/h10"
)

type OSCRelay struct {
	Sender      common.Sender
	IsConnected bool
}

func NewOSCRelay(sender common.Sender, isConnected bool) *OSCRelay {
	return &OSCRelay{
		Sender:      sender,
		IsConnected: isConnected,
	}
}

func (r *OSCRelay) Do(ctx context.Context, data <-chan h10.HeartRateMeasurement) {
	isBeating := true
	isActiveTicker := time.NewTicker(29 * time.Second)
	isBeatingTicker := time.NewTicker(5 * time.Second)
	defer isActiveTicker.Stop()
	defer isBeatingTicker.Stop()
	defer func() {
		_ = r.Sender.Zero()
	}()

	_ = r.Sender.ConnectedStatus(r.IsConnected)
	_ = r.Sender.ActiveStatus(true)

	for {
		select {
		case <-isActiveTicker.C:
			_ = r.Sender.ConnectedStatus(r.IsConnected)
			_ = r.Sender.ActiveStatus(true)
		case <-isBeatingTicker.C:
			_ = r.Sender.BeatStatus(isBeating)
			isBeating = false
		case msg, ok := <-data:
			if !ok {
				return
			}

			print("\r\033[KHR: ", msg.GetHeartRate())

			isBeating = true
			err := r.Sender.Value(msg.GetHeartRate())
			if err != nil {
				print("\r\033[KError: ", err.Error())
			}
		case <-ctx.Done():
			println()
			println("Done")
			return
		}
	}
}
