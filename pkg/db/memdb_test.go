package db_test

import (
	"testing"

	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestMemDb(t *testing.T) {
	memdb := db.NewMemDb()
	memdb.AddHit("1.1.1.1", 0)
	memdb.AddHit("1.1.1.1", 1)
	memdb.AddHit("1.1.1.1", 2)
	memdb.AddHit("1.1.1.1", 3)
	memdb.AddHit("1.1.1.1", 4)
	assert.Equal(t, &[]int64{0, 1, 2, 3, 4}, memdb.GetHits("1.1.1.1", 0))
	assert.Equal(t, &[]int64{2, 3, 4}, memdb.GetHits("1.1.1.1", 2))
	memdb.Cleanup("1.1.1.1", 3)
	assert.Equal(t, &[]int64{3, 4}, memdb.GetHits("1.1.1.1", 0))
	assert.Equal(t, &[]int64{}, memdb.GetHits("2.2.2.2", 0))
}
