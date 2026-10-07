---
文件：環境區隔（開發／測試／正式演示）
版本：v0.1
狀態：審查中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SEC-007、SEC-012、SEC-017
專案代號：SHORTURL
---

# 環境區隔：local／staging／production（v0.1）

> 三區資料與組態分離。SQLite **禁止共用同一檔案**。  
> 限流預設對齊 SEC-007 建議值；HTTPS 對齊 SEC-012；DB 非公網對齊 SEC-017。  
> 零預算託管選項僅文件化建議，**待 G5 前選定**，不強迫本閘架設。

---

## 1. 環境總表

| 環境代號 | 用途 | 誰可存取 | 對外暴露 | HTTPS（SEC-012） | 資料 |
|---|---|---|---|---|---|
| `local` | 開發者本機開發與除錯 | 各開發者本機 | 預設僅 localhost | **不要求**；文件標明「非演示對外」，**不得**以本區宣稱滿足 SEC-012 | 本機獨立 SQLite 檔 |
| `staging` | 測試／CI 可部署之整合驗證 | 專案成員；CI 部署身分 | 可為內網、Tunnel 或受控 URL | 若 URL 可被專案外造訪 → 應 HTTPS；純內網／ephemeral 須文件標明 | 獨立 SQLite；可毀即可重建 |
| `production` | 演示／正式（對外演示入口） | 維運／指定發布者；演示受眾可造訪服務 | 對外可造訪短網址入口 | **應使用 HTTPS（TLS 1.2+）**（SEC-012） | 獨立 SQLite；備份概念見 §4 |

---

## 2. 分區細則

### 2.1 `local`（開發）

| 項目 | 規定 |
|---|---|
| 用途 | 功能開發、單元／本機手動測 |
| 存取 | 開發者本人；不共享 DB 檔給他人當正式來源 |
| 組態來源 | 本機 `.env`（gitignore；見 [`03-secrets-management-v0.1.md`](./03-secrets-management-v0.1.md)）；範本僅 `.env.example` |
| `DATABASE_PATH` 示例概念 | `./data/local/shorturl.db`（路徑可自訂，勿提交 `.db`） |
| 對外暴露 | `HTTP_ADDR` 綁定 `127.0.0.1` 建議；若綁 `0.0.0.0` 僅限受信區域網並自負風險 |
| HTTPS | 不要求 |
| 限流 | 可沿用預設或本機放寬，但**不得**把放寬值帶回 staging／production 組態 |

### 2.2 `staging`（測試／CI 可部署）

| 項目 | 規定 |
|---|---|
| 用途 | PR／合併後整合測、品保驗證、預發布 |
| 存取 | 專案成員＋CI（GitHub Actions 以 Environments secrets／變數注入） |
| 組態來源 | CI／主機環境變數或 GitHub Environment `staging`；禁止把真實機密寫進版控 |
| `DATABASE_PATH` | 獨立路徑，例如 `/var/lib/shorturl/staging/shorturl.db` 或容器 volume；**禁止**與 production 同檔 |
| 對外暴露 | 依託管方式；建議非必要不公開索引 |
| HTTPS | 對外可造訪則應 HTTPS；否則文件標明範圍 |
| 限流 | **必須**套用 §3 預設（或更嚴） |

### 2.3 `production`（演示／正式）

| 項目 | 規定 |
|---|---|
| 用途 | SI-03 演示與對外短網址服務 |
| 存取 | 發布／維運最小權限；無一般開發者直接寫入主機機密 |
| 組態來源 | 主機／Paas／Tunnel 側環境變數或 GitHub Environment `production` |
| `DATABASE_PATH` | 獨立路徑；檔案權限見 §4 |
| 對外暴露 | 僅應用 HTTP(S) 服務埠／入口；**SQLite 檔與管理介面不對網際網路開放**（SEC-017） |
| HTTPS | **應啟用** TLS 1.2+；公開信任憑證或文件化之內部信任（SEC-012） |
| 限流 | **必須**套用 §3 預設（或更嚴） |
| 版控中繼 | 部署產物不得暴露 `.git`（SEC-015；部署章節實作時再驗） |

---

## 3. 限流預設（SEC-007）

具體數值可由設計／維運調整；**本文件給定預設起點**（與 SEC-007 建議一致）：

