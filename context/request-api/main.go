package main

import (
	"context"
	"log"
	"strconv"
	"time"
)

type Response struct {
	msg string
	err error
}

// fetchSlowAPI: exhibits non-deterministic behaviour with different response time
func fetchSlowAPI(respCh chan Response) {
	log.Println("fetching slow API")

	// success sceanario for handling reqs
	time.Sleep(10 * time.Millisecond)

	// req-timeout scenario
	// time.Sleep(10 * time.Millisecond)

	respCh <- Response{msg: "response from slow API", err: nil}
}

// fetchAPI: Tackles the `fetchSlowAPI` request appropriately
func fetchAPI(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Millisecond)
	defer cancel()

	reqID := ctx.Value(ctxReqID).(string)
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
				log.Printf("response for reqid=%s from slow api: %s with no error", reqID, resp.msg)
			}
			return
		}
	}
}

type ctxKey int

const (
	ctxReqID ctxKey = iota
)

func main() {
	// Just for sample use
	ctx := context.WithValue(context.Background(), ctxReqID, strconv.Itoa(int(ctxReqID)))
	fetchAPI(ctx)
}
