---
文件：設計裁定 DR-G4-01（建立請求本文上限）安全部複核意見
版本：v0.1.2
狀態：附條件同意（待補項目見 §4）；**v0.1.1：C1 已達成（main `24b8da0`）；v0.1.2：C2～C5 已達成，C3 已合併（#17，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`，2026-10-19 14:00:38 UTC+8）**
負責角色：安全部
最後更新：2026-10-19
對應：設計部 `docs/02-architecture/06-method-specs-v0.2.md` v0.2.1 §13.4.1（DR-G4-01）；SEC-007（速率限制）、NFR-003；TM（DoS 類）；PR #7（ENG-011）
專案代號：SHORTURL
---

# DR-G4-01 安全部複核意見（v0.1）

## 1. 結論

**附條件同意 DR-G4-01**（本文硬上限 8192 B → 32768 B，超限仍回 400 `invalid_url`）。

- 上限值本身：同意。32768 B 高於理論最壞值（2048 字元 × 12 B ＋外框 ≈ 24.6 KB），單請求記憶體增量（8 KiB → 32 KiB）在有限流、有讀取逾時的前提下屬可接受。
- 不改錯誤類型（不新增 413）：同意。維持 OpenAPI 四類錯誤，回應不帶上限值或內部細節。
- 但 §3 列的 5 項前提中，**2 項在現行 main 不成立**（限流仍為 stub；`MaxHeaderBytes` 未設定），所以本同意**附帶 §4 條件**。條件 C1 達成之前，DR-G4-01 的實作小 PR **不得**合併。
- 本文件只記錄安全部意見，**沒有修改**設計部文件。DR-G4-01 依設計部原文「以安全部意見一致為生效條件」，生效時點由設計部／總協調依本意見登記。

## 2. 審查對象與方法

| 項目 | 內容 |
|---|---|
| 裁定原文 | `docs/02-architecture/06-method-specs-v0.2.md` v0.2.1 §13.4.1（含 §13.4 表「建立請求本文」列） |
| 程式碼基準 A | 現行 `main` @ `161a4fdd0abf9b242ec22ff2c1cdba4d82ba9f86`（2026-10-07 14:20 UTC+8 `git ls-remote` 核對） |
| 程式碼基準 B | PR #7（ENG-011，open、未合併）head `04815e5746fed3c6133bcbc6c27669c3906efc3b` |
| 方法 | 唯讀檢視 `internal/httpapi/adapter.go`、`router.go`、`mapper.go`、`internal/ratelimit/guard.go`、`source.go`、`internal/domain/validator.go`、`internal/application/service.go`、`cmd/shorturl/main.go`；沒有改程式碼，也沒有在 PR 上留言 |

## 3. 前提逐項驗證

