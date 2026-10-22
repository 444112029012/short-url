# CR-004 第一個 PR（#24）實測紀錄 v0.1
- 文件負責：維運部；建立：2026-10-22 13:19（UTC+8）；狀態：草稿（待品保／審查部）
- 取證方式：維運部唯讀 gh 讀取，未變更任何設定

## 同步前（AC-CR004-6/7 基準）
- 時間：2026-10-22 13:02:05（UTC+8），取自 pre-pr-list.json（非同步當下重取）
- #24 head 80fc321a2d363871bfdf9b50aa89ca10a3be359d；reviewDecision REVIEW_REQUIRED；無核准

## 同步後（13:18 讀取）
- head 0379b64c892e9f664e07c879a16f9acaa164be88；state OPEN；mergeStateStatus BLOCKED；reviewDecision REVIEW_REQUIRED；reviews []
- run 39726831602 四項 checks 全 pass（pr24-checks.txt）
- BLOCKED 原因判讀：尚無核准（REVIEW_REQUIRED），與 CR-004 前相同之預期狀態；未見與 review 規則相關之錯誤訊息
- 時間軸（pr24-timeline.json）：committed 80fc321a、committed 0379b64c；無 review、無 review_dismissed 事件
- ruleset 24629854 讀取與 CR-004 變更後快照一致（除 updated_at）

## 聲明
- sync 時無核准、無可撤銷對象，reviewDecision 維持 REVIEW_REQUIRED。
- **未實測撤銷行為。**
- **未實測最後一次 push 須他人核准。**

## 尚未完成
- AC-CR004-8：#24 尚未合併（mergedAt null），bypass 合併紀錄與 head／diff 待補。
