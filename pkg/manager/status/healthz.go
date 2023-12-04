package status

import (
	"fmt"
	"net/http"
	"sync"
)

type Status struct {
	stateHandler *sync.Mutex
	stateChannel chan struct{}
	state        string
}

const (
	defaultPath string = "/healthz"
	defaultPort string = "8080"

	stateHealthy   string = "healthy"
	stateUnhealthy string = "unhealthy"
)

var status *Status

func Init() error {

	status = &Status{
		stateHandler: &sync.Mutex{},
		stateChannel: make(chan struct{}, 1),
		state:        stateHealthy}

	http.HandleFunc(defaultPath, func(w http.ResponseWriter, _ *http.Request) {
		var httpStatusCode = http.StatusOK
		if getState() == stateUnhealthy {
			httpStatusCode = http.StatusInternalServerError
		}
		w.WriteHeader(httpStatusCode)
	},
	)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", defaultPort), nil); err != nil {
		return err
	}

	go checkState(status)

	return nil
}

func getState() (current string) {
	status.stateHandler.Lock()
	current = status.state
	status.stateHandler.Unlock()
	return current
}

func checkState(s *Status) {
	<-s.stateChannel
	s.stateHandler.Lock()
	s.state = stateUnhealthy
	s.stateHandler.Unlock()
}

func SendSignal() {
	status.stateChannel <- struct{}{}
}
