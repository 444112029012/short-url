---
文件：G5 部署檢核清單
文件編號：09-cicd-env/07-g5-deploy-checklist
版本：v0.1
狀態：草稿待審查
負責角色：維運部
日期：2026-10-07
專案代號：SHORTURL
對應：SEC-007、SEC-008、SEC-011、SEC-012、SEC-015、SEC-017；SCR-008、SCR-009、SCR-010、SCR-013、SCR-015；RR-004、RR-007、RR-008、RR-009；DR-G4-01（C3）；ENG-011、ENG-012、ENG-016；D-11（使用者接受 RR-007～RR-009）
---

# G5 部署檢核清單（v0.1 草稿）

> 用途：G5 送審前，維運部逐項對**實際部署**做驗證，並把證據存檔。本清單只整理既有文件的要求，不新增門檻數字；數值都引用來源文件。  
> 範圍：`production`（演示／正式）。若 staging 與 production 採相同代理拓撲，可先在 staging 演練，但 G5 證據以 production 為準。  
> 前提：託管方式**尚未選定**（[02-environments-v0.1.md](./02-environments-v0.1.md) §5，第 98–107 行）。下表驗證指令中的 `<正式網域>`、`<容器名>`、`<服務名>`、`<DATABASE_PATH>`、`<觀測日誌路徑>`、`<已建立短碼>` 都是佔位符，選定託管後再代入實際值；本文件不填任何真實網域。  
> 行號依 2026-10-07 各來源檔現行版本（05 v0.1.7、10 v0.1、11 v0.1.2、12 v0.1.2、02 v0.1.1、03 v0.1.2、01-sec-requirements v0.2），來源改版時以章節名稱為準。

---

## 1. 檢核表

