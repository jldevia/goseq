package goseq

// Take returns a new Seq[T] with at most the first n elements.
// If n is greater than the length of the sequence, all elements are returned.
// If n is zero or negative, an empty sequence is returned.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice() // []int{1, 2, 3}
func (s Seq[T]) Take(n int) Seq[T] {
	if n <= 0 {
		return Empty[T]()
	}
	if n >= len(s.items) {
		return From(s.items)
	}
	return From(s.items[:n])
}

// Skip returns a new Seq[T] skipping the first n elements.
// If n is greater than the length of the sequence, an empty sequence is returned.
// If n is zero or negative, the full sequence is returned.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice() // []int{3, 4, 5}
func (s Seq[T]) Skip(n int) Seq[T] {
	if n <= 0 {
		return From(s.items)
	}
	if n >= len(s.items) {
		return Empty[T]()
	}
	return From(s.items[n:])
}

// TakeWhile returns a new Seq[T] with elements from the start of the sequence
// as long as the predicate is true. It stops at the first element that does
// not satisfy the predicate.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).TakeWhile(func(n int) bool { return n < 4 }).ToSlice()
//	// []int{1, 2, 3}
func (s Seq[T]) TakeWhile(predicate func(T) bool) Seq[T] {
	result := make([]T, 0)
	for _, item := range s.items {
		if !predicate(item) {
			break
		}
		result = append(result, item)
	}
	return Seq[T]{items: result}
}

// SkipWhile returns a new Seq[T] skipping elements from the start of the
// sequence as long as the predicate is true. Once an element fails the
// predicate, the rest are included — including subsequent elements that
// would match the predicate.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).SkipWhile(func(n int) bool { return n < 3 }).ToSlice()
//	// []int{3, 4, 5}
func (s Seq[T]) SkipWhile(predicate func(T) bool) Seq[T] {
	i := 0
	for i < len(s.items) && predicate(s.items[i]) {
		i++
	}
	return From(s.items[i:])
}

// Distinct returns a new Seq[T] with duplicate elements removed,
// preserving the order of first appearance.
// The type T must be comparable.
//
// Example:
//
//	goseq.From([]int{1, 2, 2, 3, 1, 4}).Distinct().ToSlice() // []int{1, 2, 3, 4}
func Distinct[T comparable](s Seq[T]) Seq[T] {
	seen := make(map[T]struct{}, len(s.items))
	result := make([]T, 0)
	for _, item := range s.items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return Seq[T]{items: result}
}
