## Concurrent Worker Example

This Go program demonstrates the use of goroutines and channels for concurrent processing of tasks using multiple workers. The code showcases how to schedule multiple wait groups and synchronize between channels effectively.

### Features:

- **Concurrent Workers**: The program creates multiple worker goroutines to perform tasks concurrently.
- **Synchronization**: It utilizes wait groups to ensure synchronization between the main goroutine and worker goroutines.
- **Channel Communication**: Communication between the main goroutine and worker goroutines is facilitated through channels.
- **Experimentation**: For experimental purposes, each worker sleeps for a duration based on its ID, showcasing the flexibility of concurrent execution.
- **Efficiency**: By employing goroutines and channels, the program maximizes CPU utilization and overall efficiency.

### Usage:

To run the program, simply execute the following command:

```bash
go run main.go
```

### Example Output:

```
worker 1 doing work
worker 2 doing work
worker 3 doing work
worker 4 doing work
response -> result from worker 1 in iteration 1 
response -> result from worker 1 in iteration 2 
response -> result from worker 2 in iteration 1 
response -> result from worker 2 in iteration 2 
response -> result from worker 3 in iteration 1 
response -> result from worker 3 in iteration 2 
response -> result from worker 4 in iteration 1 
response -> result from worker 4 in iteration 2 
Total operation took time 4.006730904 seconds
```

The output demonstrates the concurrent execution of tasks by multiple workers, each producing results that are received and printed by the main goroutine. The total operation time reflects the duration of the longest-running worker, ensuring efficient utilization of resources.