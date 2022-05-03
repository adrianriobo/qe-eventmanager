package manager

type Event interface {
	GetDestination() error
	Handler(event interface{}) error
}
