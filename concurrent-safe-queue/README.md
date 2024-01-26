# Concurrent Safe Queue

Key Points:
- Utilize two wait groups, `wgE` for enqueue and `wgD` for dequeue operations, to maintain a clear separation of concerns.
- Alternatively, if opting for a single wait group (`wg`), ensure a careful sequence of concurrent enqueue operations followed by concurrent dequeue operations.
- Caution: Failure to call `wg.Wait()` after enqueue operations in a single wait group instance may lead to deadlock.
  - For instance, if 2 dequeue operations follow 1 enqueue operation without calling `wg.Wait()`, the second dequeue operation may return prematurely, leaving the main goroutine waiting indefinitely, resulting in a deadlock.