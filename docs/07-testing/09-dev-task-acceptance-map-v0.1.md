---
文件：研發任務↔驗收／TC 對照
版本：v0.3
狀態：審查中
負責角色：品保部（骨架）／研發部（ENG 回填）
最後更新：2026-10-22
對應需求：REQ-001～REQ-012、SEC 關鍵項
專案代號：SHORTURL
備註：研發正式 ID 已回填自 `docs/08-engineering/01-task-breakdown-v0.1.md`（ENG-001～021）；ENG-022 依 CR-003（`docs/99-changelog/CR-003-branch-docs-exception-v0.1.md`）新增，待研發部核准
---

# 研發任務 ↔ 驗收／TC 對照

> **ENG 已回填（G3-QA-S3）**：下列「研發正式 ID」取自任務拆解表既有 `ENG-xxx`，未新編 ID。  
> 一對多／跨 WP 已於備註註明；骨架 `T-*` ID 保留供稽核。

---

## 1. 工作包一覽

| WP ID | 名稱 | 範圍摘要 | 研發正式 ID | 狀態 |
|---|---|---|---|---|
| WP-DOM-01 | 領域：驗證與短碼 | UrlValidator、ShortCodeGenerator | ENG-005、ENG-006 | 已對齊 |
| WP-DOM-02 | 領域：應用與導向服務 | ShortUrlApplicationService、RedirectService | ENG-007、ENG-008、ENG-009 | 已對齊 |
| WP-INF-01 | 基礎設施：Store／Counter | UrlRepository、ClickCounter、InMemoryOrDbUrlStore | ENG-003、ENG-004 | 已對齊（schema＋Store） |
| WP-API-01 | 適配：HTTP 建立 | handleCreate、createShortUrl operation | ENG-013 | 已對齊 |
| WP-API-02 | 適配：HTTP 導向 | handleRedirect、redirectShortCode | ENG-014 | 已對齊 |
| WP-API-03 | 適配：HTTP 統計 | handleStats、getUrlStats | ENG-015 | 已對齊 |
| WP-API-04 | 錯誤對映 | ErrorMapper | ENG-010 | 已對齊 |
| WP-SEC-01 | 限流守衛 | RateLimitGuard | ENG-011 | 已對齊 |
| WP-OPS-01 | 可觀測性 | ObservabilityHooks | ENG-012 | 已對齊 |
| WP-CI-01 | CI 測試與資安門檻接入 | test job＋SAST／SCA／secrets | ENG-018、ENG-019、ENG-020 | 已對齊（一對多：單測／整合／CI） |

---

## 2. 任務 ↔ 驗收 TC

