package circbrk

import (
	"math/rand"
	"testing"
	"time"

	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("Should work", func(t *testing.T) {
		cb := New(WithWindowSize(8),
			WithFailureRate(0.5),
			WithOpenDuration(200*time.Millisecond),
			WithSuccessThreshold(4),
		)
		assertfatal.Equal(cb.State(), Closed, t)

		cb.Fail()
		cb.Fail()
		cb.Fail()
		cb.Fail()
		cb.Fail()
		assertfatal.Equal(cb.State(), Open, t)

		time.Sleep(250 * time.Millisecond)

		assertfatal.Equal(cb.State(), HalfOpen, t)
		cb.Success()
		cb.Success()
		cb.Fail()
		assertfatal.Equal(cb.State(), Open, t)

		time.Sleep(250 * time.Millisecond)

		assertfatal.Equal(cb.State(), HalfOpen, t)

		cb.Success()
		cb.Success()
		cb.Success()
		cb.Success()
		assertfatal.Equal(cb.State(), Closed, t)

		cb.Success()
		cb.Success()
		cb.Success()
		cb.Success()
		cb.Fail()
		cb.Fail()
		cb.Success()
		cb.Success()
		assertfatal.Equal(cb.State(), Closed, t)
	})

	t.Run("HalfOpen state consistently allows SuccessThreshold calls before closing",
		func(t *testing.T) {
			cb := New(
				WithWindowSize(5),
				WithFailureRate(0.5),
				WithOpenDuration(50*time.Millisecond),
				WithSuccessThreshold(3),
			)

			// Force breaker to Open state by failing enough calls
			for range 3 {
				cb.Fail()
			}
			assertfatal.Equal(cb.State(), Open, t)

			// Wait for OpenDuration to elapse so it transitions to HalfOpen
			time.Sleep(60 * time.Millisecond)
			assertfatal.Equal(cb.State(), HalfOpen, t)

			// First HalfOpen period: count allowed attempts
			unblockedCount := 0
			for range 10 {
				if cb.Allow() {
					unblockedCount++
				}
			}
			assertfatal.Equal(3, unblockedCount, t)
			for range 3 {
				cb.Success()
			}
			assertfatal.Equal(cb.State(), Closed, t)

			// Repeat: cause Open again
			for range 3 {
				cb.Fail()
			}
			assertfatal.Equal(cb.State(), Open, t)

			time.Sleep(60 * time.Millisecond)
			assertfatal.Equal(cb.State(), HalfOpen, t)

			// Second HalfOpen period: count allowed attempts
			unblockedCount = 0
			for range 10 {
				if cb.Allow() {
					unblockedCount++
				}
			}
			assertfatal.Equal(3, unblockedCount, t)
			for range 3 {
				cb.Success()
			}
			assertfatal.Equal(Closed, cb.State(), t)
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
