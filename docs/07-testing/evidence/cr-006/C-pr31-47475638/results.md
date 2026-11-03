# CR-006 PR #31（新 head 474756381bd1…）品保驗收取證結果

- 執行：品保部，2026-11-02 16:29～16:36（UTC+8；§4 於 16:34 起追加）。GitHub 全程唯讀。舊目錄 `qa-evidence-pr31/` 未改動。
- 對象：head `474756381bd170beb3a5d723b63f1b8ba48b2add`、base `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`、run 41403811744（pull_request，2026-11-02 15:03:25（UTC+8）建立）。
- 基準：main 在 164b62c9 的 push run 41402171425（success，Go 1.27.2，2026-11-02 14:47:45（UTC+8）建立）。
- log：`api-job-<id>-<name>.log` 是 per-job API 原始格式，與維運副本逐位元組相同；`job-…log` 是 `gh run view --log` 格式。grep 都沒有加 `^`。

## 0. 維運說法核對（全部相符）
| 說法 | 實測 | 證據 |
|---|---|---|
| head 474756381bd1… | 未登入 ls-remote＝未登入 API `.head.sha`＝474756381bd170beb3a5d723b63f1b8ba48b2add；PR open、未合併 | `01-head-base.txt` 第 3、4 行 |
| base 164b62c9 | API `.base.sha`＝ls-remote main＝164b62c956a9…3661 | `01` 第 2、4 行 |
| run 41403811744 四項綠 | headSha＝新 head，四 job 與全部步驟 success | `03-run-41403811744.json`；`10-checks.txt`「steps」段 |
| SHA256SUMS 45cff365… | 重算相同；`sha256sum -c` 6 列 OK；四份 log 和品保從 API 下載的版本逐位元組相同 | `00-ops-evidence-verify.txt` |
| gosec 391～394 為 57／7485／0／0；govulncheck 221 為 `No vulnerabilities found.` | 屬實（API 格式與 gh 格式行號相同） | `10-checks.txt` L5、L6 段 |

## 1. main 在 ca09137 → 164b62c9 之間的變更
只有 2 個 commit：`fb270d90`（docs：CR-006 後續文件）與 merge `164b62c9`（PR #32）。改了 5 個檔案，都是文件：`docs/06-security/34-…`、`35-…`、`36-…`、`docs/06-security/README.md`、`docs/99-changelog/decision-log.md`。`.github/`、`go.mod`、`go.sum`、`*.go` 都沒有改動（`02-main-changes-and-diff.txt` 第 1～10 行）。

