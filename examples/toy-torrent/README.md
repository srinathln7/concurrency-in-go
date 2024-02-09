# Concurrent Chunk Download and Verification

This Go program demonstrates a simple concurrent chunk download and verification process. It simulates a scenario where data is downloaded from multiple peers concurrently and verifies the integrity of the downloaded chunks.

## Overview

The program defines a `Peer` structure representing a peer with an `id` and a `chunk` of data. It initializes two peers with chunks of data and concurrently downloads these chunks using goroutines. The downloaded chunks are then verified for integrity.

## Code Structure

- **Peer Struct:**
  - `id`: Identifier for the peer.
  - `chunk`: Represents a chunk of data.

- **Constants:**
  - `TOTAL_CHUNKS`: Total number of chunks to download.
  - `CHUNKS`: Slice containing simulated chunks of data.

- **Main Function:**
  - Initializes two peers with chunks.
  - Creates a buffered channel to communicate downloaded chunks.
  - Concurrently downloads chunks from peers using goroutines.
  - Waits for all downloads to complete using a `sync.WaitGroup`.
  - Closes the chunk channel after downloads are completed.
  - Verifies the integrity of downloaded chunks.

- **Download and Verification Functions:**
  - `downloadFromPeer`: Simulates downloading a chunk from a peer.
  - `verifyChunk`: Verifies the integrity of downloaded chunks.

## Execution

To run the program:

```bash
go run main.go
```

## Notes

- The program uses goroutines and channels to achieve concurrency.
- The `sync.WaitGroup` ensures that the program waits for all downloads to complete before proceeding to verification.
- The chunk channel is closed after all downloads are completed.

