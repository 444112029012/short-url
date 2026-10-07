---
文件：方法規格表
版本：v0.2
狀態：審查中（G4-PR2-R1 補件：新增 §13 Go 實作對映；待審查部複審）
負責角色：設計部
最後更新：2026-10-07
對應需求：SEC-016, SEC-017, REQ-001～REQ-012, NFR-001, NFR-006, SEC-001～SEC-008, SEC-013～SEC-014
專案代號：SHORTURL
取代：06-method-specs-v0.1.md（§1～§12 規格內容不變，僅新增 §13 與表頭註記）
對照依據：PR #2 head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc`（tree `224c14157733205d2c255d48278b141ddaa47944`）
---

# 方法規格表

> 欄位：類別｜方法｜職責｜輸入參數｜回傳｜例外／錯誤｜前置條件｜後置條件｜對應 REQ／NFR／SEC｜備註  
> 型別為語言中立偽碼。穩定錯誤碼：`invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`。  
> **Go 實作對映**（命名、`Result`／`Option` 轉換、`context.Context`、埠層 `error`／sentinel）見 **§13**；依 §13 規則產生之簽名差異屬語言對映，**不構成**本表規格差異（裁示：`08-ruling-G4-PR2-R1-go-mapping-v0.1.md`）。

---

## 1. HttpApiAdapter

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| HttpApiAdapter | handleCreate | 將建立請求轉應用服務並對映 HTTP | req: HttpRequest；body.url: string，maxLength 2048 | HttpResponse：201 含 short_code／short_url／long_url；或 400／429／500 | 經 ErrorMapper：invalid_url、rate_limited、internal_error | 服務已啟動；Content-Type 為 JSON | 成功則庫新增一筆；失敗無新短碼 | REQ-001、005、007、009～012；SEC-001、003、007 | 非冪等；先 RateLimitGuard；觸發 ObservabilityHooks；執行緒安全依賴注入元件 |
| HttpApiAdapter | handleRedirect | 處理短碼造訪並回 302 或錯誤 | req: HttpRequest；path shortCode: string | HttpResponse：302 Location＝long_url；或 400／404／429 | invalid_url（格式）、not_found、rate_limited、internal_error | 無 | 成功導向後 click_count＝N+1；失敗次數不變 | REQ-002、003、006、008、009；SEC-007、013、014 | 非冪等（成功會計數）；Location 禁止依 query 改寫 |
| HttpApiAdapter | handleStats | 查詢短碼點擊次數 | req: HttpRequest；path shortCode: string | HttpResponse：200 {short_code, click_count}；或 400／404 | invalid_url、not_found、internal_error | 無 | 無寫入 | REQ-004、006、008、009；SEC-013 | 讀取冪等；建議可不限流或與建立分開政策 |

---

## 2. ShortUrlApplicationService

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| ShortUrlApplicationService | createShortUrl | 驗證、產生唯一短碼並持久化 | longUrl: string；非空、將經 UrlValidator | Result&lt;UrlCreated, AppError&gt;：含 short_code、short_url、long_url | invalid_url；internal_error（含重試耗盡） | validator／generator／repo 已注入 | 成功：存在唯一 short_code 對應 longUrl，click_count=0；同一 longUrl 可另建不同短碼 | REQ-001、005、007、010、011、012；SEC-001、006 | 非冪等；短碼衝突時有限次重試 generate＋save；單實例下 repo 須保證 unique |
| ShortUrlApplicationService | getStats | 驗證短碼並回傳點擊次數 | shortCode: string；須符合格式 | Result&lt;UrlStats, AppError&gt;：short_code、click_count≥0 | invalid_url；not_found；internal_error | counter／validator 已注入 | 無狀態變更 | REQ-004、006、008；SEC-013 | 讀取冪等 |

---

