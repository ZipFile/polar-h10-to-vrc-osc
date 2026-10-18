package common

import "errors"

type CombinedSender []Sender

var _ Sender = (CombinedSender)(nil)

func (s CombinedSender) Zero() error {
	var errs []error
	for _, sender := range s {
		errs = append(errs, sender.Zero())
	}
	return errors.Join(errs...)
}

func (s CombinedSender) Value(hr int) error {
	var errs []error
	for _, sender := range s {
		errs = append(errs, sender.Value(hr))
	}
	return errors.Join(errs...)
}

func (s CombinedSender) ActiveStatus(value bool) error {
	var errs []error
	for _, sender := range s {
		errs = append(errs, sender.ActiveStatus(value))
	}
	return errors.Join(errs...)
}

func (s CombinedSender) ConnectedStatus(value bool) error {
	var errs []error
	for _, sender := range s {
		errs = append(errs, sender.ConnectedStatus(value))
	}
	return errors.Join(errs...)
}

func (s CombinedSender) BeatStatus(value bool) error {
	var errs []error
	for _, sender := range s {
		errs = append(errs, sender.BeatStatus(value))
	}
	return errors.Join(errs...)
}
