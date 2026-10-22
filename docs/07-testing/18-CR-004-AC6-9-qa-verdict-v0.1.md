---
文件：CR-004 AC-CR004-6～9 與 #24 首測 品保判定
版本：v0.4
日期：2026-10-22（Asia/Taipei，UTC+8）
作者：品保部（Grok Bot 代擬，唯讀查證）
專案代號：SHORTURL
依據：`16-CR-004-R4-qa-opinion-v0.1.md`，sha256 `ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c`（13:16 重算相符）
綁定 CR：`/workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.3.md`，sha256 `eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4`（13:16 重算相符）
查證時間：v0.1 2026-10-22 13:16；v0.2 13:19；v0.3 14:01；v0.4 14:22（UTC+8），repo `444112029012/short-url`，僅 gh 唯讀查詢
效力：16 號檔或 CR v0.3 或下列證據檔 sha256 改變時，本判定失效
---

# CR-004 AC-CR004-6～9 品保判定

## 000. v0.4 更新：AC-CR004-8 重判（2026-10-22 14:22 UTC+8）

### 補存證據核對（`docs/09-cicd-env/evidence/cr-004/`）

| 檔案 | sha256（重算） | 回報 | 內容核對 |
|---|---|---|---|
| pr24-merged.json | 7039591fb2f592de3d0fbda67a02ba12eb17ff064c9878dbbcdf0af5d8f6c3f4 | 相符 | 與 14:22 即時 `gh pr view 24` 同欄位 `jq -S` diff 為空：MERGED、head 0379b64c、merge 5084d6a9、mergedAt 05:59:49Z、reviews [] |
| pr24-diff-head-merge.txt | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 | 相符 | 0 bytes；與 compare 0379b64c…5084d6a9 files＝0 一致 |
| pr24-merge-record-v0.2.md | 36cb275c972d769799de5f2dc73398373003a87401a601318f08b26c76fd0700 | 相符 | head、merge、parents、diff 空、reviews []、兩項未實測聲明，均與 v0.3 查證一致；寫明合併方式 admin bypass |
| （參考）pr24-merge-record-v0.1.md | 629050c3a9052c346997d740ea2082e4fefb660d264d81cc16310afe8cc04f3e | 未回報 | 合併方式為佔位「待使用者本人確認」；v0.2 僅改該行並加變更註記 |

14:22 複核：main 仍為 5084d6a9；24629854、24629805 與變更後快照 diff 為空。

### 合併方式來源標記

- 「是，用 admin bypass 合併」——**來源：合併紀錄引述使用者原話，品保未直接查證。** 合併紀錄寫明為使用者 2026-10-22 14:20（UTC+8）於 Grok Bot 聊天之原文（回答 t257s1017），經 Grok Bot 轉述給維運，維運未直接見證；依使用者指示，引用 Grok Bot 聊天原話者可接受。
- 技術佐證（品保已查）：規則需 1 核准，合併時 reviews []、REVIEW_REQUIRED 仍合併成功，僅 bypass 可能；rule-suites API 404，無直接稽核紀錄。

### AC-CR004-8：**成立**

合併成功；合併 head＝審查部確認之 0379b64c；`git diff 0379b64c 5084d6a9` 為空，符合同步 merge 預期；合併紀錄寫明 bypass、head 與 diff；`pr24-merged.json` 已存；兩 ruleset 未變。備註：紀錄仍為「草稿」，head 一欄註「待兩部門比對簽註」——不影響本判定（審查部 13:59 PR 留言已綁定同一 head），建議定稿。

### 首測整體結論（16 號檔 §4.4）：**通過**

AC-6、7、8 成立，紀錄載明兩項未實測；無失敗條件；AC-9 不適用（未回滾）。N-1～N-3 撤銷與最後一次 push 核准仍未實測，RR-011 維持中、SCR-042 不關閉；head 綁定仍靠 D-09 人工核對。

## 00. v0.3 更新：AC-CR004-8 與首測整體結論（2026-10-22 14:01 UTC+8）

### 查證（gh 唯讀）

