---
文件：變更請求 CR-005
版本：v0.3（草稿，依審查部 2026-10-27 退回紀錄 CR5-R1～R4、S1～S7 修訂；技術內容與 v0.2 相同）
日期：2026-10-27
狀態：待複審（維運部：已提出｜企劃部：無影響｜設計部：無影響｜研發部：無影響｜品保部：無影響，附驗收條件 QA-CR005-1～6｜安全部：同意（附條件，30 號 v0.1.1）｜審查部：v0.2 退回，v0.3 待複審｜使用者：待拍板）
提出：維運部
專案代號：SHORTURL
關聯：CR-002；06-security/13（C7、C2、C3、C4、C8）；06-security/30 v0.1.1（sha256 e4579d31ab6ae9ba883b5a2e84f21d2da22b2b456f658127bca15a893a598104）；RR-010；SEC-009；上游 securego/gosec#1774；審查紀錄 /workspace/review-CR-005-2026-10-27.md（sha256 08adc77a274622589d2ca7bcc75a0d9cb69e3b34648efff7959558bb920c8cb2）；C7 週檢 v0.2（sha256 801a5ea75bb99b71673a579152d1fbef2ec970a1ef51eb021d34371f4bc2e01c）
時區：Asia/Taipei（UTC+8）；文中 Z 結尾時間為 UTC，已並列換算
---

# CR-005｜sast-gosec 由 CR-002 暫行 pseudo-version 換回發行版 gosec v2.30.0

> **v0.3 說明：技術內容與 v0.2 相同。** patch（sha256 `7ddfa0f8bfd1e9465fdb85ed8f40c5da25b1d1645211519090eb302806df1853`）、`ci-devsecops.proposed.yml`（sha256 `a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f`）、§1 變更內容、§5 驗證、§7 回滾方案均未修改（§2、§3、§6 亦未修改）。本版只補影響分析、條件原文、驗收條件、執行程序與紀錄類內容，逐處改動見 §10。
>
> **Go 1.27.2 是否可用，以 PR 的 CI 結果為準；本機實測僅為參考。**（審查部 S5）

- 提出者／日期：維運部／2026-10-27
- 類型：一般變更（CI 資安控制工具版本；不改門檻、不改 job 名稱）
- 影響的需求編號：SEC-009（門檻不變，安全部判定不受影響）；REQ：無；NFR：無（企劃部判定）；關聯驗收項 CHK-SEC-009、CHK-NFR-004、T-CI-01（品保部，判定內容不變、證據須於合併後重取）
- 期望完成時間：**最遲 2026-11-06（UTC+8）**（C7：上游發行後 14 日內；v2.30.0 發行於 2026-10-23T09:22:19Z＝2026-10-23 17:22:19（UTC+8））
- 附件：擬議 patch `CR-005-ci-devsecops-gosec-v2.30.0.patch`；證據目錄 `docs/09-cicd-env/evidence/cr-005/`（§8）

## 1. 變更內容（現況 → 期望）

只改 `.github/workflows/ci-devsecops.yml` 的 `sast-gosec` job 兩行（外加檔頭一行註解同步），其他一律不動。

| 項目 | 現況（CR-002，main 2425129） | 期望（CR-005） |
|---|---|---|
| `Install gosec` | `go install github.com/securego/gosec/v2/cmd/gosec@v2.29.1-0.20261015133444-707f2271f665` | `go install github.com/securego/gosec/v2/cmd/gosec@v2.30.0` |
| `Verify gosec module hash` 的 `expected` | `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=` | `h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=` |
| 檔頭第 8 行註解 | 「gosec 依 CR-002 釘 707f227…」 | 「gosec 依 CR-005 換回發行版 v2.30.0（取代 CR-002 暫行釘的 707f227）…」（僅註解） |
| Go 1.27.2（三處 setup-go、check-latest: false）、govulncheck v1.8.0、gitleaks 8.21.2（含 sha256）、h1 比對邏輯、`-fmt=text -no-fail=false -severity=high -tests ./...`、四個 required check 名稱（secrets-gitleaks／sast-gosec／sca-govulncheck／unit-test）、ruleset、dependency-review fail-closed | — | **不變** |

說明：步驟名稱 `Verify gosec module hash (CR-002)` 刻意不改，避免夾帶非必要變更；若審查部要求改為 `(CR-005)`，屬純名稱變更，不影響 required check（required check 為 job 名稱）。

擬議 patch（只存在 `/workspace/rebuild/cr-005/`，未套用到任何 repo）：

```diff
-        run: go install github.com/securego/gosec/v2/cmd/gosec@v2.29.1-0.20261015133444-707f2271f665
+        run: go install github.com/securego/gosec/v2/cmd/gosec@v2.30.0
...
-          expected='h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40='
+          expected='h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0='
```

`patch -p1 --dry-run` 對 main 2425129 版 workflow 可乾淨套用；PyYAML 解析前後兩版皆成功，jobs 仍為四個，結構比對只有 `sast-gosec` 第 2、3 步（Install gosec、Verify gosec module hash）不同（`patch-validation.txt`）。

## 2. 原因與效益

- **原因（C7）**：安全部 06-security/13 C7：上游發行含 golang.org/x/tools ≥ v0.50.0 且支援 Go 1.27.2 的版本後 14 日內開 CR 換回發行版。
  - gosec v2.30.0 發行：2026-10-23T09:22:19Z（gh release list）；tag 物件 `2385f2c` → commit `5eb1d104ef85ff2c905c877eb600e8c1dee33a22`；proxy.golang.org 的 Origin.Hash 一致。
  - v2.30.0 `go.mod`：`golang.org/x/tools v0.52.0`（≥ v0.50.0）；`go 1.26.0`。
  - #1774「Go 1.27.2 compatibility: v2.29.0 fails to decode export data version 5」：CLOSED／COMPLETED，2026-10-23T10:24:55Z（＝2026-10-23 18:24:55（UTC+8））。
  - 本機實測 v2.30.0 在 Go 1.27.2 下可安裝、可完整掃描（§5）。
- **效益**：結束 RR-010（CI 使用未發行 commit、無發行公告）；回到有 tag 的發行版，之後可依發行版追蹤；含 707f227 之後 19 個上游 commit 的修正（§4）。

## 3. 不變更的後果

逾 2026-11-06 即違反 C7 條件；RR-010 持續；CI 掃描工具停在未發行 commit。

## 4. 影響分析

### 4.1 各部門影響（各列照各部門影響評估原文；維運列由維運填）

各部門評估檔位於 `/workspace/rebuild/cr-005/impact/`，維運部於 2026-10-27（UTC+8）以 `sha256sum` 重算，與總協調轉達的前綴一致：

