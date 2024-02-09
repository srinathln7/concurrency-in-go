package main

import (
	"fmt"
	"runtime"
	"sync"
)

// To calculate the size of each goroutine - we leverage the fact that go does not
// garbage collect abandoned go routines i.e. go routines that block forever without
// exiting
func main() {

	// memConsumed : runs the garbage collection and returns the total number of memmory obtained from the OS
	memConsumed := func() uint64 {

		// GC runs a garbage collection and blocks the caller until the garbage collection is complete. It may also block the entire program.
		runtime.GC()

		var s runtime.MemStats

		// ReadMemStats populates m with memory allocator statistics. The returned memory allocator statistics are up to date as of the call to ReadMemStats.
		// This is in contrast with a heap profile, which is a snapshot as of the most recently completed garbage collection cycle.
		runtime.ReadMemStats(&s)

		// total bytes of memory obtained from the OS
		return s.Sys

	}

	beforeMemory := memConsumed()

	var wg sync.WaitGroup
	c := make(<-chan interface{})

	var just_block = func() {

		// adding `defer` would result in deadlock because the following line is a blocking condition and
		// then defer statement would never be executed
		wg.Done()

		// block forever
		<-c
	}

	// Create a certain number of go routines and check the memory consumed
	const num_go_routines = 1e1
	wg.Add(num_go_routines)

	for i := 0; i < num_go_routines; i++ {
		go just_block()
	}

	wg.Wait()

	afterMemory := memConsumed()

	// Div by 1000 to convert bytes to kilo bytes
	fmt.Printf("%.3fkb \n", float64(afterMemory-beforeMemory)/num_go_routines/1000)
}
