package main

import (
	"fmt"
	"sync"
)

func main() {

	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instance.")
			return struct{}{}
		},
	}

	//myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()

	// Limit memory usage example

	var numCalcsCreated int

	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// Seed the pool with 4KB i.e. 4 calcs
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024 // 1GB
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {

			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

			// Do something quick and interesting with memory
		}()
	}

	wg.Wait()

	// Only 4 calculators are created.Just used only 4KB of space.
	// Without pools, worst case sceanario is using 1GB.
	fmt.Printf("%d calculators were created. \n", numCalcsCreated)
}
