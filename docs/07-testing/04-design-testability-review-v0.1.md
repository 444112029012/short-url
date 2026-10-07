---
文件：G2 設計可測性審查
版本：v0.1.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-017
對應規格：docs/02-architecture/、docs/04-api/、docs/06-security/01-sec-requirements-v0.2.md
專案代號：SHORTURL
關卡：G2 設計＋威脅模型
---

# G2 設計可測性審查：短網址服務（SHORTURL）

> 目的：確認設計交付（方法規格、OpenAPI、設計圖、ADR）與 SEC／威脅模型是否足以推導可測案例（TC／檢查表）。  
> 判定：**可測**＝現有設計即可寫 TC／檢查項；**有條件可測**＝判準客觀但依賴尚缺之設計補件或 G3 實作／環境；**不可測**＝設計本身無客觀 I/O／例外／契約。  
> 本輪 **無不可測**；有條件項與必改見 §6～§7。

---

## 1. 統計摘要

| 檢核面向 | 判定 | 說明 |
|---|---|---|
| 方法規格表 | **可測** | 公開方法含輸入／回傳／例外／前置／後置；四穩定錯誤碼齊 |
| OpenAPI 契約 | **可測** | 建立／導向／查次數齊；L=2048；短碼 `^[A-Za-z0-9]{8}$` |
| 四錯誤碼 | **可測** | 文件穩定集合恰為四碼，可斷言 |
| 設計圖（C4／類別／循序） | **可測** | SEQ-01～07 覆蓋主路徑與主要否定路徑 |
| ADR（001～003） | **可測** | D-04 已採納；與 OpenAPI／方法規格一致 |
| SEC-001～015 | **可測／有條件** | 多數可自設計／API 推導；CI／部署類屬 G3 檢查表 |
| SEC-016～017 | **可測** | 設計部 v0.1+R1 已補追溯／方法規格／資料模型 §6／C4；**G2-QA-R1 已關閉** |
| 威脅模型／殘餘 | **不阻擋可測結論** | RR-001／TM-002 已 D-06 接受；TC 仍可寫（白名單／殘餘聲明檢查） |

| 判定 | 數量（面向列） |
|---|---:|
| 可測 | 6 |
| 可測／有條件（CI／部署，不阻送審） | 1（SEC-001～015 列） |
| 有條件可測（阻擋） | 0（原 SEC-016／017 已關閉） |
| 不可測 | 0 |

---

## 2. 逐項勾稽

### 2.1 方法規格表（`02-architecture/06-method-specs-v0.1.md`）

| 檢查點 | 結果 | 可測？ | 缺口／備註 |
|---|---|---|---|
| 欄位含前置／後置、I/O、例外 | ✅ | 可測 | 足以寫 Given／When／Then |
| 對應 REQ／NFR／SEC | ✅ | 可測 | 公開方法列齊（§12 自檢） |
| 例外路徑與錯誤碼 | ✅ | 可測 | `invalid_url`／`not_found`／`rate_limited`／`internal_error` |
| 長網址上限 L | ✅ | 可測 | `maxLength 2048`／`長度>2048` → invalid_url |
| 短碼規則 | ✅ | 可測 | `^[A-Za-z0-9]{8}$`；對齊 ADR-003 |
| 原子計數／失敗不計次 | ✅ | 可測 | RedirectService／ClickCounter 後置條件明確 |
| SEC-016 參數化查詢 | ✅ | 可測 | 方法規格明文參數化／禁拼接／原子 UPDATE；資料模型 §6；Component §5.1 |
| SEC-017 最小權限／憑證 | ✅ | 可測 | 資料模型 §6、Component §5.1、Container §5 已列最小權限／非公網／憑證不進版控 |

**結論**：核心 API／領域行為可測；SEC-016／017 已可由設計推導檢查項（G2-QA-R1 已關）。

### 2.2 OpenAPI／API notes／README（`04-api/`）

