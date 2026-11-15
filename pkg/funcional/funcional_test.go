package funcional_test

import (
	"testing"

	"github.com/IsRengel/butaca/pkg/funcional"
)

func TestMapear(t *testing.T) {
	numeros := []int{1, 2, 3}
	dobles := funcional.Mapear(numeros, func(n int) int { return n * 2 })

	if len(dobles) != 3 || dobles[2] != 6 {
		t.Fatalf("Mapear devolvió %v", dobles)
	}
	if numeros[2] != 3 {
		t.Fatal("Mapear mutó la porción original")
	}
}

func TestFiltrar(t *testing.T) {
	pares := funcional.Filtrar([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	if len(pares) != 2 {
		t.Fatalf("se esperaban 2 pares, se obtuvieron %d", len(pares))
	}
}

func TestReducir(t *testing.T) {
	suma := funcional.Reducir([]int{1, 2, 3, 4}, 0, func(acumulado, n int) int { return acumulado + n })
	if suma != 10 {
		t.Fatalf("se esperaba 10, se obtuvo %d", suma)
	}
}

func TestOrdenarPorNoMutaElOriginal(t *testing.T) {
	original := []int{3, 1, 2}
	ordenado := funcional.OrdenarPor(original, func(a, b int) bool { return a < b })

	if ordenado[0] != 1 {
		t.Fatalf("no quedó ordenado: %v", ordenado)
	}
	if original[0] != 3 {
		t.Fatal("OrdenarPor mutó la porción original")
	}
}

func TestAlgunoYTodos(t *testing.T) {
	numeros := []int{2, 4, 6}
	esPar := func(n int) bool { return n%2 == 0 }

	if !funcional.Todos(numeros, esPar) {
		t.Fatal("Todos debía ser verdadero")
	}
	if funcional.Alguno(numeros, func(n int) bool { return n > 10 }) {
		t.Fatal("Alguno debía ser falso")
	}
	if funcional.Todos([]int{}, esPar) != true {
		t.Fatal("Todos sobre una porción vacía es verdadero por definición")
	}
}

func TestCopiaEsIndependiente(t *testing.T) {
	original := []string{"a", "b"}
	copia := funcional.Copia(original)
	copia[0] = "z"

	if original[0] != "a" {
		t.Fatal("Copia no aisló la porción original")
	}
}
