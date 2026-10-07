---
文件：任務拆解（G3）
版本：v0.1
日期：2026-10-07
狀態：Draft／Ready for Review
對齊備註：已對齊 DoD／acceptance map（2026-10-07）
負責角色：研發部
作者：研發部
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 任務拆解：SHORTURL（G3 開工就緒）

## 依據文件

| 文件 | 路徑 |
|---|---|
| 方法規格 | `docs/02-architecture/06-method-specs-v0.1.md` |
| OpenAPI | `docs/04-api/openapi.yaml` |
| 需求規格 | `docs/01-requirements/02-requirements-spec-v0.1.md` |
| 範圍 | `docs/01-requirements/03-scope-v0.1.md` |
| SEC v0.2 | `docs/06-security/01-sec-requirements-v0.2.md` |
| 需求—設計追溯 | `docs/02-architecture/07-req-design-traceability-v0.1.md` |
| 品保追溯／可測 | `docs/07-testing/02-traceability-matrix-v0.1.md`、`04-design-testability-review-v0.1.md` |
| **正式 DoD／驗收對照** | `docs/07-testing/08-definition-of-done-v0.1.md`、`09-dev-task-acceptance-map-v0.1.md`（ENG 已回填） |
| **DoR** | `docs/08-engineering/03-definition-of-ready-v0.2.md` |
| CI 資安門檻 | `docs/06-security/03-ci-security-gates-v0.1.md` |
| ADR（D-04） | `docs/02-architecture/adr/ADR-001～003` |
| 初步時程 | `docs/00-charter/02-preliminary-schedule-v0.1.md` |
| 決策 | D-04、D-06（`docs/99-changelog/decision-log.md`） |

## 約束與假設

- **技術棧（D-04）**：Go＋chi／SQLite／base62×8 CSPRNG。
- **零現金預算**：禁止依賴付費雲服務；CI 採免費公開 runner／本機可跑工具；SQLite 單檔、無外部託管 DB。
- **範圍外**：OUT-05 到期／TTL、OUT-01 自訂短碼、OUT-03 帳號——本表不排實作。
- **粒度**：單一 PR 可完成為佳；估時含寫碼＋自測；不含審查等待。
- **建議負責人角色**：研發＝實作主力；維運／安全＝協作標記。

## 覆蓋自檢（方法規格 §12 × OpenAPI）

| 規格列／端點 | 覆蓋任務 |
|---|---|
| HttpApiAdapter.handleCreate；POST `/api/v1/urls` | ENG-013、ENG-017 |
| HttpApiAdapter.handleRedirect；GET `/{shortCode}` | ENG-014、ENG-017 |
| HttpApiAdapter.handleStats；GET `/api/v1/urls/{shortCode}/stats` | ENG-015、ENG-017 |
| ShortUrlApplicationService.createShortUrl／getStats | ENG-007、ENG-009 |
| RedirectService.redirect | ENG-008 |
| UrlValidator.validateLongUrl／validateShortCode | ENG-005 |
| ShortCodeGenerator.generate | ENG-006 |
| UrlRepository.save／findByShortCode／exists | ENG-004 |
| ClickCounter.incrementAtomic／getCount | ENG-004 |
| InMemoryOrDbUrlStore（SQLite＋測試雙件） | ENG-003、ENG-004、ENG-018 |
| ErrorMapper.toHttpResponse | ENG-010 |
| RateLimitGuard.checkCreate／checkRedirect | ENG-011 |
| ObservabilityHooks（Create／Redirect／Stats Success／Failure） | ENG-012 |
| 設定、日誌、CI gates、README／本機啟動 | ENG-002、ENG-012、ENG-020、ENG-021 |
| 安全標頭 SEC-008 | ENG-016 |

---

## 任務一覽

