---
文件：PR #2 安全程式碼審查
版本：v0.1.4
狀態：核准（CR-001（D-07）與合併路徑 A（D-08）已核准）
負責角色：安全部
最後更新：2026-10-07
對應需求：SEC-001～SEC-008、SEC-013、SEC-014、SEC-016、SEC-017；TM-001、TM-002、TM-006、TM-007、TM-009、TM-012、TM-014、TM-015、TM-019、TM-022
專案代號：SHORTURL
---

# PR #2 安全程式碼審查（v0.1.4）

## 安全審查摘要
- PR：#2　風險等級：**高風險**（外部輸入、重新導向、DB、CI 設定）
- PR head SHA：審查基準 `3061289d3717a29af0314d4e23aeddd06f9ec322`（短碼 `3061289`；複核前為 `ccd3588`）；現行 head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc`（空 commit 重跑 CI，tree 與 `3061289` 相同，無程式碼差異）
- 對應：ENG-002～010、013～015、017、018；REQ-001～012；SEC-001、003、004、006、009、010、012（本機非演示）、013、014、016、017
- CI 掃描（head 實測）：`secrets-gitleaks` ✅ / `sast-gosec` ✅ / `sca-govulncheck` ✅ / `unit-test` ✅（run [37573748568](https://github.com/444112029012/short-url/actions/runs/37573748568)；現行 head `2698138` run [37575622585](https://github.com/444112029012/short-url/actions/runs/37575622585) 亦四項全綠，且 dependency graph compare HTTP 200、dependency-review 實跑）
- 合併狀態：本審查於 2026-10-07 完成時 PR #2 尚未合併（`mergeable_state=blocked`）；合併依 D-08 路徑由使用者執行
- 後續：PR #2 於 2026-10-07 13:56:49 UTC+8 合併，merge commit `a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`（父 commit `d09cceb` 與 head `2698138`）
- 發現項：嚴重 0、高 0、中 1 開放（SCR-002）／1 已修正（SCR-001）、低 2、資訊 1
- 結論：**核准**（SCR-001 已於 `3061289` 關閉；無剩餘合併前資安阻擋項。CR-001（D-07）與合併路徑 A（D-08）已核准；平台合併依 D-08 補償控制由使用者執行）

## 審查範圍（實際讀取）

| 來源 | 路徑／項目 |
|---|---|
| GitHub connector | `get_pull_request`、`list_pull_request_files`、`get_pull_request_diff`、`list_check_runs_for_ref`、`get_file_contents`（ref=`ccd3588`） |
| 應用程式 | `internal/httpapi/adapter.go`、`router.go`、`mapper.go`、`router_test.go`；`internal/application/redirect.go`、`service.go`；`internal/domain/validator.go`、`generator.go`、`errors.go`；`internal/infra/sqlite/store.go`；`internal/ratelimit/guard.go`；`internal/observability/hooks.go`；`internal/config/config.go`；`cmd/shorturl/main.go` |
| 組態／CI | `go.mod`、`.github/workflows/ci-devsecops.yml`、`.env.example`、`.gitignore` |
| 專案文件 | `01-sec-requirements-v0.2.md`、`03-ci-security-gates-v0.2.md`、`04-threat-model`（對照 SEC／TM）、`09-cicd-env/05-pr2-ci-change-review-v0.1.md`、`01-ci-pipeline-v0.2.md`（§2.1／§3.3／§5.2） |

未逐行精讀全部測試輔助函式與 `go.sum` 細節；相依鎖定以 `go.mod`／`go.sum` 存在且 CI `test -f go.sum` 為準。未於 GitHub 發表 review／留言。

---

## 四項重點（符合／不符合＋證據）

### 1. 開放導向（SEC-014、SEC-002、TM-006）— **符合**
- `HandleRedirect` 僅將 `target.LongURL` 寫入 `Location`；註解明示忽略 query／header（`internal/httpapi/adapter.go` 約 L106–123）。
- `RedirectService.Redirect` 只自 DS1 載入後回傳存檔 `long_url`，請求參數非輸入（`internal/application/redirect.go` 約 L11–45）；`safeLocation` 拒絕含控制字元的值。
- 測試：`router_test.go` `TestCreateRedirectStatsHappyPath` 附加 `?url=`／`X-Redirect-Target` 後 `Location` 仍為建立時長網址。

### 2. 輸入驗證（SEC-001、SEC-013）— **符合**
- 長網址：非空、UTF-8、長度 ≤2048、無 CTL／空白、`url.Parse` 絕對 URL、scheme 僅 `http`／`https`（大小寫正規化）、須有 Hostname（`internal/domain/validator.go` 約 L22–41）。
- 短碼：`^[A-Za-z0-9]{8}$`（同檔 L44–49）。
- 本文：`maxCreateBody=8192`、`LimitReader`、`DisallowUnknownFields`、僅 `application/json`（`adapter.go` 約 L18、L151–174）。
- 測試覆蓋 `ftp://`、`javascript:`、空字串、非法短碼（`router_test.go`／validator 測試）。

