# short-url

Short URL service in Go (chi router, SQLite via `modernc.org/sqlite`, no cgo needed for the build). One process serves create, redirect, and stats over plain HTTP.

`http://localhost` is for local development only and does not satisfy the demo HTTPS requirement (SEC-012).

- API contract: [`docs/04-api/openapi.yaml`](docs/04-api/openapi.yaml), notes in [`docs/04-api/01-api-notes-v0.1.md`](docs/04-api/01-api-notes-v0.1.md)
- Method spec: [`docs/02-architecture/06-method-specs-v0.2.md`](docs/02-architecture/06-method-specs-v0.2.md)
- Environments and secrets: [`docs/09-cicd-env/02-environments-v0.1.md`](docs/09-cicd-env/02-environments-v0.1.md), [`docs/09-cicd-env/03-secrets-management-v0.1.md`](docs/09-cicd-env/03-secrets-management-v0.1.md)
- Engineering conventions: [`docs/08-engineering/02-branch-commit-conventions-v0.1.md`](docs/08-engineering/02-branch-commit-conventions-v0.1.md)

## Prerequisites

- Go 1.27.x (`go.mod` declares `go 1.27.2`; CI pins 1.27.2).
- `curl` for the examples below.

## Build and run

```bash
go build ./...                       # compile everything
go build -o bin/shorturl ./cmd/shorturl
./bin/shorturl                       # or: go run ./cmd/shorturl
```

With no environment variables set, the process listens on `127.0.0.1:8080`, creates and uses `./data/shorturl.local.db`, and builds short URLs from `http://localhost:8080`. On start it creates the database directory (mode `0700`) if missing, sets the database file to mode `0600`, and creates the `url_mapping` table if it does not exist. It logs `shorturl listening on <addr>` and stops cleanly on `SIGINT`/`SIGTERM` (5-second shutdown timeout). Invalid configuration stops the process with exit code 1 before it listens.

The process reads only its environment. It does not load `.env` or `.env.example`; export the variables yourself, for example:

```bash
set -a; . ./.env.example; set +a
go run ./cmd/shorturl
```

## Try the API

```bash
# Create (201): {"short_code":"...","short_url":"http://localhost:8080/...","long_url":"..."}
curl -sS -X POST http://127.0.0.1:8080/api/v1/urls \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/path"}'

# Redirect (302, Location is the stored long_url). Use GET: HEAD is not
# routed and returns 405, so `curl -I` does not work here.
curl -sS -o /dev/null -D - http://127.0.0.1:8080/<shortCode>

# Stats (200): {"short_code":"...","click_count":N}
curl -sS http://127.0.0.1:8080/api/v1/urls/<shortCode>/stats

# Liveness (200, text/plain body "ok")
curl -sS http://127.0.0.1:8080/health
```

Request rules enforced by the code:

- `POST /api/v1/urls` requires `Content-Type: application/json` and a body of exactly one JSON object `{"url": "..."}`. Unknown fields, trailing JSON, an empty body, or a body over **32768 bytes** are rejected with `400 invalid_url`; a body of exactly 32768 bytes is accepted. A source that is over its create limit gets `429 rate_limited` instead, even with an oversized body (the limit check runs first and the body is not read).
- Request headers larger than about 20 KiB (`MaxHeaderBytes` 16384 plus the 4096 bytes of slack net/http adds) are rejected by net/http itself with `431 Request Header Fields Too Large` (`text/plain`) before the router runs.
- `url` must be an absolute `http` or `https` URL of at most 2048 characters.
- Short codes are 8 characters of base62 generated from `crypto/rand`.
- A redirect sends `Location` equal to the stored `long_url` only; query strings and request headers cannot change it. Each successful redirect increments `click_count` by one.

Errors from the create, redirect, and stats handlers use `Content-Type: application/json` and `{"error":{"type":"...","message":"..."}}` with one of four types: `invalid_url` (400), `not_found` (404), `rate_limited` (429), `internal_error` (500). Messages come from a fixed catalog and do not include internal details.

Responses produced by the router itself, before any of those handlers runs, are not JSON:

- A path that matches no route returns `404` with `Content-Type: text/plain; charset=utf-8` and the body `404 page not found`. (`GET /<shortCode>` for an unknown but well-formed code is a handler error and returns the JSON `not_found` body.)
- A known path with an unsupported method (for example `DELETE /api/v1/urls`, `GET /api/v1/urls`, or `HEAD /<shortCode>`) returns `405` with an empty body, no `Content-Type`, and an `Allow` header naming the permitted method.

