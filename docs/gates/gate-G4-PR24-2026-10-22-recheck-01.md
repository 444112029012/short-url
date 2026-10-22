---
文件：關卡審查紀錄
關卡：G4 PR 合併前審查（PR #24｜ENG-021 README 431 門檻與例外回應措辭）——**同步後複審 #01**（含 CR-004 首次實測核對）
版本：審查紀錄 v1
日期：2026-10-22（Asia/Taipei，UTC+8；審查時段約 13:16～13:20）
負責角色：審查部
狀態：已出具（結論：**退回**；同步本身核對全部成立，但安全部書面未綁新 head、CR-004 維運證據未入 repo docs 且無 AC-1～5 維運判定原件；**不出具「可合併」，無附錄 A 定稿**）
送審物：https://github.com/444112029012/short-url/pull/24
前次紀錄：docs/gates/gate-G4-PR24-2026-10-22.md（初審：內容通過（待同步後複審）；附錄 A 同步後複審條件 1～7）
安全結論：`docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md` v0.1（sha256 `f978d1b550bfdb3762a84896365b5af58f9b1982bfc2c4aa01151b4e168eb5bc`）**只綁舊 head `80fc321`**，未綁新 head
對應：ENG-021；NFR-002、NFR-006；SEC-012；DR-G4-01；SCR-046；CR-004（AC-CR004-1～9）；D-08～D-10；RR-006
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 關卡審查紀錄：G4 PR #24（同步後複審 #01）

- **送審來源**：他人宣稱 head `0379b64c892e9f664e07c879a16f9acaa164be88`、base `2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`（#23 merge）、run 39726831602 四綠；CR-004 after.json sha256 `1404dab6…`、中間版 `9b4fffe1…`（13:04 誤關、13:05:58 勾回）。以下逐項**獨立核實**。
- **審查方式**：只審不改（未改專案檔、未 commit、未在 GitHub 留言或 review、未合併）。cursor-github 唯讀（`get_pull_request` #24、`get_workflow_run` 39726831602、`list_check_runs_for_ref` filter=all、`get_job_logs`×3）；未認證公開 REST `GET /repos/…/rulesets/24629854`、`/rules/branches/main`；`git ls-remote`；`/tmp/rv24` 唯讀 clone＋fetch；`sha256sum` 自算。
- **結論**：**退回**。初審附錄 A 條件 1～5、7 全部成立（同步為 merge、原 commit 未改寫、remerge-diff 空、README blob 不變、run 四綠恰 4 筆、merge ref tree＝head tree）；條件 6 不成立（安全部書面只綁 `80fc321`）→ R1。CR-004 首次實測：ruleset 現況符合，after.json／中間版 sha256 與宣稱相符，但兩檔**不在 repo docs**（只在 `/workspace/rebuild/cr-004/exec/`），且找不到 AC-CR004-1～5 的維運判定原件 → R2。

---

## 一、初審附錄 A 條件逐項核對

