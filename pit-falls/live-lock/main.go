package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Livelocks: Refer to a situation in which a program performs concurrent operations and progress individually
// but the combined operations do not collectively progress the overall state of the program
func main() {

	// A cadence is a rhythm, or a flow of words or music, in a sequence that is regular (or steady as it were)
	cadence := sync.NewCond(&sync.Mutex{})

	// Each go routine needs to wake up at the same
	go func() {
		for range time.Tick(1 * time.Millisecond) {

			// Broadcast wakes all goroutines waiting on c.
			// It is allowed but not required for the caller to hold c.L during the call.
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDirn := func(dirn string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, "%v", dirn)
		atomic.AddInt32(dir, 1) // declare the intention to move in a direction
		fmt.Println("No. of ppl walking in the same direction: ", *dir)
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}

		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	var tryLeft = func(out *bytes.Buffer) bool { return tryDirn("left ", &left, out) }
	var tryRight = func(out *bytes.Buffer) bool { return tryDirn("right ", &right, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}

			fmt.Fprintf(&out, "\n%v tosses his/her hands up in exasperation!", name)
		}
	}

	var pplInHallway sync.WaitGroup
	pplInHallway.Add(2)
	go walk(&pplInHallway, "Srinath")
	go walk(&pplInHallway, "Srinidhi")
	pplInHallway.Wait()
}