## 3. RedirectService

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| RedirectService | redirect | 查短碼、回傳導向目標並原子 +1 | shortCode: string | Result&lt;RedirectTarget, AppError&gt;：long_url 為庫中值 | invalid_url；not_found；internal_error | validator／repo／counter 已注入 | 成功：click_count 原子 +1；失敗：次數不變、不回任意 Location | REQ-002、003、006、008；SEC-013、014 | 非冪等；先確認紀錄存在再 increment；禁止使用請求參數覆寫 Location |

---

## 4. UrlValidator

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| UrlValidator | validateLongUrl | 驗證長網址規則 | url: string | Result&lt;ValidatedUrl, AppError&gt; | invalid_url（空、不可解析、非 http(s)、長度&gt;2048） | 無 | 無副作用 | REQ-005、007、011；SEC-001 | 純函式；冪等；執行緒安全 |
| UrlValidator | validateShortCode | 驗證短碼字元集與長度 | code: string；建議 ^[A-Za-z0-9]{8}$ | Result&lt;string, AppError&gt;：正規化後之 code | invalid_url（格式不符） | 無 | 無副作用 | REQ-006；SEC-013 | 純函式；不符者不視為有效短碼 |

---

## 5. ShortCodeGenerator

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| ShortCodeGenerator | generate | 以 CSPRNG 產生短碼 | 無 | string：長度 8、字元集 A-Za-z0-9（base62） | internal_error（CSPRNG 不可用） | 作業系統 CSPRNG 可用 | 每次呼叫獨立亂數；禁止可預測序列 | REQ-006、010；SEC-006 | 非冪等（每次不同）；熵≈47.6 bit（62^8）；見 ADR-003；執行緒安全 |

---

## 6. UrlRepository（介面）

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| UrlRepository | save | 持久化 UrlMapping | record: UrlRecord；short_code 唯一；long_url 長度≤2048；click_count≥0 | Result&lt;void, AppError&gt; | internal_error；唯一衝突由呼叫端重試 | record 已通過驗證 | 成功後 findByShortCode 可讀得；欄位無 visitor_ip／user_agent | REQ-001、010、012；SEC-004、**016** | 寫入；同一 long_url 允許多 short_code；**必須參數化／綁定；禁止字串拼接 SQL** |
| UrlRepository | findByShortCode | 依短碼查紀錄 | code: string | Option&lt;UrlRecord&gt;：Some 或 None（僅業務欄位） | 實作異常時由上層轉 internal_error | 無 | 無寫入 | REQ-002、004、008；SEC-016、**017** | 讀取冪等；**參數化**；SELECT 僅 short_code／long_url／click_count／created_at |
| UrlRepository | exists | 短碼是否已存在 | code: string | boolean | internal_error（可選） | 無 | 無寫入 | REQ-010；SEC-016 | 讀取冪等；**參數化**；供產生後碰撞檢查 |

---

## 7. ClickCounter（介面）

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| ClickCounter | incrementAtomic | 原子遞增點擊次數 | shortCode: string；須已存在 | Result&lt;int, AppError&gt;：遞增後值 | not_found；internal_error | 短碼已存在（由 RedirectService 保證或本方法檢查） | click_count 精確 +1，無遺失更新 | REQ-003；NFR-001；**SEC-016** | 非冪等；**原子寫入**（單句 `UPDATE … SET click_count = click_count + 1 WHERE short_code = ?` 或等價交易）；參數綁定；單實例假設 |
| ClickCounter | getCount | 讀取目前點擊次數 | shortCode: string | Result&lt;int, AppError&gt;：≥0 | not_found；internal_error | 無 | 無寫入 | REQ-004；SEC-016、017 | 讀取冪等；參數化；僅選 click_count（或最小欄位集） |

---

## 8. InMemoryOrDbUrlStore（Infrastructure）

