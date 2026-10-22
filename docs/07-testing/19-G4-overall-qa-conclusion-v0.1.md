---
文件：G4 整體品保結論（PR #13～#26）＋G4-QA-R2 處置
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-22（Asia/Taipei，UTC+8）
基準：main `f646beeb25eec7a20b358df0af88450c76a85af5`；CI run 39737480926
回應：審查部 `docs/gates/gate-G4-overall-2026-10-22.md`（sha256 `01648257a04a43be4f01738db304d74d12e4aeea8b434a20ef8eeb93f7cb9446`）G4-R2
對應需求：REQ-001～012、NFR-001～008、SEC-001～017
取代：`12-gate-G4-qa-precheck-v0.1.md`（#13～#20 前置，基準 db51cd5）之結論；12 號缺陷表沿用
專案代號：SHORTURL
效力：main 有新 commit 或下列證據檔 sha256 改變時，本結論失效
---

# G4 整體品保結論（PR #13～#26）

## 0. 結論

**品保觀點：有條件通過。** main `f646beeb` 上建置、測試、CT、CI 全部重跑全綠（181 PASS／0 FAIL／0 SKIP；CT 14 PASS；`-race` 全綠；覆蓋 80.9%），#13～#26 各 PR 四項 checks 全 success，無致命／嚴重缺陷。

條件（未滿足前本結論不轉為「通過」）：
1. **G4-QA-R2 尚未關閉**：本次只提出處置方式（§6），須審查部接受且使用者決定入庫方式；本結論**不視為已獲推送或開 PR 核准**。
2. 開放缺陷 DEF-002、DEF-003、DEF-005、DEF-006（皆輕微）維持 G5 前處置，不阻擋 G4。
3. 本結論只涵蓋品保專業範圍；G4-R1（安全整體結論）、G4-R4（#24 bypass 追認、RR-011 接受）、G4-R5（D-15／CR-003 生效）不由品保判定，品保相關判定依賴其結果（見 §5）。

> 本文為審查部前置之品保書面結論，不取代審查部 G4 結論。所有「通過」均附本次重跑或 gh 唯讀查詢證據。

## 1. 基準查證（gh 唯讀，2026-10-22 15:09～15:12 UTC+8）

| 項 | 回報 | 品保實查 | 判定 |
|---|---|---|---|
| main HEAD | f646bee | `gh api repos/444112029012/short-url/commits/main` → `f646beeb25eec7a20b358df0af88450c76a85af5` | ✅ |
| 開著的 PR | 無 | `gh pr list --state open` → `[]` | ✅ |
| run 39737480926 | 四項通過 | event push、headSha f646beeb…、conclusion success；sast-gosec、secrets-gitleaks、unit-test、sca-govulncheck 皆 success；完成 15:01:49（UTC+8） | ✅ |

## 2. PR 盤點（#13～#26）

| PR | 標題摘要 | 分支 | merge commit | 合併（UTC+8） | PR checks | 審查部最新紀錄 |
|---|---|---|---|---|---|---|
| #13 | G3-R1 落地證據（ENG-020） | docs/ENG-020-… | 2e54f7e | 10-19 09:05 | 4／4 success | PR13：通過 |
| #14 | 釘回 Go 1.27.1 | — | — | **CLOSED 未合併** | — | 不在範圍 |
| #15 | CR-002 Go 1.27.2＋gosec | fix/ENG-020-… | c341be1 | 10-19 10:22 | 4／4 | PR15 recheck-01：通過 |
| #16 | 安全標頭 E2E（ENG-016） | test/ENG-016-… | ed62039 | 10-19 13:17 | 4／4 | PR16：通過 |
| #17 | 本文上限 32768（ENG-007） | feat/ENG-007-… | 433140c | 10-19 14:00 | 4／4 | PR17 recheck-02：通過 |
| #18 | 限流 E2E（ENG-011） | test/ENG-011-… | 643bf8a | 10-19 14:24 | 4／4 | PR18 recheck-01：通過 |
| #19 | README（ENG-021） | docs/ENG-021-readme | 141fcda | 10-19 14:49 | 4／4 | PR19 recheck-02：通過 |
| #20 | OpenAPI 本文上限／431（ENG-007、DR-G4-01） | docs/G4-DR-G4-01-… | abbec9e | 10-19 15:15 | 4／4 | PR20 recheck-01：通過 |
| #21 | G4 第二批安全審查文件（ENG-011） | docs/ENG-011-… | db51cd5 | 10-22 08:57 | 4／4 | PR21：通過 |
| #22 | S16 裁示與 SCR-041（ENG-007、DR-G4-01） | docs/G4-DR-G4-01-… | 363f169 | 10-22 09:08 | 4／4 | PR22：通過 |
| #23 | 安全審查 08 v0.1.5（ENG-011） | docs/ENG-011-… | 2a02819 | 10-22 12:09 | 4／4 | PR23 recheck-01：通過 |
| #24 | README 431 措辭（ENG-021） | docs/ENG-021-readme-followup | 5084d6a | 10-22 13:59 | 4／4 | PR24 recheck-02：通過；admin bypass（G4-R4 待追認） |
| #25 | CI 強化 SCR-032／019／021／004（ENG-020） | **fix/ENG-022-ci-hardening-scr032** | d0749f5 | 10-22 14:29 | 4／4 | PR25 recheck-01：通過（分支名偏差列 S1） |
| #26 | OpenAPI 契約測試 CT-*（ENG-013） | test/ENG-013-… | f646bee | 10-22 14:59 | 4／4 | PR26 recheck-01：通過；綁定 head 49e29d8 |

