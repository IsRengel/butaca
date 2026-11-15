// Package funcional reúne las operaciones de orden superior que usa el resto
// del sistema. Son funciones puras y genéricas: no mutan la porción que
// reciben, no dependen de estado externo y devuelven siempre una porción nueva.
package funcional

import "sort"

// Mapear transforma cada elemento de xs aplicando f.
func Mapear[T, R any](xs []T, f func(T) R) []R {
	resultado := make([]R, 0, len(xs))
	for _, x := range xs {
		resultado = append(resultado, f(x))
	}
	return resultado
}

// Filtrar devuelve los elementos de xs que satisfacen el predicado.
func Filtrar[T any](xs []T, predicado func(T) bool) []T {
	resultado := make([]T, 0, len(xs))
	for _, x := range xs {
		if predicado(x) {
			resultado = append(resultado, x)
		}
	}
	return resultado
}

// Reducir pliega xs sobre un acumulador, empezando por el valor inicial.
func Reducir[T, A any](xs []T, inicial A, f func(A, T) A) A {
	acumulador := inicial
	for _, x := range xs {
		acumulador = f(acumulador, x)
	}
	return acumulador
}

// Alguno indica si al menos un elemento satisface el predicado.
func Alguno[T any](xs []T, predicado func(T) bool) bool {
	for _, x := range xs {
		if predicado(x) {
			return true
		}
	}
	return false
}

// Todos indica si todos los elementos satisfacen el predicado.
func Todos[T any](xs []T, predicado func(T) bool) bool {
	for _, x := range xs {
		if !predicado(x) {
			return false
		}
	}
	return true
}

// OrdenarPor devuelve una copia ordenada de xs según el criterio de orden.
// No altera la porción original, a diferencia de sort.Slice.
func OrdenarPor[T any](xs []T, menor func(a, b T) bool) []T {
	copia := make([]T, len(xs))
	copy(copia, xs)
	sort.SliceStable(copia, func(i, j int) bool { return menor(copia[i], copia[j]) })
	return copia
}

// Copia devuelve un duplicado de la porción recibida. Se usa en los métodos de
// acceso del dominio para no exponer la porción interna de una estructura.
func Copia[T any](xs []T) []T {
	copia := make([]T, len(xs))
	copy(copia, xs)
	return copia
}
