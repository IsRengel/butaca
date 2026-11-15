package almacen_test

import (
	"context"
	"errors"
	"testing"

	"github.com/IsRengel/butaca/internal/almacen"
)

type contenidoDePrueba struct {
	ID     string
	Titulo string
}

func TestGuardarYObtener(t *testing.T) {
	ctx := context.Background()
	repo := almacen.NuevoEnMemoria[contenidoDePrueba]()

	if err := repo.Guardar(ctx, "c1", contenidoDePrueba{ID: "c1", Titulo: "Uno"}); err != nil {
		t.Fatal(err)
	}

	item, err := repo.Obtener(ctx, "c1")
	if err != nil {
		t.Fatal(err)
	}
	if item.Titulo != "Uno" {
		t.Fatalf("se recuperó %q", item.Titulo)
	}
}

func TestObtenerInexistente(t *testing.T) {
	repo := almacen.NuevoEnMemoria[contenidoDePrueba]()
	if _, err := repo.Obtener(context.Background(), "nada"); !errors.Is(err, almacen.ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestListarConservaElOrdenDeAlta(t *testing.T) {
	ctx := context.Background()
	repo := almacen.NuevoEnMemoria[contenidoDePrueba]()
	for _, id := range []string{"c3", "c1", "c2"} {
		if err := repo.Guardar(ctx, id, contenidoDePrueba{ID: id}); err != nil {
			t.Fatal(err)
		}
	}

	items, err := repo.Listar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[0].ID != "c3" {
		t.Fatalf("el orden de alta no se conservó: %v", items)
	}
}

func TestEliminar(t *testing.T) {
	ctx := context.Background()
	repo := almacen.NuevoEnMemoria[contenidoDePrueba]()
	_ = repo.Guardar(ctx, "c1", contenidoDePrueba{ID: "c1"})

	if err := repo.Eliminar(ctx, "c1"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Eliminar(ctx, "c1"); !errors.Is(err, almacen.ErrNoEncontrado) {
		t.Fatalf("eliminar dos veces debía fallar, devolvió %v", err)
	}
}

// El repositorio en memoria debe satisfacer la interfaz genérica. Si deja de
// hacerlo, este archivo no compila: es una verificación en tiempo de
// compilación, no en tiempo de ejecución.
var _ almacen.Repositorio[contenidoDePrueba] = (*almacen.EnMemoria[contenidoDePrueba])(nil)
