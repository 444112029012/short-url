---
文件：CI 資安掃描阻擋門檻
版本：v0.3
狀態：審查中
負責角色：安全部
最後更新：2026-10-07
對應需求：NFR-004、SI-05；SEC-009～SEC-011
專案代號：SHORTURL
與維運對齊狀態：對齊維運 `01-ci-pipeline-v0.2`（工具鏈／dependency-review 暫行／§5.2）；一致性以本文件阻擋語意為準；**合併路徑 Rulesets 目標態待使用者拍板後生效**
對應 CR：CR-PENDING-01～04（見 [09-cr-ci-deviations-ruling-v0.1.md](./09-cr-ci-deviations-ruling-v0.1.md)；編號待總協調登記）
---

# CI 資安掃描阻擋門檻（v0.3）

> **本版預定取代** [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md)（待本版核准後生效）。  
> **零現金預算**：僅採免費／開源／GitHub 內建工具；不得以付費授權為必選。  
> **技術棧（D-04）**：Go＋chi／SQLite。  
> **門檻語意**：以本文件為準；維運實作**不得放寬**阻擋級別（High／Critical／機密語意與 v0.2 相同）。  
> **v0.3 變更摘要**：工具鏈 Go 1.27.x；gosec v2.29.0；govulncheck v1.8.0；dependency-review 附條件暫行；分支保護改寫為方案 A 目標態（待使用者拍板）。

與維運對齊路徑：
- [`docs/09-cicd-env/README.md`](../09-cicd-env/README.md)
- [`docs/09-cicd-env/01-ci-pipeline-v0.2.md`](../09-cicd-env/01-ci-pipeline-v0.2.md)
- [`docs/09-cicd-env/05-pr2-ci-change-review-v0.1.md`](../09-cicd-env/05-pr2-ci-change-review-v0.1.md)
- [`docs/09-cicd-env/02-environments-v0.1.md`](../09-cicd-env/02-environments-v0.1.md)
- [`docs/09-cicd-env/03-secrets-management-v0.1.md`](../09-cicd-env/03-secrets-management-v0.1.md)

---

## 1. 適用範圍

| 項目 | 說明 |
|---|---|
| 適用分支 | 保護 `main`；所有進入合併之 PR |
| 適用環境 | GitHub Actions CI |
| 不適用 | 本機未推送實驗分支（仍建議 pre-commit／本機自跑） |
| 預算約束 | 零現金；工具須 OSS 或 GitHub 內建 |

---

## 2. 本專案選定工具（寫死）

| 掃描 | 對應 SEC | **本專案選定（主選）** | 版本鎖定（v0.3） | 備援／併用（可選） | 授權／成本 |
|---|---|---|---|---|---|
| **機密掃描** | SEC-011 | **gitleaks**（CI job） | `8.21.2`（與維運一致） | GitHub **secret scanning**＋**push protection** | 開源／GitHub 內建 |
| **SAST** | SEC-009 | **gosec**（Go） | **`v2.29.0`**（精確釘版；禁止 `@latest`） | 可選 **CodeQL**（language: `go`） | 開源／GitHub 內建 |
| **SCA** | SEC-010 | **govulncheck** | **`v1.8.0`**（精確釘版；禁止 `@latest`） | PR 併用 **dependency-review-action**（見 §4.1 暫行）；備援 Trivy fs／OSV-Scanner | 開源／GitHub 內建 |
| **單元測試** | （品質門檻） | `go test` | — | — | 標準庫 |
| **Go toolchain** | （供應鏈） | 官方 Go | **`go.mod` `go 1.27.1`**；CI `setup-go` **`1.27.x`**（須一致） | — | 官方 |

> 選定理由：對齊 Go 技術棧、零預算、可於 GitHub Actions 直接執行；離開已停修之 1.24 線屬安全強化（CR-PENDING-01／02）。

---

## 3. 建議／定案 workflow job 名稱（required checks）

