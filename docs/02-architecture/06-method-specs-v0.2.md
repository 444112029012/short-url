---
文件：方法規格表
版本：v0.2.2
狀態：審查中（v0.2.2 隨設計部文件 PR #20 送審；§13 經 G4 PR #17 審查／複審，v0.2.1 已由設計部簽核，見 `09-signoff-DR-G4-01-method-specs-v0.2.1.md`）；DR-G4-01（§13.4.1）已於 2026-10-19 生效
負責角色：設計部
最後更新：2026-10-19
對應需求：SEC-016, SEC-017, REQ-001～REQ-012, NFR-001, NFR-006, SEC-001～SEC-008, SEC-013～SEC-014
專案代號：SHORTURL
取代：06-method-specs-v0.1.md（§1～§12 規格內容不變，僅新增 §13 與表頭註記）
對照依據：main `433140c5132778bf7cd731e3ecaca8187e50d678`（含 PR #7／#11／#16／#17；§13 初版 v0.2 原對照 PR #2 head `2698138`）
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
> **對照依據**：main `433140c5132778bf7cd731e3ecaca8187e50d678`（含 PR #7 ENG-011、#11 ENG-012、#16 ENG-016、#17 DR-G4-01；v0.2 初版原對照 PR #2 head `2698138`）。下文以函式／常數名稱引用程式碼位置。  
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

### 13.2 逐方法對映表

