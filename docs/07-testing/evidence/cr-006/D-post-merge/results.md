# CR-006 合併後品保驗收取證（PR #31 → main）

- 執行：品保部，2026-11-03 09:21～09:25（UTC+8）。GitHub 全程唯讀。
- 對象：merge commit `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`、push run 41559091290（2026-11-03 09:20:24（UTC+8）建立，09:23 前後完成）。
- log：`api-job-<id>-<name>.log`（per-job API 原始格式）。grep 都沒有加 `^`。

## 0. 說法核對
| 說法 | 實測 | 證據 |
|---|---|---|
| 09:20 合併 | merged_at `2026-11-03T01:20:21Z`＝**2026-11-03 09:20:21（UTC+8）** | `02-pr31.json`、`01-unauth.txt` |
| merge commit 1b53fcde | `merge_commit_sha`＝1b53fcde7cf5068dcf7193b92205d6e6e9adb09d；未登入 ls-remote main＝同值 | `01` 第 2、4 行 |
| 合併的 head 474756381bd1… | `.head.sha`＝474756381bd170beb3a5d723b63f1b8ba48b2add | `01`、`03-pr-view.json` |
| push run 41559091290 | event push、headSha＝merge sha、completed／success | `20-run-41559091290.json` |

## 1. 逐項判定
| 項目 | 判定 | 證據（檔：行） |
|---|---|---|
| QA-CR006-15 合併紀錄 | **記錄完成；rule suite result＝bypass，依判準由審查部／使用者判定是否接受** | 見 §2 |
| QA-CR006-13 ruleset 未變 | **PASS** | 快照 `sha256sum -c` 9 列 OK（`10-snapshot-sha256c.txt`）；三份 `jq -S` 檔與快照 `diff` 為空、`cmp` 位元組相同（含 updated_at、_links、current_user_can_bypass）（`11-ruleset-diff.txt`）|
| QA-CR006-1／L1 | **PASS** | gosec 115／133／172／180、sca 116／134／173／181、unit 116／134／173／181：spec 1.27.3、go1.27.3、GOROOT 1.27.3、GOVERSION go1.27.3；1.27.2 字樣 0；GOTOOLCHAIN auto（177／178／178）只記錄（`23-checks.txt`）|
| QA-CR006-2／L2 | **PASS** | key `…go-1.27.3-cf6b0d42…a025`：sca 214 `Cache saved with the key`；gosec 381、unit 253 `Failed to save … another job …`；還原時都是 not found（132／133／133），只記錄 |
| QA-CR006-3／L3 | **PASS** | 三個 job 都是 0 |
| QA-CR006-4／L4 | **PASS** | gosec 250 h1 行、251 `gosec h1 OK` |
| QA-CR006-5／L5 | **PASS** | 命令列 256 沒變；368～371 Files 57／Lines 7485／Nosec 0／Issues 0；與基準（164b62c9 run 41402171425 及 PR run）相同 |
| QA-CR006-6 | **PASS** | merge sha 的 check-runs：四個必要檢查都是 completed／success（另有非必要的 `update-go_modules-graph` success）（`21-check-runs-merge.json`）；名稱與 ruleset 相同（`22-required-contexts.txt`）|
| QA-CR006-7／L7 | **PASS** | `art/go-test.json`：pass 181、fail 0、skip 0；unit 202 `go test -race`、218 ENG-011 PASS |
| QA-CR006-8 | **PASS** | TestCT_* pass 14、fail 0 |
| QA-CR006-9 | **PASS** | 全新 clone `/workspace/tmp-qa-post-merge` @1b53fcde；Go 1.27.3、`GOTOOLCHAIN=local`、`-count=1 -race -coverpkg=./...`，exit 0，total 80.9%（精算 474/586＝80.887%，與 PR 前後基準相同）（`31-coverage-merge.txt`、`31-c.out`）|
| QA-CR006-10(a)／L6 | **PASS** | sca 203 `No vulnerabilities found.`（計數 1）、`This scan also found` 0、GO-ID 0、GOVERSION go1.27.3（181）、Run govulncheck（199）步驟 success。dependency-review 在 push 事件下是 skipped（步驟名稱「(PR only)」，設計如此），PR run 已確認沒有 Critical |
| QA-CR006-10(b)(c) | **PASS（合併前，head 474756381bd1…）** | 見 `qa-evidence-pr31-47475638/results.md` §4；merge 與該 head 的 tree 差異為空（`30` 第 8～9 行），結果可以沿用 |
| QA-CR006-11 | **PASS** | 基準 0／1／0（`../qa-evidence-pr31-47475638/base-run-41402171425.log` 第 398～402 行）→ 合併後 0／0／0；沒有其他差異 |
| QA-CR006-12 | **PASS** | main 上 go.mod 第 3 行 `go 1.27.3`，沒有 toolchain 行；workflow 62／91／134 `"1.27.3"`、`check-latest: false`；沒有 Dockerfile（`30-diff-gomod.txt`）|
| QA-CR006-14 | **PASS（依淨差異判準；「單一 commit」偏離仍待拍板）** | 舊 main 164b62c9 → merge：只有 workflow 4/4、go.mod 1/1；merge^1..merge 相同；`diff(474756381bd1, merge)` 空；sha256 06833a31…f245／02842cf5…8ebf，go.sum 1e2b43d7… 沒變。從 ca09137 起算的淨差異是這 2 檔加上 PR #32 的 5 個 docs 檔（`30`）|

## 2. QA-CR006-15 詳細
- merged_at：2026-11-03 09:20:21（UTC+8）；merge_commit_sha：`1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`；merged_by：`444112029012`（YSW）；merged head：`474756381bd170beb3a5d723b63f1b8ba48b2add`。
- 方式：**merge commit**（PR API 沒有直接欄位；依 commit 判斷：兩個 parent、committer GitHub、訊息「Merge pull request #31 from …」，`04-merge-commit.txt`）。
- parents：`164b62c956a9…3661`（舊 main）＋`474756381bd1…2add`，符合。
- rule suite：id **4869407317**，pushed_at 2026-11-03 09:20:21（UTC+8），actor `444112029012`，before 164b62c9 → after 1b53fcde，ref refs/heads/main，**result `bypass`**，evaluation_result `null`（`05-rule-suites.json`）。
- 細節（`06-rule-suite-4869407317.json`）：24629805 main-required-checks：required_status_checks **pass**、non_fast_forward pass、deletion pass；24629854 main-pr-review：**pull_request fail**（"New changes require approval from someone other than the last pusher."），non_fast_forward pass、deletion pass。
- 判讀：**required_status_checks＝pass，符合硬條件**；pull_request 規則未通過，是靠 24629854 的 bypass_actors（RepositoryRole 5、`bypass_mode: pull_request`）bypass 合併的，即 CR §6 R10（RR-011／D-17 已接受）所述的情況。依判準，品保只記錄；是否接受由審查部／使用者判定。

## 3. 差異與注意事項
- 與說法沒有不符。
- 合併以 bypass 完成（見 §2）。
- L2：gosec、unit-test 快取存檔失敗（並行搶存），只記錄。
- 前一輪的偏離：QA-CR006-14「單一 commit」的偏離仍待使用者拍板。
