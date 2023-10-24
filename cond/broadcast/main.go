package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	// subscribe - helper fn that allows us to register functions to handle signals
	// from a condition. Each handler is run on its own go routine and subscribe will not exit until
	// that goroutine is confirmed to be running
	subscribe := func(c *sync.Cond, msg string, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait() // Block here
			fmt.Println(msg)
			fn()
		}()

		goroutineRunning.Wait() // Block here
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, "maximizing window", func() {
		//fmt.Println("Maximizing window")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, "displaying annoying dialog box", func() {
		//fmt.Println("Displaying annoying dialog box.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, "mouse click", func() {
		//fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}
