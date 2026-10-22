---
文件：PR #7（ENG-011 依來源限制建立與轉址頻率）最終 head 合併前安全確認
版本：v0.1.2
狀態：**通過；v0.1.2：#7 已依條件以 head `fc9a7a7` 合併（2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`（parents：main `3d3a9cd`＋PR head `fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc`）），main 四綠**（head `fc9a7a7`；SCR-017 已修正，待 #7 合併後關閉；附合併條件，見 §6）。v0.1（head `99b0879`）曾判退回，必改 1 項已於 `fc9a7a7` 完成
負責角色：安全部
作者：安全部
日期：2026-10-19
最後更新：2026-10-19
對應需求：SEC-005、SEC-007；NFR-003；D-08／D-09／D-10；前次審查 [11](./11-pr-007-008-security-review-v0.1.md) v0.1.4（SCR-006／007／008／011／017／018）；[12](./12-pr-011-security-review-v0.1.md) SCR-012；[05](./05-residual-risk-register-v0.1.md) RR-007；門檻 [03](./03-ci-security-gates-v0.3.md)
專案代號：SHORTURL
---

# PR #7 最終 head 合併前安全確認（v0.1）

> **v0.1.1 摘要**：最終 head 為 `fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc`（自 `99b0879` fast-forward 1 個 commit，只加 // 註解），run 39196278197 四項 success 且 head_sha 相符，**判定：通過**。複核細節見 §8；§1～§5 為 v0.1 對 `99b0879` 的審查，仍然有效（`fc9a7a7` 沒有改動任何可執行程式）。

## 1. 範圍與方法

- 對象：`444112029012/short-url` PR #7，head `99b087951c2007595138e5ccbe6e1da11ed76634`，相對前次複核 head `ef9fcb8cd4999cc40dab990dda0d091a851c8615`（11 號檔 v0.1.2）新增 4 個 commit。
- 方法：全程唯讀。cursor-github 讀取工具（`get_pull_request`、`list_check_runs_for_ref`、`get_workflow_run`、`get_job_logs`）、GitHub 公開 REST API（ruleset、PR #11、`rules/branches/main`），以及本機唯讀 clone `/workspace/tmp/su-pr7`（`pull/7/head`、`pull/7/merge`）上的 `git diff`、`git show --remerge-diff`、`git merge-tree --write-tree`。
- **沒有**在 GitHub 留言、review、approve、合併或 push；沒有執行 PR 內程式碼，也沒有撰寫攻擊程式。本機未跑 `go test`：CI 的 unit-test 已在 merge ref 上執行，而 merge ref 的 tree 與 head tree 完全相同（見 §2），本機重跑沒有增加證據價值。
- 查證時間：2026-10-19 約 11:40～11:55（UTC+8）。

## 2. 查證的 SHA 與 run

| 項目 | 宣稱 | 查證結果 |
|---|---|---|
| PR head | `99b0879…6634` | ✅ 相符（`get_pull_request`：state open、未合併、draft false、commits 8、changed_files **13**、+935／−39） |
| PR base | main `3d3a9cd…913b` | ✅ 相符；`origin/main` 最新即 `3d3a9cd`（PR #11 的 merge commit） |
| #11 合併時間 | 11:34（UTC+8） | ✅ `merged_at` 2026-10-19T03:34:13Z＝**11:34:13（UTC+8）**，`merge_commit_sha`＝`3d3a9cd` |
| mergeable | true | ✅ `mergeable: true`、`rebaseable: true`。另：`mergeable_state: blocked`，原因是 ruleset 24629854 要求 1 個核准 review（`required_approving_review_count: 1`），不是 CI 或衝突 |
| run 39195448908 | head_sha＝head，success | ✅ workflow `ci-devsecops`（`.github/workflows/ci-devsecops.yml`），event `pull_request`，`head_sha`＝`99b0879…`，`pull_requests[0]`＝#7（head `99b0879`／base `3d3a9cd`），run_attempt 1，completed／**success**；建立 11:36:01、完成 11:37:23（UTC+8） |
| 測試的 merge ref | — | `3d4debac1f56fcf6eaa06bed8bf557d40c7eda16`＝Merge `99b0879` into `3d3a9cd`；四個 job 的 checkout 都是 `HEAD is now at 3d4deba`；merge ref tree `06b3b28…`＝head tree `06b3b28…`（**完全相同**，綠燈等同 head 本身） |
| 必要檢查（ruleset 24629805） | 四項綠燈 | ✅ `list_check_runs_for_ref(99b0879)` 共 4 筆，`head_sha` 都是 `99b0879`：`secrets-gitleaks`（job 117877121960）、`sast-gosec`（117877121811）、`sca-govulncheck`（117877121769）、`unit-test`（117877121817），皆 completed／success |
| ruleset strict | strict=true | ✅ 公開 API 讀取 ruleset 24629805 `main-required-checks`：enforcement active、`strict_required_status_checks_policy: true`、四項 context 名稱不變、`updated_at` 2026-10-07 14:33:28（UTC+8） |

