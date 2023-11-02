package main

import "fmt"

func main() {

	// `generator` function converts a discrete set of values into a stream of data on a channel
	// This is frequently used when working with pipelines because at the beginning of the pipeline
	// you'll always have some batch of data that you need to convert to a channel
	generator := func(done <-chan interface{}, integers ...int) <-chan int {

		intStream := make(chan int)
		go func() {
			defer close(intStream)

			// ranging over the incoming channel: when the incoming channel is closed, the range will exit
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()

		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {

		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)

			for i := range intStream {
				select {
				case <-done:
					// return
				case multipliedStream <- multiplier * i:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {

		addedStream := make(chan int)
		go func() {
			defer close(addedStream)

			for i := range intStream {
				select {
				case <-done:
					// return
				case addedStream <- additive + i:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})

	// This statement ensures that our program exits cleanly and never leaks goroutines
	// closing the done channel will force the pipeline stage to termiante
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)

	// This statement lets us safely execute concurrently because our inputs and outputs are safe in concurrent contexts
	// Each state of the pipeline is executing concurrently and this means that any stage only need to wait for its inputs
	// and be able to send its outputs. It allows our stages to execute independent of one another for some slice of time
	// In the context of your code, "preemptable" means that at certain points in the process, the execution can be interrupted or stopped,
	// usually to allow other concurrent tasks to run.
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	// Cannot close receive-only channel
	// defer close(pipeline)

	// The statements refer to the concept of "preemptability" at different stages of the pipeline and how it's ensured:
	// 1. **"On the other end of the pipeline, the final stage is ensured preemptability by induction."**
	//    - This means that the last stage of the pipeline, where you're ranging over the `pipeline` channel, is designed to be preemptable.
	//    - The term "induction" here refers to the fact that the preemption mechanism is propagated throughout the pipeline. If any part of
	// the pipeline is preemptable, the final stage (which relies on data from the previous stages) will be preemptable.

	// 2. **"It is preemptable because the channel we're ranging over will be closed when preempted, and therefore our range will break when this occurs."**
	//  - The final stage becomes preemptable because it's ranging over the `pipeline` channel. If this channel is closed, the `for` loop that ranges over it will exit.
	//  Closing the `pipeline` channel signals to the final stage that it should stop processing and exit.

	// 3. **"The final stage is preemptable because the stream we rely on is preemptable."**
	//    - This emphasizes that the final stage inherits its preemptability from the earlier stages in the pipeline. Since the data source (the `pipeline` channel)
	// can be preempted, the final stage, which relies on this data, is also preemptable.

	// In essence, the design of the pipeline ensures that it can be gracefully terminated at any point. Preemptability is achieved through the closing of channels
	// and the usage of `select` statements that check for a signal to stop processing (`done` channel). When preemption occurs, the subsequent stages in the pipeline,
	// including the final one, respond to the closed channels and halt their operations.

	for v := range pipeline {
		fmt.Println(v)
	}

}
