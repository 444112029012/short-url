---
文件：需求—設計追溯矩陣
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012, NFR-001～NFR-008, SEC-001～SEC-017
專案代號：SHORTURL
---

# 需求—設計追溯矩陣

> 每條 REQ／NFR／SEC 至少對應一個設計元素（模組／類別／方法／API／資料實體／ADR／DFD）。

## 1. 功能需求（REQ）

| 編號 | 設計元素（至少一） |
|---|---|
| REQ-001 | ShortUrlApplicationService.createShortUrl；POST /api/v1/urls；SEQ-01；P1；UrlRecord |
| REQ-002 | RedirectService.redirect；GET /{shortCode}；SEQ-03；P2；SEC-014 Location |
| REQ-003 | ClickCounter.incrementAtomic；RedirectService.redirect；SEQ-03；DF-14 |
| REQ-004 | ShortUrlApplicationService.getStats；ClickCounter.getCount；GET …/stats；SEQ-05；P3 |
| REQ-005 | UrlValidator.validateLongUrl；SEQ-02；P4；SEC-001 |
| REQ-006 | UrlValidator.validateShortCode；ShortCodeGenerator；ADR-003；OpenAPI pattern；SEQ-06 |
| REQ-007 | ErrorMapper.toHttpResponse；SEQ-02；P6；invalid_url |
| REQ-008 | RedirectService／getStats → not_found；SEQ-04；P6 |
| REQ-009 | ErrorMapper；OpenAPI 錯誤格式；SEC-003；P6 |
| REQ-010 | ShortCodeGenerator＋UrlRepository.exists／save 唯一約束；SEQ-01 |
| REQ-011 | UrlValidator；L=2048；OpenAPI maxLength；UrlRecord.long_url |
| REQ-012 | UrlRepository.save 不去重 long_url；資料模型註記 |

## 2. 非功能需求（NFR）

| 編號 | 設計元素 |
|---|---|
| NFR-001 | Container 單路徑；方法備註；演示環境 P95≤2s（實作／測試驗證） |
| NFR-002 | C4 Context 演示操作者；核心 API 三路徑 |
| NFR-003 | C4 Container 單實例假設；ADR-002；RateLimitGuard 進程內 |
| NFR-004 | SEC 全表追溯；安全標頭見 API notes／SEC-008；CI 屬維運／安全（設計標介面） |
| NFR-005 | 資料模型無 IP／UA；隱私生命週期；SEC-004；DFD DS1 |
| NFR-006 | ObservabilityHooks；P7；DF-20／DF-21；DS2 |
| NFR-007 | 流程文件（gates／）；設計部交出 G2 文件供關卡 |
| NFR-008 | 本手冊 02～05 章齊或 N/A；05-ui-ux 不適用聲明 |

## 3. 資安需求（SEC）

| 編號 | 設計元素 |
|---|---|
| SEC-001 | UrlValidator.validateLongUrl；SEQ-02；P4；DF-03／04 |
| SEC-002 | 協定白名單（SEC-001）；殘餘風險書面（API notes／本表）；無預覽頁 OUT-10 |
| SEC-003 | ErrorMapper；REQ-009；P6 |
| SEC-004 | UrlMapping 欄位集合；DS1 禁止 visitor_ip／user_agent |
| SEC-005 | ObservabilityHooks 約束；DF-21；TB-04 |
| SEC-006 | ShortCodeGenerator；ADR-003 熵估算 |
| SEC-007 | RateLimitGuard；SEQ-07；P5；API notes |
| SEC-008 | 01-api-notes／OpenAPI 說明安全標頭期望 |
| SEC-009 | CI／維運（設計標「實作於 09-cicd」）；追溯佔位 |
| SEC-010 | 同 SEC-009（SCA） |
| SEC-011 | 同 SEC-009（secrets） |
| SEC-012 | ReverseProxy 可選；API notes 環境區分 |
| SEC-013 | UrlValidator.validateShortCode；SEQ-06；OpenAPI pattern |
| SEC-014 | RedirectService 僅用庫內 long_url；SEQ-03；DF-13 |
| SEC-015 | 部署約束（維運）；設計註記：產物不含 .git 暴露 |
| SEC-016 | UrlRepository／InMemoryOrDbUrlStore／ClickCounter.incrementAtomic：參數化查詢或等價綁定、禁止字串拼接 SQL；原子 click_count+1（方法規格 §6～8；資料模型 §6） |
| SEC-017 | UrlStore／組裝／資料模型 §6：DB 最小權限、DS1 非公網、憑證不進版控（環境變數／機密注入）、查詢欄位最小化 |


---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | G2 追溯初稿 |
| v0.1+R1 | 2026-10-07 | 設計部 | 關閉 G2-QA-R1：納入 SEC-016／017 |
