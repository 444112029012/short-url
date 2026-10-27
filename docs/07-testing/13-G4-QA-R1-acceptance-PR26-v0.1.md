# G4-QA-R1 驗收：PR #26 契約測試 CT-*

版本：v0.1.1
日期：2026-10-22（UTC+8）
作者：品保部
對象：PR #26 `test: add OpenAPI contract tests CT-CREATE/REDIR/STATS/ERR`，head `e224253de29743a956e8c952cafc4eff8acff1d7`，base `363f1692`（= 目前 origin/main）
依據：`07-test-cases-v0.1.md` §CT、`02-traceability-matrix-v0.1.md`、`12-gate-G4-qa-precheck-v0.1.md`（G4-QA-R1、DEF-001）、`docs/04-api/openapi.yaml`

## 結論：**有條件關閉**

6 個 CT 都已實作、本機實跑全綠，突變驗證 3／3 會讓對應 CT 變紅，行為面有涵蓋 CT 定義。但 DEF-001 的預期結果是「用 `openapi.yaml` 驗證回應」，PR 的做法是把期望值**寫死成常數**（註解標行號），沒有讀檔。這比較接近 G4-QA-R1 的選項 (b)（替代方案＋登錄殘餘＋期限），不是完全符合選項 (a)。
關閉條件（見 §7）：**C1** 研發在 PR 描述或 ENG-013 補一段書面例外說明「寫死常數取代解析 openapi.yaml」的理由，並同意 DEF-005 的期限；**C2** DEF-005／DEF-006 登錄並排期限（G5 前）。兩項都不用改 PR #26 的程式。C1 完成後，DEF-001 改為「已修正」。

## 1. PR 與 CI 查證（證據：`evidence/pr26-gh-ci.txt`、`evidence/pr26-gh-pr.json`）

| 項 | 研發說法 | 品保實查 | 判定 |
|---|---|---|---|
| head SHA | e224253… | `gh pr view` headRefOid＝`e224253de29743a956e8c952cafc4eff8acff1d7` | ✅ |
| base | 363f169 | `363f1692e5…`；`git merge-base origin/main e224253`＝同值，origin/main 也是這個值（strict 規則要求的「與 main 同步」成立） | ✅ |
| 變更檔 | 只有 contract_test.go | 1 個檔，ADDED，+258／−0 | ✅ |
| CI run 39704609939 | 四項全綠 | headSha＝e224253…、event＝pull_request、conclusion＝success；sast-gosec、secrets-gitleaks、unit-test、sca-govulncheck 全部 success（完成時間 2026-10-22 09:15～09:16 UTC+8） | ✅ |
| required checks | 四項 | 傳統 branch protection 回 404（沒設定）；實際上由 ruleset 24629805 要求 `sca-govulncheck`、`sast-gosec`、`secrets-gitleaks`、`unit-test`，strict＝true；四項名稱跟這次 run 一一對應 | ✅ |
| 只用標準庫 | 是 | import 只有 net/http、httptest、url、filepath、regexp、sort、strings、testing，加上專案內的 observability、ratelimit；go.mod／go.sum 沒有改 | ✅ |
| CT-ERR-01 有 8 個子測試 | 是 | `-v` 輸出 8 個子測試 PASS | ✅ |

## 2. CT 逐案對照

| CT | 定義（07-test-cases） | 實作斷言 | 對照 openapi.yaml | 判定 |
|---|---|---|---|---|
| CT-CREATE-01 | 201＋CreateUrlResponse（short_code pattern、必填） | 201；Content-Type application/json；nosniff；欄位**剛好** {long_url, short_code, short_url}（additionalProperties:false）；short_code 符合 `^[A-Za-z0-9]{8}$`；short_url／long_url 是絕對 URI；short_url＝base＋code；long_url 原值回傳、≤2048 | L69、L248–L261 吻合 | ✅ |
| CT-CREATE-02 | 400＋type∈四碼且＝invalid_url | 空字串、ftp、非 URL、超長 2049＋ 4 種輸入都是 400／invalid_url；ErrorResponse 兩層欄位剛好；type 在 enum 裡；message 非空；沒有 Location | L79、L276–L292 吻合 | ✅ |
| CT-REDIR-01 | 302＋Location（uri） | 302；Location 是絕對 URI 且等於存起來的 long_url（帶 `?url=` 注入參數也不受影響，SEC-014）；nosniff | L123–L130、L13 吻合 | ✅ |
| CT-REDIR-02 | 404＋not_found 形狀 | 404；完整 ErrorResponse 斷言；type＝not_found | L142 吻合 | ✅ |
| CT-STATS-01 | 200＋UrlStatsResponse（click_count≥0） | 200；欄位剛好 {click_count, short_code}；pattern；click_count 是整數、≥0，導向一次後＝1 | L180、L265–L273 吻合 | ✅ |
| CT-ERR-01 | 各錯誤 type 只有四值、無多餘欄位 | 8 個子測試：invalid_url×3（create／redirect／stats）、not_found×2（redirect／stats）、rate_limited×2（create／redirect）、internal_error×1（create） | 四碼都有測到；但 openapi 也為 redirect 與 stats 定義了 500，這兩條沒測（DEF-006） | ⚠️ 涵蓋定義；路徑不完整 |

