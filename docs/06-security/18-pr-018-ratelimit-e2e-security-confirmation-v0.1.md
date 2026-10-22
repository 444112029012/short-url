# PR #18 速率限制 E2E 測試：安全確認（SEC-007）

| 欄位 | 內容 |
|------|------|
| 文件版本 | v0.1.4 |
| 日期 | 2026-10-19 |
| 作者 | 安全部（DevSecOps） |
| 狀態 | **v0.1.4：#18 已合併**（2026-10-19 14:24:23（UTC+8），merge commit `643bf8ab052ebb175e35bc365ec97863ec9f4202`，合併的 head 即核准的 `f547a6e`）；**SCR-028／029／031 已關閉**；SCR-030、032、033 仍為開放建議。v0.1.3：同步後新 head 已書面再確認，通過；核准只適用於 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`**（之後再有任何 commit 都要重新送核；SCR-028／029／031 待合併後關閉）。v0.1.2：通過，核准綁定 head `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`**；main 已前進到 `ed62039`，#18 必須同步，同步後的新 head 須經安全部重新書面確認才可合併（§6）。v0.1.1：通過（新 head `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`；SCR-028 已修正，待合併後關閉；依 §6 合併條件）。歷史：v0.1（head `9370def`）附條件通過（必改 1 項：SCR-028） |
| 對象 | PR #18 `test: 速率限制 E2E 測試（429／來源隔離／可信代理） (ENG-011)`，分支 `test/ENG-011-rate-limit-e2e` |
| 依據 | SEC-007（`01-sec-requirements-v0.2.md`）；方法規格 v0.2 §10；品保部 `07-testing/07-test-cases-v0.1.md`（UT-RL-01～02、UT-HC-03、E2E-05）；SCR-006／SCR-007／SCR-017 |
| 發現編號 | SCR-028～SCR-032（v0.1）；SCR-033（v0.1.1）；SCR-034～SCR-039 保留未使用（SCR-027 保留給 #17） |
| 方法（v0.1.1） | 只讀：GitHub API（PR、compare、run／jobs、job log）、`git fetch` 至 `/workspace/tmp/pr18-review/short-url`。**經使用者明確授權**，在本機 worktree（`/workspace/tmp/pr18-run`，head `663ee9d`）執行 `go test -race`、/64→/63 突變測試與 bind 失敗環境變數矩陣；所有暫改都已還原（`git status` 乾淨），worktree 已移除，未 push。未寫 exploit、未在 GitHub 留言／review／核准／合併／push |
| 方法（v0.1） | 只讀：GitHub API（cursor-github 讀取工具、公開 rulesets API）、`git clone` 至 `/workspace/tmp/pr18-review/`、CI job log。**未執行 PR 內任何程式碼**、未寫 exploit、未在 GitHub 留言／review／核准／合併／push |

---

## 0B. v0.1.4 合併後狀態（2026-10-19 14:59 UTC+8 查證）

- PR #18 API：`merged: true`、`merged_at` 2026-10-19T06:24:23Z＝**14:24:23（UTC+8）**（merge commit 的 commit date 為 14:24:22 +0800）、`merge_commit_sha`＝`643bf8ab052ebb175e35bc365ec97863ec9f4202`、`head.sha`＝`f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`。
- `643bf8a` 的 parents 是 `433140c`（當時的 main）與 `f547a6e`（v0.1.3 核准的 head），`git diff f547a6e 643bf8a` 為空：合併的就是核准的 SHA，符合 §6 第 1 點。
- 合併後 `internal/integration/ratelimit_e2e_test.go` blob 仍為 `0e2b6d9f1aa735c0c8b24fd75af672522af72c47`（`643bf8a` 與目前 main `141fcda` 皆同）。
- 因此 SCR-028、SCR-029、SCR-031 改為**已關閉**（§4）。SCR-030（轉研發部）、SCR-032、SCR-033 仍為開放建議，不受合併影響。
- 合併當下 ruleset 規定的審查部關卡與核准 review 由 GitHub 強制，本文件未另查。

## 0A. v0.1.3 同步後新 head 書面再確認（head `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`，2026-10-19 14:08 UTC+8 查證）

範圍依 v0.1.2 §6 第 3 點。方法：GitHub 讀取 API（PR、run／jobs、job log）、`git fetch`；**經使用者授權**在本機暫時 worktree（新 head）執行測試，執行後 `git status` 乾淨、worktree 已移除。未在 GitHub 留言／review／核准／合併／push。

| # | 項目 | 查證結果 |
|---|------|----------|
| 1 | main＝`433140c`，且為 #17 合併結果 | ✅ `git ls-remote refs/heads/main`＝`433140c5132778bf7cd731e3ecaca8187e50d678`；commit 訊息 `Merge pull request #17`，parents＝`ed62039`＋`e6d9ab1`；PR #17 API：`merged: true`、`merged_at` 2026-10-19T06:00:38Z＝**2026-10-19 14:00:38（UTC+8）**、`merge_commit_sha`＝`433140c5132778bf7cd731e3ecaca8187e50d678`。main tree 與 `e6d9ab1` tree 相同 |
| 2 | 新 head 的 parents | ✅ PR API `head.sha` 與 `refs/pull/18/head`＝`f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`（14:02:00 UTC+8）；parents＝`663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`、`433140c5132778bf7cd731e3ecaca8187e50d678`（＝當時 main）。PR 3 commits、1 檔、+838／−0，`mergeable: true`、`mergeable_state: blocked`（缺核准 review） |
| 3 | `git show --remerge-diff` 為空 | ✅ 安全部自行執行，輸出 0 byte（研發部宣稱屬實） |
| 4 | 測試檔 blob 不變 | ✅ `internal/integration/ratelimit_e2e_test.go` 在新 head＝`0e2b6d9f1aa735c0c8b24fd75af672522af72c47`，與 `663ee9d` 相同 |
| 5 | main 對 head 只差測試檔 | ✅ `git diff --name-status 433140c f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`：只有 `A internal/integration/ratelimit_e2e_test.go`（+838）。`refs/pull/18/merge`＝`56f0c5e`（parents `433140c`＋`f547a6e`），tree＝新 head tree＝`52ddd01d` |
| 6 | run 39207726472 | ✅ event `pull_request`、run `head_sha`＝`f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`、`conclusion: success`、run_attempt 1，建立 14:03:22 UTC+8。四個 job 的 head_sha 都是 `f547a6e…`、都 success：sast-gosec（117914105869）`gosec h1 OK`、`Files : 19`、`Lines : 1719`、`Nosec : 0`、**`Issues : 0`**；unit-test（117914105941）`go test ./...` 全部套件 `ok`（`internal/integration 0.341s`）；secrets-gitleaks（117914106047）`57 commits scanned.`、`no leaks found`；sca-govulncheck（117914106080）`No vulnerabilities found.`、`dependency graph compare HTTP 200`、Dependency review 無 critical／denied |
| 6 | CI 是否有 `-v` | ❌ 沒有。unit-test 仍是 `go test ./...`，log 沒有單一測試紀錄，**無法從 CI log 確認** `TestENG011_PerSourceIsolationByPeerAddress` 為 PASS；仍依 §0.2 的推論（CI 上 skip 路徑走不到＋套件 `ok`） |
| 7 | 本機 `go test -race -count=1 ./...` | ✅ Go 1.27.2 linux/amd64、未設 CI：10 個有測試的套件全部 `ok`（`internal/integration 1.622s`），`internal/ports` 無測試檔 |
| 7 | 本機 `CI=true go test -v -run TestENG011_ ./internal/integration/` | ✅ `-count=1`：**19 個 `--- PASS`（14 個頂層＋5 個子測試）、0 個 SKIP、0 個 FAIL**；`--- PASS: TestENG011_PerSourceIsolationByPeerAddress (0.01s)`；研發部「19 PASS、0 SKIP」宣稱屬實（本機結果，不等於 CI runner） |

