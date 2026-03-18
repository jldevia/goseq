package goseq_test

import (
	"strings"
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// FlatMap
// =====================

// Caso clásico: expandir cada número en una secuencia repetida.
func TestFlatMap_ExpandsEachElement(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]int{1, 2, 3}),
		func(n int) goseq.Seq[int] {
			items := make([]int, n)
			for i := range items {
				items[i] = n
			}
			return goseq.From(items)
		},
	).ToSlice()

	assertSliceEqual(t, []int{1, 2, 2, 3, 3, 3}, result)
}

// Dividir strings en palabras y aplanar.
func TestFlatMap_SplitsStringsIntoWords(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]string{"hello world", "foo bar"}),
		func(s string) goseq.Seq[string] {
			return goseq.From(strings.Split(s, " "))
		},
	).ToSlice()

	assertSliceEqual(t, []string{"hello", "world", "foo", "bar"}, result)
}

// Cambio de tipo: cada int genera sus divisores como strings.
func TestFlatMap_ChangesType(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]int{1, 2, 3}),
		func(n int) goseq.Seq[string] {
			return goseq.From([]string{strings.Repeat("x", n)})
		},
	).ToSlice()

	assertSliceEqual(t, []string{"x", "xx", "xxx"}, result)
}

// Si la función devuelve secuencias vacías, el resultado es vacío.
func TestFlatMap_InnerSeqEmpty_ReturnsEmpty(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]int{1, 2, 3}),
		func(n int) goseq.Seq[int] {
			return goseq.Empty[int]()
		},
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Secuencia vacía produce resultado vacío.
func TestFlatMap_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.FlatMap(
		goseq.Empty[int](),
		func(n int) goseq.Seq[int] {
			return goseq.From([]int{n, n * 2})
		},
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Preserva el orden: primero todos los del primer elemento, luego del segundo.
func TestFlatMap_PreservesOrder(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]int{10, 20}),
		func(n int) goseq.Seq[int] {
			return goseq.From([]int{n, n + 1, n + 2})
		},
	).ToSlice()

	assertSliceEqual(t, []int{10, 11, 12, 20, 21, 22}, result)
}

// No modifica la secuencia original.
func TestFlatMap_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3})
	_ = goseq.FlatMap(original, func(n int) goseq.Seq[int] {
		return goseq.From([]int{n * 10})
	})

	assertSliceEqual(t, []int{1, 2, 3}, original.ToSlice())
}

// Encadenamiento: FlatMap seguido de Filter.
func TestFlatMap_ChainedWithFilter(t *testing.T) {
	result := goseq.FlatMap(
		goseq.From([]int{1, 2, 3}),
		func(n int) goseq.Seq[int] {
			return goseq.From([]int{n, n * 10})
		},
	).Filter(func(n int) bool { return n > 5 }).ToSlice()

	assertSliceEqual(t, []int{10, 20, 30}, result)
}

// =====================
// Zip
// =====================

// Caso básico: dos secuencias del mismo largo.
func TestZip_SameLength_PairsAllElements(t *testing.T) {
	result := goseq.Zip(
		goseq.From([]string{"Alice", "Bob", "Charlie"}),
		goseq.From([]int{95, 87, 72}),
	).ToSlice()

	if len(result) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(result))
	}
	if result[0].First != "Alice" || result[0].Second != 95 {
		t.Errorf("unexpected pair[0]: %v", result[0])
	}
	if result[1].First != "Bob" || result[1].Second != 87 {
		t.Errorf("unexpected pair[1]: %v", result[1])
	}
	if result[2].First != "Charlie" || result[2].Second != 72 {
		t.Errorf("unexpected pair[2]: %v", result[2])
	}
}

// Cuando la primera es más corta, se usa su largo.
func TestZip_FirstShorter_TruncatesToFirst(t *testing.T) {
	result := goseq.Zip(
		goseq.From([]int{1, 2}),
		goseq.From([]int{10, 20, 30, 40}),
	).ToSlice()

	if len(result) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(result))
	}
}

// Cuando la segunda es más corta, se usa su largo.
func TestZip_SecondShorter_TruncatesToSecond(t *testing.T) {
	result := goseq.Zip(
		goseq.From([]int{1, 2, 3, 4}),
		goseq.From([]int{10, 20}),
	).ToSlice()

	if len(result) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(result))
	}
	if result[0].First != 1 || result[0].Second != 10 {
		t.Errorf("unexpected pair[0]: %v", result[0])
	}
	if result[1].First != 2 || result[1].Second != 20 {
		t.Errorf("unexpected pair[1]: %v", result[1])
	}
}

// Primera vacía produce resultado vacío.
func TestZip_FirstEmpty_ReturnsEmpty(t *testing.T) {
	result := goseq.Zip(
		goseq.Empty[int](),
		goseq.From([]int{1, 2, 3}),
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Segunda vacía produce resultado vacío.
func TestZip_SecondEmpty_ReturnsEmpty(t *testing.T) {
	result := goseq.Zip(
		goseq.From([]int{1, 2, 3}),
		goseq.Empty[int](),
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Ambas vacías produce resultado vacío.
func TestZip_BothEmpty_ReturnsEmpty(t *testing.T) {
	result := goseq.Zip(goseq.Empty[int](), goseq.Empty[string]()).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Zip con tipos distintos funciona correctamente.
func TestZip_DifferentTypes(t *testing.T) {
	result := goseq.Zip(
		goseq.From([]string{"a", "b", "c"}),
		goseq.From([]bool{true, false, true}),
	).ToSlice()

	if len(result) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(result))
	}
	if result[0].First != "a" || result[0].Second != true {
		t.Errorf("unexpected pair[0]: %v", result[0])
	}
}

// Preserva el orden de los pares.
func TestZip_PreservesOrder(t *testing.T) {
	keys := goseq.From([]int{3, 1, 2})
	vals := goseq.From([]string{"c", "a", "b"})
	result := goseq.Zip(keys, vals).ToSlice()

	if result[0].First != 3 || result[1].First != 1 || result[2].First != 2 {
		t.Errorf("order not preserved: %v", result)
	}
}
