---
文件：PR #15（CR-002 釘 Go 1.27.2 與 gosec 707f227）安全審查
版本：v0.1
狀態：核准（C1～C5 符合；附合併程序條件，見「結論」）
負責角色：安全部
最後更新：2026-10-19
對應需求：SEC-009、SEC-010、SEC-011；NFR-004、SI-05；CR-002（D-13）；D-08／D-09／D-10；安全部意見 [13](./13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md) C1～C5；門檻 [03](./03-ci-security-gates-v0.3.md) v0.3.7 §2；RR-010（[05](./05-residual-risk-register-v0.1.md)）
專案代號：SHORTURL
---

# PR #15 安全審查（v0.1）

## 安全審查摘要

| 項目 | 內容 |
|---|---|
| PR | [#15](https://github.com/444112029012/short-url/pull/15)　fix: CR-002 釘 Go 1.27.2 與 gosec 707f227 (ENG-020)（open、未合併） |
| 完整 head | `412cefafa0e6f285cc89c4ba8951de58ad8197ab`（單一 commit「fix(ci): pin Go 1.27.2 and gosec 707f227 per CR-002」，2026-10-19 09:55:53 UTC+8）；與總協調參考值一致 |
| base | main `2e54f7e2d04c891886a4e4210f5880475c9f0d18`；merge-base＝`2e54f7e`；2026-10-19 10:01:22 UTC+8 `git ls-remote` main 仍為 `2e54f7e` → 分支為最新（strict=true 成立） |
| CI run | [39187593303](https://github.com/444112029012/short-url/actions/runs/39187593303)（pull_request；2026-10-19 09:56:09 UTC+8；attempt 1；`head_sha`＝`412cefa…97ab`；四項 success） |
| 變更 | 2 檔、+21／−7：`.github/workflows/ci-devsecops.yml`（+20／−6）、`go.mod`（+1／−1）；**無產品程式、無文件、無 go.sum 變更**（go.sum 未變經驗證合理，見 C0） |
| 風險等級 | 高風險 PR 類別（CI 設定／供應鏈）；實際變更風險低 |
| 發現項 | 阻擋 0；低 1（SCR-019）、資訊 2（SCR-020、SCR-021） |
| **結論** | **核准（符合 C1～C5）**。合併程序條件：依 D-08／D-09 由審查部於 PR 留言「可合併」並附 run 39187593303 連結後由使用者合併；合併前若 main 有新 commit，依 D-10 更新分支並取得新四綠（新 run 仍須符合本檔 C1～C4 log 證據）。建議 #15 **優先於** #10／#12／#8／#11／#7 合併（main 目前 `sast-gosec` 紅燈，見 PR #14 內文 main push run 39183720869），其後各 PR 依既定順序更新分支取新四綠 |

- 審查方式：唯讀。工具：GitHub REST API（公開讀取：PR、run、jobs）、cursor-github connector（`get_pull_request`、`get_job_logs`）、本機唯讀 clone（`git fetch pull/15/head`）於 `/workspace/tmp/cr002/`。**沒有**在 PR 上留言或 review，沒有 push。
- 依 `secure-code-review` skill：前提（CI 全綠、PR 標 CR）確認；逐項檢查差異與周邊 workflow 全檔；複核掃描結果。
- SCR 編號接續（11 號檔最大 SCR-018），本檔使用 SCR-019～SCR-021。

## 維運聲稱逐項驗證

| # | 聲稱 | 驗證 | 證據 |
|---|---|---|---|
| 1 | Go 精確 1.27.2（go.mod＋三處 setup-go） | **屬實** | `go.mod:3` `go 1.27.2`；workflow `:53`、`:81`、`:127` `go-version: "1.27.2"`＋`check-latest: false` |
| 2 | gosec 以 `@v2.29.1-0.20261015133444-707f2271f665` 安裝 | **屬實** | workflow `:56` |
| 3 | h1 比對步驟輸出 OK | **屬實** | sast-gosec log 01:57:21Z（09:57:21 UTC+8）`mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`、`gosec h1 OK` |
| 4 | 門檻與檢查名稱不變 | **屬實** | 見 C4 |
| 5 | #14 已關閉 | **屬實** | #14 state closed、merged false、closed_at 2026-10-19 09:58:25 UTC+8 |
| 6 | PR 不含文件 | **屬實** | diff 只有 2 檔（workflow、go.mod） |

## 條件逐項（13 號檔 C1～C5）

### C0 變更範圍（前置）
- `git diff 2e54f7e...412cefa --stat`：`.github/workflows/ci-devsecops.yml | 26 ++++++++++++++++++++------`、`go.mod | 2 +-`，2 files changed, 21 insertions(+), 7 deletions(-)。與總協調參考值（2 檔、+21／−7）一致。
- workflow 變更僅：檔頭註解 `:6`、`:8`；三處 setup-go；gosec 安裝字串；新增 Verify 步驟 `:57-67`。其餘（gitleaks、govulncheck、dependency-review、unit-test 步驟）逐行未變。
- **go.sum 未變是否合理：合理。** (1) gosec／govulncheck 以 `go install module@version` 在模組外安裝，不寫入專案 go.sum；(2) `go` 指示 1.27.1→1.27.2 不改變依賴圖。本機於 PR head worktree 以 Go 1.27.2 執行 `go mod tidy` → **無任何差異**，`go mod verify` → `all modules verified`。

### C1 Go 精確 1.27.2 — **符合**
| 檢查 | 結果 |
|---|---|
| go.mod | `go.mod:3` `go 1.27.2`；無 `toolchain` 行（不需一致性比對） |
| setup-go | 三處皆 `go-version: "1.27.2"`、`check-latest: false`（`:51-54` sast-gosec、`:79-82` sca-govulncheck、`:125-129` unit-test）；全檔無 `1.27.x`／`stable`／`go-version-file` |
| 實際版本（log） | sast-gosec：`Acquiring 1.27.2 from …/go-1.27.2-linux-x64.tar.gz`、`Successfully set up Go version 1.27.2`、`go version go1.27.2 linux/amd64`（01:56:26Z）；sca-govulncheck：同上（01:56:28Z）；unit-test：同上（01:56:26Z）。secrets-gitleaks 不裝 Go |

### C2 gosec 安裝字串與 Go 環境 — **符合**
| 檢查 | 結果 |
|---|---|
| 安裝字串 | `:56` `go install github.com/securego/gosec/v2/cmd/gosec@v2.29.1-0.20261015133444-707f2271f665`，與 13 號檔／03 v0.3.7 完全相同；log `go: downloading github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665` |
| 禁止設定 | `git grep -iE 'GOFLAGS|GONOSUM|GOSUMDB|GOPROXY|GOPRIVATE|GOINSECURE'` 於 PR head：workflow 無任何命中（僅 evidence JSON 之 ruleset `exclude` 欄位，無關） |
| 實際 go env（log） | 三個 Go job 皆 `GOFLAGS=''`、`GONOSUMDB=''`、`GOPRIVATE=''`、`GOPROXY='https://proxy.golang.org,direct'`（Go 預設值，`direct` 僅為 proxy 失敗後備；非 `GOPROXY=direct`）、`GOSUMDB='sum.golang.org'` |

### C3 h1 比對步驟 — **符合**（附 SCR-019／SCR-020）
- 邏輯（`:57-67`）：`set -euo pipefail`；`bin="$(go env GOPATH)/bin/gosec"`；`go version -m "$bin"` 經 awk 精確比對欄位 `$1=="mod" && $2=="github.com/securego/gosec/v2"`；空行即 `::error::` ＋ `exit 1`；取第 4 欄與常數 `expected` 以 `[ != ]` 精確字串比對，不符 `exit 1`。比對對象是**實際安裝的 binary 內嵌 build info**，不是印常數。
- `:62` 的 `|| true`：只吞掉 `go version -m` 失敗的退出碼，下一行空字串檢查必定 `exit 1`，仍為 fail-closed（見下表 T4／T5；SCR-020）。
- CI 實際輸出（sast-gosec，01:57:21Z）：`mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`、`gosec h1 OK`；log `GOBIN='/home/runner/go/bin'`＝`GOPATH/bin`，所以 install 目的地即被驗證的檔案。
- 本機重跑該步驟腳本（自 PR head workflow `:59-67` 原文擷取，`/workspace/tmp/cr002/c3test/`）：

| # | 情境 | 結果 |
|---|---|---|
| T1 | 正確 binary（707f227）＋正確 expected | `gosec h1 OK`，exit 0 |
| T2 | 正確 binary＋**錯誤 expected**（`h1:AAAA…=`） | `::error::gosec h1 mismatch`，**exit 1** |
| T3 | 換成 gosec v2.29.0 binary | `::error::gosec h1 mismatch: h1:pF2HSLcn…`，**exit 1** |
| T4 | binary 不存在 | `::error::… line missing`，**exit 1** |
| T5 | 偽造 shell 腳本冒充 gosec（輸出假 mod 行） | `could not read Go build info`→`line missing`，**exit 1** |

- 後續步驟是否同一 binary：`:72` 以 PATH 呼叫 `gosec`。setup-go 將 `GOBIN`（`/home/runner/go/bin`）加入 PATH，ubuntu-latest 映像不預裝 gosec，本 run 實務上為同一 binary；但未以路徑綁定（SCR-019，低，不阻擋）。

### C4 門檻與檢查不變 — **符合**
| 檢查 | 結果 |
|---|---|
| gosec 參數 | `:72` `gosec -fmt=text -no-fail=false -severity=high ./...`（未變）；log `Including rules: default`、`Excluding rules: default` |
| govulncheck | `:86` `golang.org/x/vuln/cmd/govulncheck@v1.8.0`（未變） |
| gitleaks | `:39` `VER=8.21.2`、`:44` `--exit-code 1`（未變） |
| job 名稱 | `:30` `secrets-gitleaks`、`:47` `sast-gosec`、`:75` `sca-govulncheck`、`:121` `unit-test`（未變） |
| 新增排除／跳過 | 無 `-exclude`／`-exclude-dir`／`#nosec`／`continue-on-error`；新增步驟無 `if:`；既有 `if: github.event_name == 'pull_request'` 只在 dependency 兩步（未變） |
| dependency-review | `:93-118` 未變：compare API 非 200 即 `::error::`＋`exit 1`；`fail-on-severity: critical` |

### C5 run、base、標註、#14 — **符合**
| 檢查 | 結果 |
|---|---|
| run head | run 39187593303 `head_sha`＝`412cefafa0e6f285cc89c4ba8951de58ad8197ab`、event pull_request、`pull_requests`＝[15]、conclusion success |
| 四項 | sast-gosec success（job 117853180547）、unit-test success（117853180771）、sca-govulncheck success（117853180794，含 Probe／Dependency review 兩步 success）、secrets-gitleaks success（117853180896）；四個 job 的 head_sha 皆為 412cefa |
| base | `2e54f7e`＝現行 main HEAD（10:01:22 UTC+8 核對）；strict=true 下可合併 |
| CR 標註 | 標題含「CR-002」；內文「CR-002；決策 D-13；安全意見書 13」「取代 #14」 |
| #14 | closed、**未合併**（merged false），closed_at 09:58:25 UTC+8 |

## 掃描結果複核（run 39187593303）

| job | 證據 |
|---|---|
| sast-gosec | 有實際掃描：逐一 `Checking package`（config、ratelimit、application、health、domain、observability、httpapi、ports、sqlite、main）與 `Checking file`；Summary `Files : 17`、`Lines : 1206`、`Nosec : 0`、`Issues : 0`（與安全部本機 13 號檔 R2 完全一致）；無 `export data`／`could not import` 錯誤 |
| sca-govulncheck | `No vulnerabilities found.`（Go 1.27.2）；`dependency graph compare HTTP 200`；Dependency review：`did not detect any vulnerable packages with severity level "critical" or higher`、`did not detect any denied packages` |
| secrets-gitleaks | `8.21.2`；`38 commits scanned`、`no leaks found` |
| unit-test | `go vet` 通過；`go test` 8 個套件 `ok` |

## 發現項

### [SCR-019] gosec 執行步驟以 PATH 呼叫，未綁定到已驗證的 binary 路徑
- 嚴重度：**低**　不阻擋
- 位置：`.github/workflows/ci-devsecops.yml:61`（驗證 `$(go env GOPATH)/bin/gosec`）、`:72`（執行 `gosec`，經 PATH 解析）
- 證據：log `GOBIN='/home/runner/go/bin'`、`GOPATH='/home/runner/go'`，本 run 驗證檔即安裝檔；ubuntu-latest 不預裝 gosec，實務上 PATH 解析到同一檔。但若日後 runner 映像或前置步驟在 PATH 較前位置放入另一個 `gosec`，驗證與執行可能不是同一檔（TOCTOU 類）。
- 建議：下一次動到此 workflow 時（例如 C7 換回發行版的 CR）改為 `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high ./...`，或在驗證步驟加 `[ "$(command -v gosec)" = "$bin" ] || exit 1`。不需為此重開 PR #15。
- 狀態：開放（追蹤至 C7 換回 CR）

### [SCR-020] 驗證步驟 `:62` 的 `|| true` 與僅比對 h1 欄位
- 嚴重度：**資訊**
- 位置：`.github/workflows/ci-devsecops.yml:62`、`:65`
- 說明：`|| true` 會吞掉 `go version -m` 的失敗碼，但緊接的空字串檢查使其仍 fail-closed（本機 T4／T5 exit 1）。腳本只比對第 4 欄 h1、不比對第 3 欄版本字串；h1 已唯一綁定模組內容（sum.golang.org 記錄 `v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HN…`），安全上足夠。可選強化：同時比對 `$3` 等於 pseudo-version、移除 `|| true` 改為明確錯誤訊息。
- 狀態：資訊，不追蹤

### [SCR-021] 既有：actions 以 tag 釘選、gitleaks 下載無 checksum 驗證
- 嚴重度：**資訊**（既有、非本 PR 引入）
- 位置：`:33`、`:50`、`:78`、`:124` `actions/checkout@v4`；`:51`、`:79`、`:125` `actions/setup-go@v5`；`:116` `actions/dependency-review-action@v4`；`:40-41` gitleaks tarball 以 curl 下載後直接解壓
- 說明：03 §7.4／§10 要求「穩定後」釘完整 commit SHA；workflow `:5` 亦自註。gitleaks 未比對官方 checksums。本 PR 依最小變更原則未處理，正確。
- 建議：另開 CI 強化 PR（標準變更）：actions 釘 SHA＋Dependabot `github-actions`；gitleaks 比對 release checksums。G5 前完成。
- 狀態：開放（追蹤至 G5 前）

## 殘餘風險
- RR-010（低）於 #15 合併後開始適用（05 v0.1.9）；本審查未改變其等級。
- 13 號檔 C6（合併後 7 日內 CI 金絲雀）、C7（每週檢查上游發行；最遲 2026-12-31 重評）、C8 於 #15 合併後起算。

## 無法查證
- 無。PR／run／jobs／logs 皆可讀；main HEAD 以 `git ls-remote` 核對。審查部「可合併」留言屬合併前程序，本檔撰寫時尚未發生。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-19 | 安全部 | 初版：PR #15（head `412cefa`、run 39187593303）審查，核准；C1～C5 符合；go.sum 未變經 `go mod tidy` 驗證合理；h1 驗證步驟本機 T1～T5 驗證 fail-closed；SCR-019（低）、SCR-020（資訊）、SCR-021（資訊、既有）；建議 #15 優先合併 |
