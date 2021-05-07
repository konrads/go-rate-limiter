package leakybucket

import (
	"time"

	"github.com/konrads/go-rate-limiter/pkg/model"
	"github.com/konrads/go-rate-limiter/pkg/utils"
)

// Leaky bucket implementation, inspired by:
// https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm/#:~:text=Leaky%20bucket%20(closely%20related%20to,a%20bucket%20holding%20the%20requests.&text=This%20algorithm's%20advantage%20is%20that,at%20an%20approximately%20average%20rate.
//
// FIXME: utilizes slices which leak memory on the underlying memory
type SimpleLeakyBucket struct {
	limitRules []model.LimitRule
	byIP       map[string]map[model.LimitRule][]time.Time
}

func NewSimpleLeakyBucket(limitRules []model.LimitRule) *SimpleLeakyBucket {
	return &SimpleLeakyBucket{
		limitRules: limitRules,
		byIP:       make(map[string]map[model.LimitRule][]time.Time),
	}
}

func (lb *SimpleLeakyBucket) GetRejectionRule(ip string, now time.Time) *model.LimitRule {
	var firstRejection *model.LimitRule = nil
	// add ip if needed
	buckets, ok := lb.byIP[ip]
	if !ok {
		buckets = make(map[model.LimitRule][]time.Time, len(lb.limitRules))
		lb.byIP[ip] = buckets
	}
	for _, rule := range lb.limitRules {
		// add bucket if needed
		bucket, ok := buckets[rule]
		if !ok {
			bucket = make([]time.Time, 0, rule.Limit)
		}
		// append ts otherwise
		bucket = append(bucket, now)
		buckets[rule] = bucket

		minTs := now.Add(-rule.Duration.Duration)
		utils.DropWhileMut(&bucket, func(x time.Time) bool { return x.Before(minTs) })
		if len(bucket) > int(rule.Limit) && firstRejection == nil {
			rejectionRule := rule
			firstRejection = &rejectionRule
		}
	}
	return firstRejection
}

// scan all of the leaky buckets, prune as per timestamp,
func (lb *SimpleLeakyBucket) Cleanup(now time.Time) {
	var allBucketCnt int
	for ip, buckets := range lb.byIP {
		allBucketCnt = 0
		for rule, bucket := range buckets {
			minTs := now.Add(-rule.Duration.Duration)
			utils.DropWhileMut(&bucket, func(x time.Time) bool { return x.Before(minTs) })
			buckets[rule] = bucket
			allBucketCnt += len(bucket)
			if len(bucket) == 0 {
				delete(buckets, rule)
			}
		}

		if allBucketCnt == 0 {
			delete(lb.byIP, ip)
		}
	}
}

// counts recorded
func (lb *SimpleLeakyBucket) Stats() int {
	size := 0
	for _, buckets := range lb.byIP {
		for _, bucket := range buckets {
			size += len(bucket)
		}
	}
	return size
}
