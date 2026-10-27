---
文件：關卡審查紀錄
關卡：G4 PR 合併前審查（PR #26｜G4-QA-R1 契約測試 CT-*，ENG-013）——**同步後複審 #01**
版本：審查紀錄 v1
日期：2026-10-22（Asia/Taipei，UTC+8；審查時段約 14:35～14:40）
負責角色：審查部
狀態：已出具（結論：**通過**；可合併結論綁定 head `49e29d80d1f1af227d83c40c0ab212f4372029f0`＋base＝main `d0749f57038aa0f56e4fa26b2c2b3c97ceb71265`＋run 39734581059；附錄 A 定稿**未發布，待使用者核准**）
送審物：https://github.com/444112029012/short-url/pull/26
前次紀錄：docs/gates/gate-G4-PR26-2026-10-22.md（初審：內容通過，待同步後複審；附錄 A.1 條件 R1～R8）
安全結論：快照 `docs/06-security/27-pr-026-contract-tests-security-confirmation-v0.1.md` v0.1（sha256 `56b20589e3346b8906f9854605e6c7bff71977edf12f4ce36d355414f28acf6d`），綁定 head `49e29d80d1f1af227d83c40c0ab212f4372029f0`；SCA 意見 25 號（sha256 `0975ec992f6efc0e58c3288aad63093d5f4666afbf869987e610e767a99e31e5`）
品保結論：`docs/07-testing/13-G4-QA-R1-acceptance-PR26-v0.1.md`（sha256 `2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c`）、`12-gate-G4-qa-precheck-v0.1.md`（sha256 `85831dc8c12f6b26a1899e55b5f534f7fa11de2e6c1cce1f501420ccdb29f09b`）
對應：ENG-013（主）、ENG-010、ENG-014、ENG-015；REQ-001、002、004、005、007、008、009；SEC-001、003、014；G4-QA-R1／DEF-001；DEF-005、DEF-006；D-08～D-10
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 關卡審查紀錄：G4 PR #26（同步後複審 #01）

- **範圍**：初審附錄 A.1 條件 R1～R8；外加新 head 本機實跑、-json 比對、DEF-005 重做、安全部 27 號與 25 號、品保 12／13 號。內容面（測試是否真的驗證契約、突變 6／6）沿用初審 §四：`contract_test.go` blob 未變，仍有效。
- **審查方式**：只審不改（未改專案檔、未 commit、未在 GitHub 留言／review、未合併）。`git ls-remote`；`/tmp/g4pr26` fetch `refs/pull/26/head`、`refs/pull/26/merge`、main；`git log`／`show --remerge-diff`／`diff`／`rev-parse`／`merge-base`；拋棄式 worktree（`/tmp/g4pr26-wt` @ 49e29d8，用後移除）實跑；未認證公開 API（PR、run、check-runs、artifacts）；`sha256sum` 自算。
- **結論**：**通過**。R1～R8 全部成立；安全部已出具綁定新 head 的書面確認；**必改：無**。S1～S8 不阻擋。

## 一、與宣稱核對

| 宣稱 | 核實 | 結果 |
|---|---|---|
| head `49e29d80…` | 14:35:56 `ls-remote`：`refs/pull/26/head`＝`49e29d80d1f1af227d83c40c0ab212f4372029f0`＝API `head.sha` | ✅ |
| base `d0749f5…`（#25 merge） | `refs/heads/main`＝`d0749f57038aa0f56e4fa26b2c2b3c97ceb71265`（「Merge pull request #25 from 444112029012/fix/ENG-022-ci-hardening-scr032」）＝API `base.sha` | ✅ |
| run 39734581059 四綠 | 見 §三 | ✅ |
| trailer ENG-013 ENG-014 ENG-015 ENG-010 | 49e29d8 訊息「chore: merge main d0749f5 into test/ENG-013-openapi-contract-tests (D-10)」，空行後為「ENG-013 ENG-014 ENG-015 ENG-010」 | ✅ |
| -json 原檔 14 PASS | 自算 sha256 `121d219e3ce083a225f023c1e6c2048cef91267f9119f49355d84b19951afca6`＝宣稱值；79 行、只有 httpapi 套件；14 筆 pass（6 頂層＋CT_ERR_01 8 子測）、0 fail、0 skip；與本審查新 head 實跑逐名一致（§四） | ✅（產生時間見 S3） |
| 安全部 27 號 sha256 `56b20589…f6d`、綁 49e29d80、無條件可合併 | 自算相同；見 §五 | ✅ |
| 25 號 kin-openapi 結案 | head go.mod、go.sum `kin-openapi` 0 筆；`git grep -i kin-openapi` 排除 docs 後 0 檔 | ✅ |

