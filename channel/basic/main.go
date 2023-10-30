package main

import (
	"fmt"
	"sync"
)

func main() {

	var strStream chan interface{}
	strStream = make(chan interface{})

	go func() {
		strStream <- "srinath-srinidhi"
	}()

	fmt.Printf("value read from channel %v \n", <-strStream)

	// send only channel
	// writeStream := make(chan<- interface{})

	// read only channel
	// readStream := make(<-chan interface{})

	// writeStream <- "srinath"
	// <-readStream

	stringStream := make(chan string)

	go func() {
		stringStream <- "hello channels!"
	}()

	value, ok := <-stringStream
	fmt.Printf("(%v): %v \n", value, ok)

	intStream := make(chan int)

	// The ability to close the channel is a very desirable quality for any program
	close(intStream)

	value1, ok1 := <-intStream
	fmt.Printf("(%v): %v \n", value1, ok1)

	// range over a channel
	integerStream := make(chan int)
	go func() {
		defer close(integerStream)
		for i := 1; i <= 5; i++ {

			fmt.Println("Sending i=", i)
			integerStream <- i
		}
	}()

	for integer := range integerStream {
		fmt.Printf("Received %v\n", integer)
	}

	//simple way to unblock multiple go routines at once

	begin := make(chan interface{})
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Block all go routines here until a value has been written into the channel
			<-begin
			fmt.Printf("%v has begun \n", i)
		}(i)
	}

	// Closing a channel is one of the ways you can signal multiple goroutines simultaneously
	// If you have n goroutines waiting on a single channel, instead of writing n times to the channel to unblock each goroutines
	// you can simply close the channel. Closing is both cheaper and faster than performing `n` writes.

	close(begin) // Unblock all go routines

	wg.Wait()
}
