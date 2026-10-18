package common

type EventKind uint8
type Status uint8

type Sender interface {
	Zero() error
	Value(int) error
	ConnectedStatus(bool) error
	ActiveStatus(bool) error
	BeatStatus(bool) error
}
