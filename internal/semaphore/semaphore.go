package semaphore

type Weighted struct {
	c chan struct{}
}

func New(workers int) *Weighted {
	return &Weighted{
		c: make(chan struct{}, workers),
	}
}

func (w *Weighted) Acquire() {
	w.c <- struct{}{}
}

func (w *Weighted) Free() {
	<-w.c
}

func (w *Weighted) Close() {
	close(w.c)
}
