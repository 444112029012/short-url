# CR-006 PR #31 品保驗收取證結果（QA-CR006）

- 執行：品保部，2026-11-02 14:29～14:35（UTC+8）；GitHub 全程唯讀（未 commit／push／留言／重跑／改 ruleset）。
- 依據：CR v0.4 §5.1（L1～L7）、`impact/品保部.md` v0.4 §2（QA-CR006-1～15）。
- 對象：PR #31，run 40660317286（event pull_request，headBranch `ci/CR-006-go-1.27.3`，2026-10-28 16:55:42（UTC+8）建立）。
- 本目錄 log 有兩種格式：`api-job-<id>.log`＝per-job API 原始 log（與維運 `evidence-pr31/job-*.log` 逐位元組相同，行號與維運引用一致）；`job-<id>-<name>.log`＝`gh run view --log --job`（每行前有 job 名稱、`UNKNOWN STEP`、時間戳）。grep 皆未加 `^`。

## 0. 維運說法核對

| 說法 | 實測 | 證據 | 結果 |
|---|---|---|---|
| head＝826d5d34…1581 | 未登入 `git ls-remote` refs/pull/31/head＝826d5d34d676e11f52e204e155ea271c7d811581；未帶 token 的 API `.head.sha` 同值 | `01-ls-remote-unauth.txt` 第 2 行、`02-pr31.json` | 相符 |
| base＝ca09137 | API `.base.sha`＝ca09137ae8f3…；merge-base 同；origin/main 仍為 ca09137 | `02-pr31.json`、`13-diff-gomod-head.txt` 第 2 行、`01` 第 1 行 | 相符 |
| run 40660317286 四 job success | headSha＝826d5d34…；四 job 與所有步驟 success | `03-run-40660317286.json` | 相符 |
| SHA256SUMS sha256 b1405323…93cf | 重算 `b1405323a58c96ca541897bac1eef0ed3aae1967b38bfb96922361d01a9d93cf`；`sha256sum -c` 6 列全 OK | `00-ops-evidence-verify.txt` | 相符 |
| 維運的四個 job log | 與品保自 API 下載者逐位元組相同 | `00-ops-evidence-verify.txt` | 相符 |

## 1. 逐項判定

