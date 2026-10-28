---
文件：06-security 章節索引
版本：v0.24
狀態：審查中
負責角色：安全部
最後更新：2026-10-27
專案代號：SHORTURL
---

# 06 資安｜章節索引

> 命名規則：檔名保留主版號；小版修訂（如 v0.3.1、v0.3.2）只更新檔頭「版本」與修訂紀錄，不改檔名，以免斷鏈。

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [00-asvs-level-v0.1.md](./00-asvs-level-v0.1.md) | v0.1 | 審查中 | ASVS 5.0 **Level 1** 選定與理由 |
| [01-sec-requirements-v0.2.md](./01-sec-requirements-v0.2.md) | **v0.2** | 審查中 | SEC-001～**SEC-017**（G2 回寫 TM；新增 016／017） |
| [01-sec-requirements-v0.1.md](./01-sec-requirements-v0.1.md) | v0.1 | 歷史 | G1 初版；請以 v0.2 為準 |
| [02-g1-security-pregate-v0.1.md](./02-g1-security-pregate-v0.1.md) | v0.1 | 審查中 | G1 審查部用書面前置結論 |
| [03-ci-security-gates-v0.3.md](./03-ci-security-gates-v0.3.md) | **v0.3.9** | **已核准**（D-07／D-08）；§2 依 CR-005（D-19）**已生效**（#29，2026-10-27 16:37:12 UTC+8，merge `db380a40`） | **現行 CI 門檻**：Go **1.27.2（精確釘版）**／gosec **發行版 v2.30.0**（h1 `67aF6LHa…` 比對）／govulncheck v1.8.0；CR-002 的 gosec 暫行例外已結束（C6 見 29、K4 見 32）；dependency-review fail-closed；方案 A（D-08）；規則一 strict=true（D-10）。完整歷程見檔內修訂紀錄 |
| [03-ci-security-gates-v0.2.md](./03-ci-security-gates-v0.2.md) | v0.2.3 | 已取代 | G3 定案稿；**已由 v0.3 取代（D-07）** |
| [03-ci-security-gates-v0.1.md](./03-ci-security-gates-v0.1.md) | v0.1 | 歷史 | 已由 v0.2 取代 |
| [04a-dfd-working-v0.1.md](./04a-dfd-working-v0.1.md) | v0.1 | 審查中 | DFD 對齊說明（採設計部定稿；**EE1**） |
| [04-threat-model-v0.1.md](./04-threat-model-v0.1.md) | v0.1.1 | 審查中 | STRIDE 威脅模型（23 條 TM）；TM-002 已接受（D-06） |
| [05-residual-risk-register-v0.1.md](./05-residual-risk-register-v0.1.md) | **v0.1.13** | 審查中 | 殘餘風險：RR-001 已接受（D-06）；RR-006 已接受（D-08）；RR-007～009 已接受（D-11），RR-007 已降為低；**RR-010 已結案**（CR-005，2026-10-27）；**RR-011（中）已接受（D-17）**，D-09 人工核 head 持續；**RR-012（低）**：gosec 二進位含 GO-2026-6809，升 Go 1.27.3 的 CR 最遲 2026-11-20 提出（K6） |
| [06-g2-security-pregate-v0.1.md](./06-g2-security-pregate-v0.1.md) | v0.1.1 | 審查中 | G2 審查部用書面前置結論（後註：D-06 已接受） |
| [07-g3-security-pregate-v0.1.md](./07-g3-security-pregate-v0.1.md) | v0.1 | 審查中 | G3 資安書面前置 |
| [08-pr-002-security-review-v0.1.md](./08-pr-002-security-review-v0.1.md) | **v0.1.5** | **已核准**；SCR-002／003 已關閉（v0.1.5，以 11 號檔為準） | **PR #2 安全審查**：核准；CR-001（D-07）與合併路徑 A（D-08）已核准 |
| [09-cr-ci-deviations-ruling-v0.1.md](./09-cr-ci-deviations-ruling-v0.1.md) | **v0.1.4** | **已核准**（D-07／D-08） | **CI 偏離（CR-001）＋合併路徑（D-08）裁示**；履行複核（對象 PR #6 `37d906a`）：已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束 |
| [10-dr-g4-01-body-limit-security-opinion-v0.1.md](./10-dr-g4-01-body-limit-security-opinion-v0.1.md) | **v0.1.2** | **附條件同意；C1～C5 已達成** | **DR-G4-01 本文上限 32768 B 複核**：C1（實作 PR 須在 PR #7 合併後）**已達成**（main `24b8da0`）；**v0.1.2：C2～C5 已達成**，MaxHeaderBytes（C3）**已合併**（#17，merge `433140c`，2026-10-19 14:00:38 UTC+8） |
| [11-pr-007-008-security-review-v0.1.md](./11-pr-007-008-security-review-v0.1.md) | **v0.1.6** | **PR #7 已合併（最終 head `fc9a7a7` 經 16 通過）／PR #8 核准並已合併** | **PR #7（ENG-011 限流）／PR #8（ENG-016 標頭）安全審查**：SCR-006～011、**SCR-017（低）、SCR-018（資訊）**；v0.1.5：SCR-002／006／007／008／017 已關閉、SCR-011 已關閉（main 四綠）；**SCR-010（低，HSTS 由邊緣滿足；開放，G5 前驗證）**；**SCR-003 已關閉**（v0.1.6：#8 合併，`e9f4396`）；SCR-018 資訊 |
| [12-pr-011-security-review-v0.1.md](./12-pr-011-security-review-v0.1.md) | **v0.1.5** | **核准**；#11 已合併（2026-10-19 11:34:13 UTC+8，`3d3a9cd`） | **PR #11（ENG-012 觀測）安全審查**：SCR-012～016；short_code 入日誌可接受（附條件：日誌存取同 DB、≤14 日輪替）；**SCR-012 已關閉** |
| [13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md](./13-cr-002-go-1272-gosec-pin-security-opinion-v0.1.md) | **v0.1.3** | **安全部已同意（附條件）**；CR-002 已生效（#15）；**C6 已關閉（29）、C7 已結束（CR-005，32）** | **CR-002 安全部意見**：Go 精確 1.27.2＋gosec 暫釘 `707f227`；條件 C1～C8；C7 改為每月檢視工具版本；C8（Go 修補版走 CR）持續有效 |
| [14-pr-015-cr-002-security-review-v0.1.md](./14-pr-015-cr-002-security-review-v0.1.md) | v0.1 | **核准**；PR #15 已合併（2026-10-19 10:22:50 UTC+8，`c341be1`） | **PR #15（CR-002 實作）安全審查**：C1～C5 符合；h1 驗證步驟 fail-closed；**SCR-019（低）、SCR-020（資訊）、SCR-021（資訊、既有）** |
| [15-third-party-skills-security-review-v0.1.md](./15-third-party-skills-security-review-v0.1.md) | v0.1 | 審查中（安全部結論已出；待總協調轉使用者確認） | **第三方 AI agent skill 安裝前安全審查**（superpowers 三檔／Trail of Bits 兩項／Humanizer-zh／Google・Uber Go 指南） |
| [16-pr-007-final-head-security-confirmation-v0.1.md](./16-pr-007-final-head-security-confirmation-v0.1.md) | **v0.1.2** | **通過**；#7 已依條件以 head `fc9a7a7` 合併 | **PR #7 最終 head 合併前安全確認**：v0.1（`99b0879`）退回（SCR-017 註解）→ v0.1.1（`fc9a7a7`）通過；v0.1.2 合併後註記（`24b8da0`，main 四綠） |
| [17-pr-017-dr-g4-01-body-limit-security-confirmation-v0.1.md](./17-pr-017-dr-g4-01-body-limit-security-confirmation-v0.1.md) | **v0.1.2** | **#17 已合併**（以核准的 head `e6d9ab1` 合併，merge `433140c`）。v0.1.1：通過，僅適用於最終 head `e6d9ab1`（兩段核對）；v0.1：head `3b1aa5b` 通過 | **PR #17（DR-G4-01 本文上限 32768 B＋MaxHeaderBytes）合併前安全確認**：C1～C5 達成；**SCR-022、SCR-023、SCR-040、SCR-041 已關閉**（v0.1.2）；SCR-024、025、026（資訊）已處理；SCR-027（資訊）接受；SCR-042（流程）開放 |
| [18-pr-018-ratelimit-e2e-security-confirmation-v0.1.md](./18-pr-018-ratelimit-e2e-security-confirmation-v0.1.md) | **v0.1.4** | **#18 已合併**（2026-10-19 14:24:23 UTC+8，`643bf8a`，合併的 head 即核准的 `f547a6e`）。v0.1.3：同步後新 head 已書面再確認，通過；核准只適用於 `f547a6e`（之後任何 commit 都要重新送核）。v0.1.2：核准綁定 `663ee9d`（回應審查部 S4）；v0.1：head `9370def` 附條件通過 | **PR #18（速率限制 E2E 測試，SEC-007）安全確認**：SCR-028（低）、SCR-029／031（資訊）**已關閉**（v0.1.4）；SCR-030（既有，G5 前待辦）、SCR-032（維運範圍）開放；新增 SCR-033（資訊，CI 偵測只接受 `true`） |
| [19-pr-020-openapi-011-security-opinion-v0.1.md](./19-pr-020-openapi-011-security-opinion-v0.1.md) | **v0.1.1** | **#20 已合併**（2026-10-19 15:15:58 UTC+8，`abbec9e`，合併的 head 即核准的 `44b099c`）。v0.1：通過，僅適用於最終 head `44b099c` | **PR #20（設計部文件 PR：OpenAPI 0.1.1 本文上限／431、方法規格 v0.2.2）安全意見**：回應審查部 G4-PR20-R3 三個問題；SCR-043～045（資訊）**已關閉**（v0.1.1）；新增 SCR-046（資訊，根目錄 README 431 描述，不在本 PR 範圍；已轉研發部） |
| [20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md](./20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md) | v0.1.1 | 通過，僅適用於 head `0379b64`；#24 已合併 | **PR #24（根目錄 README 431／本文上限描述）安全確認**；SCR-046 已關閉（見 27） |
| [21-pr-025-ci-hardening-security-review-v0.1.md](./21-pr-025-ci-hardening-security-review-v0.1.md) | v0.1.1 | 有條件可合併，僅適用於 head `e54d900`；#25 已合併（`d0749f5`） | **PR #25（維運 CI 強化）安全審查**；SCR-004／019／021／032 已關閉（見 27） |
| [22-cr3-r8-stale-approval-residual-risk-v0.1.md](./22-cr3-r8-stale-approval-residual-risk-v0.1.md) | v0.1 | 評等完成：中 | **CR3-R8 殘餘風險評等**：ruleset 不強制核准綁定 SHA，登錄為 RR-011；「修補後可降為低」的說法以 26 為準 |
| [23-pr-023-s8-docs-sync-head-security-confirmation-v0.1.md](./23-pr-023-s8-docs-sync-head-security-confirmation-v0.1.md) | v0.1 | 通過，僅適用於 head `c3155b3` | **PR #23（08 v0.1.5＋README v0.22）同步後 head 安全確認** |
| [24-cr-003-branch-docs-exception-security-opinion-v0.1.md](./24-cr-003-branch-docs-exception-security-opinion-v0.1.md) | v0.1 | 同意，附三項前提 | **CR-003（純文件 PR 分支命名例外）安全意見**；綁 CR-003 sha256 `c55ad3ca…` |
| [25-kin-openapi-test-dependency-sca-opinion-v0.1.md](./25-kin-openapi-test-dependency-sca-opinion-v0.1.md) | v0.1 | **不同意**（低到中；使用者政策只接受低風險） | **第三方元件 SCA 意見**：kin-openapi（僅測試用） |
| [26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md](./26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md) | v0.1.1 | 附條件同意（CR-004 v0.3） | **CR-004（ruleset stale approval 兩旗標）安全評估**；RR-011 不降級，維持中；SCR-042 不關閉 |
| [27-pr-026-contract-tests-security-confirmation-v0.1.md](./27-pr-026-contract-tests-security-confirmation-v0.1.md) | v0.1 | PR #26 可合併，僅適用於 head `49e29d8`；#26 已合併（`f646bee`） | **PR #26（OpenAPI 契約測試）安全確認**＋#25 合併後 SCR 關閉紀錄（SCR-004／019／021／032／046） |
| [28-gate-G4-overall-security-conclusion-v0.1.md](./28-gate-G4-overall-security-conclusion-v0.1.md) | v0.1 | 安全面無未解決阻擋項 | **G4 整體安全結論**（綁 main `f646bee`）；檔內寫的「RR-011 待使用者書面接受」已由 D-17 滿足（見 05 v0.1.12） |
| [29-cr-002-c6-canary-closure-v0.1.md](./29-cr-002-c6-canary-closure-v0.1.md) | v0.1 | **C6 關閉**（逾期 1 天，記為偏差） | **CR-002 C6 CI 金絲雀關閉紀錄**：run 40430102886，G101 HIGH 被擋，分支已刪 |
| [30-cr-005-gosec-v2.30.0-security-opinion-v0.1.md](./30-cr-005-gosec-v2.30.0-security-opinion-v0.1.md) | v0.1.1 | 同意（附條件 K1～K6） | **CR-005（gosec 換回 v2.30.0）安全評估**；v0.1.1 更正作業時間 |
| [31-cr-005-pr-029-k1-k2-confirmation-v0.1.md](./31-cr-005-pr-029-k1-k2-confirmation-v0.1.md) | v0.1.1 | K1、K2 成立，僅適用於 head `1545143c`；#29 已合併（`db380a40`） | **CR-005 PR #29 K1／K2 書面確認**；v0.1.1 核對原始 log |
| [32-cr-005-k4-g404-canary-closure-v0.1.md](./32-cr-005-k4-g404-canary-closure-v0.1.md) | v0.1 | **K4 關閉；RR-010 結案** | **CR-005 K4 G404 金絲雀結案**：run 40477332321，G404 HIGH 被擋，分支已刪；只剩 K6（最遲 2026-11-20） |
| [33-pr-030-gitleaks-false-positive-ruling-v0.1.md](./33-pr-030-gitleaks-false-positive-ruling-v0.1.md) | v0.1 | 誤報；不放寬門檻 | **PR #30 gitleaks 命中裁定**：google.golang.org/api 的公開 h1 被 generic-api-key 規則誤判；命中的 3 個證據檔只保留在本機，不改 .gitleaksignore |

