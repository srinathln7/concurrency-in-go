package main

import (
	"fmt"
)

// There are 3-critical sections to this program.
// Critical sections are parts or sections of your program that require exclusive access to a shared resource
func main() {
	var data int

	go func() {
		data++
	}()

	if data == 0 {
		fmt.Println("Value of data is 0")
	} else {
		fmt.Println("Value of data is:", data)
	}
}
