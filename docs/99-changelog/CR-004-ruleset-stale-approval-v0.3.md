---
文件：變更請求 CR-004
版本：v0.3
日期：2026-10-22
狀態：已提出（草案 v0.3，回應審查部退回 docs/gates/review-CR-004-2026-10-22.md；使用者對本 CR 之簽核待取得）
提出：維運部（總協調登錄）
專案代號：SHORTURL
關聯：安全部 docs/06-security/22-cr3-r8-stale-approval-residual-risk-v0.1.md（sha256 ca0d6106e697bd4e1415be8c464016baecd3875ca65147c6e038409c38f3f316，2026-10-22 已核對）；RR-011（預定）；安全部 docs/06-security/26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md（sha256 869bc2152edf3c0336967a65cf4c4f2926e23e4e056c19b6d99d972c41c57d53，2026-10-22 09:58 重算相符；**結尾為 7d53**，送審單先前所記「3d57」有誤）；審查部 CR-004 審查紀錄 docs/gates/review-CR-004-2026-10-22.md；SCR-042；D-08／D-09／D-10；審查部 CR3-R8／S12（docs/gates/review-CR-003-2026-10-22.md 第 197、215 行）；CR-003
時區：Asia/Taipei（UTC+8）
---

# CR-004｜Ruleset main-pr-review 開啟「核准過期即撤銷」與「最後一次 push 須經核准」

## 0. 申請欄位

- 提出者／日期：維運部（總協調登錄）／2026-10-22
- 類型：一般變更（環境／規則變更；屬資安控制**強化**，非減弱）。決策：使用者本人 2026-10-22 09:51（UTC+8）決定採行（經總協調轉述）；**使用者對本 CR 文件本身的簽核尚待取得**。
- 變更內容（僅 Ruleset 24629854 `main-pr-review` 的 pull_request 規則兩個參數）：
  - `dismiss_stale_reviews_on_push`：false → **true**
  - `require_last_push_approval`：false → **true**
  - 其餘一律不變：bypass_actors（actor_id 5 RepositoryRole／admin，bypass_mode pull_request）、enforcement active、conditions ~DEFAULT_BRANCH、deletion、non_fast_forward、required_approving_review_count 1、required_reviewers []、require_code_owner_review false、require_extra_approval_for_unattributed_changes true、required_review_thread_resolution false、allowed_merge_methods [merge, squash, rebase]；Ruleset 24629805 `main-required-checks`（四項 required checks、strict=true）不在本 CR 範圍。
  - **編號更正註**：原申請寫「ruleset 24629805 main-pr-review」，經 2026-10-22 gh api 讀取，24629805 為 `main-required-checks`（required_status_checks，strict_required_status_checks_policy=true）；`main-pr-review` 正確 id 為 **24629854**，本 CR 一律採 24629854。
- 原因與效益：安全部 22 號檔評等殘餘風險為**中**：核准後再 push，原核准仍有效，「合併的 head＝核准的 head」僅靠 D-09 人工核 head（SCR-042）。開啟後由 GitHub 對「有實際審查者核准」的情境強制綁定；**本 CR 不修補 RR-011**：依安全部 26 號檔（sha256 869bc2152edf3c0336967a65cf4c4f2926e23e4e056c19b6d99d972c41c57d53），目前合併實際走 bypass，bypass 時這兩個旗標和「至少 1 人核准」都不生效，所以 RR-011 維持**中**、SCR-042 不關閉；22 號檔「修補後可降為低」以 26 號檔為準。本 CR 只是預先強化控制，有實際其他審查者時才會生效。
- 不變更的後果：RR-011 維持中，須使用者書面接受長期風險（22 號檔 §4）；審查部 S12 持續追蹤。
- 影響的需求編號：無 REQ／NFR 變更；SEC：相關殘餘風險 RR-011、SCR-042（不新增 SEC 編號）。
- 期望完成時間：本 CR 核准後，建議於 **#23 合併完成後、#24 同步 main 之前**由使用者執行（**待使用者確認**，見 §3.1）。
- 附件／參考：
  - 變更前快照：/workspace/rebuild/cr-004/ruleset-24629854-before-2026-10-22.json（2026-10-22 09:51 後讀取；ruleset updated_at 2026-10-07T13:21:54+08:00）
  - GitHub 文件：https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets ；https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/creating-rulesets-for-a-repository ；REST：https://docs.github.com/en/rest/repos/rules#update-a-repository-ruleset