| 任務 ID | 研發正式 ID | 公開方法／operationId | 驗收要點 | 主要 TC／UT | REQ／SEC | 備註 |
|---|---|---|---|---|---|---|
| T-DOM-01 | ENG-005 | validateLongUrl | 空／非 URL／非 http(s)／L／L+1 | UT-VAL-01～05 | REQ-005、011；SEC-001 | |
| T-DOM-02 | ENG-005 | validateShortCode | pattern 8；非法拒 | UT-VAL-06～07 | REQ-006；SEC-013 | 與 T-DOM-01 同 ENG |
| T-DOM-03 | ENG-006 | generate | base62×8；CSPRNG；無序列 | UT-GEN-01～02 | REQ-006、010；SEC-006 | |
| T-DOM-04 | ENG-007 | createShortUrl | 成功欄位；失敗不落庫；衝突重試 | UT-APP-01～03；IT-FLOW-02 | REQ-001、007、010 | |
| T-DOM-05 | ENG-009 | getStats | 存在／不存在 | UT-APP-04～05 | REQ-004、008 | |
| T-DOM-06 | ENG-008 | redirect | Location＝庫值；成功 +1；失敗不計 | UT-RED-01～02 | REQ-002、003、008；SEC-014 | |
| T-INF-01 | ENG-004（schema：ENG-003） | save／find／exists | 參數化；UNIQUE；無 IP／UA | UT-REPO-01～03；IT-STORE-01、03 | REQ-001、010、012；SEC-004、016 | **一對多**：表結構 ENG-003；讀寫 ENG-004 |
| T-INF-02 | ENG-004 | incrementAtomic／getCount | 原子 +1；並發正確 | UT-CLK-01～02；IT-STORE-02 | REQ-003、004；SEC-016 | |
| T-API-01 | ENG-013 | handleCreate／createShortUrl | 201／400／429；契約 | UT-HC-01～03；CT-CREATE-*；E2E-01～02、05 | REQ-001、007；SEC-001、007 | |
| T-API-02 | ENG-014 | handleRedirect／redirectShortCode | 302／400／404／429；參數改寫無效 | UT-HR-01～03；CT-REDIR-*；E2E-01、03～04、06 | REQ-002、003、006、008；SEC-013、014、007 | |
| T-API-03 | ENG-015 | handleStats／getUrlStats | 200／400／404；契約 | UT-HS-01～02；CT-STATS-01；E2E-01 | REQ-004、006、008 | |
| T-API-04 | ENG-010 | toHttpResponse | 四碼 HTTP 對映；無洩漏 | UT-ERR-01～02；CT-ERR-01 | REQ-009；SEC-003 | |
| T-SEC-01 | ENG-011 | checkCreate／checkRedirect | 門檻內外；超限不寫／不導；來源判定／可信代理／IPv6 /64／並發；日誌無位址 | UT-RL-01～02；UT-HC-03；E2E-05（E2E-05-01～11） | SEC-007；SEC-005 | S6 處置見 `11-s6-eng011-disposition-v0.1.md`（v0.3） |
| T-OPS-01 | ENG-012 | ObservabilityHooks.* | 成功失敗可觀測；無敏感原文 | UT-OBS-01；CHK-NFR-006 | NFR-006；SEC-005 | |
| T-CI-01 | ENG-020 | （流水線） | test＋SAST／SCA／secrets 門檻 | CHK-SEC-009～011；CHK-NFR-004 | NFR-004；SEC-009～011 | 測試鉤子另見 ENG-018／019 |
| T-DOC-01 | ENG-022 | （無公開方法；流程規範） | 見 §2.1 AC-022-1～6 | **CHK-ENG-022**（新增） | NFR-007（流程合規，輔助）；無 REQ／SEC | CR-003；條款生效以研發部核准 CR-003 為準 |
| T-E2E-01 | ENG-019（組裝前置：ENG-017） | 三 API 串接 | 主路徑＋關鍵否定 | E2E-01～06；IT-FLOW-01、03 | REQ-001～012 抽樣 | **一對多**：整合測 ENG-019；可啟動組裝 ENG-017 |

### 2.1 ENG-022 驗收條件（CR-003：純文件 PR 分支例外條款）

> 依據：CR-003 v0.1 §2.1（研發部 2026-10-19 14:59 更正版原文）與 §2.2（本 CR 追加規則，非原文）；先例 `docs/gates/gate-G4-PR20-2026-10-19.md`。CR-003 狀態：研發已書面核准 c55ad3ca、待審查部確認；以下 AC-022-1～6 於 CR-003 核准、ENG-022 落地時驗收。

| AC | 驗收條件 | 驗證方式 | 對應 |
|---|---|---|---|
| AC-022-1 | `08-engineering/02-branch-commit-conventions` 升 v0.2，例外條款與 CR-003 §2.1 逐字一致，含五要素：①只適用只改 docs/ 的純文件 PR；②分支名可用 DR／CR 編號取代 ENG，格式 `docs/<階段>-<DR 或 CR 編號>-<kebab>`；③PR 標題與 commit trailer 仍帶任務拆解既有的 ENG 編號（如有）及 DR 編號；④描述關聯段列涉及的 REQ／SEC；⑤改到程式或測試即不適用。§2.2 追加規則分開寫入並標明非原文：DR／CR 編號須可追溯到已登錄編號；CR 分支且無 ENG 時以 CR 編號為必要編號，有 DR 時一併帶；被程式／測試／CI 讀取的 docs/ 檔照常走資安審查 | 文件 diff 逐句比對 | CHK-ENG-022 第 1 項 |
| AC-022-2 | 條款明寫「改到程式或測試（docs/ 以外任何路徑）即不適用」 | 文件審查 | CHK-ENG-022 第 2 項 |
| AC-022-3 | `08-engineering/01-task-breakdown` 升 v0.2 新增 ENG-022；`08-engineering/README` 索引同步 | 文件審查 | CHK-ENG-022 第 3 項 |
| AC-022-4 | 條款生效後，任一使用 DR／CR 分支名的 PR：`git diff --name-only base...head` 全部位於 `docs/`；標題與 trailer 帶既有 ENG（如有）及 DR 編號（§2.1 ③）；CR 分支且無 ENG 時帶 CR 編號，有 DR 時一併帶（§2.2）；描述列涉及 REQ／SEC | PR 抽查（每個 DR／CR 分支 PR 全查） | CHK-ENG-022 第 4 項 |
| AC-022-5 | 回滾路徑（撤回條款、恢復分支名必含 ENG、不追溯）已記錄 | 文件審查 | CHK-ENG-022 第 5 項 |
| AC-022-6 | 02 規範 v0.2 寫明安全部三項前提：(1) 例外只涵蓋分支命名，四項 required checks 與審查核准規則不變（核准與 head 綁定靠合併前人工核 head，D-09／SCR-042），不得加 paths-ignore；(2)「只改 docs/」以實際變更檔判定並附 `git diff --name-only`；(3) 被程式／測試／CI 讀取的 docs/ 檔（如 `docs/04-api/openapi.yaml`）照常走資安審查 | 文件審查＋PR 抽查（含 required checks 紀錄） | CHK-ENG-022 第 6 項 |

