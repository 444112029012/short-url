---
文件：PR #7（ENG-011 速率限制）／PR #8（ENG-016 安全標頭）安全程式碼審查
版本：v0.1.6
狀態：**v0.1.5：PR #7 已合併（最終 head `fc9a7a7` 經 16 號檔 v0.1.1 通過，2026-10-19 12:11:10 UTC+8 合併，merge commit `24b8da0`）；SCR-002／006／007／008／011／017 已關閉；v0.1.6：SCR-003 已關閉（#8 已合併，`e9f4396`）**。歷史：PR #7 複核（head `ef9fcb8`）邏輯上核可——SCR-006／007／008 修正已驗證；合併前新 head 差異須再送安全部核對並取得四項新綠燈；PR #8 核准（附 D-10 條件：合併前分支須與 main 同步並取得新四綠）
負責角色：安全部
最後更新：2026-10-19
對應需求：SEC-005、SEC-007、SEC-008；D-10、D-11；NFR-003、NFR-004；TM-003、TM-008、TM-011、TM-016、TM-017；前次 SCR-002、SCR-003（見 08）；RR-004、RR-007（見 05）
專案代號：SHORTURL
---

# PR #7／PR #8 安全程式碼審查（v0.1）

## 安全審查摘要

| 項目 | PR #7 | PR #8 |
|---|---|---|
| 主題 | ENG-011 每來源 create／redirect 速率限制＋可信代理來源識別 | ENG-016 SEC-008 安全回應標頭＋HSTS 開關 |
| 完整 head | v0.1／v0.1.1 審查：`04815e5746fed3c6133bcbc6c27669c3906efc3b`；**v0.1.2 複核：`ef9fcb8cd4999cc40dab990dda0d091a851c8615`**（見「PR #7 複核」） | `9dbc246238bd8623b16ca3dc795bb88976be6d81` |
| 分支基準（merge-base） | `a8246ce`（PR #2 合併） | `a8246ce`（PR #2 合併） |
| CI run | 37579834778（pull_request；merge ref `a59870b` = head ＋ main `161a4fd`）；**複核：37581955467**（head `ef9fcb8`，四項 success） | 37580036601（pull_request；merge ref `97ae999` = head ＋ main `161a4fd`） |
| 風險等級 | 高（外部輸入、可用性控制、來源識別） | 中（全回應路徑中介層） |
| 發現項 | v0.1：中 2（阻擋）、低 1、資訊 1；**v0.1.2：SCR-006／007／008 修正已驗證；新增 SCR-017 低、SCR-018 資訊（均不阻擋）** | 低 1 |
| 跨 PR | SCR-011 中：處置已定（D-10 strict=true＋合併順序 #10 → #12 → #8 → #11 → #7）；#7 排最後，由 #7 更新分支解決衝突 | 同左；#8 合併前須在最新 main 上取得新四綠 |
| **結論** | **v0.1.5：最終 head `fc9a7a7` 已經 16 號檔 v0.1.1 通過並合併（merge commit `24b8da0`，main 四綠）。** v0.1／v0.1.1（head `04815e5`）：要求修改（SCR-006、SCR-007）。**v0.1.2（head `ef9fcb8`）：邏輯上核可**；合併前新 head 差異（更新到最新 main、修 #8／#11 簽章衝突、補 SCR-017 前置條件註解）須再送安全部核對，並取得四項新綠燈 | **核准**（附條件：依 D-10，#8 於 #10、#12 合併後、自身合併前，分支須與 main 同步並取得新的四綠 run；依合併順序 #8 先於 #7，與 #7 的簽章衝突由 #7 更新分支時解決） |

- 審查方式：唯讀。工具：cursor-github connector（`list_pull_request_files`、`get_file_contents`、`get_workflow_run`、`list_workflow_run_jobs`、`get_job_logs`）、raw.githubusercontent.com、本機唯讀 clone（`git fetch pull/7/head`、`pull/8/head`）執行 `go vet`／`go test`（含 `-race`）。**沒有**在 PR 上留言或 review，沒有 push，也沒有撰寫攻擊程式；表滿與 IPv6 的判斷依程式碼閱讀和 PR 自帶測試 `TestGuardEvictsIdleSourcesAtCapacity`。
- 本文上限（DR-G4-01）：不在本檔重做，結論見 [10-dr-g4-01-body-limit-security-opinion-v0.1.md](./10-dr-g4-01-body-limit-security-opinion-v0.1.md)（附條件同意；C1＝實作 PR 須在 PR #7 合併後）。PR #7 程式碼事實（`adapter.go:83` 限流先於 `:87` 解碼、`:164-166` 截斷讀取）與 10 號檔一致。
- `MaxHeaderBytes`（10 號檔 C3）：屬請求端資源限制，不在 SEC-007（限流）或 SEC-008（回應標頭）範圍，ENG-011／016 都沒有承接，**不列為本檔發現項**，只引用 10 號檔 C3（G5 前待補）。

## 審查範圍（實際讀取）

| PR | 檔案 |
|---|---|
| #7 | `internal/ratelimit/guard.go`、`source.go`、`guard_test.go`、`source_test.go`；`internal/httpapi/adapter.go`、`ratelimit_test.go`、`router_test.go`；`cmd/shorturl/main.go`；`internal/config/config.go`；`.env.example` |
| #8 | `internal/httpapi/security_headers.go`、`security_headers_test.go`、`router.go`；`internal/config/config.go`、`config_test.go`；`cmd/shorturl/main.go`；`.env.example` |
| 對照 | `01-sec-requirements-v0.2.md` SEC-005／007／008；`09-cicd-env/02-environments-v0.1.md`（TRUSTED_PROXY_CIDRS、HTTPS、TLS 終止）；`09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`（`strict_required_status_checks_policy=false`）（快照當時為 strict=false；D-10 後自 2026-10-07 14:33:28（UTC+8）起為 strict=true，見 `ruleset-strict-true-d10-2026-10-07.md`；維運快照 v0.3.1 亦將加註） |

## 檢查表

