package circbrk

// State represents a CircuitBreaker state.
type State int

func (s State) String() string {
	switch s {
	case Closed:
		return "Closed"
	case HalfOpen:
		return "HalfOpen"
	case Open:
		return "Open"
	default:
		return "Unknown"
	}
}

const (
	Closed State = iota
	HalfOpen
	Open
)
