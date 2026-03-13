package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// First
// =====================

func TestFirst_NonEmptySeq_ReturnsFirstElement(t *testing.T) {
	val, ok := goseq.From([]int{10, 20, 30}).First()

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != 10 {
		t.Errorf("expected 10, got %d", val)
	}
}

func TestFirst_EmptySeq_ReturnsFalse(t *testing.T) {
	val, ok := goseq.Empty[int]().First()

	if ok {
		t.Error("expected ok=false for empty seq")
	}
	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}
}

func TestFirst_SingleElement(t *testing.T) {
	val, ok := goseq.From([]string{"only"}).First()

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != "only" {
		t.Errorf("expected 'only', got '%s'", val)
	}
}

// =====================
// FirstWhere
// =====================

func TestFirstWhere_ReturnsFirstMatch(t *testing.T) {
	val, ok := goseq.From([]int{1, 2, 3, 4, 5}).
		FirstWhere(func(n int) bool { return n%2 == 0 })

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}

func TestFirstWhere_NoMatch_ReturnsFalse(t *testing.T) {
	val, ok := goseq.From([]int{1, 3, 5}).
		FirstWhere(func(n int) bool { return n%2 == 0 })

	if ok {
		t.Error("expected ok=false")
	}
	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}
}

func TestFirstWhere_EmptySeq_ReturnsFalse(t *testing.T) {
	_, ok := goseq.Empty[int]().FirstWhere(func(n int) bool { return true })

	if ok {
		t.Error("expected ok=false for empty seq")
	}
}

func TestFirstWhere_ReturnsFirstNotLast(t *testing.T) {
	val, _ := goseq.From([]int{2, 4, 6}).
		FirstWhere(func(n int) bool { return n%2 == 0 })

	if val != 2 {
		t.Errorf("expected first match (2), got %d", val)
	}
}

// =====================
// Last
// =====================

func TestLast_NonEmptySeq_ReturnsLastElement(t *testing.T) {
	val, ok := goseq.From([]int{10, 20, 30}).Last()

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != 30 {
		t.Errorf("expected 30, got %d", val)
	}
}

func TestLast_EmptySeq_ReturnsFalse(t *testing.T) {
	val, ok := goseq.Empty[int]().Last()

	if ok {
		t.Error("expected ok=false for empty seq")
	}
	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}
}

func TestLast_SingleElement(t *testing.T) {
	val, ok := goseq.From([]int{99}).Last()

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != 99 {
		t.Errorf("expected 99, got %d", val)
	}
}

// =====================
// LastWhere
// =====================

func TestLastWhere_ReturnsLastMatch(t *testing.T) {
	val, ok := goseq.From([]int{1, 2, 3, 4, 5}).
		LastWhere(func(n int) bool { return n%2 == 0 })

	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != 4 {
		t.Errorf("expected 4, got %d", val)
	}
}

func TestLastWhere_NoMatch_ReturnsFalse(t *testing.T) {
	_, ok := goseq.From([]int{1, 3, 5}).
		LastWhere(func(n int) bool { return n%2 == 0 })

	if ok {
		t.Error("expected ok=false")
	}
}

func TestLastWhere_ReturnsLastNotFirst(t *testing.T) {
	val, _ := goseq.From([]int{2, 4, 6}).
		LastWhere(func(n int) bool { return n%2 == 0 })

	if val != 6 {
		t.Errorf("expected last match (6), got %d", val)
	}
}

// =====================
// Any
// =====================

func TestAny_MatchExists_ReturnsTrue(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		Any(func(n int) bool { return n > 2 })

	if !result {
		t.Error("expected true")
	}
}

func TestAny_NoMatch_ReturnsFalse(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		Any(func(n int) bool { return n > 10 })

	if result {
		t.Error("expected false")
	}
}

func TestAny_EmptySeq_ReturnsFalse(t *testing.T) {
	result := goseq.Empty[int]().Any(func(n int) bool { return true })

	if result {
		t.Error("expected false for empty seq")
	}
}

// =====================
// All
// =====================

