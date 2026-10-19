---
文件：PR #20（設計部文件 PR：OpenAPI 0.1.1 本文上限／431、方法規格 v0.2.2）安全意見——回應審查部 G4-PR20-R3
版本：v0.1.1
狀態：**v0.1.1：#20 已合併**（2026-10-19 15:15:58（UTC+8），merge commit `abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`，合併的 head 即核准的 `44b099c`）；**SCR-043～045 已關閉**；SCR-046 資訊、開放。v0.1：通過，僅適用於最終 head `44b099c227cac2ab24380f4346654540da6535de`（附 §7 合併條件）
負責角色：安全部
作者：安全部
日期：2026-10-19
最後更新：2026-10-19
對應：安全部意見 [10](./10-dr-g4-01-body-limit-security-opinion-v0.1.md)（DR-G4-01，§1「維持 OpenAPI 四類錯誤，回應不帶上限值或內部細節」）；[17](./17-pr-017-dr-g4-01-body-limit-security-confirmation-v0.1.md) SCR-027（431 由 net/http 產生、不帶安全標頭之接受理由）；資安需求 [01 v0.2](./01-sec-requirements-v0.2.md) SEC-001、SEC-003、SEC-007、SEC-008
專案代號：SHORTURL
發現項編號區段：SCR-043～SCR-049（本文件使用 SCR-043～SCR-046）
---

# PR #20 安全意見（v0.1）：回應審查部 G4-PR20-R3

## 0. v0.1.1 合併後狀態（2026-10-19 15:18 UTC+8 查證）

- PR #20 API：`merged: true`、`merged_at` 2026-10-19T07:15:58Z＝**15:15:58（UTC+8）**（merge commit 的 commit date 為 15:15:57 +0800）、`merge_commit_sha`＝`abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`、`head.sha`＝`44b099c227cac2ab24380f4346654540da6535de`。
- `abbec9e` 的 parents 是 `141fcda38e1cb268f04aacdeba28db1c86e7ceca`（當時的 main）與 `44b099c`（本文件核准的 head）；`git diff 44b099c abbec9e` 為空（0 byte）：合併的就是核准的 SHA，符合 §7 第 1 點。main 已是 `abbec9e`。
- 因此 SCR-043、SCR-044、SCR-045 改為**已關閉**（§8）。SCR-046（資訊，根目錄 README）仍開放。
- 合併當下的審查部關卡與核准 review 由 GitHub ruleset 強制，本文件未另查。

## 1. 結論

**判定：通過，僅適用於最終 head `44b099c227cac2ab24380f4346654540da6535de`。** 核准只適用此 SHA；之後的任何 commit（包括同步 main 的 merge commit），本核准即不適用，都要重新送安全部核對（見 §7）。

- 審查部 G4-PR20-R3 三個問題的答案：(1) 仍符合 10 號意見；(2) 公開上限值與 431 源站不帶安全標頭可接受，SCR-027 可延伸到規格層（措辭已依 SCR-043 修正）；(3) 只動 `docs/`，不需另做完整安全審查，安全部對最終 head 做差異核對即可（本文件即該次核對）。
- 初步意見（head `519e030`，2026-10-19 14:19 UTC+8 回報，未落檔）的三項資訊發現 SCR-043～045 已在 `3c037a1` 修正，狀態「待 #20 合併後關閉」。
- 新發現 SCR-046（資訊，不在本 PR 範圍）：main 根目錄 `README.md`（#19）的 431 描述仍是修正前寫法。
- 方法：GitHub 讀取 API（PR、run／jobs、job log）、`git fetch` 至 `/workspace/tmp/pr18-review/short-url`；在本機以安全部自寫的探測程式（非 PR 程式碼、非 exploit，用後刪除）驗證 Go 1.27.2 net/http 的 431 行為。未在 GitHub 留言／review／核准／合併／push。

## 2. 查證過的 SHA、run 與 PR 狀態（2026-10-19 14:59（UTC+8）核對）

