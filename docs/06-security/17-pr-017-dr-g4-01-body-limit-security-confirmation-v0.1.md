---
文件：PR #17（DR-G4-01 建立請求本文上限 32768 B＋MaxHeaderBytes）合併前安全確認
版本：v0.1.2
狀態：**v0.1.2：#17 已合併**（#17 於 2026-10-19 14:00:38（UTC+8）合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`（parents `ed62039`＋`e6d9ab1`，合併的 head 即安全部核准的 `e6d9ab1`）；SCR-022／023／040／041 已關閉）。**v0.1.1：通過，僅適用於最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`**（附 §7 合併條件；SCR-022／023 已修正、SCR-040／041 已處理，皆待 #17 合併後關閉）。歷史：v0.1 對 head `3b1aa5b` 通過
負責角色：安全部
作者：安全部
日期：2026-10-19
最後更新：2026-10-19
對應：安全部意見 [10](./10-dr-g4-01-body-limit-security-opinion-v0.1.md) §4（C1～C5）；設計部 `docs/02-architecture/06-method-specs-v0.2.md` v0.2.1 §13.4.1（DR-G4-01）；SEC-001、SEC-007、NFR-003；門檻 [03](./03-ci-security-gates-v0.3.md)；前置 [16](./16-pr-007-final-head-security-confirmation-v0.1.md)（PR #7）
專案代號：SHORTURL
發現項編號區段：本文件新發現使用 **SCR-040～SCR-049**（SCR-028～SCR-039 保留給同時進行的 PR #18 審查，避免衝突）
---

# PR #17 合併前安全確認（v0.1.1）

## 0. 兩段核對（v0.1.1，2026-10-19 13:28（UTC+8）查證）

背景：v0.1 通過後，main 因 #16 合併前進到 `ed62039`。研發部在 `3b1aa5b` 之上先加 `347a027`（規格 R2／R3），再以 merge 同步 main 得到 **H**＝`5b5017b2b9109b9e547e9273787fd6c790ae92d5`；設計部接著以唯一 parent 為 H 的 commit 加入簽核紀錄，得到**最終 head** `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`。本節分兩段核對：第一段 `3b1aa5b → H`，第二段 `H → 最終 head`。

方法：GitHub 讀取 API（PR、run／jobs、job log）、`git fetch` 至 `/workspace/tmp/pr18-review/short-url`；**經使用者授權**在本機暫時 worktree（最終 head）執行 `go test -count=1 ./...`，執行後 `git status` 乾淨、worktree 已移除。未在 GitHub 留言／review／核准／合併／push，未寫 exploit。

### 0.1 第一段：`3b1aa5b → H`

| # | 宣稱 | 查證結果 |
|---|------|----------|
| 1 | main＝`ed620394b1c9986cc4ac243eb1543c8b677d92cb`（#16 已合併） | ✅ `git ls-remote refs/heads/main`＝`ed62039…`；該 commit 為 `Merge pull request #16`（2026-10-19 13:17:22 UTC+8），parents `24b8da0`＋`95e7561` |
| 2 | H 的 parents 為 `347a027`、`ed62039` | ✅ `5b5017b2b9109b9e547e9273787fd6c790ae92d5` 的 parents＝`347a02759552d8f8a1e18484eaa3b70fb5fe55e1`、`ed620394b1c9986cc4ac243eb1543c8b677d92cb`（13:18:28 UTC+8） |
| 3 | `git show --remerge-diff H` 為空 | ✅ 自行執行，輸出 0 byte：同步 merge 沒有解衝突或夾帶改動 |
| 4 | `3b1aa5b..347a027` 只改規格檔（R2、R3） | ✅ `347a027` 唯一 parent＝`3b1aa5b`；只改 `docs/02-architecture/06-method-specs-v0.2.md`（+8／−8）。內容見 §0.3 |
| 5 | 4 個程式 blob 與 `3b1aa5b` 相同 | ✅ `cmd/shorturl/main.go` `e27a120a2227e35b34d2eb1b0f0c0f1bf50c5306`、`cmd/shorturl/main_test.go` `cb9987afdaa2593dba0d65d61fd1f7b77071a199`、`internal/httpapi/adapter.go` `216e43d71680a4b6a7ddc78bb9bbd2149c3fc325`、`internal/httpapi/bodylimit_test.go` `e1dec979984273a486992e896a51f42565bea29f`，在 `3b1aa5b`、H、最終 head 三者完全相同 |
| 5 | H 與 main 只差 #17 原本 5 檔 | ✅ `git diff ed62039 H --stat`：`main.go`（28）、`main_test.go`（+76）、`06-method-specs-v0.2.md`（78）、`adapter.go`（5）、`bodylimit_test.go`（+147），共 5 檔、+292／−42。另外 `git diff 3b1aa5b H`（排除 docs）只多 `internal/integration/security_headers_test.go`，正好是 main `24b8da0..ed62039` 的唯一增量（#16） |
| 5 | 規格檔 sha256 `c89e336e…a687` | ✅ H 上 `06-method-specs-v0.2.md` sha256＝`c89e336e84dc524f74961dab9874bde1d8669cec7373e798f5edef5a35aea687`（與 `347a027` 相同；blob `110b9b384bc0c57849356799bb9a5ffc4f2f54a4`） |
| 7 | run 39203891218 四綠 | ✅ event `pull_request`、run `head_sha`＝H、`conclusion: success`，建立 13:18:50 UTC+8。四個 job 的 head_sha 都是 H、都 success：sast-gosec（117902436648）、secrets-gitleaks（117902436865）、unit-test（117902436927）、sca-govulncheck（117902437003）。gosec log：`gosec h1 OK`、`Files : 19`、`Lines : 1719`、`Nosec : 0`、`Issues : 0` |

