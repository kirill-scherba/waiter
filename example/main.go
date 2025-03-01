// This example creates a waiter with a delay of 100ms and a queue length of 10.
// It then adds 10 calls to the waiter and waits for all of them to finish.
// The waiter will only call the function after the delay time, so all the calls
// will be called with a delay of 100ms from each other.
package main

import (
	"log"
	"sync"
	"time"

	"github.com/kirill-scherba/waiter"
)

func main() {

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	// Create a waiter with a delay of 100ms and a queue length of 10
	w := waiter.New(100*time.Millisecond, 10)
	log.Println("waiter created")

	// Create a wait group to wait for all the calls to finish
	wg := sync.WaitGroup{}

	// Add 10 calls to the waiter
	for i := range 10 {
		wg.Add(1)
		// Call the function after the delay time
		w.Call(func() {
			log.Println("call", i, "executed")
			wg.Done()
		})
		// Log that the call was added to the queue
		log.Println("call", i, "added to queue")
	}

	// Log that all the calls were added to the queue
	log.Println("all calls added to queue")

	// Wait for all the calls to finish
	wg.Wait()

	// Log that all the calls were executed
	log.Println("all calls executed")

	// Output:
	/*
		10:54:10.169111 main.go:22: waiter created
		10:54:10.169158 main.go:36: call 0 added to queue
		10:54:10.169163 main.go:36: call 1 added to queue
		10:54:10.169167 main.go:36: call 2 added to queue
		10:54:10.169170 main.go:36: call 3 added to queue
		10:54:10.169174 main.go:36: call 4 added to queue
		10:54:10.169177 main.go:36: call 5 added to queue
		10:54:10.169179 main.go:36: call 6 added to queue
		10:54:10.169182 main.go:36: call 7 added to queue
		10:54:10.169184 main.go:36: call 8 added to queue
		10:54:10.169187 main.go:36: call 9 added to queue
		10:54:10.169189 main.go:40: all calls added to queue
		10:54:10.269362 main.go:32: call 0 executed
		10:54:10.369580 main.go:32: call 1 executed
		10:54:10.469782 main.go:32: call 2 executed
		10:54:10.569975 main.go:32: call 3 executed
		10:54:10.670172 main.go:32: call 4 executed
		10:54:10.770371 main.go:32: call 5 executed
		10:54:10.870579 main.go:32: call 6 executed
		10:54:10.970773 main.go:32: call 7 executed
		10:54:11.070914 main.go:32: call 8 executed
		10:54:11.171110 main.go:32: call 9 executed
		10:54:11.171146 main.go:46: all calls executed
	*/
}
