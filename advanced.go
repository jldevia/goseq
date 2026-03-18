package goseq

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// OrderBy returns a new Seq[T] sorted in ascending order by the value
// returned by the key selector function.
//
// Example — sort persons by age:
//
//	type Person struct{ Name string; Age int }
//	sorted := goseq.From(people).
//	    OrderBy(func(p Person) int { return p.Age }).
//	    ToSlice()
//
// Example — sort strings by length:
//
//	sorted := goseq.From([]string{"banana", "go", "rust"}).
//	    OrderBy(func(s string) int { return len(s) }).
//	    ToSlice()
//	// ["go", "rust", "banana"]
func OrderBy[T any, K constraints.Ordered](s Seq[T], keyFn func(T) K) Seq[T] {
	result := make([]T, len(s.items))
	copy(result, s.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFn(result[i]) < keyFn(result[j])
	})
	return Seq[T]{items: result}
}

// OrderByDescending returns a new Seq[T] sorted in descending order by the
// value returned by the key selector function.
//
// Example — sort persons by age descending:
//
//	sorted := goseq.OrderByDescending(people, func(p Person) int { return p.Age })
func OrderByDescending[T any, K constraints.Ordered](s Seq[T], keyFn func(T) K) Seq[T] {
	result := make([]T, len(s.items))
	copy(result, s.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFn(result[i]) > keyFn(result[j])
	})
	return Seq[T]{items: result}
}

// Sum returns the sum of all elements in the sequence.
// Returns the zero value of T for empty sequences.
//
// Example:
//
//	goseq.Sum(goseq.From([]int{1, 2, 3, 4, 5})) // 15
func Sum[T constraints.Ordered | constraints.Complex](s Seq[T]) T {
	var total T
	for _, item := range s.items {
		total += item
	}
	return total
}

// Min returns the minimum element in the sequence and true.
// Returns the zero value of T and false if the sequence is empty.
//
// Example:
//
//	min, ok := goseq.Min(goseq.From([]int{3, 1, 4, 1, 5}))
//	// min=1, ok=true
func Min[T constraints.Ordered](s Seq[T]) (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	min := s.items[0]
	for _, item := range s.items[1:] {
		if item < min {
			min = item
		}
	}
	return min, true
}

// Max returns the maximum element in the sequence and true.
// Returns the zero value of T and false if the sequence is empty.
//
// Example:
//
//	max, ok := goseq.Max(goseq.From([]int{3, 1, 4, 1, 5}))
//	// max=5, ok=true
func Max[T constraints.Ordered](s Seq[T]) (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	max := s.items[0]
	for _, item := range s.items[1:] {
		if item > max {
			max = item
		}
	}
	return max, true
}
