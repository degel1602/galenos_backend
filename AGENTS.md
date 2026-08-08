# AGENTS.md

Galenos Pro Appointments API — REST API en Go (Gin + SQL Server 2022), arquitectura hexagonal (Ports and Adapters). Sigue las convenciones y el árbol del `README.md`.

## Comandos

```bash
go mod tidy
go run ./cmd/api        # arranca la API (puerto 8080)
go build ./...          # compilar
```

- **No hay** test suite, linter o CI configurados. Solo `gofmt` (`gofmt -w .`). `go test ./...` no encontrará tests; verifica con `go build ./...`.
- El código y los comentarios **están en español** (tipos/variables en inglés, comentarios en español). Mantén ese estilo.
- Rama de trabajo actual: `feature/gustavo` (hay branchs `main`, `feature/carlos`). Y No hacer push ni commit salvo que se pida.

## Arquitectura Hexagonal — dirección de dependencias (ver README)

La regla central: `internal/domain` **no importa** nada de `adapters` ni frameworks. La dependencia no se infiere fácilmente:

- `domain` → entidades puras, reglas de negocio, errores sentinel.
- `ports/{input,output}` → contratos (interfaces y tipos de datos) compartidos.
- `usecase` → implementa los input ports orquestando `domain` + output ports.
- `adapters/input/http` → Gin handlers, router, DTOs, respuesta estándar.
- `adapters/output/persistence/sqlserver` → `database/sql` + driver `microsoft/go-mssqldb`.
- `cmd/api/main.go` → **único** composition root; solo aquí se instancian las implementaciones concretas.

Al agregar código, respeta la dirección: `domain` nunca importa `adapters`, los repositorios concretos solo se ven como interfaces desde `usecase`, y los handlers nunca instancian repos sino que reciben el servicio por el port.

## Configuración / entorno

- `SQLSERVER_DSN` es **requerido**: `${config.Load}` falla y la app corta sí falta. Además `main.go` abre el pool y hace **PING con timeout de 5s** al arrancar: sin VPN/red al host SQL Server o contraseña incorrecta la API falla al inicio (quick fail).
- `main.go` carga `.env` solo si existe (desarrollo local: copia `.env.example → .env`). Los valores reales nunca se suben a git.
- Variables: `SERVER_PORT` (8080), `ALLOWED_ORIGINS` (default `http://localhost:4200`, separadas por coma), `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`.

## SQL / persistencia

- `patient_repository.go` consulta la tabla `pacientes` e invoca el SP `sp_listarPaciente`. Si el esquema cambia, ajusta ambas cosas; respeta las columnas que espera la stored procedure.
- `appointment_repository.go` usa una transacción `SERIALIZABLE` con `WITH (UPDLOCK, HOLDLOCK)`: lee disponibilidad del médico antes de insertar para evitar dobles reservas concurrentes. No cambies ese patrón de bloqueo sin revisarlo —es clave para la integridad concurrente.

## Router / HTTP

- Base path: `http://localhost:8080/api/v1`. Todas las respuestas usan el wrapper JSON estándar `{ "success": true, "data": <...> }` / `{ "success": false, "error": { "code": "...", "message": "..." } }`.
- Nuevos handlers: inspírate en los existentes de `internal/adapters/input/http/` para mantener el patrón de DTOs y envoltura de respuesta (`response.go`).