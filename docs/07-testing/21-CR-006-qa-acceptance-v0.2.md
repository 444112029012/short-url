---
文件：CR-006 品保驗收結論（CI 升 Go 1.27.3；PR #31）
版本：v0.2
狀態：審查中
負責角色：品保部
最後更新：2026-11-03（Asia/Taipei，UTC+8）
依據：CR-006 v0.4（repo 預定路徑 `docs/99-changelog/CR-006-ci-go-1.27.3-v0.4.md`，待 PR33-R3 入庫；本機原件 `/workspace/rebuild/cr-006/CR-006-ci-go-1.27.3-v0.4.md`；sha256 `818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734`）§5.1 L1～L7、§5.2 M5-1～M5-8；品保影響評估 `品保部.md`（未入庫，本機 `/workspace/rebuild/cr-006/impact/品保部.md`；v0.4，sha256 `98d714da2467f8e41ce124fdad814c0598a08112a42107efbc03644aeaed03a0`）§2 QA-CR006-1～15
比對基準：合併前 main `ca09137ae8f336a8bb02ac7f20e893eb2a6ccfbe`（push run 40618331206）、同步後 main `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`（push run 41402171425）；合併後 main `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`（push run 41559091290）
專案代號：SHORTURL
證據路徑：repo `docs/07-testing/evidence/cr-006/`（A～D 四集，v0.2 起，回應審查部 PR33-R4）
---

# CR-006 品保驗收結論

## 0. 結論：**CR-006 品保驗收 PASS；缺陷 0**

QA-CR006-1～15 全部判定完成：14 項 PASS；QA-CR006-15（合併紀錄）已記錄，合併經 bypass，依判準由審查部／使用者判定是否接受。驗收過程中**未發現缺陷**（未開 DEF）。有 2 項偏離（§4）與 5 項不阻擋建議（§5）。

- 本結論由品保部以 GitHub 唯讀查詢、box 本機實跑取得證據；未 commit、未 push、未在 GitHub 留言。
- 本結論**不是合併核准**（PR #31 已於 2026-11-03 09:20:21（UTC+8）合併），也不判定 RR-012 是否結案（屬安全部）。

## 1. 追溯

### 1.1 綁定文件
| 文件 | 版本 | sha256 |
|---|---|---|
| CR-006 `CR-006-ci-go-1.27.3-v0.4.md` | v0.4 | `818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734` |
| 品保影響評估 `impact/品保部.md` | v0.4 | `98d714da2467f8e41ce124fdad814c0598a08112a42107efbc03644aeaed03a0` |
| 安全部 M5／S2 確認 36 號（#31 合併所依據的是 v0.1.1；repo 上 `docs/06-security/36-cr-006-pr31-m5-s2-confirmation-v0.1.md` 現為 v0.1.2 `2b7ade2f…`，v0.1.1 原件待安全部依 PR33-R2 另存新檔名入庫） | v0.1.1 | `2936baec4ff6d20474a4b31a30e9820e8a3a1177de17544518dd20179d58d598` |

### 1.2 需求 → 變更 → 驗收 → 證據
| 需求 | 變更（CR-006 §1） | 驗收項 | 證據集 |
|---|---|---|---|
| SEC-009（SAST，門檻不變） | workflow 三處 `go-version: "1.27.3"`；gosec v2.30.0 改由 Go 1.27.3 建置 | QA-CR006-1～5、6、10(c) | B、C、D、E |
| SEC-010（SCA，門檻不變）；RR-012／GO-2026-6809 | 執行 Go 改 1.27.3，stdlib 修補 GO-2026-6809 | QA-CR006-10、11 | B、C、D、E |
| NFR-004（CI 門檻） | 四個 required checks 名稱與門檻不變 | QA-CR006-6、13、15 | A、C、D、E |
| SI-05（工具鏈一致） | go.mod 第 3 行 `go 1.27.3`（方案甲，無 toolchain 行），與 setup-go 同 commit | QA-CR006-12、14 | B、C、D |
| 回歸（單測、CT、覆蓋率） | 無 `.go` 變更 | QA-CR006-7、8、9 | B、C、D |

