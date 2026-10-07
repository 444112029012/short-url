---
文件：測試案例目錄
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-017
對應規格：方法規格、OpenAPI、SEC v0.2
專案代號：SHORTURL
關卡：G3
---

# 測試案例目錄：SHORTURL

> 編號規則：`UT-*` 單元｜`IT-*` 整合｜`CT-*` 契約｜`E2E-*` 端到端｜`CHK-*` 檢查表。  
> Given／When／Then 為摘要；實作時可拆 table-driven。層級與追溯見矩陣 v0.3。

---

## 1. 公開方法 → 單測計畫（覆蓋勾稽，須 100%）

| 類別 | 方法 | UT 計畫 ID | 層級補充 |
|---|---|---|---|
| HttpApiAdapter | handleCreate | UT-HC-01、UT-HC-02、UT-HC-03 | — |
| HttpApiAdapter | handleRedirect | UT-HR-01、UT-HR-02、UT-HR-03 | — |
| HttpApiAdapter | handleStats | UT-HS-01、UT-HS-02 | — |
| ShortUrlApplicationService | createShortUrl | UT-APP-01、UT-APP-02、UT-APP-03 | — |
| ShortUrlApplicationService | getStats | UT-APP-04、UT-APP-05 | — |
| RedirectService | redirect | UT-RED-01、UT-RED-02 | — |
| UrlValidator | validateLongUrl | UT-VAL-01～UT-VAL-05 | — |
| UrlValidator | validateShortCode | UT-VAL-06、UT-VAL-07 | — |
| ShortCodeGenerator | generate | UT-GEN-01、UT-GEN-02 | — |
| UrlRepository | save | UT-REPO-01 | IT-STORE-01 |
| UrlRepository | findByShortCode | UT-REPO-02 | IT-STORE-01 |
| UrlRepository | exists | UT-REPO-03 | IT-STORE-01 |
| ClickCounter | incrementAtomic | UT-CLK-01 | IT-STORE-02 |
| ClickCounter | getCount | UT-CLK-02 | IT-STORE-01 |
| InMemoryOrDbUrlStore | （實作上列） | — | IT-STORE-01～03 覆蓋 |
| ErrorMapper | toHttpResponse | UT-ERR-01、UT-ERR-02 | — |
| RateLimitGuard | checkCreate | UT-RL-01 | — |
| RateLimitGuard | checkRedirect | UT-RL-02 | — |
| ObservabilityHooks | onCreate／Redirect／Stats Success／Failure | UT-OBS-01（合併輕量） | — |

**公開方法單測計畫覆蓋：100%（無遺漏）。**

---

## 2. 單元案例（UT）

