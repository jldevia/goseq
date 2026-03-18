package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// OrderBy
// =====================

func TestOrderBy_SortsIntegersAscending(t *testing.T) {
	result := goseq.OrderBy(
		goseq.From([]int{3, 1, 4, 1, 5, 9, 2, 6}),
		func(n int) int { return n },
	).ToSlice()

	assertSliceEqual(t, []int{1, 1, 2, 3, 4, 5, 6, 9}, result)
}

func TestOrderBy_SortsStringsByLength(t *testing.T) {
	result := goseq.OrderBy(
		goseq.From([]string{"banana", "go", "rust", "python"}),
		func(s string) int { return len(s) },
	).ToSlice()

	assertSliceEqual(t, []string{"go", "rust", "banana", "python"}, result)
}

func TestOrderBy_SortsStructsByField(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	result := goseq.OrderBy(
		goseq.From([]Person{{"Charlie", 40}, {"Alice", 25}, {"Bob", 30}}),
		func(p Person) int { return p.Age },
	).ToSlice()

	if result[0].Name != "Alice" || result[1].Name != "Bob" || result[2].Name != "Charlie" {
		t.Errorf("unexpected order: %v", result)
	}
}

// Estabilidad: elementos con la misma clave mantienen su orden relativo original.
func TestOrderBy_IsStable(t *testing.T) {
	type Item struct {
		Name     string
		Priority int
	}

	result := goseq.OrderBy(
		goseq.From([]Item{{"a", 2}, {"b", 1}, {"c", 2}, {"d", 1}}),
		func(i Item) int { return i.Priority },
	).ToSlice()

	// Priority 1: b antes que d (orden original)
	// Priority 2: a antes que c (orden original)
	if result[0].Name != "b" || result[1].Name != "d" {
		t.Errorf("stable sort violated for priority 1: got %v, %v", result[0].Name, result[1].Name)
	}
	if result[2].Name != "a" || result[3].Name != "c" {
		t.Errorf("stable sort violated for priority 2: got %v, %v", result[2].Name, result[3].Name)
	}
}

func TestOrderBy_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.OrderBy(
		goseq.Empty[int](),
		func(n int) int { return n },
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestOrderBy_SingleElement_ReturnsSame(t *testing.T) {
	result := goseq.OrderBy(
		goseq.From([]int{42}),
		func(n int) int { return n },
	).ToSlice()

	assertSliceEqual(t, []int{42}, result)
}

// OrderBy no modifica la secuencia original.
func TestOrderBy_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{3, 1, 2})
	_ = goseq.OrderBy(original, func(n int) int { return n })

	assertSliceEqual(t, []int{3, 1, 2}, original.ToSlice())
}

// =====================
// OrderByDescending
// =====================

func TestOrderByDescending_SortsIntegersDescending(t *testing.T) {
	result := goseq.OrderByDescending(
		goseq.From([]int{3, 1, 4, 1, 5, 9, 2, 6}),
		func(n int) int { return n },
	).ToSlice()

	assertSliceEqual(t, []int{9, 6, 5, 4, 3, 2, 1, 1}, result)
}