| 規格方法 | 規格簽名（§1～§11） | Go 簽名（檔案:行） | ctx 用途 | 規格例外 → Go 錯誤 → 對外（HTTP／type） | 行為一致性 |
|---|---|---|---|---|---|
| HttpApiAdapter.handleCreate | `handleCreate(req: HttpRequest) → HttpResponse` | `func (a *HTTPAPIAdapter) HandleCreate(r *http.Request) response`（`internal/httpapi/adapter.go`） | 取 `r.Context()` 傳入 `CreateShortURL`；本身無 ctx 參數 | rate_limited→429；invalid_url→400（含 Content-Type 非 JSON、本文空／>32768 B（`maxCreateBody`，DR-G4-01，見 13.4.1）、未知欄位、多餘 JSON：`decodeCreate`）；internal_error→500；皆經 `ErrorMapper` | 先 `CheckCreate` 後 `decodeCreate` 與 `CreateShortURL`；成功 201 `{short_code, short_url, long_url}`；觸發 Obs ✅ |
| HttpApiAdapter.handleRedirect | `handleRedirect(req) → HttpResponse` | `HandleRedirect(r *http.Request) response`（`adapter.go`） | `r.Context()` 傳入 `Redirect` | invalid_url→400；not_found→404；rate_limited→429；internal_error→500；錯誤時刪除 Location | 先 `CheckRedirect`；成功 302、Location＝`target.LongURL`（庫值）；不讀 query／header ✅ |
| HttpApiAdapter.handleStats | `handleStats(req) → HttpResponse` | `HandleStats(r *http.Request) response`（`adapter.go`） | `r.Context()` 傳入 `GetStats` | invalid_url→400；not_found→404；internal_error→500 | 200 `{short_code, click_count}`；不限流（符合 §1「建議可不限流」）；無寫入 ✅ |
| ShortUrlApplicationService.createShortUrl | `createShortUrl(longUrl: string) → Result<UrlCreated, AppError>` | `func (s *ShortURLApplicationService) CreateShortURL(ctx context.Context, longURL string) (domain.URLCreated, *domain.AppError)`（`internal/application/service.go`） | 僅下傳 `repo.Exists(ctx, …)`、`repo.Save(ctx, …)` | invalid_url：`ValidateLongURL` 結果原樣回傳；internal_error：`Generate` 失敗原樣回傳、Exists／Save 非衝突錯誤 `InternalWrap`、重試耗盡 `Internal()`；`ErrConflict`／`exists=true`→重試，不外露 | 驗證→generate→exists→save；有限次重試上限 8（`defaultGenerateAttempts`）；click_count=0；回傳 short_code／short_url（baseURL＋"/"＋code）／long_url ✅ |
| ShortUrlApplicationService.getStats | `getStats(shortCode: string) → Result<UrlStats, AppError>` | `GetStats(ctx context.Context, shortCode string) (domain.URLStats, *domain.AppError)`（`service.go`） | 僅下傳 `counter.GetCount(ctx, …)` | invalid_url（`ValidateShortCode`）；`ErrNotFound`→not_found；其他→internal_error；負值→internal_error | 無寫入；回傳 short_code、click_count≥0 ✅ |
| RedirectService.redirect | `redirect(shortCode: string) → Result<RedirectTarget, AppError>` | `func (s *RedirectService) Redirect(ctx context.Context, shortCode string) (domain.RedirectTarget, *domain.AppError)`（`internal/application/redirect.go`） | 僅下傳 `FindByShortCode(ctx, …)`、`IncrementAtomic(ctx, …)` | invalid_url（`ValidateShortCode`）；`ErrNotFound`→not_found；其他→internal_error；庫值空／含控制字元→internal_error（防禦，`safeLocation`） | 先 `FindByShortCode` 確認紀錄存在再 `IncrementAtomic` 原子 +1，成功後才回目標；任何失敗不回目標、不計數（increment 前之失敗）✅ |
| UrlValidator.validateLongUrl | `validateLongUrl(url: string) → Result<ValidatedUrl, AppError>` | `func (v *URLValidator) ValidateLongURL(raw string) (ValidatedURL, *AppError)`（`internal/domain/validator.go`） | 無（純函式） | invalid_url（MsgURLValidationFailed） | 空、長度>2048（以 Unicode 字元計，`MaxLongURLLength`、`utf8.RuneCountInString`）、不可解析、非 http(s)→invalid_url；另見 13.4 細化 ✅ |
| UrlValidator.validateShortCode | `validateShortCode(code: string) → Result<string, AppError>` | `ValidateShortCode(code string) (string, *AppError)`（`validator.go`） | 無 | invalid_url（MsgInvalidShortCode） | `^[A-Za-z0-9]{8}$`（`shortCodePattern`）；正規化＝原樣回傳 ✅ |
| ShortCodeGenerator.generate | `generate() → string`；例外 internal_error | `func (g *ShortCodeGenerator) Generate() (string, *AppError)`（`internal/domain/generator.go`） | 無 | CSPRNG 讀取失敗→`InternalWrap`；其他→`Internal()` | `crypto/rand`（`NewShortCodeGenerator`）、base62、長度 8、拒絕取樣去偏（`Generate`）✅ |
| UrlRepository.save | `save(record: UrlRecord) → Result<void, AppError>` | 介面 `Save(ctx context.Context, record domain.URLRecord) error`（`internal/ports/store.go` `URLRepository`）；實作 `(*URLStore).Save`（`internal/infra/sqlite/store.go`） | `ExecContext(ctx, insertSQL, …)` | 唯一衝突→`ports.ErrConflict`（`isUniqueViolation`；呼叫端重試）；其他→原生 error（上層轉 internal_error）；click_count<0→error | 參數化 INSERT（`insertSQL`）；schema 無 IP／UA（`createTableSQL`）✅ |
| UrlRepository.findByShortCode | `findByShortCode(code: string) → Option<UrlRecord>` | `FindByShortCode(ctx context.Context, code string) (*domain.URLRecord, error)`（`ports.URLRepository`；`(*URLStore).FindByShortCode`） | `QueryRowContext(ctx, selectByCodeSQL, code)` | None→`(nil, ports.ErrNotFound)`；實作異常→error（上層轉 internal_error） | 參數化 SELECT，僅四業務欄位（`selectByCodeSQL`）✅ |
| UrlRepository.exists | `exists(code: string) → boolean` | `Exists(ctx context.Context, code string) (bool, error)`（`ports.URLRepository`；`(*URLStore).Exists`） | `QueryRowContext(ctx, existsSQL, code)` | 不存在→`(false, nil)`；異常→error（上層轉 internal_error） | 參數化 `SELECT 1 … LIMIT 1`（`existsSQL`）✅ |
| ClickCounter.incrementAtomic | `incrementAtomic(shortCode: string) → Result<int, AppError>` | `IncrementAtomic(ctx context.Context, shortCode string) (int, error)`（`ports.ClickCounter`；`(*URLStore).IncrementAtomic`） | `QueryRowContext(ctx, incrementSQL, shortCode)` | not_found→`ports.ErrNotFound`；internal_error→error（上層 `InternalWrap`） | 單句原子 `UPDATE … SET click_count = click_count + 1 WHERE short_code = ? RETURNING click_count`（`incrementSQL`）；回傳遞增後值 ✅ |
| ClickCounter.getCount | `getCount(shortCode: string) → Result<int, AppError>` | `GetCount(ctx context.Context, shortCode string) (int, error)`（`ports.ClickCounter`；`(*URLStore).GetCount`） | `QueryRowContext(ctx, countSQL, shortCode)` | not_found→`ports.ErrNotFound`；internal_error→error | 參數化、僅選 click_count（`countSQL`）✅ |
| InMemoryOrDbUrlStore.* | 實作上列五方法 | `sqlite.URLStore`（`store.go` 之 `var _ ports.URLRepository = (*URLStore)(nil)` 等編譯期斷言實作兩介面） | 同上 | 同上 | SEC-017：DB 檔 0600、非網路監聽（`Open`）✅ |
| ErrorMapper.toHttpResponse | `toHttpResponse(error: AppError) → HttpResponse` | `func (ErrorMapper) ToHTTPResponse(err *domain.AppError) response`（`internal/httpapi/mapper.go`） | 無 | 不拋出；nil／未知 type→internal_error 500（`ToHTTPResponse`、`statusFor`） | invalid_url→400、not_found→404、rate_limited→429、internal_error→500（`statusFor`）；message 固定目錄（`safeMessage`）✅ |
| RateLimitGuard.checkCreate／checkRedirect | `(sourceId: string) → Result<void, AppError>` | `(*Guard).CheckCreate/CheckRedirect(sourceID string) *domain.AppError`（`internal/ratelimit/guard.go`；建構子 `ratelimit.NewGuard(createPerMin, redirectPerMin int)`） | 無 | 超過門檻→`domain.RateLimited(MsgTooManyRequests)`（`rate_limited`）→429；nil `Guard` 亦回 rate_limited（fail-closed） | 每來源、每方法固定 1 分鐘窗：建立 30 次／分鐘、導向 120 次／分鐘（`NewGuard` 非正值回落此預設；`cmd/shorturl` 由 `RATE_LIMIT_CREATE_PER_MIN`／`RATE_LIMIT_REDIRECT_PER_MIN` 設定）；放行扣一次額度、拒絕不扣；單程序、不跨程序共享。`sourceID` 由 adapter 以 `ratelimit.ClientSource(r, trusted)` 產生（僅連線端為可信代理時採 X-Forwarded-For／X-Real-IP；IPv4-mapped 折疊為 IPv4、IPv6 取 /64），Guard 原樣作為 bucket key；`sourceID` 為空字串時改用共用鍵 `empty-source`（常數 `emptySourceKey`，於 `allow`），所有空來源共用同一份額度、不得藉空身分繞過限流；不記錄來源位址（SEC-005）。`HandleCreate`／`HandleRedirect` 先檢查限流再解碼本文或導向。PR #7（ENG-011）實作 ✅ |
| ObservabilityHooks.on* | `(shortCode｜errorType: string) → void` | `(*Hooks).OnCreateSuccess(shortCode string)`、`OnCreateFailure(errorType string)`、`OnRedirectSuccess`、`OnRedirectFailure`、`OnStatsSuccess`、`OnStatsFailure`（`internal/observability/hooks.go`；建構子 `NewHooks()`、`NewHooksWithLogger(*slog.Logger)`） | 無 | 不拋出、不回錯誤；記錄失敗以 `recover` 吞掉（`emit`） | 每次呼叫遞增該事件之行程內計數器（`atomic.Uint64`），並以 `log/slog` 輸出一筆 JSON 紀錄（`NewHooks` 寫 stderr）：成功含 `event`、`short_code`；失敗含 `event`、`error_type`（穩定碼）；不記錄 Authorization、Cookie、本文、長網址、客戶端 IP（SEC-005）。PR #11（ENG-012）實作 ✅ |

