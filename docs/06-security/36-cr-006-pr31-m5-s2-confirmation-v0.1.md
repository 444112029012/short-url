---
文件：CR-006 PR #31 合併前安全部獨立重跑（M5-1～M5-8）與 S2 書面確認
版本：v0.1.3
狀態：已提出（v0.1.1：新 head `47475638` 通過，待品保部核對；合併由使用者決定）。v0.1 對舊 head `826d5d34` 的確認已因 head 變動作廢
負責角色：安全部（執行者）；核對者：品保部
對應 CR：CR-006 v0.4 §5.2（`/workspace/rebuild/cr-006/CR-006-ci-go-1.27.3-v0.4.md`，sha256 818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734，安全部 2026-11-02 14:29:21（UTC+8）以 `sha256sum` 核對相符）
綁定 head（v0.1.1，現行）：PR #31 `474756381bd170beb3a5d723b63f1b8ba48b2add`（merge commit，parents `826d5d34`＋main `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`；base main `164b62c9`）；**head 一有變動，本確認即作廢，須重跑**。v0.1 綁定（已作廢）：`826d5d34d676e11f52e204e155ea271c7d811581`（base `ca09137a`）
CI run（v0.1.1）：41403811744（jobs：secrets-gitleaks 125141230565、sast-gosec 125141230822、sca-govulncheck 125141230836、unit-test 125141230942）。CI run（v0.1）：40660317286（jobs：secrets-gitleaks 122721425088、sast-gosec 122721425145、sca-govulncheck 122721425254、unit-test 122721424775）
維運證據（v0.1.1）：`/workspace/rebuild/cr-006/evidence-pr31-47475638/`，SHA256SUMS sha256 45cff365e74ec3d6ba1919362c3ecbfa698132b2f42f316b6d3da0f5aa377a6b（相符），`sha256sum -c` 6 檔全部 OK（2026-11-02 16:29:20 UTC+8）。維運證據（v0.1）：`/workspace/rebuild/cr-006/evidence-pr31/`，SHA256SUMS sha256 b1405323a58c96ca541897bac1eef0ed3aae1967b38bfb96922361d01a9d93cf（相符），`sha256sum -c` 6 檔全部 OK（2026-11-02 14:29:21 UTC+8）
安全部證據（v0.1.1）：`docs/06-security/evidence/cr-006-m5/474756381bd170beb3a5d723b63f1b8ba48b2add/`，SHA256SUMS sha256 43d33bb47bed52571dff61b56776787ff4e838a8634da79565aa533122a659bd（14 檔，2026-11-02 16:32:30 UTC+8 `sha256sum -c` 全部 OK）。安全部證據（v0.1）：`docs/06-security/evidence/cr-006-m5/826d5d34d676e11f52e204e155ea271c7d811581/`，SHA256SUMS sha256 33c96bae0c0dc14d41e47bd3f0aff20f3c25c364d33edb0a8264be7ac19c2ef9（11 檔，2026-11-02 14:32:27 UTC+8 `sha256sum -c` 全部 OK）
版本來源（v0.1.3 起，回應審查部 PR33-R2）：v0.1 原文＝main `1b53fcde` 上本檔 blob（sha256 b7444593…，檔名 `-v0.1.md`）；v0.1.1 原文＝`36-cr-006-pr31-m5-s2-confirmation-v0.1.1.md`（sha256 2936baec4ff6d20474a4b31a30e9820e8a3a1177de17544518dd20179d58d598，#31 合併所依據的 S2）；v0.1.2 原文＝`36-cr-006-pr31-m5-s2-confirmation-v0.1.2.md`（sha256 2b7ade2fa8d929ffb5a242c1794605f5b4471c5a9db69ac912168d1f7b2e4b4e，PR #33 head 87943c91 上的版本）。版本快照檔一經入庫即凍結不改；本檔 `-v0.1.md` 為現行版，舊版一律另存快照
撰寫時間：v0.1.3 開始於 2026-11-03 15:51:18（UTC+8，取自 `date`）；v0.1.1 開始於 2026-11-02 16:32:34（UTC+8，取自 `date`），完成時間以檔案 mtime 為準；v0.1 開始於 2026-11-02 14:33:34（UTC+8，取自 `date`），完成時間以檔案 mtime 為準
---