### 0.2 第二段：`H → 最終 head`

| # | 宣稱 | 查證結果 |
|---|------|----------|
| 1 | PR head＝`e6d9ab12d62d7f872a7bc18e8761a2a43ca45402` | ✅ PR API `head.sha` 與 `git ls-remote refs/pull/17/head` 相同；PR `state: open`、`mergeable: true`、`mergeable_state: blocked`（缺核准 review）、10 commits、6 檔、+331／−42 |
| 2 | 唯一 parent 是 H | ✅ `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402` parents＝`5b5017b2b9109b9e547e9273787fd6c790ae92d5`（13:22:29 UTC+8） |
| 6 | 只新增 1 個簽核檔 | ✅ `git diff --name-status H e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`：只有 `A docs/02-architecture/09-signoff-DR-G4-01-method-specs-v0.2.1.md`；sha256＝`806170e05bef8639c7437ebe499d071eea4c8a649f344f0a2145c1ec0462e3bd`（blob `827e74ffa7433df7d7d46656412e1c93f5363b01`），與宣稱 `806170e0…e3bd` 相符 |
| 6 | 該檔為簽核紀錄，不含程式或設定 | ✅ 全檔為 Markdown 表格與條列（綁定 H、parents、規格 sha256、run 39203891218；必改 1～5 與 R2、R3 關閉狀態；聲明以 PR 版為準、只在「唯一 parent 為 H 且只新增本檔」時有效）；沒有程式碼區塊、沒有 workflow／設定內容 |
| 7 | run 39204252238 四綠且對應最終 head | ✅ event `pull_request`、run `head_sha`＝`e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`、`conclusion: success`，建立 13:22:52 UTC+8。四個 job 的 head_sha 都是 `e6d9ab1…`、都 success：unit-test（117903521806）、sca-govulncheck（117903521903）、secrets-gitleaks（117903521934）、sast-gosec（117903521954）。log：`gosec h1 OK`、`Files : 19`、`Lines : 1719`、`Nosec : 0`、`Issues : 0`；gitleaks `56 commits scanned.`、`no leaks found`；govulncheck `No vulnerabilities found.`、`dependency graph compare HTTP 200`、Dependency review 無 critical／denied；unit-test `go test ./...` 全部套件 `ok` |
| — | merge ref | `refs/pull/17/merge`＝`134f135`，parents＝`ed62039`＋`e6d9ab1`，tree＝最終 head tree＝`b14a40eb`（main 未再前進，合併不會帶進額外內容） |
| 8 | 本機 `go test ./...` | ✅ 最終 head、Go 1.27.2 linux/amd64、`-count=1`：10 個有測試的套件全部 `ok`（`cmd/shorturl 0.506s`、`internal/httpapi 0.066s`、`internal/integration 0.088s`…），`internal/ports` 無測試檔 |

### 0.3 R2、R3 規格差異的安全判讀（第 4 點）

