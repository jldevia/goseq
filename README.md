# goseq

> A fluent, type-safe collection processing library for Go — inspired by Java Streams and C# LINQ.

## What is goseq?

Go's standard library is intentionally minimal and doesn't include higher-order collection utilities like `map`, `filter`, `reduce`, or `groupBy`. Anyone coming from Java, C#, Kotlin, or Rust quickly notices this gap.

**goseq** fills that gap by providing a clean, chainable API for processing sequences of data. Built on top of Go generics (introduced in Go 1.18), it brings the expressiveness of Java Streams and LINQ to Go — without sacrificing type safety.

```go
result := goseq.From(numbers).
    Filter(func(n int) bool { return n%2 == 0 }).
    Map(func(n int) int { return n * 2 }).
    ToSlice()
```

## Design Principles

- **Type-safe**: built entirely with Go generics — no `interface{}`, no casting
- **Fluent API**: chainable operations for readable, expressive pipelines
- **Eager evaluation**: operations execute immediately, keeping the implementation simple and predictable
- **Zero dependencies**: only the Go standard library
- **Test-first**: every feature ships with comprehensive tests

## Requirements

- Go 1.18 or later (generics support required)

## Installation

```bash
go get github.com/your-username/goseq
```

## Quick Example

```go
import "github.com/your-username/goseq"

names := []string{"Alice", "Bob", "Charlie", "Ana", "Brian"}

result := goseq.From(names).
    Filter(func(s string) bool { return strings.HasPrefix(s, "A") }).
    Map(strings.ToUpper).
    ToSlice()

// result: ["ALICE", "ANA"]
```

---

## Roadmap

### v0.1 — Core Foundation

The goal of this version is to establish the central `Seq[T]` type and the most essential operations, along with a solid testing baseline.

- `From(slice)` constructor
- `Map` — transform each element
- `Filter` — keep elements matching a predicate
- `Reduce` — aggregate elements into a single value
- `ForEach` — iterate with side effects
- `ToSlice` — collect results back to a slice
- `ToMap` — collect into a `map[K]V`
- `GroupBy` — group elements by a key function

### v0.2 — Query Operations

Expand the API with common query and pagination operations.

- `First`, `Last` — get first or last element (with optional predicate)
- `Any`, `All`, `None` — boolean checks over the sequence
- `Count` — count elements (optionally matching a predicate)
- `Take`, `Skip` — slice the sequence by count
- `TakeWhile`, `SkipWhile` — slice by predicate
- `Contains` — check membership
- `Distinct` — remove duplicates

### v0.3 — Advanced Operations

More powerful transformations and aggregations.

- `FlatMap` — transform and flatten nested collections
- `Zip` — combine two sequences element by element
- `OrderBy`, `OrderByDescending` — sorting with key selector
- `Sum`, `Min`, `Max` — numeric aggregations with type constraints

### v1.0 — Lazy Evaluation

Refactor the internal architecture to support lazy (deferred) evaluation, where operations are only executed when the terminal result is consumed. The public API remains unchanged.

- Pipeline execution deferred until `ToSlice`, `ToMap`, etc.
- Efficient `Take(n)` — stops processing after n elements
- Support for potentially infinite sequences
- Significant memory improvements for large datasets

---

## Project Status

Currently in active early development. The API may change between minor versions until v1.0.

## Contributing

Contributions, issues, and feature requests are welcome. Feel free to open an issue to discuss what you'd like to add.

## License

MIT