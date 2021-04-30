package db

type DB interface {
	AddHit(ip string, hit int64)
	GetHits(ip string, minHit int64) *[]int64
	Cleanup(ip string, minHit int64)
	Close() error
}
