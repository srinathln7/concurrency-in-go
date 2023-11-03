package main

import "testing"

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	// ResetTimer zeroes the elapsed benchmark time and memory allocation counters and deletes user-reported metrics. It does not affect whether the timer is running.
	b.ResetTimer()

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

	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

// Type-specific stages are twice as fast, but only marginally faster in magnitude
func BenchMarkTyped(b *testing.B) {

	// repeat generator will repeat values you pass to it infinitely until you tell it to stop
	repeat := func(
		done <-chan interface{},
		val ...string) <-chan string {

		valStream := make(chan string)

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
		valStream <-chan string,
		num int) <-chan string {

		takeStream := make(chan string)

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

	// ResetTimer zeroes the elapsed benchmark time and memory allocation counters and deletes user-reported metrics. It does not affect whether the timer is running.
	b.ResetTimer()

	for range take(done, repeat(done, "a"), b.N) {
	}
}
