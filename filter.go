package goseq

// Filter returns a new Seq[T] that yields only the elements for which
// the predicate returns true. The order of elements is preserved.
//
// This is a lazy operation: no elements are processed until a terminal
// operation (ToSlice, ForEach, etc.) is called.
//
// Example:
//
//	evens := goseq.From([]int{1, 2, 3, 4, 5}).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    ToSlice() // []int{2, 4}
func (s Seq[T]) Filter(predicate func(T) bool) Seq[T] {
	return newSeq(func() func() (T, bool) {
		// Capture a fresh iterator from the upstream sequence.
		next := s.iterate()

		return func() (T, bool) {
			// Keep pulling from upstream until we find a match
			// or the upstream is exhausted.
			for {
				val, ok := next()
				if !ok {
					var zero T
					return zero, false
				}
				if predicate(val) {
					return val, true
				}
			}
		}
	})
}