| 項目 | 判定 | 證據（檔案：行） |
|---|---|---|
| QA-CR006-1／L1 | **PASS** | `api-job-122721425145.log`（sast-gosec）：131 `Setup go version spec 1.27.3`、149 `go version go1.27.3 linux/amd64`、188 GOROOT `/opt/hostedtoolcache/go/1.27.3/x64`、196 `GOVERSION='go1.27.3'`；`api-job-122721425254.log`（sca）：132／150／189／197；`api-job-122721424775.log`（unit-test）：132／150／189／197。三 job 反向字樣 `go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/` 計數皆 0。記錄：`GOTOOLCHAIN='auto'`（193／194／194）。彙整 `10-checks-L1-L3.txt` |
| QA-CR006-2／L2 | **PASS** | key＝`setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d429736099b597565b998fd2f9b82c99c65d9cffc4d8c9270df02c0a025`（前綴 go-1.27.3、雜湊段＝基準）：sca 268 `Cache saved with the key:`；gosec 397、unit-test 269 為 `Failed to save: Unable to reserve cache with key <同 key>`（並行 job 搶存，屬記錄事項）。三 job 還原時皆 `Cache is not found`（148／149／149），依判準只記錄 |
| QA-CR006-3／L3 | **PASS** | `switching to go\|downloading go1\.` 計數：三 job 在 API log 與 `gh run view --log` 兩格式皆 0（`10-checks-L1-L3.txt` 末段）；與 L1 GOROOT 合併判讀一致 |
| QA-CR006-4／L4 | **PASS** | `api-job-122721425145.log` 266 `mod	github.com/securego/gosec/v2	v2.30.0	h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`、267 `gosec h1 OK`；步驟 5 success（`03`）|
| QA-CR006-5／L5 | **PASS** | 命令列 272 `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...`（`shell: bash -e`，273）；Summary 384 Files 57（掃描的 Go 檔數）、385 Lines 7485（掃描的程式行數）、386 Nosec 0（`#nosec` 抑制註記數）、387 Issues 0（severity≥high 的發現數）。維運引用的 384～387 行在 API 格式與 `gh run view` 格式兩檔都相符。exit 0：步驟 6 conclusion success（CI 成功時不印 exit code，以步驟結論為證）。基準 run 40618331206：`base-run-40618331206.log` 639～642 同為 57／7485／0／0 → **無差異** |
| L6（＝QA-CR006-10(a)） | **PASS** | `api-job-122721425254.log`：`grep -c 'No vulnerabilities found\.'`＝1（219）；`grep -c 'This scan also found'`＝0；`Your code is affected` 0、`GO-2026-6809` 0；197 `GOVERSION='go1.27.3'`；「Run govulncheck」步驟 6 success＝exit 0（`bash -e`，217）。dependency-review（步驟 8，另列）：245「did not detect any vulnerable packages with severity level "critical" or higher」|
| QA-CR006-6 | **PASS** | check-runs API（head 826d5d34）：sca-govulncheck、sast-gosec、secrets-gitleaks、unit-test 皆 completed／success（`04-check-runs-head.json`）；名稱與 ruleset 24629805 required contexts 完全相同（`05-ruleset-required-contexts.txt`）；workflow 無 `continue-on-error`（`13` 末段 grep 0 筆）。註：「兩個 run」中的合併後 run 尚未存在 |
| QA-CR006-7／L7 | **PASS** | artifact `art-pr/go-test.json`（jq）：pass 181／fail 0／skip 0；基準 `art-base/go-test.json` 同 181／0／0；`api-job-122721424775.log` 218 `go test -race -json ./...`、225 ENG-011 以 `-race` 執行、234 `--- PASS: TestENG011_PerSourceIsolationByPeerAddress` |
| QA-CR006-8 | **PASS** | `art-pr/go-test.json`：`^TestCT` pass 14、fail 0（基準同 14／0）|
| QA-CR006-9 | **PASS（PR head 部分）** | 本機 clone `/workspace/tmp-qa-pr31` @826d5d34，`go version go1.27.3 linux/amd64`、`GOTOOLCHAIN=local`；`go test -race -coverpkg=./... -coverprofile=c.out ./...` exit 0；`go tool cover -func` total **80.9%**（`14-coverage-local.txt`，profile `14-c.out`）。精算 474/586＝80.887%，顯示四捨五入為 80.9%；同法量基準 ca09137 也是 474/586＝80.9%（`15-coverage-base-ca09137.txt`），**與基準完全相同**。注意：品保判準原寫 `-count=1`，本次依指示未加 `-count=1`，部分結果可能取自測試快取。merge commit 部分待合併後 |
| QA-CR006-10 | (a) **PASS**（見 L6）；(b) **PASS**；(c) **PASS**（皆限 head 826d5d34；合併後 run 的 (a) 待定） | 見 §4。(b)：安全部 `06-m5-6-project-verbose.txt` 第 19 行「scanned the following 11 modules and the go1.27.3 standard library」、第 32 行 `No vulnerabilities found.`、第 33 行 exit 0；品保重現 `23-qa-m5-6-project-verbose.txt` 第 19、32 行同，第 33 行 exit 0，GO-2026-6809 0。(c)：安全部 `05-m5-5-gosec-binary.txt` 第 5、12 行（exit 0）、第 77、113 行（verbose exit 0），第 24 行 go1.27.3 標準庫；品保重現 `22-qa-m5-5-gosec-binary.txt` 第 5、76 行，第 12、112 行兩次 exit 0，GO-2026-6809 0；內容與安全部逐行相同（`25`）。RR-012 結案仍由安全部判定 |
| QA-CR006-11 | **PASS** | govulncheck Symbol／package／module：基準 `base-run-40618331206.log` 1062～1068「No vulnerabilities found」「affected by 0」「also found 1 … in packages you import and 0 … in modules you require」＝**0／1／0**；PR `api-job-122721425254.log` 219 僅 `No vulnerabilities found.`＝**0／0／0**。差異只有 package 層 1→0（即 GO-2026-6809 消失，ID 由 CI 預設輸出無法直接看到，以安全部 §5.2 verbose 佐證）；無其他差異。dependency-review 無 Critical（245）。gosec 前後無差異（見 QA-CR006-5）|
| QA-CR006-12 | **PASS（PR head 部分）** | `13-diff-gomod-head.txt`：go.mod 第 3 行 `go 1.27.3`、無 `toolchain` 行；workflow 62／91／134 `go-version: "1.27.3"`、63／92／135 `check-latest: false`；無 Dockerfile／Containerfile／.go-version。merge commit 部分待合併後 |
| QA-CR006-13 | **待定** | ruleset 前後比對屬合併前後取證（品保本次未取合併前快照比對） |
| QA-CR006-14 | **PASS（PR 部分）** | `13-diff-gomod-head.txt`：merge-base＝ca09137；ca09137..head 只有 1 個 commit 826d5d3；numstat `.github/workflows/ci-devsecops.yml` 4/4、`go.mod` 1/1；兩檔 sha256＝`06833a31…f245`／`02842cf5…8ebf`；go.sum `1e2b43d7…` 不變。合併後 `merge^1..merge` 部分待合併後 |
| QA-CR006-15 | **待定** | PR 尚未合併（`02-pr31.json` state open、merged false）|