## 1. 兩個設定的作用（依 GitHub 文件）

| 設定（UI 名稱） | 文件描述（available-rules-for-rulesets） |
|---|---|
| `dismiss_stale_reviews_on_push`（Dismiss stale pull request approvals when new commits are pushed） | GitHub 記錄核准當下的 diff；若 diff 改變（例如有人 push 新 commit、按 **Update branch**，或相關 PR 合入 target branch），核准即被撤銷為 stale，須重新核准才能合併。 |
| `require_last_push_approval`（Require approval of the most recent reviewable push） | 須由「最後一次 push 者以外」的人核准，即至少有一位其他審查者核准過最新變更；stale 核准不會因此被撤銷，但需覆蓋最後一次 push。 |

文件另註：勾選任一項時，(a) 手動建立 merge commit 直接推到受保護分支會失敗，除非內容與 GitHub 產生的合併完全一致；(b) 若 merge base 在核准後引入新變更，核准也會被視為 stale。

## 2. 評估一：對使用者 admin bypass 合併流程的影響

- 現況：repo 只有一個真人帳號（444112029012），PR 作者即潛在核准者；GitHub 不允許作者核准自己的 PR，因此依 D-09，審查部回報「可合併」後由使用者本人以 admin bypass 合併。
- bypass 行為：bypass_actors 為 admin、bypass_mode=pull_request。依 creating-rulesets 文件，「For pull requests only」的 actor 必須透過 PR 變更，「然後可選擇 bypass 任何分支保護並合併該 PR」。因此開啟兩旗標後，**admin 仍可在 PR 上以 bypass 合併**；沒有核准的狀態本來就靠 bypass，故流程不變。
- 實際作用：兩旗標主要在**有真正其他審查者核准**時發揮（例如未來新增協作者，或審查部依 D-12 以其他身分留核准）。在目前單帳號＋bypass 模式下，GitHub **不會**替使用者強制「合併 head＝審查部認可的 head」，因為 bypass 本身跳過核准要求。
- 殘餘風險：bypass 合併時沒有工具強制綁定 head；審查部「可合併」結論是留言而非 GitHub approval（D-12 允許以使用者帳號留審查結論，但作者無法 approve 自己的 PR），故不會被 dismiss 機制追蹤。→ **SCR-042／D-09 合併前人工核對 head 仍必須保留**：合併紀錄寫明「合併 head＝審查部認可 SHA」及 `git diff <認可 SHA> <merge>` 結果。安全部「修補後可降為低並關閉 SCR-042」之前提在單帳號 bypass 下可能不成立，安全部 26 號檔已確認：SCR-042 不關閉、D-09 必須保留。
- 未驗證：(1) bypass 合併時 UI 是否仍顯示 dismissed 核准紀錄、bypass 是否寫入 audit log 的細節，本次未實測；(2) 文件未明載 pull_request 模式 bypass 與這兩旗標的交互細節，以上依「bypass any branch protections」推論；(3) require_extra_approval_for_unattributed_changes 的作用本次未查文件，維持原值；(4) **與 required_approving_review_count 的相依（無法查證）**：官方文件僅對 unattributed 該項明寫「count=0 時無作用」，對 `dismiss_stale_reviews_on_push`、`require_last_push_approval` 在 count=0 時是否失效未明載（推論：前者只撤銷既有核准；後者在 count=0 時可能仍形成需一人核准的效果，未查證）。本案 count=1 且本 CR 不變更 count，故不影響本次；日後若調整 count 須重新評估本 CR 效果。

## 3. 評估二：對進行中 PR #23～#26 的影響（2026-10-22 gh 讀取；main=363f1692；09:58 UTC+8 重取見下）

