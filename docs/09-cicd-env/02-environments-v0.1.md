---
文件：環境區隔（開發／測試／正式演示）
版本：v0.1.2
狀態：審查中
負責角色：維運部
最後更新：2026-10-19
對應需求：NFR-004、SEC-007、SEC-008、SEC-012、SEC-017；SCR-009、SCR-010、SCR-013、SCR-015；RR-004、RR-008、RR-009；DR-G4-01（C3）
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
| Cloudflare Tunnel（或同類） | 為本機／家用伺服器提供 HTTPS 入口 | 免費檔可達 HTTPS | Token 屬機密；勿進版控；仍須限流與檔案權限；**必須設定 `TRUSTED_PROXY_CIDRS`**＝Tunnel 連入位址的 CIDR（同機時 peer 為 127.0.0.1），否則全站共用一個限流 bucket（SCR-009，見 §6、[G5-01](./07-g5-deploy-checklist-v0.1.md)） |
| 其他免費靜態＋分離 API | 視架構 | — | 本專案為單服務 Go＋SQLite，需確認持久磁碟 |

**狀態**：文件化建議；**待 G5 前選定**。選定後更新本文件與部署手冊，並驗證 SEC-012／SEC-017。

---

## 6. 環境變數清單（無真實值）

> 本表只列**程式現行有實際使用**的變數。有載入但未使用、或尚未實作的變數一律列在 §6.2「保留／未實作」，**不得**當成程式已支援的功能。

| 變數名 | 必填 | 範例假值／格式 | 說明 |
|---|---|---|---|
| `HTTP_ADDR` | 是 | `127.0.0.1:8080` | 監聽位址 |
| `DATABASE_PATH` | 是 | `./data/local/shorturl.db` | SQLite 檔路徑（每區不同） |
| `BASE_URL` | 是 | `http://127.0.0.1:8080` | 對外短網址前綴（production 應為 `https://...`） |
| `RATE_LIMIT_CREATE_PER_MIN` | 建議 | `30` | 建立限流 |
| `RATE_LIMIT_REDIRECT_PER_MIN` | 建議 | `120` | 導向限流 |
| `TRUSTED_PROXY_CIDRS` | 條件必填 | `10.0.0.0/8`（示例 CIDR，非機密） | 受信代理 CIDR。服務只要部署在 Cloudflare Tunnel 或反向代理之後就**必填**（＝代理位址 CIDR），否則所有訪客共用一個限流 bucket（[SCR-009](../06-security/11-pr-007-008-security-review-v0.1.md)，11 號檔第 88–93 行；G5-01）；代理必須覆寫 `X-Real-IP`（RR-004；G5-02）。**只有沒有任何代理（直連）時才留空** |

> **禁止**在文件、`.env.example`、issue 中填入真實金鑰或「看起來像真的」長隨機密鑰字串。假值僅用明顯佔位，例如 `changeme-not-a-secret`、`REPLACE_ME`。

完整機密存放規則見 [`03-secrets-management-v0.1.md`](./03-secrets-management-v0.1.md)。

### 6.1 伺服器標頭上限 `MaxHeaderBytes`

`http.Server` 必須**明確設定** `MaxHeaderBytes`，不依賴 Go 預設 1 MiB（建議 ≤ 16 KiB；**超過 16 KiB 須在 PR 附理由並取得安全部書面接受**）；超大標頭回 `431`。此為 DR-G4-01 條件 C3（[10 號檔](../06-security/10-dr-g4-01-body-limit-security-opinion-v0.1.md) §4 第 54 行），必須在 G5 送審前完成；實作屬研發部，維運於部署時驗證（[G5-06](./07-g5-deploy-checklist-v0.1.md)）。目前不是環境變數，未列入上表。

### 6.2 保留／未實作變數（程式不支援，設了也無作用）

以下變數**目前沒有任何作用**，部署時不得依賴，也不得據以宣稱滿足任何資安需求。是否實作，於 G5 決定（不在本次範圍）。

