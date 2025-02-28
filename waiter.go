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
		fnCh:  make(chan func(), queueLen),
	}
	go w.run()
	return w
}

func (w *Waiter) run() {
	for fn := range w.fnCh {
		if w.closed.Load() {
			break
		}
		w.wait()
		fn()
	}
}

func (w *Waiter) wait() {
	now := time.Now()

	if w.last.IsZero() {
		w.last = now
		return
	}

	elapsed := now.Sub(w.last)
	if elapsed < w.delay {
		time.Sleep(w.delay - elapsed)
	}

	w.last = time.Now()
}

func (w *Waiter) Call(fn func()) (err error) {
	if w.closed.Load() {
		err = ErrWaiterClosed
		return
	}

	w.fnCh <- fn
	return
}

func (w *Waiter) Wait(f func()) error {
	done := make(chan error)
	go func() {
		if err := w.Call(func() { f(); done <- nil }); err != nil {
			done <- err
		}
	}()
	return <-done
}

func (w *Waiter) Len() int {
	return len(w.fnCh)
}

func (w *Waiter) Close() bool {
	if !w.closed.CompareAndSwap(false, true) {
		fmt.Println("was already closed")
		return false
	}
	fmt.Println("set closed")

	return true
}
