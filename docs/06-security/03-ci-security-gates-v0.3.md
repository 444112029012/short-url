---
文件：CI 資安掃描阻擋門檻
版本：v0.3.8
狀態：已核准（D-07：CR-001；D-08：合併路徑 A；RR-006 已接受）；**v0.3.7 §2 工具鏈／SAST 暫行例外：依 D-13（CR-002）核准、安全部同意（附條件）；v0.3.8：已生效（#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）**
負責角色：安全部
最後更新：2026-10-19
對應需求：NFR-004、SI-05；SEC-009～SEC-011
專案代號：SHORTURL
與維運對齊狀態：對齊維運 `01-ci-pipeline-v0.2`（工具鏈／dependency-review／§5.2）；一致性以本文件阻擋語意為準；**合併路徑 Rulesets 方案 A 已生效（D-08；取證見 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）**
對應 CR／決策：**CR-001**（第 1～3 項：工具鏈／掃描器／dependency-review；D-07 核准）；合併路徑依 **D-08**（decision-log 未另立 CR 編號）；**CR-002**（Go 1.27.2＋gosec 暫行例外；D-13；安全部意見見 [13](./13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md)）。裁示見 [09-cr-ci-deviations-ruling-v0.1.md](./09-cr-ci-deviations-ruling-v0.1.md)
---

# CI 資安掃描阻擋門檻（v0.3）

> **本版取代** [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md)（D-07 核准後生效，2026-10-07）。  
> **零現金預算**：僅採免費／開源／GitHub 內建工具；不得以付費授權為必選。  
> **技術棧（D-04）**：Go＋chi／SQLite。  
> **門檻語意**：以本文件為準；維運實作**不得放寬**阻擋級別（High／Critical／機密語意與 v0.2 相同）。  
> **v0.3.8 變更摘要**：CR-002 改為**已生效**（#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`；合併後 main push run [39189683566](https://github.com/444112029012/short-url/actions/runs/39189683566) 四項 success，`sast-gosec` 由合併前 main `2e54f7e` 的 failure（run 39183720869）轉綠，log 有 `gosec h1 OK`、Issues 0、Go `go1.27.2`）；移除「CI 實況仍為 1.27.1／v2.29.0」註記；C6／C7 期限寫入具體日期（§2）。  
> **v0.3.7 變更摘要（CR-002／D-13；v0.3.8 起已生效）**：Go **精確釘 1.27.2**（go.mod 與 CI 一致，禁止 `1.27.x`）；gosec 暫行例外釘上游 commit `707f227`（pseudo-version＋h1 比對，期限與退場見 §2）；govulncheck v1.8.0 不變；阻擋語意不變。  
> **v0.3 變更摘要**（歷史）：工具鏈 Go 1.27.x；gosec v2.29.0；govulncheck v1.8.0；dependency-review（CR-001 第 3 項）已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束；分支保護方案 A 已生效（D-08）。

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
| **SAST** | SEC-009 | **gosec**（Go） | **暫行例外（CR-002／D-13）**：`github.com/securego/gosec/v2@v2.29.1-0.20261015133444-707f2271f665`（上游 `securego/gosec` 預設分支 `master` commit `707f2271f66579a18e1da25794b1e8012247274b`；sum.golang.org `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`）。CI **必須**於安裝後以 `go version -m` 讀出 gosec 模組 h1 並比對上列值，不符即失敗；禁止 `@latest`／`@master`／短 SHA；不得關閉或繞過 GOSUMDB／GOPROXY 驗證。**例外期限**：上游發行支援 Go 1.27.2 之 gosec 版本後 **14 日內**以 CR 換回；最遲 **2026-12-31** 由安全部重評（延長須安全部再同意）。**退場條件**：CR 換回精確釘版之發行版（`vX.Y.Z`）並合併後，本例外與 RR-010 結案。（**已生效**：#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`；main `c341be1` 與 `24b8da0` 的 `sast-gosec` log 皆有 h1 比對輸出與 `gosec h1 OK`、Issues 0） | 可選 **CodeQL**（language: `go`） | 開源／GitHub 內建 |
| **SCA** | SEC-010 | **govulncheck** | **`v1.8.0`**（精確釘版；禁止 `@latest`） | PR 併用 **dependency-review-action**（見 §4.1；fail-closed）；備援 Trivy fs／OSV-Scanner | 開源／GitHub 內建 |
| **單元測試** | （品質門檻） | `go test` | — | — | 標準庫 |
| **Go toolchain** | （供應鏈） | 官方 Go | **`go.mod` `go 1.27.2`；CI `setup-go` `"1.27.2"`**（精確釘版，**禁止 `1.27.x`／`stable`**；兩者須一致）。後續修補版以 CR（標準變更）更新，且須先驗證 `sast-gosec`／`sca-govulncheck` 皆綠。（**已生效**：main `c341be1` 起 `go.mod` 為 `go 1.27.2`、workflow 三處 `setup-go` 皆 `"1.27.2"`；main CI log 為 `go1.27.2 linux/amd64`；main `24b8da0` 仍相同） | — | 官方 |