**v0.1.2 §6 第 3 點四項確認全部成立**（parents、remerge-diff、blob、四綠且 head_sha 相符），同步 commit 沒有夾帶 main 以外的改動。SCR-028、SCR-029、SCR-031 維持「已修正／已處理，待合併後關閉」；SCR-033 維持資訊級、開放，留待後續 PR 處理。

## 0. v0.1.1 新 head 複核（head `663ee9d`，2026-10-19 12:48 UTC+8 查證）

### 0.1 SHA、範圍與 run

| 項目 | 宣稱 | 查證結果 |
|------|------|----------|
| head | `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0` | ✅ PR API `head.sha` 與 `git ls-remote refs/pull/18/head` 相同；commit 時間 2026-10-19 12:40:51 UTC+8，parent＝`9370def` |
| base | `24b8da0` | ✅ `base.sha`＝`24b8da0252ac8cbf1c5615592eb9ad18996f3f36`；merge-base 同 |
| main | 仍為 `24b8da0` | ✅ `git ls-remote refs/heads/main`＝`24b8da0`（12:48 再確認） |
| 與 9370def 比較 | 只多 1 commit、只改測試檔 +234／−20 | ✅ compare API：`status: ahead`、`ahead_by: 1`、`behind_by: 0`、`total_commits: 1`；唯一檔案 `internal/integration/ratelimit_e2e_test.go`（modified，+234／−20）。`git diff 9370def 663ee9d -- .github go.mod go.sum ':!*_test.go'` 為空：非測試程式、workflow、go.mod、go.sum 都沒有變動。PR 整體：2 commits、1 檔、+838／−0 |
| merge ref | — | ✅ `refs/pull/18/merge`＝`3ca1caa`，parents＝`24b8da0`＋`663ee9d`；merge tree＝head tree＝`f5174e94` |
| mergeable | — | `mergeable: true`、`mergeable_state: blocked`（仍缺 1 個核准 review） |
| run | 39200695930 四項全綠 | ✅ event `pull_request`、`head_sha`＝`663ee9d…`、`conclusion: success`、run_attempt 1，建立 2026-10-19 12:41:00 UTC+8。四個 job 的 `head_sha` 都是 `663ee9d…` |

| job（id） | 結論 | log 證據 |
|-----------|------|----------|
| secrets-gitleaks（117892868741） | success | `53 commits scanned.`、`no leaks found` |
| sast-gosec（117892868834） | success | `mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`、`gosec h1 OK`；`Files : 19`、`Lines : 1704`、`Nosec : 0`、`Issues : 0`（不掃 `_test.go`，同 §1.1） |
| sca-govulncheck（117892868597） | success | `No vulnerabilities found.`；`dependency graph compare HTTP 200`；Dependency review 無 critical、無 denied packages |
| unit-test（117892868857） | success | `go vet ./...` 通過；`go test ./...`（**無 `-v`**）；`ok  github.com/444112029012/short-url/internal/integration 0.514s`；未見 `(cached)` |

### 0.2 來源隔離測試在 CI 是否實際執行（第 4 點）

- **log 本身無法證明。** unit-test 只跑 `go test ./...`，沒有 `-v`／`-json`，log 裡沒有 `=== RUN`／`--- PASS`／`--- SKIP` 任何單一測試的紀錄，也沒有印出環境變數。只看 log，不能斷言 `TestENG011_PerSourceIsolationByPeerAddress` 是 PASS 而不是 SKIP。
- **間接推論（靜態＋本機實測，不是 CI 直接證據）**：新 head 的檔案中只有一個 `t.Skipf`（`rlRequireSourceAddr`），而且只有在「錯誤是 `EADDRNOTAVAIL`」且「`CI`、`GITHUB_ACTIONS` 都不等於 `true`（不分大小寫）」時才會走到；GitHub Actions runner 預設會設定 `CI=true` 與 `GITHUB_ACTIONS=true`，workflow 也沒有覆寫這兩個變數（`ci-devsecops.yml` 的 `env:` 只用在 dependency graph probe）。因此在這個 runner 上，skip 路徑走不到：要嘛執行並 PASS，要嘛 `t.Fatalf` 讓套件 FAIL。套件結果是 `ok`，所以**推論**該測試已執行並 PASS。
- **建議佐證方式**（維運範圍，不列合併條件）：unit-test 改成 `go test -json ./... | tee test.json` 並上傳 artifact，或另加一步 `go test -v -count=1 -run 'TestENG011_PerSourceIsolationByPeerAddress' ./internal/integration/`；也可以在 workflow 加一步 `echo "CI=$CI GITHUB_ACTIONS=$GITHUB_ACTIONS"` 留下環境證據。

### 0.3 SCR-028 逐項核對（第 3 點）

