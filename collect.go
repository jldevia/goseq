package goseq

// ToMap converts a Seq[T] into a map[K]V by applying a key extractor
// and a value extractor to each element.
//
// The key type K must be comparable (a Go requirement for map keys).
// If two elements produce the same key, the last one wins.
//
// Note: ToMap is a standalone function because it introduces new type
// parameters K and V that cannot be expressed as method type parameters in Go.
//
// Example — index persons by name:
//
//	type Person struct{ Name string; Age int }
//	people := goseq.From([]Person{{"Alice", 30}, {"Bob", 25}})
//	byName := goseq.ToMap(people,
//	    func(p Person) string { return p.Name },
//	    func(p Person) int    { return p.Age  },
//	)
//	// map[string]int{"Alice": 30, "Bob": 25}
//
// Example — square numbers indexed by themselves:
//
//	squares := goseq.ToMap(goseq.From([]int{1, 2, 3}),
//	    func(n int) int { return n },
//	    func(n int) int { return n * n },
//	)
//	// map[int]int{1: 1, 2: 4, 3: 9}
func ToMap[T any, K comparable, V any](s Seq[T], keyFn func(T) K, valueFn func(T) V) map[K]V {
	result := make(map[K]V, len(s.items))
	for _, item := range s.items {
		result[keyFn(item)] = valueFn(item)
	}
	return result
}

// GroupBy groups the elements of a Seq[T] into a map[K][]T, where each key
// is produced by applying the key extractor function to each element.
// Elements that produce the same key are collected into the same slice,
// preserving their original order within each group.
//
// The key type K must be comparable (a Go requirement for map keys).
//
// Note: GroupBy is a standalone function because it introduces a new type
// parameter K that cannot be expressed as a method type parameter in Go.
//
// Example — group numbers by parity:
//
//	groups := goseq.GroupBy(goseq.From([]int{1, 2, 3, 4, 5}),
//	    func(n int) string {
//	        if n%2 == 0 { return "even" }
//	        return "odd"
//	    },
//	)
//	// map[string][]int{"odd": {1, 3, 5}, "even": {2, 4}}
//
// Example — group persons by age range:
//
//	type Person struct{ Name string; Age int }
//	people := goseq.From([]Person{{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}})
//	groups := goseq.GroupBy(people, func(p Person) string {
//	    if p.Age >= 18 { return "adult" }
//	    return "minor"
//	})
//	// map[string][]Person{"adult": {Alice, Charlie}, "minor": {Bob}}
func GroupBy[T any, K comparable](s Seq[T], keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range s.items {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}
	return result
}
