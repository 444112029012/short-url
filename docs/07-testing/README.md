---
文件：07-testing 章節索引
版本：v0.5
狀態：審查中
負責角色：品保部
最後更新：2026-10-27（Asia/Taipei，UTC+8）
專案代號：SHORTURL
---

# 07 測試與品保｜章節索引

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-requirements-testability-review-v0.1.md](./01-requirements-testability-review-v0.1.md) | v0.1 | 審查中 | G1：REQ／NFR 逐條可測確認 |
| [02-traceability-matrix-v0.1.md](./02-traceability-matrix-v0.1.md) | **v0.5** | 審查中 | 追溯矩陣（v0.5：§7 基準 main f646bee、§7.1 PR #13～#26 追溯、§7.2 開放缺陷、§8 CR-004；sha256 `8a005aefcf24a32fa769c803373e88a3bb411aa71652ff14426feb18907b9932`） |
| [03-gate-G1-qa-precheck-v0.1.md](./03-gate-G1-qa-precheck-v0.1.md) | v0.1.1 | 審查中 | 品保 G1 書面前置結論 |
| [04-design-testability-review-v0.1.md](./04-design-testability-review-v0.1.md) | v0.1.1 | 審查中 | G2 設計可測性審查（R1 已關） |
| [05-gate-G2-qa-precheck-v0.1.md](./05-gate-G2-qa-precheck-v0.1.md) | v0.1.1 | 審查中 | G2 品保前置（可送審；G2-QA-R1 已關） |
| [06-test-strategy-v0.1.md](./06-test-strategy-v0.1.md) | **v0.1** | 審查中 | **G3**：測試策略（金字塔、範圍、工具、NFR／SEC） |
| [07-test-cases-v0.1.md](./07-test-cases-v0.1.md) | **v0.2** | 審查中 | 測試案例目錄（v0.2：§5.1 E2E-05 細分、CHK-ENG-022） |
| [08-definition-of-done-v0.1.md](./08-definition-of-done-v0.1.md) | **v0.1** | 審查中 | **G3**：Definition of Done |
| [09-dev-task-acceptance-map-v0.1.md](./09-dev-task-acceptance-map-v0.1.md) | **v0.3** | 審查中 | 研發任務↔驗收 TC（v0.3：ENG-022／CR-003 AC-022-1～6；T-SEC-01 改指 E2E-05-01～11） |
| [10-gate-G3-qa-precheck-v0.1.md](./10-gate-G3-qa-precheck-v0.1.md) | **v0.1** | 審查中 | **G3**：測試策略與 DoD 前置結論（**可送審**） |
| [11-s6-eng011-disposition-v0.1.md](./11-s6-eng011-disposition-v0.1.md) | v0.1 | 審查中 | **G4 S6**：PR #18 TestENG011_* 19 項處置（11 納入、8 合併、0 不納入）；E2E-05 最終處置 |
| [12-gate-G4-qa-precheck-v0.1.md](./12-gate-G4-qa-precheck-v0.1.md) | v0.1 | 審查中 | **G4**：#13～#20 品保前置（**有條件可送審**；R1～R2） |
| `evidence/` | — | — | 實跑輸出與查證紀錄：2026-10-22 全量測試、TestENG011（S6）、PR #26（cover／env／gh-ci／gh-pr／mutation／race／vet／ct-run）、G4 整體 f646（cover／ct／env／race／test／vet），共 16 檔 |

### G4 文件（13～20 號）

檔名維持 `-v0.1` 後綴，版本以檔內「版本」欄為準；sha256 為 2026-10-27 17:10（UTC+8）本機重算值。

