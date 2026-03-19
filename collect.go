package goseq

// ToMap converts a Seq[T] into a map[K]V by applying a key extractor
// and a value extractor to each element.
// If two elements produce the same key, the last one wins.
// This is a terminal operation.
//
// Example:
//
//	type Person struct{ Name string; Age int }
//	byName := goseq.ToMap(people,
//	    func(p Person) string { return p.Name },
//	    func(p Person) int    { return p.Age  },
//	)
func ToMap[T any, K comparable, V any](s Seq[T], keyFn func(T) K, valueFn func(T) V) map[K]V {
	result := make(map[K]V)
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		result[keyFn(val)] = valueFn(val)
	}
	return result
}

// GroupBy groups the elements into a map[K][]T, where each key is produced
// by applying keyFn to each element. Elements with the same key are collected
// into the same slice, preserving their original order within each group.
// This is a terminal operation.
//
// Example:
//
//	groups := goseq.GroupBy(goseq.From([]int{1, 2, 3, 4, 5}),
//	    func(n int) string {
//	        if n%2 == 0 { return "even" }
//	        return "odd"
//	    },
//	)
func GroupBy[T any, K comparable](s Seq[T], keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	next := s.iterate()
	for {
		val, ok := next()
		if !ok {
			break
		}
		key := keyFn(val)
		result[key] = append(result[key], val)
	}
	return result
}

// Pair holds two values of potentially different types.
// It is the element type returned by Zip.
type Pair[A any, B any] struct {
	First  A
	Second B
}

// Zip combines two sequences into a Seq of Pairs, pairing elements by position.
// The resulting sequence has the length of the shorter input.
// This is a lazy operation.
//
// Example:
//
//	names := goseq.From([]string{"Alice", "Bob", "Charlie"})
//	scores := goseq.From([]int{95, 87})
//	result := goseq.Zip(names, scores)
//	// Seq[Pair[string,int]]{{"Alice",95}, {"Bob",87}}
func Zip[A any, B any](a Seq[A], b Seq[B]) Seq[Pair[A, B]] {
	return newSeq(func() func() (Pair[A, B], bool) {
		nextA := a.iterate()
		nextB := b.iterate()

		return func() (Pair[A, B], bool) {
			valA, okA := nextA()
			valB, okB := nextB()
			if !okA || !okB {
				var zero Pair[A, B]
				return zero, false
			}
			return Pair[A, B]{First: valA, Second: valB}, true
		}
	})
}
