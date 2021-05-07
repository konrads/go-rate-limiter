package leakybucket

import (
	"sync"
	"time"

	"github.com/konrads/go-rate-limiter/pkg/model"
)

// Thread-safe LeakyBucket implementation, wraps LeakBucket
type SafeLeakyBucket struct {
	wrapped *SimpleLeakyBucket
	mutex   sync.RWMutex
}

func NewSafeLeakyBucket(limitRules []model.LimitRule) *SafeLeakyBucket {
	return &SafeLeakyBucket{
		wrapped: NewSimpleLeakyBucket(limitRules),
		// mutex does not need initializing
	}
}

func (w *SafeLeakyBucket) GetRejectionRule(ip string, now time.Time) *model.LimitRule {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.wrapped.GetRejectionRule(ip, now)
}

// scan all of the leaky buckets, prune as per timestamp,
func (w *SafeLeakyBucket) Cleanup(now time.Time) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.wrapped.Cleanup(now)
}

// counts recorded
func (w *SafeLeakyBucket) Stats() int {
	// no write, won't mutex lock
	return w.wrapped.Stats()
}
