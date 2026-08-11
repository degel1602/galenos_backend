# Arquitectura del Backend (Galenos Pro)

Este documento describe de manera estricta la arquitectura adoptada en el backend de Go (`galenos_backend`). Está diseñado como guía obligatoria para cualquier desarrollador actual o futuro que toque el proyecto.

El proyecto utiliza una combinación de **Arquitectura Hexagonal (Ports & Adapters)** y principios de **Clean Architecture**, aprovechando la naturaleza de encapsulación de Go mediante la carpeta `internal/`.

---

## 1. Estructura de Directorios

La estructura de carpetas refleja de manera directa la separación de responsabilidades:

```text
galenos_backend/
├── cmd/
│   └── api/
│       └── main.go                 <- (Composition Root) Punto de entrada.
├── internal/
│   ├── domain/                     <- (Capa más interna) Entidades puras y errores de negocio.
│   │
│   ├── ports/                      <- Contratos (Interfaces).
│   │   ├── input/                  <- Lo que el sistema hace (Casos de Uso).
│   │   └── output/                 <- Lo que el sistema necesita (Repositorios/APIs externas).
│   │
│   ├── usecase/                    <- Implementación de los 'ports/input'. Orquesta el 'domain'.
│   │
│   ├── adapters/                   <- Implementaciones técnicas.
│   │   ├── input/                  <- (Primary Adapters) Controladores HTTP (Gin), DTOs de request/response.
│   │   └── output/                 <- (Secondary Adapters) Repositorios SQL Server (database/sql).
│   │
│   └── config/                     <- Configuraciones del sistema, carga de entorno (.env).
```

---

## 2. Las 4 Reglas Estrictas de Arquitectura

Si rompes alguna de estas reglas, la arquitectura se degrada. **Deben respetarse sin excepción:**

### Regla 1: La Regla de Dependencia (De afuera hacia adentro)
* El paquete `domain/` **no puede importar nada** de `usecase/`, `adapters/` o `config/`. Es 100% puro y solo depende de la librería estándar de Go.
* El paquete `usecase/` solo puede importar de `domain/` y `ports/`. **Jamás** puede importar de `adapters/` (es decir, un caso de uso no puede saber qué framework web se usa ni qué base de datos guarda la información).
* El paquete `adapters/` importa de `usecase/`, `ports/` y `domain/`. Aquí es donde viven dependencias externas como `github.com/gin-gonic/gin` o `database/sql`.

### Regla 2: La Inversión de Control mediante Puertos
* Cuando la capa de `usecase/` necesita interactuar con la base de datos, **no llama a la base de datos directamente**. En su lugar, llama a una interfaz definida en `ports/output/`.
* La inyección real (qué repositorio exacto se usará) se resuelve en ejecución desde el `cmd/api/main.go`.

### Regla 3: Pureza de los DTOs (Data Transfer Objects)
* Las estructuras en `domain/` (ej. `Paciente`) **no deben tener tags** específicos de frameworks web (`json:"nombre"`) ni de bases de datos (`db:"id"`).
* El mapeo debe ocurrir en los adaptadores. 
  * En `adapters/input/http/`, un JSON entrante se mapea a un DTO y luego se convierte al modelo de `domain/` antes de pasarlo a `usecase/`.
  * Las respuestas HTTP se mapean desde el `domain/` hacia un DTO de respuesta antes de enviarse al cliente.

### Regla 4: El "Composition Root" asume el acoplamiento
* El único lugar del sistema donde las capas se cruzan de manera sucia es en `cmd/api/main.go`. 
* Aquí se instancian los repositorios (SQL Server), se inyectan en los casos de uso, y los casos de uso se inyectan en los controladores HTTP, para finalmente encender el servidor.

---

## 3. ¿Cómo agregar un nuevo flujo (Ej: Crear Especialidad)?

Si necesitas agregar un nuevo módulo, este es el flujo de trabajo profesional paso a paso:

1. **Definir el Dominio (`internal/domain/especialidad.go`):**
   Crea la entidad (struct `Especialidad`) y cualquier regla o error de negocio (ej. `var ErrEspecialidadDuplicada = errors.New(...)`).

2. **Definir los Contratos (`internal/ports/`):**
   * En `output/especialidad_repository.go`: Define la interfaz `EspecialidadRepository` con los métodos necesarios (ej. `Guardar(ctx context.Context, e *domain.Especialidad) error`).
   * En `input/crear_especialidad.go`: Define la interfaz del caso de uso.

3. **Crear el Caso de Uso (`internal/usecase/crear_especialidad.go`):**
   Implementa el input port. Pide en su constructor (`NewCrearEspecialidadUseCase`) la interfaz del output port (`EspecialidadRepository`). Realiza las validaciones de negocio y llama al repositorio.

4. **Crear el Adaptador de Salida (`internal/adapters/output/persistence/sqlserver/especialidad_repository.go`):**
   Implementa la interfaz del output port. Escribe la query SQL, ejecuta el Stored Procedure usando `database/sql`, mapea el resultado a la entidad de dominio y devuélvelo.

5. **Crear el Adaptador de Entrada (`internal/adapters/input/http/especialidad_handler.go`):**
   Crea el endpoint (ej. `POST /especialidades`). 
   * Extrae el JSON de la petición hacia un struct local (DTO).
   * Llama al caso de uso inyectado pasándole los datos limpios.
   * Devuelve un JSON con el resultado o código HTTP correspondiente.

6. **Cablear todo en el Composition Root (`cmd/api/main.go`):**
   * Instancia el repo: `repo := sqlserver.NewEspecialidadRepository(db)`
   * Instancia el caso de uso: `uc := usecase.NewCrearEspecialidadUseCase(repo)`
   * Instancia el handler: `handler := http.NewEspecialidadHandler(uc)`
   * Registra la ruta en Gin: `router.POST("/especialidades", handler.Crear)`

Manteniendo esta disciplina, el proyecto será siempre testeable, predecible y capaz de soportar la complejidad empresarial sin degradarse a código espagueti.
