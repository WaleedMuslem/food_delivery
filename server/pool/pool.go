package pool

import (
	"fmt"
	"sync"
)

const (
	defaultBufferSize     = 10000
	defaultListenersCount = 2
)

type WorkerPool struct {
	wg          sync.WaitGroup
	queue       chan func() error
	stop        chan struct{}
	errorCh     chan error
	brokerCount int
}

func NewWorkerPool(errCh chan error) *WorkerPool {
	return &WorkerPool{
		queue:       make(chan func() error, defaultBufferSize),
		stop:        make(chan struct{}),
		brokerCount: defaultListenersCount,
		errorCh:     errCh,
	}
}

func (s *WorkerPool) WithBrokerCount(cnt int) *WorkerPool {
	s.brokerCount = cnt
	return s
}

func (s *WorkerPool) Append(job func() error) {
	s.wg.Add(1)
	s.queue <- job
}

func (s *WorkerPool) Start() {
	for i := 0; i < s.brokerCount; i++ {
		go s.listen()
	}
}

func (s *WorkerPool) Shutdown() {
	s.wg.Wait()
	for i := 0; i < s.brokerCount; i++ {
		s.stop <- struct{}{}
	}
}

func (s *WorkerPool) listen() {
	fmt.Println("run listener")
	for {
		select {
		case job := <-s.queue:
			err := job()
			s.errorCh <- err

			s.wg.Done()
		case <-s.stop:
			fmt.Println("stop listener")
			return
		}
	}
}
