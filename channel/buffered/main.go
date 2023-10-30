package main

import "fmt"

func main() {

	// Unbuffered channel - range over a channel
	intUnbufStream := make(chan int)
	go func() {
		fmt.Println("SEND TO UNBUFFERED CHANNEL")
		defer close(intUnbufStream)
		for i := 1; i <= 5; i++ {
			intUnbufStream <- i
			fmt.Println("Sending i=", i)
		}
	}()

	for integer := range intUnbufStream {
		fmt.Println("RECEIVE FROM UNBUFFERED CHANNEL")
		fmt.Printf("Received %v\n", integer)
	}

	// Buffered channel example

	intBufStream := make(chan int, 4)
	go func() {
		fmt.Println("SEND TO BUFFERED CHANNEL")
		defer close(intBufStream)
		for i := 1; i <= 5; i++ {
			intBufStream <- i
			fmt.Println("Sending i=", i)
		}
	}()

	for integer := range intBufStream {
		fmt.Println("RECEIVE FROM BUFFERED CHANNEL")
		fmt.Printf("Received %v\n", integer)
	}
}