| 項目 | 要求 | 來源 | 驗證方式 | 負責 | 狀態 |
|---|---|---|---|---|---|
| G5-01 可信代理 CIDR | 服務只要部署在 Cloudflare Tunnel 或反向代理之後，就**必須**設定 `TRUSTED_PROXY_CIDRS`＝代理位址的 CIDR（例：同機 Tunnel 時 peer 為 127.0.0.1）。未設定時，所有訪客的 peer 都是代理位址，全站共用一個限流 bucket（create 全站 30 次／分、redirect 全站 120 次／分），等於自我 DoS | 11 號檔 [SCR-009]（第 88–93 行）；05 RR-007 補償措施（§2.5 第 102 行）；02 §3 第 78 行、§6 第 121 行；限流實作 ENG-011（PR #7） | (1) 設定：`docker exec <容器名> printenv TRUSTED_PROXY_CIDRS` 或 `systemctl show <服務名> -p Environment`，值非空且等於代理位址 CIDR。(2) 行為：用戶端 A 在 1 分鐘內執行 `for i in $(seq 1 31); do curl -s -o /dev/null -w '%{http_code}\n' -X POST https://<正式網域>/api/v1/urls -H 'Content-Type: application/json' -d '{"url":"https://example.com/g5-check"}'; done`，預期前 30 次 `201`、第 31 次 `429`；隨即由不同網路的用戶端 B 送一次相同 POST，預期 `201`（若 B 也得 `429`，代表全站共用 bucket，未設或設錯）。註：此測試會建立測試短碼 | 維運部 | 待 G5 前驗證 |
| G5-02 代理轉傳標頭 | **代理必須覆寫 X-Real-IP（RR-004）**：代理設定須覆寫 `X-Real-IP`（或由研發停用 X-Real-IP 退回），且不得讓用戶端自帶的 `X-Forwarded-For`／`X-Real-IP` 繞過每來源限流。研發端修補（改用 `Header.Values` 串接後由右解析）屬研發，G5 前複核 | 11 號檔 [SCR-008]（第 80–86 行）；05 RR-004（第 25 行、§2.4 第 84–95 行）；11 號檔 SCR-008 狀態（第 86 行：維運部署文件須註記代理必須覆寫 X-Real-IP） | (1) 檢視代理設定：有覆寫 `X-Real-IP` 的設定行（截圖或設定檔片段）。(2) 行為：用戶端 A 執行 `for i in $(seq 1 31); do curl -s -o /dev/null -w '%{http_code}\n' -X POST https://<正式網域>/api/v1/urls -H 'Content-Type: application/json' -H "X-Forwarded-For: 203.0.113.$i" -H "X-Real-IP: 198.51.100.$i" -d '{"url":"https://example.com/g5-check"}'; done`，預期第 31 次仍為 `429`（偽造標頭不能換到新 bucket） | 維運部（代理設定）；研發部（程式修補） | 待 G5 前驗證 |
| G5-03 邊緣限流／單實例 | RR-007 的補償控制：邊緣（Cloudflare）速率限制／WAF；應用限流是單程序記憶體計數，部署須維持**單實例**（多實例計數不共享） | 05 RR-007（第 28 行、§2.5 第 102 行）；11 號檔「殘餘風險」RR-007（第 244 行）；D-11 | (1) 邊緣設定頁截圖（保留網址列），顯示已啟用的速率限制／WAF 規則。(2) `docker ps --filter "name=<服務名>" --format '{{.Names}}'` 只輸出一行（或託管平台實例數畫面截圖＝1） | 維運部 | 待 G5 前驗證 |
| G5-04 HSTS 由邊緣送出 | production 的 HSTS 由 Tunnel／反向代理邊緣送出：`Strict-Transport-Security`，max-age ≥ 31536000（建議設剛好 `31536000`），**不加** `includeSubDomains`、**不加** `preload`（安全部 2026-10-19 確認；日後加 `includeSubDomains` 須先取得安全部書面同意）。**應用程式的 `HSTS_ENABLED` 在現行部署下不會生效**：程式只有明文 `ListenAndServe`、沒有行程內 TLS，HSTS 只在 `r.TLS != nil` 時送出，TLS 又在邊緣終止，所以設 `HSTS_ENABLED=true` 也不會有 HSTS（SCR-010）。不得以應用層開關宣稱已滿足 SEC-008(1)。G5 前必須以**實際 production 回應**驗證 | 01-sec-requirements-v0.2 SEC-008(1)（第 31 行，max-age ≥ 31536000）；11 號檔 [SCR-010]（第 95–100 行）、「HSTS 裁定」（第 213–219 行；不含 preload 見第 216 行）；05 RR-008（第 29 行、§2.6 第 111–121 行）；ENG-016（PR #8）；D-11 | (1) 邊緣：`curl -sS -o /dev/null -D - https://<正式網域>/<已建立短碼>`，回應標頭須有 `strict-transport-security`，其 `max-age` 值 ≥ 31536000（設定值為 31536000 時即顯示該值），且該行**不含** `includeSubDomains`、**不含** `preload`。另對錯誤回應重複一次：`curl -sS -o /dev/null -D - https://<正式網域>/api/v1/urls/<不存在短碼>/stats`。(2) 佐證來源在邊緣：在主機上直連 origin `curl -sS -o /dev/null -D - http://127.0.0.1:8080/<已建立短碼>`（埠依 `HTTP_ADDR`），預期**沒有** `Strict-Transport-Security`，證明標頭不是由應用送出 | 維運部 | 待 G5 前驗證 |
| G5-05 觀測日誌存取與保存 | 觀測日誌（DS2，PR #11 寫到 stderr）含全部成功事件的短碼，短碼集合加上公開導向端點可列舉全部長網址，敏感度視同 DB 對應索引。要求：(1) 存取權限與 DB 同級（僅維運／使用者）；(2) 保存 ≤ 14 日並自動輪替，容器／服務層設日誌大小上限；(3) 不公開、不外送未過濾的觀測日誌（公開 issue、CI artifact、外部日誌服務的公開儀表板皆不可） | 12 號檔「short_code 寫進日誌是否可接受」條件 1–3（第 62–65 行）、[SCR-015]（第 97–102 行）、[SCR-013] 維運修補建議 (1)（第 87 行）；05 RR-009（第 30 行、§2.7 第 122–134 行）；ENG-012（PR #11）；D-11 | (1) 權限：`stat -c '%a %U:%G %n' <DATABASE_PATH> <觀測日誌路徑>`，日誌的權限與擁有者不得比 DB 寬。(2) 輪替：Docker 用 `docker inspect --format '{{json .HostConfig.LogConfig}}' <容器名>`，須有 `max-size`／`max-file` 設定；systemd/journald 用 `systemd-analyze cat-config systemd/journald.conf`，須有 `MaxRetentionSec` ≤ 14 天與大小上限。實際上限值由維運選定後記錄，本文件不預設數字。(3) 不外送：上一步輸出的 log driver 為本機型（非遠端或第三方）；CI workflow 不上傳日誌 artifact（檢視 `.github/workflows/`） | 維運部 | 待 G5 前驗證 |
| G5-06 MaxHeaderBytes | `http.Server` 必須**明確設定** `MaxHeaderBytes`，不依賴 Go 預設 1 MiB；超大標頭回 `431`（必守）。建議值 ≤ 16 KiB；超過 16 KiB 須敘明理由並取得安全部書面接受（安全部 2026-10-19 確認）。G5 送審前完成（可併入 DR-G4-01 實作 PR 或 ENG-016）。實作由研發部負責；維運在部署時驗證 | 10 號檔 §3 P3（第 37 行）、§4 C3（第 54 行）與生效建議（第 58 行）；11 號檔第 28 行（引用 C3，G5 前待補） | (1) 程式：在發布 commit 的 `cmd/shorturl/main.go` 找到 `MaxHeaderBytes:` 明確設定（必過）；值 ≤ 16384 符合建議，若 > 16384 須附理由與安全部書面接受紀錄。(2) 行為：在主機上直連 origin（避免邊緣先攔截）`curl -s -o /dev/null -w '%{http_code}\n' -H "X-Pad: $(printf 'a%.0s' $(seq 32768))" http://127.0.0.1:8080/api/v1/urls/<已建立短碼>/stats`，預期 `431`（若設定值 > 16 KiB，測試長度改為大於「設定值＋4 KiB」）；同一指令改用 `$(seq 1024)` 預期不是 `431`。註：Go 讀取標頭時會在 `MaxHeaderBytes` 之外多留約 4 KiB 緩衝，所以測試值取 32768 B，明顯超過 16 KiB | 研發部（實作）；維運部（部署時驗證） | 待 G5 前驗證 |
| G5-07 HTTPS／TLS | production 對外入口啟用 HTTPS，TLS 1.2 以上；憑證為公開信任或文件化的內部信任。`BASE_URL` 在 production 為 `https://...`。不得以 local 的 HTTP 宣稱滿足 SEC-012 | 01-sec-requirements-v0.2 SEC-012（第 35 行）；02 §1 第 25 行、§2.3 第 64 行、§6 第 118 行 | (1) `curl -sS -o /dev/null -w '%{http_code} %{ssl_verify_result}\n' https://<正式網域>/<已建立短碼>`（不加 `-k`），預期 `302 0`。(2) `openssl s_client -connect <正式網域>:443 -servername <正式網域> -tls1_1 </dev/null` 握手失敗；改 `-tls1_2` 握手成功。(3) `docker exec <容器名> printenv BASE_URL` 以 `https://` 開頭 | 維運部 | 待 G5 前驗證 |
| G5-08 限流預設值 | staging／production **必須**套用 SEC-007 預設（或更嚴）：`RATE_LIMIT_CREATE_PER_MIN`＝30、`RATE_LIMIT_REDIRECT_PER_MIN`＝120（變數名為 02 的建議名，以實作為準）；不得把 local 放寬值帶入 production | 02 §2.1 第 41 行、§2.3 第 65 行、§3 第 70–80 行；11 號檔檢查表 SEC-007 列（第 42 行，非正值回落預設） | `docker exec <容器名> printenv RATE_LIMIT_CREATE_PER_MIN RATE_LIMIT_REDIRECT_PER_MIN`：未設（回落預設）或值 ≤ 30／≤ 120；G5-01 的 31 次測試同時佐證 create 門檻為 30 | 維運部 | 待 G5 前驗證 |
| G5-09 SQLite 不公開、檔案權限 | production 使用獨立 SQLite 檔（不與 local／staging 共用）；檔案只有服務執行身分可讀寫（例如 `0600` 或目錄 `0700`）；資料目錄不對網際網路開放，只暴露應用 HTTP(S) | 01-sec-requirements-v0.2 SEC-017（第 40 行）；02 §2.3 第 62–63 行、§4 第 88–90 行；03 §7 第 5 點（第 119 行） | (1) `stat -c '%a %U:%G %n' <DATABASE_PATH> $(dirname <DATABASE_PATH>)`，預期檔案 `600`（或目錄 `700`）且擁有者為服務帳號。(2) `docker exec <容器名> printenv DATABASE_PATH` 與 staging 的值不同。(3) 主機上 `ss -ltnp` 只有應用（或代理）監聽埠，沒有檔案分享服務 | 維運部 | 待 G5 前驗證 |
| G5-10 不暴露 `.git` | 部署產物與對外路徑不得暴露 `.git` | 01-sec-requirements-v0.2 SEC-015（第 38 行）；02 §2.3 第 66 行 | `curl -s -o /dev/null -w '%{http_code}\n' https://<正式網域>/.git/HEAD`，預期 `404` 或 `403`，不得回傳版控內容 | 維運部 | 待 G5 前驗證 |
| G5-11 機密不進版控 | 機密（Tunnel token、TLS 私鑰、任何金鑰）不進 git，只存在 GitHub Environments secrets 或主機／容器注入；版控只允許 `.env.example` 與明顯假值 | 03 §1 第 22–23 行、§2 第 33–36 行、§3 第 48 行、§7（第 113–121 行）、§8 第 129 行；02 §5 第 104 行；SEC-011（03 引用） | (1) 在發布 commit 執行 `git ls-files -- '*.env' '.env.*' '*.db' '*.pem' '*.key' '*.p12' '*.pfx'`，預期只輸出 `.env.example`。(2) 發布 commit 的 `secrets-gitleaks` check 為綠（附 Actions run 連結）。(3) Tunnel token 的存放位置截圖（Environments secrets 或主機祕密掛載），畫面不得顯示 token 值 | 維運部 | 待 G5 前驗證 |
| G5-12 託管選定 | 零預算託管方式在 G5 前選定，選定後更新 02 與部署手冊，並驗證 SEC-012／SEC-017。本表其他項目依此代入實際值 | 02 §5 第 98–107 行；09-cicd-env README「文件狀態 vs 實作狀態」 | 02 已有選定結果的修訂紀錄；本表佔位符已代入 | 維運部 | 待 G5 前驗證 |