| PR | 標題 | head 分支／SHA | base | reviewDecision | 核准數 | 與 main 關係 | mergeState |
|---|---|---|---|---|---|---|---|
| #23 | docs: 安全審查 08 v0.1.5 與 06-security README v0.22 (ENG-011) | docs/ENG-011-security-08-v015-readme-v022／c3155b3b | main | REVIEW_REQUIRED | 0（無任何 review） | ahead 2、behind 0 | BLOCKED |
| #24 | docs: README 431 門檻與例外回應措辭對齊 API 說明 (ENG-021) | docs/ENG-021-readme-followup／80fc321a | main | REVIEW_REQUIRED | 0 | ahead 1、behind 2（需同步） | UNKNOWN（v0.2 時點；09:58 重取為 BEHIND） |
| #25 | fix: CI 強化 SCR-032／019／021／004 (ENG-022) | fix/ENG-022-ci-hardening-scr032／bf190c4c | main | REVIEW_REQUIRED | 0 | ahead 1、behind 2（需同步） | BEHIND |
| #26 | test: 補 OpenAPI 契約測試 CT-* (ENG-013) | test/ENG-013-openapi-contract-tests／e224253d | main | REVIEW_REQUIRED | 0 | ahead 1、behind 0 | BLOCKED |

- 四個 PR 目前**皆無核准**，變更當下不會撤銷任何既有核准。
- strict=true（D-10）下，後合者須先 update branch；同步 merge 是 push，開啟後會使當時存在的核准變 stale。#24、#25 已落後 main 2 個 commit；任一 PR 合併後，其餘 PR 亦會落後。
- 因目前一律以 admin bypass 合併，實務影響是：審查部「可合併」須在**最後一次同步後的 head** 上重新給出（D-10 已要求重取四綠與「可合併」），工作量不變。
- 2026-10-22 09:58（UTC+8）唯讀重取（`gh pr list -R 444112029012/short-url --json number,headRefOid,mergeStateStatus,reviewDecision`；main 仍為 363f1692e56228b39c70c4dfe1da58976a3c9871）：

| PR | headRefOid | mergeStateStatus | reviewDecision |
|---|---|---|---|
| #23 | c3155b3b6130cf8c171a62536987563cce7b45d0 | BLOCKED | REVIEW_REQUIRED |
| #24 | 80fc321a2d363871bfdf9b50aa89ca10a3be359d | BEHIND | REVIEW_REQUIRED |
| #25 | bf190c4c9292f15f7c152ad340c07d58bac2fa58 | BEHIND | REVIEW_REQUIRED |
| #26 | e224253de29743a956e8c952cafc4eff8acff1d7 | BLOCKED | REVIEW_REQUIRED |

  head 與上表一致；#24 由 UNKNOWN 變 BEHIND（時點差異，非異常）。執行前須再重取一次（§4 步驟 0）。

### 3.1 執行時機建議（**待使用者確認**）

前提：合併順序 #23 → #24 → #25 → #26；每個 PR 於審查部回報「可合併」後由使用者以 admin bypass 合併；strict=true（D-10）下，後續 PR 以一般 merge 同步 main。

- **建議方案 A：#23 合併完成後、#24 同步 main 之前執行本變更，並由 #24 擔任安全部 26 號 §4 條件 2 的「變更後第一個 PR 實測」。**
  - 理由：#24 僅改根目錄 README（ENG-021），風險最低；變更後 #24 必須同步 main（一般 merge＝push），正好產生「同步 push」事件，再以 bypass 合併，可一次走完條件 2 的流程。#23 不受影響（變更前合併）。
  - #24 實測步驟：(1) 變更與 §4 步驟 3 驗證完成；(2) #24 以一般 merge 同步 main（sync push）；(3) 以 `gh pr view 24 -R 444112029012/short-url --json reviews,reviewDecision,headRefOid` 及 PR 時間軸檢查是否有核准被 dismiss；(4) 審查部於同步後 head 重給「可合併」；(5) 使用者以 bypass 合併，依條件 3 留紀錄。
  - **可測與不可測**：#24 目前**沒有任何核准**，故「同步 push 後既有核准被撤銷」**無法在 #24 上實際觀察**，只能記錄「當時無核准、無可撤銷對象、reviewDecision 維持 REVIEW_REQUIRED、無錯誤／阻擋」；條件 2 該半句以「如果當時有核准的話」為前提，此情形依原文屬不適用，須於紀錄中明載「未實測撤銷行為」。可實測的是：(a) 開啟兩旗標後 sync push 不造成新的阻擋或錯誤；(b) bypass 合併仍可完成；(c) ruleset 其餘欄位與 24629805 不受影響。撤銷行為要實測需有第二人核准（單帳號無法自我核准），留待有真審查者時或另案。
