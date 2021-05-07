package leakybucket

import (
	"sync"
	"time"

	"github.com/konrads/go-rate-limiter/pkg/model"
)

// Thread-safe LeakyBucket implementation, wraps LeakBucket
type SafeLeakyBucket struct {
	wrapped *LeakyBucket
	mutex   sync.RWMutex
}

func NewSafeLeakyBucket(limitRules []model.LimitRule) *SafeLeakyBucket {
	return &SafeLeakyBucket{
		wrapped: NewLeakyBucket(limitRules),
		// mutex does not need initializing
	}
}

func (slb *SafeLeakyBucket) GetRejectionRule(ip string, now time.Time) *model.LimitRule {
	slb.mutex.Lock()
	defer slb.mutex.Unlock()
	return slb.wrapped.GetRejectionRule(ip, now)
}

// scan all of the leaky buckets, prune as per timestamp,
func (slb *SafeLeakyBucket) Cleanup(now time.Time) {
	slb.mutex.Lock()
	defer slb.mutex.Unlock()
	slb.wrapped.Cleanup(now)
}

// counts recorded
func (lb *SafeLeakyBucket) Stats() int {
	// no write, won't mutex lock
	return lb.wrapped.Stats()
}
