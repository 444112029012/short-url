---
文件：追溯矩陣
版本：v0.3
狀態：審查中
負責角色：品保部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-017
對應規格：docs/01-requirements/02-requirements-spec-v0.1.md；docs/06-security/01-sec-requirements-v0.2.md
對應案例：docs/07-testing/07-test-cases-v0.1.md
專案代號：SHORTURL
對應退回項：G1-R1（已關）；G2-QA-R1（已關）；G3 TC 回填
---

# 追溯矩陣：短網址服務

> G3 目標：需求＋設計／API 欄齊＋**TC 欄填入實際編號**；公開方法單測覆蓋勾稽 100%。  
> SEC 來源優先 **v0.2**（含 SEC-016、SEC-017）。

**來源規格路徑（共通）：** `docs/01-requirements/02-requirements-spec-v0.1.md`  
**資安來源：** `docs/06-security/01-sec-requirements-v0.2.md`  
**設計來源：** `docs/02-architecture/`、`docs/03-data/`、`docs/04-api/`  
**案例來源：** `docs/07-testing/07-test-cases-v0.1.md`  
**策略／DoD：** `06-test-strategy-v0.1.md`、`08-definition-of-done-v0.1.md`

**路徑縮寫：** Arch＝`docs/02-architecture/`；API＝`docs/04-api/`；Data＝`docs/03-data/`

---

## 1. 功能需求（REQ）

| 需求編號 | 需求摘要 | 優先級 | 來源規格路徑 | 設計文件／ADR | API／元件 | TC 編號 | SEC | 備註 |
|---|---|---|---|---|---|---|---|---|
| REQ-001 | 接受合法長網址、建立短碼並回傳可造訪短網址 | 必須 | `docs/01-requirements/02-requirements-spec-v0.1.md` | Arch `06-method-specs`（createShortUrl）；`05-sequence` SEQ-01；Data UrlMapping | **createShortUrl**；HttpApiAdapter.handleCreate；ShortUrlApplicationService | UT-APP-01、UT-HC-01、CT-CREATE-01、E2E-01、IT-FLOW-01 | SEC-001、SEC-006 | |
| REQ-002 | 造訪有效短網址時 HTTP 重新導向至長網址 | 必須 | 同上 | RedirectService.redirect；SEQ-03；ADR-001 | **redirectShortCode**；handleRedirect | UT-RED-01、UT-HR-01、CT-REDIR-01、E2E-01 | SEC-014 | |
| REQ-003 | 成功導向時點擊次數原子遞增 1 | 必須 | 同上 | ClickCounter.incrementAtomic；SEQ-03 | redirect 成功路徑；InMemoryOrDbUrlStore | UT-RED-01、UT-CLK-01、IT-STORE-02 | SEC-016 | 失敗不計次見 UT-RED-02 |
| REQ-004 | 查詢指定短碼目前點擊次數 | 必須 | 同上 | getStats；SEQ-05；getCount | **getUrlStats**；handleStats | UT-APP-04、UT-HS-01、CT-STATS-01、E2E-01 | — | |
| REQ-005 | 長網址驗證（非空、可解析、http／https、長度） | 必須 | 同上 | UrlValidator.validateLongUrl；SEQ-02 | createShortUrl 400 | UT-VAL-01～05、E2E-02、CT-CREATE-02 | SEC-001 | |
| REQ-006 | 短碼符合字元集與長度 | 必須 | 同上 | **ADR-003**；validateShortCode；SEQ-06 | ShortCodePath pattern | UT-VAL-06～07、UT-HR-03、E2E-06 | SEC-006、013 | |
| REQ-007 | 建立驗證失敗明確錯誤且不落庫 | 必須 | 同上 | ErrorMapper；SEQ-02 | createShortUrl 400 invalid_url | UT-APP-02、UT-HC-02、IT-FLOW-02、CT-CREATE-02 | SEC-001、003 | |
| REQ-008 | 短碼不存在明確錯誤；不導向、不計次 | 必須 | 同上 | redirect／getStats → not_found；SEQ-04 | 404 not_found | UT-RED-02、UT-APP-05、E2E-03、CT-REDIR-02 | SEC-013 | |
| REQ-009 | 對外錯誤不洩漏內部細節 | 必須 | 同上 | ErrorMapper.toHttpResponse | ErrorResponse | UT-ERR-01～02、CT-ERR-01 | SEC-003 | |
| REQ-010 | 已建立短碼全域唯一（含並發） | 必須 | 同上 | ShortCodeGenerator＋Repo；ADR-003 | createShortUrl；UNIQUE | UT-APP-03、UT-GEN-02、UT-REPO-03、IT-STORE-01 | SEC-006 | |
| REQ-011 | 長網址最大長度 L 文件化且行為一致 | 必須 | 同上 | L=**2048**；OpenAPI maxLength | CreateUrlRequest.url | UT-VAL-04～05、E2E-02 | SEC-001 | |
| REQ-012 | 同一長網址可多次建立為不同短碼 | 可以 | 同上 | save 不去重 | createShortUrl | IT-FLOW-03、UT-REPO-01 | — | |