| 設計要求 | 新 head 實作 | 判斷 |
|----------|--------------|------|
| CI（`CI` 或 `GITHUB_ACTIONS`）時一律 `t.Fatalf` | `rlInCI()`＝`strings.EqualFold(os.Getenv("CI"),"true") \|\| strings.EqualFold(os.Getenv("GITHUB_ACTIONS"),"true")`；`EADDRNOTAVAIL` 且 `rlInCI()` → `t.Fatalf` | ✅ 對 GitHub Actions 成立。偵測只接受值為 `true`，`CI=1` 之類的值會被當成本機 → 見 **SCR-033**（資訊） |
| 本機只有 `EADDRNOTAVAIL` 才可帶理由 `t.Skipf` | 非 `EADDRNOTAVAIL` 的撥號錯誤一律 `t.Fatalf`；`t.Skipf` 訊息寫明「EADDRNOTAVAIL；skipped locally only, CI fails instead」 | ✅ |
| 127.0.0.2、127.0.0.3 都用真實 bind 預檢 | `rlRequireSourceAddr` 對 peer A、peer B 各以 `net.Dialer{LocalAddr}` 實際撥號；之後再各做 `/health` 預檢，錯誤或非 200 一律 `t.Fatalf` | ✅ |
| 保留真實 TCP bind | 測試本體仍用 `s.client(t, ip)`（`Dialer.LocalAddr`），RemoteAddr 由 net/http 從 socket 填入 | ✅ |
| 分類器用 `errors.Is`／syscall，不用字串比對 | `rlSourceAddrUnavailable`＝`errors.Is(err, syscall.EADDRNOTAVAIL)`；分類器測試：包裝後的 `EADDRNOTAVAIL` 必須為 true；`ECONNREFUSED`、`EADDRINUSE`、`ETIMEDOUT`、內容相同的純字串 error、`context.DeadlineExceeded`、nil 必須為 false | ✅ 測試有效：改成字串比對（純字串案例會被誤判）或改成 `err != nil`（其他 errno 會被誤判）都會失敗 |
| 環境變數能否被空字串／`"false"` 繞過 | 本機實測（暫把 peer B 改成無法 bind 的 `192.0.2.1`，已還原）：未設、`CI=`、`CI=false`、`CI=1` → `--- SKIP`；`CI=true`、`CI=TRUE`、`GITHUB_ACTIONS=true` → `--- FAIL`（`t.Fatalf`） | ⚠️ 空字串與 `false` 被視為本機，這符合預期（非 CI）；要在 GitHub Actions 上走到 skip，必須在 workflow 同時把 `CI` 與 `GITHUB_ACTIONS` 都改成非 `true`，而 workflow 變更本身要經 PR 與 ruleset（審查＋必要檢查）才能進 main。`CI=1` 被當成本機是較寬的缺口 → SCR-033 |
| CI 上還有沒有其他可能 skip 的路徑 | 全檔只有 1 處 `t.Skipf`、沒有 `testing.Short()`、沒有 build tag；workflow 用 `go test ./...`，沒有 `-run`／`-short`／`-tags` | ✅ 無其他 skip 路徑 |
| 失敗訊息不印位址 | `rlRedactErr` 只輸出 op、errno 文字或錯誤型別；實測訊息為 `dial: cannot assign requested address` | ✅ |

**SCR-028：已修正，待 #18 合併後關閉。**

### 0.4 SCR-029（第 5 點）

- 測試專用路由 `/__test/errorlog`：只存在於測試檔 `newRateLimitServer` 包在 production router 外層的 `http.HandlerFunc`；`httpapi.NewRouterWithOptions` 等非測試程式都沒有變動（§0.1），所以**不會進入正式路由表** ✅。
- 該路由連續呼叫兩次 `WriteHeader`，讓 net/http 自己寫 `superfluous response.WriteHeader` 到擷取的 `ErrorLog`；測試以**區分大小寫**的 `strings.Contains` 斷言確實收到，否則失敗 ✅（ErrorLog 擷取路徑確定是活的）。
- 不分大小寫：只用在「禁止字串」檢查；`all` 與 needle 都轉小寫。正向斷言（`superfluous …`、`rate_limited` 建立 2／導向 1）沒有轉小寫，且保證日誌非空，所以不分大小寫**沒有讓斷言失去意義**；轉小寫只會讓禁止比對更嚴，不會更寬 ✅。
- IPv6：needle 涵蓋壓縮（`2001:db8`、`2001:db8:77::77`）、補零（`2001:0db8`）、完整 8 組展開、/64 bucket key（`2001:db8:77::/64`）、IPv4-mapped 前綴（`::ffff:`）、無冒號 hex（`20010db8`）；並實際送出大寫完整展開形式的 XFF 請求 ✅。另涵蓋小寫標頭名稱請求與 `x-forwarded-for`／`x-real-ip`／`remote_addr`／`remoteaddr`／`client_ip` ✅。
- 限制（已在測試註解與 PR 說明中寫明）：觸發的 ErrorLog 訊息本身不含位址，所以對 net/http panic／TLS 握手這類**會帶位址**的訊息仍沒有覆蓋，這部分仍由 SCR-030 追蹤。
- **SCR-029：已處理（資訊），待 #18 合併後關閉。**

### 0.5 SCR-031（第 6 點）

- `TestENG011_AdjacentIPv6Slash64sDoNotShareQuota`：`2001:db8:0:0::/64` 與 `2001:db8:0:1::/64` 只差第 64 位元；同 /64 的最高位址與完整展開寫法共用 ✅。
- `TestENG011_TrustedProxyXFFEdgeCases` 有 3 個子測試：整串可信 → 採最右側有效項；客戶端右側有無效項 → 跳過無效項；XFF 全無效 → 改用 X-Real-IP，兩者皆無則用 peer ✅。
- **本機突變測試（已還原，未 push）**：把 `internal/ratelimit/source.go` 的 `net.CIDRMask(64, 128)` 改成 `63` → `TestENG011_IPv6ForwardedClientsShareSlash64` 仍 PASS（印證 v0.1 的缺口），`TestENG011_AdjacentIPv6Slash64sDoNotShareQuota` **FAIL**（`ratelimit_e2e_test.go:632: create status 429`）；還原後 `git status` 乾淨 ✅。
- **SCR-031：已處理（資訊），待 #18 合併後關閉。**

### 0.6 本機 `go test -race`

- Go 1.27.2 linux/amd64，head `663ee9d`：
  - `env -u CI -u GITHUB_ACTIONS go test -race -count=1 -v ./internal/integration/...` → `ok … internal/integration 1.579s`，29 個 `--- PASS`（含子測試），無 SKIP／FAIL。
  - `CI=true go test -race -count=1 -v -run ENG011 ./internal/integration/...` → `ok … 1.542s`，14 個 `TestENG011_*` 與 5 個子測試全部 PASS（19 個 `--- PASS`），`TestENG011_PerSourceIsolationByPeerAddress` 為 `--- PASS (0.03s)`。
- 這是安全部 box 上的結果，不等於 CI runner 的實測。

### 0.7 新發現

| 編號 | 嚴重度 | 標題 | 處置 |
|------|--------|------|------|
| SCR-033 | 資訊（測試品質） | CI 偵測只接受 `true`，其他 CI 慣用值（如 `CI=1`）會走本機 skip 路徑 | 建議，不阻擋 |

#### SCR-033（資訊）CI 偵測值過窄

- 位置：`internal/integration/ratelimit_e2e_test.go:416-418`（commit `663ee9d`）
- 類別：其他（測試的 fail-closed 條件）；對應 SEC-007、SCR-028
- 影響：設計要求是「`CI` 或 `GITHUB_ACTIONS` 有設定」就 fail-closed，實作只接受值為 `true`。在把 `CI` 設成 `1`、`yes` 等值的 CI 系統，或 workflow 把變數覆寫成非 `true` 時，來源隔離 E2E 遇到 `EADDRNOTAVAIL` 會 skip，又回到 SCR-028 的「看不出來」問題。目前 GitHub Actions 兩個變數都是 `true`，**對現行 CI 沒有影響**。
- 重現：暫把 peer B 改成無法 bind 的位址，分別以 `CI=1`、`CI=false`、`CI=` 執行 `go test -v -run PerSource ./internal/integration/` → 都是 `--- SKIP`；`CI=true` → `--- FAIL`。
- 修補建議：改成「變數有設定且值不是明確的 `false`／`0`」就視為 CI（例如 `v, ok := os.LookupEnv("CI"); ok && v != "" && !strings.EqualFold(v,"false") && v != "0"`），或改成反向設計：只有明確設定 `SHORTURL_ALLOW_SOURCE_SKIP=1` 時才允許 skip。可在後續 PR 處理。
- 狀態：開放（資訊，不阻擋）

