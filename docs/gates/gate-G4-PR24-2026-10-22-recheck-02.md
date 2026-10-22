---
文件：關卡審查紀錄
關卡：G4 PR 合併前審查（PR #24｜ENG-021 README 431 門檻與例外回應措辭）——**同步後複審 #02**（只核 G4-PR24-R1、R2）
版本：審查紀錄 v1
日期：2026-10-22（Asia/Taipei，UTC+8；審查時段約 13:21～13:24）
負責角色：審查部
狀態：已出具（結論：**通過**；可合併結論綁定 head `0379b64c892e9f664e07c879a16f9acaa164be88`＋base＝main `2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`＋run 39726831602；附錄 A 定稿**未發布，待使用者核准**）
送審物：https://github.com/444112029012/short-url/pull/24
前次紀錄：docs/gates/gate-G4-PR24-2026-10-22.md（初審）、docs/gates/gate-G4-PR24-2026-10-22-recheck-01.md（複審 #01：退回，R1、R2）
安全結論：快照 `docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md` v0.1.1（sha256 `b10707d0c28eb2631e85f65047332adeab7541e0ba12ace18466ca79b6c719e9`），綁定 head `0379b64c892e9f664e07c879a16f9acaa164be88`
對應：ENG-021；NFR-002、NFR-006；SEC-012；SCR-046；CR-004（AC-CR004-1～5）；D-08～D-10；RR-006
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 關卡審查紀錄：G4 PR #24（同步後複審 #02）

- **範圍**：只核複審 #01 的 R1、R2，並重核 head／main 未變。其餘項目（同步 merge、remerge-diff、blob、run 四綠恰 4 筆、merge ref tree）沿用複審 #01 §一、§二，head 未變故仍有效。
- **審查方式**：只審不改（未改專案檔、未 commit、未在 GitHub 留言或 review、未合併）。`git ls-remote`；`/tmp/rv24` fetch main 與 `refs/pull/24/head` 後 `git ls-tree`；`sha256sum` 自算；逐行閱讀兩份原件。
- **存放位置（重要）**：R1、R2 所列檔案**全部只在 `/workspace/projects/short-url/` 快照**；main `2a02819` 與 #24 head `0379b64` 的 tree 中**都沒有** `docs/06-security/20-*` 與 `docs/09-cicd-env/evidence/cr-004/`（`git ls-tree -r` 0 筆）。未入版控 → S1。
- **結論**：**通過**。R1、R2 已處理；無必改。

---

## 一、head／main

2026-10-22 13:21:05（UTC+8）`git ls-remote`：`refs/heads/main`＝`2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`、`refs/pull/24/head`＝`0379b64c892e9f664e07c879a16f9acaa164be88`、`refs/pull/24/merge`＝`100dc6ddc0f4a529a00d562d6c471b424baf1d55`，與複審 #01 相同 ✅。

## 二、R1（安全部）

| 檢查 | 結果 |
|---|---|
| 位置 | 快照 `docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md`（只在快照，不在 repo） |
| sha256（自算） | `b10707d0c28eb2631e85f65047332adeab7541e0ba12ace18466ca79b6c719e9`；宣稱 `b10707d0…9e9` ✅ |
| 安全部名義 | frontmatter「負責角色：安全部」「作者：安全部」；修訂紀錄 v0.1.1 作者「安全部」 ✅ |
| 綁定新 head | 狀態「**v0.1.1：通過，僅適用於同步後最終 head `0379b64c892e9f664e07c879a16f9acaa164be88`**」；§4 結論「安全部同意合併 head `0379b64c892e9f664e07c879a16f9acaa164be88`。被合併的 head 必須就是這個 SHA；若再有 push，此確認即失效。」 ✅ |
| SCR-046 狀態 | 「**SCR-046 狀態：開放，已轉研發部、修補已就緒。** #24 以 head `0379b64c892e9f664e07c879a16f9acaa164be88` 合併後就關閉」 ✅ |
| 內容與審查部實測 | §4 表 parents `80fc321`＋`2a02819`、remerge-diff 空、只 `M README.md`、blob `fbdfb496…`、run 39726831602 四 success，與複審 #01 §一、§二相符 ✅ |

- 判定：**R1 已處理**。註：§1「head：`80fc321…`」為 v0.1 原文保留，已由 §4 補充，不構成矛盾（S2）。

## 三、R2（維運部）

### 3.1 證據檔（快照 `docs/09-cicd-env/evidence/cr-004/`，只在快照，不在 repo）

