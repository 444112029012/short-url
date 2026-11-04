**繁體中文** | [English](README.en.md)

# short-url

以 Go 撰寫的短網址服務（chi router、透過 `modernc.org/sqlite` 使用 SQLite，建置不需要 cgo）。單一行程以純 HTTP 提供建立、轉址與統計功能。

`http://localhost` 僅供本機開發使用，不符合展示環境的 HTTPS 要求（SEC-012）。

- API 契約：[`docs/04-api/openapi.yaml`](docs/04-api/openapi.yaml)，說明見 [`docs/04-api/01-api-notes-v0.1.md`](docs/04-api/01-api-notes-v0.1.md)
- 方法規格：[`docs/02-architecture/06-method-specs-v0.2.md`](docs/02-architecture/06-method-specs-v0.2.md)
- 環境與機密：[`docs/09-cicd-env/02-environments-v0.1.md`](docs/09-cicd-env/02-environments-v0.1.md)、[`docs/09-cicd-env/03-secrets-management-v0.1.md`](docs/09-cicd-env/03-secrets-management-v0.1.md)
- 工程規範：[`docs/08-engineering/02-branch-commit-conventions-v0.1.md`](docs/08-engineering/02-branch-commit-conventions-v0.1.md)

## 前置需求

- Go 1.27.x（`go.mod` 宣告 `go 1.27.3`；CI 固定使用 1.27.3）。
- `curl`，用於下方範例。

## 建置與執行

```bash
go build ./...                       # compile everything
go build -o bin/shorturl ./cmd/shorturl
./bin/shorturl                       # or: go run ./cmd/shorturl
```

未設定任何環境變數時，行程會監聽 `127.0.0.1:8080`，建立並使用 `./data/shorturl.local.db`，並以 `http://localhost:8080` 組成短網址。啟動時，若資料庫目錄不存在會先建立（權限 `0700`），將資料庫檔案權限設為 `0600`，並在 `url_mapping` 資料表不存在時建立它。行程會記錄 `shorturl listening on <addr>`，收到 `SIGINT`/`SIGTERM` 時正常停止（關閉逾時 5 秒）。設定無效時，行程會在開始監聽前以結束碼 1 停止。

行程只讀取自身的環境變數，不會載入 `.env` 或 `.env.example`；請自行 export 變數，例如：

```bash
set -a; . ./.env.example; set +a
go run ./cmd/shorturl
```

## 試用 API

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

程式強制執行的請求規則：

- `POST /api/v1/urls` 必須帶 `Content-Type: application/json`，且主體必須恰好是一個 JSON 物件 `{"url": "..."}`。未知欄位、尾隨的 JSON、空主體，或超過 **32768 位元組**的主體，都會以 `400 invalid_url` 拒絕；剛好 32768 位元組的主體可接受。主體以原始位元組計算，包含 JSON 跳脫字元與空白。過大的主體會收到標準的 `400 invalid_url` 錯誤封裝，絕不會是 `413`。若來源已超過建立次數限制，即使主體過大，也會改為收到 `429 rate_limited`（限流檢查先執行，且不會讀取主體）。
- 請求列（含路徑與查詢字串）與所有請求標頭合併計算。總計不超過 16384 位元組（`MaxHeaderBytes`）保證可接受；介於 16384 與約 20480 位元組之間則不保證。總計超過約 20480 位元組（16384 加上 net/http 額外保留的 4096 位元組餘裕），或超過 500 個標頭值（標頭欄位行數，含 `Host`；net/http `DefaultMaxHeaderValueCount`），會在路由器執行前由 net/http 本身以 `431 Request Header Fields Too Large` 拒絕。該回應為 `Content-Type: text/plain; charset=utf-8`，主體是純文字 `431 Request Header Fields Too Large`，而非 JSON 錯誤封裝。來源伺服器不保證這個由 net/http 產生的回應帶有安全標頭；部署邊緣層可能會加上，用戶端不得依賴它們是否存在。
- `url` 必須是絕對的 `http` 或 `https` URL，長度最多 2048 個字元。
- 短碼為 8 個 base62 字元，由 `crypto/rand` 產生。
- 轉址送出的 `Location` 只會等於儲存的 `long_url`；查詢字串與請求標頭都無法改變它。每次成功轉址會讓 `click_count` 加一。

