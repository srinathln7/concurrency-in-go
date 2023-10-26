package main

import "fmt"

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
}
