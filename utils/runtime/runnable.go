package runtime

type Runnable interface {
	Run(done <-chan struct{})
}
