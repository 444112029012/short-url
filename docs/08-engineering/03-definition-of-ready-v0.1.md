---
文件：Definition of Ready（DoR）
版本：v0.1
日期：2026-10-07
狀態：已取代（請改用 v0.2）
負責角色：研發部
作者：研發部
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
對齊：已由 v0.2 取代；正式 DoD 見 docs/07-testing/08-definition-of-done-v0.1.md
---

# Definition of Ready（DoR）：SHORTURL

> ⚠️ **本文件已由 [`03-definition-of-ready-v0.2.md`](./03-definition-of-ready-v0.2.md) 取代**（2026-10-07）。v0.2 已對齊品保正式 DoD；下列「尚未發布 DoD／待品保確認」為歷史表述，**請勿再依本稿開工**。

> **目的**：任務「可開工」的入場條件。  
> **對齊說明**：品保於 `docs/07-testing/` **尚未**發布獨立正式 **DoD** 文件。本 DoR 自 G2 設計可測審查、追溯矩陣、G2 品保前置與方法規格／OpenAPI **推導**「對齊品保期望的入場條件」，並標出 **待品保確認** 項。  
> 任務完成後之 **DoD** 仍以品保後續文件為準；研發自測最低線見 §4。

## 依據文件

| 文件 | 路徑 |
|---|---|
| 任務拆解 | `docs/08-engineering/01-task-breakdown-v0.1.md` |
| 分支／Commit | `docs/08-engineering/02-branch-commit-conventions-v0.1.md` |
| 方法規格 | `docs/02-architecture/06-method-specs-v0.1.md` |
| OpenAPI | `docs/04-api/openapi.yaml` |
| 需求／範圍 | `docs/01-requirements/02-requirements-spec-v0.1.md`、`03-scope-v0.1.md` |
| SEC v0.2 | `docs/06-security/01-sec-requirements-v0.2.md` |
| 追溯矩陣 | `docs/07-testing/02-traceability-matrix-v0.1.md` |
| 設計可測審查 | `docs/07-testing/04-design-testability-review-v0.1.md` |
| G2 品保前置 | `docs/07-testing/05-gate-G2-qa-precheck-v0.1.md` |
| CI 門檻 | `docs/06-security/03-ci-security-gates-v0.1.md` |
| G2 通過 | `docs/gates/gate-G2-2026-10-07-recheck-01.md` |

---

## 1. 任務可開工檢查清單（DoR）

每一 `ENG-xxx` 開工前，負責人須全部勾選（或標 N/A＋理由）：

### 1.1 規格與追溯連結

- [ ] **任務 ID 穩定**：已列於 `01-task-breakdown-v0.1.md`，可被分支／PR 引用
- [ ] **方法規格列齊**：對應「類別.方法」已指向方法規格具體列（或標基礎設施／文件例外）
- [ ] **OpenAPI 對齊**（若為 API 任務）：operationId／路徑／狀態碼／schema 已引用，無與規格矛盾
- [ ] **REQ／NFR ID 真實存在**：僅使用規格中之 REQ-001～012、NFR-001～008；無對應則標「—」並註記
- [ ] **SEC ID 真實存在**：僅使用 SEC-001～017；不適用標「—」
- [ ] **範圍內**：不實作 OUT-01～11 未納入項（如 TTL＝OUT-05）

### 1.2 驗收可測

- [ ] **驗收標準可觀察**：Given／When／Then 或等價可測敘述已寫在任務卡／拆解表
- [ ] **錯誤碼可斷言**（若涉錯誤路徑）：僅使用四穩定碼 `invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`
- [ ] **與品保可測結論一致**：不依賴「不可測」設計缺口（G2 設計可測已判定核心契約可測）
- [ ] **測試策略入口已知**：單元（ENG-018）／整合（ENG-019）／檢查表（CI／部署類 SEC）何者適用已標明

### 1.3 相依與環境

- [ ] **相依任務就緒**：拆解表「相依」列之 ENG 已合併或本任務明確允許平行介面開發
- [ ] **技術決策已拍板**：D-04（Go＋chi／SQLite／base62×8 CSPRNG）無未決替代方案阻礙本任務
- [ ] **本機可建置假設成立**：Go 工具鏈與零預算約束下無付費服務依賴
- [ ] **設定鍵已知**（若需要）：BASE_URL、DB 路徑、限流門檻等已於 ENG-002 或文件定義