| 檢查點 | 結果 | 可測？ | 證據 |
|---|---|---|---|
| `openapi.yaml` 存在 | ✅ | — | `docs/04-api/openapi.yaml` |
| `01-api-notes-v0.1.md` | ✅ | — | 錯誤／限流／標頭／殘餘 |
| `README.md` | ✅ | — | 索引＋交接單 |
| 建立 `POST /api/v1/urls` | ✅ | 可測 | `operationId: createShortUrl`；201／400／429／500 |
| 導向 `GET /{shortCode}` | ✅ | 可測 | `operationId: redirectShortCode`；302＋Location；400／404／429／500 |
| 查次數 `GET …/stats` | ✅ | 可測 | `operationId: getUrlStats`；200／400／404／500 |
| 長網址上限 L 寫死 | ✅ | 可測 | `maxLength: 2048`；notes §9 **L = 2048** |
| 短碼與 ADR-003 一致 | ✅ | 可測 | `pattern: '^[A-Za-z0-9]{8}$'`；min／maxLength 8 |
| 錯誤 `type` enum | ✅ | 可測 | 恰四碼（見 §2.3） |
| 速率門檻文件化 | ✅ | 可測 | 建立 ≤30／分；導向 ≤120／分（info／notes） |

**結論**：OpenAPI 齊且與方法規格、ADR-003、REQ-011 一致；可直接據此起草 API TC。

### 2.3 四錯誤碼（穩定錯誤碼集合）

| type | HTTP | 出現處 | 可斷言？ |
|---|---|---|---|
| `invalid_url` | 400 | 方法規格、OpenAPI enum、api-notes、ErrorMapper、SEQ-02／06 | ✅ |
| `not_found` | 404 | 同上；SEQ-04 | ✅ |
| `rate_limited` | 429 | 同上；SEQ-07 | ✅ |
| `internal_error` | 500 | 同上 | ✅ |

> 文件實際穩定集合**正好四個**（非「約四個」）。短碼格式失敗對外統一 `invalid_url`（SEQ-06 約定）；若未來拆 `invalid_short_code` 須 CR——TC 以現行四碼為準。

**結論**：**可測**；G1-QA-O3 已因設計關閉。

### 2.4 設計圖與 ADR

| 交付物 | 路徑 | 結果 | 可測對齊 |
|---|---|---|---|
| C4 Context | `01-c4-context-v0.1.md` | ✅ | 三核心路徑；EE1 |
| C4 Container | `02-c4-container-v0.1.md` | ✅ | Go+chi／SQLite；單實例 NFR-003 |
| C4 Component | `03-c4-component-v0.1.md` | ✅ | 元件名可填追溯「API／元件」欄 |
| 類別圖 | `04-class-diagram-v0.1.md` | ✅ | ErrorType 四碼；公開方法與規格一致 |
| 循序圖 SEQ-01～07 | `05-sequence-diagrams-v0.1.md` | ✅ | 成功／驗證失敗／導向＋計數／not_found／stats／非法短碼／限流 |
| 需求—設計追溯 | `07-req-design-traceability-v0.1.md` | ✅ | REQ／NFR／SEC-001～017 有列（含 016／017） |
| ADR-001～003 | `adr/` | ✅ **已採納**（D-04） | Go+chi；SQLite；base62×8 CSPRNG≈47.6 bit |
| UI／UX | `05-ui-ux/README.md` | ✅ N/A | 理由充分（API＋導向、無獨立 UI）→ NFR-008 |
| 資料模型／DFD | `03-data/` | ✅／⚠ | L=2048、禁 IP／UA；正文仍有「持久化待 ADR-002」殘句（建議修） |

**結論**：設計圖足以對應 TC 場景；SEC-016／017 追溯已補齊（G2-QA-R1 已關）。

### 2.5 SEC／威脅模型可測對齊

