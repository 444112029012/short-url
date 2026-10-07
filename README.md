# short-url

Short URL service (Go, chi, SQLite). `http://localhost` is for local development only and does not satisfy the demo HTTPS requirement (SEC-012).

## Run locally

Go 1.27.x is required. With no environment variables set, the process listens on `127.0.0.1:8080`, stores `./data/shorturl.local.db`, and builds short URLs from `http://localhost:8080`. Names and examples are in `.env.example` (the process does not load that file; export the variables yourself).

```bash
go run ./cmd/shorturl
```

```bash
curl -sS -X POST http://127.0.0.1:8080/api/v1/urls \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/path"}'
curl -sSI http://127.0.0.1:8080/<shortCode>
curl -sS http://127.0.0.1:8080/api/v1/urls/<shortCode>/stats
```

`GET /health` returns `ok`.

## Test

```bash
go test ./...
go build ./...
```

Rate limiting (ENG-011), structured logs (ENG-012), the rest of the security-header set (ENG-016), the full end-to-end suite (ENG-019), CI changes (ENG-020), and a full operator README (ENG-021) are not in this batch.
