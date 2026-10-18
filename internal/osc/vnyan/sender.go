package vnyan_osc

import (
	"github.com/ZipFile/polar-h10-to-osc/internal/common"
	"github.com/hypebeast/go-osc/osc"
)

type Sender struct {
	Client *osc.Client
}

var _ common.Sender = (*Sender)(nil)

func (s *Sender) Zero() error {
	return s.Value(0)
}

func (s *Sender) Value(hr int) error {
	return s.Client.Send(
		osc.NewMessage(
			"/VNyan/Param/Float",
			"heartrate",
			float32(hr),
		),
	)
}

func (*Sender) ActiveStatus(bool) error {
	return nil
}

func (*Sender) ConnectedStatus(bool) error {
	return nil
}

func (*Sender) BeatStatus(bool) error {
	return nil
}
