# Go Async/Await Implementation

A lightweight implementation of async/await pattern in Go using generics and goroutines.

## Overview

This project implements JavaScript-style async/await functionality in Go, allowing asynchronous operations to be handled with a cleaner, more sequential-looking syntax while maintaining Go's concurrent execution model.

## Core Components

### `Future[T]` Type
```go
type Future[T any] struct {
    await func() (T, error)
}
```
A generic container that wraps an asynchronous operation's result. It holds a function that blocks until the result is ready.

### `Async` Function
```go
func Async[T any](f func() (T, error)) *Future[T]
```
- Wraps any function that returns `(T, error)` into a `Future`
- Executes the function in a separate goroutine
- Returns immediately with a `Future` object
- Includes panic recovery to convert panics into errors

### `Await` Method
```go
func (f *Future[T]) Await() (T, error)
```
- Blocks until the async operation completes
- Returns the result and any error (including recovered panics)

## Implementation Logic

1. **Async Execution**: When `Async()` is called, it:
   - Creates a channel for synchronization (`done`)
   - Launches a goroutine to execute the provided function
   - Returns a `Future` immediately without blocking

2. **Panic Handling**: The goroutine includes a deferred recovery function that:
   - Catches any panics from the worker function
   - Converts panics to formatted errors with stack traces
   - Ensures the `done` channel is always closed

3. **Result Retrieval**: The `Await()` method:
   - Blocks on the `done` channel until the goroutine completes
   - Returns the stored result and error

## Example Usage

### Basic Async Operation
```go
futureUser := Async(fetchUserData)
// Do other work while data is being fetched
userData, err := futureUser.Await()
```

### Panic Recovery
```go
fut := Async(workerThatPanics)
val, err := fut.Await()
// err will contain the panic message and stack trace
```

## Demo Functions

- **`fetchUserData()`**: Simulates a 2-second API call
- **`workerThatPanics()`**: Demonstrates panic recovery mechanism

## Key Features

-  Generic type support for any return type
-  Non-blocking async execution
-  Automatic panic recovery with stack traces
-  Clean error propagation
-  Simple, intuitive API similar to JavaScript async/await

## Running the Example

```bash
go run main.go
```

