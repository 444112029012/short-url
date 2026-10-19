---
文件：CI 資安掃描阻擋門檻
版本：v0.2.3
狀態：已取代（由 v0.3 取代；D-07）
負責角色：安全部
最後更新：2026-10-07
對應需求：NFR-004、SI-05；SEC-009～SEC-011
專案代號：SHORTURL
與維運對齊狀態：**一致性簽核通過（2026-10-07，維運部）**；`docs/09-cicd-env/` 已改指本 v0.2；job／阻擋語意未放寬；**GitHub 實作生效仍待維運交付（repo 就緒後）**
---
> **已由 v0.3 取代（D-07，2026-10-07）**：請改以 [03-ci-security-gates-v0.3.md](./03-ci-security-gates-v0.3.md) 為準（CR-001；合併路徑 D-08）。本版僅供歷史追溯。


# CI 資安掃描阻擋門檻（v0.2）

> **本版取代** [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md)。  
> **零現金預算**：僅採免費／開源／GitHub 內建工具；不得以付費授權為必選。  
> **技術棧（D-04）**：Go＋chi／SQLite。  
> **門檻語意**：以本文件為準；維運實作**不得放寬**阻擋級別。  
> **與維運對齊**：**一致性簽核通過（2026-10-07）**；維運已改指本 v0.2；job 名與阻擋語意一致且未放寬。其餘路徑：
> - [`docs/09-cicd-env/README.md`](../09-cicd-env/README.md)
> - [`docs/09-cicd-env/01-ci-pipeline-v0.1.md`](../09-cicd-env/01-ci-pipeline-v0.1.md)
> - [`docs/09-cicd-env/02-environments-v0.1.md`](../09-cicd-env/02-environments-v0.1.md)
> - [`docs/09-cicd-env/03-secrets-management-v0.1.md`](../09-cicd-env/03-secrets-management-v0.1.md)
> - （選）[`docs/09-cicd-env/examples/ci-devsecops.yml.example`](../09-cicd-env/examples/ci-devsecops.yml.example)

---

## 1. 適用範圍

| 項目 | 說明 |
|---|---|
| 適用分支 | 保護 `main`；所有進入合併之 PR |
| 適用環境 | GitHub Actions CI（G3 過關前須**實作生效**） |
| 不適用 | 本機未推送實驗分支（仍建議 pre-commit／本機自跑） |
| 預算約束 | 零現金；工具須 OSS 或 GitHub 內建 |

---

## 2. 本專案選定工具（寫死）

| 掃描 | 對應 SEC | **本專案選定（主選）** | 備援／併用（可選） | 授權／成本 |
|---|---|---|---|---|
| **機密掃描** | SEC-011 | **gitleaks**（CI job） | GitHub **secret scanning**＋**push protection**（倉庫設定啟用） | 開源／GitHub 內建 |
| **SAST** | SEC-009 | **gosec**（Go） | 可選 **CodeQL**（language: `go`，security-extended）作補強 | 開源／GitHub 內建 |
| **SCA** | SEC-010 | **govulncheck**（Go 官方漏洞資料庫） | PR 可併用 **dependency-review-action**（對齊維運建議 `fail-on-severity: critical`；主擋仍以 govulncheck＋§4）；備援 Trivy fs／OSV-Scanner | 開源／GitHub 內建 |
| **單元測試** | （品質門檻；G3／DoD） | `go test` | — | 標準庫 |

> 選定理由：對齊 Go 技術棧、零預算、可於 GitHub Actions 直接執行；與維運定案一致。

---

## 3. 建議／定案 workflow job 名稱（required checks）

供分支保護「必要狀態檢查」使用；名稱須與維運流水線一致：

| Required check（定案） | 職責 | 失敗即擋合併 |
|---|---|---|
| `secrets-gitleaks` | gitleaks 機密掃描 | 是（任何真實機密發現） |
| `sast-gosec` | gosec SAST | 是（High／Critical） |
| `sca-govulncheck` | govulncheck SCA | 是（依 §4 阻擋規則） |
| `unit-test` | `go test ./...`（或等價） | 是（測試失敗） |