### 2.1 掃描 log 重點

| 檢查 | log 證據 |
|---|---|
| sast-gosec | Go 1.27.2；`go install …gosec@v2.29.1-0.20261015133444-707f2271f665`；h1 比對步驟輸出 `mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=` 與 **`gosec h1 OK`**；`gosec -fmt=text -no-fail=false -severity=high ./...` 結果 Files 19、Lines 1692、Nosec 0、**Issues 0** |
| sca-govulncheck | Go 1.27.2；govulncheck v1.8.0 **`No vulnerabilities found.`**；`dependency graph compare HTTP 200`（fail-closed 探測通過）；dependency-review：未偵測到 critical 以上漏洞套件、無 denied 套件 |
| secrets-gitleaks | gitleaks 8.21.2；`gitleaks detect --source . --verbose --exit-code 1`：41 commits scanned、**no leaks found** |
| unit-test | Go 1.27.2；`go vet ./...` 通過；`go test ./...` 全部 ok（application、config、domain、health、httpapi、infra/sqlite、integration、observability、ratelimit；cmd/shorturl 與 ports 無測試檔）。CI 未加 `-race`（既知，與前次相同） |

## 3. 新增 commit 查證（ef9fcb8..99b0879）

| commit | 研發宣稱 | 查證結果 |
|---|---|---|
| `b08417f` | 合入 main `e9f4396` | ✅ parents＝`ef9fcb8`＋`e9f4396`（#8 merge）。`git show --remerge-diff` **無輸出**；`git merge-tree --write-tree ef9fcb8 e9f4396` 重算 tree `77b9d15…`＝實際 tree `77b9d15…`。**沒有文字衝突，也沒有夾帶額外改動**。`.env.example` 兩邊都有改（#7 的 TRUSTED_PROXY_CIDRS 註解、#8 的 HSTS_ENABLED），自動合併後兩者並存 |
| `92545cf` | security_headers_test.go 3 處 adapter 呼叫；integration/api_test.go:294 改 `NewGuard(30,120)`、可信代理傳 nil | ✅ 只動 2 個測試檔（+5／−4）。`security_headers_test.go`：第 87 行（`blockingLimiter{}` 的 adapter 加 `nil`）、第 99～102 行（`NewGuard(30, 120)`＋`nil`）、第 156 行（`NewGuard(30, 120)`＋`nil`），共 3 處；`integration/api_test.go` **第 294 行** `NewGuard(30, 120)`＋`nil` |
| `a5442a2` | 合入 main `3d3a9cd` | ✅ parents＝`92545cf`＋`3d3a9cd`（#11 merge）。`--remerge-diff` **無輸出**；`merge-tree` 重算 tree `e97b47d…`＝實際 tree。沒有文字衝突、沒有夾帶 |
| `99b0879` | SCR-012：observability_test.go 第 53、123 行改 NewGuard，第 198 行 adapter 第 6 參數傳 nil | ✅ 只動 `internal/httpapi/observability_test.go`（+3／−3）：**第 53、123 行** `ratelimit.NewGuard()`→`NewGuard(30, 120)`；**第 198 行** `NewHTTPAPIAdapter(…, hooks, nil)` |

### 3.1 非測試程式、公開簽章、workflows

