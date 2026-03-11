package goseq_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/jldevia/goseq"
)

// --- Map: casos básicos ---

// Transformación numérica simple: duplicar cada elemento.
func TestMap_DoublesIntegers(t *testing.T) {
	result := goseq.Map(
		goseq.From([]int{1, 2, 3, 4}),
		func(n int) int { return n * 2 },
	).ToSlice()

	expected := []int{2, 4, 6, 8}
	assertSliceEqual(t, expected, result)
}

// Transformación que mantiene el mismo tipo pero cambia el valor.
func TestMap_UppercasesStrings(t *testing.T) {
	result := goseq.Map(
		goseq.From([]string{"hello", "world"}),
		strings.ToUpper,
	).ToSlice()

	expected := []string{"HELLO", "WORLD"}
	assertSliceEqual(t, expected, result)
}

// Transformación que cambia el tipo: int → string.
func TestMap_IntToString(t *testing.T) {
	result := goseq.Map(
		goseq.From([]int{1, 2, 3}),
		strconv.Itoa,
	).ToSlice()

	expected := []string{"1", "2", "3"}
	assertSliceEqual(t, expected, result)
}

// Transformación que cambia el tipo: string → int (longitud).
func TestMap_StringToLength(t *testing.T) {
	result := goseq.Map(
		goseq.From([]string{"go", "rust", "python"}),
		func(s string) int { return len(s) },
	).ToSlice()

	expected := []int{2, 4, 6}
	assertSliceEqual(t, expected, result)
}

// --- Map: casos borde ---

// Map sobre secuencia vacía debe devolver secuencia vacía sin errores.
func TestMap_EmptySeq_ReturnsEmptySeq(t *testing.T) {
	result := goseq.Map(
		goseq.Empty[int](),
		func(n int) int { return n * 2 },
	).ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// Map sobre un único elemento.
func TestMap_SingleElement(t *testing.T) {
	result := goseq.Map(
		goseq.From([]int{7}),
		func(n int) int { return n * n },
	).ToSlice()

	expected := []int{49}
	assertSliceEqual(t, expected, result)
}

// Función identidad: la secuencia resultante debe ser igual a la original.
func TestMap_IdentityFunction_ReturnsSameValues(t *testing.T) {
	input := []int{10, 20, 30}
	result := goseq.Map(
		goseq.From(input),
		func(n int) int { return n },
	).ToSlice()

	assertSliceEqual(t, input, result)
}

// --- Map: preservación de orden ---

// Los elementos transformados deben mantener el mismo orden que los originales.
func TestMap_PreservesOrder(t *testing.T) {
	result := goseq.Map(
		goseq.From([]int{5, 3, 1, 4, 2}),
		func(n int) int { return n * 10 },
	).ToSlice()

	expected := []int{50, 30, 10, 40, 20}
	assertSliceEqual(t, expected, result)
}

// --- Map: no mutación ---

// Map no debe modificar la secuencia original.
func TestMap_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3})
	_ = goseq.Map(original, func(n int) int { return n * 100 })

	result := original.ToSlice()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

// --- Map: con structs ---

// Extraer un campo de una struct (string).
func TestMap_ExtractsNameFromStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := goseq.From([]Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 40},
	})

	names := goseq.Map(people, func(p Person) string { return p.Name }).ToSlice()

	assertSliceEqual(t, []string{"Alice", "Bob", "Charlie"}, names)
}

// Extraer un campo de una struct (int).
func TestMap_ExtractsAgeFromStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := goseq.From([]Person{
		{"Alice", 30},
		{"Bob", 25},
	})

	ages := goseq.Map(people, func(p Person) int { return p.Age }).ToSlice()

	assertSliceEqual(t, []int{30, 25}, ages)
}

// Transformar una struct en otra struct de distinto tipo.
func TestMap_TransformsStructToAnotherStruct(t *testing.T) {
	type Input struct{ Value int }
	type Output struct{ Doubled int }

	result := goseq.Map(
		goseq.From([]Input{{1}, {2}, {3}}),
		func(i Input) Output { return Output{Doubled: i.Value * 2} },
	).ToSlice()

	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
	if result[0].Doubled != 2 || result[1].Doubled != 4 || result[2].Doubled != 6 {
		t.Errorf("unexpected result: %v", result)
	}
}

// --- Map: encadenamiento con Filter ---

// Map y Filter deben poder combinarse en un pipeline.
func TestMap_ChainedWithFilter(t *testing.T) {
	// Filtrar pares, luego duplicarlos
	result := goseq.Map(
		goseq.From([]int{1, 2, 3, 4, 5, 6}).
			Filter(func(n int) bool { return n%2 == 0 }),
		func(n int) int { return n * 2 },
	).ToSlice()

	expected := []int{4, 8, 12}
	assertSliceEqual(t, expected, result)
}

// Aplicar Map dos veces seguidas.
func TestMap_ChainedWithAnotherMap(t *testing.T) {
	// Primero duplicar, luego sumar 1
	step1 := goseq.Map(
		goseq.From([]int{1, 2, 3}),
		func(n int) int { return n * 2 },
	)
	result := goseq.Map(step1, func(n int) int { return n + 1 }).ToSlice()

	expected := []int{3, 5, 7}
	assertSliceEqual(t, expected, result)
}

// --- Map: el resultado es un Seq funcional ---

// Len debe funcionar sobre el resultado de Map.
func TestMap_ResultIsSeq_LenWorks(t *testing.T) {
	mapped := goseq.Map(
		goseq.From([]int{1, 2, 3}),
		func(n int) int { return n * 2 },
	)

	if mapped.Len() != 3 {
		t.Errorf("expected Len() = 3, got %d", mapped.Len())
	}
}

// IsEmpty debe funcionar sobre el resultado de Map.
func TestMap_ResultIsSeq_IsEmptyWorks(t *testing.T) {
	mapped := goseq.Map(
		goseq.Empty[int](),
		func(n int) int { return n * 2 },
	)

	if !mapped.IsEmpty() {
		t.Error("expected IsEmpty() = true")
	}
}
