package utils_test

import (
	"testing"

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
