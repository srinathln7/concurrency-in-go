package main

import "fmt"

func main() {

	// orDone channel pattern:  Useful pattern for managing the lifecycle of goroutines that read from multiple channels.
	// It's especially handy when you want to coordinate multiple goroutines and stop them when any of them is done.
	// The pattern is particularly useful for scenarios like aggregating results from multiple sources or waiting for the first goroutine to complete.
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})

		go func() {
			defer close(valStream)

			for {
				select {
				case <-done:
					return
				// Read from the given channel
				case v, ok := <-c:
					// if it is a nil channel i.e. reading default value
					if !ok {
						return
					}
					select {
					case <-done:
					case valStream <- v:
					}
				}
			}
		}()

		return valStream
	}

	// tee (duplicate) channel pattern : Involves duplicating or "teeing" the data from one channel into two or more separate channels.
	// 	This can be useful in various scenarios, such as parallel processing or creating fan-out and fan-in patterns.
	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {

		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)

			for val := range orDone(done, in) {
				// We use local versions of `out1` and `out2` to shadow the outer scope variables
				var out1, out2 = out1, out2

				// Use one select statement so that writes to `out1` and `out2` dont block each other
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						// once a val is written to channel, we set the copy to `nil` so that further writes will block
						out1 = nil
					case out2 <- val:
						// once a val is written to channel, we set the copy to `nil` so that further writes will block
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	// repeat generator will repeat values you pass to it infinitely until you tell it to stop
	repeat := func(
		done <-chan interface{},
		val ...interface{}) <-chan interface{} {

		valStream := make(chan interface{})

		go func() {
			defer close(valStream)
			// Infinite loop
			for {
				for _, v := range val {
					select {
					case <-done:
						return
					case valStream <- v:
					}
				}
			}
		}()

		return valStream
	}

	take := func(
		done <-chan interface{},
		valStream <-chan interface{},
		num int) <-chan interface{} {

		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valStream:
				}
			}

		}()

		return takeStream
	}

	// bridge: defines a function that can destructure the channel of channels into a simple channel
	bridge := func(done <-chan interface{},
		chanStream <-chan <-chan interface{}) <-chan interface{} {

		valStream := make(chan interface{})

		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if !ok {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}

				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()

		return valStream
	}

	// genVals: creates a series of 10 channels, each with one element written to them, and passes
	// these channels into the bridge function
	genVals := func() <-chan <-chan interface{} {

		chanStream := make(chan (<-chan interface{}))

		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	// Output
	// out1: 1, out2: 1
	// out1: 2, out2: 2
	// out1: 1, out2: 1
	// out1: 2, out2: 2
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v \n", val1, <-out2)
	}

	for v := range bridge(done, genVals()) {
		fmt.Printf("%v\t", v)
	}

}
