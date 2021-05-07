package utils

import (
	"time"
)

func DropWhile(xs []int64, predicate func(x int64) bool) []int64 {
	firstKeptInd := 0
	keepSome := false
	for i, x := range xs {
		firstKeptInd = i
		if !predicate(x) {
			keepSome = true
			break
		}
	}
	if keepSome {
		return xs[firstKeptInd:]
	} else {
		return []int64{}
	}
}

func DropWhileMut(xs *[]time.Time, predicate func(x time.Time) bool) {
	firstKeptInd := 0
	keepSome := false
	for i, x := range *xs {
		firstKeptInd = i
		if !predicate(x) {
			keepSome = true
			break
		}
	}
	if keepSome {
		*xs = (*xs)[firstKeptInd:]
	} else {
		*xs = []time.Time{}
	}
}