| ID | 層級 | Given／When／Then 摘要 | REQ／NFR／SEC | 方法／operationId |
|---|---|---|---|---|
| UT-VAL-01 | unit | Given 合法 https URL 長度≤2048；When validateLongUrl；Then Ok | REQ-005；SEC-001 | validateLongUrl |
| UT-VAL-02 | unit | Given 空字串；When validateLongUrl；Then invalid_url | REQ-005、007；SEC-001 | validateLongUrl |
| UT-VAL-03 | unit | Given `ftp://` 或 `javascript:`；When validateLongUrl；Then invalid_url | REQ-005；SEC-001 | validateLongUrl |
| UT-VAL-04 | unit | Given 長度＝2048 合法 http(s)；When validateLongUrl；Then Ok | REQ-011；SEC-001 | validateLongUrl |
| UT-VAL-05 | unit | Given 長度＝2049；When validateLongUrl；Then invalid_url | REQ-011、007；SEC-001 | validateLongUrl |
| UT-VAL-06 | unit | Given `Ab12Cd39`（^[A-Za-z0-9]{8}$）；When validateShortCode；Then Ok | REQ-006；SEC-013 | validateShortCode |
| UT-VAL-07 | unit | Given 過短／非法字元／過長；When validateShortCode；Then invalid_url | REQ-006；SEC-013 | validateShortCode |
| UT-GEN-01 | unit | When generate；Then 長度 8 且字元集 A-Za-z0-9 | REQ-006、010；SEC-006 | generate |
| UT-GEN-02 | unit | When 連續 generate N 次（建議≥100）；Then 無明顯遞增序列；來源為 CSPRNG（碼審／介面） | REQ-010；SEC-006 | generate |
| UT-APP-01 | unit | Given 合法 longUrl＋依賴就緒；When createShortUrl；Then 含 short_code／short_url／long_url，click=0 | REQ-001 | createShortUrl |
| UT-APP-02 | unit | Given 非法 longUrl；When createShortUrl；Then invalid_url 且 repo.save 未呼叫 | REQ-007；SEC-001 | createShortUrl |
| UT-APP-03 | unit | Given save 首次唯一衝突後成功；When createShortUrl；Then 有限次重試後成功 | REQ-010 | createShortUrl |
| UT-APP-04 | unit | Given 短碼存在；When getStats；Then click_count≥0 | REQ-004 | getStats |
| UT-APP-05 | unit | Given 短碼不存在；When getStats；Then not_found | REQ-008 | getStats |
| UT-RED-01 | unit | Given 短碼存在次數 N；When redirect；Then long_url＝庫值且 increment 一次 | REQ-002、003；SEC-014 | redirect |
| UT-RED-02 | unit | Given 短碼不存在；When redirect；Then not_found 且不 increment | REQ-008；SEC-013 | redirect |
| UT-REPO-01 | unit | Given 已驗證 UrlRecord；When save；Then 可 find；無 IP／UA 欄寫入 | REQ-001、012；SEC-004 | save |
| UT-REPO-02 | unit | Given 已存短碼；When findByShortCode；Then Some 僅業務欄位 | REQ-002、004；SEC-017 | findByShortCode |
| UT-REPO-03 | unit | Given 碼存在／不存在；When exists；Then true／false | REQ-010 | exists |
| UT-CLK-01 | unit | Given 短碼存在；When incrementAtomic；Then 回傳值＝原+1 | REQ-003；SEC-016 | incrementAtomic |
| UT-CLK-02 | unit | Given 短碼存在；When getCount；Then ≥0 且與狀態一致 | REQ-004 | getCount |
| UT-ERR-01 | unit | Given 四種 AppError type；When toHttpResponse；Then 400／404／429／500 且 type 正確 | REQ-007～009；SEC-003 | toHttpResponse |
| UT-ERR-02 | unit | Given 含堆疊／路徑／SQL 之內含錯誤；When toHttpResponse；Then 本文無洩漏 | REQ-009；SEC-003 | toHttpResponse |
| UT-RL-01 | unit | Given 超過建立門檻（≤30／分）；When checkCreate；Then rate_limited；門檻內 Ok | SEC-007 | checkCreate |
| UT-RL-02 | unit | Given 超過導向門檻（≤120／分）；When checkRedirect；Then rate_limited | SEC-007 | checkRedirect |
| UT-HC-01 | unit | Given 合法 JSON body；When handleCreate；Then 201＋欄位齊；先限流再服務 | REQ-001；createShortUrl | handleCreate |
| UT-HC-02 | unit | Given 非法 url；When handleCreate；Then 400 invalid_url | REQ-007；SEC-001 | handleCreate |
| UT-HC-03 | unit | Given 限流觸發；When handleCreate；Then 429 rate_limited 無新碼 | SEC-007 | handleCreate |
| UT-HR-01 | unit | Given 有效短碼；When handleRedirect；Then 302 Location＝庫值 | REQ-002；SEC-014；redirectShortCode | handleRedirect |
| UT-HR-02 | unit | Given 未登記短碼；When handleRedirect；Then 404 not_found | REQ-008 | handleRedirect |
| UT-HR-03 | unit | Given 非法短碼格式；When handleRedirect；Then 400 invalid_url 不計次 | REQ-006；SEC-013 | handleRedirect |
| UT-HS-01 | unit | Given 有效短碼；When handleStats；Then 200＋click_count | REQ-004；getUrlStats | handleStats |
| UT-HS-02 | unit | Given 非法／不存在；When handleStats；Then 400／404 | REQ-006、008 | handleStats |
| UT-OBS-01 | unit | Given 成功／失敗路徑；When 對應 hook；Then 被呼叫且參數無 Authorization／Cookie 原文、無業務 IP 明細 | NFR-006；SEC-005 | ObservabilityHooks（合併） |

---

## 3. 整合案例（IT）

