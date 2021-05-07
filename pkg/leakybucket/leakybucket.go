package leakybucket

import (
	"time"

	"github.com/konrads/go-rate-limiter/pkg/model"
)

type LeakyBucket interface {
	GetRejectionRule(ip string, now time.Time) *model.LimitRule
	Cleanup(now time.Time)
	Stats() int
}
