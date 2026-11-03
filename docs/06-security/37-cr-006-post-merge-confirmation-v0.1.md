---
文件：CR-006（CI 升 Go 1.27.3）合併後確認——35 號 S1～S7 狀態彙整
版本：v0.1.3
狀態：已提出（S1～S5 完成；S6 待到期（G5 前）；S7 未觸發；RR-012 已結案）
負責角色：安全部
專案代號：SHORTURL
對應：35 號 v0.1.3（sha256 1000dd9f151f1c33f1c7a2e7b5a4480e093ea1eee1c136a2ff495311b583f38e）§5 S1～S7；36 號 v0.1.1（`36-cr-006-pr31-m5-s2-confirmation-v0.1.1.md`，sha256 2936baec…d598）；CR-006 v0.4（`/workspace/rebuild/cr-006/CR-006-ci-go-1.27.3-v0.4.md`，sha256 818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734，2026-11-03 09:21 核對相符，只讀）
綁定：merge commit `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`（parents `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`＋PR head `474756381bd170beb3a5d723b63f1b8ba48b2add`）；main push run 41559091290
證據：`docs/06-security/evidence/cr-006-s4/`，SHA256SUMS sha256 `799617802037b294708ba02a35f9526137ca594abef74c148aaa42e1ac53a9ba`（12 檔，2026-11-03 09:46:24（UTC+8）`sha256sum -c` 全部 OK；v0.1／v0.1.1 時為 `5aeb30e6…20b6`，11 檔）
版本來源：v0.1.2 原文＝`37-cr-006-post-merge-confirmation-v0.1.2.md`（sha256 9cacd4682464d8e5e3e87e72e9e047919da012826acbb13d6de5a2f093b2d4e4）；v0.1、v0.1.1 未入庫，只有 sha256 紀錄（335f9591…、c495a6c5…）
撰寫時間：v0.1.3 開始於 2026-11-03 15:51:18（UTC+8，取自 `date`）；v0.1.2 開始於 2026-11-03 09:46:56（UTC+8，取自 `date`），完成時間以檔案 mtime 為準；v0.1.1 開始於 2026-11-03 09:29:43（UTC+8，取自 `date`），完成時間以檔案 mtime 為準；v0.1 開始於 2026-11-03 09:25:26（UTC+8，取自 `date`），完成時間以檔案 mtime 為準
---

# CR-006 合併後確認

## 1. 結論
**v0.1.2 更新**：S4 的 CI 金絲雀已完成並經安全部獨立核對（§4.1），S5 已依原文完成（03 v0.3.11、05 v0.1.15 RR-012 結案、13 v0.1.5）。以下 v0.1 原段落保留作歷史。
CR-006 已合併且 main CI 四綠（S3 成立）；修正已在 CI 生效。**S4 的 CI 紅燈部分需要推分支，未執行，改列為待使用者核准的執行計畫（§4），截止 2026-11-10 09:20:21（UTC+8）。** 因為 S5 要求 S3、S4 都完成，**RR-012 這次不結案**；03、05、13 號只出狀態註記小版。全程 GitHub 只讀，未修改 `/workspace/rebuild/**`，未做 exploit。

## 2. 合併事實（安全部自行查證）
| 項目 | 值 | 來源（時間 UTC+8，取自 `date`） |
|---|---|---|
| main | `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d` | `git ls-remote` 與未登入 API `commits/main` 一致（09:21:36；evidence 01 09:22:07） |
| parents | `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`、`474756381bd170beb3a5d723b63f1b8ba48b2add` | 同上；第二 parent＝36 號綁定的 head |
| 合併時間 | 2026-11-03 09:20:21（UTC+8；API merged_at `01:20:21Z` 換算）；commit 時間 09:20:20 | API `pulls/31` |
| main 檔案 | go.mod `go 1.27.3`（02842cf5…8ebf）；workflow 三處 `"1.27.3"`（06833a31…f245），與 CR 擬議值逐位元相同 | evidence 01 |
| run 41559091290 | push、main、head_sha＝merge、completed／success、attempt 1；建立 09:20:24、完成 09:22:19（UTC+8，由 Z 換算） | evidence 07（09:23:55） |

轉述值（09:20、1b53fcde、47475638…2add、41559091290）與查證值**全部相符**。

