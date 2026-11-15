// Package catalogo modela los contenidos que el sistema publica y las
// operaciones de búsqueda sobre ellos.
//
// Go no tiene clases ni herencia: la orientación a objetos se expresa con
// estructuras, composición (embedding) e interfaces. Pelicula y Serie no
// heredan de Contenido, lo embeben, y ambas satisfacen Reproducible porque
// implementan sus métodos, sin declararlo en ninguna parte.
package catalogo

import (
	"errors"
	"fmt"
	"time"

	"github.com/IsRengel/butaca/pkg/funcional"
)

// Genero clasifica un contenido. Es un tipo propio y no una cadena suelta para
// que el compilador impida pasar cualquier texto donde se espera un género.
type Genero string

const (
	GeneroAccion     Genero = "acción"
	GeneroDrama      Genero = "drama"
	GeneroComedia    Genero = "comedia"
	GeneroDocumental Genero = "documental"
	GeneroInfantil   Genero = "infantil"
	GeneroSuspenso   Genero = "suspenso"
)

var (
	ErrTituloVacio  = errors.New("el título no puede estar vacío")
	ErrAnioInvalido = errors.New("el año debe estar entre 1900 y 2100")
	ErrEdadInvalida = errors.New("la edad mínima debe estar entre 0 y 21")
)

// Reproducible es el contrato que cumple todo lo que un perfil puede ver.
// Es la pieza que da polimorfismo al catálogo: una porción de Reproducible
// puede contener películas y series mezcladas.
type Reproducible interface {
	ID() string
	Titulo() string
	Duracion() time.Duration
	EsAptoPara(edad int) bool
	Generos() []Genero
}

// Contenido agrupa los atributos comunes a películas y series. Sus campos son
// privados: solo se construye a través de NuevoContenido, que valida, y solo se
// lee a través de métodos, de modo que ninguna instancia puede quedar en un
// estado inválido.
type Contenido struct {
	id         string
	titulo     string
	sinopsis   string
	anio       int
	generos    []Genero
	edadMinima int
}

// NuevoContenido es el constructor. Devuelve un error en lugar de entrar en
// pánico, siguiendo la convención de Go de tratar los errores como valores.
func NuevoContenido(id, titulo string, anio, edadMinima int, generos ...Genero) (Contenido, error) {
	if titulo == "" {
		return Contenido{}, ErrTituloVacio
	}
	if anio < 1900 || anio > 2100 {
		return Contenido{}, fmt.Errorf("%w: %d", ErrAnioInvalido, anio)
	}
	if edadMinima < 0 || edadMinima > 21 {
		return Contenido{}, fmt.Errorf("%w: %d", ErrEdadInvalida, edadMinima)
	}
	return Contenido{
		id:         id,
		titulo:     titulo,
		anio:       anio,
		edadMinima: edadMinima,
		generos:    funcional.Copia(generos),
	}, nil
}

func (c Contenido) ID() string       { return c.id }
func (c Contenido) Titulo() string   { return c.titulo }
func (c Contenido) Sinopsis() string { return c.sinopsis }
func (c Contenido) Anio() int        { return c.anio }
func (c Contenido) EdadMinima() int  { return c.edadMinima }

// Generos devuelve una copia: si devolviera la porción interna, quien la
// recibe podría modificar el contenido desde fuera del paquete.
func (c Contenido) Generos() []Genero { return funcional.Copia(c.generos) }

// EsAptoPara aplica la regla de control parental.
func (c Contenido) EsAptoPara(edad int) bool { return edad >= c.edadMinima }

// ConSinopsis devuelve una copia con la sinopsis cambiada. El receptor es por
// valor y no por puntero: el original queda intacto. Es la forma de mantener
// inmutabilidad en un lenguaje que no la impone.
func (c Contenido) ConSinopsis(sinopsis string) Contenido {
	c.sinopsis = sinopsis
	return c
}

// Pelicula compone Contenido y le agrega su duración.
type Pelicula struct {
	Contenido
	duracion time.Duration
}

func NuevaPelicula(base Contenido, duracion time.Duration) Pelicula {
	return Pelicula{Contenido: base, duracion: duracion}
}

func (p Pelicula) Duracion() time.Duration { return p.duracion }

// Episodio es la unidad reproducible de una serie.
type Episodio struct {
	Numero   int
	Titulo   string
	Duracion time.Duration
}

// Temporada agrupa episodios.
type Temporada struct {
	Numero    int
	Episodios []Episodio
}

// Serie compone Contenido y agrega sus temporadas.
type Serie struct {
	Contenido
	temporadas []Temporada
}

func NuevaSerie(base Contenido, temporadas ...Temporada) Serie {
	return Serie{Contenido: base, temporadas: funcional.Copia(temporadas)}
}

func (s Serie) Temporadas() []Temporada { return funcional.Copia(s.temporadas) }

// Duracion suma la de todos los episodios de todas las temporadas usando
// Reducir, en lugar de dos bucles anidados con un acumulador mutable.
func (s Serie) Duracion() time.Duration {
	return funcional.Reducir(s.temporadas, time.Duration(0), func(total time.Duration, t Temporada) time.Duration {
		return total + funcional.Reducir(t.Episodios, time.Duration(0), func(parcial time.Duration, e Episodio) time.Duration {
			return parcial + e.Duracion
		})
	})
}

// CantidadEpisodios informa cuántos episodios tiene la serie en total.
func (s Serie) CantidadEpisodios() int {
	return funcional.Reducir(s.temporadas, 0, func(total int, t Temporada) int {
		return total + len(t.Episodios)
	})
}