- **R2（SCR-040 修正）**：§13.2 兩列改寫。
  - `RateLimitGuard`：刪除「現 stub 一律放行／行為延後 ENG-011」，改為現行實作：超限回 `rate_limited`→429、**nil Guard 也回 rate_limited（fail-closed）**、建立 30／導向 120 次／分鐘、`NewGuard` 非正值回落預設、放行才扣額度、單程序；`sourceID` 由 `ClientSource` 產生（只有可信代理才採 XFF／X-Real-IP，IPv4-mapped 折疊、IPv6 取 /64）；不記錄來源位址；先限流再解碼。逐項對照 main `ed62039` 的 `internal/ratelimit/guard.go`（`NewGuard` 的 `<= 0` 回落 30／120、`CheckCreate`／`CheckRedirect`／`allow` 在 nil 時回 `RateLimited`）與 `adapter.go`（`CheckCreate`／`CheckRedirect` 在解碼／導向前），**描述與程式一致**。
  - `ObservabilityHooks`：刪除「no-op／延後 ENG-012」，改為六個方法、`atomic.Uint64` 計數、`slog` JSON、`emit` 以 `recover` 吞錯、只記 `event`＋`short_code`／`error_type`，不記 Authorization／Cookie／本文／長網址／客戶端 IP。與 `internal/observability/hooks.go` 一致。
- **R3**：§13.3 欄名「檔案:行」改為「函式」，各列行號改成函式／分支名稱；§13.5 `service.go:25` 改成欄位與建構子參數名。只是引用方式改變，錯誤對映（404／500／400／429 與對外 message）**沒有變動**。H 的規格檔中已沒有 `*.go:行號` 形式的引用，也沒有 stub／no-op 字樣。
- 修訂紀錄 v0.2.1 列補記 G4-PR17-R2／R3。
- **結論**：R2 是 SCR-040 的修正，R3 是同類的行號清理；兩者都只把文件對齊已審過的程式，**沒有放寬任何安全規格**（上限、fail-closed、來源識別、日誌欄位限制都和程式一致或更明確）。

### 0.4 兩段核對結論

最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402` 相對 v0.1 已確認的 `3b1aa5b`：程式與測試 blob 完全相同；新增內容只有 (a) 規格 R2／R3（文件）、(b) 同步 main 帶進的 #16 測試檔（remerge-diff 為空）、(c) 設計部簽核紀錄（文件）。C1～C5 結論不變，四項必要檢查在最終 head 上全綠。


## 1. 結論

**v0.1.1 判定：通過，僅適用於最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`。** 核准只適用於這個 SHA；之後若有任何 commit（包括同步 main 的 merge commit），本核准即不適用，必須重新送安全部書面確認（見 §7）。

**v0.1 判定（歷史）：通過。** 新 head `3b1aa5b68ea2251a6551f8901e657c95b3657410` 已修正舊 head `7155a48` 的兩項必改（SCR-022、SCR-023），並處理資訊項 SCR-024、SCR-025；C1～C5 全部達成。新發現 SCR-040～SCR-042 都是「資訊」等級，不阻擋合併。合併前必須遵守 §7 的合併條件。

- 審查方式：唯讀。`git clone` 公開 repo 到 `/workspace/tmp/pr17-review/`，加上 GitHub 讀取 API 與公開 rules API。**沒有**執行 PR 內的程式碼，**沒有**在 GitHub 留言、review、approve、合併或 push，也沒有修改其他部門或本目錄既有的檔案。
- C4(a) 本文長度 24366 B 是安全部用 Python 在本機獨立組出相同字串算出來的，沒有執行 PR 程式碼。

## 2. 查證過的 SHA、run 與 PR 狀態（2026-10-19 12:34（UTC+8）核對）

