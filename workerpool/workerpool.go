package workerpool

import (
	"fmt"
	"sync"
)

type (
	Pool struct {
		jobs  <-chan string
		stop  chan struct{}
		wg    *sync.WaitGroup
		count int
		id    int
	}
)

func (p *Pool) Close() {
	close(p.stop)
	p.count = 0
	p.wg.Wait()
}

func (p *Pool) StopOne() {
	if p.count > 0 {
		p.stop <- struct{}{}
		p.count--
	}
}

func (p *Pool) Add() {
	p.wg.Add(1)
	go p.worker(p.id)
	fmt.Println("Worker ", p.id, " start")
	p.count++
	p.id++
}

func NewPool(jobs <-chan string) *Pool {
	return &Pool{jobs: jobs, stop: make(chan struct{}), wg: &sync.WaitGroup{}}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	defer fmt.Println("Worker ", id, "delete")

	for {
		select {
		case j, ok := <-p.jobs:
			if !ok {
				return
			}
			fmt.Println("Worker", id, " ", j)
		case <-p.stop:
			return
		}
	}
}
