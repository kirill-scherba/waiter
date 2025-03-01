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
	/*
		waiter created
		call 0 added to queue
		call 1 added to queue
		call 2 added to queue
		call 3 added to queue
		call 4 added to queue
		call 5 added to queue
		call 6 added to queue
		call 7 added to queue
		call 8 added to queue
		call 9 added to queue
		all calls added to queue
		call 0 executed
		call 1 executed
		call 2 executed
		call 3 executed
		call 4 executed
		call 5 executed
		call 6 executed
		call 7 executed
		call 8 executed
		call 9 executed
		all calls executed
	*/
}
