package main

import (
	"fmt"
	"sync"
	"time"
)

// Task function type
type Task func()

// Worker pool structure
type WorkerPool struct {
	tasks    chan Task
	capacity int
	wg       sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(capacity int) *WorkerPool {
	return &WorkerPool{
		tasks:    make(chan Task),
		capacity: capacity,
	}
}

// Run starts the worker pool
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.capacity; i++ {
		go wp.worker()
	}
}

// worker processes tasks from the task channel
func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task()
		wp.wg.Done()
	}
}

// AddTask adds a task to the task channel
func (wp *WorkerPool) AddTask(task Task) {
	wp.wg.Add(1)
	wp.tasks <- task
}

// Wait waits for all tasks to be completed
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.tasks)
}

// Main function
func main() {
	capacity := 5 // Number of worker goroutines
	workerPool := NewWorkerPool(capacity)

	// Start the worker pool
	workerPool.Run()

	// Example tasks
	for i := 0; i < 50; i++ {
		taskID := i
		workerPool.AddTask(func() {
			fmt.Printf("Processing task %d\n", taskID)
			time.Sleep(1 * time.Second) // Simulate work
		})
	}

	// Wait for all tasks to be completed
	workerPool.Wait()
	fmt.Println("All tasks completed")
}
