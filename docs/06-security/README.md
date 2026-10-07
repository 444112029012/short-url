---
文件：06-security 章節索引
版本：v0.7
狀態：審查中
負責角色：安全部
最後更新：2026-10-07
專案代號：SHORTURL
---

# 06 資安｜章節索引

> 命名規則：檔名保留主版號；小版修訂（如 v0.3.1、v0.3.2）只更新檔頭「版本」與修訂紀錄，不改檔名，以免斷鏈。

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [00-asvs-level-v0.1.md](./00-asvs-level-v0.1.md) | v0.1 | 審查中 | ASVS 5.0 **Level 1** 選定與理由 |
| [01-sec-requirements-v0.2.md](./01-sec-requirements-v0.2.md) | **v0.2** | 審查中 | SEC-001～**SEC-017**（G2 回寫 TM；新增 016／017） |
| [01-sec-requirements-v0.1.md](./01-sec-requirements-v0.1.md) | v0.1 | 歷史 | G1 初版；請以 v0.2 為準 |
| [02-g1-security-pregate-v0.1.md](./02-g1-security-pregate-v0.1.md) | v0.1 | 審查中 | G1 審查部用書面前置結論 |
| [03-ci-security-gates-v0.3.md](./03-ci-security-gates-v0.3.md) | **v0.3.3** | **已核准**（D-07／D-08） | **現行 CI 門檻**：Go 1.27.x／gosec v2.29.0／govulncheck v1.8.0（CR-001）；dependency-review 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束；方案 A 已生效（D-08） |
| [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) | v0.2.3 | 已取代 | G3 定案稿；**已由 v0.3 取代（D-07）** |
| [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md) | v0.1 | 歷史 | 已由 v0.2 取代 |
| [04a-dfd-working-v0.1.md](./04a-dfd-working-v0.1.md) | v0.1 | 審查中 | DFD 對齊說明（採設計部定稿；**EE1**） |
| [04-threat-model-v0.1.md](./04-threat-model-v0.1.md) | v0.1.1 | 審查中 | STRIDE 威脅模型（23 條 TM）；TM-002 已接受（D-06） |
| [05-residual-risk-register-v0.1.md](./05-residual-risk-register-v0.1.md) | v0.1.5 | 審查中 | 殘餘風險；RR-001 已接受（D-06）；**RR-006 已接受（D-08；D-09 維持並補強流程）**；含 PR #6 補償控制偏離紀錄 |
| [06-g2-security-pregate-v0.1.md](./06-g2-security-pregate-v0.1.md) | v0.1.1 | 審查中 | G2 審查部用書面前置結論（後註：D-06 已接受） |
| [07-g3-security-pregate-v0.1.md](./07-g3-security-pregate-v0.1.md) | v0.1 | 審查中 | G3 資安書面前置 |
| [08-pr-002-security-review-v0.1.md](./08-pr-002-security-review-v0.1.md) | **v0.1.4** | **已核准** | **PR #2 安全審查**：核准；CR-001（D-07）與合併路徑 A（D-08）已核准 |
| [09-cr-ci-deviations-ruling-v0.1.md](./09-cr-ci-deviations-ruling-v0.1.md) | **v0.1.4** | **已核准**（D-07／D-08） | **CI 偏離（CR-001）＋合併路徑（D-08）裁示**；履行複核（對象 PR #6 `37d906a`）：已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束 |

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
| 殘餘中（合併治理） | **RR-006 已接受**（D-08；G5 前再評） |
| CI 門檻書面 | **v0.3 已核准**（D-07；v0.2 已取代）；阻擋語意未放寬 |
| PR #2 安全審查 | **核准**（SCR-001 已關）；PR #2 已於 2026-10-07 13:56:49 UTC+8 合併（`a8246ce`） |
| CI 實作／分支保護生效 | 四 checks 已綠於 PR #2；Rulesets 方案 A **已生效**（D-08；取證 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)） |
| dependency-review fail-closed | Dependency graph 已開（compare 200）；已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束（原期限 min(2026-10-14, G5 送審前)，已於期限內履行） |

## 編號清單（供品保追溯）

- SEC-001～SEC-017  
- TM-001～TM-023  
- RR-001～**RR-006**  
- SCR-001～SCR-005（見 PR #2 審查）  

## 交接單：PR #2 安全產出 → 總協調／審查／維運／研發

- 交出角色 → 接收角色：安全部 → 總協調（D-07／D-08 已登記；CR-001）；審查部（PR 可合併結論）；維運／研發（fail-closed workflow）
- 日期：2026-10-07（v0.7 更新）
- 交付物：08 審查、09 裁示、03 門檻 v0.3、RR-006、本索引
- 已完成：使用者核准 CR-001（D-07）、選 A 並接受 RR-006（D-08）；Rulesets 已設定（取證 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）；Dependency graph 已開；dependency-review 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束
- 已合併：PR #2（2026-10-07 13:56:49 UTC+8，`a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`）；PR #6 fail-closed（2026-10-07 14:04:02 UTC+8，`161a4fdd0abf9b242ec22ff2c1cdba4d82ba9f86`；合併時無審查部「可合併」結論，審查部事後補審，偏離記於 05 RR-006）
- 剩餘：文件 PR B（`docs/ENG-020-g4-cr001-docs-sync`）；合併時須依 D-08 於 PR 留言附審查部「可合併」結論＋綠燈連結

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 資安章節初建 |
| v0.2 | 2026-10-07 | 安全部 | G2：STRIDE／殘餘／前置；SEC v0.2 |
| v0.3 | 2026-10-07 | 安全部 | G3：CI 門檻 v0.2；G3 前置 |
| v0.4 | 2026-10-07 | 安全部 | PR #2：審查 08、裁示 09、門檻 v0.3、RR-006 |
| v0.5 | 2026-10-07 | 安全部 | 依 D-07／D-08 同步索引：門檻 v0.3 已核准、v0.2 已取代、08／09 已核准、RR-006 已接受；記 PR #4 fail-closed 履行狀態 |
| v0.6 | 2026-10-07 | 安全部 | 維運交叉比對修正：Rulesets 證據改指 repo 檔；fail-closed 改為統一句型「已恢復 fail-closed，暫行條款結束」（含 fail-closed PR 連結佔位字串）；新增命名規則（小版修訂不改檔名（檔名保留主版號））；索引版本同步（03 v0.3.2／v0.2.3、05 v0.1.3、06 v0.1.1、08 v0.1.3、09 v0.1.3） |
| v0.7 | 2026-10-07 | 安全部 | 佔位字串替換為 PR #6；索引同步（03 v0.3.3、05 v0.1.4、08 v0.1.4、09 v0.1.4）；交接單更新：PR #2、PR #6 已合併，剩文件 PR B；註記 PR #6 事後補審 |
