package main

import (
	"fmt"
	"sync"
	"time"
)

// Program that concurrently prints a greeting and a farewell msg while ignoring the race condition.
// Closing the done channel at any point in `main`, both branches will be cancelled.
// At each stack frame, a function can affect the entirey of the call stack below it
func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	// notice the use of `switch` block and not `select`
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	// notice the use of `switch` block and not `select`
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("cancelled")
	case <-time.After(10 * time.Second):
	}

	return "EN/US", nil
}