| 檔案 | sha256（自算） | 要求／宣稱 | 判定 |
|---|---|---|---|
| `after.json` | `1404dab6a8767b6e6235d7c1f9ccadf8dd309cc97606f15d53c2a3952b7ccefa` | `1404dab6…` | ✅（＝複審 #01 暫存區值，未變） |
| `after-1-partial.json`（中間版） | `9b4fffe12602411d24e41258a8e02075f4f634a1636419f0e5311ac5cf6d4717` | `9b4fffe1…` | ✅ |
| `change-payload.json` | `c6b261ac4aee12e4c1f99f1b65d3ab27f7daaff12931743b8e0854a60676e02e` | change `c6b261ac` | ✅ |
| `rollback-payload.json` | `dfcae65fd02fd144ca84b84ff01c4e453f57c46461270aa5722bb527a1aa8683` | rollback `dfcae65f` | ✅ |
| `ruleset-24629854-before-2026-10-22.json` | `af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12` | snapshot `af1c9e4a` | ✅（＝`before.json`） |
| `ac1-5-record-v0.2.md` | `040fced08f6e09e32060cb7ffdf717913b31b8ed4e8bdb7ae9c2da7ea9905b55` | `040fced0…5b55` | ✅ |
| 其他 | `diff-step4.txt` `6f680c07…`、`diff-step4-1-partial.txt` `3f0536e9…`、`805-before/after.json` `ba4cad93…`（相同）、`ui-screenshot.png` `c8d1f4eb…`、`ui-screenshot-2.png` `368c9e7e…`、`pre-pr-list.json` `3d3610b6…` | 與複審 #01 暫存區值相同 | ✅ |

- 另見 `ac1-5-record-v0.1.md`（`1e12b7f3…`，被 v0.2 取代）、`pr24-*.json`／`pr24-first-test-record-v0.1.md`（AC-6～9 用，本次不核）。

### 3.2 維運判定原件 `ac1-5-record-v0.2.md`

| 檢查 | 原文 | 判定 |
|---|---|---|
| 維運部名義 | 「判定者：維運部（本檔為維運部 AC-CR004-1～5 逐項判定原件…）」 | ✅ |
| AC-1 | 「pre-pr-list.json（13:02:05）、before.json af1c9e4a…（＝附件快照）、diff-step2.txt 空 → 通過」 | ✅ |
| AC-2～4 | 「after.json（13:05:58 版）jq -S diff 僅兩欄位 false→true；enforcement active、bypass_actors 不變；24629805 .rules 前後相同 → 通過」 | ✅（三項合列一行，各有對應證據；S3） |
| AC-5 | 「ui-screenshot-2.png 含網址列與勾選 → 通過」 | ✅ |
| 13:04:53 誤關 | 「2026-10-22 13:04:53 首次存檔時 UI 將 require_extra_approval_for_unattributed_changes 一併改為 false…**違反 AC-CR004-3**」 | ✅（與 `after-1-partial.json` updated_at 13:04:53.307 相符） |
| 13:05:58 勾回 | 「13:05:58 勾回修正，誤關約 65 秒」 | ✅（與 `after.json` updated_at 13:05:58.537、API 現況相符） |

- 判定：**R2 已處理**。

## 四、必改清單

無（R1、R2 已處理）。

## 五、S 項（不阻擋）

| 編號 | 內容 | 建議負責 |
|---|---|---|
| S1 | 20 號 v0.1.1 與 `evidence/cr-004/` 全部只在快照，未入 repo；合併後另開文件 PR 收入，sha256 不得變 | 安全部；維運部 |
| S2 | 20 號 §1 仍寫舊 head（v0.1 原文），下版可加註「見 §4」 | 安全部 |
| S3 | ac1-5-record v0.2 無 frontmatter／版本表；AC-2～4 合列；未寫 after.json sha256；「判定時間 13:06」早於檔案時間 13:20（v0.1 13:19），建議註明判定與成文時間 | 維運部 |
| S4 | AC-CR004-6～9 待品保判定（不列合併前提） | 品保部 |
| S5 | 複審 #01 S7：仍須使用者本人 admin bypass；被合併 head 須為 `0379b64` | 總協調；使用者 |
| S6 | 初審 S1～S3、S5、S6；複審 #01 S8、S9 沿用 | 研發部；維運部 |

## 六、最後核對

13:21:05（UTC+8）`ls-remote` 見 §一。引號稽核：「」內文字逐字取自 20 號 v0.1.1、ac1-5-record-v0.2.md。

---

## 附錄 A：「可合併」留言定稿

> **未發布，待使用者核准。** 發布前一刻須再以 `ls-remote`／API 確認 head 與 main 均未變；任一變動即作廢。

