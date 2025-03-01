# Waiter

Go package waiter implements an object for waiting a specified delay
time since the last call before calling the next function. This is useful
when needing to call some code with a rate limit.

[![GoDoc](https://godoc.org/github.com/kirill-scherba/waiter?status.svg)](https://godoc.org/github.com/kirill-scherba/waiter/)
[![Go Report Card](https://goreportcard.com/badge/github.com/kirill-scherba/waiter)](https://goreportcard.com/report/github.com/kirill-scherba/waiter)

![waiter](img/waiter.jpg)

## Where to use

For example, if you need to call a function at a rate of 1 per second, you
can use a Waiter with a delay of 1 second. This can be useful in cases where
you need to call some API endpoint at a certain rate set by the API provider.

## Usage example

This example creates a waiter with a delay of 100ms and a queue length of 10.
It then adds 10 calls to the waiter and waits for all of them to finish.
The waiter will only call the function after the delay time, so all the calls
will be called with a delay of 100ms from each other.

```go
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
}
```

Output:

```bash
go run ./example
10:48:18.077676 main.go:22: waiter created
10:48:18.077727 main.go:36: call 0 added to queue
10:48:18.077732 main.go:36: call 1 added to queue
10:48:18.077738 main.go:36: call 2 added to queue
10:48:18.077742 main.go:36: call 3 added to queue
10:48:18.077745 main.go:36: call 4 added to queue
10:48:18.077749 main.go:36: call 5 added to queue
10:48:18.077752 main.go:36: call 6 added to queue
10:48:18.077755 main.go:36: call 7 added to queue
10:48:18.077759 main.go:36: call 8 added to queue
10:48:18.077762 main.go:36: call 9 added to queue
10:48:18.077766 main.go:40: all calls added to queue
10:48:18.177956 main.go:32: call 0 executed
10:48:18.278170 main.go:32: call 1 executed
10:48:18.378393 main.go:32: call 2 executed
10:48:18.478591 main.go:32: call 3 executed
10:48:18.578798 main.go:32: call 4 executed
10:48:18.678997 main.go:32: call 5 executed
10:48:18.779192 main.go:32: call 6 executed
10:48:18.879396 main.go:32: call 7 executed
10:48:18.979600 main.go:32: call 8 executed
10:48:19.079798 main.go:32: call 9 executed
10:48:19.079824 main.go:46: all calls executed
```

Execute this example on Go playground: [https://go.dev/play/p/mne0Z8E9UjK](https://go.dev/play/p/mne0Z8E9UjK)

Or run it on your local machine: `go run ./example`

## Licence

[BSD](LICENSE)
