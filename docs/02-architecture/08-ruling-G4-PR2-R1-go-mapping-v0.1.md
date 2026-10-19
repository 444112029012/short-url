---
文件：設計裁示：G4-PR2-R1 公開方法簽名之 Go 語言對映
版本：v0.1
狀態：已裁示（設計部）；待審查部複審關閉 G4-PR2-R1
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-001～REQ-012；SEC-001、003、006、013、014、016、017；NFR-001
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
參考：PR #2 head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc`（tree `224c14157733205d2c255d48278b141ddaa47944`）；審查紀錄 `docs/gates/gate-G4-PR2-2026-10-07.md` §一-2、§二-3、§三；方法規格 `06-method-specs-v0.2.md` §13
---

# 設計裁示：G4-PR2-R1（路徑 a）公開方法簽名之 Go 語言對映

## 一、緣由

審查部 G4 PR #2 初審（`gate-G4-PR2-2026-10-07.md`）以條件 2 退回，開立 **G4-PR2-R1**：實作相對方法規格 v0.1 於下列方法**新增 `context.Context` 參數**，且埠介面以 **Go `error`＋sentinel（`ports.ErrNotFound`／`ports.ErrConflict`）** 取代 `Result<_, AppError>`／`Option`，查無 CR 或設計裁示：

- `ShortURLApplicationService.CreateShortURL`／`GetStats`
- `RedirectService.Redirect`
- `URLRepository.Save`／`FindByShortCode`／`Exists`
- `ClickCounter.IncrementAtomic`／`GetCount`

審查部要求擇一：(a) 設計部書面裁示屬「Go 語言繫結、非規格差異」並於方法規格新增「Go 實作對映」節、升版；或 (b) 提 CR-002。本文件為路徑 (a) 之書面裁示。

## 二、依據

| 依據 | 內容 |
|---|---|
| 方法規格 v0.1 | `06-method-specs-v0.1.md`：表頭「型別為語言中立偽碼」；§6 `findByShortCode`「實作異常時由上層轉 internal_error」、`save`「唯一衝突由呼叫端重試」 |
| D-04 | Go＋chi、SQLite、base62 長度 8 CSPRNG（ADR-001～003 已採納）→ 規格須有 Go 繫結 |
| 循序圖 v0.1 | SEQ-01～07 之呼叫順序（特別是 SEQ-03 先 find 再 increment、成功後 302；SEQ-04 not_found 不 increment） |
| OpenAPI v0.1.0 | 三端點狀態碼、欄位、錯誤四碼 `{error:{type,message}}`、`maxLength: 2048`、`additionalProperties: false` |
| API 設計說明 v0.1 §3 | 錯誤類型 ↔ HTTP：invalid_url 400／not_found 404／rate_limited 429／internal_error 500 |
| SEC-016／017 | 參數化查詢、原子寫入、DB 最小權限（方法規格 §6～§8） |

## 三、比對範圍與方法