來源優先：`01-sec-requirements-v0.2.md`（含 SEC-016／017）；威脅：`04-threat-model-v0.1.md`；殘餘：`05-residual-risk-register-v0.1.md`（RR-001 **已接受**，D-06）；安全前置：`06-g2-security-pregate-v0.1.md`。

| SEC | 自設計／API 推導測試？ | 建議測法 | 備註 |
|---|---|---|---|
| SEC-001 | ✅ 可測 | API：空／非 URL／ftp／js／L／L+1 | UrlValidator；OpenAPI |
| SEC-002 | ✅ 可測（文件＋抽查） | 檢查殘餘聲明；抽查僅 http(s) | RR-001／D-06 已接受；不阻可測 |
| SEC-003 | ✅ 可測 | 錯誤路徑本文無堆疊／路徑／SQL | ErrorMapper；OpenAPI |
| SEC-004 | ✅ 可測 | 模型／遷移無 visitor_ip／UA | 資料模型 §禁止欄 |
| SEC-005 | ✅／有條件 | 程式＋日誌抽查 | ObservabilityHooks 約束已寫；實作於 G3 |
| SEC-006 | ✅ 可測 | ADR 熵＋碼審 CSPRNG；抽樣無序列 | ADR-003；ShortCodeGenerator |
| SEC-007 | ✅ 可測 | 超門檻 429＋rate_limited | RateLimitGuard；notes |
| SEC-008 | ✅／有條件 | 標頭擷取 | api-notes／OpenAPI info；演示環境 |
| SEC-009～011 | 有條件 | CI 設定檢查表 | 設計標「09-cicd」；屬 G3 |
| SEC-012 | 有條件 | TLS 連線／環境文件 | servers 已區分 demo／localhost |
| SEC-013 | ✅ 可測 | 非法短碼 400、不導向、不計次 | pattern；SEQ-06 |
| SEC-014 | ✅ 可測 | Location＝庫值；參數無法改寫 | RedirectService；SEQ-03 |
| SEC-015 | 有條件 | 部署路徑 `/.git` | 設計註記；G3／維運 |
| SEC-016 | **可測** | 碼審參數化；並發計數；方法規格 SQL 例 | 設計已明示；R1 關 |
| SEC-017 | **可測**（部署類於 G3 举证） | 環境／權限／機密掃描 | 設計已明示約束；G3 举证 |

**威脅模型影響**：TM／RR 不改變「可從契約寫 TC」之結論。RR-001 接受後，SEC-002 驗收以「白名單＋書面殘餘」為準，不要求預覽頁（OUT-10）。

---

## 3. 檢查重點問答（審查部可直接引用）

| # | 問題 | 答案 |
|---|---|---|
| 1 | OpenAPI 是否定義建立／導向／查次數？ | **是**（`createShortUrl`／`redirectShortCode`／`getUrlStats`） |
| 2 | 長網址上限 L 是否寫死？ | **是，L = 2048**（OpenAPI `maxLength`＋api-notes §9） |
| 3 | 短碼規則是否與 ADR-003（base62×8）一致？ | **是**（`^[A-Za-z0-9]{8}$`；CSPRNG；熵≈47.6 bit） |
| 4 | 四錯誤碼是否文件化且可斷言？ | **是**；穩定集合恰為 `invalid_url`／`not_found`／`rate_limited`／`internal_error` |
| 5 | 方法規格是否足以寫 TC？ | **是**（前置／後置／I/O／例外齊；含 SEC-016／017） |
| 6 | SEC-001～017 是否可從設計／API 推導？ | **001～017 可推導**；009～012／015／017 部署举证偏 G3 有條件執行，不阻擋設計可測 |
| 7 | 威脅模型／殘餘是否影響可測結論？ | **不阻擋**；D-06 已接受 RR-001；TC 仍可寫 |
| 8 | G1-QA-O1～O7 哪些已因設計關閉？ | 見 §4 |

---

