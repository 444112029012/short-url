---
文件：01-requirements 章節索引
版本：v0.2
狀態：審查中
負責角色：企劃部
最後更新：2026-10-07
專案代號：SHORTURL
對應退回項：G1-R2
---

# 01 需求規格｜章節索引

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-requirements-comprehensiveness-v0.1.md](./01-requirements-comprehensiveness-v0.1.md) | v0.1 | 審查中 | 角色、情境、邊界、情境→REQ |
| [02-requirements-spec-v0.1.md](./02-requirements-spec-v0.1.md) | v0.1 | 審查中 | REQ-001～012、NFR-001～008；§4 對照安全文件 |
| [03-scope-v0.1.md](./03-scope-v0.1.md) | v0.1 | 審查中 | 範圍內外（對齊章程並細化） |
| [04-privacy-compliance-v0.1.md](./04-privacy-compliance-v0.1.md) | v0.1 | 審查中 | 隱私合規（C-05） |

## 資安需求正式文件（安全部）

| 文件 | 說明 |
|---|---|
| [../06-security/00-asvs-level-v0.1.md](../06-security/00-asvs-level-v0.1.md) | ASVS 驗證等級選定（OWASP ASVS 5.0 Level 1） |
| [../06-security/01-sec-requirements-v0.1.md](../06-security/01-sec-requirements-v0.1.md) | 資安需求清單 SEC-001～SEC-015 |

## 編號清單（供品保追溯矩陣需求欄）

- REQ-001、REQ-002、REQ-003、REQ-004、REQ-005、REQ-006、REQ-007、REQ-008、REQ-009、REQ-010、REQ-011、REQ-012
- NFR-001、NFR-002、NFR-003、NFR-004、NFR-005、NFR-006、NFR-007、NFR-008
- SEC-001～SEC-015（詳見 [`docs/06-security/01-sec-requirements-v0.1.md`](../06-security/01-sec-requirements-v0.1.md)；ASVS 見 [`docs/06-security/00-asvs-level-v0.1.md`](../06-security/00-asvs-level-v0.1.md)）

## 交接單：G1 需求 → 總協調／安全／品保／審查

- 交出角色 → 接收角色：企劃部 → 總協調（串安全前置結論、品保可測與追溯、送審查）
- 日期：2026-10-07
- 交付物清單：上表四份需求文件＋本索引；資安正式文件見 `docs/06-security/`（ASVS、SEC）
- 版本：索引 v0.2／需求正文 v0.1（狀態：審查中）
- 對應需求編號：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-015
- 已知問題：長網址上限具體數字 L 待設計寫入 OpenAPI（REQ-011）；G1-R2 已關閉（改引用 06-security 正式文件）
- 需要下游注意：品保追溯矩陣需求欄須含 REQ／NFR／SEC 上列編號；審查依 G1 四層檢查；安全前置結論見 `docs/06-security/`
- 待決問題：D-08 保存期若使用者本輪不拍板，採隱私文件預設值並留追溯

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 企劃部 | G1 初稿索引與交接單 |
| v0.2 | 2026-10-07 | 企劃部 | 關閉 G1-R2：改引用 06-security ASVS／SEC 正式文件 |
