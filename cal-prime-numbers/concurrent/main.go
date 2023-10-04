package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var N uint64 = 10
var COUNT_PRIME uint64 = 1 // Exclude 2
var TOTAL_GOROUTINE_COUNT = 10
var CURRENT_NUM uint64 = 2

func isPrime(n uint64) bool {

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

// Want to use goroutines to spd up calc
func doWork(wg *sync.WaitGroup, thread_id int) {

	start_time := time.Now()
	defer wg.Done()

	for {
		x := atomic.AddUint64(&CURRENT_NUM, 1)

		fmt.Printf("\n go routine = %v, CURRENT_NUM= %v \n", thread_id, CURRENT_NUM)
		// Breaking condition
		if x > N {
			break
		}

		if isPrime(x) {
			atomic.AddUint64(&COUNT_PRIME, 1)
		}
	}

	fmt.Printf("\n GO ROUTINE %v TOOK %v \n", thread_id, time.Since(start_time))
}

func main() {

	start_time := time.Now()

	var wg sync.WaitGroup
	for i := 1; i <= TOTAL_GOROUTINE_COUNT; i++ {
		wg.Add(1)
		go doWork(&wg, i)
	}

	wg.Wait()

	fmt.Printf("Counted %v prime numbers between 1 and %v in %v \n", COUNT_PRIME, N, time.Since(start_time))
}
