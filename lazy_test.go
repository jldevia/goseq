package goseq_test

import (
	"testing"

	"github.com/jldevia/goseq"
)

// =====================
// Pruebas de lazyness
// Estos tests verifican que la evaluación diferida funciona realmente:
// que las operaciones no se ejecutan hasta que se consume la secuencia,
// y que Take(n) no procesa más elementos de los necesarios.
// =====================

// Take(n) sobre una secuencia grande solo debe procesar n elementos.
// Lo verificamos contando cuántas veces se invoca la función de transformación.
func TestLazy_TakeOnlyProcessesNElements(t *testing.T) {
	processed := 0

	// Creamos una secuencia de 1000 elementos
	input := make([]int, 1000)
	for i := range input {
		input[i] = i + 1
	}

	result := goseq.Map(
		goseq.From(input),
		func(n int) int {
			processed++
			return n * 2
		},
	).Take(5).ToSlice()

	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
	}

	// Con lazy evaluation solo se procesan 5 elementos, no 1000
	if processed != 5 {
		t.Errorf("expected 5 elements processed (lazy), but processed %d", processed)
	}
}

// Any debe detenerse en el primer match — no recorrer toda la secuencia.
func TestLazy_AnyShortCircuits(t *testing.T) {
	processed := 0

	input := make([]int, 1000)
	for i := range input {
		input[i] = i + 1
	}

	found := goseq.From(input).
		Filter(func(n int) bool {
			processed++
			return n == 3
		}).
		Any(func(n int) bool { return true })

	if !found {
		t.Error("expected true")
	}

	// Solo debe haber procesado hasta el elemento 3
	if processed != 3 {
		t.Errorf("expected 3 elements processed (short-circuit), but processed %d", processed)
	}
}

// All debe detenerse en el primer elemento que no cumple el predicado.
func TestLazy_AllShortCircuits(t *testing.T) {
	processed := 0

	result := goseq.From([]int{2, 4, 5, 8, 10}).
		All(func(n int) bool {
			processed++
			return n%2 == 0
		})

	if result {
		t.Error("expected false — 5 is not even")
	}

	// Debe detenerse en el tercer elemento (5), no recorrer los 5
	if processed != 3 {
		t.Errorf("expected 3 elements processed (short-circuit), but processed %d", processed)
	}
}

// Una misma Seq puede consumirse múltiples veces independientemente.
func TestLazy_SeqIsReusable(t *testing.T) {
	s := goseq.From([]int{1, 2, 3, 4, 5})

	first := s.Take(3).ToSlice()
	second := s.ToSlice()
	third := s.Filter(func(n int) bool { return n%2 == 0 }).ToSlice()

	assertSliceEqual(t, []int{1, 2, 3}, first)
	assertSliceEqual(t, []int{1, 2, 3, 4, 5}, second)
	assertSliceEqual(t, []int{2, 4}, third)
}

// Filter + Take solo procesa lo necesario para obtener n elementos filtrados.
func TestLazy_FilterTakeProcessesMinimum(t *testing.T) {
	processed := 0

	input := make([]int, 1000)
	for i := range input {
		input[i] = i + 1
	}

	result := goseq.From(input).
		Filter(func(n int) bool {
			processed++
			return n%2 == 0
		}).
		Take(3).
		ToSlice()

	assertSliceEqual(t, []int{2, 4, 6}, result)

	// Para obtener 3 pares (2, 4, 6) solo necesita procesar hasta el 6 → 6 elementos
	if processed != 6 {
		t.Errorf("expected 6 elements processed, but processed %d", processed)
	}
}

// Las operaciones intermedias no tienen efectos secundarios hasta ToSlice.
func TestLazy_NoSideEffectsUntilTerminal(t *testing.T) {
	executed := false

	// Construir el pipeline completo sin ejecutarlo
	pipeline := goseq.From([]int{1, 2, 3}).
		Filter(func(n int) bool {
			executed = true
			return n > 1
		})

	// Nada debe haberse ejecutado todavía
	if executed {
		t.Error("filter should not have executed before terminal operation")
	}

	// Ahora sí ejecutamos
	pipeline.ToSlice()

	if !executed {
		t.Error("filter should have executed after ToSlice")
	}
}

// TakeWhile lazy: se detiene sin procesar el resto de la secuencia.
func TestLazy_TakeWhileShortCircuits(t *testing.T) {
	processed := 0

	result := goseq.From([]int{1, 2, 3, 10, 11, 12}).
		TakeWhile(func(n int) bool {
			processed++
			return n < 5
		}).
		ToSlice()

	assertSliceEqual(t, []int{1, 2, 3}, result)

	// Procesa 1, 2, 3 (pasan) y 10 (falla) → 4 evaluaciones, no 6
	if processed != 4 {
		t.Errorf("expected 4 elements processed, but processed %d", processed)
	}
}
