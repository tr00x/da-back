# da-back — Mashynbazar Backend

> Go backend for a car marketplace platform. REST API + WebSocket real-time chat, structured as a clean layered architecture with Swagger documentation.

## Tech Stack

- **Language**: Go 1.21+
- **Architecture**: Layered (`cmd` / `internal` / `pkg`)
- **Transport**: HTTP REST + WebSocket
- **Docs**: Swagger / OpenAPI (`/docs`)
- **Build**: Makefile

## Project Structure

```
.
├── cmd/http/          # Entrypoint — HTTP server bootstrap
├── internal/          # Core business logic (handlers, services, repositories)
├── pkg/               # Shared utilities and helpers
├── docs/              # Swagger / OpenAPI specs
├── Makefile           # Build, run, and tooling commands
└── .env.example       # Environment variable reference
```

## Getting Started

### Prerequisites

- Go 1.21+
- Make

### Setup

```bash
cp .env.example .env
# Fill in DB connection, JWT secret, etc.
```

### Run

```bash
make run
```

Or build and run manually:

```bash
go build -o da-backend ./cmd/http
./da-backend
```

### Available Make targets

```bash
make run      # Run in development mode
make build    # Compile binary
make docs     # Regenerate Swagger docs (requires swaggo)
```

## API Documentation

Swagger UI is available at:

```
http://localhost:<PORT>/docs/index.html
```

## WebSocket

Real-time chat is exposed via WebSocket. A minimal test client is included at `socketClient.html` for local development and debugging.

## Environment Variables

See `.env.example` for the full list. Key variables:

| Variable | Description |
|---|---|
| `PORT` | HTTP server port |
| `DB_DSN` | Database connection string |
| `JWT_SECRET` | JWT signing secret |

## Notes

This repository contains the backend service for the Mashynbazar car marketplace platform. The frontend React client lives at [tr00x/offercar](https://github.com/tr00x/offercar).

A security audit of this codebase was conducted by [AmriTech](https://amritech.us), identifying SQL injection vectors, authentication weaknesses, and architectural issues — with remediation recommendations delivered as a formal technical report.