# CR-006 PR #31：M5 獨立重跑與 S2 書面確認

> **v0.1.2 證據收錄注意（2026-11-03 11:44:57 UTC+8 起草）**：依安全部比照 33 號的裁定（獨立裁定檔：38 號 `38-cr-006-m5-gitleaks-false-positive-ruling-v0.1.md`），兩個 M5 證據目錄的 `01-m5-1-head.txt` 與 `04-m5-4-gosec-buildinfo.txt` 被 gitleaks `generic-api-key` 誤判（內容分別是公開 commit SHA 與 `google.golang.org/api v0.304.0` 的 h1 校驗和，非機密），**不收進 repo，原檔只留在安全部本機** `/workspace/projects/short-url/docs/06-security/evidence/cr-006-m5/`。SHA256SUMS 不改（43d33bb4…、33c96bae…）。在 repo 內驗證**必須用 `sha256sum -c --ignore-missing SHA256SUMS`**；一般 `sha256sum -c` 會因缺這 2 檔報錯，屬預期。缺檔以本機原檔對 SHA256SUMS 驗證為準。不遮蔽、不改寫證據，不改 `.gitleaksignore` 或 CI 設定。

> **v0.1.1 注意**：PR #31 已於 2026-11-02 14:55:52（UTC+8，commit 時間）以一般 merge 同步 main，head 改為 `474756381bd170beb3a5d723b63f1b8ba48b2add`。**現行結論見 §8**。§1～§7 是 v0.1 對舊 head `826d5d34` 的內容，原文保留，**只對舊 head 有效，已作廢，不得用於新 head**。

## 1. 結論（v0.1，僅對舊 head `826d5d34`，已作廢）

**通過。** M5-1～M5-8 全部 PASS；L6（含 dependency-review 無 Critical）依可取得的證據成立。安全部據此對 head `826d5d34d676e11f52e204e155ea271c7d811581` 出具 S2 書面確認。

- 本確認**只適用於上述 head**；PR 有任何新 commit（含同步 main）即作廢，須重跑 M5。
- 本確認不是合併核准。依 §5.2 須先經品保部核對 M5 輸出與 SHA256SUMS，再由審查部判「可合併」；**是否合併由使用者決定**。
- 全程只讀：未 push、未留言、未 review；未修改 `/workspace/rebuild/**`；未做任何 exploit。

## 2. 執行環境

- 工具鏈：安全部自行從 go.dev 下載 `go1.27.3.linux-amd64.tar.gz`，對照 `https://go.dev/dl/?mode=json&include=all` 的官方 sha256 驗證後，解壓到安全部自己的暫存目錄（`/tmp/sec-m5.*`），未沿用其他部門的工具目錄。
- 環境比照 CI：PATH 只含該工具鏈與系統 bin、未設 `GOTOOLCHAIN`（`go env GOTOOLCHAIN`＝`auto`）、`GOFLAGS` 空、`GOPROXY=https://proxy.golang.org,direct`、`GOSUMDB=sum.golang.org`、全新 GOPATH／GOCACHE。
- head 取得：公開 `git ls-remote`／`git clone` 與未登入 `api.github.com`；gh 未登入、未使用。
- 結束後已刪除暫存 clone 與工具鏈（`08-m5-8-cleanup.txt`）。

## 3. M5 逐項結果（時間皆 UTC+8，取自證據檔內 `date` 輸出）