- **替代方案 B：#26 合併完成後執行。**
  - 優點：不干擾 #23～#26 合併序列，四個 PR 合併期間零變數。
  - 缺點：第一個實測 PR 要等到下一個新 PR，時間未定，條件 2 驗證延後；且其合併前已無同步窗口可測 sync push（除非新 PR 恰好落後 main）；強化控制延後生效（雖目前 bypass 下實益有限）。
- 維運部建議 **A**；若使用者希望 #23～#26 合併期間不引入任何設定變動，則採 B。**待使用者確認。**
- 若執行前 PR 狀態與上表不同（例如 #24 已同步、#23 尚未合併或出現新核准），停止並重新評估時機。

- 註：#25 標 ENG-022，與 CR-003 擬新增之 ENG-022（分支規範例外）編號衝突，非本 CR 範圍，提請總協調確認。

## 4. 評估三：上線步驟（由使用者本人執行；維運部不變更 ruleset）

0. **執行前重取 PR 狀態**（唯讀）：`gh pr list -R 444112029012/short-url --json number,headRefOid,mergeStateStatus,reviewDecision`，並確認 #23 已合併、#24 尚未同步 main（方案 A）；與 §3 表不符則停止，回報維運部。
1. **變更前快照**（使用者或任一代理以唯讀呼叫）：
   ```bash
   gh api repos/444112029012/short-url/rulesets/24629854 > ruleset-24629854-before-$(date +%F).json
   sha256sum ruleset-24629854-before-*.json
   ```
   與本 CR 附檔比對 rules／bypass_actors 是否一致（`jq -S '{bypass_actors,rules}'` diff）；不一致則停止，回報維運部更新 payload。
2. **變更**（擇一）：
   - ⚠ **警語：PUT 會覆寫整個 ruleset。** 一律以**當下最新 GET** 為基礎產生**完整** payload（name、target、enforcement、conditions、bypass_actors、`rules` 陣列三條全列），不得只送 pull_request 片段或省略頂層欄位；官方未明載省略頂層欄位（例如 bypass_actors）時是保留還是清空（無法查證）。若步驟 1 之 GET 與快照不一致，不得使用本附件 payload。
   - UI：Settings → Rules → Rulesets → main-pr-review → Require a pull request before merging → 勾選「Dismiss stale pull request approvals when new commits are pushed」與「Require approval of the most recent reviewable push」→ Save changes。其他項目不動。
   - gh：`gh api -X PUT repos/444112029012/short-url/rulesets/24629854 --input change-payload.json`，payload 全文（rules 陣列完整，僅兩欄位改變）：
```json
{
  "name": "main-pr-review",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "exclude": [],
      "include": [
        "~DEFAULT_BRANCH"
      ]
    }
  },
  "bypass_actors": [
    {
      "actor_id": 5,
      "actor_type": "RepositoryRole",
      "bypass_mode": "pull_request"
    }
  ],
  "rules": [
    {
      "type": "deletion"
    },
    {
      "type": "non_fast_forward"
    },
    {
      "type": "pull_request",
      "parameters": {
        "required_approving_review_count": 1,
        "dismiss_stale_reviews_on_push": true,
        "required_reviewers": [],
        "require_code_owner_review": false,
        "require_last_push_approval": true,
        "required_review_thread_resolution": false,
        "require_extra_approval_for_unattributed_changes": true,
        "allowed_merge_methods": [
          "merge",
          "squash",
          "rebase"
        ]
      }
    }
  ]
}
```
3. **變更後驗證**：
   ```bash
   gh api repos/444112029012/short-url/rulesets/24629854 \
     -q '[.enforcement, (.bypass_actors|tostring), (.rules[]|select(.type=="pull_request").parameters|tostring)]|.[]'
   ```
   預期輸出：
   ```
   active
   [{"actor_id":5,"actor_type":"RepositoryRole","bypass_mode":"pull_request"}]
   {"required_approving_review_count":1,"dismiss_stale_reviews_on_push":true,"required_reviewers":[],"require_code_owner_review":false,"require_last_push_approval":true,"required_review_thread_resolution":false,"require_extra_approval_for_unattributed_changes":true,"allowed_merge_methods":["merge","squash","rebase"]}
   ```
   另：`gh api repos/444112029012/short-url/rulesets/24629805 -q '.rules|tostring'` 應與變更前一致（四項 checks、strict true）；並以 `jq -S` diff 前後快照，只允許兩欄位與 updated_at 不同。
