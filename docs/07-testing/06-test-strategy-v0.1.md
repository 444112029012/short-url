---
文件：測試策略
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-017
對應規格：docs/02-architecture/06-method-specs-v0.1.md；docs/04-api/openapi.yaml；docs/06-security/
專案代號：SHORTURL
關卡：G3 測試策略與 DoD
---

# 測試策略：短網址服務（SHORTURL）

> 零現金預算；開源工具手寫。技術棧（D-04）：**Go＋chi／SQLite／base62×8 CSPRNG**。  
> 穩定錯誤碼：`invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`。長網址上限 **L＝2048**。

---

## 1. 目標與範圍

| 項目 | 說明 |
|---|---|
| 目標 | 以可追溯 TC／檢查表驗證 REQ／NFR／SEC；公開方法皆有單測計畫；三 API 有契約＋E2E 主路徑 |
| 範圍內 | 單元、整合（含 SQLite／記憶體雙件）、OpenAPI 契約、E2E 主／否定路徑、安全關鍵 TC、NFR 抽樣／檢查表、CI 資安門檻檢查表 |
| 範圍外（本輪排除） | 多區／水平擴展負載（NFR-003 明確不承諾）；預覽頁／帳號體系（OUT）；正式滲透外包；跨實例限流一致性 |
| 對齊來源 | 方法規格公開方法清單；OpenAPI `createShortUrl`／`redirectShortCode`／`getUrlStats`；SEC v0.2；CI 門檻文件 |

---

## 2. 測試金字塔

```
        ┌─────────┐
        │  E2E    │  少：三 API 主路徑＋關鍵否定（建立失敗／not_found）
        ├─────────┤
        │ Contract│  中：OpenAPI 回應碼／schema／錯誤 type enum
        ├─────────┤
        │Integrat.│  中：Store（SQLite）＋ create→redirect→stats 流程
        ├─────────┤
        │  Unit   │  多：各公開方法 Given／When／Then；Validator／ErrorMapper 為主
        └─────────┘
```

| 層級 | 比例建議（案例數） | 工具建議（開源） | 執行時機 |
|---|---|---|---|
| Unit | ~55–65% | Go `testing`、`testify`（可選）、table-driven | 每次 commit／PR |
| Integration | ~15–20% | `testing`＋SQLite（`:memory:` 或暫存檔）、httptest | PR／合併前 |
| Contract | ~10–15% | OpenAPI 對照（手寫 assert 或 `kin-openapi`／spectral 等開源）；httptest | PR |
| E2E | ~5–10% | 本機／CI 起服務＋`curl`／Go HTTP client | 合併前＋演示前 |
| 檢查表 | 另計 | Markdown 勾選（NFR／SEC CI／部署） | G3～G5 關卡 |

---

## 3. 層級範圍細則

### 3.1 單元測試（Unit）

- **必測公開方法**（各至少 1 個 UT 計畫 ID，見 `07-test-cases` §公開方法覆蓋）：  
  `HttpApiAdapter`（handleCreate／handleRedirect／handleStats）、`ShortUrlApplicationService`（createShortUrl／getStats）、`RedirectService.redirect`、`UrlValidator`（validateLongUrl／validateShortCode）、`ShortCodeGenerator.generate`、`UrlRepository`（save／findByShortCode／exists）、`ClickCounter`（incrementAtomic／getCount）、`ErrorMapper.toHttpResponse`、`RateLimitGuard`（checkCreate／checkRedirect）、`ObservabilityHooks`（Success／Failure 可合併輕量 UT）。
- **InMemoryOrDbUrlStore**：介面行為以整合測為主；單元可用記憶體 fake／stub。
- 依賴注入以 interface mock／fake；禁止單測打真實公網。

### 3.2 整合測試（Integration）

- SQLite schema：無 `visitor_ip`／`user_agent`（SEC-004）。
- `save`／`findByShortCode`／`exists`／`incrementAtomic`／`getCount`：**參數化**、唯一約束、原子 +1（SEC-016）。
- 應用流程：建立成功 → 導向 302＋計數 +1 → stats 一致；失敗路徑不落庫／不計次。

### 3.3 契約測試（Contract）

- 對齊 `docs/04-api/openapi.yaml`：狀態碼、必要欄位、`ErrorResponse.type` 四碼 enum、`ShortCodePath` pattern、`maxLength: 2048`。
- 三 operationId 各 ≥1 契約成功案例；各端點至少覆蓋一種錯誤回應形狀。