## 2. 逐項判定
| 項目 | 判定 | 證據 |
|---|---|---|
| QA-CR006-1／L1 | **PASS** | 三個 job 都有 `spec 1.27.3`（gosec 131、sca 132、unit 132）、`go1.27.3 linux/amd64`（157）、GOROOT 1.27.3（196）、`GOVERSION='go1.27.3'`（204）；1.27.2 反向字樣都是 0；`GOTOOLCHAIN='auto'`（201）只記錄（`10-checks.txt`）|
| QA-CR006-2／L2 | **PASS** | key＝`setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d42…a025`。這次三個 job 都 `Cache hit`（gosec 148、sca 149、unit 149），依判準只記錄 |
| QA-CR006-3／L3 | **PASS** | 三個 job 在兩種格式下都是 0 |
| QA-CR006-4／L4 | **PASS** | sast-gosec 273 `mod … v2.30.0 h1:67aF6LHa…`、274 `gosec h1 OK` |
| QA-CR006-5／L5 | **PASS** | 命令列 279 沒變；391 Files 57（掃描檔數）、392 Lines 7485（掃描行數）、393 Nosec 0（`#nosec` 抑制數）、394 Issues 0（high 以上發現數）；步驟 success。基準 41402171425 第 801～804 行也是 57／7485／0／0，沒有差異 |
| L6（10(a)） | **PASS** | sca：`No vulnerabilities found.` 1（221）、`This scan also found` 0、GO-2026-6809 0、GOVERSION go1.27.3（204）、Run govulncheck（217）步驟 success＝exit 0；dependency-review 247 沒有 Critical |
| QA-CR006-6 | **PASS** | check-runs（新 head）四個名稱都是 success（`04-check-runs-head.json`），與 ruleset 24629805 的 required contexts 相同（`05`）；沒有 continue-on-error（`11-gomod-workflow.txt`）|
| QA-CR006-7／L7 | **PASS** | `art-pr/go-test.json`：pass 181、fail 0、skip 0；基準 `art-base` 也是 181／0／0；unit-test 215 `go test -race`、231 `--- PASS: TestENG011…` |
| QA-CR006-8 | **PASS** | TestCT_* pass 14、fail 0（基準相同）|
| QA-CR006-9 | **PASS（PR head）** | 新 head 用 Go 1.27.3、`GOTOOLCHAIN=local`、`go test -count=1 -race -coverpkg=./... -coverprofile=c.out ./...`，exit 0，total 80.9%（精算 474/586＝80.887%）（`12-coverage-head.txt`、`12-c.out`）；基準 164b62c9 用同樣方法量也是 474/586（`13-coverage-base-164b62c9.txt`）|
| QA-CR006-10 | (a) **PASS**；(b) **PASS**；(c) **PASS**（限 head 474756381bd1…；合併後 run 待定） | 見 §4。(b)：品保 `33-qa-m5-6-project-verbose.txt` 第 19、32、33 行；(c)：品保 `32-qa-m5-5-gosec-binary.txt` 第 5、12、23、76、112 行 |
| QA-CR006-11 | **PASS** | govulncheck 三層：基準 41402171425 第 398～402 行（also found 1 in packages、0 in modules）＝0／1／0 → PR＝0／0／0；只差 package 層 1→0（ID 要由安全部 verbose 佐證），沒有其他差異；gosec 前後相同 |
| QA-CR006-12 | **PASS（PR head）** | go.mod 第 3 行 `go 1.27.3`，沒有 toolchain 行；workflow 62／91／134 是 `"1.27.3"`，`check-latest: false`；沒有 Dockerfile（`11`）|
| QA-CR006-14 | **PASS（依你指定的改用判準）／與原文「單一 commit」不符，請拍板** | 見 §3 |
| QA-CR006-13、15、合併後部分 | **待定** | 尚未合併 |

## 3. QA-CR006-14 與偏離說明
- PR 現在有 2 個 commit：`826d5d34`（實際變更）＋ merge `47475638`（parents 826d5d34、164b62c9；訊息 `Merge remote-tracking branch 'origin/main' …`）。
- merge-base(164b62c9, head)＝164b62c9；淨差異只有 2 個檔案：workflow 4/4、go.mod 1/1。head 上兩檔 sha256 為 `06833a31…f245`、`02842cf5…8ebf`，go.sum `1e2b43d7…` 沒變。
- merge commit 本身沒有夾帶 main 以外的內容：`git show --cc` 是空的（沒有衝突解決改動）；`diff(826d5d34, head)` 和 `diff(ca09137, 164b62c9)` 逐位元組相同，都是那 5 個文件檔（`02` 末段）。
- **偏離**：品保部 v0.4 QA-CR006-14 原文要求兩檔「在同一個 commit」，這點仍然成立（兩檔都在 826d5d34）。但 PR 已不是「單一 commit」，多了一個同步 main 的 merge commit。這屬於判準文字的偏離，**要不要接受，請使用者拍板**；建議下一版 CR／品保文件把判準改成「淨差異＋merge 不帶額外內容」。

## 4. QA-CR006-10(b)(c)：安全部 36 號 v0.1.1 與新 head M5 證據核對（2026-11-02 16:34～16:36（UTC+8））

（16:32 那一輪時證據尚未完成，快照保留在 `20-security-m5-status.txt`；本節取代當時的「待定」。）

E＝`/workspace/projects/short-url/docs/06-security/evidence/cr-006-m5/474756381bd170beb3a5d723b63f1b8ba48b2add/`

