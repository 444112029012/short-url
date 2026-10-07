---
文件：API 設計說明
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
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
| X-Content-Type-Options | nosniff（API／錯誤 JSON） |
| Access-Control-Allow-Origin | 若啟用 CORS，不得對敏感回應濫用未驗證 `*` |
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

**L = 2048**（REQ-011）；OpenAPI `maxLength: 2048`。

## 10. 短碼格式

建議（ADR-003）：`^[A-Za-z0-9]{8}$`；與 path parameter pattern 一致。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | 錯誤／限流／標頭／殘餘風險 |