| 項目 | 指令 | 結果 |
|---|---|---|
| 合併狀態 | `gh pr view 24 --json state,mergedAt,mergeCommit,headRefOid,mergedBy,reviewDecision,reviews` | MERGED；mergedAt 2026-10-22T05:59:49Z（13:59:49 UTC+8）；mergeCommit `5084d6a95a0ac1d1a3884e990089fa742822c3b1`；headRefOid `0379b64c892e9f664e07c879a16f9acaa164be88`；mergedBy 444112029012；reviewDecision REVIEW_REQUIRED；reviews [] |
| merge commit | `gh api repos/…/commits/5084d6a9…` | parents＝`2a028193…`（#23 合併後 main，與審查部綁定 base 相同）＋`0379b64c…`；tree `2fc3c900…` |
| head tree | `gh api repos/…/commits/0379b64c…` | tree `2fc3c9008500f37e4487c6adcaad7b9a99b9d00d`，**與 merge commit 相同** |
| diff | `gh api repos/…/compare/0379b64c…5084d6a9…` | ahead 1、behind 0、**files 0**（等同 `git diff 0379b64c 5084d6a9` 為空） |
| main | `gh api repos/…/commits/main` | `5084d6a9…`（合併後 main 未再變動） |
| main checks | `gh api …/commits/5084d6a9…/check-runs` | sast-gosec、unit-test、sca-govulncheck、secrets-gitleaks 皆 success（14:00:48 時 sast-gosec 仍 in_progress，約 20 秒後複查為 success） |
| bypass | 由狀態推定 | 規則需 1 核准，合併時 reviews 為空且 REVIEW_REQUIRED，仍合併成功 → 只可能是 admin bypass。`rule-suites` API 回 404，**無直接 bypass 稽核證據** |
| 「可合併」 | issue comments | 13:59:00 PR 留言「審查部 G4 書面結論：可合併」，綁定 head 0379b64c、base 2a028193、四綠 run 39726831602 |
| ruleset | 24629854 vs after.json、24629805 .rules vs 805-before.json（jq -S） | 均 diff 為空，未變 |
| timeline | `…/issues/24/timeline` | committed×2、commented 13:59:00、merged／closed 13:59:49；無 review_dismissed |

### AC-CR004-8：有條件成立（v0.3；已由 §000 v0.4 改為成立）

- 成立部分：以（推定）admin bypass 合併成功；合併 head＝審查部確認之 0379b64c；merge 與確認 SHA 之 diff 為空（tree 相同），符合純同步 merge 預期；合併後兩 ruleset 未變（§4.2 步驟 7）。
- 條件（未滿足前不得視為完全成立）：
  1. `/workspace/rebuild/cr-004/exec/` 及 `docs/09-cicd-env/evidence/cr-004/` 下**沒有** `pr24-merged.json`、合併紀錄或 `git diff <確認 SHA> <merge>` 輸出（16 號檔 AC-8 證據欄、§4.2 步驟 6、7）。須補存。
  2. 合併紀錄須明寫「以 admin bypass 合併」（16 號檔 AC-8、N-5、D-09）；目前 PR 留言只有審查部「可合併」，沒有寫明 bypass 的合併紀錄。
  3. 因此無合併紀錄可算 sha256。

### 首測整體結論（16 號檔 §4.4）（v0.3，已由 §000 取代）

**有條件通過。** AC-6、AC-7 成立；AC-8 有條件成立（補存合併紀錄＋pr24-merged.json＋diff 輸出後即為成立）；紀錄已載明兩項未實測。未見任何失敗條件，不需回滾，AC-9 不適用。撤銷與最後一次 push 核准（N-1～N-3）仍未實測，RR-011 維持中、SCR-042 不關閉。

### 補充核對
- 審查部留言引用之 `docs/09-cicd-env/evidence/cr-004/ac1-5-record-v0.2.md` sha256 `040fced08f6e09e32060cb7ffdf717913b31b8ed4e8bdb7ae9c2da7ea9905b55` 重算相符；`exec/ac1-5-record-v0.2.md` 同值。v0.2 與 v0.1 差異僅標題與判定者／時間一行。

## 0. v0.2 更新（2026-10-22 13:19 UTC+8）

維運補存證據後重算並複核：

