package main

import (
	"fmt"
	"sync"
)

func sayHello(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Hello from go routine %v \n", id)
}

// Waitgroups: Use waitgroups only when you don't care about the intermediate results of the concurrent operations
// or have other means to get the results of those concurrent operations
func main() {
	var wg sync.WaitGroup

	const num_routines = 5
	wg.Add(num_routines)

	for i := 0; i < num_routines; i++ {
		go sayHello(&wg, i)
	}

	wg.Wait()
}