- **標的**：`https://github.com/444112029012/short-url/pull/2`，head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc`（與安全審查版本 `3061289` 同 tree `224c141…`）。
- **檔案**：`internal/ports/store.go`、`internal/application/service.go`、`internal/application/redirect.go`、`internal/domain/{errors,types,validator,generator}.go`、`internal/httpapi/{adapter,mapper,response,router}.go`、`internal/infra/sqlite/store.go`、`internal/ratelimit/guard.go`、`internal/observability/hooks.go`、`cmd/shorturl/main.go`。
- **方法**：本機唯讀 checkout 至該 head 逐行閱讀；以全文搜尋確認 ctx 之全部使用點（含有無 `context.WithValue`／`ctx.Value`）；逐方法比對規格之參數、回傳、前置／後置條件、例外，與實作之簽名、輸入驗證、重試、回傳欄位、錯誤對映（HTTP 狀態／error type）、計數原子性、302 與計數順序。
- **未做**：未修改程式、未 commit／push、未變更 head、未於 PR 留言或 review。

## 四、逐項結論

| # | 比對項 | 規格 v0.1 | PR #2 實作（檔案:行） | 結論 |
|---|---|---|---|---|
| 1 | ctx 用途 | （無） | 全部 ctx 皆源自 HTTP `r.Context()`（`adapter.go:87、113、128`），經應用服務原樣下傳至 `database/sql` 之 `ExecContext`（`store.go:123`）／`QueryRowContext`（`store.go:136、155、167、182`）；`Migrate` 用 `context.Background()`（`main.go:40-41`）。全 repo **無** `context.WithValue`／`ctx.Value`；ctx 錯誤落入「其他 error」→`internal_error` | **僅取消／期限傳遞，未攜帶業務值、未改語意** → 對映（GO-MAP-03） |
| 2 | 埠層錯誤契約 | save：internal_error；唯一衝突呼叫端重試／findByShortCode：Option，異常上層轉 internal_error／exists：boolean，internal_error 可選／incrementAtomic、getCount：not_found、internal_error | `ports.ErrConflict`（`store.go:125-127`）↔ 唯一衝突；`ports.ErrNotFound`（`store.go:139-141、168-170、183-185`）↔ None／not_found；其他 error ↔ 實作異常。應用層一對一轉換：`ErrNotFound`→`not_found`（`redirect.go:31-32、40-42`；`service.go:107-108`）、`ErrConflict`→重試（`service.go:85-87`）、其他→`InternalWrap`→`internal_error`（`service.go:69、88、110`；`redirect.go:34、43`） | **sentinel 一對一對應規格例外**；對外錯誤碼不變 → 對映（GO-MAP-04／05／06） |
| 3 | createShortUrl 行為 | 驗證→generate→exists→save；有限次重試；invalid_url／internal_error（含重試耗盡）；click_count=0；回 short_code／short_url／long_url | `service.go:50-97`：同序；上限 8（:13）；耗盡 `Internal()`（:96）；generator 失敗即回 internal_error（:60-63）；click_count 0（:81）；回三欄（:90-94） | 一致 |
| 4 | getStats 行為 | 驗證短碼；invalid_url／not_found／internal_error；無寫入 | `service.go:100-116` | 一致 |
| 5 | redirect 行為與順序 | 先確認存在再原子 +1；失敗次數不變、不回任意 Location；Location＝庫值 | `redirect.go:25-46`：validate→find→increment→回庫值；adapter 錯誤時刪 Location（`adapter.go:117`）、成功 302（:121-123）；計數於 302 回應寫出前完成，與 SEQ-03 同序 | 一致 |
| 6 | 計數原子性（SEC-016） | 單句 `UPDATE … click_count + 1 WHERE short_code = ?` | `store.go:51-55` 單句 UPDATE … RETURNING；參數綁定 | 一致 |
| 7 | 參數化／欄位最小化（SEC-016／017） | 全部參數化；SELECT 僅業務欄位；DB 不對外監聽 | `store.go:38-58` 固定 SQL＋`?`；DB 檔 0600、無網路監聽（:60-102） | 一致 |
| 8 | 輸入驗證 | 空、不可解析、非 http(s)、>2048→invalid_url；短碼 `^[A-Za-z0-9]{8}$` | `validator.go:23-50`（另見 §五-1、2） | 一致 |
| 9 | 錯誤 → HTTP | invalid_url 400／not_found 404／rate_limited 429／internal_error 500；未知→internal_error；message 不洩漏 | `mapper.go:28-77` | 一致 |
| 10 | 命名與 Result | camelCase、Result／Option | PascalCase；`(T, *domain.AppError)` | 語言慣例 → 對映（GO-MAP-01／02） |
| 11 | 未列差異之方法 | validateLongUrl／validateShortCode／generate／toHttpResponse／checkCreate／checkRedirect／on* | 簽名一致（審查紀錄 §二-3 已判一致） | 一致（RateLimitGuard 行為延後 ENG-011、Obs 延後 ENG-012，屬既有追蹤，非本裁示範圍） |

**綜合**：所有差異僅為 Go 語言對映——ctx 僅用於取消／期限傳遞；sentinel error 一對一對應規格例外，且由應用層轉回規格之 AppError 穩定碼；輸入、輸出、錯誤類型、HTTP 狀態、副作用及其順序與規格 v0.1、循序圖、OpenAPI 一致。**未發現語意變更。**

## 五、次要觀察（不影響本裁示，記錄供後續）

| # | 觀察 | 說明 | 建議 |
|---|---|---|---|
| 1 | 長網址驗證細化 | `ValidateLongURL` 另拒 ASCII 控制字元／空白、非 UTF-8、空 host（`validator.go:24-40`） | 對齊 OpenAPI `format: uri`（RFC 3986），屬「不可解析」之細化；已記入方法規格 v0.2 §13.4。品保可視需要補 TC（含空白之 URL→400） |
| 2 | 建立本文上限 8192 bytes | `adapter.go:18、160-162`。2048 字元之 URL 若客戶端以大量 `\uXXXX` 跳脫（或多位元組字元加跳脫）編碼，JSON 本文可能超過 8192 bytes 而回 400；常規編碼（含 2048 ASCII 邊界）不受影響 | 研發部可評估調高上限（例如 ≥ 2048×12＋封包）或於 API 說明註記；非 G4 阻擋 |
| 3 | 類別圖 v0.1 落差 | `ShortUrlApplicationService` 未列 counter 相依；`generate()` 未列 internal_error | 以方法規格為準（實作已依之）；類別圖下次修訂對齊（設計部） |
| 4 | 測試缺口 | 重試耗盡→internal_error、generator 失敗經 CreateShortURL→internal_error 無單元測試 | 同審查紀錄 S3（研發部） |
| 5 | Redirect 防禦檢查 | 庫值空或含控制字元→internal_error（`redirect.go:36-38`）；正常建立路徑不可達 | 符合「失敗不回任意 Location、不計數」；無需動作 |
| 6 | ctx 無業務逾時 | 僅依 HTTP 請求 ctx（斷線／Shutdown 取消）；DB 層另有 `busy_timeout=5000`（`store.go:87`） | 如日後加 `context.WithTimeout`，屬 GO-MAP-03 對映範圍，逾時仍須對映 internal_error |
| 7 | 下游引用檔名 | `08-engineering/01-task-breakdown`、`03-definition-of-ready`、`07-testing/06-test-strategy`、`07-testing/04-design-testability-review` 等仍引用 `06-method-specs-v0.1.md` | v0.1 已加註「請以 v0.2 為準」；各部門下次修訂時更新引用 |

## 六、裁示

1. G4-PR2-R1 所列之 `context.Context` 新增參數與埠層 Go `error`＋`ports.ErrNotFound`／`ports.ErrConflict` 契約，**屬 Go 語言對映，不改變行為語意，不構成方法規格差異，不需 CR**（不走 CR-002）。
2. 對映規則已入方法規格 **v0.2 §13「Go 實作對映」**（GO-MAP-01～09＋逐方法對映表＋錯誤轉換總表），自本日起為方法規格之規範性內容；後續 PR 之 G4 條件 2 比對，以「方法規格 §1～§12 經 §13 對映後之簽名」為準。
3. 本裁示之效力以 PR #2 head `2698138b6c779a2ce8657bdb5b82c43b3b7209dc` 為對照基準；若 head 變更，僅在新 head 仍符合 §13 規則時沿用，由審查部於複審時確認。
4. PR #2 程式碼與 head **無需變更**。

## 七、後續

| 項目 | 負責 |
|---|---|
| 審查部依本文件與方法規格 v0.2 複審，關閉 G4-PR2-R1 | 審查部 |
| 退回項追蹤表 G4-PR2-R1 狀態更新、轉知審查部 | 總協調 |
| §五 次要觀察 2、4 評估／補測 | 研發部 |
| §五 次要觀察 3 類別圖對齊 | 設計部 |

## 簽署

| 角色 | 簽署 | 日期 |
|---|---|---|
| 設計部 | 設計部 | 2026-10-07（Asia/Taipei） |

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 設計部 | G4-PR2-R1 路徑 (a) 書面裁示：屬 Go 語言對映、不改語意、不需 CR；對應方法規格 v0.2 §13 |
