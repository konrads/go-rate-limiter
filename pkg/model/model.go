package model

import (
	"encoding/json"
	"errors"
	"time"
)

// using duration marshalling from https://stackoverflow.com/questions/48050945/how-to-unmarshal-json-into-durations
type Duration struct {
	time.Duration
}

func NewDuration(secs uint) Duration {
	return Duration{
		time.Duration(secs) * time.Second,
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type LimitRule struct {
	Limit    int64    `json:"Limit" binding:"required"`
	Duration Duration `json:"Duration" binding:"required"`
}

type LimitRules []LimitRule

// for sortability
func (a LimitRules) Len() int           { return len(a) }
func (a LimitRules) Less(i, j int) bool { return a[i].Duration.Duration < a[j].Duration.Duration }
func (a LimitRules) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type LimitRecord struct {
	Rule *LimitRule
	Hits []int64
}

type LimitRecords []*LimitRecord
