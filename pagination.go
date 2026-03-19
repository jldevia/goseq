package goseq

// Take returns a new Seq[T] yielding at most the first n elements.
// If n <= 0, an empty sequence is returned.
//
// With lazy evaluation, Take(5) over a million elements only processes
// the first 5 — the rest are never touched.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice() // []int{1, 2, 3}
func (s Seq[T]) Take(n int) Seq[T] {
	return newSeq(func() func() (T, bool) {
		next := s.iterate()
		remaining := n

		return func() (T, bool) {
			if remaining <= 0 {
				var zero T
				return zero, false
			}
			val, ok := next()
			if !ok {
				var zero T
				return zero, false
			}
			remaining--
			return val, true
		}
	})
}

// Skip returns a new Seq[T] that skips the first n elements.
// If n <= 0, the full sequence is returned.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice() // []int{3, 4, 5}
func (s Seq[T]) Skip(n int) Seq[T] {
	return newSeq(func() func() (T, bool) {
		next := s.iterate()
		skipped := false

		return func() (T, bool) {
			// Skip the first n elements exactly once per iterator lifetime.
			if !skipped {
				skipped = true
				for i := 0; i < n; i++ {
					_, ok := next()
					if !ok {
						var zero T
						return zero, false
					}
				}
			}
			return next()
		}
	})
}

// TakeWhile returns a new Seq[T] yielding elements from the start of the
// sequence as long as the predicate is true. Stops at the first non-match.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).
//	    TakeWhile(func(n int) bool { return n < 4 }).
//	    ToSlice() // []int{1, 2, 3}
func (s Seq[T]) TakeWhile(predicate func(T) bool) Seq[T] {
	return newSeq(func() func() (T, bool) {
		next := s.iterate()
		done := false

		return func() (T, bool) {
			if done {
				var zero T
				return zero, false
			}
			val, ok := next()
			if !ok || !predicate(val) {
				done = true
				var zero T
				return zero, false
			}
			return val, true
		}
	})
}

// SkipWhile returns a new Seq[T] that skips elements from the start as long
// as the predicate is true, then yields the rest — including elements that
// would have matched the predicate later.
//
// Example:
//
//	goseq.From([]int{1, 2, 3, 4, 5}).
//	    SkipWhile(func(n int) bool { return n < 3 }).
//	    ToSlice() // []int{3, 4, 5}
func (s Seq[T]) SkipWhile(predicate func(T) bool) Seq[T] {
	return newSeq(func() func() (T, bool) {
		next := s.iterate()
		skipping := true

		return func() (T, bool) {
			for {
				val, ok := next()
				if !ok {
					var zero T
					return zero, false
				}
				if skipping && predicate(val) {
					// Still in the skipping phase — discard this element.
					continue
				}
				// First non-match ends the skipping phase permanently.
				skipping = false
				return val, true
			}
		}
	})
}

// Distinct returns a new Seq[T] with duplicate elements removed,
// preserving the order of first appearance.
// The type T must be comparable.
//
// Example:
//
//	goseq.Distinct(goseq.From([]int{1, 2, 2, 3, 1})).ToSlice() // []int{1, 2, 3}
func Distinct[T comparable](s Seq[T]) Seq[T] {
	return newSeq(func() func() (T, bool) {
		next := s.iterate()
		seen := make(map[T]struct{})

		return func() (T, bool) {
			for {
				val, ok := next()
				if !ok {
					var zero T
					return zero, false
				}
				if _, exists := seen[val]; !exists {
					seen[val] = struct{}{}
					return val, true
				}
			}
		}
	})
}