- 兩個非 merge commit（`92545cf`、`99b0879`）**只改 `_test.go`**。✅
- `git diff ef9fcb8 99b0879` 中的非測試變動（`security_headers.go`、`observability/hooks.go`、`config.go`、`main.go` 1 行、`.github/workflows/ci-devsecops.yml` 等）**全部來自 main**（#8／#10／#11／#12／#15），經上述 remerge／merge-tree 檢查證實沒有額外改動。
- **#7 自身的非測試 patch 未變**：把 `git diff <merge-base 63b373d> ef9fcb8 -- ':!*_test.go'` 與 `git diff 3d3a9cd 99b0879 -- ':!*_test.go'` 逐行比對，差異只有 blob index 行與 `.env.example` 的上下文行（main 帶進的 HSTS_ENABLED 兩行）。`guard.go`、`source.go`、`adapter.go` 自 `ef9fcb8` 起 **0 行變更**。✅
- 公開方法簽章：自 `ef9fcb8` 起沒有變。`NewGuard(createPerMin, redirectPerMin int)` 與 6 參數 `NewHTTPAPIAdapter(…, trusted []*net.IPNet)` 是 #7 原本的設計（相對 main 是變更，已於 11 號檔審過），這次只是把 main 帶進的舊呼叫端對齊。✅
- `.github/workflows`：`git diff 3d3a9cd 99b0879 -- .github` **為空**。✅
- `go.mod`／`go.sum`：相對 main 無變更。✅
- 因為以上都沒有新的非測試變動，不需要逐項的新增程式安全審查。

## 4. 最終 head 相對 main `3d3a9cd` 的完整複核

完整 diff 共 13 檔（+935／−39）：`.env.example`、`cmd/shorturl/main.go`、`internal/config/config.go`（只改註解）、`internal/httpapi/adapter.go`、`internal/ratelimit/guard.go`、`internal/ratelimit/source.go`，其餘 7 檔為測試。

