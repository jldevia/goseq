package goseq

import "slices"

func (s Seq[T]) Filter(predicate func(T) bool) Seq[T] {
	var filtered []T
	for i := range s.items {
		if predicate(s.items[i]) {
			filtered = append(filtered, s.items[i])
		}
	}
	return Seq[T]{items: slices.Clip(filtered)}
}
