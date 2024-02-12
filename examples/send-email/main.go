package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

var (
	CONTENT       = "Test email from server"
	MAX_TIMEOUT   = 150
	NETWORK_DELAY = 50
)

type server struct {
	src string
	dir []string
	ch  chan []byte
}

type client struct {
}

type email struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

func newServer(src string, dir []string) *server {
	return &server{
		src: src,
		dir: dir,
		ch:  make(chan []byte),
	}
}

func (s *server) sendEmail(from, to, body string) error {
	log.Printf("sending email from %s to %s\n", from, to)

	// Simulate sending email over the network
	email := email{From: from, To: to, Body: body}
	sPayload, err := json.Marshal(email)
	if err != nil {
		return err
	}

	// Mock delay to send msg over the network
	time.Sleep(time.Duration(NETWORK_DELAY) * time.Millisecond)

	s.ch <- sPayload
	return nil
}

func (s *server) doWork() {
	var wg sync.WaitGroup
	for _, toAddr := range s.dir {
		wg.Add(1)
		go func(toAddr string) {
			defer wg.Done()
			if err := s.sendEmail(s.src, toAddr, CONTENT); err != nil {
				log.Fatalf("error %s sending email from %s to %s \n", err.Error(), s.src, toAddr)
			} else {
				log.Printf("succesfully sent email from %s to %s\n", s.src, toAddr)
			}
		}(toAddr)
	}

	wg.Wait()
}

func (c *client) recvEmail(rCh chan []byte) <-chan email {
	log.Println("invoking recvEmail func")
	recvCh := make(chan email)
	go func() {
		rPayload := <-rCh
		// Deserialize logic
		var email email
		err := json.Unmarshal(rPayload, &email)
		if err != nil {
			log.Fatal("error receiving email.")
		}
		recvCh <- email
	}()

	log.Println("returning client's email")
	return recvCh
}

func main() {
	startTime := time.Now()
	dir := []string{"abc@example.com", "def@exmaple.com", "ghij@example.com", "klm@example.com"}
	server := newServer("127.0.0.1", dir)

	client := client{}
	// Simulate sending msgs over N/w => serialize/deserailize
	// Do work concurrently, handle timeouts etc
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(MAX_TIMEOUT)*time.Millisecond)
	defer cancel()

	// Spin up workers
	go server.doWork()

	for range dir {
		select {
		case <-ctx.Done():
			log.Fatal("request timed out.")
		case email := <-client.recvEmail(server.ch):
			log.Printf("%s received email from %s with body %s\n", email.To, email.From, email.Body)
		}
	}

	log.Printf("All Emails sent succesfully in %vs \n", time.Since(startTime))
}
