package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	writerStream := func() <-chan int {
		randStream := make(chan int)
		go func() {

			// The defer statmene will never get run in this case. After third iteration of the loop
			// in the main goroutine, our goroutine blocks trying to send the next random integer to a channel
			// that is no longer being read from. We have now no way of telling the producer to STOP.
			defer fmt.Println("Closure of writerStream")

			// not necessary in this example
			defer close(randStream)

			for {
				randStream <- rand.Int()
			}
		}()

		return randStream
	}

	randReadStream := writerStream()
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randReadStream)
	}

	writerStreamCorrect := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("Closure of writerStreamCorrect")

			// not necessary in this example
			defer close(randStream)

			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done := make(chan interface{})
	randReadStreamCorrect := writerStreamCorrect(done)
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randReadStreamCorrect)
	}
	close(done)

	// Simulate some work
	time.Sleep(1 * time.Second)
}
