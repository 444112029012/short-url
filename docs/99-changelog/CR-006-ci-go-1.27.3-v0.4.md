---
文件：變更申請 CR-006（CI 升 Go 1.27.3，處理 GO-2026-6809／RR-012）
版本：v0.4
狀態：草稿
負責角色：維運部
最後更新：2026-10-28
對應需求：SEC-009、SEC-010、NFR-004、SI-05；RR-012；06-security/13 C8；06-security/30 K6；06-security/35（v0.1.2）S1～S7、§6.2、§7
專案代號：SHORTURL
變更紀錄：v0.1｜2026-10-28｜維運部｜初稿（sha256 e4722b6f4cbaae711c91c23a0425bb7895483209ab3de88a81b87fae2369e2a8）；v0.2｜2026-10-28｜維運部｜整合五份影響分析和安全部 35 號裁示（sha256 5d233a71da15d9385b4eb035ad9f5979cbfed3f397d0c2bb5f1fc2e21124e83b）；v0.3｜2026-10-28｜維運部｜依五部門對 v0.2 的回覆修訂（sha256 e4d9ca21e9193f2bbde9fcf14dada29982e7cbe96259db9099a1121e4edf4fa8）；v0.4｜2026-10-28｜維運部｜QA-M10 與 QA-S7～S9
---

# CR-006：CI 升 Go 1.27.3（處理 GO-2026-6809／RR-012）

> 時區：Asia/Taipei（UTC+8）；以 Z 結尾的時間為 UTC，文中已換算。
> 本草稿只做唯讀查詢與本機實測：**沒有開分支、沒有推送、沒有 commit、沒有改 CI 或 GitHub 上任何東西**；`docs/` 底下所有檔案只讀不改；各部門影響分析檔只引用，不修改。本機結果僅供參考，**Go 1.27.3 是否可用以 PR 的 CI 結果為準**。
>
> **v0.2 說明：技術內容與 v0.1 相同。** patch `CR-006-ci-go-1.27.3.patch`（sha256 `c3939e62c4a3a9081332fbdd7fd06bb056085bd40f732fe97b25acbcac56eaec`）與 `ci-devsecops.proposed.yml`（sha256 `06833a314722d4bc7b436533a6268b2be29d5c56d028eeb1bfd20f9151def245`）**未修改**：35 號裁示採方案甲、不加斷言步驟，都不需要改 patch。本版只改寫判準、條件、驗收、執行程序與文件分工，逐處改動見 §10 與 `CR-006-v0.1-to-v0.2.diff`。
>
> **v0.3 說明：技術內容仍與 v0.1 相同**，patch 與 proposed.yml 的 sha256 未變（同上兩值，`evidence-v0.3/18` 重算）。本版依五部門對 v0.2 的回覆修訂：更新影響分析引用（企劃部 v0.2、設計部 v2、品保部 v0.2、35 號 v0.1.1）；QA-M5 改由安全部獨立重跑（35 號 v0.1.1 §7）；L6 依 QA-M6、QA-M7 修正；採納 QA-S5、QA-S6。逐處改動見 §10 與 `CR-006-v0.2-to-v0.3.diff`。
>
> **v0.4 說明：技術內容仍與 v0.1 相同**，patch 與 proposed.yml 的 sha256 未變（`evidence-v0.4/19` 重算）。本版只收錄品保部 v0.3 的 QA-M10（L6 計數不加 `^`）與 QA-S7～S9，並在 §9 簽核表、§9.5 寫入各部門對 v0.3 的確認紀錄。逐處改動見 §10 與 `CR-006-v0.3-to-v0.4.diff`。

- 提出者／日期：維運部／2026-10-28（v0.1）；v0.2、v0.3 同日
- 類型：**標準變更**（依 13 號 C8：「Go 後續修補版（1.27.3…）同樣以精確釘版＋CR（標準變更）更新」，`docs/06-security/13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md` 第 115 行；03 號 §2 Go toolchain 列同旨，`docs/06-security/03-ci-security-gates-v0.3.md` 第 52 行）。核准與執行程序依使用者指定，比照 CR-005 §9.4 走完整流程（§9.4）。
- 變更內容：見 §1
- 原因與效益：見 §2
- 不變更的後果：見 §3
- 影響的需求編號：SEC-009（SAST，門檻不變）、SEC-010（SCA，門檻不變）、NFR-004、SI-05（取自 `docs/09-cicd-env/01-ci-pipeline-v0.2.md` 第 7 行「對應需求」）；REQ：無（企劃部判定，`impact/企劃部.md` v0.2 第 31 行）；RR-012
- 期望完成時間：**最遲 2026-11-20（UTC+8）提出**（30 號 K6）；**最遲 2026-11-20 合併**（35 號 S1）；合併後 7 日內完成 G404 金絲雀（S4）；G5 前在部署清單加上建置用 Go 版本（S6，無具體日期）。執行日期待使用者拍板後排定。企劃部附註：**提出與合併期限同為 2026-11-20，建議使用者拍板時一併排定合併日**（`impact/企劃部.md` v0.2 第 38 行原文：「建議使用者拍板時一併排定合併日期，並留出 CI 重跑的時間（main 前進時 head 會變，要重做確認，CR §6 R3）」）
- 附件／參考：§8

## 0. 基準（v0.3 於 2026-10-28 14:03（UTC+8）複查）

`git fetch origin` 後 origin/main 仍為 **`ca09137ae8f336a8bb02ac7f20e893eb2a6ccfbe`**（以公開 https 未登入 `git ls-remote` 同值；未登入 API 回 HTTP 200，`evidence-v0.3/18`）；workflow sha256 仍為 `a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f`；go.mod 仍為 `31c3fd2790509a6ce56c44831b5a609746c1a39656ad95e2b4892ae98f6bbb49`（第 3 行 `go 1.27.2`）；go.sum 仍為 `1e2b43d75b3352825d6227289fc6f3f42f6d883b4d171c7db1ada91afd49dd72`（`evidence-v0.2/17-v0.2-baseline-recheck.txt`；v0.3 複查 `evidence-v0.3/18-v0.3-baseline-recheck.txt`，全部與 v0.2 相同）。

| 項目 | 實際值 | 來源 |
|---|---|---|
| origin/main | `ca09137ae8f336a8bb02ac7f20e893eb2a6ccfbe`（PR #30 docs/cr-005-followups 的 merge，2026-10-28 10:00:05（UTC+8））。v0.1 時已確認不是原先預期的 `db380a40…`；兩者之間 `.github/`、`go.mod`、`go.sum` 無差異 | `git fetch origin`、`git rev-parse origin/main` |
| main 上的 `ci-devsecops.yml` | blob `7d9cf80cb3d51ccf49ee1167f36520bdfb4bc6eb`，sha256 `a67a0be4…e8f`（＝CR-005 proposed.yml） | 副本 `ci-devsecops.main-ca09137.yml` |
| main 最新 push run（＝本 CR 現時的 merge-base 基準） | run 40618331206（head ca09137，2026-10-28T02:00:08Z＝10:00:08（UTC+8）），四項 success。三個 job：`Setup go version spec 1.27.2`、`go version go1.27.2 linux/amd64`、`GOVERSION='go1.27.2'`、`GOTOOLCHAIN='auto'`；快取 key `setup-go-Linux-x64-ubuntu24-go-1.27.2-cf6b0d429736099b597565b998fd2f9b82c99c65d9cffc4d8c9270df02c0a025`；三個 job 全 log 對 `switching to go\|downloading go1\.` 計數皆 0；sast-gosec Files 57／Lines 7485／Nosec 0／Issues 0；sca-govulncheck 有「This scan also found 1 vulnerability in packages you import」；unit-test artifact 181 pass／0 fail／0 skip、TestCT_* 14 | `evidence/15`、`evidence-v0.2/17`（sast-gosec 計數）、品保 `impact/qa-evidence/10`（sca、unit-test） |
| ruleset 24629805 `main-required-checks` | active、bypass_actors 空、required contexts＝sca-govulncheck／sast-gosec／secrets-gitleaks／unit-test、`strict_required_status_checks_policy: true` | `evidence/14` |
| ruleset 24629854 `main-pr-review` | bypass_actors：RepositoryRole 5，`bypass_mode: pull_request`（QA-S3）；基準快照 `jq -S 'del(._links, .current_user_can_bypass)'` | `evidence/14` 第 53～57 行；品保 `impact/qa-evidence/08` |

## 1. 變更內容（現況 → 期望）

### 1.1 Workflow 中 Go 版本與 setup-go 的所有位置（origin/main ca09137 `.github/workflows/ci-devsecops.yml`）

| 行號 | job | 內容（現況） | 要改？ | 期望 |
|---|---|---|---|---|
| 6 | （檔頭註解） | `# 工具鏈釘 Go 1.27.2（CR-002；check-latest: false）。…` | **改（僅註解）** | `# 工具鏈釘 Go 1.27.3（CR-006 由 CR-002 的 1.27.2 升修補版，修 GO-2026-6809；check-latest: false）。…`（同行後半不動） |
| 60 | sast-gosec | `- uses: actions/setup-go@40f1582b2485089dde7abd97c1529aa768e1baff # v5` | 不改 | — |
| 62 | sast-gosec | `go-version: "1.27.2"` | **改** | `go-version: "1.27.3"` |
| 63 | sast-gosec | `check-latest: false` | 不改 | — |
| 89 | sca-govulncheck | `- uses: actions/setup-go@40f1582…` | 不改 | — |
| 91 | sca-govulncheck | `go-version: "1.27.2"` | **改** | `go-version: "1.27.3"` |
| 92 | sca-govulncheck | `check-latest: false` | 不改 | — |
| 132 | unit-test | `- uses: actions/setup-go@40f1582…` | 不改 | — |
| 134 | unit-test | `go-version: "1.27.2"` | **改** | `go-version: "1.27.3"` |
| 135、136 | unit-test | `check-latest: false`、`cache: true` | 不改 | — |
| 29–53 | secrets-gitleaks | 無 setup-go、無 Go 版本 | 不適用 | — |

- 三處都要改：13 號 C1 要求三個 job 精確同版（第 108 行），03 號 §2 要求「兩者須一致」（第 52 行）。
- **三處 setup-go 與 go.mod 必須放在同一個 commit 改**（研發部條件 2，`impact/研發部.md` 第 54 行；35 號 S1）。理由：CI 是 `GOTOOLCHAIN=auto`，只要有一處漏改，該 job 會在 job 內默默下載 1.27.3（研發部 (b)，`rd-evidence/41`），工具鏈就不再由 setup-go 決定（§6 R2）。
- **其他 workflow**：origin/main 的 `.github/` 底下只有這一個 workflow。repo 其他地方的 Go 版本：`go.mod` 第 3 行、`README.md` 第 14 行；沒有 Dockerfile、Containerfile、`.go-version`、`GOTOOLCHAIN` 設定（維運 v0.1 查證；品保 `impact/品保部.md` 第 38 行同結論）。
- 其他遠端分支上的 workflow 有 `1.27.x`、`1.24.x`、`1.27.1` 等舊值，屬歷史分支，不在本 CR 範圍。

### 1.2 go.mod（**go.mod 寫法只在本節說明；裁示見 §9.1 第 1 項**）