| 判準 | 結果 | 實際時間 | 證據檔 | 要點 |
|---|---|---|---|---|
| M5-1 取得 head | **PASS** | 14:29:38～14:29:41 | `01-m5-1-head.txt` | (i) ls-remote `refs/pull/31/head`＝(ii) API `pulls/31 .head.sha`＝(iii) 維運值＝`826d5d34…1581`；PR open、未合併、非 draft、mergeable_state clean、base `ca09137a…`；checkout 後 `git rev-parse HEAD`＝同值；merge-base＝origin/main＝`ca09137`；go.mod `02842cf5…8ebf`、workflow `06833a31…f245`，與 CR 預期相符 |
| M5-2 工具鏈 | **PASS** | 14:29:56～14:29:59 | `02-m5-2-toolchain.txt` | JSON 官方值 `65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40`＝CR 值；`go1.27.3.linux-amd64.tar.gz: OK`；`go version go1.27.3 linux/amd64`；GOVERSION `go1.27.3`、GOTOOLCHAIN `auto`；環境中無 GOTOOLCHAIN 變數 |
| M5-3 安裝 | **PASS** | 14:30:08～14:30:45 | `03-m5-3-install.txt` | `go install …gosec@v2.30.0` exit 0；`go install …govulncheck@v1.8.0` exit 0；只有模組下載訊息 |
| M5-4 建置版本 | **PASS** | 14:30:45 | `04-m5-4-gosec-buildinfo.txt` | 第一行 `…/gosec: go1.27.3`；`mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`（h1 OK）；另記錄 govulncheck `golang.org/x/vuln v1.8.0 h1:clG4qBU6zH5VKjti8n5j8BBuYzoSha392xXMkXS351U=` |
| M5-5 gosec 二進位 | **PASS** | 14:30:45～14:30:56 | `05-m5-5-gosec-binary.txt`（補充：`05b-m5-5-project-binary.txt`） | `-mode=binary` 與 `-mode=binary -show verbose` 皆 exit 0；Symbol 層 `No vulnerabilities found.`；掃描標準庫 go1.27.3；**GO-2026-6809 計數 0**；module 層恰為預期 3 筆 GO-2026-6866、GO-2026-6860（x/net v0.61.0，修於 v0.62.0）、GO-2026-5932（x/crypto openpgp，無修補），記錄、不判失敗（SCR-054／055）。補充：另以 head 建出專案 binary `./cmd/shorturl`（`go1.27.3`），`-mode=binary -show verbose` exit 0、`No vulnerabilities found.`、GO-2026-6809 計數 0 |
| M5-6 專案 stdlib | **PASS** | 14:30:59～14:31:07 | `06-m5-6-project-verbose.txt` | `govulncheck -show verbose ./...`（未使用 `-scan module`）exit 0；「scanned the following 11 modules and the go1.27.3 standard library」；`No vulnerabilities found.`；GO-2026-6809 計數 0；專案本身三層皆 0 筆（module 層 3 筆只出現在 gosec 二進位，見 M5-5） |
| M5-7 無工具鏈切換 | **PASS** | 14:31:26 | `07-m5-7-no-toolchain-switch.txt` | 對 02～06 全部輸出 `grep -cE 'switching to go\|downloading go1\.'`（不加 `^`）：各檔 0、合計 0。註：07 檔本身因標頭印出 pattern 字串，自計為 1，不屬受檢輸出 |
| M5-8 留證 | **PASS** | 14:32:21～14:32:27 | `08-m5-8-cleanup.txt`、`SHA256SUMS` | 暫存 clone 已刪除；`sha256sum *.txt > SHA256SUMS`，`sha256sum -c` 11 檔全部 OK；SHA256SUMS sha256 `33c96bae0c0dc14d41e47bd3f0aff20f3c25c364d33edb0a8264be7ac19c2ef9` |

## 4. CI 狀態與 L6（未登入 API，`09-api-checks-l6.txt`，2026-11-02 14:31:36～14:31:49 UTC+8）

