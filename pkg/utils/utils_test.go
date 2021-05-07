package utils_test

import (
	"testing"
	"time"

	"github.com/konrads/go-rate-limiter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestDropWhile(t *testing.T) {
	xs := []int64{2, 3, 4, 10, 11, 12}

	assert.Equal(t, []int64{10, 11, 12}, utils.DropWhile(xs, func(x int64) bool { return x < 9 }))
	// ensure no mutation
	assert.Equal(t, xs, []int64{2, 3, 4, 10, 11, 12})
}

func TestDropWhile_takeAll(t *testing.T) {
	xs := []int64{2, 3, 4, 10, 11, 12}
	assert.Equal(t, []int64{2, 3, 4, 10, 11, 12}, utils.DropWhile(xs, func(x int64) bool { return x < 0 }))
}

func TestDropWhile_takeNone(t *testing.T) {
	xs := []int64{2, 3, 4, 10, 11, 12}
	assert.Equal(t, []int64{}, utils.DropWhile(xs, func(x int64) bool { return x > 0 }))
}

func ts(ts int) time.Time {
	return time.Date(2020, time.Month(1), 1, 1, 1, ts, 0, time.UTC)
}

func TestDropWhileMut(t *testing.T) {
	xs := []time.Time{ts(2), ts(3), ts(4), ts(10), ts(11), ts(12)}
	utils.DropWhileMut(&xs, func(x time.Time) bool { return x.Before(ts(9)) })
	assert.Equal(t, []time.Time{ts(10), ts(11), ts(12)}, xs)
}

func TestDropWhileMut_takeAll(t *testing.T) {
	xs := []time.Time{ts(2), ts(3), ts(4), ts(10), ts(11), ts(12)}
	utils.DropWhileMut(&xs, func(x time.Time) bool { return x.Before(ts(0)) })
	assert.Equal(t, []time.Time{ts(2), ts(3), ts(4), ts(10), ts(11), ts(12)}, xs)
}

func TestDropWhileMut_takeNone(t *testing.T) {
	xs := []time.Time{ts(2), ts(3), ts(4), ts(10), ts(11), ts(12)}
	utils.DropWhileMut(&xs, func(x time.Time) bool { return x.After(ts(0)) })
	assert.Equal(t, []time.Time{}, xs)
}