> **SEC-016（必須，可驗證）**：凡對 DS1 之讀寫，一律使用**參數化查詢**或等價驅動／ORM 綁定（Go：`database/sql` 的 `?` 佔位符／具名參數）。**禁止**以字串拼接組合 SQL 或未轉義查詢片段。`incrementAtomic` 必須以**原子寫入**完成（見下表 SQL 例）。  
> **SEC-017（必須，可驗證）**：連線使用最小權限帳號；DS1 不對公網監聽；連線字串／憑證不進版控；查詢僅選業務欄位（見資料模型 §6、組裝說明）。

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| InMemoryOrDbUrlStore | save | 實作 UrlRepository.save | 同 UrlRepository.save | 同左 | 同左；唯一索引衝突 | StoreBackend 就緒（ADR-002A SQLite） | DS1 新增列 | REQ-001、010、012；SEC-004、**016**、**017**；NFR-003 | SQLite（演示）；測試可用記憶體雙件；schema 禁止 IP／UA；**參數化 INSERT** |
| InMemoryOrDbUrlStore | findByShortCode | 實作查詢 | 同介面 | 同介面 | 同介面 | 同左 | 同左 | REQ-002、004、008；SEC-016、017 | **參數化 SELECT**；欄位最小化 |
| InMemoryOrDbUrlStore | exists | 實作存在檢查 | 同介面 | 同介面 | 同介面 | 同左 | 同左 | REQ-010；SEC-016 | **參數化**（例如 `SELECT 1 … WHERE short_code = ?`） |
| InMemoryOrDbUrlStore | incrementAtomic | 實作原子 +1 | 同 ClickCounter | 同介面 | 同介面 | 同左 | 同左 | REQ-003；**SEC-016** | 參數化原子更新例：`UPDATE url_mapping SET click_count = click_count + 1 WHERE short_code = ?`；禁止先讀後寫無交易 |
| InMemoryOrDbUrlStore | getCount | 實作讀次數 | 同介面 | 同介面 | 同介面 | 同左 | 同左 | REQ-004；SEC-016、017 | 參數化；僅選 `click_count` |

---

## 9. ErrorMapper

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| ErrorMapper | toHttpResponse | 將 AppError 對映為對外 HTTP 本文 | error: AppError；type ∈ 四穩定碼 | HttpResponse：JSON `{error:{type,message}}`；message 無堆疊／路徑／SQL | 不拋出；未知 type 對映 internal_error | 無 | 對外不洩漏內部細節 | REQ-007～009；SEC-003 | 純對映；冪等；HTTP：invalid_url→400、not_found→404、rate_limited→429、internal_error→500 |

---

## 10. RateLimitGuard

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| RateLimitGuard | checkCreate | 檢查建立端點速率 | sourceId: string；非空（IP 或可信轉傳位址） | Result&lt;void, AppError&gt; | rate_limited | 門檻已設定（建議 30／分／來源） | 通過則消耗 1 次配額 | SEC-007；NFR-003 | 有狀態（計數器）；單實例進程內；非跨實例共享 |
| RateLimitGuard | checkRedirect | 檢查導向端點速率 | sourceId: string | Result&lt;void, AppError&gt; | rate_limited | 門檻建議 120／分／來源 | 通過則消耗 1 次配額 | SEC-007 | 同左；超限不導向、不計數 |

---

## 11. ObservabilityHooks

| 類別 | 方法 | 職責 | 輸入參數（名稱＋型別＋約束） | 回傳（型別＋語意） | 例外／錯誤 | 前置條件 | 後置條件 | 對應 REQ／NFR／SEC | 備註（副作用／冪等／執行緒） |
|---|---|---|---|---|---|---|---|---|---|
| ObservabilityHooks | onCreateSuccess | 記錄建立成功事件 | shortCode: string | void | 不向呼叫端拋；內部失敗僅降級日誌 | 無 | DS2／指標多一筆成功事件 | NFR-006；SEC-005 | 副作用：日誌／計數；禁止寫 Authorization／Cookie 原文、禁止業務 IP 明細 |
| ObservabilityHooks | onCreateFailure | 記錄建立失敗 | errorType: string（穩定碼） | void | 同左 | 無 | 失敗事件可觀測 | NFR-006；SEC-005 | 僅 errorType，不含敏感 payload |
| ObservabilityHooks | onRedirectSuccess | 記錄導向成功 | shortCode: string | void | 同左 | 無 | 成功事件可觀測 | NFR-006；SEC-005 | |
| ObservabilityHooks | onRedirectFailure | 記錄導向失敗 | errorType: string | void | 同左 | 無 | 失敗事件可觀測 | NFR-006；SEC-005 | |
| ObservabilityHooks | onStatsSuccess | 記錄查統計成功 | shortCode: string | void | 同左 | 無 | 事件可觀測 | NFR-006 | |
| ObservabilityHooks | onStatsFailure | 記錄查統計失敗 | errorType: string | void | 同左 | 無 | 事件可觀測 | NFR-006；SEC-005 | |