| 檔名 | 標題 | 版本 | sha256 |
|---|---|---|---|
| [13-G4-QA-R1-acceptance-PR26-v0.1.md](./13-G4-QA-R1-acceptance-PR26-v0.1.md) | G4-QA-R1 驗收：PR #26 契約測試 CT-* | v0.1.1 | `8f692c33ba5a88ee191b1320f84c5dd9fd1a27981a7c04b2595f1bf3376b4c43` |
| [14-CR-003-qa-confirmation-v0.1.md](./14-CR-003-qa-confirmation-v0.1.md) | CR-003 品保部確認書（CR3-R11 原件） | v0.1 | `3fb5d45727a337d9d9901ed7a5399922457d166b135cba10845b408f68b30ae9` |
| [15-CR-004-R3-qa-impact-row-v0.1.md](./15-CR-004-R3-qa-impact-row-v0.1.md) | CR-004 §6 影響分析｜品保列（回應 CR4-R3） | v0.1 | `4c3684fbfb92ed4463a24cdd4a6d0e14349ce4418bfb7315b440a49f1617ade2` |
| [16-CR-004-R4-qa-opinion-v0.1.md](./16-CR-004-R4-qa-opinion-v0.1.md) | CR-004 品保書面意見（回應 CR4-R4） | v0.1 | `ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c` |
| [17-CR-004-signoff-2-2-qa-confirmation-v0.1.md](./17-CR-004-signoff-2-2-qa-confirmation-v0.1.md) | CR-004 簽核彙整檔 §2-2 品保確認（CR4-RR02） | v0.1 | `e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587` |
| [18-CR-004-AC6-9-qa-verdict-v0.1.md](./18-CR-004-AC6-9-qa-verdict-v0.1.md) | CR-004 AC-CR004-6～9 與 #24 首測 品保判定 | v0.4 | `543763dad54eb6dd1bbd3b68afcdfa4c80e0b375cfea912fb24a02ef02d61288` |
| [19-G4-overall-qa-conclusion-v0.1.md](./19-G4-overall-qa-conclusion-v0.1.md) | G4 整體品保結論（PR #13～#26）＋G4-QA-R2 處置 | v0.1 | `86606bd3cfb059c52dc1ac9984161d9c8cc879eeeebec9c46795f4d4960b0d70` |
| [20-G4-QA-R2-closure-note-v0.1.md](./20-G4-QA-R2-closure-note-v0.1.md) | G4-QA-R2 處置紀錄（方案 B：經 PR #27、PR #28 進 main） | **v0.2** | `8918b3c958cbd6537b156a69acb813c59554e35c9e08574cec679841708b826f` |

註：20 號 v0.2 與本索引 v0.5 為本次修改，尚未進 main（main `36f4717c` 上仍為 20 號 v0.1 `dcc63af8…`、本索引 v0.4 `25c1807b…`），需下一個純文件 PR。本索引無法記錄自身的 sha256。

## 編號覆蓋

- REQ-001～REQ-012：可測審查＋追溯＋**TC 實號**（v0.3）
- NFR-001～NFR-008：追溯＋CHK／TC
- SEC-001～SEC-017：追溯＋關鍵 TC／CHK
- 公開方法：單測計畫覆蓋 **100%**（見案例目錄 §1、矩陣 §4）

## 交接單：品保 G4-QA-R2 結案 → 總協調／審查（v0.5）

- 日期：2026-10-27（UTC+8）；基準 main `36f4717c`（PR #28 merge commit，2026-10-27 12:01:01 UTC+8 合併），並以最新 main `db380a40`（PR #29，docs 無變更）複核
- 查證：PR #28 狀態、merge sha、合併時間、合併方式皆以 gh／GitHub API 唯讀查證（見 20 號 §2）
- 比對：20 號 v0.1 原列 16 缺漏檔＋13 號＋20 號共 18 檔 sha256 全一致；07-testing 全目錄 37／37 一致；main 原列 6 處簡體字已為 0
- 結論：**G4-QA-R2 結案**；DEF-004 main 全數修正（詳見 20 號 v0.2）
- 待辦：20 號 v0.2 與本索引 v0.5 需以純文件 PR 進 main；開 PR／推送須使用者本人同意

## 交接單（歷史）：品保 G4 前置 → 總協調／審查（v0.4）

- 日期：2026-10-22（UTC+8）；基準 main `db51cd5`
- 實跑：`go test ./...` 167 PASS／0 FAIL／0 SKIP；TestENG011 19 PASS；覆蓋 80.9%；`-race` 全綠
- 結論：**有條件可送審**；必改 G4-QA-R1（CT-* 未實作）、G4-QA-R2（07-testing 文件更新進 main）；建議 S1～S7
- S6：已處置，見 11 號文件
- CR-003：ENG-022 驗收對應已加入任務對照 v0.3（CHK-ENG-022），待研發部核准 CR-003

## 交接單（歷史）：品保 G3 前置 → 總協調／審查

- 交出角色 → 接收角色：品保部 → 總協調（並請審查部引用前置結論）
- 日期：2026-10-07（Asia/Taipei）
- 版本：章節索引 v0.3；狀態：審查中
- 對應需求編號：REQ-001～012、NFR-001～008、SEC-001～017
- 已知問題：無 G3-QA 阻擋必改；建議項 G3-QA-S1～S5（環境基線、事件鍵、研發正式 ID、CI 待辦、G2 建議殘句）
- 需要下游注意：審查依 `10-gate-G3-qa-precheck`；實作綠燈與 CI required checks 屬後續舉證；研發任務對照表待正式清單回填
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
| v0.4 | 2026-10-22 | 品保部 | G4：新增 11（S6）、12（G4 前置）、evidence；矩陣 v0.4、案例 v0.2、任務對照 v0.3（ENG-022）；簡體字修正 |
| v0.5 | 2026-10-27 | 品保部 | 矩陣列改標 v0.5；新增 13～20 號索引列（檔名、標題、版本、sha256）；evidence 列補齊 16 檔；新增 G4-QA-R2 結案交接單（PR #28 merge 36f4717c 查證，依 20 號 v0.2） |
