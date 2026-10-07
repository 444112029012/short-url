---
文件編號：EVD-G3-ADDENDUM-PR2-20261007
文件：G3 證據補註：PR #2 工具鏈升版的影響
版本：v0.2
狀態：草稿待審查
負責角色：維運部
日期：2026-10-07
專案代號：SHORTURL
對應：ENG-020；NFR-004（SI-05）；SEC-009、SEC-010、SEC-011；D-07、D-08；`g3r1-landing-evidence-2026-10-07.md`
時區：Asia/Taipei（UTC+8）
---

# G3 證據補註：PR #2 工具鏈升版的影響（2026-10-07）

對應：[`g3r1-landing-evidence-2026-10-07.md`](./g3r1-landing-evidence-2026-10-07.md)（G3-R1 落地證據）、D-07、D-08。

## 1. 工具鏈升版

- G3-R1 證據取自 Go **1.24.x**、gosec v2.21.4、govulncheck v1.1.3 的基準（main 當時版本）。
- ENG-007／PR #2 依 D-07（CR-001 核准）改為 Go 1.27.x、gosec v2.29.0、govulncheck v1.8.0；四個 required job 名稱與 High／Critical 阻擋語意不變，required checks 設定不需修改，**G3-R1 證據仍有效**。
- PR #2 合併後，維運以合併後 main 的 push run 補一筆綠燈證據（待合併後補）。

## 2. dependency review 恢復 fail-closed

- D-07 核准的 dependency review 暫行略過（須 `::warning::`；期限 min(2026-10-14, G5 送審前)）已結束：Dependency graph 已開啟，已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 **fail-closed**——compare API 探測非 200 時以 `::error::` 標示並使 job 失敗；回 200 才執行 `dependency-review-action`（`fail-on-severity: critical`）。**暫行略過已關閉**。

## 3. 合併路徑改為 Rulesets

- 依 D-08 改採 Rulesets（方案 A），G3-R1 證據中的 classic branch protection 已刪除。Rulesets 快照（API＋截圖）**已取得**：[`rulesets-snapshot-2026-10-07.md`](./rulesets-snapshot-2026-10-07.md)（ruleset 24629805 `main-required-checks`、24629854 `main-pr-review`）。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | 補註 PR #2 工具鏈升版對 G3-R1 證據的影響（無 frontmatter） |
| v0.2 | 2026-10-07 | 維運部 | 依 G4-PR3 審查建議：S2 補 frontmatter（文件編號、版本、狀態、負責角色、日期、對應）；S1 註明 dependency review 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed、暫行略過關閉；註明 Rulesets 快照已取得並連結；新增修訂紀錄 |