| 項目 | 值 | 核對方式 | 結果 |
|---|---|---|---|
| PR #17 head | `3b1aa5b68ea2251a6551f8901e657c95b3657410` | `get_pull_request`；`git ls-remote refs/pull/17/head` | 與宣稱相符 |
| PR #17 base | `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`（#7 merge commit） | `get_pull_request`；`git merge-base 24b8da0 3b1aa5b` = `24b8da0` | 相符；非堆疊 PR |
| main | `24b8da0`（尚未前進） | `git ls-remote refs/heads/main`；`get_branch main` | main = base |
| PR #18 | open、未合併（head `9370def`） | `get_pull_request #18` | main 沒有因 #18 前進 |
| mergeable | `true`；`mergeable_state: blocked`；`rebaseable: true` | `get_pull_request` | blocked 的原因與「尚缺 1 個核准 review」一致（reviews 目前為空） |
| merge ref | `6cba901a414682784220d6bb97f1c7e0c87d4d61`（parents `24b8da0`＋`3b1aa5b`）；tree `d3437ab…` 與 head tree 相同 | `git rev-parse` | CI 測的內容就是 head 的內容 |
| 舊 head → 新 head | `7155a48..3b1aa5b`：4 個 commit（`afafe09`、`8910c1d`、`95e7871`、`3b1aa5b`），一般 fast-forward | `git log` | 沒有改寫歷史 |
| Workflow run | `39199412428`（ci-devsecops #48，event `pull_request`，attempt 1），`head_sha` = `3b1aa5b…`，conclusion success；12:25:42～12:27:10（UTC+8） | `get_workflow_run` | 相符 |
| secrets-gitleaks | success，head_sha 相符；gitleaks 8.21.2，`fetch-depth: 0`，`52 commits scanned`、`no leaks found` | `list_check_runs_for_ref`＋job log | 通過 |
| sast-gosec | success，head_sha 相符；`gosec h1 OK`（`h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`，`v2.29.1-0.20261015133444-707f2271f665`）；Files 19、Lines 1719、Nosec 0、**Issues 0**；Go 1.27.2 | 同上 | 通過 |
| sca-govulncheck | success，head_sha 相符；govulncheck v1.8.0 `No vulnerabilities found.`；dependency review 沒有 critical 以上 | 同上 | 通過 |
| unit-test | success，head_sha 相符；`go vet ./...`、`go test ./...` 全部 `ok`（含 `cmd/shorturl`、`internal/httpapi`） | 同上 | 通過（CI 沒有加 `-race`、`-v`，研發部宣稱的 -race 是本機結果，CI 沒有佐證） |
| `.github`、`go.mod`、`go.sum` | `24b8da0..3b1aa5b` 與 `7155a48..3b1aa5b` 都沒有變動 | `git diff --stat` | 沒有變動 |
| 變動檔案（對 base） | `cmd/shorturl/main.go`、`cmd/shorturl/main_test.go`、`internal/httpapi/adapter.go`、`internal/httpapi/bodylimit_test.go`、`docs/02-architecture/06-method-specs-v0.2.md`（+287／−37） | `git diff --stat` | 範圍符合 |
| main ruleset | required checks 4 項（gitleaks／gosec／govulncheck／unit-test），`strict_required_status_checks_policy: true`；PR 需 1 個核准 review；`dismiss_stale_reviews_on_push: false`、`require_last_push_approval: false` | 公開 `GET /repos/…/rules/branches/main` | strict=true 確認（另見 SCR-042） |

## 3. 舊 head `7155a48` 初步發現摘要

- C1 已達成（base 含 #7 merge `24b8da0`，guard 是真正的實作）；C2、C3、C5 通過；C4 (b)(c)(d)(e) 通過，(a) 缺測試。
- 實作：`io.LimitReader(body, 32769)` 讀完後判斷 `len > 32768`；限流在讀本文之前；`newHTTPServer` 用在 `run()`；四個 Timeout（5／15／15／60 秒）不變。
- 發現：SCR-022（低，必改，缺 C4(a) 測試）、SCR-023（低，必改，429 測試的計數器無效）、SCR-024（資訊，容許值太寬）、SCR-025（資訊，MaxHeaderBytes 額外寬放 4096）、SCR-026（資訊，設計部範圍：repo 方法規格缺 §13.4.1）、SCR-027（資訊，431 由 net/http 回應、不帶安全標頭）。

## 4. `7155a48..3b1aa5b` 逐行審查

### 4.1 非測試程式

- **只有** `cmd/shorturl/main.go` 的 `newHTTPServer` 前面多了 3 行 `//` 註解（說明 4096 寬放、實際約 20480 B、超過時由 net/http 在 router 前回 431 text/plain）。沒有可執行程式變動。
- `internal/httpapi/adapter.go` 在這個區間沒有變動，對 base 的差異仍只有 `maxCreateBody` 8192→32768 加上註解。`decodeCreate` 仍是 `io.ReadAll(io.LimitReader(body, maxCreateBody+1))`，接著判斷 `len(raw) > maxCreateBody` 回 400，之後才 `json.NewDecoder`。`HandleCreate` 仍是先 `CheckCreate` 再 `decodeCreate`。
- **沒有**改用 `MaxBytesReader`（宣稱屬實，C2 接受）。

