package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// --- Filter: casos básicos ---

// El predicado más simple: dejar pasar todos los elementos.
func TestFilter_PredicateAlwaysTrue_ReturnsAllElements(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := goseq.From(input).
		Filter(func(n int) bool { return true }).
		ToSlice()

	if len(result) != len(input) {
		t.Errorf("expected %d elements, got %d", len(input), len(result))
	}
}

// El predicado más restrictivo: no dejar pasar ningún elemento.
func TestFilter_PredicateAlwaysFalse_ReturnsEmptySeq(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		Filter(func(n int) bool { return false }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// Caso típico: filtrar por una condición real.
func TestFilter_EvenNumbers_ReturnsOnlyEvens(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToSlice()

	expected := []int{2, 4, 6}
	assertSliceEqual(t, expected, result)
}

func TestFilter_OddNumbers_ReturnsOnlyOdds(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5}).
		Filter(func(n int) bool { return n%2 != 0 }).
		ToSlice()

	expected := []int{1, 3, 5}
	assertSliceEqual(t, expected, result)
}

// --- Filter: casos borde ---

// Aplicar Filter sobre una secuencia vacía debe devolver una secuencia vacía,
// sin errores ni panics.
func TestFilter_EmptySeq_ReturnsEmptySeq(t *testing.T) {
	result := goseq.Empty[int]().
		Filter(func(n int) bool { return true }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// Secuencia con un solo elemento que cumple el predicado.
func TestFilter_SingleElement_Matches(t *testing.T) {
	result := goseq.From([]int{42}).
		Filter(func(n int) bool { return n > 0 }).
		ToSlice()

	expected := []int{42}
	assertSliceEqual(t, expected, result)
}

// Secuencia con un solo elemento que NO cumple el predicado.
func TestFilter_SingleElement_DoesNotMatch(t *testing.T) {
	result := goseq.From([]int{42}).
		Filter(func(n int) bool { return n < 0 }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// Solo el primer elemento cumple el predicado.
func TestFilter_OnlyFirstElementMatches(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		Filter(func(n int) bool { return n == 1 }).
		ToSlice()

	expected := []int{1}
	assertSliceEqual(t, expected, result)
}

// Solo el último elemento cumple el predicado.
func TestFilter_OnlyLastElementMatches(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		Filter(func(n int) bool { return n == 3 }).
		ToSlice()

	expected := []int{3}
	assertSliceEqual(t, expected, result)
}

// --- Filter: preservación de orden ---

// Los elementos que pasan el filtro deben mantener su orden relativo original.
func TestFilter_PreservesOrder(t *testing.T) {
	result := goseq.From([]int{5, 1, 4, 2, 3}).
		Filter(func(n int) bool { return n > 2 }).
		ToSlice()

	// Deben aparecer en el mismo orden que en el input: 5, 4, 3
	expected := []int{5, 4, 3}
	assertSliceEqual(t, expected, result)
}

// --- Filter: no mutación ---

// Filter no debe modificar la secuencia original.
func TestFilter_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3, 4, 5})
	_ = original.Filter(func(n int) bool { return n%2 == 0 })

	// La secuencia original debe seguir intacta
	result := original.ToSlice()
	expected := []int{1, 2, 3, 4, 5}
	assertSliceEqual(t, expected, result)
}

// --- Filter: tipos distintos ---

// Filter debe funcionar con strings.
func TestFilter_WithStrings_ByLength(t *testing.T) {
	result := goseq.From([]string{"go", "rust", "c", "python", "java"}).
		Filter(func(s string) bool { return len(s) > 3 }).
		ToSlice()

	expected := []string{"rust", "python", "java"}
	assertSliceEqual(t, expected, result)
}

func TestFilter_WithStrings_ByPrefix(t *testing.T) {
	result := goseq.From([]string{"Alice", "Bob", "Ana", "Brian", "Carlos"}).
		Filter(func(s string) bool { return s[0] == 'A' }).
		ToSlice()

	expected := []string{"Alice", "Ana"}
	assertSliceEqual(t, expected, result)
}

// Filter debe funcionar con structs.
func TestFilter_WithStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 17},
		{"Charlie", 25},
		{"Diana", 15},
	}

	result := goseq.From(people).
		Filter(func(p Person) bool { return p.Age >= 18 }).
		ToSlice()

	if len(result) != 2 {
		t.Fatalf("expected 2 adults, got %d", len(result))
	}
	if result[0].Name != "Alice" || result[1].Name != "Charlie" {
		t.Errorf("unexpected result: %v", result)
	}
}

// --- Filter: encadenamiento ---

// Filter debe poder encadenarse consigo mismo.
func TestFilter_Chained_MultipleFilters(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Filter(func(n int) bool { return n%2 == 0 }). // pares: 2,4,6,8,10
		Filter(func(n int) bool { return n > 5 }).    // mayores a 5: 6,8,10
		ToSlice()

	expected := []int{6, 8, 10}
	assertSliceEqual(t, expected, result)
}

// Filter devuelve un Seq[T], por lo que IsEmpty y Len deben funcionar
// sobre el resultado sin necesidad de llamar ToSlice.
func TestFilter_ResultIsSeq_LenWorks(t *testing.T) {
	filtered := goseq.From([]int{1, 2, 3, 4}).
		Filter(func(n int) bool { return n%2 == 0 })

	if filtered.Len() != 2 {
		t.Errorf("expected Len() = 2, got %d", filtered.Len())
	}
}

func TestFilter_ResultIsSeq_IsEmptyWorks(t *testing.T) {
	filtered := goseq.From([]int{1, 3, 5}).
		Filter(func(n int) bool { return n%2 == 0 })

	if !filtered.IsEmpty() {
		t.Error("expected IsEmpty() = true after filtering out all elements")
	}
}

// --- helpers ---

// assertSliceEqual compara dos slices elemento a elemento.
// Es un helper genérico para evitar repetición en los tests.
func assertSliceEqual[T comparable](t *testing.T, expected, got []T) {
	t.Helper()

	if len(expected) != len(got) {
		t.Errorf("length mismatch: expected %d, got %d\n  expected: %v\n  got:      %v",
			len(expected), len(got), expected, got)
		return
	}

	for i := range expected {
		if expected[i] != got[i] {
			t.Errorf("mismatch at index %d: expected %v, got %v\n  expected: %v\n  got:      %v",
				i, expected[i], got[i], expected, got)
		}
	}
}
