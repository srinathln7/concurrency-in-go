package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func isPrime(n int) bool {
	// check if n is even
	if n&1 == 0 {
		return false
	}

	for i := 3; i <= int(math.Sqrt(float64(n))); i++ {
		if int(n)%i == 0 {
			return false
		}
	}

	return true
}

func main() {

	rand := func() interface{} { return rand.Intn(50000000000) }

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})

		go func() {
			defer close(valStream)

			for {
				select {
				case <-done:
					return
				case valStream <- fn():
				}
			}
		}()

		return valStream
	}

	toInt := func(
		done <-chan interface{},
		valStream <-chan interface{},
	) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)

			for v := range valStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()

		return intStream
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

	primeFinder := func(done <-chan interface{},
		intStream <-chan int) <-chan interface{} {

		resultStream := make(chan interface{})

		go func() {
			defer close(resultStream)

			for num := range intStream {
				if !isPrime(num) {
					continue
				}
				select {
				case <-done:
					return
				case resultStream <- num:
				}

			}
		}()

		return resultStream
	}

	// Fan-in: Process of combining multiple resluts into one channel
	fanIn := func(done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {

		// We wait untill all channels have been drained out
		var wg sync.WaitGroup

		multiplexedStream := make(chan interface{})

		// `multiplex` reads from the channel and pass the value onto the `multiplexedStream` channel
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// select from all channels
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		// Wait for all reads to complete and then close the channel
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		// Results in a deadlock
		// wg.Wait()
		// close(multiplexedStream)

		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	randIntStream := toInt(done, repeatFn(done, rand))

	// Fan-out stage: Process of starting multiple goroutines to handle input
	// from the pupeline.
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders. \n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	// Use Fan-in stage:
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	// Search took: 10.017054ms
	fmt.Printf("Search took: %v\n", time.Since(start))
}