| 項目 | 宣稱 | 查證結果 |
|------|------|----------|
| head | `44b099c227cac2ab24380f4346654540da6535de` | ✅ PR API `head.sha` 與 `git ls-remote refs/pull/20/head` 相同；commit 14:50:35 UTC+8 |
| parents | `3c037a1b9b1e24c924a922bac436fe393abb4e2a`、`141fcda38e1cb268f04aacdeba28db1c86e7ceca` | ✅ 相同 |
| main | `141fcda`（#19 合併後） | ✅ `refs/heads/main`＝`141fcda38e1cb268f04aacdeba28db1c86e7ceca`；PR #19 API：`merged: true`、`merged_at` 2026-10-19T06:49:27Z＝**14:49:27（UTC+8）**、`merge_commit_sha`＝`141fcda…` |
| #18 | — | **已合併**：PR #18 `merged_at` 2026-10-19T06:24:23Z＝**14:24:23（UTC+8）**，merge commit `643bf8ab052ebb175e35bc365ec97863ec9f4202`（parents `433140c`＋`f547a6e`，合併的 head 即安全部 18 號檔 v0.1.3 核准的 `f547a6e`；`git diff f547a6e 643bf8a` 為空） |
| remerge-diff | 空 | ✅ 安全部自行執行 `git show --remerge-diff 44b099c`，輸出 0 byte |
| 與 main 差異 | 6 檔，全在 `docs/` | ✅ `git diff --stat 141fcda 44b099c`：`docs/02-architecture/06-method-specs-v0.2.md`（17）、`docs/02-architecture/08-ruling-G4-PR2-R1-go-mapping-v0.1.md`（5）、`docs/02-architecture/README.md`（8）、`docs/04-api/01-api-notes-v0.1.md`（22）、`docs/04-api/README.md`（9）、`docs/04-api/openapi.yaml`（39），+77／−23。沒有程式碼、`.github`、`go.mod`、`go.sum` |
| 08-ruling 是否為新增 | 「本次新增，519e030 時還沒有」 | ❌ **宣稱不正確**：該檔自 `6c7104c`（2026-10-07 13:58:53 UTC+8）起就在 main，`433140c`、`519e030`、`141fcda` 的 blob 都是 `47a42a96e02b3fd9bed459193baec8ab2175d7d6`。本 PR 是**修改**（`M`，blob 改為 `a6da63628a70e6afecf90f8a781ff67e14ba67fb`，於 `3c037a1` 改動），不是新增；`519e030` 時未改動此檔。內容見 §5 |
| merge ref | — | `refs/pull/20/merge`＝`0bb4e12`，parents `141fcda`＋`44b099c`，tree＝head tree＝`ea78dbfc` |
| PR 狀態 | — | `open`、`mergeable_state: blocked`（缺核准 review）、3 commits |

### 2.1 run 39212359855（四個 job 的 head_sha 皆為 `44b099c`）

event `pull_request`、run `head_sha`＝`44b099c…`、`conclusion: success`、run_attempt 1，建立 2026-10-19 14:51:37（UTC+8）。

| job（id） | 結論 | log 證據 |
|-----------|------|----------|
| secrets-gitleaks（117928385818） | success | `59 commits scanned.`、`no leaks found` |
| sca-govulncheck（117928386014） | success | `No vulnerabilities found.`；`dependency graph compare HTTP 200`；Dependency review 無 critical、無 denied packages |
| sast-gosec（117928386043） | success | `gosec h1 OK`；`Files : 19`、`Lines : 1719`、`Nosec : 0`、**`Issues : 0`** |
| unit-test（117928386061） | success | `go test ./...` 全部套件 `ok` |

### 2.2 `519e030 → 44b099c` 的每個 commit

| commit | parents | 內容 |
|--------|---------|------|
| `3c037a1` | `519e030` | docs only，5 檔（方法規格、08-ruling、02-architecture README、API 說明、openapi.yaml），+31／−17：SCR-043～045 修正、R4、審查部 S1～S5／S9／S10（`x-max-body-bytes`、根層 `security: []`、400 描述、§5 補列實際安全標頭、08-ruling 8192 後續註記）；commit 訊息末行 `REQ-011 SEC-001 SEC-007 SEC-008 ENG-007 DR-G4-01 G4-PR20-R2 G4-PR20-R4 SCR-043 SCR-044 SCR-045`，正文註明補記 `519e030` 缺漏的關聯行 |
| `44b099c` | `3c037a1`、`141fcda` | 同步 main 的 merge，remerge-diff 為空；帶進的 main 增量是 #18（`internal/integration/ratelimit_e2e_test.go`）與 #19（根目錄 `README.md`），兩者都已在 main |