---

> 以下 §1～§3 為 v0.1（head `9370def`）的原始審查紀錄，保留不改；§4～§7 已依 v0.1.1 更新。

## 1. 查證過的 SHA 與 run（v0.1，head `9370def`）

| 項目 | 宣稱 | 查證結果 |
|------|------|----------|
| head | `9370def57416a57344cfd8f3caab7e422788c3a4` | ✅ PR API `head.sha` 相同；PR 只有 1 個 commit |
| base | main `24b8da0`（#7 merge commit） | ✅ `base.sha`＝`24b8da0252ac8cbf1c5615592eb9ad18996f3f36`；merge-base 也是 `24b8da0` |
| main 是否前進 | — | ✅ 未前進：main 分支 head 仍為 `24b8da0`（#7 merge，2026-10-19 12:11:09 UTC+8） |
| mergeable | — | ✅ `mergeable: true`、`rebaseable: true`；`mergeable_state: blocked`，原因為 ruleset 24629854 要求 1 個核准 review（目前 reviews 為 0），不是 CI 或衝突 |
| merge ref | — | ✅ `refs/pull/18/merge`＝`fed562f`，parents＝`24b8da0`＋`9370def`；merge tree＝head tree＝`8598e51a`（main 未前進，沒有額外改動） |
| ruleset | strict=true | ✅ 24629805 `main-required-checks` active、`strict_required_status_checks_policy: true`、四項 context 不變；24629854 `main-pr-review` active、`required_approving_review_count: 1`、`dismiss_stale_reviews_on_push: false` |
| run | 39198983049 | ✅ event `pull_request`、`head_sha`＝`9370def…`、`conclusion: success`、run_attempt 1；建立 2026-10-19 12:20:32 UTC+8。該分支只有這一個 run |
| 變更範圍 | 只新增 `internal/integration/ratelimit_e2e_test.go` | ✅ `changed_files: 1`、+624／−0；`git diff 24b8da0..9370def -- .github go.mod go.sum` 為空 |

### 1.1 四項必要檢查（run 39198983049，四個 job 的 head_sha 皆為 `9370def`）

| job（id） | 結論 | log 證據 |
|-----------|------|----------|
| secrets-gitleaks（117887782976） | success | gitleaks 8.21.2：`47 commits scanned`、`no leaks found` |
| sast-gosec（117887782769） | success | `mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`、`gosec h1 OK`；Summary `Files: 19`、`Lines: 1704`、`Nosec: 0`、`Issues: 0` |
| sca-govulncheck（117887782909） | success | govulncheck v1.8.0：`No vulnerabilities found.`；`dependency graph compare HTTP 200`；Dependency review：無 critical、無 denied packages |
| unit-test（117887783023） | success | `go vet ./...` 通過；`go test ./...`（無 `-v`）各套件 `ok`；`internal/integration 0.629s`；未見 `(cached)` |

注意：gosec 預設不掃 `_test.go`。`Files: 19` 剛好等於 repo 的非測試 `.go` 檔數（測試檔 18 個），所以本 PR 新增的測試檔**不在 gosec 掃描範圍內**。本次由安全部人工靜態審閱補足（§3.8）。

---

## 2. SEC-007 原文與本 PR 的關係

> SEC-007：對「建立短碼」與「短碼導向」端點實施基本速率限制（依來源識別，如連線 IP 或反向代理轉傳之可信客戶端位址）；超過門檻回傳明確限流錯誤且不建立短碼／不導向。具體數值由設計／維運文件化（建議起點：建立 ≤ 30 次／分／來源；導向 ≤ 120 次／分／來源）。驗收：測試環境對同一來源短時間超過文件化門檻，收到限流回應；門檻內請求仍成功；門檻寫入設計或維運設定文件。

本 PR 只加測試，不動正式程式。正式程式已在 #7（`24b8da0`）合併，並經安全部 11 號、16 號檔審查。本檔的範圍是：**這些 E2E 測試是否真的證明了宣稱的行為**，以及測試本身有沒有安全問題。

---

## 3. 逐項對照（11 個測試，數量查證相符）

| # | SEC-007 要素／宣稱 | 測試 | 判斷 |
|---|------|------|------|
| 1 | 超過門檻回明確限流錯誤，不建立短碼 | `CreateQuotaExhaustedReturns429` | ✅ 門檻 3：前 3 次 201，接著 3 次 429；資料列前後都是 3（拒絕不建碼、不重置視窗）；超限時送無效本文仍回 429（先限流再解碼，符合方法規格順序） |
| 2 | 門檻文件化、預設值 | `DefaultLimitsFromConfig` | ✅ env 空 → 第 31 次建立回 429，證明 30／分的預設值經 `config.Load`→`NewGuard` 生效。導向預設 120 沒有 E2E（`NewGuard` 單元測試覆蓋），可接受 |
| 3 | 超過門檻不導向、不計數 | `RedirectQuotaExhaustedReturns429WithoutCounting` | ✅ 門檻 2：2 次 302 且 Location 正確，接著 3 次 429；未知短碼超限也回 429；`click_count`＝2 |
| 4 | 建立與導向額度分開計算 | `CreateAndRedirectQuotasAreIndependent` | ✅ 兩個方向都測（建立用盡不擋導向；導向用盡不擋建立），各自回 429 |
| 5 | 依來源識別：不同來源各自計算 | `PerSourceIsolationByPeerAddress` | ⚠️ 斷言本身正確（127.0.0.2／127.0.0.3 真實 TCP 來源，A 用盡不影響 B，建立與導向都測）。但 bind 失敗時會 `t.Skipf`，CI 沒開 `-v`，看不出這次是否實際執行 → **SCR-028** |
| 6 | 不受信任的來源偽造 XFF／X-Real-IP 無效 | `UntrustedPeerCannotSpoofForwardingHeaders` | ✅ 兩種設定（空清單；清單不含 peer、但含 XFF 裡出現的 192.0.2.0/24）都測，每次輪換 XFF／X-Real-IP 仍共用 peer 額度；建立與導向都測；資料列＝3 |
| 7 | 受信任代理採用 XFF | `TrustedProxyForwardedForIsHonored` | ✅ 同一條代理連線後的兩個客戶端額度分開；無轉傳標頭時以代理本身計；無 XFF 時採 X-Real-IP 並對應到同一 bucket |
| 8 | 由右往左、取最右側「非可信」項目 | `TrustedProxyIgnoresClientPrependedEntries` | ✅（見 §3.1）。註：研發部說「只採用最右邊可信的項目」用字有誤，實際行為與測試都是「最右側**非**可信項目」 |
| 9 | IPv6 同一個 /64 共用；IPv4-mapped 併入 IPv4 | `IPv6ForwardedClientsShareSlash64` | ✅ 有不同 /64 的反向案例（`2001:db8:0:2::1` 可通過）；含 `[addr]:443` 形式；`::ffff:198.51.100.20` 與 IPv4 共用。邊界精度有小缺口 → **SCR-031** |
| 10 | 並發時放行數量剛好等於上限 | `ConcurrentCreatesAdmitExactlyLimit` | ✅ 40 個並發、門檻 5 → 恰 5 個 201、35 個 429（每個 429 都做完整斷言）、資料列恰 5（見 §3.3） |
| 11 | 日誌不出現用戶端位址 | `LogsDoNotContainClientAddresses` | ✅（附註）觀測 hooks 的所有輸出都有抓到；`rate_limited` 事件數（建立 2、導向 1）也有斷言。但 server ErrorLog 路徑在測試中沒有被觸發，部分比對字串區分大小寫 → **SCR-029**；正式程式的 panic 日誌另見 **SCR-030** |

