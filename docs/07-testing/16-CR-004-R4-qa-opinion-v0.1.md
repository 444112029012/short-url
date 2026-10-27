---
文件：CR-004 品保書面意見（回應 CR4-R4）
版本：v0.1
日期：2026-10-22（Asia/Taipei，UTC+8）
作者：品保部
專案代號：SHORTURL
綁定對象：`/workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.3.md`，sha256 `eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4`（品保部 2026-10-22 10:01 以 `sha256sum` 實際核對相符）
效力：**本意見只對上述 sha256 版本有效；CR-004 內容一旦變更（sha256 不同），本意見即失效，須重新出具。** 附件 change-payload.json（c6b261ac…e02e）、rollback-payload.json（dfcae65f…8683）、ruleset-24629854-before-2026-10-22.json（af1c9e4a…ee12）同日重算相符；附件 sha256 改變亦視為失效。
配套：`15-CR-004-R3-qa-impact-row-v0.1.md`（CR4-R3 品保列）
---

# CR-004 品保書面意見

## 1. 結論

**有條件同意。** 條件：(1) 依 §3 驗收條件 AC-CR004-1～9 執行並留存證據；(2) #24 實測紀錄須明載「未實測撤銷行為」與「未實測最後一次 push 核准」；(3) D-09／SCR-042 人工核 head 照舊，不得因本 CR 放寬；(4) §6 所列 CR 問題中，「§4 步驟 3 預期輸出以字串比對可能誤判」須於執行前改用 `jq -S` 比對（或由執行者知悉並以本意見 AC-CR004-2 為準）。

## 2. 現況（唯讀實查，2026-10-22 10:01 UTC+8）

| 項目 | 實查結果 | 指令 |
|---|---|---|
| Ruleset 24629854 | enforcement active；`dismiss_stale_reviews_on_push=false`、`require_last_push_approval=false`、count 1；bypass actor_id 5／RepositoryRole／pull_request；updated_at 2026-10-07T13:21:54+08:00 | `gh api repos/444112029012/short-url/rulesets/24629854` |
| Ruleset 24629805 | 四項 checks（sca-govulncheck、sast-gosec、secrets-gitleaks、unit-test），strict=true；updated_at 2026-10-07T14:33:28+08:00 | `gh api repos/444112029012/short-url/rulesets/24629805` |
| PR #23～#26 | 皆 OPEN、REVIEW_REQUIRED；#23 c3155b3b BLOCKED、#24 80fc321a BEHIND、#25 bf190c4c BEHIND、#26 e224253d BLOCKED（與 CR §3 一致） | `gh pr list --state all --json …` |
| PR #24 | reviews 0、commits 1 | `gh pr view 24 --json reviews,headRefOid,commits` |
| main | 363f1692e56228b39c70c4dfe1da58976a3c9871 | `git log -1 origin/main` |

### 2.1 啟用時機

- CR-004 v0.3 §3.1 **建議**方案 A（#23 合併完成後、#24 同步 main 之前執行），但明寫「**待使用者確認**」，替代方案 B 為 #26 合併後。
- 總協調轉達「已選定 #23 合併後啟用」：CR v0.3 內文**未記載此決定**，本意見標為「**總協調轉述，待確認**」（須使用者書面確認或 CR 改版記載；CR 改版則本意見失效）。以下 §4 依方案 A 撰寫。

## 3. 驗收條件（AC）

| AC | 條件（可測） | 證據 | 通過標準 |
|---|---|---|---|
| AC-CR004-1 | 執行前狀態：#23 已合併；#24 尚未同步 main、head＝80fc321a…（或與執行前重取一致）、無核准；變更前 GET 之 `{bypass_actors,rules}` 與附件快照 `jq -S` diff 為空 | `pre-pr-list.json`、`before.json`＋sha256、diff 輸出 | diff 為空且狀態符合；否則**停止**（CR §4 步驟 0／1） |
| AC-CR004-2 | 變更後 24629854：`dismiss_stale_reviews_on_push=true`、`require_last_push_approval=true`，enforcement active，bypass_actors 不變 | `after.json`＋sha256；`jq -S` 比對輸出 | 以 `jq -S` 比對值相符（**不以 CR §4 步驟 3 字串逐字比對**，因 gh 實際輸出鍵序為字母序，見 §6-1） |
| AC-CR004-3 | 前後快照差異只有兩欄位與 updated_at（含 name、target、conditions、其他 rules） | `jq -S 'del(.updated_at,…)'` diff 輸出 | diff 恰為兩欄位 false→true，其他無差異 |
| AC-CR004-4 | 24629805 `.rules` 與變更前一致 | 前後 `.rules` 輸出 diff | diff 為空 |
| AC-CR004-5 | 證據完整：UI 截圖含網址列 `https://github.com/444112029012/short-url/settings/rules/24629854` 及兩勾選（若以 UI 執行）；gh 輸出；前後快照 sha256 | 截圖檔＋檔案清單 | 缺任一項即不通過 |
| AC-CR004-6 | #24 sync push 後不產生新阻擋／錯誤：mergeStateStatus 不因本變更出現預期外狀態；四項 checks 在新 head 重跑 | `pr24-after-sync.json`、`gh pr checks 24` 輸出 | 新 head 四綠；無與 review 規則相關之錯誤訊息 |
| AC-CR004-7 | #24 實測紀錄明載：sync 時無核准、無可撤銷對象、reviewDecision 維持 REVIEW_REQUIRED；「**未實測撤銷行為**」「**未實測最後一次 push 須他人核准**」 | 實測紀錄 md | 字句存在 |
| AC-CR004-8 | #24 以 admin bypass 合併成功；合併紀錄寫明 bypass、合併 head＝審查部／安全部確認之 SHA，附 `git diff <確認 SHA> <merge>` 結果（安全部條件 3、D-09） | `pr24-merged.json`、合併紀錄、diff 輸出 | 合併成功且 head 相符；diff 結果與預期一致（預期內容：待確認，依合併方式而定） |
| AC-CR004-9 | 回滾（若觸發）：兩欄位回 false，其餘與變更前快照 `jq -S` diff 為空；通知安全部、審查部 | `rollback-after.json`＋diff | diff 為空（updated_at 除外） |