| 檔案（`/workspace/rebuild/cr-004/exec/`） | sha256（重算） | 維運所稱開頭 | 內容核對 |
|---|---|---|---|
| pr24-before-sync.json | 6e0010e8d29b851433ea2a8bc9d58141ed8871d8d49fe32e8a863a8b1c74a408 | 相符 | head 80fc321a、REVIEW_REQUIRED；**非步驟 2 同步前當下重取**，係由 pre-pr-list.json（mtime 13:02:05.98）轉錄，檔內已自行註明；無 `reviews` 欄（以 reviewDecision 推得無核准） |
| pr24-after-sync.json | 32885366068d07a3229309ba7ce33115b6f1884f2f94df5b53c8869c06895831 | 相符 | 與品保 13:19 即時查詢一致；缺步驟 4 之 `commits` 欄，檔內無取得時間（紀錄 md 記 13:18） |
| pr24-checks.txt | 2fe7f4c3dde400deeba9671a09fe507d3ffb679bb01d48b88350a0185c8563dc | 相符 | 四項 pass，與 13:19 `gh pr checks 24` 輸出 diff 為空 |
| pr24-timeline.json | 6e16984e0181bd7b5bf83cb5cea51c3db007479fa31925cdd14de12fa3ff2d28 | 相符 | 僅兩筆 committed，無 review_dismissed；與即時查詢一致 |
| ruleset-24629854-at-pr24-test.json | 1404dab6a8767b6e6235d7c1f9ccadf8dd309cc97606f15d53c2a3952b7ccefa | 相符 | 與 after.json 位元相同；與即時 GET `jq -S` diff 為空 |
| pr24-first-test-record-v0.1.md | 3cbe705e5bbc1c8e5c00a5f959dc9c416bb0b3f096931e99da564be7a6ebf99b | 相符 | 含「sync 時無核准、無可撤銷對象，reviewDecision 維持 REVIEW_REQUIRED」「未實測撤銷行為」「未實測最後一次 push 須他人核准」；狀態草稿 |
| ac1-5-record-v0.1.md | 1e12b7f3ac25b05ca2dfe802d11c2c42b242fd0a7d60e44c05599482bd0529b1 | 相符 | 含 13:04:53～13:05:58 誤關約 65 秒、違反 AC-3、兩份快照皆保留 |

13:19 複核 #24：OPEN、head 0379b64c892e9f664e07c879a16f9acaa164be88、REVIEW_REQUIRED、BLOCKED、reviews []、mergedAt null；24629854 與 after.json 一致；24629805 .rules 與變更前一致。

小瑕疵（不影響判定，建議下一版修）：紀錄 md 寫 ruleset「一致（除 updated_at）」，實際連 updated_at 都相同；before-sync 為轉錄而非重取，after-sync 缺 commits 欄與取得時間。

## 1. 結論（v0.2，AC-8 已由 §00 v0.3 更新）

**#24 首測仍未完成，待 #24 合併。** AC-6、AC-7 已成立；AC-8 待 #24 以 bypass 合併後判定；AC-9 未觸發回滾，不適用。未見 §4.4 失敗條件。

| AC | 判定 | 摘要 |
|---|---|---|
| AC-CR004-6 | **成立** | 新 head 四綠（pr24-checks.txt）；BLOCKED 為 REVIEW_REQUIRED 之既有預期狀態；pr24-after-sync.json 已補存且與即時查詢一致 |
| AC-CR004-7 | **成立** | pr24-first-test-record-v0.1.md 三句聲明字句存在，與 timeline／after-sync 事實相符 |
| AC-CR004-8 | **無法驗證（待 #24 合併）** | mergedAt null |
| AC-CR004-9 | **無法驗證（不適用）** | 未觸發回滾 |

（以下 §1 舊表、§2～§7 為 v0.1 內容，保留供追溯；AC-6、7 以 v0.2 為準。）

## 1-v0.1. 結論（已被 v0.2 取代）

| AC | v0.1 判定 |
|---|---|
| AC-CR004-6 | 有條件成立（證據檔缺） |
| AC-CR004-7 | 不成立（紀錄 md 缺） |
| AC-CR004-8 | 無法驗證 |
| AC-CR004-9 | 無法驗證（不適用） |

## 2. 重新執行之驗證（13:16 UTC+8）

