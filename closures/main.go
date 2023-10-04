package main

import (
	"fmt"
	"sync"
)

// Go routines operate within the same address space as each other.
func main() {

	var wg sync.WaitGroup
	var name string = "Srinath"

	wg.Add(1)

	go func() {
		defer wg.Done()

		// go routines execute within the same address space they were created in
		name = name + "-Srinidhi"
	}()

	for _, salutation := range []string{"Hello", "Greetings", "Good day!"} {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// go runtime is observant enough to notice that the reference to `salutation` is maintained
			// and hence transfers this variable to the heap. Notice only the last reference of `salutation`
			// -> `goodbye` is maintained
			fmt.Println(" ", salutation)
		}()
	}

	for _, day := range []string{"Sunday", "Monday", "Tuesday"} {
		wg.Add(1)

		go func(day string) {
			defer wg.Done()

			fmt.Println(" ", day)
		}(day)
	}

	// Join point
	wg.Wait()

	fmt.Println(name)
}