### 3.4 E2E

- 主路徑：POST 建立 → GET 導向 → GET stats（NFR-002／演示）。
- 否定：非法 URL 建立、不存在短碼、非法短碼格式、（可選）限流 429。
- Location 必須等於庫內 long_url；附加 query 不得改寫（SEC-014）。

---

## 4. 環境與資料

| 環境 | 用途 | 備註 |
|---|---|---|
| 單測 | 純記憶體／fake | 無 I/O 或最少 I/O |
| 整合／契約 | SQLite `:memory:` 或 CI 暫存 DB | 每次測試隔離；遷移腳本可重入 |
| E2E／演示 | localhost:8080（dev）或演示 HTTPS | localhost **不宣稱** SEC-012；演示須 HTTPS（SEC-012） |
| NFR-001 抽樣 | **約定基線**（建議關閉 G1-QA-O5）：單實例、無人工負載、暖機後；CPU／RAM 於測試報告註記 | P95 ≤ 2s；建立／導向各 ≥10 次；冷啟動第一次可文件標明排除 |

---

## 5. NFR 策略

| NFR | 測法 |
|---|---|
| NFR-001 | 抽樣腳本＋報告表（P95）；見 TC-NFR-001 |
| NFR-002 | E2E 演示檢查表：建立→導向→stats 完整 1 次 |
| NFR-003 | 文件審查：單實例假設已寫入（C4／ADR-002） |
| NFR-004 | 檢查表：ASVS／SEC／CI 門檻／G5 高風險（跨關卡） |
| NFR-005 | 模型／整合：無 IP／UA 欄；對照隱私文件 |
| NFR-006 | ObservabilityHooks UT＋演示指出日誌／指標證據 |
| NFR-007 | `docs/gates/` 關卡紀錄檢查表 |
| NFR-008 | `docs/README`／各章狀態檢查表（含 UI N/A） |

---

## 6. 安全測試策略（SEC）

| 優先 | SEC | 策略 |
|---|---|---|
| 關鍵 TC | 001、003、004、006、007、013、014、016 | 具名 TC（見案例目錄）；自動化為主 |
| 檢查表／抽查 | 002、005、008、009～012、015、017 | 文件＋環境／CI／碼審；G3～G5 举证 |
| 殘餘 | SEC-002／RR-001 | 已接受（D-06）；TC 驗證白名單＋殘餘聲明存在即可 |

注入／原子（SEC-016）：特製字元不得改變 SQL 語意；並發導向後 click_count 與成功次數一致。

---

## 7. 與 OpenAPI／方法規格對齊原則

1. HTTP 行為以 OpenAPI 為契約準繩；領域後置條件以方法規格為準。  
2. 錯誤僅允許四穩定碼；短碼格式失敗對外為 `invalid_url`。  
3. 追溯矩陣「TC 欄」填入本策略對應之 TC／UT／CHK ID（見矩陣 v0.3）。  
4. 新增公開方法或 operation 須同步補 UT／契約／追溯（CR）。

---

## 8. 缺陷與退出準則（摘要）

| 級別 | 定義（MVP） | G3／合併前 |
|---|---|---|
| Blocker | 主路徑不可用、資料損毀、SEC 必須項失敗 | 不得合併／不得送後續實作驗收 |
| Major | REQ 必須項失敗、契約嚴重不符 | 須修或書面接受 |
| Minor | 文案、非關鍵標頭、建議項 | 可追蹤 |

**測試退出（DoD 細節見 `08-definition-of-done`）**：公開方法 UT 計畫 100% 有對應案例 ID；REQ-001～012 各 ≥1 TC；關鍵 SEC 有 TC／檢查；三 API 契約＋E2E 主路徑各 ≥1；追溯無空白必填 TC。

---

## 9. 風險與假設

| 風險／假設 | 緩解 |
|---|---|
| 研發任務清單尚未產出 | 品保先以 WP／任務骨架對帳（`09-dev-task-acceptance-map`），標「待研發對齊」 |
| NFR-001 環境基線未定 | 本策略 §4 暫定建議基線；維運／設計補正式文件後更新報告 |
| Observability 事件鍵未齊（G1-QA-O7） | Hooks 介面級 UT 先過；事件鍵清單列建議項 |
| 限流為進程內單實例 | 測試不宣稱跨實例；與 NFR-003 一致 |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G3：測試策略初版（金字塔、範圍、工具、NFR／SEC、排除項） |
