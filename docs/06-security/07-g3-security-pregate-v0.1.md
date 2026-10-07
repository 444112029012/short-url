---
文件：G3 資安書面前置結論
版本：v0.1
狀態：審查中
負責角色：安全部
最後更新：2026-10-07
對應需求：NFR-004、SI-05；SEC-009～SEC-011；關卡 G3
專案代號：SHORTURL
---

# G3 資安書面前置結論（供審查部）

> 依階段關卡檢查表 G3：**前置＝安全確認 CI 掃描門檻**。  
> 本前置僅證明**門檻書面已定案**；**不**等同流水線已在 GitHub 上跑過或分支保護已生效。

---

## 1. 結論（摘要）

**建議審查勾選：安全前置已附（CI 掃描門檻書面已齊）；可進入 G3 細審（附條件）。**

| 檢核項 | 結論 |
|---|---|
| CI 掃描門檻文件（SAST／SCA／機密阻擋條件） | **通過** — [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) |
| 零預算工具已選定（Go） | **通過** — gitleaks；gosec；govulncheck（＋約定備援） |
| Required job 名已定案 | **通過** — `secrets-gitleaks`｜`sast-gosec`｜`sca-govulncheck`｜`unit-test` |
| 與 SEC-009～011／SI-05／NFR-004 對照 | **通過** — 見門檻文件 §9 |
| 分支保護／workflow **已在 repo 生效** | **未完成本前置範圍** — 屬維運交付；**G3 過關條件仍須驗證** |
| D-06／RR-001 | **已接受**（G5 前再評）— 不阻擋本前置 |

**一句結論：**CI 資安掃描門檻**書面已齊且工具／job／阻擋語意已定案**，建議審查勾選「安全前置已附」；惟 **GitHub Actions 實作與分支保護生效尚未視為完成**，G3 正式過關之「CI 已啟用 SAST、SCA、機密掃描且分支保護生效」須待維運關閉並由審查核對證據。

### 明確區分

| 層次 | 狀態 | 負責 |
|---|---|---|
| **門檻文件定案（本前置）** | **已完成** | 安全部 |
| **分支保護／job 實作生效（G3 過關條件）** | **待關閉** | 維運部交付；審查驗證 |

**請勿將本前置解讀為「CI 已在 GitHub 上跑過」。**

---

## 2. 交付物連結

| 交付物 | 路徑 |
|---|---|
| CI 門檻 v0.2（準據） | [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) |
| CI 門檻 v0.1（歷史） | [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md) |
| SEC 清單 | [01-sec-requirements-v0.2.md](./01-sec-requirements-v0.2.md) |
| 殘餘風險（D-06 已接受） | [05-residual-risk-register-v0.1.md](./05-residual-risk-register-v0.1.md) |
| 維運 CI／環境（對齊路徑） | [../09-cicd-env/README.md](../09-cicd-env/README.md)、[01-ci-pipeline-v0.1.md](../09-cicd-env/01-ci-pipeline-v0.1.md)、[02-environments-v0.1.md](../09-cicd-env/02-environments-v0.1.md)、[03-secrets-management-v0.1.md](../09-cicd-env/03-secrets-management-v0.1.md) |
| 範例 workflow（選） | [../09-cicd-env/examples/ci-devsecops.yml.example](../09-cicd-env/examples/ci-devsecops.yml.example) |

---

## 3. 選定工具與阻擋規則（摘要）

| 掃描 | 選定工具 | 阻擋 |
|---|---|---|
| 機密 | gitleaks（輔 GitHub secret scanning＋push protection） | **一律擋**真實機密 |
| SAST | gosec（可選 CodeQL go） | **High／Critical** 擋 |
| SCA | govulncheck（PR 可併用 dependency-review） | **Critical** 擋；**High** 有修補逾 **14 日**未升擋 |
| 測試 | `unit-test`（go test） | 失敗擋 |

---

## 4. 待維運關閉項（G3 過關前）

- [x] 設計稿：`docs/09-cicd-env/README.md`、`01-ci-pipeline-v0.1.md`、`02-environments-v0.1.md`、`03-secrets-management-v0.1.md` 已存在；job／工具與 01 一致（安全部已 Read README／01）  
- [x] 維運文件交叉引用改為門檻 **v0.2**（一致性簽核通過，2026-10-07）  
- [ ] 可選 `examples/ci-devsecops.yml.example`  
- [ ] 於 repo 落地 `.github/workflows/`，job 名：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`  
- [ ] `main` 分支保護：禁直推、僅 PR、禁 force push、≥1 核准、上列 required checks；管理員繞過須留紀錄  
- [ ] 啟用 GitHub secret scanning＋push protection（及 Dependabot alerts 建議）  
- [ ] 以測試分支證實門檻會阻擋，之後刪除該分支  
- [ ] 門檻語意與安全部 v0.2 **一致且未放寬**；例外／allowlist 存放位置與安全部對齊  

---

## 5. 建議審查部勾選用語

- 前置結論：安全 □ **已附**（本文件＋門檻 v0.2）  
- CI 掃描門檻書面：□ **已齊**（工具／job／阻擋／例外／分支保護建議）  
- CI 流水線與分支保護**生效**：□ **待維運關閉／審查核對證據** — 不阻擋「前置已附進細審」，但 G3 過關條件勾選前須關閉  
- 可進入 G3 細審（安全前置已附）：□ **是（附條件：知悉實作未生效）**

---

## 6. 已知限制

| 限制 | 處理 |
|---|---|
| 維運 README／01／02／03 設計稿已齊；example 可選；job 名與 01 一致 | **一致性簽核通過（2026-10-07）**；`09-cicd-env` 已改指門檻 v0.2 |
| 維運稿標「實作待 repo」 | 與本前置區分一致：書面≠生效 |
| 無 Docker／Terraform | IaC／容器標「G3 不適用／出現即納入」（見門檻 §8） |
| D-06／RR-001 | 已接受；G5 前再評（與本前置無衝突） |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G3 安全前置：門檻書面已齊；實作生效待維運；建議附條件進細審 |
| v0.1.1 | 2026-10-07 | 安全部 | 維運一致性簽核通過；待關閉項僅剩 workflow／分支保護／擋測實作 |