## 3. S1～S7 狀態
| # | 原文要點（35 號 §5） | 狀態 | 實際時間（UTC+8） | 依據／證據 |
|---|---|---|---|---|
| S1 | 只改 workflow 與 go.mod、與擬議值一致、同一 commit；最遲 2026-11-20 合併 | **完成（附註）** | 合併 2026-11-03 09:20:21 | 檔案 sha256 相符（evidence 01）；36 號 §8.3：`164b62c..47475638` 只剩兩檔。附註：品保部將「單一 commit」的偏離（PR 期間以 merge 同步 main）列為待拍板，見 §6 異常 2 |
| S2 | PR 最新 head 四綠、L1/L3/L4/L5、安全部書面確認 | **完成** | 2026-11-02 16:32 | 36 號 v0.1.1 §8（head `47475638`）；合併的 head 即此值 |
| S3 | 合併後 main push run 四綠，log 印出 `go1.27.3` | **完成** | 查證 09:23:55～09:24:58 | 四 job success（evidence 07）；未登入 job log 仍 403，log 內容以維運副本抽查（`/workspace/rebuild/cr-006/evidence-main-1b53fcde/`，SHA256SUMS -c 6 檔 OK）：gosec／sca／unit 各有 `go version go1.27.3 linux/amd64`、`GOVERSION='go1.27.3'`，1.27.2 字樣 0、`switching to go\|downloading go1\.` 0；gosec h1 OK、Files 57／Lines 7485／Nosec 0／Issues 0；sca `No vulnerabilities found.`（evidence 10）。品保部合併後取證結論相同 |
| S4 | G404 金絲雀，合併後 7 日內；推分支前要使用者核准；拿到 CI 紅燈後刪分支留證 | **完成**（v0.1.2） | CI run 2026-11-03 09:42:03～09:43:48；安全部核對 09:45:21～09:45:53；截止 2026-11-10 09:20:21 前 | run 41561438299 sast-gosec failure：G404 HIGH、Issues 1、exit 1；分支已刪（§4.1；evidence 11）。核准為**轉述**（§4.1）。本機部分見 evidence 02～06 |
| S5 | S3、S4 完成後：03 §2 Go 列改 1.27.3；05 RR-012 結案；13 C8 補註 | **完成**（v0.1.2） | 2026-11-03 09:46:20 | 03 v0.3.11（§2 Go 列 1.27.3）、05 v0.1.15（RR-012 已結案，綁 merge `1b53fcde`＋run 41561438299）、13 v0.1.5（C8 補註「1.27.3 由 CR-006 執行」）；備份 `/workspace/tmp/sec-s5-20261103-094559/` |
| S6 | G5 部署清單加「正式建置用 Go 版本＝CI（1.27.3）」並附證據 | **待到期**（G5 前，維運部） | 查證 09:24:38 | `docs/09-cicd-env/07-g5-deploy-checklist-v0.1.md` v0.1（7f697ce9…700b）目前**沒有**此列（1.27.3 計數 0，evidence 09）。安全部不改維運部檔案，請維運部在 G5 送審前補上 |
| S7 | 回滾時 RR-012 恢復、CI 回 1.27.2、7 日內重訂期限；禁止繞過 | **未觸發** | 查證 09:24:10～09:24:38 | main 仍為 `1b53fcde`；近 10 個 PR 沒有 revert／回滾（evidence 08）。條款持續有效 |

GO-2026-6809：CI gosec 二進位已改用 Go 1.27.3 建置（main run L1、h1 OK；36 號 M5-5 symbol 層 0），**修正已生效**；登錄結案依 RR-012 隨 S5 處理。GOTOOLCHAIN（SCR-053）照 35 號 §3 第 5 項另開後續項，不在本檔處理，最遲 2026-11-30 每月檢視時決定。

## 4. S4 待使用者核准的執行計畫（v0.1.1 改寫：以維運部準備的 commit 為準）
**截止：2026-11-10 09:20:21（UTC+8）**（合併時間加 7 日）。比照 CR-005 K4（32 號）。執行：維運部；覆核：安全部。

| 項目 | 值 | 安全部核對 |
|---|---|---|
| 分支 | `canary/cr006-g404-red` | 2026-11-03 09:28:23～09:28:27（UTC+8，取自 `date`）執行 `git ls-remote https://github.com/444112029012/short-url 'refs/heads/canary/*'`：輸出 0 行、exit 0，**遠端沒有此分支**，與維運部「尚未推送」的說法一致 |
| commit | `790b1c3aeef41ae5c80d6831b4afed27cebe926e`，父 commit main `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d` | 2026-11-03 09:29:19（UTC+8）以唯讀 `git log`／`git show` 讀維運部本機 repo `/workspace/ops-cr002`：父 commit 相符；`--stat` 只有 1 檔、+15 行 |
| 新增檔 | 只新增 `internal/cr006canary/canary.go`，package `cr006canary` | 相符（第 2 行 `package cr006canary`） |
| 檔案 sha256 | `b3df0e7c51f6141e1404d8362a17ae93601dc86842cec0d72e6f81deb189c8a6` | `git show 790b1c3…:internal/cr006canary/canary.go \| sha256sum` 相符 |
| 與 CR-005 11 號樣本 | 據維運部說法，第 3 行以下逐字相同 | 與 `docs/09-cicd-env/evidence/cr-005/11-canary-g404-token.go.txt` 的第 3 行以下 `diff` 為空，相符；`#nosec` 計數 0。推送後再以遠端內容核對一次 sha256 |