> 可選非 required（建議排程或另 job）：`sast-codeql-go`、`sca-dependency-review`（僅 PR）。  
> 範例 workflow 對齊：[`docs/09-cicd-env/examples/ci-devsecops.yml.example`](../09-cicd-env/examples/ci-devsecops.yml.example)。

---

## 4. 掃描類型與阻擋規則（沿用 v0.1，不得無故放寬）

| 掃描 | 對應 SEC | **阻擋規則（定案）** | 追蹤／不擋 | 例外核准 |
|---|---|---|---|---|
| **機密掃描** | SEC-011 | 偵測疑似**真實機密**（雲端金鑰、私鑰、密碼、權杖、連線字串）→ **一律擋** | 明顯測試用假值且經審核之 allowlist 路徑／規則 | **原則不准**關閉整條規則；單一路徑例外須安全部書面＋紀錄；發現真實外洩須立即輪替 |
| **SAST** | SEC-009 | **High／Critical**（gosec 對應高／嚴重級，或等價 severity）→ **失敗並阻擋合併** | Medium／Low：開缺陷追蹤；G5 前應清或書面接受 | 誤判：書面理由、替代控制、安全部核准；高風險另呈使用者（D-06 級） |
| **SCA** | SEC-010 | 已知漏洞 **Critical** → **擋**；**High** 且有可用修補版本超過 **14 日**未升級 → **擋** | Medium／Low：追蹤；無修補版本：記錄並評估緩解；`go.sum`／鎖定檔須存在 | 同左；EOL 套件另案；高風險接受仍走 D-06 |

**一句阻擋規則：**機密一律擋；SAST High／Critical 擋；SCA Critical 擋、High 有修補逾 14 日未升擋。

---

## 5. 修補時限政策（對齊 ASVS 5.0 15.1.1；沿用 v0.1）

| 嚴重度 | 目標修復時限（發現起算） |
|---|---|
| Critical | **7 日**內修復，或下線受影響功能／使用者接受殘餘 |
| High | **14 日** |
| Medium | **30 日**（追蹤） |
| Low | 下一個里程碑前評估 |

---

## 6. 例外核准流程（誤判／暫緩）

1. **開缺陷單**：標嚴重度、工具／規則 ID、影響路徑、發現日期。  
2. **誤判**：附重現證據與為何為誤判；安全部核准後可 allowlist（單一路徑／規則，禁止整條關閉）；設覆核日。  
3. **暫緩修補**：寫替代控制、到期日、負責人；安全部核准。  
4. **高風險以上接受殘餘**：須使用者書面（對齊章程 **D-06**）；不得僅由部門內部默許。  
5. **禁止**：無紀錄關閉掃描規則；以編碼／更名繞過機密掃描；「先合併再補」而無核准單。

例外紀錄建議存放：缺陷追蹤系統＋本專案變更／決策日誌交叉引用（維運／安全共用欄位：核准人、到期日、替代控制）。

---

## 7. 分支保護建議（對齊維運；G3 過關須生效）

| 設定 | 定案 |
|---|---|
| `main` 禁直接 push | **是**；所有變更僅能經 PR |
| 核准數 | 至少 **1** 位核准（練手專案可 1；新 commit 推送後撤銷舊核准） |
| Required checks | `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` |
| Force push／刪除受保護分支 | **禁止** |
| 管理員 | 規則亦適用；若允許繞過，**每次須留書面紀錄**（誰、為何、何時） |
| 其他建議 | 要求分支為最新；必須解決審查對話；啟用 Dependabot alerts（`gomod`＋`github-actions`）、secret scanning＋push protection；可選 CODEOWNERS（`.github/workflows/`、資安相關路徑） |

細節實作見維運 [`02-environments-v0.1.md`](../09-cicd-env/02-environments-v0.1.md)／流水線文件；**門檻語意以本文件為準**。

---

## 8. IaC／容器掃描

| 類型 | 本專案現況 | G3 門檻 |
|---|---|---|
| **Dockerfile／容器映像** | MVP **暫無** Docker 交付假設 | **不適用（G3 可選）**；一旦出現 Dockerfile／映像建置，即納入：Trivy image（或等價）對 Critical／High 且有修補者**擋部署／擋合併建置 job** |
| **Terraform／K8s 等 IaC** | **暫無** | **不適用（G3 可選）**；一旦出現 IaC，即納入：Checkov 或 Trivy misconfig；公開存取、過寬權限、未加密、以 root 執行等高風險規則**擋** |

