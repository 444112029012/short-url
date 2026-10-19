---
文件：CR-002（Go 1.27.2＋gosec 釘未發行 commit）安全部意見——對 D-13 之同意／不同意
版本：v0.1.2
狀態：安全部已同意（附條件），2026-10-19 已回覆總協調；使用者已於 D-13 核准甲（經總協調轉述）；**v0.1.2：CR-002 已生效（CR-002 PR #15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`；PR 審查見 [14](./14-pr-015-cr-002-security-review-v0.1.md)）**
負責角色：安全部
最後更新：2026-10-19
對應：SEC-009（SAST）、SEC-010（SCA）、NFR-004、SI-05；門檻 [03-ci-security-gates-v0.3.md](./03-ci-security-gates-v0.3.md) v0.3.6 §2；上一個 CR：CR-001（[09](./09-cr-ci-deviations-ruling-v0.1.md)）；PR #14（Go 釘回 1.27.1，open）
專案代號：SHORTURL
---

# CR-002 安全部意見（v0.1）

## 1. 結論

**安全部同意（附條件）D-13 方案甲**：Go 精確釘 `1.27.2`（`go.mod` 與 CI `setup-go` 一致，不再用 `1.27.x`），gosec 暫釘上游預設分支 commit `707f2271f66579a18e1da25794b1e8012247274b`（pseudo-version `v2.29.1-0.20261015133444-707f2271f665`），`-severity=high`、`-no-fail=false` 與四項 required checks 不變。

- 查證未發現阻擋級問題：該 commit 在 `securego/gosec` 預設分支 `master` 上（截至查證時即 master HEAD），由維護者合併；與 v2.29.0 之間沒有停用規則、沒有 `replace`，新增依賴只有 3 個間接依賴（隨既有 AI SDK 升版帶入）；五個漏洞在 1.27.2 已修、1.27.1 受影響，本機重現一致。
- 同意以 §5 條件 C1～C8 為前提；其中 **C1～C5 為阻擋**（未滿足前，CR-002 的 CI 變更 PR 不得合併）。
- **PR #14（釘回 1.27.1）不可接受**：會讓 `sca-govulncheck` 持續紅燈（五個可達標準函式庫漏洞），唯一的「綠燈」路徑只剩放寬 SCA 門檻，違反 03 §4、§12（見 §6.2）。
- v0.1 時本文件只新增於 `docs/06-security/`、未改其他檔；**v0.1.1：安全部已於 2026-10-19 回覆總協調「同意（附條件）」，§7 草稿已落實為 03 v0.3.7（待 CR-002 PR 合併後生效），RR-010 已登錄於 05 v0.1.9。** **v0.1.2：CR-002 PR 為 #15，安全部審查見 [14-pr-015-cr-002-security-review-v0.1.md](./14-pr-015-cr-002-security-review-v0.1.md)（核准，C1～C5 符合）；CR-002 PR #15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`，合併後 main push run [39189683566](https://github.com/444112029012/short-url/actions/runs/39189683566) 四項 success、`sast-gosec` 由 `2e54f7e` 的 failure 轉綠且 log 有 `gosec h1 OK`；03 v0.3.8、05 v0.1.10 已改為已生效。C6 期限 2026-10-26（UTC+8）；C7 每週檢查（第一次最遲 2026-10-26），發行後 14 日內換回，最遲 2026-12-31 重評。**

## 2. 背景

| 項目 | 內容 |
|---|---|
| 觸發 | CI `setup-go` `1.27.x` 於 Go 1.27.2 發行（2026-10-08）後解析為 1.27.2；gosec v2.29.0 依賴 `golang.org/x/tools v0.49.0`，無法解碼 Go 1.27.2 的 export data v5，`sast-gosec` 全 repo 紅燈（不是有發現，是載入失敗後 exit 1） |
| 替代嘗試 | PR #14 釘回 1.27.1 → gosec 綠，但 `govulncheck` 擋下 5 個可達漏洞 GO-2026-6603／6611／6612／6613／6617 |
| 決策 | 使用者已核准方案甲（D-13，經總協調轉述），以安全部同意為執行前提 |
| 方案乙（不採用） | gosec v2.29.0＋`replace golang.org/x/tools => v0.50.0` 自行編譯：偏離上游測試組合（上游從未以該組合測試發行），需自維護 fork／replace，誤判或漏判無人背書，可重現性要靠自家建置腳本；甲所用 commit 則已由上游 CI 以 Go 1.27.2 測過（PR #1772 將 CI matrix 改為 1.26.9／1.27.2）。故不採用 |
| 維運查證報告 | canvas 連結 WebFetch 逾時、curl 回 307 轉址（需登入），**本文件未能讀取**，以下全為安全部獨立查證 |