| # | 前提（DR-G4-01／總協調指定） | main `161a4fd` | PR #7 `04815e5` | 判定 |
|---|---|---|---|---|
| P1 | 本文在**解碼前**以硬上限截斷讀取（如 `http.MaxBytesReader`），不先讀入全量 | `decodeCreate`：`io.ReadAll(io.LimitReader(body, maxCreateBody+1))`，`len(raw) > maxCreateBody` 回 400，之後才 `json.NewDecoder`（同 PR #7） | `adapter.go:164-166` 截斷讀取，`:167` 起才解碼；`maxCreateBody = 8192`（`:20`） | **成立（等效）**。用的是 `io.LimitReader(max+1)`，不是 `MaxBytesReader`，但已經做到「截斷讀取、不全量讀入、解碼前判斷」。DR 原文寫的是「如」`MaxBytesReader`，所以兩者都接受；建議改用 `MaxBytesReader`（見 C2） |
| P2 | 限流（RateLimitGuard）在讀取／解碼本文**之前** | 呼叫順序成立：`adapter.go:78` `CheckCreate` 在 `:82` `decodeCreate` 之前。但 `internal/ratelimit/guard.go` 仍是 **always-allow stub**（`Guard struct{}`，TODO(ENG-011)），**實際上沒有限流** | 順序成立：`adapter.go:83` `CheckCreate` 在 `:87` `decodeCreate` 之前；`router.go` 的 POST 路由前沒有任何中介層先讀本文；Guard 已實作每來源固定 1 分鐘窗（預設 create 30／redirect 120；非正值回落預設；來源表上限 8192，滿了 fail-closed；被拒不扣額度） | **main 不成立／PR #7 成立**。DR 原文「建立端點另有 30 次／分鐘／來源之速率限制」在 main 上還不是事實，要等 PR #7 合併才成立。另外 DR 引用的 `adapter.go:78` 在 PR #7 後會變成 `:83`（行號漂移，語意不變） |
| P3 | 伺服器設有 `ReadHeaderTimeout`／`ReadTimeout`／`WriteTimeout`／`MaxHeaderBytes` | `cmd/shorturl/main.go:60-63`：ReadHeaderTimeout 5s、ReadTimeout 15s、WriteTimeout 15s、IdleTimeout 60s；**沒有設 `MaxHeaderBytes`** | `main.go:65-68` 同上（PR #7 沒改 server 區塊）；**沒有設 `MaxHeaderBytes`** | **部分成立**。三項逾時已設（ReadTimeout 也限制慢速送本文）。`MaxHeaderBytes` 未設時採 Go 預設 `http.DefaultMaxHeaderBytes` = 1 MiB，比 32 KiB 本文上限大得多，所以標頭是比本文更大的單請求記憶體面 → 待補（C3） |
| P4 | 解碼後仍檢查 `long_url` ≤ 2048 字元 | 同 PR #7 | `service.go:51` `ValidateLongURL` → `validator.go:24` `utf8.RuneCountInString(raw) > 2048` 回 `invalid_url`；這一步發生在 adapter 解碼之後。JSON `\uXXXX`／surrogate pair 解碼後以 rune 計，emoji 計 1 字元 | **成立**。上限放寬不影響業務長度規則 |
| P5 | 超限錯誤回應不洩漏內部資訊 | 同 PR #7 | 超限、讀取錯誤、空本文都回 `domain.InvalidURL(MsgURLValidationFailed)`（`adapter.go:165-166`）；`mapper.go` 只輸出 `{error:{type,message}}`，message 取固定目錄（`URL validation failed`），不帶上限值、實際長度、`err` 字串或堆疊 | **成立** |

### 3.1 DoS 增量評估

- 單請求：讀取緩衝最多 32769 B，加上解碼後字串（≤ 2048 rune，UTF-8 最多約 8 KiB）。每個請求的記憶體峰值仍是 KiB 等級。
- 被限流的請求（PR #7 合併後）在 `decodeCreate` 前就回 429，不讀本文。handler 結束後，`net/http` 為了重用連線最多會排空約 256 KiB 未讀本文，而且受 ReadTimeout 15s 限制，屬有界成本。
- 沒有全域併發上限：限流是單程序、每來源計數。分散來源（大量 IP）的總量攻擊不在本 DR 範圍，維持既有 DoS 類殘餘風險的評估，不因 32 KiB 新增等級。
- 判斷：**限流生效（PR #7）＋現有逾時**這個前提下，風險增量為「低」，可接受。限流是 stub（main 現況）時放寬上限，等於放大無限流的建立端點，不接受（C1）。

## 4. 同意條件與待補項目

