# Butaca — Sistema de Gestión de Streaming

Proyecto de la asignatura **Programación Orientada a Objetos** (Semestre 3, UIDE).
Etapa 1: planeación del software y esqueleto ejecutable del dominio.

Butaca gestiona el catálogo de una plataforma de streaming, las cuentas y perfiles
de sus suscriptores, y el progreso de reproducción de cada perfil. Está escrito en
**Go**, con el dominio modelado mediante estructuras, composición e interfaces, y
la lógica de consulta resuelta con funciones puras y de orden superior.

## Estado

Etapa 1 completada: dominio modelado, reglas de negocio implementadas y probadas,
y un comando de demostración que recorre los cuatro módulos del alcance. La
persistencia real (MySQL) y la interfaz de usuario corresponden a la Etapa 2.

| Métrica | Valor |
|---|---|
| Paquetes | 7 |
| Código de producción | 807 líneas |
| Código de pruebas | 484 líneas |
| Cobertura del dominio | 92,7 % |
| Dependencias de terceros | 0 en esta etapa |

## Cómo ejecutarlo

Requiere Go 1.24 o superior.

```bash
git clone https://github.com/IsRengel/butaca.git
cd butaca
go run ./cmd/butaca     # demostración del dominio
go test ./...           # suite de pruebas
make cobertura          # porcentaje de cobertura
```

## Estructura del repositorio

```
butaca/
├── cmd/butaca/            Punto de entrada: demostración ejecutable
├── internal/              Dominio. No importable desde fuera del módulo
│   ├── catalogo/          Contenidos, géneros, búsqueda por criterios
│   ├── usuarios/          Cuentas, perfiles y control parental
│   ├── suscripciones/     Planes comerciales y sus reglas
│   ├── reproduccion/      Progreso de visualización y "continuar viendo"
│   └── almacen/           Interfaz Repositorio e implementación en memoria
├── pkg/funcional/         Map, Filter, Reduce y utilidades genéricas
├── docs/                  Decisiones de arquitectura
└── Makefile               Tareas habituales
```

`internal/` y `pkg/` no son carpetas decorativas: Go impide que cualquier módulo
externo importe lo que está bajo `internal/`, de modo que el dominio queda cerrado
por el compilador, mientras que `pkg/funcional` sí es reutilizable por otros
proyectos.

## Módulos y alcance

| Módulo | Qué resuelve | Etapa 1 |
|---|---|---|
| Catálogo | Alta de películas y series, géneros, búsqueda por criterios combinables | Completo |
| Usuarios y perfiles | Cuentas, perfiles, límites por plan, control parental | Completo |
| Suscripciones | Planes, precios, pantallas simultáneas y calidad habilitada | Completo |
| Reproducción | Progreso por perfil, "continuar viendo", tiempo consumido | Completo |
| Persistencia | Contrato de repositorio e implementación en memoria | Completo (MySQL en Etapa 2) |
| Recomendaciones | Listas por género, tendencias, afinidad | Etapa 2 |
| Administración y reportes | Panel de métricas y altas masivas de contenido | Etapa 2 |

## Decisiones de diseño

- **Composición en lugar de herencia.** `Pelicula` y `Serie` embeben `Contenido`;
  ninguna hereda de nada, porque Go no tiene herencia.
- **Interfaces implícitas.** `Pelicula` y `Serie` satisfacen `Reproducible` por
  implementar sus métodos, sin declararlo. El polimorfismo aparece al guardar
  ambas en una misma porción de `Reproducible`.
- **Encapsulamiento por visibilidad.** Los campos en minúscula son privados al
  paquete; el estado solo cambia por constructores y métodos que validan.
- **Inmutabilidad deliberada.** Los métodos que "modifican" reciben el receptor
  por valor y devuelven una copia: `Agregar`, `Avanzar`, `CambiarPlan`,
  `ConSinopsis`. El valor anterior nunca se pierde.
- **Filtros como valores.** `Criterio` es un tipo función; `PorGenero`, `AptoPara`
  y `TituloContiene` son closures que se combinan con `Y` y `O`.
- **Errores como valores.** Nada entra en pánico: los constructores devuelven
  `error` y las comparaciones usan `errors.Is` contra errores centinela.
- **Inversión de dependencias.** El dominio depende de la interfaz `Repositorio`,
  nunca de un motor concreto. Cambiar memoria por MySQL no toca el dominio.

Los detalles están en [docs/arquitectura.md](docs/arquitectura.md).

## Dependencias

Esta etapa no usa ninguna biblioteca de terceros: todo se resuelve con la
biblioteca estándar (`sort`, `strings`, `errors`, `fmt`, `sync`, `context`,
`time`, `testing`). Las dependencias previstas para la Etapa 2 están documentadas
en el documento de planeación, con su licencia y su justificación.

## Licencia

MIT. Ver [LICENSE](LICENSE).
