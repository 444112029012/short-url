---
文件：main 分支 Rulesets 設定步驟（D-08 方案 A）
版本：v0.1
狀態：已完成（使用者於 2026-10-07 設定；ruleset 24629805 main-required-checks、24629854 main-pr-review；舊 classic 保護與舊 ruleset main 已刪；證據 `evidence/rulesets-snapshot-2026-10-07.md`）
負責角色：維運部
最後更新：2026-10-07
專案代號：SHORTURL
對應：D-08、RR-006、09-cr-ci-deviations-ruling-v0.1、01-ci-pipeline-v0.2 §5.2；ENG-020；NFR-004（SI-05）；SEC-009、SEC-010、SEC-011
---

# main 分支 Rulesets 設定步驟

> **狀態：已完成**。使用者已於 2026-10-07 依本指引完成設定：ruleset `24629805`（`main-required-checks`，Bypass list 空）、`24629854`（`main-pr-review`，僅 Repository admin、限 PR）；舊 classic branch protection 與舊 ruleset `main` 已刪除。唯讀 API 快照與截圖證據見 [`evidence/rulesets-snapshot-2026-10-07.md`](./evidence/rulesets-snapshot-2026-10-07.md)。以下步驟保留作為設定紀錄與日後重建參考。

目的：四項掃描仍然誰都不能繞過；只有「1 人核准」允許 repo 管理員經由 PR 繞過，讓單人擁有者能合併 PR。

順序很重要：**先建好兩條規則並確認生效，最後才刪舊的 Branch protection rule**，中間不會有空窗。

> GitHub 介面用詞偶爾調整，以下以目前英文介面為準；名稱略有出入時找意思相同的選項即可。

---

## 規則一：main-required-checks（不可繞過）

1. 開 [Rulesets 設定頁](https://github.com/444112029012/short-url/settings/rules)，按 **New ruleset** → **New branch ruleset**。
2. **Ruleset Name**：`main-required-checks`
3. **Enforcement status**：選 **Active**。
4. **Bypass list**：**保持空白**，什麼都不要加。
5. **Target branches**：按 **Add target** → **Include default branch**（即 `main`）。
6. **Branch rules** 勾選：
   - ✅ **Restrict deletions**
   - ✅ **Block force pushes**
   - ✅ **Require status checks to pass**，展開後按 **Add checks**，逐一輸入並加入下面四項（名稱一字不差）：
     - `secrets-gitleaks`
     - `sast-gosec`
     - `sca-govulncheck`
     - `unit-test`
     
     若每項旁邊可選來源，選 **GitHub Actions**。
     「Require branches to be up to date before merging」與「Do not require status checks on creation」**都不要勾**。
   - ❌ **不要**勾 Require a pull request before merging（放在規則二）。
7. 按最下方 **Create**。

## 規則二：main-pr-review（admin 可經 PR 繞過）

1. 同一頁再按 **New ruleset** → **New branch ruleset**。
2. **Ruleset Name**：`main-pr-review`
3. **Enforcement status**：**Active**。
4. **Bypass list**：按 **Add bypass** → 勾 **Repository admin**（角色）→ 加入後，把它旁邊的模式從 **Always allow** 改成 **For pull requests only**（僅限 PR）。
5. **Target branches**：**Add target** → **Include default branch**。
6. **Branch rules** 勾選：
   - ✅ **Restrict deletions**、✅ **Block force pushes**（與規則一重複無妨，多一層保險）
   - ✅ **Require a pull request before merging**，展開後 **Required approvals** 設 **1**；其他子選項維持預設。
   - ❌ 不要勾 status checks（放在規則一）。
7. 按 **Create**。

## 確認兩條規則生效（刪舊規則前）

1. 回 [Rulesets 頁](https://github.com/444112029012/short-url/settings/rules)，應看到兩條都顯示 **Active**、目標 `main`。
2. 打開 [PR #2](https://github.com/444112029012/short-url/pull/2) 底部合併區：此時舊規則還在，仍會卡住，這是正常的。

## 最後：刪除舊的 Branch protection rule

1. 開 [Branches 設定頁](https://github.com/444112029012/short-url/settings/branches)。
2. 在 **Branch protection rules** 找到 `main` 那條 → **Delete** → 確認。
3. 再看 PR 底部合併區：四項 check 綠燈時，應出現可勾選的「繞過規則合併」選項（類似 *Merge without waiting for requirements to be met (bypass rules)*）。四項 check 若有任何一項沒過，仍然不能合併——這代表規則一正確生效。

## 每次用管理員繞過合併時（補償控制，RR-006）

1. 先確認 PR 留言裡有審查部的書面「可合併」結論。
2. 在 PR 留言貼上四項 check 綠燈的 Actions run 連結。
3. 由你本人按合併。

## 完成後

告訴總協調「Rulesets 設好了」。維運會用 API 取兩條 rulesets 與 main 生效規則的快照，更新 G3 證據。

## 萬一卡住

- 刪了舊規則後 PR 仍無繞過選項：檢查規則二的 bypass 是否為 Repository admin，且模式是 **For pull requests only** 而非空白。
- 綠燈 PR 顯示「required check expected」：檢查規則一的四個名稱是否一字不差，來源是否為 GitHub Actions。
- 任何一步不確定就先停，回報總協調，不要為了合併去刪規則一。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | D-08 方案 A 設定步驟（待使用者操作） |
| v0.1 | 2026-10-07 | 維運部 | 狀態同步（依 G4-PR3-R3）：狀態改「已完成」，記錄 ruleset 24629805／24629854、舊 classic 與舊 ruleset `main` 已刪，指向 `evidence/rulesets-snapshot-2026-10-07.md`；對應補 ENG-020、NFR-004、SEC-009～011；新增修訂紀錄節 |
