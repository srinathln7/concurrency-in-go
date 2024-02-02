package main

import (
	"context"
	"log"
	"time"
)

type Response struct {
	msg string
	err error
}

// fetchSlowAPI: exhibits non-deterministic behaviour with different response time
func fetchSlowAPI(respCh chan Response) {
	log.Println("fetching slow API")
	time.Sleep(101 * time.Millisecond)

	respCh <- Response{msg: "response from slow API", err: nil}
}

// fetchAPI: Tackles the `fetchSlowAPI` request appropriately
func fetchAPI(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	respCh := make(chan Response)
	go func() {
		fetchSlowAPI(respCh)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Fatal("request timed out")
		case resp := <-respCh:
			if resp.err != nil {
				log.Printf("response from slow api - resp %s with err %s", resp.msg, resp.err.Error())
			} else {
				log.Printf("response from slow api - resp %s with no error", resp.msg)
			}
			return
		}
	}
}

func main() {
	ctx := context.Background()
	fetchAPI(ctx)
}