| ID | 標題 | 估時（人時） | 建議角色 | 相依 |
|---|---|---:|---|---|
| ENG-001 | 專案骨架（Go module／目錄／chi 掛載空殼） | 3 | 研發 | — |
| ENG-002 | 設定與環境變數（BASE_URL、DB 路徑、限流門檻） | 2 | 研發；維運協作 | ENG-001 |
| ENG-003 | SQLite schema／遷移與 UrlMapping 約束 | 3 | 研發 | ENG-001、ENG-002 |
| ENG-004 | InMemoryOrDbUrlStore：Repository＋ClickCounter（參數化／原子） | 6 | 研發；安全協作（碼審） | ENG-003 |
| ENG-005 | UrlValidator（長網址＋短碼） | 3 | 研發 | ENG-001 |
| ENG-006 | ShortCodeGenerator（base62×8 CSPRNG） | 3 | 研發；安全協作 | ENG-001 |
| ENG-007 | ShortUrlApplicationService.createShortUrl | 4 | 研發 | ENG-004～006 |
| ENG-008 | RedirectService.redirect | 3 | 研發 | ENG-004、ENG-005 |
| ENG-009 | ShortUrlApplicationService.getStats | 2 | 研發 | ENG-004、ENG-005 |
| ENG-010 | ErrorMapper.toHttpResponse（四穩定錯誤碼） | 2 | 研發 | ENG-001 |
| ENG-011 | RateLimitGuard（建立／導向） | 4 | 研發 | ENG-001、ENG-002 |
| ENG-012 | ObservabilityHooks 與結構化日誌約束 | 3 | 研發；維運協作 | ENG-001 |
| ENG-013 | HttpApiAdapter.handleCreate（對接 OpenAPI） | 3 | 研發 | ENG-007、ENG-010～012 |
| ENG-014 | HttpApiAdapter.handleRedirect | 3 | 研發 | ENG-008、ENG-010～012 |
| ENG-015 | HttpApiAdapter.handleStats | 2 | 研發 | ENG-009、ENG-010、ENG-012 |
| ENG-016 | 安全回應標頭中介層（SEC-008 可設子集） | 2 | 研發；安全協作 | ENG-001 |
| ENG-017 | 組裝（DI／路由／main）與本機可啟動 | 3 | 研發 | ENG-002～016 |
| ENG-018 | 單元測試鉤子（Validator／Generator／ErrorMapper／Service） | 5 | 研發；品保對齊 | ENG-005～010 |
| ENG-019 | 整合測試鉤子（三 API＋SQLite 主路徑／錯誤路徑） | 6 | 研發；品保對齊 | ENG-017 |
| ENG-020 | CI 基本檢查對接 SEC-009～011 門檻 | 4 | 研發；維運；安全協作 | ENG-018、ENG-019 |
| ENG-021 | README／本機啟動與開發說明 | 2 | 研發；維運協作 | ENG-017 |

**任務數**：21  
**總估時**：約 **68 人時**（約 **8.5 人天**，以 8 人時／天計）  
**對齊時程**：G4 實作窗（粗估 2026-10-20～22）可單人串行完成核心路徑；CI／文件可與尾段平行。

---

## 任務明細

### ENG-001｜專案骨架（Go module／目錄／chi 空殼）

| 欄位 | 內容 |
|---|---|
| **描述** | 建立 Go module、標準目錄（`cmd/`、`internal/` 依方法規格分層：adapter／application／domain／infra）、掛載 chi 空 router 與 health 或空 main 可編譯。 |
| **對應規格** | 組裝基礎；HttpApiAdapter／各介面之套件邊界（方法規格 §1～11） |
| **REQ／NFR／SEC** | NFR-003（單實例）；ADR-001A |
| **驗收標準** | Given 乾淨環境；When `go build ./...`；Then 成功。目錄可對應 Adapter／Application／Ports／Infra。 |
| **相依** | — |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-002｜設定與環境變數

| 欄位 | 內容 |
|---|---|
| **描述** | 讀取：`BASE_URL`（組 short_url）、SQLite 檔路徑、建立限流 ≤30／分、導向 ≤120／分、listen address。預設適合本機；憑證／連線字串不進版控。 |
| **對應規格** | RateLimitGuard 門檻前置；InMemoryOrDbUrlStore 就緒；OpenAPI servers 區分 |
| **REQ／NFR／SEC** | SEC-007、SEC-012（文件標 localhost 非演示）、SEC-017（憑證不進版控）；NFR-003 |
| **驗收標準** | Given 未設必填時有安全預設或明確錯誤；When 改環境變數；Then 行為改變且範例 `.env.example`（無真實密鑰）存在。 |
| **相依** | ENG-001 |
| **估時** | 2 人時 |
| **角色** | 研發；維運協作 |

### ENG-003｜SQLite schema／遷移與 UrlMapping

