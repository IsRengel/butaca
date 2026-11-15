package usuarios_test

import (
	"errors"
	"testing"

	"github.com/IsRengel/butaca/internal/suscripciones"
	"github.com/IsRengel/butaca/internal/usuarios"
)

func TestCorreoInvalido(t *testing.T) {
	if _, err := usuarios.NuevoUsuario("u1", "josue-ejemplo", suscripciones.Basico); !errors.Is(err, usuarios.ErrCorreoInvalido) {
		t.Fatalf("se esperaba ErrCorreoInvalido, se obtuvo %v", err)
	}
}

func TestLimiteDePerfilesSegunPlan(t *testing.T) {
	cuenta, err := usuarios.NuevoUsuario("u1", "josue@ejemplo.ec", suscripciones.Basico)
	if err != nil {
		t.Fatal(err)
	}

	cuenta, _ = cuenta.AgregarPerfil(usuarios.Perfil{ID: "p1", Nombre: "Uno", Edad: 20})
	cuenta, _ = cuenta.AgregarPerfil(usuarios.Perfil{ID: "p2", Nombre: "Dos", Edad: 20})

	if _, err := cuenta.AgregarPerfil(usuarios.Perfil{ID: "p3", Nombre: "Tres", Edad: 20}); !errors.Is(err, usuarios.ErrLimiteDePerfiles) {
		t.Fatalf("el plan Básico debía rechazar el tercer perfil, devolvió %v", err)
	}
}

func TestPerfilDuplicadoSinDistinguirMayusculas(t *testing.T) {
	cuenta, _ := usuarios.NuevoUsuario("u1", "josue@ejemplo.ec", suscripciones.Premium)
	cuenta, _ = cuenta.AgregarPerfil(usuarios.Perfil{ID: "p1", Nombre: "Josue", Edad: 24})

	if _, err := cuenta.AgregarPerfil(usuarios.Perfil{ID: "p2", Nombre: "josue", Edad: 24}); !errors.Is(err, usuarios.ErrPerfilDuplicado) {
		t.Fatalf("se esperaba ErrPerfilDuplicado, se obtuvo %v", err)
	}
}

func TestPerfilInfantilPorEdad(t *testing.T) {
	cuenta, _ := usuarios.NuevoUsuario("u1", "josue@ejemplo.ec", suscripciones.Estandar)
	cuenta, _ = cuenta.AgregarPerfil(usuarios.Perfil{ID: "p1", Nombre: "Adulto", Edad: 30})
	cuenta, _ = cuenta.AgregarPerfil(usuarios.Perfil{ID: "p2", Nombre: "Niña", Edad: 7})

	if len(cuenta.PerfilesInfantiles()) != 1 {
		t.Fatalf("se esperaba 1 perfil infantil, hay %d", len(cuenta.PerfilesInfantiles()))
	}
}

func TestCambiarPlanDevuelveCuentaNueva(t *testing.T) {
	original, _ := usuarios.NuevoUsuario("u1", "josue@ejemplo.ec", suscripciones.Basico)
	mejorada := original.CambiarPlan(suscripciones.Premium)

	if original.Plan() != suscripciones.Basico {
		t.Fatal("CambiarPlan mutó la cuenta original")
	}
	if mejorada.Plan() != suscripciones.Premium {
		t.Fatalf("la cuenta nueva debía tener plan Premium, tiene %s", mejorada.Plan())
	}
}