來源：`gh pr list --state all --json number,state,headRefOid,mergeCommit,mergedAt,title`、`gh pr view <n> --json statusCheckRollup,headRefName`。合併 head 與審查部綁定 head 之吻合，採審查部 G4 整體紀錄 §二 表（品保未逐筆重算第二 parent）。

程式／測試變更範圍（相對 #20 後）：#26 只新增 `internal/httpapi/contract_test.go`（+258），內容與 13 號驗收之 head `e224253` 相同（`git diff e224253 49e29d8 -- internal` 無差異）；#25 只改 `.github/workflows/ci-devsecops.yml`；#21～#24 純文件。

## 3. main `f646beeb` 重跑（box，go1.27.3；證據 `evidence/g4-overall-f646-*`）

| 指令 | 結果 |
|---|---|
| `go vet ./...` | exit 0，無輸出 |
| `go test ./... -count=1 -json` | **181 PASS／0 FAIL／0 SKIP**（頂層 106＋子測試 75）；＝12 號 167＋CT 14 |
| `go test ./... -race -count=1` | 全部 ok，exit 0 |
| `go test ./internal/httpapi/ -run CT -v -count=1`（**留待之 CT 重跑**） | **14 PASS／0 FAIL**（頂層 6＋CT-ERR-01 子測試 8），exit 0 → 13 號「合併後在 main 重跑 CT」**完成** |
| `go test ./internal/integration/ -run TestENG011 -v` | 19 PASS |
| 覆蓋 `-coverpkg=./...` | **80.9%**；套件：application 82.0、config 84.2、domain 81.1、health 100、httpapi 90.4、infra/sqlite 69.9、observability 87.5、ratelimit 87.3、cmd/shorturl 2.8 |
| DEF-005 複驗（暫存複本，已刪） | openapi.yaml L288 enum 加 `gone` → CT 仍 14 PASS → **DEF-005 仍開放** |

CI 以 run 39737480926（Go 1.27.2）為準；本機版本差異不影響結論。

## 4. 缺陷狀態（截至 f646bee）

| ID | 嚴重度 | 狀態 | 本次依據 |
|---|---|---|---|
| DEF-001 | 一般 | **關閉**（main 重跑確認） | §3 CT 14 PASS |
| DEF-002 | 輕微 | 開放（G5 前排期） | IT-FLOW-03 仍只有 `application/TestCreateShortURLAllowsDuplicateLongURL` |
| DEF-003 | 輕微 | 開放 | cmd/shorturl 覆蓋仍 2.8% |
| DEF-004 | 輕微 | 本機已修正；**main 上仍為舊版**（隨 G4-QA-R2） | — |
| DEF-005 | 輕微 | 開放，G5 前 | §3 突變仍 PASS |
| DEF-006 | 輕微 | 開放，G5 前 | contract_test.go internal_error 只有 `internal_error/create` |

無致命／嚴重缺陷。

## 5. G4 過關條件 × 品保觀點

| 條件 | 判定 | 依據 |
|---|---|---|
| PR 標示 REQ／SEC | ✅ | §2；#25 分支名 ENG-022 與標題 ENG-020 不一致（審查部 S1） |
| 簽名一致 | ✅（採審查部） | 品保未重做 |
| 單測涵蓋且通過 | ✅ | §3；DEF-005／006 G5 前 |
| SAST／SCA／機密掃描 | ✅ | run 39737480926 與各 PR 四綠 |
| 安全審查完成 | —（安全部，G4-R1） | 品保不判定 |
| 介面同步 OpenAPI | ✅（有缺口） | CT 已比對現行 openapi.yaml；不偵測漂移（DEF-005） |
| commit 規範與 DoD | ⚠️ | DoD D2（矩陣／案例進 main）待 G4-QA-R2；#24、#25 規範項依 G4-R4、G4-R5 |

