package main

import (
	"sync"
	"testing"
)

// Objective is to calculate the context switching time between go routines and
// compare it to OS threads
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	data := make(chan struct{})

	send := func() {
		defer wg.Done()

		// We block until the `begin` channel is closed
		<-begin

		for i := 0; i < b.N; i++ {

			// send an empty anonymous struct (takes no memory) to the data channel
			data <- struct{}{}
		}
	}

	receive := func() {
		defer wg.Done()
		<-begin

		for i := 0; i < b.N; i++ {
			<-data
		}
	}

	wg.Add(2)
	go send()
	go receive()

	// begin the performance timer
	b.StartTimer()

	// close the begin channel and unblock `send` and `receive`
	close(begin)

	wg.Wait()
}