4. **證據**：UI 截圖須**含瀏覽器網址列**（https://github.com/444112029012/short-url/settings/rules/24629854）與兩個勾選狀態；附 gh 驗證輸出與前後快照 sha256。

## 5. 評估四：回滾

- 方式：`gh api -X PUT repos/444112029012/short-url/rulesets/24629854 --input rollback-payload.json`（由變更前快照產生，兩欄位為 false），或 UI 取消兩個勾選。回滾後以 §4 步驟 3 驗證，預期兩欄位為 false。
- 快照保存：/workspace/rebuild/cr-004/ruleset-24629854-before-2026-10-22.json（sha256 見 §9），正式入庫位置由總協調決定。
- 觸發條件：(1) admin bypass 無法合併已獲審查部「可合併」的 PR；(2) 驗證發現 bypass_actors、required checks 或其他參數被意外改變；(3) 合併流程出現文件未載明的阻擋（例如 GitHub 產生之 merge 被拒）且當日無法排除；(4) 使用者指示回滾。回滾須通知安全部（RR-011 回復為中）與審查部。

## 6. 影響分析（由各角色填寫）

| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | 待確認 | 待確認 | 待確認 | 企劃 |
| 架構／資料／API／UI | 待確認 | 待確認 | 待確認 | 設計 |
| 實作 | 待確認（預期同步 main 後需重取「可合併」，與 D-10 現行一致） | 待確認 | 待確認 | 研發 |
| 資安（威脅模型、SEC） | 附條件同意（26 號檔 v0.1）；RR-011 維持中，SCR-042 不關閉 | 26 號檔 v0.1 | 0 | 安全 |
| 環境／部署／監控 | Ruleset 24629854 兩參數 false→true；無程式、CI、部署變更；bypass 與 required checks 不變；admin bypass 合併仍可用；提供快照、payload、驗證與回滾步驟 | 本 CR；變更前後快照；維運文件（如有 ruleset 設定紀錄）補記 CR-004 | 0.25 人日（文件），執行由使用者約 10 分鐘 | 維運 |
| 測試與追溯矩陣 | 待確認 | 待確認 | 待確認 | 品保 |
| 審查（一致性） | 待確認 | 待確認 | 待確認 | 審查 |

## 7. 風險評估

- 單帳號 bypass 下控制效果有限，可能被誤認為「已自動綁定 head」而放鬆人工核對：中；緩解：SCR-042 保留，直到安全部複評。
- 未來有真審查者時，每次同步 main 都需重新核准，合併節奏變慢：低（安全部 22 號檔已列代價）。
- 手動建立 merge commit 直推被拒（文件 note）：低；現行流程一律以 PR 合併。
- PUT 誤改其他欄位（payload 錯誤）：低；以快照比對與變更後 diff 驗證，可即時回滾。

## 8. 建議與核准

- 建議：**接受**。成本低、可回滾，屬控制強化；但須保留 SCR-042 人工核 head。不宣稱修補 RR-011。
- 安全部同意條件（以下逐字照錄自安全部 26 號檔 §4，原編號 1～5，未改寫）：
  1. 由使用者本人執行變更；變更後以 GET 取得 ruleset，留存快照與 sha256，並確認只有兩個旗標改變、bypass_actors 不變。
  2. 變更後第一個 PR 實測：bypass 合併可行；同步 push 後，既有核准會被撤銷（如果當時有核准的話）。
  3. D-09 人工核 head 照舊。每次用 bypass 合併時，合併紀錄要寫明 bypass、被合併的 head＝安全部確認的 SHA，並附 `git diff <確認 SHA> <merge>` 的結果。
  4. 不得在 CR-004 裡順帶改動 bypass_actors、required checks（ruleset 24629805）或任何其他欄位。
  5. 若要真正降低 RR-011，可另案評估以下兩種做法，並由使用者拍板：(a) 設第二個人工核准帳號、移除或收窄 bypass；(b) 以 CI 檢查合併前的 head 是否等於安全部確認的 SHA。

  來源：docs/06-security/26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md §4（sha256 869bc2152edf3c0336967a65cf4c4f2926e23e4e056c19b6d99d972c41c57d53）。
