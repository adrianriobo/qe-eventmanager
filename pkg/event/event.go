package event

type Event interface {
	GetDestination() error
	Handler(event interface{}) error
}
