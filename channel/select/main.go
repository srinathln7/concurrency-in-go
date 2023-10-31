package main

import (
	"fmt"
	"time"
)

func main() {

	// What happens when mutliple channels have something to read?
	c1 := make(chan interface{})
	c2 := make(chan interface{})

	close(c1)
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d \nc2Count: %d\n", c1Count, c2Count)

	// What happens if they are never any channels that become ready?
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

	// What if we want to do something in the meanwhile no channels are ready?

	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("achieved %v cycles of work before signaled to stop. \n", workCounter)
}
