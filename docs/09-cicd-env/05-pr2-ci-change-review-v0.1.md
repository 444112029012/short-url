---
文件：PR #2 CI 變更維運意見
版本：v0.1
狀態：會簽完成（維運同意 (a)(b)）
負責角色：維運部
最後更新：2026-10-07
專案代號：SHORTURL
對應：PR #2（ENG-007）；SEC-009、SEC-010；01-ci-pipeline v0.2
---

# PR #2 CI 變更維運意見

PR：https://github.com/444112029012/short-url/pull/2（head `ccd3588`）  
CI：[run 37573066229](https://github.com/444112029012/short-url/actions/runs/37573066229)，四項 required check 皆 success。

## 1. 結論：有條件同意（維運面）

| 變更 | 判定 | 說明 |
|---|---|---|
| Go 1.24.x → 1.27.x（`go.mod` 1.27.1） | 同意 | 1.24 線已停止維護、標準庫修補僅在受支援版本；避免 govulncheck 長期紅燈或被迫放寬 |
| gosec v2.21.4 → v2.29.0 | 同意 | 舊版無法以 Go 1.25+ 編譯；`-severity=high` 不變 |
| govulncheck v1.1.3 → v1.8.0 | 同意 | 同上；`govulncheck ./...` 非零即擋不變 |
| dependency review 依賴圖探測後略過 | 有條件同意 | 門檻 v0.2 已將 dependency review 列為非 required 輔助，主擋仍 govulncheck，**未放寬門檻**；但屬 fail-open，須限期：依賴圖開啟後驗證實跑，再改回 fail-closed |
| 四個 job 名稱 | 未變 | 分支保護不需改 |

安全部裁示（維運已同意，2026-10-07）：
1. Go／掃描器升版核准；`go.mod` 與 `setup-go` 一致；掃描器精確釘版（非 @latest）；job 名與 High／Critical 不變。
2. dependency review 略過**不接受常態**。暫行條件：略過須發可見 `::warning::`；`sca-govulncheck` 仍 required；使用者開 Dependency graph 後改回 fail-closed；**暫行最遲 G5 送審前失效**。
3. 跟進：PR #2 目前略過僅 `echo` 說明，尚非 `::warning::`——請研發補一行 `echo "::warning::Dependency graph unavailable; dependency-review skipped; govulncheck remains blocking."`（或等價），再合併。
4. 維運 `01-ci-pipeline` v0.2 與 G3 證據補註已起草，裁示文件定稿後推入 repo。

## 2. 開啟 Dependency graph（使用者操作）

Cursor GitHub App 無 Administration 權限，需使用者在 GitHub 網頁操作：

1. 開 https://github.com/444112029012/short-url/settings/security_analysis
2. 左側選單為 Settings 下的「Advanced Security」（舊版介面叫「Code security and analysis」）。
3. 找到 **Dependency graph**，按 **Enable**（若已顯示啟用則免）。
4. 建議同頁一併啟用 **Dependabot alerts**（免費）。
5. 完成後在 PR #2 重新執行 workflow，或推一個新 commit，確認 `sca-govulncheck` job 中 dependency review 步驟有實際執行。

> GitHub 介面名稱可能隨版本調整，以頁面上的 Dependency graph 項目為準。

## 3. 合併路徑（零預算方案）

現況：`main` 要求 PR＋1 核准＋管理員不可繞過；PR 由擁有者帳號開出，作者不能自核，PR #1、#2 皆 blocked。

建議方案 A（Rulesets 拆兩條），詳見 `01-ci-pipeline-v0.2.md` §5.2：
- 規則 1「main-required-checks」：四項 required checks、禁 force push、禁刪除，**不設 bypass**。
- 規則 2「main-pr-review」：需 PR、1 核准；bypass 名單加 Repository admin，模式「僅限 PR」。
- 刪除舊 Branch protection rule，避免兩邊同時生效仍卡住。
- 補償控制：合併前 PR 留言附審查部書面「可合併」結論＋綠燈連結，由使用者本人合併。

備案 B：核准數改 0（其餘不變）＋同樣補償控制。  
備案 C：邀請第二位真人 collaborator 核准。

放寬「核准」屬控制調整，須安全部會簽、使用者拍板後由使用者在 Settings 操作；設定完成後維運以 API 取證更新 G3 證據。


## 4. 正式會簽確認（2026-10-07）

維運**同意**安全部 (a)(b) 全文：工具升版核准；dependency-review 暫行至 **2026-10-14 與 G5 送審前取較早者**；合併路徑 **C 理想、無第二人則 A、不建議 B**；補償控制與取證義務如上。裁示文件定稿後由總協調呈使用者拍板 Settings。
