package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// --- From ---

func TestFrom_CreatesSeqWithCorrectElements(t *testing.T) {
	input := []int{1, 2, 3}
	s := goseq.From(input)

	if s.Len() != 3 {
		t.Errorf("expected Len() = 3, got %d", s.Len())
	}
}

func TestFrom_DoesNotModifyOriginalSlice(t *testing.T) {
	input := []int{1, 2, 3}
	s := goseq.From(input)
	result := s.ToSlice()

	// Modify the result and check the original is untouched
	result[0] = 999
	if input[0] != 1 {
		t.Errorf("original slice was modified: expected input[0] = 1, got %d", input[0])
	}
}

func TestFrom_EmptySlice(t *testing.T) {
	s := goseq.From([]string{})

	if !s.IsEmpty() {
		t.Error("expected IsEmpty() = true for empty slice")
	}
}

func TestFrom_WorksWithStrings(t *testing.T) {
	s := goseq.From([]string{"hello", "world"})

	if s.Len() != 2 {
		t.Errorf("expected Len() = 2, got %d", s.Len())
	}
}

// --- Empty ---

func TestEmpty_IsEmpty(t *testing.T) {
	s := goseq.Empty[int]()

	if !s.IsEmpty() {
		t.Error("expected IsEmpty() = true")
	}

	if s.Len() != 0 {
		t.Errorf("expected Len() = 0, got %d", s.Len())
	}
}

func TestEmpty_ToSliceReturnsEmptySlice(t *testing.T) {
	result := goseq.Empty[string]().ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// --- Len / IsEmpty ---

func TestLen_ReturnsCorrectCount(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"empty", []int{}, 0},
		{"one element", []int{42}, 1},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goseq.From(tt.input).Len()
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestIsEmpty_ReturnsTrueOnlyWhenEmpty(t *testing.T) {
	if !goseq.From([]int{}).IsEmpty() {
		t.Error("expected true for empty seq")
	}

	if goseq.From([]int{1}).IsEmpty() {
		t.Error("expected false for non-empty seq")
	}
}

// --- ToSlice ---

func TestToSlice_ReturnsAllElements(t *testing.T) {
	input := []int{10, 20, 30}
	result := goseq.From(input).ToSlice()

	if len(result) != len(input) {
		t.Fatalf("expected len %d, got %d", len(input), len(result))
	}

	for i, v := range input {
		if result[i] != v {
			t.Errorf("at index %d: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestToSlice_ReturnsCopy(t *testing.T) {
	s := goseq.From([]int{1, 2, 3})
	r1 := s.ToSlice()
	r1[0] = 999

	r2 := s.ToSlice()
	if r2[0] == 999 {
		t.Error("ToSlice should return a copy, but internal state was modified")
	}
}

// --- ForEach ---

func TestForEach_VisitsAllElements(t *testing.T) {
	input := []int{1, 2, 3}
	visited := []int{}

	goseq.From(input).ForEach(func(n int) {
		visited = append(visited, n)
	})

	if len(visited) != len(input) {
		t.Fatalf("expected %d visits, got %d", len(input), len(visited))
	}

	for i, v := range input {
		if visited[i] != v {
			t.Errorf("at index %d: expected %d, got %d", i, v, visited[i])
		}
	}
}

func TestForEach_EmptySeqDoesNothing(t *testing.T) {
	called := false
	goseq.Empty[int]().ForEach(func(n int) {
		called = true
	})

	if called {
		t.Error("ForEach should not call fn for an empty sequence")
	}
}