| ID | 層級 | Given／When／Then 摘要 | REQ／NFR／SEC | 方法／operation |
|---|---|---|---|---|
| IT-STORE-01 | integration | Given SQLite schema；When save／find／exists／getCount；Then 參數化行為正確、無訪客 IP／UA 欄 | REQ-001、004、010；SEC-004、016、017 | InMemoryOrDbUrlStore |
| IT-STORE-02 | integration | Given 同短碼並發 N 次成功導向路徑 increment；When 完成；Then click_count＝N（無遺失更新） | REQ-003；SEC-016；NFR-001 輔助 | incrementAtomic |
| IT-STORE-03 | integration | Given 特製 short_code／long_url 含 `'`／`--` 等；When 讀寫；Then 不改變查詢語意、無注入 | SEC-016 | Store 讀寫 |
| IT-FLOW-01 | integration | Given 空庫；When 建立→redirect→stats；Then 次數一致、long_url 一致 | REQ-001～004 | 跨服務 |
| IT-FLOW-02 | integration | Given 驗證失敗建立；When 完成；Then 庫無新列 | REQ-007；SEC-001 | createShortUrl |
| IT-FLOW-03 | integration | Given 同一 long_url 建立兩次；When 比對短碼；Then 兩不同短碼皆可導向 | REQ-012 | createShortUrl |

---

## 4. 契約案例（CT）

| ID | 層級 | Given／When／Then 摘要 | REQ／SEC | operationId |
|---|---|---|---|---|
| CT-CREATE-01 | contract | When POST /api/v1/urls 合法；Then 201 符合 CreateUrlResponse schema（short_code pattern、欄位必填） | REQ-001、006、011 | createShortUrl |
| CT-CREATE-02 | contract | When POST 非法 url；Then 400＋ErrorResponse.type∈四碼且＝invalid_url | REQ-007；SEC-001 | createShortUrl |
| CT-REDIR-01 | contract | When GET /{shortCode} 有效；Then 302＋Location header（uri） | REQ-002；SEC-014 | redirectShortCode |
| CT-REDIR-02 | contract | When GET 不存在；Then 404＋not_found 形狀 | REQ-008 | redirectShortCode |
| CT-STATS-01 | contract | When GET …/stats 有效；Then 200 符合 UrlStatsResponse（click_count≥0） | REQ-004 | getUrlStats |
| CT-ERR-01 | contract | When 觸發各錯誤；Then type enum 僅四值、無 extra properties（對齊 OpenAPI） | REQ-009；SEC-003 | ErrorResponse |

---

## 5. E2E 案例

| ID | 層級 | Given／When／Then 摘要 | REQ／NFR／SEC | operationId |
|---|---|---|---|---|
| E2E-01 | e2e | Given 服務啟動；When 建立→導向（不跟隨或檢查 Location）→stats；Then 全成功且次數≥1 | REQ-001～004；NFR-002 | 三 API |
| E2E-02 | e2e | Given 服務啟動；When 建立空／ftp／L+1；Then 400、庫無碼 | REQ-005、007、011；SEC-001 | createShortUrl |
| E2E-03 | e2e | Given 無此短碼；When 導向／stats；Then 404、不導向任意站 | REQ-008 | redirect／stats |
| E2E-04 | e2e | Given 已建立；When 導向並附加 `?url=https://evil.example`；Then Location 仍為庫值 | SEC-014 | redirectShortCode |
| E2E-05 | e2e | Given 同來源快速超過建立門檻；When POST；Then 429 rate_limited | SEC-007 | createShortUrl |
| E2E-06 | e2e | Given 非法短碼 path；When GET；Then 400、不計次 | REQ-006；SEC-013 | redirectShortCode |

---

## 6. 檢查表（CHK）— NFR／SEC 非自動化為主

