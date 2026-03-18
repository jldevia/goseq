package goseq

// FlatMap transforms each element of the sequence into a Seq[U] using the
// provided function, then flattens the results into a single Seq[U].
//
// Note: FlatMap is a standalone function because it introduces a new type
// parameter U that cannot be expressed as a method type parameter in Go.
//
// Example — flatten a slice of slices:
//
//	words := goseq.FlatMap(
//	    goseq.From([]string{"hello world", "foo bar"}),
//	    func(s string) goseq.Seq[string] {
//	        return goseq.From(strings.Split(s, " "))
//	    },
//	)
//	// Seq[string]{"hello", "world", "foo", "bar"}
//
// Example — expand each number into a range:
//
//	result := goseq.FlatMap(
//	    goseq.From([]int{1, 2, 3}),
//	    func(n int) goseq.Seq[int] {
//	        items := make([]int, n)
//	        for i := range items { items[i] = n }
//	        return goseq.From(items)
//	    },
//	)
//	// Seq[int]{1, 2, 2, 3, 3, 3}
func FlatMap[T any, U any](s Seq[T], fn func(T) Seq[U]) Seq[U] {
	result := make([]U, 0)
	for _, item := range s.items {
		inner := fn(item)
		result = append(result, inner.items...)
	}
	return Seq[U]{items: result}
}

// Zip combines two sequences into a single Seq of pairs, pairing elements
// by position. The resulting sequence has the length of the shorter input.
// Surplus elements from the longer sequence are ignored.
//
// Note: Zip is a standalone function because it introduces new type
// parameters A and B that cannot be expressed as method type parameters in Go.
//
// Example:
//
//	names := goseq.From([]string{"Alice", "Bob", "Charlie"})
//	scores := goseq.From([]int{95, 87})
//
//	result := goseq.Zip(names, scores)
//	// Seq[Pair[string, int]]{{"Alice", 95}, {"Bob", 87}}
//	// "Charlie" is dropped — shorter sequence wins
func Zip[A any, B any](a Seq[A], b Seq[B]) Seq[Pair[A, B]] {
	minLen := len(a.items)
	if len(b.items) < minLen {
		minLen = len(b.items)
	}

	result := make([]Pair[A, B], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = Pair[A, B]{First: a.items[i], Second: b.items[i]}
	}
	return Seq[Pair[A, B]]{items: result}
}

// Pair holds two values of potentially different types.
// It is the element type returned by Zip.
type Pair[A any, B any] struct {
	First  A
	Second B
}
