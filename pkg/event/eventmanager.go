package event

type EventManager interface {
	Init() error
	Process() error
	Close() error
}