---

## 12. 覆蓋自檢

| Component 公開類型 | 公開方法已列 |
|---|---|
| HttpApiAdapter | handleCreate、handleRedirect、handleStats |
| ShortUrlApplicationService | createShortUrl、getStats |
| RedirectService | redirect |
| UrlValidator | validateLongUrl、validateShortCode |
| ShortCodeGenerator | generate |
| UrlRepository | save、findByShortCode、exists |
| ClickCounter | incrementAtomic、getCount |
| InMemoryOrDbUrlStore | 上述 Repository＋Counter 實作方法 |
| ErrorMapper | toHttpResponse |
| RateLimitGuard | checkCreate、checkRedirect |
| ObservabilityHooks | onCreate／Redirect／Stats Success／Failure |

DTO（UrlRecord、UrlCreated、UrlStats、RedirectTarget、AppError）為資料載體，無行為方法，不列入本表。

---

## 13. Go 實作對映（v0.2 新增；G4-PR2-R1）

> **性質**：本節為 §1～§12 語言中立規格之 **Go 語言繫結（language binding）**，屬規範性內容。依本節規則產生之簽名差異**不構成**方法規格差異、**不需 CR**；凡改變輸入約束、回傳欄位語意、錯誤類型／HTTP 狀態對映、副作用或其順序者，仍屬規格差異，須走 CR。  
> **對照依據**：PR #2（`feat/ENG-007-core-create-redirect-stats`）head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc`；下文「檔案:行」皆指該 commit 之 repo 路徑與行號。  
> **書面裁示**：`./08-ruling-G4-PR2-R1-go-mapping-v0.1.md`。

### 13.1 通用對映規則

