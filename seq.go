// Package goseq provides a fluent, type-safe API for processing collections,
// inspired by Java Streams and C# LINQ.
//
// As of v1.0, evaluation is lazy: operations are not executed until a terminal
// operation (ToSlice, ToMap, ForEach, etc.) is called. This means that
// pipelines like Take(5) over a million elements only process what they need.
package goseq

// Seq represents a lazy sequence of elements of type T.
//
// Internally it holds a factory function (iterate) that creates a fresh
// iterator each time a terminal operation is called. This allows the same
// Seq[T] to be consumed multiple times safely.
//
// An iterator is a function that returns the next element and true on each
// call, or the zero value of T and false when the sequence is exhausted.
type Seq[T any] struct {
	// iterate is a factory: each call returns a brand new iterator.
	// This is what makes sequences reusable.
	iterate func() func() (T, bool)
}

// newSeq is the internal constructor. It wraps an iterator factory into a Seq.
func newSeq[T any](iterate func() func() (T, bool)) Seq[T] {
	return Seq[T]{iterate: iterate}
}

// From creates a new Seq[T] from an existing slice.
// The original slice is not modified.
//
// Example:
//
//	s := goseq.From([]int{1, 2, 3, 4, 5})
func From[T any](items []T) Seq[T] {
	// Copy once at construction time to protect against external mutation.
	copied := make([]T, len(items))
	copy(copied, items)

	return newSeq(func() func() (T, bool) {
		// Each call to iterate() creates a fresh index for this iterator.
		// This is the key to reusability: every terminal operation gets
		// its own independent cursor into the slice.
		i := 0
		return func() (T, bool) {
			if i >= len(copied) {
				var zero T
				return zero, false
			}
			val := copied[i]
			i++
			return val, true
		}
	})
}

// Empty creates a new empty Seq[T] that yields no elements.
//
// Example:
//
//	s := goseq.Empty[int]()
func Empty[T any]() Seq[T] {
	return newSeq(func() func() (T, bool) {
		return func() (T, bool) {
			var zero T
			return zero, false
		}
	})
}

// ToSlice collects all elements of the sequence into a plain Go slice.
// This is a terminal operation: it triggers the execution of the pipeline.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Filter(...).Map(...).ToSlice()
func (s Seq[T]) ToSlice() []T {
	result := make([]T, 0)
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		result = append(result, val)
	}
	return result
}

// ForEach executes fn for each element in the sequence.
// This is a terminal operation.
//
// Example:
//
//	goseq.From([]string{"a", "b"}).ForEach(func(s string) {
//	    fmt.Println(s)
//	})
func (s Seq[T]) ForEach(fn func(T)) {
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		fn(val)
	}
}

// Len returns the number of elements in the sequence.
// This is a terminal operation: it consumes the sequence to count elements.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Len() // 3
func (s Seq[T]) Len() int {
	count := 0
	next := s.iterate()
	for {
		_, ok := next()
		if !ok {
			break
		}
		count++
	}
	return count
}

// IsEmpty returns true if the sequence contains no elements.
// This is a terminal operation.
//
// Example:
//
//	goseq.Empty[int]().IsEmpty() // true
func (s Seq[T]) IsEmpty() bool {
	next := s.iterate()
	_, ok := next()
	return !ok
}