**「以 openapi.yaml 為準」查證**：
- 所有測試都**沒有讀** `openapi.yaml`（`rg openapi --glob '*_test.go'` 只找到這個檔的註解）。期望值寫死在常數裡。
- 品保逐一核對 29 處行號註解（L11、L13、L63、L69、L72、L79、L82、L101、L123、L126、L130、L142、L152、L180、L228、L248、L249、L253、L257、L259、L261、L265、L266、L270、L272、L273、L276、L278、L284、L288、L292）：**全部對得上** head 版 openapi.yaml，常數值目前也一致。head 版的 openapi.yaml 跟 `/workspace/projects/short-url/docs/04-api/openapi.yaml` 完全相同（diff 沒有差異）。
- 結論：**目前一致，但沒有自動防漂移機制** → DEF-005。

## 3. 執行證據（2026-10-22 UTC+8，box go1.27.3，工作目錄 `/workspace/tmp/qa-pr26` @ e224253）

| 指令 | 結果 | 證據 |
|---|---|---|
| `go test ./internal/httpapi/ -run CT -v -count=1` | 頂層 6 個＋子測試 8 個，共 **14 PASS／0 FAIL**，exit 0 | `evidence/pr26-ct-run.txt` |
| `go test ./... -race -count=1` | 10 個套件 ok（ports 沒有測試檔），exit 0 | `evidence/pr26-race.txt` |
| `go test ./... -v` 計數 | **181 PASS／0 FAIL／0 SKIP**（G4 基準 167＋這次 14） | 本次實跑 |
| `go vet ./...` | 沒有輸出，exit 0 | `evidence/pr26-vet.txt` |
| 覆蓋率 | 總計 **78.7%**；httpapi 90.4%；只跑 CT 時 httpapi 82.2%；cmd/shorturl 2.8%（DEF-003 既有） | `evidence/pr26-cover.txt` |

## 4. 突變驗證（只在 `/workspace/tmp/qa-pr26-mut`；證據：`evidence/pr26-mutation.txt`）

| # | 突變 | 預期 | 實際 | 判定 |
|---|---|---|---|---|
| M1 | `mapper.go` L53 not_found 404→410 | CT-REDIR-02、CT-ERR-01 not_found 變紅 | REDIR-02 失敗；ERR-01/not_found/redirect、/stats 失敗（`status = 410, want 404`） | ✅ 抓到 |
| M2 | `mapper.go` L55 rate_limited→internal_error（狀態碼維持 429） | CT-ERR-01 rate_limited 變紅 | ERR-01/rate_limited/create、/redirect 失敗（`error.type = "internal_error", want "rate_limited"`） | ✅ 抓到 |
| M3 | `adapter.go` L135 302→301 | CT-REDIR-01 變紅 | `status = 301, want 302` | ✅ 抓到 |
| 還原 | `git checkout -- .` | 全綠 | `git status` 乾淨；CT 全部 ok，exit 0；原本的 worktree `git status` 也乾淨 | ✅ |

這可以補「事後才補、沒有紅燈階段」的證據：CT 對狀態碼和錯誤碼的變動有偵測能力。

## 5. 研發殘餘風險評估

