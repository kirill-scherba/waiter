// Copyright 2025 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package waiter implements an object for waiting a specified delay
// time since the last call before calling the next function. This is useful
// when needing to call some code with a rate limit.
//
// For example, if you need to call a function at a rate of 1 per second, you
// can use a Waiter with a delay of 1 second. This can be useful in cases where
// you need to call some API endpoint at a certain rate set by the API provider.
package waiter

import (
	"fmt"
	"sync/atomic"
	"time"
)

var ErrWaiterClosed = fmt.Errorf("waiter is closed")

// Waiter represents an object for waiting a specified delay time since the
// last call before calling the next function. This is useful when needing to
// call some code with a rate limit.
type Waiter struct {
	// delay is the time to wait between calls.
	delay time.Duration

	// last is the time of the last call.
	last time.Time

	// fnCh is a channel of functions to call.
	fnCh chan func()

	// closed is a flag to indicate if the waiter is closed.
	closed atomic.Bool
}

// New creates a new Waiter object.
//
// The function creates a new Waiter object and starts running it in a goroutine.
// The Waiter object is used to wait a specified delay time since the last call
// before calling the next function. This is useful when needing to call some code
// with a rate limit.
func New(delay time.Duration, queueLen int) *Waiter {
	w := &Waiter{
		delay: delay,
		last:  time.Now(),
		fnCh:  make(chan func(), queueLen),
	}
	go w.run()
	return w
}

// RateLimit returns the time to wait between calls based on the specified
// quantity and time delay. The return value is the time to wait between calls,
// it uses in the New function call.
//
// The rate limit is calculated as follows:
//
//	delay / quantity
//
// For example, if you need to call a function at a rate of 100 per second, you
// can use RateLimit(100, 1*time.Second) in the New function call.
func RateLimit(quantity int, delay time.Duration) time.Duration {
	return delay / time.Duration(quantity)
}

// Call calls the specified function after waiting the specified delay time
// since the last call.
//
// If the Waiter is closed, the function will return ErrWaiterClosed.
func (w *Waiter) Call(fn func()) (err error) {
	if w.closed.Load() {
		// If the Waiter is closed, return ErrWaiterClosed
		err = ErrWaiterClosed
		return
	}

	// Add the function to the channel of functions to call
	w.fnCh <- fn
	return
}

// Wait calls the specified function after waiting the specified delay time
// since the last call.
//
// This function is similar to Call but it waits until the specified function
// is called and returns any error that occurred.
func (w *Waiter) Wait(f func()) error {
	// Create a channel to receive the error
	done := make(chan error)

	// Start a new goroutine to call the function
	go func() {
		// Call the function with the specified delay
		if err := w.Call(func() {
			// Call the function
			f()

			// Send the error to the channel
			done <- nil
		}); err != nil {
			// If there is an error, send it to the channel
			done <- err
		}
	}()

	// Wait until the f function is called and error is received
	// from the done channel
	return <-done
}

// Len returns the number of functions currently waiting in the channel.
func (w *Waiter) Len() int {
	return len(w.fnCh)
}

// Close closes the Waiter and stops it from calling any more functions.
//
// If the Waiter is already closed, the function will return ErrWaiterClosed
// error.
func (w *Waiter) Close() (err error) {
	// Set the closed flag to true
	if !w.closed.CompareAndSwap(false, true) {
		// If the flag is already true, return ErrWaiterClosed
		err = ErrWaiterClosed
	}
	return
}

// run starts a new goroutine to run the Waiter object.
// It loops through the channel of functions to call and calls them with the
// specified delay.
func (w *Waiter) run() {
	// Loop through the channel of functions to call
	for fn := range w.fnCh {
		// If the Waiter is closed, exit the loop
		if w.closed.Load() {
			break
		}

		// Wait the specified delay before calling the function
		w.wait()

		// Call the function
		fn()
	}
}

// wait waits the specified delay time since the last call before calling the
// next function.
func (w *Waiter) wait() {
	// Get the current time
	now := time.Now()

	// If the last call time is zero, set it to the current time
	if w.last.IsZero() {
		w.last = now
		return
	}

	// Calculate the time elapsed since the last call
	elapsed := now.Sub(w.last)

	// If the elapsed time is less than the delay, sleep for the difference
	if elapsed < w.delay {
		time.Sleep(w.delay - elapsed)
	}

	// Update the last call time
	w.last = time.Now()
}
