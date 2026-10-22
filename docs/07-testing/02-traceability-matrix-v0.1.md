---
文件：追溯矩陣
版本：v0.5
狀態：審查中
負責角色：品保部
最後更新：2026-10-22
對應需求：REQ-001～REQ-012、NFR-001～NFR-008、SEC-001～SEC-017
對應規格：docs/01-requirements/02-requirements-spec-v0.1.md；docs/06-security/01-sec-requirements-v0.2.md
對應案例：docs/07-testing/07-test-cases-v0.1.md（內容版 v0.2）
專案代號：SHORTURL
對應退回項：G1-R1（已關）；G2-QA-R1（已關）；G3 TC 回填；G4 審查部 S6（PR #18，v0.4 處置）；G4 整體退回 G4-R2（v0.5，PR #13～#26、main f646bee）
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
| REQ-003 | 成功導向時點擊次數原子遞增 1 | 必須 | 同上 | ClickCounter.incrementAtomic；SEQ-03 | redirect 成功路徑；InMemoryOrDbUrlStore | UT-RED-01、UT-CLK-01、IT-STORE-02、E2E-05-03 | SEC-016 | 失敗不計次見 UT-RED-02；限流被拒不計次見 E2E-05-03（v0.4） |
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
| NFR-003 | 不承諾水平擴展；單實例假設文件化 | 可以 | 同上 | ADR-002；RateLimitGuard 進程內 | UrlStore；RateLimitGuard | **CHK-NFR-003**、E2E-05-10（輔助） | — | 文件審查；E2E-05-10 僅證單實例進程內限流並發一致，不宣稱跨實例（S6，v0.4） |
| NFR-004 | ASVS、SEC、CI 門檻、G5 高風險 | 必須 | 同上；`03-ci-security-gates` | SEC 全表；CI 文件 | CI jobs | **CHK-NFR-004**、CHK-SEC-009～011 | SEC-001～017 | |
| NFR-005 | 隱私合規蒐集／保存邊界 | 必須 | 同上；隱私文件 | Data 禁 IP／UA；DFD | UrlMapping | **CHK-NFR-005**、IT-STORE-01、UT-REPO-01 | SEC-004、005 | |
| NFR-006 | 建立／導向成功失敗可觀察 | 應該 | 同上 | ObservabilityHooks | Hooks.* | **UT-OBS-01**、**CHK-NFR-006**、E2E-05-11 | SEC-005 | 事件鍵建議 G3-QA-S2；rate_limited 失敗事件見 E2E-05-11 |
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
| SEC-005 | 日誌不得列印 Authorization／Cookie 原文 | 必須 | 同上 | ObservabilityHooks 約束 | Hooks | UT-OBS-01、**CHK-SEC-005**、E2E-05-11 | SEC-005 | E2E-05-11：日誌無 peer／轉送位址與標頭名（v0.4） |
| SEC-006 | 短碼 CSPRNG；禁止可預測序列 | 必須 | 同上 | **ADR-003**；generate | createShortUrl | UT-GEN-01～02 | SEC-006 | 關鍵 |
| SEC-007 | 建立／導向基本速率限制 | 應該 | 同上 | RateLimitGuard；notes §4 | checkCreate／Redirect；429 | UT-RL-01～02、UT-HC-03、E2E-05（E2E-05-01～11） | SEC-007 | 關鍵；S6 處置見 `11-s6-eng011-disposition-v0.1.md` |
| SEC-008 | ASVS L1 合理安全標頭子集 | 應該 | 同上 | api-notes §5 | Http 標頭 | **CHK-SEC-008** | SEC-008 | |
| SEC-009 | CI SAST；高／嚴重擋合併 | 必須 | 同上；CI 門檻 | 09-cicd／CI 文件 | CI | **CHK-SEC-009** | SEC-009 | |
| SEC-010 | CI SCA | 必須 | 同上 | 同上 | CI | **CHK-SEC-010** | SEC-010 | |
| SEC-011 | CI 機密掃描 | 必須 | 同上 | 同上 | CI | **CHK-SEC-011** | SEC-011 | |
| SEC-012 | 演示／對外 HTTPS；localhost 標非演示 | 應該 | 同上 | OpenAPI servers；notes §6 | 部署 | **CHK-SEC-012** | SEC-012 | |
| SEC-013 | 路徑短碼須符合字元集與長度 | 必須 | 同上 | validateShortCode；SEQ-06 | ShortCodePath；400 | UT-VAL-07、UT-HR-03、E2E-06 | SEC-013 | 關鍵 |
| SEC-014 | Location＝庫內 long_url；禁止參數改寫 | 必須 | 同上 | RedirectService；SEQ-03 | redirectShortCode | UT-RED-01、UT-HR-01、CT-REDIR-01、E2E-04 | SEC-014 | 關鍵 |
| SEC-015 | 不得暴露 `.git` 等 | 必須 | 同上 | 部署註記 | 對外路徑 | **CHK-SEC-015** | SEC-015 | |
| SEC-016 | DS1 參數化；點擊原子遞增 | 必須 | 同上 | 方法規格 §6～8；資料模型 §6 | Store；incrementAtomic | UT-CLK-01、IT-STORE-02～03 | SEC-016 | 關鍵；G2-R1 已關 |
| SEC-017 | 最小權限；非公網；憑證不進版控；欄位最小化 | 必須 | 同上 | 資料模型 §6；Container §5 | 部署／Store | **CHK-SEC-017**、UT-REPO-02、IT-STORE-01 | SEC-017 | G3 舉證 |

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
| rate_limited | 429 | 超過來源速率門檻 | UT-RL-*、E2E-05-01～11 |
| internal_error | 500 | 未預期錯誤（無內部細節） | UT-ERR-01 |