| # | 條件 | 結果 | 證據 |
|---|---|---|---|
| 0 | head／main | ✅ `refs/pull/24/head`＝API `head.sha`＝`0379b64c892e9f664e07c879a16f9acaa164be88`；`refs/heads/main`＝API `base.sha`＝`2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`（「Merge pull request #23 from 444112029012/docs/ENG-011-security-08-v015-readme-v022」，parents `363f169`＋`c3155b3`，12:09:17 UTC+8）；`rev-list --count head..main`＝0 | `ls-remote`；API |
| 1 | 同步為 merge、remerge-diff 空 | ✅ `0379b64` parents＝`80fc321a2d363871bfdf9b50aa89ca10a3be359d`（第一 parent）＋`2a02819…`；訊息「Merge remote-tracking branch 'origin/main' into docs/ENG-021-readme-followup」，13:07:07 UTC+8；`git show --remerge-diff` 只輸出 commit 行（41 bytes），無 diff；非 rebase | 本機 |
| 2 | 原 commit 未改寫 | ✅ `80fc321…` 同 SHA 為 head 祖先（`merge-base --is-ancestor` 成立）；PR commits＝2 | 本機；API |
| 3 | README blob 不變、只 1 檔 | ✅ `git diff --name-status main...head`：只 `M README.md`；+3／−3；blob `fbdfb496b2ff44a31907abbc8d572447691cf1cd`＝初審；API changed_files 1 | 本機；API |
| 4 | run 四綠恰 4 筆 | ✅ 見 §二 | API |
| 5 | merge ref tree＝head tree | ✅ `refs/pull/24/merge`＝API `merge_commit_sha`＝`100dc6ddc0f4a529a00d562d6c471b424baf1d55`；parents `2a02819`、`0379b64`；tree `2fc3c9008500f37e4487c6adcaad7b9a99b9d00d`＝head tree＝run `head_commit.tree_id` | 本機；API |
| 6 | 安全部對新 head 書面 | ❌ 見 §三 → **R1** | 本機檔 |
| 7 | head／main 未變 | ✅（審查當下）見 §六 | `ls-remote` |

## 二、CI（run 39726831602）

| 項目 | 結果 |
|---|---|
| run | `ci-devsecops` #73、event=`pull_request`、attempt 1、completed／**success**；`head_sha`＝`0379b64…` ✅；`head_branch`＝`docs/ENG-021-readme-followup`；`pull_requests` 只含 #24（head `0379b64`、base `2a02819`）；check_suite 107664109343；13:13:55～13:15:10（UTC+8） |
| check runs（filter=all） | total_count **4**：119652549727 `secrets-gitleaks`、119652549970 `sast-gosec`、119652550021 `sca-govulncheck`、119652549955 `unit-test`；皆 completed／success、head_sha＝head、app 15368、check_suite 107664109343、details_url 皆屬 run 39726831602 ✅ |
| checkout | 「HEAD is now at 100dc6d Merge 0379b64c892e9f664e07c879a16f9acaa164be88 into 2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5」 |
| gitleaks log | 「68 commits scanned」、「no leaks found」 |
| gosec log | 「go version go1.27.2」、「gosec h1 OK」（初審 S4 部分處理；Issues 數未擷取到，以 job success 判定） |
| govulncheck log | 「No vulnerabilities found.」；「BASE_SHA: 2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5」「HEAD_SHA: 0379b64c892e9f664e07c879a16f9acaa164be88」；「dependency graph compare HTTP 200」 |
| unit-test | success（log 本次未讀） |

## 三、安全部書面原件

- `docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md` v0.1，審查部自算 sha256＝`f978d1b550bfdb3762a84896365b5af58f9b1982bfc2c4aa01151b4e168eb5bc`。
- 原文：狀態「通過，僅適用於 head `80fc321a2d363871bfdf9b50aa89ca10a3be359d`；同步 main 後須針對最終 head 另行補確認」；§3「若因 strict 須同步 main，新 head 須符合以下條件，安全部才會針對該 head 另出確認」。
- 檢索 `docs/` 全部檔：**無**任何含 `0379b64` 的檔案 → 安全部未對新 head 出確認，SCR-046 關閉條件亦未就新 head 成立 → **R1**。（26 號檔為 CR-004 安全意見，不是 #24 head 確認。）
- 依 §一、§二，安全部 20 號檔 §3 所列四條件於新 head 皆已成立，補確認應無實質障礙。

## 四、CR-004 首次實測核對（#24 為第一個實測）

### 4.1 ruleset 24629854 現況（公開 REST，13:17 UTC+8 讀取）

