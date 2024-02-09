package main

import (
	"log"
	"sync"
	"time"
)

type job func()

type pool struct {
	wg        sync.WaitGroup
	workqueue chan job
}

func newPool(count int) *pool {
	pool := pool{
		workqueue: make(chan job),
	}

	pool.wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			log.Printf("worker %d reading from the work queue\n", i)
			defer pool.wg.Done()
			for job := range pool.workqueue {
				job()
			}
		}(i)
	}
	return &pool
}

func (p *pool) AddJob(job func(), id int) {
	log.Printf("pushing job:%d to the work queue\n", id)
	p.workqueue <- job
}

func (p *pool) Wait() {
	close(p.workqueue)
	p.wg.Wait()
}

func main() {
	jobPool := newPool(5)
	for i := 0; i < 10; {
		job := func() {
			log.Printf("executing job:%d\n", i)
			time.Sleep(1 * time.Second)
			log.Printf("job %d completed\n", i)
		}

		jobPool.AddJob(job, i)

	}

	jobPool.Wait()
}