---

## 7. 實作對應與執行結果（v0.4 新增；v0.5 更新基準）

> **v0.5 基準**：main `f646beeb25eec7a20b358df0af88450c76a85af5`（#26 合併後）；box go1.27.3；2026-10-22 15:10～15:16（UTC+8）實跑 `go test ./... -count=1`：**181 PASS／0 FAIL／0 SKIP**（頂層 106＋子測試 75）；`-race` 全綠；CT 14 PASS；覆蓋 80.9%；main CI run 39737480926 success。原始輸出：`evidence/g4-overall-f646-*`；結論：`19-G4-overall-qa-conclusion-v0.1.md`。下方 v0.4 基準說明保留供追溯；除 CT 列外，其餘各列在 f646bee 重跑結果相同（✅ 通過）。


> 基準：main `db51cd5`；box go1.27.3；2026-10-22（UTC+8）實跑 `go test ./... -count=1`：**167 PASS／0 FAIL／0 SKIP**（頂層 100＋子測試 67）；`-race` 全綠；main CI run（`db51cd5`）success。原始輸出：`evidence/g4-full-run-2026-10-22.txt`。  
> 「標記」＝測試名稱或註解含 TC ID；「未標記」＝行為已測到但名稱沒有 TC ID（追溯靠本表）。

