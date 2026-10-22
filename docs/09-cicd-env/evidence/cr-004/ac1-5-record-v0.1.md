# CR-004 AC-CR004-1～5 執行紀錄 v0.1
- 執行：使用者本人以 GitHub UI；核對：維運部唯讀 gh＋jq -S
- AC-1：pre-pr-list.json（13:02:05）、before.json af1c9e4a…（＝附件快照）、diff-step2.txt 空 → 通過
- AC-2～4：after.json（13:05:58 版）jq -S diff 僅兩欄位 false→true；enforcement active、bypass_actors 不變；24629805 .rules 前後相同 → 通過
- AC-5：ui-screenshot-2.png 含網址列與勾選 → 通過
## 偏差註記
- 2026-10-22 13:04:53 首次存檔時 UI 將 require_extra_approval_for_unattributed_changes 一併改為 false（UI 顯示未勾、Preview 欄位），**違反 AC-CR004-3**；13:05:58 勾回修正，誤關約 65 秒。
- 兩份快照皆保留：after-1-partial.json（誤關版）、diff-step4-1-partial.txt、ui-screenshot.png；after.json（修正版）、diff-step4.txt、ui-screenshot-2.png。
- 誤關期間 #24 未同步、無合併，無 PR 受影響（以 pre-pr-list 與時間軸判讀）。
