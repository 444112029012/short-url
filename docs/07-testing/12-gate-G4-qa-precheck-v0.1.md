---
文件：G4 品保前置結論（已合併 PR #13～#20 實作盤點）
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-22
對應需求：REQ-001～012、NFR-001～008、SEC-001～017；ENG-022（CR-003）
盤點基準：main `db51cd5`（#21 合併後；含 #13、#15～#20；#14 未合併、已關閉）
專案代號：SHORTURL
關卡：G4（PR 合併）—品保前置
---

# G4 品保前置結論：PR #13～#20

## 0. 結論

**有條件可送審**：建置、測試、CI、逐 PR 審查與安全結論都有證據且全綠；但 6 個契約案例（CT-*）完全沒有實作（DEF-001），且 S6／追溯 v0.4 的文件更新還沒進 main。必改 G4-QA-R1～R2 關閉前，品保不出具「通過」。

> 本文只是品保前置，不取代審查部的 G4 結論。所有「通過」都附本次實跑或可查證的紀錄；沒有證據的項目一律標「未驗證」。

## 1. 盤點範圍

| PR | 標題（摘要） | 類型 | 合併時間（UTC+8） | PR checks（4 項） | 審查部 G4 紀錄（最新） | 安全結論 |
|---|---|---|---|---|---|---|
| #13 | G3-R1 落地證據（ENG-020） | 文件 | 2026-10-19 09:05 | 全 success | `gate-G4-PR13-2026-10-07.md` | 安全 README 已引用 |
| #14 | 釘回 Go 1.27.1 | — | **未合併**（sca-govulncheck FAILURE） | — | — | 不在盤點內 |
| #15 | CR-002 Go 1.27.2＋gosec 釘版（ENG-020） | CI | 10-19 10:22 | 全 success | PR15 recheck-01：通過 | 14 號文件 |
| #16 | 安全標頭 E2E（ENG-016） | 測試 | 10-19 13:17 | 全 success | PR16：通過 | 審查部判定不需安全書面結論（只改測試） |
| #17 | 建立本文上限 32768（ENG-007） | 程式 | 10-19 14:00 | 全 success | PR17 recheck-02：通過 | 17 號文件 |
| #18 | 限流 E2E（ENG-011） | 測試 | 10-19 14:24 | 全 success | PR18 recheck-01：通過（S6 交品保） | 18 號文件 |
| #19 | README（ENG-021） | 文件 | 10-19 14:49 | 全 success | PR19 recheck-02：通過 | 不需 |
| #20 | OpenAPI 本文上限／431（ENG-007、DR-G4-01） | 文件／契約 | 10-19 15:15 | 全 success | PR20 recheck-01：通過 | 19 號文件 |

PR checks 來源：`gh pr view <n> --json statusCheckRollup`（2026-10-22 查詢）。main `db51cd5` 的 ci-devsecops run：success（建立於 2026-10-22 08:57 UTC+8）。

## 2. 測試實跑結果（2026-10-22，box，go1.27.3 linux/amd64，detached worktree @ `db51cd5`）

| 指令 | 結果 |
|---|---|
| `go vet ./...` | exit 0 |
| `go test ./... -cover -count=1` | 10 個套件 ok、1 個無測試（ports）；0 FAIL |
| `go test ./... -count=1 -json` 統計 | **167 PASS／0 FAIL／0 SKIP**（頂層 100、子測試 67） |
| `go test ./... -race -count=1` | 全部 ok |
| `go test ./internal/integration/ -run TestENG011 -v` | **19 PASS／0 FAIL／0 SKIP** |
| `CI=true … -run TestENG011_PerSource` | PASS（未 skip） |
| 全專案語句覆蓋（`-coverpkg=./...`） | **80.9%** |

各套件覆蓋（`-cover`，套件自身）：application 82.0%｜config 84.2%｜domain 81.1%｜health 100%｜httpapi 90.4%｜infra/sqlite 69.9%｜observability 87.5%｜ratelimit 87.3%｜cmd/shorturl 2.8%｜integration（只有測試）—。

原始輸出：`evidence/g4-full-run-2026-10-22.txt`、`evidence/s6-eng011-run-2026-10-22.txt`。
註：本機是 Go 1.27.3，CI 釘 1.27.2（CR-002）；CI 證據以 main run 為準。

## 3. TC 實作覆蓋（依案例目錄 v0.2，明細見矩陣 v0.4 §7）

| 類別 | 目錄數 | 已實作且通過 | 未實作 | 備註 |
|---|---:|---:|---:|---|
| UT | 34 | 34 | 0 | 5 個未在名稱標 TC ID（UT-HC-02／03、UT-HR-02／03、UT-HS-02），靠矩陣 §7 追溯 |
| IT | 6 | 6 | 0 | IT-FLOW-03 只在單元層實作（DEF-002） |
| CT | 6 | **0** | **6** | DEF-001 |
| E2E（含 E2E-05-01～11） | 16 | 16 | 0 | |
| CHK | 18 | —（G5 舉證） | — | 不屬 G4；CHK-ENG-022 待 CR-003 |
| **自動化 TC 合計** | **62** | **56（90.3%）** | **6** | |