### 1.3 證據集與 SHA256SUMS（2026-11-03 品保重算，`sha256sum -c` 全部 OK）
| 代號 | repo 路徑（`docs/07-testing/evidence/cr-006/` 底下；括號內為本機原始目錄，`/workspace/rebuild/cr-006/impact/` 底下） | 範圍 | SHA256SUMS sha256（現值） | 轉述值 |
|---|---|---|---|---|
| A | `A-pre-push-ruleset/`（`qa-evidence-acceptance/pre-push-ruleset-snapshot/`）；9 檔 | 推送前 ruleset 快照（2026-10-28 16:35:19（UTC+8））；**沒有 results.md**，判定寫在 E | `594ac3c7dea19698959b3d8c7122fa93df23838d297778eea4810010e039c152` | 594ac3c7…，相符 |
| B | `B-pr31-826d5d34/`（`qa-evidence-pr31/`）；33 檔 | 舊 head 826d5d34、run 40660317286；含 10(b)(c) 核對 | `eff8f2e6ce0b762cc472357ca32cbb26e60f23bfdabe14e9de917044839f91a9` | eff8f2e6…，相符 |
| B0 | 同 B（第一輪，已被 B 取代） | B 追加 §4 前的版本 | 無法重算（檔案已重產） | 45418aa4…（品保當時回報值） |
| C | `C-pr31-47475638/`（`qa-evidence-pr31-47475638/`）；32 檔 | 新 head 474756381bd1…、run 41403811744；含 10(b)(c) 核對 | `1af3cf98cc9e852cf5f8fc2e82b3a729259380762532b4da515377b0dfe61dab` | 1af3cf98…，相符 |
| C0 | 同 C（第一輪，已被 C 取代） | 10(b)(c) 待定時的版本 | 無法重算（檔案已重產） | 78c55146…（品保當時回報值） |
| D | `D-post-merge/`（`qa-evidence-post-merge/`）；27 檔 | merge 1b53fcde、run 41559091290、rule suite、ruleset 比對 | `efe052caf3ba344ee3b98031b070018e9744784419b5e02a18b87deb8e26344e` | efe052ca…，相符 |

> 入庫時四集內容與本機原始目錄逐位元組相同（`diff -r` 一致），沒有排除、遮蔽或改寫任何檔案，所以 SHA256SUMS 重產後與 v0.1 所記的值位元組相同。各集 `results.md` 內文仍寫本機路徑（證據原件凍結不改），對照如下：`qa-evidence-acceptance/pre-push-ruleset-snapshot/`→`A-pre-push-ruleset/`、`qa-evidence-pr31/`→`B-pr31-826d5d34/`、`qa-evidence-pr31-47475638/`→`C-pr31-47475638/`、`qa-evidence-post-merge/`→`D-post-merge/`；D `results.md` 第 31 行的 `../qa-evidence-pr31-47475638/base-run-41402171425.log` 即 `C-pr31-47475638/base-run-41402171425.log`。
> 安全部 M5 證據在 repo `docs/06-security/evidence/cr-006-m5/<head SHA>/`；其中 01、04 號檔依安全部裁定（gitleaks generic-api-key 誤判公開資料）不入庫，本文引用的 05、06 號都在 repo 內。
> 45418aa4…、78c55146… 是同一目錄較早一輪的 SHA256SUMS；後來追加內容並重產，舊清單檔沒有保留，所以**無法重算，只能記為歷史回報值**。現行有效值是 B、C 兩列。

## 2. 逐項判定（QA-CR006-1～15）

行號：「B:22」表示 `docs/07-testing/evidence/cr-006/B-pr31-826d5d34/results.md` 第 22 行（A～D 同理）；log 行號指同一集內的 per-job API 原始 log（`api-job-*.log`）。v0.2 已對入庫後的檔案逐一重核，行號不變。

