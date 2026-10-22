---
文件：S6 處置紀錄｜PR #18 限流 E2E（TestENG011_*）列入測試矩陣決定
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-22
對應需求：SEC-007（主）；SEC-005、SEC-004、NFR-003、NFR-006、REQ-003（輔）
對應程式：`internal/integration/ratelimit_e2e_test.go`（main `db51cd5`；隨 PR #18 合併，merge commit `643bf8a`）
對應案例：`07-test-cases-v0.1.md`（內容版 v0.2）E2E-05、E2E-05-01～E2E-05-11
對應矩陣：`02-traceability-matrix-v0.1.md`（內容版 v0.4）
專案代號：SHORTURL
關卡：G4（審查部 S6 請品保決定）
---

# S6 處置紀錄：TestENG011_* 列入測試矩陣

## 1. 背景與問題

審查部在 PR #18（限流 E2E）提出 S6：`ratelimit_e2e_test.go` 內 19 個 `TestENG011_*` 測試項目是否列入測試矩陣、對應哪些 TC（SEC-007、NFR-003）。

**「19 個」的組成（實測）**：14 個頂層測試函式＋5 個子測試（`UntrustedPeerCannotSpoofForwardingHeaders` 2 個、`TrustedProxyXFFEdgeCases` 3 個）＝19 個 PASS 結果列。下表逐項處置。

原 E2E-05（v0.1 案例目錄）只有一句：「同來源快速超過建立門檻；When POST；Then 429 rate_limited」，涵蓋不了來源判定、可信代理、IPv6 聚合、並發與日誌等已實作的行為，矩陣也看不到這些證據。

## 2. 處置原則

| 決定 | 條件 |
|---|---|
| **納入** | 驗證產品對外可觀察行為，且可對應 REQ／NFR／SEC → 給正式 TC ID |
| **合併** | 與另一測試驗同一性質（資料列、邊界或前提），併入同一 TC，不另給 ID |
| **不納入** | 不驗產品行為（純測試工具自檢）且無追溯價值 |

E2E-05 保留為**總 ID**（向下相容 DoD P7、任務對照 T-SEC-01、安全部 18 號文件的引用），細分為 E2E-05-01～E2E-05-11。

## 3. 逐項處置（19 項）

| # | 測試（`go test -v` 名稱） | 實際驗證內容（讀碼摘要） | 決定 | TC ID | 對應需求 |
|---|---|---|---|---|---|
| 1 | `TestENG011_CreateQuotaExhaustedReturns429` | 建立門檻 3：前 3 次 201；第 4～6 次 429 `rate_limited`、無 Location／Retry-After／明文 HSTS；被拒不落庫（列數維持 3）；超限時非法 body 仍回 429（先限流再解析） | **納入** | **E2E-05-01**（＝原 E2E-05 主斷言） | SEC-007 |
| 2 | `TestENG011_DefaultLimitsFromConfig` | 未設定環境變數時預設建立門檻 30／分：第 31 次 429 | **納入** | **E2E-05-02** | SEC-007（門檻值與 UT-RL-01 文件值一致） |
| 3 | `TestENG011_RedirectQuotaExhaustedReturns429WithoutCounting` | 導向門檻 2：前 2 次 302；之後 429；超限時未知短碼也是 429；stats 不限流，click_count＝2（被拒不計次） | **納入** | **E2E-05-03** | SEC-007；REQ-003（失敗不計次） |
| 4 | `TestENG011_CreateAndRedirectQuotasAreIndependent` | 建立配額用盡不影響導向，反之亦然 | **納入** | **E2E-05-04** | SEC-007 |
| 5 | `TestENG011_PerSourceIsolationByPeerAddress` | 兩個真實 TCP 來源位址（127.0.0.2／.3），A 用盡不影響 B；CI 下來源位址不可用時 fail-closed | **納入** | **E2E-05-05** | SEC-007 |
| 6 | `TestENG011_SourceAddrUnavailableClassifier` | 只有 `EADDRNOTAVAIL` 可讓 #5 在本機 skip，其餘錯誤一律失敗 | **合併** | 併入 E2E-05-05（防假綠前提） | 不驗產品行為，但沒有它 #5 可能默默 skip；不另給 ID |
| 7 | `TestENG011_UntrustedPeerCannotSpoofForwardingHeaders` | 頂層容器（2 子測試） | **納入** | **E2E-05-06** | SEC-007 |
| 8 | `…/no_trusted_proxies` | 未設可信代理：輪換 XFF／X-Real-IP 仍共用 peer 配額（建立＋導向），列數＝3 | **合併** | 併入 E2E-05-06（資料列 a） | SEC-007 |
| 9 | `…/trusted_list_excludes_peer` | 可信清單不含 peer：同上 | **合併** | 併入 E2E-05-06（資料列 b） | SEC-007 |
| 10 | `TestENG011_TrustedProxyForwardedForIsHonored` | peer 為可信代理：依 XFF 分配獨立配額；無 XFF 用 X-Real-IP；無轉送標頭時算代理本身 | **納入** | **E2E-05-07** | SEC-007 |
| 11 | `TestENG011_TrustedProxyIgnoresClientPrependedEntries` | 可信代理後方：取最右側非可信項，用戶端前置的值（同行或另一行）不能換身分 | **納入** | **E2E-05-08** | SEC-007 |
| 12 | `TestENG011_TrustedProxyXFFEdgeCases` | 頂層容器（3 子測試，SCR-031） | **合併** | 併入 E2E-05-08 | SEC-007 |
| 13 | `…/chain_entirely_trusted_uses_right-most_valid_entry` | 全鏈可信時取最右側合法項 | **合併** | 併入 E2E-05-08（邊界 a） | SEC-007 |
| 14 | `…/invalid_entry_right_of_the_client_is_skipped` | 用戶端右側的非法項略過 | **合併** | 併入 E2E-05-08（邊界 b） | SEC-007 |
| 15 | `…/no_valid_XFF_entry_falls_back_to_X-Real-IP` | XFF 無合法項時退回 X-Real-IP，再退回 peer | **合併** | 併入 E2E-05-08（邊界 c） | SEC-007 |
| 16 | `TestENG011_IPv6ForwardedClientsShareSlash64` | 同 /64 共用配額、他 /64 不共用；IPv4-mapped 與 IPv4 共用（SCR-007） | **納入** | **E2E-05-09** | SEC-007 |
| 17 | `TestENG011_AdjacentIPv6Slash64sDoNotShareQuota` | 相鄰 /64 不合併（/63 實作會失敗）；同 /64 首尾共用（SCR-031） | **合併** | 併入 E2E-05-09（邊界） | SEC-007 |
| 18 | `TestENG011_ConcurrentCreatesAdmitExactlyLimit` | 40 個並發建立、門檻 5：恰 5 個 201、35 個 429，庫內恰 5 列 | **納入** | **E2E-05-10** | SEC-007；NFR-003（輔助：單實例、進程內限流的並發一致性） |
| 19 | `TestENG011_LogsDoNotContainClientAddresses` | 允許／拒絕／偽造／可信轉送流量後：rate_limited 事件只記 error_type（建立 2、導向 1），觀測日誌與 `http.Server.ErrorLog` 都沒有 peer／轉送位址（IPv4、IPv6 各寫法、/64 鍵）及標頭名 | **納入** | **E2E-05-11** | SEC-005、SEC-004（不留訪客位址）、NFR-006（失敗可觀察）；SEC-007 |