## 3. 查證結果（查證時間 2026-10-19 09:35～09:40 UTC+8）

### 3.1 gosec 上游 issue #1774 與發行狀態

| 項 | 結果 | 來源 |
|---|---|---|
| issue #1774 | 標題「Go 1.27.2 compatibility: v2.29.0 fails to decode export data version 5」；**open**；建立 2026-10-11 19:38 UTC+8；2 則留言：一則指出 master 自 PR #1772 起可用於 Go 1.27.2，一則詢問何時發行；無維護者回覆、無 milestone | https://github.com/securego/gosec/issues/1774 （GitHub REST API） |
| 最新發行版 | **v2.29.0**（2026-08-26 16:34 UTC+8）；其後無新 tag／release（v2.28.0、v2.27.1 更舊） | https://github.com/securego/gosec/releases |
| 結論 | **目前沒有支援 Go 1.27.2 的 gosec 發行版**，與維運結論一致 | — |

### 3.2 commit `707f2271f66579a18e1da25794b1e8012247274b`

| 項 | 結果 |
|---|---|
| 所在分支 | `securego/gosec` 預設分支為 `master`；compare `707f227...master` 回 `identical`（查證時即 master HEAD）→ **在預設分支上，不是 fork 或未合併 PR 的 commit** |
| 內容 | `fix(G203): treat constant template arguments as safe (#1770)`；作者 Luis Figueroa（LuisFigueroaG，外部貢獻者，PR 來自其 fork），2026-10-15 21:34:44 UTC+8；committer GitHub web-flow，簽章 verified；**PR #1770 由維護者 ccojocar 合併**，merge commit 即本 SHA |
| 相容性來源 | PR #1772「Update to go version 1.27.2 and 1.26.9」（ccojocar，2026-10-09 20:08 UTC+8 合併，commit `7b1b5ce`）；本 SHA 之後 |
| pseudo-version | `v2.29.1-0.20261015133444-707f2271f665`；proxy `.info` Origin Hash＝本 SHA |
| sum.golang.org | `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`（go.mod：`h1:sXj2muZ/CL9FWlyhZQZh8ya4oQmSmpfLGm7A8atCh/w=`）；本機下載 ziphash 一致 |
| 與 v2.29.0 差異 | 35 個 commit、45 個檔案（compare `v2.29.0...707f227`：ahead 35、behind 0） |

**差異重點：**

| 類別 | 內容 | 判定 |
|---|---|---|
| go.mod | `golang.org/x/tools v0.49.0 → v0.52.0`（修正 export data v5 的關鍵）；x/crypto 0.58.0、x/net 0.61.0、x/sys 0.49.0 等例行升版；AI SDK（anthropic 1.80.0、openai 3.77.0、genai 1.73.0）升版；`go 1.25.0 → 1.26.0`；**無 `replace`** | 正常（Renovate 例行） |
| 新增間接依賴 | `github.com/coder/websocket v1.8.15`、`github.com/pb33f/go-yaml v0.1.1`、`go.opentelemetry.io/otel/log v1.47.0`；移除 `go.yaml.in/yaml/v4 rc.6` | 皆為已知專案，隨 AI SDK／otel 升版帶入；gosec 預設掃描不呼叫 AI（autofix 需另設金鑰）。非阻擋，列 C7 觀察 |
| 規則行為「收窄」（減少誤判） | G203：編譯期常數視為安全（本 commit）；G404：`Shuffle` 不再報；G702：`exec.CommandContext` 只檢查 name／argv 參數；`#nosec` 抑制改為與註解順序無關、巢狀範圍各自保留 | 非停用；對本專案影響低（未用 html/template、os/exec）。`#nosec` 語意變更已由 C6 控管 |
| 規則行為「擴大」 | G102 偵測 `[::]`；G402 追蹤經 `*tls.Config` 指派；G403／G124／檔案權限規則解析常數；G407、G602 強化；G121 追蹤 closure | 偵測力提升 |
| 其他 | SARIF 規則排序、HTML 報表 SRI 修正、上游 CI matrix 改 1.26.9／1.27.2 | 無安全疑慮 |
| 停用規則 | **未發現**任何規則被移除或預設排除 | — |