- check-runs（head `826d5d34`）共 4 筆，皆 `completed/success`：sca-govulncheck、sast-gosec、secrets-gitleaks、unit-test（2026-10-28 16:55:46 開始，最晚 16:57:59 完成，UTC+8，由 API 的 Z 時間換算）。
- actions runs（head_sha）只有 1 筆：40660317286，pull_request，completed/success；jobs 四項 ID 與任務所列一致。
- check-suites 3 筆：github-actions completed/success；cursor、render 兩個 app 為 queued（非 required checks，記錄於 SCR-062）。combined status：pending、total_count 0（無 commit status）。
- **L6 dependency-review**：check-runs 名稱不錨定 `grep -ci dependency-review` 計數 **0**，即**沒有獨立的 dependency-review check-run**。查 head 的 workflow：dependency-review 是 `sca-govulncheck` job 內「Probe dependency graph」（fail-closed）與「Dependency review (PR only)」兩個步驟（`actions/dependency-review-action@2031cfc0…`，fail-on-severity: critical）。該 job 經 API 確認為 success；維運 log 副本顯示 `dependency graph compare HTTP 200`、「Dependency review did not detect any vulnerable packages with severity level "critical" or higher.」、「did not detect any denied packages」。**結論：L6 dependency-review 存在（以步驟形式），無 Critical，成立。**
- L6 govulncheck 部分（維運 log 副本，不錨定）：`No vulnerabilities found.` 1、`This scan also found` 0、GO-ID 0；同 job GOVERSION `go1.27.3`。成立。

## 5. 維運所述 L1／L3／L4／L5 抽查（`10-ops-log-spotcheck.txt`）

公開 job log 以未登入 API 取得回 **403「Must have admin rights to Repository.」**，故無法獨立取得 log，只能抽查維運副本（其 SHA256SUMS 已驗證）。

| 項目 | sast-gosec | sca-govulncheck | unit-test | 判讀 |
|---|---|---|---|---|
| L1 `spec 1.27.3`／`go1.27.3 linux/amd64`／`GOVERSION='go1.27.3'`／GOROOT 1.27.3 | 1/1/1/1 | 1/1/1/1 | 1/1/1/1 | 成立 |
| L1 1.27.2 字樣（應 0） | 0 | 0 | 0 | 成立；GOTOOLCHAIN 記錄為 `auto` |
| L3 `switching to go\|downloading go1\.`（不加 `^`） | 0 | 0 | 0 | 成立 |
| L4 h1 行與 `gosec h1 OK` | 有 | — | — | 成立 |
| L5 命令列不變；Files 57、Lines 7485、Nosec 0、Issues 0 | 相符 | — | — | 成立（與 merge-base 基準 57／7485／0／0 相同） |

## 6. 發現項

| 編號 | 嚴重度 | 內容／影響 | 重現 | 修補建議 |
|---|---|---|---|---|
| SCR-060 | 資訊 | 公開 job log 未登入不可取（403），L1～L6 的 log 內容只能依維運副本，無法由安全部獨立取得；check-run 結論則已由 API 獨立確認。影響：L 系列 log 判讀仍有單一來源 | `curl -s -o /dev/null -w '%{http_code}' https://api.github.com/repos/444112029012/short-url/actions/jobs/122721425254/logs` → 403 | 品保核對時，由有權限的角色另取一次 log 與維運副本比對 sha256；或於使用者同意下讓安全部以唯讀權杖取得 |
| SCR-061 | 資訊 | CR §5.1 L6 寫「PR 上 dependency-review 不得有 Critical」，但 dependency-review 不是獨立 check-run，而是 sca-govulncheck 的步驟；只查 check-runs 會得到「不存在」的誤判 | 見 §4 | 建議維運部於 CR 下一版 L6 註明 dependency-review 的位置（sca-govulncheck job 的步驟）與判讀方式 |
| SCR-062 | 資訊 | head 的 check-suites 有 cursor、render 兩個 app 停在 queued，combined status 為 pending（0 筆）。不屬 required checks，不影響本判定 | `curl …/commits/826d5d34…/check-suites` | 審查部判「可合併」時確認分支保護 required checks 只含四項；如有需要再處理這兩個 app |