| 規則 | 語言中立規格 | Go 對映 | 約束（不得改變之語意） |
|---|---|---|---|
| GO-MAP-01 命名 | `camelCase` 方法、`PascalCase` 類別；縮寫首字大寫（Url、Http、Id） | 匯出識別字 `PascalCase`；縮寫依 Go 慣例全大寫：`createShortUrl`→`CreateShortURL`、`HttpApiAdapter`→`HTTPAPIAdapter`、`UrlRepository`→`URLRepository`、`sourceId`→`sourceID`；`InMemoryOrDbUrlStore` 之 SQLite 實作為 `sqlite.URLStore` | 僅命名；方法集、職責不變。DTO 欄位 `snake_case`→Go 結構欄位 `PascalCase`；**對外 JSON 欄位名仍為 `snake_case`**（以 OpenAPI 為準） |
| GO-MAP-02 Result | `Result<T, AppError>` | `(T, *domain.AppError)`；`nil`＝Ok；非 `nil` 時 `T` 為零值且呼叫端不得使用 | `AppError.Type` 仍僅四穩定碼；`Result<void, AppError>` → `*domain.AppError` |
| GO-MAP-03 context | （規格無對應參數） | 凡會觸及 DS1 I/O 之**應用服務方法**與**埠方法**，以 `ctx context.Context` 為**第一參數** | ctx **僅**傳遞取消／逾時／期限（來源：HTTP `r.Context()`，下傳 `database/sql` 之 `ExecContext`／`QueryRowContext`）；**禁止**以 `context.WithValue` 攜帶業務輸入（長網址、短碼、來源 ID 等）或據以改變行為；ctx 取消／逾時所致錯誤一律屬「實作異常」→ `internal_error`，不得對映為 `not_found` 等業務錯誤。純函式（UrlValidator、ShortCodeGenerator、ErrorMapper）、RateLimitGuard、ObservabilityHooks 現不帶 ctx；日後如需加入，依本規則辦理屬對映 |
| GO-MAP-04 埠層錯誤契約 | `UrlRepository`／`ClickCounter`：`Result<_, AppError>`、`Option`、`boolean`；「實作異常時由上層轉 internal_error」；「唯一衝突由呼叫端重試」 | 埠方法回傳 Go `error`（不回 `AppError`），以 sentinel 表示規格之業務結果：`ports.ErrNotFound` ⇔ 規格 None／`not_found`；`ports.ErrConflict` ⇔ 「唯一衝突（呼叫端重試）」信號；其他非 `nil` error ⇔ 「實作異常」。比對一律 `errors.Is`（允許實作以 `%w` 包裝） | **應用層為唯一轉換點**：`ErrNotFound`→`domain.NotFound`（`not_found`）；`ErrConflict`→重試，不外露；其他→`domain.InternalWrap(err)`（`internal_error`；cause 僅供 `errors.Unwrap` 診斷，`Error()` 只回穩定碼，**不得序列化**）。埠層不決定對外錯誤碼 |
| GO-MAP-05 Option | `Option<T>` | `(*T, error)`：Some＝(非 nil 指標, nil)；None＝(nil, `ports.ErrNotFound`) | None 不得以 (nil, nil) 表示 |
| GO-MAP-06 boolean＋可選例外 | `boolean`；`internal_error（可選）` | `(bool, error)`；不存在＝(false, nil) | 「不存在」不得以 error 表示 |
| GO-MAP-07 HTTP 型別 | `HttpRequest`／`HttpResponse` | `*http.Request`／`httpapi` 套件內 `response{status, header, body}`，由 router 寫入 `http.ResponseWriter` | 狀態碼、標頭（含 Location）、本文依 §1／§9 與 OpenAPI |
| GO-MAP-08 消費端介面 | 依賴注入之元件型別 | 依 Go 慣例可由使用端定義最小介面（`application.CodeGenerator`、`httpapi.RateLimitChecker`、`httpapi.EventHooks`） | 方法集與簽名須與本表（經本節對映後）一致 |
| GO-MAP-09 生命週期／組裝 | 不在本表 | 建構子 `New*`、`sqlite.Open`／`Migrate`／`Close`、測試注入 `NewShortCodeGeneratorWithReader` | 非業務公開方法；不得承載本表以外之業務規則；生產碼 generator 必須使用 `crypto/rand` |

### 13.2 逐方法對映表（PR #2 head `2698138`）