### 4.2 SCR-022／C4(a)：`TestCreateC4aSurrogatePairEscaped2048CharsAccepted`

| 查核點 | 結果 |
|---|---|
| 2028＋20＝2048 rune | `n := 2048 - utf8.RuneCountInString("https://example.com/")` = 2048−20 = 2028；fixture 另外斷言 `RuneCountInString(want) == 2048`。正確 |
| 本文位元組數 | `{"url":"`（8）＋前綴（20）＋2028×12（24336）＋`"}`（2）= **24366 B**。大於 8192，也不超過 32768（測試本身也斷言）；距理論最壞值 2048×12＋外框 ≈ 24.6 KB 約 0.2 KB（前綴是 ASCII） |
| 跳脫形式與欄位 | `\uD83D\uDE00`（U+1F600，BMP 外，surrogate pair 跳脫）；欄位是實作用的 `url`。符合 SCR-022 定義 |
| 讀回來源 | **從儲存讀**：用回應裡的 `short_code` 呼叫同一個 SQLite store（`openStore`＋`routerWith` 注入同一個 `*sqlite.URLStore`）的 `FindByShortCode`，斷言 `RuneCountInString(rec.LongURL) == 2048`，而且 `rec.LongURL == want`（逐字元相同） |
| 回應斷言 | `res.Code == 201` |
| 錯誤實作下是否仍會通過 | 上限仍是 8192 → 400，失敗；validator 以 byte 計（8132 B > 2048）→ 400，失敗；以 UTF-16 code unit 計（4076）→ 400，失敗；surrogate pair 被解成兩個 U+FFFD → rune 數 4076 ≠ 2048，等值比對也失敗；儲存截斷或改寫 → 等值比對失敗；回應沒有 short_code → 查無資料，失敗。**沒有找到會讓錯誤實作通過的情況**。沒有斷言 long_url 回應欄位，但讀回儲存是更強的斷言 |

結論：**SCR-022 已修正。**

### 4.3 SCR-023：`TestRateLimitedOversizedCreateGets429WithoutReadingBody`

- `countingReader` 加了 `r io.Reader` 欄位；`Read` 委派給 `r` 並累加實際讀到的位元組數。
- 測試改成 `body := &countingReader{r: bytes.NewReader(paddedCreateBody(t, createBodyLimit+1))}`，再把 `body` **整個**交給 `postRawBody`。計數器包住整個 32769 B 本文，不再接在後面。
- `httptest.NewRequest` 碰到非 `*bytes.Reader` 型別時用 `io.NopCloser` 包裝，不會預讀；`ResponseRecorder` 也不會排空本文。所以 `n == 0` 確實證明 adapter 和 router 中介層都沒有讀本文。
- 限流條件：`newLimitedServer(t, 1, 120, "")`，同一個 RemoteAddr 的第一次建立回 201，第二次斷言 429 `rate_limited`，而且 `n == 0`。
- 結論：**SCR-023 已修正。**（研發部宣稱「限流分支讀 1 byte 時測試會失敗」，與上述機制一致。）

### 4.4 SCR-024：`TestCreateBodyReadIsBounded`

- 容許值改成 `got > createBodyLimit+1` 才失敗，也就是不超過 32769。無窮本文模式（`r == nil`）下，`io.LimitReader` 會先把 `p` 截到剩餘 N 再呼叫下層 `Read`，所以實際讀到的量剛好不超過 32769，不會因為緩衝分塊而誤判失敗。結論：**已處理。**

### 4.5 SCR-025：`TestNewHTTPServerRejectsOversizedHeaders`

- 這個測試用 `newHTTPServer` 的實際設定（`ts.Config = cfg`），handler 是**測試內的 `http.HandlerFunc`，一律回 204**，不是應用程式 router。所以 204 代表「標頭通過 net/http 的上限檢查並到達 handler」，用途是驗證伺服器設定，這個設計合理。
- 16000 B：`X-Pad` 值 16000 B，加上請求列與 `Host`、`Connection` 等約 60 B，總計約 16060 B，低於 16384（也低於 20480）→ 204。合理。
- 21 KiB = 21504 B，大於 16384＋4096 = 20480 B（超出約 1 KiB）→ 431。32 KiB → 431。兩者都斷言 `Content-Type` 以 `text/plain` 開頭，與 net/http 回 431 的行為一致。
- `newHTTPServer` 已有註解寫明 4096 寬放與實際約 20 KiB 上限；方法規格 §13.4「標頭上限」列也寫明了。結論：**已處理。**（沒有測 16384～20480 之間的邊界，不要求。）

