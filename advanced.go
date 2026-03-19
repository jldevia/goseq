package goseq

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func OrderBy[T any, K constraints.Ordered](s Seq[T], keyFn func(T) K) Seq[T] {
	items := s.ToSlice()
	sort.SliceStable(items, func(i, j int) bool {
		return keyFn(items[i]) < keyFn(items[j])
	})
	return From(items)
}

func OrderByDescending[T any, K constraints.Ordered](s Seq[T], keyFn func(T) K) Seq[T] {
	items := s.ToSlice()
	sort.SliceStable(items, func(i, j int) bool {
		return keyFn(items[i]) > keyFn(items[j])
	})
	return From(items)
}

func Sum[T constraints.Ordered | constraints.Complex](s Seq[T]) T {
	var total T
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		total += val
	}
	return total
}

func Min[T constraints.Ordered](s Seq[T]) (T, bool) {
	next := s.iterate()
	first, ok := next()
	if !ok {
		var zero T
		return zero, false
	}
	min := first
	for {
		val, ok := next()
		if !ok {
			break
		}
		if val < min {
			min = val
		}
	}
	return min, true
}

func Max[T constraints.Ordered](s Seq[T]) (T, bool) {
	next := s.iterate()
	first, ok := next()
	if !ok {
		var zero T
		return zero, false
	}
	max := first
	for {
		val, ok := next()
		if !ok {
			break
		}
		if val > max {
			max = val
		}
	}
	return max, true
}
