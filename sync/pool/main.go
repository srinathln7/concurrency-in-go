package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create a sync.Pool with a New function to initialize new []byte slices
	bytePool := sync.Pool{
		New: func() interface{} {
			fmt.Println("Allocating new byte slice")
			return make([]byte, 1024) // Allocate a 1KB buffer
		},
	}

	// Get a []byte slice from the pool
	buffer := bytePool.Get().([]byte)

	// Use the buffer for some work
	for i := 0; i < len(buffer); i++ {
		buffer[i] = byte(i % 256)
	}
	fmt.Println("Using buffer for some work")

	// Once done, put the buffer back into the pool for reuse
	bytePool.Put(buffer)

	// Get the buffer again from the pool (this time it should be reused)
	buffer2 := bytePool.Get().([]byte)

	// Verify that the buffer is the same
	if &buffer[0] == &buffer2[0] {
		fmt.Println("Reused the same buffer")
	} else {
		fmt.Println("Allocated a new buffer")
	}
}
