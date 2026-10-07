---
文件：08-engineering 章節索引
版本：v0.2
日期：2026-10-07
狀態：Draft／Ready for Review
負責角色：研發部
作者：研發部
專案代號：SHORTURL
---

# 08 工程實作｜章節索引

> G3 開工就緒：任務拆解、分支／Commit 規範、Definition of Ready。  
> **已對齊品保 DoD／acceptance map**（2026-10-07）。  
> 技術選型（D-04）：Go＋chi／SQLite／base62×8 CSPRNG；零現金預算。

## 文件列表

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-task-breakdown-v0.1.md](./01-task-breakdown-v0.1.md) | v0.1 | Draft／Ready for Review | 任務拆解 ENG-001～021；對應方法規格／REQ／SEC／估時；已對齊 DoD／acceptance map |
| [02-branch-commit-conventions-v0.1.md](./02-branch-commit-conventions-v0.1.md) | v0.1 | Draft／Ready for Review | trunk＋feature；分支／Commit／PR；禁直接 push main |
| [03-definition-of-ready-v0.2.md](./03-definition-of-ready-v0.2.md) | **v0.2** | Draft／Ready for Review | DoR 檢查清單；**已對齊正式 DoD**（`07-testing/08-definition-of-done-v0.1.md`） |
| [03-definition-of-ready-v0.1.md](./03-definition-of-ready-v0.1.md) | v0.1 | 已取代 | 歷史稿；請改用 v0.2 |

## 摘要數據

| 項目 | 值 |
|---|---|
| 任務數 | 21（ENG-001～ENG-021） |
| 總估時 | 約 68 人時（≈ 8.5 人天） |
| 覆蓋 | 方法規格公開方法全列＋OpenAPI 三端點＋CI／README |
| 工作流 | 簡化 trunk＋feature |
| DoD／對帳 | 正式 DoD v0.1；acceptance map 已回填 ENG（見 `docs/07-testing/09-…`） |

## 交接單：研發 G3 工程包 → 總協調／品保／審查

| 欄位 | 內容 |
|---|---|
| 交出角色 → 接收 | 研發部 → 總協調（並請品保確認 DoR↔DoD、安全抽看 CI／SEC 任務） |
| 日期 | 2026-10-07（Asia/Taipei） |
| 已知缺口 | acceptance map 骨架未單列 ENG-001／002／016／017／021（已於 map §4.1 註明掛載）；G1-QA-O5／O7 建議追蹤；矩陣檔名仍為 `02-…-v0.1.md`（內容 v0.3） |
| 下游注意 | G3 審查引用本目錄三份主文件＋品保 DoD／acceptance map；實作自 G4 起依 ENG 開支與 PR |
| 對齊狀態 | **已對齊 DoD／acceptance map**（DoR v0.2；map v0.2 ENG 回填） |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 研發部 | G3 工程章節初版索引 |
| v0.2 | 2026-10-07 | 研發部 | 指向 DoR v0.2；註記已對齊 DoD／acceptance map；更新已知缺口 |