### 3.1 XFF 測試涵蓋（逐項）

`trusted = 127.0.0.0/8,10.0.0.0/8`，門檻 2，四種標頭預期都對應到 `198.51.100.9`：

| 情境 | 案例 | 涵蓋 |
|------|------|------|
| 同一行前置偽造值 | `198.51.100.1, 198.51.100.9` | ✅ |
| 最右側是可信代理 | `198.51.100.2, 198.51.100.9, 10.1.2.3` | ✅ 跳過 10.1.2.3 |
| 多個 XFF 標頭行 | `"198.51.100.3"`＋`"198.51.100.9"` | ✅ 若實作改用 `Header.Get`（只取第一行）會得到 .3 而不回 429，測試會失敗 |
| 含空白／無效項＋多行＋可信尾端 | `"garbage, 198.51.100.4"`＋`"198.51.100.9, 10.0.0.1"` | ✅ |

靜態推演：「取最左」「取最右但不跳過可信」「只讀第一行」三種錯誤實作都會讓第 3 或第 4 個請求不回 429，測試會失敗。
未涵蓋（E2E 與單元測試都沒有）：整串都是可信項目時的行為（實作改用最右側有效項）、無效項位在客戶端右側（例如 `198.51.100.9, garbage`）。XFF 全部無效 → 改用 X-Real-IP 由單元測試 `TestClientSourceTrustedProxyFallsBackToXRealIP` 覆蓋。列入 SCR-031（資訊）。

### 3.2 IPv6

- 「完全不彙總（/128）」或「/65 以上」的錯誤實作：`2001:db8:0:1::1` 與 `2001:db8:0:1:ffff::2` 會分屬不同 bucket → 第 3 次不回 429 → 測試失敗 ✅。
- 「/48、/56、/60」等過寬的實作：`0:1` 與 `0:2` 會被併在一起 → `v6-other` 回 429 → 測試失敗 ✅。
- 但 `/63` 的實作會把 `…:0:1` 與 `…:0:2` 分開，所以**仍會通過**。反向案例改用只差第 64 位元的前綴（例如 `2001:db8:0:0::1` 對 `2001:db8:0:1::1`）才能把邊界釘在 /64 → SCR-031。

### 3.3 並發

- 底層確實並發：40 個 goroutine 共用一個 `http.Transport`（`MaxConnsPerHost` 預設不限），會開多條 TCP 連線；`httptest.Server` 每條連線一個 goroutine；guard 以單一 mutex 原子化「檢查＋計數」。
- 穩定性：放行數由 guard 決定，與排程無關；SQLite `SetMaxOpenConns(1)`＋`busy_timeout=5000`，5 筆寫入不會遇到 BUSY。唯一的時間依賴是 1 分鐘固定視窗：測試要跨過 60 秒才會出錯，而 CI 中整個 integration 套件只花 0.629 秒，flaky 風險可忽略。
- 限制：CI 沒有跑 `-race`。若 guard 日後失去互斥，本測試只會**機率性**抓到超放（資訊，列入 SCR-032；要改 workflow 屬維運範圍）。

### 3.4 日誌

- 擷取範圍：`NewHooksWithLogger` 的 JSON 輸出（與 main 的 `NewHooks()` 同格式，main 寫 stderr）＋ `httptest` server 的 `ErrorLog`。`rg` 確認 `internal/httpapi`、`application`、`infra` 的非測試程式沒有其他 `log`／`slog` 呼叫，所以請求路徑上的應用日誌都在擷取範圍內 ✅。
- 不足（SCR-029）：(a) 測試沒有觸發任何 server 錯誤，對 ErrorLog 的「不含位址」斷言沒有實際內容；(b) `"X-Real-Ip"`、`"X-Forwarded-For"`、`"remote_addr"` 是區分大小寫的比對，小寫或 `X-Real-IP` 寫法不會被抓到；IPv6 展開形式（`2001:0db8:…`）也不會被抓到。
- 測試訊息只印標籤，不印位址（`needleLabel`）✅。

### 3.5 429 斷言（`rlAssertLimited`）

每個 429（包括並發測試中的 35 個）都斷言：狀態 429；本文**完全等於** `{"error":{"type":"rate_limited","message":"Too many requests"}}`，並另外驗證信封只有 `error` 一個鍵、內層只有兩個鍵；`Content-Type: application/json`；SEC-008 四個標頭（nosniff、no-referrer、DENY、CSP）各**恰好一個值**；`Cache-Control` 恰為 `[no-store]`；沒有 `Location`；純 HTTP 沒有 HSTS（`HSTS_ENABLED` 設為空）；沒有 `Retry-After`。斷言是精確比對，不算太寬 ✅。研發部宣稱屬實。

### 3.6 「三種突變都會被抓到」（只做靜態推論，未實際執行）

| 突變 | 靜態推論 |
|------|----------|
| (1) `ClientSource` 不論 peer 是否可信都採用 XFF | 合理：#6 每次輪換 XFF，額度不會用盡 → 預期 429 處失敗 |
| (2) guard `b.count >= limit` 改 `>` | 合理：多放行 1 次 → #1～#4、#10 都會失敗 |
| (3) 導向改用建立的計數表 | 合理：#4 建立用盡後導向會被擋 → 失敗 |

三項推論都成立。但研發部所述「本機、未提交」的突變執行結果，安全部**無法查證**（依規定不執行 PR 程式碼）。另見 §3.1／§3.2：有兩種突變（/63 彙總、整串可信 XFF 的處理）目前測不到。

### 3.7 skip 與 CI 是否實際執行（重點殘留問題）