## 二、同步結構（R1～R3、R7）

| 條件 | 實測 | 結果 |
|---|---|---|
| R1 #23～#25 已合併 | main merge 序列：2a02819「Merge pull request #23」→5084d6a「#24」→d0749f5「#25」（parents `5084d6a`＋`e54d9003…`）；新 main＝`d0749f5` | ✅ |
| R2 merge 非 rebase | 49e29d8 parents＝`e224253de29743a956e8c952cafc4eff8acff1d7`＋`d0749f57…`；`main..head` 只有 49e29d8、e224253 兩筆；e224253 SHA 不變、parent 仍 `363f169`、時間 09:11:45、訊息與 trailer 與初審 §1.1 同；`merge-base`＝`d0749f5`；`rev-list --count head..main`＝0 | ✅ |
| R3 remerge-diff／blob／範圍 | `git show --remerge-diff --format= 49e29d8` 0 bytes；`git diff --stat d0749f5 49e29d8` 只有 `internal/httpapi/contract_test.go`（+258）；blob `d774172f383708a7da591ba7e2cb2c911efc77ea`＝e224253 版；`go.mod`／`go.sum`／`.github` diff 0 行；API changed_files=1、additions=258 | ✅ |
| R7 merge ref | `refs/pull/26/merge`＝`e656b9c1c760934dec768fc91fff12138675afe1`＝API `merge_commit_sha`；parents `d0749f5`＋`49e29d8`；tree `9d0085070b472af5e27a4b95b6cf3fa8af4479d4`＝head tree＝run `head_commit.tree_id` | ✅ |

## 三、CI 與 required checks（R4）

- run 39734581059：`ci-devsecops` #77、event `pull_request`、path `.github/workflows/ci-devsecops.yml`、run_attempt 1、completed／**success**；head_sha＝`49e29d80…`；head_branch `test/ENG-013-openapi-contract-tests`；pull_requests 只含 #26（head 49e29d8、base d0749f5）；建立 14:32:31、更新 14:34:31（UTC+8）。✅
- check-runs（ref＝49e29d8，`filter=all`）total_count＝**4**：119677182706 `secrets-gitleaks`、119677182891 `sast-gosec`、119677182560 `sca-govulncheck`、119677182927 `unit-test`；皆 completed／success、app 15368、details_url 皆屬 run 39734581059；完成時間 14:32:38～14:34:30（UTC+8）。名稱與 ruleset 24629805 required contexts 逐字對應（PR25 複審 #01 §三）。✅
- workflow 為 #25 合併後版本（main 上 blob 與 #25 相同，本 PR 未動 `.github`），unit-test 即 `-race -json`、gosec 含 `-tests`。run 有 artifact `go-test-json`（16057 bytes）；下載需認證，本審查未開檔，改以 §四本機實跑與宣稱 -json 雙重比對（S4）。
- PR：open、merged=false、`mergeable_state: clean`；comments 0、review comments 0。

## 四、新 head 本機實跑與 DEF-005 重做（worktree @ 49e29d8；go1.27.3，CI 1.27.2）

| 指令 | 結果 |
|---|---|
| `go vet ./...` | 無輸出 |
| `go test -count=1 -race ./...` | 10 套件 ok（ports 無測試檔） |
| `go test -count=1 -race -json -run 'TestCT_' ./internal/httpapi/` | exit 0；**14 pass**（6 頂層＋CT_ERR_01 8 子測：internal_error/create、invalid_url/create、invalid_url/redirect_bad_code、invalid_url/stats_bad_code、not_found/redirect、not_found/stats、rate_limited/create、rate_limited/redirect），0 fail／skip |
| 與宣稱 -json 比對 | (Test, Action) 集合完全相同 ✅ |
| DEF-005 重做：`openapi.yaml` L69 `'201'`→`'200'` 後再跑 `-json -run TestCT_` | 仍 **14 pass** → DEF-005 在新 head 仍存在（已知、已接受，期限 G5 前）。`git checkout -- .` 後 status 空、HEAD＝49e29d8，worktree 已移除 |

