package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// ToMap
// =====================

// --- ToMap: casos básicos ---

// Indexar enteros por sí mismos, con el cuadrado como valor.
func TestToMap_IntKeyIntValue(t *testing.T) {
	result := goseq.ToMap(
		goseq.From([]int{1, 2, 3}),
		func(n int) int { return n },
		func(n int) int { return n * n },
	)

	expected := map[int]int{1: 1, 2: 4, 3: 9}
	assertMapEqual(t, expected, result)
}

// Indexar strings por su longitud.
func TestToMap_StringKeyIntValue(t *testing.T) {
	result := goseq.ToMap(
		goseq.From([]string{"go", "rust", "java"}),
		func(s string) string { return s },
		func(s string) int { return len(s) },
	)

	expected := map[string]int{"go": 2, "rust": 4, "java": 4}
	assertMapEqual(t, expected, result)
}

// Extraer campos de una struct para construir el mapa.
func TestToMap_StructToMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	result := goseq.ToMap(
		goseq.From([]Person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 40}}),
		func(p Person) string { return p.Name },
		func(p Person) int { return p.Age },
	)

	expected := map[string]int{"Alice": 30, "Bob": 25, "Charlie": 40}
	assertMapEqual(t, expected, result)
}

// --- ToMap: casos borde ---

// Secuencia vacía produce mapa vacío.
func TestToMap_EmptySeq_ReturnsEmptyMap(t *testing.T) {
	result := goseq.ToMap(
		goseq.Empty[int](),
		func(n int) int { return n },
		func(n int) int { return n },
	)

	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

// Un único elemento produce mapa de un solo par.
func TestToMap_SingleElement(t *testing.T) {
	result := goseq.ToMap(
		goseq.From([]int{7}),
		func(n int) int { return n },
		func(n int) string { return "seven" },
	)

	if len(result) != 1 {
		t.Fatalf("expected map of length 1, got %d", len(result))
	}
	if result[7] != "seven" {
		t.Errorf("expected result[7] = 'seven', got '%s'", result[7])
	}
}

// Cuando hay claves duplicadas, el último elemento gana.
func TestToMap_DuplicateKeys_LastWins(t *testing.T) {
	type Item struct {
		Key   string
		Value int
	}

	result := goseq.ToMap(
		goseq.From([]Item{{"a", 1}, {"b", 2}, {"a", 99}}),
		func(i Item) string { return i.Key },
		func(i Item) int { return i.Value },
	)

	if result["a"] != 99 {
		t.Errorf("expected result['a'] = 99 (last wins), got %d", result["a"])
	}
	if result["b"] != 2 {
		t.Errorf("expected result['b'] = 2, got %d", result["b"])
	}
}

// --- ToMap: no mutación ---

func TestToMap_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3})
	_ = goseq.ToMap(original,
		func(n int) int { return n },
		func(n int) int { return n * 2 },
	)

	assertSliceEqual(t, []int{1, 2, 3}, original.ToSlice())
}

// --- ToMap: encadenamiento con Filter ---

func TestToMap_AfterFilter(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	result := goseq.ToMap(
		goseq.From([]Person{{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}}).
			Filter(func(p Person) bool { return p.Age >= 18 }),
		func(p Person) string { return p.Name },
		func(p Person) int { return p.Age },
	)

	if _, ok := result["Bob"]; ok {
		t.Error("Bob should have been filtered out")
	}
	if result["Alice"] != 30 {
		t.Errorf("expected Alice=30, got %d", result["Alice"])
	}
	if result["Charlie"] != 25 {
		t.Errorf("expected Charlie=25, got %d", result["Charlie"])
	}
}

// =====================
// GroupBy
// =====================

// --- GroupBy: casos básicos ---

// Agrupar enteros por paridad.
func TestGroupBy_ByParity(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{1, 2, 3, 4, 5, 6}),
		func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		},
	)

	assertSliceEqual(t, []int{1, 3, 5}, result["odd"])
	assertSliceEqual(t, []int{2, 4, 6}, result["even"])
}