### 3.3 五個標準函式庫漏洞（vuln.go.dev／pkg.go.dev/vuln；皆於 2026-10-09 06:31 UTC+8 發布）

| ID（CVE） | 摘要 | 受影響（stdlib） | 修正 |
|---|---|---|---|
| [GO-2026-6603](https://pkg.go.dev/vuln/GO-2026-6603)（CVE-2026-78659） | HTTP/2 伺服器以 Trailer 宣告大量欄位，繞過 MaxHeaderBytes 等限制造成記憶體耗盡；不適用 HTTP/1 | < 1.26.9；1.27.0～<1.27.2 | 1.26.9、**1.27.2**（x/net 0.60.0） |
| [GO-2026-6611](https://pkg.go.dev/vuln/GO-2026-6611)（CVE-2026-78669） | HTTP/2 對端大量 SETTINGS_INITIAL_WINDOW_SIZE 造成 CPU 耗盡（client／server） | 同上 | 同上 |
| [GO-2026-6612](https://pkg.go.dev/vuln/GO-2026-6612)（CVE-2026-78663） | HTTP/2 伺服器連線層流量控制重複退還，可繞過 MaxReceiveBufferPerConnection | 同上 | 同上 |
| [GO-2026-6613](https://pkg.go.dev/vuln/GO-2026-6613)（CVE-2026-94439） | HTTP/1 伺服器對 CONNECT 回 2xx 但未 hijack 後仍繼續解析，可能與前置代理不同步（request smuggling） | 同上 | 1.26.9、**1.27.2** |
| [GO-2026-6617](https://pkg.go.dev/vuln/GO-2026-6617)（CVE-2026-97032） | HTTP/2 伺服器 HPACK encoder 競態，惡意 client 可使伺服器崩潰 | 同上 | 同上 |

- Go 官方發行紀錄：go1.27.2 於 2026-10-08 發行，含 `net/http` 等安全修正（https://go.dev/doc/devel/release#go1.27.2）。
- **確認：1.27.1 受影響、1.27.2 已修**（本機重現見 §3.4）。
- 對本專案影響：`cmd/shorturl/main.go` 以 `srv.ListenAndServe()` 提供明文 HTTP/1（TLS 在邊緣終止，見 RR-008），未設定 h2c／`Server.Protocols`，四個 HTTP/2 漏洞在現行設定下**實際觸發面低**，但 govulncheck 以符號可達判定（`http.Server.ListenAndServe` 等），且日後一旦啟用 HTTP/2 即直接暴露；GO-2026-6613 屬 HTTP/1，前置代理（Cloudflare Tunnel／反向代理）情境下有 smuggling 面，雖 chi 對未註冊 CONNECT 回 405、影響有限，仍屬真實暴露。有修補版本可用，依 03 §4／§5 應升級而非例外。

### 3.4 本機重現（唯讀 clone，`/workspace/tmp/cr002/`，不提交、不推送）

基準：`main` @ `2e54f7e2d04c891886a4e4210f5880475c9f0d18`（Merge PR #13，2026-10-19 09:05:26 UTC+8）；`go.mod` `go 1.27.1`。工具以 `go install`（GOPROXY＝proxy.golang.org、GOSUMDB＝sum.golang.org）安裝於暫存 GOPATH。

| # | 組合 | 結果 |
|---|---|---|
| R1 | Go 1.27.2＋gosec v2.29.0（`-fmt=text -no-fail=false -severity=high ./...`） | **exit 1**；`could not import context … export data version 5 is greater than maximum supported version 4`；Issues 0（重現 #1774） |
| R2 | Go 1.27.2＋gosec @`707f227`（同參數） | **exit 0**；Files 17、Lines 1206、Nosec 0、Issues 0；無 export data 錯誤（與維運一致） |
| R3 | Go 1.27.2＋govulncheck v1.8.0 | **exit 0**，`No vulnerabilities found.` |
| R4 | Go 1.27.1（GOTOOLCHAIN=go1.27.1）＋govulncheck v1.8.0 | **exit 3**；GO-2026-6617／6613／6612／6611／6603，「Your code is affected by 5 vulnerabilities from the Go standard library」（重現 PR #14 結果） |
| R5 | 金絲雀：repo 外暫存模組（含常數假密碼與 `crypto/md5`）＋gosec @`707f227` | `-severity=high`：**exit 1**，G101（HIGH）被抓；不加 severity：G101、G401、G501 共 3 筆。證明新 gosec 確實在掃且門檻會擋 |

- gosec 二進位 `go version -m`：`mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HN…`、`x/tools v0.52.0`、`go1.27.2`。注意 `go install` 版 gosec 的 Summary 顯示 `Gosec : dev`，版本須以 `go version -m` 確認（C3）。
- 旁註（非本 CR 範圍）：現行 workflow `actions/setup-go@v5` 為 tag 釘選、非完整 SHA（03 §10 建議 SHA）；建議另案處理。

## 4. 方案比較（摘要）

| 面向 | 甲（D-13 已核准） | 乙（不採用） | PR #14（1.27.1） |
|---|---|---|---|
| 供應鏈 | 未發行 commit，無 release 公告／簽章；但 SHA 不可變、上游維護者合併、模組 h1 經 sum.golang.org 透明日誌驗證 | 自編 replace 組合，sumdb 只能驗證各模組，組合本身無人背書 | 發行版 |
| 上游測試覆蓋 | 上游 CI 以 Go 1.27.2 測過 | 無 | 有 |
| 可重現性 | `go install …@<pseudo-version>` 完全決定 | 依自家 go.mod／腳本 | 高 |
| 維護成本 | 低；上游發行後一次 CR 換回 | 高（維護 fork／replace） | 低 |
| 安全結果 | gosec＋govulncheck 皆綠 | 未驗證 | **5 個可達漏洞未修**，只能放寬 SCA 才綠 |
| 判定 | **同意（附條件）** | 不採用 | **不可接受** |

## 5. 同意條件

| # | 條件 | 阻擋 | 負責 | 驗證方式 |
|---|---|---|---|---|
| **C1** | Go 精確釘 **1.27.2**：`go.mod` `go 1.27.2`（或 `go 1.27.1`＋`toolchain go1.27.2`，擇一，須與 CI 一致）；三個 job 的 `setup-go` `go-version: "1.27.2"`，**禁止 `1.27.x`／`stable`** | 是 | 研發部／維運部 | PR diff；CI log `go version go1.27.2` |
| **C2** | gosec 以**完整 pseudo-version 或完整 40 碼 SHA** 安裝：`go install github.com/securego/gosec/v2/cmd/gosec@v2.29.1-0.20261015133444-707f2271f665`（或 `@707f2271f66579a18e1da25794b1e8012247274b`），禁止 `@master`／`@latest`／短 SHA；不得設 `GOFLAGS=-insecure`、`GONOSUMDB`／`GONOSUMCHECK`／`GOSUMDB=off`、`GOPROXY=direct` | 是 | 維運部 | PR diff；CI env 檢視 |
| **C3** | CI 在 gosec 安裝後輸出 `go version -m "$(go env GOPATH)/bin/gosec"`，並比對模組 h1 等於 `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`（不符即 `exit 1`）；log 保留 | 是 | 維運部 | CI log 有該行且比對步驟成功 |
| **C4** | 門檻不變：`gosec -no-fail=false -severity=high ./...`；`govulncheck` v1.8.0 不變；四項 required checks 名稱、ruleset 24629805（strict=true、無 bypass）不變；**本期間禁止放寬 sast 門檻、禁止把 `sast-gosec` 改為非必要或 `continue-on-error`、禁止新增 `-exclude`／`#nosec`** | 是 | 維運部／安全部 | PR diff；合併後以 API 唯讀讀 ruleset |
| **C5** | 合併前四項 required checks 在 PR 最新 head 全綠（D-10），並依 D-08／D-09 留言程序合併；PR 標 `CR-002`；PR #14 關閉不合併 | 是 | 研發部／總協調 | run 連結；PR 狀態 |
| C6 | 金絲雀驗證：已由安全部本機完成（§3.4 R5，不提交）；CI 變更合併後，維運於一個不合併的測試分支或 `workflow_dispatch` 以等價假樣本再驗一次 `sast-gosec` 會紅，之後刪除分支（對應 DevSecOps skill 步驟 4）。另：`#nosec` 語意已變，日後新增 `#nosec` 一律須安全部核准 | 否（合併後 7 日內） | 維運部 | run 連結（紅燈）＋分支已刪 |
| C7 | 退場監控：每週（或下次 CR，取早者）檢查 gosec releases 與 #1774；上游發行含 x/tools ≥ v0.50.0 且支援 Go 1.27.2 的版本後 **14 日內**開 CR 換回發行版；**最遲 2026-12-31** 仍未發行則由安全部重評（延長須再經安全部同意並記 RR-010）。順帶觀察新增間接依賴（coder/websocket 等）有無公告漏洞 | 否 | 維運部（執行）／安全部（重評） | 週檢紀錄；CR 編號 |
| C8 | Go 後續修補版（1.27.3…）同樣以精確釘版＋CR（標準變更）更新，先在 PR 驗證 gosec／govulncheck 兩者皆綠才換，避免再次自動漂移 | 否 | 維運部 | PR／CR 紀錄 |

## 6. 殘餘風險建議

### 6.1 建議新增 RR-010（低）

| 欄位 | 建議內容 |
|---|---|
| 編號 | RR-010（關聯 SEC-009、NFR-004、SI-05） |
| 殘餘描述 | SAST 工具使用上游未發行 commit（gosec `707f227`），無發行公告／簽章；相對 v2.29.0 有規則行為變更（部分收窄）與依賴升級；工具版本偏離 03 原核准發行版 |
| 緩解 | 上游預設分支、維護者合併、SHA＋sum.golang.org h1 鎖定（C2／C3）；門檻不變（C4）；金絲雀驗證（C6）；退場期限（C7） |
| 建議等級 | **低**（理由：完整性由 SHA 與透明日誌保證；上游 CI 已以 Go 1.27.2 測過；本機實測可擋 HIGH；只影響 CI 工具、不進產品二進位；有期限與退場條件） |
| 接受者 | 安全部（低／中可由安全部接受）；使用者已於 D-13 知悉 |
| 重新評估 | 上游發行時或 2026-12-31，取早者 |
| 結案條件 | 換回支援 Go 1.27.2（或當時版本）的 gosec 發行版並經 CR 合併 |

### 6.2 PR #14（Go 1.27.1）不可接受的理由

- 五個漏洞皆被 govulncheck 判為可達、皆有修正版 1.27.2；03 §4／§5 要求有修補即升級（High 14 日內）。
- 停在 1.27.1 只能讓 `sca-govulncheck` 長期紅燈，或放寬 SCA 門檻／加例外——後者違反 03 §4「不得無故放寬」與 §12 禁止事項，等同用削弱 SCA 換 SAST，安全淨值為負。
- 若採 #14 而接受五個漏洞，殘餘（含 GO-2026-6613 HTTP/1 smuggling、HTTP/2 DoS 啟用即暴露）需另立風險並由使用者書面接受；甲可完全避免，故不建議。

## 7. 03 §2 改版草稿（v0.1.1：已落實為 03 v0.3.7，以 03 正文為準）

**版號建議：v0.3.7**（小版、不改檔名）。理由：屬 CR-002 核准下的工具鏈版本調整，阻擋語意、job 名、required checks 均未變；依 change-management 為「一般變更」，非技術選型或門檻語意變更，不需升 v0.4。檔頭「對應 CR／決策」補「CR-002（D-13）」。

§2 表格替換列（草稿）：

| 掃描 | 對應 SEC | 本專案選定（主選） | 版本鎖定（v0.3.7） | 備援／併用（可選） | 授權／成本 |
|---|---|---|---|---|---|
| **SAST** | SEC-009 | **gosec**（Go） | **暫行例外（CR-002／D-13）**：`github.com/securego/gosec/v2@v2.29.1-0.20261015133444-707f2271f665`（上游 master commit `707f2271f66579a18e1da25794b1e8012247274b`；sum.golang.org `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`）；CI 須以 `go version -m` 比對 h1；禁止 `@latest`／`@master`／短 SHA。**例外期限**：上游發行支援 Go 1.27.2 之版本後 14 日內換回，最遲 2026-12-31 重評。**退場條件**：CR 換回精確釘版之發行版（`vX.Y.Z`）並合併後，本例外與 RR-010 結案 | 可選 **CodeQL**（language: `go`） | 開源 |
| **Go toolchain** | （供應鏈） | 官方 Go | **`go.mod` `go 1.27.2`；CI `setup-go` `"1.27.2"`（精確釘版，禁止 `1.27.x`／`stable`；須一致）**。後續修補版以 CR（標準變更）更新，且須先驗證 gosec／govulncheck 皆綠 | — | 官方 |

§2 註記（草稿新增）：
> CR-002（D-13，2026-10-19）：Go 1.27.2 修正 GO-2026-6603／6611／6612／6613／6617；gosec v2.29.0 無法解碼 Go 1.27.2 export data v5（上游 issue #1774，截至 2026-10-19 無發行版），故暫釘上游預設分支 commit。門檻（§4）、job 名（§3）、`-severity=high`／`-no-fail=false` 均未放寬；條件見 [13](./13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md)。

連帶：檔頭摘要「工具鏈 Go 1.27.x」改「Go 1.27.2（精確釘版）」；§9 SEC-009 列 `gosec v2.29.0` 改「gosec（CR-002 暫行 pseudo-version）」；§12 新增「禁止在 CR-002 例外期間放寬 SAST 門檻或將 `sast-gosec` 改為非必要檢查」；修訂紀錄加 v0.3.7。

## 8. 同步狀態（v0.1.1）

- 05 殘餘風險表 → v0.1.9：已新增 RR-010（低）。
- 03 → v0.3.7：已落實 §7（待 CR-002 PR 合併後生效）。
- README → v0.14：已同步。
- `docs/99-changelog/`：待總協調登記 CR-002／D-13（非安全部檔案，未改）。
- 待辦：CR-002 PR（研發部／維運部）依 C1～C5 實作；PR #14 關閉。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-19 | 安全部 | 初版：對 D-13（方案甲）同意（附條件）；查證 #1774（open、無發行版）、commit `707f227` 在 master 且由維護者合併、v2.29.0 差異（無停用規則、無 replace、3 個新間接依賴）、五漏洞 1.27.1 受影響／1.27.2 已修；本機重現 R1～R5；條件 C1～C8（C1～C5 阻擋）；建議 RR-010（低）；PR #14 不可接受；03 §2 v0.3.7 草稿 |
| v0.1.1 | 2026-10-19 | 安全部 | 狀態改為「安全部已同意（附條件），2026-10-19 已回覆總協調」；§1、§7、§8 註記草稿已落實為 03 v0.3.7、RR-010 已登錄於 05 v0.1.9、README v0.14 已同步；結論與條件 C1～C8 不變；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-19 | 安全部 | 註明 CR-002 已生效（CR-002 PR #15 於 2026-10-19 10:22:50（UTC+8）合併，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`），引用 14 號檔（PR #15 審查，核准）；記錄合併後 main run 39189683566 四項 success、sast-gosec 轉綠、h1 OK；寫入 C6（2026-10-26）、C7（每週、發行後 14 日內、最遲 2026-12-31）具體期限；結論與條件 C1～C8 不變；小版修訂不改檔名（檔名保留主版號） |