---

## 2. 非功能需求（NFR）

| 需求編號 | 需求摘要 | 優先級 | 來源規格路徑 | 設計文件／ADR | API／元件 | TC 編號 | SEC | 備註 |
|---|---|---|---|---|---|---|---|---|
| NFR-001 | 演示／小流量下建立與導向 P95 ≤ 2 秒 | 應該 | 同上 | C4 Container；策略 §4 基線建議 | create／redirect | **CHK-NFR-001** | — | 環境基線建議補（G3-QA-S1） |
| NFR-002 | 演示期間核心路徑至少完整成功 1 次 | 應該 | 同上 | 三 API 路徑 | 三 operationId | **CHK-NFR-002**、E2E-01 | — | |
| NFR-003 | 不承諾水平擴展；單實例假設文件化 | 可以 | 同上 | ADR-002；RateLimitGuard 進程內 | UrlStore；RateLimitGuard | **CHK-NFR-003** | — | 文件審查 |
| NFR-004 | ASVS、SEC、CI 門檻、G5 高風險 | 必須 | 同上；`03-ci-security-gates` | SEC 全表；CI 文件 | CI jobs | **CHK-NFR-004**、CHK-SEC-009～011 | SEC-001～017 | |
| NFR-005 | 隱私合規蒐集／保存邊界 | 必須 | 同上；隱私文件 | Data 禁 IP／UA；DFD | UrlMapping | **CHK-NFR-005**、IT-STORE-01、UT-REPO-01 | SEC-004、005 | |
| NFR-006 | 建立／導向成功失敗可觀察 | 應該 | 同上 | ObservabilityHooks | Hooks.* | **UT-OBS-01**、**CHK-NFR-006** | SEC-005 | 事件鍵建議 G3-QA-S2 |
| NFR-007 | G0–G6 關卡通過紀錄 | 必須 | 同上 | `docs/gates/` | — | **CHK-NFR-007** | — | |
| NFR-008 | G3 前各章簽核或不適用 | 必須 | 同上 | UI N/A；02～04 章 | — | **CHK-NFR-008** | — | |

---

## 3. 安全需求（SEC）

| 需求編號 | 需求摘要 | 優先級 | 來源規格路徑 | 設計文件／ADR | API／元件 | TC 編號 | SEC | 備註 |
|---|---|---|---|---|---|---|---|---|
| SEC-001 | 建立時伺服端驗證長網址；失敗不建碼 | 必須 | `01-sec-requirements-v0.2.md` | validateLongUrl；SEQ-02 | createShortUrl；UrlValidator | UT-VAL-02～05、UT-APP-02、UT-HC-02、E2E-02、CT-CREATE-02 | SEC-001 | 關鍵 |
| SEC-002 | 開放重新導向殘餘：白名單＋書面接受 | 必須 | 同上 | api-notes §7；RR-001／D-06 | 僅 http(s) | **CHK-SEC-002**、UT-VAL-03 | SEC-002 | 已接受殘餘 |
| SEC-003 | 對外錯誤不洩漏堆疊／路徑／SQL | 必須 | 同上 | ErrorMapper | ErrorResponse | UT-ERR-01～02、CT-ERR-01 | SEC-003 | 關鍵 |
| SEC-004 | 不得持久化訪客 IP／UA | 必須 | 同上 | Data 禁止欄；UrlRecord | Store schema | UT-REPO-01、IT-STORE-01、CHK-NFR-005 | SEC-004 | 關鍵 |
| SEC-005 | 日誌不得列印 Authorization／Cookie 原文 | 必須 | 同上 | ObservabilityHooks 約束 | Hooks | UT-OBS-01、**CHK-SEC-005** | SEC-005 | |
| SEC-006 | 短碼 CSPRNG；禁止可預測序列 | 必須 | 同上 | **ADR-003**；generate | createShortUrl | UT-GEN-01～02 | SEC-006 | 關鍵 |
| SEC-007 | 建立／導向基本速率限制 | 應該 | 同上 | RateLimitGuard；notes §4 | checkCreate／Redirect；429 | UT-RL-01～02、UT-HC-03、E2E-05 | SEC-007 | 關鍵 |
| SEC-008 | ASVS L1 合理安全標頭子集 | 應該 | 同上 | api-notes §5 | Http 標頭 | **CHK-SEC-008** | SEC-008 | |
| SEC-009 | CI SAST；高／嚴重擋合併 | 必須 | 同上；CI 門檻 | 09-cicd／CI 文件 | CI | **CHK-SEC-009** | SEC-009 | |
| SEC-010 | CI SCA | 必須 | 同上 | 同上 | CI | **CHK-SEC-010** | SEC-010 | |
| SEC-011 | CI 機密掃描 | 必須 | 同上 | 同上 | CI | **CHK-SEC-011** | SEC-011 | |
| SEC-012 | 演示／對外 HTTPS；localhost 標非演示 | 應該 | 同上 | OpenAPI servers；notes §6 | 部署 | **CHK-SEC-012** | SEC-012 | |
| SEC-013 | 路徑短碼須符合字元集與長度 | 必須 | 同上 | validateShortCode；SEQ-06 | ShortCodePath；400 | UT-VAL-07、UT-HR-03、E2E-06 | SEC-013 | 關鍵 |
| SEC-014 | Location＝庫內 long_url；禁止參數改寫 | 必須 | 同上 | RedirectService；SEQ-03 | redirectShortCode | UT-RED-01、UT-HR-01、CT-REDIR-01、E2E-04 | SEC-014 | 關鍵 |
| SEC-015 | 不得暴露 `.git` 等 | 必須 | 同上 | 部署註記 | 對外路徑 | **CHK-SEC-015** | SEC-015 | |
| SEC-016 | DS1 參數化；點擊原子遞增 | 必須 | 同上 | 方法規格 §6～8；資料模型 §6 | Store；incrementAtomic | UT-CLK-01、IT-STORE-02～03 | SEC-016 | 關鍵；G2-R1 已關 |
| SEC-017 | 最小權限；非公網；憑證不進版控；欄位最小化 | 必須 | 同上 | 資料模型 §6；Container §5 | 部署／Store | **CHK-SEC-017**、UT-REPO-02、IT-STORE-01 | SEC-017 | G3 举证 |

