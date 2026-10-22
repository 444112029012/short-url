---
文件：PR #26（OpenAPI 契約測試 CT-*，ENG-013）安全確認，以及 PR #25 合併後的 SCR 關閉紀錄
版本：v0.1
狀態：**PR #26 可合併**，只適用於 head `49e29d80d1f1af227d83c40c0ab212f4372029f0`
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22（14:40）
對應：PR #26、PR #25、檔 21 v0.1.1、檔 25
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

## 1. PR #25 合併後關閉 SCR
| 檢查 | 結果 |
|---|---|
| merge commit | `d0749f57038aa0f56e4fa26b2c2b3c97ceb71265`，parents 為 `5084d6a` 和 `e54d9003…` |
| 合併的 head | 與檔 21 v0.1.1 綁定的 head 相同 |
| workflow blob | `287a9b88…`，與審查時相同 |
| main push run | 39734240648，completed / success |

結論如下：
- **SCR-004、SCR-019、SCR-021 關閉。**
- **SCR-032 關閉**，但附一個條件：部分證據是維運部轉述的，G5 前由安全部自行下載 artifact 核對，不符就重開。
- SCR-050（低）仍開放，G5 前要修補。
- SCR-051、SCR-052（資訊）仍開放。

## 2. PR #26 新 head 核對（公開 clone 加未驗證身分的 API）
| 檢查 | 結果 |
|---|---|
| head | `49e29d80d1f1af227d83c40c0ab212f4372029f0`，parents 為 `e224253` 和 `d0749f5`（main） |
| remerge-diff | 空 |
| 與 main 的差異 | 只新增 `internal/httpapi/contract_test.go`（+258） |
| go.mod / go.sum | 沒有變更 |
| CI run 39734581059 | 四項 success，head_sha 相符 |
| `-json` 結果 | `/workspace/dev/pr-drafts/pr26-ct.json`，sha256 為 `121d219e3ce083a225f023c1e6c2048cef91267f9119f49355d84b19951afca6`，內含 14 個 pass，沒有 fail 或 skip。此檔由研發部或總協調提供，屬轉述證據 |

## 3. 安全審查
- **import 範圍：** 只用 stdlib（net/http、httptest、net/url、path/filepath、regexp、sort、strings、testing）和專案內部的 observability、ratelimit 套件。
- **外部存取：** 沒有網路、exec、環境變數或機密存取。
- **檔案存取：** 只在 `t.TempDir()` 裡建立資料庫。
- **正式程式碼：** 沒有改動，也沒有動 workflow。
- **安全相關斷言：** CT-REDIR-01 驗證附帶 `?url=evil` 時 Location 不會被改寫（對應 SEC-003 開放重導），也驗證 nosniff、ErrorResponse 沒有多餘欄位和 internal_error。這些對 SEC-001、003、014 有正面幫助。
- **新發現：** 無。

## 4. 檔 25（kin-openapi SCA）現況
- **結論維持不同意。**
- 研發部改用 stdlib 手寫斷言，沒有引入 `github.com/getkin/kin-openapi`。在 head `49e29d8` 的 go.mod、go.sum 和整個 tree 裡都查不到這個套件。
- 因此不需要風險接受。檔 25 的申請在 PR #26 中結案。
- 往後如果要再引入，必須重新送 SCA。

## 5. 結論
- PR #26 **可合併**，沒有附加條件，只適用於 head `49e29d80d1f1af227d83c40c0ab212f4372029f0`。
- 合併前如果再有 push，此確認即失效。
- 是否合併由使用者拍板。安全部沒有在 GitHub 上留言或核准。

## 6. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版 |