| 項目 | 判準摘要 | 判定 | 證據 |
|---|---|---|---|
| QA-CR006-1（L1） | 三個 Go job 為 go1.27.3；無 1.27.2 字樣 | **PASS** | B:22；C:23；D:20（`api-job-125665062413-sast-gosec.log` 115／133／172／180 等，`23-checks.txt`） |
| QA-CR006-2（L2） | key 前綴 `go-1.27.3-`、雜湊段 `cf6b0d42…a025` | **PASS** | B:23；C:24；D:21（sca log 214 `Cache saved with the key`） |
| QA-CR006-3（L3） | `switching to go\|downloading go1\.` 每個 job 0 | **PASS** | B:24；C:25；D:22 |
| QA-CR006-4（L4） | gosec v2.30.0 h1 `67aF6LHa…` 與 `gosec h1 OK` | **PASS** | B:25；C:26；D:23（sast log 250、251） |
| QA-CR006-5（L5） | 命令列不變；Files／Lines／Nosec／Issues＝57／7485／0／0＝基準 | **PASS** | B:26（log 384～387）；C:27（391～394）；D:24（368～371）；基準 40618331206 第 639～642 行、41402171425 第 801～804 行 |
| QA-CR006-6 | 四個 required checks success、名稱不變 | **PASS** | B:28；C:29；D:25（`21-check-runs-merge.json`、`22-required-contexts.txt`） |
| QA-CR006-7（L7） | 181／0／0、`-race`、ENG-011 PASS | **PASS** | B:29；C:30；D:26（unit log 202、218；`art/go-test.json`） |
| QA-CR006-8（L7） | TestCT_* 14 pass | **PASS** | B:30；C:31；D:27 |
| QA-CR006-9 | `-coverpkg=./...` 合計 ≥ 80.9% | **PASS** | B:31；C:32（`-count=1`）；D:28（merge sha 全新 clone、Go 1.27.3、`-count=1`：80.9%，精算 474/586＝80.887%，與基準相同） |
| QA-CR006-10 | (a) L6；(b) 專案 `-show verbose` go1.27.3 標準庫、無 6809；(c) gosec 二進位無 6809、exit 0 | **PASS** | (a) B:27、C:28、D:29（sca log 203）；(b)(c) C §4（C:45～78；安全部證據 `docs/06-security/evidence/cr-006-m5/474756381bd170beb3a5d723b63f1b8ba48b2add/` 05、06；品保重現 `C-pr31-47475638/32-qa-m5-5-gosec-binary.txt` 第 12、112 行、`C-pr31-47475638/33-qa-m5-6-project-verbose.txt` 第 19、32、33 行）；D:30（merge tree＝該 head） |
| QA-CR006-11 | SCA 三層 0／1／0 → 0／0／0；其他差異記錄 | **PASS** | B:33；C:34；D:31（基準 `C-pr31-47475638/base-run-41402171425.log` 第 398～402 行）；沒有 GO-2026-6809 以外的差異；dependency-review 無 Critical（C:28） |
| QA-CR006-12 | go.mod `go 1.27.3`、無 toolchain 行；三處 `"1.27.3"`、`check-latest: false`；無 Dockerfile | **PASS** | B:34；C:35；D:32（`30-diff-gomod.txt`） |
| QA-CR006-13 | ruleset 前後 `jq -S` diff 為空（含 updated_at） | **PASS** | A（`sha256sum -c` 9 列 OK）；D:19（`11-ruleset-diff.txt`：24629805、24629854、清單三檔 `cmp` 位元組相同） |
| QA-CR006-14 | 只有 workflow 4/4、go.mod 1/1；內容 `06833a31…f245`／`02842cf5…8ebf`；go.sum 不變 | **PASS（偏離見 §4.1）** | B:36；C:36、C §3（C:39～43）；D:33（`30-diff-gomod.txt`：舊 main 164b62c9→merge 只有 2 檔；`diff(474756381bd1, merge)` 空） |
| QA-CR006-15 | 合併事實可追溯；required_status_checks 必須 pass | **已記錄（required_status_checks pass；result bypass，見 §4.2）** | D:18、D §2（D:35～41；`02-pr31.json`、`04-merge-commit.txt`、`05-rule-suites.json`、`06-rule-suite-4869407317.json`） |

合併紀錄摘要：merged_at 2026-11-03 09:20:21（UTC+8）；merge_commit_sha `1b53fcde7cf5068dcf7193b92205d6e6e9adb09d`；merged_by `444112029012`；方式＝merge commit（依兩個 parent `164b62c9`＋`474756381bd1` 與訊息判斷，PR API 沒有方式欄位）；未登入 ls-remote main＝merge sha。

## 3. 安全部 M5／S2 核對摘要（QA-CR006-10(b)(c)）
- 舊 head 826d5d34：36 號 v0.1（`b7444593…`）與證據 SHA256SUMS `33c96bae…`，M5-1～M5-8 PASS（B §4）；head 已變，**作廢**。
- 新 head 474756381bd1…：36 號 v0.1.1（`2936baec…`）與證據 SHA256SUMS `43d33bb4…`（14 檔，`-c` OK），M5-1～M5-8 PASS；品保在全新 clone 重現 M5-5、M5-6，都 exit 0、GO-2026-6809 0，內容與安全部相同（C §4）。module 層 GO-2026-6866、6860、5932 屬已登錄的 SCR-054／055，不啟動 M7；沒有以 `-scan module` 作判準。

