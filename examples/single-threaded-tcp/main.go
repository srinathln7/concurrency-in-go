package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type server struct {
	numOfConn uint
	maxConn   uint
	file      *os.File
	writer    *csv.Writer
	mu        sync.RWMutex
	shutdown  chan os.Signal
}

// Initializes the CSV file
func (s *server) initCSV() {
	file, err := os.Create("status_code.csv")
	if err != nil {
		log.Fatal("error creating json file")
	}

	s.writer = csv.NewWriter(file)
	s.writer.Write([]string{"url", "http_status_code"})
}

// Handles TCP connection
func (s *server) handerConn(conn net.Conn) {
	log.Println("starting to run scheduled tcp handler")
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read the input from the client
		url, err := reader.ReadString('\n')
		if err != nil {
			// If an error occurs, check if it's an EOF, indicating the client closed the connection
			if err == io.EOF {
				log.Println("client has terminated the connection")
				// Decrement the number of connected clients in a thread-safe manner
				s.mu.Lock()
				if s.numOfConn > 0 {
					s.numOfConn--
				}
				log.Printf("number of connected clients: %d", s.numOfConn)
				s.mu.Unlock()
			} else {
				log.Printf("error reading from client: %v", err)
			}
			return
		}

		url = strings.TrimSpace(url)

		// Process the received URL
		go s.fetchURL(url)
	}
}

// Fetches URL and writes its status code to the CSV file
func (s *server) fetchURL(url string) {
	log.Println("starting to run SCHEDULED fetchURL from the internet")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error %s while fetching URL %s from the internet", err, url)
		return
	}

	defer resp.Body.Close()

	// Write the response to the CSV file. Locks ensure concurrent-safe writes
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writer.Write([]string{url, fmt.Sprintf("%d", resp.StatusCode)})
	s.writer.Flush()
}

func main() {
	server := server{
		maxConn:  10,
		shutdown: make(chan os.Signal, 1), // buffered channel to avoid deadlock
	}
	server.initCSV()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("error starting TCP server")
	}

	log.Printf("started TCP server at %s\n", ln.Addr().String())
	defer ln.Close()

	// Create a wait group to wait for all connections to finish
	var wg sync.WaitGroup
	go func() {
		for {

			server.mu.RLock()
			if server.numOfConn >= server.maxConn {
				server.mu.RUnlock()
				// Optionally, you can add a sleep here to prevent a tight loop when max connections are reached
				continue
			}
			server.mu.RUnlock()

			conn, err := ln.Accept() // Blocking call
			if err != nil {
				select {
				case <-server.shutdown:
					log.Println("aborting all client connections and initiate graceful server shutdown")
					wg.Wait()
					os.Exit(0)
				default:
					log.Fatal("error connecting to TCP server")
				}
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				server.handerConn(conn)
			}()

			server.mu.Lock()
			server.numOfConn++
			log.Printf("number of connected clients: %d", server.numOfConn)
			server.mu.Unlock()
		}
	}()

	// Implement graceful shutdown
	signal.Notify(server.shutdown, os.Interrupt, syscall.SIGTERM)
	<-server.shutdown
	close(server.shutdown)
}
