# Mashynbazar — Backend API

> Production Go backend for a multi-vertical vehicle marketplace. Covers cars, motorcycles, and commercial transport — with real-time WebSocket chat, multi-role auth, Firebase push notifications, and a full admin panel API.

**Frontend:** [tr00x/offercar](https://github.com/tr00x/offercar)

---

## Stack

| Layer | Technology |
|---|---|
| Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database | PostgreSQL via [pgx/v5](https://github.com/jackc/pgx) |
| Auth | JWT (access + refresh), Google OAuth, Phone OTP |
| Real-time | WebSocket (`gofiber/contrib/websocket`) |
| Push | Firebase Cloud Messaging |
| Validation | `go-playground/validator` |
| Logging | `zerolog` |
| Docs | Swagger (`swaggo/swag`) |
| OTP/SMS | Twilio |
| Images | EXIF stripping, resize, format normalization |
| Reports | Excel export (`excelize`) |

---

## Architecture

Clean layered architecture — Handler → Service → Repository. No frameworks dictating structure.

```
.
├── cmd/http/              # Entrypoint
├── internal/
│   ├── config/            # Config loader
│   ├── delivery/http/     # HTTP handlers (one file per domain)
│   ├── model/             # Domain models, DTOs, constants
│   ├── repository/        # DB queries (pure SQL via pgx)
│   ├── route/             # Route registration (one file per domain)
│   ├── service/           # Business logic
│   ├── storage/postgres/  # Migrations and seed SQL
│   └── utils/             # Helpers: response, email, OTP, migrate
├── pkg/
│   ├── auth/              # JWT, guards, CORS, validators
│   ├── files/             # File upload handling
│   ├── firebase/          # FCM push notifications
│   └── logger/            # Logger setup
├── docs/                  # Swagger generated output
└── Makefile               # Build and deploy commands
```

---

## API Overview

Base path: `/api/v1`

### Auth — `/auth`
| Method | Endpoint | Description |
|---|---|---|
| POST | `/admin-login` | Admin credentials login |
| POST | `/user-login-google` | Google OAuth login |
| POST | `/user-login-email` | Email + password login |
| POST | `/user-login-phone` | Phone OTP initiation |
| POST | `/user-phone-confirmation` | OTP verification |
| POST | `/user-forget-password` | Password reset request |
| POST | `/user-reset-password` | Password reset confirm |
| POST | `/user-email-confirmation` | Email confirmation |
| POST | `/third-party-login` | Dealer / broker / logist login |
| POST | `/send-application` | New partner application |
| POST | `/send-application-document` | Upload application docs `🔒` |
| POST | `/user-register-device` | FCM device registration `🔒` |
| DELETE | `/account/:id` | Account deletion `🔒` |

### Users & Catalog — `/users`
Cars, brands, models, generations, body types, transmissions, engines, drivetrains, fuel types, colors, cities, countries — full CRUD for listings with image/video upload, likes, price recommendations, and profile management.

### Motorcycles — `/motorcycles`
Category-based listings with dynamic parameters per category. Full CRUD with media management.

### Commercial Transport — `/comtrans`
Same pattern as motorcycles — categories, brands, models, dynamic parameters, full listing lifecycle.

### Third-Party — `/third-party`
Dedicated routes per role:
- **Dealers** — create/manage car listings with dealer-specific guards
- **Logists** — manage delivery destinations
- **Brokers** / **Car Services** — planned

### Admin — `/admin` `🔒 Admin only`
Full back-office API: users, countries, cities, regions, brands, models, generations, modifications, body types, transmissions, engines, drivetrains, fuel types, colors, moto/comtrans catalogs, company types, activity fields, partner applications (accept/reject).

### Real-time Chat — `/ws`
| Endpoint | Description |
|---|---|
| `GET /ws/conversations` | List user conversations `🔒` |
| `GET /ws/conversations/:id/messages` | Message history `🔒` |
| `GET /ws` | WebSocket upgrade — live messaging `🔒` |

WebSocket handler manages concurrent connections with per-connection write mutexes, graceful disconnects, and Firebase push delivery for offline users.

---

## Getting Started

### Prerequisites
- Go 1.24+
- PostgreSQL
- `swag` CLI (optional, for doc regeneration)

### Setup

```bash
git clone https://github.com/tr00x/da-back.git
cd da-back

cp .env.example .env
# Fill in DB credentials, JWT secrets, Twilio, Firebase config
```

### Run

```bash
go run ./cmd/http/main.go
```

### Build

```bash
make deploy   # Cross-compile to Linux amd64 and deploy via SCP + systemd restart
```

### Regenerate Swagger docs

```bash
make swag
```

Swagger UI available at: `http://localhost:8080/swagger/`

---

## Environment Variables

| Variable | Description |
|---|---|
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` | PostgreSQL connection |
| `PORT` | HTTP server port (e.g. `:8080`) |
| `ACCESS_KEY` / `REFRESH_KEY` | JWT signing secrets |
| `ACCESS_TIME` / `REFRESH_TIME` | Token TTLs (e.g. `1h`, `72h`) |
| `UPLOAD_PATH` | Local path for uploaded files |
| `LOGGER_FOLDER_PATH` / `LOGGER_FILENAME` | Log output config |
| `APP_MODE` | `dev` or `production` |
| `APP_VERSION` | Shown in logs and responses |

---

