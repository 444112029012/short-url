---
文件：CI 偏離與合併路徑變更裁示（安全部）
版本：v0.1.4
狀態：已核准（D-07、D-08）
負責角色：安全部
最後更新：2026-10-07
對應需求：SEC-009、SEC-010、NFR-004、SI-05；對齊 `03-ci-security-gates-v0.2`→v0.3
專案代號：SHORTURL
參考：PR #2 head `ccd35887ab954a8d213640ee74ea9e7750030bcf`；維運 `09-cicd-env/05-pr2-ci-change-review-v0.1.md`、`01-ci-pipeline-v0.2.md`
---

# CI 偏離與合併路徑變更裁示（安全部 v0.1）

> CR 編號：(a)(b) 三項登記為 **CR-001**（[`docs/99-changelog/CR-001-go-toolchain-ci-scanners-v0.1.md`](../99-changelog/CR-001-go-toolchain-ci-scanners-v0.1.md)；原暫編 01～03 → CR-001 第 1～3 項），使用者 王昱殊 於 **D-07（2026-10-07）** 核准；合併路徑（原暫編 04）decision-log 未另立 CR 編號，依 **D-08（2026-10-07）**：選 A 並接受 RR-006（中）。  
> 僅填「資安」影響分析列與安全建議／裁示；其他角色列標「待各角色」。  
> 維運已於 2026-10-07 同意本文件 (a)(b) 全文（見維運 05 意見書）。

---

# CR-001 第 1 項：Go toolchain 1.24.x → 1.27.1
- 提出者／日期：研發（PR #2）／2026-10-07
- 類型：一般變更（未改技術選型 D-04 Go＋chi，僅版本線）
- 變更內容：現況任務／文件假設 Go 1.24.x → 期望 `go.mod` `go 1.27.1`、CI `setup-go` `1.27.x`
- 原因與效益：1.24 已停止修補；govulncheck 對舊線標準庫仍標可觸及漏洞；修補僅在受支援版本
- 不變更的後果：SCA 長期紅燈或被迫放寬門檻
- 影響的需求編號：NFR-004、SEC-010；SI-05
- 期望完成時間：隨 PR #2
- 附件／參考：https://go.dev/doc/devel/release（go1.27.1 released 2026-09-01；現行支援線含 1.26／1.27）
---
## 影響分析（由各角色填寫）
| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | 待各角色 | 待各角色 | 待各角色 | 企劃 |
| 架構／資料／API／UI | 待各角色 | ADR／技術文件版本字樣 | 待各角色 | 設計 |
| 實作 | 待各角色 | `go.mod`、CI setup-go | 待各角色 | 研發 |
| 資安（威脅模型、SEC） | **強化**：離開已停修線，降低標準庫已知漏洞暴露；不改變 SEC-009／010 阻擋語意 | `03-ci-security-gates`→v0.3；維運／ADR 版本字樣；PR 標本 CR | 低 | 安全 |
| 環境／部署／監控 | 待各角色 | 維運 CI／環境文件 | 待各角色 | 維運 |
| 測試與追溯矩陣 | 待各角色 | 待各角色 | 待各角色 | 品保 |
---
- 風險評估：低（安全強化）；須保持 `go.mod` 與 CI 一致
- 回滾方式：改回受支援線之另一受支援版本（不可退回已 EOL 的 1.24 作為長期態）
- 建議：**接受／核准**
- 核准者／日期／結論：安全部 2026-10-07 **核准**；核准矩陣＝環境／工具鏈（維運當責＋安全評估）。條件：`go.mod` go（／toolchain）指示與 CI `setup-go` 一致（已核：`go 1.27.1` ↔ `1.27.x`）；維運文件與 ADR／技術文件版本字樣更新；PR 標 CR 編號。維運同意。**使用者 王昱殊 D-07（2026-10-07）核准。**

---

