package vrc_osc

import (
	"time"

	"github.com/ZipFile/polar-h10-to-osc/internal/common"
	"github.com/hypebeast/go-osc/osc"
)

type Sender struct {
	Params Params
	Client *osc.Client
}

var _ common.Sender = (*Sender)(nil)

func (s *Sender) Zero() error {
	return s.send(
		fmt(s.Params.StatusConnected, false),
		fmt(s.Params.StatusActive, false),
		fmt(s.Params.StatusBeat, false),
		fmt(s.Params.ValueAbsolute, float32(0)),
		fmt(s.Params.ValuePercent, float32(0)),
		fmt(s.Params.ValuePercentFull, float32(0)),
	)
}

func (s *Sender) Value(hr int) error {
	return s.send(
		fmt(s.Params.ValueAbsolute, int32(hr)),
		fmt(s.Params.ValuePercent, float32(s.Params.Percent(hr))),
		fmt(s.Params.ValuePercentFull, float32(s.Params.PercentFull(hr))),
	)
}

func (s *Sender) ActiveStatus(value bool) error {
	return s.send(fmt(s.Params.StatusActive, value))
}

func (s *Sender) ConnectedStatus(value bool) error {
	return s.send(fmt(s.Params.StatusConnected, value))
}

func (s *Sender) BeatStatus(value bool) error {
	return s.send(fmt(s.Params.StatusBeat, value))
}

func (s *Sender) send(messages ...*osc.Message) error {
	n := 0

	for _, msg := range messages {
		if msg != nil {
			messages[n] = msg
			n++
		}
	}

	messages = messages[:n]

	switch len(messages) {
	case 0:
		return nil
	case 1:
		return s.Client.Send(messages[0])
	default:
		bundle := osc.Bundle{Messages: messages}

		bundle.Timetag.SetTime(time.Time{})

		return s.Client.Send(&bundle)
	}
}