---

## 4. 公開方法 → 單測覆蓋勾稽

| 類別 | 方法 | UT 計畫 ID | 覆蓋 |
|---|---|---|---|
| HttpApiAdapter | handleCreate | UT-HC-01、02、03 | ✅ |
| HttpApiAdapter | handleRedirect | UT-HR-01、02、03 | ✅ |
| HttpApiAdapter | handleStats | UT-HS-01、02 | ✅ |
| ShortUrlApplicationService | createShortUrl | UT-APP-01、02、03 | ✅ |
| ShortUrlApplicationService | getStats | UT-APP-04、05 | ✅ |
| RedirectService | redirect | UT-RED-01、02 | ✅ |
| UrlValidator | validateLongUrl | UT-VAL-01～05 | ✅ |
| UrlValidator | validateShortCode | UT-VAL-06、07 | ✅ |
| ShortCodeGenerator | generate | UT-GEN-01、02 | ✅ |
| UrlRepository | save | UT-REPO-01 | ✅ |
| UrlRepository | findByShortCode | UT-REPO-02 | ✅ |
| UrlRepository | exists | UT-REPO-03 | ✅ |
| ClickCounter | incrementAtomic | UT-CLK-01 | ✅ |
| ClickCounter | getCount | UT-CLK-02 | ✅ |
| InMemoryOrDbUrlStore | 實作 repo／counter | IT-STORE-01～03 | ✅（整合） |
| ErrorMapper | toHttpResponse | UT-ERR-01、02 | ✅ |
| RateLimitGuard | checkCreate | UT-RL-01 | ✅ |
| RateLimitGuard | checkRedirect | UT-RL-02 | ✅ |
| ObservabilityHooks | on* Success／Failure | UT-OBS-01 | ✅（合併） |

**公開方法單測計畫覆蓋：100%。**

---

## 5. 覆蓋勾稽總表

| 集合 | 需求欄 | 設計／API | TC 欄 | 狀態 |
|---|---|---|---|---|
| REQ-001～012 | 12 | 12 | 12 已填 | ✅ |
| NFR-001～008 | 8 | 已填 | 8 已填（CHK／TC） | ✅ |
| SEC-001～017 | 17 | 17 | 17 已填 | ✅ |
| 公開方法 UT | — | — | 無遺漏 | ✅ 100% |
| 三 API 契約＋E2E 主路徑 | — | create／redirect／stats | CT-*＋E2E-01 | ✅ |

---

## 6. 穩定錯誤碼（供 TC 斷言）

| type | HTTP | 主要觸發 | 代表 TC |
|---|---|---|---|
| invalid_url | 400 | 長網址驗證失敗；短碼格式非法 | UT-VAL-*、CT-CREATE-02、E2E-06 |
| not_found | 404 | 短碼未登記 | UT-RED-02、E2E-03 |
| rate_limited | 429 | 超過來源速率門檻 | UT-RL-*、E2E-05 |
| internal_error | 500 | 未預期錯誤（無內部細節） | UT-ERR-01 |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G1：需求欄齊；設計／TC TBD |
| v0.1.1 | 2026-10-07 | 品保部 | G1-R1：SEC-001～015 |
| v0.2 | 2026-10-07 | 品保部 | G2：設計／API 欄；SEC-016／017 |
| v0.2.1 | 2026-10-07 | 品保部 | 關閉 G2-QA-R1；016／017 設計欄回填 |
| v0.3 | 2026-10-07 | 品保部 | **G3**：TC 欄填入實際 ID；公開方法覆蓋勾稽 100%；對齊案例目錄 v0.1 |