| # | 條件／待補 | 類型 | 負責 | 驗收 |
|---|---|---|---|---|
| **C1** | DR-G4-01 實作小 PR（`maxCreateBody` → 32768＋測試）**必須在 PR #7（ENG-011 限流）合併之後**才合併；也可以在同一 PR 內一起帶入限流實作。PR #7 合併前，main 維持 8192 | **阻擋（待補）** | 研發部／總協調 | 實作 PR 的 base 已包含 PR #7 合併 commit；`guard.go` 已不是 always-allow stub。**v0.1.1：已達成**——PR #7 於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`；`guard.go` 已不是 always-allow stub。實作 PR 是否以含此 commit 的 main 為 base，於實作 PR 複核時確認 |
| C2 | 維持「解碼前截斷讀取」；超限回 400 `invalid_url`，message 固定 `URL validation failed`。**建議**改用 `http.MaxBytesReader(w, r.Body, 32768)`（超限時伺服器關閉連線，不再排空剩餘本文）；沿用 `io.LimitReader(max+1)` 也接受 | 必守＋建議 | 研發部 | 程式碼檢視；`ReadAll` 不得無上限。**v0.1.2：已達成**——#17（merge `433140c5132778bf7cd731e3ecaca8187e50d678`）維持 `io.ReadAll(io.LimitReader(body, maxCreateBody+1))`，`maxCreateBody`＝32768，超限 400 `invalid_url`；`TestCreateBodyReadIsBounded`（最多讀 32769 B）、`TestCreateBody32769RejectedInvalidURL`（見 17 號檔 §5） |
| C3 | 在 `http.Server` 明確設定 `MaxHeaderBytes`（建議 ≤ 16 KiB），不依賴 1 MiB 預設 | **待補**（不阻擋 C1 的上限調整；必須在 G5 送審前完成，可併入 DR 實作 PR 或 ENG-016） | 研發部／維運 | `main.go` 有設定值；超大標頭回 431。**v0.1.2：已合併**——#17 於 2026-10-19 14:00:38（UTC+8）合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`；`main.go` `maxHeaderBytes = 16 << 10`、`newHTTPServer` 設定 `MaxHeaderBytes`，21／32 KiB 標頭回 431（`main_test.go`）；安全部 17 號檔 v0.1.1 已確認（歷史：v0.1.1 記為「#17 實作中」） |
| C4 | 測試至少涵蓋 DR §13.4.1 (a)(b)(c)，再加：(d) 恰好 32768 B 的合法本文不會因大小被拒；(e) 被限流來源送超大本文 → 429（不是 400），證明限流在讀本文之前。測試本文使用實作的實際欄位名（`createRequest` 目前是 `json:"url"`） | 必守 | 研發部 | `unit-test` check 綠，測試名可追溯。**v0.1.2：已達成**——#17（merge `433140c5132778bf7cd731e3ecaca8187e50d678`）：(a) `TestCreateC4aSurrogatePairEscaped2048CharsAccepted`、(b) `TestCreateBody32769RejectedInvalidURL`、(c) `TestE2E_02_Length2048Boundary`、(d) `TestCreateBodyExactly32768Accepted`、(e) `TestRateLimitedOversizedCreateGets429WithoutReadingBody`；最終 head run 39204252238 unit-test success（見 17 號檔 §0、§5） |
| C5 | 2048 字元（rune）驗證、Content-Type、未知欄位、多餘 JSON 值規則不得一起放寬 | 必守 | 研發部 | 程式碼檢視。**v0.1.2：已達成**——#17（merge `433140c5132778bf7cd731e3ecaca8187e50d678`）對 `adapter.go` 只改 `maxCreateBody` 常數；`validator.go`、`jsonContentType`、`DisallowUnknownFields`、多餘值檢查都沒有變動（見 17 號檔 §5） |

> 生效建議：設計部可在 C1 達成時，把 DR-G4-01 登記為「安全部附條件同意，C1 已滿足」後生效；C3 列入 G5 前待辦。實作小 PR 由安全部複核 C2／C4／C5。

## 5. 備註（資訊性，不影響結論）

- DR 原文 §13.4.1 以 `{"long_url":""}` 估算外框；實作欄位是 `url`（短 5 B），對 32768 的餘裕判斷沒有影響。欄位名稱與 OpenAPI 是否一致，屬設計部／研發部範圍，本意見不判定。
- DR 引用的 `adapter.go:78`（main）在 PR #7 後變為 `:83`；建議設計部下次修訂時改為引用函式名（`HandleCreate` → `CheckCreate` 先於 `decodeCreate`），避免行號漂移。

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 安全部 | 初版：DR-G4-01 附條件同意；P1 成立（LimitReader 等效）、P2 main 不成立（Guard stub）／PR #7 成立、P3 部分成立（缺 MaxHeaderBytes）、P4／P5 成立；條件 C1～C5 |
| v0.1.1 | 2026-10-19 | 安全部 | 狀態與 §4：C1 已達成（PR #7 於 2026-10-19 12:11:10（UTC+8）合併，merge commit `24b8da0252ac8cbf1c5615592eb9ad18996f3f36`；`guard.go` 已不是 always-allow stub）；C3 註記「#17 實作中」（未合併、未判定）；P2「main 不成立」為 v0.1 當時事實，保留不改；C2／C4／C5 待實作 PR 複核；小版修訂不改檔名（檔名保留主版號） |
| v0.1.2 | 2026-10-19 | 安全部 | #17 合併後同步：C2～C5 改為已達成；C3 由「#17 實作中」改為「已合併」（#17 於 2026-10-19 14:00:38（UTC+8）合併，merge commit `433140c5132778bf7cd731e3ecaca8187e50d678`），依據安全部 17 號檔 v0.1.1 對最終 head `e6d9ab1` 的確認；小版修訂不改檔名（檔名保留主版號） |
