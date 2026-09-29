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

		It("should be safe for concurrent Wait calls across keys", func() {
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
	})
})