- 各條落實位置：條件 1 → §0 執行者、§4 步驟 3（GET、jq diff、sha256）；條件 2 → §3.1（#24 實測，含可測／不可測說明）；條件 3 → §2 殘餘風險、§3.1 步驟 (5)；條件 4 → §0 變更內容、§4 步驟 2 警語與步驟 3 對照 24629805；條件 5 → 不在本 CR 範圍，另案由使用者拍板。
- 核准鏈：安全部（資安評估／RR-011 複評）→ 審查部（一致性確認）→ 使用者簽核本 CR → 使用者本人執行變更。

| 角色 | 狀態 | 備註 |
|---|---|---|
| 維運部（提出／當責） | 已提出（2026-10-22） | 本草案 |
| 安全部 | 附條件同意（26 號檔 v0.1，針對 CR v0.1） | 須對 v0.3 書面再確認（見 §11） |
| 審查部 | 待確認 | |
| 使用者（王昱殊） | 決定採行（2026-10-22 09:51，經總協調轉述）；**本 CR 簽核待取得** | 執行者 |

- 核准者／日期／結論：待安全部對 v0.3 再確認、審查部確認及使用者簽核。

## 9. 附件雜湊

- ruleset-24629854-before-2026-10-22.json：sha256 af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12
- change-payload.json：sha256 c6b261ac4aee12e4c1f99f1b65d3ab27f7daaff12931743b8e0854a60676e02e
- rollback-payload.json：sha256 dfcae65fd02fd144ca84b84ff01c4e453f57c46461270aa5722bb527a1aa8683
- 安全部 26 號檔 v0.1：sha256 869bc2152edf3c0336967a65cf4c4f2926e23e4e056c19b6d99d972c41c57d53（結尾 7d53；先前送審單所記 3d57 有誤）
- 以上皆於 2026-10-22 09:58（UTC+8）以 sha256sum 重算。
- 安全部 22 號檔：sha256 ca0d6106e697bd4e1415be8c464016baecd3875ca65147c6e038409c38f3f316（已核對）

## 10. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 維運部 | 初稿；更正 ruleset id 24629805→24629854 |
| v0.2 | 2026-10-22 | 維運部 | 依安全部 26 號檔：不宣稱修補 RR-011（維持中，SCR-042 不關閉）；寫入 5 項同意條件與另案方向 |
| v0.3 | 2026-10-22 | 維運部 | 回應審查部退回（CR4-R2/R5/R6/R7 及 S1～S4、S6），詳下 |

變更摘要：
- v0.1 → v0.2：依安全部 26 號檔改寫原因與效益，不再宣稱修補 RR-011（維持中、SCR-042 不關閉、D-09 保留）；§8 加入同意條件（但為改寫重編，經審查部判定不符）與另案方向。
- v0.2 → v0.3：
  - CR4-R2：§8 改為逐字照錄 26 號 §4 條件 1～5，並標示各條落實位置。
  - CR4-R5：新增 §3.1，建議 #23 合併後、#24 同步前執行，#24 為第一個實測 PR，並比較 #26 合併後方案（待使用者確認）。
  - CR4-R6：§8 核准表安全部改為附條件同意（26 號檔 v0.1）；§6 資安列標明 v0.1。
  - CR4-R7：§9 補 change／rollback payload 與快照全值 sha256。
  - S1：註明 26 號 sha256 結尾為 7d53；S2：§4 加 PUT 全量覆寫警語；S3：§2 補 required_approving_review_count 相依（無法查證）；S4：§3 重取 PR 狀態、§4 加步驟 0；S6：本摘要。
  - CR4-R1/R3/R4 屬他部門，列於 §11。

## 11. 尚待補齊事項（非維運部所屬）

| 編號 | 待辦 | 負責 | 狀態 |
|---|---|---|---|
| CR4-R1 | 安全部對本最終版（v0.3，引用其 sha256）出具書面再確認；26 號檔評估對象為 v0.1 | 安全部 | 待確認 |
| CR4-R3 | 填寫 §6 影響列（可填「不適用」並署名）：研發（同步後重取「可合併」、bypass 合併紀錄格式）、品保、企劃、設計、審查 | 研發部、品保部、企劃、設計、審查部（維運彙整） | 待確認 |
| CR4-R4 | 品保書面意見：變更後驗證（§4 步驟 3）與第一個 PR 實測（§3.1）的通過標準與證據要求 | 品保部 | 待確認 |
| — | 執行時機（§3.1 方案 A／B） | 使用者 | 待使用者確認 |
| S5 | #25 ENG-022 編號與 CR-003 衝突 | 總協調 | 另案 |
