---
文件：07-testing 章節索引
版本：v0.3
狀態：審查中
負責角色：品保部
最後更新：2026-10-07
專案代號：SHORTURL
---

# 07 測試與品保｜章節索引

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-requirements-testability-review-v0.1.md](./01-requirements-testability-review-v0.1.md) | v0.1 | 審查中 | G1：REQ／NFR 逐條可測確認 |
| [02-traceability-matrix-v0.1.md](./02-traceability-matrix-v0.1.md) | **v0.3** | 審查中 | 追溯矩陣（需求＋設計／API＋**TC 實號**；公開方法 UT 100%） |
| [03-gate-G1-qa-precheck-v0.1.md](./03-gate-G1-qa-precheck-v0.1.md) | v0.1.1 | 審查中 | 品保 G1 書面前置結論 |
| [04-design-testability-review-v0.1.md](./04-design-testability-review-v0.1.md) | v0.1.1 | 審查中 | G2 設計可測性審查（R1 已關） |
| [05-gate-G2-qa-precheck-v0.1.md](./05-gate-G2-qa-precheck-v0.1.md) | v0.1.1 | 審查中 | G2 品保前置（可送審；G2-QA-R1 已關） |
| [06-test-strategy-v0.1.md](./06-test-strategy-v0.1.md) | **v0.1** | 審查中 | **G3**：測試策略（金字塔、範圍、工具、NFR／SEC） |
| [07-test-cases-v0.1.md](./07-test-cases-v0.1.md) | **v0.1** | 審查中 | **G3**：測試案例目錄（UT／IT／CT／E2E／CHK） |
| [08-definition-of-done-v0.1.md](./08-definition-of-done-v0.1.md) | **v0.1** | 審查中 | **G3**：Definition of Done |
| [09-dev-task-acceptance-map-v0.1.md](./09-dev-task-acceptance-map-v0.1.md) | **v0.2** | 審查中 | **G3**：研發任務↔驗收 TC（ENG 已回填，對齊任務拆解） |
| [10-gate-G3-qa-precheck-v0.1.md](./10-gate-G3-qa-precheck-v0.1.md) | **v0.1** | 審查中 | **G3**：測試策略與 DoD 前置結論（**可送審**） |

## 編號覆蓋

- REQ-001～REQ-012：可測審查＋追溯＋**TC 實號**（v0.3）
- NFR-001～NFR-008：追溯＋CHK／TC
- SEC-001～SEC-017：追溯＋關鍵 TC／CHK
- 公開方法：單測計畫覆蓋 **100%**（見案例目錄 §1、矩陣 §4）

## 交接單：品保 G3 前置 → 總協調／審查

- 交出角色 → 接收角色：品保部 → 總協調（並請審查部引用前置結論）
- 日期：2026-10-07（Asia/Taipei）
- 版本：章節索引 v0.3；狀態：審查中
- 對應需求編號：REQ-001～012、NFR-001～008、SEC-001～017
- 已知問題：無 G3-QA 阻擋必改；建議項 G3-QA-S1～S5（環境基線、事件鍵、研發正式 ID、CI 待辦、G2 建議殘句）
- 需要下游注意：審查依 `10-gate-G3-qa-precheck`；實作綠燈與 CI required checks 屬後續举证；研發任務對照表待正式清單回填
- **結論**：**可送審**；建議總協調送 G3 審查
- 歷史：G1-R1、G2-QA-R1／G2-R1 已關閉；G2 複審通過

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G1 前置初版 |
| v0.1.1 | 2026-10-07 | 品保部 | G1-R1：追溯納入 SEC-001～015 |
| v0.2 | 2026-10-07 | 品保部 | G2：設計可測、G2 前置、追溯 v0.2 |
| v0.2.1 | 2026-10-07 | 品保部 | 關閉 G2-QA-R1；矩陣 v0.2.1 |
| v0.3 | 2026-10-07 | 品保部 | **G3**：策略／案例／DoD／任務對照／G3 前置；矩陣 v0.3 |
