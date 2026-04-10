## Prerequisites

- Go 1.25.2 or higher
- Make (optional, but recommended)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd queue
```

### 2. Setup Git Hooks

Install Lefthook and setup Git hooks:

```bash
make setup
```

This installs Lefthook and configures the commit-msg validation hook.

### 3. Run the Server

**Using Make:**
```bash
make run
```

**Or directly with Go:**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Available Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the server in development mode |
| `make build` | Build the server binary (outputs to `bin/server`) |
| `make lint` | Run linter (golangci-lint) |
| `make test` | Run all tests |
| `make setup` | Install Lefthook and Git hooks |

## API Endpoints

### Health Check

```
GET /health
```

Returns the health status of the server.

**Response:**
```json
{
  "status": "ok"
}
```

## Project Structure

```
.
├── cmd/server/
│   └── main.go           # Application entry point with graceful shutdown
├── internal/
│   ├── handler/
│   │   └── health.go     # Health check handler
│   ├── middleware/
│   │   └── logger.go     # Request logging middleware
│   └── router/
│       └── router.go     # Route configuration
├── Makefile              # Common development tasks
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── .golangci.yml         # Linter configuration
├── lefthook.yml          # Git hooks configuration
└── .gitignore
```

## Development

### Building for Production

```bash
make build
```

This creates a binary at `bin/server` that you can run:

```bash
./bin/server
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

## Configuration

The server currently runs on port `8080` by default. To change the port, modify the `Addr` field in [cmd/server/main.go](cmd/server/main.go):

```go
srv := &http.Server{
    Addr:    ":8080",  // Change this to your desired port
    Handler: r,
}
```

## Graceful Shutdown

The server handles graceful shutdown on receiving `SIGINT` or `SIGTERM` signals:
- Stops accepting new connections
- Waits for existing requests to complete (up to 5 seconds)
- Exits cleanly

## Git Hooks

This project uses [Lefthook](https://github.com/evilmartians/lefthook) for Git hooks management. Currently configured:
- **commit-msg**: Validates commit message format

## Dependencies

Main dependencies:
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework

See [go.mod](go.mod) for complete dependency list.