| 部門 | 評估檔 | sha256 | 評估綁定的 CR 版本 | 結論（原文） | 具名／日期 |
|---|---|---|---|---|---|
| 企劃部 | `impact/企劃部.md` | 28b207d727865b51c2ed8bd94000af15f4ba51bf5fcc696e889a3b97524ffb4a | v0.2（5c6b6f56…25fd） | 「結論：無影響」 | 企劃部（2026-10-27） |
| 設計部 | `impact/設計部.md` | fa5a0b13689344a5fff2279f9def9bb898d30bce32f2b52604bc66d17bab0967 | v0.2（5c6b6f56…25fd）＋patch（7ddfa0f8…1853） | 「**無影響**」 | 設計部 2026-10-27（UTC+8） |
| 研發部 | `impact/研發部.md`（證據 `impact/rd-evidence/` 6 檔，維運重算與該檔所列 sha256 全部相符） | b81b8d193ced8c35edd73860dab8cee260f3b1ba1431a8eaccef217ba52ac8a8 | v0.2（5c6b6f56…25fd） | 「結論：無影響」 | 研發部（2026-10-27） |
| 品保部 | `impact/品保部.md` | 2244c51df983a7298ec7fe13f0b6c18f9e70924ebd44114d2dd0085f02b489b3 | v0.2（5c6b6f56…25fd） | 「結論：無影響（測試策略、追溯矩陣、DoD 均不需改版）；附驗收條件 QA-CR005-1～6」 | 品保部（2026-10-27） |
| 安全部 | `impact/安全部.md`（依據 30 號 v0.1.1，sha256 e4579d31…8598104） | d8ea948d9f404748ce4ff087f57a2d870596caf3962783342b15dd5336c09d48 | v0.1（74f5d76c…df21） | 「結論：**同意（附條件 K1～K6）**，詳見 30 號 §3」 | 安全部（2026-10-27） |

各評估檔均註明「CR-005 再修改時須重新確認」或同義條款；v0.3 未改技術內容，但 sha256 已變，**須各部門對 v0.3 重新確認**（見 §9.5 待決）。

| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍 | 無影響：只把 CI 的 gosec 從 pseudo-version 707f227 換回發行版 v2.30.0，不新增或刪除功能，SI／IN／OUT 都不變 | 無（03-scope、00-charter 不必改版） | 0 | 企劃 |
| 時程 | 無影響：里程碑與 G4／G5 日期不變。CR 的完成期限 2026-11-06（UTC+8）是 C7 條件「上游發行後 14 日內」的要求，不是新的里程碑；維運估 0.5 人日，不占關鍵路徑 | 無（02-preliminary-schedule 不必改版） | 0 | 企劃 |
| 需求 | 無影響：REQ／NFR 不變。SEC-009 的門檻不變（CR 第 15 行「不改門檻、不改 job 名稱」），是否完全不受影響以安全部的確認為準 | 無 | 0 | 企劃 |
| 架構／資料／API／UI | **無影響**（設計部原文要點：只改 `.github/workflows/ci-devsecops.yml` 的 3 行；C4、ADR-001A／002A／003A、類別圖、循序圖、方法規格、資料模型、DFD 與 TB-01～04、OpenAPI 皆不變；UI 原本 N/A；docs/02～05 無工具版本需同步） | 無 | 0 | 設計 |
| 實作 | 無影響。main 2425129 上，研發部用 gosec v2.30.0 以 CI 參數、不加 -severity、只開變動規則三種方式掃描，結果都是 0 發現、Nosec 0、exit 0；程式碼沒有 #nosec；G125 在現有程式碼沒有可觸發的呼叫 | 不適用 | 0 | 研發 |
| 資安（威脅模型、SEC） | 沒有新的資料流或信任邊界，威脅模型不受影響。SEC-009 門檻不變（`-severity=high -tests`、絕對路徑、h1 驗證）。上游 19 個 commit 的規則收窄只限於編譯期常數或已確認設有非零 timeout 的情況，可以接受；G303、G102、G112（timeout 設為 0）屬於擴大偵測。gosec 沒有新增依賴，只有既有依賴升版。GO-2026-6809 新舊版都有，依 C8 另開 CR。結論：**同意（附條件 K1～K6）**，詳見 30 號 §3 | 03 號 CI 門檻文件的工具版本欄；05 號 RR-010 結案，並登錄 GO-2026-6809；13 號補註 C7 由 CR-005 達成；安全部對 PR head 出具書面確認（K2）；合併後金絲雀確認（K4） | 安全部約 0.5 人日（K2 head 確認、K4 覆核、03／05／13 改版） | 安全部 |
| 環境／部署／監控 | 只改 CI `sast-gosec` 安裝版本與 h1 期望值；Go／govulncheck／gitleaks 版本、門檻、job 名、ruleset 不變；不影響執行環境與部署。本機以 Go 1.27.2 官方 tarball 實測：安裝成功、h1 與 sum.golang.org 一致、main 掃描 0 發現 exit 0、G404 金絲雀 exit 1；以擬議 workflow 原文在全新 GOPATH 模擬三步皆通過（僅供參考，以 PR CI 為準） | `.github/workflows/ci-devsecops.yml`（2 行＋1 行註解，CR PR 內）；`docs/09-cicd-env/01-ci-pipeline` 工具版本表（合併後另以文件變更處理，見 §4.3）；合併後依 K4 做 G404 金絲雀 | 約 0.5 人日（PR、CI 觀察、金絲雀、文件） | 維運 |
| 測試策略（06 v0.1） | 無影響：策略只寫「CI 資安門檻檢查表」，未指定 gosec 版本；門檻、層級、工具類別不變 | 無 | 0 | 品保 |
| 追溯矩陣（02 v0.5） | 無影響：全檔**無任何列引用 707f227**、pseudo-version 或 gosec 版本（`rg -n '707f227\|v2\.29\|pseudo' docs/07-testing` 0 筆）。SEC-009 列只對應 CHK-SEC-009，工具版本不在矩陣內 | 不需改；建議下次改版時於 §8 流程類加一列 CR-005（非必要，見 §4） | 0 | 品保 |
| DoD（08 v0.1） | 無影響：S3「PR 觸發 SAST／SCA；High/Critical 依門檻處理」條文與門檻不變 | 無 | 0 | 品保 |
| 由 SAST 對應之驗收項 | CHK-SEC-009（CI SAST job 存在且 High/Critical 擋合併）、CHK-NFR-004、09 號 T-CI-01：判定內容不變，但**證據須在合併後重取**（見 QA-CR005） | 無文件改版；CHK-SEC-009 G5 舉證時改引用 CR-005 後的 run | 約 0.25 人日（驗收比對） | 品保 |
| 程式行為／單元測試 | 無：不改 Go 程式與測試；unit-test 結果不應變 | 無 | 0 | 品保 |

