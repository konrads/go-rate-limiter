package limiter

import (
	"log"
	"sort"
	"time"

	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/konrads/go-rate-limiter/pkg/model"
)

// Deprecated: attempt at mimiter that should be able to use pluggable DB.
// Not memory efficient...
type Limiter struct {
	templateRules model.LimitRules
	db            *db.DB
	maxDuration   *time.Duration
}

func NewLimiter(rules model.LimitRules, db *db.DB) Limiter {
	sort.Sort(rules)
	maxDuration := rules[len(rules)-1].Duration.Duration
	return Limiter{
		templateRules: rules,
		db:            db,
		maxDuration:   &maxDuration,
	}
}

// add a hit, return false if breached a rule
func (l *Limiter) GetRejectionRule(ip string, now time.Time) *model.LimitRule {
	(*l.db).AddHit(ip, now.Unix())
	var res *model.LimitRule = nil
	for _, rule := range l.templateRules {
		hits4Rule := (*l.db).GetHits(ip, now.Add(-rule.Duration.Duration).Unix())
		if len(*hits4Rule) > int(rule.Limit) {
			res = &rule
			log.Printf("...got rejection: %v", rule)
			break
		}
	}

	(*l.db).Cleanup(ip, now.Add(-*l.maxDuration).Unix())

	return res
}
