package goseq

// Reduce aggregates all elements of the sequence into a single value of type U,
// by applying the accumulator function to each element starting from an initial value.
//
// Note: Reduce is a standalone function rather than a method because Go does not
// allow methods to introduce new type parameters. The output type U can be
// completely different from the input type T.
//
// The accumulator function receives the current accumulated value and the current
// element, and returns the new accumulated value.
//
// If the sequence is empty, the initial value is returned as-is.
//
// Example — sum of integers:
//
//	sum := goseq.Reduce(goseq.From([]int{1, 2, 3, 4, 5}), 0,
//	    func(acc, n int) int { return acc + n },
//	) // 15
//
// Example — concatenate strings:
//
//	result := goseq.Reduce(goseq.From([]string{"a", "b", "c"}), "",
//	    func(acc, s string) string { return acc + s },
//	) // "abc"
//
// Example — collect even numbers into a slice (T → U of different type):
//
//	evens := goseq.Reduce(goseq.From([]int{1, 2, 3, 4}), []int{},
//	    func(acc []int, n int) []int {
//	        if n%2 == 0 {
//	            return append(acc, n)
//	        }
//	        return acc
//	    },
//	) // []int{2, 4}
func Reduce[T any, U any](s Seq[T], initial U, fn func(U, T) U) U {
	acc := initial
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		acc = fn(acc, val)
	}
	return acc
}