編號說明：SCR-050～059 區段原由 21 號檔保留（35 號已用 053～055），為避免衝突，本文件使用 SCR-060～062。

本次無 Low 以上發現項；module 層 3 筆屬 35 號已登錄的 SCR-054／055，不新增。

## 7. 與任務說明的差異（照實記錄）

1. 任務說明 M5-5 寫「在 head checkout 建出 binary」；CR v0.4 §5.2 M5-5 的對象是 gosec 二進位。本次兩者都做，以 gosec 二進位為主要判準，專案 binary 為補充，兩者皆 PASS。
2. 任務說明 M5-6 提到 module 層 3 筆屬預期；實際上專案 `./...` 掃描三層皆 0 筆，3 筆只出現在 M5-5 的 gosec 二進位。不影響判定。
3. 10 號證據檔中的 dependency-review critical 自動計數原先誤列為 1（把設定值 `fail-on-severity: critical` 計入），已於 2026-11-02 14:32:27（UTC+8）更正為文字判讀，之後才產生最終 SHA256SUMS。

## 8. v0.1.1：新 head `47475638` 重跑與 S2 書面確認（現行）

### 8.1 結論

**通過。** 對 head `474756381bd170beb3a5d723b63f1b8ba48b2add` 重跑 M5-1～M5-8，全部 PASS；合併差異審查顯示從 main 帶進的只有 `docs/` 下 5 個文件檔，go.mod、go.sum、workflow 與 CR patch／proposed.yml 一致、未被改動；L1、L3～L6 抽查成立。安全部據此對新 head 出具 S2 書面確認。

- 只適用於 `474756381bd170beb3a5d723b63f1b8ba48b2add`；再有任何 commit 即作廢。
- 不是合併核准：須品保部核對、審查部判「可合併」，**是否合併由使用者決定**。
- 全程只讀：未 push、未留言、未 review；未修改 `/workspace/rebuild/**`；未做 exploit。
- 前置：CR-006 v0.4 sha256 仍為 `818f48eb…9734`（16:29:20 核對）；維運證據 SHA256SUMS 相符且 `-c` 全 OK。

### 8.2 M5 逐項結果（2026-11-02，UTC+8，取自證據檔內 `date`）

證據目錄：`evidence/cr-006-m5/474756381bd170beb3a5d723b63f1b8ba48b2add/`

