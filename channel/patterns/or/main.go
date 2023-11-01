package main

import (
	"fmt"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	// `or()` enables to combine any number of channels together into a single channel that will close
	// as soon as any of its component channels are closed, or wirtten to.
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		// Fork for the special case where length >=2
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				for {
					select {
					case <-channels[0]:
					case <-channels[1]:
					}
				}

			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:

				// Recurse continously
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	// This function simple creates a channel that will close when the time specified in the `after` elpases
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