| 欄位 | 值 | 判定 |
|---|---|---|
| name／enforcement | `main-pr-review`／`active`；conditions `~DEFAULT_BRANCH` | ✅ |
| `dismiss_stale_reviews_on_push` | `true` | ✅ |
| `require_last_push_approval` | `true` | ✅ |
| `required_approving_review_count` | `1` | ✅ |
| `require_extra_approval_for_unattributed_changes` | `true`（與變更前相同，13:04 誤關後已勾回） | ✅ |
| 其他 rules | `deletion`、`non_fast_forward` | ✅ |
| `updated_at` | `2026-10-22T05:05:58.537Z`＝13:05:58（UTC+8），與宣稱勾回時間相符 | ✅ |
| bypass_actors | 未認證讀取回 `null`（不可見）；after.json 記 `actor_id 5 RepositoryRole pull_request` | 以 after.json 為準 |

- `/rules/branches/main` 另含 24629805 的 `deletion`、`non_fast_forward`、`required_status_checks`（strict=true）。

### 4.2 維運證據檔 sha256（審查部自算）

| 檔案（實際位置） | sha256 | 宣稱 | 內容 |
|---|---|---|---|
| `/workspace/rebuild/cr-004/exec/after.json`（1074 B，13:06） | `1404dab6a8767b6e6235d7c1f9ccadf8dd309cc97606f15d53c2a3952b7ccefa` | `1404dab6…` ✅ | 兩欄位 true、unattributed true、updated_at 13:05:58.537 |
| `/workspace/rebuild/cr-004/exec/after-1-partial.json`（中間版，1075 B，13:05） | `9b4fffe12602411d24e41258a8e02075f4f634a1636419f0e5311ac5cf6d4717` | `9b4fffe1…` ✅ | 兩欄位 true，但 `require_extra_approval_for_unattributed_changes: false`、updated_at 13:04:53.307（誤關） |
| `before.json` | `af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12` | — | — |
| `diff-step4.txt` | `6f680c074d0c8273d773b323d6b334e92db259afcdad87901673d97c0468194b` | — | 恰兩欄位 false→true（AC-3 相符） |
| `805-before.json`＝`805-after.json` | `ba4cad93614911f45d7bf21fe7db94f265c8f5817e8716e6a163906c596b7984`；`diff-805.txt` 空 | — | AC-4 相符 |
| `ui-screenshot.png`／`ui-screenshot-2.png` | `c8d1f4eb…1260`／`368c9e7e…db29` | — | AC-5 截圖存在（內容未逐張判讀） |

- **問題**：上列檔案**不在** repo（main `2a02819`、#24 head 皆無）也**不在**快照 `docs/`（含 `09-cicd-env/evidence/`，該處只有 2026-10-07 的舊快照）；只在 box 暫存 `/workspace/rebuild/cr-004/exec/`。

### 4.3 AC 證據狀態

| AC | 狀態 |
|---|---|
| AC-CR004-1～5 | 原始證據檔存在（暫存區），審查部抽核 AC-2、3、4 值相符；**但 `docs/` 內無維運部對 AC-1～5 的判定原件**（檢索 `09-cicd-env/` 與全 docs 之 `AC-CR004`，只出現在 CR、品保 15～17 號與審查部 review 檔）→ **R2** |
| AC-CR004-6～9 | 待品保判定（`pr24-after-sync.json`、實測紀錄、`pr24-merged.json` 均未見；8 須合併後才有）；**不列合併前提、非必改**（S4） |

- 文本核對：CR-004 v0.3 `docs/99-changelog/CR-004-ruleset-stale-approval-v0.3.md` sha256 `eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4` ✅；品保 16 號 `ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c` ✅；17 號 `e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587` ✅（均與宣稱前綴相符）。

## 五、必改清單

