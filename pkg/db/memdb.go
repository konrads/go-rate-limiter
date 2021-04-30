package db

import (
	"log"

	"github.com/konrads/go-rate-limiter/pkg/utils"
)

type MemDb struct {
	byIP map[string][]int64
}

func NewMemDb() *MemDb {
	return &MemDb{byIP: map[string][]int64{}}
}

func (db *MemDb) AddHit(ip string, hit int64) {
	if _, ok := db.byIP[ip]; !ok {
		db.byIP[ip] = []int64{}
	}
	hits := db.byIP[ip]
	hits = append(hits, hit)
	log.Printf("DB.AddHit: IP: %s, hit %d", ip, hit)
	db.byIP[ip] = hits
}

func (db *MemDb) GetHits(ip string, minHit int64) *[]int64 {
	if hits, ok := db.byIP[ip]; ok {
		res := utils.DropWhile(hits, func(x int64) bool { return x < minHit })
		log.Printf("DB.GetHits: IP: %s, minHit %d, res: %v", ip, minHit, res)
		return &res
	}
	res := []int64{}
	log.Printf("DB.GetHits: IP: %s, minHit %d, res: <empty>", ip, minHit)
	return &res
}

func (db *MemDb) Cleanup(ip string, minHit int64) {
	if hits, ok := db.byIP[ip]; ok {
		truncated := utils.DropWhile(hits, func(x int64) bool { return x < minHit })
		db.byIP[ip] = truncated
	}
}

// noop for memory
func (db *MemDb) Close() error {
	return nil
}
