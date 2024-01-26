package main

import (
	"log"
	"sync"
)

const NUM_OF_WORKERS = 10000

type Queue struct {
	mu    sync.Mutex
	queue []int
}

func newQueue() *Queue {
	return &Queue{}
}

func (q *Queue) enqueue(wg *sync.WaitGroup, val int) {
	// Remember defer oeprations are like stack - LIFO.
	// So second defer line will return first (Line 25)

	defer wg.Done()
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, val)
	log.Printf("enqueued item %d to the queue \n", val)
}

func (q *Queue) dequeue(wg *sync.WaitGroup) {

	defer wg.Done()
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return
	} else {
		item := q.queue[0]
		if len(q.queue) == 1 {
			q.queue = []int{}
		} else {
			q.queue = q.queue[1:]
		}
		log.Printf("dequeued item %d from the queue \n", item)
	}
}

func main() {

	//var wg sync.WaitGroup
	var wgE, wgD sync.WaitGroup

	qu := newQueue()

	wgE.Add(NUM_OF_WORKERS)
	for i := 0; i < NUM_OF_WORKERS; i++ {
		go qu.enqueue(&wgE, i)

	}

	wgE.Wait()

	log.Printf("concurrent enqueue operations complete and current queue size is %d", len(qu.queue))

	wgD.Add(NUM_OF_WORKERS)
	for i := 0; i < NUM_OF_WORKERS; i++ {
		go qu.dequeue(&wgD)
	}

	wgD.Wait()

	log.Printf("concurrent dequeue operations complete and current queue size is %d", len(qu.queue))
}