| 項目 | PR #7 | PR #8 |
|---|---|---|
| SEC-007 門檻與行為 | 符合：create 30／redirect 120（env 可調，非正值回落預設，`NewGuard`）；超限回 429 `rate_limited` 固定訊息；`CheckRedirect`（`adapter.go:114`）先於 `Redirect`（`:118`），被限流時不導向、不增加 click_count（`TestAdapterRateLimit429AndRedirectDoesNotCount`）；被拒不扣額度；空來源共用單一 bucket；`-race` 通過 | — |
| 來源識別／防偽造 | 預設只用 peer IP；僅 peer 在 `TRUSTED_PROXY_CIDRS` 內才讀 XFF（取最右不可信）／X-Real-IP（`source.go:46-62`）；偽造 XFF 不能繞過（`TestSpoofedXFFDoesNotBypassCreateLimit`）。**缺口**：IPv6 未聚合（SCR-007）、XFF 只讀第一個標頭行（SCR-008）（v0.1.2：兩者於 `ef9fcb8` 修正，已驗證） | — |
| 有界記憶體 | 有上限 8192／每表（`guard.go:15`），但表滿處理方式造成全域拒絕（**SCR-006**；v0.1.2：`ef9fcb8` 改為 FIFO 淘汰最舊並放行，已驗證） | — |
| SEC-005 日誌 | 不記錄來源位址（`guard.go`、`source.go` 無 log 呼叫） | 無日誌 |
| SEC-008 (2) nosniff | — | 符合：所有回應（含 chi 404／405、429、500、302）`X-Content-Type-Options: nosniff`（`security_headers.go:25-37`；`TestSecurityHeadersOnSuccessRedirectAndErrors`、`TestSecurityHeadersOn429And500`） |
| SEC-008 (1) HSTS | — | 條件式：`HSTS_ENABLED` 且 `r.TLS != nil` 才送 `max-age=31536000`（`:32`）；見 §HSTS 裁定、SCR-010 |
| SEC-008 (3) CORS | — | 符合：不送 `Access-Control-Allow-Origin`（測試斷言） |
| SEC-008 (4) 版本標頭 | — | 符合：不送 `Server`（測試斷言）；Go `net/http` 預設不送 |
| SEC-008 3xx | — | 符合：302 帶全部基線標頭，`Location` 不被中介層修改，僅為庫中長網址 |
| 額外標頭 | — | `Referrer-Policy: no-referrer`、`X-Frame-Options: DENY`、`CSP: default-src 'none'; frame-ancestors 'none'`、JSON 回應 `Cache-Control: no-store`；皆為縱深防禦，無功能副作用（302 不加 Cache-Control） |
| 機密／相依 | `.env.example` 無機密；無新增 module（`go.mod` 未變） | 同左 |
| CI workflow | 未改 `.github/` | 未改 `.github/` |

---

## 發現項

### [SCR-006] 來源表滿時對所有新來源回 429（全域拒絕）
- 嚴重度：**中**　**阻擋 PR #7**
- 位置：`internal/ratelimit/guard.go:94-100`（`allow`：`len(buckets) >= maxSources` → `evictIdle` → 仍滿 → `RateLimited`）；create 與 redirect 各一張表，各 8192
- 類別：可用性（DoS）；SEC-007 控制本身引入的新失效模式；TM-003、TM-011
- 影響：一分鐘窗內只要有 ≥ 8192 個仍在窗內的不同來源，**所有尚未在表內的來源**（包含正常訪客）的 create 與 **redirect**（核心功能）都會收到 429，最長持續到最舊窗過期（≤ 60 秒），攻擊者持續送就能一直維持。觸發門檻：每分鐘 8192 個不同來源（約 137 個／秒）；IPv6 下單一主機即可達成（見 SCR-007）。PR 自帶測試 `TestGuardEvictsIdleSourcesAtCapacity` 也明確斷言「表滿時新來源 c 被拒」，可見這是設計行為，不是偶發。另外，每次插入新來源時表都滿，會觸發 O(n) 的 `evictIdle` 全表掃描，而且是在持有 mutex 時執行，攻擊下有 CPU 放大。
- 裁定：**「表滿 fail-closed（全域拒絕）」不接受**。限流是保護可用性的控制，表滿時拒絕所有新來源，等於讓控制本身成為最便宜的 DoS 開關，失敗方向錯誤。表滿時應**犧牲限流精度、保住可用性**（fail-soft），記憶體仍要維持有界。
- 修補（必要）：
  1. 表滿且無過期項時，**淘汰窗起點最早的 bucket**，再放入新來源，不回 429。固定窗下插入順序等於 `windowStart` 順序，可用 `container/list`（或 ring）維持 FIFO，達到 O(1) 淘汰，取代 O(n) `evictIdle` 全表掃描。被淘汰來源的計數歸零，這是可接受的精度損失：只有活躍來源超過容量時才發生，影響上限是「該來源在窗內多拿一輪額度」。
  2. 與 SCR-007 一起修：IPv6 依 /64 聚合，降低單一主機填表能力。
  3. 不記錄被淘汰或被拒的來源位址（SEC-005）；需要時只計數（可留待 ENG-012）。
  4. 測試：把 `TestGuardEvictsIdleSourcesAtCapacity` 的期待改為「表滿且全部在窗內時，新來源仍放行、表大小不超過上限」；另補「被淘汰的是最舊來源」與 redirect 表的同類測試。
- 狀態：**已關閉（v0.1.5：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8））**。歷史：已修正（commit `a5b3d2d`，含於 head `ef9fcb8`）；修正已驗證（v0.1.2 複核）

### [SCR-007] IPv6 來源未依前綴聚合，可繞過每來源限流並加速填表
- 嚴重度：**中**　**阻擋 PR #7**（與 SCR-006 同修）
- 位置：`internal/ratelimit/source.go:46-62`（`ClientSource` 回傳完整位址字串作為 bucket key）
- 類別：限流規避；TM-003、TM-017
- 影響：`HTTP_ADDR=:8080` 在 Go 下是雙棧監聽。取得一個 IPv6 /64 的用戶端可以逐請求換來源位址，每次都拿到新 bucket，所以每來源限流對 IPv6 用戶端**實質無效**，同時可以單機填滿來源表（觸發 SCR-006）。
- 修補（必要）：IPv6 來源 key 正規化為 /64 前綴（`ip.Mask(net.CIDRMask(64,128))`）；IPv4 維持 /32；IPv4-mapped IPv6 先轉 IPv4。可信代理路徑解析出的位址套用同一正規化。補測試：同一 /64 內不同位址共用額度、不同 /64 獨立。
- 狀態：**已關閉（v0.1.5：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8））**。歷史：已修正（commit `ef9fcb8`）；修正已驗證（v0.1.2 複核）。殘餘：/64 聚合下，持有 ≥ /51（如常見 /48 配置）的單一用戶端仍可造出 > 8192 個 /64 使 FIFO 循環，詳 RR-007