### 2.3 main `433140c → 141fcda` 是否改過程式

`git diff --stat 433140c 141fcda`：只有根目錄 `README.md`（#19）與 `internal/integration/ratelimit_e2e_test.go`（#18）。**`cmd/shorturl/main.go`、`internal/httpapi/adapter.go`、`internal/httpapi/mapper.go` 都沒有變動**，方法規格「對照依據 main `433140c`」與文件／程式的對照仍成立。

## 3. 回應審查部 G4-PR20-R3 三個問題

### 3.1 問題 1：是否仍符合 10 號意見？——**符合**

- 10 號意見（10 號檔第 18 行）要求的是**回應**不帶上限值或內部細節，並維持 OpenAPI 四類錯誤；不禁止**規格文件**寫出上限。repo 公開，常數本來就可見（`adapter.go:23` `maxCreateBody = 32768`、`main.go:22-23` `maxHeaderBytes = 16 << 10`）。
- 本文超限：main 沒有 413。`adapter.go:171-172` 以 `LimitReader(maxCreateBody+1)` 截斷，超限回 `InvalidURL(MsgURLValidationFailed)`；`mapper.go:48-75` 只輸出 `{error:{type,message}}`，message 取固定目錄。回應不帶上限值。
- 四類錯誤：`openapi.yaml:288`（`44b099c`）`enum: [invalid_url, not_found, rate_limited, internal_error]` 未變；本 PR 未新增 413 或第五類錯誤，413 只出現在「不回 413」的說明（`openapi.yaml:59`）。程式（`mapper.go`、`adapter.go`）在 `433140c → 141fcda` 間未變，錯誤訊息目錄不變。
- 431：Go 1.27.2 `net/http/server.go:2084-2096` 寫死回應：`HTTP/1.1 431 Request Header Fields Too Large`、`Content-Type: text/plain; charset=utf-8`、`Connection: close`、本文固定字串 `431 Request Header Fields Too Large`；不經 router，不反射請求內容。安全部本機探測的原始回應與此完全相同。
- 標頭上限：`server.go:960-962` `initialReadLimitSize = MaxHeaderBytes + 4096`。本機實測：請求列＋標頭 20336 B → 200、20516 B → 431；請求列約 21 KB（標頭很小）→ 431；Host＋499 個標頭（500 行）→ 204、Host＋500 個（501 行）→ 431。最終 head 的描述與此一致（§4 SCR-044）。

### 3.2 問題 2：公開上限值、431 不帶安全標頭是否可接受？SCR-027 能否延伸到規格層？——**可接受；可延伸**

- **公開上限值不協助 DoS**：攻擊者可直接讀常數或以二分法試出，公開不增加能力；資源已受限——單請求本文最多讀 32769 B、建立端點 30 次／分鐘／來源且先限流不讀本文（`adapter.go:90` 先 `CheckCreate`）、標頭讀取上限約 20 KiB、`ReadHeaderTimeout` 5 s／`ReadTimeout` 15 s（`main.go:97-100`）。公開反而讓合法用戶端能遵守。
- **431 不帶安全標頭的風險可忽略**：本文是固定純文字、無標記、無使用者可控內容——即使被 MIME 嗅探為 HTML 也只是一行固定文字，沒有 XSS 面；頁面沒有可操作 UI，被嵌入框架也沒有點擊劫持標的；`Connection: close` 後連線即結束。HSTS 由邊緣送出（SCR-010／RR-008），瀏覽器快取 HSTS 政策，單一 431 缺少 HSTS 不使政策失效。
- **SCR-027 延伸到規格層**：17 號檔 SCR-027 的接受理由（net/http 直接產生、固定狀態文字、不洩漏資訊）同樣適用於 OpenAPI 0.1.1 的 `HeaderTooLarge` 與 API 說明 §11。前提是規格以「描述源站行為」而非「承諾／要求不帶標頭」的方式寫——初步意見指出的 SCR-043 已修正（`openapi.yaml:26`、`01-api-notes-v0.1.md:98`、`HeaderTooLarge` description）。將來在邊緣替 431 補上 nosniff／HSTS 不會「不符規格」。