## 4. #24 首次實測驗證程序（方案 A）

### 4.1 前置條件

1. **#23 已合併**（`gh pr view 23 --json state,mergeCommit` 為 MERGED）。
2. **Ruleset 已啟用**：以 `gh api repos/444112029012/short-url/rulesets/24629854` 讀取設定，`dismiss_stale_reviews_on_push=true`、`require_last_push_approval=true` 作為啟用證據（AC-CR004-2～4 已通過）。
3. #24 尚未同步 main；CR-004 已經使用者簽核（目前**待取得**）。
4. 啟用時機「#23 合併後」：總協調轉述，待確認。

### 4.2 步驟（全部唯讀查詢由任一代理執行；同步與合併由使用者執行）

建議證據目錄：`docs/07-testing/evidence/cr-004/`（入庫位置待總協調確認）。

1. `gh api repos/444112029012/short-url/rulesets/24629854 > ruleset-24629854-after-<日期>.json && sha256sum …`
2. `gh pr view 24 -R 444112029012/short-url --json reviews,reviewDecision,headRefOid,mergeStateStatus > pr24-before-sync.json`
3. 使用者以一般 merge 同步 main（Update branch 或本機 merge 後 push）。
4. `gh pr view 24 … --json reviews,reviewDecision,headRefOid,mergeStateStatus,commits > pr24-after-sync.json`；`gh api repos/444112029012/short-url/issues/24/timeline > pr24-timeline.json`（查是否有 `review_dismissed` 事件）；`gh pr checks 24 > pr24-checks.txt`。
5. 審查部於同步後 head 重給「可合併」（D-10）。
6. 使用者以 admin bypass 合併；記錄合併 head、`git diff <確認 SHA> <merge>`。
7. `gh pr view 24 --json state,mergeCommit,headRefOid > pr24-merged.json`；再 GET 兩個 ruleset，確認未變。

### 4.3 預期結果

- 步驟 4：headRefOid 變為新 SHA；reviews 仍為 0；timeline **無** `review_dismissed`（因無核准可撤銷）；reviewDecision REVIEW_REQUIRED；四項 checks 於新 head 綠燈。
- 步驟 6：bypass 合併成功。
- 步驟 7：兩 ruleset 與步驟 1／變更前一致。

### 4.4 判定

- **通過**：AC-CR004-6～8 全部成立，且紀錄載明兩項未實測。
- **失敗**：bypass 無法合併已獲「可合併」之 PR；sync 後出現與 review 規則相關之新阻擋或錯誤；任一 ruleset 欄位被意外改變；證據缺漏。→ 依 CR §5 觸發條件評估回滾。
- **不適用（非失敗）**：未觀察到核准撤銷、未驗最後一次 push 核准——單帳號下無核准可撤銷，作者亦無法核准自己的 PR。

### 4.5 否定案例（應發生而不得被誤判為通過的情形）

| 編號 | 情境 | 預期 | 單帳號可否實測 |
|---|---|---|---|
| N-1 | 他人核准後再 push 新 commit | 核准轉為 dismissed／stale，未重新核准前非 bypass 合併被擋 | **否**，待第二位審查者（安全部條件 5(a) 另案） |
| N-2 | 最後一次 push 者為唯一核准者 | 該核准不滿足要求，非 bypass 合併被擋 | **否**（作者本就不能 approve 自己 PR） |
| N-3 | 核准後 merge base 引入新變更（他 PR 合入 main） | 核准視為 stale（CR §1 文件註 b） | **否**，同 N-1 |
| N-4 | 手動建 merge commit 直推 main | 被拒（CR §1 文件註 a；另 non_fast_forward／PR 規則） | 不建議實測（會動到 main），列為不測 |
| N-5 | 有人誤以為「已自動綁定 head」而省略 D-09 | 合併紀錄缺 head／diff → AC-CR004-8 不通過 | 可（文件檢查） |