註：企劃、研發、安全、品保各列為原文逐字貼入；設計部評估檔未提供表格列，該列「影響描述」為其「結論」與「核對內容」第 1～5 點的摘要，結論「無影響」與「需更新項目：無。工作量 0。」為原文。品保列中「見 §4」「見 QA-CR005」指品保部評估檔的章節，QA-CR005-1～6 收錄於本 CR §9.3。

### 4.2 影響分析清單（逐項判定）

| 清單項 | 判定 | 判定者（依據） | 要更新的檔案與章節 |
|---|---|---|---|
| 需求規格（SEC-009） | **不受影響**。門檻、阻擋條件和 job 名稱都不變 | 安全部（impact/安全部.md §4.2）；企劃部（REQ／NFR 無影響） | 無 |
| 範圍內外、時程、里程碑 | 無影響 | 企劃部 | 無 |
| 系統架構（C4）與 ADR | 無影響 | 設計部 | 無 |
| 資料模型與資料流程圖 | 無影響 | 設計部 | 無 |
| API 規格（OpenAPI） | 無影響 | 設計部 | 無 |
| UI/UX | 無影響（N/A） | 設計部 | 無 |
| 威脅模型 | **不受影響**。沒有新的資料流，只是工具規則集有變動，已在 30 號 §2.1 評估為可接受 | 安全部 | 無（04 號威脅模型不改版） |
| 資安需求／03 號 CI 門檻文件 | **要更新**。門檻文字不變，只改工具版本欄 | 安全部 | 03 號的 sast-gosec 工具版本列由 707f227 改為 `v2.30.0`（h1 `67aF6LHa…`），合併後由安全部改版 |
| 殘餘風險登錄（05 號） | **要更新** | 安全部 | 05 號 RR-010：K2～K4 完成後改為結案，註明「由 CR-005 換回發行版」。另新增一筆「gosec 工具二進位含 GO-2026-6809（低），依 C8 升 Go 1.27.3」，最遲 2026-11-20 提出 CR（K6）。合併後由安全部改版 |
| 13 號 CR-002 安全意見 | **要補註** | 安全部 | 13 號 C7 列補註「已由 CR-005 達成」，合併後由安全部改版 |
| 隱私合規 | 無影響（不涉個資、蒐集目的、保存期限） | 維運初判；無部門提出異議 | 無 |
| 測試策略、測試案例與追溯矩陣 | 無影響；驗收條件 QA-CR005-1～6 | 品保部 | 無文件改版（矩陣列為建議，非必要） |
| CI/CD、IaC、環境設定、機密 | **要更新** | 維運部 | `.github/workflows/ci-devsecops.yml` sast-gosec 3 行（CR PR）；`docs/09-cicd-env/01-ci-pipeline-v0.2.md` §2.1 |
| 部署與回滾手冊、資料庫遷移可逆性 | 不影響；本 CR §7 即回滾；無資料庫遷移 | 維運部 | 無 |
| 監控告警、備份、事故 runbook | 無影響 | 維運部 | 無 |
| 使用者文件與版本說明 | 無影響 | 維運部 | 無 |
| 變更紀錄（99-changelog） | **要登記** | 總協調 | 見 §4.3；維運不改該目錄 |

上表安全部判定的五項（需求規格、威脅模型、03 號、05 號、13 號）判定與更新內容照 `impact/安全部.md` §4.2 原文。

上游 707f227 → v2.30.0 之間 19 個 commit（`gh api repos/securego/gosec/compare/707f227...5eb1d10`，`status=ahead`, `ahead_by=19`, `behind_by=0`，證據 16）。含規則行為變動，**供安全部評估是否有收窄**：
- 新增：G125（url.JoinPath 未跳脫元素）、G706 擴充 log/slog 方法覆蓋。
- 可能減少誤報（收窄）：G112「honor read timeouts assigned after the http.Server literal」「assigned through a struct field」、G118「cancel sent on a channel as handed off」、G303「resolve constant temp file paths」、G102「resolve concatenated and constant listen addresses」、resolve「constants from other packages as resolved」、G602 兩項。
- 可能增加：G112「report read timeouts explicitly set to zero」。
- 其他：依賴更新、`-concurrency` 預設為 GOMAXPROCS、autofix（AI）相關修正（CI 未啟用 autofix）、上游自身 CI 改用 Go 1.27.3／1.26.10（#1790）。

本專案 main 2425129 上兩版結果相同（皆 0 發現），故上述差異在目前程式碼上沒有可觀察的結果差異；未來程式碼的判定可能不同。

安全部評估結果（30 號 v0.1.1 §2.1）：**可接受**；收窄只限於編譯期常數或已確認有效的設定；G303、G102 與 G112（timeout 為 0）屬擴大偵測；不需新增金絲雀樣本，沿用 G404 金絲雀（K4）。

### 4.3 文件更新分工（誰更新、何時更新）

依 K1，CR PR 只能改 `ci-devsecops.yml` 一個檔，下列文件一律**不放進 CR PR**，於各自時點另行改版。

| 檔案 | 章節／要改的內容 | 負責 | 時點 | 依據 |
|---|---|---|---|---|
| `docs/06-security/03-ci-security-gates-v0.3.md` | §2 SAST 列：工具版本由 707f227 暫行例外改為 `v2.30.0`（h1 `67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`），門檻文字不變 | 安全部 | CR-005 合併後 | 30 號 §4；impact/安全部.md §4.2 |
| `docs/06-security/05-residual-risk-register-v0.1.md` | RR-010（§1 一覽表、§2.8）改結案，註明「由 CR-005 換回發行版」 | 安全部 | K2～K4 完成後 | 30 號 §4；impact/安全部.md §4.2 |
| 同上 | 新增一筆「gosec 工具二進位含 GO-2026-6809（低），依 C8 升 Go 1.27.3」（RR-010 後繼項） | 安全部 | 合併後由安全部改版（impact/安全部.md）；K6 原文：「最遲 2026-11-20 提出；在那之前此項記入 05 號檔」 | 30 號 §2.3、§3 K6 |
| `docs/06-security/13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md` | §5 C7 列補註「已由 CR-005 達成」 | 安全部 | CR-005 合併後 | 30 號 §4 |
| `docs/09-cicd-env/01-ci-pipeline-v0.2.md` | §2.1 現行版本鎖定表：CR-002 列改「已由 CR-005 取代」，新增 CR-005 列（`v2.30.0`、h1 `67aF6…`） | 維運部 | CR-005 合併、main 四綠（K3）後，以文件變更提出（推送同樣須使用者核准） | 本 CR §4.1 維運列 |
| `docs/99-changelog/CR-002-go-1272-gosec-pin-v0.1.md` | 加註「C7 已由 CR-005 達成；gosec 已換回 v2.30.0」 | **總協調** | CR-005 合併後 | 審查部 CR5-R1、S4 |
| `docs/99-changelog/`（變更紀錄表、decision-log） | 登記 CR-005（狀態依 change-management：已提出 → 分析中 → 待核准 → …） | **總協調** | 即時（目前未登記） | 審查部 S4；change-management 一般流程 1 |

