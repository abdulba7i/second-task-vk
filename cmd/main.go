package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id          int
	stopChannel chan struct{}
	stopped     chan struct{}
}

type Pool struct {
	mu         sync.Mutex
	workers    []*Worker
	jobChannel chan string
	nextID     int
}

func NewPool() *Pool {
	return &Pool{
		jobChannel: make(chan string),
	}
}

func (p *Pool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	id := p.nextID
	p.nextID++

	worker := &Worker{
		id:          id,
		stopChannel: make(chan struct{}),
		stopped:     make(chan struct{}),
	}

	p.workers = append(p.workers, worker)

	go func() {
		defer close(worker.stopped)
		for {
			select {
			case job := <-p.jobChannel:
				time.Sleep(800 * time.Millisecond)
				fmt.Printf("Worker %d: обработал -> %s\n", worker.id, job)
			case <-worker.stopChannel:
				fmt.Printf("Worker %d: остановлен\n", worker.id)
				return
			}
		}
	}()
}

func (p *Pool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) == 0 {
		fmt.Println("Все воркеры завершены")
		return
	}

	worker := p.workers[len(p.workers)-1]
	p.workers = p.workers[:len(p.workers)-1]

	close(worker.stopChannel)
	<-worker.stopped
}

func (p *Pool) Submit(job string) {
	p.jobChannel <- job
}

func main() {
	pool := NewPool()

	pool.AddWorker()
	pool.AddWorker()

	for i := 0; i <= 10; i++ {
		msg := fmt.Sprintf("сообщение: %d", i)
		pool.Submit(msg)
	}

	time.Sleep(1 * time.Second)

	pool.RemoveWorker()

	time.Sleep(1 * time.Second)

	pool.RemoveWorker()
	pool.RemoveWorker()
}