| 項目 | 指令 | 結果 |
|---|---|---|
| #24 狀態 | `gh pr view 24 -R 444112029012/short-url --json state,headRefOid,reviewDecision,mergeStateStatus,reviews,commits,statusCheckRollup,mergeCommit` | OPEN；headRefOid `0379b64c892e9f664e07c879a16f9acaa164be88`；reviewDecision REVIEW_REQUIRED；mergeStateStatus BLOCKED；reviews []；commits 2（80fc321a、0379b64c merge origin/main）；mergeCommit null |
| checks | `gh pr checks 24 -R …` | sast-gosec pass、sca-govulncheck pass、secrets-gitleaks pass、unit-test pass（皆 run 39726831602） |
| run | `gh run view 39726831602 --json headSha,conclusion,createdAt` | headSha 0379b64c…；conclusion success；created 13:13:55（UTC+8） |
| timeline | `gh api repos/…/issues/24/timeline` | 只有兩筆 committed 事件，**無 review_dismissed** |
| Ruleset 24629854 | `gh api repos/…/rulesets/24629854` | `dismiss_stale_reviews_on_push=true`、`require_last_push_approval=true`、`require_extra_approval_for_unattributed_changes=true`；updated_at 2026-10-22T13:05:58.537+08:00；與 `after.json` 做 `jq -S` diff **為空** |
| 前後差異 | `jq -S 'del(.updated_at)'` before.json vs 即時 GET | 只有兩欄位 false→true |
| Ruleset 24629805 | `gh api …/rulesets/24629805 -q .rules` vs `805-before.json`（jq -S） | 相同 |
| #23 | `gh pr view 23 --json state,mergedAt` | MERGED，12:09:18（UTC+8） |
| 誤改期間活動 | `gh api repos/…/commits?since=…05:00Z`、pulls 依 updated 排序 | 13:00 後 main 無新 commit；13:04:53～13:05:58 期間無 PR 更新 |

## 3. 誤改事件（require_extra_approval_for_unattributed_changes）判定

- `after-1-partial.json`（9b4fffe1…）：updated_at 13:04:53.307，該欄位 true→**false**；`diff-step4-1-partial.txt` 顯示三欄位變更，違反 AC-3「只動兩欄位」。
- `after.json`（1404dab6…）：updated_at 13:05:58.537，該欄位回到 true；`diff-step4.txt` 只剩兩欄位；即時 GET 與 after.json 一致。**最終設定正確。**
- 影響：誤設窗口約 65 秒；期間無 push、無 PR 更新、無合併（查證如上），無實際受影響之變更。對驗收影響：AC-2～4 以最終狀態判定仍可成立（由維運核定，非本檔範圍）；但此事件須寫入執行紀錄（含兩份快照、時間、更正方式），不得只保留最終快照。建議維運在 AC-1～5 紀錄補註本事件。

## 4. 逐條理由

### AC-CR004-6（有條件成立）
新 head 四項 checks 於 run 39726831602 全部 success；mergeStateStatus 由 BEHIND（16 號檔 §2）變為 BLOCKED，原因為 reviewDecision REVIEW_REQUIRED（需 1 核准，變更前就存在），屬預期狀態，非本變更造成之新阻擋；未見 review 規則相關錯誤。條件：16 號檔要求之證據檔 `pr24-after-sync.json`、`pr24-checks.txt`、`pr24-timeline.json`、以及步驟 1、2 之 `ruleset-24629854-after-<日期>.json`、`pr24-before-sync.json` 均未存於 `exec/` 或 `docs/07-testing/evidence/`，須補存（補存為事後快照，須註明取得時間）。

### AC-CR004-7（不成立）
通過標準為「字句存在」於實測紀錄 md。目前無此檔。事實面已由品保查證：同步時 reviews 為空、timeline 無 review_dismissed、reviewDecision 為 REVIEW_REQUIRED。補上紀錄並載明兩句未實測後即可重判。

### AC-CR004-8（無法驗證）
#24 未合併。另提醒：合併前須由審查部對 0379b64c 重給「可合併」（D-10），合併紀錄須寫 bypass、head 與 `git diff <確認 SHA> <merge>`（D-09，N-5）。

### AC-CR004-9（無法驗證，不適用）
未觸發回滾。

## 5. 與回報不符或未能證實之處

