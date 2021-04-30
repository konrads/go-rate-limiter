package limiter

import (
	"log"
	"sort"

	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/konrads/go-rate-limiter/pkg/model"
)

type Limiter struct {
	templateRules model.LimitRules
	db            *db.DB
	maxDuration   int64
}

func NewLimiter(rules model.LimitRules, db *db.DB) Limiter {
	sort.Sort(rules)
	maxDuration := rules[len(rules)-1].Duration
	return Limiter{
		templateRules: rules,
		db:            db,
		maxDuration:   maxDuration,
	}
}

// add a hit, return false if breached a rule
func (l *Limiter) GetRejectionRule(ip string, now int64) *model.LimitRule {
	(*l.db).AddHit(ip, now)
	var res *model.LimitRule = nil
	for _, rule := range l.templateRules {
		hits4Rule := (*l.db).GetHits(ip, now-rule.Duration)
		if len(*hits4Rule) > int(rule.Limit) {
			res = &rule
			log.Printf("...got rejection: %v", rule)
			break
		}
	}

	(*l.db).Cleanup(ip, now-l.maxDuration)

	return res
}
