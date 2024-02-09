package main

import (
	"fmt"
	"math"
	"time"
)

var N = 1000000
var COUNT_PRIME = 1 // Exclude 2
var CONCURRENCY = 10

func isPrime(n int) bool {

	// check if n is even
	if n&1 == 0 {
		return false
	}

	for i := 3; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func main() {

	start_time := time.Now()

	for i := 3; i <= N; i++ {
		if isPrime(i) {
			COUNT_PRIME++
		}
	}

	fmt.Printf("Counted %v prime numbers between 1 and %v in %v \n", COUNT_PRIME, N, time.Since(start_time))
}