步驟：
1. 經使用者核准後，維運部把 `790b1c3a…` 推送為 `canary/cr006-g404-red`。不開 PR、不合併。推送後安全部以 `git ls-remote` 確認分支 head＝`790b1c3aeef41ae5c80d6831b4afed27cebe926e`，並以公開 clone 重算該檔 sha256。
2. 在該分支以 `workflow_dispatch` 觸發 `ci-devsecops.yml`（main 的 workflow 有 `workflow_dispatch`）。
3. 預期：sast-gosec **failure**；log 有 `internal/cr006canary/canary.go:12 - G404 (CWE-338) … Severity: HIGH`、`Issues : 1`、`Process completed with exit code 1`；L1 為 `go1.27.3`、切換訊息 0；Files 58／Lines 7500（本機值，供比對）。若 sast-gosec 綠燈或因載入失敗而紅（沒有 G404 那一行），判定 S4 **不通過**，交安全部重評 R4。
4. 留證：run ID、API run／jobs JSON、log 副本與 SHA256SUMS（維運放 `docs/09-cicd-env/evidence/cr-006/`，安全部覆核結果放本目錄）。
5. 刪除分支，安全部以 `git ls-remote origin 'refs/heads/canary/*'` 確認 0 行、branches API 404。
6. 安全部出 37 號小版：S4 改完成；接著依 S5 出 03（§2 Go 列改 1.27.3）、05（RR-012 結案）、13（C8 補註）。
若截止前無法完成，請使用者決定是否延期；延期期間 RR-012 維持「待結案」。

**審核問句**（核准由使用者在維運部的聊天中直接給，安全部不另外徵求核准）：
> 「是否核准維運部把 commit `790b1c3aeef41ae5c80d6831b4afed27cebe926e`（父 commit main `1b53fcde`，只新增 `internal/cr006canary/canary.go`，package `cr006canary`，sha256 `b3df0e7c…c8a6`，G404 假樣本、無 `#nosec`）推送為分支 `canary/cr006-g404-red`，並在該分支以 workflow_dispatch 觸發 ci-devsecops？不開 PR、不合併；預期 sast-gosec 紅燈（G404 HIGH、Issues 1、exit 1）；取得 run 與 log 後立即刪除分支並留證，最遲 2026-11-10 09:20:21（UTC+8）前完成。核准／不核准？」

### 4.1 S4 執行結果與安全部核對（v0.1.2 新增）
**核准（轉述）**：據維運部轉述，使用者在 2026-11-03 09:41 於維運部聊天 t116u 核准，原文為「核准 G404 金絲雀：推送 canary/cr006-g404-red（commit 790b1c3a，base 1b53fcde，只新增 G404 假樣本，不開 PR、不合併），以 workflow_dispatch 跑 CI，驗證紅燈後刪除分支」。**安全部沒有看到原文。**

核對時間 2026-11-03 09:45:21～09:45:53（UTC+8，取自 `date`），輸出存 `evidence/cr-006-s4/11-s4-ci-canary-verify.txt`。

| # | 核對項 | 結果 |
|---|---|---|
| 1 | 維運證據 `/workspace/rebuild/cr-006/evidence-canary-g404/` | SHA256SUMS sha256 `20a382d1c50a69dfd44458b020f0370f50ca9ea8a38f0ae9f2442178d3d9c536`，與維運提供值相符；`sha256sum -c` 4 檔全部 OK（只讀） |
| 2 | 未登入 API run 41561438299 | event `workflow_dispatch`、head_branch `canary/cr006-g404-red`、head_sha `790b1c3aeef41ae5c80d6831b4afed27cebe926e`、completed／**failure**、attempt 1；created 2026-11-03 09:42:03、updated 09:43:48（UTC+8，由 Z 換算）。jobs：sast-gosec 125672736697 **failure**（09:42:07～09:43:47）；secrets-gitleaks 125672736446、unit-test 125672736705、sca-govulncheck 125672736799 皆 success。job log 未登入取得回 **403**（照實記錄，log 內容依維運副本） |
| 3 | log 指定行（維運副本 `sast-gosec.log`） | 140 `go version go1.27.3 linux/amd64`；257 `gosec h1 OK`；375 `[…/internal/cr006canary/canary.go:12] - G404 (CWE-338) … Severity: HIGH`；387 `Issues : 1`；389 `##[error]Process completed with exit code 1.`，五行全部相符。補充：187 `GOVERSION='go1.27.3'`；384～386 Files 58／Lines 7500／Nosec 0（與本機值相同）；切換訊息 0、1.27.2 字樣 0；第 95、107 行 checkout 為 `790b1c3a…` |
| 4 | 遠端分支與 main | `git ls-remote 'refs/heads/canary/*'` 0 行；branches API 404；main 仍為 `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`。維運部回報的刪除時間 09:44:26 是轉述 |
| 5 | 只有一筆 G404 | 全 log 的 `] - G### ` 發現行只有 1 筆（第 375 行，canary.go:12），Issues 1，沒有其他 Issue |

