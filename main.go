package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

type Future[T any] struct {
	await func() (T, error)
}

func (f *Future[T]) Await() (T, error) {
	return f.await()
}

func Async[T any](f func() (T, error)) *Future[T] {
	var result T
	var err error

	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					err = fmt.Errorf(
						"panic in worker: %w\n%s", x, debug.Stack(),
					)

				default:
					err = fmt.Errorf(
						"panic in worker: %v\n%s", x, debug.Stack(),
					)
				}
			}
			close(done)
		}()

		result, err = f()
	}()

	return &Future[T]{
		await: func() (T, error) {
			<-done
			return result, err
		},
	}
}

func fetchUserData() (string, error) {
	fmt.Println("=> Starting to fetch user data...")
	time.Sleep(2 * time.Second)
	fmt.Println("=> Finished fetching user data...")
	return "John Doe", nil
}

func workerThatPanics() (string, error) {
	fmt.Println("=> Starting risky work (still work)...")
	time.Sleep(300 * time.Millisecond)
	panic("something went wrong")
	return "", nil
}

func main() {
	futureUser := Async(fetchUserData)

	fmt.Println("Doing other work while fetching data...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Doing work...")

	userData, err := futureUser.Await()
	if err != nil {
		fmt.Printf("Error while getting user data: %v\n", err)
	} else {
		fmt.Printf("User data: %s\n", userData)
	}

	fut := Async(workerThatPanics)
	val, err := fut.Await()

	if err != nil {
		fmt.Printf("Awaited Error: %v\n", err)
	} else {
		fmt.Printf("Success: %s\n", val)
	}
}