### 4.6 文件（方法規格 v0.2.1 §13.4.1）

- `docs/02-architecture/06-method-specs-v0.2.md` 已在 PR 內升為 v0.2.1，含 §13.4.1，以及 §13.4 的「建立請求本文」列和新增的「標頭上限」列。**SCR-026 已處理。**
- 與本機設計部版本 `diff`：**不完全一致**（本機 sha256 `4190e41a…`、PR 版 `6769f80d…`）。差異包括：狀態與日期改為 2026-10-19 生效；§13.2／§13.4 的行號引用改為函式名；外框由 `{"long_url":""}` 改為實際欄位 `{"url":""}`（PR 版較正確）；§13.4.1 補上實作方式（`io.LimitReader`、未改用 `MaxBytesReader`）、C4 (d)(e) 測試、`MaxHeaderBytes` 與生效敘述；修訂紀錄作者寫「設計部（研發部草擬）」。**裁定本身的實質內容一致**：上限 32768 B、超限 400 `invalid_url`、解碼前截斷、2048 rune 規則不變、HTTP 契約不變。與安全條件 C1～C5 一致（見 SCR-041）。
- 漏改的 8192：PR 版中剩下的 8192 只出現在「舊上限 8192」與修訂紀錄「8192→32768」，屬歷史描述，**沒有漏改**。程式碼中的 `defaultMaxSources = 8192`（guard 來源表上限）與本文上限無關；測試中的 8192 是 C4(a) 的「舊上限」斷言。`08-ruling-G4-PR2-R1-go-mapping-v0.1.md`（設計部歷史裁示）與本部 08 號檔的 8192 屬歷史紀錄，PR 說明已列為殘餘。
- 規格內容正確性屬設計部範圍，安全部只確認與安全條件一致。

## 5. C1～C5 結論

| # | 條件 | 結論 | 依據 |
|---|---|---|---|
| C1 | 實作 PR 的 base 含 #7，guard 不是 stub | **已達成** | base = main = `24b8da0`（#7 merge commit）；merge-base 相符 |
| C2 | 解碼前截斷，超限回 400 `invalid_url`／`URL validation failed`；`LimitReader(max+1)` 可接受 | **已達成** | `decodeCreate` 沒有變動；`TestCreateBodyReadIsBounded` 證明最多讀 32769 B；`TestCreateBody32769RejectedInvalidURL` |
| C3 | `http.Server` 明確設定 `MaxHeaderBytes`（建議 ≤ 16 KiB），超大標頭回 431 | **已達成（#17 合併後生效）** | `maxHeaderBytes = 16 << 10`；`run()` 使用 `newHTTPServer`；四個 Timeout 不變（測試斷言）；21／32 KiB → 431 |
| C4 | (a)(b)(c)(d)(e) 測試，使用實際欄位 `url` | **已達成** | (a) `TestCreateC4aSurrogatePairEscaped2048CharsAccepted`；(b) `TestCreateBody32769RejectedInvalidURL`；(c) `TestE2E_02_Length2048Boundary`（2049 → invalid_url，而且不寫入）；(d) `TestCreateBodyExactly32768Accepted`；(e) `TestRateLimitedOversizedCreateGets429WithoutReadingBody`（n==0）；unit-test 為綠 |
| C5 | 2048 rune、Content-Type、未知欄位、多餘 JSON 規則不得一起放寬 | **已達成** | `validator.go`、`jsonContentType`、`DisallowUnknownFields`、多餘值檢查都沒有變動；adapter 對 base 只改常數 |

## 6. 發現項

### 6.1 前次發現（SCR-022～SCR-027）狀態

| ID | 嚴重度 | 狀態 | 說明 |
|---|---|---|---|
| SCR-022 | 低（必改） | **已關閉**（v0.1.2：merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`） | §4.2 |
| SCR-023 | 低（必改） | **已關閉**（v0.1.2：merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`） | §4.3 |
| SCR-024 | 資訊 | **已處理** | §4.4，容許值 ≤ 32769 |
| SCR-025 | 資訊 | **已處理** | §4.5，註解＋規格列＋21 KiB／32 KiB 測試 |
| SCR-026 | 資訊（設計部範圍） | **已處理**（repo 版 v0.2.1 含 §13.4.1）；與本機版的差異另列 SCR-041 | §4.6 |
| SCR-027 | 資訊 | **已記錄、接受** | 方法規格 §13.4「標頭上限」列寫明 431 由 net/http 回 text/plain，不經 router／ErrorMapper／安全標頭中介層；內容只是固定狀態文字，不洩漏資訊。不需修改 |