| 變數名 | 現況 | 現行做法 | 依據 |
|---|---|---|---|
| `APP_ENV` | 保留：config 有讀入，程式**未使用** | 環境區隔靠各區獨立的 `DATABASE_PATH`、`BASE_URL` 與部署設定，不靠此變數 | [12 號檔](../06-security/12-pr-011-security-review-v0.1.md) SCR-013；研發部 README（PR #19） |
| `LOG_LEVEL` | 保留：config 有讀入，觀測 hook **未使用**，設了也無法調降 | 日誌量以 §7 保存與大小上限控管；勿記機密／個資 | 12 號檔 SCR-013 第 83–87 行；接線屬研發 |
| `TLS_CERT_FILE`／`TLS_KEY_FILE` | **未實作**：交付的 binary 只跑明文 HTTP（`ListenAndServe`） | **TLS 在 G5 邊緣層終止**（Cloudflare Tunnel／反向代理，見 §5） | [11 號檔](../06-security/11-pr-007-008-security-review-v0.1.md) SCR-010 第 95–100 行；RR-008 |
| `HSTS_ENABLED` | **無作用**：服務只回明文 HTTP，瀏覽器不採信 HTTP 回應中的 HSTS | **HSTS 由 G5 邊緣層送出**：`max-age` ≥ 31536000，不加 `includeSubDomains`、不加 `preload` | SCR-010；RR-008；[G5-04](./07-g5-deploy-checklist-v0.1.md) |

> G5 邊緣層（反向代理 TLS 與 HSTS）是滿足 SEC-008(1) 的**唯一**現行途徑；驗證見 G5-04。

---

## 7. 觀測日誌（DS2）

觀測日誌（PR #11，寫到 stderr）含全部成功事件的短碼，短碼集合加上公開導向端點可列舉全部長網址，敏感度視同 DB 對應索引（[12 號檔](../06-security/12-pr-011-security-review-v0.1.md) 條件 1–3 第 62–65 行、SCR-015 第 97–102 行；[05 RR-009](../06-security/05-residual-risk-register-v0.1.md)）。staging／production 規定：

| 項目 | 規定 |
|---|---|
| 存取 | 權限與 DB 同級（僅維運／使用者）；檔案權限與擁有者不得比 `DATABASE_PATH` 寬 |
| 保存 | ≤ 14 日並自動輪替；容器／服務層設日誌大小上限（SCR-013 維運修補建議 (1)，第 87 行） |
| 外送 | 未過濾的觀測日誌**不得公開或外送**（公開 issue、CI artifact、外部日誌服務的公開儀表板皆不可） |
| 驗證 | G5 前依 [G5-05](./07-g5-deploy-checklist-v0.1.md) 驗證 |

---

## 8. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3：三區區隔、限流預設、SQLite／SEC-017、託管待選、環境變數表 |
| v0.1.1 | 2026-10-19 | 維運部 | 補前次審查缺口：§6 `TRUSTED_PROXY_CIDRS` 改條件必填（Tunnel／反向代理後必填，僅無代理時留空；SCR-009），§5 Cloudflare Tunnel 列註明須設定；`TLS_CERT_FILE`／`TLS_KEY_FILE` 註明行程內 TLS 未實作、TLS 於邊緣終止、HSTS 由邊緣送出、`HSTS_ENABLED` 現行無效（SCR-010、RR-008）；`LOG_LEVEL` 註明未接線（SCR-013）；新增 §6.1 `MaxHeaderBytes`（DR-G4-01 C3）與 §7 觀測日誌（SCR-015、RR-009）；連結 07 G5 檢核清單；修訂紀錄改為 §8；§1～§6 既有行號不變；小版修訂不改檔名 |
| v0.1.2 | 2026-10-19 | 維運部 | 依總協調轉達（研發部 PR #19 README 發現、SCR-010／SCR-013）：§6 主表只留程式實際使用的變數；`APP_ENV`、`LOG_LEVEL`、`TLS_CERT_FILE`／`TLS_KEY_FILE`、`HSTS_ENABLED` 移到新增 §6.2「保留／未實作」，寫明現況與 G5 邊緣層做法（TLS 與 HSTS 由反向代理／Tunnel 處理）；§6.1 補「超過 16 KiB 須附理由並取得安全部書面接受」。是否實作 TLS／LOG_LEVEL 留待 G5。§6 主表行號有變動 |
