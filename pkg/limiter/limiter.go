package limiter

import (
	"github.com/mrkish/crawly/pkg/cache"
	"github.com/mrkish/crawly/pkg/semaphore"
)

type Request struct {
	connections int
	cache       *cache.Hashed[string, *semaphore.Weighted]
}

func New(connections int) *Request {
	return &Request{
		connections: connections,
		cache:       cache.NewHashed[string, *semaphore.Weighted](nil, nil),
	}
}

func (l *Request) Request(url string, fn func()) {
	if sem := l.cache.Get(url); sem != nil {
		defer sem.Free()
		sem.Acquire()
		fn()
	} else {
		sem := semaphore.New(l.connections)
		l.cache.Add(url, sem)

		defer sem.Free()
		sem.Acquire()
		fn()
	}
}
