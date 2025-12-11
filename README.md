# Go API Service

A production-ready Go API service with unified Docker setup and VSCode debugging support.

## Features

- **Unified Docker Setup** - Single Dockerfile with multi-stage builds (dev + production)
- **VSCode Integration** - One-click debug that starts Docker and attaches debugger
- **Hot Reload** - Code changes auto-restart the container in development
- **Chi Router** - Lightweight, composable HTTP routing
- **Clean Architecture** - Handler → Service → Repository layers
- **PostgreSQL** - Database with standard library `database/sql`
- **Structured Logging** - Go's built-in `slog`
- **Production Ready** - Optimized, minimal Docker images

## Quick Start

### Development Mode (VSCode)

1. Open project in VSCode
2. Press **F5**
3. Select **"Start/Restart Docker Services & Debug"**

That's it! This will:
- Start PostgreSQL database (port 5432)
- Start Go API with hot reload (port 8080)
- Start Delve debugger (port 2345)
- Automatically attach the debugger

### Alternative: Manual Start

```bash
docker compose up -d
```

This uses `docker-compose.yml` (production base) + `docker-compose.override.yml` (development extensions)

### Test the API

```bash
# Health check
curl http://localhost:8080/health

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe"}'

# Get all users
curl http://localhost:8080/api/v1/users
```

## Project Structure

```
.
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── domain/                  # Business entities and DTOs
│   ├── handler/                 # HTTP handlers
│   ├── service/                 # Business logic
│   ├── repository/              # Data access
│   └── middleware/              # HTTP middleware
├── pkg/
│   ├── database/                # Database utilities
│   └── logger/                  # Logger setup
├── migrations/                  # SQL migrations
│
├── Dockerfile                   # Multi-stage: development + production
├── docker-compose.yml           # Base configuration (production)
├── docker-compose.override.yml  # Development extensions (auto-applied)
├── .air.toml                    # Hot reload configuration
│
└── .vscode/
    ├── launch.json              # VSCode debug configuration
    └── tasks.json               # Docker start/restart task
```

## Architecture

This boilerplate follows a 3-layer clean architecture:

1. **Handler Layer** - HTTP request handling, validation, response formatting
2. **Service Layer** - Business logic and orchestration
3. **Repository Layer** - Data access and persistence

Dependencies flow inward: Handler → Service → Repository

## Development Workflow

### 1. Start Development & Debug

Press **F5** in VSCode → Select **"Start/Restart Docker Services & Debug"**

This automatically:
- Stops any existing containers
- Rebuilds Docker images
- Starts PostgreSQL + API services
- Attaches the debugger

### 2. Set Breakpoints

- Open any Go file in VSCode
- Click in gutter (left of line numbers) to set breakpoints 🔴
- Breakpoints work immediately

### 3. Edit Code

- Edit any Go file
- Changes auto-reload in container ♻️
- Debugger automatically reconnects

### 4. Debug

- Make API request (curl, Postman, etc.)
- Breakpoint triggers in VSCode
- Step through: F10 (step over), F11 (step into), F5 (continue)
- Inspect variables in debug panel

### 5. View Logs

```bash
docker compose logs -f api
```

### 6. Stop Development

```bash
docker compose down
```

## API Endpoints

### Health Check
```
GET /health
```

### Users

```
POST   /api/v1/users      # Create user
GET    /api/v1/users      # Get all users
GET    /api/v1/users/{id} # Get user by ID
PUT    /api/v1/users/{id} # Update user
DELETE /api/v1/users/{id} # Delete user
```

## Production Deployment

### Deploy Production

```bash
# Use base docker-compose.yml with production env vars
docker compose --env-file .env.prod -f docker-compose.yml up -d
```

This builds the **production** stage from the Dockerfile:
- Optimized, minimal Alpine image
- No debugger or development tools
- Automatic health checks
- Auto-restart on failure

### Production Configuration

Create `.env.prod`:

```bash
ENVIRONMENT=production
LOG_LEVEL=info
SERVER_PORT=8080
DB_USER=produser
DB_PASSWORD=strongpassword
DB_NAME=go_api_prod
DB_SSLMODE=require
```

### How It Works

The same `Dockerfile` has multiple stages:
- **development** stage: Includes Air + Delve for debugging
- **production** stage: Minimal image with compiled binary

`docker-compose.yml` targets `production` by default.
`docker-compose.override.yml` switches to `development` stage for local work.

## Adding New Features

### Example: Add New Endpoint

1. **Create Domain Model** (`internal/domain/product.go`)
2. **Create Repository** (`internal/repository/product_repository.go`)
3. **Create Service** (`internal/service/product_service.go`)
4. **Create Handler** (`internal/handler/product_handler.go`)
5. **Register Routes** (`internal/handler/router.go`)
6. **Create Migration** (`migrations/002_create_products_table.sql`)

### Testing Your Changes

- Set breakpoints in VSCode
- Press F5 to start debugging
- Make requests, debug in real-time
- Hot reload handles rebuilds automatically

## Technology Stack

- **Go 1.23** - Latest stable Go version
- **Chi v5** - HTTP router
- **PostgreSQL 16** - Database
- **Docker & Docker Compose** - Unified multi-stage setup
- **Air** - Hot reload for development
- **Delve** - Go debugger
- **VSCode** - Integrated debugging

## Requirements

- Go 1.23+
- Docker & Docker Compose
- VSCode (for integrated debugging)

## Common Tasks

| Task | Command |
|------|---------|
| Start dev (VSCode) | Press `F5` |
| Start dev (manual) | `docker compose up -d` |
| Stop dev | `docker compose down` |
| View logs | `docker compose logs -f api` |
| Restart API | `docker compose restart api` |
| Deploy prod | `docker compose --env-file .env.prod -f docker-compose.yml up -d` |
| Rebuild | `docker compose up -d --build` |

## Cloud Deployment

Ready for cloud deployment with Pulumi (Step 2):
- AWS (ECS + RDS)
- GCP (Cloud Run + Cloud SQL)

Docker images are optimized and production-ready.

## License

MIT

---

**Get Started:** Open in VSCode and press `F5` 🚀
