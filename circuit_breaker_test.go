package circbrk

import (
	"math/rand"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {

	t.Run("Should work", func(t *testing.T) {
		cb := New(WithWindowSize(8),
			WithFailureRate(0.5),
			WithOpenDuration(200*time.Millisecond),
			WithSuccessThreshold(4),
		)
		if cb.State() != Closed {
			t.Fatalf("expected Closed, got %v", cb.State())
		}
		cb.Fail()
		cb.Fail()
		cb.Fail()
		cb.Fail()
		cb.Fail()
		if cb.State() != Open {
			t.Fatalf("expected Open, got %v", cb.State())
		}

		time.Sleep(250 * time.Millisecond)

		if cb.State() != HalfOpen {
			t.Fatalf("expected HalfOpen, got %v", cb.State())
		}
		cb.Success()
		cb.Success()
		cb.Fail()
		if cb.State() != Open {
			t.Fatalf("expected Open, got %v", cb.State())
		}

		time.Sleep(250 * time.Millisecond)

		if cb.State() != HalfOpen {
			t.Fatalf("expected HalfOpen, got %v", cb.State())
		}

		cb.Success()
		cb.Success()
		cb.Success()
		cb.Success()
		if cb.State() != Closed {
			t.Fatalf("expected Closed, got %v", cb.State())
		}
		cb.Success()
		cb.Success()
		cb.Success()
		cb.Success()
		cb.Fail()
		cb.Fail()
		cb.Success()
		cb.Success()
		if cb.State() != Closed {
			t.Fatalf("expected Closed, got %v", cb.State())
		}
	})

	t.Run("Concurrent", func(t *testing.T) {
		cb := New(WithWindowSize(8),
			WithFailureRate(0.5),
			WithOpenDuration(200*time.Millisecond),
			WithSuccessThreshold(4),
			// WithChangeStateCallback(func(state State) {
			// 	fmt.Printf("new state %v\n", state)
			// }),
			// WithSuccessCallback(
			// 	func(state State) {
			// 		fmt.Printf("s %v\n", state)
			// 	},
			// ),
			// WithFailCallback(
			// 	func(state State) {
			// 		fmt.Printf("f %v\n", state)
			// 	},
			// ),
		)

		for range 4 {
			go func() {
				for range time.NewTicker(time.Millisecond * 10).C {
					if rand.Intn(3) == 0 {
						cb.Fail()
					} else {
						cb.Success()
					}
				}
			}()
		}

		time.Sleep(5 * time.Second)
	})
}
