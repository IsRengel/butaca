package catalogo

import (
	"strings"
	"time"

	"github.com/IsRengel/butaca/pkg/funcional"
)

// Criterio es un predicado sobre un contenido. Al ser un tipo función, los
// filtros se componen como valores: se guardan en variables, se pasan como
// argumento y se combinan entre sí.
type Criterio func(Reproducible) bool

// PorGenero devuelve un criterio que acepta los contenidos de ese género.
// Es un closure: la función devuelta captura el valor de genero.
func PorGenero(genero Genero) Criterio {
	return func(r Reproducible) bool {
		return funcional.Alguno(r.Generos(), func(g Genero) bool { return g == genero })
	}
}

// AptoPara devuelve un criterio de control parental.
func AptoPara(edad int) Criterio {
	return func(r Reproducible) bool { return r.EsAptoPara(edad) }
}

// TituloContiene busca por coincidencia parcial, sin distinguir mayúsculas.
func TituloContiene(fragmento string) Criterio {
	buscado := strings.ToLower(fragmento)
	return func(r Reproducible) bool {
		return strings.Contains(strings.ToLower(r.Titulo()), buscado)
	}
}

// MasCortaQue filtra por duración, útil para la sección "para ver en una hora".
func MasCortaQue(limite time.Duration) Criterio {
	return func(r Reproducible) bool { return r.Duracion() < limite }
}

// Y combina varios criterios en uno solo que exige que se cumplan todos.
func Y(criterios ...Criterio) Criterio {
	return func(r Reproducible) bool {
		return funcional.Todos(criterios, func(c Criterio) bool { return c(r) })
	}
}

// O combina varios criterios en uno que se cumple si alguno se cumple.
func O(criterios ...Criterio) Criterio {
	return func(r Reproducible) bool {
		return funcional.Alguno(criterios, func(c Criterio) bool { return c(r) })
	}
}

// Catalogo es la colección de contenidos publicados. Todas sus operaciones
// devuelven un catálogo nuevo en lugar de modificar el receptor.
type Catalogo struct {
	items []Reproducible
}

func Nuevo(items ...Reproducible) Catalogo {
	return Catalogo{items: funcional.Copia(items)}
}

// Agregar devuelve un catálogo nuevo con el contenido incorporado.
func (c Catalogo) Agregar(items ...Reproducible) Catalogo {
	return Catalogo{items: append(funcional.Copia(c.items), items...)}
}

func (c Catalogo) Cantidad() int { return len(c.items) }

func (c Catalogo) Todos() []Reproducible { return funcional.Copia(c.items) }

// Buscar aplica los criterios recibidos y devuelve los contenidos que los
// satisfacen todos.
func (c Catalogo) Buscar(criterios ...Criterio) []Reproducible {
	return funcional.Filtrar(c.items, Y(criterios...))
}

// Titulos proyecta el catálogo a una lista de títulos ordenada alfabéticamente.
func (c Catalogo) Titulos() []string {
	ordenados := funcional.OrdenarPor(c.items, func(a, b Reproducible) bool {
		return a.Titulo() < b.Titulo()
	})
	return funcional.Mapear(ordenados, func(r Reproducible) string { return r.Titulo() })
}

// DuracionTotal suma la duración de todo el catálogo.
func (c Catalogo) DuracionTotal() time.Duration {
	return funcional.Reducir(c.items, time.Duration(0), func(total time.Duration, r Reproducible) time.Duration {
		return total + r.Duracion()
	})
}
