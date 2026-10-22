---
文件：PR #25（維運部 CI 強化）安全審查意見
版本：v0.1.1
狀態：**v0.1.1：有條件可合併，僅適用於同步後最終 head `e54d9003f303e7e1a33082a70308d048caf0fdc1`**；v0.1：head `bf190c4` 有條件可合併
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22（v0.1.1，14:25）
對應：SCR-004、SCR-019（[14](./14-pr-015-cr-002-security-review-v0.1.md) 第 108 行）、SCR-021、SCR-032（[18](./18-pr-018-ratelimit-e2e-security-confirmation-v0.1.md)）
專案代號：SHORTURL
發現項編號區段：SCR-050～SCR-059（本文件使用 SCR-050～SCR-052；SCR-047～049 為 19 號檔保留）
時區：Asia/Taipei（UTC+8）
---

# PR #25 安全審查意見

## 1. 範圍
- PR：https://github.com/444112029012/short-url/pull/25
- head：`bf190c4c9292f15f7c152ad340c07d58bac2fa58`；唯一 parent：`db51cd53d775c398dc2fbfcd9892f2052ba9bf36`。
- 變更只有 `.github/workflows/ci-devsecops.yml`（+48／−20），沒有改到程式碼。
- CI run 39703496002：四個 job 的 head_sha 都是 `bf190c4`，conclusion 都是 success（未認證 API 查得）。

## 2. 判定
| SCR | 判定 | 依據 |
|---|---|---|
| SCR-004 | 達成 | 每個 `uses:` 都是 40 碼 SHA，已用 `git ls-remote --tags` 核對：checkout `11d5960a…` 對應 v4.4.0、setup-go `40f1582b…` 對應 v5.6.0、dependency-review `2031cfc0…` 對應 v4.9.0、upload-artifact `ea165f8d…` 對應 v4.6.2 |
| SCR-021 | 達成 | gitleaks 8.21.2 linux_x64 的 sha256 `5bc41815…e3ba` 與官方 checksums.txt 一致；`set -euo pipefail`、`curl -f`、先校驗再解壓，校驗不符就失敗 |
| SCR-019 | 達成 | 改以 `"$(go env GOPATH)/bin/gosec"` 呼叫 |
| SCR-032 | 達成（部分證據為轉述） | 安全部自查：有 -race、gosec -tests、go-test-json artifact（14679 bytes），ENG-011 不符 `--- PASS` 就 exit 1，SKIP 會判失敗；本機 go1.27.3 跑 `go test -race -json -run ENG011` 為 pass（本機結果，非 CI）。維運部提供（安全部尚未獨立取得）：日誌 01:05:13Z 的 `--- PASS: TestENG011_PerSourceIsolationByPeerAddress (0.04s)`；go-test.json（sha256 78bc49d76dc9318f12d1f23919a97f88cf919caf4f28f56ddf3c7e4d83e1b10a）終態 Action 為 pass。G5 前安全部會自行下載 artifact 核對，不符則重開 SCR-032 |
| 權限與注入 | 符合 | 權限只有 `contents: read`、`pull-requests: read`；沒有 pull_request_target；run 區塊沒有 `${{ }}` 插值 |

## 3. 新發現項
- **SCR-050（低）**：ENG-011 的檢查是對 `-v` 的文字輸出做 grep，eng011.log 沒有上傳；artifact 設 `if-no-files-found: warn`。影響：日誌過期後就沒有「PASS 不是 SKIP」的證據，且結果檔遺失也不會讓 job 失敗。重現：列 artifacts 只有 go-test-json；把 tee 的路徑改壞，job 仍是綠的。修補建議：上傳 eng011.log、改成 `error`、對 go-test.json 用 jq 斷言 `"Action":"pass"`。期限：G5 前。
- **SCR-051（資訊）**：沒有 github-actions 的 dependabot 設定，SHA 釘選後收不到 action 的安全更新。
- **SCR-052（資訊）**：dependency-review 那行註解寫 `# v4 (v4.9.0)`，但 v4 tag 不指向該 SHA，應改寫 `# v4.9.0`。

## 4. 結論
有條件可合併，僅適用於 head `bf190c4c9292f15f7c152ad340c07d58bac2fa58`。條件：SCR-050 在 G5 前以後續 PR 關閉；若因 strict 須同步 main，須依「只同步」範圍（parent、remerge-diff 為空、blob 不變、四項綠且 head_sha 相符）重新核對。合併由使用者拍板。安全部未在 GitHub 留言或核准。

## 5. v0.1.1：同步後最終 head 核對（2026-10-22 14:25，公開 clone）
| 檢查 | 結果 |
|---|---|
| head | `e54d9003f303e7e1a33082a70308d048caf0fdc1`（`git fetch origin pull/25/head` 結果相符） |
| parents | `bf190c4`（v0.1 審查 head）＋`5084d6a95a0ac1d1a3884e990089fa742822c3b1`（main） |
| `git show --remerge-diff e54d900` | 空（0 bytes） |
| `git diff --name-status 5084d6a e54d900` | 只有 `M .github/workflows/ci-devsecops.yml` |
| workflow blob | `287a9b88a72f2bbae5a6d8f9c5a91fb58c27312e`，與 `bf190c4` 相同，內容沒有變 |
| CI run 39733542195 | 四項都是 success，head_sha 都是 `e54d9003…`；artifact go-test-json 為 14806 bytes，未過期 |

**各 SCR 狀態（以此 head 合併為前提）**
| SCR | 狀態 |
|---|---|
| SCR-004 | 達成，合併後關閉 |
| SCR-019 | 達成，合併後關閉 |
| SCR-021 | 達成，合併後關閉 |
| SCR-032 | 達成，合併後關閉。其中 CI 實際輸出（`--- PASS` 那一行、go-test.json 的終態事件）仍是維運部提供的轉述；G5 前安全部會自行下載 artifact 核對，不符就重開 |
| SCR-050（低） | 開放，G5 前以後續 PR 關閉 |
| SCR-051、SCR-052（資訊） | 開放 |

**資訊（不阻擋）**：合併 commit 訊息和分支名用的是 `ENG-022`（`fix/ENG-022-ci-hardening-scr032`），但總協調轉達時寫的是 ENG-020；ENG-022 在 CR-003 已指派給分支規範的文件工作。請研發部或維運部確認 ENG 編號，以免追溯錯置。這不影響安全結論。

**結論**：有條件可合併，僅適用於 head `e54d9003f303e7e1a33082a70308d048caf0fdc1`；條件同 §4。被合併的 head 必須就是這個 SHA，若再有 push，此確認即失效。合併由使用者拍板，安全部未在 GitHub 留言或核准。

## 6. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版，綁定 head bf190c4 |
| v0.1.1 | 2026-10-22 | 安全部 | 綁定同步後 head e54d900，列出各 SCR 狀態；小版修訂不改檔名（檔名保留主版號） |
