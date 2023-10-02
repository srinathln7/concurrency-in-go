package main

import (
	"fmt"
	"sync"
	"time"
)

// Demonstrating the concept of starvation: Starvation refers to any situation where a concurrent process cannot get all the resources
// it requires to perform work. Refer to live-lock example where each goroutine was starved of a "shared lock" resource.

// Key take away: Opt for fairness (fine-grained sync) over course-grained sync (performance) to start with
func main() {
	var wg sync.WaitGroup
	var count int
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	// course-grained sync - performance
	// broaden critical section
	greedyWorker := func() {
		defer wg.Done()
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("greedy worker executed %v work loops \n", count)
	}

	// fine-grained sync - fairness
	// narrow critical section and mem sync as much as possible
	politeWorker := func() {
		defer wg.Done()
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("polite worker executed %v work loops \n", count)
	}

	wg.Add(2)

	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
