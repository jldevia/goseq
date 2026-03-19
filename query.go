package goseq

// First returns the first element of the sequence and true.
// Returns the zero value and false if the sequence is empty.
// This is a terminal operation.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3}).First()
//	// val=1, ok=true
func (s Seq[T]) First() (T, bool) {
	next := s.iterate()
	return next()
}

// FirstWhere returns the first element satisfying the predicate and true.
// Returns the zero value and false if no element matches.
// This is a terminal operation.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3, 4}).
//	    FirstWhere(func(n int) bool { return n%2 == 0 })
//	// val=2, ok=true
func (s Seq[T]) FirstWhere(predicate func(T) bool) (T, bool) {
	return s.Filter(predicate).First()
}

// Last returns the last element of the sequence and true.
// Returns the zero value and false if the sequence is empty.
// This is a terminal operation.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3}).Last()
//	// val=3, ok=true
func (s Seq[T]) Last() (T, bool) {
	var last T
	found := false
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		last = val
		found = true
	}
	return last, found
}

// LastWhere returns the last element satisfying the predicate and true.
// Returns the zero value and false if no element matches.
// This is a terminal operation.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3, 4}).
//	    LastWhere(func(n int) bool { return n%2 == 0 })
//	// val=4, ok=true
func (s Seq[T]) LastWhere(predicate func(T) bool) (T, bool) {
	return s.Filter(predicate).Last()
}

// Any returns true if at least one element satisfies the predicate.
// Returns false for empty sequences.
// This is a terminal operation that short-circuits on the first match.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Any(func(n int) bool { return n > 2 }) // true
func (s Seq[T]) Any(predicate func(T) bool) bool {
	_, ok := s.FirstWhere(predicate)
	return ok
}

// All returns true if every element satisfies the predicate.
// Returns true for empty sequences (vacuous truth).
// This is a terminal operation that short-circuits on the first non-match.
//
// Example:
//
//	goseq.From([]int{2, 4, 6}).All(func(n int) bool { return n%2 == 0 }) // true
func (s Seq[T]) All(predicate func(T) bool) bool {
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		if !predicate(val) {
			return false
		}
	}
	return true
}

// None returns true if no element satisfies the predicate.
// Returns true for empty sequences.
// This is a terminal operation.
//
// Example:
//
//	goseq.From([]int{1, 3, 5}).None(func(n int) bool { return n%2 == 0 }) // true
func (s Seq[T]) None(predicate func(T) bool) bool {
	return !s.Any(predicate)
}

// Count returns the number of elements that satisfy the predicate.
// This is a terminal operation.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4}).Count(func(n int) bool { return n%2 == 0 }) // 2
func (s Seq[T]) Count(predicate func(T) bool) int {
	return s.Filter(predicate).Len()
}

// Contains returns true if the sequence contains the given value.
// The type T must be comparable.
// This is a terminal operation that short-circuits on the first match.
//
// Example:
//
//	goseq.Contains(goseq.From([]int{1, 2, 3}), 2) // true
func Contains[T comparable](s Seq[T], value T) bool {
	return s.Any(func(item T) bool { return item == value })
}
