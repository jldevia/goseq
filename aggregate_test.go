package goseq_test

import (
	"strings"
	"testing"

	"github.com/jldevia/goseq"
)

// --- Reduce: casos básicos ---

// El caso más clásico: sumar una secuencia de enteros.
func TestReduce_SumOfIntegers(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3, 4, 5}),
		0,
		func(acc, n int) int { return acc + n },
	)

	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

// Multiplicar todos los elementos (producto).
func TestReduce_ProductOfIntegers(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3, 4, 5}),
		1,
		func(acc, n int) int { return acc * n },
	)

	if result != 120 {
		t.Errorf("expected 120, got %d", result)
	}
}

// Concatenar strings.
func TestReduce_ConcatenatesStrings(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]string{"hello", " ", "world"}),
		"",
		func(acc, s string) string { return acc + s },
	)

	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

// Contar elementos que cumplen una condición.
func TestReduce_CountsEvenNumbers(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3, 4, 5, 6}),
		0,
		func(acc, n int) int {
			if n%2 == 0 {
				return acc + 1
			}
			return acc
		},
	)

	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

// --- Reduce: casos borde ---

// Sobre secuencia vacía debe devolver el valor inicial.
func TestReduce_EmptySeq_ReturnsInitialValue(t *testing.T) {
	result := goseq.Reduce(
		goseq.Empty[int](),
		42,
		func(acc, n int) int { return acc + n },
	)

	if result != 42 {
		t.Errorf("expected initial value 42, got %d", result)
	}
}

// Sobre un único elemento debe aplicar la función una sola vez.
func TestReduce_SingleElement(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{7}),
		0,
		func(acc, n int) int { return acc + n },
	)

	if result != 7 {
		t.Errorf("expected 7, got %d", result)
	}
}

// El valor inicial importa: sumar partiendo de 100.
func TestReduce_InitialValueAffectsResult(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3}),
		100,
		func(acc, n int) int { return acc + n },
	)

	if result != 106 {
		t.Errorf("expected 106, got %d", result)
	}
}

// --- Reduce: cambio de tipo (T → U) ---

// Reducir []int a un string (tipo completamente distinto).
func TestReduce_IntSeqToString(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3}),
		"nums:",
		func(acc string, n int) string {
			return acc + " " + strings.TrimSpace(strings.Repeat("x", n))
		},
	)

	// Verifica que el resultado empieza con el prefijo correcto
	if !strings.HasPrefix(result, "nums:") {
		t.Errorf("expected result to start with 'nums:', got '%s'", result)
	}
}

// Reducir []string a un int (conteo de caracteres totales).
func TestReduce_StringSeqToTotalLength(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]string{"go", "is", "great"}),
		0,
		func(acc int, s string) int { return acc + len(s) },
	)

	// "go"=2, "is"=2, "great"=5 → 9
	if result != 9 {
		t.Errorf("expected 9, got %d", result)
	}
}

// Reducir a un slice: recolectar solo los elementos que cumplen una condición.
func TestReduce_CollectsEvenNumbersIntoSlice(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3, 4, 5, 6}),
		[]int{},
		func(acc []int, n int) []int {
			if n%2 == 0 {
				return append(acc, n)
			}
			return acc
		},
	)

	expected := []int{2, 4, 6}
	assertSliceEqual(t, expected, result)
}

// Reducir structs a un valor agregado.
func TestReduce_SumsStructField(t *testing.T) {
	type Product struct {
		Name  string
		Price float64
	}

	products := goseq.From([]Product{
		{"Apple", 1.50},
		{"Bread", 2.30},
		{"Milk", 0.99},
	})

	total := goseq.Reduce(products, 0.0,
		func(acc float64, p Product) float64 { return acc + p.Price },
	)

	expected := 4.79
	if total < expected-0.001 || total > expected+0.001 {
		t.Errorf("expected %.2f, got %.2f", expected, total)
	}
}

// --- Reduce: no mutación ---

// Reduce no debe modificar la secuencia original.
func TestReduce_DoesNotModifyOriginalSeq(t *testing.T) {
	original := goseq.From([]int{1, 2, 3, 4, 5})
	_ = goseq.Reduce(original, 0, func(acc, n int) int { return acc + n })

	result := original.ToSlice()
	assertSliceEqual(t, []int{1, 2, 3, 4, 5}, result)
}

// --- Reduce: encadenamiento con Filter y Map ---

// Sumar solo los números pares después de filtrar.
func TestReduce_AfterFilter_SumsOnlyEvens(t *testing.T) {
	result := goseq.Reduce(
		goseq.From([]int{1, 2, 3, 4, 5, 6}).
			Filter(func(n int) bool { return n%2 == 0 }),
		0,
		func(acc, n int) int { return acc + n },
	)

	// 2 + 4 + 6 = 12
	if result != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}

// Sumar los valores ya transformados por Map.
func TestReduce_AfterMap_SumsDoubledValues(t *testing.T) {
	result := goseq.Reduce(
		goseq.Map(
			goseq.From([]int{1, 2, 3}),
			func(n int) int { return n * 2 },
		),
		0,
		func(acc, n int) int { return acc + n },
	)

	// (1*2) + (2*2) + (3*2) = 2 + 4 + 6 = 12
	if result != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}