**否定案例（必查）**：(a) DR／CR 分支名的 PR 只要改到 `docs/` 以外的任一檔（如 `internal/`、`cmd/`、`.github/`、`go.mod`）→ CHK-ENG-022 判定「不符」，分支須改回含 ENG。(b) 分支名使用未登錄的 DR／CR 編號 → 判定「不符」（CR-003 §2.2）。

---

## 3. 驗收對帳方式（建議）

1. 研發開 PR 時於描述填：`ENG-xxx`＋`TC IDs`（對齊 DoD R2／§9）。  
2. 品保依本表抽查：方法是否有對應綠燈測試。  
3. 骨架 `T-*` ID 保留於左欄以利稽核；以「研發正式 ID」為合併／追蹤主鍵。  
4. DoD（`08-definition-of-done`）任務級迷你清單須全部勾選。  
5. DoR（`docs/08-engineering/03-definition-of-ready-v0.2.md`）開工前勾選 §1。

---

## 4. 覆蓋摘要

| 檢核 | 結果 |
|---|---|
| 公開方法皆落入至少一任務 | ✅ |
| 三 operationId 各有 API 任務＋CT／E2E | ✅ |
| 骨架列皆已填研發正式 ENG ID | ✅（17／17 驗收列；v0.3 新增 T-DOC-01／ENG-022） |
| WP 皆已填 ENG ID | ✅（10／10 WP） |

### 4.1 任務拆解有、本骨架未單列之 ENG（缺口／掛載說明）

| ENG ID | 標題 | 說明 |
|---|---|---|
| ENG-001 | 專案骨架 | 基礎設施前置；無獨立 WP／T 列；支撐全 WP |
| ENG-002 | 設定與環境變數 | 同上；限流門檻／DB 路徑支撐 WP-SEC-01、WP-INF-01 |
| ENG-016 | 安全回應標頭（SEC-008） | 骨架無對應 WP；驗收以 CHK-SEC-008 為主；建議後續補 WP-SEC-02 或掛 WP-API-* |
| ENG-017 | 組裝／本機可啟動 | 已註於 T-E2E-01 前置；非獨立驗收列 |
| ENG-018 | 單元測試鉤子 | 已併入 WP-CI-01；覆蓋 T-DOM／T-API 相關 UT |
| ENG-022 | 分支規範純文件例外（CR-003） | v0.3 已單列 T-DOC-01，驗收見 §2.1；不屬缺口 |
| ENG-021 | README／本機啟動 | 文件任務；對齊 DoD D3／C3；無獨立 TC 列 |

以上缺口**不阻擋** G3 對帳；ENG-016／021 建議品保後續可選補列。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 品保部 | G3：WP／任務骨架；待研發對齊；方法／API↔TC |
| v0.2 | 2026-10-07 | 研發部 | 回填 ENG-001～021 對應；10 WP＋16 驗收列已對齊；註明一對多與未單列 ENG 缺口 |
| v0.3 | 2026-10-22 | 品保部 | CR-003：新增 T-DOC-01／ENG-022 與 §2.1 AC-022-1～6、新 CHK-ENG-022；S6：T-SEC-01 改指 E2E-05-01～11、補 UT-HC-03 |