## 4. QA-CR006-10(b)(c)：安全部 M5 重跑核對（2026-11-02 14:35～14:37（UTC+8）追加）

對象：安全部 36 號 v0.1 `/workspace/projects/short-url/docs/06-security/36-cr-006-pr31-m5-s2-confirmation-v0.1.md`；證據 `docs/06-security/evidence/cr-006-m5/826d5d34d676e11f52e204e155ea271c7d811581/`（下稱 E）。

| 核對 | 結果 | 證據 |
|---|---|---|
| 36 號 sha256 | `b7444593cb2e8f5cafe5f0f1836c96dec15585eee1570e7facafb23bb8a5f3ec`，與轉述 b7444593… 相符 | `26-security-evidence-verify.txt` 第 2 行 |
| E/SHA256SUMS sha256 | `33c96bae0c0dc14d41e47bd3f0aff20f3c25c364d33edb0a8264be7ac19c2ef9`，相符；`sha256sum -c` 11 列全 OK、exit 0；目錄 11 檔皆列入 | `26` 第 4、6～18 行 |
| 目錄名 | 完整 40 碼 head SHA | `26` 第 5 行 |
| head 未變 | 14:35:37 重取：未登入 ls-remote 與未登入 API `.head.sha` 皆 826d5d34…1581；PR open、未合併、updated_at 未變 | `20-head-refetch.txt` 第 3、6 行 |
| M5-1 | PASS：ls-remote＝API＝維運值，checkout＝同值，go.mod `02842cf5…`、workflow `06833a31…` | E/`01-m5-1-head.txt` 第 4、12、18～21、26、29～30 行 |
| M5-2 | PASS：官方 JSON 值 65efbc80…9b40、tarball OK、go1.27.3、GOTOOLCHAIN auto | E/`02` 第 4、7、15、18～19 行 |
| M5-3 | PASS：gosec、govulncheck 安裝皆 exit 0，只有模組下載訊息 | E/`03` 第 52、61 行 |
| M5-4 | PASS：`gosec: go1.27.3`、h1 67aF6LHa… | E/`04` 第 3、5 行 |
| M5-5 | PASS：`-mode=binary` exit 0、verbose exit 0；GO-2026-6809 0（檔內唯一一次為檔尾自註「grep -c …: 0」）；module 層只有 GO-2026-6866、6860、5932（預期 3 筆，不算不通過） | E/`05` 第 5、12、24、77、85～107、113、115 行 |
| M5-6 | PASS：go1.27.3 標準庫、`No vulnerabilities found.`、exit 0、GO-2026-6809 0 | E/`06` 第 19、32、33、35 行 |
| M5-7 | PASS：02～06 計數皆 0（品保重算同；07 自身 1 次為標頭字串） | E/`07` 第 3～9 行；`26` 註 |
| M5-8 | PASS：暫存已刪、SHA256SUMS 可驗 | E/`08` 第 3～4 行；`26` |
| 35 號 §6.2 (a)(b)(c)（第 80～83 行第 4、5、7 項：40 碼子目錄／只用公開 ls-remote 與未登入 API／M5-5「必須 0」＝GO-2026-6809 0 筆） | 皆符合 | 同上 |
| QA-M7（6809 以外項目） | 專案 `./...` 三層皆 0；gosec 二進位只有預期的 module 層 3 筆，沒有 symbol 層的新項目，所以不必啟動 M7 | E/`05`、`06` |

