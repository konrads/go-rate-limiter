package limiter_test

import (
	"testing"
	"time"

	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/konrads/go-rate-limiter/pkg/limiter"
	"github.com/konrads/go-rate-limiter/pkg/model"
	"github.com/stretchr/testify/assert"
)

func ts(ts int) time.Time {
	return time.Date(2020, time.Month(1), 1, 1, 1, ts, 0, time.UTC)
}

func TestSingleIp(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: model.NewDuration(10)},
			{Limit: 7, Duration: model.NewDuration(20)},
			{Limit: 9, Duration: model.NewDuration(30)},
		},
		&db,
	)

	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(1)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(2)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(3)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(4)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(5)))
	assert.Equal(t, model.LimitRule{Limit: 5, Duration: model.NewDuration(10)}, *l.GetRejectionRule("1.1.1.1", ts(6)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(15)))
	assert.Equal(t, model.LimitRule{Limit: 7, Duration: model.NewDuration(20)}, *l.GetRejectionRule("1.1.1.1", ts(16)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(27)))
	assert.Equal(t, model.LimitRule{Limit: 9, Duration: model.NewDuration(30)}, *l.GetRejectionRule("1.1.1.1", ts(28)))
}

func TestMultipleIp(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: model.NewDuration(10)},
			{Limit: 7, Duration: model.NewDuration(20)},
			{Limit: 9, Duration: model.NewDuration(30)},
		},
		&db,
	)

	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(1)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(2)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(3)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(4)))
	assert.Nil(t, l.GetRejectionRule("1.1.1.1", ts(5)))
	assert.Nil(t, l.GetRejectionRule("2.2.2.2", ts(6))) // Note: different IP, not causing rejection
}

func TestGC(t *testing.T) {
	var db db.DB = db.NewMemDb()
	l := limiter.NewLimiter(
		[]model.LimitRule{
			{Limit: 5, Duration: model.NewDuration(10)},
			{Limit: 7, Duration: model.NewDuration(20)},
			{Limit: 9, Duration: model.NewDuration(30)},
		},
		&db,
	)
	l.GetRejectionRule("1.1.1.1", ts(1))
	l.GetRejectionRule("1.1.1.1", ts(2))
	l.GetRejectionRule("1.1.1.1", ts(3))
	l.GetRejectionRule("1.1.1.1", ts(4))
	l.GetRejectionRule("1.1.1.1", ts(5))
	l.GetRejectionRule("1.1.1.1", ts(55))
}