> 選定理由：對齊 Go 技術棧、零預算、可於 GitHub Actions 直接執行；離開已停修之 1.24 線屬安全強化（CR-001 第 1、2 項；D-07）。
>
> **CR-002（D-13，2026-10-19）**：Go 1.27.2 修正 GO-2026-6603／6611／6612／6613／6617（1.27.1 受影響，govulncheck 判可達）；gosec v2.29.0 無法解碼 Go 1.27.2 export data v5（上游 issue securego/gosec#1774，截至 2026-10-19 無支援之發行版），故暫釘上游預設分支 commit。門檻（§4）、job 名（§3）、`-severity=high`／`-no-fail=false`、govulncheck `v1.8.0` 均未放寬或變更；安全部同意（附條件 C1～C8）見 [13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md](./13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md)；殘餘風險 RR-010（低）見 [05](./05-residual-risk-register-v0.1.md)。**已生效（#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）。** 依生效日起算的期限：C6 金絲雀驗證須於 **2026-10-26（UTC+8）前**完成（合併後 7 日內）；C7 每週檢查上游發行（第一次最遲 2026-10-26，其後每 7 日：11-02、11-09、11-16、11-23、11-30、12-07、12-14、12-21、12-28），上游發行支援版本後 **14 日內**以 CR 換回，最遲 **2026-12-31** 由安全部重評。

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

### 4.1 dependency-review 條款（CR-001 第 3 項；D-07；暫行條款已結束）

| 項 | 定案 |
|---|---|
| 定位 | 於 required job `sca-govulncheck` 內執行；失敗即使該 required check 失敗；主擋仍為 `govulncheck` |
| 常態 | **不接受**靜默略過 |
| 暫行（歷史，已結束） | 原條款：Dependency graph 未啟用時，可探測後略過，但必須輸出 Actions 可見 `::warning::`；不得使 `sca-govulncheck` 因略過而失敗 |
| **現行（fail-closed）** | 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束：compare API 非 200 即失敗（`::error::`＋`exit 1`）；dependency-review 每個 PR 實跑；`fail-on-severity: critical`；Dependency graph 已由使用者開啟 |
| 期限（歷史） | 原期限 min(2026-10-14, G5 送審前)（兩者取較早），已於期限內履行 |
| 合併前 | PR #2 須先補 `::warning::`（見 [08-pr-002-security-review-v0.1.md](./08-pr-002-security-review-v0.1.md) SCR-001；**已於 `3061289` 關閉**） |
| **履行狀態** | 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束。複核對象：實際合併之 PR #6 head `37d906a96f14d855dc33b742c273a0686d2a9956`（merge commit `161a4fdd0abf9b242ec22ff2c1cdba4d82ba9f86`，2026-10-07 14:04:02 UTC+8）；以實際合併 head 為準，與前次複核 `cc75d35` 之差異（env: 傳入 REPO／BASE_SHA／HEAD_SHA、註解）經安全部書面接受——無放寬、env: 為改善。已驗證：run [37579097511](https://github.com/444112029012/short-url/actions/runs/37579097511)（pull_request、head＝37d906a）四項 required checks success；log `dependency graph compare HTTP 200`；`Dependency review (PR only)` 實跑 success、critical 0。**已滿足事實：**fail-closed PR 已先於本文件合併，CI 已驗證，連結已為實際網址。安全部書面結論為事後補登；PR #6 合併時無審查部「可合併」結論（見 [05](./05-residual-risk-register-v0.1.md) RR-006 偏離紀錄）。詳見 [09](./09-cr-ci-deviations-ruling-v0.1.md) CR-001 第 3 項履行複核。 |

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

