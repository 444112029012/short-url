---
文件：ADR-002 持久化
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-001～REQ-004, REQ-010, REQ-012, NFR-003, SEC-004
專案代號：SHORTURL
---

# ADR-002：持久化（UrlStore／DS1）

| 欄位 | 內容 |
|---|---|
| 狀態 | **已採納**（D-04，2026-10-07，方案 A：SQLite） |
| 日期 | 2026-10-07 |

## 背景

需持久化 UrlMapping（short_code、long_url、click_count、created_at）。MVP 小流量、單實例（NFR-003）、無訪客 IP／UA 明細（SEC-004）。

## 選項

### 方案 A：SQLite 單檔

- 優：零額外部務、備份＝拷貝檔、契合單實例練手
- 缺：多寫併發有限；未來水平擴展需遷移

### 方案 B：PostgreSQL

- 優：併發與 SQL 能力強、易接維運監控
- 缺：需額外部署；對 MVP 偏重

## 比較

| 維度 | A SQLite | B PostgreSQL |
|---|---|---|
| 單實例契合 | 高 | 中（偏重） |
| 原子 click_count+1 | 支援 | 支援 |
| 運維成本 | 低 | 較高 |
| 未來擴展 | 需 CR 遷移 | 較順 |

## 建議傾向

**已採納方案 A（SQLite 單檔）**（D-04，2026-10-07）。InMemory 僅可作測試雙件；演示／正式路徑須可持久化至 SQLite。

## 後果

- Schema 必須：short_code UNIQUE；無 visitor_ip／user_agent。
- 多實例若未來需要，改集中式 DB＋限流，走 CR。

## 對應需求

REQ-001～004、010、012；NFR-003；SEC-004；DFD DS1。
| v0.1+D04 | 2026-10-07 | 設計部核對 | D-04 採納方案 A；清除待拍板措辭 |
