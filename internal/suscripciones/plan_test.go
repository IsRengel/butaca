package suscripciones_test

import (
	"testing"

	"github.com/IsRengel/butaca/internal/suscripciones"
)

func TestReglasDeCadaPlan(t *testing.T) {
	casos := []struct {
		plan       suscripciones.Plan
		nombre     string
		pantallas  int
		perfiles   int
		precio     float64
		altaDefini bool
	}{
		{suscripciones.Basico, "Básico", 1, 2, 5.99, false},
		{suscripciones.Estandar, "Estándar", 2, 5, 9.99, true},
		{suscripciones.Premium, "Premium", 4, 5, 14.99, true},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if caso.plan.String() != caso.nombre {
				t.Errorf("nombre: se esperaba %q, se obtuvo %q", caso.nombre, caso.plan)
			}
			if caso.plan.PantallasSimultaneas() != caso.pantallas {
				t.Errorf("pantallas: se esperaban %d, se obtuvieron %d", caso.pantallas, caso.plan.PantallasSimultaneas())
			}
			if caso.plan.PerfilesMaximos() != caso.perfiles {
				t.Errorf("perfiles: se esperaban %d, se obtuvieron %d", caso.perfiles, caso.plan.PerfilesMaximos())
			}
			if caso.plan.PrecioMensual() != caso.precio {
				t.Errorf("precio: se esperaba %.2f, se obtuvo %.2f", caso.precio, caso.plan.PrecioMensual())
			}
			if caso.plan.PermiteAltaDefinicion() != caso.altaDefini {
				t.Errorf("alta definición: se esperaba %t", caso.altaDefini)
			}
		})
	}
}

func TestPlanDesconocido(t *testing.T) {
	desconocido := suscripciones.Plan(99)
	if desconocido.PantallasSimultaneas() != 0 || desconocido.PrecioMensual() != 0 {
		t.Fatal("un plan inexistente no debe habilitar nada")
	}
	if desconocido.String() != "Plan(99)" {
		t.Fatalf("representación inesperada: %s", desconocido)
	}
}