### 3. 參數化查詢與原子 click_count（SEC-016）— **符合**
- 所有 SQL 為固定字串＋`?` 綁定（`internal/infra/sqlite/store.go`：`insertSQL`／`selectByCodeSQL`／`existsSQL`／`incrementSQL`／`countSQL`）。
- `IncrementAtomic`：單句 `UPDATE … SET click_count = click_count + 1 … RETURNING`（約 L50–54、L164–177）。無字串拼接 SQL。

### 4. CSPRNG（SEC-006、REQ-010）— **符合**
- 僅 `crypto/rand`（`internal/domain/generator.go`）；無 `math/rand`。
- 取模偏差：拒絕 `[248,255]` 後再 `%62`（`alphabetUnbiasedLimit`）。
- 碰撞：應用層最多 8 次（`Exists`／`ErrConflict` 重試，`service.go` `defaultGenerateAttempts=8`）。

---

## 其他清單項（摘要）

| 項目 | 結果 | 證據／說明 |
|---|---|---|
| 錯誤不洩漏（SEC-003） | 符合 | `ErrorMapper` 固定四類型＋訊息目錄；`InternalWrap` 之 cause 不序列化（`mapper.go`、`errors.go`） |
| 日誌／個資（SEC-004／005） | 符合（本批） | schema 無 `visitor_ip`／`user_agent`；obs hooks 為 no-op 且 TODO 禁止記 Auth／Cookie／IP |
| 限流（SEC-007） | **未實作** | `ratelimit.Guard` 一律放行；adapter 已接埠；屬 ENG-011，見 SCR-002 |
| 安全標頭（SEC-008） | **未實作** | `router.go` TODO(ENG-016)；見 SCR-003 |
| 機密／.gitignore | 符合 | `.env` 已 ignore；`.env.example` 無真實機密 |
| 相依 | 可接受 | `chi/v5`、`modernc.org/sqlite`；`go.sum` 已提交 |
| CI workflow 安全 | 部分不足 | 頂層 `permissions: contents: read`＋`pull-requests: read`；無 `pull_request_target`；無不可信輸入內插 `run:`；action 仍用 major tag（SCR-004）；審查時 dependency-review 為暫行 fail-open 且已發 `::warning::`（SCR-001 已關；其後由獨立 PR 恢復 fail-closed，見 [09](./09-cr-ci-deviations-ruling-v0.1.md) CR-001 第 3 項） |

---

## 發現項

### [SCR-001] dependency-review 略過未發 Actions 可見警告
- 嚴重度：中
- 位置：`.github/workflows/ci-devsecops.yml`（Probe dependency graph 步驟；commit `ccd3588`）
- 類別：IaC／CI；對應 SEC-010；CR 暫行條件
- 影響：依賴圖未開時靜默略過輔助 SCA，維運／審查不易察覺控制未生效
- 重現：1. 於 dependency graph 未啟用倉庫對 PR 跑 `sca-govulncheck` 2. 檢視 job log：僅一般 echo、無 `::warning::` 註解 3. dependency-review 步驟被跳過且 job 仍 success
- 修補建議：在 `enabled=false` 分支加入 `echo "::warning::Dependency graph unavailable; dependency-review skipped; govulncheck remains blocking."`（或等價）；開啟 graph 後改 fail-closed（見 CR 裁示）
- 狀態：**已修正**（commit `3061289`；run 37573748568 之 `sca-govulncheck` Annotations 可見 warning：`dependency-review skipped: Dependency graph disabled; fail-open until enabled; expires before G5`）。複核 2026-10-07。