1. 回報「推送前（13:07）reviewDecision=REVIEW_REQUIRED」：無 13:07 之快照檔（`pr24-before-sync.json` 不存在）；可查到的是 `pre-pr-list.json`（13:02）中 #24 為 REVIEW_REQUIRED、head 80fc321a。另 0379b64c 的 committedDate 為 13:07:07，CI 於 13:13:55 觸發，推送時間依 run 判斷約 13:13。
2. 回報「13:13 同步」：與 run 建立時間相符；本機 merge commit 建於 13:07:07。
3. 其餘（head、四項 check、BLOCKED、reviews 空、ruleset 兩欄位 true、after.json／after-1-partial.json sha256 開頭）均相符。
4. 單帳號無法驗證撤銷與最後一次 push 核准：同意，屬 16 號檔 §4.4「不適用（非失敗）」、N-1～N-3；RR-011 維持中、SCR-042 不關閉。

## 6. 證據檔 sha256（13:16 重算）

| 檔案（`/workspace/rebuild/cr-004/exec/`） | sha256 |
|---|---|
| before.json | af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12 |
| after-1-partial.json | 9b4fffe12602411d24e41258a8e02075f4f634a1636419f0e5311ac5cf6d4717 |
| after.json | 1404dab6a8767b6e6235d7c1f9ccadf8dd309cc97606f15d53c2a3952b7ccefa |
| 805-before.json | ba4cad93614911f45d7bf21fe7db94f265c8f5817e8716e6a163906c596b7984 |
| 805-after.json | ba4cad93614911f45d7bf21fe7db94f265c8f5817e8716e6a163906c596b7984 |
| diff-805.txt（空檔） | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| diff-step2.txt（空檔） | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| diff-step4-1-partial.txt | 3f0536e989746fb88db6298a40f6d3d946e63bb3b7b6d819f4fac6bbf4db1cff |
| diff-step4.txt | 6f680c074d0c8273d773b323d6b334e92db259afcdad87901673d97c0468194b |
| pre-pr-list.json | 3d3610b6bc631855c2a70010a25a065fd35ba007d3d66527d7b6e846212aeef5 |
| ui-screenshot.png | c8d1f4ebdb6061cf449a56e1ffe1f6504b1682747688e5daf603579b15251260 |
| ui-screenshot-2.png | 368c9e7ef7b360a00b747b97c8061fc809dcd4eb698db4548c8114586976db29 |
| 16-CR-004-R4-qa-opinion-v0.1.md | ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c |

## 7. S15（不阻擋）：16 號檔待改轉述標記（本版不改 16 號檔，留下一版）

| 行號 | 內容 | 建議 |
|---|---|---|
| 31 | 「總協調轉達『已選定 #23 合併後啟用』」 | 已標「總協調轉述，待確認」，統一改為「agent 轉述，未經使用者原話確認」 |
| 54 | 「啟用時機『#23 合併後』：總協調轉述，待確認」 | 同上統一措辭 |
| 58 | 「入庫位置待總協調確認」 | 確認後須註明來源（使用者原話或轉述） |
| 115 | 「已決定由 #25 改掛 ENG-020」 | 有標「待確認」但未標明為轉述，補「agent 轉述，未經使用者原話確認」 |
| 116 | 「總協調轉述已選 #23 合併後」 | 同 31 行統一 |
| 122 | 「CR 記為『王昱殊』，品保無法核對」 | 標明來源僅為 CR 記載，未經使用者確認 |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | 初版；AC-6 有條件成立、AC-7 不成立、AC-8 無法驗證、AC-9 不適用；#24 首測未完成 |
| v0.2 | 2026-10-22 | 品保部 | 核對維運補存 7 檔（sha256 全相符）；AC-6、AC-7 改為成立；AC-8 待 #24 合併；檔名不變 |
| v0.3 | 2026-10-22 | 品保部 | #24 合併查證；AC-8 有條件成立（缺合併紀錄／pr24-merged.json／diff 輸出、bypass 未書面記載）；首測有條件通過 |
| v0.4 | 2026-10-22 | 品保部 | 核對 pr24-merged.json／diff／合併紀錄 v0.2（sha256 相符）；AC-8 成立；首測通過；bypass 來源標「合併紀錄引述使用者原話，品保未直接查證」 |