### 13.3 錯誤轉換總表（埠層 → 應用層 → 對外）

| 來源（埠層 Go error） | 應用層轉換（函式） | AppError.Type | HTTP | 對外 message |
|---|---|---|---|---|
| `ports.ErrNotFound`（FindByShortCode／IncrementAtomic／GetCount） | `domain.NotFound(MsgShortCodeNotFound)`（`RedirectService.Redirect` 之 `FindByShortCode` 與 `IncrementAtomic` 錯誤分支；`ShortURLApplicationService.GetStats` 之 `GetCount` 錯誤分支） | `not_found` | 404 | Short code not found |
| `ports.ErrConflict`（Save） | 重試，不外露（`ShortURLApplicationService.CreateShortURL` 之 `Save` 錯誤分支）；達上限 8 次（`defaultGenerateAttempts`）→`domain.Internal()`（`CreateShortURL` 迴圈結束後） | （耗盡時）`internal_error` | 500 | An unexpected error occurred |
| 其他非 nil error（含 ctx 取消／逾時、驅動錯誤） | `domain.InternalWrap(err)`（`CreateShortURL` 之 `Exists`／`Save` 錯誤分支、`GetStats` 之 `GetCount` 錯誤分支；`Redirect` 之 `FindByShortCode`／`IncrementAtomic` 錯誤分支） | `internal_error` | 500 | An unexpected error occurred |
| validator 失敗 | 原樣回傳 | `invalid_url` | 400 | URL validation failed／Invalid short code format |
| RateLimitGuard 拒絕 | 原樣回傳（adapter 層） | `rate_limited` | 429 | Too many requests |