| TC | 實作（套件／測試） | 標記 | 結果 |
|---|---|---|---|
| UT-VAL-01～07 | domain／validator_test | 標記 | ✅ 通過 |
| UT-GEN-01～02 | domain／generator_test | 標記 | ✅ 通過 |
| UT-APP-01～05 | application／service_test | 標記 | ✅ 通過 |
| UT-RED-01～02 | application／redirect_test | 標記 | ✅ 通過 |
| UT-REPO-01～03、UT-CLK-01～02 | infra/sqlite／store_test | 標記 | ✅ 通過 |
| UT-ERR-01～02 | httpapi／mapper_test | 標記 | ✅ 通過 |
| UT-RL-01～02 | ratelimit／guard_test | 標記 | ✅ 通過 |
| UT-HC-01、UT-HR-01、UT-HS-01 | httpapi／router_test | 標記 | ✅ 通過 |
| UT-HC-02 | httpapi／`TestCreateValidationError` | 未標記 | ✅ 通過 |
| UT-HC-03 | httpapi／`TestAdapterRateLimit429AndRedirectDoesNotCount` | 未標記 | ✅ 通過 |
| UT-HR-02～03、UT-HS-02 | httpapi／`TestRedirectAndStatsErrors` | 未標記 | ✅ 通過（UT-HR-03「不計次」由 E2E-06 斷言） |
| UT-OBS-01 | httpapi／observability_test；observability／hooks_test | 標記 | ✅ 通過 |
| IT-STORE-01～03 | integration／api_test；infra/sqlite | 標記 | ✅ 通過 |
| IT-FLOW-01～02 | integration／api_test | 標記 | ✅ 通過 |
| IT-FLOW-03 | application／`TestCreateShortURLAllowsDuplicateLongURL`（單元層，非整合層） | 未標記 | ⚠️ 通過但層級不符（DEF-002） |
| CT-CREATE-01～02、CT-REDIR-01～02、CT-STATS-01、CT-ERR-01 | httpapi／contract_test（`TestCT_*` 6 項＋CT-ERR-01 子測試 8 項；PR #26 已合併，merge `f646bee`，綁定 head `49e29d8`，測試檔與 `e224253` 相同） | 標記 | ✅ main f646bee 實跑 **14 PASS**（v0.5）；DEF-001 **關閉**。殘餘：期望值寫死、不偵測 openapi.yaml 漂移（DEF-005，v0.5 突變複驗仍 PASS）；redirect／stats 的 500 沒測（DEF-006）；皆 G5 前 |
| E2E-01～04、E2E-06 | integration／api_test | 標記 | ✅ 通過 |
| E2E-05-01～11 | integration／ratelimit_e2e_test（TestENG011_* 19 項） | 未標記（以 S6 處置表追溯） | ✅ 通過 |
| ENG-016 標頭 E2E（CHK-SEC-008 自動化部分） | integration／security_headers_test | — | ✅ 通過（G5 仍需演示環境擷取） |
| CHK-* | 檢查表類，屬 G5 舉證 | — | ⏳ 未執行（不屬 G4） |

> **附註（2026-10-22 UTC+8，PR #26 驗收；檔內版本不升）**：CT 列已更新為 PR #26 head `e224253` 的結果。DEF-001 改為「有條件關閉」，條件是研發補書面例外並登錄 DEF-005／DEF-006（期限 G5 前）。PR 合併後，main 的結果要重跑才算數。證據：`evidence/pr26-*.txt`。

---

## 7.1 PR → 任務 → 需求 → 測試 → 結果（v0.5 新增；PR #13～#26）

| PR | 任務／DR | 對應 REQ／NFR／SEC | 變更類型 | 對應 TC／檢查 | PR checks | main f646bee 結果 |
|---|---|---|---|---|---|---|
| #13 | ENG-020 | NFR-004；SEC-009～011 | 文件（G3-R1 落地證據） | CHK（G5 舉證） | 4／4 ✅ | 不適用（文件） |
| #14 | ENG-020 | — | CLOSED 未合併 | — | — | 不在範圍 |
| #15 | ENG-020（CR-002） | NFR-004；SEC-009～011 | CI | CI 四項 check | 4／4 ✅ | run 39737480926 ✅ |
| #16 | ENG-016 | SEC-008 | 測試 | ENG-016 標頭 E2E（CHK-SEC-008 自動化部分） | 4／4 ✅ | ✅ 通過 |
| #17 | ENG-007 | REQ-001；NFR-006；SEC-012 | 程式 | UT-HC-*、bodylimit 測試 | 4／4 ✅ | ✅ 通過 |
| #18 | ENG-011 | SEC-007；NFR-003（輔助） | 測試 | E2E-05-01～11（TestENG011_* 19 項） | 4／4 ✅ | ✅ 19 PASS |
| #19 | ENG-021 | NFR-002；NFR-006 | 文件 | 不適用 | 4／4 ✅ | 不適用 |
| #20 | ENG-007、DR-G4-01 | REQ-001；NFR-006；SEC-012 | 文件／契約（OpenAPI） | CT-*（#26 驗證） | 4／4 ✅ | ✅ CT 14 PASS |
| #21 | ENG-011 | SEC（安全審查文件） | 文件 | 不適用 | 4／4 ✅ | 不適用 |
| #22 | ENG-007、DR-G4-01 | SEC-012（SCR-041） | 文件 | 不適用 | 4／4 ✅ | 不適用 |
| #23 | ENG-011 | SEC（安全審查 08 v0.1.5） | 文件 | 不適用 | 4／4 ✅ | 不適用 |
| #24 | ENG-021 | NFR-002；NFR-006；SEC-012 | 文件（README 431） | CR-004 首測 AC-6～8（18 號） | 4／4 ✅ | 不適用；admin bypass 合併（G4-R4 待追認） |
| #25 | ENG-020（分支名 ENG-022，偏差） | NFR-004；SEC-009～011 | CI 強化 | CI 四項 check | 4／4 ✅ | run 39737480926 ✅ |
| #26 | ENG-013 | REQ-001／002／004／007／009；SEC-001、003、014 | 測試 | CT-CREATE-01～02、CT-REDIR-01～02、CT-STATS-01、CT-ERR-01 | 4／4 ✅ | ✅ 14 PASS（DEF-005／006） |

