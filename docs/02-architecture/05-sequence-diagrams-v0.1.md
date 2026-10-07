---
文件：循序圖
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012, SEC-007, SEC-013
專案代號：SHORTURL
---

# 循序圖：關鍵流程

> 參與者對應 Component。錯誤碼：`invalid_url`｜`not_found`｜`rate_limited`｜`internal_error`。

---

## SEQ-01 成功建立短碼

**對應**：REQ-001、REQ-005、REQ-010、REQ-011、REQ-012；SEC-001、SEC-006

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant App as ShortUrlApplicationService
  participant Val as UrlValidator
  participant Gen as ShortCodeGenerator
  participant Repo as UrlRepository
  participant Obs as ObservabilityHooks

  EE1->>API: POST /api/v1/urls {url}
  API->>RL: checkCreate(sourceId)
  RL-->>API: Ok
  API->>App: createShortUrl(longUrl)
  App->>Val: validateLongUrl(longUrl)
  Val-->>App: Ok(ValidatedUrl) 長度≤2048 http(s)
  App->>Gen: generate()
  Gen-->>App: short_code (CSPRNG base62×8)
  App->>Repo: exists(short_code)?
  Repo-->>App: false
  App->>Repo: save(UrlRecord)
  Repo-->>App: Ok
  App-->>API: Ok(UrlCreated)
  API->>Obs: onCreateSuccess(short_code)
  API-->>EE1: 201 {short_code, short_url, long_url}
```

---

## SEQ-02 建立驗證失敗

**對應**：REQ-007、REQ-005、REQ-011；SEC-001

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant App as ShortUrlApplicationService
  participant Val as UrlValidator
  participant Map as ErrorMapper
  participant Obs as ObservabilityHooks

  EE1->>API: POST /api/v1/urls {url=非法}
  API->>RL: checkCreate(sourceId)
  RL-->>API: Ok
  API->>App: createShortUrl(longUrl)
  App->>Val: validateLongUrl(longUrl)
  Val-->>App: Err(invalid_url)
  App-->>API: Err(invalid_url)
  API->>Map: toHttpResponse(invalid_url)
  Map-->>API: 400 {error.type=invalid_url}
  API->>Obs: onCreateFailure("invalid_url")
  API-->>EE1: 400 無堆疊／無內部細節
  Note over App: 不呼叫 save；庫無新短碼
```

---

## SEQ-03 成功導向＋原子計數

**對應**：REQ-002、REQ-003；SEC-014

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant Red as RedirectService
  participant Val as UrlValidator
  participant Repo as UrlRepository
  participant Cnt as ClickCounter
  participant Obs as ObservabilityHooks

  EE1->>API: GET /{shortCode}
  API->>RL: checkRedirect(sourceId)
  RL-->>API: Ok
  API->>Red: redirect(shortCode)
  Red->>Val: validateShortCode(shortCode)
  Val-->>Red: Ok
  Red->>Repo: findByShortCode(shortCode)
  Repo-->>Red: Some(UrlRecord)
  Red->>Cnt: incrementAtomic(shortCode)
  Cnt-->>Red: Ok(N+1)
  Red-->>API: Ok(RedirectTarget{long_url})
  API->>Obs: onRedirectSuccess(shortCode)
  API-->>EE1: 302 Location=long_url（庫中已驗證值）
```

---

## SEQ-04 短碼不存在導向

**對應**：REQ-008；SEC-013

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant Red as RedirectService
  participant Val as UrlValidator
  participant Repo as UrlRepository
  participant Map as ErrorMapper
  participant Obs as ObservabilityHooks

  EE1->>API: GET /{shortCode}
  API->>RL: checkRedirect(sourceId)
  RL-->>API: Ok
  API->>Red: redirect(shortCode)
  Red->>Val: validateShortCode(shortCode)
  Val-->>Red: Ok
  Red->>Repo: findByShortCode(shortCode)
  Repo-->>Red: None
  Red-->>API: Err(not_found)
  API->>Map: toHttpResponse(not_found)
  API->>Obs: onRedirectFailure("not_found")
  API-->>EE1: 404 {error.type=not_found}
  Note over Red: 不呼叫 incrementAtomic；次數不增
```

---

## SEQ-05 查詢點擊次數

**對應**：REQ-004

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant App as ShortUrlApplicationService
  participant Val as UrlValidator
  participant Cnt as ClickCounter
  participant Obs as ObservabilityHooks

  EE1->>API: GET /api/v1/urls/{shortCode}/stats
  API->>App: getStats(shortCode)
  App->>Val: validateShortCode(shortCode)
  Val-->>App: Ok
  App->>Cnt: getCount(shortCode)
  Cnt-->>App: Ok(click_count)
  App-->>API: Ok(UrlStats)
  API->>Obs: onStatsSuccess(shortCode)
  API-->>EE1: 200 {short_code, click_count}
```

---

## SEQ-06 短碼格式非法

**對應**：REQ-006；SEC-013

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant Red as RedirectService
  participant Val as UrlValidator
  participant Map as ErrorMapper

  EE1->>API: GET /{非法短碼}
  API->>RL: checkRedirect(sourceId)
  RL-->>API: Ok
  API->>Red: redirect(code)
  Red->>Val: validateShortCode(code)
  Val-->>Red: Err(invalid_url 或等價格式錯誤)
  Note over Val: 建議對外統一 invalid_url 或保留穩定碼；<br/>不當作有效短碼、不導向、不計數
  Red-->>API: Err(invalid_url)
  API->>Map: toHttpResponse
  API-->>EE1: 400 {error.type=invalid_url}
```

> 實作約定：短碼格式失敗對外 `type=invalid_url`（與 OpenAPI 400 對齊）；若未來拆 `invalid_short_code` 須走 CR。查統計路徑同理。

---

## SEQ-07 速率限制觸發

**對應**：SEC-007

```mermaid
sequenceDiagram
  actor EE1 as EE1 客戶端
  participant API as HttpApiAdapter
  participant RL as RateLimitGuard
  participant Map as ErrorMapper
  participant Obs as ObservabilityHooks

  EE1->>API: POST /api/v1/urls（超過 30／分／來源）
  API->>RL: checkCreate(sourceId)
  RL-->>API: Err(rate_limited)
  API->>Map: toHttpResponse(rate_limited)
  API->>Obs: onCreateFailure("rate_limited")
  API-->>EE1: 429 {error.type=rate_limited}
  Note over API: 不進入 Application；不建立短碼
```

導向端點同理：`checkRedirect` 門檻 120／分／來源 → 429，不導向、不計數。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | SEQ-01～07；覆蓋主要 REQ |
