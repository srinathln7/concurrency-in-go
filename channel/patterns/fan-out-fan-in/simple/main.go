package main

import (
	"fmt"
	"math"
	"math/rand"
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

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	// Search took: 47.187403ms
	fmt.Printf("Search took: %v\n", time.Since(start))
}