## 準據 DFD

| 項目 | 路徑／內容 |
|---|---|
| 設計部 DFD 正式稿（編號凍結） | [`docs/03-data/02-dfd-v0.1.md`](../03-data/02-dfd-v0.1.md) |
| 外部實體 | **EE1**（匿名客戶端） |
| 信任邊界 | TB-01～TB-04 |
| 處理／儲存／流 | P1～P7；DS1／DS2；DF-L0-01～06、DF-01～21 |

## 本章狀態（G3＋PR #2）

| 項目 | 狀態 |
|---|---|
| ASVS 等級 | 已選定 L1（審查中） |
| SEC 清單 | **v0.2**（SEC-001～017） |
| STRIDE／TM-xxx | 已完成 v0.1（23 條） |
| 殘餘高以上 | RR-001／TM-002：**D-06 已接受**（G5 前再評） |
| 殘餘中（合併治理） | **RR-006 已接受**（D-08；G5 前再評） |
| 殘餘（#7／#8／#11） | **RR-007（低；v0.1.10 由中降為低：#7，2026-10-19 12:11:10 UTC+8，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`）、RR-008（低）、RR-009（低）已接受**（D-11；G5 前再評） |
| CI 門檻書面 | **v0.3 已核准**（D-07；v0.2 已取代）；**v0.3.7 依 D-13（CR-002）更新 §2，v0.3.8 已生效（#15，2026-10-19 10:22:50 UTC+8，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）**；阻擋語意未放寬 |
| PR #2 安全審查 | **核准**（SCR-001 已關）；PR #2 已於 2026-10-07 13:56:49 UTC+8 合併（`a8246ce`） |
| CI 實作／分支保護生效 | 四 checks 已綠於 PR #2；Rulesets 方案 A **已生效**（D-08；取證 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）；規則一 **strict=true 已生效（D-10）**（維運 14:34:10 UTC+8 取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md`，隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；安全部以公開 API 獨立讀取一致；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）） |
| PR #7／#8 安全審查 | **#7 已合併**（#7，2026-10-19 12:11:10 UTC+8，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`；最終 head 經 16 號檔 v0.1.1 通過；main push run 39198199210 四項 success）；SCR-002／006／007／008／017 已關閉；SCR-011 已關閉（依序合併後 main 四綠）；**#8 已合併**（2026-10-19 11:21:48 UTC+8，`e9f43968236878eb605e8e3048bae59a35815e70`）；**SCR-003 已關閉**；**SCR-010（低）開放：HSTS 由邊緣（TLS 終止端）送出，應用層 `HSTS_ENABLED` 無行程內 TLS 時無效；文件修補＋G5 前以實際回應驗證（定義見 11 號檔；RR-008）** |
| PR #11 安全審查 | **核准；已合併**（2026-10-19 11:34:13 UTC+8，`3d3a9cd`）；**SCR-012 已關閉**（#7 修正測試簽章並合併，main 四綠）；日誌輪替／存取控制列維運 G5 前待辦 |
| 設計裁定 DR-G4-01（本文上限） | **附條件同意**（10 v0.1.1）；限流已隨 #7 合併生效（main `24b8da0`，不再是 stub），C1～C5 已達成（10 v0.1.2）；MaxHeaderBytes（C3）**已合併**（#17 於 2026-10-19 14:00:38 UTC+8 合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`） |
| PR #17 安全確認（17） | **已合併**：2026-10-19 14:00:38（UTC+8），merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`（parents `ed62039`＋`e6d9ab1`；合併的 head 即安全部核准的 `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`，`git diff e6d9ab1 433140c` 為空）。SCR-022／023／040／041 已關閉（17 v0.1.2）；SCR-042（流程）開放 |
| PR #18 安全確認（18） | **已合併**：2026-10-19 14:24:23（UTC+8），merge commit `643bf8ab052ebb175e35bc365ec97863ec9f4202`（parents `433140c`＋`f547a6e`；合併的 head 即安全部核准的 `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f`，`git diff f547a6e 643bf8a` 為空；測試檔 blob 仍為 `0e2b6d9f…`）。SCR-028／029／031 已關閉（18 v0.1.4）；SCR-030（G5 前待辦）、SCR-032、SCR-033 開放 |
| PR #20 安全意見（19） | **已合併**：2026-10-19 15:15:58（UTC+8），merge commit `abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`（parents `141fcda`＋`44b099c`；合併的 head 即安全部核准的 `44b099c227cac2ab24380f4346654540da6535de`，`git diff 44b099c abbec9e` 為空）。SCR-043～045 已關閉（19 v0.1.1）；SCR-046（資訊，根目錄 README）已轉研發部 |
| G5 前待辦（安全部追蹤） | **SCR-010**（低）：HSTS 由邊緣送出，G5 前以實際回應驗證（RR-008）；**SCR-030**（資訊，既有）：main 未設 `ErrorLog`、沒有 recover，panic 時 net/http 預設日誌會帶出 peer 位址，轉研發部，後續登錄殘餘風險（見 18）；**MaxHeaderBytes（DR-G4-01 C3）：已合併**（#17，`433140c5132778bf7cd731e3ecaca8187e50d678`；G5 時以部署環境回應驗證 431）；觀測日誌保存期／存取控制（SCR-013／015，見 12） |
| dependency-review fail-closed | Dependency graph 已開（compare 200）；已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束（原期限 min(2026-10-14, G5 送審前)，已於期限內履行） |
| CR-002（Go 1.27.2＋gosec 暫行例外） | 使用者 D-13 核准方案甲；安全部同意（附條件）（13）；PR 審查核准（14）；**已生效**（#15，2026-10-19 10:22:50 UTC+8，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`；main push run 39189683566 四項 success，`sast-gosec` 由 `2e54f7e` 的 failure 轉綠、log 有 `gosec h1 OK`）；03 → v0.3.8、05 → v0.1.10（RR-010 已生效）；**C6 金絲雀期限 2026-10-26（UTC+8）**；**C7 每週檢查上游（第一次最遲 2026-10-26），發行後 14 日內換回，最遲 2026-12-31 重評** |