## 7. 分支保護（方案 A；已生效——D-08；Rulesets 取證見 `docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`）

> 對齊維運 `01-ci-pipeline-v0.2` §5.2；使用者 王昱殊 於 **D-08（2026-10-07）** 選定方案 A 並接受 RR-006。Rulesets 已由使用者設定；維運取證見 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)（`main-required-checks` id 24629805，無 bypass；`main-pr-review` id 24629854，bypass 僅 Repository admin、僅限 PR）。

### 7.1 Rulesets 現行設定（方案 A）

| 規則 | 內容 | Bypass |
|---|---|---|
| **規則一「main-required-checks」** | 四項 required checks：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`；**要求合併前分支與 main 同步（`strict_required_status_checks_policy=true`，D-10；已生效）**；**禁止 force push**；**禁止刪除**受保護分支 | **無任何 bypass**（資安控制核心，必須保住） |
| **規則二「main-pr-review」** | 需要 Pull Request；至少 **1** 核准 | Repository admin，模式**僅限 PR**（不得允許直接 push） |

操作注意：刪除舊 Branch protection rule，避免與 Rulesets 雙軌卡住；設定完成後維運以 API 取證。

> **D-10（王昱殊，2026-10-07，使用者回覆 t83s391）**：規則一改 `strict=true`——合併前分支須與 main 同步，舊綠燈不得用於合併（安全部建議，源自 11 號檔 SCR-011／12 號檔 SCR-012；屬 D-08 範圍）。前述 2026-10-07 Rulesets 快照記錄的是 D-10 之前的 `strict=false`。
>
> **strict=true 已生效（D-10）**：使用者約 14:33（UTC+8）於 GitHub UI 啟用；維運 14:34:10（UTC+8）唯讀讀取 ruleset 24629805，`strict_required_status_checks_policy=true`、`enforcement=active`、四項 check 名稱不變、`deletion`／`non_fast_forward` 仍在、`updated_at` 14:33:28（UTC+8）。取證檔：`docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`；原始 JSON：`docs/09-cicd-env/evidence/rulesets-2026-10-07/ruleset-24629805-d10-public-read.json`、`docs/09-cicd-env/evidence/rulesets-2026-10-07/ruleset-24629805-d10-auth-read.json`。**取證檔隨 PR #13 提交（head `6674a59`；已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`，取證檔連結已生效）。** 安全部約 14:36（UTC+8）另以公開 REST API 唯讀讀取同一 ruleset，`strict=true`、四項名稱、`updated_at` 一致。**bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）**（取證檔 `ruleset-strict-true-d10-2026-10-07.md` §3 第 6 項，依 C 列；A 列公開讀取與 B 列 `metadata=read` 讀取皆未回傳 `bypass_actors`，屬權限不足，不代表清單為空或非空）。據此，規則一「無任何 bypass」於 D-10 後仍成立，證據強度為使用者 UI 畫面、API 未能佐證；若日後以具權限讀取發現非空，依 RR-006 重評並走變更管理（見 05 v0.1.8）。（v0.3.5 時本機取證檔 B 列與第 6 項仍為佔位；v0.3.6 依隨 PR #13 提交之取證檔 v0.1 更新。）

