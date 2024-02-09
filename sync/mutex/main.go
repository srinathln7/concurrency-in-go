package main

import (
	"fmt"
	"sync"
)

// `mutex`: refers to mutual exclusion that are used to gaurd the critical sections of our programs.
//
//	This is essential to synchronize access to the memory
func main() {

	var count int
	var wg sync.WaitGroup
	var mu sync.Mutex

	inc := func() {
		mu.Lock()
		defer mu.Unlock()

		// It is important to specify wg.Done() here to decrement the wait group counter
		defer wg.Done()
		count++
		fmt.Println("incrementing count to", count)
	}

	dec := func() {
		mu.Lock()
		defer mu.Unlock()
		defer wg.Done()
		count--
		fmt.Println("decrementing count to", count)
	}

	const num_of_iterations = 5

	wg.Add(2 * num_of_iterations)

	// increment go routines
	for i := 0; i < num_of_iterations; i++ {

		// defer wg.Done() --> No use specifying the wait group Done() here. This will only result in DEADLOCK
		go inc()
	}

	// decrement go routines
	for i := 0; i < num_of_iterations; i++ {

		// defer wg.Done() --> No use specifying the wait group Done() here. This will only result in DEADLOCK
		go dec()
	}

	wg.Wait()

	fmt.Println("Final value of count", count)
}
