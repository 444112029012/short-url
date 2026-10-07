---
文件：Definition of Ready（DoR）
版本：v0.2
日期：2026-10-07
狀態：Draft／Ready for Review
負責角色：研發部
作者：研發部
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
對齊：品保正式 DoD（`docs/07-testing/08-definition-of-done-v0.1.md`）
---

# Definition of Ready（DoR）：SHORTURL

> **目的**：任務「可開工」的入場條件。  
> **對齊說明**：本 DoR 與品保正式 **Definition of Done**（`docs/07-testing/08-definition-of-done-v0.1.md`）成對使用——**DoR＝開工前**；**DoD＝宣告 Done 前**。  
> 研發任務 ID（ENG-xxx）與驗收 TC 對帳見 `docs/07-testing/09-dev-task-acceptance-map-v0.1.md`。

## 依據文件

| 文件 | 路徑 |
|---|---|
| **正式 DoD** | `docs/07-testing/08-definition-of-done-v0.1.md` |
| 研發任務↔驗收對照 | `docs/07-testing/09-dev-task-acceptance-map-v0.1.md` |
| 任務拆解 | `docs/08-engineering/01-task-breakdown-v0.1.md` |
| 分支／Commit | `docs/08-engineering/02-branch-commit-conventions-v0.1.md` |
| 方法規格 | `docs/02-architecture/06-method-specs-v0.1.md` |
| OpenAPI | `docs/04-api/openapi.yaml` |
| 需求／範圍 | `docs/01-requirements/02-requirements-spec-v0.1.md`、`03-scope-v0.1.md` |
| SEC v0.2 | `docs/06-security/01-sec-requirements-v0.2.md` |
| 追溯矩陣（內容 v0.3；TC 欄已填） | `docs/07-testing/02-traceability-matrix-v0.1.md` |
| 測試案例 | `docs/07-testing/07-test-cases-v0.1.md` |
| 測試策略 | `docs/07-testing/06-test-strategy-v0.1.md` |
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
- [ ] **驗收對照已知**：本任務在 `09-dev-task-acceptance-map` 有對應列（或標文件／骨架例外並註明）

### 1.2 驗收可測

- [ ] **驗收標準可觀察**：Given／When／Then 或等價可測敘述已寫在任務卡／拆解表
- [ ] **錯誤碼可斷言**（若涉錯誤路徑）：僅使用四穩定碼 `invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`
- [ ] **與品保可測結論一致**：不依賴「不可測」設計缺口（G2 設計可測已判定核心契約可測）
- [ ] **主要 TC／UT 已知**：已對照 acceptance map 與追溯矩陣 TC 欄；測試策略入口（單元 ENG-018／整合 ENG-019／CI ENG-020）何者適用已標明

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

## 2. 完成側對齊正式 DoD（開工前須知）

> 任務／PR **宣告 Done** 時，適用區塊須符合品保正式 DoD（`08-definition-of-done-v0.1.md`）全部勾選（或不適用且附理由）。  
> 下表為研發自檢速查，**對應 DoD 區塊編號**；細節以 DoD 本文為準。

| DoD 區塊 | 研發自檢要點（摘要） | DoD 準則 |
|---|---|---|
| §2 程式 | 對齊方法規格／OpenAPI／L=2048／base62×8／參數化 Store／無 IP·UA／ErrorMapper／限流 | P1～P8 |
| §3 測試 | 影響方法有 UT；矩陣 TC 非空白；UT／IT／CT／E2E 綠；否定路徑有覆蓋 | T1～T5 |
| §4 文件 | 契約變更同步 OpenAPI／規格；追溯／案例已更新；繁中；PR 含 REQ／SEC／TC | D1～D4 |
| §5 安全 | 無 SEC 必須項迴歸；無機密進版控；SAST／SCA 門檻；Location 不可改寫；日誌無敏感原文 | S1～S5 |
| §6 CI／維運 | 測試失敗不得無紀錄合併；不擅自放寬資安門檻；localhost≠演示標示 | C1～C3 |
| §7 追溯 | REQ／SEC→設計→TC 可追；**ENG ID ↔ TC** 已對帳（acceptance map） | R1～R3 |
| §8 里程碑 | G4／G5 前另檢：E2E-01、NFR-001／006、高風險 SEC | M1～M4 |
| §9 迷你清單 | PR 描述建議貼 DoD 任務級勾選清單 | 全文 §9 |

**PR 最低自檢（對齊 DoD §9）**：

```text
[ ] 程式：對齊規格／OpenAPI／四錯誤碼／L=2048／base62×8
[ ] 測試：相關 UT/IT/CT/E2E 綠燈；否定路徑有覆蓋
[ ] 文件：契約或追溯若需改已改
[ ] 安全：無機密；掃描未無紀錄關閉
[ ] 追溯：任務 ID ↔ TC 已對帳
```

建議／殘餘（不阻本 DoR 開工，屬 DoD／里程碑追蹤）：NFR-001 環境基線（G1-QA-O5）、可觀測事件鍵（G1-QA-O7）、G5 前高風險關閉證明。

---

## 3. 各類任務適用速查

| 任務類型 | 額外 DoR 重點 | 主要對應 DoD |
|---|---|---|
| 領域／應用服務 | 方法規格例外／後置；REQ 驗收 | P1、T1、T4 |
| Store／SQL | SEC-016／017／004；資料模型 | P4、P5、S1 |
| HTTP Adapter | OpenAPI operationId；限流與 Hooks | P2、P7、S4 |
| Generator | SEC-006；ADR-003 熵 | P3、S1 |
| CI | SEC-009～011 門檻表；零預算工具 | S2、S3、C1～C2 |
| 文件／README | SEC-012 環境區分；無假承諾 HTTPS | D3、C3 |

---

## 4. DoR 簽核方式（G3）

| 步驟 | 角色 | 動作 |
|---|---|---|
| 1 | 研發 | 任務卡勾選 §1 清單 |
| 2 | 品保 | 確認本 DoR 與正式 DoD／acceptance map 無衝突；抽查 ENG↔TC |
| 3 | 安全 | 抽查含 SEC 之任務（尤其 006／007／016／017／009～011） |
| 4 | 審查／總協調 | G3 關卡將本文件＋任務拆解＋分支規範＋DoD 列為開工就緒證據 |

**未通過 DoR 不得開實作 PR 合併至 `main`。**  
**未通過適用 DoD 不得標任務／PR Done。**

---

## 5. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 研發部 | G3 初稿；自品保可測／追溯推導；當時正式 DoD 尚未產出 |
| v0.2 | 2026-10-07 | 研發部 | **對齊品保 DoD**（`08-definition-of-done-v0.1.md`）；去除「待品保確認／尚無 DoD」表述；§2 改為正式 DoD 速查；補 acceptance map／矩陣 v0.3 依據 |
