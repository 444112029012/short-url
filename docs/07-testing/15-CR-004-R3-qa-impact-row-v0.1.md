---
文件：CR-004 §6 影響分析｜品保列（回應 CR4-R3）
版本：v0.1
日期：2026-10-22（Asia/Taipei，UTC+8）
作者：品保部
專案代號：SHORTURL
綁定對象：`/workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.3.md`，sha256 `eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4`（品保部 2026-10-22 10:01 以 `sha256sum` 實際核對相符）
效力：**本列只對上述 sha256 版本有效；CR-004 內容一旦變更（sha256 不同），本列即失效，須重新填寫。**
配套：品保書面意見 `16-CR-004-R4-qa-opinion-v0.1.md`（CR4-R4）
---

# CR-004 §6 影響分析：品保列

## 1. 品保列（依 CR-004 v0.3 §6 表格欄位格式）

| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 測試與追溯矩陣 | 無 REQ／NFR／SEC 變更（CR §0），**不新增、修改或作廢任何 TC**，產品回歸範圍不變（無程式、CI、部署變更，CR §6 維運列）。屬流程控制驗證：需驗收 §4 步驟 3 變更後設定、§3.1 第一個 PR（#24，方案 A）實測、§5 回滾驗證；驗收條件 AC-CR004-1～9 與程序見 16 號檔。限制：單帳號下「核准被撤銷」「最後一次 push 者以外須核准」**無法實測**，只能驗「不造成新阻擋、bypass 合併可完成、其他欄位不變」。D-09／SCR-042 人工核 head 仍為合併 DoD 必要項。 | (1) 16 號檔（本次新增，品保意見）；(2) 建議新增 CHK-CR-004（僅於 16 號檔提出，**既有檔未改**），待 CR-004 核准後再決定是否寫入 `07-test-cases`（CHK 表）、`02-traceability-matrix` §8 流程類任務；(3) `08-definition-of-done` 不需改條文（C1／R3 現行可涵蓋），建議於合併紀錄沿用 D-09 格式；(4) 證據存放：建議 `docs/07-testing/evidence/cr-004/`（待總協調確認入庫位置，CR §5 亦寫「由總協調決定」） | **品保估計**：0.5 人日（意見書與本列 0.25；執行日依 AC 收證、判定與回填 0.25）；不含研發／維運／使用者執行時間 | 品保 |

## 2. 依據（唯讀實查，2026-10-22 10:01 UTC+8）

- Ruleset 24629854：`dismiss_stale_reviews_on_push=false`、`require_last_push_approval=false`、`required_approving_review_count=1`、bypass_actors `[{"actor_id":5,"actor_type":"RepositoryRole","bypass_mode":"pull_request"}]`、updated_at 2026-10-07T13:21:54+08:00（與 CR §0 一致）。
- Ruleset 24629805：四項 required checks、strict=true，updated_at 2026-10-07T14:33:28+08:00。
- PR #23～#26 head 與 mergeState 與 CR §3 09:58 重取表一致；#24 reviews 0、commits 1；main＝363f1692e56228b39c70c4dfe1da58976a3c9871。
- 附件 sha256（change／rollback payload、變更前快照）與 CR §9 相符。

## 3. 待確認

- 證據正式入庫位置（總協調）。
- CHK-CR-004 是否正式登錄（CR 核准後由品保改版 07／02，另走版本遞增）。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | 初版，綁定 CR-004 v0.3（eae83378…7ed4） |
