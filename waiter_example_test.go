package waiter

import (
	"fmt"
	"sync"
	"time"
)

// ExampleWaiter creates a waiter with a delay of 100ms and a queue length of 10.
// It then adds 10 calls to the waiter and waits for all of them to finish.
// The waiter will only call the function after the delay time, so all the calls
// will be called with a delay of 100ms from each other.
func ExampleWaiter() {

	// Create a waiter with a delay of 100ms and a queue length of 10
	w := New(100*time.Millisecond, 10)
	fmt.Println("waiter created")

	// Create a wait group to wait for all the calls to finish
	wg := sync.WaitGroup{}

	// Add 10 calls to the waiter
	for i := range 10 {
		wg.Add(1)

		// Call the function after the delay time
		w.Call(func() {
			fmt.Println("call", i, "executed")
			wg.Done()
		})

		// Log that the call was added to the queue
		fmt.Println("call", i, "added to queue")
	}

	// Log that all the calls were added to the queue
	fmt.Println("all calls added to queue")

	// Wait for all the calls to finish
	wg.Wait()

	// Log that all the calls were executed
	fmt.Println("all calls executed")

	// Output:
	// waiter created
	// call 0 added to queue
	// call 1 added to queue
	// call 2 added to queue
	// call 3 added to queue
	// call 4 added to queue
	// call 5 added to queue
	// call 6 added to queue
	// call 7 added to queue
	// call 8 added to queue
	// call 9 added to queue
	// all calls added to queue
	// call 0 executed
	// call 1 executed
	// call 2 executed
	// call 3 executed
	// call 4 executed
	// call 5 executed
	// call 6 executed
	// call 7 executed
	// call 8 executed
	// call 9 executed
	// all calls executed
}

// ExampleNew creates a waiter using the New function, with a delay of 100ms and
// queue length 10.
func ExampleNew() {

	// Create a waiter using the New function
	w := New(100*time.Millisecond, 10)
	fmt.Println("waiter created", w.delay)

	// Output:
	// waiter created 100ms
}

// ExampleRateLimit creates a waiter using the the New function and RateLimit
// function in parameters, with a delay equal to 100 calls per second, and queue
// length 10.
func ExampleRateLimit() {

	// Create a waiter using the RateLimit function
	w := New(RateLimit(100, 1*time.Second), 10)
	fmt.Println("waiter created", w.delay)

	// Output:
	// waiter created 10ms
}

// ExampleWaiter_Call calls the Waiter.Call function.
//
// This code  demonstrates the usage of a Waiter object's Call function. It
// creates a Waiter with a 100ms delay and a queue length of 10, then schedules
// a function call that prints "call executed" after the delay. The code uses a
// WaitGroup to wait for the scheduled function to complete before exiting.
func ExampleWaiter_Call() {

	// Create a waiter
	w := New(100*time.Millisecond, 10)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Call the Waiter.Call function
	w.Call(func() {
		fmt.Println("call executed")
		wg.Done()
	})

	wg.Wait()

	// Output:
	// call executed
}

// ExampleWaiter_Wait calls the Waiter.Wait function.
//
// This code demonstrates the usage of a Waiter object's Call function. It
// creates a Waiter with a 100ms delay and a queue length of 10, then uses it
// to schedule a function call that prints "call executed" after the delay. The
// code uses a WaitGroup to wait for the scheduled function to complete before exiting.
func ExampleWaiter_Wait() {

	// Create a waiter
	w := New(100*time.Millisecond, 10)

	// Call the Waiter.Wait function
	w.Wait(func() {
		fmt.Println("call executed")
	})

	// Output:
	// call executed
}

// ExampleWaiter_Len calls the Waiter.Len function.
//
// This code  demonstrates the usage of the Waiter.Len function. It creates a
// new Waiter with a 100ms delay and a queue length of 10, and then prints the
// length of the waiter's queue, which is initially 0.
func ExampleWaiter_Len() {

	// Create a waiter
	w := New(100*time.Millisecond, 10)

	// Call the Waiter.Len function
	fmt.Println("len =", w.Len())

	// Output:
	// len = 0
}
