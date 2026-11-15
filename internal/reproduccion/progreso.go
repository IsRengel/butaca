// Package reproduccion registra hasta dónde vio cada perfil cada contenido y
// deriva de ese historial la lista de "continuar viendo".
package reproduccion

import (
	"time"

	"github.com/IsRengel/butaca/pkg/funcional"
)

// Umbrales que definen cuándo un contenido se considera empezado y terminado.
const (
	UmbralIniciado  = 0.05
	UmbralTerminado = 0.95
)

// Progreso es un registro inmutable: avanzar no modifica el valor existente,
// devuelve uno nuevo. Así el historial puede guardarse como una secuencia de
// hechos en lugar de un dato que se sobrescribe.
type Progreso struct {
	PerfilID    string
	ContenidoID string
	Posicion    time.Duration
	Total       time.Duration
	Actualizado time.Time
}

// Avanzar devuelve el progreso resultante de ver d más, sin pasarse del total.
func (p Progreso) Avanzar(d time.Duration, reloj time.Time) Progreso {
	p.Posicion += d
	if p.Posicion > p.Total {
		p.Posicion = p.Total
	}
	p.Actualizado = reloj
	return p
}

// Porcentaje devuelve la fracción vista, entre 0 y 1.
func (p Progreso) Porcentaje() float64 {
	if p.Total <= 0 {
		return 0
	}
	return float64(p.Posicion) / float64(p.Total)
}

func (p Progreso) Iniciado() bool  { return p.Porcentaje() >= UmbralIniciado }
func (p Progreso) Terminado() bool { return p.Porcentaje() >= UmbralTerminado }

// ContinuarViendo filtra el historial de un perfil y devuelve lo que está
// empezado pero no terminado, ordenado por lo más reciente.
func ContinuarViendo(historial []Progreso, perfilID string) []Progreso {
	delPerfil := funcional.Filtrar(historial, func(p Progreso) bool {
		return p.PerfilID == perfilID && p.Iniciado() && !p.Terminado()
	})
	return funcional.OrdenarPor(delPerfil, func(a, b Progreso) bool {
		return a.Actualizado.After(b.Actualizado)
	})
}

// TiempoVisto suma cuánto consumió un perfil en total.
func TiempoVisto(historial []Progreso, perfilID string) time.Duration {
	delPerfil := funcional.Filtrar(historial, func(p Progreso) bool { return p.PerfilID == perfilID })
	return funcional.Reducir(delPerfil, time.Duration(0), func(total time.Duration, p Progreso) time.Duration {
		return total + p.Posicion
	})
}