**99-changelog 底下所有登記與加註由總協調負責；維運部不改該目錄。**（審查部 S4）。CR-005 本文、30 號、`evidence/cr-005/`、`evidence/c7-weekly/` 目前都只在本機，入庫方式（隨文件 PR）由總協調安排。

## 5. 驗證（本機已執行；全部輸出見 §8）

環境：box（linux/amd64）；Go 1.27.2 官方 tarball `go1.27.2.linux-amd64.tar.gz`，sha256 `ecbadb99091a3f46e31f5f934b068b1864eafa7995211b39eaddf76996045fe5`（與 go.dev/dl JSON 官方值一致，`sha256sum -c` OK），解壓於 `/workspace/tools/go1.27.2`；`GOTOOLCHAIN=local`、`GOFLAGS=` 空、`GOPROXY=https://proxy.golang.org,direct`、`GOSUMDB=sum.golang.org`（未關 sumdb，符合 C2）。每個 gosec 版本用獨立 GOPATH／GOMODCACHE／GOCACHE。

| # | 驗證 | 結果 |
|---|---|---|
| V1 | `go install …/gosec@v2.30.0`（Go 1.27.2） | exit 0；`go version -m`：built with go1.27.2，`mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`，`dep golang.org/x/tools v0.52.0` |
| V2 | h1 三方比對 | sum.golang.org lookup `h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`（go.mod `h1:qQn/kd9spDGy4yCVezLnzPkUuHPhm6jiIcBANad01Sw=`）＝`go mod download -json` 的 `Sum`／`GoModSum`＝暫存模組 `go.sum` 行＝`.ziphash`＝二進位 `go version -m` mod 行 → **全部一致** |
| V3 | 以 CI「Verify gosec module hash」同段邏輯重現 | v2.30.0＋新 expected：`gosec h1 OK` exit 0；負向（expected 對調）：`::error::gosec h1 mismatch` exit 1。pseudo-version＋現行 expected 同樣 OK／負向 exit 1 |
| V4 | main 2425129 乾淨副本（git archive 到 /tmp，未動任何工作樹）以 CI 參數 `-fmt=text -no-fail=false -severity=high -tests ./...` | v2.30.0：exit 0，Files 57、Lines 7485、Nosec 0、Issues 0；無 export data／type error |
| V5 | 同條件跑現行 pseudo-version 對照 | exit 0，Files 57、Lines 7485、Nosec 0、Issues 0；JSON 報告兩版位元組相同（sha256 同為 `92ef6150…`）；不加 `-severity` 兩版亦皆 0 發現 |
| V6 | G404 金絲雀（/tmp 副本加 `internal/canaryg404/token.go`，`math/rand.Intn` 產 token） | v2.30.0 以 CI 參數：報 `G404 (CWE-338) … (Confidence: MEDIUM, Severity: HIGH)`，Issues 1，**exit 1**；pseudo-version 同樣 exit 1。樣本未提交 |
| V7 | 擬議 workflow 模擬：取 `ci-devsecops.proposed.yml` sast-gosec 三個 run 步驟原文，在全新 GOPATH／GOCACHE 依序執行 | Install exit 0；Verify `gosec h1 OK` exit 0；Run Issues 0 exit 0 |
| V8 | 負向對照：v2.29.0 於 Go 1.27.2 | 重現 #1774：`export data version 5 is greater than maximum supported version 4`、30 行 SSA 錯誤，exit 1 |

**合併前仍須完成（未做，不得以本機結果代替）**：PR 最新 head 四項 required checks 全綠（D-10），`sast-gosec` log 有 `go version go1.27.2`、`mod … v2.30.0 h1:67aF6…` 與 `gosec h1 OK`；合併後 main push run 四綠；合併後依 C6 做法以不合併測試分支或 workflow_dispatch 再驗一次 G404 金絲雀會紅，之後刪分支。這些都會改到 GitHub，本 CR 起草時一律未執行。

## 6. 風險評估

| # | 風險 | 可能性／影響 | 緩解 |
|---|---|---|---|
| R1 | 新版規則收窄（G112／G118／G303／G102／constant resolve 等）造成漏報 | 低／中 | main 兩版結果相同；G404 金絲雀仍擋；請安全部評估是否需加金絲雀樣本 |
| R2 | gosec 二進位以 Go 1.27.2 建置，govulncheck `-mode=binary` 報標準庫 **GO-2026-6809**（crypto/x509 `Certificate.Verify`，1.27.3 修正）；pseudo-version 二進位同樣有，**非本 CR 新增**。另 module 層 GO-2026-5932（x/crypto/openpgp，無修補，未呼叫）兩版同樣存在 | 低（CI 內工具，短時執行）／低 | 不屬本 CR；依 C8 另開 CR 評估 Go 1.27.3（一次一件事）。專案本身 `govulncheck ./...`（Go 1.27.2）symbol 層 0 漏洞 exit 0 |
| R3 | 上游重打 tag／proxy 內容被換 | 極低／高 | sum.golang.org h1 鎖定＋CI 比對，不符即 exit 1（已以負向測試證明會擋） |
| R4 | v2.30.0 上游 CI 已改用 Go 1.27.3，對 1.27.2 的支援是本機實測得知，非上游宣告 | 低／中 | 本機 V1、V4、V7 已實測；PR 四綠再確認 |
| R5 | 間接依賴（coder/websocket v1.8.15、pb33f/go-yaml v0.1.1、otel/log v1.47.0）有漏洞 | 低 | **更正（v0.2，依安全部 30 號 K5）**：v0.1 稱「新增間接依賴」有誤。依安全部核對，這三個依賴在 CR-002 時即已存在，CR-005 前後依賴數皆為 48，僅為升版，非本 CR 新增。漏洞查詢：2026-10-27 查 govulncheck（DB 2026-10-21 版）、OSV、GitHub Advisory 皆 0 筆（見 C7 週檢紀錄） |

## 7. 回滾方案

觸發條件：PR 或合併後 `sast-gosec` 因 v2.30.0 失敗（安裝失敗、h1 不符、誤報無法在期限內處理）、或安全部判定規則收窄不可接受。

做法：以一個 PR（標 `CR-005 回滾`）把 sast-gosec 兩行改回現行值，其他不動：