| 規格方法 | 規格簽名（§1～§11） | Go 簽名（檔案:行） | ctx 用途 | 規格例外 → Go 錯誤 → 對外（HTTP／type） | 行為一致性 |
|---|---|---|---|---|---|
| HttpApiAdapter.handleCreate | `handleCreate(req: HttpRequest) → HttpResponse` | `func (a *HTTPAPIAdapter) HandleCreate(r *http.Request) response`（`internal/httpapi/adapter.go:77-104`） | 取 `r.Context()` 傳入 `CreateShortURL`（:87）；本身無 ctx 參數 | rate_limited→429；invalid_url→400（含 Content-Type 非 JSON、本文空／>8192 B、未知欄位、多餘 JSON：:152-175）；internal_error→500；皆經 `ErrorMapper` | 先 `CheckCreate`（:78）後解碼與建立；成功 201 `{short_code, short_url, long_url}`（:92-103）；觸發 Obs ✅ |
| HttpApiAdapter.handleRedirect | `handleRedirect(req) → HttpResponse` | `HandleRedirect(r *http.Request) response`（`adapter.go:108-124`） | `r.Context()` 傳入 `Redirect`（:113） | invalid_url→400；not_found→404；rate_limited→429；internal_error→500；錯誤時刪除 Location（:117） | 先 `CheckRedirect`（:109）；成功 302、Location＝`target.LongURL`（庫值，:121-123）；不讀 query／header ✅ |
| HttpApiAdapter.handleStats | `handleStats(req) → HttpResponse` | `HandleStats(r *http.Request) response`（`adapter.go:127-141`） | `r.Context()` 傳入 `GetStats`（:128） | invalid_url→400；not_found→404；internal_error→500 | 200 `{short_code, click_count}`；不限流（符合 §1「建議可不限流」）；無寫入 ✅ |
| ShortUrlApplicationService.createShortUrl | `createShortUrl(longUrl: string) → Result<UrlCreated, AppError>` | `func (s *ShortURLApplicationService) CreateShortURL(ctx context.Context, longURL string) (domain.URLCreated, *domain.AppError)`（`internal/application/service.go:50-97`） | 僅下傳 `repo.Exists(ctx, …)`（:67）、`repo.Save(ctx, …)`（:84） | invalid_url：validator 原樣回傳（:51-54）；internal_error：generator 失敗原樣回傳（:60-63）、Exists／Save 非衝突錯誤 `InternalWrap`（:68-70、:88）、重試耗盡 `Internal()`（:96）；`ErrConflict`／`exists=true`→重試（:71-73、:85-87），不外露 | 驗證→generate→exists→save；有限次重試上限 8（`defaultGenerateAttempts`，:13）；click_count=0（:81）；回傳 short_code／short_url（baseURL＋"/"＋code）／long_url（:90-94）✅ |
| ShortUrlApplicationService.getStats | `getStats(shortCode: string) → Result<UrlStats, AppError>` | `GetStats(ctx context.Context, shortCode string) (domain.URLStats, *domain.AppError)`（`service.go:100-116`） | 僅下傳 `counter.GetCount(ctx, …)`（:105） | invalid_url（:101-104）；`ErrNotFound`→not_found（:107-108）；其他→internal_error（:110）；負值→internal_error（:112-114） | 無寫入；回傳 short_code、click_count≥0 ✅ |
| RedirectService.redirect | `redirect(shortCode: string) → Result<RedirectTarget, AppError>` | `func (s *RedirectService) Redirect(ctx context.Context, shortCode string) (domain.RedirectTarget, *domain.AppError)`（`internal/application/redirect.go:25-46`） | 僅下傳 `FindByShortCode(ctx, …)`（:29）、`IncrementAtomic(ctx, …)`（:39） | invalid_url（:26-28）；`ErrNotFound`→not_found（:31-32、:40-42）；其他→internal_error（:34、:43）；庫值空／含控制字元→internal_error（防禦，:36-38） | 先確認紀錄存在（:29）再原子 +1（:39），成功後才回目標（:45）；任何失敗不回目標、不計數（increment 前之失敗）✅ |
| UrlValidator.validateLongUrl | `validateLongUrl(url: string) → Result<ValidatedUrl, AppError>` | `func (v *URLValidator) ValidateLongURL(raw string) (ValidatedURL, *AppError)`（`internal/domain/validator.go:23-42`） | 無（純函式） | invalid_url（MsgURLValidationFailed） | 空、長度>2048（以 Unicode 字元計，:11、:24）、不可解析、非 http(s)（:34-37）→invalid_url；另見 13.4 細化 ✅ |
| UrlValidator.validateShortCode | `validateShortCode(code: string) → Result<string, AppError>` | `ValidateShortCode(code string) (string, *AppError)`（`validator.go:45-50`） | 無 | invalid_url（MsgInvalidShortCode） | `^[A-Za-z0-9]{8}$`（:13）；正規化＝原樣回傳 ✅ |
| ShortCodeGenerator.generate | `generate() → string`；例外 internal_error | `func (g *ShortCodeGenerator) Generate() (string, *AppError)`（`internal/domain/generator.go:35-57`） | 無 | CSPRNG 讀取失敗→`InternalWrap`（:44-46）；其他→`Internal()`（:36-38、:53-55） | `crypto/rand`（:22）、base62、長度 8、拒絕取樣去偏（:12-13、:47-49）✅ |
| UrlRepository.save | `save(record: UrlRecord) → Result<void, AppError>` | 介面 `Save(ctx context.Context, record domain.URLRecord) error`（`internal/ports/store.go:20`）；實作 `internal/infra/sqlite/store.go:118-131` | `ExecContext(ctx, insertSQL, …)`（:123） | 唯一衝突→`ports.ErrConflict`（:125-127；呼叫端重試）；其他→原生 error（上層轉 internal_error）；click_count<0→error（:119-121） | 參數化 INSERT（:38-40）；schema 無 IP／UA（:28-36）✅ |
| UrlRepository.findByShortCode | `findByShortCode(code: string) → Option<UrlRecord>` | `FindByShortCode(ctx context.Context, code string) (*domain.URLRecord, error)`（`ports/store.go:21`；`store.go:133-151`） | `QueryRowContext(ctx, selectByCodeSQL, code)`（:136） | None→`(nil, ports.ErrNotFound)`（:139-141）；實作異常→error（上層轉 internal_error） | 參數化 SELECT，僅四業務欄位（:42-45）✅ |
| UrlRepository.exists | `exists(code: string) → boolean` | `Exists(ctx context.Context, code string) (bool, error)`（`ports/store.go:22`；`store.go:153-163`） | `QueryRowContext(ctx, existsSQL, code)`（:155） | 不存在→`(false, nil)`（:156-158）；異常→error（上層轉 internal_error） | 參數化 `SELECT 1 … LIMIT 1`（:47-48）✅ |
| ClickCounter.incrementAtomic | `incrementAtomic(shortCode: string) → Result<int, AppError>` | `IncrementAtomic(ctx context.Context, shortCode string) (int, error)`（`ports/store.go:28`；`store.go:165-178`） | `QueryRowContext(ctx, incrementSQL, shortCode)`（:167） | not_found→`ports.ErrNotFound`（:168-170）；internal_error→error（上層 `InternalWrap`） | 單句原子 `UPDATE … SET click_count = click_count + 1 WHERE short_code = ? RETURNING click_count`（:51-55）；回傳遞增後值 ✅ |
| ClickCounter.getCount | `getCount(shortCode: string) → Result<int, AppError>` | `GetCount(ctx context.Context, shortCode string) (int, error)`（`ports/store.go:29`；`store.go:180-193`） | `QueryRowContext(ctx, countSQL, shortCode)`（:182） | not_found→`ports.ErrNotFound`（:183-185）；internal_error→error | 參數化、僅選 click_count（:57-58）✅ |
| InMemoryOrDbUrlStore.* | 實作上列五方法 | `sqlite.URLStore`（`store.go:19-26` 編譯期斷言實作兩介面） | 同上 | 同上 | SEC-017：DB 檔 0600、非網路監聽（:60-102）✅ |
| ErrorMapper.toHttpResponse | `toHttpResponse(error: AppError) → HttpResponse` | `func (ErrorMapper) ToHTTPResponse(err *domain.AppError) response`（`internal/httpapi/mapper.go:28-46`） | 無 | 不拋出；nil／未知 type→internal_error 500（:29-31、:58-59） | invalid_url→400、not_found→404、rate_limited→429、internal_error→500（:48-61）；message 固定目錄（:63-77）✅ |
| RateLimitGuard.checkCreate／checkRedirect | `(sourceId: string) → Result<void, AppError>` | `CheckCreate/CheckRedirect(sourceID string) *domain.AppError`（`internal/ratelimit/guard.go:16-25`） | 無 | rate_limited（現 stub 一律放行） | 簽名一致；行為延後 ENG-011（SCR-002 追蹤，非本節範圍） |
| ObservabilityHooks.on* | `(shortCode｜errorType: string) → void` | `OnCreateSuccess(shortCode string)` 等六方法（`internal/observability/hooks.go:11-16`） | 無 | 不拋出 | no-op；行為延後 ENG-012 |