**統計**：納入 11 項（11 個 TC）；合併 8 項；**不納入 0 項**。

## 4. 對 SEC-007 與 NFR-003 的結論

- **SEC-007**：TC 改為 `UT-RL-01～02、UT-HC-03、E2E-05（E2E-05-01～11）`。E2E-05-01 滿足原 E2E-05 驗收句；其餘 10 個補上來源判定、防偽造、IPv6、並發與日誌等證據。
- **NFR-003**（不承諾水平擴展；單實例假設文件化）：主 TC 仍是 **CHK-NFR-003**（文件審查）。TestENG011_* 都在單一進程、單一伺服器實例內執行，只能證明「進程內限流在單實例下正確」，不能也不應宣稱跨實例一致性。因此只把 **E2E-05-10 列為輔助證據**，其餘不對應 NFR-003，以免過度宣稱。

## 5. 執行結果（實測）

| 項目 | 結果 |
|---|---|
| 指令 | `go test ./internal/integration/ -run TestENG011 -v -count=1` |
| 環境 | box，go1.27.3 linux/amd64；worktree detached at `db51cd5`（＝origin/main） |
| 時間 | 2026-10-22 約 09:05（UTC+8） |
| 結果 | **19 PASS／0 FAIL／0 SKIP**（14 頂層＋5 子測試）；`ok … internal/integration 0.145s` |
| 補充 | `CI=true` 單跑 `TestENG011_PerSourceIsolationByPeerAddress`：PASS（未 skip）；`go test ./... -race -count=1` 全部套件 ok |
| 原始輸出 | `evidence/s6-eng011-run-2026-10-22.txt` |

> 註：box 的 Go 是 1.27.3，CI 釘 1.27.2（CR-002）；版本差只影響本機證據，CI 證據以 main run（`db51cd5`，ci-devsecops success，2026-10-22 08:57 UTC+8 建立）為準。

## 6. E2E-05 最終處置

**E2E-05：保留總 ID，細分為 E2E-05-01～E2E-05-11；狀態＝已自動化、已通過（main `db51cd5`）；實作＝`internal/integration/ratelimit_e2e_test.go`。** 審查部 S6 可以關閉：

1. 19 個測試項目全部有處置（11 納入、8 合併、0 不納入）。
2. 追溯矩陣 v0.4 的 SEC-007、SEC-005、NFR-003、NFR-006、REQ-003 列已回填。
3. 案例目錄 v0.2 新增 §5.1 E2E-05 細分表。
4. **採認範圍說明**（審查部 #18 S6 提問「不經 `cmd/shorturl` 程序與 `ListenAndServe` 是否可採認」）：品保**採認**。理由：`newRateLimitServer` 用真實 `config.Load`（環境變數）、真實 SQLite、真實 TCP（`httptest` 監聽埠，#5 還做來源位址 bind），組裝順序與 `cmd/shorturl/main.go` `run()` 逐行對照一致（2026-10-22 讀碼）。差異只有日誌輸出目的地（測試導入緩衝區）與 `http.Server` 逾時／`MaxHeaderBytes`（由 `cmd/shorturl/main_test.go` 另測）。組裝碼重複有漂移風險，列 G4 建議 G4-QA-S3（DEF-003）。
5. 剩餘限制（不阻擋）：測試只覆蓋單實例（NFR-003 範圍外）；`Retry-After` 不在現行契約內，測試鎖定「不得出現」，契約若改要同步改測試；安全部開放建議 SCR-030（net/http panic／TLS 訊息含 peer 位址）不在 E2E-05-11 覆蓋範圍內。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 品保部 | S6 初版：19 項逐項處置；E2E-05 細分 01～11；實跑 19 PASS |