| 核對 | 結果 | 證據 |
|---|---|---|
| 36 號 v0.1.1 sha256 | `2936baec4ff6d20474a4b31a30e9820e8a3a1177de17544518dd20179d58d598`，與轉述相符 | `35-security-v011-verify.txt` 第 2 行 |
| E/SHA256SUMS sha256 | `43d33bb47bed52571dff61b56776787ff4e838a8634da79565aa533122a659bd`，相符；`-c` 14 列全 OK、exit 0；14 個檔案都有列入 | `35` |
| 目錄名 | 完整 40 碼新 head | `35` |
| head 未變 | 16:35:02 未登入 ls-remote＝未登入 API＝474756381bd1…2add；base 164b62c9；PR open | `30-head-refetch.txt` |
| M5-1 | PASS：(i) ls-remote＝(ii) API＝(iii) 維運值；rev-parse、parents 826d5d34＋164b62c9、merge-base 164b62c9；go.mod 02842cf5…、workflow 06833a31…、go.sum 1e2b43d7… | E/`01` 第 4、11、16～19、24～27、32～34 行 |
| M5-2 | PASS：官方 JSON 65efbc80…；tarball OK；go1.27.3；GOTOOLCHAIN auto | E/`02` 第 3、6、14、17～18 行 |
| M5-3 | PASS：兩個 install 都 exit 0 | E/`03` 第 52、61 行 |
| M5-4 | PASS：`gosec: go1.27.3`、h1 67aF6LHa… | E/`04` 第 3、5 行 |
| M5-5 | PASS：兩次 exit 0（第 12、113 行）、go1.27.3 標準庫（24）、symbol 層無發現（5、77）、GO-2026-6809 0 筆；module 層只有 6866／6860／5932（85、93、101），是預期項目 | E/`05` |
| M5-6 | PASS：go1.27.3 標準庫（19）、`No vulnerabilities found.`（32）、exit 0（33）、6809 0 筆 | E/`06` |
| M5-7 | PASS：02～06、05b 各 0 | E/`07` 第 3～9 行 |
| M5-8 | PASS：暫存已刪；SHA256SUMS 可驗 | E/`08` 第 3～4 行 |
| 35 號 §6.2 (a)(b)(c)（40 碼子目錄／只用公開 ls-remote 與未登入 API／「必須 0」＝6809 0 筆） | 皆符合 | 同上 |
| M7 | 沒有 6809 以外的 symbol 層或 package 層項目；gosec 二進位 module 層 3 筆屬已登錄的 SCR-054／055，不啟動 M7。專案 `./...` 三層皆 0 | E/`05`、`06`；品保 `32`、`33` |

品保獨立重現（全新 clone `/workspace/tmp-qa-pr31-47475638` @474756381bd1…；box 預裝 Go 1.27.3、GOTOOLCHAIN auto、全新 GOPATH／GOCACHE；自行 `go install` gosec v2.30.0、govulncheck v1.8.0，建置版本皆 go1.27.3；`31-qa-repro-setup.txt`）：
- M5-5：`32-qa-m5-5-gosec-binary.txt` 第 12、112 行兩次 exit 0；第 23 行 go1.27.3 標準庫；GO-2026-6809 0；只有 6866／6860／5932。
- M5-6：`33-qa-m5-6-project-verbose.txt` 第 19 行 go1.27.3 標準庫、第 32 行 `No vulnerabilities found.`、第 33 行 exit 0；GO-2026-6809 0。
- 與安全部輸出比對（去除路徑、時間、指令列）：兩者**相同**（`34-qa-vs-security-diff.txt`）。
- 沒有使用 `-scan module` 作判準。
- 限制：box 預裝 Go 的 tarball 這次沒有再對官方 JSON 核對；獨立核驗以安全部 M5-2 為準。

與安全部說法的差異／注意事項：
1. 36 號 v0.1.1 **沿用 `-v0.1.md` 檔名就地覆寫**，本機已沒有 v0.1 的檔案（v0.1 原文 `b7444593…` 只留在 repo，head 版經品保重算相符）。這不影響判定，但不符合「舊版另存」的慣例，建議安全部改用 `-v0.1.1.md` 檔名或另存備份。
2. E/`01` 沒有記錄 ls-remote 和 curl 的 exit code（舊 head 那一輪有記錄），屬證據格式上的小缺漏，值本身經品保重取相符。
3. SCR-060 仍屬實：未登入取新 run 的 job log 回 403（`35` 末行）。品保已用登入的唯讀權限取得 log，並確認與維運副本逐位元組相同（`00`）。

## 5. main 自身變更對數字的影響
沒有影響：測試數 181、CT 14、覆蓋率 474/586、gosec 57／7485／0／0、govulncheck 基準 0／1／0，在 ca09137 與 164b62c9 兩個基準之間都相同（main 只改了文件）。

## 6. 與維運說法的差異
沒有差異。另外記錄：這次三個 job 都命中快取（上一輪是 not found），依 L2 判準只記錄。
