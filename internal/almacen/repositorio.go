// Package almacen define cómo persiste el sistema y aísla al dominio de esa
// decisión. El dominio depende de la interfaz Repositorio, nunca de MySQL ni
// de ningún otro motor concreto: es la inversión de dependencias que permite
// probar el sistema sin base de datos.
package almacen

import (
	"context"
	"errors"
	"sync"

	"github.com/IsRengel/butaca/pkg/funcional"
)

var ErrNoEncontrado = errors.New("no existe un registro con ese identificador")

// Repositorio es genérico: el mismo contrato sirve para contenidos, usuarios
// y progresos, sin duplicar una interfaz por entidad.
type Repositorio[T any] interface {
	Guardar(ctx context.Context, id string, item T) error
	Obtener(ctx context.Context, id string) (T, error)
	Listar(ctx context.Context) ([]T, error)
	Eliminar(ctx context.Context, id string) error
}

// EnMemoria es la implementación que se usa en las pruebas y en la demo de la
// Etapa 1. En la Etapa 2 se agrega una implementación equivalente sobre MySQL
// sin tocar una línea del dominio.
type EnMemoria[T any] struct {
	mu    sync.RWMutex
	datos map[string]T
	orden []string
}

func NuevoEnMemoria[T any]() *EnMemoria[T] {
	return &EnMemoria[T]{datos: make(map[string]T)}
}

func (r *EnMemoria[T]) Guardar(_ context.Context, id string, item T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, existe := r.datos[id]; !existe {
		r.orden = append(r.orden, id)
	}
	r.datos[id] = item
	return nil
}

func (r *EnMemoria[T]) Obtener(_ context.Context, id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, existe := r.datos[id]
	if !existe {
		var vacio T
		return vacio, ErrNoEncontrado
	}
	return item, nil
}

func (r *EnMemoria[T]) Listar(_ context.Context) ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return funcional.Mapear(r.orden, func(id string) T { return r.datos[id] }), nil
}

func (r *EnMemoria[T]) Eliminar(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, existe := r.datos[id]; !existe {
		return ErrNoEncontrado
	}
	delete(r.datos, id)
	r.orden = funcional.Filtrar(r.orden, func(existente string) bool { return existente != id })
	return nil
}