| Required check（定案） | 職責 | 失敗即擋合併 |
|---|---|---|
| `secrets-gitleaks` | gitleaks 機密掃描 | 是（任何真實機密發現） |
| `sast-gosec` | gosec SAST | 是（High／Critical） |
| `sca-govulncheck` | govulncheck SCA | 是（依 §4 阻擋規則） |
| `unit-test` | `go test ./...`（或等價） | 是（測試失敗） |

> 可選非 required：`sast-codeql-go`、`sca-dependency-review`（僅 PR；見 §4.1）。  
> **job 名稱與 High／Critical 阻擋語意較 v0.2 未放寬。**

---

## 4. 掃描類型與阻擋規則（沿用 v0.2，不得無故放寬）

| 掃描 | 對應 SEC | **阻擋規則（定案）** | 追蹤／不擋 | 例外核准 |
|---|---|---|---|---|
| **機密掃描** | SEC-011 | 偵測疑似**真實機密** → **一律擋** | 明顯測試用假值且經審核之 allowlist | **原則不准**關閉整條規則 |
| **SAST** | SEC-009 | **High／Critical** → **失敗並阻擋合併** | Medium／Low：開缺陷追蹤 | 誤判須書面＋安全部核准 |
| **SCA** | SEC-010 | 已知漏洞 **Critical** → **擋**；**High** 且有可用修補版本超過 **14 日**未升級 → **擋** | Medium／Low：追蹤；`go.sum`／鎖定檔須存在 | 同左；高風險接受仍走 D-06 |

**一句阻擋規則：**機密一律擋；SAST High／Critical 擋；SCA Critical 擋、High 有修補逾 14 日未升擋。

### 4.1 dependency-review 暫行條款（CR-PENDING-03）

| 項 | 定案 |
|---|---|
| 定位 | **非 required** 輔助；主擋仍為 `govulncheck`（required） |
| 常態 | **不接受**靜默略過 |
| **暫行** | Dependency graph 未啟用時，可探測後略過，但**必須**輸出 Actions 可見 `::warning::`；**不得**使 `sca-govulncheck` 因略過而失敗 |
| 恢復 | 使用者開啟 Dependency graph（**使用者 repo 設定動作**）後，改 **fail-closed**：compare API 非 200 即失敗，並驗證 review 實跑 |
| 期限 | 暫行最遲至 **min(2026-10-14, G5 送審前)** 失效（兩者取較早） |
| 合併前 | PR #2 須先補 `::warning::`（見 [08-pr-002-security-review-v0.1.md](./08-pr-002-security-review-v0.1.md) SCR-001） |

---

## 5. 修補時限政策（對齊 ASVS 5.0 15.1.1；沿用 v0.2）

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
4. **高風險以上接受殘餘**：須使用者書面（對齊章程 **D-06**）。  
5. **禁止**：無紀錄關閉掃描規則；以編碼／更名繞過機密掃描；「先合併再補」而無核准單。

---

## 7. 分支保護（目標態＝方案 A；待使用者拍板後生效）

> 對齊維運 `01-ci-pipeline-v0.2` §5.2 與安全部 CR-PENDING-04。拍板前仍以現況 Settings 為準；**本節為核准後目標態**。

### 7.1 Rulesets 目標態（方案 A）