```yaml
      - name: Install gosec
        run: go install github.com/securego/gosec/v2/cmd/gosec@v2.29.1-0.20261015133444-707f2271f665
...
          expected='h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40='
```

- 回滾值已在本機驗證：sum.golang.org lookup、`go mod download -json` Sum、`go version -m` 皆為 `h1:KvLOP3HNC+bEcwPGBQ4mVI70uhZ/scyFgiF1IVooj40=`，CI 同法比對 OK（證據 01b、03、04、06 pseudo707f227）。
- 等同 `git revert <CR-005 merge commit>`；回滾 PR 同樣須四綠與 D-08／D-09 程序。
- **禁止**以放寬門檻、`continue-on-error`、移除 required check、加 `-exclude`／`#nosec`、退回 v2.29.0（Go 1.27.2 下無法分析，證據 13）或退回 Go 1.27.1（有可達漏洞，見 CR-002）作為回滾。
- 回滾後 RR-010 恢復有效，C7 監控繼續，安全部重評期限仍為 2026-12-31。

## 8. 證據清單（`docs/09-cicd-env/evidence/cr-005/`，2026-10-27（UTC+8）產生；完整清單見同目錄 `SHA256SUMS`）

| 檔案 | 內容 | sha256 |
|---|---|---|
| 00-go1.27.2-download.txt | Go 1.27.2 tarball 官方 sha256 與本機校驗 | 3aeed08b0d57fa85474a6727069289b623834201685e6d1127a03eca10b0f567 |
| 01-sumdb-lookup-gosec-v2.30.0.txt | sum.golang.org lookup v2.30.0 | 3d904c6faa312f9a089eb50381932138c6214ac9aa17fb93d6cae9bd0a42f7ed |
| 01b-sumdb-lookup-gosec-pseudo-707f227.txt | sum.golang.org lookup 現行 pseudo-version（回滾值） | 3471b73a4a17613ef365ddc5f445b24535365a940226395214dd40b48b609100 |
| 02-install-gosec-v2300.txt | go install v2.30.0 | fa0d2c467d04e61ff5c6d11a0212b92eea0a7e197844637f7f208584bf811172 |
| 02-install-gosec-pseudo707f227.txt | go install pseudo-version | 25a158286e67ca89ecd075dbf80de4cbe9b327f369756651e5b49443332113c6 |
| 03-version-m-gosec-v2300.txt | go version -m（v2.30.0） | a60760a4c2bf1ea8b05caeb886968b59c70bd0ad22651f9bdca16e1e93586126 |
| 03-version-m-gosec-pseudo707f227.txt | go version -m（pseudo） | 197f2ec8832b8a6723ab4e2b06002fbaec6749947224a689e7cd55e9b2ef2b34 |
| 04-mod-download-json-v2300.txt | go mod download -json＋ziphash（v2.30.0） | d6d2b287794f67a0670e4d11b6302d0e56e6eecf0fc247d4a19255c08549e14b |
| 04-mod-download-json-pseudo707f227.txt | 同上（pseudo） | 936fc67e00bd49d6fdad7fe772cc588317f9dd0543d4b3876eae51bcc9334347 |
| 05-tempmodule-go-sum-v2300.txt | 暫存模組 go.sum 行 | 6b233775ff6dd3f98ca395b61aace49abbc609c5dbd77d33dfeb5e8d9c6f4248 |
| 06-ci-h1-check-repro-v2300.txt | CI h1 比對重現＋負向 | b93d24a8f570ac5fcfe46ffe9091494f139098a7e46ebe0e5eba88aaabdcf9db |
| 06-ci-h1-check-repro-pseudo707f227.txt | 同上（pseudo） | 2f02310885778a7d52a8408cd5201ca3507dde81eb9b6ca1aa5ab0874637d94a |
| 07-main-checkout.txt | main 2425129 乾淨副本來源 | 4c4b65c4a36f2cd3a34e3d345666c096b9b646e9d36404f404cfb90fb5d92d9a |
| 08-gosec-main-ci-args-v2300.txt | main 掃描（CI 參數，text） | 915def2c5e4714861514fba820969ee0b0890c173af2c92b440bc2eb2d938594 |
| 08-gosec-main-ci-args-pseudo707f227.txt | 同上（pseudo） | 221565e4b29424c00fded64a39ee099274970b11d0cd930742832a28cfcb0a80 |
| 09-gosec-main-ci-args-v2300.json | main 掃描（CI 參數，JSON） | 92ef6150380848c75c5412eb2525b4b3c93beb9a249edf9c82b15986f62733f8 |
| 09-gosec-main-ci-args-v2300.json.log | 其執行 log | 57691f4eb03dee550769c8e4389768ef98a439a8c7f25bbde6b06b68f46b87e6 |
| 09-gosec-main-ci-args-pseudo707f227.json | 同上（pseudo） | 92ef6150380848c75c5412eb2525b4b3c93beb9a249edf9c82b15986f62733f8 |
| 09-gosec-main-ci-args-pseudo707f227.json.log | 其執行 log | 190a84a398e0d4e713598a4a0b8f34c6221b01ad210b6573ea3f06cce8b551cb |
| 10-gosec-main-allsev-v2300.json | main 掃描不加 -severity | 92ef6150380848c75c5412eb2525b4b3c93beb9a249edf9c82b15986f62733f8 |
| 10-gosec-main-allsev-v2300.json.log | 其執行 log | c5f2f7663b5d7b6260f4b6561fa3e3a7fa025dbf4567b0810fbd787eb10d2702 |
| 10-gosec-main-allsev-pseudo707f227.json | 同上（pseudo） | 92ef6150380848c75c5412eb2525b4b3c93beb9a249edf9c82b15986f62733f8 |
| 10-gosec-main-allsev-pseudo707f227.json.log | 其執行 log | ecb88f7652b288cbf1fb3d8b11d87aad916e89b47b1d64ddddc21b9cbbcc2589 |
| 11-canary-g404-token.go.txt | 金絲雀樣本原始碼 | 34e1eaa832737090288abd40d85db040662fee1fd36ac02a8d3f0f34de0c55b0 |
| 12-canary-g404-v2300.txt | 金絲雀結果（v2.30.0） | 2c10272d637a5fcaae387d66ff28b69a43494f1563a766c2c89b6ea7239e4f2a |
| 12-canary-g404-pseudo707f227.txt | 金絲雀結果（pseudo） | 55de6be9e04bcdc1cc0dcf6eb0928c09bef7ff40dff775bf2556113bc002ef09 |
| 13-negative-control-gosec-v2.29.0.txt | v2.29.0 負向對照（#1774 重現） | 54b117fde7f2f590ba4335e622d8e11fd45d40dfc1302373d056d95174f8f757 |
| 14-proposed-workflow-sast-gosec-sim.txt | 擬議 workflow 步驟模擬 | 4b85ec9a263dd45f003655f42e906f41be23466b61f0ca5b265f88ee64479ea4 |
| 15-diff-pseudo-vs-v2300.txt | 兩版依賴與結果差異 | f27236cece2657d23e627409c4e5890c215d3b239c8e2e9f8305faae78f56061 |
| 16-upstream-compare-707f227-to-v2.30.0.txt | 上游 19 個 commit 清單 | 9ea941a39ff7bc3346aade91e3c64cb3f38c734cd06e1696523b02afeb7bed59 |