func TestAll_AllMatch_ReturnsTrue(t *testing.T) {
	result := goseq.From([]int{2, 4, 6}).
		All(func(n int) bool { return n%2 == 0 })

	if !result {
		t.Error("expected true")
	}
}

func TestAll_OneDoesNotMatch_ReturnsFalse(t *testing.T) {
	result := goseq.From([]int{2, 4, 5}).
		All(func(n int) bool { return n%2 == 0 })

	if result {
		t.Error("expected false")
	}
}

func TestAll_EmptySeq_ReturnsTrue(t *testing.T) {
	// Vacuous truth: all elements of an empty set satisfy any predicate
	result := goseq.Empty[int]().All(func(n int) bool { return false })

	if !result {
		t.Error("expected true for empty seq (vacuous truth)")
	}
}

// =====================
// None
// =====================

func TestNone_NoMatchExists_ReturnsTrue(t *testing.T) {
	result := goseq.From([]int{1, 3, 5}).
		None(func(n int) bool { return n%2 == 0 })

	if !result {
		t.Error("expected true")
	}
}

func TestNone_MatchExists_ReturnsFalse(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		None(func(n int) bool { return n%2 == 0 })

	if result {
		t.Error("expected false")
	}
}

func TestNone_EmptySeq_ReturnsTrue(t *testing.T) {
	result := goseq.Empty[int]().None(func(n int) bool { return true })

	if !result {
		t.Error("expected true for empty seq")
	}
}

// None debe ser consistente con Any: None == !Any
func TestNone_ConsistentWithAny(t *testing.T) {
	seq := goseq.From([]int{1, 2, 3, 4, 5})
	predicate := func(n int) bool { return n > 3 }

	if seq.None(predicate) == seq.Any(predicate) {
		t.Error("None and Any should always return opposite values")
	}
}

// =====================
// Count
// =====================

func TestCount_CountsMatchingElements(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5, 6}).
		Count(func(n int) bool { return n%2 == 0 })

	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

func TestCount_NoneMatch_ReturnsZero(t *testing.T) {
	result := goseq.From([]int{1, 3, 5}).
		Count(func(n int) bool { return n%2 == 0 })

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestCount_AllMatch(t *testing.T) {
	result := goseq.From([]int{2, 4, 6}).
		Count(func(n int) bool { return n%2 == 0 })

	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

func TestCount_EmptySeq_ReturnsZero(t *testing.T) {
	result := goseq.Empty[int]().Count(func(n int) bool { return true })

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

// Count con predicado siempre true debe ser igual a Len
func TestCount_AllTruePredicate_EqualsLen(t *testing.T) {
	seq := goseq.From([]int{1, 2, 3, 4, 5})

	if seq.Count(func(n int) bool { return true }) != seq.Len() {
		t.Error("Count(always true) should equal Len()")
	}
}

// =====================
// Contains
// =====================

func TestContains_ElementExists_ReturnsTrue(t *testing.T) {
	result := goseq.Contains(goseq.From([]int{1, 2, 3}), 2)

	if !result {
		t.Error("expected true")
	}
}

func TestContains_ElementNotExists_ReturnsFalse(t *testing.T) {
	result := goseq.Contains(goseq.From([]int{1, 2, 3}), 99)

	if result {
		t.Error("expected false")
	}
}

func TestContains_EmptySeq_ReturnsFalse(t *testing.T) {
	result := goseq.Contains(goseq.Empty[int](), 1)

	if result {
		t.Error("expected false for empty seq")
	}
}

func TestContains_WorksWithStrings(t *testing.T) {
	result := goseq.Contains(goseq.From([]string{"go", "rust", "java"}), "rust")

	if !result {
		t.Error("expected true")
	}
}

func TestContains_FirstElement(t *testing.T) {
	result := goseq.Contains(goseq.From([]int{10, 20, 30}), 10)

	if !result {
		t.Error("expected true for first element")
	}
}

func TestContains_LastElement(t *testing.T) {
	result := goseq.Contains(goseq.From([]int{10, 20, 30}), 30)

	if !result {
		t.Error("expected true for last element")
	}
}
