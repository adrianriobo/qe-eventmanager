package api

type ClientInterface interface {
	Disconnect()
	Subscribe(destination string, handlers []func(event interface{}) error) (SubscriptionInterface, error)
	Send(destination string, message interface{}) error
}

type SubscriptionInterface interface {
	Read() ([]byte, error)
	Unsubscribe() error
}