| 端點類別 | 環境變數（建議名） | 預設值 | 說明 |
|---|---|---|---|
| 建立短碼 | `RATE_LIMIT_CREATE_PER_MIN` | **30** | ≤ 30 次／分／來源 |
| 短碼導向 | `RATE_LIMIT_REDIRECT_PER_MIN` | **120** | ≤ 120 次／分／來源 |
| 來源識別 | `TRUSTED_PROXY_CIDRS` | 空＝直連 peer IP | 僅在受信反向代理後填入 CIDR；防偽造 `X-Forwarded-For` |

超過門檻：回傳明確限流錯誤；**不**建立短碼、**不**導向。若 G3 前應用層尚未實作，須升殘餘並呈報（不得默默省略）——實作屬研發；本文件提供組態預設。

---

## 4. SQLite 與 SEC-017

| 項目 | 規定 |
|---|---|
| 隔離 | 每環境**獨立檔案**；禁止 local／staging／production 共用同一 `.db` |
| 檔案權限 | 僅服務執行身分可讀寫（例如 `0600` 或目錄 `0700`）；勿世界可讀 |
| 網路 | SQLite 為本機檔案；**不得**將資料目錄掛載為公網檔案服務；應用只暴露 HTTP(S) |
| 帳號 | 無獨立 DB 伺服器帳號時，以 OS 使用者最小權限等同滿足「最小權限」精神；若未來改遠端 DB 須另開帳號 CRUD-only |
| 憑證 | 連線字串／金鑰不進版控（見 03） |
| 備份（概念） | production：定期複製 DB 檔至**非公網**備份位置（本機另一路徑、加密物件儲存或離線媒體）；staging／local 可依需要；詳細 RPO／步驟屬 `12-backup-dr`（後續章） |
| 欄位 | 查詢僅業務所需欄位（應用層；對齊 SEC-017） |

---

## 5. 零預算託管選項（待 G5 前選定）

| 選項 | 適用環境 | 優點 | 注意 |
|---|---|---|---|
| 本機 Docker Compose | local／staging／演示練習 | 零費用、可重現 | 對外演示需另解 HTTPS／穿透 |
| 免費 PaaS 試用（如限時免費額度） | staging／production | 少維護 | 注意休眠、額度、SQLite 檔系統是否持久 |
| Cloudflare Tunnel（或同類） | 為本機／家用伺服器提供 HTTPS 入口 | 免費檔可達 HTTPS | Token 屬機密；勿進版控；仍須限流與檔案權限 |
| 其他免費靜態＋分離 API | 視架構 | — | 本專案為單服務 Go＋SQLite，需確認持久磁碟 |

**狀態**：文件化建議；**待 G5 前選定**。選定後更新本文件與部署手冊，並驗證 SEC-012／SEC-017。

---

## 6. 環境變數清單（無真實值）

| 變數名 | 必填 | 範例假值／格式 | 說明 |
|---|---|---|---|
| `APP_ENV` | 是 | `local`／`staging`／`production` | 環境代號 |
| `HTTP_ADDR` | 是 | `127.0.0.1:8080` | 監聽位址 |
| `DATABASE_PATH` | 是 | `./data/local/shorturl.db` | SQLite 檔路徑（每區不同） |
| `BASE_URL` | 是 | `http://127.0.0.1:8080` | 對外短網址前綴（production 應為 `https://...`） |
| `RATE_LIMIT_CREATE_PER_MIN` | 建議 | `30` | 建立限流 |
| `RATE_LIMIT_REDIRECT_PER_MIN` | 建議 | `120` | 導向限流 |
| `TRUSTED_PROXY_CIDRS` | 選 | `10.0.0.0/8`（示例 CIDR，非機密） | 受信代理；無代理則留空 |
| `LOG_LEVEL` | 選 | `info` | 日誌層級；勿記機密／個資 |
| `TLS_CERT_FILE`／`TLS_KEY_FILE` | 選 | 路徑字串 | 若行程內終止 TLS；用 Tunnel／反向代理時可空 |

> **禁止**在文件、`.env.example`、issue 中填入真實金鑰或「看起來像真的」長隨機密鑰字串。假值僅用明顯佔位，例如 `changeme-not-a-secret`、`REPLACE_ME`。

完整機密存放規則見 [`03-secrets-management-v0.1.md`](./03-secrets-management-v0.1.md)。

---

## 7. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3：三區區隔、限流預設、SQLite／SEC-017、託管待選、環境變數表 |
