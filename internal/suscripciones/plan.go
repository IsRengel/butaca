// Package suscripciones modela los planes comerciales y las reglas que de
// ellos se derivan, como cuántas pantallas simultáneas admite una cuenta.
package suscripciones

import "fmt"

// Plan es un tipo enumerado. Al implementar Stringer, cualquier función que
// imprima un plan obtiene el nombre legible sin conversiones manuales.
type Plan int

const (
	Basico Plan = iota
	Estandar
	Premium
)

func (p Plan) String() string {
	switch p {
	case Basico:
		return "Básico"
	case Estandar:
		return "Estándar"
	case Premium:
		return "Premium"
	default:
		return fmt.Sprintf("Plan(%d)", int(p))
	}
}

// PantallasSimultaneas es la regla de negocio que distingue a los planes.
func (p Plan) PantallasSimultaneas() int {
	switch p {
	case Basico:
		return 1
	case Estandar:
		return 2
	case Premium:
		return 4
	default:
		return 0
	}
}

// PerfilesMaximos limita cuántos perfiles admite la cuenta.
func (p Plan) PerfilesMaximos() int {
	if p == Basico {
		return 2
	}
	return 5
}

// PrecioMensual devuelve el precio en dólares.
func (p Plan) PrecioMensual() float64 {
	switch p {
	case Basico:
		return 5.99
	case Estandar:
		return 9.99
	case Premium:
		return 14.99
	default:
		return 0
	}
}

// PermiteAltaDefinicion indica si el plan habilita calidad HD o superior.
func (p Plan) PermiteAltaDefinicion() bool { return p != Basico }