- 現況（origin/main `go.mod` 第 3 行）：`go 1.27.2`；沒有 `toolchain` 行。
- **定案：方案甲，只改 go 行 `go 1.27.2` → `go 1.27.3`，不加 `toolchain` 行。**（35 號 §3 第 1 項：「採甲：只改 `go 1.27.3`，不加 toolchain 行。乙和丙都不採用」；研發部同意甲，`impact/研發部.md` 第 34 行。）方案乙、丙不再列為選項。
- 依據：01-ci-pipeline §2.1 維護規則（第 55 行）；03 號 §2 Go toolchain 列（第 52 行）；CR-002 前例（commit `412cefa`）。研發部補充：官方升版方式 `go get go@1.27.3` 與 `go mod tidy` 的產出與 patch 完全相同（go.mod sha256 `02842cf5…`），不會自動加上 toolchain 行（`rd-evidence/42`）。
- 不夾帶依賴升級：`go mod tidy` 後 go.sum 位元組不變（`1e2b43d7…`）、模組圖 27 個模組不變；vet、`test -race`、govulncheck 通過（`evidence/11`；研發部 `rd-evidence/12`～`18`、`42`）。不做 `go get -u`、不改 require。
- 對 govulncheck 的影響：govulncheck 掃描時用的 stdlib 版本取決於執行的 Go 工具鏈，不是 go 行（`evidence/10`；研發部 `rd-evidence/25`：Go 1.27.3 掃未改的 go.mod 也沒有 GO-2026-6809）。
- **對 gosec 的影響（v0.2 依研發部 (b) 更正）**：gosec 的**安裝**（`go install …@v2.30.0`）在空的 main module 中進行，**不讀本 repo 的 go.mod**；但 gosec **掃描時**透過 `go list` 載入套件，**會讀本 repo 的 go.mod**。因此 go.mod 改成 `go 1.27.3` 後，若 PATH 上是 Go 1.27.2 且 `GOTOOLCHAIN=local`，gosec 會對每個套件報 `go.mod requires go >= 1.27.3`，**exit 1，Files 0**（`rd-evidence/31` 第 75～99 行）；govulncheck 也會 exit 1（`rd-evidence/24`）。失敗是明確的紅燈，**不會出現假綠燈**。CI 是 auto，若有漏改則會切換到 1.27.3（由 L3 偵測，§5.1）。（v0.1 寫「改不改 go.mod，gosec 二進位都由 setup-go 的 Go 建置」只對安裝成立，對掃描不完整，已更正。）
- 副作用：本機 Go 1.27.2＋`GOTOOLCHAIN=local` 的 build／vet／test／gosec／govulncheck 都會失敗；`GOTOOLCHAIN=auto` 會自動下載 go1.27.3（需能連 proxy.golang.org；離線且無快取時 exit 1）（研發部 (a)～(b')，`rd-evidence/40`、`41`）。處理見 §4.3（`env-go1273.sh`）與 §6 R6。

### 1.3 擬議 patch（只存在 `/workspace/rebuild/cr-006/`，未套用到任何 repo；v0.2 未修改）

```diff
--- a/.github/workflows/ci-devsecops.yml
+++ b/.github/workflows/ci-devsecops.yml
-# 工具鏈釘 Go 1.27.2（CR-002；check-latest: false）。1.22 標準庫有 GO-2025-3750；1.24.13 是 1.24 最後一版，
+# 工具鏈釘 Go 1.27.3（CR-006 由 CR-002 的 1.27.2 升修補版，修 GO-2026-6809；check-latest: false）。1.22 標準庫有 GO-2025-3750；1.24.13 是 1.24 最後一版，
-          go-version: "1.27.2"        # sast-gosec（第 62 行）
+          go-version: "1.27.3"
-          go-version: "1.27.2"        # sca-govulncheck（第 91 行）
+          go-version: "1.27.3"
-          go-version: "1.27.2"        # unit-test（第 134 行）
+          go-version: "1.27.3"
--- a/go.mod
+++ b/go.mod
-go 1.27.2
+go 1.27.3
```
（行尾 `#` 註記只為說明，不在 patch 內。）

`patch-validation.txt`：`git apply --check` exit 0；`git diff --numstat`：`4	4	.github/workflows/ci-devsecops.yml`、`1	1	go.mod`；套用後 workflow＝`06833a31…f245`（blob `af4c394`）、go.mod＝`02842cf5c5b0a74916d55ffa85fcb05db9afc800e4609e8b1c5c35e045f28ebf`。研發部與品保部各自在暫存 worktree／clone 重跑，結果相同（`rd-evidence/02`；`impact/品保部.md` 第 37 行）。

### 1.4 不變（一次只改一件事）

gosec `v2.30.0` 與 h1 `h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`、步驟名稱 `Verify gosec module hash (CR-002)`、govulncheck `v1.8.0`、gitleaks `8.21.2`（含 sha256）、`-fmt=text -no-fail=false -severity=high -tests ./...`、`check-latest: false`、`cache` 設定、所有 action SHA、四個 job 名稱、ruleset、dependency-review fail-closed、go.sum、所有 require：**全部不動**。**不加**建置版本斷言步驟、**不設** `GOTOOLCHAIN`（35 號 §3 第 2、5 項，留給 SCR-053 的後續強化 CR）。

## 2. 原因與效益

- **來由**：30 號 §3 K6 原文：「依 C8 另開 CR 升 Go 1.27.3（處理 GO-2026-6809），最遲 2026-11-20 提出；在那之前此項記入 05 號檔｜合併後｜維運（提出）／安全（評估）」（`docs/06-security/30-cr-005-gosec-v2.30.0-security-opinion-v0.1.md` 第 61 行，v0.1.1，sha256 `e4579d31ab6ae9ba883b5a2e84f21d2da22b2b456f658127bca15a893a598104`）。
- **風險登錄**：RR-012（`docs/06-security/05-residual-risk-register-v0.1.md` 第 33 行，v0.1.13，sha256 `df8636d892cc9a3faf85c1fc01a6db08a6c87a2855face0ffca346090ad9f7c2`）。
- **流程依據**：13 號 C8（第 115 行，v0.1.3，sha256 `32f7f8c29f21d8fbbe244d6b38d4330b8c21663a5ebe5b04658245543976343e`）。
- **安全部評估**：35 號 **v0.1.2**（檔名未改，仍為 `docs/06-security/35-cr-006-go-1273-security-opinion-v0.1.md`，sha256 `dc94d52b5c7475a21c78d6fd377c3a3725bf9d836bcd6060141b5f3b8895df0a`；v0.1.1 為 `a557a914d0435e98fb8055f8da7b5ed153a3bb1925dca123a15578647f74ceb2`；v0.1 為 `da44487464c8a1858b2a9c31082cd3371461da98c0aaac0cea202f1f74a4e8f2`）§1 第 14 行：「**同意（附條件 S1～S7）。** 風險等級：低。」（註：35 號目前只在本機 `/workspace/projects/short-url/docs/`，origin/main 上尚無此檔，`evidence-v0.2/17`。）
- **上游事實**（v0.1 查證，不變）：Go 1.27.3 於 2026-10-19 發行（`evidence/01` 第 12 行）；`go1.27.3.linux-amd64.tar.gz` sha256 `65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40`（`evidence/03`；研發部重新下載核對相同，`rd-evidence/01`）；GO-2026-6809 為 stdlib `crypto/x509` `Certificate.Verify`，修正版本 1.26.10、1.27.3（`evidence/02`；安全部 13:26:52～13:28:09（UTC+8）重查一致，35 號 §2 第 2 項）。
- **效益**：gosec 二進位改用 1.27.3 建置後，`govulncheck -mode=binary` symbol 層由「GO-2026-6809，exit 3」變成 0（`evidence/09`），RR-012 可在 S3、S4 完成後由安全部結案（35 號 S5）；sca-govulncheck 的 package 層提示由 1 變 0（§5.1 L6）；go.mod 與 CI 一致。

## 3. 不變更的後果

- 30 號 K6 逾 2026-11-20 未提出即違反條件；35 號 S1 要求最遲 2026-11-20 合併，未合併亦屬違反條件。企劃部附註：「如果 2026-11-20 前沒有提出，K6 條件會逾期；沒有合併，S1 條件會逾期。兩者都須由總協調呈報使用者。」（`impact/企劃部.md` v0.2 第 37 行）。
- CI 三個 job 持續使用含 GO-2026-6809 的 1.27.2 標準庫；RR-012 無法結案。
- **14 日規則（v0.2 依 35 號 §3 第 4 項改寫）**：GO-2026-6809 情境嚴重度為**低**（一般評等：中），14 日規則只針對 High 且有修補的漏洞，**不適用、不會阻擋**。v0.1 寫的「可能開始擋」已刪除。

## 4. 影響分析

### 4.1 各部門影響（各列依各部門原文；檔案均於 2026-10-28 14:03（UTC+8）以 `sha256sum` 核對，與總協調轉達值相符，`evidence-v0.3/18`）

| 部門 | 評估檔 | sha256（完整） | 綁定的 CR 版本 | 結論（原文） | 具名／日期 |
|---|---|---|---|---|---|
| 企劃部 | `impact/企劃部.md` **v0.2**（v0.1 為 0279ff38…34f，已由 v0.2 取代） | 9387362bdd3d20d96be638d7615ca7431ab60833416b10a5869191cfc0b9f1c7 | v0.2（5d233a71…e83b） | 「結論：範圍無影響；時程沒有新增或改動任何已排定的項目（v0.2 新增的期限都是條件，不是里程碑）」（第 21 行） | 企劃部（2026-10-28） |
| 設計部 | `impact/設計部.md` **v2**（v1 為 0c708efc…4c5d，已由 v2 取代） | 32ef61411cc7ba94396027cb3b244cf783cf4acf8db2ee77cbad6bb0e5aa845b | v0.2（5d233a71…e83b）＋patch（c3939e62…eaec） | 「**無影響**。設計文件沒有寫到 Go 版本，不需同步更新」（第 9 行；v2 修訂紀錄第 43 行：「結論不變（無影響）」） | 設計部 2026-10-28（UTC+8） |
| 研發部 | `impact/研發部.md` v0.1（證據 `impact/rd-evidence/`，清單 `SHA256SUMS` sha256 0b97a8018d10e23d5a408eb3f6d1a7144e59855a73ea08ec3ffb7b025cfcf6b5，維運以 `sha256sum -c` 驗 32 檔全 OK） | 5275dd2cf436100ef9ee4f875aa4c1a2ff8c192aab74e88e4e6053f4a9b944af | v0.1（e4722b6f…）＋patch（c3939e62…） | 「低影響，不必改程式碼（程式碼工作量 0）；§六 第 1 點（go.mod 方案甲）有條件可行」 | 研發部 2026-10-28 |
| 安全部 | `impact/安全部.md` v0.1（沿用；依據 35 號，現為 **v0.1.2**，sha256 dc94d52b5c7475a21c78d6fd377c3a3725bf9d836bcd6060141b5f3b8895df0a；35 號 §6 第 61 行：該檔「對 v0.2 也一樣有效，所以那份檔案不改版」；§6.2 第 78 行：「`impact/安全部.md` v0.1（1dc296b2…）仍有效，不改版」） | 1dc296b2f4e02cff4e30ec20b1f3234216660edff523d5ce59be0a5a758c4432 | v0.1（e4722b6f…）、patch c3939e62…、proposed.yml 06833a31… | 「結論：**同意（附條件 S1～S7）**，詳見 35 號 §5」 | 安全部（具名）2026-10-28 |
| 維運部 | 本 CR §4.1 下表維運列 | （本檔） | v0.2 | 只改 CI 版本字串與 go.mod go 行；不影響執行環境與部署；本機實測僅供參考 | 維運部 2026-10-28 |
| 品保部 | `impact/品保部.md` **v0.3**（v0.2 為 f1bdc807…2e86、v0.1 為 44917b51…aae3c，均已取代；`impact/qa-evidence-v0.3/` 清單 14b88a4400813f3afacf2ebdb87f73329e0e54a79efac24925fe4be56712d159 驗 3 檔 OK；證據 `impact/qa-evidence/` 清單 748b8d57996875827133d3dc17da2c444bc85169ab3c17b4e67b0ed3790664c5 驗 11 檔 OK；`impact/qa-evidence-v0.2/` 清單 5199c03ee3b987dc926e3ec680fc314b7dbc59fcf31f3725202299c4615faf09 驗 7 檔 OK） | 5b8869cbb5f6de401e8ea747453cc29f88f38e418f6351d8ab97594253d621ba | v0.3（e4d9ca21…4fa8）＋SHA256SUMS 7ae77e92…e8ba | 「結論：CR PR 內無影響（同 v0.1）；v0.3 對 QA-M1、M5～M9、S5、S6 的改法全部放行，但新增 1 項實質必改 QA-M10（L6「開頭」字樣在 CI log 上會漏判），**本版不放行**（見 §8）」（第 25 行；QA-M10 已於 v0.4 收錄） | 品保部 2026-10-28 |
| 審查部 | — | — | — | **待審查** | — |

**效力（v0.4 更新）**：品保部 v0.3 與 35 號 v0.1.2 綁定 CR **v0.3**（`e4d9ca21…`）；企劃部 v0.2、設計部 v2 檔案仍綁 v0.2，但兩部門已對 v0.3 書面確認沿用（§9 簽核表）。以下為 v0.3 當時的說明：企劃部 v0.2、設計部 v2、品保部 v0.2、35 號 v0.1.1 都綁定 CR **v0.2**（`5d233a71…`）；研發部與安全部 impact 檔沿用（仍綁 v0.1 文字，研發部第 129 行只在技術內容改變時重新確認，安全部依 35 號 v0.1.1 §6 第 61 行沿用）。v0.3 未改 patch 與 proposed.yml，但 CR 文件 sha256 又變了：企劃部（v0.2 第 19 行）、設計部（第 36 行）、品保部（v0.2 第 23 行）的效力條款要求**對 v0.3 重新確認**；35 號 v0.1.1 的效力條款（§6 第 64 行）寫明「S1～S7 的落實方式有任何改變，就要重新確認」，v0.3 依安全部自己的 §7 改了 QA-M5 的落實方式，**仍須安全部對 v0.3 確認**（§9.5 第 1 項）。

各面向原文列（照各評估檔 §4.1 列逐字貼入；維運列由維運填）：

| 面向 | 影響描述 | 需更新項目 | 工作量 | 填寫角色 |
|---|---|---|---|---|
| 範圍與時程 | **範圍**：無影響。CR-006 只改 CI 三個 job 的 setup-go `go-version`、workflow 第 6 行註解與 `go.mod` 第 3 行，而且放在同一個 commit（CR v0.2 §1.1～1.3；35 號 S1）；不新增或刪除功能；IN-01～IN-06、OUT-01～OUT-06（`docs/01-requirements/03-scope-v0.1.md`）都不變，屬 IN-06「資安／維運章節」既有範圍內的維護。**時程**：只列已確定的時點，不填估算。（1）最遲 2026-11-20（UTC+8）提出（30 號 K6；CR v0.2 第 25 行）；（2）**v0.2 新增**：最遲 2026-11-20 合併（35 號 S1，第 46 行；CR v0.2 第 25 行）；合併由使用者本人執行（CR v0.2 §9.4 ⑪），所以這也是使用者動作的期限；（3）G404 金絲雀已改為必做，合併後 7 日內完成（35 號 S4，第 49 行）；（4）G5 前要在部署清單加上建置用 Go 版本（35 號 S6，第 51 行），沒有具體日期。以上都是條件期限，不是里程碑。執行、合併的實際日期 CR v0.2 第 25 行仍寫「待使用者拍板後排定」，目前**沒有排定日期**，本列不填。CR-006 不改 `docs/00-charter/02-preliminary-schedule-v0.1.md` 的任何一列（該檔日期本來就標明是粗估目標，不是實際排程） | 無（`03-scope`、`00-charter`、`02-preliminary-schedule` 都不必改版） | 企劃部 0（不需改文件）；其他部門工作量由各部門自填 | 企劃 |
| 架構／資料／API／UI | 無影響。設計文件（docs/02-architecture、03-data、04-api、05-ui-ux）都沒有寫 Go 版本；ADR-001A 只定「Go＋chi」，沒有釘版本。設計所依賴的 net/http 傳輸層行為（標頭上限、431），在 1.27.2 與 1.27.3 之間原始碼完全相同 | 無 | 0 | 設計 |
| 實作 | **低影響，不必改 .go。** 在 main ca09137 的 detached worktree 套用 patch 後，以官方 Go 1.27.3（sha256 已核對）執行 `go mod verify`／`go build`／`go vet`／`go test`／`go test -race`／`go test -count=1 -json` 全部 exit 0：10 個套件 pass、1 個無測試檔（internal/ports）、181 個測試（含 75 個子測試）全部 pass，fail 0、skip 0；`go mod tidy` 後 go.mod／go.sum 位元組不變，也沒有加上 toolchain 行；模組圖 27 個模組不變；專案 govulncheck 已不再列出 GO-2026-6809；gosec v2.30.0 能正常載入 go 1.27.3 的 module（Files 57、Issues 0、無 Golang errors）。副作用：Go 1.27.2＋`GOTOOLCHAIN=local` 的 build／vet／test／gosec／govulncheck 都會因 `go.mod requires go >= 1.27.3` 失敗；預設 `GOTOOLCHAIN=auto` 則會自動下載 go1.27.3（需能連 proxy.golang.org） | 程式碼：不適用（0 個 .go 檔）。開發環境：本機 Go 須升到 ≥ 1.27.3，或改用 `GOTOOLCHAIN=auto`（合併前由維運通知，見 §6 R6）。`README.md` 第 14 行的 Go 版本說明合併後會過時（負責角色依 §4.3 待確認） | 程式碼 0；開發者本機升級 Go 屬環境作業，不計入程式碼工作量 | 研發 |
| 資安（威脅模型、SEC） | 沒有新的資料流或信任邊界，威脅模型不受影響。SEC-009 和 SEC-010 的門檻不變。換成 Go 1.27.3 可以修掉 gosec 二進位裡的 GO-2026-6809（情境嚴重度低，不適用 14 日阻擋規則），RR-012 可以結案。go.mod 採甲案。這次不加斷言步驟，GOTOOLCHAIN 留到後續 CR（SCR-053）。合併後要做 G404 金絲雀。x/net 和 x/crypto 屬於範圍外的資訊級發現（SCR-054、SCR-055），納入每月檢視。結論：**同意（附條件 S1～S7）**，詳見 35 號 §5 | 03 號 §2 的 Go 列；05 號 RR-012 結案；13 號 C8 補註；安全部對 PR head 出具書面確認（S2）；金絲雀覆核（S4） | 安全部約 0.5 人日 | 安全部 |
| 環境／部署／監控 | 只改 CI 三個 job 的 setup-go `go-version`（`"1.27.2"`→`"1.27.3"`）、檔頭第 6 行註解，以及 go.mod 第 3 行 `go 1.27.3`（方案甲，35 號裁示）；三處與 go.mod 在同一 commit（S1、研發條件 2）。action SHA、`check-latest: false`、cache 設定、gosec／govulncheck／gitleaks 版本、h1、門檻、job 名稱、ruleset 都不變；go.sum 與依賴圖不變。不影響執行環境與部署（repo 無 Dockerfile；G5 部署清單須補「正式建置用 Go 版本＝CI（1.27.3）」，S6）。CI 第一次跑 1.27.3 時快取 key 改為 `…go-1.27.3-cf6b0d42…a025`，可能未命中、時間變長，不影響結果（§6 R1）。本機實測結果同 v0.1（僅供參考，以 PR CI 為準）。合併前由安全部依 §5.2 對 PR head 獨立重跑、品保部核對（QA-M5，35 號 v0.1.1 §7）；維運只提供 head SHA 與 CI 取證（v0.3 改） | CR PR：`.github/workflows/ci-devsecops.yml`（4 行）、`go.mod`（1 行）。合併後：`docs/09-cicd-env/01-ci-pipeline-v0.2.md` §2.1；`docs/09-cicd-env/evidence/cr-006/`（只放維運 CI 取證，隨文件 PR；§5.2 證據改由安全部存 `docs/06-security/evidence/cr-006-m5/<head SHA>/`）；G5 部署清單（S6，G5 前）；本機 `/workspace/tools/env-go1273.sh` 新增、`env-go1272.sh` 檔頭加註（§4.3，不進 repo）；G404 金絲雀（S4） | 約 1 人日（PR、CI 取證 L1～L7、提供 head SHA、金絲雀、§2.1 文件、環境檔；§5.2 重跑改由安全部執行） | 維運 |
| 測試與追溯矩陣 | **CR PR 內無影響**：矩陣 02 號 v0.5（main `ca09137` 與本機 sha256 皆 `8a005aef…9932`）沒有任何列引用 Go 1.27.2、setup-go 版本或 GO-2026-6809（`git grep -nE '1\.27\.[0-9]\|GO-2026\|RR-012\|CR-005\|CR-006' ca09137 -- docs/07-testing`，只命中 02 號第 140、143 行的 box 實跑環境說明，以及 11、12、13、19 號與 evidence 的歷史紀錄，都屬當時事實，不改）。SEC-009、SEC-010、NFR-004 列對應 CHK-SEC-009／010、CHK-NFR-004，判定內容不變，**但 G5 舉證須改引用 CR-006 合併後的 run**。合併後在矩陣下一版 §8 流程類補上 CR-005、CR-006 兩列（**非阻擋**），另出 CR-006 驗收結論檔（見 §4） | CR PR 內：無。合併後：02 號下一版 §8 加兩列；新增 07-testing 21 號（CR-006 驗收結論）；README 索引同步 | 驗收比對與結論檔約 0.5 人日；矩陣兩列約 0.1 人日 | 品保 |
| 測試策略（06 v0.1） | 無影響：策略未指定 Go 版本，層級、工具類別、門檻不變 | 無 | 0 | 品保 |
| DoD（08 v0.1） | 無影響：S1～S5（全綠＋安全文件）條文不變 | 無 | 0 | 品保 |
| 程式行為／單元測試 | 不改 Go 程式與測試；以 Go 1.27.3 建置與執行測試後結果應相同（品保本機參考：套用 patch 後 181 PASS、CT 14 PASS、覆蓋 80.9%，與基準相同，見 §1） | 無 | 0 | 品保 |
| 文件一致性與追溯性 | 待審查 | 待審查 | 待審查 | 審查 |

註：企劃（v0.2 第 25 行）、設計（v2 第 17 行）、研發（第 21 行）、安全（第 13 行）、品保（v0.3 第 29～32 行）各列為原文逐字貼入（品保部原檔有四列，全部照貼；v0.2、v0.3 與 v0.1 逐字相同，品保部 v0.3 第 34 行）。

### 4.2 影響分析清單（逐項判定）

| 清單項 | 判定 | 判定者（依據） | 要更新的檔案與章節 |
|---|---|---|---|
| 需求規格 | REQ 無影響；NFR-004 驗收條件仍成立；SEC-009／SEC-010 **不受影響**（門檻和 job 名稱都不變） | 企劃部（v0.2 第 31 行）；安全部（`impact/安全部.md` §4.2 第 18 行） | 無 |
| 範圍內外、時程、里程碑 | 企劃部原文：「同上表：範圍無影響；沒有新的里程碑；目前確定的日期是 2026-11-20 提出期限（K6）與 2026-11-20 合併期限（S1），加上合併後 7 日內的金絲雀（S4）」；同檔第 25 行另列 S6（G5 前，無具體日期），並附註提出與合併同日、建議拍板時一併排定合併日（第 38 行） | 企劃部（v0.2 第 32 行） | 無（`03-scope`、`00-charter`、`02-preliminary-schedule` 不必改版） |
| 系統架構（C4）與 ADR | 無影響；不需新 ADR，ADR-001A 不改 | 設計部（核對內容第 3、4 點） | 無 |
| 資料模型與資料流程圖 | 無影響 | 設計部（第 7 點） | 無 |
| API 規格（OpenAPI） | 無影響（net/http 標頭上限相關原始碼兩版相同，431 契約不變） | 設計部（第 6 點） | 無 |
| UI/UX | N/A | 設計部 | 無 |
| 威脅模型 | **不受影響** | 安全部（第 19 行） | 無（04 號不改版） |
| 資安需求／03 號 | **要更新**，只改工具版本欄 | 安全部（第 20 行） | 03 號 §2 Go toolchain 列改 1.27.3，S3、S4 完成後由安全部改版（S5） |
| 殘餘風險（05 號） | **要更新** | 安全部（第 21 行） | RR-012 在 S3、S4 完成後結案；回滾時恢復有效（S7）。SCR-053～055 不開 RR（35 號 §4 第 41 行） |
| 13 號 | **要補註** | 安全部（第 22 行） | C8 補註「1.27.3 由 CR-006 執行」（S5） |
| 隱私合規 | 無影響 | 企劃部（v0.2 第 33 行） | 無（`04-privacy-compliance-v0.1.md` 不必改） |
| 測試策略、測試案例與追溯矩陣 | **CR PR 內不改 07-testing；合併後補件，非阻擋**（QA-S4 採納，直接引用 `impact/品保部.md` v0.3 §5 結論；v0.1 的 §4 在 v0.2 起改為 §5，品保部 v0.3 第 34 行）：矩陣 02 號下一版 §8 加 CR-005、CR-006 兩列；另出 07-testing 21 號 CR-006 驗收結論；README 索引同步；06／07／08／09 號不改版 | 品保部（v0.3 §5，第 165～176 行） | 合併後由品保部處理 |
| CI/CD、IaC、環境設定、機密 | **要更新** | 維運部 | `.github/workflows/ci-devsecops.yml`、`go.mod`（CR PR，同一 commit）；`docs/09-cicd-env/01-ci-pipeline-v0.2.md` §2.1（合併後）；本機環境檔（§4.3，不進 repo） |
| 部署與回滾手冊、資料庫遷移可逆性 | 回滾見 §7；無資料庫遷移；G5 部署清單須補建置用 Go 版本（S6） | 維運部；安全部（S6） | G5 部署清單（G5 前） |
| 監控告警、備份、事故 runbook | 無影響 | 維運部 | 無 |
| 使用者文件與版本說明 | `README.md` 第 14 行合併後會過時 | **待審查部裁示**；維運部建議合併後另外提（研發部表示若被指派可提 1 行文件變更，`impact/研發部.md` 第 128 行） | README.md 第 14 行 |
| 變更紀錄（99-changelog） | 要登記 CR-006 | 總協調 | 維運不改該目錄 |

### 4.3 文件與環境更新分工（都不放進 CR PR）

| 檔案 | 內容 | 負責 | 時點 | 依據 |
|---|---|---|---|---|
| `docs/09-cicd-env/01-ci-pipeline-v0.2.md` | §2.1：CR-002 的 Go toolchain 列（第 49 行）改「已由 CR-006 取代」；新增 CR-006 列 | 維運部 | 合併、main 四綠（S3）後，以文件變更提出（推送須使用者核准） | 本 CR |
| `docs/09-cicd-env/evidence/cr-006/` | **只放維運的 CI 取證**（L1～L7 擷取、run 連結、全 log sha256）＋SHA256SUMS | 維運部 | 合併前（⑦）與合併後（⑫）產生，隨文件 PR 入庫（推送須使用者核准） | 35 號 v0.1.2 §7 第 91 行 |
| `docs/06-security/evidence/cr-006-m5/<head SHA>/` | §5.2 PR head 獨立重跑證據（M5-1～M5-8 輸出、exit code）＋安全部自己的 SHA256SUMS | **安全部**（執行）；品保部（核對，記錄該 SHA256SUMS 的 sha256） | 合併前（⑧）就要存在；入庫隨文件 PR（推送須使用者核准，入庫由總協調安排） | 35 號 v0.1.2 §7 第 91 行、§6.2 第 4 項（第 80 行）；品保部 v0.3 §3.3 第 5 點（第 134 行） |
| `docs/06-security/03-ci-security-gates-v0.3.md` §2 Go 列 | 改 1.27.3 | 安全部 | S3、S4 之後 | 35 號 S5 |
| `docs/06-security/05-residual-risk-register-v0.1.md` RR-012 | 結案 | 安全部 | S3、S4 之後 | 35 號 S5 |
| `docs/06-security/13-…-v0.1.md` C8 | 補註「1.27.3 由 CR-006 執行」 | 安全部 | S3、S4 之後 | 35 號 S5 |
| G5 部署清單（`docs/09-cicd-env/07-g5-deploy-checklist-v0.1.md`） | 加「正式建置用的 Go 版本與 CI 相同（1.27.3）」並附證據 | 維運部 | G5 前 | 35 號 S6 |
| `docs/07-testing/02` 矩陣下一版 §8、新增 21 號驗收結論、README 索引 | 見 §4.2 | 品保部 | 合併後（非阻擋） | `impact/品保部.md` v0.2 §5 |
| `README.md` 第 14 行 | Go 版本說明 | **待審查部裁示**（維運建議合併後另外提） | 合併後 | — |
| `docs/99-changelog/` | 登記 CR-006、decision-log | 總協調 | 即時 | — |
| `/workspace/tools/env-go1273.sh`（**本機工具，不進 repo**） | 新增，內容比照 `env-go1272.sh`，改指向 `/workspace/tools/go1.27.3`（`GOTOOLCHAIN=local`、`GOFLAGS=` 空、`GOPROXY=https://proxy.golang.org,direct`、`GOSUMDB=sum.golang.org`） | 維運部 | 合併、main 四項全綠（S3）之後 | 研發部殘餘風險第 1 點（第 124 行），(c) |
| `/workspace/tools/env-go1272.sh`（本機工具） | **保留**作為回滾用；檔頭加註「已由 CR-006 取代（僅供 CR-006 回滾時使用）」。**本版不改此檔** | 維運部 | 同上 | (c) |
| 研發部本機工具 `/workspace/tools/gp-v2300/bin/{gosec,govulncheck}` | 以 go1.27.2 建置，二進位仍帶 GO-2026-6809 的標準庫；**只影響本機，與 CI 無關**（CI 每次以 setup-go 的 Go 重新 `go install`）。**建議研發部合併後用 1.27.3 重建** | 研發部（建議） | 合併後 | 研發部殘餘風險第 3 點（第 126 行），(d) |

## 5. 驗證

### 5.1 CI log 判準 L1～L7（PR 最新 head 的 run 與合併後 main push run **都要成立**，QA-S1）

**比較基準**：PR merge-base 那個 main commit 的 push run（QA-M4）。以目前 main 計，merge-base＝`ca09137`，基準 run＝40618331206（§0）；推送時若 main 已前進，以當時 merge-base 的 push run 為準，並先依 §9.4 ⑥ 停止條件處理。

| # | 範圍 | 判準（變更後必須成立） | 變更前（基準 run 40618331206） | 對應 QA |
|---|---|---|---|---|
| L1 | sast-gosec、sca-govulncheck、unit-test 三個 job | 各 job 都有 `Setup go version spec 1.27.3`、`go version go1.27.3 linux/amd64`、`GOVERSION='go1.27.3'`、`GOROOT='/opt/hostedtoolcache/go/1.27.3/x64'`；任一 job 的 log 以 `grep -E 'go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/'` 命中即不通過（**v0.3 依 QA-S5 縮小比對字樣**，避免 dependency-review 等其他輸出的版本字樣誤中）。另**記錄**（不判定）`GOTOOLCHAIN='auto'`（QA-S2） | 三個 job 皆 1.27.2；`GOTOOLCHAIN='auto'` | QA-CR006-1 |
| L2 | 同上三個 job | 快取 key 以前綴 **`setup-go-Linux-x64-ubuntu24-go-1.27.3-`** 開頭，且雜湊段＝**`cf6b0d429736099b597565b998fd2f9b82c99c65d9cffc4d8c9270df02c0a025`**（go.sum 不變故不變）。命中或未命中**只記錄，不作判定**（QA-M2）。若 runner image 改版使 `ubuntu24` 段不同，記錄並轉品保判定 | 前綴 `…go-1.27.2-`，雜湊段同 | QA-CR006-2 |
| L3 | 同上三個 job 的**完整 log** | `grep -cE 'switching to go\|downloading go1\.'` 每個 job 都是 **0**（QA-M3；涵蓋 Install gosec、go install govulncheck、go vet、go test 等所有步驟） | 三個 job 皆 0（`evidence-v0.2/17`；品保 `qa-evidence/10`） | QA-CR006-3 |
| L4 | sast-gosec「Verify gosec module hash (CR-002)」 | `mod	github.com/securego/gosec/v2	v2.30.0	h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=` 與 `gosec h1 OK` 皆在。不以 `gosec -version` 為證，也不以 h1 證明建置用 Go 版本 | 相同 | QA-CR006-4 |
| L5 | sast-gosec「Run gosec」 | 命令列仍為 `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...`；**Issues 0、Nosec 0、exit 0 為硬條件**；Files／Lines 必須與 merge-base 基準相同（本 CR 不改 `.go` 檔）。任何差異記錄並轉安全部判定（QA-M4） | 57／7485／0／0 | QA-CR006-5 |
| L6 | sca-govulncheck「Run govulncheck」（v0.2 新增，QA-M1；**v0.3 依 QA-M6、QA-M7 修正**） | 對 sca-govulncheck「Run govulncheck」步驟的 log：`grep -c 'No vulnerabilities found\.'` ≥ 1；`grep -c 'This scan also found'` ＝ **0**（**不加 `^` 錨點**，v0.4 依 QA-M10）；exit 0；同 job 的 L1 為 `GOVERSION='go1.27.3'`。**為什麼不加 `^`**：`gh run view --log` 的 CI log 每一行前面都有 job 名稱、步驟名稱和時間戳，`grep '^This scan also found'` 對變更前的 run 也得到 0（不錨定才是 1，品保 `qa-evidence-v0.3/16`），等於永遠判通過；變更後若 vuln DB 新增 package／module 層項目，這個硬條件會被誤判為通過，QA-M7 的處理也不會啟動（品保部 v0.3 §8.4 QA-M10，第 258 行）（govulncheck 掃描時用的 stdlib＝執行的 Go）。**不要求** `Your code is affected by 0 vulnerabilities.`（QA-M6：govulncheck v1.8.0 在三層皆 0 時只印 `No vulnerabilities found.` 一行，品保部 v0.3 §3.2 第 115 行 B1、B2；研發 `rd-evidence/20` 同）。CI 預設輸出不印 stdlib 版本與 ID，stdlib 版本由 §5.2 M5-6 佐證。PR 上 dependency-review 不得有 Critical。**非 GO-2026-6809 項目的處理（QA-M7）**：出現 GO-2026-6809 以外的項目時，L6 直接判不通過，轉安全部依 03 號 §3 門檻裁定，不得為了讓本 CR 過而放寬。安全部原文（35 號 v0.1.1 §6.1(b)，第 69 行）：「L6 是硬條件。如果 PR 的 sca-govulncheck 出現 GO-2026-6809 以外的項目，就先判定 L6 不通過，交給安全部依 03 號 §3 的門檻裁定。裁定結果只有兩種：一是不屬於本 CR 造成，記錄後依門檻處理；二是要擋就擋。不得因為要讓本 CR 通過而放寬門檻（對應 R7）。」品保部 v0.3 §3.4 QA-M7（第 144 行）另要求：由安全部以 `-show verbose` 確認 ID 後裁定 | 變更前（品保 `qa-evidence-v0.2/11` A1、A2 與 CI 相同）：`=== Symbol Results ===`、`No vulnerabilities found.`、`Your code is affected by 0 vulnerabilities.`、「This scan also found 1 vulnerability in packages you import and 0 vulnerabilities in modules you require」（GO-2026-6809 的 package 層，`qa-evidence/10`）；exit 0；GOVERSION go1.27.2 | QA-CR006-10(a)、11 |
| L7 | unit-test（**v0.2 新增，QA-M1**） | 同 job L1 為 `go1.27.3`；artifact `go-test-json` 以 `jq` 統計：**0 fail、0 skip，pass 數＝merge-base 基準**（ca09137 為 181，含子測試）；TestCT_* pass 14；log 有 `--- PASS: TestENG011_PerSourceIsolationByPeerAddress`；指令仍含 `-race` | 同樣 181／0／0、CT 14、ENG-011 PASS，但在 go1.27.2 下（前後差別由 L1 的版本看出，測試數必須相同） | QA-CR006-7、8 |

**為什麼 h1 不會變**（同 v0.1）：`h1:` 是 module 內容雜湊（dirhash Hash1），只跟 v2.30.0 原始碼有關，跟編譯器版本無關（`evidence/05`、`13`）。所以 h1 只能證明原始碼，建置版本要看 L1＋L3＋§5.2。

**不加斷言步驟**（35 號 §3 第 2 項）：v0.1 §5.1 C 的可選斷言**本次不加**，留給 SCR-053 的後續強化 CR。替代證據＝**L1＋L3＋本機重現**：同一 job、同一 PATH 的 `go` 是 1.27.3（L1），全 log 沒有切換或下載工具鏈（L3），因此 `go install` 建置 gosec 用的就是 1.27.3；本機重現（v0.1 `evidence/05`、`09`）與合併前由安全部獨立執行的 PR head 重跑（§5.2）直接證明 `go version -m` 第一行為 go1.27.3、gosec 二進位無 GO-2026-6809。

**L1 與 L3 必須合併判讀（v0.3 依 QA-S6）**：L3 偵測**有訊息**的工具鏈切換（`switching to go…`、`go: downloading go1.x.y (os/arch)`）；工具鏈若已在模組快取中，切換可能不印訊息，這個缺口由 L1 的 `GOROOT='/opt/hostedtoolcache/go/1.27.3/x64'` 偵測（**沒有訊息**的切換會讓 GOROOT 或 `go version` 不符）。兩者任一不符即不通過，不得單獨以 L3＝0 推論建置版本（品保部 v0.3 §3.1 第 104 行、§3.5 第 152 行）。

### 5.2 合併前置條件：安全部對 PR head 獨立重跑（QA-M5；**v0.3 依 35 號 v0.1.1 §7 改寫，以 §7 原文為準**）

35 號 v0.1.2 §7 原文（第 87～93 行，逐字引用；與 v0.1.1 第 72～78 行逐字相同）：

> ## 7. QA-M5 和 QA-CR006-10(b)：改由安全部獨立重跑（v0.1.1 新增）
> **安全部接下這兩項。** 這樣做是為了職責分離：提出者是維運部，如果也由維運部驗證，等於自己驗自己。
> - **執行者**：安全部。**核對者**：品保部。維運部負責提供 PR head 和 CI 取證，不執行這次重跑。
> - **步驟**：照 CR v0.2 §5.2 的 1～7 步，並補兩點。第一，Go 1.27.3 的 tarball 由安全部自己對照 go.dev/dl 的 JSON 核對 sha256（65efbc80…9b40），不直接沿用其他部門的工具目錄結果。第二，只用公開 clone 和未登入的 API 取得 head，全程只讀，不推送，也不碰 GitHub 寫入。
> - **證據位置**：因為安全部不改其他部門的目錄，證據改放在 `docs/06-security/evidence/cr-006-m5/`，並附上安全部自己的 SHA256SUMS。請維運部把 CR §5.2 第 8 步、§4.3 和 §9.4 ⑧⑨ 的執行者和證據路徑改成這個；`docs/09-cicd-env/evidence/cr-006/` 只放維運的 CI 取證就好。
> - **時點**：PR 最新 head 四項全綠、L1～L7 成立之後才開始。head 一變就重跑。重跑結果和 S2 的書面確認會綁定同一個 head，一起出具。
> - **判準**：gosec 第一行是 go1.27.3，h1 是 67aF6LHa…；`-mode=binary` 掃描中 GO-2026-6809 為 0 筆；專案 `-show verbose` 顯示 go1.27.3 的標準庫，而且沒有 GO-2026-6809；所有輸出中 `switching to go|downloading go1\.` 的計數為 0。任何一項不符，就不出具 S2 確認。

- **執行者**：**安全部**（具名）。**核對者**：**品保部**。**維運部只提供 PR head SHA 和 CI 取證**（§5.1 L1～L7），不執行本重跑。
- **對象**：PR 最新 head SHA。維運部提供 head SHA；安全部以公開 clone 或未登入 API 獨立取得並比對（35 號 §7 第 90 行），寫進證據檔第一行。**head 一有變動就重跑，舊證據作廢**；重跑結果與 S2 書面確認**綁定同一個 head**，一起出具（35 號 §7 第 92 行）。
- **時點**：PR 最新 head 四項全綠、L1～L7 成立之後（35 號 §7 第 92 行）；審查部「可合併」之前。
- **證據位置**：`docs/06-security/evidence/cr-006-m5/<head SHA>/`，附安全部自己的 `SHA256SUMS`；合併前就要存在（box 上 `/workspace/projects/short-url/docs/06-security/evidence/cr-006-m5/<head SHA>/`），品保核對時記錄該 SHA256SUMS 的 sha256；入庫隨文件 PR，推送須使用者核准。`docs/09-cicd-env/evidence/cr-006/` 只放維運的 CI 取證（35 號 §7 第 91 行）。
- **全程只讀**：不推送、不碰 GitHub 寫入（35 號 §7 第 90 行）。

**判準 M5-1～M5-8**（指令沿用 v0.2 §5.2 第 1～7 步，加上 35 號 §7 兩點補充；任一判準不符，安全部就不出具 S2 確認）：

| 判準 | 指令（安全部執行） | 預期輸出 | exit |
|---|---|---|---|
| M5-1 取得 head | **兩種方式都要做**（v0.4 依 QA-S8）：(i) 公開 `git ls-remote https://github.com/444112029012/short-url refs/pull/<N>/head`；(ii) 未登入 `curl https://api.github.com/repos/444112029012/short-url/pulls/<N>` 的 `.head.sha`；`git clone https://github.com/444112029012/short-url <tmp>`、`git -C <tmp> fetch origin <head SHA>`、`git -C <tmp> checkout --detach <head SHA>`；`sha256sum go.mod .github/workflows/ci-devsecops.yml` | (i)＝(ii)＝維運提供的 head SHA，三個值都寫進證據檔。**任一不一致就停止（不進 M5-2）、記錄三個值、等 head 穩定後重取**；PR 仍開啟時以 ls-remote 為主（API 可能有快取延遲）；未登入 API 限額 60 次／小時（品保部 v0.3 §8.3 第 2 點，第 248 行；`qa-evidence-v0.3/14`）。維運提供的值只供比對，不能取代安全部自取的值（35 號 §6.2 第 5 項，第 81 行）。`git rev-parse HEAD`＝該 SHA；go.mod＝`02842cf5c5b0a74916d55ffa85fcb05db9afc800e4609e8b1c5c35e045f28ebf`、workflow＝`06833a314722d4bc7b436533a6268b2be29d5c56d028eeb1bfd20f9151def245` | 0 |
| M5-2 工具鏈 | 安全部自行下載 `https://go.dev/dl/go1.27.3.linux-amd64.tar.gz`，**對照 `https://go.dev/dl/?mode=json&include=all` 的官方值**以 `sha256sum -c` 核對，解壓到安全部自己的暫存目錄（**不直接沿用** `/workspace/tools/go1.27.3` 等其他部門的工具目錄結果，35 號 §7 第 90 行）。環境比照 CI：PATH／GOROOT 指向該目錄、不設 `GOTOOLCHAIN`（沿用 go.env 的 auto）、`GOFLAGS=` 空、`GOPROXY=https://proxy.golang.org,direct`、`GOSUMDB=sum.golang.org`、全新暫存 GOPATH／GOCACHE；`go version`；`go env GOVERSION GOTOOLCHAIN` | `…linux-amd64.tar.gz: OK`（官方值 `65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40`）；`go version go1.27.3 linux/amd64`；`go1.27.3`、`auto` | 0 |
| M5-3 安裝 | 在 checkout 目錄內照 CI：`go install github.com/securego/gosec/v2/cmd/gosec@v2.30.0`；`go install golang.org/x/vuln/cmd/govulncheck@v1.8.0` | 只有模組下載訊息，無錯誤 | 0、0 |
| M5-4 建置版本 | `go version -m "$(go env GOPATH)/bin/gosec"` | **第一行 `…/gosec: go1.27.3`**；`mod	github.com/securego/gosec/v2	v2.30.0	h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=` | 0 |
| M5-5 gosec 二進位（**主要指令一**） | `govulncheck -mode=binary "$(go env GOPATH)/bin/gosec"`，另跑一次 `-show verbose`；`grep -c GO-2026-6809` | `No vulnerabilities found.`（symbol 層）；**GO-2026-6809 計數 0**；module 層預期只有 SCR-054／055 的 3 筆（GO-2026-6866、6860、5932），只記錄、不算不通過（品保部 v0.3 §3.2 E，第 118 行；§8.2，第 233 行） | **govulncheck exit 0 且 GO-2026-6809 計數為 0（兩者都要）**（v0.4 依 QA-S7）。此欄指 govulncheck 的 exit code，不是 `grep -c` 的（`grep -c` 計數 0 時本身 exit 1）。module 層 3 筆是預期內的：35 號 §6.2 第 7 項（第 83 行）原文「M5-5 的「必須 0」是指 GO-2026-6809 必須 0 筆；module 層預期出現的 SCR-054／055 共 3 筆，不算不通過」。**不得改用 `-scan module`**（會因這 3 筆必然 exit 3，品保部 v0.3 §8.2 第 235、243 行）。若 govulncheck exit 3：是 GO-2026-6809 即不通過；是其他 symbol 層項目，比照 QA-M7 由安全部依 03 號 §3 裁定，不得放寬 |
| M5-6 專案 stdlib（**主要指令二**，QA-CR006-10(b)） | `govulncheck -show verbose ./...` | 「Govulncheck scanned the following … modules and the go1.27.3 standard library」；`No vulnerabilities found.`；**無 GO-2026-6809**（品保部 v0.3 §3.2 C，第 116 行；§8.2 第 238 行） | **必須 0** |
| M5-7 無工具鏈切換 | 對 M5-2～M5-6 全部輸出 `grep -cE 'switching to go\|downloading go1\.'` | 0 | — |
| M5-8 留證 | 每步輸出含指令、時間（`date`）、exit code 存到 `docs/06-security/evidence/cr-006-m5/<head SHA>/`，`sha256sum * > SHA256SUMS`；刪除暫存 clone | SHA256SUMS 可 `sha256sum -c` 全 OK | 0 |

- **品保核對**：品保部核對 M5-1～M5-8 的輸出與 exit code、記錄 SHA256SUMS 的 sha256，確認 head SHA＝安全部 S2 書面確認所載 SHA，才可進入審查部「可合併」（§9.4 ⑨）。

### 5.3 本機實測結果（v0.1，僅供參考；與研發部、品保部獨立實測一致）

| # | 驗證 | 結果 | 證據 |
|---|---|---|---|
| V1 | 下載 go1.27.3 tarball 並核對 sha256 | OK；`go version go1.27.3 linux/amd64` | `evidence/03`；`rd-evidence/01` |
| V2 | Go 1.27.3 安裝 gosec v2.30.0、govulncheck v1.8.0 | exit 0 | `evidence/04`；`rd-evidence/32` |
| V3 | `go version -m` 比較 | 1.27.3 版第一行 `go1.27.3`、1.27.2 版 `go1.27.2`；h1 相同 | `evidence/05`；`rd-evidence/32` 第 53～54 行 |
| V4 | 擬議 sast-gosec 三步模擬 | `gosec h1 OK`；57／7485／0／0，exit 0 | `evidence/07`；`rd-evidence/33`、`34` |
| V5 | 加擬議 go.mod 後掃描 | Issues 0，exit 0 | `evidence/07b` |
| V6 | 可選斷言正／負向（**本次不採用**，留作 SCR-053 參考） | 正向 exit 0；負向 exit 1 | `evidence/08` |
| V7 | govulncheck `-mode=binary` 掃兩個 gosec | 1.27.3 版 symbol 層 0；1.27.2 版 GO-2026-6809，exit 3 | `evidence/09` |
| V8 | 專案 `govulncheck ./...` | 1.27.3：無；1.27.2：package 層 1 | `evidence/10`；`rd-evidence/20`～`25` |
| V9 | go.mod 改 `go 1.27.3` 的影響 | go.sum／模組圖不變；vet、test -race 通過；1.27.2＋local 拒絕建置 | `evidence/11`；`rd-evidence/12`～`18`、`40`～`42`；`qa-evidence/06`、`07`（181 pass、CT 14、`-coverpkg` 80.9%） |
| V10 | 本機 G404 金絲雀 | G404 HIGH，Issues 1，exit 1 | `evidence/16` |

**合併前仍須完成（不得以本機結果代替）**：PR head 四綠且 L1～L7 成立（S2）；§5.2 由安全部對 PR head 獨立重跑（M5-1～M5-8）並經品保部核對；安全部書面確認（S2，與重跑綁同一 head）。合併後：main push run 四綠且 L1～L7 成立（S3）；G404 金絲雀（S4）。這些都會動到 GitHub，本 CR 起草時一律未執行。

## 6. 風險評估（v0.2 依 35 號裁示改寫）

> 編號說明（研發部 (a)；v0.3 沿用）：本節 R1～R10 是**風險**編號（v0.4 依 QA-S9 更正），R1＝setup-go 快取。**go.mod 寫法不屬本節**，只在 §1.2 說明、§9.1 第 1 項記錄裁示；v0.1 §9.4（待確認）第 2 項已因裁示結案移除。對照表見 §10。

| # | 風險 | 可能性／影響 | 緩解 |
|---|---|---|---|
| R1 | **setup-go 快取**：key 含 versionSpec；改 1.27.3 後可能未命中，job 時間變長；PR ref 的快取不給 main 用 | 高／低（只影響時間；研發部同意，`impact/研發部.md` 第 32 行） | 預期中；判準 L2 只看前綴與雜湊段，命中與否只記錄；不改 `cache` 設定 |
| R2 | **工具鏈自動下載（GOTOOLCHAIN）**：CI 是 `GOTOOLCHAIN='auto'`；若 go.mod 與三處 setup-go 不同步，job 內會默默下載 go1.27.3（研發部 (b)，`rd-evidence/41`） | 低／中 | **S1＋研發條件 2：三處 setup-go 與 go.mod 同一個 commit**；L3 在三個 job 全 log 計數須為 0。根本處理（workflow 設 `GOTOOLCHAIN: local`＋斷言步驟）**另開 CR，登錄 SCR-053（資訊級），最遲在 2026-11-30 的每月工具檢視時決定是否提出**（35 號 §3 第 5 項） |
| R3 | **main ruleset strict=true**：main 前進時 PR 要 update，head 變動，四綠、§5.2 安全部重跑、安全部確認、審查部「可合併」都要重做 | 中／低 | 推送前複查 main（§9.4 ⑥ 停止條件）；head 變動即回到 ⑦ |
| R4 | **新工具鏈改變 gosec 分析行為**（#1774 前例） | 低／中 | 本機 V4、V10；L5 Files／Lines 必須與基準相同；**合併後 7 日內做 G404 金絲雀（S4），推分支前要使用者核准**（35 號 §3 第 3 項） |
| R5 | **Go 發行檔來源**：setup-go 先從 actions/go-versions 下載，不驗 sha256 | 低／中 | 既有風險，本 CR 不改變；L1 確認版本 |
| R6 | **研發本機環境**：1.27.2＋local 會失敗（不會假綠，§1.2）；auto 需連 proxy.golang.org | 中／低 | 合併前通知各工作樹；合併後新增 `env-go1273.sh`，`env-go1272.sh` 保留供回滾並加註（§4.3，(c)）；研發部本機工具建議以 1.27.3 重建（(d)） |
| R7 | **govulncheck DB 更新**使結果改變（與本 CR 無關但可能同時發生）；CI 預設輸出不印 ID，無法從 CI log 判斷新增項目是不是 GO-2026-6809 | 低／中 | **與 L6 一致（v0.3 依 QA-M7）**：出現 GO-2026-6809 以外的項目時，L6 直接判不通過，轉安全部依 03 號 §3 門檻裁定，不得為了讓本 CR 過而放寬（35 號 v0.1.1 §6.1(b) 第 69 行）。安全部以 `-show verbose` 確認 ID；裁定只有兩種：不屬於本 CR 造成，記錄後依門檻處理；或要擋就擋。失敗時另比對同時間 main 的結果供安全部參考 |
| R8 | **範圍外發現**：gosec 依賴 x/net v0.61.0 的 GO-2026-6866、GO-2026-6860 → **SCR-054**；x/crypto v0.58.0 的 GO-2026-5932 → **SCR-055**；均為資訊級、module 層未呼叫 | —／資訊 | **不在 CR-006 處理、不開 RR**，納入 13 號 C7 之後的每月檢視（35 號 §4） |
| R9 | **GO-2026-6809 本身**：情境嚴重度**低**（一般評等：中），14 日規則**不適用、不阻擋**；但 S1 要求最遲 2026-11-20 合併 | 低／低 | 依 §9.4 時程，於 2026-11-20 前合併 |
| R10 | **合併方式為 bypass**：使用者本人合併且無其他核准者時，main-pr-review 的 pull_request 規則預期為 bypass（RepositoryRole 5，`bypass_mode: pull_request`；#28、#29 前例，QA-S3） | 高（預期發生）／低（RR-011 已由 D-17 接受） | D-09 合併前人工核對 head＝安全部確認與審查部「可合併」所載 SHA；品保以 QA-CR006-15 記錄 rule suite（required_status_checks 必須為 pass）；是否接受由審查部／使用者判定 |

## 7. 回滾方式

- 觸發條件：PR 或合併後任一 required check 因 1.27.3 失敗且無法在期限內處理，或安全部判定不可接受。
- 做法：開一個標 `CR-006 回滾` 的 revert PR（等同 `git revert <CR-006 merge commit>`），把兩個檔在**同一個 commit**還原：
  - `.github/workflows/ci-devsecops.yml` 回到 sha256 **`a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f`**；
  - `go.mod` 回到 sha256 **`31c3fd2790509a6ce56c44831b5a609746c1a39656ad95e2b4892ae98f6bbb49`**（`go 1.27.2`）。
  - 回滾 PR 同樣要四綠，走 D-08／D-09；推送須使用者核准。
- 本機：改用保留的 `/workspace/tools/env-go1272.sh`（§4.3）。
- **S7（35 號 §5 原文）**：「如果回滾，RR-012 恢復有效，CI 回到 1.27.2，安全部在 7 日內重訂期限。禁止 CR §7 列出的各種繞過方式」。
- **禁止**：改 `1.27.x`／`stable`／`check-latest: true`、設 `GOTOOLCHAIN` 繞過、放寬門檻、`continue-on-error`、移除 required check、加 `-exclude`／`#nosec`。

## 8. 附件／參考

`/workspace/rebuild/cr-006/`（完整 sha256 見同目錄 `SHA256SUMS`）：

| 檔案 | 內容 | sha256 |
|---|---|---|
| CR-006-ci-go-1.27.3-v0.1.md | v0.1（保留不動） | e4722b6f4cbaae711c91c23a0425bb7895483209ab3de88a81b87fae2369e2a8 |
| CR-006-ci-go-1.27.3-v0.2.md | v0.2（保留不動） | 5d233a71da15d9385b4eb035ad9f5979cbfed3f397d0c2bb5f1fc2e21124e83b |
| CR-006-v0.1-to-v0.2.diff | v0.1→v0.2 差異 | a067df9bc3ed208eede9810b8d7ad9e58a3af372a322a368e85b61b0dde59fe1 |
| CR-006-ci-go-1.27.3-v0.3.md | v0.3（保留不動） | e4d9ca21e9193f2bbde9fcf14dada29982e7cbe96259db9099a1121e4edf4fa8 |
| CR-006-v0.2-to-v0.3.diff | v0.2→v0.3 差異 | cad99bf815b6066e81e8903a129f7295c6f960d2fd8568663379f264a96251a8 |
| CR-006-v0.3-to-v0.4.diff | v0.3→v0.4 差異 | 見 `SHA256SUMS` |
| CR-006-ci-go-1.27.3.patch | 擬議 patch（v0.2 未改） | c3939e62c4a3a9081332fbdd7fd06bb056085bd40f732fe97b25acbcac56eaec |
| ci-devsecops.proposed.yml | 套用後 workflow（v0.2 未改） | 06833a314722d4bc7b436533a6268b2be29d5c56d028eeb1bfd20f9151def245 |
| ci-devsecops.main-ca09137.yml | main 原文（回滾值） | a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f |
| patch-validation.txt | `git apply --check`、numstat、YAML 比對 | e5dcb4e3c5f592d1f531019fcf709f225b0af35aba6a92756c4e458033dcbdf9 |
| evidence/SHA256SUMS | v0.1 證據 19 檔（未改） | 41cbbce966f477edfe4da419f58bcf27ef3f53958a01edc8a1713c8057934c33 |
| evidence-v0.2/17-v0.2-baseline-recheck.txt | v0.2 基準複查、sast-gosec L3 基準計數、影響分析檔核對（當時版本） | d6404be851d4303e1c4f164367edbe0cf6e1bc04a36da98b406c8ba0aea6dd3a |
| evidence-v0.3/18-v0.3-baseline-recheck.txt | v0.3 基準複查、公開／未登入存取確認、影響分析與 35 號核對（當時版本） | 007a76513896fa0418749447bb134358c840e632303e3dc74a21eb86ae1df417 |
| evidence-v0.4/19-v0.4-baseline-recheck.txt | v0.4 基準複查、品保 v0.3 與 35 號 v0.1.2 核對 | 見 `SHA256SUMS` |
| impact/企劃部.md | 企劃影響評估 v0.2 | 9387362bdd3d20d96be638d7615ca7431ab60833416b10a5869191cfc0b9f1c7 |
| impact/設計部.md | 設計影響評估 v2 | 32ef61411cc7ba94396027cb3b244cf783cf4acf8db2ee77cbad6bb0e5aa845b |
| impact/研發部.md | 研發影響評估 | 5275dd2cf436100ef9ee4f875aa4c1a2ff8c192aab74e88e4e6053f4a9b944af |
| impact/rd-evidence/SHA256SUMS | 研發證據 32 檔 | 0b97a8018d10e23d5a408eb3f6d1a7144e59855a73ea08ec3ffb7b025cfcf6b5 |
| impact/安全部.md | 安全影響意見 | 1dc296b2f4e02cff4e30ec20b1f3234216660edff523d5ce59be0a5a758c4432 |
| impact/品保部.md | 品保影響評估 v0.3 | 5b8869cbb5f6de401e8ea747453cc29f88f38e418f6351d8ab97594253d621ba |
| impact/qa-evidence/SHA256SUMS | 品保證據 v0.1，11 檔 | 748b8d57996875827133d3dc17da2c444bc85169ab3c17b4e67b0ed3790664c5 |
| impact/qa-evidence-v0.2/SHA256SUMS | 品保證據 v0.2，7 檔 | 5199c03ee3b987dc926e3ec680fc314b7dbc59fcf31f3725202299c4615faf09 |
| impact/qa-evidence-v0.3/SHA256SUMS | 品保證據 v0.3，3 檔（14、15、16） | 14b88a4400813f3afacf2ebdb87f73329e0e54a79efac24925fe4be56712d159 |

引用的專案文件（讀自 `/workspace/projects/short-url/docs/`；除 35 號外 sha256 與 origin/main ca09137 相同）：

| 檔案 | 版本 | 引用行 | sha256 |
|---|---|---|---|
| 06-security/35-cr-006-go-1273-security-opinion-v0.1.md（**origin/main 上尚無**；檔名未改，內容為 v0.1.2） | v0.1.2 | 6、14、27～31、37～41、44～52、61、64、69、72～85（§6.2）、87～93（§7） | dc94d52b5c7475a21c78d6fd377c3a3725bf9d836bcd6060141b5f3b8895df0a |
| 06-security/30-cr-005-gosec-v2.30.0-security-opinion-v0.1.md | v0.1.1 | 46～51、61 | e4579d31ab6ae9ba883b5a2e84f21d2da22b2b456f658127bca15a893a598104 |
| 06-security/05-residual-risk-register-v0.1.md | v0.1.13 | 32（RR-011）、33（RR-012） | df8636d892cc9a3faf85c1fc01a6db08a6c87a2855face0ffca346090ad9f7c2 |
| 06-security/13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md | v0.1.3 | 108、115 | 32f7f8c29f21d8fbbe244d6b38d4330b8c21663a5ebe5b04658245543976343e |
| 06-security/03-ci-security-gates-v0.3.md | v0.3.9 | 52 | 624f08d73b5d162e922c4b95fda62bf27257607330d7cf82044df0873882e25d |
| 09-cicd-env/01-ci-pipeline-v0.2.md | v0.2.3 | 7、43～55 | 0ca19522a0d02e3946b4d108d1702468f7d7f940916745caac52eb4733954904 |

---

## 9. 建議與核准

- 風險評估：見 §6。安全部判定風險等級低（35 號 §1）。
- 回滾方式：見 §7。
- 維運建議：**接受**。理由：滿足 30 號 K6／13 號 C8；修正 RR-012；改動 5 行、2 檔且同一 commit；go.sum 與依賴不變；五份影響分析無阻擋項；安全部同意（附條件 S1～S7）。（本機實測僅供參考，以 PR CI 為準。）
- 核准者／日期／結論：待填（各部門對 v0.4 確認 → 審查部 → 使用者拍板）

| 角色 | 狀態 | 日期 | 備註 |
|---|---|---|---|
| 維運部（提出） | 草稿 v0.4 | 2026-10-28 | 本檔 |
| 企劃部（影響） | 範圍無影響；時程無新增（期限 K6、S1、S4、S6 皆為條件，不是里程碑） | 2026-10-28 14:09 | `impact/企劃部.md` v0.2（9387362bdd3d20d96be638d7615ca7431ab60833416b10a5869191cfc0b9f1c7）**對 v0.3 確認沿用**，2026-10-28 14:09（UTC+8）書面確認（經總協調轉達；檔案本身未改版） |
| 設計部（影響） | 無影響 | 2026-10-28 14:09 | `impact/設計部.md` v2（32ef61411cc7ba94396027cb3b244cf783cf4acf8db2ee77cbad6bb0e5aa845b）**對 v0.3 確認沿用**，2026-10-28 14:09（UTC+8）書面確認（經總協調轉達；檔案本身未改版） |
| 研發部（影響） | 低影響，程式碼工作量 0；方案甲有條件可行 | 2026-10-28 | `impact/研發部.md`（5275dd2cf436100ef9ee4f875aa4c1a2ff8c192aab74e88e4e6053f4a9b944af）沿用；條款只在技術內容改變時重新確認（第 129 行），v0.2～v0.4 技術內容未改 |
| 品保部（影響、驗收條件） | CR PR 內無影響；v0.3 的 QA-M1、M5～M9、S5、S6 改法全部放行；新增必改 QA-M10，v0.3 不放行（v0.4 已改，§5.1 L6、§9.3） | 2026-10-28 | `impact/品保部.md` v0.3（5b8869cbb5f6de401e8ea747453cc29f88f38e418f6351d8ab97594253d621ba）；綁定 CR v0.3 |
| 安全部（評估） | 同意（附條件 S1～S7） | 2026-10-28 14:09 | **35 號 v0.1.2（dc94d52b5c7475a21c78d6fd377c3a3725bf9d836bcd6060141b5f3b8895df0a）§6.2 對 CR v0.3 確認**：第 78 行「§1 的結論、§3 的裁示、§5 的 S1～S7、§6.1、§7 對 v0.3 全部沿用」；`impact/安全部.md`（1dc296b2f4e02cff4e30ec20b1f3234216660edff523d5ce59be0a5a758c4432）仍有效、不改版（同行）。**注意**：原文**沒有**「唯一一次書面確認，之後一律沿用」的寫法；§6.2 效力條款（第 85 行）寫「本節只對 CR-006 v0.3（e4d9ca21…）、patch c3939e62…、proposed.yml 06833a31… 有效，重新確認的條件和 §6 相同」，即依 §6 第 64 行：只改 CR 文字敘述時，安全部比對差異檔後以小版修訂確認。v0.4 依原文為準，列為待安全部比對 v0.3→v0.4 差異（§9.5 第 1 項） |
| 審查部（確認） | 待審查 | | |
| 使用者（拍板） | 待拍板 | | 提出與合併期限皆為 2026-11-20 |

### 9.1 安全部裁示與條件（35 號 v0.1.2，sha256 dc94d52b5c7475a21c78d6fd377c3a3725bf9d836bcd6060141b5f3b8895df0a；§1、§3、§5 與 v0.1 相同，35 號 §6 第 61 行、§6.2 第 78 行）

**(一) 對 v0.1 §9.1 各項的裁示**（35 號 §3，摘要；原文見該檔第 27～32 行）

| # | 事項 | 裁示 | 本 CR 落實位置 |
|---|---|---|---|
| 1 | go.mod 寫法 | **採甲：只改 go 行 `go 1.27.3`，不加 toolchain 行**；乙、丙不採用 | §1.2（唯一說明處）、§1.3（patch 不變） |
| 2 | 建置版本斷言步驟 | **這次不加**；改用 **L1、L3 加上本機重現**作為證據；CI log 必須符合 S2 | §5.1（不加斷言說明、L1＋L3 合併判讀）、§5.2（安全部 PR head 獨立重跑）、§1.4 |
| 3 | 合併後 G404 金絲雀 | **要做**：合併後 7 日內，**推分支前要使用者核准**（S4） | §6 R4、§9.4 ⑬ |
| 4 | GO-2026-6809 嚴重度與 14 日規則 | **情境嚴重度低**（一般評等：中）；**14 日規則不適用、不阻擋**；但 S1 要求**最遲 2026-11-20 合併** | §3、§6 R9、檔頭期望完成時間 |
| 5 | GOTOOLCHAIN | **另開 CR，不放在 CR-006**；登錄 **SCR-053**（資訊級），**最遲在 2026-11-30 的每月工具檢視時決定要不要提出**（內容：`GOTOOLCHAIN: local`＋第 2 項的斷言步驟） | §6 R2、§9.5 |
| — | 範圍外發現 | **SCR-054**：x/net v0.61.0 的 GO-2026-6866、GO-2026-6860（修正於 v0.62.0）；**SCR-055**：x/crypto v0.58.0 的 GO-2026-5932（無修正版）；皆資訊級，**不開 RR**，納入每月檢視 | §6 R8 |

**(二) 研發部條件 2（對應 §6 R2）**：三處 setup-go 要和 go.mod **放在同一個 commit** 改（`impact/研發部.md` 第 54 行「條件 2」；與 35 號 S1「而且在同一個 commit 裡」一致）。落實：§1.1、§6 R2、§9.4 ⑥（推送前以 `git show --stat HEAD` 確認單一 commit 含兩檔）、QA-CR006-14。

**(三) S1～S7 原文（逐字引用 35 號 v0.1.2 §5，第 44～52 行，與 v0.1 逐字相同）與落實方式**

> | # | 條件 | 時點 | 負責 |
> |---|---|---|---|
> | S1 | PR 只能改 `.github/workflows/ci-devsecops.yml` 和 `go.mod` 兩個檔，內容必須和 proposed.yml（06833a31…）及擬議 go.mod（02842cf5…）完全一致，而且在同一個 commit 裡。最遲 2026-11-20 合併 | PR | 維運 |
> | S2 | PR 最新 head 的四項 required checks 都要綠。sast-gosec 的 log 要符合 CR §5.1 的 L1、L3、L4、L5，其中 L3 是完全沒有 `switching to go`；sca-govulncheck 和 unit-test 的 L1 要是 `go1.27.3`，govulncheck 要 exit 0。安全部會對這個 head 出具書面確認，head 一變就要重做 | 合併前 | 維運（取證）／安全 |
> | S3 | 合併後 main push run 四項全綠，log 印出 `go1.27.3` | 合併後 | 維運 |
> | S4 | G404 金絲雀：比照 K4（用 workflow_dispatch 或不合併的分支，樣本沿用 CR-005），在合併後 7 日內完成。推金絲雀分支前要先經使用者核准；拿到 sast-gosec 紅燈（G404 HIGH、Issues 1、exit 1）的 run 和 log 後，刪除分支並留證 | 合併後 7 日內 | 維運／安全覆核 |
> | S5 | S3 和 S4 完成後，安全部改版這幾份文件：03 號 §2 的 Go 列改成 1.27.3；05 號 RR-012 結案；13 號 C8 補註「1.27.3 由 CR-006 執行」 | S4 之後 | 安全 |
> | S6 | G5 部署清單要加上「正式建置用的 Go 版本與 CI 相同（1.27.3）」，並附證據 | G5 前 | 維運 |
> | S7 | 如果回滾，RR-012 恢復有效，CI 回到 1.27.2，安全部在 7 日內重訂期限。禁止 CR §7 列出的各種繞過方式 | 回滾時 | 維運／安全 |

| 條件 | 本 CR 怎麼滿足 | 證據 | 現況 |
|---|---|---|---|
| S1 | 推送內容＝proposed.yml（06833a31…）＋擬議 go.mod（02842cf5…），**單一 commit** 只含這兩檔（§9.4 ⑥）；排程使合併不晚於 2026-11-20 | `git show --stat <commit>`、兩檔 blob sha256、`git diff --numstat <merge-base> <head>`（QA-CR006-14） | 未執行 |
| S2 | PR head 四綠；sast-gosec 符合 L1、L3（`switching to go` 0 筆）、L4、L5；sca-govulncheck、unit-test 的 L1 為 go1.27.3、govulncheck exit 0（v0.2 另加 L6、L7）；安全部依 §5.2（35 號 §7）對同一 head 獨立重跑 M5-1～M5-8，結果與書面確認綁同一 head 一起出具，head 變動即重做（§9.4 ⑦～⑨） | run 連結、L1～L7 擷取（維運，`docs/09-cicd-env/evidence/cr-006/`）、M5 證據（安全部，`docs/06-security/evidence/cr-006-m5/<head SHA>/`）、安全部書面確認 | 未執行 |
| S3 | 合併後 main push run 四綠，L1～L7 也成立（QA-S1）（§9.4 ⑫） | main run 連結與擷取 | 未執行 |
| S4 | 合併後 7 日內，樣本沿用 CR-005（`evidence/cr-005/11`），**推金絲雀分支前取得使用者核准**；取得紅燈（G404 HIGH、Issues 1、exit 1）run 與 log 後刪分支並留證（§9.4 ⑬） | 紅燈 run、log、分支已刪證明 | 未執行 |
| S5 | 由安全部於 S3、S4 後改版 03、05、13 號（§4.3） | 安全部改版紀錄 | 未執行（安全部） |
| S6 | 維運於 G5 前在 G5 部署清單加「正式建置用的 Go 版本與 CI 相同（1.27.3）」並附證據（§4.3） | G5 部署清單改版 | 未執行（G5 前） |
| S7 | §7 回滾段已寫入原文；回滾時通知安全部 7 日內重訂期限 | 回滾 PR、安全部重訂紀錄 | 僅回滾時適用 |

**(四) 35 號 v0.1.1／v0.1.2 新增 §6、§6.2、§7 的落實**：§6 效力條款（第 64 行）→ §4.1 效力段、§9.5 第 1 項；§6.1(a)（第 68 行）→ §1.2（不會假綠）；§6.1(b)（第 69 行）→ §5.1 L6、§6 R7；§6.2（第 72～85 行）→ §9 簽核表、§9.5；§7（第 87～93 行）→ §5.2、§4.3、§9.4 ⑧⑨、QA-CR006-10(b)(c)。

### 9.2 品保部必改項 QA-M1～M10 的處理（品保部 v0.3 §8.1 第 207～222 行：M1、M5～M9 全部放行；§8.4 第 258 行 QA-M10；v0.1 §3.1 第 85～91 行；v0.2 §3.1 第 97～103 行、§3.4 第 137～142 行）

| # | 必改內容（摘要） | 品保部對 v0.2 的結論 | v0.3 處理 | 位置 |
|---|---|---|---|---|
| QA-M1 | sca-govulncheck、unit-test 各加編號判準，要能區分變更前後 | 部分退回：L7 放行；L6 退回（見 M6、M7） | L7 不變；L6 依 M6、M7 修正 | §5.1 |
| QA-M2 | L2 不用「第一次」 | 放行 | 不變 | §5.1 L2 |
| QA-M3 | L3 可 grep、涵蓋全部步驟 | 放行（補充見 QA-S6） | 不變；加 L1＋L3 合併判讀（QA-S6） | §5.1 L3 與表後說明 |
| QA-M4 | L5 要有基準 | 放行 | 不變 | §5.1 L5 |
| QA-M5 | 重跑執行者與判準 | 有條件（接受安全部獨立重跑，列 6 項前提，第 126～132 行） | **改由安全部獨立重跑、品保核對**（35 號 v0.1.1 §7）；6 項前提逐項寫入：執行者／核對者、head 綁定與重跑、指令（v0.2 第 1～7 步＋§7 兩點）、每步預期輸出與兩個主要指令 exit 0、證據位置、合併前置條件 | §5.2、§9.4 ⑧⑨ |
| QA-M6 | L6 刪「`Your code is affected by 0 vulnerabilities.`」 | 新必改 | **已刪**，並註明不要求的理由 | §5.1 L6 |
| QA-M7 | L6／R7 一致；收錄 35 號 §6.1(b) | 新必改 | **已收錄**：L6 逐字寫入「出現 GO-2026-6809 以外的項目時，L6 直接判不通過…」並附 35 號第 69 行原文；R7 改為一致 | §5.1 L6、§6 R7 |
| QA-M8 | §5.2、§4.3、§9.4 ⑧⑨、§9.5 第 3、4 項、§8 的執行者、證據位置、exit code、35 號版本 | 新必改 | **已改**：執行者安全部、品保核對；證據 `docs/06-security/evidence/cr-006-m5/<head SHA>/`；M5-5、M5-6 寫明 exit 0 與非 6809 項目處理；§8 改引用 35 號 v0.1.1；QA-CR006-10(b)(c) 執行者改安全部 | §5.2、§4.3、§8、§9.3、§9.4 ⑧⑨、§9.5 |
| QA-M10 | L6 與 QA-CR006-10(a) 的「開頭」字樣在 CI log 上會漏判 | 品保部 v0.3 §8.4 新必改（第 258 行），v0.3 不放行 | **已改（v0.4）**：§5.1 L6 改為 `grep -c 'No vulnerabilities found\.'` ≥ 1、`grep -c 'This scan also found'`＝0（不加 `^`）並寫明理由；§9.3 改逐字引用品保部 v0.3 §2（10(a) 已同步）；全文檢查 L6 已無「開頭」或 `^` 寫法 | §5.1 L6、§9.3 |
| QA-M9 | §9.3 改逐字引用品保部 v0.2 §2；對應表與原文一致；重產 SHA256SUMS | 新必改 | **已改**：§9.3 逐字引用 v0.2 第 58～89 行；證據對應表依第 89 行對應說明改寫；SHA256SUMS 重產 | §9.3、`SHA256SUMS` |

### 9.3 驗收條件（逐字引用品保部 `impact/品保部.md` **v0.3** §2，第 59～92 行，sha256 5b8869cbb5f6de401e8ea747453cc29f88f38e418f6351d8ab97594253d621ba；對應 QA-M9、QA-M10）

> ## 2. 驗收條件（QA-CR006；與 L1～L7 對應）
>
> 共通規則：每一條都要在 **PR 最新 head 的 run** 與 **合併後 main push run** 各取證一次（下稱「兩個 run」）；head 一變就重取（D-10）。本機結果不得替代 CI 結果。
>
> v0.3 只改兩處（其餘與 v0.2 逐字相同）：(i) QA-CR006-1 的反向字樣改為與 CR v0.3 L1（QA-S5）相同；(ii) QA-CR006-10(a) 的「開頭」改為不錨定的計數寫法（QA-M10）。兩處都是判準文字，CR §9.3 須重新引用（見 §8.4）。
>
> v0.2 只改下列四處，其餘條文與 v0.1 逐字相同：(1) 對應欄依 CR v0.2 新增的 L6、L7 調整；(2) QA-CR006-10 的 (a) 依 QA-M6 修正、(c) 的章節引用改為 CR §5.2；(3) QA-CR006-14 加入「單一 commit」（採納 CR v0.2 §9.3 的補充，與 35 號 S1、研發條件 2 一致），停止推送的引用改為 CR §9.4 ⑥；(4) 標題與對應說明。
>
> ### 2.1 對應 L1～L7
>
> | ID | 對應 | 條件 | 證據 | 通過標準 |
> |---|---|---|---|---|
> | QA-CR006-1 | L1 | 三個有 setup-go 的 job（sast-gosec、sca-govulncheck、unit-test）都跑在 Go 1.27.3 | 兩個 run 三個 job 的 log 片段 | 每個 job 都有 `Setup go version spec 1.27.3`、`go version go1.27.3 linux/amd64`、`GOVERSION='go1.27.3'`、`GOROOT='/opt/hostedtoolcache/go/1.27.3/x64'`；任一 job 的 log 以 `grep -E 'go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/'` 命中即不通過（v0.3 與 CR L1 同步；縮小後對基準 run 每個 job 仍命中 14 處，`qa-evidence-v0.3/16`） |
> | QA-CR006-2 | L2 | 快取 key 改為 1.27.3 且 go.sum 雜湊不變 | 三個 job 的 `Cache hit for:`／`Cache is not found`／`Cache saved with key:` 行 | key 為 `setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d429736099b597565b998fd2f9b82c99c65d9cffc4d8c9270df02c0a025`（雜湊段＝基準，因 go.sum 不變）；命中或未命中**只記錄，不作判定** |
> | QA-CR006-3 | L3 | 沒有工具鏈自動切換或下載 | 兩個 run 三個 job **全 log** 以 `grep -cE 'switching to go\|downloading go1\.'` 計數 | 0 筆（基準亦為 0） |
> | QA-CR006-4 | L4 | gosec 原始碼仍是 v2.30.0 | sast-gosec「Verify gosec module hash (CR-002)」步驟 | `mod	github.com/securego/gosec/v2	v2.30.0	h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=` 與 `gosec h1 OK` 皆在；**不以 `gosec -version` 為證**，也**不以 h1 證明建置用的 Go 版本** |
> | QA-CR006-5 | L5 | gosec 指令與發現前後對照 | 變更前（PR merge-base 那個 main commit 的 push run）與變更後（兩個 run）的命令列與 Summary 並列 | 命令列仍為 `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...`；Issues 0、Nosec 0、exit 0；Files／Lines 與 merge-base 基準相同（ca09137 為 57／7485）。**任何差異（增或減）都記錄並轉安全部判定** |
> | QA-CR006-7 | L7 | 全部單元測試（含 -race）結果不變 | 兩個 run 的 artifact `go-test-json`（`gh run download -n go-test-json`），以 `jq` 統計；unit-test log 的 ENG-011 步驟 | 0 fail、0 skip；pass 數＝merge-base 基準（ca09137 為 **181**，含子測試；若期間 main 有新測試，以當時 merge-base 的 main run 為準）；`--- PASS: TestENG011_PerSourceIsolationByPeerAddress` 存在；CI 指令仍含 `-race` |
> | QA-CR006-8 | L7 | 14 個 CT 測試全過 | 同上 artifact 以 `jq` 篩 `^TestCT` 的 pass 事件 | TestCT_* pass **14**（6 個頂層＋ERR-01 的 8 個子測試），0 fail |
> | QA-CR006-10 | L6＋CR §5.2 | govulncheck 不再報 GO-2026-6809 | (a) 兩個 run 的 sca-govulncheck log；(b) 在 PR head 以 Go 1.27.3 本機跑 `govulncheck -show verbose ./...`；(c) 依 CR §5.2 對同版建置的 gosec 二進位跑 `govulncheck -mode=binary`（RR-012 所指對象）。(b)(c) 的執行者依 35 號 v0.1.1 §7 為安全部、品保核對（待 CR 改版收錄，見 §3.3） | (a) 有 `No vulnerabilities found.`、`grep -c 'This scan also found'`（**不加 `^` 錨點**；CI log 每行前面有 job 名稱與時間戳）為 0、同 job `GOVERSION='go1.27.3'`、exit 0；**不要求** `Your code is affected by 0 vulnerabilities.`（三層皆 0 時 govulncheck 不印這行，見 §3.2）；(b) 印出「the go1.27.3 standard library」且無 `GO-2026-6809`；(c) 無 `GO-2026-6809`、exit 0。RR-012 是否結案由安全部判定 |
> | QA-CR006-11 | L6 | SCA findings 前後對照 | 基準（merge-base main run）與兩個 run 的 govulncheck Symbol／package／module 三層計數並列；PR 上 dependency-review 結果 | 預期由「0／1／0」變「0／0／0」；**除 GO-2026-6809 消失以外的任何差異**都記錄並轉安全部；dependency-review 不得出現 Critical |
>
> ### 2.2 L1～L7 對不上的條件（獨立列出）
>
> | ID | 條件 | 證據 | 通過標準 |
> |---|---|---|---|
> | QA-CR006-6 | 四個必要檢查全綠，名稱不變 | `gh pr checks <PR>`、`gh run view <run> --json jobs`（兩個 run） | secrets-gitleaks、sast-gosec、sca-govulncheck、unit-test 四項 success；名稱與 ruleset 24629805 的 required contexts 完全相同；不得有 `continue-on-error` 或新增 skip |
> | QA-CR006-9 | 覆蓋率不低於現值 | CI 不產生覆蓋率，由品保在 PR head 與 merge commit 以 Go 1.27.3、`GOTOOLCHAIN=local` 實跑 `go test ./... -count=1 -coverpkg=./... -coverprofile=cov.out && go tool cover -func=cov.out \| tail -1` | 合計 ≥ **80.9%**（同一算法；不得拿不加 `-coverpkg` 的 78.7% 比較）；各套件值並列 |
> | QA-CR006-12 | go.mod 與 setup-go 一致 | PR head 與 merge commit 的 `go.mod`、workflow 三處 `go-version`；`git ls-tree -r --name-only` 確認仍無 Dockerfile | 方案甲：go.mod 第 3 行 `go 1.27.3`、無 `toolchain` 行；三處 `go-version: "1.27.3"`、`check-latest: false` 不變。若安全部／研發部改採方案乙或丙，本條須依裁示重寫。Dockerfile：repo 無，**不適用**（如 PR 新增 Dockerfile 即屬範圍外變更，不通過） |
> | QA-CR006-13 | ruleset 未變 | 合併前（推送前）與合併後各取一次 `gh api repos/444112029012/short-url/rulesets/{24629805,24629854} \| jq -S 'del(._links, .current_user_can_bypass)'`，與 qa-evidence 08 基準 `diff` | diff 為空（**含 updated_at**，比 QA-CR005-5 更嚴；`_links`、`current_user_can_bypass` 依呼叫者而變，故排除） |
> | QA-CR006-14 | diff 範圍與 CR 一致 | `git diff --numstat <merge-base> <head>`；`git show --stat <commit>`；兩檔 sha256；合併後 `git diff --numstat <merge^1> <merge>` | 只有 `.github/workflows/ci-devsecops.yml`（4／4）與 `go.mod`（1／1），且兩檔在**同一個 commit**；兩檔內容＝`06833a31…f245`／`02842cf5…8ebf`（若 main 前進導致基準 workflow 不再是 `a67a0be4…`，依 CR §9.4 ⑥ 停止推送，本條不通過）；go.sum 不變（`1e2b43d7…`） |
> | QA-CR006-15 | 合併事實可追溯 | `gh pr view <PR> --json state,mergedAt,mergeCommit,mergedBy`；merge commit parent；合併那次 push 的 rule suite（`gh api …/rulesets/rule-suites?ref=refs/heads/main` 找 after_sha＝merge sha，再查 `…/rule-suites/{id}`） | merge sha、時間（換算 UTC+8）、合併方式、rule suite 的 result（pass 或 bypass）與各規則評估都記錄；required_status_checks 必須為 pass。若為 bypass，品保只記錄，是否接受由審查部／使用者判定 |
>
> > L1～L7 與 QA-CR006 的對應：L1→1、L2→2、L3→3、L4→4、L5→5、L6→10(a)、11、L7→7、8；QA-CR006-10(b)(c)→CR §5.2；QA-CR006-6、9、12、13、14、15 沒有對應的 L 判準（9 由品保本機實跑，12 對 CR §1.2 方案甲，14 對 35 號 S1）。CR v0.2 §9.3 維運整理的「取證者／時點」欄，品保同意沿用，唯 10(b)(c) 執行者改依 35 號 v0.1.1 §7。

**逐條對應證據（依品保部 v0.3 第 92 行對應說明；取證者分工品保部同意沿用，唯 10(b)(c) 改依 35 號 §7；證據編號 E1～E12 為品保部 v0.3 §6 原編號，第 180～193 行）**

| ID | 對應（品保部 v0.3 第 92 行） | 取證者／時點 | 證據 |
|---|---|---|---|
| QA-CR006-1 | L1 | 維運；PR head run、main push run | E2：三 job log 片段（Setup go version spec／go version／GOVERSION／GOROOT／GOTOOLCHAIN）與全 log sha256；反向字樣 `go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/` 的計數（v0.3 起與 L1 同步） |
| QA-CR006-2 | L2 | 維運；同上 | E2：快取 key 行 |
| QA-CR006-3 | L3 | 維運；同上 | E3：三 job 全 log 計數 |
| QA-CR006-4 | L4 | 維運；同上 | E4：Verify 步驟 |
| QA-CR006-5 | L5 | 維運；merge-base run＋兩個 run | E4：命令列與 Summary 並列 |
| QA-CR006-6 | 無 L 判準 | 維運；兩個 run | E1：`gh pr checks`、`gh run view --json jobs` |
| QA-CR006-7 | L7 | 品保；兩個 run | E6：`go-test-json` jq 統計、ENG-011 片段 |
| QA-CR006-8 | L7 | 品保；兩個 run | E6：TestCT_* 統計 |
| QA-CR006-9 | 無 L 判準（品保本機實跑） | 品保；PR head 與 merge commit 本機（Go 1.27.3、`GOTOOLCHAIN=local`） | E7：`-coverpkg=./...` 合計 ≥ 80.9% |
| QA-CR006-10 | (a)→L6；(b)(c)→CR §5.2 | (a) 維運（兩個 run）；**(b)(c) 安全部執行（§5.2 M5-6、M5-5）、品保核對**（35 號 v0.1.1 §7；QA-M8） | E5：sca-govulncheck log，(a) 以 `grep -c 'No vulnerabilities found\.'`、`grep -c 'This scan also found'`（不加 `^`）計數；E12：`docs/06-security/evidence/cr-006-m5/<head SHA>/` 的 M5-5、M5-6 輸出 |
| QA-CR006-11 | L6 | 品保；merge-base run＋兩個 run＋PR dependency-review | E5：三層計數並列（預期 0／1／0 → 0／0／0） |
| QA-CR006-12 | 無 L 判準（對 CR §1.2 方案甲） | 品保；PR head、merge commit | E8：go.mod、三處 go-version、`git ls-tree` |
| QA-CR006-13 | 無 L 判準 | 品保；推送前（§9.4 ⑥）、合併後（⑫） | E9：`jq -S` 快照 diff，對 `qa-evidence/08` |
| QA-CR006-14 | 無 L 判準（對 35 號 S1） | 品保；PR head、merge commit | E10：numstat、`git show --stat`（單一 commit）、兩檔與 go.sum sha256 |
| QA-CR006-15 | 無 L 判準（對 CR §6 R10） | 品保；合併後 | E11：merge sha、時間（UTC+8）、rule suite |

**品保建議項 QA-S1～S9（不阻擋）的處理**

| # | 建議（摘要） | 處理 | 位置 |
|---|---|---|---|
| QA-S1 | L1～L7 在合併後 main push run 也要成立 | **採納**（品保部 v0.3 第 150 行確認處理正確） | §5.1 標題、§9.1 S3 落實、§9.4 ⑫ |
| QA-S2 | L1 加記 `GOTOOLCHAIN='auto'` | **採納**（只記錄、不判定） | §5.1 L1 |
| QA-S3 | 預先寫明合併為 bypass 時如何處理 | **採納**（RR-011 已由 D-17 接受，05 號第 32 行） | §0、§6 R10、§9.4 ⑪ |
| QA-S4 | §4.2 測試列直接引用品保結論 | **採納**（v0.3 改引用 v0.2 §5） | §4.2「測試策略…」列 |
| QA-S5 | 縮小 L1 比對字樣為 `go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/` | **採納**（v0.3） | §5.1 L1 |
| QA-S6 | 寫明 L1、L3 合併判讀 | **採納**（v0.3） | §5.1 表後「L1 與 L3 必須合併判讀」 |
| QA-S7 | M5-5 exit 欄改「govulncheck exit 0 且 GO-2026-6809 計數 0」；module 層 3 筆只記錄；不得用 `-scan module`（品保部 v0.3 第 264 行） | **採納**（v0.4；引用 35 號 §6.2 第 7 項，第 83 行） | §5.2 M5-5 |
| QA-S8 | M5-1 兩種取法或與維運值不一致時「停止、記錄、重取」（第 265 行） | **採納**（v0.4；兩種取法都做） | §5.2 M5-1 |
| QA-S9 | §8 改引用 35 號 v0.1.2；§6 編號說明改 R1～R10（第 266 行） | **採納**（v0.4；另為一致，§2、§4.1、§4.3、§5.2、§9.1 的 35 號版本與 §7 行號一併改為 v0.1.2） | §8、§6、§2、§4.1、§4.3、§5.2、§9.1 |

### 9.4 執行步驟（比照 CR-005 §9.4）

| 步驟 | 內容 | 執行者 | 前置條件 | 產出 |
|---|---|---|---|---|
| ① | 各部門對 v0.3 重新確認影響分析（企劃 v0.2、設計 v2、品保 v0.2 綁 v0.2；研發、安全 impact 檔沿用） | 企劃／設計／研發／安全／品保 | 本 CR v0.3 送出 | 綁定 v0.3 sha256 的確認 |
| ② | 安全部評估（依 35 號 §6 第 64 行效力條款，確認 v0.3 的 QA-M5 落實方式符合 §7） | 安全部 | ① | 35 號小版修訂或確認 |
| ③ | 審查部審查（含 README 第 14 行裁示） | 審查部 | ② | 審查紀錄 |
| ④ | 使用者拍板 | 使用者 | ③ | 拍板紀錄（總協調登記） |
| ⑤ | 使用者在**維運部聊天裡親口同意推送** | 使用者 | ④ | 聊天中的同意訊息 |
| ⑥ | 推送前：`git fetch`，確認 main workflow 仍為 `a67a0be4…`、go.mod 仍為 `31c3fd27…`（不符即**停止推送**回報）；品保取 ruleset 推送前快照（QA-CR006-13）。從 origin/main 新建工作樹（提議 `/workspace/ops-cr006`），分支（提議）`ci/CR-006-go-1.27.3`；**以單一 commit 同時改 workflow 三處 setup-go＋第 6 行註解與 go.mod**（S1、研發條件 2），推送前核對兩檔 sha256＝`06833a31…`／`02842cf5…`、`git show --stat HEAD` 只有這兩檔；推送並開 PR，標 `CR-006` | 維運部 | ⑤ | PR 連結、head SHA、diff stat |
| ⑦ | PR 最新 head 四項 required checks 全綠，L1～L7 取證（存 `docs/09-cicd-env/evidence/cr-006/`）；**維運把 head SHA 提供給安全部** | CI（維運取證） | ⑥ | run 連結、log 擷取（QA-CR006-1～8、10(a)、11）、head SHA |
| ⑧ | **§5.2 PR head 獨立重跑**（QA-M5，35 號 §7）：M5-1～M5-8；證據存 `docs/06-security/evidence/cr-006-m5/<head SHA>/`＋安全部 SHA256SUMS；head 一變就重跑 | **安全部** | ⑦ | 重跑證據（M5-5、M5-6 皆 exit 0、GO-2026-6809 0 筆） |
| ⑨ | **品保部核對 ⑧ 證據**（記錄 SHA256SUMS 的 sha256、確認 head SHA 一致）；安全部依 S2 對**同一 head** 出具書面確認，與 ⑧ 結果一起出具（任一 M5 判準不符就不出具）；品保完成合併前的驗收條件（含 QA-CR006-9 PR head 覆蓋率、10(b)(c) 核對、12、14）。**⑧⑨ 都完成才可進入 ⑩（合併前置條件）** | 品保部（核對）、安全部（書面確認） | ⑧ | 品保核對紀錄；安全部書面確認（載明 head SHA） |
| ⑩ | 審查部對**同一 head** 說「可合併」；head 一變就失效，回到 ⑦ | 審查部 | ⑨ | 「可合併」紀錄 |
| ⑪ | **使用者本人合併**（D-09 人工核對 head＝⑨⑩ 所載 SHA；預期 bypass，見 §6 R10）；最遲 2026-11-20（S1） | 使用者 | ⑩ | merge commit |
| ⑫ | 合併後 main push run 四項全綠且 L1～L7 成立（S3、QA-S1）；品保取合併後 ruleset 快照與 QA-CR006-9、13～15 | CI（維運取證）、品保 | ⑪ | main run 連結；品保紀錄 |
| ⑬ | G404 金絲雀（S4）：合併後 7 日內；**推金絲雀分支前取得使用者核准**；取得紅燈 run 後刪分支 | 維運部（執行）、安全部（覆核）、使用者（核准推分支） | ⑫＋使用者核准 | 紅燈 run／log、分支已刪證明 |
| ⑭ | 更新 `docs/09-cicd-env/01-ci-pipeline-v0.2.md` §2.1，`docs/09-cicd-env/evidence/cr-006/` 一併入庫（文件變更，推送須使用者核准；`docs/06-security/evidence/cr-006-m5/` 的入庫由總協調安排）；本機新增 `env-go1273.sh`、`env-go1272.sh` 檔頭加註（§4.3） | 維運部 | ⑫ | 文件 PR；環境檔 |
| 後續 | 安全部 S5 改版 03／05／13 號（⑬ 之後）；品保 07-testing 補件；總協調 99-changelog；README 第 14 行（依審查部裁示）；S6 G5 部署清單（G5 前）；SCR-053 於 2026-11-30 每月檢視決定 | 各負責角色 | 各自時點 | — |

- ⑤ 之前維運部不推送、不開 PR、不改 GitHub 任何設定。
- 期限：**最遲 2026-11-20（UTC+8）提出（K6）並合併（S1）**。

### 9.5 待確認事項（v0.4 依各部門回覆與原文更新狀態）

| 項 | 內容 | v0.4 狀態 |
|---|---|---|
| 1 | 各部門對最新版的確認 | **大致解決**：企劃部 v0.2、設計部 v2 對 v0.3 確認沿用（14:09 書面確認，經總協調轉達）；研發部沿用；品保部 v0.3 綁 v0.3（QA-M10 v0.4 已改）；安全部 35 號 v0.1.2 §6.2（第 72～85 行）對 v0.3 確認全部沿用。**仍待**：依 35 號 §6.2 效力條款（第 85 行）與 §6 第 64 行，v0.4 屬文字修訂（patch、proposed.yml 未變），須安全部比對 `CR-006-v0.3-to-v0.4.diff` 後以小版修訂確認；品保部確認 QA-M10 改法；企劃、設計是否需對 v0.4 再確認，依其效力條款（企劃部第 19 行、設計部第 36 行）由兩部門判定 |
| 2 | `README.md` 第 14 行 | **未解決，待審查部裁示**（維運部建議合併後另外提） |
| 3 | QA-M5 執行者「維運自己驗自己」 | **已解決**（v0.3）：安全部獨立重跑、品保部核對；品保部 v0.3 §8.1 第 221 行「QA-M5 整體 … 放行」 |
| 4 | QA-CR006-10(b) 執行者 | **已解決**（v0.3）：安全部執行、品保核對（品保部 v0.3 第 78 行） |
| 5 | 35 號尚未入庫 origin/main | **未解決**：v0.4 複查時 origin/main 仍為 ca09137，該檔不在其上（`evidence-v0.4/19`）；入庫方式由總協調安排 |
| 6 | SCR-053 於 2026-11-30 每月檢視決定 | **維持**（未到期）；品保建議的 CI 步驟一併交 SCR-053 評估（品保部 v0.3 第 137 行） |
| 7 | `docs/` 下證據與 G5 清單入庫須使用者核准 | **維持** |
| 8 | M5 證據子目錄名 | **已解決**：採完整 40 碼 head SHA，`docs/06-security/evidence/cr-006-m5/<40 碼 head SHA>/`（35 號 §6.2 第 4 項，第 80 行；品保部 v0.3 §3.3 第 5 點第 134 行同步改為 40 碼，§8.3 第 1 點第 247 行接受） |
| 9 | M5-1 取得 head 的方式 | **已解決**：公開 ls-remote 或未登入 API，登入版 `gh pr view` 不採用，維運的值只供比對（35 號 §6.2 第 5 項，第 81 行；品保部 v0.3 第 131 行）；v0.4 依 QA-S8 改為兩種都做、不一致即停止、記錄、重取 |
| 10 | QA-M7 句子非逐字原文 | **已解決**：35 號 §6.2 第 3 項（第 79 行）「這句是總協調的轉述，不是本檔逐字原文，但意思和本檔 §6.1(b)（第 69 行）一致。CR 兩者並列、以本檔為準，安全部同意」；品保部 v0.3 §8.1 第 210 行放行 |
| 11 | QA-CR006-1 字樣與 L1 縮小字樣不一致（v0.3 回報） | **已解決**：品保部 v0.3 §2 第 71 行已同步為 `go1\.27\.2\|go-1\.27\.2\|spec 1\.27\.2\|/go/1\.27\.2/`（§8.3 第 3 點，第 249 行；§9.3 已逐字引用） |
| 12 | 安全部「唯一一次書面確認，之後一律沿用」的轉述 | **與原文不符，以原文為準**：35 號 §6.2 沒有這句；效力條款（第 85 行）要求依 §6 條件重新確認。請總協調向安全部確認是否要另外寫入 |

## 10. 修訂紀錄

| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-28 | 維運部 | 初稿：基準 ca09137；workflow 與 go.mod 位置盤點；Go 1.27.3／GO-2026-6809 查證；本機實測 V1～V10；patch 草稿（numstat 4/4、1/1）；影響分析只填維運列 |
| v0.2 | 2026-10-28 | 維運部 | 整合五份影響分析和安全部 35 號裁示。**技術內容不變**（patch c3939e62…、proposed.yml 06833a31… 未改）。逐處：(1) metadata 加變更紀錄、對應需求加 35 號；(2) §0 複查基準（13:39，仍 ca09137／a67a0be4／31c3fd27），加 ruleset 24629854 bypass 設定；(3) §1.1 加「三處 setup-go 與 go.mod 同一 commit」（研發條件 2、S1）；(4) §1.2 定案方案甲、刪乙丙選項，更正 gosec 安裝與掃描對 go.mod 的關係（研發 (b)）；(5) §1.4 註明不加斷言、不設 GOTOOLCHAIN；(6) §3 依 35 號改寫 14 日規則（不適用），加 S1 合併期限；(7) §4.1 改為七部門彙整表（檔名、完整 sha256、綁定版本、原文結論）＋各面向原文列，審查部「待審查」；§4.2 依各部門判定改寫，測試列引用品保 §4（QA-S4）；§4.3 加 S5、S6、evidence/cr-006、07-testing、env-go1273.sh／env-go1272.sh（研發 (c)）、研發本機工具（研發 (d)）、README 待審查部裁示；(8) §5.1 判準改 L1～L7（QA-M1～M4、QA-S1、QA-S2），刪可選斷言 C 改為不加說明；新增 §5.2 PR head 重跑（QA-M5）；原 §5.2 改為 §5.3；(9) §6 依 35 號改寫 R2、R4，新增 R8（SCR-054／055）、R9（嚴重度）、R10（bypass，QA-S3），加編號說明（研發 (a)）；(10) §7 加 S7 原文與本機回滾；(11) §9.1 改為安全部裁示表、研發條件 2、S1～S7 原文與落實；新增 §9.2 QA-M1～M5 處理；§9.3 改為品保 QA-CR006-1～15 原文＋證據對應＋QA-S1～S4 處理；§9.4 執行步驟 ①～⑭；§9.5 待確認重寫 |
| v0.3 | 2026-10-28 | 維運部 | 依五部門對 v0.2 的回覆修訂。**技術內容不變**（patch c3939e62…、proposed.yml 06833a31… 未改）。逐處：(1) metadata 版本 v0.3、變更紀錄加 v0.3、對應需求改 35 號 v0.1.1；開頭加 v0.3 說明；(2) 期望完成時間補 S4、S6 與企劃部「提出與合併同日、建議拍板時一併排定合併日」附註；(3) §0 複查（14:03，仍 ca09137／a67a0be4／31c3fd27；公開未登入可讀）；(4) §2、§3、§4.1、§4.2、§8、§9 簽核表的企劃部、設計部、品保部、35 號引用與行號改為新版；§4.1 企劃列換成 v0.2 原文；效力段重寫；(5) §4.3 證據列拆為維運 CI 取證（09-cicd-env）與安全部 M5（06-security/evidence/cr-006-m5/<head SHA>/）；(6) §5.1 L1 縮小比對字樣（QA-S5）、L6 刪 `Your code is affected…`（QA-M6）並寫入非 6809 項目處理（QA-M7，附 35 號第 69 行原文）、加 L1＋L3 合併判讀（QA-S6）；(7) §5.2 改為安全部獨立重跑，逐字引用 35 號 §7，判準 M5-1～M5-8 含預期輸出與 exit code（QA-M5、QA-M8）；(8) §6 R3、R7 改寫（R7 與 L6 一致）；(9) §9.1 標題改 35 號 v0.1.1、(二) 標題改「研發部條件 2（對應 §6 R2）」、S2 落實改安全部重跑、新增 (四) 35 號 §6、§7 落實；(10) §9.2 擴為 QA-M1～M9；(11) §9.3 改逐字引用品保部 v0.2 第 58～89 行，證據對應表依第 89 行改寫，10(b)(c) 執行者改安全部（QA-M9、QA-M8），建議項擴為 QA-S1～S6；(12) §9.4 ①②⑦⑧⑨⑭ 改寫；(13) §9.5 逐條更新狀態並新增第 8～10 項；(14) 新增 `evidence-v0.3/18`，重產 SHA256SUMS |
| v0.4 | 2026-10-28 | 維運部 | QA-M10 與 QA-S7～S9。**技術內容不變**（patch c3939e62…、proposed.yml 06833a31… 未改）。逐處：(1) metadata 版本 v0.4、變更紀錄、對應需求改 35 號 v0.1.2；開頭加 v0.4 說明；(2) §5.1 L6 改為 `grep -c 'No vulnerabilities found\.'` ≥ 1、`grep -c 'This scan also found'`＝0（不加 `^`）並寫明理由（QA-M10）；(3) §5.2 M5-1 兩種取法、不一致即停止記錄重取（QA-S8）；M5-5 exit 欄改「govulncheck exit 0 且 GO-2026-6809 計數 0」、不得用 `-scan module`、引用 35 號 §6.2 第 7 項（QA-S7）；(4) §6 編號說明改 R1～R10；§8 改引用 35 號 v0.1.2，並同步 §2、§4.1、§4.3、§5.2、§9.1 的 35 號版本與 §7 行號（72～78→87～93）（QA-S9）；(5) 品保部引用改 v0.3（5b8869cb…）與新行號：§4.1、§4.2、§5.1、§5.2、§8、§9.2、§9.3；(6) §9 簽核表寫入企劃、設計 14:09 對 v0.3 確認沿用、安全部 35 號 v0.1.2 §6.2、研發沿用、品保 v0.3；(7) §9.2 加 QA-M10；§9.3 改逐字引用品保部 v0.3 第 59～92 行，證據對應表與建議項 S7～S9 更新；(8) §9.5 逐條更新；(9) 新增 `evidence-v0.4/19`，重產 SHA256SUMS |

**v0.1 → v0.2 章節／編號對照（研發部 (a)）**

| v0.1 | v0.2 |
|---|---|
| §1.2 方案甲／乙／丙 | §1.2 只留甲（定案） |
| §5.1 C 可選斷言 | 刪除；§5.1「不加斷言步驟」說明 |
| §5.2 本機實測 | §5.3 |
| — | §5.2 PR head 重跑（QA-M5） |
| §9.1 待安全部裁示 1～6 | §9.1(一) 裁示表 1～5（go.mod 改列第 1 項，與 35 號 §3 第 1 項一致） |
| §9.2 執行步驟 | §9.4 |
| §9.3 驗收條件（預留） | §9.3（品保原文） |
| §9.4 待確認 1～8 | §9.5 待確認 1～7（原第 2 項 go.mod、第 4 項 14 日規則、第 6 項嚴重度、第 7 項裁示事項、第 8 項 x/net 均已由 35 號結案移除） |
| §6 R1～R8 | §6 R1～R10（R1 仍為快取；go.mod 不在 §6） |
