package circbrk

import (
	"errors"
	"time"
)

const (
	DefaultWindowSize       = 50
	DefaultFailureRate      = 0.5
	DefaultOpenDuration     = 5 * time.Second
	DefaultSuccessThreshold = 3
)

type (
	ChangeStateCallback func(state State)
	SuccessCallback     func(state State)
	FailCallback        func(state State)
)

type Options struct {
	WindowSize          int
	FailureRate         float64
	OpenDuration        time.Duration
	SuccessThreshold    int
	ChangeStateCallback ChangeStateCallback
	SuccessCallback     SuccessCallback
	FailCallback        FailCallback
}

func (o Options) Validate() error {
	if o.WindowSize <= 0 {
		return errors.New("Options.WindowSize must be greater than 0")
	}
	if o.FailureRate <= 0.0 || o.FailureRate > 1.0 {
		return errors.New("Options.FailureRate must be in (0.0, 1.0]")
	}
	if o.OpenDuration <= 0 {
		return errors.New("Options.OpenDuration must be greater than 0")
	}
	if o.SuccessThreshold <= 0 {
		return errors.New("Options.SuccessThreshold must be greater than 0")
	}
	return nil
}

type SetOption func(o *Options)

func Apply(o *Options, ops ...SetOption) {
	for i := range ops {
		ops[i](o)
	}
}

// WithWindowSize sets the number of recent calls to track when calculating the
// failure rate.
func WithWindowSize(windowSize int) SetOption {
	return func(o *Options) {
		o.WindowSize = windowSize
	}
}

// WithFailureRate sets the failure rate threshold (from 0.0 to 1.0) at which
// the circuit opens. For example, 0.5 means the circuit opens if at least 50%
// of recent calls failed.
func WithFailureRate(failureRate float64) SetOption {
	return func(o *Options) {
		o.FailureRate = failureRate
	}
}

// WithOpenDuration sets the time the circuit remains open before transitioning
// to half-open.
func WithOpenDuration(duration time.Duration) SetOption {
	return func(o *Options) {
		o.OpenDuration = duration
	}
}

// WithSuccessThreshold sets the number of successful calls required in the
// half-open state to transition back to the closed state.
func WithSuccessThreshold(threshold int) SetOption {
	return func(o *Options) {
		o.SuccessThreshold = threshold
	}
}

// WithChangeStateCallback sets a callback function that is invoked whenever the
// circuit changes its state (e.g., Closed -> Open, Open -> HalfOpen).
// The callback runs while the breaker mutex is held. To prevent deadlocks,
// avoid calling CircuitBreaker methods inside it.
func WithChangeStateCallback(callback ChangeStateCallback) SetOption {
	return func(o *Options) {
		o.ChangeStateCallback = callback
	}
}

// WithSuccessCallback sets a callback that is triggered when the Success method
// is called.
// The callback runs while the breaker mutex is held. To prevent deadlocks,
// avoid calling CircuitBreaker methods inside it.
func WithSuccessCallback(callback SuccessCallback) SetOption {
	return func(o *Options) {
		o.SuccessCallback = callback
	}
}

// WithFailCallback sets a callback that is triggered when the Fail method is
// called.
// The callback runs while the breaker mutex is held. To prevent deadlocks,
// avoid calling CircuitBreaker methods inside it.
func WithFailCallback(callback FailCallback) SetOption {
	return func(o *Options) {
		o.FailCallback = callback
	}
}