func TestOrderByDescending_SortsStructsByFieldDescending(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	result := goseq.OrderByDescending(
		goseq.From([]Person{{"Alice", 25}, {"Bob", 30}, {"Charlie", 40}}),
		func(p Person) int { return p.Age },
	).ToSlice()

	if result[0].Name != "Charlie" || result[1].Name != "Bob" || result[2].Name != "Alice" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestOrderByDescending_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.OrderByDescending(
		goseq.Empty[int](),
		func(n int) int { return n },
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// OrderBy y OrderByDescending producen resultados inversos entre sí.
func TestOrderBy_AndDescending_AreInverse(t *testing.T) {
	input := goseq.From([]int{3, 1, 4, 2, 5})

	asc := goseq.OrderBy(input, func(n int) int { return n }).ToSlice()
	desc := goseq.OrderByDescending(input, func(n int) int { return n }).ToSlice()

	for i := range asc {
		if asc[i] != desc[len(desc)-1-i] {
			t.Errorf("asc and desc are not inverse at index %d", i)
		}
	}
}

// =====================
// Sum
// =====================

func TestSum_SumsIntegers(t *testing.T) {
	result := goseq.Sum(goseq.From([]int{1, 2, 3, 4, 5}))

	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestSum_SumsFloats(t *testing.T) {
	result := goseq.Sum(goseq.From([]float64{1.5, 2.5, 3.0}))

	if result != 7.0 {
		t.Errorf("expected 7.0, got %f", result)
	}
}

func TestSum_EmptySeq_ReturnsZero(t *testing.T) {
	result := goseq.Sum(goseq.Empty[int]())

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestSum_SingleElement(t *testing.T) {
	result := goseq.Sum(goseq.From([]int{42}))

	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestSum_SumsStrings(t *testing.T) {
	result := goseq.Sum(goseq.From([]string{"hello", " ", "world"}))

	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

// =====================
// Min
// =====================

func TestMin_ReturnsMinimum(t *testing.T) {
	result, ok := goseq.Min(goseq.From([]int{3, 1, 4, 1, 5, 9}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

func TestMin_EmptySeq_ReturnsFalse(t *testing.T) {
	_, ok := goseq.Min(goseq.Empty[int]())

	if ok {
		t.Error("expected ok=false for empty seq")
	}
}

func TestMin_SingleElement(t *testing.T) {
	result, ok := goseq.Min(goseq.From([]int{42}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestMin_WorksWithStrings(t *testing.T) {
	result, ok := goseq.Min(goseq.From([]string{"banana", "apple", "cherry"}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != "apple" {
		t.Errorf("expected 'apple', got '%s'", result)
	}
}

func TestMin_AllSameValue(t *testing.T) {
	result, ok := goseq.Min(goseq.From([]int{5, 5, 5}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

// =====================
// Max
// =====================

func TestMax_ReturnsMaximum(t *testing.T) {
	result, ok := goseq.Max(goseq.From([]int{3, 1, 4, 1, 5, 9}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != 9 {
		t.Errorf("expected 9, got %d", result)
	}
}

func TestMax_EmptySeq_ReturnsFalse(t *testing.T) {
	_, ok := goseq.Max(goseq.Empty[int]())

	if ok {
		t.Error("expected ok=false for empty seq")
	}
}

func TestMax_SingleElement(t *testing.T) {
	result, ok := goseq.Max(goseq.From([]int{42}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestMax_WorksWithStrings(t *testing.T) {
	result, ok := goseq.Max(goseq.From([]string{"banana", "apple", "cherry"}))

	if !ok {
		t.Fatal("expected ok=true")
	}
	if result != "cherry" {
		t.Errorf("expected 'cherry', got '%s'", result)
	}
}

// Min y Max son consistentes: Min <= Max siempre.
func TestMin_AndMax_AreConsistent(t *testing.T) {
	seq := goseq.From([]int{3, 1, 4, 1, 5, 9, 2, 6})
	min, _ := goseq.Min(seq)
	max, _ := goseq.Max(seq)

	if min > max {
		t.Errorf("Min (%d) should never be greater than Max (%d)", min, max)
	}
}

// =====================
// Pipelines combinados
// =====================

// FlatMap + OrderBy + Sum en un pipeline completo.
func TestAdvanced_FullPipeline(t *testing.T) {
	type Order struct {
		Items []int
	}

	orders := goseq.From([]Order{
		{Items: []int{10, 20}},
		{Items: []int{5, 15, 30}},
	})

	// Aplanar todos los items, filtrar los mayores a 10, sumar
	total := goseq.Sum(
		goseq.FlatMap(orders, func(o Order) goseq.Seq[int] {
			return goseq.From(o.Items)
		}).Filter(func(n int) bool { return n > 10 }),
	)

	// 20 + 15 + 30 = 65
	if total != 65 {
		t.Errorf("expected 65, got %d", total)
	}
}
