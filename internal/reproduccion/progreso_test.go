package reproduccion_test

import (
	"testing"
	"time"

	"github.com/IsRengel/butaca/internal/reproduccion"
)

func TestAvanzarNoMutaElOriginal(t *testing.T) {
	ahora := time.Now()
	inicial := reproduccion.Progreso{PerfilID: "p1", ContenidoID: "c1", Total: time.Hour}

	avanzado := inicial.Avanzar(30*time.Minute, ahora)

	if inicial.Posicion != 0 {
		t.Fatal("Avanzar modificó el progreso original")
	}
	if avanzado.Posicion != 30*time.Minute {
		t.Fatalf("posición incorrecta: %v", avanzado.Posicion)
	}
}

func TestAvanzarNoSePasaDelTotal(t *testing.T) {
	p := reproduccion.Progreso{Total: time.Hour}.Avanzar(90*time.Minute, time.Now())
	if p.Posicion != time.Hour {
		t.Fatalf("la posición superó la duración: %v", p.Posicion)
	}
}

func TestContinuarViendoFiltraYOrdena(t *testing.T) {
	ahora := time.Now()
	historial := []reproduccion.Progreso{
		{PerfilID: "p1", ContenidoID: "viejo", Total: time.Hour, Posicion: 30 * time.Minute, Actualizado: ahora.Add(-72 * time.Hour)},
		{PerfilID: "p1", ContenidoID: "reciente", Total: time.Hour, Posicion: 20 * time.Minute, Actualizado: ahora},
		{PerfilID: "p1", ContenidoID: "terminado", Total: time.Hour, Posicion: 59 * time.Minute, Actualizado: ahora},
		{PerfilID: "p1", ContenidoID: "apenas", Total: time.Hour, Posicion: time.Minute, Actualizado: ahora},
		{PerfilID: "p2", ContenidoID: "de otro perfil", Total: time.Hour, Posicion: 30 * time.Minute, Actualizado: ahora},
	}

	lista := reproduccion.ContinuarViendo(historial, "p1")

	if len(lista) != 2 {
		t.Fatalf("se esperaban 2 contenidos por continuar, se obtuvieron %d", len(lista))
	}
	if lista[0].ContenidoID != "reciente" {
		t.Fatalf("el orden debía ser por fecha descendente, llegó %q primero", lista[0].ContenidoID)
	}
}

func TestTiempoVistoSoloCuentaElPerfilPedido(t *testing.T) {
	historial := []reproduccion.Progreso{
		{PerfilID: "p1", Posicion: 30 * time.Minute, Total: time.Hour},
		{PerfilID: "p1", Posicion: 15 * time.Minute, Total: time.Hour},
		{PerfilID: "p2", Posicion: 45 * time.Minute, Total: time.Hour},
	}

	if visto := reproduccion.TiempoVisto(historial, "p1"); visto != 45*time.Minute {
		t.Fatalf("se esperaban 45 min, se obtuvieron %v", visto)
	}
}
