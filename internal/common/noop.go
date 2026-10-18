package common

type NoopSender struct {
	Err error
}

var NOOP Sender = (*NoopSender)(nil)

func (s *NoopSender) e() error {
	if s == nil || s.Err == nil {
		return nil
	}
	return s.Err
}

func (s *NoopSender) Zero() error {
	return s.e()
}

func (s *NoopSender) Value(int) error {
	return s.e()
}

func (s *NoopSender) ActiveStatus(bool) error {
	return s.e()
}

func (s *NoopSender) ConnectedStatus(bool) error {
	return s.e()
}

func (s *NoopSender) BeatStatus(bool) error {
	return s.e()
}
