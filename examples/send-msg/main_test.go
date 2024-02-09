package main

import (
	"context"
	"testing"
	"time"
)

func TestSendSMS(t *testing.T) {
	dir := []string{"0649571937", "0623789027", "0634181723", "0676542312", "06234787125"}
	system := system{
		src:       "127.0.0.1",
		directory: dir,
		msgch:     make(chan []byte),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	system.run(ctx)

	// Add assertions here if needed
}
