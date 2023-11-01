package main

import (
	"fmt"
	"time"
)

func main() {

	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()

		return completed
	}

	// Passing a nil channel into doWork go routine, the strings channel will never
	// actually get any strings written onto it, and the goroutine containins doWork
	// will remain in memory for the lifetime of this process
	doWork(nil)

	fmt.Println("Done.")

	// To signal completion, we establish a signal between the parent goroutine and the children
	// that allows the parent to signal cancellation to its children. We achieve this through convention
	// of using the `done` channel which is usually a read-only channel named done.

	doWorkCorrect := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {

		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")

			// Without closing this channel, the program would result in deadlock
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWorkCorrect(done, nil)

	// Before we create a join point, we create a third goroutine that cancel the goroutine within
	// doWork after a second to successfully eliminate our goroutine leak
	go func() {
		// Cancel the operation after one second
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling doWork goroutine")
		close(done)
	}()

	// Join point
	<-terminated
	fmt.Println("Correctly done!")
}
