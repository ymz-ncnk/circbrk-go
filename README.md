# circbrk-go
circbrk-go is a Circuit Breaker Pattern implementation for Go.

# How To
```go
package main

import (
	"sync"
	"time"

	assert "github.com/ymz-ncnk/assert/panic"
	circbrk "github.com/ymz-ncnk/circbrk-go"
)

func init() {
	assert.On = true
}

const OpenDuration = time.Second

func main() {
	cb := circbrk.New(circbrk.WithWindowSize(4),
		circbrk.WithFailureRate(0.5),
		circbrk.WithOpenDuration(OpenDuration),
		circbrk.WithSuccessThreshold(2),
	)
	// The circuit breaker starts in the closed state.
	assert.Equal(cb.State(), circbrk.Closed)

	cb.Fail()
	cb.Success()
	cb.Success()
	cb.Fail()

	// After several calls the failure rate (in this case, 2 failures out of 4
	// calls = 0.5) meets the configured threshold, so the circuit transitions to
	// the open state and becomes blocked.
	assert.Equal(cb.State(), circbrk.Open)
	assert.Equal(cb.Blocked(), true)

	// While in the open state, Fail() and Success() calls have no effect.
	cb.Fail()
	cb.Success()

	time.Sleep(OpenDuration + 200*time.Millisecond)

	// After the open duration, the circuit transitions to the half-open state.
	assert.Equal(cb.State(), circbrk.HalfOpen)

	// In the half-open state, 2 successful calls (as configured by
	// WithSuccessThreshold) are required to transition back to closed.
	assert.Equal(cb.Blocked(), false)
	// The breaker is not blocked, so we can execute the request and
	// notify it of a successful outcome.
	cb.Success()

	assert.Equal(cb.Blocked(), false)
	// The number of allowed attempts equals the number of required successes.
	// Any additional calls to cb.Blocked() in half-open state will return true
	// (blocked).
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		assert.Equal(cb.Blocked(), true)
		assert.Equal(cb.Blocked(), true)
		wg.Done()
	}()
	wg.Wait()
	cb.Success()

	assert.Equal(cb.State(), circbrk.Closed)
	assert.Equal(cb.Blocked(), false)

	// On any failure in the half-open state, the circuit re-enters the open
  // state.
	// cb.Fail()
	// assert.Equal(cb.Open(), true)
	// assert.Equal(cb.Blocked(), true)
}
```