## 4. G1 開放項關閉狀態

| 編號 | G1 狀態 | G2 判定 | 證據 |
|---|---|---|---|
| G1-QA-O1 | 建議 | **已關閉** | OpenAPI 導向為 **GET** `/{shortCode}` |
| G1-QA-O2 | 追蹤 | **已關閉** | ADR-003 已採納 base62×8；OpenAPI pattern 一致 |
| G1-QA-O3 | 建議 | **已關閉** | 四穩定錯誤碼於 OpenAPI enum／notes／方法規格／ErrorMapper |
| G1-QA-O4 | 追蹤 | **已關閉** | L=2048 寫死於 OpenAPI 與資料模型 |
| G1-QA-O5 | 追蹤 | **仍開放（建議）** | 單實例／P95≤2s 已於 Container 標明；**機器規格／演示環境基線**仍建議於測試計畫補（不阻 G2 送審） |
| G1-QA-O6 | 部分關閉 | **大部分關閉** | ASVS／SEC v0.2／CI 門檻文件已存在；品保矩陣本輪升至含 SEC-016／017 與設計欄 |
| G1-QA-O7 | 追蹤 | **部分關閉** | ObservabilityHooks 已列成功／失敗事件；維運「事件 ID／訊息鍵」清單仍建議於 G3 前補齊 |

---

## 5. REQ／NFR 設計可測摘要（相對 G1）

| 編號 | G1 | G2（設計後） | 說明 |
|---|---|---|---|
| REQ-001～005、007～010、012 | 可測 | **可測** | 方法規格＋OpenAPI＋SEQ 齊 |
| REQ-006 | 有條件 | **可測** | ADR-003＋pattern |
| REQ-011 | 有條件 | **可測** | L=2048 |
| NFR-001 | 有條件 | **有條件** | 門檻清楚；環境基線建議補（O5） |
| NFR-002～003、005、007～008 | 可測 | **可測** | 文件／演示／模型可查 |
| NFR-004 | 有條件 | **有條件** | SEC／威脅已齊；CI 實作屬 G3 |
| NFR-006 | 有條件 | **有條件** | Hooks 已設計；事件鍵建議補（O7） |

---

## 6. 必改項（G2-QA-R\*）

| 編號 | 問題 | 必改內容 | 負責角色 |
|---|---|---|---|
| **G2-QA-R1** | （原）SEC-016／017 未落入設計 | 設計部已於追溯、方法規格 §6～8、資料模型 §6、Component §5.1、Container §5 補齊可驗證描述 | **設計部** → **已關閉（2026-10-07 品保核對）** |

### 建議項（非必改，不阻擋有條件送審）

| 編號 | 對象 | 建議 |
|---|---|---|
| G2-QA-S1 | 設計 | 資料模型「持久化引擎待 ADR-002」改為 **已採納 SQLite（ADR-002A／D-04）** |
| G2-QA-S2 | 設計／維運 | 補 NFR-001 演示測試環境基線（關閉 G1-QA-O5） |
| G2-QA-S3 | 設計 | 查統計限流：api-notes「可不強制」與 TM-011 對齊——於 notes 或 RateLimitGuard 加一句正式策略 |
| G2-QA-S4 | 企劃 | 規格 §4／開頭 SEC 引用可升註至 `01-sec-requirements-v0.2.md`（含 016／017） |

---

## 7. 一句話結論

**設計契約（方法規格＋OpenAPI＋四錯誤碼＋ADR-003／L=2048＋SEQ）可測；G1-QA-O1～O4 已關；SEC-016／017 已落入設計 → G2-QA-R1 已關閉；品保 G2 設計可測性判定為「可送審」。**

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G2 設計可測性初審；G2-QA-R1；G1-O1～O4 關閉 |
| v0.1.1 | 2026-10-07 | 品保部 | 核對設計補件；**關閉 G2-QA-R1**；結論改可送審 |
