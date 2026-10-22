---
文件：關卡審查紀錄
關卡：G4 PR 合併前審查（PR #25｜維運部 CI 強化：SCR-032／019／021／004，ENG-020）——**同步後複審 #01**
版本：審查紀錄 v1
日期：2026-10-22（Asia/Taipei，UTC+8；審查時段約 14:25～14:30）
負責角色：審查部
狀態：已出具（結論：**通過**；可合併結論綁定 head `e54d9003f303e7e1a33082a70308d048caf0fdc1`＋base＝main `5084d6a95a0ac1d1a3884e990089fa742822c3b1`＋run 39733542195；附錄 A 定稿**未發布，待使用者核准**）
送審物：https://github.com/444112029012/short-url/pull/25
前次紀錄：docs/gates/gate-G4-PR25-2026-10-22.md（初審：內容通過，待同步後複審；§十 複審條件 1～7）
安全結論：快照 `docs/06-security/21-pr-025-ci-hardening-security-review-v0.1.md` v0.1.1（sha256 `08fe6783e0c76f7a2fffa8be090e09fef7b1076a16b0f22f2bf065f636b56365`），綁定 head `e54d9003f303e7e1a33082a70308d048caf0fdc1`
對應：ENG-020；NFR-004；SEC-009～011；SCR-004／019／021／032／050～052；CR-002；CR-003；CR-004；D-08～D-10
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 關卡審查紀錄：G4 PR #25（同步後複審 #01）

- **範圍**：初審 §十 條件 1～7，加上 job 名稱／ruleset 對應、merge trailer、PR 標題與描述、CR-004 一致性。內容面（workflow 逐項、SCR 達成）沿用初審 §四、§五：workflow blob 未變，所以仍有效。
- **審查方式**：只審不改（沒有改專案檔、沒有 commit、沒有在 GitHub 留言或 review、沒有合併）。`git ls-remote`；`/tmp/su25` fetch `refs/pull/25/*` 和 main；`git cat-file`／`show --remerge-diff`／`diff`／`rev-parse`；cursor-github 唯讀（PR、check-runs）；未認證公開 API（run、rulesets 24629805／24629854、reviews、comments）；`sha256sum` 自算。
- **結論**：**通過**。條件 1～7 全部成立；無必改。ENG 編號相關問題只涉及編號歸屬，不影響程式正確性和安全結論，列 S 項。

## 一、與宣稱核對

| 宣稱 | 核實 | 結果 |
|---|---|---|
| head `e54d9003…` | 14:26:50 `ls-remote`：`refs/pull/25/head`＝`e54d9003f303e7e1a33082a70308d048caf0fdc1`＝API `head.sha` | ✅ |
| base `5084d6a…`（#24 merge） | `refs/heads/main`＝`5084d6a95a0ac1d1a3884e990089fa742822c3b1`（「Merge pull request #24」）＝API `base.sha` | ✅ |
| run 39733542195 四綠 | 見 §三 | ✅ |
| merge commit trailer ENG-020 | e54d900 訊息「chore: merge main 5084d6a into fix/ENG-022-ci-hardening-scr032 (D-10)」，空行後為「ENG-020」 | ✅ |
| 描述已補分支名偏差 | 已補「分支名偏差說明」段落，但同一描述「任務：ENG-022（CI 強化）」沒改 | ⚠️ S1 |
| 標題 | 「fix: CI 強化 SCR-032／019／021／004 (ENG-020)」 | ✅ |

## 二、同步結構（條件 1～3、5、7）

| 條件 | 實測 | 結果 |
|---|---|---|
| 1 merge 非 rebase | e54d900 parents＝`bf190c4c9292f15f7c152ad340c07d58bac2fa58`＋`5084d6a95a0ac1d1a3884e990089fa742822c3b1`；`5084d6a..e54d900` 只有 e54d900、bf190c4 兩筆；bf190c4 SHA 不變，唯一 parent 仍是 `db51cd5` | ✅ |
| 2 remerge-diff | `git show --remerge-diff --format= e54d900` 輸出 0 bytes | ✅ |
| 3 blob | `.github/workflows/ci-devsecops.yml`：head 與 bf190c4 都是 `287a9b88a72f2bbae5a6d8f9c5a91fb58c27312e`，和初審 §七相同；`git diff --stat 5084d6a e54d900` 只有此檔（+48／−20）；API changed_files=1；`bf190c4..e54d900` 的差異（README、02-architecture、06-security 共 6 檔）全部來自 main | ✅ |
| 5 merge ref | `refs/pull/25/merge`＝`d035f3678b46db666d651914535b5edbd2d5041b`，parents `5084d6a`＋`e54d900`；tree `76579c34f21585a0a642e89263e0770851ad7efc`＝head tree | ✅ |
| 7 順序 | #24 已合併（main＝#24 merge），之後才同步（e54d900 committer 時間 14:00:48 UTC+8） | ✅ |

