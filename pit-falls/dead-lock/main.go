package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

// Deadlocks refer to a situation in which two or more concurrent processes are waiting on one anther and stall the program.
// Simulating a dead-lock situation. Deadlocks can be detected with the help of Coffman conditions
func main() {

	wg.Add(2)

	var a, b value
	go printSum(&a, &b)
	go printSum(&b, &a)

	wg.Wait()
}

type value struct {
	mu  sync.Mutex
	val int
}

var printSum = func(v1, v2 *value) {
	defer wg.Done()

	v1.mu.Lock()
	defer v1.mu.Lock()

	// Simulate work
	//time.Sleep(2 * time.Second)

	v2.mu.Lock()
	defer v2.mu.Unlock()

	fmt.Println("Sum: ", v1.val+v2.val)
}