| 欄位 | 內容 |
|---|---|
| **描述** | 建立 `url_mapping`（或等價）：`short_code` UNIQUE PK、`long_url`≤2048、`click_count`≥0 預設 0、`created_at`。禁止 `visitor_ip`／`user_agent`。啟動時遷移或 embed schema。 |
| **對應規格** | InMemoryOrDbUrlStore；資料模型 UrlMapping |
| **REQ／NFR／SEC** | REQ-001、010、011、012；NFR-005；SEC-004、SEC-017 |
| **驗收標準** | Given 空 DB；When 啟動／遷移；Then 表存在、UNIQUE 約束生效、schema 無禁止欄。 |
| **相依** | ENG-001、ENG-002 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-004｜InMemoryOrDbUrlStore（Repository＋ClickCounter）

| 欄位 | 內容 |
|---|---|
| **描述** | 實作 `save`、`findByShortCode`、`exists`、`incrementAtomic`、`getCount`：一律 `database/sql` 參數綁定；原子 `UPDATE … SET click_count = click_count + 1 WHERE short_code = ?`；SELECT 欄位最小化。提供測試用記憶體雙件（介面相容）可選。 |
| **對應規格** | UrlRepository §6；ClickCounter §7；InMemoryOrDbUrlStore §8 |
| **REQ／NFR／SEC** | REQ-001～004、010、012；SEC-004、**016**、**017**；NFR-003 |
| **驗收標準** | （1）碼審無字串拼接 SQL；（2）並發導向後次數與成功次數一致；（3）特製輸入不得改變查詢語意；（4）無 IP／UA 寫入。 |
| **相依** | ENG-003 |
| **估時** | 6 人時 |
| **角色** | 研發；安全協作（碼審） |

### ENG-005｜UrlValidator

| 欄位 | 內容 |
|---|---|
| **描述** | `validateLongUrl`：非空、可解析、僅 http／https、長度 ≤2048。`validateShortCode`：`^[A-Za-z0-9]{8}$`。純函式、無副作用。 |
| **對應規格** | UrlValidator §4 |
| **REQ／NFR／SEC** | REQ-005、006、007、011；SEC-001、013 |
| **驗收標準** | 單元：空、非 URL、`ftp://`／`javascript:`、L 與 L+1、合法 http(s)；非法短碼長度／字元 → `invalid_url`。 |
| **相依** | ENG-001 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-006｜ShortCodeGenerator

| 欄位 | 內容 |
|---|---|
| **描述** | 以 OS CSPRNG（如 `crypto/rand`）產生長度 8、字元集 A-Za-z0-9；禁止 `math/rand` 作為唯一來源。 |
| **對應規格** | ShortCodeGenerator §5；ADR-003 |
| **REQ／NFR／SEC** | REQ-006、010；SEC-006 |
| **驗收標準** | 產出符合 pattern；連續抽樣無遞增序列；CSPRNG 失敗回 `internal_error`。 |
| **相依** | ENG-001 |
| **估時** | 3 人時 |
| **角色** | 研發；安全協作 |

### ENG-007｜ShortUrlApplicationService.createShortUrl

| 欄位 | 內容 |
|---|---|
| **描述** | 驗證→產生短碼→exists／save；衝突有限次重試；成功 click_count=0；同一 long_url 可多碼。 |
| **對應規格** | ShortUrlApplicationService.createShortUrl §2 |
| **REQ／NFR／SEC** | REQ-001、005、007、010、011、012；SEC-001、006 |
| **驗收標準** | 合法 URL → 唯一短碼可回查；非法不落庫；並發不產生重複短碼（UNIQUE＋重試）。 |
| **相依** | ENG-004、ENG-005、ENG-006 |
| **估時** | 4 人時 |
| **角色** | 研發 |

### ENG-008｜RedirectService.redirect

| 欄位 | 內容 |
|---|---|
| **描述** | 驗證短碼→查庫→原子 +1→回 long_url；失敗不改次數、不回任意 Location；禁止依請求參數覆寫目標。 |
| **對應規格** | RedirectService §3 |
| **REQ／NFR／SEC** | REQ-002、003、006、008；SEC-013、014、016 |
| **驗收標準** | 存在：Location＝庫值、次數 N→N+1；不存在：not_found、次數不變；附加 query 無法改寫目標。 |
| **相依** | ENG-004、ENG-005 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-009｜ShortUrlApplicationService.getStats

| 欄位 | 內容 |
|---|---|
| **描述** | 驗證短碼→getCount；無寫入。 |
| **對應規格** | getStats §2；ClickCounter.getCount |
| **REQ／NFR／SEC** | REQ-004、006、008；SEC-013 |
| **驗收標準** | 存在回非負 click_count；不存在 not_found；非法短碼 invalid_url。 |
| **相依** | ENG-004、ENG-005 |
| **估時** | 2 人時 |
| **角色** | 研發 |