### 13.4 規格留白處之實作定值（記錄用，非規格變更）

| 項目 | 規格 v0.1 | 實作 | 判定 |
|---|---|---|---|
| 碰撞重試上限 | 「有限次重試」 | 8 次（`exists=true` 與 `ErrConflict` 皆計一次；`defaultGenerateAttempts`、`CreateShortURL`） | 滿足「有限次」；耗盡→internal_error 與規格一致 |
| 長度單位 | 長度≤2048 | Unicode 字元數（`utf8.RuneCountInString`，`ValidateLongURL`） | 與 OpenAPI `maxLength`（字元）、SQLite `CHECK(length(long_url) <= 2048)`（`createTableSQL`）一致 |
| 「不可解析」細化 | 不可解析→invalid_url | 另拒 ASCII 控制字元／空白（`hasCTLOrSpace`）、非 UTF-8、非絕對／opaque、空 host（皆於 `ValidateLongURL`） | 對齊 OpenAPI `long_url: format: uri`（RFC 3986 不允許未編碼空白／控制字元）；屬「不可解析」之實作細化 |
| 建立請求本文 | Content-Type 為 JSON（前置條件） | 非 `application/json`、空本文、未知欄位、多餘 JSON 值→400 invalid_url（`decodeCreate`）；**本文硬上限 32768 bytes（32 KiB），超過→400 invalid_url**（`maxCreateBody`；`decodeCreate` 於 JSON 解碼前以 `io.LimitReader(body, maxCreateBody+1)` 截斷讀取，最多讀 32769 B；`HandleCreate` 先 `CheckCreate` 後 `decodeCreate`，被限流請求回 429 且不讀本文；設計裁定 DR-G4-01，見 13.4.1） | 對齊 OpenAPI 僅定義 400 與 `additionalProperties: false`；上限為傳輸層防護，不改變 `long_url` ≤2048 字元之業務規則 |
| 標頭上限 | — | `http.Server.MaxHeaderBytes = 16384`（`cmd/shorturl/main.go` `newHTTPServer`，常數 `maxHeaderBytes`；未設 `MaxHeaderValueCount`，採 net/http 預設 500）。請求列（request line，含路徑與查詢字串）＋全部標頭合計 **16384 B 以下保證接受**；合計超過約 20480 B（`net/http` 讀取上限＝`MaxHeaderBytes`＋4096 B bufio 緩衝），或標頭值（header field line，含 Host）超過 **500 個**（`DefaultMaxHeaderValueCount`）時，於讀取階段直接回 **431**（`Content-Type: text/plain; charset=utf-8`），不經 router、`ErrorMapper` 與安全標頭中介層；源站不保證帶安全標頭（部署邊緣可能加上）。HTTP/1.1 行為 | 傳輸層防護（安全部 DR-G4-01 意見 C3）。431 已在 OpenAPI 0.1.1 以 `HeaderTooLarge`（text/plain）標示，不屬四類錯誤；業務錯誤契約（`ErrorResponse` 四類）不變；不破壞相容（不升 v2） |
| ctx 來源與期限 | — | HTTP 請求 ctx（客戶端斷線／Server Shutdown 時取消）；未另設業務逾時 | 符合 GO-MAP-03 |