```markdown
## 審查部 G4 書面結論：可合併

- PR：#24 `docs: README 431 門檻與例外回應措辭對齊 API 說明 (ENG-021)`（ENG-021；NFR-002、NFR-006；SEC-012；DR-G4-01；G4-PR19 S8～S10；SCR-046）
- **綁定 head：`0379b64c892e9f664e07c879a16f9acaa164be88`**（原 commit `80fc321a2d363871bfdf9b50aa89ca10a3be359d` 未改寫＋同步 merge commit）
- **綁定 base（＝審查當下 main）：`2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`**（#23 合併後）
- 四項 required checks 綠燈 run：https://github.com/444112029012/short-url/actions/runs/39726831602 （head_sha＝上列 head；check runs 恰 4 筆）
  - `secrets-gitleaks`：success（68 commits scanned，no leaks found）
  - `sast-gosec`：success（Go 1.27.2；gosec h1 OK）
  - `sca-govulncheck`：success（No vulnerabilities found；dependency graph compare HTTP 200）
  - `unit-test`：success
- 差異核對：
  - main...head 只有 `README.md`（+3／−3，blob `fbdfb496b2ff44a31907abbc8d572447691cf1cd`，與初審相同）；未動程式碼、測試、`.github/`、`go.mod`／`go.sum`。
  - 同步 merge commit `0379b64` parents＝`80fc321`＋`2a02819`；remerge-diff 為空。
  - merge ref `100dc6d` 的 tree＝head tree `2fc3c9008500f37e4487c6adcaad7b9a99b9d00d`。
  - README 第 55、56、101 行與 main 程式及 API 說明 §11 一致（初審實跑）。
- 前置：CI 全綠；純文件 PR，非高風險；安全部書面確認 `docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md` v0.1.1（sha256 `b10707d0c28eb2631e85f65047332adeab7541e0ba12ace18466ca79b6c719e9`）同意合併、僅適用於上列 head；SCR-046 合併後關閉。
- CR-004：ruleset 24629854 `dismiss_stale_reviews_on_push`／`require_last_push_approval`＝true；維運部 AC-CR004-1～5 判定 `docs/09-cicd-env/evidence/cr-004/ac1-5-record-v0.2.md`（sha256 `040fced08f6e09e32060cb7ffdf717913b31b8ed4e8bdb7ae9c2da7ea9905b55`）通過（含 13:04:53 誤關、13:05:58 勾回註記）；AC-6～9 待品保判定，非合併前提。
- G4 七條件：1、4、5、6、7 達成；2、3 不適用。必改：無。非必改見紀錄，不阻擋。
- 時效：本結論只適用於上列 head 與 base。**head 或 main 任一變動即失效**（D-10：規則一 strict=true），須重新與 main 同步、重新取得四綠後再送審查部。
- 合併程序（引用 `docs/99-changelog/decision-log.md` 與 `docs/06-security/05-residual-risk-register-v0.1.md` 原文）：
  - D-08（decision-log）：「選 A（Rulesets）並接受 RR-006（中），G5 前再評」
  - RR-006 接受條件（05 §2.3）：「規則一四 required checks＋禁 force push／刪除**無 bypass**；規則二僅 Repository admin、僅限 PR；每筆 admin bypass 合併須於 PR 留言附審查部「可合併」結論＋綠燈連結；使用者本人合併」
  - D-09（decision-log）：「記偏差；確認由使用者本人合併；RR-006 維持接受並補強流程（總協調先確認 PR 上已有審查部「可合併」才通知使用者「現在可以按合併了」）；G5 前再評」
  - D-10（decision-log）：「合併前分支須與 main 同步，舊綠燈不得用於合併」；「後合者先 update branch 並重取四綠與審查部「可合併」」
- 紀錄：`docs/gates/gate-G4-PR24-2026-10-22.md`（初審）、`docs/gates/gate-G4-PR24-2026-10-22-recheck-01.md`（同步後複審 #01）、`docs/gates/gate-G4-PR24-2026-10-22-recheck-02.md`（同步後複審 #02；2026-10-22，Asia/Taipei）
```

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v1 | 2026-10-22 | 審查部 | 同步後複審 #02：head 0379b64／main 2a02819 未變；R1 安全部 20 號 v0.1.1（b10707d0…）綁 0379b64、SCR-046 狀態已寫；R2 evidence/cr-004 after.json 1404dab6、中間版 9b4fffe1、change c6b261ac、rollback dfcae65f、snapshot af1c9e4a、ac1-5-record-v0.2 040fced0 均自算相符，維運部逐項判定且載明 13:04:53／13:05:58；所有檔只在快照未入 repo（S1）；通過；無必改；附錄 A 定稿（未發布） |