### [SCR-002] 速率限制為 always-allow stub（SEC-007 未滿足）
- 嚴重度：中
- 位置：`internal/ratelimit/guard.go`（commit `ccd3588`）
- 類別：大量請求；對應 SEC-007、TM-003／008／011
- 影響：建立／導向端點尚無來源級限流，濫用與枚舉風險殘餘升高（已知產品範圍延後）
- 重現：1. 單元／整合路徑確認 `CheckCreate`／`CheckRedirect` 恒回 nil 2. 對照 PR 正文標 ENG-011 延後
- 修補建議：ENG-011 實作門檻與可信代理邏輯；完成前殘餘風險持續追蹤（條件式殘餘見殘餘風險表 §3）
- 狀態：開放（追蹤 ENG-011；**不擋本批**——ENG 範圍未含 011）

### [SCR-003] 缺少 SEC-008 安全回應標頭
- 嚴重度：低
- 位置：`internal/httpapi/router.go` TODO(ENG-016)（commit `ccd3588`）
- 類別：輸出編碼／安全標頭；對應 SEC-008、TM-016
- 影響：缺 `X-Content-Type-Options` 等縱深防禦標頭
- 重現：1. 對建立／導向／錯誤回應擷取標頭 2. 確認尚無 nosniff／HSTS（本機 HTTP 本就不應宣稱 HSTS）
- 修補建議：ENG-016 依 SEC-008 子集實作
- 狀態：開放（追蹤 ENG-016；**不擋本批**）

### [SCR-004] 第三方 GitHub Actions 未釘選完整 commit SHA
- 嚴重度：低
- 位置：`.github/workflows/ci-devsecops.yml`（`actions/checkout@v4`、`actions/setup-go@v5`、`actions/dependency-review-action@v4`；commit `ccd3588`）
- 類別：IaC／CI；對應門檻文件 §10
- 影響：major tag 浮動可能引入非預期 action 變更
- 重現：1. 檢視 workflow `uses:` 行 2. 確認為 `@v4`／`@v5` 而非 40 字元 SHA
- 修補建議：穩定後釘完整 SHA；Dependabot `github-actions` 更新
- 狀態：開放（追蹤；不擋本批；workflow 註解已自承待改）

### [SCR-005] 短碼 CSPRNG 取模偏差已處理（記錄用）
- 嚴重度：資訊
- 位置：`internal/domain/generator.go`（commit `ccd3588`）
- 類別：加密；對應 SEC-006
- 影響：無負面影響；採 rejection sampling 避免 base62 偏差
- 重現：程式碼審查確認 `alphabetUnbiasedLimit=248`
- 修補建議：無
- 狀態：已接受（資訊）

---

## 結論

**核准**（資安合併前阻擋項已清）。

| 類型 | 項目 |
|---|---|
| 合併前資安阻擋 | **無**（SCR-001 已關） |
| 不擋本批（追蹤） | SCR-002 → ENG-011；SCR-003 → ENG-016；SCR-004 釘 SHA |
| 四項核心 | 開放導向／輸入驗證／參數化／CSPRNG **皆符合** |
| 平台合併 | CR-001（D-07）與合併路徑 A（D-08）已核准；RR-006 已由使用者接受。本審查於 2026-10-07 完成時 PR #2 尚未合併（後續：13:56:49 UTC+8 合併，merge commit `a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`）；合併時須依 D-08 於 PR 留言附審查部「可合併」＋綠燈連結，由使用者本人合併 |

無嚴重／高發現項；無機密進版控；四 required checks 全綠（含 Annotations 可見 dependency-review 略過警告）。高風險 PR 本審查由安全部執行。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | PR #2 高風險完整審；附條件核准（::warning::） |
| v0.1.1 | 2026-10-07 | 安全部 | 複核 head `3061289`／run 37573748568；SCR-001 關閉；結論改核准 |
| v0.1.2 | 2026-10-07 | 安全部 | 依 D-07／D-08：「合併待使用者拍板」改為 CR-001（D-07）與合併路徑 A（D-08）已核准；記錄 PR #2 現行 head `2698138`（同 tree）與 run 37575622585 全綠；PR #2 截至 13:43 UTC+8 未合併 |
| v0.1.3 | 2026-10-07 | 安全部 | 合併狀態改為不過時寫法（「本審查於 2026-10-07 完成時 PR #2 尚未合併」）；workflow 安全列標明 fail-open 為審查時狀態；小版修訂不改檔名（檔名保留主版號） |
| v0.1.4 | 2026-10-07 | 安全部 | 補後續：PR #2 於 2026-10-07 13:56:49 UTC+8 合併，merge commit `a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`；小版修訂不改檔名（檔名保留主版號） |
