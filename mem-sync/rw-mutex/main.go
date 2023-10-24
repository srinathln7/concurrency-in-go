package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// RWMutex: Similar to Mutex but is less stringent. Multiple readers (reader go routines) can hold a shared reader
// lock as long as the lock is not held by writer (writer go routine). This is in contrast to mutex where only
// one go routine can hold the lock at a given moment
func main() {

	// Locker interface has two methods `Lock()` and `Unlock()`
	// Both pointer receivers of Mutex and RWMutex implement the Locker interface
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1) // Sleep to illustrate that it is indeed less active than observer go routines
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,

			// The two lines you provided call the `test` function with different locks (`sync.Locker`) for synchronization. The difference lies in the types of locks used:
			// In summary, the first line uses `RWMutex` with read-only access (allowing multiple concurrent readers), while the second line uses a regular `Mutex` with exclusive access
			// (only one goroutine can hold the lock at a time). These lines test the performance difference between read-heavy workloads (line 1) and workloads with potential contention between readers and writers (line 2).

			// 1. `test(count, &m, m.RLocker())`
			// - This line uses `m.RLocker()` to obtain a read lock from the `sync.RWMutex` `m`. `RLocker()` returns a read-only view of the `sync.RWMutex`, allowing multiple readers to access it concurrently. It's intended for read-only operations and does not provide write capabilities.
			// - This simulates multiple reader goroutines attempting to read data protected by an `RWMutex`. Multiple readers can hold a read lock concurrently, allowing for parallel reads but no writes.
			test(count, &m, m.RLocker()),

			// 2. `test(count, &m, &m)`
			// - This line uses the `&m` reference, which refers to the original `sync.RWMutex` `m`. It passes the mutex itself as the lock to be used for synchronization.
			// - This simulates reader and writer goroutines using a traditional `sync.Mutex`. With a `Mutex`, only one goroutine (either reader or writer) can hold the lock at any given time. This provides strict mutual exclusion, ensuring exclusive access to the protected data.
			test(count, &m, &m),
		)
	}
}