### 13.3 錯誤轉換總表（埠層 → 應用層 → 對外）

| 來源（埠層 Go error） | 應用層轉換（檔案:行） | AppError.Type | HTTP | 對外 message |
|---|---|---|---|---|
| `ports.ErrNotFound`（FindByShortCode／IncrementAtomic／GetCount） | `domain.NotFound(MsgShortCodeNotFound)`（`redirect.go:31-32、40-42`；`service.go:107-108`） | `not_found` | 404 | Short code not found |
| `ports.ErrConflict`（Save） | 重試，不外露（`service.go:85-87`）；達上限 8 次→`domain.Internal()`（:96） | （耗盡時）`internal_error` | 500 | An unexpected error occurred |
| 其他非 nil error（含 ctx 取消／逾時、驅動錯誤） | `domain.InternalWrap(err)`（`service.go:69、88、110`；`redirect.go:34、43`） | `internal_error` | 500 | An unexpected error occurred |
| validator 失敗 | 原樣回傳 | `invalid_url` | 400 | URL validation failed／Invalid short code format |
| RateLimitGuard 拒絕 | 原樣回傳（adapter 層） | `rate_limited` | 429 | Too many requests |

### 13.4 規格留白處之實作定值（記錄用，非規格變更）

| 項目 | 規格 v0.1 | PR #2 實作 | 判定 |
|---|---|---|---|
| 碰撞重試上限 | 「有限次重試」 | 8 次（`exists=true` 與 `ErrConflict` 皆計一次；`service.go:13、59`） | 滿足「有限次」；耗盡→internal_error 與規格一致 |
| 長度單位 | 長度≤2048 | Unicode 字元數（`utf8.RuneCountInString`，`validator.go:24`） | 與 OpenAPI `maxLength`（字元）、SQLite `CHECK(length(long_url) <= 2048)`（`store.go:34`）一致 |
| 「不可解析」細化 | 不可解析→invalid_url | 另拒 ASCII 控制字元／空白（`validator.go:27-29`）、非 UTF-8（:24）、非絕對／opaque（:31）、空 host（:38-40） | 對齊 OpenAPI `long_url: format: uri`（RFC 3986 不允許未編碼空白／控制字元）；屬「不可解析」之實作細化 |
| 建立請求本文 | Content-Type 為 JSON（前置條件） | 非 `application/json`、空本文、>8192 bytes、未知欄位、多餘 JSON 值→400 invalid_url（`adapter.go:18、152-175`） | 對齊 OpenAPI 僅定義 400 與 `additionalProperties: false` |
| ctx 來源與期限 | — | HTTP 請求 ctx（客戶端斷線／Server Shutdown 時取消）；未另設業務逾時 | 符合 GO-MAP-03 |

