---
文件：機密管理
版本：v0.1
狀態：審查中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05、SEC-011、SEC-017
專案代號：SHORTURL
---

# 機密管理（v0.1）

> 原則：**機密永不進入 git**。版控僅允許 `.env.example` 與文件中的**明顯假值／佔位符**。  
> CI 機密掃描門檻對齊 [`../06-security/03-ci-security-gates-v0.3.md`](../06-security/03-ci-security-gates-v0.3.md)（現行，D-07；SEC-011：**真實機密一律擋**）。

---

## 1. 原則

| 規則 | 說明 |
|---|---|
| 不進版控 | API 金鑰、私鑰、密碼、權杖、TLS 私鑰、Tunnel token、雲端憑證等皆不可 commit |
| 允許進版控 | `.env.example`（僅鍵名＋假值）、文件表格中的佔位說明、公開設定（如預設限流數字） |
| 最小權限 | 誰需要誰持有；CI 優先使用內建 `GITHUB_TOKEN`，避免另發長期 PAT |
| 可審計 | 例外／allowlist 須有紀錄與安全部審核 |

---

## 2. 存放位置

| 場所 | 用途 | 注意 |
|---|---|---|
| **GitHub Actions Secrets**（repo） | 全 workflow 共用、非環境特有之 CI 機密 | 勿在 log 印出 |
| **GitHub Environments secrets**（`staging`／`production`） | 依環境隔離之部署機密 | 建議 production 加審批保護 |
| **本機 `.env`** | 開發者 local | **必須**列入 `.gitignore`；勿複製進 PR 描述 |
| 主機／容器環境變數或祕密掛載 | 執行期注入 | 檔案權限受限；備份勿含明文散落 |
| 密碼管理器／離線保管（建議） | 團隊共享之演示憑證 | 零預算可用既有免費方案或實體保管 |

**MVP 說明**：本專案為練手／零現金。若演示僅本機或 Tunnel 且**無**雲端 API 金鑰，文件標明「**MVP 可能無雲端金鑰**」即可；一經產生任何 token／私鑰，一律走本文件存放規則。

---

## 3. 機密分類

| 分類 | 示例 | MVP 預期 | 存放 |
|---|---|---|---|
| 執行期 | 若未來有 signing key、admin token、加密鍵 | **可能無**；有則必管 | Environment secrets／主機注入；不進 git |
| 傳輸 | TLS 私鑰、Cloudflare Tunnel token | 選定託管後才有 | Environments secrets 或主機；輪替見 §5 |
| CI 用 token | 推送映像、通知 webhook 等 | **盡量無**；優先 `GITHUB_TOKEN` | Actions Secrets；範圍最小、短期 |
| 非機密組態 | `APP_ENV`、`HTTP_ADDR`、限流數字、公開 `BASE_URL` | 有 | 環境變數／`.env.example` |

---

## 4. Allowlist（僅假值）

| 項目 | 規定 |
|---|---|
| 用途 | 僅排除**明顯測試用假值**，避免 gitleaks 誤判阻擋合法文件／測試 |
| 允許內容 | 如 `REPLACE_ME`、`changeme-not-a-secret`、文件中的空白／星號遮罩說明；**禁止**把真實機密加入 allowlist |
| 建議路徑 | `.gitleaks.toml` 或 `.gitleaksignore`（實作時擇一；放 repo 根目錄） |
| 審核 | **新增／修改 allowlist 須安全部審核**並留 PR 紀錄；指向門檻文件例外條款 |
| 禁止 | 關閉整條機密掃描規則；以編碼／更名繞過掃描（門檻 v0.3 §12「禁止事項」） |

例外流程：不得無紀錄關閉規則 → 見 [`../06-security/03-ci-security-gates-v0.3.md`](../06-security/03-ci-security-gates-v0.3.md) §6 與 [`01-ci-pipeline-v0.2.md`](./01-ci-pipeline-v0.2.md) §8。

---

## 5. 輪替與外洩應變（對齊 SEC-011）

| 情況 | 行動 |
|---|---|
| CI／掃描發現疑似**真實**機密進版控或 PR | **擋合併**；不得要求關閉規則過關 |
| 已推送或已外洩 | 視為事故：立即**輪替／作廢**該機密；自版控歷史移除（rewrite 或 GitHub 官方移除指引）；通知安全部／總協調 |
| 事故流程 | 詳細 runbook 屬 `docs/13-incident/`（後續）；G3 最低要求為「擋合併＋輪替＋紀錄」 |
| 定期輪替 | 有長期 token 者建議訂週期（例如 90 日）或於人員異動時輪替；MVP 無金鑰則 N/A |

---

## 6. `.gitignore` 建議條目

```gitignore
# 機密與本機環境
.env
.env.*
!.env.example

# SQLite 與本地資料（勿把正式／測試 DB 推進版控）
*.db
*.db-journal
*.db-wal
*.db-shm
/data/

# 憑證與金鑰材料
*.pem
*.key
*.p12
*.pfx
id_rsa
id_rsa.*
!.well-known/  # 若未來有公開挑戰檔再收斂規則

# OS／IDE
.DS_Store
.idea/
.vscode/
```

> 實作時由研發／維運將條目納入 repo 根目錄 `.gitignore`；本段為建議，非要求在本閘建立真實檔。

---

## 7. 禁止事項清單

1. 將真實機密 commit、強制推送「之後再刪」蒙混過關。  
2. 在 issue、PR、聊天、文件中貼上真實金鑰或「看起來像真的」長隨機密鑰。  
3. 為合併而關閉 `secrets-gitleaks` 或 GitHub push protection。  
4. 以 base64／拆字／假副檔名等方式繞過掃描。  
5. 將 production SQLite 或含機密之備份放在可公網讀取位置（SEC-017）。  
6. 在 CI log 印出 secrets（`echo`、除錯 dump）。  
7. 共用同一組 production 憑證給所有個人本機長期使用（改為每人最小權限或僅注入部署環境）。

---

## 8. 與 CI 的銜接

| Job | 關係 |
|---|---|
| `secrets-gitleaks` | 合併閘；真實機密 → 失敗 |
| （建議）GitHub secret scanning＋push protection | 推送前／後雙層防護 |
| Allowlist 變更 | 須安全部審核；變更本身可被 gitleaks 設定檔 diff 審查 |

流水線設計見 [`01-ci-pipeline-v0.2.md`](./01-ci-pipeline-v0.2.md)。

---

## 9. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3：存放、分類、allowlist、輪替、gitignore、禁止事項 |
| v0.1.1 | 2026-10-07 | 維運部 | 一致性簽核：門檻引用改指 `03-ci-security-gates-v0.2.md`；job／阻擋語意未放寬 |
| v0.1.2 | 2026-10-07 | 維運部 | 引用同步（檔名與 frontmatter 維持 v0.1，沿用本檔 v0.1.1 列之小版號慣例）：CI 資安門檻引用由 `03-ci-security-gates-v0.2.md` 改指現行 `03-ci-security-gates-v0.3.md`（v0.2 已由 v0.3 取代，D-07；禁止事項改引 §12、例外改引 §6）；流水線引用改指 `01-ci-pipeline-v0.2.md`；內容與阻擋語意未變 |
