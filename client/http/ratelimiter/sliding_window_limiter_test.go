package ratelimiter

// import (
// 	"context"
// 	"sync"
// 	"sync/atomic"
// 	"testing"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// func TestSlidingWindowLimiter(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "SlidingWindowLimiter Suite")
// }

// var _ = Describe("SlidingWindowLimiter", func() {
// 	var limiter *SlidingWindowLimiter

// 	Describe("NewSlidingWindowLimiter", func() {
// 		It("should create a new limiter with correct configuration", func() {
// 			maxRequests := 10
// 			window := time.Minute

// 			limiter = NewSlidingWindowLimiter(maxRequests, window)

// 			Expect(limiter).ToNot(BeNil())
// 			Expect(limiter.maxRequests).To(Equal(maxRequests))
// 			Expect(limiter.window).To(Equal(window))
// 			Expect(limiter.requests).To(HaveLen(0))
// 			Expect(limiter.requests).To(HaveCap(maxRequests))
// 		})
// 	})

// 	Describe("Allow", func() {
// 		It("should allow requests up to the limit", func() {
// 			limiter = NewSlidingWindowLimiter(3, time.Second)
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeFalse())
// 		})

// 		It("should track request timestamps", func() {
// 			limiter = NewSlidingWindowLimiter(3, time.Second)

// 			start := time.Now()

// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())

// 			Expect(len(limiter.requests)).To(Equal(2))
// 			Expect(limiter.requests[0]).To(BeTemporally(">=", start))
// 			Expect(limiter.requests[1]).To(BeTemporally(">=", limiter.requests[0]))
// 		})

// 		It("should allow new requests after window expires", func() {
// 			limiter = NewSlidingWindowLimiter(2, 100*time.Millisecond)

// 			// Fill the limit
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeFalse())

// 			// Wait for window to expire
// 			time.Sleep(150 * time.Millisecond)

// 			// Should allow new requests
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeFalse())
// 		})

// 		It("should handle partial window expiration", func() {
// 			limiter = NewSlidingWindowLimiter(3, 200*time.Millisecond)

// 			// Make first request
// 			Expect(limiter.Allow()).To(BeTrue())

// 			// Wait a bit
// 			time.Sleep(50 * time.Millisecond)

// 			// Make second and third requests
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeFalse())

// 			// Wait for first request to expire (but not the second and third)
// 			time.Sleep(160 * time.Millisecond)

// 			// Should allow one more request (first expired, second and third still active)
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeFalse())
// 		})
// 	})

// 	Describe("Wait", func() {
// 		BeforeEach(func() {
// 			limiter = NewSlidingWindowLimiter(2, 100*time.Millisecond)
// 		})

// 		It("should return immediately when under limit", func() {
// 			ctx := context.Background()
// 			start := time.Now()

// 			err := limiter.Wait(ctx)

// 			Expect(err).To(BeNil())
// 			Expect(time.Since(start)).To(BeNumerically("<", 10*time.Millisecond))
// 		})

// 		It("should wait when at limit", func() {
// 			ctx := context.Background()

// 			// Fill the limit
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())

// 			start := time.Now()
// 			err := limiter.Wait(ctx)
// 			duration := time.Since(start)

// 			Expect(err).To(BeNil())
// 			Expect(duration).To(BeNumerically(">=", 90*time.Millisecond))
// 			Expect(duration).To(BeNumerically("<", 200*time.Millisecond))
// 		})

// 		It("should respect context cancellation", func() {
// 			// Fill the limit
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())

// 			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
// 			defer cancel()

// 			start := time.Now()
// 			err := limiter.Wait(ctx)
// 			duration := time.Since(start)

// 			Expect(err).To(Equal(context.DeadlineExceeded))
// 			Expect(duration).To(BeNumerically(">=", 45*time.Millisecond))
// 			Expect(duration).To(BeNumerically("<", 70*time.Millisecond))
// 		})

// 		It("should handle context cancellation during wait", func() {
// 			// Fill the limit
// 			Expect(limiter.Allow()).To(BeTrue())
// 			Expect(limiter.Allow()).To(BeTrue())

// 			ctx, cancel := context.WithCancel(context.Background())

// 			var err error
// 			var wg sync.WaitGroup
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				err = limiter.Wait(ctx)
// 			}()

// 			// Cancel after a short delay
// 			time.Sleep(25 * time.Millisecond)
// 			cancel()