公開方法單測計畫：19 個方法都有通過的 UT 或 IT（矩陣 §4、§7）→ G4 條件「單元測試涵蓋公開方法正常 I/O、邊界、例外且通過」：**達成**。

## 4. 追溯（需求→設計→測試→程式）

| 段 | 狀態 | 證據／缺口 |
|---|---|---|
| 需求→設計 | ✅ | 矩陣 §1～3 設計欄 37／37 已填 |
| 設計→TC | ✅ | 矩陣 TC 欄 37／37 非空 |
| TC→程式（測試實作） | ⚠️ | CT-* 6 個無實作；5 個 UT 無 TC 標記（已在 §7 人工補對應） |
| 任務／PR→需求 | ✅ | #13～#20 描述都帶 REQ／NFR／SEC（#13、#15、#16 為 NFR-004／SEC-008～011） |
| ENG→TC | ✅ | 任務對照 v0.3：17／17 驗收列；ENG-022 已列 |
| 文件已進 main | ❌ | 矩陣 v0.4、案例 v0.2、任務對照 v0.3、S6 處置與本文件目前只在 `/workspace/projects/short-url/docs/`，**尚未提 PR**（G4-QA-R2） |

## 5. G4 過關條件 × 品保觀點

| G4 條件（原文） | 判定 | 依據 |
|---|---|---|
| PR 標示對應 REQ／SEC 編號 | ✅ | §1、§4 |
| 公開方法簽名與方法規格一致（差異附 CR） | ✅（採審查部結論） | 各 PR 審查紀錄；品保未重做 |
| 單元測試涵蓋正常 I/O、邊界、規格例外且通過 | ✅ | §2、§3 |
| SAST／SCA／機密掃描無超門檻 | ✅ | 各 PR 與 main 的 4 項 check success |
| 安全審查清單完成、阻擋項已修 | ✅ | 06-security 14、17、18、19 號；SCR-030／032／033 為開放建議，不阻擋 |
| 介面變更同步 OpenAPI 與相關文件 | ⚠️ | #20 已同步 OpenAPI；但沒有任何 CT 自動比對 OpenAPI（DEF-001），無法自動證明一致 |
| 符合 commit 規範與 DoD | ⚠️ | DoD P2（「HTTP 對齊 OpenAPI」驗證方式＝CT／E2E）只有 E2E 部分；D2（矩陣與案例更新）本文件產出後才達成，待進 main |

## 6. DoD 達成（`08-definition-of-done`，針對 #13～#20 合計）

| 項 | 判定 | 說明 |
|---|---|---|
| P1、P3～P7 | ✅ | UT／IT／E2E 全綠；P7 由 E2E-05-01～11 證明 |
| P2 | ⚠️ | 缺 CT（DEF-001） |
| P8 | ✅ | 無已知 Blocker |
| T1～T5 | ✅（T2 依 v0.4） | |
| D1 | ✅ | #20 |
| D2 | ⚠️ | 待 R2 |
| D3 | ✅ | 07-testing 內簡體字「举证」「绿」已改為「舉證」「綠」（本次） |
| D4 | ✅ | PR 描述有 ENG／REQ／SEC；但 TC ID 只有 #20（UT-03）有寫，列建議 S1 |
| S1～S5 | ✅ | 全綠＋安全文件 |
| C1 | ✅ | strict=true 已生效（D-10） |

## 7. 缺陷清單