| 編號 | 條件 | 問題 | 必改內容 | 負責角色 |
|---|---|---|---|---|
| G4-PR24-R1 | 初審附錄 A 條件 6；G4 七條件 5 | 安全部 20 號檔只綁 `80fc321`，明寫同步後須另補確認；docs 內無綁 `0379b64` 的書面 | 出具綁定 head `0379b64c892e9f664e07c879a16f9acaa164be88`（base `2a02819`）的書面確認（含 SCR-046 狀態），提供路徑與 sha256 | 安全部 |
| G4-PR24-R2 | CR-004 首次實測 AC-CR004-1～5（16 號檔 §3） | after.json、中間版等證據只在 box 暫存區，未入 repo `docs/`；無維運部 AC-1～5 判定原件；中間版誤關一事無書面說明 | 將 `/workspace/rebuild/cr-004/exec/` 證據（sha256 不得變）收入 `docs/09-cicd-env/evidence/`，並出具 AC-CR004-1～5 逐項判定原件（含 13:04:53 誤關 unattributed、13:05:58 勾回的經過與中間版 sha256 `9b4fffe1…`） | 維運部 |

## 六、S 項（不阻擋）

| 編號 | 內容 | 建議負責 |
|---|---|---|
| S1～S3、S5、S6 | 初審沿用（POST /health 重複 Allow、501 未列、go 版本差、trailer 出處、SCR-043／044 比照） | 研發部 |
| S4 | AC-CR004-6～9 待品保判定（未出列為合併前提）；初審 S4 job log 已讀 3 份 | 品保部 |
| S7 | `mergeable_state: blocked`；CR-004 後 `require_last_push_approval=true`，sync push（13:07:07）後亦無核准，仍須 D-08／RR-006 使用者本人 admin bypass，被合併 head 須＝審查部／安全部確認之 SHA | 總協調；使用者 |
| S8 | 公開 API 讀不到 bypass_actors，建議維運部以認證 `gh api` 重取現況快照附入 R2 證據 | 維運部 |
| S9 | PR 描述測試證據仍寫 head `80fc321`、未提同步 merge 與 run 39726831602（只改描述） | 研發部 |

## 七、與宣稱比對

| 宣稱 | 核實 |
|---|---|
| head `0379b64…` | ✅ |
| base `2a02819…`（#23 merge） | ✅ |
| run 39726831602 四綠 | ✅ head_sha 相符、恰 4 筆 |
| after.json `1404dab6…` | ✅ sha256 相符；⚠ 不在 repo docs |
| 中間版 `9b4fffe1…`，13:04 誤關、13:05:58 勾回 | ✅ sha256、updated_at 13:04:53／13:05:58 相符；⚠ 不在 repo docs |
| CR-004 eae83378／16 號 ad44316f／17 號 e34b76f7 | ✅ |

## 八、最後核對

2026-10-22 13:18:09（Asia/Taipei，UTC+8）`git ls-remote`：`refs/heads/main`＝`2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`、`refs/pull/24/head`＝`0379b64c892e9f664e07c879a16f9acaa164be88`、`refs/pull/24/merge`＝`100dc6ddc0f4a529a00d562d6c471b424baf1d55`，與審查開始一致。引號稽核：「」內文字逐字取自 commit 訊息、job log、20 號檔、16 號檔。

---

## 附錄 A：「可合併」定稿

> **不出具。** 本次結論為退回。R1、R2 完成後，若 head／main 未變，由審查部複審 #02 核對 R1、R2 後出具定稿（預計綁 head `0379b64`、base `2a02819`、run 39726831602）；任一變動即須重取四綠。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v1 | 2026-10-22 | 審查部 | 同步後複審 #01：head 0379b64（parents 80fc321＋2a02819）／base＝main 2a02819；remerge-diff 空；80fc321 未改寫；README blob fbdfb49 不變；run 39726831602 四綠恰 4 筆；merge ref 100dc6d tree＝head tree 2fc3c90；安全部 20 號只綁 80fc321（R1）；CR-004 ruleset 現況符合、after.json／中間版 sha256 相符但未入 repo docs、無 AC-1～5 維運判定（R2）；退回；無附錄 A 定稿 |