## 三、CI 與 required checks（條件 4）

- check-runs（ref＝e54d900）total_count＝**4**：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`。四筆都是 completed／success，head_sha 都是 `e54d9003…`，app id 都是 15368，都屬 run 39733542195（job 119673867189／119673866850／119673867044／119673867074；完成時間 14:22:44～14:24:30 UTC+8）。✅
- run 39733542195：head_sha＝e54d900、event `pull_request`、path `.github/workflows/ci-devsecops.yml`、run_attempt 1、conclusion success。✅
- **job 名稱與 ruleset 對應**：ruleset 24629805「main-required-checks」（active，updated_at 2026-10-07 14:33 UTC+8，本 PR 沒有動到）的 required contexts 是 `sca-govulncheck`、`sast-gosec`、`secrets-gitleaks`、`unit-test`，integration_id 都是 15368，strict=true。head workflow 的 job key 和 `name:` 是 secrets-gitleaks（L29–30）、sast-gosec（L55–56）、sca-govulncheck（L84–85）、unit-test（L127–128），逐字一一對應。job 層沒有 `if:`；step 層只有 dependency-review 兩步 `if: github.event_name == 'pull_request'`（初審已有）和 upload-artifact `if: always()`，沒有跳過整個 job 的條件，也沒有 `continue-on-error`。✅

## 四、安全部原件（條件 6）

| 檢查 | 結果 |
|---|---|
| 位置 | 快照 `docs/06-security/21-pr-025-ci-hardening-security-review-v0.1.md`，mtime 14:25:54；只在快照，不在 repo |
| sha256（自算） | `08fe6783e0c76f7a2fffa8be090e09fef7b1076a16b0f22f2bf065f636b56365`，與宣稱一致 ✅（初審 v0.1 是 `517b4e56…`，已被 v0.1.1 取代） |
| 安全部名義 | frontmatter「負責角色：安全部」「作者：安全部」；修訂紀錄 v0.1.1 作者「安全部」 ✅ |
| 綁定新 head | 狀態「v0.1.1：有條件可合併，僅適用於同步後最終 head `e54d9003f303e7e1a33082a70308d048caf0fdc1`」；§5 結論「被合併的 head 必須就是這個 SHA，若再有 push，此確認即失效」 ✅ |
| 與審查部實測一致 | §5 的 parents、remerge-diff 空、只改 workflow、blob `287a9b88…`、run 39733542195 四 success，都和本紀錄 §二、§三相符 ✅ |
| SCR 狀態 | SCR-004／019／021／032「達成，合併後關閉」；SCR-050（低）「開放，G5 前以後續 PR 關閉」；SCR-051、052（資訊）「開放」 |
| 開放項是否阻擋 | 不阻擋。安全部條件是「SCR-050 在 G5 前以後續 PR 關閉」，期限是 G5，不是 G4 合併前提；051／052 是資訊級。SCR-032 部分證據仍是維運轉述（初審 S5 沿用） |

→ 新 head 已有綁定的書面確認，**不列必改**。

## 五、可追溯性（ENG 編號）

- trailer：bf190c4 是原 commit，不能改寫，仍帶 ENG-022（初審 §三）；同步 merge e54d900 帶 ENG-020；標題 ENG-020；分支名 ENG-022（決定保留，描述已記偏差）。
- 專案紀錄裡，ENG-020 歷來用於 CI／維運類工作（例如 CR-002 對應「ENG-020」、PR #6「(ENG-020)」）。改掛 ENG-020 在內容上說得通。
- 缺口：改掛 ENG-020、以及 ENG-022 保留給 CR-003 的決定，沒有找到企劃部書面原件；PR 描述「任務：ENG-022」和偏差說明互相矛盾；安全部 21 號 §5 寫「合併 commit 訊息…用的是 `ENG-022`」，但實際上 e54d900 只有 subject 裡的分支名是 ENG-022，trailer 是 ENG-020，描述不精確。
- 判定：這些都只涉及編號歸屬，不影響程式正確性和安全結論；標題、trailer、描述偏差段已經能從 PR 追到 ENG-020，所以 **G4 可追溯性仍成立**，列 S1～S3。

## 六、CR-004 一致性（只記 S）

- ruleset 24629854「main-pr-review」：active，`required_approving_review_count: 1`、`dismiss_stale_reviews_on_push: true`、`require_last_push_approval: true`、`require_extra_approval_for_unattributed_changes: true`，updated_at 13:05:58 UTC+8（與 PR24 複審 #02 的 CR-004 紀錄相同）。
- PR #25：reviews 0 筆、comments 0 筆、mergeable_state `blocked`。和 CR-004 的設計一致：沒有核准就只能走使用者本人 admin bypass（D-08／D-09／RR-006）。沒有舊核准可被 dismiss，不構成不一致。
- 注意：bypass 合併時，規則二的 review 條件被繞過，但 ruleset 24629805 的四項 required checks（無 bypass、strict）仍然適用；PR 上必須先有審查部「可合併」留言（S4）。

## 七、必改清單

無。

## 八、S 項（不阻擋）

| 編號 | 內容 | 建議 | 負責角色 |
|---|---|---|---|
| G4-PR25-S1 | PR 描述「任務：ENG-022（CI 強化）」和同一描述的分支名偏差說明（改掛 ENG-020）矛盾 | 只改描述，改成「任務：ENG-020（原掛 ENG-022，見偏差說明）」；改描述不改 head，不影響本結論 | 維運部 |
| G4-PR25-S2 | 改掛 ENG-020、ENG-022 保留給 CR-003 的決定缺企劃部書面原件 | 企劃部補決定紀錄（decision-log 或任務表），G5 前完成 | 企劃部 |
| G4-PR25-S3 | 安全部 21 號 §5 說合併 commit 訊息用 ENG-022，與實際（trailer 是 ENG-020）不精確；21 號 v0.1.1 只在快照，未入庫（承初審 S1） | 下版修正措辭；入庫時另記 sha256 | 安全部 |
| G4-PR25-S4 | 合併須使用者本人 admin bypass（reviews 0）；被合併的 head 必須是 `e54d900`；合併前 PR 上要先有附錄 A 留言 | 依 D-09 流程 | 總協調；使用者 |
| G4-PR25-S5 | 初審 S2～S6、S8 沿用（SCR-050 G5 前、051、052、SCR-032 artifact 自核、gitleaks 釘值、Node 20）；初審 S7 已處理（描述已寫同步至 `5084d6a`） | 照初審 | 維運部；安全部 |

## 九、最後核對

14:26:50（UTC+8）`ls-remote`：main `5084d6a…`、`refs/pull/25/head` `e54d900…`、`refs/pull/25/merge` `d035f36…`，和 §一、§二相同。「」內文字逐字取自 PR、commit 訊息、21 號 v0.1.1。

---

## 附錄 A：「可合併」留言定稿

**未發布，待使用者核准。** 發布前一刻須再以 `ls-remote`／API 確認 head 與 main 都沒變；任一變動，本定稿即作廢。

## 審查部 G4 書面結論：可合併

- PR：#25 「fix: CI 強化 SCR-032／019／021／004 (ENG-020)」（ENG-020；NFR-004；SEC-009～011；SCR-004／019／021／032）
- **綁定 head：`e54d9003f303e7e1a33082a70308d048caf0fdc1`**（原 commit `bf190c4c9292f15f7c152ad340c07d58bac2fa58` 未改寫，加上同步 merge commit，trailer ENG-020）
- **綁定 base（＝審查當下 main）：`5084d6a95a0ac1d1a3884e990089fa742822c3b1`**（#24 合併後）
- 四項 required checks 綠燈 run：https://github.com/444112029012/short-url/actions/runs/39733542195 （head_sha＝上列 head；check runs 恰 4 筆：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` 都是 success，與 ruleset 24629805 required contexts 逐字對應）
- 差異核對：main...head 只有 `.github/workflows/ci-devsecops.yml`（blob `287a9b88a72f2bbae5a6d8f9c5a91fb58c27312e`，與初審相同）；remerge-diff 為空；merge ref `d035f36` 的 tree＝head tree `76579c34f21585a0a642e89263e0770851ad7efc`
- 前置：高風險（CI）。安全部書面確認 `docs/06-security/21-pr-025-ci-hardening-security-review-v0.1.md` v0.1.1（sha256 `08fe6783e0c76f7a2fffa8be090e09fef7b1076a16b0f22f2bf065f636b56365`）為有條件可合併，只適用於上列 head；SCR-004／019／021／032 合併後關閉；SCR-050（低）G5 前以後續 PR 關閉；SCR-051、052（資訊）開放，不阻擋。
- 必改：無。非必改見紀錄 S1～S5，不阻擋。
- 時效：本結論只適用於上列 head 與 base。**head 或 main 任一變動即失效**（D-10：strict=true）。
- 合併程序：依 D-08／D-09／RR-006，由使用者本人以 admin bypass 合併，本留言即審查部「可合併」結論與綠燈連結。
- 紀錄：`docs/gates/gate-G4-PR25-2026-10-22.md`（初審）、`docs/gates/gate-G4-PR25-2026-10-22-recheck-01.md`（同步後複審 #01；2026-10-22，Asia/Taipei）

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v1 | 2026-10-22 | 審查部 | 同步後複審 #01：head e54d900／base 5084d6a；merge 同步、remerge-diff 空、blob 287a9b88 不變、run 39733542195 恰 4 success 且 job 名稱與 ruleset 24629805 對應、merge tree＝head tree；安全部 21 號 v0.1.1（08fe6783…）綁定新 head；通過；無必改；S1～S5；附錄 A 定稿（未發布） |
