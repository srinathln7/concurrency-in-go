package main

import (
	"fmt"
	"sync"
	"time"
)

const QUEUE_FIX_LENGTH = 7
const TOT_CAP_SIZE = 50

func main() {

	c := sync.NewCond(&sync.Mutex{})

	queue := make([]interface{}, 0, TOT_CAP_SIZE)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Printf("Removing element %v from queue \n", queue[0])
		queue = queue[1:]
		c.L.Unlock()
		c.Signal()
	}

	startTime := time.Now()

	for i := 0; i < TOT_CAP_SIZE; i++ {

		c.L.Lock()

		for len(queue) == QUEUE_FIX_LENGTH {
			c.Wait()
		}

		fmt.Printf("Adding element %v to the queue\n", i)
		queue = append(queue, i)
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()

	}

	// Wait for all dequeue operations to finish
	for len(queue) > 0 {
		c.L.Lock()
		for len(queue) > 0 {
			c.Wait()
		}
		c.L.Unlock()
	}

	elapsedTime := time.Since(startTime)
	fmt.Println("Total time taken: ", elapsedTime)
	fmt.Println("Final queue: ", queue)
}