## 6. G4-QA-R2 處置（**未關閉；僅提出方案**）

現況（gh 唯讀 15:11 查證）：main `docs/07-testing/` 只有 01～10 與 README，`02-traceability-matrix` 為 **v0.3**；本機為 v0.4（本次升 v0.5），11～19 號與 evidence 只在本機。

本次限制：使用者指明只做本機文件、不推送；開 PR／推送須使用者本人同意。品保**未**開 PR、未推送、未 commit。

提議（擇一，交使用者決定、審查部確認）：

| 方案 | 內容 | 關閉 G4-QA-R2？ |
|---|---|---|
| **A（建議，暫行）本機入庫＋sha256 綁定** | 以本文 §8 清單作為 07-testing 文件之正式基線，審查部與 G4 整體送審以 sha256 引用；任何改動須改版並更新清單 | **轉為「處置中（暫行）」**，不視為關閉；DoD D2「進 main」仍未達成，記為偏差，G5 前須以方案 B 補 |
| B 純文件 PR 進 main | 併入審查部 G4-R3 文件 PR（分支命名依 CR-003／D-15 結果）；PR 內容須與 §8 sha256 一致 | 合併後關閉 |

品保建議：G4 複審先採 A（審查部 G4-R3 亦已要求入庫，A 只作過渡）；B 的開立時間由使用者在 G4-R3 一併決定。**本文任何字句不構成推送或開 PR 之核准。**

## 7. 交接

- 品保部 → 總協調、審查部（G4-R2 複審）
- 結論：有條件通過（條件見 §0）
- 已知問題：DEF-002／003／005／006 開放；G4-QA-R2 處置中

## 8. 綁定清單（sha256，15:16 UTC+8 重算）

| 檔案（`docs/07-testing/`） | sha256 |
|---|---|
| evidence/g4-overall-f646-vet.txt | 00a5f1d995f1697b6f134bcca3e00d503345a5241662ca72a16f8b0ac5768294 |
| evidence/g4-overall-f646-test.json | 71353d9a8d74dea36ca6fb9916b07cfd93bfba0f6f2dafcd50e0867a5bcb89a7 |
| evidence/g4-overall-f646-race.txt | 42f375378befa50f51b9255f994c4bde7e54a00b04804f0d78d86abdca21ad86 |
| evidence/g4-overall-f646-ct.txt | e363d42c357983b5f34a117d508e459345af12b641ef9f020e5b4cd2c486e990 |
| evidence/g4-overall-f646-cover.txt | 2d4a41869fcc9f106babebcc8eb4bb95f2953bb0d4d5efbe0c0c31397c269c29 |
| evidence/g4-overall-f646-env.txt | 73138c31bfb06fb5d5699641b93b029d364cd2a65824f0d2d431a583e9eec537 |
| 12-gate-G4-qa-precheck-v0.1.md | 85831dc8c12f6b26a1899e55b5f534f7fa11de2e6c1cce1f501420ccdb29f09b |
| 13-G4-QA-R1-acceptance-PR26-v0.1.md | 2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c |
| 14-CR-003-qa-confirmation-v0.1.md | 3fb5d45727a337d9d9901ed7a5399922457d166b135cba10845b408f68b30ae9 |
| 16-CR-004-R4-qa-opinion-v0.1.md | ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c |
| 17-CR-004-signoff-2-2-qa-confirmation-v0.1.md | e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587 |
| 18-CR-004-AC6-9-qa-verdict-v0.1.md（檔內 v0.4） | 543763dad54eb6dd1bbd3b68afcdfa4c80e0b375cfea912fb24a02ef02d61288 |
| 02-traceability-matrix-v0.1.md（檔內 v0.5） | 見回報（本文寫入後才定稿，不自我引用） |

CR 狀態：CR-003 本文 c55ad3ca 已核實（審查部），但依 review-CR-003 R10 待使用者確認 D-15 才生效（G4-R5）；CR-004 AC-1～9：品保 18 號 v0.4 判定首測通過（AC-8 之 bypass 來源為合併紀錄引述使用者原話，品保未直接查證；G4-R4 待使用者親自追認）。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | G4 整體品保結論：#13～#26、基準 f646bee 重跑 181 PASS、CT 14 PASS；有條件通過；G4-QA-R2 提出方案 A／B，未關閉 |