| 判準 | 結果 | 實際時間 | 證據檔 | 要點 |
|---|---|---|---|---|
| M5-1 取得 head | **PASS** | 16:29:30～16:29:41 | `01-m5-1-head.txt` | ls-remote＝API＝維運值＝`47475638…2add`；PR open、未合併、clean、base `164b62c9…`；`rev-parse HEAD` 相符；parents `826d5d34`＋`164b62c9`；merge-base＝origin/main＝`164b62c9`；go.mod `02842cf5…8ebf`、workflow `06833a31…f245`（皆與 CR 預期相符）、go.sum `1e2b43d7…dd72`（與舊 head 相同） |
| M5-2 工具鏈 | **PASS** | 16:30:16～16:30:19 | `02-m5-2-toolchain.txt` | 官方 JSON 值＝CR 值 `65efbc80…9b40`；tarball OK；`go version go1.27.3 linux/amd64`；GOTOOLCHAIN 未設，`go env`＝`auto` |
| M5-3 安裝 | **PASS** | 16:30:19～16:30:55 | `03-m5-3-install.txt` | gosec v2.30.0、govulncheck v1.8.0 皆 exit 0 |
| M5-4 建置版本 | **PASS** | 16:30:55 | `04-m5-4-gosec-buildinfo.txt` | 第一行 `…/gosec: go1.27.3`；h1 `67aF6LHa…Dx0=` OK |
| M5-5 gosec 二進位 | **PASS** | 16:30:55～16:31:06 | `05-m5-5-gosec-binary.txt`（補充 `05b-…`，16:31:06～16:31:25） | 兩次掃描 exit 0；Symbol 層 `No vulnerabilities found.`；stdlib go1.27.3；**GO-2026-6809 計數 0**；module 層恰為預期 3 筆 GO-2026-6866／6860／5932（記錄、不判失敗，SCR-054／055）。補充：專案 binary（go1.27.3）exit 0、無發現、6809 計數 0 |
| M5-6 專案 stdlib | **PASS** | 16:31:25～16:31:34 | `06-m5-6-project-verbose.txt` | `-show verbose ./...`（未用 `-scan module`）exit 0；11 modules＋go1.27.3 standard library；`No vulnerabilities found.`；6809 計數 0 |
| M5-7 無工具鏈切換 | **PASS** | 16:31:34 | `07-m5-7-no-toolchain-switch.txt` | 02～06（含 05b）各檔 0、合計 0（不加 `^`） |
| M5-8 留證 | **PASS** | 16:32:29～16:32:30 | `08-m5-8-cleanup.txt`、`SHA256SUMS` | 暫存 clone 已刪；14 檔 `sha256sum -c` 全 OK；SHA256SUMS sha256 `43d33bb47bed52571dff61b56776787ff4e838a8634da79565aa533122a659bd` |

### 8.3 合併差異審查（`11-merge-diff-stat.txt`、`12-merge-diff-full.txt`、`13-merge-diff-security-review.txt`，16:31:46～16:31:54）

- `git diff 826d5d34..47475638 --stat`：5 檔、+264／−2，**全部在 `docs/` 下**：`docs/06-security/34-…`、`35-…`、`36-…`（新增）、`docs/06-security/README.md`（v0.24→v0.30）、`docs/99-changelog/decision-log.md`（+1 列 D-20）。來源為 main 的 PR #32（`fb270d9`，merge `164b62c`）。與 `git diff ca09137..164b62c --stat` 完全相同，merge 沒有引入衝突解決的額外變更。
- `git diff 164b62c..47475638 --stat`：只剩 CR 本身的 `.github/workflows/ci-devsecops.yml`（+4／−4）與 `go.mod`（+1／−1）。
- `git diff 826d5d34..47475638 -- go.mod go.sum .github/`：**空**。
- 與 CR 一致性：head workflow 與 `ci-devsecops.proposed.yml` 逐位元相同（sha `06833a31…f245`）；於 `ca09137` 套用 `CR-006-ci-go-1.27.3.patch`（sha `c3939e62…eaec`）exit 0，結果的 workflow、go.mod 與 head 逐位元相同；go.sum 與舊 head 相同。
- **安全判讀**：帶進的變更無任何程式碼、依賴、CI、建置設定，不影響攻擊面與 M5 結果；新增行機密樣式掃描 0 筆；新增外部連結只有 api.github.com、go.dev、proxy.golang.org。帶進的 34／35／36 號與 06-security README 與本機副本 sha256 相同（36 號為 v0.1 原文 `b7444593…`）。**無安全疑慮。**

### 8.4 CI 狀態與 L1、L3～L6（`09-api-checks-l6.txt` 16:32:06～16:32:18；`10-ops-log-spotcheck.txt`）

- 未登入 API：head 的 check-runs 4 筆皆 completed/success（2026-11-02 15:03:28～15:05:20 UTC+8，由 Z 時間換算）；runs 只有 41403811744（pull_request，success）；jobs ID 與任務所列一致。check-suites：github-actions success；cursor、render queued；combined status pending／0 筆。
- L6 dependency-review：check-run 名稱不錨定計數 0（非獨立 check-run，見 SCR-061）；在 sca-govulncheck 步驟中，`dependency graph compare HTTP 200`、`fail-on-severity: critical`、「did not detect any vulnerable packages with severity level "critical" or higher」、「did not detect any denied packages」。**存在且無 Critical，成立。**
- 指定行抽查（維運副本；公開 log 未登入仍 403）：

