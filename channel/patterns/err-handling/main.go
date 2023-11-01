package main

import (
	"fmt"
	"net/http"
)

// This program states the importance of error handling in the program
func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)

		// Here we see goroutine has been given no choice in the matter.
		go func() {
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}

				select {
				case <-done:
					return
				case responses <- resp:

				}
			}
		}()

		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}

	// Upon running this, we see the program hangs due to bad host
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}

	// Now, let us make sure our concurrent process send their errors to another part of the program
	// that has complete info abt the state of our program which can make more informed deicison about
	// what to todo.

	type Result struct {
		err error
		res *http.Response
	}

	checkStatusCorrect := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)

		go func() {
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{err: err, res: resp}

				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	doneCorrect := make(chan interface{})
	defer close(doneCorrect)

	errCount := 0
	urlsT := []string{"a", "https://www.google.com", "b", "c", "d"}

	for result := range checkStatusCorrect(doneCorrect, urlsT...) {
		if result.err != nil {
			fmt.Printf("error: %v\n", result.err)

			// error handling logic defered to the main goroutine and
			// errors are handled more gracefully
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors")
				break
			}

			continue
		}

		fmt.Printf("Response: %v\n", result.res.Status)
	}
}