## 五、安全部（27 號、25 號）

| 檢查 | 結果 |
|---|---|
| 27 號 sha256（自算） | `56b20589e3346b8906f9854605e6c7bff71977edf12f4ce36d355414f28acf6d`＝宣稱 ✅；只在快照（mtime 14:36），未入庫 |
| 名義 | frontmatter「負責角色：安全部」「作者：安全部」；修訂紀錄 v0.1 作者「安全部」✅ |
| 綁新 head | 狀態「**PR #26 可合併**，只適用於 head `49e29d80d1f1af227d83c40c0ab212f4372029f0`」；§5「沒有附加條件」「合併前如果再有 push，此確認即失效」✅ |
| 與審查部實測一致 | §2 parents、remerge-diff 空、只新增 contract_test.go、go.mod／go.sum 無變更、run 39734581059 四 success、-json sha256 皆與本紀錄相符 ✅ |
| 25 號 | sha256 `0975ec99…31e5` 自算；狀態「不同意」；27 號 §4「結論維持不同意」「檔 25 的申請在 PR #26 中結案」；本審查核 head go.mod／go.sum 無 kin-openapi ✅ → 初審 S1（docs 找不到安全部原文）**已解決** |

→ 新 head 已有綁定的安全部書面確認，**不列必改**。（初審 §七已判非高風險，此確認屬加強證據。）

## 六、品保（R5、R6）

| 檢查 | 結果 |
|---|---|
| 13 號 sha256 | `2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c`（初審為 `e688bd6d…157d`，已變更：新增「附錄：關閉條件複核（2026-10-22 09:5x）」；版本號仍 v0.1、修訂紀錄未加列 → S2） |
| R5／C1 | 附錄「C1：關閉 ✅」，依據 PR 描述原文「安全部 SCA 不同意引入 kin-openapi，所以以標準函式庫手寫斷言」「常數與 openapi.yaml 為人工對應…」；本審查以 API 核 PR 描述第 33、36 行確有此兩句 ✅ |
| R6／C2 | 附錄「C2：關閉 ✅」；12 號 sha256 `85831dc8…f09b` 自算，缺陷表 DEF-005、DEF-006「開放；G5 前」，DEF-001「關閉」；G4-QA-R1 最終判定「關閉」✅ |
| 綁定 head | 13 號對象仍為 `e224253`；新 head 對 contract_test.go blob 無變更（§二 R3），品保結論內容仍適用 → 不阻擋（S2） |

## 七、必改清單

無。

## 八、S 項（不阻擋）

| 編號 | 內容 | 建議 | 負責角色 |
|---|---|---|---|
| G4-PR26-S1 | PR 描述仍寫「基底：origin/main `363f169`…commit `e224253`」，未記同步至 d0749f5／head 49e29d8 | 只改描述；不改 head，不影響本結論 | 研發部 |
| G4-PR26-S2 | 13 號加附錄後版本仍 v0.1、修訂紀錄未加列；對象仍綁 e224253 | 下版升版並註明新 head；12、13 號仍未入庫（承初審 S7） | 品保部 |
| G4-PR26-S3 | 宣稱 -json 檔首筆時間 14:30:26（UTC+8），早於同步 commit 49e29d8 的 14:30:35，應是在 e224253（同 blob）跑的；27 號也標為轉述證據 | 不影響（blob 相同且本審查已在新 head 重跑一致）；往後證據檔請在最終 head 產生 | 研發部；總協調 |
| G4-PR26-S4 | CI artifact `go-test-json` 本審查未下載（需認證） | 品保合併後在 main 重跑 CT 時一併核 artifact | 品保部 |
| G4-PR26-S5 | 27 號「最後更新：2026-10-22（14:40）」晚於檔案 mtime 14:36 與本審查讀取時間 | 下版修正時間戳；27 號入庫時另記 sha256 | 安全部 |
| G4-PR26-S6 | DEF-005（新 head 重做仍 14 PASS）、DEF-006 開放，G5 前 | 照 12 號期限；初審 S5 三處 spec 突變作驗收 | 研發部；品保部 |
| G4-PR26-S7 | 合併須使用者本人 admin bypass（reviews 0）；被合併 head 須為 49e29d8；合併前 PR 上要先有附錄 A 留言 | 依 D-08／D-09 | 總協調；使用者 |
| G4-PR26-S8 | 初審 S2～S4、S6、S8、S9 沿用；S1 已解決（§五）；S5 併入本 S6；S7 併入本 S2 | 照初審 | 各該角色 |

