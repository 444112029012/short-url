---
文件：CR-006 PR #31 合併前安全部獨立重跑（M5-1～M5-8）與 S2 書面確認
版本：v0.1
狀態：已提出（待品保部核對；合併由使用者決定）
負責角色：安全部（執行者）；核對者：品保部
對應 CR：CR-006 v0.4 §5.2（`/workspace/rebuild/cr-006/CR-006-ci-go-1.27.3-v0.4.md`，sha256 818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734，安全部 2026-11-02 14:29:21（UTC+8）以 `sha256sum` 核對相符）
綁定 head：PR #31 `826d5d34d676e11f52e204e155ea271c7d811581`（base main `ca09137ae8f336a8bb02ac7f20e893eb2a6ccfbe`）；**head 一有變動，本確認即作廢，須重跑**
CI run：40660317286（jobs：secrets-gitleaks 122721425088、sast-gosec 122721425145、sca-govulncheck 122721425254、unit-test 122721424775）
維運證據：`/workspace/rebuild/cr-006/evidence-pr31/`，SHA256SUMS sha256 b1405323a58c96ca541897bac1eef0ed3aae1967b38bfb96922361d01a9d93cf（相符），`sha256sum -c` 6 檔全部 OK（2026-11-02 14:29:21 UTC+8）
安全部證據：`docs/06-security/evidence/cr-006-m5/826d5d34d676e11f52e204e155ea271c7d811581/`，SHA256SUMS sha256 33c96bae0c0dc14d41e47bd3f0aff20f3c25c364d33edb0a8264be7ac19c2ef9（11 檔，2026-11-02 14:32:27 UTC+8 `sha256sum -c` 全部 OK）
撰寫時間：開始於 2026-11-02 14:33:34（UTC+8，取自 `date`），完成時間以檔案 mtime 為準
---

# CR-006 PR #31：M5 獨立重跑與 S2 書面確認

## 1. 結論

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

## 8. 修訂紀錄

| 版本 | 日期 | 角色 | 內容 |
|---|---|---|---|
| v0.1 | 2026-11-02 | 安全部 | 初版：對 PR #31 head `826d5d34` 獨立重跑 M5-1～M5-8，全部 PASS；L6 成立；S2 書面確認通過；SCR-060～062（資訊） |
