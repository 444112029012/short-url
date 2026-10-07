---
文件：方法規格表
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：SEC-016, SEC-017, REQ-001～REQ-012, NFR-001, NFR-006, SEC-001～SEC-008, SEC-013～SEC-014
專案代號：SHORTURL
---

# 方法規格表

> 欄位：類別｜方法｜職責｜輸入參數｜回傳｜例外／錯誤｜前置條件｜後置條件｜對應 REQ／NFR／SEC｜備註  
> 型別為語言中立偽碼。穩定錯誤碼：`invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`。

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

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | 公開方法規格齊；狀態審查中 |
| v0.1+R1 | 2026-10-07 | 設計部 | G2-QA-R1：SEC-016／017 參數化／原子寫入明文 |