出現即納入時：新增 required check 名稱由維運提案、安全部確認門檻語意後寫入本文件升版。

---

## 9. 與 SEC-009～011、SI-05、NFR-004 對照

| 需求／SEC | 本文件對應 |
|---|---|
| **SEC-009**（SAST 擋高／嚴重） | §2 gosec；§4 SAST 列；job `sast-gosec` |
| **SEC-010**（SCA＋修補時限） | §2 govulncheck；§4 SCA 列；§5 時限；job `sca-govulncheck` |
| **SEC-011**（機密掃描一律擋） | §2 gitleaks＋GitHub push protection；§4 機密列；job `secrets-gitleaks` |
| **SI-05** | 章程成功指標：CI 資安掃描門檻存在並可驗證阻擋 |
| **NFR-004** | 須有 ASVS／SEC／CI 門檻且流水線依門檻阻擋（G3 驗證實作；G5 關閉高風險或使用者接受） |

驗證方式（供品保／審查）：（1）本文件已定案；（2）維運 workflow 存在且 job 名一致；（3）分支保護勾選上列 required checks；（4）以測試用假密鑰樣式／已知高危 gosec 規則／govulncheck 可偵測案例證實會失敗擋 PR（驗證後刪除測試分支）。

---

## 10. Workflow 安全守則（摘要，供維運對齊）

- 頂層 `permissions: contents: read`；上傳 SARIF 等才加 `security-events: write`。  
- 第三方 action **釘選完整 commit SHA**；Dependabot 更新 `github-actions`。  
- 不在 `pull_request_target` 執行不可信 fork 程式碼。  
- 機密僅放 GitHub Secrets／Environments（見 [`03-secrets-management-v0.1.md`](../09-cicd-env/03-secrets-management-v0.1.md)）。  
- 不把 PR 標題等不可信輸入直接內插進 `run:`。

---

## 11. 與維運對齊檢查清單

| 項 | 狀態 |
|---|---|
| 工具選定（gitleaks／gosec／govulncheck） | **已定案且與維運 01 一致** |
| Required job 名 | **已定案且與維運 01／README 一致**（§3） |
| 門檻語意不得放寬 | **已定案**（§4；維運聲明不放寬） |
| 維運 `docs/09-cicd-env/` | **設計稿齊**：README＋01＋02＋03 已存在；安全已 Read README／01（job／工具一致）；example 可選待補；**實作待 repo** |
| 維運引用門檻版本 | **已簽核** — `09-cicd-env` 已改指 v0.2（2026-10-07） |
| GitHub workflow 實作＋分支保護生效 | **待維運關閉**（G3 過關條件；文件標「實作待 repo」） |
| 測試分支證實會擋 | **待維運關閉** |
| 例外／allowlist 存放位置 | 見維運 03；**啟用時與安全共同關閉／簽核** |

---

## 12. 禁止事項

- 禁止為「先合併」而無紀錄關閉掃描規則。  
- 禁止將真實機密寫入版控後以「之後再刪」通過。  
- 禁止以編碼／更名／分割檔案繞過機密掃描當作例外。  
- 禁止將付費授權工具列為必選（零預算）。  
- 禁止維運實作放寬本文件阻擋級別而未走本文件升版＋安全部核准。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 初版門檻；標 G3 與維運對齊 |
| v0.2 | 2026-10-07 | 安全部 | G3：寫死 Go 零預算工具；定案 job 名與分支保護；對齊 `docs/09-cicd-env/`；IaC／容器標不適用／出現即納入；不放寬阻擋語意 |
| v0.2.1 | 2026-10-07 | 安全部 | 記錄維運一致性簽核通過；引用已改指 v0.2 |
| v0.2.2 | 2026-10-07 | 安全部 | 依 D-07／D-08：頂部註記改「已由 v0.3 取代（D-07）」；狀態改已取代；內文未改 |
| v0.2.3 | 2026-10-07 | 安全部 | 版本與檔名說明：小版修訂不改檔名（檔名保留主版號）；內文未改 |
