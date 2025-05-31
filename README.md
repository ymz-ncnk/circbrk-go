# circbrk
circbrk is a Circuit Breaker pattern implementation for Golang.

# How To
```go
package main

import (
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
	// the open state.
	assert.Equal(cb.State(), circbrk.Open)
	// You can also check this with the cb.Open() method:
	assert.Equal(cb.Open(), true)

	// While in the open state, Fail and Success calls have no effect.
	cb.Fail()
	cb.Success()

	time.Sleep(OpenDuration + 200*time.Millisecond)

	// After the open duration, the circuit transitions to the half-open state.
	assert.Equal(cb.State(), circbrk.HalfOpen)

	// In the half-open state, it requires 2 successful call (as configured by
	// WithSuccessThreshold) to return to the closed state.
	cb.Success()
	cb.Success()
	assert.Equal(cb.State(), circbrk.Closed)

	// On failure in the half-open state, the circuit re-enters the open state.
	// cb.Fail()
	// assert.Equal(cb.Open(), true)
}
```