| 規則 | 內容 | Bypass |
|---|---|---|
| **規則一「main-required-checks」** | 四項 required checks：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`；**禁止 force push**；**禁止刪除**受保護分支 | **無任何 bypass**（資安控制核心，必須保住） |
| **規則二「main-pr-review」** | 需要 Pull Request；至少 **1** 核准 | Repository admin，模式**僅限 PR**（不得允許直接 push） |

操作注意：刪除舊 Branch protection rule，避免與 Rulesets 雙軌卡住；設定完成後維運以 API 取證。

### 7.2 補償控制（A／B 必做）

- 每個 PR 合併前：於 PR 留言貼上審查部（必要時安全部）書面「可合併」結論與四項 check 綠燈連結。  
- 合併人＝使用者本人。  
- 每筆 admin bypass 合併須可稽核（留言含審查結論連結）。  
- **G5 前**重新評估是否仍需 admin bypass（見殘餘風險 **RR-006**）。

### 7.3 方案比較（安全立場）

| 方案 | 安全立場 |
|---|---|
| A（建議） | 四 checks 無 bypass；僅「1 核准」對 admin 可繞過＋書面補償 |
| B（不建議） | 核准數 0；平台無核准紀錄，補償偏弱 |
| C（理想） | 第二位真人 collaborator；若可行則優先於 A |

### 7.4 其他建議（沿用）

| 設定 | 定案 |
|---|---|
| 線性歷史 | 建議 |
| Dependabot alerts | 建議（`gomod`＋`github-actions`） |
| secret scanning＋push protection | 建議啟用 |
| 第三方 action | **釘選完整 commit SHA**（穩定後）；Dependabot 更新 `github-actions` |

---

## 8. IaC／容器掃描

| 類型 | 本專案現況 | G3 門檻 |
|---|---|---|
| **Dockerfile／容器映像** | MVP **暫無** | **不適用（可選）**；出現即納入 Trivy image 等 |
| **Terraform／K8s 等 IaC** | **暫無** | **不適用（可選）**；出現即納入 |

---

## 9. 與 SEC-009～011、SI-05、NFR-004 對照

| 需求／SEC | 本文件對應 |
|---|---|
| **SEC-009** | §2 gosec v2.29.0；§4 SAST；job `sast-gosec` |
| **SEC-010** | §2 govulncheck v1.8.0；§4 SCA；§4.1 暫行；job `sca-govulncheck` |
| **SEC-011** | §2 gitleaks；§4 機密；job `secrets-gitleaks` |
| **SI-05** | CI 資安掃描門檻存在並可驗證阻擋 |
| **NFR-004** | 門檻＋流水線依門檻阻擋；分支保護目標態 §7 |

---

## 10. Workflow 安全守則（摘要）

- 頂層 `permissions` 最小化（至少 `contents: read`；PR 輔助可加 `pull-requests: read`）；上傳 SARIF 等才加 `security-events: write`。  
- 第三方 action **釘選完整 commit SHA**；Dependabot 更新 `github-actions`。  
- 不在 `pull_request_target` 執行不可信 fork 程式碼。  
- 機密僅放 GitHub Secrets／Environments。  
- 不把 PR 標題等不可信輸入直接內插進 `run:`。

---

## 11. 與維運對齊檢查清單

| 項 | 狀態 |
|---|---|
| 工具選定（gitleaks／gosec／govulncheck） | **v0.3 版本已定**（與維運 01 v0.2 一致） |
| Required job 名 | **未變** |
| 門檻語意不得放寬 | **已定案**（§4） |
| dependency-review 暫行 | **已寫入 §4.1**；維運同意 |
| 合併路徑方案 A | **目標態已寫 §7**；**待使用者拍板** |
| GitHub 設定生效 | **待使用者／維運關閉** |

---

## 12. 禁止事項

- 禁止為「先合併」而無紀錄關閉掃描規則。  
- 禁止將真實機密寫入版控後以「之後再刪」通過。  
- 禁止以編碼／更名／分割檔案繞過機密掃描當作例外。  
- 禁止將付費授權工具列為必選（零預算）。  
- 禁止維運實作放寬本文件阻擋級別而未走本文件升版＋安全部核准。  
- 禁止對規則一（required checks／禁 force push）設定 bypass。  
- 禁止 dependency-review 常態靜默略過（暫行須 `::warning::` 且限期）。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 初版門檻 |
| v0.2 | 2026-10-07 | 安全部 | G3：寫死 Go 零預算工具；定案 job 名與分支保護 |
| v0.2.1 | 2026-10-07 | 安全部 | 記錄維運一致性簽核通過 |
| v0.3 | 2026-10-07 | 安全部 | CR-PENDING：Go 1.27.1／gosec v2.29.0／govulncheck v1.8.0；dependency-review 暫行至 min(2026-10-14, G5)；分支保護方案 A 目標態；阻擋語意未放寬 |