附件（`/workspace/rebuild/cr-005/`）：

| 檔案 | sha256 |
|---|---|
| CR-005-ci-devsecops-gosec-v2.30.0.patch | 7ddfa0f8bfd1e9465fdb85ed8f40c5da25b1d1645211519090eb302806df1853 |
| ci-devsecops.main-2425129.yml（main 現行原文） | da73b404f17c4e3c273b3d6e973c4f98d3a5f1ea2128709c58e3eaf022c27fe9 |
| ci-devsecops.proposed.yml | a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f |
| patch-validation.txt | b44f5dba358793a03a916074c58c2ab007f9f0f94d08b358cdc8b25a283c077d |

C7 週檢紀錄：`docs/09-cicd-env/evidence/c7-weekly/c7-check-2026-10-27.md`（v0.1，sha256 334fcc3160e0907c7c64af837f7d5f0b1c7978e4404fdc3e829965a8dab24430）；**現行版為 `c7-check-2026-10-27-v0.2.md`（sha256 801a5ea75bb99b71673a579152d1fbef2ec970a1ef51eb021d34371f4bc2e01c）**，只改用語（「新增三個間接依賴」改為「CR-002 時即已存在、僅升版」，依 K5／審查部 S3），查詢結果與證據檔不變；原始輸出清單 `SHA256SUMS-c7-2026-10-27-raw`（sha256 6f6df06dbad0cc7d1441b414f4276656cd457a441a21f708e9c79463d8048239）。

v0.3 新增引用（維運部 2026-10-27 重算）：

| 檔案 | 用途 | sha256 |
|---|---|---|
| `/workspace/rebuild/cr-005/CR-005-gosec-release-v2.30.0-v0.2.md` | 本版底稿（保留不動） | 5c6b6f56d8bbaec28016a12f1daf6faec8743a18ed95ae863d792efe007c25fd |
| `/workspace/review-CR-005-2026-10-27.md` | 審查部初審（退回） | 08adc77a274622589d2ca7bcc75a0d9cb69e3b34648efff7959558bb920c8cb2 |
| `docs/06-security/30-cr-005-gosec-v2.30.0-security-opinion-v0.1.md`（v0.1.1） | 安全部評估 | e4579d31ab6ae9ba883b5a2e84f21d2da22b2b456f658127bca15a893a598104 |
| `impact/企劃部.md` | 企劃影響評估 | 28b207d727865b51c2ed8bd94000af15f4ba51bf5fcc696e889a3b97524ffb4a |
| `impact/設計部.md` | 設計影響評估 | fa5a0b13689344a5fff2279f9def9bb898d30bce32f2b52604bc66d17bab0967 |
| `impact/研發部.md` | 研發影響評估 | b81b8d193ced8c35edd73860dab8cee260f3b1ba1431a8eaccef217ba52ac8a8 |
| `impact/rd-evidence/rd-cr005-{ci,all,rules}.json` | 研發掃描報告（三份相同） | 92ef6150380848c75c5412eb2525b4b3c93beb9a249edf9c82b15986f62733f8 |
| `impact/rd-evidence/rd-cr005-ci.log` | 研發掃描 log | 93f37db5d1a2eb45b1eb20e0fc5dea79cc3e98a6a6594614b56dc9ecdcd488eb |
| `impact/rd-evidence/rd-cr005-all.log` | 同上 | 737a94de5d3ef8f4b8d9b77149e012173b0e32e6aa2ddb3c0165090daf3d9172 |
| `impact/rd-evidence/rd-cr005-rules.log` | 同上 | ead69a18759f1c4a0d319bc12f20015724b045a9ca8bd4eeb4892a5e62771e41 |
| `impact/品保部.md` | 品保影響評估（QA-CR005-1～6） | 2244c51df983a7298ec7fe13f0b6c18f9e70924ebd44114d2dd0085f02b489b3 |
| `impact/安全部.md` | 安全 §4.1／§4.2 影響意見 | d8ea948d9f404748ce4ff087f57a2d870596caf3962783342b15dd5336c09d48 |


## 9. 建議與核准

- 維運建議：**接受**。理由：符合 C7、只改兩行、h1 三方一致、本機實測 main 結果與現行相同且金絲雀仍擋、回滾值已驗證。（本機實測僅為參考；Go 1.27.2 是否可用以 PR 的 CI 結果為準。）
- 核准路徑：各部門影響評估（已收，待對 v0.3 重新確認）→ 安全部評估（同意，附條件 K1～K6，30 號 v0.1.1）→ 審查部複審 → **使用者拍板** → 依 §9.4 流程 ①～⑧ 執行 → 結案。

| 角色 | 狀態 | 日期 | 備註 |
|---|---|---|---|
| 維運部（提出） | 已提出（v0.3） | 2026-10-27 | 本草稿 |
| 企劃部（影響） | 無影響 | 2026-10-27 | impact/企劃部.md（28b207d7…）；綁定 v0.2，待對 v0.3 重新確認 |
| 設計部（影響） | 無影響 | 2026-10-27 | impact/設計部.md（fa5a0b13…）；綁定 v0.2＋patch，待對 v0.3 重新確認 |
| 研發部（影響） | 無影響 | 2026-10-27 | impact/研發部.md（b81b8d19…）；綁定 v0.2，待對 v0.3 重新確認 |
| 品保部（影響） | 無影響，附 QA-CR005-1～6 | 2026-10-27 | impact/品保部.md（2244c51d…）；綁定 v0.2，待對 v0.3 重新確認 |
| 安全部（評估） | 同意（附條件，30 號 v0.1.1） | 2026-10-27 | 30 號 v0.1.1（e4579d31…）；impact/安全部.md（d8ea948d…）；兩者綁定 CR v0.1 |
| 審查部（確認） | v0.2 退回，v0.3 待複審 | 2026-10-27 | /workspace/review-CR-005-2026-10-27.md（08adc77a…） |
| 使用者（拍板） | 待拍板 | | 期限 2026-11-06 |

