---
文件：API 設計說明
版本：v0.1.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-19
對應需求：REQ-001～REQ-012, NFR-004, SEC-002, SEC-007, SEC-008, SEC-012
專案代號：SHORTURL
---

# API 設計說明

## 1. 版本策略

- 路徑前綴 `/api/v1/` 用於 JSON API（建立、統計）。
- 短碼導向 `GET /{shortCode}` 置於根路徑，便於短網址分享（非 `/api/v1` 下）。
- 破壞性變更升主版（v2）並走 CR；欄位新增採相容擴充。

## 2. 無帳號／無認證（OUT-03）

- 全部端點 **無** API Key／OAuth／Session。
- 依賴速率限制（SEC-007）與輸入驗證降低濫用；殘餘風險見 SEC-002。

## 3. 錯誤原則（REQ-007～009、SEC-003）

統一本文：

```json
{ "error": { "type": "invalid_url|not_found|rate_limited|internal_error", "message": "..." } }
```

| type | HTTP | 語意 |
|---|---|---|
| invalid_url | 400 | 長網址驗證失敗或短碼格式非法 |
| not_found | 404 | 短碼未登記 |
| rate_limited | 429 | 超過來源速率門檻 |
| internal_error | 500 | 未預期錯誤（訊息不洩漏內部） |

禁止：堆疊、內部路徑、SQL、原始例外字串、框架除錯頁。

## 4. 速率限制（SEC-007，應該）

| 端點 | 建議門檻 |
|---|---|
| POST /api/v1/urls（建立） | ≤ **30** 次／分鐘／來源 |
| GET /{shortCode}（導向） | ≤ **120** 次／分鐘／來源 |
| GET …/stats | MVP 可不強制；若實作建議 ≤ 60／分／來源 |

來源識別：連線 IP，或反向代理轉傳之**可信**客戶端位址。超限：429＋`rate_limited`，不建立、不導向、不計數。

## 5. 安全標頭期望（SEC-008）

| 標頭 | 期望 |
|---|---|
| Strict-Transport-Security | 演示／對外 HTTPS 時：max-age ≥ 31536000 |
| X-Content-Type-Options | nosniff（所有經應用之回應，含 API／錯誤 JSON） |
| Access-Control-Allow-Origin | 若啟用 CORS，不得對敏感回應濫用未驗證 `*` |
| Content-Security-Policy | `default-src 'none'; frame-ancestors 'none'`（實作：`internal/httpapi/security_headers.go`） |
| X-Frame-Options | `DENY` |
| Referrer-Policy | `no-referrer` |
| Cache-Control | JSON 回應 `no-store` |
| Server／框架版本 | 能關閉則不揭露細節 |
| Location（302） | **僅**庫內已驗證 long_url（SEC-014） |

## 6. HTTPS（SEC-012）

- 演示／對外入口應 HTTPS（TLS 1.2+）。
- 本機 localhost HTTP 開發須標「非演示對外」，不得宣稱滿足 SEC-012。

## 7. 開放重新導向殘餘（SEC-002）

MVP 以 http(s) 協定與格式白名單為主控；**不**實作預覽頁（OUT-10）。攻擊者仍可能縮短已通過格式驗證的惡意 https 連結——殘餘風險須於 G2／G5 風險表列示。

## 8. 與 REQ／SEC 對照（摘要）

| API | REQ | SEC |
|---|---|---|
| POST /api/v1/urls | 001、005、007、010～012 | 001、006、007 |
| GET /{shortCode} | 002、003、006、008 | 007、013、014 |
| GET …/stats | 004、006、008 | 013 |
| 錯誤格式 | 007～009 | 003 |

## 9. 長網址上限

**L = 2048**（REQ-011）；OpenAPI `maxLength: 2048`。以 JSON 解碼後之 Unicode 字元計（`\uXXXX` 跳脫與 surrogate pair 各計 1 字元）。

## 10. 短碼格式

建議（ADR-003）：`^[A-Za-z0-9]{8}$`；與 path parameter pattern 一致。

## 11. 傳輸層限制（DR-G4-01）

| 限制 | 值 | 超過時 |
|---|---|---|
| 建立請求本文（POST /api/v1/urls） | 32768 bytes（原始位元組，含 JSON 跳脫與空白） | 400 `invalid_url`（標準錯誤信封；不回 413）。被限流請求先回 429，不讀本文 |
| 請求列＋請求標頭（所有端點） | 請求列（含路徑與查詢字串）＋全部標頭合計 **16384 B 以下保證接受**（`MaxHeaderBytes`）；16384～約 20480 B 不保證 | 合計超過約 20480 B（16384＋net/http 4096 B 緩衝），或標頭值（含 Host）超過 **500 個**（net/http `DefaultMaxHeaderValueCount`）→ **431** Request Header Fields Too Large |

431 由 Go net/http 伺服器在讀取請求列與標頭階段、進入應用（router／ErrorMapper／安全標頭中介層）之前直接回應：`Content-Type: text/plain; charset=utf-8`，本文為純文字 `431 Request Header Fields Too Large`。**不**使用 `{error:{type,message}}` 信封、**不屬** §3 四類錯誤。源站由 net/http 直接回應，不保證帶有 §5 安全標頭；部署邊緣可能加上；用戶端不得依賴這些標頭的有無；用戶端應依狀態碼處理，勿解析 JSON。設計依據：方法規格 v0.2.2 §13.4、§13.4.1。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | 錯誤／限流／標頭／殘餘風險 |
| v0.1.1 | 2026-10-19 | 設計部 | DR-G4-01（PR #17，main `433140c`）：新增 §11 傳輸層限制（建立本文 32768 B→400 invalid_url；標頭 16384 B→431 text/plain 非信封）；§9 補長度以解碼後字元計；OpenAPI 0.1.0→0.1.1（requestBody 說明、各 operation 431、info「傳輸層限制」）；G4-PR20 S1：§11 標頭上限改為請求列＋標頭總量（門檻約 20480 B）；S10：§5 補列實際送出之 CSP、X-Frame-Options、Referrer-Policy、Cache-Control（431 不帶，見 §11）；SCR-043：431 安全標頭改為「不保證帶有、邊緣可能加上、用戶端不得依賴」；SCR-044：標頭上限改為「16384 B 以下保證接受；超過約 20480 B 或標頭值超過 500 個回 431」 |