### 1.4 安全／SEC 已知

- [ ] **相關 SEC 已讀**：任務表列出之 SEC 條文驗證方式已知
- [ ] **SEC-016／017**（若動 Store）：參數化、禁拼接、原子 +1、欄位最小化、憑證不進版控已納入實作注意
- [ ] **SEC-004／005**：無 IP／UA 明細欄；日誌無 Authorization／Cookie 原文
- [ ] **殘餘風險**：SEC-002／RR-001 已 D-06 接受者，不重新擴大為預覽頁需求
- [ ] **CI 門檻**：若任務改 CI，已對齊 `03-ci-security-gates-v0.1.md` 阻擋語意

### 1.5 無阻擋開放問題

- [ ] **無未關閉必改**擋本任務（G2-R1／G2-QA-R1 已關）
- [ ] **開放建議項不阻開工**：G1-QA-O5（環境基線）、G1-QA-O7（事件鍵）已知為建議；若本任務強依賴則先開協作單
- [ ] **估時與角色已填**：人時與研發／維運／安全協作標記完整

---

## 2. 對齊品保期望的「完成側」最低線（非正式 DoD）

> **待品保確認**：以下由研發自 G2 可測／追溯推導，供 PR 自檢；**不取代**品保正式 DoD／TC 簽核。

| # | 研發自檢項（對齊品保期望） | 來源推導 | 品保確認 |
|---|---|---|---|
| D1 | 行為符合任務驗收與方法規格前後置條件 | 方法規格；可測審查 §2.1 | 待確認 |
| D2 | API 任務通過對應 OpenAPI 狀態碼／body／Location 斷言 | OpenAPI；追溯 API 欄 | 待確認 |
| D3 | 錯誤路徑僅四穩定 `type`，無堆疊／路徑／SQL 外洩 | ErrorMapper；SEC-003；REQ-009 | 待確認 |
| D4 | 短碼 `^[A-Za-z0-9]{8}$`、L=2048 行為與文件一致 | ADR-003；REQ-011 | 待確認 |
| D5 | Store 變更：參數化＋原子計數；schema 無 visitor_ip／UA | SEC-016／017／004 | 待確認 |
| D6 | 相關 `go test`（單元／整合鉤子）綠 | ENG-018／019；追溯 TC 仍 TBD | 待確認（TC 編號後補） |
| D7 | PR 含 ENG＋REQ／SEC；CI required checks 綠 | 分支規範；SEC-009～011 | 待確認 |
| D8 | 不擴大 OUT 範圍；殘餘僅文件化已接受者 | 範圍；D-06 | 待確認 |

品保正式 DoD 建議涵蓋：追溯矩陣 TC 欄由 TBD→實號、NFR-001 環境基線（O5）、可觀測事件鍵（O7）、G5 前高風險關閉證明——**非本 DoR 宣告已完成**。

---

## 3. 各類任務適用速查

| 任務類型 | 額外 DoR 重點 |
|---|---|
| 領域／應用服務 | 方法規格例外／後置；REQ 驗收 |
| Store／SQL | SEC-016／017／004；資料模型 §6 |
| HTTP Adapter | OpenAPI operationId；限流與 Hooks |
| Generator | SEC-006；ADR-003 熵 |
| CI | SEC-009～011 門檻表；零預算工具 |
| 文件／README | SEC-012 環境區分；無假承諾 HTTPS |

---

## 4. DoR 簽核方式（G3）

| 步驟 | 角色 | 動作 |
|---|---|---|
| 1 | 研發 | 任務卡勾選 §1 清單 |
| 2 | 品保 | 確認 §2 與未來 DoD 無衝突（**待確認**項關閉或接受） |
| 3 | 安全 | 抽查含 SEC 之任務（尤其 006／007／016／017／009～011） |
| 4 | 審查／總協調 | G3 關卡將本文件＋任務拆解＋分支規範列為開工就緒證據 |

**未通過 DoR 不得開實作 PR 合併至 `main`。**

---

## 5. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 研發部 | G3 初稿；自品保可測／追溯推導；標待品保確認 DoD |
