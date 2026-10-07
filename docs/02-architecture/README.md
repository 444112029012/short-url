---
文件：系統架構章節索引與交接單
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012, NFR-001～NFR-008, SEC-001～SEC-015
專案代號：SHORTURL
---

# 02-architecture：系統架構與詳細設計

## 章節索引

| 檔案 | 說明 | 版本 | 狀態 |
|---|---|---|---|
| [01-c4-context-v0.1.md](./01-c4-context-v0.1.md) | C4 Context | v0.1 | 審查中 |
| [02-c4-container-v0.1.md](./02-c4-container-v0.1.md) | C4 Container | v0.1 | 審查中 |
| [03-c4-component-v0.1.md](./03-c4-component-v0.1.md) | C4 Component（必交） | v0.1 | 審查中 |
| [04-class-diagram-v0.1.md](./04-class-diagram-v0.1.md) | 類別圖 | v0.1 | 審查中 |
| [05-sequence-diagrams-v0.1.md](./05-sequence-diagrams-v0.1.md) | 循序圖 | v0.1 | 審查中 |
| [06-method-specs-v0.1.md](./06-method-specs-v0.1.md) | 方法規格表 | v0.1 | 審查中 |
| [07-req-design-traceability-v0.1.md](./07-req-design-traceability-v0.1.md) | 需求—設計追溯 | v0.1 | 審查中 |
| [adr/](./adr/) | 架構決策紀錄 | v0.1 | 審查中（決策：已採納） |

---

## 交接單

| 項目 | 內容 |
|---|---|
| **交出** | 設計部 → 總協調／安全部／品保部／審查部 |
| **交付物** | C4 三層、類別圖、循序圖（≥7）、方法規格表、追溯表、ADR-001～003（D-04 已採納：Go+chi／SQLite／base62×8 CSPRNG）、資料模型、DFD（見 `03-data/`）、OpenAPI（見 `04-api/`）、UI N/A（見 `05-ui-ux/`） |
| **版本** | 全部 v0.1 |
| **狀態** | 審查中 |
| **對應需求** | REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-015 |
| **已知問題** | G2-QA-R1（SEC-016／017）設計側已補；servers placeholder；無認證（OUT-03） |
| **下游注意** | 安全部請依 `03-data/02-dfd-v0.1.md` 信任邊界做 STRIDE；研發須等 D-04 與 G2 簽核後依方法規格表實作；品保可依 OpenAPI／方法規格建 TC |
| **待決** | 無（D-04 已於 2026-10-07 拍板）；後續待 G2 審查通過 |

### 依賴方向摘要（OOP）

```
HttpApiAdapter → ShortUrlApplicationService / RedirectService
               → RateLimitGuard / ErrorMapper / Observability
Application → Domain Ports（UrlRepository、ClickCounter）← Infrastructure（InMemoryOrDbUrlStore）
```

禁止循環依賴；領域不依賴框架細節。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | G2 初稿，狀態審查中 |
