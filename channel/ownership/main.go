package main

import "fmt"

func main() {

	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)

		// Start a go routine that performs write operations on the resultStream channel
		// Notice, creation of goroutine is encapsulated within the surrounding function
		go func() {
			// close the channel inside the go routine -> CHANNEL OWNER responsbility
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()

		// implicit conversion to read-only channel
		return resultStream
	}

	// CHANNEL-utilizer
	chanUtilizer := chanOwner()
	for result := range chanUtilizer {
		fmt.Println("Received ", result)
	}
}