// 			wg.Wait()
// 			Expect(err).To(Equal(context.Canceled))
// 		})

// 		It("should handle multiple concurrent waiters", func() {
// 			limiter = NewSlidingWindowLimiter(1, 100*time.Millisecond)

// 			// Fill the limit
// 			Expect(limiter.Allow()).To(BeTrue())

// 			ctx := context.Background()
// 			var wg sync.WaitGroup
// 			results := make([]error, 3)

// 			for i := 0; i < 3; i++ {
// 				wg.Add(1)
// 				go func(idx int) {
// 					defer wg.Done()
// 					results[idx] = limiter.Wait(ctx)
// 				}(i)
// 			}

// 			wg.Wait()

// 			// All should succeed eventually
// 			for i := 0; i < 3; i++ {
// 				Expect(results[i]).To(BeNil())
// 			}
// 		})
// 	})

// 	Describe("Concurrent Access", func() {
// 		BeforeEach(func() {
// 			limiter = NewSlidingWindowLimiter(10, time.Second)
// 		})

// 		It("should be safe for concurrent Allow calls", func() {
// 			var wg sync.WaitGroup
// 			successCount := int32(0)

// 			for i := 0; i < 20; i++ {
// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					if limiter.Allow() {
// 						atomic.AddInt32(&successCount, 1)
// 					}
// 				}()
// 			}

// 			wg.Wait()

// 			// Should have exactly 10 successful calls
// 			Expect(successCount).To(Equal(int32(10)))
// 		})

// 		It("should be safe for concurrent Wait calls", func() {
// 			limiter = NewSlidingWindowLimiter(3, 100*time.Millisecond)

// 			var wg sync.WaitGroup
// 			ctx := context.Background()
// 			errors := make([]error, 5)

// 			for i := 0; i < 5; i++ {
// 				wg.Add(1)
// 				go func(idx int) {
// 					defer wg.Done()
// 					errors[idx] = limiter.Wait(ctx)
// 				}(i)
// 			}

// 			wg.Wait()

// 			// All should succeed eventually
// 			for i := 0; i < 5; i++ {
// 				Expect(errors[i]).To(BeNil())
// 			}
// 		})

// 		It("should handle mixed Allow and Wait calls", func() {
// 			limiter = NewSlidingWindowLimiter(5, 200*time.Millisecond)

// 			var wg sync.WaitGroup
// 			ctx := context.Background()

// 			// Start some Allow calls
// 			allowResults := make([]bool, 10)
// 			for i := 0; i < 10; i++ {
// 				wg.Add(1)
// 				go func(idx int) {
// 					defer wg.Done()
// 					allowResults[idx] = limiter.Allow()
// 				}(i)
// 			}

// 			// Start some Wait calls
// 			waitResults := make([]error, 3)
// 			for i := 0; i < 3; i++ {
// 				wg.Add(1)
// 				go func(idx int) {
// 					defer wg.Done()
// 					waitResults[idx] = limiter.Wait(ctx)
// 				}(i)
// 			}

// 			wg.Wait()

// 			// Count successful Allow calls
// 			allowSuccess := 0
// 			for _, result := range allowResults {
// 				if result {
// 					allowSuccess++
// 				}
// 			}

// 			// All Wait calls should succeed
// 			for i := 0; i < 3; i++ {
// 				Expect(waitResults[i]).To(BeNil())
// 			}

// 			// Total successful requests should not exceed limit + Wait calls
// 			Expect(allowSuccess).To(BeNumerically("<=", 5))
// 		})
// 	})

// 	Describe("Edge Cases", func() {
// 		It("should handle very high request rates", func() {
// 			limiter = NewSlidingWindowLimiter(100, time.Second)

// 			var wg sync.WaitGroup
// 			successCount := int32(0)

// 			for i := 0; i < 1000; i++ {
// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					if limiter.Allow() {
// 						atomic.AddInt32(&successCount, 1)
// 					}
// 				}()
// 			}

// 			wg.Wait()

// 			Expect(successCount).To(Equal(int32(100)))
// 		})

// 		It("should handle extremely small windows", func() {
// 			limiter = NewSlidingWindowLimiter(1, time.Nanosecond)

// 			// First request should succeed
// 			Expect(limiter.Allow()).To(BeTrue())

// 			// Should be able to make another request immediately
// 			// as the window is so small it expires instantly
// 			Expect(limiter.Allow()).To(BeTrue())
// 		})
// 	})
// })
