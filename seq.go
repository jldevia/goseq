// Package goseq provides a fluent, type-safe API for processing collections,
// inspired by Java Streams and C# LINQ.
package goseq

// Seq represents a sequence of elements of type T.
// It wraps a slice and provides a chainable API for transforming,
// filtering, and aggregating data.
//
// All operations are eager: they execute immediately and return a new Seq[T].
type Seq[T any] struct {
	items []T
}

// From creates a new Seq[T] from an existing slice.
// The original slice is not modified; a copy is made internally.
//
// Example:
//
//	s := goseq.From([]int{1, 2, 3, 4, 5})
func From[T any](items []T) Seq[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return Seq[T]{items: copied}
}

// Empty creates a new empty Seq[T].
//
// Example:
//
//	s := goseq.Empty[int]()
func Empty[T any]() Seq[T] {
	return Seq[T]{items: []T{}}
}

// Len returns the number of elements in the sequence.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Len() // 3
func (s Seq[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the sequence contains no elements.
//
// Example:
//
//	goseq.Empty[int]().IsEmpty() // true
func (s Seq[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// ToSlice returns the elements of the sequence as a plain Go slice.
// This is the primary way to "collect" results after chaining operations.
//
// Example:
//
//	result := goseq.From([]int{1, 2, 3}).ToSlice() // []int{1, 2, 3}
func (s Seq[T]) ToSlice() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// ForEach executes the given function for each element in the sequence.
// It is a terminal operation used for side effects (e.g. printing, logging).
//
// Example:
//
//	goseq.From([]string{"a", "b"}).ForEach(func(s string) {
//	    fmt.Println(s)
//	})
func (s Seq[T]) ForEach(fn func(T)) {
	for _, item := range s.items {
		fn(item)
	}
}
