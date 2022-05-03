package actions

type Runnable interface {
	Run(actionInformation interface{}) error
}

type Action struct {
}
