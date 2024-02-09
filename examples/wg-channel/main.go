package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func doWork(wg *sync.WaitGroup, ch chan string, id, iter int) error {
	defer wg.Done()
	//log.Printf(" doing work on behalf of worker %d  -> BEFORE pushing value into channel\n", id)
	ch <- fmt.Sprintf("result from worker %d in iteration %d \n", id, iter)
	//log.Printf(" doing work on behalf of worker %d  -> AFTER pushing value into channel\n", id)
	return nil
}

func worker(wgM *sync.WaitGroup, ch chan string, id int) {
	defer wgM.Done()
	log.Printf("worker %d doing work\n", id)

	// For EXPERIMENT purpose, this will make worker with id `i` sleep for `i` seconds
	time.Sleep(time.Duration(id) * time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go doWork(&wg, ch, id, i+1)
	}
	wg.Wait()
}

func main() {
	startTime := time.Now()

	// Notice the use of only one channel to sync access among multiple goroutines
	ch := make(chan string)

	// When you have to schedule multiple wait groups and synchronise between channels
	// schedule this function collectively in the background
	go func() {
		defer close(ch)
		var wgM sync.WaitGroup
		numOfWorkers := 4
		wgM.Add(numOfWorkers)
		for i := 1; i <= numOfWorkers; i++ {
			// If we do not schedule this worker in a go routine and just
			// call `worker(...)`, it is equivalent to calling it sequentially.
			// Total time would be 10s and then there is LITERALLY no use
			// writing concurrent code.
			go worker(&wgM, ch, i)
		}
		wgM.Wait()
	}()

	for val := range ch {
		fmt.Printf("response -> %s \n", val)
	}

	// Total time would always take the weakest link which in this case worker4 sleeping for 4s.
	fmt.Printf("Total operation took time %v seconds\n", time.Since(startTime))
}