**結論：S4 成立。** 換 Go 1.27.3 後 gosec 仍能分析並擋下 HIGH，R4 已排除。


## 5. 本次文件
| 檔案 | 現行版 | 內容 | 歷程 |
|---|---|---|---|
| 03-ci-security-gates-v0.3.md | v0.3.11 | §2 Go 列 1.27.3 | v0.3.10 合併狀態註記（v0.1） |
| 05-residual-risk-register-v0.1.md | v0.1.15 | RR-012 結案 | v0.1.14 合併註記、待結案（v0.1） |
| 13-cr-002-…-v0.1.md | v0.1.5 | C8 補註 | v0.1.4 C8 狀態註記（v0.1） |
| 36-…-v0.1.md | v0.1.3 | 版本來源、改引 38 號 | v0.1.1、v0.1.2 另存快照檔 |
| 38-cr-006-m5-gitleaks-false-positive-ruling-v0.1.md | v0.1 | M5 證據 gitleaks 誤判獨立裁定 | 新增（PR33-S3） |
| README.md | v0.36 | 索引同步 | v0.32（v0.1）→ v0.34（v0.1.2）→ v0.35 → v0.36 |
備份：`/workspace/tmp/sec-cr006-post-20261103-092458/`（v0.1）、`/workspace/tmp/sec-pr33-20261103-155039/`（v0.1.3）。

## 6. 異常與記錄
1. 未登入 API 取 job log 仍回 403（SCR-060 延續），S3 log 內容依維運副本（SHA256SUMS 已驗）與品保取證交叉比對。
2. 品保部合併後取證：rule suite 4869407317 result `bypass`（main-pr-review 的 pull_request 規則 fail：「New changes require approval from someone other than the last pusher.」）。屬 RR-011／D-17 已接受、CR §6 R10 已預期的情況，本檔不重新評估；另品保部把「單一 commit」偏離列為待拍板，由審查部／使用者判定。
2a. **單一 commit 偏離（轉述，v0.1.1 新增）**：據 Grok Bot 轉述，使用者在 2026-11-02 17:06 於 Grok Bot 的聊天中，以 widget t375s1422 的同一句核准接受了這個偏離。**安全部沒有看到原文**，本列只記錄轉述。
2b. **v0.1.2**：推送來源是 worktree `/workspace/ops-cr006-canary`，屬於 `/workspace/ops-cr002` 這個 repo（維運部說法；與 v0.1.1 看到的 commit 所在 repo 一致）。目錄名稱和 CR 編號不一致，不影響結果。
3. run 41559099969（event dynamic，success）與 check-run `update-go_modules-graph` 為 GitHub 依賴圖自動作業，非 required checks。
4. S6 目前未落實（G5 前到期），非本次阻擋。

## 修訂紀錄
| 版本 | 日期 | 負責 | 說明 |
|---|---|---|---|
| v0.1 | 2026-11-03 | 安全部 | 初版：S1～S7 狀態彙整；S4 本機部分完成、CI 部分列待核准計畫；小版修訂不改檔名 |
| v0.1.1 | 2026-11-03 | 安全部 | §4 依維運部準備的 commit `790b1c3a…`（分支 `canary/cr006-g404-red`、`internal/cr006canary/canary.go`）改寫並核對（ls-remote 0 行、sha256 相符）；審核問句改名並註明核准由使用者在維運部聊天中直接給；§6 新增單一 commit 偏離的轉述紀錄；小版修訂不改檔名 |
| v0.1.2 | 2026-11-03 | 安全部 | S4 完成（run 41561438299 紅燈，安全部獨立核對 5 項相符；核准為轉述）；S5 完成（03 v0.3.11、05 v0.1.15、13 v0.1.5）；新增 §4.1；證據 SHA256SUMS 更新為 79961780…；小版修訂不改檔名 |
| v0.1.3 | 2026-11-03 | 安全部 | 回應審查部 PR33-S5：§5 改列現行版本與歷程；v0.1.2 原文另存快照；新增 36 號 v0.1.3 與 38 號。S1～S7 狀態不變 |
