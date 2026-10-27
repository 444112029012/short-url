---
文件：PR #24（研發部：根目錄 README.md 431／本文上限描述）安全書面確認
版本：v0.1.1
狀態：**v0.1.1：通過，僅適用於同步後最終 head `0379b64c892e9f664e07c879a16f9acaa164be88`**（回應審查部 G4-PR24-R1）；v0.1：head `80fc321` 通過
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22（v0.1.1，13:20）
對應：[19](./19-pr-020-openapi-011-security-opinion-v0.1.md) SCR-046；OpenAPI 0.1.1（docs/04-api/openapi.yaml）
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# PR #24 安全書面確認

## 1. 範圍與處理項目
- PR：https://github.com/444112029012/short-url/pull/24
- head：`80fc321a2d363871bfdf9b50aa89ca10a3be359d`；唯一 parent：`db51cd53d775c398dc2fbfcd9892f2052ba9bf36`（main）。
- **本 PR 實際處理的是 SCR-046**（19 號檔，資訊級：根目錄 README.md 第 56、101 行 431 描述過時）。研發部 PR 描述寫「SCR-046，併入 SCR-043／044」；SCR-043／044 已於 #20（merge commit `abbec9e`）關閉，本 PR 不重開也不再關閉它們，只與其措辭對齊。

## 2. 核對方式與結果（2026-10-22 09:05，公開 clone `git fetch origin pull/24/head`）
- `git diff --name-status db51cd5 80fc321`：只有 `M README.md`，3 行增、3 行刪（第 55、56、101 行）。
- 措辭：與 OpenAPI 0.1.1 一致。上限計算範圍為請求列加標頭合計；16384 B 以下保證接受；超過約 20480 B 或標頭值超過 500 個（含 Host）時回 431；net/http 自行產生的 400／417／431／505 回應「不保證」帶安全標頭，用戶端不得依賴其有無。沒有寫成「不帶」的承諾。本文以原始位元組計算、超限回 400 invalid_url、不回 413。
- CI（未認證 API，run 39703382114）：secrets-gitleaks、sca-govulncheck、sast-gosec、unit-test 四項 conclusion 都是 success，head_sha 都是 `80fc321`。

## 3. 結論
- 安全部同意合併 head `80fc321a2d363871bfdf9b50aa89ca10a3be359d`。無新發現項。
- #24 以此 head 合併後，SCR-046 關閉。
- 若因 strict 須同步 main，新 head 須符合以下條件，安全部才會針對該 head 另出確認：parent 為本 head 與當時 main；`git show --remerge-diff` 為空；README.md blob 不變；四項檢查都綠且 head_sha 相符。
- 合併由使用者拍板。安全部未在 GitHub 留言或核准。

## 4. v0.1.1：同步後最終 head 核對（2026-10-22 13:20，公開 clone）
| 檢查 | 結果 |
|---|---|
| head | `0379b64c892e9f664e07c879a16f9acaa164be88`（`git fetch origin pull/24/head` 結果相符） |
| parents | `80fc321`（v0.1 核對 head）＋`2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5`（main） |
| `git show --remerge-diff 0379b64` | 空（0 bytes） |
| `git diff --name-status 2a02819 0379b64` | 只有 `M README.md` |
| README.md blob | `fbdfb496b2ff44a31907abbc8d572447691cf1cd`，與 v0.1 核對的 head `80fc321` 相同，內容沒有變 |
| CI run 39726831602 | secrets-gitleaks、sca-govulncheck、sast-gosec、unit-test 都是 success，head_sha 都是 `0379b64c…` |

- **SCR-046 狀態：開放，已轉研發部、修補已就緒。** #24 以 head `0379b64c892e9f664e07c879a16f9acaa164be88` 合併後就關閉；合併後安全部會在 19 號檔與 README 登錄關閉。SCR-043／044 已在 #20 關閉，本 PR 與其無關。
- **結論：** 安全部同意合併 head `0379b64c892e9f664e07c879a16f9acaa164be88`。被合併的 head 必須就是這個 SHA；若再有 push，此確認即失效。合併由使用者拍板，安全部未在 GitHub 留言或核准。

## 5. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版，綁定 head 80fc321 |
| v0.1.1 | 2026-10-22 | 安全部 | 回應 G4-PR24-R1：綁定同步後 head 0379b64，並寫明 SCR-046 狀態；小版修訂不改檔名（檔名保留主版號） |