// Agrupar strings por su primera letra.
func TestGroupBy_StringsByFirstLetter(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]string{"Alice", "Bob", "Ana", "Brian", "Charlie"}),
		func(s string) byte { return s[0] },
	)

	assertSliceEqual(t, []string{"Alice", "Ana"}, result['A'])
	assertSliceEqual(t, []string{"Bob", "Brian"}, result['B'])
	assertSliceEqual(t, []string{"Charlie"}, result['C'])
}

// Agrupar structs por un campo.
func TestGroupBy_StructsByField(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}, {"Diana", 15},
	}

	result := goseq.GroupBy(
		goseq.From(people),
		func(p Person) string {
			if p.Age >= 18 {
				return "adult"
			}
			return "minor"
		},
	)

	if len(result["adult"]) != 2 {
		t.Errorf("expected 2 adults, got %d", len(result["adult"]))
	}
	if len(result["minor"]) != 2 {
		t.Errorf("expected 2 minors, got %d", len(result["minor"]))
	}
}

// --- GroupBy: preservación de orden dentro de cada grupo ---

func TestGroupBy_PreservesOrderWithinGroups(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{5, 2, 8, 1, 4, 9, 6}),
		func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		},
	)

	// Los pares deben aparecer en el orden en que estaban en el input
	assertSliceEqual(t, []int{2, 8, 4, 6}, result["even"])
	assertSliceEqual(t, []int{5, 1, 9}, result["odd"])
}

// --- GroupBy: casos borde ---

// Secuencia vacía produce mapa vacío.
func TestGroupBy_EmptySeq_ReturnsEmptyMap(t *testing.T) {
	result := goseq.GroupBy(
		goseq.Empty[int](),
		func(n int) string { return "key" },
	)

	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

// Un único elemento produce un mapa con un solo grupo de un elemento.
func TestGroupBy_SingleElement(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{42}),
		func(n int) string { return "only" },
	)

	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}
	assertSliceEqual(t, []int{42}, result["only"])
}

// Todos los elementos en el mismo grupo.
func TestGroupBy_AllElementsSameKey(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{1, 2, 3, 4}),
		func(n int) string { return "all" },
	)

	if len(result) != 1 {
		t.Errorf("expected 1 group, got %d", len(result))
	}
	assertSliceEqual(t, []int{1, 2, 3, 4}, result["all"])
}

// Cada elemento en su propio grupo (todos distintos).
func TestGroupBy_AllElementsDifferentKey(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{1, 2, 3}),
		func(n int) int { return n },
	)

	if len(result) != 3 {
		t.Errorf("expected 3 groups, got %d", len(result))
	}
	for _, n := range []int{1, 2, 3} {
		if len(result[n]) != 1 || result[n][0] != n {
			t.Errorf("expected group[%d] = [%d], got %v", n, n, result[n])
		}
	}
}

// --- GroupBy: no mutación ---

func TestGroupBy_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3, 4, 5})
	_ = goseq.GroupBy(original, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})

	assertSliceEqual(t, []int{1, 2, 3, 4, 5}, original.ToSlice())
}

// --- GroupBy: encadenamiento con Filter ---

func TestGroupBy_AfterFilter(t *testing.T) {
	result := goseq.GroupBy(
		goseq.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
			Filter(func(n int) bool { return n > 5 }),
		func(n int) string {
			if n%2 == 0 {
				return "even"
			}
			return "odd"
		},
	)

	// Solo números > 5: 6, 7, 8, 9, 10
	assertSliceEqual(t, []int{6, 8, 10}, result["even"])
	assertSliceEqual(t, []int{7, 9}, result["odd"])
}

// =====================
// helpers
// =====================

// assertMapEqual compara dos mapas entrada a entrada.
func assertMapEqual[K comparable, V comparable](t *testing.T, expected, got map[K]V) {
	t.Helper()

	if len(expected) != len(got) {
		t.Errorf("map length mismatch: expected %d, got %d\n  expected: %v\n  got:      %v",
			len(expected), len(got), expected, got)
		return
	}

	for k, v := range expected {
		if gotV, ok := got[k]; !ok {
			t.Errorf("missing key %v in result", k)
		} else if gotV != v {
			t.Errorf("for key %v: expected %v, got %v", k, v, gotV)
		}
	}
}