- `PerSourceIsolationByPeerAddress` 先用 peerA（127.0.0.2）`GET /health`，**任何錯誤**都會 `t.Skipf`；peerB（127.0.0.3）沒有預檢（失敗時會 `t.Fatal`，這部分是 fail-closed）。
- CI 的 `go test ./...` 沒有 `-v`；Go 在非 verbose 模式下，跳過的測試不會輸出 `--- SKIP` 或 Skipf 訊息。unit-test log 全文搜尋 `SKIP`／`skip`：0 筆。**這只代表看不到，不能證明沒有 skip。**
- 耗時推論：沒有這個檔案的同 base run（main push run 39198199210、docs 分支 `docs/ENG-021-readme` 的 PR run 39199211879）中，`internal/integration` 分別是 0.152 秒和 0.125 秒；本 run 是 0.629 秒（約多 0.48～0.50 秒），可以確認新加的測試**整體有執行**（也沒有 `(cached)`）。但單一測試大約只發 10 個請求、耗時是毫秒等級，**無法從套件耗時判斷這一個測試是否被 skip**。
- 環境：runner 是 `ubuntu-24.04`（image 20261012.334）。Linux 預設把 127.0.0.0/8 整段路由到 lo，可以 bind 127.0.0.2；安全部在自己的 Linux box 上用一般 socket 驗證可以 bind 127.0.0.2（不是執行 PR 程式碼），但這不等於該 runner 的實測結果。
- **結論：本次 CI 很可能有執行該測試，但無法從 log 查證。**
- 減輕因素：`internal/httpapi/ratelimit_test.go` 的 `TestSpoofedXFFDoesNotBypassCreateLimit` 已在 handler 層設定不同 `RemoteAddr`，證明「不同 peer 各自計算」，且沒有 skip。所以正式程式的來源隔離仍有 CI 證據；缺少的是「經真實 TCP、由 net/http 填入 RemoteAddr」的 E2E 證據。

### 3.8 測試本身的安全性

- 沒有寫死的憑證、token 或金鑰；位址都在文件保留段（RFC 5737／RFC 3849）；長網址用 `example.com`，正式程式不會解析或連線長網址（`rg` 確認非測試程式沒有 `Lookup`／`Dial`／`http.Get`）✅。
- 只連 loopback 上的 `httptest` server，沒有對外連線 ✅。
- 檔案：SQLite 放在 `t.TempDir()`，會自動清除；server、store、transport 都有 `t.Cleanup` 關閉 ✅。
- `t.Setenv` 會自動還原環境變數，也不能和 `t.Parallel` 一起用，所以不會互相污染 ✅。
- gitleaks 無洩漏；gosec 不掃測試檔（§1.1），以本節人工審閱代替。

---

## 4. 發現項

| 編號 | 嚴重度 | 標題 | 處置（v0.1） | 狀態（v0.1.1，head `663ee9d`） |
|------|--------|------|------|------|
| SCR-028 | 低（測試品質；不影響正式程式安全） | 來源隔離 E2E 在 bind 失敗時靜默 skip，CI 無法證明有執行 | **必改**（合併條件） | **已修正，待合併後關閉**（§0.3）→ **v0.1.4：已關閉**（#18 合併，`643bf8a`） |
| SCR-029 | 資訊 | 日誌測試的 ErrorLog 路徑未被觸發；比對字串區分大小寫 | 建議 | **已處理，待合併後關閉**（§0.4）→ **v0.1.4：已關閉**（#18 合併，`643bf8a`） |
| SCR-030 | 資訊（既有，非本 PR 引入） | 正式程式 panic 時，net/http 預設日誌會帶出 peer 位址 | 建議；轉研發部，後續由安全部登錄殘餘風險 | 開放（既有，非本 PR 範圍；G5 前待辦） |
| SCR-031 | 資訊 | IPv6 邊界與 XFF 邊緣情境的突變測不到 | 建議 | **已處理，待合併後關閉**（§0.5，/63 突變實測會失敗）→ **v0.1.4：已關閉**（#18 合併，`643bf8a`） |
| SCR-032 | 資訊 | CI 沒有 `-race`、gosec 不掃測試檔，並發與測試碼的保證較弱 | 建議（維運範圍） | 開放（維運範圍；本次 run 仍無 `-race`、gosec 仍不掃測試檔） |
| SCR-033 | 資訊 | CI 偵測只接受 `true`（§0.7） | — | 開放（建議，不阻擋） |

### SCR-028（低，必改）來源隔離 E2E 會靜默 skip

- **影響**：在無法 bind 第二個 loopback 位址的環境，或 `/health` 預檢因任何其他原因失敗時（例如 server 回應異常、逾時），測試會 skip 而不是失敗。CI 只顯示套件 `ok`，所以「不同來源經真實 TCP 各自計算」這項 SEC-007 證據可能悄悄消失，審查者也看不出來。正式程式不受影響（handler 層有不會 skip 的單元測試）。
- **重現方式（不需攻擊程式）**：在 macOS 預設環境（只有 127.0.0.1）跑 `go test ./internal/integration/`，結果是 `ok`；加上 `-v` 才會看到 `--- SKIP: TestENG011_PerSourceIsolationByPeerAddress`。在 CI log 中看不出差別。
- **安全部立場（兩種做法的比較）**：
  - 改成「在 handler 層直接設定 RemoteAddr」：等於重做 `internal/httpapi/ratelimit_test.go` 已有的測試，也失去這個 E2E 唯一新增的價值（net/http 從真實 socket 填入 RemoteAddr）。**不採用為主要做法。**
  - **採用**：保留真實 bind，但在 CI 上 fail-closed：
    1. 偵測到 `CI=true` 或 `GITHUB_ACTIONS=true` 時，預檢失敗一律 `t.Fatalf`（訊息寫明原因，不印位址）；
    2. 本機開發環境仍可 skip，但必須 `t.Logf`／`t.Skipf` 寫明原因（現有的 Skipf 已帶原因，保留即可）；
    3. 把 skip 條件縮小到「無法使用來源位址」這類錯誤（例如先用 `net.Dialer{LocalAddr:…}` 撥號、檢查 `EADDRNOTAVAIL`／`*net.OpError` 的 dial 階段），其他錯誤一律 `t.Fatal`；
    4. 預檢同時涵蓋 peerB（127.0.0.3）。
  - 列為**必改**的理由：只動測試檔、幾行就能完成，而且能補上 CI 對 SEC-007 來源隔離 E2E 證據的可驗證性。嚴重度維持「低」，因為正式程式的行為已有其他測試證明。
  - **不要求** CI 改開 `-v`。若想在 workflow 加 `-v` 或 `-json` 來保留測試明細，屬於維運範圍，只列為建議。
- **驗收**：新 head 只改 `internal/integration/ratelimit_e2e_test.go`；靜態可見 CI 環境變數判斷與 `t.Fatalf` 分支；新 head 的 run 四項綠燈。

### SCR-029（資訊）日誌測試的覆蓋缺口

