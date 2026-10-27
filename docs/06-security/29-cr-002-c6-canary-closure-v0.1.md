---
文件：CR-002 條件 C6（CI 金絲雀驗證）關閉紀錄
版本：v0.1
狀態：**C6 關閉**（逾期 1 天，記為偏差，不影響控制有效性）
負責角色：安全部
作者：安全部
日期：2026-10-27
最後更新：2026-10-27（09:15）
對應：13 號檔 §7 C6；14 號檔；維運部證據 docs/09-cicd-env/evidence/cr-002-c6/
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

## 1. 安全部獨立核對
核對方式：未驗證身分的 GitHub API、公開 clone，以及本機證據檔。

| 檢查 | 結果 |
|---|---|
| run 40430102886 | 由 workflow_dispatch 觸發，分支 `canary/c6-gosec-red`，head 為 `99fa4bd6176ecd61af4ac9c744e9356843c0b2f3`，建立時間 2026-10-27 09:01:52 |
| jobs | sast-gosec 為 **failure**；secrets-gitleaks、sca-govulncheck、unit-test 為 success |
| 金絲雀 commit | 訊息為「test(ci): C6 金絲雀假樣本（G404），不合併、驗證後刪分支」。parent 是 main `2425129`，只新增 `internal/c6canary/canary.go` |
| log：完整性驗證 | `mod github.com/securego/gosec/v2 v2.29.1-0.20261015133444-707f2271f665 h1:KvLOP3HN…oj40=`，接著出現 `gosec h1 OK`。C2、C3 的鎖定仍然有效 |
| log：掃描命令 | `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...`（SCR-019、SCR-032 的修正仍然生效） |
| log：偵測 | `canary.go:11 - G404 (CWE-338) … Severity: HIGH`、`Issues : 1`、`exit code 1`。這證明門檻會擋 HIGH |
| 分支 | `git ls-remote origin 'refs/heads/canary/*'` 結果是空的，branches API 也回 404，**分支已刪除** |
| main | 目前是 `2425129`（#27），與金絲雀的 parent 相同，沒有被這次驗證動到。`.github/` 從 f646bee 到 main 沒有變更 |
| PR | 金絲雀 commit 不在 main 上，沒有合併 |

## 2. 證據檔 sha256（本機）
| 檔案 | sha256 |
|---|---|
| run-40430102886.json | `c256707511c0f6a9d51edac162197bbd330963b15d8be06250cbd40df0175303` |
| run-40430102886.log | `c330ddf6831c8601a83ea1ef0b341e539d70d90adc1f82c9e26a0bc6b979aa2c` |
| branch-deleted-note.txt | `8d992d5680d9d66dbe718dea3c86b594c1ef5618a0a380ee7636eb8d1c343d44` |
| branch-deleted-ls-remote.txt | `e3b0c442…b855`（空檔，代表 ls-remote 查無結果） |

## 3. 結論與偏差
- **C6 關閉。** 依 13 號檔，要求的「run 連結（紅燈）＋分支已刪」都已具備，而且安全部已自行核對 log。
- 偏差：C6 期限是 2026-10-26，實際在 2026-10-27 09:01 執行，**逾期 1 天**。在這段期間 main 的 sast-gosec 每次 push 都正常執行並通過，13 號檔 §3.4 R5 的本機金絲雀也已先驗證過，所以逾期沒有造成控制缺口。不另開 SCR。
- 使用者核准這次推送：根據總協調與維運部轉述，時間是 10/27 09:01，在維運部的聊天裡。這是轉述，安全部沒有見到原話。
- 後續 C7：gosec v2.30.0 已發行（轉述）。維運部在起草 CR-005，安全部收到後會評估。依 13 號檔 C7，要在發行後 14 日內換回發行版；換回時一併處理 SCR-019 的後續追蹤。13 號檔會在下一個小版補註 C6 已關閉。

## 4. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 安全部 | 初版 |