Both still carry the security headers listed under [HTTPS and HSTS](#https-and-hsts), but no `Cache-Control`.

## Configuration

All variables are optional. An unset, empty, or whitespace-only value uses the default.

| Variable | Default | Behaviour in the current code |
|---|---|---|
| `HTTP_ADDR` | `127.0.0.1:8080` | Listen address passed to `http.Server`. |
| `DATABASE_PATH` | `./data/shorturl.local.db` | SQLite file. Use a separate file per environment. |
| `BASE_URL` | `http://localhost:8080` | Prefix for `short_url` (independent of `HTTP_ADDR`; change both together). Must be an absolute `http`/`https` URL with a host and no user info, path prefix, query, or fragment; a trailing `/` is removed. |
| `RATE_LIMIT_CREATE_PER_MIN` | `30` | Create requests allowed per source per minute. Must be a positive integer, otherwise startup fails. |
| `RATE_LIMIT_REDIRECT_PER_MIN` | `120` | Redirect requests allowed per source per minute. Must be a positive integer, otherwise startup fails. |
| `TRUSTED_PROXY_CIDRS` | empty | Comma-separated CIDRs of immediate peers whose `X-Forwarded-For` / `X-Real-IP` are trusted. Empty means those headers are never used. An invalid CIDR fails startup. |
| `HSTS_ENABLED` | `false` | Accepts `true/false`, `1/0`, `yes/no`, `on/off` (case-insensitive); anything else fails startup. See [HTTPS and HSTS](#https-and-hsts). |
| `LOG_LEVEL` | `info` | Read into config but **not used yet**: event logs are always written at info level. |
| `APP_ENV` | `local` | Read into config but **not used yet** by any behaviour. |

`TLS_CERT_FILE` / `TLS_KEY_FILE`, listed in the environments document, are **not implemented**; the process does not terminate TLS.

HTTP server timeouts are fixed in code: read header 5 s, read 15 s, write 15 s, idle 60 s.

## Rate limiting (SEC-007)

- `POST /api/v1/urls` and `GET /{shortCode}` are limited; `GET /api/v1/urls/{shortCode}/stats` and `GET /health` are not.
- The check runs before anything else in the handler. Over the limit the response is `429` with `{"error":{"type":"rate_limited","message":"Too many requests"}}`; no short code is created, no redirect is sent, and `click_count` does not change. No `Retry-After` header is sent.
- Create and redirect have separate counters. Each source gets a fixed one-minute window that starts at its first request; rejected requests do not consume quota.
- Counters live in process memory: they reset on restart and are not shared between instances. Each counter table holds up to 8192 sources; when it is full the oldest window is dropped (a full table never by itself returns 429).
- Source identity is the TCP peer address. Only when that peer is inside `TRUSTED_PROXY_CIDRS` does the service read `X-Forwarded-For` (all header lines, taking the rightmost address that is not itself a trusted proxy) and, if there is no usable `X-Forwarded-For`, `X-Real-IP`. IPv4-mapped IPv6 counts as the IPv4 address; IPv6 addresses are grouped by `/64`.
- Behind a reverse proxy, set `TRUSTED_PROXY_CIDRS` to the proxy's address range, or every client shares the proxy's quota. Configure the proxy to append to `X-Forwarded-For` and to overwrite (not forward) `X-Real-IP`. Do not list ranges that untrusted clients can connect from.

## HTTPS and HSTS

Every response carries `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY`, and `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`. JSON responses also get `Cache-Control: no-store`. Exception: `400` and `431` responses that net/http writes itself before the request reaches the router (for example a malformed request, or request headers over about 20 KiB) are `text/plain` and carry none of these security headers.

`Strict-Transport-Security: max-age=31536000` is sent only when `HSTS_ENABLED` is true **and** the request reached this process over TLS. Because `cmd/shorturl` listens on plain HTTP only, the shipped binary never sends HSTS, whatever `HSTS_ENABLED` says. When TLS ends at a reverse proxy, tunnel, or edge, HSTS has to be added there. Per SCR-010, HSTS at the proxy/edge layer is handled by operations in G5.

## Logs and privacy

- Create, redirect, and stats outcomes are written as one `log/slog` JSON line each to stderr. Besides slog's `time`, `level`, and `msg`, a success line carries `event` and `short_code`, and a failure line carries `event` and `error_type` (for example `"event":"create_failure","error_type":"rate_limited"`).
- The service does not log client IP addresses, forwarding headers, `Authorization`, cookies, request bodies, or long URLs. The database stores only `short_code`, `long_url`, `click_count`, and `created_at`.
- Do not expose the SQLite file or its directory to the network (SEC-017); only the HTTP port should be reachable.

## Test

```bash
gofmt -l .                      # must print nothing
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...    # race detector needs cgo and a C toolchain
```

Unit tests sit next to each package under `internal/`. `internal/integration` runs the full router, adapter, rate limiter, and a temporary SQLite database over real HTTP connections. Tests use temporary directories and do not touch `./data`.

CI (`.github/workflows/ci-devsecops.yml`) runs four required checks on pull requests to `main` and on pushes to `main`: `secrets-gitleaks`, `sast-gosec`, `sca-govulncheck`, and `unit-test` (`go test ./...`).

## Layout

| Path | Contents |
|---|---|
| `cmd/shorturl` | Process entry point: config, store, wiring, HTTP server |
| `internal/config` | Environment loading and validation |
| `internal/httpapi` | Router, handlers, error mapping, security headers |
| `internal/application` | Create, redirect, and stats services |
| `internal/domain` | URL validator, short code generator, error types |
| `internal/infra/sqlite` | SQLite store and schema |
| `internal/ratelimit` | Per-source rate limiter and client source resolution |
| `internal/observability` | Structured event logging and counters |
| `internal/health` | Health status |
| `internal/integration` | End-to-end tests |
| `docs/` | Requirements, architecture, API, security, testing, engineering, CI/environment documents |