| ID | 嚴重度 | 標題 | 環境 | 對應 | 重現步驟 | 預期結果 | 實際結果 | 狀態 |
|---|---|---|---|---|---|---|---|---|
| DEF-001 | 一般 | 契約案例 CT-* 沒有實作 | main `db51cd5` | REQ-001／002／004／007／009；SEC-001、003、014；DoD P2；CT-CREATE-01～02、CT-REDIR-01～02、CT-STATS-01、CT-ERR-01 | 1. `git checkout db51cd5` 2. `rg -n 'CT-|openapi' --glob '*_test.go'` 3. 看 `go test ./... -v` 的測試名稱 | 6 個 CT 各有自動測試，用 `docs/04-api/openapi.yaml` 驗證回應（狀態碼、schema、ErrorResponse type 只有 4 種值、無多餘欄位） | 沒有結果，也沒有任何測試讀 OpenAPI | **有條件關閉**（2026-10-22 附註：PR #26 `e224253` 已補 6 個 CT，實跑 14 PASS，突變 3／3 抓到；條件是研發補書面例外，並登錄 DEF-005／DEF-006。詳見 `13-G4-QA-R1-acceptance-PR26-v0.1.md`）。10-22 複核：C1／C2 均已關閉，DEF-001 **關閉**（main 合併後重跑確認） |
| DEF-002 | 輕微 | IT-FLOW-03 只在單元層 | 同上 | REQ-012；IT-FLOW-03 | 1. `rg -n 'IT-FLOW-03\|DuplicateLongURL' --glob '*_test.go'` | 整合層（真實 SQLite＋HTTP）建立同一 long_url 兩次，得到兩個不同短碼且都能導向 | 只有 `application/TestCreateShortURLAllowsDuplicateLongURL`（假 store），不驗導向 | 開放 |
| DEF-003 | 輕微 | 整合／E2E 測試重複組裝，沒有執行 `cmd/shorturl` `run()` | 同上 | ENG-017；E2E-*；NFR-002 | 1. `go test ./cmd/shorturl -cover` → 2.8% 2. 比對 `internal/integration/*_test.go` 的 `newRateLimitServer`／`newServer` 與 `main.go` `run()` | 測試與正式程式共用同一個組裝函式，`run()` 的組裝路徑有測試覆蓋 | 組裝碼重複兩份；目前逐行一致（2026-10-22 讀碼），但未來改 `run()` 時測試不會發現 | 開放 |
| DEF-004 | 輕微 | 07-testing 文件含簡體字 | docs | DoD D3 | `rg -n '举证\|绿' docs/07-testing` | 0 筆 | 原有 7 處（04×2、05、06、10、02、README） | **已修正**（本次；待 R2 進 main） |
| DEF-005 | 輕微 | CT 期望值寫死常數，未讀 `openapi.yaml`，規格變更時 CT 不會失敗 | internal/httpapi/contract_test.go | G4-QA-R1 | 改 openapi.yaml 某 enum／狀態碼後跑 `go test ./internal/httpapi -run CT` | CT 失敗 | 審查部實測：openapi.yaml 改三處（L69 201→200、L248 required 加 created_at、L288 enum 加 gone）後 CT 仍 14 PASS（gate-G4-PR26-2026-10-22.md §四.3） | 開放；G5 前（見 13 號） |
| DEF-006 | 輕微 | 導向與統計的 500 internal_error 未有 CT | internal/httpapi/contract_test.go | G4-QA-R1 | 檢視 CT-ERR-01 子測試 | 含 redirect／stats 500 | 未涵蓋 | 開放；G5 前（見 13 號） |

嚴重度定義依品保規範：致命／嚴重／一般／輕微。**沒有致命或嚴重缺陷。**

## 8. 必改項

| ID | 條件 | 問題 | 必改內容 | 負責 |
|---|---|---|---|---|
| G4-QA-R1 | 介面同步 OpenAPI；DoD P2 | DEF-001：6 個 CT 沒有實作 | 二擇一：(a) 研發補 CT-CREATE-01～02、CT-REDIR-01～02、CT-STATS-01、CT-ERR-01，以 openapi.yaml 為準（零預算：可用 Go 開源 OpenAPI 驗證套件，經 SCA）；或 (b) 研發提書面例外，品保同意以 E2E＋handler 測試替代，並登錄殘餘與 G5 前補齊期限 | 研發部（品保驗收） |
| G4-QA-R2 | 可追溯性；DoD D2／D3 | 矩陣 v0.4、案例 v0.2、任務對照 v0.3、S6 處置、本文件、README v0.4 與簡體字修正都還沒進 main | 以純文件 PR 提交 `docs/07-testing/` 變更（分支命名依 CR-003 核准結果，或用 ENG-011／ENG-022） | 品保部／總協調 |

## 9. 建議項（不阻擋）

| ID | 建議 |
|---|---|
| G4-QA-S1 | 未標記的 5 個 UT 和 TestENG011_* 在名稱或註解加 TC ID（如 `// TC: E2E-05-03`）；PR 描述補「TC IDs」欄（DoD D4） |
| G4-QA-S2 | 補 IT-FLOW-03 整合層測試（DEF-002） |
| G4-QA-S3 | 把 `run()` 的組裝抽成可匯出的建構函式，供 integration 測試共用（DEF-003），提高 cmd/shorturl 覆蓋 |
| G4-QA-S4 | 安全部開放建議 SCR-030／032／033 在 G5 前收斂；SCR-030 需要時另開 TC |
| G4-QA-S5 | G5 前補齊 CHK-*（NFR-001 基線、SEC-008 演示 HTTPS 擷取、SEC-012／015／017） |
| G4-QA-S6 | CR-003 核准、ENG-022 落地後，品保依 CHK-ENG-022 驗收；之前別用 DR 分支名提含程式的 PR |
| G4-QA-S7 | infra/sqlite 覆蓋 69.9% 是各套件最低，補錯誤路徑測試 |

## 10. 交接單

- 交出 → 接收：品保部 → 總協調（並請審查部引用）
- 日期：2026-10-22（UTC+8）；版本 v0.1；狀態：審查中
- 結論：**有條件可送審**；R1、R2 關閉後品保複核改為「可送審」
- 已知問題：DEF-001 已關閉（13 號，C1／C2 已複核）；DEF-002、003、005、006 開放；DEF-004 已修正

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | G4 品保前置初版：#13～#20 盤點、實跑 167 PASS、80.9% 覆蓋、DEF-001～004、R1～R2、S1～S7 |