### 7.2 補償控制（A／B 必做）

- 每個 PR 合併前：於 PR 留言貼上審查部（必要時安全部）書面「可合併」結論與四項 check 綠燈連結。  
- 合併人＝使用者本人。  
- 每筆 admin bypass 合併須可稽核（留言含審查結論連結）。  
- **G5 前**重新評估是否仍需 admin bypass（見殘餘風險 **RR-006**）。

### 7.3 方案比較（安全立場）

| 方案 | 安全立場 |
|---|---|
| **A（已採；D-08）** | 四 checks 無 bypass；僅「1 核准」對 admin 可繞過＋書面補償 |
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
| **SEC-009** | §2 gosec（CR-002 暫行例外：pseudo-version `v2.29.1-0.20261015133444-707f2271f665`，已生效（#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）；之前為 v2.29.0）；§4 SAST；job `sast-gosec` |
| **SEC-010** | §2 govulncheck v1.8.0；§4 SCA；§4.1 dependency-review（fail-closed）；job `sca-govulncheck` |
| **SEC-011** | §2 gitleaks；§4 機密；job `secrets-gitleaks` |
| **SI-05** | CI 資安掃描門檻存在並可驗證阻擋 |
| **NFR-004** | 門檻＋流水線依門檻阻擋；分支保護方案 A 已生效 §7 |

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
| dependency-review | §4.1（CR-001 第 3 項，D-07）；維運同意；已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束 |
| 合併路徑方案 A | **已採（D-08）**；§7 |
| GitHub 設定生效 | **Rulesets 已生效**（使用者設定；維運取證 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）；規則一 strict=true **已生效（D-10）**（取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）） |

---

## 12. 禁止事項

