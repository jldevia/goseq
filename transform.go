package goseq

// Map transforms each element of the sequence using the provided function,
// returning a new Seq[U] with the results.
//
// Note: Map is a standalone function rather than a method because Go does not
// allow methods to introduce new type parameters. This means the output type U
// can be completely different from the input type T.
//
// Example — transform ints to strings:
//
//	result := goseq.Map(goseq.From([]int{1, 2, 3}), strconv.Itoa)
//	// Seq[string]{"1", "2", "3"}
//
// Example — extract a field from a struct:
//
//	type Person struct{ Name string; Age int }
//	people := goseq.From([]Person{{"Alice", 30}, {"Bob", 25}})
//	names := goseq.Map(people, func(p Person) string { return p.Name })
//	// Seq[string]{"Alice", "Bob"}
func Map[T any, U any](s Seq[T], fn func(T) U) Seq[U] {
	result := make([]U, len(s.items))
	for i, item := range s.items {
		result[i] = fn(item)
	}
	return Seq[U]{items: result}
}
