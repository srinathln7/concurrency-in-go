package main

import (
	"fmt"
	"sync"
	"time"
)

const QUEUE_FIX_LENGTH = 2
const TOT_CAP_SIZE = 10

// This program demonstrates how to efficiently enqueue and dequeue elements from a queue
// with a fixed length
func main() {

	c := sync.NewCond(&sync.Mutex{})

	// init a slice of length 0 and total capacity of 10
	queue := make([]interface{}, 0, TOT_CAP_SIZE)

	// Dequeue operations
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Printf("Removing element %v from queue \n", queue[0])
		queue = queue[1:]
		c.L.Unlock()
		c.Signal() // Signals to the other go routine that something has occured
	}

	// Enqueue the operations
	for i := 0; i < TOT_CAP_SIZE; i++ {

		c.L.Lock()

		// This condition is the most important one
		for len(queue) == QUEUE_FIX_LENGTH {
			// suspend the current go routine until a signal on the condition has been sent
			c.Wait()
		}

		fmt.Printf("Adding element %v to the queue\n", i)
		queue = append(queue, i)
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()

	}

	// Dequeue the remaining elements
	for j := 0; j < QUEUE_FIX_LENGTH; j++ {
		fmt.Printf("Removing element %v from the queue\n", TOT_CAP_SIZE-QUEUE_FIX_LENGTH+j)
		queue = queue[1:]
	}

	fmt.Println("Final queue: ", queue)
}
