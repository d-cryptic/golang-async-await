package main

import (
	"fmt"
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

	return &Future[T]{
		await: func() (T, error) {
			return result, err
		}
	}
}

func fetchUserData() (string, error) {
	fmt.Println("=> Starting to fetch user data...")
	time.Sleep(2 * time.Second)
	fmt.Println("=> Finished fetching user data...")
	return "User Data: John Doe", nil
}

func main() {
	futureUser := Async(fetchUserData)

	fmt.Println("Doing other work while fetching data...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Doing work...")

	userData, err := futureUser.Await()
	if err != nil {
		fmt.Println("Error while getting user data: %v\n", err)
	} else {
		fmt.Println("User data: %s\n", userData)
	}
}