#### 13.4.1 設計裁定 DR-G4-01：建立請求本文上限（裁定 2026-10-07；生效 2026-10-19）

- **問題**：舊上限 8192 bytes 會誤拒合法輸入。`long_url` 以 Unicode 字元計上限 2048；JSON 允許任意字元以 `\uXXXX` 跳脫，用戶端程式庫（如 Go `encoding/json` 將 `&`、`<`、`>` 跳脫為 `\u0026` 等；Python `json.dumps` 預設 `ensure_ascii=True`）會使本文膨脹。研發部實測：純 & 約 12263 B、非 ASCII 跳脫約 12298 B、emoji surrogate pair 跳脫約 24544 B。
- **理論最壞值**：每個字元最多跳脫為 12 bytes（BMP 外字元＝兩段 `\uXXXX`），2048 × 12 = 24576 B，加上 JSON 外框 `{"url":""}` 與合理空白，約 24.6 KB。
- **裁定**：本文硬上限 **32768 bytes（32 KiB）**。高於理論最壞值約 8 KB 餘裕（容納排版空白），且為 2 的次方易於設定。超過上限→400 `invalid_url`（不新增 413／新錯誤類型，維持 OpenAPI 四類錯誤）；讀取須以硬上限截斷，不得先讀入全量再判斷。實作採 `io.LimitReader(body, maxCreateBody+1)`（安全部意見 C2 認可之等效作法，未改用 `http.MaxBytesReader`）。
- **不變**：`long_url` ≤2048 字元、Content-Type、未知欄位、多餘 JSON 值之規則皆不變；超限仍為 400 `invalid_url`。業務錯誤契約（`ErrorResponse` 四類）不變、無新增錯誤類型；不破壞相容（不升 v2）。OpenAPI 0.1.1 補記本文上限 32768 B，並以 `HeaderTooLarge`（text/plain）標示傳輸層 431，不屬四類錯誤。
- **DoS 考量**：單請求記憶體上限由 8 KiB 升至 32 KiB；建立端點另有 30 次／分鐘／來源之速率限制，且 `HandleCreate` 先呼叫 `CheckCreate` 再呼叫 `decodeCreate`，被限流請求不讀本文。另 `http.Server.MaxHeaderBytes` 設為 16384（見 13.4「標頭上限」），標頭不再是大於本文上限之單請求記憶體面。
- **測試**：(a) `url` 為 2048 字元（含 `https://example.com/` 前綴），其餘字元皆為 BMP 外 emoji 且全以 surrogate pair 跳脫之本文→201；(b) 32768 B＋1 之本文→400 invalid_url；(c) 2049 字元（未跳脫）→400 invalid_url；另依安全部 C4：(d) 恰好 32768 B 之合法本文不因大小被拒；(e) 被限流來源送超大本文→429。
- **生效**：安全部複核意見 v0.1（`10-dr-g4-01-body-limit-security-opinion-v0.1`）**附條件同意**（條件 C1～C5）；C1（PR #7 ENG-011 限流先合併）已由 PR #7 滿足（main `24b8da0`）；由 PR #17 實作（`maxCreateBody` = 32768、`MaxHeaderBytes` = 16384 與上述測試），C2～C5 於該 PR 對照。本裁定自 2026-10-19 生效（PR #17 合併，merge commit `433140c`）。