### 3.3 問題 3：是否需要其他審查？——**不需另做完整安全審查**

- 最終 head 對 main 只動 6 個 `docs/` 檔（§2），未動程式碼、`.github`、`go.mod`、`go.sum`；文件描述與 main 程式一致（§2.3、§3.1）。
- 安全部只需對最終 head 做一次差異核對，即本文件。仍需審查部關卡（R1、R2 流程項）、1 個核准 review，並在合併的 head 上四項綠燈。

### 3.4 初步意見結論摘要（head `519e030`，2026-10-19 14:19 UTC+8）

初步意見已對三題給出相同答案（符合 10 號意見；公開上限與 431 可接受、SCR-027 可延伸；只動 docs、只需最終 head 差異核對），並提出三項資訊發現：SCR-043（「不帶安全標頭」寫成規格承諾）、SCR-044（標頭上限描述不精確：未含請求列、門檻應為範圍、漏列標頭值超過 500 個）、SCR-045（方法規格寫「HTTP 契約不變／不需改版」與 OpenAPI 0.1.1 新增 431 不一致）。另指出前提更正：main 沒有 413，本文超限回 400 `invalid_url`。

## 4. SCR-043～045 修正核對（`3c037a1`，最終 head 仍保留）

| 編號 | 要求 | 最終 head 內容 | 判斷 |
|------|------|----------------|------|
| SCR-043 | 「不帶安全標頭」改為描述性寫法，不得變成禁止加上 | `openapi.yaml:26`「源站由 net/http 直接回應，不保證帶有上述安全標頭；部署邊緣可能加上；用戶端不得依賴這些標頭的有無」；`HeaderTooLarge` description「不保證帶安全標頭」；`01-api-notes-v0.1.md:98` 同義；方法規格 §13.4「源站不保證帶安全標頭（部署邊緣可能加上）」 | ✅ 已修正。殘留：API 說明 v0.1.1 修訂列摘要寫「431 不帶，見 §11」，屬修訂紀錄的歷史摘要，不構成規格承諾，不另立發現 |
| SCR-044 | 上限為請求列＋標頭；16384 保證接受、16384～約 20480 為範圍；標頭值超過 500 個也觸發 | `openapi.yaml:20-23`、`01-api-notes-v0.1.md:96`、`HeaderTooLarge` description、方法規格 §13.4：請求列（含路徑與查詢字串）＋全部標頭合計 16384 B 以下保證接受；16384～約 20480 B 不保證；超過約 20480 B 或標頭值（含 Host）超過 500 個回 431；方法規格另註「HTTP/1.1 行為」 | ✅ 已修正，與安全部本機實測一致（§3.1） |
| SCR-045 | 方法規格與 OpenAPI 0.1.1 一致 | 方法規格 §13.4「標頭上限」列：「431 已在 OpenAPI 0.1.1 以 `HeaderTooLarge`（text/plain）標示，不屬四類錯誤；業務錯誤契約（`ErrorResponse` 四類）不變；不破壞相容（不升 v2）」；§13.4.1「不變」項同義並寫明 OpenAPI 0.1.1 補記本文上限 32768 B 與 431；v0.2.2 修訂列附 v0.2.1 列（已簽核歷史、保留原文）之更正註記 | ✅ 已修正；與 R4 一致 |

R4（審查部）：§13.4、§13.4.1、v0.2.2 修訂列均已改為「不破壞相容、OpenAPI 0.1.1 補記上限與 431」。`git diff -U0 433140c 44b099c -- 06-method-specs-v0.2.md` 的 hunk 只在 front matter（第 3～4、10 行）與 §13 之後（第 152、189、210、218、221、239 行；§13 起於第 149 行），**§1～§12 未動**。

## 5. 08-ruling 檔（`docs/02-architecture/08-ruling-G4-PR2-R1-go-mapping-v0.1.md`）

