package limiter_test

import (
	"testing"

	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/konrads/go-rate-limiter/pkg/limiter"
	"github.com/konrads/go-rate-limiter/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSingleIp(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: 10},
			{Limit: 7, Duration: 20},
			{Limit: 9, Duration: 30},
		},
		&db,
	)

	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 1))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 2))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 3))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 4))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 5))
	assert.Equal(t, model.LimitRule{Limit: 5, Duration: 10}, *l.GetRejectionRule("1.1.1.1", 6))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 15))
	assert.Equal(t, model.LimitRule{Limit: 7, Duration: 20}, *l.GetRejectionRule("1.1.1.1", 16))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 27))
	assert.Equal(t, model.LimitRule{Limit: 9, Duration: 30}, *l.GetRejectionRule("1.1.1.1", 28))
}

func TestMultipleIp(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: 10},
			{Limit: 7, Duration: 20},
			{Limit: 9, Duration: 30},
		},
		&db,
	)

	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 1))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 2))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 3))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 4))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", 5))
	assert.Nil(t, l.GetRejectionRule("2.2.2.2", 6)) // Note: different IP, not causing rejection
}

func TestGC(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: 10},
			{Limit: 7, Duration: 20},
			{Limit: 9, Duration: 30},
		},
		&db,
	)
	l.GetRejectionRule("1.1.1.1", 1)
	l.GetRejectionRule("1.1.1.1", 2)
	l.GetRejectionRule("1.1.1.1", 3)
	l.GetRejectionRule("1.1.1.1", 4)
	l.GetRejectionRule("1.1.1.1", 5)
	l.GetRejectionRule("1.1.1.1", 55)
}
