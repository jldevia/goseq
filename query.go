package goseq

// First returns the first element of the sequence and true.
// If the sequence is empty, it returns the zero value of T and false.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3}).First()
//	// val=1, ok=true
//
//	val, ok := goseq.Empty[int]().First()
//	// val=0, ok=false
func (s Seq[T]) First() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[0], true
}

// FirstWhere returns the first element that satisfies the predicate and true.
// If no element matches, it returns the zero value of T and false.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3, 4}).FirstWhere(func(n int) bool { return n%2 == 0 })
//	// val=2, ok=true
func (s Seq[T]) FirstWhere(predicate func(T) bool) (T, bool) {
	for _, item := range s.items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Last returns the last element of the sequence and true.
// If the sequence is empty, it returns the zero value of T and false.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3}).Last()
//	// val=3, ok=true
func (s Seq[T]) Last() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// LastWhere returns the last element that satisfies the predicate and true.
// If no element matches, it returns the zero value of T and false.
//
// Example:
//
//	val, ok := goseq.From([]int{1, 2, 3, 4}).LastWhere(func(n int) bool { return n%2 == 0 })
//	// val=4, ok=true
func (s Seq[T]) LastWhere(predicate func(T) bool) (T, bool) {
	for i := len(s.items) - 1; i >= 0; i-- {
		if predicate(s.items[i]) {
			return s.items[i], true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if at least one element satisfies the predicate.
// Returns false for empty sequences.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Any(func(n int) bool { return n > 2 }) // true
func (s Seq[T]) Any(predicate func(T) bool) bool {
	for _, item := range s.items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All returns true if every element satisfies the predicate.
// Returns true for empty sequences (vacuous truth).
//
// Example:
//
//	goseq.From([]int{2, 4, 6}).All(func(n int) bool { return n%2 == 0 }) // true
func (s Seq[T]) All(predicate func(T) bool) bool {
	for _, item := range s.items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// None returns true if no element satisfies the predicate.
// Returns true for empty sequences.
//
// Example:
//
//	goseq.From([]int{1, 3, 5}).None(func(n int) bool { return n%2 == 0 }) // true
func (s Seq[T]) None(predicate func(T) bool) bool {
	return !s.Any(predicate)
}

// Count returns the number of elements that satisfy the predicate.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).Count(func(n int) bool { return n%2 == 0 }) // 2
func (s Seq[T]) Count(predicate func(T) bool) int {
	count := 0
	for _, item := range s.items {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Contains returns true if the sequence contains the given value.
// The type T must be comparable.
//
// Example:
//
//	goseq.From([]int{1, 2, 3}).Contains(2) // true
//	goseq.From([]int{1, 2, 3}).Contains(9) // false
func Contains[T comparable](s Seq[T], value T) bool {
	for _, item := range s.items {
		if item == value {
			return true
		}
	}
	return false
}