### 13.5 與類別圖 v0.1 之已知出入（不影響本節判定）

- 類別圖 `ShortUrlApplicationService` 未列 `counter` 相依；本表 §2 `getStats` 前置條件為「counter／validator 已注入」，實作依本表注入 `ports.ClickCounter`（`ShortURLApplicationService` 之 `counter` 欄位，經 `NewShortURLApplicationService` 參數 `counter` 注入）。
- 類別圖 `ShortCodeGenerator.generate()` 回傳 `string`；本表 §5 列 `internal_error`，實作依本表回傳 `(string, *AppError)`。
- 以本表為準；類別圖待下次修訂對齊（非 G4 阻擋）。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | 公開方法規格齊；狀態審查中 |
| v0.1+R1 | 2026-10-07 | 設計部 | G2-QA-R1：SEC-016／017 參數化／原子寫入明文 |
| v0.2 | 2026-10-07 | 設計部 | G4-PR2-R1（路徑 a）：新增 §13「Go 實作對映」（命名 PascalCase、`Result`→`(T, *AppError)`、`Option`→`(*T, error)`、ctx 為首參且僅傳遞取消／逾時、埠層 `error`＋`ports.ErrNotFound`／`ErrConflict` 由應用層轉 AppError）；逐方法對照 PR #2 head `2698138`；表頭加 §13 指引；§1～§12 規格內容未變。裁示見 `08-ruling-G4-PR2-R1-go-mapping-v0.1.md` |
| v0.2.1 | 2026-10-19 | 設計部（研發部草擬） | 設計裁定 DR-G4-01（2026-10-07 裁定，2026-10-19 生效）：建立請求本文硬上限 8192→32768 B，超限仍 400 invalid_url，解碼前以 `io.LimitReader` 截斷；新增 §13.4.1；§13.4 新增「標頭上限」列（`MaxHeaderBytes`＝16384，超過由 net/http 回 431）；安全部意見 v0.1 附條件同意，C1～C5 已滿足（C1：PR #7；實作：PR #17）；§13.2／§13.4 行號引用改為函式名；G4-PR17-R2：§13.2 RateLimitGuard／ObservabilityHooks 列更新為現行實作（PR #7、#11）；G4-PR17-R3：§13.3／§13.5 行號引用改為函式名；HTTP 契約不變 |
| v0.2.2 | 2026-10-19 | 設計部 | S10：metadata 與 §13 引言之對照依據改為 main `433140c`（含 PR #7／#11／#16／#17）；S13：§13.2 RateLimitGuard 列補空 `sourceID` 使用共用鍵 `empty-source`（`emptySourceKey`）；§13.4.1 生效註記補 PR #17 merge commit `433140c`（生效日 2026-10-19 不變）；G4-PR20-R4／SCR-045：§13.4「標頭上限」列與 §13.4.1「不變」項改為「431 已在 OpenAPI 0.1.1 以 `HeaderTooLarge`（text/plain）標示，不屬四類錯誤；業務錯誤契約不變；不破壞相容（不升 v2）」——**更正**：v0.2.1 列「HTTP 契約不變」應讀作「業務錯誤契約不變、不破壞相容」，OpenAPI 已於 0.1.1 補記本文上限與 431（v0.2.1 列為已簽核歷史，保留原文）；G4-PR20 S1／SCR-044：§13.4 標頭上限改為「請求列＋標頭合計 16384 B 以下保證接受；超過約 20480 B 或標頭值超過 500 個回 431」；SCR-043：431 源站不保證帶安全標頭；G4-PR20 S4：狀態列更新；§1～§12 未動 |