品保重現（/workspace/tmp-qa-pr31 @826d5d34；Go：box 預裝 `/usr/local/go` go1.27.3，GOTOOLCHAIN auto、全新 GOPATH／GOCACHE；gosec v2.30.0、govulncheck v1.8.0 由品保自行 `go install`，建置版本 go1.27.3）：`21-qa-repro-setup.txt`。
- M5-5：`22-qa-m5-5-gosec-binary.txt` 兩次皆 exit 0（第 12、112 行 `exit 0`），GO-2026-6809 0，module 層同 3 筆。
- M5-6：`23-qa-m5-6-project-verbose.txt` exit 0，go1.27.3 標準庫，GO-2026-6809 0。
- 與安全部輸出比對（去除路徑、時間、指令列）：兩檔皆**相同**（`25-qa-vs-security-diff.txt`）。
- 補充 `-scan module`（非判準）：`-mode=binary -scan module` 對 gosec 為 exit 3，項目恰為 GO-2026-6866、6860、5932，符合預期，不算不通過；不加 `-mode=binary` 直接掃檔案為 exit 2（用法錯誤），不採計（`24-qa-scan-module-supplement.txt`）。
- 限制：品保使用的 Go 1.27.3 是 box 預裝版，這次沒有再對 go.dev 官方 JSON 核對 tarball；工具鏈的獨立核驗以安全部 M5-2 為準。

## 2. 與維運說法的差異
- 無實質差異：head、base、四個 job 結論、SHA256SUMS 值與內容、L4／L5 引用行號及數值都相符。
- 補充（非差異）：L2 有兩個 job 的快取存檔失敗（`Failed to save … another job may be creating this cache`），不影響判定；`gh run view --log` 的步驟欄顯示 `UNKNOWN STEP`，因此「步驟 exit 0」改以 run JSON 的步驟 conclusion 判定。
- 安全部 36 號：未發現與證據不符之處。SCR-060（未登入取 job log 回 403）經品保實測屬實（`26` 末行 403）；品保第一輪已用登入的唯讀權限取得 log，並確認與維運副本逐位元組相同，可補足 SCR-060 所說的單一來源問題。
- QA-CR006-9 覆蓋率實為 80.887%，以 `go tool cover` 顯示值 80.9% 達標，與基準 ca09137 精算值相同。

## 3. 檔案
`00-ops-evidence-verify.txt`、`01`～`05`（PR／run／check-runs／ruleset）、`10`～`13`（逐項 grep 與 diff）、`14`～`15`（覆蓋率）、`20`～`26`（10(b)(c) 核對與重現）、`api-job-*.log`、`job-*.log`、`base-run-40618331206.{json,log}`、`art-pr/`、`art-base/`（go-test-json artifact）、`SHA256SUMS`。