### 9.1 安全部條件（逐字引用 06-security/30 號 v0.1.1 §3，sha256 e4579d31ab6ae9ba883b5a2e84f21d2da22b2b456f658127bca15a893a598104）

> | # | 條件 | 時點 | 負責 |
> |---|---|---|---|
> | K1 | PR 只能改 `ci-devsecops.yml` 這 3 行，內容必須和 proposed.yml（sha256 `a67a0be4…`）完全一致；diff 對 main 只能有這一個檔案 | PR | 維運／研發 |
> | K2 | PR 最新 head 的四項 required checks 都要綠（D-10）。sast-gosec 的 log 要出現 `go1.27.2`、`mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa…`、`gosec h1 OK`，且命令列仍是 `-severity=high -tests` 加絕對路徑。安全部會以這個 head 出具書面確認 | 合併前 | 維運／安全 |
> | K3 | 依 D-09 人工核對 head；合併後 main push run 四項都要綠 | 合併時、合併後 | 維運 |
> | K4 | 合併後 7 日內，比照 C6 做一次 G404 金絲雀：用不合併的分支或 workflow_dispatch，sast-gosec 要紅，之後刪除分支，並提供 run 和 log 證據 | 合併後 | 維運 |
> | K5 | CR v0.2 修正 §2.2 的措辭：沒有新增依賴，只有升版 | 送審前 | 維運 |
> | K6 | 依 C8 另開 CR 升 Go 1.27.3（處理 GO-2026-6809），最遲 2026-11-20 提出；在那之前此項記入 05 號檔 | 合併後 | 維運（提出）／安全（評估） |

以上為原文；若與 v0.2 §9.1 的轉錄不同，一律以本原文為準（v0.2 轉錄的差異見 §10 v0.3 第 7 點）。

**建議（非條件）**：30 號 §2.3 正文建議 K6 的新 CR「在 CR-005 合併後 14 日內提出，最遲 2026-11-20」。依安全部（impact/安全部.md 末段）說明，以 §3 原文為準，「合併後 14 日內」是建議，不是條件。

### 9.2 條件落實對照（維運部；不改寫 §9.1 原文）

| 條件 | 落實步驟（§9.4） | 證據 | 現況 |
|---|---|---|---|
| K1 | ③：PR 只含 `.github/workflows/ci-devsecops.yml`，`git diff origin/main...HEAD --stat` 僅此一檔、3 行；PR head 該檔 blob sha256 須＝a67a0be4…（proposed.yml） | diff stat、blob sha256 | 未執行（須 ①② 後） |
| K2 | ④⑤：head 四綠；sast-gosec log 三項字串與命令列；安全部對該 head 出具書面確認 | run 連結、log 擷取、安全部書面確認 | 未執行 |
| K3 | ⑥⑦⑧：審查部與安全部確認的 head SHA＝合併時 head（D-09，由使用者本人合併時人工核對；維運取證）；合併後 main push run 四綠 | 合併紀錄、main run 連結 | 未執行 |
| K4 | ⑧：合併後 7 日內 G404 金絲雀，推金絲雀分支前須取得使用者核准；取得紅燈 run 後刪分支 | 紅燈 run、log、分支已刪證明 | 未執行 |
| K5 | 已於 v0.2 §6 R5 更正（30 號 K5 所稱「§2.2」為 30 號自身章節，CR 中對應處為 §6 R5）；C7 週檢 v0.2 同步更正 | §6 R5；C7 週檢 v0.2（801a5ea7…） | 已完成 |
| K6 | 合併後由維運提出新 CR（編號由總協調指派），最遲 2026-11-20；在那之前由安全部記入 05 號檔（§4.3） | 新 CR 編號、05 號檔改版 | 未執行 |

### 9.3 驗收條件（逐字引用品保部 `impact/品保部.md` §4，sha256 2244c51df983a7298ec7fe13f0b6c18f9e70924ebd44114d2dd0085f02b489b3）

> | ID | 條件 | 證據 | 通過標準 |
> |---|---|---|---|
> | QA-CR005-1 | PR 最新 head 四項 required checks 全綠（D-10）；合併後 main push run 四綠 | run 連結＋`gh pr checks`、`gh run view --json jobs` | 四項 success；check 名稱仍為 sca-govulncheck／sast-gosec／secrets-gitleaks／unit-test |
> | QA-CR005-2 | gosec 版本為 v2.30.0 | sast-gosec log：`go version go1.27.2`、`mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`、`gosec h1 OK` | 三段字樣皆在；**不以 `gosec -version` 為證** |
> | QA-CR005-3 | 發現差異有記錄 | 變更前 main run 與變更後 run 的 gosec Summary（Files／Lines／Nosec／Issues）並列 | 數值並列；不同時須有安全部判定；預期 57／7485／0／0（行數依當時 main 而定） |
> | QA-CR005-4 | 範圍只 3 行 | `git diff <base> <head> -- .github/` | 只有 3 行；與 proposed.yml blob 相同；其他檔不變 |
> | QA-CR005-5 | ruleset 未變 | 合併前後 `gh api …/rulesets/24629805`、`24629854` 以 `jq -S` 比對 | diff 為空（updated_at 除外） |
> | QA-CR005-6 | unit-test 結果不變 | unit-test job 的 go-test.json 或品保 main 重跑 | 0 FAIL；PASS 數與 CR 前 main 相同（f646bee 基準 181，若期間 main 有新測試以當時為準） |
>
> K4（G404 金絲雀紅燈）由維運執行、安全部判定；品保只在 CHK-SEC-009 舉證時引用其 run。

品保部同檔另載驗收要點（原文）：「若審查時要求改名，只能改步驟名，**不得改 job id 或 job name**」；「`gosec -version` 對 `go install` 建置的二進位輸出 `Version: dev`……**不能用來證明版本**」；「本機結果不得替代 CI 結果」。

**Go 1.27.2 是否可用，以 PR 的 CI 結果為準；本機實測僅為參考。** QA-CR005-1～6 與 K1～K4 全數通過，CR-005 才可結案。

### 9.4 執行程序（①～⑧，依審查部 CR5-R3）