### [SCR-008] 受信代理路徑只讀第一行 X-Forwarded-For；X-Real-IP 退回依賴代理覆寫
- 嚴重度：低　不阻擋
- 位置：`internal/ratelimit/source.go:55`（`r.Header.Get("X-Forwarded-For")`）、`:58`（`X-Real-IP`）
- 類別：來源偽造；TM-017、RR-004
- 影響：只在 peer 屬 `TRUSTED_PROXY_CIDRS` 時才會發生（預設空，不受影響）。若代理以**另一行**標頭附加 XFF（而不是合併成同一行），`Header.Get` 只取第一行，也就是用戶端自帶的值，最右不可信邏輯會被繞過。若代理不設 XFF 也不覆寫 X-Real-IP，則會採用用戶端自帶的 X-Real-IP。
- 修補建議：改用 `r.Header.Values("X-Forwarded-For")`，依序串接後再從右側解析；維運文件註明代理必須覆寫 X-Real-IP（或停用 X-Real-IP 退回）。可在本 PR 修，或另開小 PR，G5 前完成。
- 狀態：**已關閉（v0.1.5：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8））**。歷史：已修正（commit `ef9fcb8`）；修正已驗證（v0.1.2 複核）。程式註解（`source.go:44-46`）已註明代理必須覆寫 `X-Real-IP`；維運部署文件的同等註記仍待補（RR-004，G5 前）

### [SCR-009] 部署於 Tunnel／反向代理後若未設 TRUSTED_PROXY_CIDRS，全站共用一個 bucket
- 嚴重度：資訊　不阻擋
- 位置：設定（`TRUSTED_PROXY_CIDRS` 預設空；`09-cicd-env/02-environments-v0.1.md` 已列）
- 影響：經 Cloudflare Tunnel 或本機反向代理時，peer 都是代理位址（如 127.0.0.1）。未設可信 CIDR 時，所有訪客共用額度：create 全站 30 次／分、redirect 120 次／分，等於自我 DoS。這是正確的「預設不信任」行為，但上線時必須設定。
- 建議：維運在 G5 部署檢查表加入「代理後必設 `TRUSTED_PROXY_CIDRS`＝代理位址 CIDR」。另記觀察：固定一分鐘窗在窗交界最多允許 2 倍突發，SEC-007 只要求「基本限流」，可接受。
- 狀態：資訊

### [SCR-010] `HSTS_ENABLED` 在現行程式中不會生效（無行程內 TLS）
- 嚴重度：低　不阻擋
- 位置：`internal/httpapi/security_headers.go:32、39-41`（`r.TLS != nil`）；`cmd/shorturl/main.go:69` 只有 `ListenAndServe`（明文），repo 中沒有 `ListenAndServeTLS`
- 影響：即使設 `HSTS_ENABLED=true`，請求永遠不是 TLS，所以永遠不送 HSTS。開關名稱容易讓維運誤以為 production 已有 HSTS。
- 裁定與修補：見下節「HSTS 裁定」。要求在 `.env.example` 註解與維運文件寫明「TLS 在邊緣終止時，HSTS 由邊緣送出，本開關無效」。
- 狀態：開放（文件修補，可在本 PR 或維運文件完成；G5 前驗證）

### [SCR-011] PR #7 與 PR #8 語意衝突；Rulesets 不要求分支最新，後合併者會讓 main 變紅
- 嚴重度：**中**　**阻擋後合併的那一個 PR**（流程／CI 完整性）
- 位置：`internal/httpapi/security_headers_test.go:87、96-102、155`（PR #8 用舊簽章 `ratelimit.NewGuard()` 與 5 參數 `NewHTTPAPIAdapter`）；PR #7 把兩者改為 `NewGuard(int,int)` 與 6 參數
- 證據：本機唯讀模擬 `main`＋PR #7＋PR #8，文字合併無衝突，但 `go vet` 出現 `not enough arguments in call to httpapi.NewHTTPAPIAdapter`，`go test ./internal/httpapi` 是 build failed。兩個 run 都只驗證「各自 head ＋ main `161a4fd`」。Rulesets 規則一 `strict_required_status_checks_policy=false`（取證檔），所以第一個 PR 合併後，第二個 PR 的舊綠燈不會失效（v0.1 當時狀態；D-10 後自 2026-10-07 14:33:28（UTC+8）起為 strict=true，見 `ruleset-strict-true-d10-2026-10-07.md`；維運快照 v0.3.1 亦將加註；另見下方處置與狀態），仍可合併，合併後 main 的 `unit-test` 會變紅（正式二進位仍可編譯，失敗在測試）。
- 修補（必要）：後合併的 PR 必須更新分支（rebase／merge main），修正另一 PR 引入的呼叫簽章，並取得**新的四綠 run**，D-08 留言附新 run 連結後才合併。建議維運評估把規則一改為 `strict=true`（需使用者決定，屬 D-08 範圍）。
- 處置（2026-10-07 定案）：
  1. **D-10（王昱殊，2026-10-07，使用者回覆 t83s391）**：規則一 `main-required-checks` 改 `strict=true`，合併前分支須與 main 同步，舊綠燈不得用於合併。因此任何 PR 在前一個 PR 合併後，都必須先更新分支（rebase／merge main）並重新取得四綠，才能合併。strict=true 的設定取證由維運更新 `docs/09-cicd-env/evidence/` 後再引用。**（v0.1.2 更新）已生效**：使用者約 14:33（UTC+8）於 GitHub UI 啟用；維運 14:34:10（UTC+8）唯讀讀取 ruleset 24629805，`strict_required_status_checks_policy=true`、四項 check 名稱不變、`updated_at` 14:33:28（UTC+8），取證檔 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，原始 JSON `docs/09-cicd-env/evidence/rulesets-2026-10-07/ruleset-24629805-d10-public-read.json`、`ruleset-24629805-d10-auth-read.json`（隨 PR #13 提交（head `6674a59`；已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`，取證檔連結已生效））。安全部約 14:36（UTC+8）以公開 REST API 唯讀 GET 同一 ruleset 獨立查證（經 WebFetch；本機 curl 當時受未驗證速率限制回 403），`strict=true`、四項名稱與 `updated_at` 一致。**bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）**（取證檔 §3 第 6 項，依 C 列；v0.1.2 時本機取證檔 B 列與第 6 項仍為佔位，v0.1.3 依隨 PR #13 提交之取證檔 v0.1 更新）。
  2. **已定合併順序：#10 → #12 → #8 → #11 → #7。** PR #8 先於 #7 合併，所以 #8 不需改簽章，只需在 #10、#12 合併後更新分支並取得新四綠。**#7 排最後**：#7 必須在 #8、#11 合併後更新分支，修正 #8（`security_headers_test.go`）與 #11（`observability_test.go`，見 12 號檔 SCR-012）測試中的舊簽章，同時修好 SCR-006、SCR-007，在最新 main 上取得新的四綠 run，再送安全部複核；複核通過並依 D-08／D-09 留言程序後才合併。