- 性質：設計部 2026-10-07 的 G4-PR2-R1 書面裁示（Go 語言對映不需 CR），既有檔。
- 本 PR 改動（+3／−2）：「最後更新」改 2026-10-19；§五 觀察 2（建立本文上限 8192）在建議欄末尾補「〔2026-10-19 後續註記：已由設計裁定 DR-G4-01 調高為 32768 bytes（PR #17，main `433140c`），見方法規格 v0.2.2 §13.4.1；本列保留裁示當時原文〕」；修訂紀錄加一列。
- 不含程式碼區塊或設定；沒有變更資安契約，也沒有削弱驗證、限流、錯誤不洩漏等機制（只是事實註記，32768 為已生效的 DR-G4-01）。
- 與 06-security 文件沒有衝突：與 10 號檔（附條件同意、C1～C5 已達成）、17 號檔（#17 確認）一致；安全部 08 號檔的 8192 屬歷史紀錄，PR 說明也寫明未動 `06-security/08-pr-002`。
- 沒有代安全部做出結論。

## 6. 其他核對

- R1／R2 的 SEC 編號（只核對是否屬實）：PR 描述「## 關聯」與 `3c037a1` 末行列 SEC-001、SEC-007、SEC-008，對照 01 號檔 v0.2：SEC-001（長網址驗證，含「長度 ≤ 文件化上限 L」）、SEC-007（速率限制）、SEC-008（安全回應標頭）**皆存在且與本 PR 內容相關**。另可考慮補列 SEC-003（錯誤回應不洩漏內部細節，對應「回應不帶上限值」），屬建議，不影響判定。PR 描述對 SEC-008 的註解仍寫「431 不帶安全標頭之例外」，那是 PR 描述而非 repo 文件，不影響規格；流程項目由審查部判定。
- API 說明 §5 新增的 CSP、X-Frame-Options、Referrer-Policy、Cache-Control 列，與 main `internal/httpapi/security_headers.go` 一致（`headerCSP = "default-src 'none'; frame-ancestors 'none'"`、`DENY`、`no-referrer`、JSON 回應 `no-store`）。
- 根層 `security: []`（`openapi.yaml:34`）只是明示無認證（OUT-03），沒有改變任何端點的實際行為；`x-max-body-bytes: 32768` 是規格擴充欄位，不出現在回應中。

## 7. 判定與合併條件

**判定：通過，僅適用於最終 head `44b099c227cac2ab24380f4346654540da6535de`。**

1. **核准綁定**：本核准只適用於 `44b099c227cac2ab24380f4346654540da6535de`。合併當下 PR head 必須等於這個 SHA。
2. **任何新 commit 都要重新送核**：只要 head 改變（包括同步 main 的 merge commit、只改措辭的 commit），本核准即不適用，新 head 須經安全部重新書面確認後才可合併。若只是「只帶進 main 改動的同步 merge commit」，確認範圍為：parents 是 `44b099c` 與當時的 main；`git show --remerge-diff` 為空；以下 blob 不變——`06-method-specs-v0.2.md` `f7e13313d3c56ba7a509ad22b1173766ba58ebd3`、`08-ruling-G4-PR2-R1-go-mapping-v0.1.md` `a6da63628a70e6afecf90f8a781ff67e14ba67fb`、`02-architecture/README.md` `203c23d62800e083d349ab917b37e359d0a85e12`、`01-api-notes-v0.1.md` `94c2688bd37bd90347cb87da4303094c71929da6`、`04-api/README.md` `e6dcdab72338d1501500a9b434305d370387c013`、`openapi.yaml` `0f69153fc8df8f1f3f3499e5522089241573d23b`；當時 main 若改過 `main.go`、`adapter.go`、`mapper.go`，須另核文件與程式的對照；新 head 四項 success 且每個 job 的 head_sha 相符。除此之外的任何變更都要重新做安全審查。
3. **strict=true 與第 1、2 點的關係**：目前 main＝`141fcda`，就是最終 head 的第二個 parent，PR 已與 main 同步，可在不變更 head 的情況下合併。若合併前 main 再前進，依第 2 點處理，不另成條件。
4. **其他關卡**：審查部關卡（R1、R2 流程項由審查部判定）＋1 個核准 review（ruleset 24629854；目前 `mergeable_state: blocked`）。ruleset 不會在 push 後撤銷核准，合併執行者須在按下合併前確認 head 仍是 `44b099c`。
5. 合併當下四項必要檢查在合併的 head 上都必須是 success。
6. SCR-046（資訊）不阻擋本 PR。

## 8. 發現項

