package goseq

// Map transforms each element of the sequence using fn,
// returning a new Seq[U] with the results.
//
// This is a lazy operation: no elements are processed until a terminal
// operation is called.
//
// Note: Map is a standalone function because Go does not allow methods
// to introduce new type parameters. The output type U can differ from T.
//
// Example:
//
//	result := goseq.Map(goseq.From([]int{1, 2, 3}), strconv.Itoa)
//	// Seq[string]{"1", "2", "3"}
func Map[T any, U any](s Seq[T], fn func(T) U) Seq[U] {
	return newSeq(func() func() (U, bool) {
		next := s.iterate()

		return func() (U, bool) {
			val, ok := next()
			if !ok {
				var zero U
				return zero, false
			}
			return fn(val), true
		}
	})
}

// FlatMap transforms each element into a Seq[U] and flattens the results
// into a single Seq[U].
//
// This is a lazy operation.
//
// Example:
//
//	result := goseq.FlatMap(
//	    goseq.From([]string{"hello world", "foo bar"}),
//	    func(s string) goseq.Seq[string] {
//	        return goseq.From(strings.Split(s, " "))
//	    },
//	)
//	// Seq[string]{"hello", "world", "foo", "bar"}
func FlatMap[T any, U any](s Seq[T], fn func(T) Seq[U]) Seq[U] {
	return newSeq(func() func() (U, bool) {
		outerNext := s.iterate()

		// innerNext holds the iterator for the current inner sequence.
		// Starts as nil — we fetch the first outer element on the first call.
		var innerNext func() (U, bool)

		return func() (U, bool) {
			for {
				// If we have an active inner iterator, try to pull from it.
				if innerNext != nil {
					val, ok := innerNext()
					if ok {
						return val, true
					}
					// Inner exhausted — move to next outer element.
					innerNext = nil
				}

				// Advance the outer iterator.
				outerVal, ok := outerNext()
				if !ok {
					var zero U
					return zero, false
				}

				// Start iterating the new inner sequence.
				innerNext = fn(outerVal).iterate()
			}
		}
	})
}
