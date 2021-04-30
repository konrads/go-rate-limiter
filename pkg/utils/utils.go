package utils

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
