package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// Take
// =====================

func TestTake_ReturnsFirstNElements(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTake_NGreaterThanLen_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Take(10).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTake_NEqualsLen_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Take(3).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTake_ZeroN_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Take(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestTake_NegativeN_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Take(-5).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestTake_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.Empty[int]().Take(3).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestTake_TakeOne(t *testing.T) {
	result := goseq.From([]int{10, 20, 30}).Take(1).ToSlice()
	assertSliceEqual(t, []int{10}, result)
}

// =====================
// Skip
// =====================

func TestSkip_SkipsFirstNElements(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice()
	assertSliceEqual(t, []int{3, 4, 5}, result)
}

func TestSkip_NGreaterThanLen_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Skip(10).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestSkip_NEqualsLen_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Skip(3).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestSkip_ZeroN_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Skip(0).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestSkip_NegativeN_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).Skip(-5).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestSkip_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.Empty[int]().Skip(3).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// Take y Skip son complementarios: Take(n) + Skip(n) == secuencia original
func TestTakeAndSkip_AreComplementary(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	seq := goseq.From(input)

	taken := seq.Take(3).ToSlice()
	skipped := seq.Skip(3).ToSlice()
	combined := append(taken, skipped...)

	assertSliceEqual(t, input, combined)
}

// =====================
// TakeWhile
// =====================

func TestTakeWhile_TakesWhilePredicateTrue(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5}).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToSlice()

	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTakeWhile_PredicateFalseFromStart_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{5, 6, 7}).
		TakeWhile(func(n int) bool { return n < 3 }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestTakeWhile_PredicateAlwaysTrue_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		TakeWhile(func(n int) bool { return true }).
		ToSlice()

	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTakeWhile_StopsAtFirstFalse(t *testing.T) {
	// El 4 rompe la condición, aunque el 2 posterior la cumpliría
	result := goseq.From([]int{1, 2, 4, 2, 1}).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToSlice()

	assertSliceEqual(t, []int{1, 2}, result)
}

func TestTakeWhile_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.Empty[int]().
		TakeWhile(func(n int) bool { return true }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// =====================
// SkipWhile
// =====================

func TestSkipWhile_SkipsWhilePredicateTrue(t *testing.T) {
	result := goseq.From([]int{1, 2, 3, 4, 5}).
		SkipWhile(func(n int) bool { return n < 3 }).
		ToSlice()

	assertSliceEqual(t, []int{3, 4, 5}, result)
}

func TestSkipWhile_PredicateFalseFromStart_ReturnsAll(t *testing.T) {
	result := goseq.From([]int{5, 6, 7}).
		SkipWhile(func(n int) bool { return n < 3 }).
		ToSlice()

	assertSliceEqual(t, []int{5, 6, 7}, result)
}

func TestSkipWhile_PredicateAlwaysTrue_ReturnsEmpty(t *testing.T) {
	result := goseq.From([]int{1, 2, 3}).
		SkipWhile(func(n int) bool { return true }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// SkipWhile no vuelve a saltar elementos que cumplan el predicado después
// de haber encontrado uno que no lo cumple
func TestSkipWhile_DoesNotSkipAfterFirstFalse(t *testing.T) {
	result := goseq.From([]int{1, 2, 4, 1, 2}).
		SkipWhile(func(n int) bool { return n < 4 }).
		ToSlice()

	// Salta 1 y 2, se detiene en 4, incluye 4, 1, 2
	assertSliceEqual(t, []int{4, 1, 2}, result)
}

func TestSkipWhile_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.Empty[int]().
		SkipWhile(func(n int) bool { return true }).
		ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

// =====================
// Distinct
// =====================

func TestDistinct_RemovesDuplicates(t *testing.T) {
	result := goseq.Distinct(goseq.From([]int{1, 2, 2, 3, 1, 4})).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3, 4}, result)
}

func TestDistinct_PreservesOrderOfFirstAppearance(t *testing.T) {
	result := goseq.Distinct(goseq.From([]int{3, 1, 2, 1, 3})).ToSlice()
	assertSliceEqual(t, []int{3, 1, 2}, result)
}

func TestDistinct_NoDuplicates_ReturnsSame(t *testing.T) {
	result := goseq.Distinct(goseq.From([]int{1, 2, 3})).ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestDistinct_AllDuplicates_ReturnsSingleElement(t *testing.T) {
	result := goseq.Distinct(goseq.From([]int{5, 5, 5, 5})).ToSlice()
	assertSliceEqual(t, []int{5}, result)
}

func TestDistinct_EmptySeq_ReturnsEmpty(t *testing.T) {
	result := goseq.Distinct(goseq.Empty[int]()).ToSlice()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestDistinct_WorksWithStrings(t *testing.T) {
	result := goseq.Distinct(
		goseq.From([]string{"go", "rust", "go", "java", "rust"}),
	).ToSlice()
	assertSliceEqual(t, []string{"go", "rust", "java"}, result)
}

func TestDistinct_SingleElement_ReturnsSame(t *testing.T) {
	result := goseq.Distinct(goseq.From([]int{42})).ToSlice()
	assertSliceEqual(t, []int{42}, result)
}
