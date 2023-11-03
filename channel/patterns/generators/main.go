package main

import (
	"fmt"
	"math/rand"
)

func main() {

	// repeat generator will repeat values you pass to it infinitely until you tell it to stop
	repeat := func(
		done <-chan interface{},
		val ...interface{}) <-chan interface{} {

		valStream := make(chan interface{})

		go func() {
			defer close(valStream)
			// Infinite loop
			for {
				for _, v := range val {
					select {
					case <-done:
						return
					case valStream <- v:
					}
				}
			}
		}()

		return valStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})

		go func() {
			defer close(valStream)

			for {
				select {
				case <-done:
					return
				case valStream <- fn():
				}
			}
		}()

		return valStream
	}

	toString := func(
		done <-chan interface{},
		valStream <-chan interface{},
	) <-chan string {
		strStream := make(chan string)

		go func() {
			defer close(strStream)

			for v := range valStream {
				select {
				case <-done:
					return
				case strStream <- v.(string):
				}
			}
		}()

		return strStream
	}

	// This pipeline stage will only take the first `num` items off of its incoming valueStream and then exit
	take := func(
		done <-chan interface{},
		valStream <-chan interface{},
		num int) <-chan interface{} {

		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valStream:
				}
			}

		}()

		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	// combine the repeat and take generators
	for val := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v \n", val)
	}

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}

	var msg string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		msg += token
	}

	fmt.Printf("message: %s...\n", msg)
}