> REQ／NFR／SEC 欄依 PR 標題、12 號 §1／§4 與審查部逐 PR 紀錄；#21～#25 對應欄為品保依 PR 標題與審查紀錄整理，細項以審查部紀錄為準。

### 7.2 開放缺陷（v0.5，截至 f646bee）

DEF-002、DEF-003、DEF-005、DEF-006 開放（輕微，G5 前）；DEF-001 關閉；DEF-004 本機已修，待 G4-QA-R2。明細見 12 號 §7、19 號 §4。

## 8. 流程類任務（非 REQ／NFR／SEC；v0.4 新增）

| 任務 | 來源 | 設計／規範文件 | TC | 對應 | 結果 |
|---|---|---|---|---|---|
| ENG-022 純文件 PR 分支例外 | CR-003 v0.1（本文 c55ad3ca 已核實；依 R10 待使用者確認 D-15 才生效，G4-R5） | `08-engineering/02-branch-commit-conventions`（待升 v0.2）；`01-task-breakdown`（待升 v0.2） | **CHK-ENG-022** | NFR-007（輔助） | ⏳ 待 CR-003 生效與文件落地後驗收 |
| CR-004 ruleset stale approval | CR-004 v0.3（eae83378…） | ruleset 24629854 | CHK-CR-004（AC-CR004-1～9） | SCR-042／RR-011（輔助） | ✅ 首測通過（18 號 v0.4；AC-9 不適用）；撤銷行為未實測 |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G1：需求欄齊；設計／TC TBD |
| v0.1.1 | 2026-10-07 | 品保部 | G1-R1：SEC-001～015 |
| v0.2 | 2026-10-07 | 品保部 | G2：設計／API 欄；SEC-016／017 |
| v0.2.1 | 2026-10-07 | 品保部 | 關閉 G2-QA-R1；016／017 設計欄回填 |
| v0.3 | 2026-10-07 | 品保部 | **G3**：TC 欄填入實際 ID；公開方法覆蓋勾稽 100%；對齊案例目錄 v0.1 |
| v0.4 | 2026-10-22 | 品保部 | **G4／S6**：SEC-007 TC 細分 E2E-05-01～11；REQ-003、NFR-003（輔助）、NFR-006、SEC-005 回填；新增 §7 實作對應與執行結果；修正簡體字「舉證」；新增 §8 流程類 ENG-022／CHK-ENG-022（CR-003） |
| v0.5 | 2026-10-22 | 品保部 | **G4 整體／G4-R2**：§7 基準改 main f646bee（181 PASS、CT 14 PASS）；CT 列改通過、DEF-001 關閉；新增 §7.1 PR #13～#26 追溯、§7.2 開放缺陷；§8 更新 CR-003 狀態、新增 CR-004 列 |