- 禁止為「先合併」而無紀錄關閉掃描規則。  
- 禁止將真實機密寫入版控後以「之後再刪」通過。  
- 禁止以編碼／更名／分割檔案繞過機密掃描當作例外。  
- 禁止將付費授權工具列為必選（零預算）。  
- 禁止維運實作放寬本文件阻擋級別而未走本文件升版＋安全部核准。  
- 禁止對規則一（required checks／禁 force push）設定 bypass。  
- 禁止 dependency-review 略過或改回 fail-open（暫行條款已結束；任何放寬須走本文件升版＋安全部核准）。
- **CR-002 例外期間**：禁止放寬 SAST 門檻（`-severity=high`、`-no-fail=false` 不得變更）、禁止將 `sast-gosec` 改為非必要檢查或 `continue-on-error`、禁止新增 `-exclude` 或未經安全部核准的 `#nosec`。
- 禁止以 `@latest`／`@master`／短 SHA 安裝掃描器，或關閉 GOSUMDB／改 `GOPROXY=direct` 安裝；禁止 Go 版本使用 `1.27.x`／`stable` 等浮動寫法。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 初版門檻 |
| v0.2 | 2026-10-07 | 安全部 | G3：寫死 Go 零預算工具；定案 job 名與分支保護 |
| v0.2.1 | 2026-10-07 | 安全部 | 記錄維運一致性簽核通過 |
| v0.3 | 2026-10-07 | 安全部 | CR-PENDING（後登記為 CR-001）：Go 1.27.1／gosec v2.29.0／govulncheck v1.8.0；dependency-review 暫行至 min(2026-10-14, G5)；分支保護方案 A 目標態；阻擋語意未放寬 |
| v0.3.1 | 2026-10-07 | 安全部 | 依 D-07／D-08 同步：狀態改已核准；CR-PENDING-01～03 改 CR-001、合併路徑改引 D-08；§7 分支保護改已生效（Rulesets 取證 PR #5）；§4.1 補履行狀態（graph 已開、compare 200；fail-closed 已提交 PR #4 未合併；期限仍 min(2026-10-14, G5)）；檔名不變 |
| v0.3.2 | 2026-10-07 | 安全部 | 維運交叉比對修正：Rulesets 證據改指 repo 檔 `docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`（保留 ruleset id）；§4.1 改寫為統一句型「已恢復 fail-closed，暫行條款結束」（含 fail-closed PR 連結佔位字串）並加合併前提；暫行／期限改列歷史；小版修訂不改檔名（檔名保留主版號） |
| v0.3.3 | 2026-10-07 | 安全部 | §4.1 同步 09：複核對象改為實際合併之 PR #6 head `37d906a`（差異經書面接受、無放寬）；前提句改寫為已滿足事實；佔位字串替換為 PR #6；小版修訂不改檔名（檔名保留主版號） |
| v0.3.4 | 2026-10-07 | 安全部 | 依 D-10：§7.1 規則一補「合併前分支須與 main 同步（strict=true）」，註明既有快照為 D-10 前之 strict=false、strict=true 取證待維運更新 `docs/09-cicd-env/evidence/` 後再引用；§11 狀態同步；小版修訂不改檔名（檔名保留主版號） |
| v0.3.5 | 2026-10-07 | 安全部 | §7.1 規則一 strict=true 改為「已生效（D-10）」：引用維運取證 `ruleset-strict-true-d10-2026-10-07.md` 與原始 JSON（public-read／auth-read），註明取證檔隨維運 PR A 提交、合併後連結生效，bypass 清單待具權限讀取確認；安全部以公開 API 獨立讀取一致；§11 狀態同步；小版修訂不改檔名（檔名保留主版號） |
| v0.3.6 | 2026-10-19 | 安全部 | §7.1 與 §11 bypass 寫法由「待具權限讀取確認」改為「bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）」，引用維運取證 `ruleset-strict-true-d10-2026-10-07.md`（§3 第 6 項，C 列）；取證檔改註「隨 PR #13 提交（截至 2026-10-19 仍 open、未合併）」，取代「維運 PR A」；規則一「無任何 bypass」於 D-10 後仍成立（證據強度：UI 畫面）；小版修訂不改檔名（檔名保留主版號） |
| v0.3.7 | 2026-10-19 | 安全部 | CR-002／D-13（安全部同意（附條件），見 13）：§2 Go toolchain 改精確 `1.27.2`（go.mod 與 setup-go 一致，禁止 `1.27.x`／`stable`）；SAST 改暫行例外——gosec `v2.29.1-0.20261015133444-707f2271f665`（commit `707f2271f66579a18e1da25794b1e8012247274b`、`h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`），CI 必須以 `go version -m` 比對 h1，期限（上游發行後 14 日內換回、最遲 2026-12-31 重評）與退場條件；§2 補 CR-002 註記；檔頭狀態標明「依 D-13，待 CR-002 PR 合併後生效」、摘要與對應 CR 補 CR-002；§9 SEC-009 同步；§12 新增例外期間禁止事項與浮動版本禁止；阻擋語意未放寬；§7.1／§11 PR #13 改為已合併事實（2026-10-19 09:05:27 UTC+8，merge commit `2e54f7e`，取證檔連結已生效）；小版修訂不改檔名（檔名保留主版號） |
| v0.3.8 | 2026-10-19 | 安全部 | CR-002 已生效：#15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`；狀態列、§2 SAST／Go toolchain 兩列、CR-002 註記與 §9 對照表改「已生效」，移除「CI 實況仍為 v2.29.0／go 1.27.1」註記（已以 cursor-github 唯讀查證 main `c341be1`、`24b8da0` 的 go.mod、workflow 與 CI log 皆為 1.27.2）；記錄合併後 main run 39189683566 四項 success、`sast-gosec` 由 `2e54f7e` failure 轉綠且 log 有 h1 OK；寫入 C6（2026-10-26 前）、C7（每週、發行後 14 日內、最遲 2026-12-31）具體期限；阻擋語意不變；小版修訂不改檔名（檔名保留主版號） |