### 6.2 新發現（SCR-040 起）

**v0.1.2 狀態**：SCR-040、SCR-041 **已關閉**（#17 於 2026-10-19 14:00:38（UTC+8）合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`（parents `ed62039`＋`e6d9ab1`，合併的 head 即安全部核准的 `e6d9ab1`）；main tree 與 `e6d9ab1` tree 相同，`git diff e6d9ab1 433140c` 為空）。SCR-041 的後續（設計部本機檔同步為 repo 版、生效日校正）屬設計部範圍。SCR-042 仍開放（流程）。

**v0.1.1 狀態**：SCR-040 **已修正（`347a027` R2），待 #17 合併後關閉**；SCR-041 **已處理，待 #17 合併後關閉**（設計部簽核紀錄聲明以 PR 版 v0.2.1 為準、本機版不另行提交，合併後同步；合併日若非 2026-10-19 由設計部文件 PR 校正生效日）；SCR-042 開放（流程，本版 §7 以「綁定完整 SHA＋重新送核」處理）。

**SCR-040（資訊；設計部範圍）：方法規格 §13.2 仍寫 RateLimitGuard 是 stub**
- 影響：PR 版 §13.2 的 `RateLimitGuard.checkCreate／checkRedirect` 列仍寫「rate_limited（現 stub 一律放行）／行為延後 ENG-011」，`ObservabilityHooks.on*` 列仍寫「no-op；行為延後 ENG-012」。#7（ENG-011）、#11（ENG-012）都已合併。這與同一份文件 §13.4.1「建立端點另有 30 次／分鐘／來源之速率限制」互相矛盾，容易讓後續審查誤判 C1 前提。不影響程式行為。
- 重現：`git show 3b1aa5b:docs/02-architecture/06-method-specs-v0.2.md | grep -n "stub\|延後"`（第 189、190 行）。
- 修補建議：設計部下次修訂規格時，更新這兩列為已實作的現況（固定窗、NewGuard(create, redirect)、fail-closed 等）。不必為此改動 #17 的 head。

**SCR-041（資訊；設計部範圍）：本機設計部規格與 PR 版 v0.2.1 不一致**
- 影響：同為 v0.2.1，但本機 `docs/02-architecture/06-method-specs-v0.2.md`（2026-10-07 版）與 PR 版（2026-10-19，研發部草擬）內容不同（§4.6）。PR 版的修訂紀錄已寫「安全部意見 v0.1 附條件同意，C1～C5 已滿足」與「本裁定自 2026-10-19 生效」。這些敘述以本文件判定為前提，且要在 #17 合併後才成立；如果合併日延後，日期會不準。實質安全內容一致，沒有安全影響。
- 重現：`diff <本機檔> <(git show 3b1aa5b:docs/02-architecture/06-method-specs-v0.2.md)`。
- 修補建議：設計部在自己的複核關卡確認 PR 版為準，並在 #17 合併後同步本機檔；如果合併日不是 2026-10-19，生效日期改以合併日為準。

**SCR-042（資訊；流程）：ruleset 不會在新 push 後撤銷核准**
- 影響：main ruleset 是 `dismiss_stale_reviews_on_push: false`、`require_last_push_approval: false`。如果核准之後 head 又有新 commit，核准仍然有效，GitHub 不會擋下。「head 不得變動」只能靠流程把關。
- 重現：`curl -s https://api.github.com/repos/444112029012/short-url/rules/branches/main`。
- 修補建議：合併執行者在按下合併前，確認 PR head 仍是 `3b1aa5b68ea2251a6551f8901e657c95b3657410`（見 §7）。長期可評估開啟 `dismiss_stale_reviews_on_push`（屬總協調／維運決定，不阻擋本 PR）。

沒有發現新的程式或測試缺陷。

## 7. 判定與合併條件（v0.1.1）