| 位置 | 內容 | 對應 | 判讀 |
|---|---|---|---|
| gosec:157 | `go version go1.27.3 linux/amd64` | L1 | 成立 |
| gosec:204 | `GOVERSION='go1.27.3'` | L1 | 成立 |
| gosec:273 | `mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa…Dx0=` | L4 | 成立 |
| gosec:274 | `gosec h1 OK` | L4 | 成立 |
| gosec:279 | `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...` | L5 命令列 | 成立 |
| gosec:391～394 | Files 57／Lines 7485／Nosec 0／Issues 0 | L5 | 成立（帶進只有 docs，與基準相同） |
| govulncheck:221 | `No vulnerabilities found.` | L6 | 成立 |

- 全檔補充：三個 job 的 L1 四項各 1、1.27.2 字樣 0、GOTOOLCHAIN `auto`；L3 各 0；L6 `This scan also found` 0（不錨定）、GO ID 0。未抽查 L2、L7（不在本次指示範圍）。

### 8.5 發現項狀態

| 編號 | 嚴重度 | 狀態（v0.1.1） | 說明 |
|---|---|---|---|
| SCR-060 | 資訊 | **仍開放** | 新 run 125141230836 的 log 未登入仍 403，L 系列 log 仍只有維運單一來源 |
| SCR-061 | 資訊 | **仍開放** | dependency-review 仍不是獨立 check-run；CR v0.4 未改 |
| SCR-062 | 資訊 | **仍開放** | 新 head 的 cursor、render check-suite 仍 queued、combined status pending／0 筆 |
| SCR-063 | 資訊 | 新增 | 本機工作副本 `docs/99-changelog/decision-log.md`（mtime 2026-10-22 10:09）落後 repo：缺 D-15 確認、D-17～D-20 共 5 列。非安全問題，也不影響本 PR；建議總協調以 repo 版同步本機副本，避免後續文件依舊版判讀。重現：比對 head 版與本機檔 sha256（`13-…` 檔） |

無 Low 以上新發現。

## 9. 修訂紀錄

| 版本 | 日期 | 角色 | 內容 |
|---|---|---|---|
| v0.1 | 2026-11-02 | 安全部 | 初版：對 PR #31 head `826d5d34` 獨立重跑 M5-1～M5-8，全部 PASS；L6 成立；S2 書面確認通過；SCR-060～062（資訊） |
| v0.1.1 | 2026-11-02 | 安全部 | PR #31 同步 main 後 head 改為 `474756381bd170beb3a5d723b63f1b8ba48b2add`：新增 §8（修訂紀錄改為 §9；M5-1～M5-8 全 PASS、合併差異審查只含 docs 5 檔、L1／L3～L6 抽查成立、S2 通過綁定新 head）；檔頭改綁新 head；§1～§7 保留並註明只對舊 head `826d5d34` 有效、已作廢；SCR-060～062 仍開放，新增 SCR-063（資訊）。小版修訂不改檔名 |
| v0.1.2 | 2026-11-03 | 安全部 | 新增證據收錄注意：兩個 M5 目錄的 01、04 號檔因 gitleaks 誤判只留本機，repo 內驗證須加 `--ignore-missing`；結論與綁定 head 不變。小版修訂不改檔名 |
| v0.1.3 | 2026-11-03 | 安全部 | 回應審查部 PR33-R2：v0.1.1、v0.1.2 原文另存快照檔入庫，檔頭新增「版本來源」；回應 PR33-S3：證據收錄注意改引 38 號獨立裁定。結論與綁定 head 不變。依審查部要求，36 號自本版起舊版一律另存快照檔，與 README「小版不改檔名」規則的差異待使用者拍板 |