補充說明：

- G5-05 與 G5-06 中的研發端事項（`LOG_LEVEL` 接線、`rate_limited` 改計數或取樣，見 12 號檔 SCR-013 第 87 行；`MaxHeaderBytes` 實作）不屬維運，本清單只驗證部署結果。
- **安全部已確認（2026-10-19 09:03，UTC+8）**：
  1. 邊緣 HSTS **不加** `includeSubDomains`，理由同 11 號檔「HSTS 裁定」第 216 行（demo 網域範圍未定）；日後若改用專屬網域且確認所有子網域皆僅 HTTPS，加入前須先取得安全部書面同意。**不加** `preload`。
  2. max-age：設定剛好 `31536000` 可接受；G5-04 通過條件為以**實際回應標頭**驗證 max-age ≥ 31536000。
  3. `MaxHeaderBytes`：必守為在 `http.Server` **明確設定**（不依賴 1 MiB 預設），超大標頭回 `431`，於 G5 送審前完成；≤ 16 KiB 為建議值，超過 16 KiB 須敘明理由並取得安全部書面接受（10 號檔 C3 第 54 行）。

---

## 2. 證據要求

1. **截圖**：一律保留瀏覽器**網址列**（能看出網域與頁面路徑），不可只截內容區。畫面中若有 token、金鑰或其他機密，必須先遮蔽再存（03 §7 第 2 點）。
2. **API／curl 擷取**：每一筆都記錄三項：
   - 實際執行的完整指令（含參數；機密以佔位符取代）；
   - 執行時間，**UTC+8**，例如先執行 `TZ=Asia/Taipei date '+%Y-%m-%d %H:%M:%S %z'` 並一起貼上；
   - **原始輸出**，不刪減、不改寫（標頭擷取要保留完整回應標頭）。
