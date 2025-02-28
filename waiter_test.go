package waiter

import (
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	w := New(50*time.Millisecond, 1000)
	start := time.Now()

	w.Wait(func() {
		t.Log("call 1")
	})

	elapsed1 := time.Since(start)
	t.Log("elapsed 1=", elapsed1)

	w.Wait(func() {
		t.Log("call 2")
	})

	elapsed2 := time.Since(start)
	t.Log("elapsed 2=", elapsed2)
	if elapsed2 < 50*time.Millisecond {
		t.Errorf("elapsed=%v, want > 100ms", elapsed2)
	}
}

func TestCall(t *testing.T) {

	const queueLen, rounds, numInRound = 3, 3, 11

	w := New(50*time.Millisecond, queueLen)

	wg := sync.WaitGroup{}
	start := time.Now()
	lastCall := start
	stopCh := make(chan struct{}, 1)

	for j := range rounds {

		for i := range numInRound {
			wg.Add(1)
			if err := w.Call(func() {
				t.Log("call", i+1, "after", time.Since(lastCall))
				lastCall = time.Now()
				wg.Done()

				// Stop processing after 5 calls
				if j == 1 && i+1 == 5 {
					t.Log("stop after", i+1, "calls")
					if w.Close() {
						t.Log("send to stopCh")
						stopCh <- struct{}{}
					}
				}
			}); err != nil {
				t.Log("call error:", err)
				wg.Done()
				go func() { t.Log("send to stopCh"); stopCh <- struct{}{} }()
				break
			}
		}
		t.Log("all added, in channel", w.Len(), "funcs to call")

		// Wait for all calls processed
		go func() {
			wg.Wait()
			stopCh <- struct{}{}
		}()

		// Wait for all calls processed or Waiter was closed
		<-stopCh

		// Wait for next round
		if j < rounds-1 {
			time.Sleep(100 * time.Millisecond)
			t.Log("\nwaited 100ms ---")
		}
	}

	total := time.Since(start)
	t.Log("done, total time", total)
}
