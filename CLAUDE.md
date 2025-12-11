# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A production-ready Go REST API service using clean architecture with unified Docker setup and VSCode debugging support. Built with Go 1.23, Chi router, PostgreSQL, and includes hot reload with Air and remote debugging with Delve.

## Development Commands

### Starting/Stopping Development

**VSCode (Recommended):**
```
Press F5 → Select "Start/Restart Docker Services & Debug"
```

**Manual:**
```bash
# Start development environment (uses docker-compose.yml + docker-compose.override.yml)
docker compose up -d

# Stop development environment
docker compose down

# View API logs
docker compose logs -f api

# Restart API container
docker compose restart api

# Rebuild containers (after dependency changes)
docker compose up -d --build
```

### Local Development (without Docker)
```bash
# Build the application
make build
# Or: go build -o bin/main ./cmd/api

# Run locally (requires PostgreSQL running on localhost:5432)
make run
# Or: go run ./cmd/api/main.go

# Run tests
make test
# Or: go test -v ./...

# Format code
make fmt
# Or: go fmt ./...

# Lint code (requires golangci-lint)
make lint
# Or: golangci-lint run

# Download and tidy dependencies
make deps
# Or: go mod download && go mod tidy
```

### Production Deployment
```bash
# Deploy production environment (uses only docker-compose.yml with production env vars)
docker compose --env-file .env.prod -f docker-compose.yml up -d

# Stop production
docker compose --env-file .env.prod -f docker-compose.yml down
```

## Architecture

### Clean Architecture (3-Layer)

The application follows clean architecture with strict dependency flow:

```
Handler → Service → Repository → Database
```

**Dependency Direction:** Each layer only depends on the layer below it. Never import upward (e.g., Service should never import Handler).

### Layer Responsibilities

1. **Handler Layer** (`internal/handler/`)
   - HTTP request handling and routing
   - Request parsing and validation
   - Response formatting (JSON)
   - HTTP status codes
   - Does NOT contain business logic

2. **Service Layer** (`internal/service/`)
   - Business logic and orchestration
   - Validation of business rules
   - Coordinates between multiple repositories if needed
   - Error handling and logging
   - Does NOT know about HTTP or request/response formats

3. **Repository Layer** (`internal/repository/`)
   - Data access and persistence
   - SQL queries and database operations
   - Raw database error handling
   - Does NOT contain business logic

### Key Architectural Patterns

**Dependency Injection:** All dependencies are injected through constructors in `cmd/api/main.go`. The initialization order is:
1. Config → Logger → Database
2. Repositories (receive database)
3. Services (receive repositories + logger)
4. Handlers (receive services + logger)
5. Router (receives handlers + logger)

**Interface-Based Design:** Repository and Service layers are defined as interfaces, allowing for easy testing and mocking.

**Context Propagation:** All repository and service methods accept `context.Context` as the first parameter for cancellation and timeout support.

## Project Structure

### Core Directories

- `cmd/api/main.go` - Application entry point and dependency wiring
- `internal/config/` - Configuration management (loads from env vars)
- `internal/domain/` - Business entities and request/response DTOs
- `internal/handler/` - HTTP handlers and router setup
- `internal/service/` - Business logic layer
- `internal/repository/` - Data access layer (interfaces + implementations)
- `internal/middleware/` - Custom HTTP middleware
- `pkg/database/` - Database connection utilities
- `pkg/logger/` - Structured logging setup (using Go's slog)
- `migrations/` - SQL migration files (manually run)

### Configuration

Configuration is loaded from environment variables with defaults in `internal/config/config.go`. The `.env` file is used in development, loaded via `godotenv`.

**Key Environment Variables:**
- `SERVER_PORT` (default: 8080)
- `DB_HOST` (default: localhost, in Docker: db)
- `DB_PORT` (default: 5432)
- `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `LOG_LEVEL` (default: info)

### Database Migrations

Migrations are in `migrations/` directory but are NOT automatically applied. They must be run manually:

```bash
# Connect to database and run migrations
docker compose exec postgres psql -U postgres -d go_api_db -f /migrations/001_create_users_table.sql
```

Or copy-paste SQL from migration files into a database client.

## Adding New Features

When adding a new entity/resource, follow this order:

1. **Domain Model** - Create entity and DTOs in `internal/domain/`
2. **Repository** - Define interface and implement in `internal/repository/`
3. **Service** - Define interface and implement in `internal/service/`
4. **Handler** - Create HTTP handler in `internal/handler/`
5. **Router** - Register routes in `internal/handler/router.go`
6. **Migration** - Create SQL file in `migrations/`
7. **Wire Dependencies** - Update `cmd/api/main.go` to initialize new components

### Example Pattern (from existing User implementation)

Repository interface defines data operations:
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
    GetByID(ctx context.Context, id int64) (*domain.User, error)
    // ...
}
```

Service interface defines business operations:
```go
type UserService interface {
    CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
    GetUser(ctx context.Context, id int64) (*domain.User, error)
    // ...
}
```

Handler receives service and exposes HTTP endpoints.

## Debugging

### Visual Debugging in VSCode

The development Docker setup includes Delve debugger listening on port 2345:

1. Open project in VSCode
2. Press F5 → Select "Start/Restart Docker Services & Debug"
3. This automatically starts Docker and attaches the debugger
4. Set breakpoints by clicking in the gutter
5. Make API requests → breakpoints will trigger

**Debug configuration is in `.vscode/launch.json`**

### Hot Reload

Air is configured (`.air.toml`) to watch for `.go` file changes and automatically restart the Delve debugger. Changes to code trigger rebuilds while keeping the debugger attached.

## API Endpoints

All API routes are versioned under `/api/v1`:

- `GET /health` - Health check
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

Routes are defined in `internal/handler/router.go` using Chi router.

## Database

Using standard library `database/sql` with `lib/pq` PostgreSQL driver. Connection pooling is configured in `cmd/api/main.go`:

- MaxOpenConns: 25
- MaxIdleConns: 5
- ConnMaxLifetime: 5 minutes

Database initialization and connection logic is in `pkg/database/postgres.go`.

## Logging

Structured logging using Go's built-in `log/slog` package. Logger is initialized in `pkg/logger/logger.go` and passed to services for consistent log formatting.

Custom logging middleware in `internal/middleware/logger.go` logs all HTTP requests with method, path, status, and duration.

## Module Path

The Go module is currently named `github.com/yourusername/go-api-service`. When cloning or forking, you may want to update this in `go.mod` and all import statements throughout the codebase.
