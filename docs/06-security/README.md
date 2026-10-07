---
文件：06-security 章節索引
版本：v0.4
狀態：審查中
負責角色：安全部
最後更新：2026-10-07
專案代號：SHORTURL
---

# 06 資安｜章節索引

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [00-asvs-level-v0.1.md](./00-asvs-level-v0.1.md) | v0.1 | 審查中 | ASVS 5.0 **Level 1** 選定與理由 |
| [01-sec-requirements-v0.2.md](./01-sec-requirements-v0.2.md) | **v0.2** | 審查中 | SEC-001～**SEC-017**（G2 回寫 TM；新增 016／017） |
| [01-sec-requirements-v0.1.md](./01-sec-requirements-v0.1.md) | v0.1 | 歷史 | G1 初版；請以 v0.2 為準 |
| [02-g1-security-pregate-v0.1.md](./02-g1-security-pregate-v0.1.md) | v0.1 | 審查中 | G1 審查部用書面前置結論 |
| [03-ci-security-gates-v0.3.md](./03-ci-security-gates-v0.3.md) | **v0.3** | 審查中 | **CI 門檻（審查中）**：Go 1.27.x／gosec v2.29.0／govulncheck v1.8.0；dependency-review 暫行；方案 A 目標態 |
| [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) | v0.2 | 待取代 | G3 定案稿；**待 v0.3 核准後取代** |
| [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md) | v0.1 | 歷史 | 已由 v0.2 取代 |
| [04a-dfd-working-v0.1.md](./04a-dfd-working-v0.1.md) | v0.1 | 審查中 | DFD 對齊說明（採設計部定稿；**EE1**） |
| [04-threat-model-v0.1.md](./04-threat-model-v0.1.md) | v0.1 | 審查中 | STRIDE 威脅模型（23 條 TM） |
| [05-residual-risk-register-v0.1.md](./05-residual-risk-register-v0.1.md) | v0.1.1 | 審查中 | 殘餘風險；RR-001 已接受；**RR-006 待使用者拍板** |
| [06-g2-security-pregate-v0.1.md](./06-g2-security-pregate-v0.1.md) | v0.1 | 審查中 | G2 審查部用書面前置結論 |
| [07-g3-security-pregate-v0.1.md](./07-g3-security-pregate-v0.1.md) | v0.1 | 審查中 | G3 資安書面前置 |
| [08-pr-002-security-review-v0.1.md](./08-pr-002-security-review-v0.1.md) | **v0.1** | 審查中 | **PR #2 安全審查**（附條件核准） |
| [09-cr-ci-deviations-ruling-v0.1.md](./09-cr-ci-deviations-ruling-v0.1.md) | **v0.1** | 審查中 | **CI 偏離＋合併路徑裁示**（維運已同意 a／b） |

## 準據 DFD

| 項目 | 路徑／內容 |
|---|---|
| 設計部 DFD 正式稿（編號凍結） | [`docs/03-data/02-dfd-v0.1.md`](../03-data/02-dfd-v0.1.md) |
| 外部實體 | **EE1**（匿名客戶端） |
| 信任邊界 | TB-01～TB-04 |
| 處理／儲存／流 | P1～P7；DS1／DS2；DF-L0-01～06、DF-01～21 |

## 本章狀態（G3＋PR #2）

| 項目 | 狀態 |
|---|---|
| ASVS 等級 | 已選定 L1（審查中） |
| SEC 清單 | **v0.2**（SEC-001～017） |
| STRIDE／TM-xxx | 已完成 v0.1（23 條） |
| 殘餘高以上 | RR-001／TM-002：**D-06 已接受**（G5 前再評） |
| 殘餘中（合併治理） | **RR-006** 待使用者拍板（方案 A／B／C） |
| CI 門檻書面 | **v0.3 審查中**（v0.2 待取代）；阻擋語意未放寬 |
| PR #2 安全審查 | **附條件核准**（補 `::warning::`） |
| CI 實作／分支保護生效 | 四 checks 已綠於 PR #2；Rulesets 方案 A **待拍板** |

## 編號清單（供品保追溯）

- SEC-001～SEC-017  
- TM-001～TM-023  
- RR-001～**RR-006**  
- SCR-001～SCR-005（見 PR #2 審查）  

## 交接單：PR #2 安全產出 → 總協調／審查／維運／研發

- 交出角色 → 接收角色：安全部 → 總協調（呈使用者拍板合併路徑；登記 CR 編號）；審查部（PR 可合併結論）；維運（Rulesets／graph）；研發（補 `::warning::`）
- 日期：2026-10-07
- 交付物：08 審查、09 裁示、03 門檻 v0.3、RR-006、本索引
- 已知限制：dependency graph 須**使用者**開啟；合併路徑須**使用者**拍板並於 Settings 操作
- 待決：使用者選 A／B／C；CR 編號登記

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 資安章節初建 |
| v0.2 | 2026-10-07 | 安全部 | G2：STRIDE／殘餘／前置；SEC v0.2 |
| v0.3 | 2026-10-07 | 安全部 | G3：CI 門檻 v0.2；G3 前置 |
| v0.4 | 2026-10-07 | 安全部 | PR #2：審查 08、裁示 09、門檻 v0.3、RR-006 |