| 步驟 | 內容 | 執行者 | 前置條件 | 產出 |
|---|---|---|---|---|
| ① | 使用者拍板 CR-005 | 使用者 | 審查部複審通過；各部門對 v0.3 確認 | 拍板紀錄（總協調登記於 99-changelog） |
| ② | 使用者在**維運部聊天**親自核准推送 | 使用者 | ① | 聊天中的核准訊息 |
| ③ | **維運部**（單一指定推送者）推送指定分支並開 PR：PR 標題與內文標 `CR-005`，**只動 `.github/workflows/ci-devsecops.yml` 一個檔**，內容＝proposed.yml（a67a0be4…） | 維運部 | ② | PR 連結；head SHA；diff stat |
| ④ | PR 最新 head 四項 required checks 全綠（secrets-gitleaks、sast-gosec、sca-govulncheck、unit-test） | CI（維運觀察取證） | ③ | run 連結、log 擷取（K2、QA-CR005-1／2／3） |
| ⑤ | 安全部依 K2 對該 head 出具書面確認 | 安全部 | ④ | 書面確認（載明 head SHA） |
| ⑥ | 審查部對**同一 head** 給「可合併」；head 變動即失效，須回到 ④ | 審查部 | ⑤ | 「可合併」紀錄（載明 head SHA） |
| ⑦ | **使用者本人合併**（依 D-09 人工核對 head＝⑤⑥ 所載 SHA） | 使用者 | ⑥ | merge commit |
| ⑧ | 合併後 main push run 四項全綠（K3、QA-CR005-1）；再做 K4 G404 金絲雀，**推金絲雀分支前須取得使用者核准**；取得紅燈 run 後刪分支 | 維運部（金絲雀）；使用者（核准推金絲雀分支） | ⑦ | main run 連結；金絲雀紅燈 run／log；分支已刪證明 |

推送方式（維運部提議，於 ② 時一併請使用者確認）：
- 工作樹：以 `git worktree add` 從當時 origin/main 新建專用工作樹（提議 `/workspace/ops-cr005`），**不沿用** `/workspace/ops-cr002`、`/workspace/dev/short-url` 等既有工作樹（品保部已指出後者停在舊分支）。
- 分支名（提議）：`ci/CR-005-gosec-v2.30.0`。
- 套用：將 `ci-devsecops.proposed.yml` 覆蓋 `.github/workflows/ci-devsecops.yml`，推送前核對該檔 sha256＝a67a0be4…、`git diff --stat origin/main` 只有一檔。若屆時 origin/main 的 workflow 已非 da73b404…，**停止推送**，回報後重新評估（不得自行調整 patch）。
- 維運部在 ② 之前不推送、不開 PR、不改 GitHub 任何設定。

### 9.5 待決事項

1. 企劃、設計、研發、品保四部門的評估綁定 v0.2，安全部 impact 檔與 30 號綁定 v0.1；v0.3 技術內容未變，但依各檔效力條款須對 v0.3（sha256 見交付紀錄）重新確認。
2. 審查部 CR5-R4 要求「安全部並確認 v0.2 的 §9.1 轉錄無誤」：v0.3 §9.1 已改為逐字引用 30 號 v0.1.1 §3，請安全部確認 v0.3 §9.1 與原文一致。
3. 99-changelog 登記 CR-005（總協調）。
4. 審查部 S6：上游 master 已比 v2.30.0 多 9 個 commit，尚無新發行。維運建議：CR-005 合併且 K4 完成後，C7 週檢結束，改由每月工具版本檢視（devsecops-ci-pipeline 步驟 6）追蹤發行版；在合併前，C7 週檢照常（下次最遲 2026-11-03）。是否結束 C7 由安全部決定。
5. 步驟名稱 `Verify gosec module hash (CR-002)` 保留（審查部 S7、安全部 30 號 §3 末段同意），日後再改此 workflow 時一併改中性名稱。

## 10. 修訂紀錄

| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 維運部 | 初稿；本機驗證 V1～V8 |
| v0.2 | 2026-10-27 | 維運部 | 依安全部 30 號：更正 R5「新增間接依賴」說法（K5）；新增 §9.1 條件 K1～K6；簽核表更新安全部狀態。其他內容不變 |
| v0.3 | 2026-10-27 | 維運部 | 依審查部退回紀錄（08adc77a…）修訂；**技術內容與 v0.2 相同**（§1、§2、§3、§5、§6、§7、patch 7ddfa0f8…、proposed.yml a67a0be4… 皆未改）。逐處改動：(1) 檔頭版本列改 v0.3；(2) 檔頭狀態列由「已提出（待安全部評估 → 審查部 → 使用者拍板）」改為各角色現況（含安全部「同意（附條件，30 號 v0.1.1）」、審查部「v0.2 退回，v0.3 待複審」）（S1）；(3) 檔頭關聯列加 30 號 v0.1.1、審查紀錄、C7 週檢 v0.2 及其 sha（S3）；(4) 標題下新增「技術內容與 v0.2 相同」說明框與「Go 1.27.2 是否可用，以 PR 的 CI 結果為準；本機實測僅為參考」（S5）；(5) 新增「影響的需求編號」一列（S2）；(6) §4.1 改為彙整各部門原文列，含評估檔路徑、sha256、綁定版本、具名日期（CR5-R1）；§4.2 由勾選清單改為逐項判定表，安全部五項照 impact/安全部.md §4.2，並附 30 號 §2.1 結論；新增 §4.3 文件更新分工（03、05〔RR-010、GO-2026-6809〕、13、01-ci-pipeline、CR-002 紀錄、99-changelog；負責與時點）（CR5-R1、S4）；19 個 commit 清單與 main 兩版結果說明原文保留；(7) §9.1 改為逐字引用 30 號 v0.1.1 §3 原表，引用 sha 由 0be01bb5… 改為 e4579d31…（CR5-R2、CR5-R4）。v0.2 轉錄與原文差異：K1 漏「diff 對 main 只能有這一個檔案」、負責漏「研發」；K2 漏「安全部會以這個 head 出具書面確認」、命令列條件寫成「gosec 命令列不變」、log 字串寫成 `v2.30.0 h1:67aF6…`（缺 `mod github.com/securego/gosec/v2` 前綴）、漏「(D-10)」、負責漏「安全」；K3 負責原文為「維運」（v0.2 寫「使用者（合併）／維運（取證）」）；K4 v0.2 多寫「推分支前須使用者核准」（改列為 §9.4 ⑧ 程序）；K5 文字與原文不同；K6 漏「在那之前此項記入 05 號檔」、負責漏「安全（評估）」；另刪除 v0.2 表外「推送與開 PR……」一句，併入 §9.4；「合併後 14 日內」只列為建議；(8) 新增 §9.2 條件落實對照；(9) 新增 §9.3 驗收條件，逐字引用品保 QA-CR005-1～6；(10) 新增 §9.4 執行程序 ①～⑧ 與推送方式（CR5-R3）；(11) 新增 §9.5 待決事項（含 S6、S7）；(12) §9 維運建議句末加註本機實測僅供參考；核准路徑改寫；簽核表加四部門列、更新安全部與審查部列；(13) §8 C7 週檢引用改為 v0.1／v0.2 並列及其 sha，新增 v0.3 引用檔 sha 表（S3）；(14) 本修訂紀錄新增 v0.3 列 |