| 項目 | 結論 | 依據 |
|---|---|---|
| SCR-006（表滿 fail-soft） | ✅ 修正仍在 | `guard.go`：`container/list`＋map；`allow` 新來源先 `dropExpiredFront`，滿了 `evictFront`（O(1)）再 `PushBack`，不回 429；重置窗口 `MoveToBack`；淘汰不記錄位址。測試 `TestGuardFullTableEvictsOldestAndAdmits`、`TestGuardFullRedirectTableEvictsOldestAndAdmits`、`TestGuardWindowResetMovesBucketToBack`、`TestGuardRedirectWindowResetMovesBucketToBack`、`TestGuardConcurrentExactQuota` 仍在 |
| SCR-007（位址聚合） | ✅ 修正仍在 | `source.go` `bucketKey`：`To4()` 先處理 IPv4／IPv4-mapped，IPv6 `Mask(CIDRMask(64,128))`＋`/64`；可信代理解析出的位址同樣經 `bucketKey`。測試 `TestIPv6PrefixSharesQuota`、`TestIPv4MappedEqualsPlainIPv4`、`TestTrustedProxyIPv6IsAggregated` 仍在 |
| SCR-008（XFF 多行） | ✅ 修正仍在 | `joinedForwardedFor` 用 `Header.Values` 依序以逗號接起；`clientFromForwarded` 由右往左取第一個非可信位址；`X-Real-IP` 只在對端可信且 XFF 無有效位址時使用，註解註明代理須覆寫。測試 `TestForwardedForSeparateHeaderLineDoesNotLetClientWin` 仍在 |
| 合併是否沖掉修正 | ✅ 沒有 | §3.1：三個檔案自 `ef9fcb8` 起 0 行變更，兩次合併 tree 可由自動合併完全重現 |
| 與 #8／#11 的簽章衝突（SCR-011／SCR-012） | ✅ 已解決 | 全 repo 測試共 9 處 `NewHTTPAPIAdapter(`，都是 6 參數；沒有殘留 `NewGuard()`；CI `go vet`、`go test` 在 merge ref 上通過 |
| 正式呼叫點 | ✅ 不變 | 非測試程式中 `CheckCreate`／`CheckRedirect` 只在 `adapter.go:83`、`:114`，都傳 `a.sourceID(r)`→`ratelimit.ClientSource(r, a.trusted)`；`main.go:53-60` 以 `ParseTrustedCIDRs(cfg.TrustedProxyCIDRs)` 的結果傳入，解析錯誤時啟動失敗（fail-closed） |
| **SCR-017 前置條件註解與測試** | ❌ **未補** | `guard.go:92-107`（`CheckCreate`／`CheckRedirect`）與 `adapter.go:22-26`（`RateLimitChecker`）的註解跟 `ef9fcb8` 一字不差，**沒有**寫明「sourceID 必須是 `ratelimit.ClientSource` 的輸出，guard 不做位址正規化」；也沒有 `SourceKey` 型別，`internal/httpapi` 沒有新增「記錄 sourceID 的假 limiter」測試。新增 4 個 commit 都只改呼叫端，PR 說明也未提 SCR-017 |
| 可信代理傳 nil 的行為 | ✅ 安全（預設不信任 XFF） | `ClientSource`：`ipInNets(peerIP, nil)` 恆為 false，直接回 `bucketKey(peerIP)`，**完全不讀 XFF／X-Real-IP**；nil 與 `ParseTrustedCIDRs("")`（回 nil, nil）行為一致，也是正式環境 `TRUSTED_PROXY_CIDRS` 未設時的預設。peer 無法解析時回原字串，空字串由 guard 對應到共用 bucket，無法繞過。adapter 層 `TestSpoofedXFFDoesNotBypassCreateLimit`（trusted 為空）與 `TestClientSourceIgnoresSpoofedHeadersWhenUntrusted` 覆蓋此路徑。測試傳 nil 只代表「不信任任何代理」，比可信設定更保守，不會讓測試掩蓋偽造問題 |
| `NewGuard(30,120)` 含義 | ✅ 正確 | 第 1 參數＝create 每分鐘每來源 30 次，第 2 參數＝redirect 每分鐘每來源 120 次，等於 `RATE_LIMIT_CREATE_PER_MIN`／`RATE_LIMIT_REDIRECT_PER_MIN` 的預設；非正值回落 30／120（`TestNewGuardDefaultsNonPositiveLimits`），零值無法關閉控制。正式環境由 `main.go:56` 從 config 帶入 |
| 新問題 | 無新增 | 合併帶進的 #8 安全標頭（含 429 路徑，`TestSecurityHeadersOn429And500`）、#11 觀測 hooks 都不接觸來源位址；SEC-005（不記錄來源 IP）維持 |

## 5. 發現項

### [SCR-017] Guard 的 sourceID 前置條件仍未標明，也沒有型別或測試保護（延續 11 號檔）

- 嚴重度：**低**　**本次列為必改**（11 號檔 v0.1.2 已「要求」在 #7 下一個 head 補註解，並列為合併前核對範圍；此 head 沒有補）
- 位置：`internal/ratelimit/guard.go:92-107`；`internal/httpapi/adapter.go:22-26`、`:148-154`（head `99b0879`）
- 影響：目前唯一的正式呼叫路徑正確，**沒有現行可利用的弱點**。日後若有人新增呼叫點（新端點、中介層）直接把 `RemoteAddr` 或原始 IPv6 字串交給 guard，IPv6 /64 聚合與 IPv4-mapped 正規化會無聲失效，SCR-007 問題（單一 /64 換位址繞過每來源限流、加速填表）會復發，現有測試抓不到。
- 重現方式（靜態、不需執行）：讀 `guard.go` `CheckCreate`／`CheckRedirect` 與 `adapter.go` `RateLimitChecker` 的註解，確認沒有任何「須先經 `ClientSource` 正規化」的敘述；`git diff ef9fcb8 99b0879 -- internal/ratelimit/guard.go internal/httpapi/adapter.go` 為空。
- 修補建議：
  1. **必改**：在 `Guard.CheckCreate`、`Guard.CheckRedirect` 與 `RateLimitChecker` 介面的註解寫明「sourceID 必須是 `ratelimit.ClientSource` 的輸出；guard 不做位址正規化，直接傳 `RemoteAddr` 或原始 IP 會讓 IPv6 /64 聚合失效」。
  2. **建議**（二擇一，可在後續 PR）：(a) 定義 `type SourceKey string`，只由 `ClientSource` 產生，介面改收 `SourceKey`；或 (b) 在 `internal/httpapi` 加測試，用會記錄 sourceID 的假 limiter，斷言同一 /64 的兩個 `RemoteAddr`、以及 `::ffff:a.b.c.d` 與 `a.b.c.d`，傳進 limiter 的值相同。