3. **觀測日誌**：證據中不得貼上未過濾的觀測日誌內容（含短碼清單），只記錄設定與權限的輸出（12 號檔條件 3，第 65 行）。
4. **存放位置**：沿用 `docs/09-cicd-env/evidence/`，建議檔名 `g5-deploy-evidence-<YYYY-MM-DD>.md`；每筆證據註明對應本表的項目編號（G5-01～G5-12）。

---

## 3. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | 初版草稿：G5 部署檢核 12 項（可信代理 CIDR、代理轉傳標頭、邊緣限流／單實例、邊緣 HSTS 與 `HSTS_ENABLED` 不生效、觀測日誌存取／保存／不外送、MaxHeaderBytes、HTTPS／TLS、限流預設、SQLite、`.git`、機密、託管選定）；證據要求；狀態皆「待 G5 前驗證」 |
| v0.1 | 2026-10-19 | 維運部 | 草稿修訂（未入 repo，版本維持 v0.1）：G5-02 要求明寫「代理必須覆寫 X-Real-IP（RR-004）」並補 11 號檔 SCR-008 出處（安全部要求）；補充說明新增「待安全部確認」（邊緣 `includeSubDomains`／不加 `preload`、max-age 剛好 31536000 vs 來源 ≥、MaxHeaderBytes ≤ 16 KiB 屬 C3 建議而必守為明確設定）；來源行號依 05 v0.1.7、11 v0.1.2、12 v0.1.2、02 v0.1.1 重新核對（11 號檔 HSTS 裁定 120–125→213–219、第 123→216 行、RR-007 第 151→244 行；05 第 101→102 行、§2.4／§2.6／§2.7 範圍）；依安全部 2026-10-19 09:03（UTC+8）回覆：G5-04 改為不加 `includeSubDomains`／`preload`、通過條件 max-age ≥ 31536000 以實際回應標頭驗證；G5-06 必守為明確設定＋431，≤ 16 KiB 為建議、超過須理由＋安全部書面接受；「待安全部確認」改為「安全部已確認（2026-10-19）」 |