# CR-001 第 2 項：govulncheck／gosec 釘版升級
- 提出者／日期：研發（PR #2）／2026-10-07
- 類型：標準變更（掃描器修補／相容升版；job 名與 High／Critical 阻擋不變）
- 變更內容：govulncheck `v1.1.3`→`v1.8.0`；gosec `v2.21.4`→`v2.29.0`
- 原因與效益：舊版無法以 Go 1.25+ 編譯；配合 toolchain 升版
- 不變更的後後果：無法在 CI 執行 SAST／SCA
- 影響的需求編號：SEC-009、SEC-010
- 期望完成時間：隨 PR #2
- 附件／參考：查證存在——govulncheck v1.8.0（pkg.go.dev／go.googlesource.com/vuln）；gosec v2.29.0（github.com/securego/gosec releases，2026-08-26）
---
## 影響分析（由各角色填寫）
| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | 待各角色 | 待各角色 | 待各角色 | 企劃 |
| 架構／資料／API／UI | 待各角色 | 待各角色 | 待各角色 | 設計 |
| 實作 | 待各角色 | workflow 釘版 | 待各角色 | 研發 |
| 資安（威脅模型、SEC） | 工具升級；**阻擋語意不變**（gosec `-severity=high`；govulncheck 非零即擋） | 門檻文件 v0.3 記錄新版本；維運同步 | 低 | 安全 |
| 環境／部署／監控 | 待各角色 | `01-ci-pipeline` 版本表 | 待各角色 | 維運 |
| 測試與追溯矩陣 | 待各角色 | 待各角色 | 待各角色 | 品保 |
---
- 風險評估：低
- 回滾方式：改回可編譯之精確舊版（若 toolchain 允許）或改釘另一精確版
- 建議：**接受／核准**
- 核准者／日期／結論：安全部 2026-10-07 **核准**。條件：仍釘精確版本（不可 `@latest`）；若改以 action 形式引用建議釘 SHA；安全門檻升 v0.3；維運同步。維運同意。**使用者 王昱殊 D-07（2026-10-07）核准。**

---

# CR-001 第 3 項：dependency-review 在 dependency graph 未開時略過
- 提出者／日期：研發（PR #2）／2026-10-07
- 類型：一般變更（輔助控制暫行 fail-open）
- 變更內容：現況期望 PR 上 dependency-review 可跑 → 實際因 graph 未開失敗；改為 compare API 200 才跑，否則略過且不讓 required job 失敗
- 原因與效益：避免輔助步驟拖垮 `sca-govulncheck` required check；主擋仍 govulncheck
- 不變更的後果：required job 因「Dependency review is not supported」失敗，阻礙合法合併
- 影響的需求編號：SEC-010（輔助路徑）；門檻 v0.2 已列 dependency-review 為非 required
- 期望完成時間：暫行至 **min(2026-10-14, G5 送審前)**
- 附件／參考：維運 05；pipeline v0.2 §3.3
---
## 影響分析（由各角色填寫）
| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | 待各角色 | 待各角色 | 待各角色 | 企劃 |
| 架構／資料／API／UI | 待各角色 | 待各角色 | 待各角色 | 設計 |
| 實作 | 待各角色 | workflow 探測＋`::warning::`；開啟 graph 後 fail-closed | 待各角色 | 研發 |
| 資安（威脅模型、SEC） | **不接受常態略過**。輔助控制 fail-open 有靜默失效風險；主 required `sca-govulncheck` 未放寬故非高風險控制減弱。裁示：**附條件暫行** | 門檻 v0.3 暫行條款；缺陷／追蹤項；使用者開 Dependency graph（**使用者動作**） | 低～中 | 安全 |
| 環境／部署／監控 | 待各角色 | 驗證 PR log 出現 HTTP 200 且 review 實跑 | 待各角色 | 維運 |
| 測試與追溯矩陣 | 待各角色 | 待各角色 | 待各角色 | 品保 |
---
- 風險評估：中以下（輔助非 required；主擋仍在）
- 回滾方式：開啟 graph 後改 API 非 200 即 fail（fail-closed）
- 建議：**附條件暫行接受**（非法常態）
- 核准者／日期／結論：安全部 2026-10-07 **附條件同意暫行**。條件：
  1. **(a)** 由**使用者**在 repo 開啟 Dependency graph（免費；Settings → Advanced Security／Code security）。開啟後 dependency-review 恢復 **fail-closed**（API 非 200 即失敗）。
  2. **(b)** 暫行期略過必須產生可見 **`::warning::`**（不可靜默）。**PR #2 已於 commit `3061289` 補上**（SCR-001 已關；run 37573748568 Annotations 可見）。`sca-govulncheck` 仍為 required 且不得放寬。
  3. **(c)** 期限：暫行條款最遲至 **min(2026-10-14, G5 送審前)** 失效（兩者取較早）；到期未改回 fail-closed 則阻擋 G5／升級追蹤。
  4. **(d)** 記缺陷／追蹤項；不屬控制減弱至高風險，安全可核准暫行。  
  維運同意；建議暫行例外至 2026-10-14 與本裁示一致。**使用者 王昱殊 D-07（2026-10-07）核准**（decision-log 原文：「dependency-review 暫行略過至 min(2026-10-14,G5前)；須 ::warning::」）。

