package actions

type Runnable interface {
	Run(event []byte) error
}
