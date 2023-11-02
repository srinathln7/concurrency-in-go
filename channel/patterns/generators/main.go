package main

import "fmt"

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
}