- 狀態：v0.1 開放（必改）。**v0.1.1：已修正（commit `fc9a7a7`，註解三處齊備、內容正確，見 §8），待 #7 合併後關閉**。修補建議 2（`SourceKey` 型別或記錄 sourceID 的假 limiter 測試）未做，降為建議項追蹤，不阻擋

### [SCR-018] 取鎖前讀時間；zone 位址未正規化（延續 11 號檔）

- 嚴重度：**資訊**　不阻擋
- 現況：`guard.go` `allow` 仍在 `g.mu.Lock()` 前讀 `now`；`parseIP` 仍不接受 `%zone`。與 `ef9fcb8` 相同，沒有惡化。
- 處置：維持追蹤，不要求本 PR 修。

### 程序觀察（非安全發現）

- `mergeable_state: blocked` 來自 ruleset 24629854 的「1 個核准 review」，不是 CI 問題；合併仍須依 D-08／D-09 留言與核准程序。
- 必改項修好後會產生新 head；依 strict=true，舊 run 39195448908 的綠燈不能用於新 head，必須在新 head 上重新取得四項綠燈。

## 6. 判定

- **v0.1.1 判定（head `fc9a7a7`）：通過。**
  - v0.1 唯一必改項（SCR-017 前置條件註解）已於 `fc9a7a7` 完成；該 commit 只加 12 行 // 註解，可執行程式、公開簽章、workflows 都沒有變。
  - 新 head 在 run 39196278197 取得四項 success，四個 check run 的 head_sha 都是 `fc9a7a7`；merge ref tree 與 head tree 相同。
  - 無嚴重、高、中發現項未修正；SCR-018（資訊）與 SCR-017 建議 2 只做追蹤。
- **合併條件：**
  1. 合併時 PR head 必須仍是 `fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc`。head 若有任何變動（新 commit、rebase、force-push），本通過即失效，須再送安全部核對。
  2. 若 main 在合併前再前進（離開 `3d3a9cd`），依 strict=true（ruleset 24629805），#7 必須重新與 main 同步，並在同步後的新 head 上重新取得四項綠燈（run 的 head_sha 要等於新 head）。同步只能帶進 main 的改動（merge commit 經 remerge-diff／merge-tree 檢查沒有額外改動，或 rebase 後相對 main 的 patch 與 `fc9a7a7` 相同），並把差異送安全部核對。
  3. 安全部通過不是唯一條件：審查部關卡與 ruleset 24629854 要求的 1 個核准 review 仍是合併前的其他條件（目前 `mergeable_state: blocked` 即因缺核准 review），另依 D-08／D-09 留言程序合併。
- v0.1 判定（head `99b0879`）：退回，必改 1 項（SCR-017 註解）。保留作為紀錄。

## 7. 合併後的狀態處置（供 11／05／README 待改參考）

- SCR-006／007／008：#7 合併至 main 後關閉（修正於 `99b0879`／`fc9a7a7` 仍完整）。
- SCR-017：已修正（`fc9a7a7`），#7 合併至 main 後關閉；建議 2 另行追蹤。
- SCR-002：#7 合併時關閉。
- SCR-011（11 號檔）、SCR-012（12 號檔）：#7 依順序最後合併且 main 四綠後關閉。
- RR-007：合併前維持**中**；#7 合併至 main 後降為**低**。

## 8. v0.1.1 複核（head fc9a7a7）

查證時間：2026-10-19 約 11:50～11:58（UTC+8）。方法與 §1 相同（唯讀；cursor-github 讀取工具、公開 REST API、本機唯讀 clone `/workspace/tmp/su-pr7`）。沒有在 GitHub 留言、review、approve、合併或 push，也沒有執行 PR 內程式碼。

### 8.1 SHA、run 與檢查