### ENG-010｜ErrorMapper

| 欄位 | 內容 |
|---|---|
| **描述** | AppError → HTTP JSON `{error:{type,message}}`；四碼對映 400／404／429／500；message 無堆疊／路徑／SQL。 |
| **對應規格** | ErrorMapper §9；OpenAPI ErrorResponse |
| **REQ／NFR／SEC** | REQ-007～009；SEC-003 |
| **驗收標準** | 四 type 對映正確；故意內部錯誤對外僅 `internal_error`、無敏感細節。 |
| **相依** | ENG-001 |
| **估時** | 2 人時 |
| **角色** | 研發 |

### ENG-011｜RateLimitGuard

| 欄位 | 內容 |
|---|---|
| **描述** | 進程內依來源（IP 或可信轉傳）：建立 ≤30／分、導向 ≤120／分；超限 `rate_limited`，不建碼／不導向／不計數。 |
| **對應規格** | RateLimitGuard §10 |
| **REQ／NFR／SEC** | SEC-007；NFR-003 |
| **驗收標準** | 同來源短時間超門檻得 429；門檻內成功；文件化門檻與設定一致。 |
| **相依** | ENG-001、ENG-002 |
| **估時** | 4 人時 |
| **角色** | 研發 |

### ENG-012｜ObservabilityHooks 與日誌約束

| 欄位 | 內容 |
|---|---|
| **描述** | 實作 onCreate／Redirect／Stats Success／Failure；僅穩定錯誤碼與 short_code；禁止 Authorization／Cookie 原文、禁止業務 IP 明細入業務庫。 |
| **對應規格** | ObservabilityHooks §11 |
| **REQ／NFR／SEC** | NFR-006；SEC-005 |
| **驗收標準** | 成功／失敗路徑可於日誌辨識事件；抽查日誌無完整 Authorization／Cookie。 |
| **相依** | ENG-001 |
| **估時** | 3 人時 |
| **角色** | 研發；維運協作（事件鍵建議對齊 G1-QA-O7） |

### ENG-013｜HttpApiAdapter.handleCreate

| 欄位 | 內容 |
|---|---|
| **描述** | POST `/api/v1/urls`：先限流→解析 JSON→應用服務→201 含 short_code／short_url／long_url；錯誤經 ErrorMapper；觸發 Hooks。 |
| **對應規格** | handleCreate §1；OpenAPI `createShortUrl` |
| **REQ／NFR／SEC** | REQ-001、005、007、009～012；SEC-001、003、007 |
| **驗收標準** | 對齊 OpenAPI：201／400／429／500；成功可導向同一 long_url。 |
| **相依** | ENG-007、ENG-010、ENG-011、ENG-012 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-014｜HttpApiAdapter.handleRedirect

| 欄位 | 內容 |
|---|---|
| **描述** | GET `/{shortCode}`：限流→RedirectService→302 Location；400／404／429／500。 |
| **對應規格** | handleRedirect §1；OpenAPI `redirectShortCode` |
| **REQ／NFR／SEC** | REQ-002、003、006、008、009；SEC-007、013、014 |
| **驗收標準** | 302 Location＝庫值；失敗不導向、不計數。 |
| **相依** | ENG-008、ENG-010、ENG-011、ENG-012 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-015｜HttpApiAdapter.handleStats

| 欄位 | 內容 |
|---|---|
| **描述** | GET `/api/v1/urls/{shortCode}/stats`→200 `{short_code,click_count}`。 |
| **對應規格** | handleStats §1；OpenAPI `getUrlStats` |
| **REQ／NFR／SEC** | REQ-004、006、008、009；SEC-013 |
| **驗收標準** | 對齊 200／400／404／500；讀取冪等。 |
| **相依** | ENG-009、ENG-010、ENG-012 |
| **估時** | 2 人時 |
| **角色** | 研發 |

### ENG-016｜安全回應標頭中介層

| 欄位 | 內容 |
|---|---|
| **描述** | 盡力設定 `X-Content-Type-Options: nosniff`；演示 HTTPS 時 HSTS max-age≥31536000；不主動洩漏詳細 Server／框架版本；CORS 若啟用勿對敏感回應用未驗證 `*`。 |
| **對應規格** | OpenAPI info SEC-008；api-notes |
| **REQ／NFR／SEC** | SEC-008；NFR-004 |
| **驗收標準** | 擷取建立／錯誤回應標頭符合文件化預期；localhost 不宣稱滿足 SEC-012。 |
| **相依** | ENG-001 |
| **估時** | 2 人時 |
| **角色** | 研發；安全協作 |

