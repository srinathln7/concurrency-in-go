package main

import (
	"fmt"
	"sync"
)

func main() {

	var count int

	increment := func() {
		count++
	}

	var wg sync.WaitGroup
	var once sync.Once

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			once.Do(increment)
		}()
	}

	wg.Wait()

	// count: 1
	fmt.Println("count:", count)

	decrement := func() {
		count--
	}

	// sync.Once only counts the number of times Do is called, and not how many times unique functions passed into Do are called.
	// In this way, copies of sync.Once are tightly coupled to the functions they are intended to be called with.
	once.Do(decrement)

	// count: 1
	fmt.Println("count:", count)

	var onceA sync.Once
	onceA.Do(decrement)

	// count: 0
	fmt.Println("count:", count)
}
