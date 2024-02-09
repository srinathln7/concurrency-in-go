package main

import (
	"log"
	"sync"
)

type Peer struct {
	id    string
	chunk []byte
}

var TOTAL_CHUNKS = 2
var CHUNKS = [][]byte{
	[]byte("Hello there!"),
	[]byte("Welcome to Torrent"),
}

var chunksDownloaded int

func main() {
	peer1 := Peer{id: "1", chunk: CHUNKS[0]}
	peer2 := Peer{id: "2", chunk: CHUNKS[1]}
	var peers []Peer = []Peer{peer1, peer2}

	var wg sync.WaitGroup

	// Create a buffered channel with the capacity of total chunks
	chunkChannel := make(chan []byte, TOTAL_CHUNKS)

	log.Println("beginning the download process")
	// Download the data from the other peers
	for _, peer := range peers {
		log.Printf("adding one wait group to download from peer %s \n", peer.id)
		wg.Add(1)
		go func(peer Peer) {
			log.Printf("starting download from peer %s \n", peer.id)
			defer wg.Done()
			downloadFromPeer(peer, chunkChannel)
		}(peer)
	}

	go func() {
		log.Println("holding until all the wait groups are completed")
		wg.Wait()

		log.Println("closing chunk channel")
		close(chunkChannel)
	}()

	// log.Println("holding until all the wait groups are completed")
	// wg.Wait()
	// log.Println("defering closing chunk channel")
	// close(chunkChannel)

	log.Println("starting from verify the chunks")
	verifyChunk(chunkChannel)
	log.Println("downloaded and verfication successfully completed!!!")
}

func downloadFromPeer(peer Peer, chunkCh chan []byte) {
	log.Printf("downloading from peer: %v \n", peer.id)
	chunkCh <- peer.chunk
	log.Printf("downloaded from peer: %v \n", peer.id)
}

func verifyChunk(chunkCh chan []byte) {

	// Blocking call
	for recvChunk := range chunkCh {
		log.Printf("verified integrity of chunk \"%s\"", string(recvChunk))
	}
}