- **影響**：如果日後 ErrorLog 路徑或其他寫法（小寫標頭名稱、`X-Real-IP` 寫法、IPv6 展開形式）把位址寫進日誌，現有測試抓不到。目前正式程式沒有這種輸出，所以不是現存缺陷。
- **重現**：靜態閱讀 `LogsDoNotContainClientAddresses` 第 606～613 行（`strings.Contains` 區分大小寫）；測試流程中沒有任何會觸發 `http.Server.ErrorLog` 的請求。
- **修補建議**：比對前先把 `out` 轉小寫，needle 也改小寫（`x-forwarded-for`、`x-real-ip`、`remote_addr`）；加入 `2001:0db8` 展開形式；在註解中說明 ErrorLog 只用來攔截意外輸出。

### SCR-030（資訊，既有）panic 日誌帶出 peer 位址

- **影響**：`cmd/shorturl/main.go` 的 `http.Server` 沒有設定 `ErrorLog`，router 也沒有 recover middleware。handler 一旦 panic，net/http 會用標準 log 寫 `http: panic serving <RemoteAddr>` 到 stderr。這和本 PR「日誌不含用戶端位址」的宣稱範圍不同（測試沒有涵蓋 panic 路徑）。依隱私文件（`01-requirements/04-privacy-compliance-v0.1.md`、TB-04），基礎設施日誌含 IP 時要保存 ≤14 日且僅供維運，所以目前不違反 SEC-005（SEC-005 規範的是 Authorization／Cookie）。
- **重現**：靜態閱讀 `main.go`（`srv` 未設 `ErrorLog`）、`internal/httpapi/router.go`（沒有 Recoverer），以及 Go 標準庫 `net/http` 的 conn.serve panic 處理。
- **修補建議**（研發部，另開 PR）：加 recover middleware，回 `internal` 錯誤信封、只記事件類型；或者給 `ErrorLog` 一個會遮罩位址的 handler。安全部之後更新殘餘風險登錄時一併登錄（本次不改既有檔）。

### SCR-031（資訊）IPv6 /64 邊界與 XFF 邊緣情境

- **影響**：用 /63 彙總，或整串 XFF 都是可信項目時處理錯誤的實作，仍可通過本 PR 與現有單元測試。
- **重現**：§3.1、§3.2 的靜態推演。
- **修補建議**：反向案例改用只差第 64 位元的前綴（`2001:db8:0:0::1` 對 `2001:db8:0:1::1`）；加入「XFF 全為可信項目」與「`198.51.100.9, garbage`」兩個案例（可放單元測試）。

### SCR-032（資訊）CI 保證的範圍

- **影響**：(a) CI 沒有 `-race`，guard 的並發正確性只靠並發測試機率性偵測；(b) gosec 預設不掃 `_test.go`，測試碼不在 SAST 範圍內。
- **修補建議**：屬維運範圍，只列建議：考慮在 unit-test job 加 `go test -race ./...`（或另開非必要 job）；是否對測試碼加 `gosec -tests` 由維運評估。本 PR 不因此退回。

---

## 5. 判定

**v0.1.4：#18 已依核准 head `f547a6e` 合併（`643bf8a`，14:24:23 UTC+8），判定不變；本文件的合併條件已執行完畢，僅作歷史保留（§0B）。**

