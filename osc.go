package main

import (
	"context"
	"errors"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/siiimooon/go-polar/pkg/h10"
)

func FormatAvatarParam[V any](name string, value V) (msg *osc.Message) {
	msg = osc.NewMessage("/avatar/parameters/" + name)
	msg.Append(value)
	return msg
}

type OSCRelay struct {
	Client      *osc.Client
	MinHR       int
	MaxHR       int
	IsConnected bool
}

const DefaultMinHR = 32
const DefaultMaxHR = 192

func (r *OSCRelay) GetHRPercent(hr int) float64 {
	if r.MaxHR == r.MinHR || hr <= r.MinHR {
		return 0
	} else if hr >= r.MaxHR {
		return 1
	}

	return float64(hr-r.MinHR) / float64(r.MaxHR-r.MinHR)
}

func (r *OSCRelay) SendHR(hr int) error {
	percent := r.GetHRPercent(hr)

	return errors.Join(
		r.Client.Send(FormatAvatarParam("HR", int32(hr))),
		r.Client.Send(FormatAvatarParam("HRPercent", float32(percent))),
		r.Client.Send(FormatAvatarParam("FullHRPercent", float32(2*percent-1))),
	)
}

func (r *OSCRelay) SendActiveStatus(value bool) error {
	return r.Client.Send(FormatAvatarParam("isHRActive", value))
}

func (r *OSCRelay) SendConnectedStatus(value bool) error {
	return r.Client.Send(FormatAvatarParam("isHRConnected", value))
}

func (r *OSCRelay) SendIsBeating(value bool) error {
	return r.Client.Send(FormatAvatarParam("isHRBeat", value))
}

func (r *OSCRelay) SendZero() error {
	return errors.Join(
		r.Client.Send(FormatAvatarParam("isHRConnected", false)),
		r.Client.Send(FormatAvatarParam("isHRActive", false)),
		r.Client.Send(FormatAvatarParam("isHRBeat", false)),
		r.Client.Send(FormatAvatarParam[int32]("HR", 0)),
		r.Client.Send(FormatAvatarParam[float32]("HRPercent", 0)),
		r.Client.Send(FormatAvatarParam[float32]("FullHRPercent", 0)),
	)
}

func (r *OSCRelay) Do(ctx context.Context, data <-chan h10.HeartRateMeasurement) {
	isBeating := true
	isActiveTicker := time.NewTicker(29 * time.Second)
	isBeatingTicker := time.NewTicker(5 * time.Second)
	defer isActiveTicker.Stop()
	defer isBeatingTicker.Stop()
	defer r.SendZero()

	r.SendConnectedStatus(r.IsConnected)
	r.SendActiveStatus(true)

	for {
		select {
		case <-isActiveTicker.C:
			r.SendConnectedStatus(r.IsConnected)
			r.SendActiveStatus(true)
		case <-isBeatingTicker.C:
			r.SendIsBeating(isBeating)
			isBeating = false
		case msg, ok := <-data:
			if !ok {
				return
			}

			print("\r\033[KHR: ", msg.GetHeartRate())

			isBeating = true
			err := r.SendHR(msg.GetHeartRate())

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
