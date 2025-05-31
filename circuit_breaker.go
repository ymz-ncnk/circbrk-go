package circbrk

import (
	"sync"
	"time"
)

// New creates a new CircuitBreaker instance with the provided options, used to
// configure parameters such as window size, failure rate, open duration, and
// success threshold.
func New(ops ...SetOption) *CircuitBreaker {
	o := Options{}
	Apply(ops, &o)

	window := make([]bool, o.WindowSize)
	for i := range window {
		window[i] = true
	}
	return &CircuitBreaker{
		window:  window,
		state:   Closed,
		mu:      sync.Mutex{},
		options: o,
	}
}

// CircuitBreaker implements the Circuit Breaker pattern.
type CircuitBreaker struct {
	state State

	window   []bool
	index    int
	failures int
	timer    *time.Timer

	trialCount int

	mu      sync.Mutex
	options Options
}

// Open returns true if the circuit breaker is in Open state.
// This indicates that operations should fail fast without attempting execution.
func (b *CircuitBreaker) Open() bool {
	return b.State() == Open
}

// State returns the current state of the circuit breaker.
func (b *CircuitBreaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.state
}

// Success records a successful operation and updates the circuit breaker state
// accordingly. If configured, it triggers the success callback.
func (b *CircuitBreaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.options.SuccessCallback != nil {
		b.options.SuccessCallback(b.state)
	}

	switch b.state {
	case Closed:
		b.sucessOnClosed()
	case HalfOpen:
		b.trialCount++
		b.checkHalfOpen()
	case Open:
		return
	}
}

// Fail records a failed operation and updates the circuit breaker state
// accordingly. If configured, it triggers the fail callback.
func (b *CircuitBreaker) Fail() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.options.FailCallback != nil {
		b.options.FailCallback(b.state)
	}

	switch b.state {
	case Closed:
		b.failOnClosed()
	case HalfOpen:
		b.trialCount = 0
		b.trip()
	case Open:
		return
	}
}

// ResetNow immediately resets the circuit breaker to Closed state. It stops
// any pending timer and clears all failure counts. Useful for manual reset
// scenarios.
func (b *CircuitBreaker) ResetNow() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.timer != nil {
		b.timer.Stop()
	}
	b.reset()
}

func (b *CircuitBreaker) sucessOnClosed() {
	oldOk := b.window[b.index]
	if !oldOk {
		b.failures--
		b.window[b.index] = true
	}
	b.index = b.nextIndex()
}

func (b *CircuitBreaker) failOnClosed() {
	oldOk := b.window[b.index]
	if oldOk {
		b.failures++
		b.window[b.index] = false
	}
	b.index = b.nextIndex()

	if b.failureRatio() >= b.options.FailureRate {
		b.trip()
	}
}

func (b *CircuitBreaker) trip() {
	if b.state == Open {
		return
	}
	b.setState(Open)
	b.timer = time.AfterFunc(b.options.OpenDuration, func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		b.setState(HalfOpen)
	})
}

func (b *CircuitBreaker) checkHalfOpen() {
	if b.trialCount >= b.options.SuccessThreshold {
		b.reset()
	}
}

func (b *CircuitBreaker) reset() {
	b.setState(Closed)
	for i := range b.window {
		b.window[i] = true
	}
	b.index = 0
	b.failures = 0
	b.trialCount = 0
}

func (b *CircuitBreaker) setState(state State) {
	b.state = state
	if b.options.ChangeStateCallback != nil {
		b.options.ChangeStateCallback(b.state)
	}
}

func (b *CircuitBreaker) failureRatio() float64 {
	return float64(b.failures) / float64(len(b.window))
}

func (b *CircuitBreaker) nextIndex() int {
	return (b.index + 1) % len(b.window)
}