### ENG-017｜組裝（DI／路由／main）與本機可啟動

| 欄位 | 內容 |
|---|---|
| **描述** | 注入各元件；註冊三路由；啟動 SQLite；graceful 可選。`go run` 本機可走建立→導向→查次數。 |
| **對應規格** | 全公開方法組裝；C4 Container |
| **REQ／NFR／SEC** | REQ-001～004；NFR-002；SEC-017（DS1 非公網＝本機檔） |
| **驗收標準** | 本機一次完整核心路徑成功；DB 檔在設定路徑、非監聽公網埠作為 DB。 |
| **相依** | ENG-002～016 |
| **估時** | 3 人時 |
| **角色** | 研發 |

### ENG-018｜單元測試鉤子

| 欄位 | 內容 |
|---|---|
| **描述** | 為 Validator、Generator、ErrorMapper、Application／Redirect（可用記憶體雙件）建立 `go test` 套件；對應未來 TC 入口，不需等正式 TC 編號。 |
| **對應規格** | §4～5、§9、§2～3 |
| **REQ／NFR／SEC** | REQ-005～011；SEC-001、003、006、013 |
| **驗收標準** | `go test ./internal/...`（或約定路徑）通過；覆蓋主要肯定／否定案例。 |
| **相依** | ENG-005～010（可與實作平行補測） |
| **估時** | 5 人時 |
| **角色** | 研發；品保對齊（追溯 TC 欄後填） |

### ENG-019｜整合測試鉤子

| 欄位 | 內容 |
|---|---|
| **描述** | httptest＋暫存 SQLite：三 operationId 主路徑；驗證失敗、not_found、限流、Location 不可改寫、原子計數。 |
| **對應規格** | HttpApiAdapter 全；SEQ-01～07 場景 |
| **REQ／NFR／SEC** | REQ-001～012 核心；SEC-001、003、007、013、014、016 |
| **驗收標準** | CI／本機一鍵跑過；斷言四錯誤碼與 HTTP 狀態。 |
| **相依** | ENG-017 |
| **估時** | 6 人時 |
| **角色** | 研發；品保對齊 |

### ENG-020｜CI 基本檢查對接資安 gates

| 欄位 | 內容 |
|---|---|
| **描述** | 流水線至少：`go test`、格式／靜態檢查；對接安全部 `03-ci-security-gates`：SAST／SCA／secrets（**零預算**選開源或 GitHub 免費能力）。分支保護 required checks 與維運／安全共同定案。 |
| **對應規格** | 設計佔位「09-cicd」；NFR-004 |
| **REQ／NFR／SEC** | SEC-009、010、011；NFR-004、NFR-007 |
| **驗收標準** | PR 上 SAST 高／嚴重、SCA Critical（及門檻內 High）、secrets 真實機密樣式會失敗擋合併（或文件化與工具對齊之等價證明）。 |
| **相依** | ENG-018、ENG-019 |
| **估時** | 4 人時 |
| **角色** | 研發；維運；安全協作 |

### ENG-021｜README／本機啟動

| 欄位 | 內容 |
|---|---|
| **描述** | 根 README：前置（Go 版本）、環境變數、啟動、三 API 示例 curl、測試指令、指向 OpenAPI／工程規範；標明 localhost ≠ 演示 HTTPS（SEC-012）。 |
| **對應規格** | 維運／開發交接；NFR-006／008 文件精神 |
| **REQ／NFR／SEC** | NFR-002、006；SEC-012 |
| **驗收標準** | 陌生人依 README 可於 15 分鐘內完成本機核心路徑。 |
| **相依** | ENG-017 |
| **估時** | 2 人時 |
| **角色** | 研發；維運協作 |

---

## 非本迭代（範圍外／不排程）

| 項目 | 理由 |
|---|---|
| 短碼 TTL／自動過期 | OUT-05 |
| 自訂短碼、帳號、儀表板 | OUT-01～03 |
| 預覽頁 | OUT-10；SEC-002 殘餘已 D-06 接受 |
| 付費雲監控／託管 DB | 零預算；NFR-003 單實例 |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 研發部 | G3 初稿：ENG-001～021；覆蓋方法規格＋OpenAPI |
| v0.1.1 | 2026-10-07 | 研發部 | 依據補 DoD／acceptance map／DoR v0.2；**已對齊 DoD／acceptance map**（ENG ID 供 map 回填） |