- 狀態：**已關閉（v0.1.5）**：已依順序合併（#10 → #12 → #8 → #11 → #7；#8 於 2026-10-19 11:21:48 UTC+8 合併，merge commit `e9f4396`；#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`），合併後 main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8），main 未變紅。歷史：**處置已定（D-10）；strict=true 已生效（取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補））；待依順序（#10 → #12 → #8 → #11 → #7）合併後關閉**（未關閉；#7 依上述條件合併且 main 四綠後，由安全部確認關閉）

---

## PR #7 複核（head ef9fcb8）

> v0.1.2 新增。複核方式：唯讀。cursor-github connector（`get_pull_request`、`get_workflow_run`、`list_workflow_run_jobs`；公開 REST API 受速率限制）＋本機唯讀 clone `/workspace/tmp/short-url-pr7-rv2`（`pull/7/head`、`pull/7/merge`）執行 `go vet`／`go test`／`go test -race`。**沒有**在 PR 上留言或 review，沒有 push，也沒有撰寫攻擊程式；另寫的不變式驗證測試只在本機執行，執行後即刪除，未提交。

### 基本資料與 CI

| 項目 | 結果 |
|---|---|
| PR head | `ef9fcb8cd4999cc40dab990dda0d091a851c8615`（PR 狀態 open、未合併；`get_pull_request` 14:35 UTC+8 讀取） |
| 本輪 commit | `8b36767`（merge main 進分支，非 docs 變更只有 PR #6 的 `.github/workflows/ci-devsecops.yml`，與 main 相同）→ `a5b3d2d`（SCR-006：`guard.go`、`guard_test.go`）→ `ef9fcb8`（SCR-007／008：`source.go`、`source_test.go`） |
| 「已合入 main 63b373d」的意思 | 查證：是 **main `63b373d`（PR #9 合併）被 merge 進 PR #7 分支**（`8b36767`），**不是** #7 已合併進 main。`63b373d` 是 `ef9fcb8` 的祖先；PR base sha 與 main 皆為 `63b373d`；`git diff origin/main ef9fcb8` 沒有動到 `.github/`、`go.mod`、`go.sum`、`docs/` |
| CI run | [37581955467](https://github.com/444112029012/short-url/actions/runs/37581955467)：event `pull_request`，`head_sha`＝`ef9fcb8…` ✅，status completed／conclusion **success**；建立 14:31:19（UTC+8），完成 14:32:41（UTC+8） |
| 四項必要檢查 | `secrets-gitleaks` ✅ success、`sast-gosec` ✅ success、`sca-govulncheck` ✅ success（含 Probe dependency graph、Dependency review 兩步 success）、`unit-test` ✅ success（`go vet`、`go test` 兩步 success）；四個 job 的 `head_sha` 都是 `ef9fcb8` |
| 測試的 merge ref | `b074cb8`＝`ef9fcb8` ＋ main `63b373d`；tree `c355cf3` 與 head tree **完全相同**（head 已含當時 main），所以這次綠燈等同 head 本身 |
| 尚未涵蓋 | 與 #8（`NewHTTPAPIAdapter` 6 參數）、#11（`NewGuard(int,int)`）的簽章衝突**仍在**；#7 排最後合併，將在更新到最新 main 時修正（SCR-011），屆時是新 head，須再送核 |

### SCR-006（表滿處理）：修正已驗證

| 檢查點 | 結論與依據 |
|---|---|
| 表滿是否淘汰「窗口最舊」那筆 | **是**。`guard.go:146-149`：先 `dropExpiredFront` 移除已過期的前端，若仍 `Len() >= capSources` 就 `evictFront`（`:169-179`）移除 `order.Front()`，再 `PushBack` 新來源（`:150-151`），**不回 429**。被淘汰的是 list 前端，也就是 `windowStart` 最早的那筆，不是最近插入的那筆，也不是任意一筆。 |
| list 是否維持依窗口起點排序 | **是**。新 bucket 的 `windowStart=now` 且放到隊尾（`:150-151`）；窗口重置時 `windowStart=now`、`count=1`，並 `MoveToBack`（`:133-137`）；窗口內命中只 `count++`，位置不動（`:140-142`）。三條路徑都維持「前舊後新」。本機隨機操作 20 萬步（假時鐘、cap 16、40 個 key、create／redirect 交錯）逐步檢查：list 依 `windowStart` 非遞減、`map` 與 `list` 一一對應、大小 ≤ cap，全部通過。 |
| map／list 同步與並發 | **同步**。新增：`byKey[k]=PushBack(...)`（`:151`）；淘汰：`Remove(el)`＋`delete(byKey, b.key)`（`:175-177`），`bucket.key` 與 map key 相同；防禦路徑 `b == nil` 時兩邊一起移除（`:130-132`）；重置只移動同一個 element，map 指標不變。兩張表的所有變更都在同一個 `g.mu` 下（`:125-126`）。`-race` 全套件通過；本機 16 goroutine × 5000 次（真實時鐘、cap 64）結束後 list 與 map 長度一致且 ≤ 64。觀察：`now` 在取鎖前讀取（`:116-119`），並發時插入順序與 `windowStart` 可能出現微秒級逆序，影響只是「近似最舊」，列 SCR-018（資訊）。 |
| 記憶體上限 | **有界**。每表上限 `defaultMaxSources=8192`（`:16`），每次插入前若已滿就先淘汰一筆，list 與 map 都不超過 8192；兩表合計 16384 個 bucket。key 只來自 `ClientSource` 的輸出（IPv4 字串、`<prefix>/64`、或 `RemoteAddr` 主機字串，後者由 `net/http` 依連線填入，用戶端無法塞入任意長字串）。沒有發現破口。 |
| 拿掉 `evictIdle` 後過期項目是否回收 | **會**。每次插入新來源前，`dropExpiredFront`（`:155-167`）從前端持續移除已過期的 bucket；list 有序，過期的 bucket 必定集中在前端，攤銷 O(1)。本機驗證：1000 個來源過期後再插入一個新來源，表大小變成 1。若之後完全沒有新來源，過期項目會留在表內，但仍受 8192 上限約束，回訪的來源則會就地重置。原本持鎖 O(n) 全表掃描已移除，攻擊下的 CPU 放大也隨之消除。 |
| 攻擊者大量偽造來源能否持續擠掉合法來源 | **不會造成拒絕服務**，但精度會下降（殘餘風險，不阻擋）。合法來源被淘汰後只是計數歸零，再來時重新放行，不會被拒。代價是：一分鐘內不同 key 超過 8192 時 FIFO 會循環，每個 key 在被淘汰前可能還沒累積到門檻，所以**在洪泛期間，每來源限流對所有輪替的 key（含攻擊者）實質失效**。門檻：IPv4 需要 > 8192 個不同位址（botnet）；IPv6 需要 > 8192 個不同 /64（約一個 /51），單一取得 /48 配置的用戶端就做得到。此項已寫入 05 RR-007（補償：邊緣速率限制／WAF；G5 前重評；未來可評估 /48 二層聚合或全域上限）。 |
| 測試是否真的覆蓋 | `TestGuardFullTableEvictsOldestAndAdmits`（`guard_test.go:108-137`）：cap 2，插入第 3 個來源時**放行**、list 與 map 都 ≤ 2、`oldest` 被淘汰、`newer` 與 `admitted` 留下、前端是 `newer`。redirect 表同類測試（`:139-164`）。`TestGuardWindowResetMovesBucketToBack`（`:166-198`）與 redirect 版（`:200-225`）：重置後的 bucket 在隊尾，表滿時淘汰的是較舊且仍在窗內的 `still-open`，不是剛重置的 bucket。`TestGuardConcurrentExactQuota`（`:227-246`）：320 個並發請求只放行 80 個。這些都是有實際斷言的測試，不是空測試。 |

### SCR-007（IPv6 聚合）：修正已驗證

| 檢查點 | 結論與依據 |
|---|---|
| 遮罩實作 | **正確**。`source.go:87-100` `bucketKey`：先 `ip.To4()`，非 nil 就回傳 IPv4 字串（/32，`:91-93`）；否則 `To16()` 後 `Mask(net.CIDRMask(64, 128))`，key 為 `<prefix>/64`（`:94-99`）。`net.IP.Mask` 對 16 位元組位址和 128 位元遮罩的用法正確。 |
| IPv4-mapped | `::ffff:a.b.c.d` 經 `To4()` 轉成 IPv4，與一般 IPv4 共用額度（`TestIPv4MappedEqualsPlainIPv4`，`source_test.go:100-111`；本機另驗 XFF 內的 mapped 位址也轉成 IPv4）。 |
| zone、無效位址 | `parseIP`（`:135-145`）使用 `net.ParseIP`，不接受 `%zone`。帶 zone 的 peer（只可能是同鏈路 link-local）會以原字串作為 key，不聚合，也不會被當成可信代理；XFF 中帶 zone 或無效的項目會被略過。影響可忽略，列 SCR-018（資訊）。peer 無法解析時回傳原主機字串（`:60-63`），空字串由 guard 對應到共用 bucket（`guard.go:113-115`）。 |
| 經可信代理取得的位址 | 套用同一個 `bucketKey`（XFF `:66-67`、X-Real-IP `:69-70`、退回 peer `:72`）。測試 `TestTrustedProxyIPv6IsAggregated`（`source_test.go:113-138`）用三個不同的可信 peer，斷言同一 /64 的 XFF 共用額度、不同 /64 各自獨立。 |
| 測試 | `TestIPv6PrefixSharesQuota`（`:73-98`）：同一 /64 的兩個位址 key 相同，create 與 redirect 第二次都 429；不同 /64 放行。`TestClientSourceIPv6AndEmpty`（`:59-71`）已改為斷言 `2001:db8::/64`。 |

### SCR-008（XFF 多行／X-Real-IP）：修正已驗證

| 檢查點 | 結論與依據 |
|---|---|
| 多行合併 | `joinedForwardedFor`（`source.go:78-83`）用 `Header.Values("X-Forwarded-For")`，依序以逗號串接，不再只讀第一行。 |
| 由右往左、跳過可信代理 | `clientFromForwarded`（`:114-133`）從最右邊往左逐一解析，跳過無效項目和屬於 `trusted` 的項目，回傳第一個非可信的有效位址；如果全部都是可信位址，則回傳最右邊的有效位址。正確。 |
| 空值／垃圾值 | 空白標頭直接回 nil（`:115-117`）；空項目和無效項目略過（`:121-124`）；含 port 的 `IPv4:port`、`[IPv6]:port` 可正確解析（`:140-143`）。本機驗證：可信 peer 帶垃圾 XFF、沒有 X-Real-IP 時，退回 peer 位址（代理位址，共用一個 bucket，屬保守方向）。 |
| X-Real-IP 退回條件 | 只在 peer 可信**而且** XFF 沒有任何有效位址時才使用（`:66-71`）。正常會附加 XFF 的代理一定會在最右邊附加一個有效位址，所以用戶端無法讓程式走到退回分支；用戶端自帶的 `X-Real-IP` 在這種情況下不會被採用（`TestClientSourceTrustedProxyUsesRightmostUntrustedForwarded`，`source_test.go:24-40`；`TestForwardedForSeparateHeaderLineDoesNotLetClientWin`，`:140-157`，兩者都設定了偽造的 X-Real-IP）。只有在代理既不附加 XFF、也不覆寫 X-Real-IP 的錯誤設定下，才會採用用戶端的值；程式註解（`:44-46`）已寫明，維運部署文件仍需註記（RR-004）。 |
| 未設定 `TRUSTED_PROXY_CIDRS` | `ipInNets` 在 nets 為空時回傳 false（`:147-157`），所以完全不讀任何轉傳標頭，只用 peer（`:60-65`）；`.env.example` 已註明預設為空、永不信任。測試 `TestClientSourceIgnoresSpoofedHeadersWhenUntrusted`（`:8-22`）。 |
| 多行測試 | `TestForwardedForSeparateHeaderLineDoesNotLetClientWin`：用戶端行在前、代理附加行在後時取 `203.0.113.50`；代理附加的是可信 hop 時，往左取 `198.51.100.8`。有實際斷言。 |

### 研發附註（/64 聚合放在 ClientSource、guard 不聚合）判定：**可接受，附條件**

- 呼叫點查證（`ef9fcb8` 全 repo，排除測試）：`CheckCreate`／`CheckRedirect` 的正式呼叫點只有 `internal/httpapi/adapter.go:83` 與 `:114`，兩者都是傳 `a.sourceID(r)`；`sourceID`（`:148-154`）只呼叫 `ratelimit.ClientSource(r, trusted)`。`cmd/shorturl/main.go:56` 只建立 `NewGuard`。**沒有任何正式呼叫點把原始 IP 直接傳進 guard**，因此在 guard 前做聚合可以接受，不阻擋。
- 條件：`guard.go:92-107`（`CheckCreate`／`CheckRedirect`）與 `adapter.go:22-26`（`RateLimitChecker` 介面）的註解都**沒有寫明**「sourceID 必須是 `ClientSource` 的正規化輸出」，也沒有型別保護。列為 SCR-017（低），要求在 #7 下一個 head 補上註解；並建議加型別或測試保護（見 SCR-017）。

### 新發現項

#### [SCR-017] Guard 對 sourceID 的前置條件（須先經 ClientSource 正規化）未標明，也沒有型別或測試保護
- 嚴重度：低　不阻擋（列入 #7 新 head 送核時的核對項）
- 位置：`internal/ratelimit/guard.go:92-107`、`internal/httpapi/adapter.go:22-26`、`:148-154`（commit `ef9fcb8`）
- 類別：其他（縱深防禦／限流）；TM-003、TM-017；CWE-1288
- 影響：目前所有呼叫點都正確。但如果日後有人新增呼叫點（例如新端點或中介層），直接把 `RemoteAddr` 或原始 IPv6 字串傳給 guard，IPv6 /64 聚合就會失效，SCR-007 的問題會無聲復發。現有測試沒辦法發現這種情況。
- 重現：（僅邏輯說明）`g.CheckCreate("2001:db8::1")` 與 `g.CheckCreate("2001:db8::2")` 會落在兩個不同的 bucket。
- 修補建議：(1) **要求**：在 `CheckCreate`／`CheckRedirect` 與 `RateLimitChecker` 的註解寫明前置條件「sourceID 必須是 `ratelimit.ClientSource` 的輸出，guard 不做位址正規化」。(2) **建議**（二擇一）：(a) 定義 `type SourceKey string`，只由 `ClientSource` 產生，guard 介面改收 `SourceKey`；或 (b) 在 `internal/httpapi` 加一個測試，用會記錄 sourceID 的假 limiter，斷言同一 /64 的兩個 `RemoteAddr` 傳進 limiter 的值相同、IPv4-mapped 與一般 IPv4 相同。
- 狀態：**已關閉（v0.1.5：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8））**。修補建議 (1) 已於 commit `fc9a7a7` 完成（三處註解，見 16 號檔 v0.1.1 §8）；建議 (2) 未做，列建議追蹤

#### [SCR-018] 取鎖前讀時間造成 list 順序可能微幅偏離；zone 位址未正規化
- 嚴重度：資訊　不阻擋
- 位置：`internal/ratelimit/guard.go:116-119`（`now` 在 `g.mu.Lock()` 之前讀取）；`internal/ratelimit/source.go:135-145`（`net.ParseIP` 不接受 zone）（commit `ef9fcb8`）
- 類別：其他（業務邏輯精度）／輸入驗證
- 影響：(a) 並發時，先讀時間但後取得鎖的請求，插入位置會在 `windowStart` 稍晚的 bucket 之後，造成微秒級逆序。後果只是 `dropExpiredFront` 可能晚一點才回收、FIFO 淘汰是「近似最舊」，不影響上限與正確性。(b) 帶 zone 的 link-local peer 以原字串作為 key（不聚合），只有同鏈路可達，影響可忽略；XFF 內帶 zone 的項目會被略過。
- 修補建議：(a) 把 `now` 的讀取移到取鎖之後；(b) 可改用 `netip.ParseAddr`，再用 `WithZone("")` 與 `Unmap()` 正規化。可在 #7 下一個 head 或之後處理。
- 狀態：開放（建議）

### 本機測試（唯讀 clone，Go 1.27.1，於 head `ef9fcb8`）

| 指令 | 結果 |
|---|---|
| `go vet ./...` | 通過 |
| `go test -count=1 ./...` | 全部 ok（application、config、domain、health、httpapi、infra/sqlite、observability、ratelimit） |
| `go test -race -count=1 ./...` | 全部 ok |
| 研發宣稱的 8 個新增或改寫測試，加上 `TestGuardConcurrentExactQuota`（`-v`） | 全部 PASS |
| 安全部臨時不變式測試（只在本機執行，執行後刪除） | 隨機 20 萬步的順序／同步／上限檢查、淘汰最舊（不是最近命中的）、過期回收、並發同步、來源邊界案例（zone、mapped、垃圾 XFF、帶 port、IPv6 代理鏈、未設可信 CIDR）全部符合預期 |

### PR #7 複核結論

- **v0.1.5（合併後）**：最終 head `fc9a7a7` 已經 16 號檔 v0.1.1 通過並合併（#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8））。SCR-002、SCR-006、SCR-007、SCR-008、SCR-017 關閉；SCR-011 依 main 四綠關閉；RR-007 降為低（05 v0.1.10）。以下為 v0.1.2 當時結論，保留為紀錄。
- **SCR-006**：修正已驗證，待 PR #7 合併後關閉。**SCR-007**：修正已驗證，待 PR #7 合併後關閉。**SCR-008**：修正已驗證，待 PR #7 合併後關閉（維運文件註記另由 RR-004 追蹤）。
- **新發現**：SCR-017（低）、SCR-018（資訊），都不阻擋。
- **整體：邏輯上核可（head `ef9fcb8`）。** 依 skill 阻擋條件，目前沒有未修正的嚴重或高發現項，四項掃描皆綠，SCR-006／007 兩個中級阻擋項已修正。**合併前條件**：(1) #7 依順序排最後，更新到最新 main 並修正 #8（`security_headers_test.go`）與 #11（`observability_test.go`）的舊簽章後，**新 head 的差異須再送安全部核對**（核對範圍：衝突解決是否只改簽章或呼叫，沒有動到限流與來源識別邏輯；以及 SCR-017 的註解）；(2) 在新 head 上取得**四項新綠燈**（strict=true 已生效，舊綠燈不能用來合併）；(3) 依 D-08／D-09 留言程序合併。
- **SCR-011**：處置已定（D-10）；strict=true 已生效（取證檔隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補））；待依順序合併後關閉。**SCR-002**：維持開放，#7 合併時關閉。
- **RR-007**：#7 修正經安全部複核通過（head `ef9fcb8`）；**合併至 main 後降為低**，合併前仍為中（見 05 v0.1.7）。

---

## 表滿造成全域拒絕：裁定與修法（摘要）

- **裁定**：現行「表滿→新來源一律 429」**不接受**，阻擋 PR #7（SCR-006）。表滿時應 fail-soft：淘汰最舊窗的 bucket、放行新來源，記憶體仍有界；精度損失只在活躍來源超過容量時出現，影響有界。
- **修法**：(1) FIFO 淘汰最舊 bucket（O(1)），取代 O(n) `evictIdle` 掃描與拒絕；(2) IPv6 /64 聚合（SCR-007）；(3) 不記錄位址；(4) 改寫與補充測試（表滿仍放行、淘汰最舊、redirect 表同測）。
- **不採用的替代方案**：(a) 表滿時 fail-open 且不計數：效果與淘汰相近，但滿表期間所有新來源完全不受限，比淘汰更寬鬆；(b) 表滿時新來源共用一個溢出 bucket：攻擊者可以耗盡溢出額度，正常新訪客仍被拒，只是把全域拒絕縮小一點，不足；(c) 只把容量調大：提高門檻，但沒有改變失敗方向，IPv6 下仍可填滿。

## HSTS 裁定

- **裁定：PR #8 的 HSTS 行為可接受，不阻擋。** 「只在明確開啟且請求本身是 TLS 時送出」是正確的保守設計：本機 HTTP 不會宣稱 HSTS（符合 SEC-008(1)「若使用 HTTPS」的條件語意）；也**不**依據 `X-Forwarded-Proto` 判斷，避免用戶端偽造標頭讓明文回應帶 HSTS，這是正確的。
- `max-age=31536000` 符合 SEC-008(1) 下限。不含 `includeSubDomains`／`preload`：demo 網域範圍未定，這是正確的保守選擇。
- 依現行架構，production 由 Cloudflare Tunnel／反向代理終止 TLS（`02-environments`、C4 container），應用層永遠看不到 TLS，**SEC-008(1) 必須由邊緣滿足**：邊緣設定 HSTS（max-age ≥ 31536000），G5 前以實際回應標頭驗證。應用層開關只適用將來行程內 TLS（`TLS_CERT_FILE`／`TLS_KEY_FILE`，目前未實作）。文件化要求見 SCR-010。
- 不要求本 PR 改為信任 `X-Forwarded-Proto`。若將來要由應用層送出，只能在 peer 屬 `TRUSTED_PROXY_CIDRS` 時採信，並另案審查。

## SCR-002／SCR-003 處置

| 編號 | 處置 | 理由 |
|---|---|---|
| SCR-002（限流 stub） | **已關閉**（v0.1.5：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`） | v0.1.2：SCR-006／007 修正已於 head `ef9fcb8` 複核通過，但 #7 尚未合併，SCR-002 維持開放。原說明：PR #7 實作了 SEC-007，但 SCR-006／007 阻擋。#7 依合併順序排最後：於 #8、#11 合併後更新分支、修正後的新 head 經安全部複核、四綠，並在 PR #7 合併時關閉 |
| SCR-003（缺 SEC-008 標頭） | **已關閉**（v0.1.6：#8 於 2026-10-19 11:21:48（UTC+8）合併，merge commit `e9f43968236878eb605e8e3048bae59a35815e70`）。原處置：PR #8 合併時關閉（SEC-008 內容不需修正） | PR #8 已滿足 SEC-008 (2)(3)(4) 與 3xx 要求；(1) HSTS 依上節裁定由邊緣滿足（G5 驗證，SCR-010 追蹤）。依 D-10，#8 合併前須與 main 同步並取得新四綠（合併順序 #8 先於 #7，見 SCR-011 處置） |

## CI 對應確認

| 項目 | PR #7 run 37579834778 | PR #8 run 37580036601 |
|---|---|---|
| run head_sha＝PR head | ✅ `04815e5…`（event `pull_request`，conclusion success） | ✅ `9dbc246…`（event `pull_request`，conclusion success） |
| 建立時間（UTC+8） | 2026-10-07 14:07:40 | 2026-10-07 14:09:59 |
| 測試的 merge ref | `a59870b`＝head ＋ main `161a4fd`（含 PR #6 fail-closed workflow） | `97ae999`＝head ＋ main `161a4fd` |
| 四 required checks | secrets-gitleaks ✅、sast-gosec ✅、sca-govulncheck ✅、unit-test ✅ | 同左全 ✅ |
| Probe／dependency review | `dependency graph compare HTTP 200`（BASE 161a4fd／HEAD 04815e5）；`Dependency review (PR only)` 實跑，`fail-on-severity: critical`，無漏洞 | `HTTP 200`（BASE 161a4fd／HEAD 9dbc246）；實跑，無漏洞 |
| GITHUB_TOKEN 權限 | Contents／Metadata／PullRequests: read | 同左 |
| govulncheck | No vulnerabilities found | No vulnerabilities found |
| SEC 驗收在 CI 的對應 | SEC-007：`internal/ratelimit` 單元測試（門檻、窗重置、獨立性、並發、預設值）＋`internal/httpapi/ratelimit_test.go`（429、redirect 不計數、XFF 偽造、可信代理）皆在 `unit-test` job 的 `go test` 內 | SEC-008：`security_headers_test.go`（201／302／400／404／405／429／500 標頭、HSTS 僅 opt-in＋TLS）在 `unit-test` 內 |
| 本機複核 | `go test -race` ratelimit／httpapi 通過 | `go test ./...` 通過 |
| 限制 | 兩個 run 都沒有涵蓋「#7＋#8 同時在 main」的組合；組合會 build failed（SCR-011） | 同左 |

## 殘餘風險（v0.1 建議；v0.1.1 起已依 D-11 登記於 05）

- **RR-007（中→修正後低；已由使用者於 D-11 接受，登記於 05）**：限流是單程序記憶體、固定窗、每來源計數。分散來源（大量 IP／botnet）的總量攻擊、多實例部署（計數不共享）、容量淘汰時的精度損失都無法由本控制消除。補償：邊緣（Cloudflare）速率限制／WAF、單實例部署假設寫入維運文件；G5 前再評。
- **RR-008（低；已由使用者於 D-11 接受，登記於 05）**：SEC-008(1) HSTS 依賴邊緣設定，應用層無法保證（SCR-010）。補償：G5 部署檢查以實際回應驗證 `Strict-Transport-Security`。
- RR-004（既有）：SCR-008 屬其範圍，已補記於 05 一覽表與 §2.4。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | 初版：PR #7 要求修改（SCR-006 表滿全域拒絕、SCR-007 IPv6 未聚合）；PR #8 核准（SCR-010 HSTS 文件化）；跨 PR SCR-011（語意衝突＋strict=false）；SCR-008 低、SCR-009 資訊；HSTS 裁定；SCR-002 待修正後關、SCR-003 PR #8 合併時關；本文上限引用 10 號檔；MaxHeaderBytes 不在 #7／#8 範圍，只引用 10 號檔 C3 |
| v0.1.1 | 2026-10-07 | 安全部 | 依 D-10／已定合併順序（#10 → #12 → #8 → #11 → #7）：SCR-011 處置改為 strict=true＋#7 最後更新分支並修 SCR-006／007 後送複核，狀態改「處置已定（D-10）；strict=true 取證待維運」；PR #8 條件句、SCR-002／SCR-003 處置同步；RR-007／008 標註已依 D-11 接受、RR-004 已補記 SCR-008；修正 SCR-011 行號（159→155）；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-07 | 安全部 | 新增「PR #7 複核（head ef9fcb8）」：run 37581955467 head_sha＝ef9fcb8、四項 success、merge ref tree 同 head；「已合入 main 63b373d」查證為 main 併入分支（非 #7 已合併）；SCR-006／007／008 修正已驗證、待 #7 合併後關閉；研發附註（聚合在 ClientSource）判定可接受（呼叫點僅 adapter.go:83／114 經 ClientSource），附條件補前置條件註解；新增 SCR-017 低、SCR-018 資訊；PR #7 邏輯上核可，合併前新 head 差異須再送核並取得四項新綠燈；SCR-011 改「strict=true 已生效（取證檔隨維運 PR A，bypass 待具權限讀取確認）；待依順序合併後關閉」；SCR-002 維持開放；小版修訂不改檔名（檔名保留主版號） |
| v0.1.3 | 2026-10-19 | 安全部 | SCR-011 證據／處置／狀態與結論：bypass 寫法由「待具權限讀取確認」改為「bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）」，引用 `ruleset-strict-true-d10-2026-10-07.md`；「維運 PR A、待合併」改為「隨 PR #13 提交（截至 2026-10-19 仍 open、未合併）」；對照列與 SCR-011 證據句之 rulesets 快照 strict=false 保留為快照當時事實，加註 D-10 後自 2026-10-07 14:33:28（UTC+8）起為 strict=true、維運快照 v0.3.1 亦將加註；歷史修訂列不改；小版修訂不改檔名（檔名保留主版號） |
| v0.1.4 | 2026-10-19 | 安全部 | SCR-011 證據／狀態／結論：取證檔「隨 PR #13 提交（截至 2026-10-19 仍未合併／09:02 仍 open）」改為已合併事實（2026-10-19 09:05:27 UTC+8，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`，取證檔連結已生效；cursor-github 唯讀查證）；SCR-011 處置與關閉條件不變；小版修訂不改檔名（檔名保留主版號） |
| v0.1.5 | 2026-10-19 | 安全部 | #7 合併後更新：#7 已於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`（最終 head 經 16 號檔 v0.1.1 通過）；main push run [39198199210](https://github.com/444112029012/short-url/actions/runs/39198199210) 四項 success（12:11:12～12:12:50 UTC+8）；SCR-006／007／008／017 與 SCR-002 改「已關閉」；SCR-011 依順序合併完成且 main 四綠改「已關閉」（#8 merge `e9f4396` 11:21:48 UTC+8）；狀態列、摘要表結論列、PR #7 複核結論加註合併結果；RR-007 降低見 05 v0.1.10；歷史敘述保留；小版修訂不改檔名（檔名保留主版號） |
| v0.1.6 | 2026-10-19 | 安全部 | SCR-003 改「已關閉」：#8 於 2026-10-19 11:21:48（UTC+8）合併，merge commit `e9f43968236878eb605e8e3048bae59a35815e70`（公開 API／git 唯讀查證）；狀態列同步；SEC-008(1) HSTS 仍由 SCR-010（低、開放、G5 前驗證）追蹤，不隨 SCR-003 關閉；小版修訂不改檔名（檔名保留主版號） |