| 編號 | 嚴重度 | 標題 | 狀態 |
|------|--------|------|------|
| SCR-043 | 資訊 | 431「不帶安全標頭」被寫成規格承諾 | **已修正（`3c037a1`），待 #20 合併後關閉** → **v0.1.1：已關閉**（#20 合併，`abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`） |
| SCR-044 | 資訊 | 標頭上限描述不精確（未含請求列、門檻應為範圍、漏列標頭值 500 個） | **已修正（`3c037a1`），待 #20 合併後關閉** → **v0.1.1：已關閉**（#20 合併，`abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`） |
| SCR-045 | 資訊 | 方法規格「HTTP 契約不變／不需改版」與 OpenAPI 0.1.1 新增 431 不一致 | **已修正（`3c037a1`），待 #20 合併後關閉** → **v0.1.1：已關閉**（#20 合併，`abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`） |
| SCR-046 | 資訊（不在本 PR 範圍） | 根目錄 `README.md` 的 431 描述仍是修正前寫法 | 開放（建議，不阻擋） |

SCR-043～045 的影響、重現與修補建議見初步意見（§3.4）；修正內容見 §4。

### SCR-046（資訊）根目錄 README 的 431 描述與 OpenAPI 0.1.1 不一致

- 位置：`README.md:56`、`README.md:101`（main `141fcda`，PR #19／ENG-021 合併）。
- 類別：文件一致性；對應 SEC-008、DR-G4-01。
- 影響：第 56 行寫「Request headers larger than about 20 KiB … 431」，沒有寫請求列也計入、16384～約 20480 B 是範圍、標頭值超過 500 個也觸發；第 101 行寫 net/http 直接回的 400／431「carry none of these security headers」，是對源站的事實描述，但沒有註明邊緣可能加上、用戶端不得依賴。與 PR #20 修正後的 OpenAPI 0.1.1／API 說明 §11 不一致，可能讓讀者以 README 為準。沒有安全影響。
- 重現：`git show 141fcda:README.md`，對照 `git show 44b099c:docs/04-api/01-api-notes-v0.1.md` §11。
- 修補建議：研發部（ENG-021）後續文件 PR 以 OpenAPI 0.1.1／API 說明 §11 為準改寫，或直接連結 §11；不必為此變更 PR #20 的 head。

## 9. 無法查證之處

- 審查部關卡（R1、R2）與核准 review 的狀態。
- 部署邊緣（Cloudflare Tunnel／反向代理）自己的標頭上限，以及它是否會替 431 補標頭。
- 源站若日後啟用 HTTP/2（例如 h2c），431 改由 http2 程式碼路徑產生，格式與門檻可能不同；目前 `main.go:70` 以明文 `ListenAndServe` 啟動、未設定 `Protocols`／h2c，只有 HTTP/1.1；方法規格也已註明「HTTP/1.1 行為」。
- 研發部書面同意分支名例外的紀錄（`gate-G4-PR20-2026-10-19.md`）不在 repo 內，屬審查部流程範圍，未查證。

## 10. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-19 | 安全部 | 初版：回應審查部 G4-PR20-R3 三個問題（符合 10 號意見；公開上限與 431 源站不帶安全標頭可接受、SCR-027 延伸到規格層；只動 docs、只需最終 head 差異核對）；收入初步意見（head `519e030`）結論摘要；查證最終 head `44b099c`（parents `3c037a1`＋main `141fcda`，#19 於 14:49:27 UTC+8 合併）、remerge-diff 為空、只動 6 個 docs 檔（08-ruling 為修改而非新增，協調者宣稱不正確）、run 39212359855 四項 success 且 head_sha 相符、gosec Issues 0；main `433140c → 141fcda` 未改 main.go／adapter.go／mapper.go；SCR-043～045 已修正、待合併後關閉；新增 SCR-046（資訊）；判定通過並綁定完整 head SHA。小版修訂不改檔名（檔名保留主版號） |
| v0.1.1 | 2026-10-19 | 安全部 | 合併後狀態（§0）：#20 於 2026-10-19 15:15:58（UTC+8）合併，merge commit `abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`（parents `141fcda`＋核准 head `44b099c`，`git diff 44b099c abbec9e` 為空）；SCR-043～045 關閉；SCR-046 仍開放；其他內容不變。小版修訂不改檔名（檔名保留主版號） |