**v0.1.3：通過，核准只適用於 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`。** 同步後新 head 依 v0.1.2 §6 第 3 點完成書面再確認（§0A）。之後再有任何 commit（包括再次同步 main），本核准即不適用，都要重新送安全部確認。

**v0.1.1：通過**（head `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`）。

- 唯一必改 SCR-028 已依設計修正（§0.3）：CI 上 fail-closed、本機只在 `EADDRNOTAVAIL` 時帶理由 skip、兩個來源位址都真實 bind 預檢、分類器用 `errors.Is` 且有有效測試；CI 上沒有其他 skip 路徑。
- 新 head 相對 `9370def` 只改 1 個測試檔；非測試程式、`.github`、`go.mod`、`go.sum` 無變動；run 39200695930 四項 success、四個 job 的 head_sha 都是 `663ee9d`，gosec `Issues : 0`。
- SCR-029、SCR-031 已處理；新發現 SCR-033 為資訊級，不阻擋。
- 保留事項：CI log 無 `-v`，來源隔離測試 PASS 是靜態＋本機實測推論，不是 log 直接證據（§0.2）。

**v0.1（head `9370def`，歷史）：附條件通過。**

- 正式程式未變；四項必要檢查在 head `9370def` 全綠，且 run 的 head_sha 相符；只改一個測試檔，`.github`、`go.mod`、`go.sum` 都沒有變動。
- 11 個測試中，有 10 個能有效證明宣稱的 SEC-007 行為，斷言精確；429 的各項斷言宣稱屬實。
- 唯一必改是 SCR-028（低）：CI 上無法查證來源隔離 E2E 是否實際執行，需改為 CI 上 fail-closed。其餘都是資訊級建議，不阻擋合併。
- 不影響正式程式安全，所以不退回。

## 6. 合併條件（v0.1.3）

1. **核准綁定**：安全部核准只適用於 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`。合併當下 PR head 必須等於這個 SHA。
2. **之後任何 commit 都要重新送核**：只要 head 再改變（包括 main 再前進後的同步 merge commit），本核准即不適用。若只是「只帶進 main 改動的同步 merge commit」，確認範圍同 v0.1.2 第 3 點，但 parents 改為 `f547a6e` 與當時的 main：remerge-diff 為空、`ratelimit_e2e_test.go` blob 仍為 `0e2b6d9f1aa735c0c8b24fd75af672522af72c47`、新 head 四項 success 且每個 job 的 head_sha 相符。其他任何變更重新安全審查。
3. **strict=true 與第 1、2 點的關係**：目前 main＝`433140c`，就是新 head 的第二個 parent，PR 已與 main 同步，可在不變更 head 的情況下合併。若合併前 main 再前進，依第 2 點處理，不另成條件。
4. **其他關卡**：審查部關卡＋1 個核准 review（ruleset 24629854；目前 `mergeable_state: blocked`）。ruleset 不會在 push 後撤銷核准，合併執行者須在按下合併前確認 head 仍是 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`。
5. 合併當下四項必要檢查在合併的 head 上都必須是 success。
6. 資訊級發現（SCR-030、SCR-032、SCR-033）不阻擋合併；SCR-030 轉研發部。

### 6.1 v0.1.2 合併條件（歷史；第 3 點已於 v0.1.3 對 `f547a6e` 完成確認）

v0.1.1 的條件 2「head 必須是 `663ee9d`、不得變動」與條件 3「依 strict 同步」互相矛盾（同步必然改變 head）。v0.1.2 改寫如下，取代 v0.1.1 的條件 1～6。

1. **SCR-028**：✅ 已在 head `663ee9d` 修正並經安全部複核（§0.3）。
2. **核准綁定**：安全部核准只適用於 `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`。
3. **唯一允許的 head 變更**：只有「只帶進 main 改動的同步 merge commit」可以把 head 從 `663ee9d` 改成新 head。新 head 必須經安全部**重新書面確認**後才可合併，確認範圍為：
   - parents 是 `663ee9d` 與當時的 main；
   - `git show --remerge-diff <新 head>` 為空；
   - 測試檔 blob 不變：`internal/integration/ratelimit_e2e_test.go` 在新 head 仍是 `0e2b6d9f1aa735c0c8b24fd75af672522af72c47`（安全部 2026-10-19 13:28 UTC+8 以 `git rev-parse 663ee9d:internal/integration/ratelimit_e2e_test.go` 查證；協調者所述 `0e2b6d9f` 相符）；
   - 新 head 的 run 四項必要檢查（secrets-gitleaks、sast-gosec 含 `gosec h1 OK` 與 `Issues : 0`、sca-govulncheck、unit-test）都是 success，且每個 job 的 head_sha 等於新 head。
4. **其他任何變更都要重新安全審查**：同步 merge commit 以外的 commit、同步時夾帶的任何改動、或任何會讓第 3 點某項不成立的情況，一律重新做安全審查，不適用第 3 點的簡化確認。
5. **現況：同步為必然**：main 已從 `24b8da0` 前進到 `ed620394b1c9986cc4ac243eb1543c8b677d92cb`（#16 合併，2026-10-19 13:17:22 UTC+8），PR API 顯示 `mergeable_state: behind`。依 strict=true（ruleset 24629805），#18 不同步就無法合併，所以必須走第 3 點。若 #17 先合併，main 會再前進，屆時以最新 main 同步即可，確認範圍相同。
6. **其他關卡**：還要通過審查部關卡，並依 ruleset 24629854 取得 1 個核准 review。ruleset 不會在 push 後撤銷核准（`dismiss_stale_reviews_on_push: false`），合併執行者須在按下合併前確認 head 等於安全部書面確認過的那個 SHA。
7. 合併當下四項必要檢查在合併的 head 上都必須是 success。
8. 資訊級發現（SCR-030、SCR-032、SCR-033）不阻擋合併；SCR-030 轉研發部。

## 7. 無法查證之處

- （v0.1.3）run 39207726472 仍無 `-v`／`-json`：CI 上 `TestENG011_PerSourceIsolationByPeerAddress` 為 PASS 只能推論；本機 `CI=true` 的 19 PASS／0 SKIP 不等於 runner 結果。
- （v0.1.3）審查部關卡與核准 review 狀態。
- （v0.1.1）CI 中 `TestENG011_PerSourceIsolationByPeerAddress` 是 PASS 的 **log 直接證據**：run 39200695930 無 `-v`／`-json`，也沒有印出 `CI`／`GITHUB_ACTIONS` 的值；目前是依「GitHub runner 預設設定兩者為 `true`＋CI 上 skip 路徑走不到＋套件 `ok`」推論（§0.2）。
- （v0.1.1）研發部 PR 說明所列其他突變（/65、全可信鏈回傳 nil／取最左、遇無效項即停止、拿掉 ErrorLog 擷取）與 `-count=3`：安全部只實測了 /63 突變與 `-race -count=1`。
- （v0.1，已由 v0.1.1 部分補足）研發部本機的 `-race -count=5`、第一輪三種突變的實際執行結果。
- （v0.1，已由 v0.1.1 間接補足）ubuntu-24.04 runner 實際能否 bind 127.0.0.2：若不能，CI 會 `t.Fatalf` 而不是 `ok`。
- 審查部關卡的狀態（不在安全部查證範圍）。

## 8. 修訂紀錄

| 版本 | 日期 | 作者 | 說明 |
|------|------|------|------|
| v0.1 | 2026-10-19 | 安全部 | 初版：查證 PR #18 head `9370def`／base `24b8da0`（main 未前進）、mergeable true（blocked 為缺 review）、merge tree＝head tree、run 39198983049 四項 success 且 head_sha 相符（gosec h1 OK、Issues 0；govulncheck 無漏洞；gitleaks 無洩漏）、只新增 1 個測試檔；SEC-007 逐項對照 11 個測試；新增 SCR-028（低，必改）、SCR-029～SCR-032（資訊）；判定**附條件通過** |
| v0.1.1 | 2026-10-19 | 安全部 | 新 head 複核：head `663ee9da5f5f3543a18bb6a2b47a52a3ab8f04b0`／base `24b8da0`（main 未前進）；compare `9370def...663ee9d` 只多 1 commit、只改 `internal/integration/ratelimit_e2e_test.go`（+234／−20），非測試程式／`.github`／go.mod／go.sum 無變動；run 39200695930 四項 success、四個 job head_sha 相符、gosec h1 OK、`Issues : 0`；CI log 無 `-v`，來源隔離 PASS 只能推論（§0.2）；經使用者授權本機執行 `go test -race`、/63 突變（Adjacent 測試失敗）、bind 失敗環境變數矩陣（皆已還原）；SCR-028 已修正、SCR-029／031 已處理，皆待合併後關閉；新增 SCR-033（資訊）；判定改為**通過**，合併條件更新；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-19 | 安全部 | 回應審查部 S4：§6 合併條件改寫，解除「head 不得變動」與「依 strict 同步」的矛盾。核准綁定 `663ee9d`；只有「只帶進 main 改動的同步 merge commit」可變更 head，新 head 須經安全部重新書面確認（parents＝`663ee9d`＋當時 main、remerge-diff 為空、測試檔 blob `0e2b6d9f1aa735c0c8b24fd75af672522af72c47` 不變、四項 success 且 head_sha 相符）；其他任何變更重新安全審查；記錄 main 已前進到 `ed62039`（#16），#18 為 `behind`、必須同步；小版修訂不改檔名（檔名保留主版號） |
| v0.1.3 | 2026-10-19 | 安全部 | 同步後新 head 書面再確認（§0A）：head `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`，parents `663ee9d`＋main `433140c`（#17 merge，2026-10-19 14:00:38 UTC+8）；remerge-diff 0 byte；測試檔 blob `0e2b6d9f1aa735c0c8b24fd75af672522af72c47` 不變；main 對 head 只差該測試檔；run 39207726472 四項 success、head_sha 相符、gosec Issues 0（CI 仍無 `-v`）；本機 `go test -race ./...` 全部 ok、`CI=true -v -run TestENG011_` 19 PASS／0 SKIP；判定通過並綁定 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`（之後任何 commit 都要重新送核），§6 改為 v0.1.3 條件、v0.1.2 條件留作歷史；SCR-028／029／031 待合併後關閉，SCR-033 維持資訊；小版修訂不改檔名（檔名保留主版號） |
| v0.1.4 | 2026-10-19 | 安全部 | 合併後狀態（§0B）：#18 於 2026-10-19 14:24:23（UTC+8）合併，merge commit `643bf8ab052ebb175e35bc365ec97863ec9f4202`（parents `433140c`＋核准 head `f547a6e`，與核准 head 的 tree 無差異）；測試檔 blob 仍為 `0e2b6d9f…`；SCR-028／029／031 關閉；其他內容不變。小版修訂不改檔名（檔名保留主版號） |
