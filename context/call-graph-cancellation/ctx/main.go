package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Program demonstrates the use of `context` pkg to provide an API for cancelling brances of your call-graph.
// We can view the `context` pkg as a wrapper around the `done` channel pattern used in the `simple/main.go` file.
func main() {

	var wg sync.WaitGroup

	//`WithCancel` fns return a new `context` that closes its `done` channel
	// when the returned `cancel` function is called
	ctx, cancel := context.WithCancel(context.Background())

	// Call the `cancel` function before the `main` exits
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
			cancel()
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(ctx context.Context) (string, error) {

	// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)) and closes the `done` channel after the given `timeOut` duration.
	// NOTICE the use of the new context in this function call unlike others to ensure that `genGreeting` timesout after 1s.
	// This is the beauty of the context pkg. Children can use its own context pkg.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// notice the use of `switch` block and not `select`
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	// notice the use of `switch` block and not `select`
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

// `locale` is the lowest call stack in the program
func locale(ctx context.Context) (string, error) {

	// cannot print greeting: context deadline exceeded
	// cannot print farewell: context canceled
	timeout := 10 * time.Second

	// hello world!
	// goodbye world!
	//timeout := 1 * time.Second

	select {
	// `Done()` returns a closed channel when work done on behalf of this context should be cancelled.
	// Done may return `nil` if this context can never be cancelled. Successive calls to done return the same value.
	// Rem: `read` from a closed channel will return `default` value while `read` from a `nil` channel will BLOCK.
	// This cond'n gets UNBLOCKED because the parent call `genGreeting` cancels the context after a deadline of 1s.
	case <-ctx.Done():
		return "", ctx.Err()

	case <-time.After(timeout):
	}

	return "EN/US", nil
}
