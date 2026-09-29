package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyedLimiter", func() {
	var limiter *KeyedLimiter

	perKey := func(limit int) func(string) int {
		return func(string) int { return limit }
	}

	Describe("Allow", func() {
		It("should not let one key's requests use up another key's budget", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(2))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())
			Expect(limiter.Allow("b")).To(BeTrue())
		})

		It("should apply the limit limitFor returns for the key", func() {
			limiter = NewKeyedLimiter(10, time.Minute, func(key string) int {
				if key == "critical" {
					return 1
				}

				return 5
			})

			Expect(limiter.Allow("critical")).To(BeTrue())
			Expect(limiter.Allow("critical")).To(BeFalse())
			Expect(limiter.Allow("other")).To(BeTrue())
			Expect(limiter.Allow("other")).To(BeTrue())
		})

		It("should hold every key to the total limit", func() {
			limiter = NewKeyedLimiter(3, time.Minute, perKey(10))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("b")).To(BeTrue())
			Expect(limiter.Allow("c")).To(BeTrue())
			Expect(limiter.Allow("d")).To(BeFalse())
		})

		It("should not record a request the key denied against the total", func() {
			limiter = NewKeyedLimiter(4, time.Minute, perKey(2))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())
			Expect(limiter.Allow("a")).To(BeFalse())

			Expect(limiter.Allow("b")).To(BeTrue())
			Expect(limiter.Allow("c")).To(BeTrue())
			Expect(limiter.Allow("d")).To(BeFalse())
		})

		It("should not record a request the total denied against the key", func() {
			limiter = NewKeyedLimiter(1, time.Minute, perKey(1))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("b")).To(BeFalse())
			Expect(limiter.keys["b"].requests).To(BeEmpty())
		})

		It("should allow a key again after its window expires", func() {
			limiter = NewKeyedLimiter(10, 100*time.Millisecond, perKey(2))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())

			time.Sleep(150 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())
		})

		It("should free only the requests that left the window", func() {
			limiter = NewKeyedLimiter(10, 200*time.Millisecond, perKey(3))

			Expect(limiter.Allow("a")).To(BeTrue())

			time.Sleep(50 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())

			time.Sleep(160 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())
		})

		It("should handle extremely small windows", func() {
			limiter = NewKeyedLimiter(1, time.Nanosecond, perKey(1))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
		})
	})

	Describe("Pause", func() {
		It("should hold back only the paused key", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(10))

			limiter.Pause("a", 100*time.Millisecond)

			Expect(limiter.Allow("a")).To(BeFalse())
			Expect(limiter.Allow("b")).To(BeTrue())

			time.Sleep(150 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
		})

		It("should ignore a non positive pause", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(10))

			limiter.Pause("a", 0)
			limiter.Pause("a", -time.Minute)

			Expect(limiter.Allow("a")).To(BeTrue())
		})

		It("should not consume the key's or the total budget while paused", func() {
			limiter = NewKeyedLimiter(2, time.Minute, perKey(2))

			limiter.Pause("a", 50*time.Millisecond)
			Expect(limiter.Allow("a")).To(BeFalse())
			Expect(limiter.Allow("a")).To(BeFalse())

			time.Sleep(100 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("a")).To(BeFalse())
		})

		// A pause is one shared deadline, so without a spread every waiter resumes in the same
		// instant and the burst the caller's jittered backoff avoided happens anyway.
		It("should spread the wake-up of requests waiting out a pause", func() {
			const waiters = 10

			limiter = NewKeyedLimiter(waiters, time.Minute, perKey(waiters))
			limiter.Pause("a", 200*time.Millisecond)

			var (
				mu        sync.Mutex
				wakeTimes []time.Time
				wg        sync.WaitGroup
			)

			for range waiters {
				wg.Go(func() {
					defer GinkgoRecover()

					Expect(limiter.Wait(context.Background(), "a")).To(BeNil())

					mu.Lock()
					defer mu.Unlock()

					wakeTimes = append(wakeTimes, time.Now())
				})
			}

			wg.Wait()

			earliest, latest := wakeTimes[0], wakeTimes[0]

			for _, wakeTime := range wakeTimes {
				if wakeTime.Before(earliest) {
					earliest = wakeTime
				}

				if wakeTime.After(latest) {
					latest = wakeTime
				}
			}

			Expect(latest.Sub(earliest)).To(BeNumerically(">", 5*time.Millisecond))
		})

		It("should extend an existing pause but never shorten it", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(10))

			limiter.Pause("a", 200*time.Millisecond)
			limiter.Pause("a", 10*time.Millisecond)

			time.Sleep(100 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeFalse())

			time.Sleep(150 * time.Millisecond)

			Expect(limiter.Allow("a")).To(BeTrue())
		})
	})

	Describe("Wait", func() {
		It("should block a full key while another key goes out immediately", func() {
			limiter = NewKeyedLimiter(10, 100*time.Millisecond, perKey(1))

			Expect(limiter.Allow("a")).To(BeTrue())

			start := time.Now()
			err := limiter.Wait(context.Background(), "b")

			Expect(err).To(BeNil())
			Expect(time.Since(start)).To(BeNumerically("<", 10*time.Millisecond))

			start = time.Now()
			err = limiter.Wait(context.Background(), "a")

			Expect(err).To(BeNil())
			Expect(time.Since(start)).To(BeNumerically(">=", 90*time.Millisecond))
		})

		It("should block on the total limit when the key has room", func() {
			limiter = NewKeyedLimiter(1, 100*time.Millisecond, perKey(10))

			Expect(limiter.Allow("a")).To(BeTrue())

			start := time.Now()
			err := limiter.Wait(context.Background(), "b")

			Expect(err).To(BeNil())
			Expect(time.Since(start)).To(BeNumerically(">=", 90*time.Millisecond))
		})

		It("should block until the key's pause expires", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(10))

			limiter.Pause("a", 100*time.Millisecond)

			start := time.Now()
			err := limiter.Wait(context.Background(), "a")

			Expect(err).To(BeNil())
			Expect(time.Since(start)).To(BeNumerically(">=", 90*time.Millisecond))
		})

		It("should respect context cancellation while the key is paused", func() {
			limiter = NewKeyedLimiter(10, time.Minute, perKey(10))

			limiter.Pause("a", time.Minute)

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			Expect(limiter.Wait(ctx, "a")).To(Equal(context.DeadlineExceeded))
		})

		It("should respect context cancellation while the key is full", func() {
			limiter = NewKeyedLimiter(10, 100*time.Millisecond, perKey(1))

			Expect(limiter.Allow("a")).To(BeTrue())

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			start := time.Now()
			err := limiter.Wait(ctx, "a")
			duration := time.Since(start)

			Expect(err).To(Equal(context.DeadlineExceeded))
			Expect(duration).To(BeNumerically(">=", 45*time.Millisecond))
			Expect(duration).To(BeNumerically("<", 70*time.Millisecond))
		})

		It("should return when the context is canceled during a wait", func() {
			limiter = NewKeyedLimiter(10, 100*time.Millisecond, perKey(1))

			Expect(limiter.Allow("a")).To(BeTrue())

			ctx, cancel := context.WithCancel(context.Background())

			var (
				err error
				wg  sync.WaitGroup
			)

			wg.Go(func() {
				err = limiter.Wait(ctx, "a")
			})

			time.Sleep(25 * time.Millisecond)
			cancel()

			wg.Wait()
			Expect(err).To(Equal(context.Canceled))
		})
	})

	Describe("Idle keys", func() {
		It("should evict keys with no requests in the last window", func() {
			limiter = NewKeyedLimiter(10, 50*time.Millisecond, perKey(10))

			Expect(limiter.Allow("a")).To(BeTrue())
			Expect(limiter.Allow("b")).To(BeTrue())
			Expect(limiter.keys).To(HaveLen(2))

			time.Sleep(100 * time.Millisecond)

			Expect(limiter.Allow("c")).To(BeTrue())
			Expect(limiter.keys).To(HaveLen(1))
			Expect(limiter.keys).To(HaveKey("c"))
		})

		It("should not evict a paused key before its pause expires", func() {
			limiter = NewKeyedLimiter(10, 50*time.Millisecond, perKey(10))

			limiter.Pause("a", 200*time.Millisecond)

			time.Sleep(100 * time.Millisecond)

			Expect(limiter.Allow("b")).To(BeTrue())
			Expect(limiter.keys).To(HaveKey("a"))
			Expect(limiter.Allow("a")).To(BeFalse())
		})
	})

	Describe("Concurrent Access", func() {
		It("should be safe for concurrent Allow calls across keys", func() {
			limiter = NewKeyedLimiter(15, time.Minute, perKey(5))

			keys := []string{"a", "b", "c", "d"}
			counts := make([]int32, len(keys))

			var wg sync.WaitGroup

			for i := range keys {
				for range 20 {
					wg.Go(func() {
						if limiter.Allow(keys[i]) {
							atomic.AddInt32(&counts[i], 1)
						}
					})
				}
			}

			wg.Wait()

			total := int32(0)

			for i := range keys {
				Expect(counts[i]).To(BeNumerically("<=", 5))

				total += counts[i]
			}

			Expect(total).To(Equal(int32(15)))
		})

		It("should hold concurrent Wait calls to the key and total limits", func() {
			limiter = NewKeyedLimiter(4, time.Minute, perKey(2))

			keys := []string{"a", "b", "c"}
			counts := make([]int32, len(keys))

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			var wg sync.WaitGroup

			for i := range 9 {
				wg.Go(func() {
					if limiter.Wait(ctx, keys[i%len(keys)]) == nil {
						atomic.AddInt32(&counts[i%len(keys)], 1)
					}
				})
			}

			wg.Wait()

			total := int32(0)

			for i := range keys {
				Expect(counts[i]).To(BeNumerically("<=", 2))

				total += counts[i]
			}

			Expect(total).To(Equal(int32(4)))
		})

		It("should let every concurrent waiter through once the window frees up", func() {
			limiter = NewKeyedLimiter(4, 100*time.Millisecond, perKey(2))

			var wg sync.WaitGroup

			errors := make([]error, 8)

			for i := range errors {
				wg.Go(func() {
					errors[i] = limiter.Wait(context.Background(), []string{"a", "b", "c"}[i%3])
				})
			}

			wg.Wait()

			for i := range errors {
				Expect(errors[i]).To(BeNil())
			}
		})

		It("should handle very high request rates", func() {
			limiter = NewKeyedLimiter(100, time.Second, perKey(1000))

			var wg sync.WaitGroup

			successCount := int32(0)

			for i := range 1000 {
				wg.Go(func() {
					if limiter.Allow([]string{"a", "b", "c"}[i%3]) {
						atomic.AddInt32(&successCount, 1)
					}
				})
			}

			wg.Wait()

			Expect(successCount).To(Equal(int32(100)))
		})
	})
})
