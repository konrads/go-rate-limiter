package model

type LimitRule struct {
	Limit    int64 `json:"Limit" binding:"required"`
	Duration int64 `json:"Duration" binding:"required"`
}

type LimitRules []LimitRule

// for sortability
func (a LimitRules) Len() int           { return len(a) }
func (a LimitRules) Less(i, j int) bool { return a[i].Duration < a[j].Duration }
func (a LimitRules) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type LimitRecord struct {
	Rule *LimitRule
	Hits []int64
}

type LimitRecords []*LimitRecord