**判定：通過，僅適用於最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`。**

1. **核准綁定**：本核准只適用於 `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`。合併當下 PR head 必須等於這個 SHA。
2. **任何新 commit 都要重新送核**：只要 head 改變（包括為了 strict=true 同步 main 的 merge commit、只改註解或文件的 commit），本核准即不適用；新 head 必須經安全部**重新書面確認**後才可合併。若變更只是「只帶進 main 改動的同步 merge commit」，確認範圍為：
   - parents 是 `e6d9ab1` 與當時的 main；
   - `git show --remerge-diff <新 head>` 為空；
   - 以下 blob 不變：`main.go` `e27a120a2227e35b34d2eb1b0f0c0f1bf50c5306`、`main_test.go` `cb9987afdaa2593dba0d65d61fd1f7b77071a199`、`adapter.go` `216e43d71680a4b6a7ddc78bb9bbd2149c3fc325`、`bodylimit_test.go` `e1dec979984273a486992e896a51f42565bea29f`、`06-method-specs-v0.2.md` `110b9b384bc0c57849356799bb9a5ffc4f2f54a4`、`09-signoff-…-v0.2.1.md` `827e74ffa7433df7d7d46656412e1c93f5363b01`；
   - 新 head 的 run 四項 success、每個 job 的 head_sha 等於新 head、gosec `Issues : 0`。
   除此之外的任何變更，都要重新做安全審查。設計部簽核紀錄也聲明 H 之後有其他 commit 或 main 前進時須重新複核與重簽，該關卡由設計部自行判斷。
3. **strict=true 與第 1、2 點的關係**：目前 main＝`ed62039`，就是 H 的第二個 parent，PR 已與 main 同步，可在不變更 head 的情況下合併。若合併前 main 再前進（例如 #18 先合併），GitHub 會要求同步，此時依第 2 點處理，不另成條件。
4. **審查部關卡**：至少 1 個核准 review（ruleset 24629854；目前 `mergeable_state: blocked` 即因缺核准）。ruleset 不會在 push 後撤銷核准（SCR-042），合併執行者須在按下合併前確認 head 仍是 `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`。
5. **設計部關卡**：設計部已以 `e6d9ab1` 內的簽核紀錄簽核 H 的方法規格 v0.2.1；其有效性由設計部負責，本文件不代替。
6. 合併當下四項必要檢查（secrets-gitleaks、sast-gosec、sca-govulncheck、unit-test）在合併的 head 上都必須是 success。

v0.1 對 `3b1aa5b` 的合併條件已由本節取代（`3b1aa5b` 不再是 PR head）。

合併後的後續（不阻擋）：SCR-022、SCR-023、SCR-040、SCR-041 關閉；10 號意見書的 C3、C4（以及 C2、C5）改為已達成，由安全部另行修訂。

## 8. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-19 | 安全部 | 初版：對 PR #17 head `3b1aa5b` 的合併前確認；run 39199412428 四項 success、head_sha 相符；SCR-022／023 已修正、SCR-024／025／026 已處理、SCR-027 已記錄；C1～C5 達成；新增 SCR-040～042（資訊）；判定通過，附合併條件 |
| v0.1.1 | 2026-10-19 | 安全部 | 新增 §0「兩段核對」：`3b1aa5b → H`（`5b5017b`，parents `347a027`＋`ed62039`、remerge-diff 為空、4 個程式 blob 不變、H 與 main 只差原 5 檔、規格 sha256 `c89e336e…a687`、run 39203891218 四項 success 且 head_sha 相符、gosec Issues 0；R2／R3 只對齊程式，未放寬安全規格）與 `H → 最終 head`（唯一 parent H、只新增簽核紀錄 sha256 `806170e0…e3bd`、不含程式或設定、run 39204252238 對應 `e6d9ab1` 四項 success、本機 `go test ./...` 全部 ok）；SCR-040 已修正、SCR-041 已處理，皆待合併後關閉；判定綁定完整最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`，§7 合併條件改寫為不互相矛盾的形式（任何新 commit 含同步都要重新送核）；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-19 | 安全部 | 合併後同步：#17 於 2026-10-19 14:00:38（UTC+8）合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`（parents `ed62039`＋`e6d9ab1`，合併的 head 即安全部核准的 `e6d9ab1`）；`git diff e6d9ab1 433140c` 為空（main tree＝`b14a40eb`）；SCR-022、SCR-023、SCR-040、SCR-041 改為已關閉；SCR-042 仍開放；小版修訂不改檔名（檔名保留主版號） |