### 履行複核（安全部；複核對象：PR #6 head `37d906a96f14d855dc33b742c273a0686d2a9956`）

> 複核時序：安全部於 PR #6 合併前即開始對 head `37d906a` 進行技術複核（複核查詢時 PR #6 仍為 open、未合併），技術結論已送維運與總協調；**本書面結論為事後補登**（2026-10-07）。PR #6 於 2026-10-07 14:04:02 UTC+8 合併，合併時 PR 上無審查部「可合併」結論；審查部於 14:06:36 UTC+8 留言事後補審（G4-PR6），不構成「可合併」結論。本文件不表示 PR #6 係依 D-08 程序合併。  
> 複核原則：以**實際合併 head** 為準；與前次複核版本的差異須經安全部書面接受（見下方差異表）。

#### 證據

| 檢查項 | 結果 | 證據 |
|---|---|---|
| PR 與 base | 符合 | [PR #6](https://github.com/444112029012/short-url/pull/6)：base `main` `a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`（PR #2 merge commit，2026-10-07 13:56:49 UTC+8）；head `37d906a96f14d855dc33b742c273a0686d2a9956`（1 commit）；只改 `.github/workflows/ci-devsecops.yml`（+11/−11） |
| CI run 對應 head | 符合 | run [37579097511](https://github.com/444112029012/short-url/actions/runs/37579097511)：event `pull_request`、head_sha＝`37d906a…`、attempt 1、2026-10-07 13:59:39–14:00:51 UTC+8、conclusion success |
| 四項 required checks | 全部 success | `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`（head `37d906a` check runs 共 4 筆，皆屬 run 37579097511） |
| compare API＝200 | 是 | `sca-govulncheck` job 112654480938 log：`dependency graph compare HTTP 200`（14:00:01 UTC+8）；env 印出 `BASE_SHA: a8246ce…`、`HEAD_SHA: 37d906a…` |
| dependency-review 實跑 | 是 | `Dependency review (PR only)` step success（14:00:01–14:00:02 UTC+8，非 skipped）；「did not detect any vulnerable packages with severity level "critical" or higher」（critical 0） |
| PR #6 合併資訊 | 已合併 | merged_at 2026-10-07 14:04:02 UTC+8；merge commit `161a4fdd0abf9b242ec22ff2c1cdba4d82ba9f86`；merged_by `444112029012`（倉庫擁有者帳號）；PR reviews 0 筆；PR 留言 1 筆（審查部事後補審，14:06:36 UTC+8）。合併後 `main` 之 workflow 與 `37d906a` 相同；main push run [37579507811](https://github.com/444112029012/short-url/actions/runs/37579507811) 四 job success（push 事件下 probe／review 依設計 skipped） |
| fail-closed 實作 | 符合 | 非 200 → `::error::Dependency graph compare API returned HTTP ${code}; dependency-review is fail-closed (CR-001 ruling).`＋`exit 1`；fail-open `::warning::` 路徑已移除；review `if: github.event_name == 'pull_request'`（不再依 `enabled` 條件略過）；`fail-on-severity: critical`；仍在 `sca-govulncheck` job 內、四 job 名不變；`set -euo pipefail` 下 curl 失敗亦使 job 失敗 |
| 未放寬其他掃描 | 符合 | govulncheck `@v1.8.0`；gosec `@v2.29.0` `-no-fail=false -severity=high`；gitleaks 8.21.2 `--exit-code 1`；unit-test vet／test 未變 |
| Workflow 安全 | 符合 | top-level `permissions` 仍僅 `contents: read`、`pull-requests: read`；無 `pull_request_target`；`REPO`／`BASE_SHA`／`HEAD_SHA`／`GH_TOKEN` 經 `env:` 傳入，`run:` 內無 `${{ }}` 內插 |

#### 差異表：前次複核 `cc75d35`（PR #4）→ 實際合併 `37d906a`（PR #6）

比對方式：以 cursor-github connector 實際取得兩版 `.github/workflows/ci-devsecops.yml` 全文逐行比對（`cc75d35` 仍可取得）；兩 PR diff 之 base blob 同為 `74d1e10`（PR #4 → `494833b`、PR #6 → `04108e1`）。

| # | 差異 | `cc75d35` | `37d906a` | 資安影響 | 判定 |
|---|---|---|---|---|---|
| 1 | Probe 步驟 `env:` | 僅 `GH_TOKEN` | 新增 `REPO`、`BASE_SHA`、`HEAD_SHA`（`GH_TOKEN` 兩版皆已經 env 傳入） | 降低 `run:` 內插 github context 的指令注入面 | **改善** |
| 2 | curl URL | `run:` 內直接內插 `${{ github.repository }}`、`${{ github.event.pull_request.base.sha }}`、`${{ …head.sha }}` | 改用 shell 變數 `${REPO}`／`${BASE_SHA}`／`${HEAD_SHA}` | 同上；請求目標與語意不變 | **改善** |
| 3 | 檔頭門檻註解 | 引 `03-ci-security-gates-v0.2.md` | 引 `03-ci-security-gates-v0.3.md` | 註解，無執行影響 | 無影響 |
| 4 | probe 前註解 | 2 行 | 3 行（新增「SHA 與 repository 經 env 傳入…」；標點微調） | 註解，無執行影響 | 無影響 |
| — | `::error::` 訊息文字 | — | — | 兩版**逐字相同**（實際比對未見差異） | 無差異 |

其餘內容（觸發條件、permissions、四 job、掃描器釘版與參數、fail-closed 分支、review 條件、`fail-on-severity`）兩版相同。**未發現任何放寬。**

#### 結論

差異可接受（無放寬、env: 為改善）；PR #6（`37d906a`）符合 CR-001 第 3 項裁示；已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束（原期限 min(2026-10-14, G5)，已於期限內履行）。

**滿足事實（取代前版之前提句）：** fail-closed PR（[PR #6](https://github.com/444112029012/short-url/pull/6)）已於本文件之前合併入 `main`（`161a4fd`）；其 CI 四項全綠、log 含 compare HTTP 200、dependency-review 實跑，均已於上表驗證；本文件內 fail-closed PR 連結已為實際網址。

#### 發現項

| 嚴重度 | 項目 |
|---|---|
| 嚴重／高／中 | 無 |
| 低 | 第三方 action 仍用 major tag（`actions/checkout@v4`、`setup-go@v5`、`dependency-review-action@v4`）——沿用 SCR-004，不擋 |
| 資訊 | 檔頭「設計」連結仍指 `docs/09-cicd-env/01-ci-pipeline-v0.1.md`（現行 v0.2） |
| 資訊 | probe 於 200 時仍寫 `enabled=true`，已無步驟讀取（無害殘留輸出） |
| 資訊 | 非 200 → `exit 1` 路徑為程式碼審閱確認，本次 run 走 200 未實際觸發 |

#### 前次複核紀錄（保留）

前次（2026-10-07，PR #4 head `cc75d3536602f569bbddece8d46cda65036d72e1`，base 為 PR #2 分支、未觸發 CI）：workflow 內容符合裁示；Dependency graph 已開之證據取自 PR #2 run [37575622585](https://github.com/444112029012/short-url/actions/runs/37575622585)（13:18:29 UTC+8 `dependency graph compare HTTP 200`、review 實跑 success）。PR #4 未合併，已由 PR #6 取代。

簽署：安全部，2026-10-07。

---

# D-08：單人擁有者合併路徑（Rulesets 方案 A／B／C；原暫編 04）
- 提出者／日期：維運／總協調／2026-10-07
- 類型：一般變更（分支保護「1 核准」對 admin 之可繞過／補償；屬資安控制調整）
- 變更內容：現況 `main` 需 PR＋1 核准＋enforce_admins → 作者＝擁有者無法自核，PR 永久 blocked。期望可合併之路徑 A／B／C（見維運 pipeline §5.2）
- 原因與效益：零預算單人 repo 解卡；保留四 required checks
- 不變更的後果：所有 PR（含 #1、#2）無法合併
- 影響的需求編號：NFR-004、SI-05；分支保護定案；殘餘風險 RR-006
- 期望完成時間：使用者於 Settings 操作（**已完成**；Rulesets 取證見 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）
- 附件／參考：`09-cicd-env/05-pr2-ci-change-review-v0.1.md` §3；`01-ci-pipeline-v0.2.md` §5.2
---
## 影響分析（由各角色填寫）
| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | 待各角色 | 待各角色 | 待各角色 | 企劃 |
| 架構／資料／API／UI | 待各角色 | 待各角色 | 待各角色 | 設計 |
| 實作 | 待各角色 | 無程式碼；GitHub Settings／Rulesets | 待各角色 | 研發 |
| 資安（威脅模型、SEC） | 見下方裁示；新增 **RR-006**（中） | 門檻 v0.3 §7 分支保護（方案 A）；殘餘風險表；CI／分支保護文件 | 中 | 安全 |
| 環境／部署／監控 | 待各角色 | Rulesets 實作＋G3 證據 | 待各角色 | 維運 |
| 測試與追溯矩陣 | 待各角色 | 待各角色 | 待各角色 | 品保 |
---
- 風險評估：中（admin 可無同儕平台核准即合併）；四 checks 無 bypass 則核心掃描控制未減弱
- 回滾方式：移除 admin bypass 或改回 enforce 全規則；或改採方案 C
- 建議：**建議方案 A**
- 核准者／日期／結論：安全部 2026-10-07 裁示如下；**使用者 王昱殊 D-08（2026-10-07）核准：已採 A，已生效**，並接受 RR-006（中），G5 前再評。Rulesets 由使用者設定，維運取證見 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)（`main-required-checks` id 24629805 無 bypass；`main-pr-review` id 24629854 僅 Repository admin、僅限 PR）。

### 合併路徑安全裁示（已採 A，已生效——D-08）

| 方案 | 安全立場 |
|---|---|
| **A. Rulesets 拆兩條（已採；D-08）** | **已採納並生效**。規則一＝四 required checks＋禁 force push／刪除，**無 bypass**（資安控制核心）。規則二＝需 PR＋1 核准，bypass＝Repository admin、**僅限 PR**。補償：合併前 PR 留言附審查部「可合併」＋綠燈連結，使用者本人按合併。足夠作為暫行補償。 |
| **B. 核准數改 0** | **不建議**。平台無核准紀錄，補償偏弱，控制減弱較大。 |
| **C. 第二位真人 collaborator** | **理想**但非零預算必達；若使用者能邀 collaborator 則**優先於 A**，否則採 A。 |

**方案 A 必備條件：**
1. 寫入 CI／分支保護文件與殘餘風險表（**RR-006**，嚴重度**中**：admin 可無同儕平台核准即合併）。
2. bypass **僅限 PR 合併路徑**；不得允許直接 push／force push（規則一無 bypass 必須保住）。
3. 每筆 admin bypass 合併須在 PR 留言留審查部結論連結（可稽核）。
4. **G5 前**重新評估是否仍需 bypass。

A／B 屬「資安控制減弱」中風險以下可由安全核准；因涉及分支保護定案與使用者操作，已呈使用者，**使用者於 D-08 選 A**。

---

## 總表（安全裁示一覽）

| 項 | 裁示 | 維運 |
|---|---|---|
| (a) Go 1.27.1＋掃描器升版（CR-001 第 1、2 項） | **核准**；使用者 D-07 核准 | 已同意 |
| (b) dependency-review 暫行略過（CR-001 第 3 項） | 原裁示附條件暫行（::warning::；原期限 min(2026-10-14, G5)）；使用者 D-07 核准；Dependency graph 已開，已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束 | 已同意 |
| (c) 合併路徑（D-08） | **已採 A，已生效**；RR-006 使用者接受 | 維運首選 A |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | 三項 CI 偏離＋合併路徑 A／B／C 裁示；維運同意 (a)(b) |
| v0.1.1 | 2026-10-07 | 安全部 | SCR-001／::warning:: 已於 PR head `3061289` 關閉 |
| v0.1.2 | 2026-10-07 | 安全部 | 依 D-07／D-08：狀態改已核准；CR-PENDING-01～03 改 CR-001 第 1～3 項、CR-PENDING-04 改引 D-08；各項補使用者核准；合併路徑改已採 A 已生效；CR-001 第 3 項新增 PR #4 fail-closed 履行複核（內容符合、尚未合併） |
| v0.1.3 | 2026-10-07 | 安全部 | 維運交叉比對修正：Rulesets 證據改指 repo 檔 `docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`（保留 ruleset id）；CR-001 第 3 項履行複核結論改為統一句型「已恢復 fail-closed，暫行條款結束」（含 fail-closed PR 連結佔位字串）並加合併前提（保留已驗證事實）；小版修訂不改檔名（檔名保留主版號） |
| v0.1.4 | 2026-10-07 | 安全部 | CR-001 第 3 項履行複核改以實際合併之 PR #6 head `37d906a` 為對象（保留 `cc75d35` 為前次紀錄）；「內容須相同」改為以實際合併 head 為準、差異須經安全部書面接受；新增 cc75d35→37d906a 差異表（實際比對，無放寬）；前提句改寫為已滿足事實；佔位字串替換為 PR #6；註明書面結論為事後補登、PR #6 合併時無審查部「可合併」結論；小版修訂不改檔名（檔名保留主版號） |