| ID | 層級 | 檢查摘要 | 對應 |
|---|---|---|---|
| CHK-NFR-001 | checklist＋抽樣 | 暖機後建立／導向各≥10 次，P95≤2s；記錄環境基線 | NFR-001 |
| CHK-NFR-002 | checklist | 演示時段核心路徑完整成功≥1 | NFR-002 |
| CHK-NFR-003 | checklist | 設計文件標單實例假設 | NFR-003 |
| CHK-NFR-004 | checklist | ASVS 文件＋SEC 清單＋CI 門檻存在；G5 高風險關閉或接受 | NFR-004 |
| CHK-NFR-005 | checklist | 對照隱私文件：模型無姓名電郵；點擊無 IP 明細 | NFR-005；SEC-004 |
| CHK-NFR-006 | checklist | 演示可指出建立／導向成功失敗觀測證據 | NFR-006 |
| CHK-NFR-007 | checklist | docs/gates 關卡紀錄；必改關閉或接受 | NFR-007 |
| CHK-NFR-008 | checklist | 各章已簽核或 N/A 附理由（含 UI） | NFR-008 |
| CHK-SEC-002 | checklist | 殘餘聲明＋僅 http(s) 抽查；RR-001／D-06 已登錄 | SEC-002 |
| CHK-SEC-005 | checklist | 日誌抽查無完整 Authorization／Cookie | SEC-005 |
| CHK-SEC-008 | checklist | 擷取 HSTS（演示 HTTPS）／nosniff 等標頭 | SEC-008 |
| CHK-SEC-009 | checklist | CI SAST job 存在且 High/Critical 擋合併 | SEC-009 |
| CHK-SEC-010 | checklist | CI SCA 門檻與鎖檔／SBOM | SEC-010 |
| CHK-SEC-011 | checklist | 機密掃描擋合併；假密鑰樣式可失敗 | SEC-011 |
| CHK-SEC-012 | checklist | 演示基底 HTTPS；localhost 標非演示 | SEC-012 |
| CHK-SEC-015 | checklist | 演示環境 `/.git/HEAD` 不可取得 | SEC-015 |
| CHK-SEC-017 | checklist | DS1 非公網、最小權限、憑證不進版控 | SEC-017 |

---

## 7. 安全具名 TC 對照（關鍵 SEC）

| SEC | 主要案例 ID |
|---|---|
| SEC-001 | UT-VAL-02～05、UT-APP-02、UT-HC-02、E2E-02、CT-CREATE-02 |
| SEC-003 | UT-ERR-01、UT-ERR-02、CT-ERR-01 |
| SEC-004 | UT-REPO-01、IT-STORE-01、CHK-NFR-005 |
| SEC-006 | UT-GEN-01、UT-GEN-02 |
| SEC-007 | UT-RL-01、UT-RL-02、UT-HC-03、E2E-05 |
| SEC-013 | UT-VAL-07、UT-HR-03、E2E-06 |
| SEC-014 | UT-RED-01、UT-HR-01、CT-REDIR-01、E2E-04 |
| SEC-016 | UT-CLK-01、IT-STORE-02、IT-STORE-03 |

其餘 SEC：見 §6 CHK。

---

## 8. REQ 覆蓋最低要求勾稽

| REQ | ≥1 TC | 代表 ID（含否定／邊界） |
|---|---|---|
| REQ-001 | ✅ | UT-APP-01、CT-CREATE-01、E2E-01 |
| REQ-002 | ✅ | UT-RED-01、CT-REDIR-01、E2E-01 |
| REQ-003 | ✅ | UT-RED-01、IT-STORE-02 |
| REQ-004 | ✅ | UT-APP-04、CT-STATS-01 |
| REQ-005 | ✅ | UT-VAL-02～03、E2E-02 |
| REQ-006 | ✅ | UT-VAL-06～07、E2E-06 |
| REQ-007 | ✅ | UT-APP-02、CT-CREATE-02 |
| REQ-008 | ✅ | UT-RED-02、E2E-03 |
| REQ-009 | ✅ | UT-ERR-02、CT-ERR-01 |
| REQ-010 | ✅ | UT-APP-03、UT-GEN-02、IT-STORE-01 |
| REQ-011 | ✅ | UT-VAL-04～05、E2E-02 |
| REQ-012 | ✅ | IT-FLOW-03 |

---

## 9. 數量粗估

| 類型 | 數量 |
|---|---:|
| UT | 34 |
| IT | 6 |
| CT | 6 |
| E2E | 6 |
| CHK | 17 |
| **合計（具名 ID）** | **約 69** |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G3：案例目錄初版；公開方法 UT 100%；REQ／關鍵 SEC 覆蓋 |
