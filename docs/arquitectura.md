# Decisiones de arquitectura

## 1. Por qué Go para una materia de orientación a objetos

Go no tiene clases, ni herencia, ni constructores, ni excepciones. Eso no lo deja
fuera del paradigma: lo obliga a expresarlo con los mecanismos que sí tiene.

| Concepto de POO | Cómo se expresa en Go | Dónde se ve en este repositorio |
|---|---|---|
| Clase | `struct` más sus métodos | `catalogo.Contenido`, `usuarios.Usuario` |
| Objeto | Valor de ese tipo | `catalogo.NuevaPelicula(...)` |
| Encapsulamiento | Identificadores en minúscula, privados al paquete | Campos de `Contenido` y accesores |
| Constructor | Función `NuevoX` que valida y devuelve `(T, error)` | `catalogo.NuevoContenido` |
| Herencia | **No existe.** Se reemplaza por composición (embedding) | `Pelicula` embebe `Contenido` |
| Interfaz | `interface` satisfecha implícitamente | `catalogo.Reproducible` |
| Polimorfismo | Porción de un tipo interfaz con implementaciones distintas | `Catalogo.items []Reproducible` |
| Genéricos | Parámetros de tipo desde Go 1.18 | `almacen.Repositorio[T]`, `funcional.Mapear[T, R]` |

La ausencia de herencia no es una carencia. Evita el acoplamiento de las
jerarquías profundas: si mañana aparece un tipo de contenido nuevo —un evento en
vivo, por ejemplo— basta con que implemente `Reproducible`, sin tocar ni
`Pelicula` ni `Serie` ni el catálogo.

## 2. Dónde entra la programación funcional

El enfoque funcional no reemplaza al orientado a objetos, lo complementa: los
objetos modelan el dominio y las funciones resuelven las consultas sobre él.

- **Funciones puras.** Todo `pkg/funcional` depende solo de sus argumentos y no
  tiene efectos secundarios. Dos llamadas con la misma entrada dan la misma
  salida, lo que las hace triviales de probar.
- **Funciones de orden superior.** `Mapear`, `Filtrar` y `Reducir` reciben
  funciones como argumento. `Catalogo.DuracionTotal` es un `Reducir`;
  `Catalogo.Titulos` es un `OrdenarPor` seguido de un `Mapear`.
- **Closures.** `PorGenero(GeneroInfantil)` devuelve una función que captura el
  género. El filtro deja de ser un `if` dentro de un bucle y pasa a ser un valor
  que se guarda, se pasa y se combina.
- **Composición.** `Y(...)` y `O(...)` construyen criterios nuevos a partir de
  criterios existentes, sin que `Buscar` sepa cuántos ni cuáles son.
- **Inmutabilidad.** Los métodos que cambian estado devuelven una copia. Nada en
  el dominio se modifica en el lugar, de modo que no hay aliasing ni estado
  compartido que sincronizar.
- **Errores como valores.** En lugar de excepciones que interrumpen el flujo, el
  error es un valor de retorno que se compone y se inspecciona con `errors.Is`.

## 3. Capas y dirección de las dependencias

```
cmd/butaca  ──────▶  internal/{catalogo, usuarios, reproduccion, suscripciones}
                             │
                             ▼
                     internal/almacen (interfaz Repositorio)
                             │
                             ▼
                        pkg/funcional
```

Las flechas apuntan siempre hacia adentro. El dominio no conoce la base de datos
ni la interfaz de usuario: declara qué necesita (`Repositorio[T]`) y recibe una
implementación desde afuera. Es la regla que permite que las pruebas corran en
milisegundos, sin levantar un motor de base de datos.

## 4. Qué se decidió dejar fuera de la Etapa 1

| Decisión | Motivo |
|---|---|
| Sin base de datos real | La Etapa 1 es planeación; el contrato ya está definido y la implementación en memoria lo satisface |
| Sin dependencias de terceros | El dominio no las necesita; agregarlas antes de tiempo ata el diseño a una biblioteca |
| Sin interfaz gráfica ni API HTTP | El valor de esta etapa está en el modelo, no en la presentación |
| Sin autenticación | Depende del módulo de usuarios, que en esta etapa solo modela la cuenta |
