package catalogo_test

import (
	"errors"
	"testing"
	"time"

	"github.com/IsRengel/butaca/internal/catalogo"
)

func pelicula(t *testing.T, id, titulo string, edad int, dur time.Duration, generos ...catalogo.Genero) catalogo.Pelicula {
	t.Helper()
	base, err := catalogo.NuevoContenido(id, titulo, 2024, edad, generos...)
	if err != nil {
		t.Fatalf("no se pudo construir el contenido: %v", err)
	}
	return catalogo.NuevaPelicula(base, dur)
}

func TestNuevoContenidoValida(t *testing.T) {
	casos := []struct {
		nombre   string
		titulo   string
		anio     int
		edad     int
		esperado error
	}{
		{"título vacío", "", 2024, 12, catalogo.ErrTituloVacio},
		{"año fuera de rango", "Algo", 1750, 12, catalogo.ErrAnioInvalido},
		{"edad fuera de rango", "Algo", 2024, 99, catalogo.ErrEdadInvalida},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, err := catalogo.NuevoContenido("x", caso.titulo, caso.anio, caso.edad)
			if !errors.Is(err, caso.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", caso.esperado, err)
			}
		})
	}
}

func TestCatalogoEsInmutable(t *testing.T) {
	original := catalogo.Nuevo(pelicula(t, "c1", "Uno", 0, time.Hour))
	ampliado := original.Agregar(pelicula(t, "c2", "Dos", 0, time.Hour))

	if original.Cantidad() != 1 {
		t.Fatalf("Agregar mutó el catálogo original: %d", original.Cantidad())
	}
	if ampliado.Cantidad() != 2 {
		t.Fatalf("el catálogo ampliado debía tener 2, tiene %d", ampliado.Cantidad())
	}
}

func TestBuscarCombinaCriterios(t *testing.T) {
	cat := catalogo.Nuevo(
		pelicula(t, "c1", "Acción adulta", 18, 2*time.Hour, catalogo.GeneroAccion),
		pelicula(t, "c2", "Acción familiar", 7, 90*time.Minute, catalogo.GeneroAccion),
		pelicula(t, "c3", "Dibujos", 0, 45*time.Minute, catalogo.GeneroInfantil),
	)

	resultados := cat.Buscar(catalogo.PorGenero(catalogo.GeneroAccion), catalogo.AptoPara(10))
	if len(resultados) != 1 || resultados[0].ID() != "c2" {
		t.Fatalf("la combinación de criterios falló: %v", resultados)
	}
}

func TestCriterioOEsDisyuntivo(t *testing.T) {
	cat := catalogo.Nuevo(
		pelicula(t, "c1", "Drama", 12, time.Hour, catalogo.GeneroDrama),
		pelicula(t, "c2", "Comedia", 12, time.Hour, catalogo.GeneroComedia),
		pelicula(t, "c3", "Documental", 12, time.Hour, catalogo.GeneroDocumental),
	)

	criterio := catalogo.O(catalogo.PorGenero(catalogo.GeneroDrama), catalogo.PorGenero(catalogo.GeneroComedia))
	if resultados := cat.Buscar(criterio); len(resultados) != 2 {
		t.Fatalf("se esperaban 2 resultados, se obtuvieron %d", len(resultados))
	}
}

func TestSerieSumaLaDuracionDeSusEpisodios(t *testing.T) {
	base, err := catalogo.NuevoContenido("s1", "Serie", 2025, 12, catalogo.GeneroDrama)
	if err != nil {
		t.Fatal(err)
	}
	serie := catalogo.NuevaSerie(base,
		catalogo.Temporada{Numero: 1, Episodios: []catalogo.Episodio{
			{Numero: 1, Duracion: 30 * time.Minute},
			{Numero: 2, Duracion: 30 * time.Minute},
		}},
		catalogo.Temporada{Numero: 2, Episodios: []catalogo.Episodio{
			{Numero: 1, Duracion: 40 * time.Minute},
		}},
	)

	if serie.Duracion() != 100*time.Minute {
		t.Fatalf("duración incorrecta: %v", serie.Duracion())
	}
	if serie.CantidadEpisodios() != 3 {
		t.Fatalf("se esperaban 3 episodios, hay %d", serie.CantidadEpisodios())
	}
}

// TestPolimorfismo comprueba que películas y series conviven en una misma
// porción de Reproducible, que es lo que en Go reemplaza a la herencia.
func TestPolimorfismo(t *testing.T) {
	base, err := catalogo.NuevoContenido("s1", "Serie", 2025, 0)
	if err != nil {
		t.Fatal(err)
	}
	var contenidos []catalogo.Reproducible
	contenidos = append(contenidos,
		pelicula(t, "c1", "Película", 0, time.Hour),
		catalogo.NuevaSerie(base, catalogo.Temporada{Numero: 1, Episodios: []catalogo.Episodio{{Duracion: time.Hour}}}),
	)

	cat := catalogo.Nuevo(contenidos...)
	if cat.DuracionTotal() != 2*time.Hour {
		t.Fatalf("duración total incorrecta: %v", cat.DuracionTotal())
	}
}

func TestGenerosNoSeFiltranDesdeFuera(t *testing.T) {
	p := pelicula(t, "c1", "Uno", 0, time.Hour, catalogo.GeneroAccion)
	generos := p.Generos()
	generos[0] = catalogo.GeneroInfantil

	if p.Generos()[0] != catalogo.GeneroAccion {
		t.Fatal("Generos() expuso la porción interna de la estructura")
	}
}

func TestTituloContieneIgnoraMayusculas(t *testing.T) {
	cat := catalogo.Nuevo(pelicula(t, "c1", "La Última Frontera", 0, time.Hour))
	if len(cat.Buscar(catalogo.TituloContiene("última"))) != 1 {
		t.Fatal("la búsqueda debía ignorar mayúsculas")
	}
	if len(cat.Buscar(catalogo.TituloContiene("inexistente"))) != 0 {
		t.Fatal("no debía encontrar nada")
	}
}

func TestConSinopsisDevuelveUnaCopia(t *testing.T) {
	base, err := catalogo.NuevoContenido("c1", "Uno", 2024, 0)
	if err != nil {
		t.Fatal(err)
	}
	conTexto := base.ConSinopsis("Un relato breve")

	if base.Sinopsis() != "" {
		t.Fatal("ConSinopsis mutó el contenido original")
	}
	if conTexto.Sinopsis() != "Un relato breve" {
		t.Fatalf("sinopsis incorrecta: %q", conTexto.Sinopsis())
	}
}

func TestTitulosDevuelveOrdenAlfabetico(t *testing.T) {
	cat := catalogo.Nuevo(
		pelicula(t, "c1", "Zapatos", 0, time.Hour),
		pelicula(t, "c2", "Abanico", 0, time.Hour),
	)
	titulos := cat.Titulos()
	if titulos[0] != "Abanico" || titulos[1] != "Zapatos" {
		t.Fatalf("orden incorrecto: %v", titulos)
	}
}
