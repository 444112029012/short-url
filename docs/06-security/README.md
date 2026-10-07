---
文件：06-security 章節索引
版本：v0.3
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
| [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) | **v0.2** | 審查中 | **G3 CI 門檻定案**（gosec／govulncheck／gitleaks；對齊 `docs/09-cicd-env/`） |
| [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md) | v0.1 | 歷史 | 已由 v0.2 取代 |
| [04a-dfd-working-v0.1.md](./04a-dfd-working-v0.1.md) | v0.1 | 審查中 | DFD 對齊說明（採設計部定稿；**EE1**） |
| [04-threat-model-v0.1.md](./04-threat-model-v0.1.md) | v0.1 | 審查中 | STRIDE 威脅模型（23 條 TM） |
| [05-residual-risk-register-v0.1.md](./05-residual-risk-register-v0.1.md) | v0.1 | 審查中 | 殘餘風險；RR-001 **已接受**（D-06；G5 前再評） |
| [06-g2-security-pregate-v0.1.md](./06-g2-security-pregate-v0.1.md) | v0.1 | 審查中 | G2 審查部用書面前置結論 |
| [07-g3-security-pregate-v0.1.md](./07-g3-security-pregate-v0.1.md) | **v0.1** | 審查中 | **G3 資安書面前置**：門檻書面已齊；實作待維運 |

## 準據 DFD

| 項目 | 路徑／內容 |
|---|---|
| 設計部 DFD 正式稿（編號凍結） | [`docs/03-data/02-dfd-v0.1.md`](../03-data/02-dfd-v0.1.md) |
| 外部實體 | **EE1**（匿名客戶端） |
| 信任邊界 | TB-01～TB-04 |
| 處理／儲存／流 | P1～P7；DS1／DS2；DF-L0-01～06、DF-01～21 |

## 本章狀態（G3）

| 項目 | 狀態 |
|---|---|
| ASVS 等級 | 已選定 L1（審查中） |
| SEC 清單 | **v0.2**（SEC-001～017） |
| STRIDE／TM-xxx | 已完成 v0.1（23 條） |
| 殘餘高以上 | RR-001／TM-002：**D-06 已接受**（G5 前再評） |
| CI 門檻書面 | **v0.2 已定案**（對齊維運 `docs/09-cicd-env/` 路徑與 job 名） |
| 維運 CI 設計稿 | 已對齊 `docs/09-cicd-env/README.md`＋`01-ci-pipeline-v0.1.md`（job 名一致） |
| CI 實作／分支保護生效 | **待維運關閉**（非本前置範圍；維運標「實作待 repo」） |
| G3 安全前置結論 | **附條件通過**（書面已齊；生效待驗證） |

## 編號清單（供品保追溯）

- SEC-001～SEC-017  
- TM-001～TM-023（見威脅模型；編號連續使用至 TM-023）  
- RR-001～RR-005  

## 交接單：G3 安全前置 → 總協調／審查／維運／品保

- 交出角色 → 接收角色：安全部 → 總協調（審查部作 G3 前置勾選；維運實作流水線與分支保護；品保知悉 required checks）
- 日期：2026-10-07
- 交付物清單：CI 門檻 v0.2、G3 前置結論、本索引
- 版本：上表（狀態：審查中）
- 對應需求編號：SEC-009～011；SI-05；NFR-004
- 已知問題／限制：門檻書面已齊；維運 README／01 已對齊；**GitHub workflow／分支保護尚未視為生效**；請維運改引 v0.2
- 需要下游注意：維運 job 名須為 `secrets-gitleaks`｜`sast-gosec`｜`sca-govulncheck`｜`unit-test`；**不得放寬**門檻語意
- 待決問題：無（D-06 已關）；待維運關閉實作項

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 資安章節初建 |
| v0.2 | 2026-10-07 | 安全部 | G2：STRIDE／殘餘／前置；SEC v0.2；對齊設計部 DFD EE1 |
| v0.3 | 2026-10-07 | 安全部 | G3：CI 門檻 v0.2；G3 前置；RR-001 已接受；索引對齊 09-cicd-env |
