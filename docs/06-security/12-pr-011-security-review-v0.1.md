---
文件：PR #11（ENG-012 觀測事件）安全程式碼審查
版本：v0.1.5
狀態：核准（附 D-10 條件：#8 合併後、自身合併前在最新 main 上取得新四綠）；**v0.1.5：#11 已合併；SCR-012 已關閉**
負責角色：安全部
最後更新：2026-10-19
對應需求：SEC-005、SEC-006、SEC-013；D-10、D-11；NFR-005、NFR-006；TM-007、TM-009、TM-014；隱私合規 `01-requirements/04-privacy-compliance-v0.1.md` §2、§4
專案代號：SHORTURL
---

# PR #11 安全程式碼審查（v0.1）

## 安全審查摘要

| 項目 | 內容 |
|---|---|
| PR | [#11](https://github.com/444112029012/short-url/pull/11)　feat: record structured create redirect and stats events (ENG-012)（open） |
| 完整 head | `5e6de8b440b53903ed10308f5c1d64558f5603c9`（分支基準＝main `161a4fd`） |
| CI run | 37580886113（pull_request；2026-10-07 14:19:35 UTC+8；merge ref `3fd3116`＝head ＋ main `161a4fd`） |
| 變更 | 3 檔：`internal/observability/hooks.go`（+93／−15）、`hooks_test.go`、`internal/httpapi/observability_test.go`；**沒有改** adapter、router、go.mod |
| 風險等級 | 中（日誌輸出、隱私邊界） |
| 發現項 | 中 1（SCR-012，處置已定：D-10＋合併順序，由 #7 更新分支解決；**v0.1.5 已關閉**）、低 2、資訊 2 |
| **結論** | **核准**。條件：依 D-10 與已定合併順序（#10 → #12 → #8 → #11 → #7），#11 須於 #8 合併後、自身合併前更新分支，在最新 main 上取得新的四綠 run；#11 先於 #7，與 #7 的測試簽章衝突由 #7 更新分支時解決（SCR-012）。SCR-013／SCR-015 列為維運 G5 前待辦 |

- 審查方式：唯讀。工具：cursor-github connector（`list_workflow_run_jobs`、`get_job_logs`、`list_pull_requests`）、本機唯讀 clone（`git fetch pull/11/head`）執行 `go test ./...`，並模擬 main＋PR #7＋PR #11 的組合。**沒有**在 PR 上留言或 review，沒有 push。
- SCR 編號接續 11 號檔（最大號 SCR-011），本檔使用 SCR-012～SCR-016。

## 研發聲稱逐項驗證

| # | 聲稱 | 驗證 | 證據 |
|---|---|---|---|
| 1 | 成功紀錄只寫 event＋short_code；失敗紀錄只寫 event＋error_type | **成立**。每筆 JSON 只有 slog 固定欄位 `time`／`level`／`msg`，加上 `event`，再加一個 attr（`short_code` 或 `error_type`） | `hooks.go:41-68`、`:79-82`；`TestFailureLogsOnlyErrorType` 用欄位白名單驗證失敗紀錄 |
| 2 | 不記長網址、IP、標頭、本文 | **成立（結構保證）**。hook 方法只接受一個 string 參數，拿不到 request、URL、err；adapter 呼叫點只傳 `created.ShortCode`／`stats.ShortCode`／路由參數（成功時）或 `string(appErr.Type)` | `adapter.go:79-139`（12 個呼叫點）；`EventHooks` 介面未變 |
| 3 | 錯誤與 panic 不往上丟 | **成立**。方法沒有回傳值；`emit` 在 `LogAttrs` 前 `defer recover()`；`slog` 的 handler 錯誤本來就不回傳給呼叫端 | `hooks.go:70-83`；`TestHooksDoNotPanic`（nil logger、panicking handler） |
| 4 | 沒開 metrics 端點 | **成立**。計數器是未匯出的 `atomic.Uint64` 欄位，沒有 getter；`router.go` 未變，沒有新增路由 | diff 只有 3 檔 |
| 5 | 有確認紀錄不含敏感值的測試 | **成立**。`TestUT_OBS_01_LogsOmitSensitiveValues` 對 create／redirect／stats／bad create 帶入長網址、Authorization、Cookie、RemoteAddr、XFF，再對輸出做字串黑名單比對 | `observability_test.go:112-178`；涵蓋缺口見 SCR-016 |
| 6 | slog 內部用 `context.Background()`，hook 簽名不變 | **成立** | `hooks.go:79`；`adapter.go:28-33` 介面未變（`git diff 161a4fd..5e6de8b` 不含 adapter.go） |

## 重點審查

| 項目 | 判定 | 說明 |
|---|---|---|
| SEC-005／`err.Error()` 帶出長網址或 SQL | **不會** | hook 拿不到 `error`。失敗只傳 `string(appErr.Type)`；`AppError.Error()` 也只回傳 type（`domain/errors.go`），`cause`（可能含 SQL／驅動訊息）不序列化、不記錄。repo 中唯一的 `err.Error()` 日誌是 `cmd/shorturl/main.go:24` 啟動失敗訊息（設定／開檔錯誤），不在本 PR 範圍，也不經過請求路徑 |
| error_type 是否封閉列舉 | **實際封閉；hook 端未強制** | 所有 `AppError` 只由 `domain` 的 5 個建構函式產生（`grep "AppError{"` 只命中 `errors.go:46-62`），Type 只會是 `invalid_url`／`not_found`／`rate_limited`／`internal_error`。但 `OnXxxFailure(string)` 不驗證輸入；ErrorMapper 對 HTTP 有正規化，hook 端沒有（SCR-014） |
| 日誌注入：short_code 是否先通過 `^[A-Za-z0-9]{8}$` | **是** | create：`service.go` 產生後 `ValidateShortCode`，不合格回 internal；redirect：只有 `Redirect` 成功才記錄，而 `Redirect` 第一步就是 `ValidateShortCode`（`redirect.go:26`），所以記錄的路由參數必定已驗證（建議改傳已驗證值，見 SCR-016）；stats：`stats.ShortCode` 是驗證後的值。失敗路徑**不寫**任何使用者輸入，只寫 error_type。JSONHandler 也會跳脫控制字元，屬第二道防線 |
| recover 是否造成 fail-open | **不會** | `recover` 只包住 `LogAttrs`；hook 呼叫一律在 adapter 已決定回應之後，沒有回傳值，不參與任何驗證或授權決策；計數在 recover 之前完成。驗證失敗仍回 400，不會「照樣繼續」 |
| 日誌量灌爆磁碟 | **有風險（低）** | 每個請求（含 429、404、400）都寫一行 INFO（約 90～110 B）；被限流的請求仍記錄，所以限流不會限制日誌量。`LOG_LEVEL` 在 config 有載入，但 hook 用 `slog.NewJSONHandler(os.Stderr, nil)`，**沒有接線**，無法調降。輸出是 stderr，輪替與上限取決於部署（SCR-013） |
| 測試是否涵蓋所有輸出路徑 | **大致涵蓋** | 6 個 hook 方法都有單元測試；adapter 層涵蓋 create 成功／invalid_url／rate_limited、redirect 成功／not_found／invalid_url／rate_limited、stats 成功／not_found／invalid_url。缺 adapter 層 `internal_error` 路徑（如關閉 store）、成功紀錄的欄位白名單（SCR-016） |
| 新增依賴 | **無** | `go.mod`／`go.sum` 未變；只用標準庫 `log/slog`、`sync/atomic` |

## short_code 寫進日誌是否可接受

**裁定：可接受（附條件）。** 初步立場「可接受」維持，但第 3 點前提**只部分成立**，所以補上條件。

| 前提 | 判定 | 依據 |
|---|---|---|
| SEC／隱私文件未把短碼定義為機敏 | **成立** | 隱私合規 §2 資料項目表：「短碼｜是否個資：**否**」；SEC-005 只禁止 Authorization／Cookie 原文；ENG-012 任務描述（`08-engineering/01-task-breakdown-v0.1.md`）明定「僅穩定錯誤碼與 short_code」 |
| 日誌保存期限與存取控制由維運負責 | **部分成立（待補）** | 隱私合規 §2／§4 只對「基礎設施日誌含 IP／UA」定了 ≤14 日與「僅供維運除錯／資安」；**應用觀測日誌（DS2，C4 `ObsLog` 標「待維運」）的保存期與存取控制，目前沒有維運文件明定** |
| 短碼與長網址只能在 DB 內對應 | **部分成立** | 日誌本身不含長網址，所以日誌**不**直接形成對應（成立）。但公開導向端點會為任何持有短碼的人解析長網址（產品語意），所以「日誌裡全部 `create_success` 的短碼清單」加上公開端點，就等於可列舉全部長網址。這會讓日誌讀取者繞過 TM-007／SEC-006 對「枚舉短碼得知他人長網址」的緩解（CSPRNG 只防猜，不防有清單的人） |

- 結論：**單一短碼不是個資也不是憑證**，寫入日誌可以接受；但**日誌中的短碼集合，敏感度應視同 DB 的對應索引**。
- 條件（維運，G5 前；不阻擋 PR #11 合併）：
  1. 觀測日誌的存取權限與 DB 同級（僅維運／使用者）；
  2. 保存 ≤ 14 日並自動輪替（對齊隱私合規日誌預設）；
  3. 不把未過濾的觀測日誌放到公開或第三方位置（公開 issue、CI artifact、外部日誌服務的公開儀表板）。
- 追蹤：SCR-015。

---

## 發現項

### [SCR-012] PR #11 與 PR #7 語意衝突（測試簽章），後合併者會讓 main 變紅
- 嚴重度：**中**　**阻擋「後於 PR #7 合併」的那一方**（流程／CI 完整性；與 11 號檔 SCR-011 同類）
- 位置：`internal/httpapi/observability_test.go:53、123`（`ratelimit.NewGuard()`）、`:198`（5 參數 `NewHTTPAPIAdapter`）
- 證據：本機唯讀模擬 main＋PR #7＋PR #11，文字合併無衝突，但 `go vet` 出現 `not enough arguments in call to ratelimit.NewGuard`，httpapi 測試 build failed。Rulesets `strict_required_status_checks_policy=false`，舊綠燈不會失效（v0.1 當時狀態；D-10 後自 2026-10-07 14:33:28（UTC+8）起為 strict=true，見 `ruleset-strict-true-d10-2026-10-07.md`；維運快照 v0.3.1 亦將加註；另見處置第 1 項）。與 PR #8 組合沒有簽章衝突（`routerWith` 用 `NewRouter`）。
- 處置（2026-10-07 定案，與 11 號檔 SCR-011 同一套）：
  1. **D-10（王昱殊，2026-10-07，使用者回覆 t83s391）**：規則一 `main-required-checks` 改 `strict=true`，合併前分支須與 main 同步，舊綠燈不得用於合併。strict=true 的設定取證由維運更新 `docs/09-cicd-env/evidence/` 後再引用。**（v0.1.2 更新）已生效**：維運 14:34:10（UTC+8）唯讀讀取 ruleset 24629805，`strict_required_status_checks_policy=true`、四項 check 名稱不變、`updated_at` 14:33:28（UTC+8）；取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，原始 JSON `docs/09-cicd-env/evidence/rulesets-2026-10-07/ruleset-24629805-d10-public-read.json`、`ruleset-24629805-d10-auth-read.json`（隨 PR #13 提交（head `6674a59`；已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`，取證檔連結已生效））。安全部約 14:36（UTC+8）以公開 REST API 唯讀讀取，結果一致。bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）（取證檔 §3 第 6 項，依 C 列）。
  2. **已定合併順序：#10 → #12 → #8 → #11 → #7。** #11 排在 #8 之後、#7 之前：
     - #11 本身：在 #8 合併後、自己合併前，更新分支並在最新 main 上重新取得四綠 run，D-08／D-09 留言附新 run 後才合併。#11 與 #8 沒有簽章衝突（`routerWith` 用 `NewRouter`）。
     - #11 與 #7 的衝突：由 #7 在 #8、#11 合併後更新分支時修正 `observability_test.go:53、123、198` 的舊簽章解決（見 11 號檔 SCR-011 處置）。
- 狀態：**已關閉（v0.1.5）**：#11 於 2026-10-19 11:34:13（UTC+8）合併（merge commit `3d3a9cddca0c99ee263905d7b6ef847f7f6b913b`）；#7 以 commit `99b0879` 修正 `observability_test.go` 舊簽章，最終 head `fc9a7a7` 經 16 號檔 v0.1.1 通過，於 2026-10-19 12:11:10（UTC+8）合併（merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`）；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success，main 未變紅。歷史：**處置已定（D-10）；strict=true 已生效（取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補））；待依順序（#10 → #12 → #8 → #11 → #7）合併後關閉**（未關閉；#7 依條件合併且 main 四綠後，由安全部確認關閉）

### [SCR-013] 觀測日誌量無上限：每請求一行（含被限流者），LOG_LEVEL 未接線
- 嚴重度：低　不阻擋（維運 G5 前條件）
- 位置：`internal/observability/hooks.go:30-32`（`slog.NewJSONHandler(os.Stderr, nil)`）；`internal/config/config.go:36` 載入 `LOG_LEVEL`，但沒有任何程式使用
- 影響：攻擊者可以用便宜的 429／404 請求產生日誌行，速度只受網路與 CPU 限制，限流也擋不住日誌量。日誌寫到 stderr，若部署（例如 Docker 預設 json-file 驅動）沒有設大小上限，磁碟可能被灌滿。另外，stderr 若接到阻塞的 pipe，寫入可能拖慢請求 goroutine。
- 修補建議：(1) 維運：容器／服務層設定日誌輪替（大小上限＋保存 ≤ 14 日），G5 部署檢查驗證；(2) 研發（可另開小 PR）：把 `LOG_LEVEL` 接到 `HandlerOptions.Level`；`rate_limited` 失敗改為只計數或取樣記錄。
- 狀態：開放（追蹤）

### [SCR-014] hook 端未強制 error_type 封閉集合
- 嚴重度：資訊　不阻擋
- 位置：`hooks.go:46-68`（`OnXxxFailure(errorType string)` 原樣寫入）
- 影響：目前所有呼叫端都傳 domain 常數，所以實際上是封閉的。但將來若有呼叫端傳入自由字串（例如 `err.Error()`），hook 會照寫，失去縱深防禦；ErrorMapper 對 HTTP 有白名單，hook 沒有。
- 修補建議：`emit` 前以白名單正規化（不在四值內→`internal_error`），與 `mapper.go` `statusFor` 對齊；補一條測試。
- 狀態：建議

### [SCR-015] 觀測日誌中的短碼集合等同對應索引；保存期與存取控制尚無維運文件
- 嚴重度：低　不阻擋（維運 G5 前條件）
- 位置：`hooks.go:42、52、62`（`short_code`）；維運文件缺漏（DS2／`ObsLog` 待維運）
- 影響：見上節。單一短碼不是機敏資料，但日誌集合加上公開導向端點，可列舉全部長網址（TM-007 緩解對日誌讀取者失效）。
- 修補建議：維運在環境／部署文件明定觀測日誌的存取權限（同 DB）、保存 ≤ 14 日並輪替、不得公開或外送未過濾日誌。
- 狀態：開放（追蹤）

### [SCR-016] 測試涵蓋缺口與小型強化
- 嚴重度：資訊　不阻擋
- 內容：(1) adapter 層沒有 `internal_error` 路徑的 hook 測試（例如關閉 store 後 create／redirect／stats），敏感值測試也沒有涵蓋含 SQL 錯誤的路徑；結構上 hook 拿不到 err，風險低，但應有回歸測試鎖住。(2) 成功紀錄沒有欄位白名單斷言（只有失敗紀錄有）。(3) `adapter.go:120` 成功時記錄的是路由參數 `chi.URLParam`；值已經過 `Redirect` 驗證，等價但不直觀，建議改由 `RedirectTarget` 回傳已驗證短碼後再記錄，避免將來重構時誤記未驗證輸入。
- 狀態：建議

---

## CI 對應確認

| 項目 | run 37580886113 |
|---|---|
| run head_sha＝PR head | ✅ `5e6de8b…`（event `pull_request`） |
| 測試的 merge ref | `3fd3116`＝head ＋ main `161a4fd`（含 PR #6 fail-closed workflow） |
| 四 required checks | secrets-gitleaks ✅、sast-gosec ✅、sca-govulncheck ✅、unit-test ✅ |
| Probe／dependency review | `dependency graph compare HTTP 200`（BASE 161a4fd／HEAD 5e6de8b）；`Dependency review (PR only)` 實跑，`fail-on-severity: critical`，無漏洞、無 denied |
| govulncheck | No vulnerabilities found |
| GITHUB_TOKEN 權限 | Contents／Metadata／PullRequests: read |
| SEC 驗收在 CI 的對應 | SEC-005／NFR-006：`internal/observability/hooks_test.go`（UT-OBS-01、欄位白名單、計數、不 panic）與 `internal/httpapi/observability_test.go`（adapter 呼叫順序、敏感值黑名單）都在 `unit-test` 的 `go test` 內 |
| 本機複核 | `go test ./...` 於 head 通過；main＋PR #7＋PR #11 組合 build failed（SCR-012） |

## 殘餘風險（v0.1 建議；v0.1.1 起已依 D-11 登記於 05）

- **RR-009（低；已由使用者於 D-11 接受，登記於 05）**：觀測日誌（DS2）含全部短碼，加上公開導向端點可列舉長網址；日誌量無應用層上限。補償：日誌存取同 DB、保存 ≤ 14 日並輪替、容器日誌大小上限（SCR-013、SCR-015）；G5 前驗證。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | 初版：PR #11 核准（附 SCR-012 合併順序條件）；研發 6 項聲稱均成立；short_code 入日誌可接受（附條件：日誌存取同 DB、≤14 日輪替、不公開）；SCR-012～016 |
| v0.1.1 | 2026-10-07 | 安全部 | 依 D-10／已定合併順序（#10 → #12 → #8 → #11 → #7）：SCR-012 處置改為 strict=true＋#11 於 #8 合併後在最新 main 重取四綠、與 #7 衝突由 #7 更新分支解決，狀態改「處置已定（D-10）；strict=true 取證待維運」；摘要與結論條件同步；RR-009 標註已依 D-11 接受；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-07 | 安全部 | SCR-012 狀態改為「處置已定（D-10）；strict=true 已生效（取證 `ruleset-strict-true-d10-2026-10-07.md`，隨維運 PR A 提交、待合併；bypass 清單待具權限讀取確認）；待依順序合併後關閉」（未關閉）；證據句補註 strict=false 為 v0.1 當時狀態；小版修訂不改檔名（檔名保留主版號） |
| v0.1.3 | 2026-10-19 | 安全部 | SCR-012 證據／處置／狀態：bypass 寫法由「待具權限讀取確認」改為「bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）」，引用 `ruleset-strict-true-d10-2026-10-07.md`；「維運 PR A、待合併」改為「隨 PR #13 提交（截至 2026-10-19 仍 open、未合併）」；證據句 strict=false 註記補 D-10 生效時點 14:33:28（UTC+8）與取證檔、維運快照 v0.3.1 亦將加註；歷史修訂列不改；小版修訂不改檔名（檔名保留主版號） |
| v0.1.4 | 2026-10-19 | 安全部 | SCR-012 證據／狀態：取證檔「隨 PR #13 提交（截至 2026-10-19 仍未合併／09:02 仍 open）」改為已合併事實（2026-10-19 09:05:27 UTC+8，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`，取證檔連結已生效；cursor-github 唯讀查證）；SCR-012 處置與關閉條件不變；小版修訂不改檔名（檔名保留主版號） |
| v0.1.5 | 2026-10-19 | 安全部 | SCR-012 改「已關閉」：#11 於 2026-10-19 11:34:13（UTC+8）合併（merge commit `3d3a9cddca0c99ee263905d7b6ef847f7f6b913b`）；#7 以 commit `99b0879` 修正 `observability_test.go` 舊簽章，最終 head `fc9a7a7` 經 16 號檔 v0.1.1 通過，於 2026-10-19 12:11:10（UTC+8）合併（merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`）；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success，main 未變紅；狀態列與摘要表發現項列同步；小版修訂不改檔名（檔名保留主版號） |
