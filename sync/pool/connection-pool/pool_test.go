package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
)

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}

	for i := 0; i < 10; i++ {
		p.Put(p.New)
	}

	return p
}

func startNetworkDaemonPool() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}

		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}

			// Rather than opening a connection to another service for every request accepted, use the connection pool instead
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()

	return &wg

}

func init() {
	daemonStarted := startNetworkDaemonPool()
	daemonStarted.Wait()
}

func BenchmarkNetworkRequestWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", ":8081")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}

		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}

		conn.Close()
	}
}