### 13.5 與類別圖 v0.1 之已知出入（不影響本節判定）

- 類別圖 `ShortUrlApplicationService` 未列 `counter` 相依；本表 §2 `getStats` 前置條件為「counter／validator 已注入」，實作依本表注入 `ports.ClickCounter`（`service.go:25`）。
- 類別圖 `ShortCodeGenerator.generate()` 回傳 `string`；本表 §5 列 `internal_error`，實作依本表回傳 `(string, *AppError)`。
- 以本表為準；類別圖待下次修訂對齊（非 G4 阻擋）。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | 公開方法規格齊；狀態審查中 |
| v0.1+R1 | 2026-10-07 | 設計部 | G2-QA-R1：SEC-016／017 參數化／原子寫入明文 |
| v0.2 | 2026-10-07 | 設計部 | G4-PR2-R1（路徑 a）：新增 §13「Go 實作對映」（命名 PascalCase、`Result`→`(T, *AppError)`、`Option`→`(*T, error)`、ctx 為首參且僅傳遞取消／逾時、埠層 `error`＋`ports.ErrNotFound`／`ErrConflict` 由應用層轉 AppError）；逐方法對照 PR #2 head `2698138`；表頭加 §13 指引；§1～§12 規格內容未變。裁示見 `08-ruling-G4-PR2-R1-go-mapping-v0.1.md` |