| 項目 | 宣稱 | 查證結果 |
|---|---|---|
| PR head | `fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc` | ✅ 相符（state open、未合併、draft false、commits **9**、changed_files 13、+947／−39） |
| 與 `99b0879` 的關係 | fast-forward 1 個 commit | ✅ `fc9a7a7` 唯一 parent＝`99b087951c2007595138e5ccbe6e1da11ed76634`；`git rev-list --count 99b0879..fc9a7a7`＝1；`99b0879` 是祖先（FF）。commit 時間 11:46:23（UTC+8），訊息 `docs: clarify sourceID must come from ClientSource (no normalization in guard)` |
| base | main `3d3a9cd…913b` | ✅ PR base sha＝`3d3a9cd`；`origin/main` 仍為 `3d3a9cd`（main 未前進） |
| mergeable | — | ✅ `mergeable: true`、`rebaseable: true`；`mergeable_state: blocked`（ruleset 24629854 需 1 個核准 review） |
| run 39196278197 | head_sha＝新 head，success | ✅ workflow `ci-devsecops`，event `pull_request`，`head_sha`＝`fc9a7a7…`，`pull_requests[0]`＝#7（head `fc9a7a7`／base `3d3a9cd`），run_attempt 1，completed／**success**；建立 11:46:30、完成 11:47:51（UTC+8） |
| merge ref | — | `5c0215f6d43df887cc1ab48a692bfca7355a34f4`＝Merge `fc9a7a7` into `3d3a9cd`；四個 job 都 checkout `5c0215f`；merge ref tree `041749b…`＝head tree `041749b…`（完全相同） |
| 四項必要檢查 | success | ✅ `list_check_runs_for_ref(fc9a7a7)` 共 4 筆，head_sha 都是 `fc9a7a7`：`secrets-gitleaks`（job 117879630435）、`sast-gosec`（117879630466）、`sca-govulncheck`（117879630452）、`unit-test`（117879630185），皆 completed／success |
| ruleset 24629805 | strict=true | ✅ 公開 API 重讀：`strict_required_status_checks_policy: true`，`updated_at` 未變（2026-10-07 14:33:28 UTC+8） |
| sast-gosec log | h1 OK、Issues 0 | ✅ Go 1.27.2；`mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`、**`gosec h1 OK`**；Files 19、Nosec 0、**Issues 0** |
| sca-govulncheck log | — | ✅ Go 1.27.2；**`No vulnerabilities found.`**；`dependency graph compare HTTP 200`；dependency-review 無 critical 以上漏洞套件、無 denied 套件 |
| secrets-gitleaks log | — | ✅ gitleaks 8.21.2；42 commits scanned、**no leaks found** |
| unit-test log | — | ✅ Go 1.27.2；go vet 通過，go test 全部套件 ok |

### 8.2 diff 逐行檢查（99b0879..fc9a7a7）

- 範圍：`internal/httpapi/adapter.go` +4／−0、`internal/ratelimit/guard.go` +8／−0，共 +12／−0，與宣稱相符。
- **全是註解**：`git diff -U0` 中所有新增行都以 `//` 開頭，沒有刪除行；把兩個檔案去除 `//` 行後，`99b0879` 與 `fc9a7a7` 的內容完全相同；`git diff --check` 無空白問題。沒有可執行程式、型別或簽章變動。✅
- **三處齊備**：
  - `RateLimitChecker` 介面（`adapter.go:22-26`）：「Implementations receive sourceID exactly as produced by ratelimit.ClientSource and must not be handed raw RemoteAddr or header values; the guard does not normalize addresses itself.」
  - `Guard.CheckCreate`（`guard.go:92-96`）與 `Guard.CheckRedirect`（`guard.go:104-109`）：「sourceID must be the output of ClientSource. The guard uses it verbatim as the bucket key and does not normalize addresses (no IPv4-mapped folding or IPv6 /64 aggregation); passing a raw address can split one client's quota.」
  - 三處都明確寫出「sourceID 必須是 `ClientSource` 的輸出」與「guard 不做位址正規化」，符合 SCR-017 修補建議 1。註解緊接在宣告之前，是 godoc 文件註解。✅