建立、轉址與統計 handler 的錯誤使用 `Content-Type: application/json` 與 `{"error":{"type":"...","message":"..."}}`，type 為以下四種之一：`invalid_url`（400）、`not_found`（404）、`rate_limited`（429）、`internal_error`（500）。訊息來自固定的目錄，不含內部細節。

在上述 handler 執行前由路由器本身產生的回應，不是 JSON：

- 不符合任何路由的路徑回傳 `404`，帶 `Content-Type: text/plain; charset=utf-8`，主體為 `404 page not found`。（對格式正確但不存在的短碼發出 `GET /<shortCode>`，屬於 handler 錯誤，回傳 JSON 的 `not_found` 主體。）
- 已知路徑搭配不支援的方法（例如 `DELETE /api/v1/urls`、`GET /api/v1/urls` 或 `HEAD /<shortCode>`）回傳 `405`，主體為空、沒有 `Content-Type`，並帶有列出允許方法的 `Allow` 標頭。

兩者仍會帶有 [HTTPS 與 HSTS](#https-與-hsts) 中列出的安全標頭，但沒有 `Cache-Control`。

## 設定

所有變數皆為選填。未設定、空值或只有空白的值會使用預設值。

| 變數 | 預設值 | 目前程式中的行為 |
|---|---|---|
| `HTTP_ADDR` | `127.0.0.1:8080` | 傳給 `http.Server` 的監聽位址。 |
| `DATABASE_PATH` | `./data/shorturl.local.db` | SQLite 檔案。每個環境請使用不同的檔案。 |
| `BASE_URL` | `http://localhost:8080` | `short_url` 的前綴（與 `HTTP_ADDR` 無關；兩者要一起修改）。必須是含主機的絕對 `http`/`https` URL，且不含使用者資訊、路徑前綴、查詢或片段；結尾的 `/` 會被移除。 |
| `RATE_LIMIT_CREATE_PER_MIN` | `30` | 每個來源每分鐘允許的建立請求數。必須是正整數，否則啟動失敗。 |
| `RATE_LIMIT_REDIRECT_PER_MIN` | `120` | 每個來源每分鐘允許的轉址請求數。必須是正整數，否則啟動失敗。 |
| `TRUSTED_PROXY_CIDRS` | 空 | 以逗號分隔的 CIDR，列出其 `X-Forwarded-For` / `X-Real-IP` 可信任的直接對端。空值表示永不使用這些標頭。CIDR 無效會導致啟動失敗。 |
| `HSTS_ENABLED` | `false` | 接受 `true/false`、`1/0`、`yes/no`、`on/off`（不分大小寫）；其他值會導致啟動失敗。見 [HTTPS 與 HSTS](#https-與-hsts)。 |
| `LOG_LEVEL` | `info` | 會讀入設定但**尚未使用**：事件日誌一律以 info 等級寫入。 |
| `APP_ENV` | `local` | 會讀入設定，但**尚未**用於任何行為。 |

環境文件中列出的 `TLS_CERT_FILE` / `TLS_KEY_FILE` **尚未實作**；行程不會終結 TLS。

HTTP 伺服器逾時固定寫在程式中：讀取標頭 5 秒、讀取 15 秒、寫入 15 秒、閒置 60 秒。

## 限流（SEC-007）

- `POST /api/v1/urls` 與 `GET /{shortCode}` 有限流；`GET /api/v1/urls/{shortCode}/stats` 與 `GET /health` 沒有。
- 檢查會在 handler 中最先執行。超過限制時回應 `429` 與 `{"error":{"type":"rate_limited","message":"Too many requests"}}`；不會建立短碼、不會送出轉址，`click_count` 也不會改變。不會送出 `Retry-After` 標頭。
- 建立與轉址使用各自的計數器。每個來源有一個固定的一分鐘視窗，從其第一個請求開始計算；被拒絕的請求不會消耗額度。
- 計數器存在行程記憶體中：重新啟動時會重設，且不會在多個執行個體間共用。每個計數器表最多容納 8192 個來源；表滿時會丟棄最舊的視窗（表滿本身絕不會回傳 429）。
- 來源識別依據 TCP 對端位址。只有當該對端位於 `TRUSTED_PROXY_CIDRS` 內時，服務才會讀取 `X-Forwarded-For`（所有標頭行，取最右邊且本身不是可信任代理的位址），若沒有可用的 `X-Forwarded-For`，則讀取 `X-Real-IP`。IPv4 對應的 IPv6 位址視為 IPv4 位址；IPv6 位址以 `/64` 分組。
- 位於反向代理後方時，請將 `TRUSTED_PROXY_CIDRS` 設為代理的位址範圍，否則所有用戶端會共用代理的額度。請將代理設定為附加到 `X-Forwarded-For`，並覆寫（而非轉傳）`X-Real-IP`。不要列入不受信任的用戶端可以連線進來的範圍。

## HTTPS 與 HSTS

每個回應都帶有 `X-Content-Type-Options: nosniff`、`Referrer-Policy: no-referrer`、`X-Frame-Options: DENY` 與 `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`。JSON 回應另外帶有 `Cache-Control: no-store`。例外：在請求抵達路由器前由 net/http 自行寫出的 `400`、`417`、`431` 與 `505` 回應（例如格式錯誤的請求、`Expect` 標頭不是 `100-continue`、請求列加標頭超過約 20480 位元組或超過 500 個標頭值，或不支援的 HTTP 版本）不會經過本服務的 handler。`400`、`431` 與 `505` 為 `text/plain; charset=utf-8`；`417` 的主體為空。來源伺服器不保證這些由 net/http 產生的回應帶有安全標頭；部署邊緣層可能會加上，用戶端不得依賴它們存在或不存在。

只有在 `HSTS_ENABLED` 為 true **且**請求是透過 TLS 抵達本行程時，才會送出 `Strict-Transport-Security: max-age=31536000`。由於 `cmd/shorturl` 只監聽純 HTTP，無論 `HSTS_ENABLED` 如何設定，發行的執行檔都不會送出 HSTS。當 TLS 在反向代理、通道或邊緣層終結時，HSTS 必須在那裡加上。依 SCR-010，代理／邊緣層的 HSTS 由維運於 G5 處理。

## 日誌與隱私

- 建立、轉址與統計的結果各寫成一行 `log/slog` JSON 輸出到 stderr。除了 slog 的 `time`、`level` 與 `msg`，成功的日誌行帶有 `event` 與 `short_code`，失敗的日誌行帶有 `event` 與 `error_type`（例如 `"event":"create_failure","error_type":"rate_limited"`）。
- 服務不會記錄用戶端 IP 位址、轉送標頭、`Authorization`、cookie、請求主體或長網址。資料庫只儲存 `short_code`、`long_url`、`click_count` 與 `created_at`。
- 不要將 SQLite 檔案或其目錄暴露在網路上（SEC-017）；應該只有 HTTP 連接埠可以連線。

## 測試

```bash
gofmt -l .                      # must print nothing
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...    # race detector needs cgo and a C toolchain
```

單元測試放在 `internal/` 下各套件旁。`internal/integration` 透過真實的 HTTP 連線執行完整的路由器、轉接器、限流器與暫存的 SQLite 資料庫。測試使用暫存目錄，不會動到 `./data`。

CI（`.github/workflows/ci-devsecops.yml`）在對 `main` 的 pull request 以及推送到 `main` 時，執行四項必要檢查：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck` 與 `unit-test`（`go test ./...`）。

## 目錄結構

| 路徑 | 內容 |
|---|---|
| `cmd/shorturl` | 行程進入點：設定、儲存層、元件串接、HTTP 伺服器 |
| `internal/config` | 環境變數載入與驗證 |
| `internal/httpapi` | 路由器、handler、錯誤對應、安全標頭 |
| `internal/application` | 建立、轉址與統計服務 |
| `internal/domain` | URL 驗證器、短碼產生器、錯誤類型 |
| `internal/infra/sqlite` | SQLite 儲存層與 schema |
| `internal/ratelimit` | 依來源的限流器與用戶端來源解析 |
| `internal/observability` | 結構化事件日誌與計數器 |
| `internal/health` | 健康狀態 |
| `internal/integration` | 端對端測試 |
| `docs/` | 需求、架構、API、資安、測試、工程、CI／環境文件 |
