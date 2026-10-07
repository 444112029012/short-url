# G3 證據補註：PR #2 工具鏈升版的影響（2026-10-07）

對應：`g3r1-landing-evidence-2026-10-07.md`（PR #1）、D-07、D-08。


- 本證據取自 Go **1.24.x**、gosec v2.21.4、govulncheck v1.1.3 的基準（main 當時版本）。
- PR #2 改為 Go 1.27.x、gosec v2.29.0、govulncheck v1.8.0；四個 required job 名稱與 High／Critical 阻擋語意不變，分支保護設定不需修改，**本證據仍有效**。
- PR #2 合併後，維運以合併後 main 的 push run 補一筆綠燈證據；若合併路徑改用 Rulesets（01-ci-pipeline v0.2 §5.2 方案 A），另以 API 取 rulesets 快照更新本證據。