## 4. 偏離
### 4.1 QA-CR006-14：多了一個同步 main 的 merge commit
- 原文要求兩檔「在同一個 commit」；PR 最後是 `826d5d34`（兩檔，同一 commit，**仍成立**）＋ `474756381bd1`（`Merge remote-tracking branch 'origin/main'`）。
- 品保改以**淨差異**判定，並經使用者指示採用：merge-base→head 只有 2 檔；merge 本身沒有夾帶 main 以外的內容（`git show --cc` 空；帶入內容與 main `ca09137..164b62c9` 逐位元組相同，只有 PR #32 的 5 個 docs 檔）。**品保接受此偏離。**

### 4.2 合併經 bypass（rule suite 4869407317）
- result `bypass`、evaluation_result `null`；24629805：required_status_checks／non_fast_forward／deletion 全 pass；24629854：**pull_request fail**（"New changes require approval from someone other than the last pusher."），其餘 pass。
- bypass 依據是 24629854 的 bypass_actors（RepositoryRole 5、`bypass_mode: pull_request`）。品保**只記錄**；是否接受由**審查部／使用者**判定。
- 「此情況已由 RR-011／D-17 接受（CR §6 R10）」：品保仍**未直接查證**。v0.2 補記：審查部 PR #33 審查紀錄（`/workspace/review-PR33-2026-11-03.md` 第 72 行；入 `docs/gates/` 由總協調依 PR33-R4 處理）判定「這次 bypass 屬 RR-011／D-17 已接受的類型，且 D-09 人工核對 head 成立（merge parent＝47475638）……可以接受」。QA-CR006-15 據此**結案**，品保判定不變（只記錄）。

## 5. 不阻擋建議
| # | 內容 |
|---|---|
| QA-S10 | CR L6 寫明取 log 的方式：`gh run view --log` 步驟欄是 `UNKNOWN STEP`，要用整個 job log 或 `##[group]Run govulncheck` 區段（品保部 v0.4 第 290 行） |
| QA-S11 | CR L6 格內「（govulncheck 掃描時用的 stdlib＝執行的 Go）」移回 `GOVERSION='go1.27.3'` 之後（品保部 v0.4 第 305 行） |
| QA-S12 | QA-CR006-14 判準改寫為「merge-base→head 淨差異只有 2 檔、內容 sha256 相符；如有同步 main 的 merge，merge 不得帶入 main 以外的內容」，取代「同一個 commit」的字面要求（下一次 CR 或品保文件改版時處理） |
| 註 1（36 號命名） | 36 號 v0.1.1 沿用 `-v0.1.md` 檔名就地覆寫，本機已沒有 v0.1 檔案（v0.1 原文 `b7444593…` 只留在 repo）；建議安全部改版時另存或改檔名 |
| 註 2（exit code） | 安全部新 head 證據 `01-m5-1-head.txt` 沒有記錄 ls-remote 和 curl 的 exit code（舊 head 那一輪有）；值本身經品保重取相符 |

## 6. 未能直接查證的事項
- RR-011／D-17 是否涵蓋本次 bypass（§4.2）：品保未直接查證；審查部已判定可接受。
- 36 號 v0.1.1（`2936baec…`）目前無法從 repo 驗證，待安全部依 PR33-R2 入庫；CR-006 v0.4 本文待 PR33-R3 入庫；`品保部.md` 不在 repo（審查部 CR6-S4 列為 S 項）。
- 品保重現用的 box 預裝 Go 1.27.3 沒有再對 go.dev 官方 JSON 核對 tarball；工具鏈獨立核驗以安全部 M5-2 為準。
- 未登入取 job log 回 403（SCR-060）；品保用登入的唯讀權限取得 log，並與維運副本逐位元組比對相同。
- B0、C0 的 SHA256SUMS 無法重算（§1.3）。

## 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.2 | 2026-11-03 | 品保部 | 回應審查部 PR #33 退回 PR33-R4：A～D 證據集入庫至 `docs/07-testing/evidence/cr-006/`（內容與 SHA256SUMS 位元組不變，gitleaks 8.21.2 掃描 0 筆）；所有證據引用改為 repo 路徑並加本機路徑對照；行號對入庫檔重核；§1.1、依據欄註明 36 號 v0.1.1 與 CR-006 本文尚未入庫（PR33-R2、R3）；§4.2、§6 補記審查部對 bypass 的判定；結論不變（PASS、缺陷 0） |
| v0.1 | 2026-11-03 | 品保部 | 初版：QA-CR006-1～15 判定（14 PASS、15 已記錄）；證據集 A～D 的 SHA256SUMS 重算；偏離 2 項（14 淨差異、合併 bypass）；建議 S10～S12 與兩項註記；結論 PASS、缺陷 0 |
