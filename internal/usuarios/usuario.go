// Package usuarios modela las cuentas del sistema y los perfiles que cuelgan
// de cada una, con las reglas que impone el plan contratado.
package usuarios

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IsRengel/butaca/internal/suscripciones"
	"github.com/IsRengel/butaca/pkg/funcional"
)

var (
	ErrCorreoInvalido   = errors.New("el correo no tiene un formato válido")
	ErrPerfilDuplicado  = errors.New("ya existe un perfil con ese nombre")
	ErrLimiteDePerfiles = errors.New("el plan no admite más perfiles")
)

// Perfil es cada uno de los espacios independientes dentro de una cuenta.
type Perfil struct {
	ID       string
	Nombre   string
	Edad     int
	Infantil bool
}

// EsInfantil combina las dos condiciones que activan el modo restringido.
func (p Perfil) EsInfantil() bool { return p.Infantil || p.Edad < 13 }

// Usuario es la cuenta titular. Los campos son privados para que el plan y
// la lista de perfiles solo cambien por los métodos que validan las reglas.
type Usuario struct {
	id       string
	correo   string
	plan     suscripciones.Plan
	perfiles []Perfil
}

func NuevoUsuario(id, correo string, plan suscripciones.Plan) (Usuario, error) {
	if !strings.Contains(correo, "@") || !strings.Contains(correo, ".") {
		return Usuario{}, fmt.Errorf("%w: %q", ErrCorreoInvalido, correo)
	}
	return Usuario{id: id, correo: correo, plan: plan}, nil
}

func (u Usuario) ID() string               { return u.id }
func (u Usuario) Correo() string           { return u.correo }
func (u Usuario) Plan() suscripciones.Plan { return u.plan }
func (u Usuario) Perfiles() []Perfil       { return funcional.Copia(u.perfiles) }

// AgregarPerfil devuelve una cuenta nueva con el perfil incorporado, o un
// error si se rompe alguna regla del plan.
func (u Usuario) AgregarPerfil(p Perfil) (Usuario, error) {
	if len(u.perfiles) >= u.plan.PerfilesMaximos() {
		return u, fmt.Errorf("%w: %s admite %d", ErrLimiteDePerfiles, u.plan, u.plan.PerfilesMaximos())
	}
	duplicado := funcional.Alguno(u.perfiles, func(existente Perfil) bool {
		return strings.EqualFold(existente.Nombre, p.Nombre)
	})
	if duplicado {
		return u, fmt.Errorf("%w: %s", ErrPerfilDuplicado, p.Nombre)
	}
	u.perfiles = append(funcional.Copia(u.perfiles), p)
	return u, nil
}

// CambiarPlan devuelve una cuenta nueva con otro plan contratado.
func (u Usuario) CambiarPlan(nuevo suscripciones.Plan) Usuario {
	u.plan = nuevo
	return u
}

// PerfilesInfantiles proyecta los perfiles que deben ver catálogo restringido.
func (u Usuario) PerfilesInfantiles() []Perfil {
	return funcional.Filtrar(u.perfiles, func(p Perfil) bool { return p.EsInfantil() })
}
