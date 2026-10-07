---
文件：分支與 Commit 規範
版本：v0.1
日期：2026-10-07
狀態：Draft／Ready for Review
負責角色：研發部
作者：研發部
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# 分支與 Commit 規範：SHORTURL

## 依據文件

| 文件 | 路徑 |
|---|---|
| 任務拆解 | `docs/08-engineering/01-task-breakdown-v0.1.md` |
| 需求規格 | `docs/01-requirements/02-requirements-spec-v0.1.md` |
| SEC v0.2 | `docs/06-security/01-sec-requirements-v0.2.md` |
| CI 資安門檻 | `docs/06-security/03-ci-security-gates-v0.1.md` |
| DoR | `docs/08-engineering/03-definition-of-ready-v0.2.md` |
| 決策 | D-04（技術棧） |

## 1. 工作流程（選定）

採 **簡化 trunk＋feature**（適合小專案、單／少人數）：

| 規則 | 說明 |
|---|---|
| 主分支 | `main`＝唯一長期分支；隨時應可建置／測試通過 |
| 開發 | 自 `main` 開短生命週期 feature／fix 分支 |
| 合併 | **僅能經 Pull Request**；禁止直接 push `main` |
| 刪除 | 合併後刪遠端 feature 分支 |
| 不採用 | 完整 Git Flow（無 develop／release／hotfix 常態分支）；緊急修補仍用 `fix/…`→PR |

## 2. 分支命名

```text
<type>/<ENG-ID>-<short-kebab>
```

| type | 用途 | 範例 |
|---|---|---|
| `feat` | 新功能／規格實作 | `feat/ENG-006-short-code-generator` |
| `fix` | 缺陷修復 | `fix/ENG-014-location-rewrite` |
| `test` | 僅補測試 | `test/ENG-019-api-integration` |
| `chore` | 建置、CI、雜項（非功能行為） | `chore/ENG-020-ci-security-gates` |
| `docs` | 僅文件 | `docs/ENG-021-readme-local` |
| `refactor` | 重構（行為不變） | `refactor/ENG-004-sql-bind` |

**規則：**

- 必須含 **ENG-xxx**（對應任務拆解）；一分支主攻一任務為佳。
- `short-kebab`：英文小寫、簡短、可讀。
- 禁止：`main` 上直接開發、個人名當唯一分支名、無 ENG-ID 的長期分支。

## 3. Commit message 規範

採 **Conventional Commits**，並強制關聯 **REQ**（或 NFR／SEC）與 **ENG**：

```text
<type>(<optional-scope>): <摘要>

<可選正文>

REQ-xxx ENG-xxx
```

多編號時空白分隔，例如：`REQ-001 REQ-010 ENG-007`；純基礎設施可寫 `NFR-004 ENG-020` 或 `SEC-009 ENG-020`。

### 3.1 type 與 scope

| type | 何時用 |
|---|---|
| `feat` | 使用者可見／API 行為新增 |
| `fix` | 修復錯誤 |
| `test` | 測試 |
| `docs` | 文件 |
| `chore` | CI、依賴鎖定、工具 |
| `refactor` | 無行為變更大改 |
| `perf` | 效能（本 MVP 少用） |

scope 建議對應套件：`api`、`store`、`validator`、`generator`、`redirect`、`ci` 等。

### 3.2 範例

```text
feat(generator): add CSPRNG base62 length-8 short codes

Implements ADR-003; rejects math/rand as sole source.

REQ-006 REQ-010 SEC-006 ENG-006
```

```text
fix(redirect): keep Location equal to stored long_url

Ignore query/header rewrite attempts (SEC-014).

REQ-002 SEC-014 ENG-014
```

```text
chore(ci): add SAST SCA and secrets jobs as required checks

SEC-009 SEC-010 SEC-011 ENG-020
```

### 3.3 禁止

- 無 `ENG-xxx` 的功能／修復 commit（純 typo 於 docs 可例外，仍建議標 `docs`）。
- 一次 commit 塞多個無關 ENG（應拆）。
- Commit 含真實機密、連線字串、私鑰（SEC-011）。

## 4. Pull Request 最低要求

### 4.1 標題

```text
<type>: <摘要> (ENG-xxx)
```

例：`feat: implement createShortUrl application service (ENG-007)`

### 4.2 描述（模板）

PR 描述須含下列區塊（可複製）：

```markdown
## 關聯
- 任務：ENG-xxx
- 需求：REQ-xxx／NFR-xxx
- 資安：SEC-xxx（若適用；否則寫「—」）
- 規格：方法規格〈類別.方法〉；OpenAPI operationId（若適用）

## 做了什麼
- …

## 驗收對照（DoR／任務驗收）
- [ ] 任務驗收標準已自測通過
- [ ] 新增／更新測試（單元／整合）
- [ ] 無新增禁止欄位（visitor_ip／user_agent）
- [ ] SQL 皆參數化（若動到 Store）

## 測試證據
- 指令：`go test …`
- 結果摘要：…

## 風險／殘餘
- （無／或指向 RR／D-06）
```

### 4.3 合併條件

| 條件 | 說明 |
|---|---|
| 審查 | ≥1 人 Code Review（小專案可由另一角色或總協調指定）；涉及 SEC-016／017／006 建議安全部抽看 |
| CI | 全部 required checks 綠（含 SEC-009～011 門檻，見 CI 文件） |
| 衝突 | 已 rebase／merge 最新 `main` |
| 範圍 | PR 對應任務 DoR 已勾且不超出 OUT 範圍 |

## 5. 保護 `main`（強制）

| 規則 | 狀態 |
|---|---|
| 禁止直接 push `main` | **必須**（分支保護） |
| 禁止 force-push `main` | **必須** |
| PR 必經 CI required checks | **必須**（G3／G4 與維運對齊後生效） |
| 禁止無紀錄關閉資安掃描規則 | **必須**（SEC-009～011） |

設定由**維運**執行、**安全**確認門檻語意、**研發**確保 job 可綠。

## 6. 與 REQ／任務追溯

| 產出 | 必須出現的 ID |
|---|---|
| 分支名 | `ENG-xxx` |
| Commit trailer／末行 | `REQ-…`（或 NFR／SEC）＋`ENG-xxx` |
| PR 標題／描述 | `ENG-xxx`＋對應 REQ／SEC |
| 品保追溯矩陣 TC 欄 | 後續由品保填；研發於 PR 註明可對應場景即可 |

無法對應既有 REQ／SEC 時：標「—」並於 PR 說明（僅限 chore／純文件且不改行為）；**不得發明**新 REQ／SEC 編號。

## 7. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 研發部 | G3 初稿：trunk＋feature；Conventional Commits＋REQ／ENG |