- **說法正確性**：
  - 「guard 不做 IPv4-mapped 折疊或 IPv6 /64 聚合」：正確，這兩項只在 `source.go` `bucketKey` 處理，`guard.go` `allow` 不解析位址。
  - 「傳原始位址會拆散同一用戶端的配額」：正確。例如帶埠號的 `RemoteAddr` 每個連線埠都會變成新 key，同一 /64 的不同 IPv6 位址也會各自成為 key，這正是 SCR-007 的失效模式。
  - 「不得傳原始 RemoteAddr 或標頭值」：正確，也涵蓋偽造 XFF 的風險。
  - 小落差（資訊級，不需改）：「uses it verbatim as the bucket key」沒提到空字串例外：`allow` 會把 `""` 對應到共用的 `emptySourceKey`。這是既有、已審過的行為，不影響安全結論；日後若再修改註解可順帶補上。
- 結論：SCR-017 必改項已正確完成。

## 9. 合併後註記（v0.1.2）

| 項目 | 查證結果 |
|---|---|
| 合併 | PR #7 state closed、merged true，merged_at 2026-10-19T04:11:10Z＝**12:11:10（UTC+8）**；merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36` |
| 條件 1：head 仍為 `fc9a7a7` | ✅ 合併時 PR head＝`fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc`；merge commit 的第二個 parent 也是 `fc9a7a7`，head 未再變動 |
| 條件 2：main 未前進／strict=true | ✅ merge commit 第一個 parent＝`3d3a9cd`，與本檔複核時的 base 相同，不需重新同步 |
| 條件 3：審查部關卡與核准 review | 合併已完成，表示 ruleset 要求已滿足；審查部關卡紀錄不在安全部查證範圍 |
| 合併後 main CI | main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210)（head_sha `24b8da0`）四項必要檢查 `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` 皆 success（12:11:12～12:12:50 UTC+8）；`sast-gosec` log 有 `gosec h1 OK`、Issues 0 |
| 後續 | SCR-006／007／008／017、SCR-002 已關閉，SCR-011 已關閉（11 號檔 v0.1.5）；SCR-012 已關閉（12 號檔 v0.1.5）；RR-007 降為低（05 v0.1.10） |

## 修訂紀錄

| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-19 | 安全部 | 初版：PR #7 最終 head `99b0879` 合併前安全確認。查證 base `3d3a9cd`、mergeable true（mergeable_state blocked 為 review 規則）、run 39195448908 四項 success 且 head_sha 相符、merge ref tree＝head tree、gosec h1 OK／Issues 0、govulncheck 無漏洞、gitleaks 無洩漏；`b08417f`／`a5442a2` 以 remerge-diff 與 merge-tree 證實無夾帶；`92545cf`／`99b0879` 只改測試；SCR-006／007／008 修正仍在；#8／#11 簽章衝突已解決；可信代理 nil 時不讀 XFF；SCR-017 註解未補→**退回**，必改 1 項 |
| v0.1.1 | 2026-10-19 | 安全部 | 複核新 head `fc9a7a7`（自 `99b0879` FF 1 個 commit，只加 12 行 // 註解）：base 仍 `3d3a9cd`、mergeable true（blocked 為 review 規則）、run 39196278197 四項 success 且 head_sha 相符、merge ref tree＝head tree、gosec h1 OK／Issues 0、govulncheck 無漏洞、gitleaks 無洩漏；`RateLimitChecker`／`CheckCreate`／`CheckRedirect` 三處前置條件註解齊備且正確；SCR-017 改「已修正，待 #7 合併後關閉」；判定由退回改為**通過**，附合併條件（head 須仍為 `fc9a7a7`；main 前進須重新同步、只帶進 main 改動並重新四綠；審查部關卡與 1 個核准 review 仍為其他條件）；新增 §8；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-19 | 安全部 | 新增 §9 合併後註記：#7 已依條件以 head `fc9a7a7` 合併（2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`（parents：main `3d3a9cd`＋PR head `fc9a7a7f89f6e1cf585615b28bd4c22c560e44cc`））；main push run 39198199210 四項 success、gosec h1 OK；後續關閉／降級見 11 v0.1.5、12 v0.1.5、05 v0.1.10；狀態列同步；小版修訂不改檔名（檔名保留主版號） |
