package circbrk

import "time"

type ChangeStateCallback func(state State)
type SuccessCallback func(state State)
type FailCallback func(state State)

type Options struct {
	WindowSize          int
	FailureRate         float64
	OpenDuration        time.Duration
	SuccessThreshold    int
	ChangeStateCallback ChangeStateCallback
	SuccessCallback     SuccessCallback
	FailCallback        FailCallback
}

type SetOption func(o *Options)

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
func WithChangeStateCallback(callback ChangeStateCallback) SetOption {
	return func(o *Options) {
		o.ChangeStateCallback = callback
	}
}

// WithSuccessCallback sets a callback that is triggered when the Success method
// is called.
func WithSuccessCallback(callback SuccessCallback) SetOption {
	return func(o *Options) {
		o.SuccessCallback = callback
	}
}

// WithFailCallback sets a callback that is triggered when the Fail method is
// called.
func WithFailCallback(callback FailCallback) SetOption {
	return func(o *Options) {
		o.FailCallback = callback
	}
}

func Apply(ops []SetOption, o *Options) {
	for i := range ops {
		ops[i](o)
	}
	if o.FailureRate < 0.0 || o.FailureRate > 1.0 {
		panic("Options.FailureRate must be between 0.0 and 1.0")
	}
	if o.WindowSize <= 0 {
		panic("Options.WindowSize must be greater than 0")
	}
	if o.OpenDuration <= 0 {
		panic("Options.OpenDuration must be greater than 0")
	}
	if o.SuccessThreshold <= 0 {
		panic("Options.TrialMax must be greater than 0")
	}
}
