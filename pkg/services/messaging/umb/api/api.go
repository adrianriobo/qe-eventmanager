package api

type ClientInterface interface {
	Disconnect()
	Subscribe(destination string, handlers []MessageHandler) (SubscriptionInterface, error)
	Send(destination string, message []byte) error
}

type SubscriptionInterface interface {
	Read() ([]byte, error)
	Unsubscribe() error
}

type MessageHandler interface {
	// Match(event interface{}, filters []string) error
	Match(event []byte) error
	// Handle(event interface{}) error
	Handle(event []byte) error
}