| # | 研發自記 | 品保查證 | 可接受？ | 處置 |
|---|---|---|---|---|
| R-a | 事後補，沒有紅燈階段 | §4 突變 3／3 都抓到 | ✅ 可接受 | 不另外登錄；本文 §4 作為證據 |
| R-b | stats 429 沒測 | openapi.yaml 的 getUrlStats **沒有定義 429**（只有 200／400／404／431／500），所以不是契約缺口；本身也不在 CT 範圍內 | ✅ 可接受（不算風險） | 不登錄。如果之後 stats 要加限流，要先改 openapi.yaml 再補 CT |
| R-c | 431 不在範圍 | openapi 三個操作都有 431（HeaderTooLarge），但 CT 定義沒有列 431；431 由 ENG-007／bodylimit 相關測試涵蓋 | ✅ 可接受 | 不登錄；建議 S-2 |
| R-d | 常數可能跟 openapi.yaml 漂移 | 屬實；目前一致（§2），但沒有自動防護；也跟 DEF-001 的「用 openapi.yaml 驗證」不同 | ⚠️ 有條件接受 | 登錄 **DEF-005**，期限 G5 前 |

## 6. 新缺陷

| ID | 嚴重度 | 標題 | 環境 | 重現步驟 | 預期結果 | 實際結果 | 狀態 |
|---|---|---|---|---|---|---|---|
| DEF-005 | 輕微 | CT 期望值寫死，沒有跟 openapi.yaml 自動同步 | PR #26 `e224253` | 1. 在臨時目錄把 `docs/04-api/openapi.yaml` L288 enum 加一個值，或 L248 required 加一個欄位 2. `go test ./internal/httpapi/ -run CT` | CT 失敗，或至少有守門測試指出常數跟 openapi.yaml 不一致 | 從程式碼可以推知仍會 PASS（測試不讀檔；本次沒有實際執行這個突變） | 開放；期限 G5 前。建議做法：stdlib 讀 openapi.yaml，用字串或正則比對關鍵行（enum、required、pattern、狀態碼），不必引入第三方套件 |
| DEF-006 | 輕微 | CT-ERR-01 的 internal_error 只測 create | 同上 | 1. 看 contract_test.go 的 CT-ERR-01 子測試 2. 比對 openapi.yaml 中 redirectShortCode 與 getUrlStats 的 '500' | redirect／stats 的 500 也用 ErrorResponse＋internal_error 斷言 | 只有 internal_error/create | 開放；期限 G5 前（可以沿用關閉 store 的手法） |

沒有致命或嚴重缺陷。

## 7. 必改與建議

| ID | 類型 | 內容 | 負責 |
|---|---|---|---|
| C1 | 關閉條件 | 研發補書面例外：用寫死常數＋行號註解取代讀檔驗證，並說明理由（安全決策：不引入第三方驗證器） | 研發部 |
| C2 | 關閉條件 | DEF-005、DEF-006 登錄並排期限（G5 前） | 研發部／品保部 |
| S-1 | 建議 | PR 描述補「TC IDs」欄（DoD D4／G4-QA-S1） | 研發部 |
| S-2 | 建議 | 之後如果擴充 CT，把 431 HeaderTooLarge 納入契約斷言 | 研發部 |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | PR #26 驗收：有條件關閉 G4-QA-R1 |
| v0.1.1 | 2026-10-27 | 品保部 | S2（G4 整體 S 項）：sha256 變動說明。gate-G4-PR26 初審引用之 `e688bd6d48790093c4f2ed427ac4c68d99970651f4e3855baeb6ad8e5b86157d` 為初版值；其後補「附錄：關閉條件複核」使 sha256 變為 `2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c`（亦為 PR #27 合併進 main `2425129b` 之版本）。**2627c2da 是本次修訂前的值，本次修訂後的新值以品保回報為準**（本檔不自我引用）。本次只改版本欄與本列，內容不變 |

## 附錄：關閉條件複核（2026-10-22 09:5x，Asia/Taipei，品保部）

- **C1：關閉 ✅**。書面依據為 PR #26 描述（以 `gh pr view 26` 實查）：「安全部 SCA 不同意引入 kin-openapi，所以以標準函式庫手寫斷言」，以及「常數與 openapi.yaml 為人工對應…行號註解協助比對」。品保先前已逐一核對 31 處行號註解，全部對得上。PR 描述沒有寫到「殘餘缺口由 DEF-005 補上」，這一點由 12 號的 DEF-005 登錄承接，不影響 C1 關閉。
- **C2：關閉 ✅**。DEF-005、DEF-006 都已登錄在 `12-gate-G4-qa-precheck-v0.1.md` 缺陷表（輕微，期限 G5 前）。
- G4-QA-R1 最終判定：**關閉**；DEF-001 改為關閉。合併後品保會在 main 上重跑 CT，作為確認。