## 編號清單（供品保追溯）

- SEC-001～SEC-017  
- TM-001～TM-023  
- RR-001～**RR-010**  
- SCR-001～SCR-005（見 PR #2 審查）；SCR-006～SCR-011（含 **SCR-010 低、開放、G5 前待辦**）、**SCR-017～SCR-018**（見 11）；SCR-012～SCR-016（見 12）；**SCR-019～SCR-021**（見 14）；**SCR-022～SCR-027、SCR-040～SCR-042**（見 17；022／023／040／041 已關閉）；**SCR-028～SCR-033**（見 18；028／029／031 已關閉）；**SCR-043～SCR-046**（見 19；043～045 已關閉；046 已轉研發部）；**SCR-034～SCR-039 為保留編號，未使用**；SCR-047～SCR-049 為 19 號檔區段的保留編號，未使用  

## 交接單：PR #2 安全產出 → 總協調／審查／維運／研發

- 交出角色 → 接收角色：安全部 → 總協調（D-07／D-08 已登記；CR-001）；審查部（PR 可合併結論）；維運／研發（fail-closed workflow）
- 日期：2026-10-07（v0.12 更新）；2026-10-19（v0.14 更新：CR-002／D-13）；2026-10-19（v0.15 更新：#15、#7 合併後同步）
- 交付物：08 審查、09 裁示、03 門檻 v0.3（現 v0.3.8）、RR-006、RR-010、10～18 號檔、本索引
- 已完成：使用者核准 CR-001（D-07）、選 A 並接受 RR-006（D-08）；Rulesets 已設定（取證 [`docs/09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md`](../09-cicd-env/evidence/rulesets-snapshot-2026-10-07.md)）；Dependency graph 已開；dependency-review 已於 [PR #6](https://github.com/444112029012/short-url/pull/6) 恢復 fail-closed，暫行條款結束
- 已合併：PR #2（2026-10-07 13:56:49 UTC+8，`a8246ceec37d44b7b4441825f75f8cb9ab4b6e03`）；PR #6 fail-closed（2026-10-07 14:04:02 UTC+8，`161a4fdd0abf9b242ec22ff2c1cdba4d82ba9f86`；合併時無審查部「可合併」結論，審查部事後補審，偏離記於 05 RR-006）；PR #9 文件同步（2026-10-07 14:28:12 UTC+8，`63b373db390d053331e045438d059b90b84522f0`）
- DR-G4-01：安全部附條件同意（見 10）→ 設計部登記生效條件；研發部實作小 PR 依 C1～C5
- 使用者決定：D-10（規則一 strict=true；**已生效**，維運取證 `docs/09-cicd-env/evidence/ruleset-strict-true-d10-2026-10-07.md` 隨 PR #13 提交（已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e`，連結已生效）；bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補））；D-11（接受 RR-007～009，G5 前再評）
- **已定合併順序：#10 → #12 → #8 → #11 → #7**；每個 PR 在前一個合併後都須更新分支並取得新四綠（D-10），並依 D-08／D-09 留言程序合併
- PR #7／#8：見 11；**#7 head `ef9fcb8` 已複核：邏輯上核可**（SCR-006／007／008 修正已驗證）；#7 仍排最後，須在 #8、#11 合併後更新分支、修正兩者測試舊簽章並補 SCR-017 前置條件註解，新 head 差異送安全部核對、取得四項新綠燈後，依 D-08／D-09 留言程序合併；合併後 SCR-002／006／007／008 關閉、RR-007 降為低。**v0.15：已完成**——#7，2026-10-19 12:11:10 UTC+8，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`，合併時 head `fc9a7a7`，main 四綠；上述 SCR 與 SCR-017、SCR-011 已關閉，RR-007 已降為低
- PR #11：見 12；#8 合併後在最新 main 重取四綠；與 PR #7 的測試簽章衝突由 #7 解決（SCR-012）；維運補觀測日誌保存期／存取控制（SCR-013／015）
- **CR-002（D-13）**：安全部同意（附條件），見 13；研發部／維運部依 C1～C5 開 CR-002 PR（go.mod `go 1.27.2`、setup-go `"1.27.2"`、gosec 完整 pseudo-version＋`go version -m` h1 比對、門檻不變），四綠並依 D-08／D-09 留言程序合併；PR #14 關閉；合併後 7 日內維運做 CI 金絲雀（C6）；每週檢查上游 gosec 發行（C7）；總協調於 `docs/99-changelog/` 登記 CR-002。**v0.15：已合併**（#15，2026-10-19 10:22:50 UTC+8，merge commit `c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）；C6 期限 **2026-10-26（UTC+8）**；C7 第一次週檢最遲 2026-10-26，其後每 7 日，發行後 14 日內換回，最遲 2026-12-31 重評
- 已合併（補）：PR #13（strict=true 取證檔，head `6674a59`）已於 2026-10-19 09:05:27 UTC+8 合併，merge commit `2e54f7e2d04c891886a4e4210f5880475c9f0d18`（cursor-github 唯讀查證），取證檔連結已生效
- 剩餘：本輪 06-security 更新（10 號檔、11 v0.1.4、12 v0.1.4、13 v0.1.1、05 v0.1.9、03 v0.3.7、README v0.14）尚未進版控，需維運另開文件 PR；bypass 清單 API 層佐證（具 admin 讀取權限）待補；文件 PR B（`docs/ENG-020-g4-cr001-docs-sync`）已以 PR #9 合併
- 已合併（v0.15 補，cursor-github／公開 API 唯讀查證）：PR #15（2026-10-19 10:22:50 UTC+8，`c341be1a390a0a0eeed057a4e862d6cd1a9bd38b`）、PR #10（10:35:26，`f8c4011206069714d0a37a5bd6c5ba4daad3c26f`）、PR #12（10:49:28，`28fc4f7c8cd4816da177e7dbdd9c8f7c5f681786`）、PR #8（11:21:48，`e9f43968236878eb605e8e3048bae59a35815e70`）、PR #11（11:34:13，`3d3a9cddca0c99ee263905d7b6ef847f7f6b913b`）、PR #7（12:11:10，`24b8da0252ac8cbf1c5615592eb9ad18996f3f36`）、PR #16（13:17:22，`ed620394b1c9986cc4ac243eb1543c8b677d92cb`，v0.18 補）、PR #17（14:00:38，`433140c5132778bf7cd731e3ecaca8187e50d678`，v0.19 補）、PR #18（14:24:23，`643bf8ab052ebb175e35bc365ec97863ec9f4202`，v0.20 補）、PR #19（14:49:27，`141fcda38e1cb268f04aacdeba28db1c86e7ceca`，根目錄 README，v0.20 補）、PR #20（15:15:58，`abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`，v0.21 補）；時間皆 UTC+8
- **版控現況（v0.22，2026-10-22 08:58 UTC+8 查證）**：專案已於 2026-10-22 恢復。文件 PR C（#21）已於 2026-10-22 08:57:03（UTC+8）合併，merge commit `db51cd53d775c398dc2fbfcd9892f2052ba9bf36`（parents `abbec9e`＋`3f634f1`）；repo main 為 `db51cd53d775c398dc2fbfcd9892f2052ba9bf36`。**06-security 文件已進版控**：合併當下 repo 的 `docs/06-security/` 25 個檔與本機（v0.22 修改前）sha256 逐檔一致。**本機待進版控**：08 號檔 v0.1.5 與本索引 v0.22（回應審查部 S8）；其餘檔與 repo 一致。bypass 清單 API 層佐證待補
- 歷史（v0.21）：**剩餘（v0.21，版控現況照實記錄）**：截至 2026-10-19 15:18（UTC+8），repo main `abbec9e`（#20 合併，15:15:58；專案已暫停）的 `docs/06-security/` 最後一次變更是 commit `3eb69e6`（2026-10-07）；**10～19 號檔都不在 repo**；03（repo 為 v0.3.3）、05（repo 為 v0.1.5）、README（repo 為 v0.7）在 repo 中是舊版；00～02、03 v0.1／v0.2、04、04a、06～09 與 repo 一致。本輪更新後的 03 v0.3.8、05 v0.1.10、10 v0.1.2、11 v0.1.6、12 v0.1.5、13 v0.1.2、14 v0.1、15 v0.1、16 v0.1.2、17 v0.1.2、18 v0.1.4、19 v0.1.1 與本索引 v0.21，**全部待維運以文件 PR C 進版控**（11 號檔的 SCR-010 定義會隨 PR C 一起進 repo）；bypass 清單 API 層佐證待補

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | G1 資安章節初建 |
| v0.2 | 2026-10-07 | 安全部 | G2：STRIDE／殘餘／前置；SEC v0.2 |
| v0.3 | 2026-10-07 | 安全部 | G3：CI 門檻 v0.2；G3 前置 |
| v0.4 | 2026-10-07 | 安全部 | PR #2：審查 08、裁示 09、門檻 v0.3、RR-006 |
| v0.5 | 2026-10-07 | 安全部 | 依 D-07／D-08 同步索引：門檻 v0.3 已核准、v0.2 已取代、08／09 已核准、RR-006 已接受；記 PR #4 fail-closed 履行狀態 |
| v0.6 | 2026-10-07 | 安全部 | 維運交叉比對修正：Rulesets 證據改指 repo 檔；fail-closed 改為統一句型「已恢復 fail-closed，暫行條款結束」（含 fail-closed PR 連結佔位字串）；新增命名規則（小版修訂不改檔名（檔名保留主版號））；索引版本同步（03 v0.3.2／v0.2.3、05 v0.1.3、06 v0.1.1、08 v0.1.3、09 v0.1.3） |
| v0.7 | 2026-10-07 | 安全部 | 佔位字串替換為 PR #6；索引同步（03 v0.3.3、05 v0.1.4、08 v0.1.4、09 v0.1.4）；交接單更新：PR #2、PR #6 已合併，剩文件 PR B；註記 PR #6 事後補審 |
| v0.8 | 2026-10-07 | 安全部 | 新增 10：DR-G4-01 本文上限安全複核（附條件同意，C1～C5）；本章狀態與交接單同步 |
| v0.9 | 2026-10-07 | 安全部 | 新增 11：PR #7／#8 安全審查（SCR-006～011）；本章狀態、編號清單、交接單同步。補記：05 索引列於 v0.7 之後已由安全部更新為 v0.1.5（D-09），先前漏記於本修訂紀錄 |
| v0.10 | 2026-10-07 | 安全部 | 新增 12：PR #11 安全審查（SCR-012～016）；本章狀態、編號清單、交接單同步 |
| v0.11 | 2026-10-07 | 安全部 | 依 D-10／D-11：索引同步（03 v0.3.4、05 v0.1.6、11 v0.1.1、12 v0.1.1）；RR-007～009 已接受（D-11）、規則一 strict=true 已決定（D-10，取證待維運）；交接單補合併順序 #10 → #12 → #8 → #11 → #7、PR #9 已合併（`63b373d`）；小版修訂不改檔名（檔名保留主版號） |
| v0.12 | 2026-10-07 | 安全部 | PR #7 複核（head `ef9fcb8`）：11 → v0.1.2（邏輯上核可；SCR-006／007／008 修正已驗證、待合併後關閉；新增 SCR-017、SCR-018）；05 → v0.1.7（RR-007 複核通過、合併至 main 後降為低，FIFO 殘餘限制）；D-10 strict=true 已生效：03 → v0.3.5（§7.1 引用維運取證）、11 SCR-011／12 → v0.1.2 SCR-012 改「strict=true 已生效；待依順序合併後關閉」；本章狀態、編號清單、交接單同步；小版修訂不改檔名（檔名保留主版號） |
| v0.13 | 2026-10-19 | 安全部 | bypass 寫法同步：03 → v0.3.6（§7.1／§11）、05 → v0.1.8（RR-006 回應審查部 #13 S1：D-10 後規則一無 bypass 仍成立，證據強度為使用者 UI 畫面、API 未能佐證，接受狀態不變）、11 → v0.1.3、12 → v0.1.3：「待具權限讀取確認」改為「bypass 清單：PASS（依使用者 2026-10-07 14:33 UTC+8 的 UI 畫面，Bypass list 為空；API 因權限不足未回傳 bypass_actors，未能佐證，API 層佐證待補）」，引用 `ruleset-strict-true-d10-2026-10-07.md`；「維運 PR A、待合併」改為「隨 PR #13 提交（截至 2026-10-19 仍 open、未合併）」；11／12 strict=false 敘述加註 D-10 生效時點；索引、交接單同步；10 未變；小版修訂不改檔名（檔名保留主版號） |
| v0.14 | 2026-10-19 | 安全部 | CR-002／D-13：新增 13 號檔索引列（v0.1.1，安全部已同意（附條件））；03 → v0.3.7（Go 精確 1.27.2、gosec 暫行例外 pseudo-version＋h1 比對、期限與退場；待 CR-002 PR 合併後生效）、05 → v0.1.9（新增 RR-010 低）；本章狀態加 CR-002 列、CI 門檻列同步；編號清單 RR-001～RR-010；交接單補 CR-002 待辦；PR #13 已於 2026-10-19 09:05:27 UTC+8 合併（merge commit `2e54f7e`，cursor-github 查證），03、05、11、12 與本索引中「截至 2026-10-19（09:02）仍 open／未合併」改為已合併事實、取證檔連結已生效，11、12 → v0.1.4；歷史修訂列不改；小版修訂不改檔名（檔名保留主版號） |
| v0.15 | 2026-10-19 | 安全部 | #15、#7 合併後同步：03 → v0.3.8（CR-002 已生效）、05 → v0.1.10（RR-007 降為低、RR-010 已生效，C6 2026-10-26／C7 期限）、10 → v0.1.1（C1 已達成、C3 #17 實作中）、11 → v0.1.5（SCR-002／006／007／008／011／017 已關閉）、12 → v0.1.5（SCR-012 已關閉）、13 → v0.1.2（引用 14、CR-002 已生效）、16 → v0.1.2（合併後註記）；索引新增 14、15、16 號檔；狀態總表 PR #7／#8、PR #11、RR-007、CI 門檻、CR-002 改已生效／已關閉，DR-G4-01 更正（stub 不再成立、MaxHeaderBytes #17 實作中）；新增 SCR-010（低、開放、G5 前待辦，定義於 11 號檔）與 SCR-019～021 至狀態總表／編號清單；交接單補合併紀錄與版控現況（10～16 不在 repo，待 PR C）；小版修訂不改檔名（檔名保留主版號） |
| v0.16 | 2026-10-19 | 安全部 | 索引新增 17（PR #17 安全確認，v0.1，通過，head `3b1aa5b`）、18（PR #18 安全確認，v0.1，附條件通過，head `9370def`，必改 SCR-028）；狀態總表新增 PR #17、PR #18 結論與合併條件摘要、新增「G5 前待辦」列（SCR-010、SCR-030、MaxHeaderBytes、SCR-013／015）；DR-G4-01 列與 10 號檔索引列的 MaxHeaderBytes 改為「#17 已實作，待合併」；11 → v0.1.6（SCR-003 已關閉，#8 merge `e9f43968236878eb605e8e3048bae59a35815e70`），索引與 PR #7／#8 列同步；編號清單加入 SCR-022～032、SCR-040～042，註明 SCR-033～039 為保留編號、未使用；交接單交付物與「剩餘」行同步為 10～18 號檔、11 v0.1.6、本索引 v0.16；17、18 號檔內容未修改；小版修訂不改檔名（檔名保留主版號） |
| v0.17 | 2026-10-19 | 安全部 | PR #18 新 head 複核：18 → v0.1.1（head `663ee9d`，**通過**；SCR-028 已修正、SCR-029／031 已處理，皆待合併後關閉；新增 SCR-033 資訊）；索引 18 列、狀態總表 PR #18 列（結論與合併條件）、編號清單（SCR-028～033，保留編號改為 SCR-034～039）、交接單 18 版本與查證時間同步；11 號檔未提及 #15，依指示不改；小版修訂不改檔名（檔名保留主版號） |
| v0.18 | 2026-10-19 | 安全部 | 17 → v0.1.1（PR #17 兩段核對：`3b1aa5b → H 5b5017b → e6d9ab1`；判定通過，僅適用於最終 head `e6d9ab12d62d7f872a7bc18e8761a2a43ca45402`；SCR-040 已修正、SCR-041 已處理，待合併後關閉）；18 → v0.1.2（回應審查部 S4：合併條件改寫為核准綁定 `663ee9d`＋同步 merge 須重新書面確認，測試檔 blob `0e2b6d9f1aa735c0c8b24fd75af672522af72c47`）；索引 10／17／18 列、狀態總表 DR-G4-01／PR #17／PR #18 列、交接單（main 現況改為 `ed62039`、補 #16 合併、17／18 版本）同步；重新查證 main `ed62039` 上 `docs/06-security/` 無變動（最後一次仍為 `3eb69e6`，檔案清單相同），PR C 分類不變；小版修訂不改檔名（檔名保留主版號） |
| v0.19 | 2026-10-19 | 安全部 | #17 合併後同步與 #18 再確認：17 → v0.1.2（#17 於 2026-10-19 14:00:38 UTC+8 合併，merge `433140c5132778bf7cd731e3ecaca8187e50d678`；SCR-022／023／040／041 已關閉）；10 → v0.1.2（C2～C5 已達成，C3 改為已合併）；18 → v0.1.3（同步後新 head `f547a6ebc2f0512c2ef7b11f739edb03ffb50f0f` 書面再確認通過，核准只適用此 SHA）；05 無 DR-G4-01／#17 相關項目，未改；索引 10／17／18 列、狀態總表 DR-G4-01／PR #17／PR #18／G5 前待辦列、編號清單、交接單（main 改為 `433140c`、補 #17 合併）同步；重新查證 main `433140c` 上 `docs/06-security/` 無變動（最後一次仍為 `3eb69e6`，檔案清單相同），PR C 分類不變；小版修訂不改檔名（檔名保留主版號） |
| v0.20 | 2026-10-19 | 安全部 | 索引新增 19（PR #20 安全意見，v0.1，通過，僅適用於最終 head `44b099c227cac2ab24380f4346654540da6535de`；回應審查部 G4-PR20-R3）；18 → v0.1.4（#18 於 2026-10-19 14:24:23 UTC+8 合併，merge `643bf8ab052ebb175e35bc365ec97863ec9f4202`，SCR-028／029／031 已關閉）；狀態總表 PR #18 改為已合併、新增 PR #20 列；編號清單加 SCR-043～046（見 19）與 SCR-047～049 保留；交接：main 改為 `141fcda`（#19 合併，14:49:27 UTC+8），已合併清單補 #18、#19，PR C 待進版控清單加 19；repo main `docs/06-security/` 仍未變（最後變更 `3eb69e6`），PR C 分類不變（19 屬新進 repo）。小版修訂不改檔名（檔名保留主版號） |
| v0.21 | 2026-10-19 | 安全部 | 19 → v0.1.1（#20 於 2026-10-19 15:15:58 UTC+8 合併，merge `abbec9ebd4ef18ddbfb7f0b9f7cfeba8f7ba9326`，parents `141fcda`＋`44b099c`，與核准 head 無差異；SCR-043～045 已關閉）；狀態總表 PR #20 改為已合併；編號清單同步；交接：main 改為 `abbec9e`，已合併清單補 #20，PR C 待進版控清單的 19 改為 v0.1.1。專案暫停，未開始新審查。小版修訂不改檔名（檔名保留主版號） |
| v0.22 | 2026-10-22 | 安全部 | 回應審查部 S8：08 → v0.1.5（SCR-002、SCR-003 改已關閉，以 11 號檔為準：#7 `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`、#8 `e9f43968236878eb605e8e3048bae59a35815e70`）；索引 08 列同步；SCR-046 改「已轉研發部」；交接：#21（PR C）已於 2026-10-22 08:57:03 UTC+8 合併、文件已進版控、main `db51cd5`、專案 10/22 恢復，08 v0.1.5 與 README v0.22 為本機待進版控。小版修訂不改檔名（檔名保留主版號） |
| v0.23 | 2026-10-27 | 安全部 | 索引補上 20～32；03 → v0.3.9、05 → v0.1.13、13 → v0.1.3 三列同步（CR-005 已生效、RR-010 結案、RR-011 已接受（D-17）、RR-012 新增）；小版修訂不改檔名 |
| v0.24 | 2026-10-27 | 安全部 | 索引新增 33（PR #30 gitleaks 誤報裁定）；小版修訂不改檔名 |