## 九、最後核對

14:36:56（UTC+8）`ls-remote`：main `d0749f57…`、`refs/pull/26/head` `49e29d80…`、`refs/pull/26/merge` `e656b9c1…`，與 §一、§二相同（R8 於發布前一刻須再核）。「」內文字逐字取自 PR、commit 訊息、27／25／12／13 號、本機輸出。

---

## 附錄 A：「可合併」留言定稿

**未發布，待使用者核准。** 發布前一刻須再以 `ls-remote`／API 確認 head 與 main 都沒變；任一變動，本定稿即作廢。

## 審查部 G4 書面結論：可合併

- PR：#26 「test: 補 OpenAPI 契約測試 CT-* (ENG-013)」（ENG-013、ENG-010、ENG-014、ENG-015；REQ-001、002、004、005、007、008、009；SEC-001、003、014；關閉 G4-QA-R1／DEF-001）
- **綁定 head：`49e29d80d1f1af227d83c40c0ab212f4372029f0`**（原 commit `e224253de29743a956e8c952cafc4eff8acff1d7` 未改寫，加上同步 merge commit，trailer ENG-013 ENG-014 ENG-015 ENG-010）
- **綁定 base（＝審查當下 main）：`d0749f57038aa0f56e4fa26b2c2b3c97ceb71265`**（#25 合併後）
- 四項 required checks 綠燈 run：https://github.com/444112029012/short-url/actions/runs/39734581059 （head_sha＝上列 head；check runs 恰 4 筆：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` 都是 success）
- 差異核對：main...head 只有 `internal/httpapi/contract_test.go`（+258，blob `d774172f383708a7da591ba7e2cb2c911efc77ea`，與初審相同）；未動正式程式、`.github/`、`go.mod`／`go.sum`；remerge-diff 為空；merge ref `e656b9c` 的 tree＝head tree `9d0085070b472af5e27a4b95b6cf3fa8af4479d4`
- 測試：審查部於上列 head 實跑 `go test -race -json -run TestCT_`，TestCT_* 14 項 pass（6 頂層＋CT-ERR-01 8 子測），與送審 -json 一致
- 前置：品保 13 號 C1、C2 已關閉，G4-QA-R1／DEF-001 關閉；安全部 `docs/06-security/27-pr-026-contract-tests-security-confirmation-v0.1.md`（sha256 `56b20589e3346b8906f9854605e6c7bff71977edf12f4ce36d355414f28acf6d`）為可合併、無附加條件，只適用於上列 head；kin-openapi 未引入（25 號結案）
- 範圍說明：CT 以常數驗證實作符合現行 openapi.yaml，不偵測 spec 變更（DEF-005，期限 G5 前）；redirect／stats 500 未測（DEF-006，期限 G5 前）
- G4 七條件：1、3、4、5、7 達成；2、6 不適用。必改：無。非必改見紀錄 S1～S8，不阻擋。
- 時效：本結論只適用於上列 head 與 base。**head 或 main 任一變動即失效**（D-10：strict=true）。
- 合併程序：依 D-08／D-09／RR-006，由使用者本人以 admin bypass 合併，本留言即審查部「可合併」結論與綠燈連結。
- 紀錄：`docs/gates/gate-G4-PR26-2026-10-22.md`（初審）、`docs/gates/gate-G4-PR26-2026-10-22-recheck-01.md`（同步後複審 #01；2026-10-22，Asia/Taipei）

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v1 | 2026-10-22 | 審查部 | 同步後複審 #01：head 49e29d8／base d0749f5；merge 同步、e224253 未改寫、remerge-diff 空、blob d774172f 不變、run 39734581059 恰 4 success、merge tree＝head tree；新 head 實跑 14 PASS 與 -json（121d219e…）一致；DEF-005 重做仍 14 PASS；安全部 27 號（56b20589…）綁新 head、25 號 kin-openapi 結案已核；品保 13 號（2627c2da…）C1／C2 關閉、12 號（85831dc8…）；通過；無必改；S1～S8；附錄 A 定稿（未發布） |