N-1～N-3 須於日後有真審查者時實測，另立紀錄。

### 4.6 回滾後驗證

1. 以變更前快照產生之 rollback-payload（sha256 dfcae65f…8683）或 UI 取消勾選。
2. GET 24629854：兩欄位為 false；`jq -S` 與變更前快照 diff 僅 updated_at 不同。
3. GET 24629805 `.rules` 與變更前一致。
4. 任一開放 PR 以 bypass 合併能力不受影響（以下一次合併觀察，或檢視 mergeStateStatus 無異常）。
5. 留存 `rollback-after.json`＋sha256、diff 輸出；通知安全部（RR-011 回復為中，實際原本即為中）與審查部。

## 5. 對既有品保文件的影響（建議，既有檔未改）

| 文件 | 影響 | 建議 |
|---|---|---|
| `08-definition-of-done`（v0.1） | 條文不需改；C1（測試失敗不得無紀錄合併）、R3 仍適用 | 合併紀錄沿用 D-09 格式；不得因本 CR 刪除人工核 head |
| `07-test-cases`（內容版 v0.2） | 無 TC 新增／修改／作廢 | 建議新增 **CHK-CR-004**（checklist）：「Ruleset 24629854 兩旗標為 true 且其他欄位不變；24629805 不變；第一個 PR 實測紀錄含兩項未實測聲明；bypass 合併附 head／diff」，對應 AC-CR004-1～9 |
| `02-traceability-matrix`（v0.4） | 無 REQ／NFR／SEC 變更 | 建議 §8 流程類任務新增一列：CR-004｜Ruleset stale approval｜TC＝CHK-CR-004｜對應 SCR-042／RR-011（輔助）｜⏳ 待 CR 核准 |
| `09-dev-task-acceptance-map`、`12-G4 precheck` | 無直接影響 | 若 G4 引用合併流程，備註 CR-004 生效後之 #24 實測結果 |

上述新增須待 CR-004 核准後由品保改版（版本遞增）執行。

## 6. 對 CR-004 v0.3 的發現與待確認

1. **§4 步驟 3 預期輸出鍵序**：品保實查 `gh api … -q` 輸出參數鍵為字母序（`allowed_merge_methods` 在前），與 CR 所列預期字串順序不同，逐字比對會誤判失敗；建議改以 `jq -S` 比較值。
2. **§0／§9 引用之安全部 26 號檔 sha256（869bc215…7d53）已過時**：26 號檔已於 2026-10-22 10:01 更新為 v0.1.1（品保重算 `8f4f3c62684b6683e8d2b54c449b87a6e44e1d26195829b6cc36e99dcfe39a8b`），內含對 v0.3 之再確認（附條件同意）。CR §8 核准表與 §11 CR4-R1 仍寫「待確認」，待 CR 改版同步（改版則本意見失效）。
3. **§11 S5**：CR 內文寫「#25 ENG-022 編號與 CR-003 衝突｜總協調｜另案」，**未記載「已決定由 #25 改掛 ENG-020」**；該說法本意見不採用，標「待確認」。
4. **啟用時機**：CR §3.1 為「待使用者確認」；總協調轉述已選 #23 合併後，待確認。
5. **CR 文件簽核**：使用者對本 CR 之簽核待取得（CR §0、§8）。
6. **快照產出位置**：§4 步驟 1 以 `$(date +%F)` 寫在當前目錄，未指定存放路徑；待確認（建議同 §4.2 證據目錄）。
7. **回滾預期輸出**：§5 僅寫「預期兩欄位為 false」，未列完整預期；以 AC-CR004-9 補足。
8. **撤銷／最後一次 push 核准行為**：單帳號不可測（CR §3.1 已載明）；實測延至有第二審查者，時程待確認。
9. **§2 未驗證事項**（bypass 是否寫入 audit log、UI 是否顯示 dismissed、count=0 相依）：本意見不要求驗證，維持待確認。
10. **§8 使用者姓名**：CR 記為「王昱殊」，品保無法核對，待確認。

## 7. 殘餘風險

- RR-011 維持**中**、SCR-042 不關閉（CR §0、安全部 26 號檔）；目前合併走 bypass，兩旗標不生效，head 綁定仍僅靠 D-09 人工核對。
- 品保驗收僅能證明「設定已改、不造成阻擋、可回滾」，**不能證明撤銷機制有效**（N-1～N-3 未測）。
- 誤認「已自動綁定 head」而放寬人工核對之風險（CR §7，中）→ 以 AC-CR004-8／N-5 檢查。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | 初版，綁定 CR-004 v0.3（eae83378…7ed4）；有條件同意 |
