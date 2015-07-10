package main

type devEventType int

const (
	evRemoved devEventType = iota
	evAdded
	evError
)

type devEvent struct {
	EvType devEventType
	Device *device
	Error  error
}

type device struct {
	Name string
	ID   string
}

func ErrorEvent(err error) *devEvent {
	return &devEvent{
		EvType: evError,
		Device: nil,
		Error:  err,
	}
}
