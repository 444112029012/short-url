# CR-005 K4 G404 金絲雀結案（安全部）

- 文件編號：06-security/32
- 版本：v0.1（Draft）
- 負責：安全部
- 關聯：CR-005 v0.3（bfaa6d1d…）、D-19；30 號 v0.1.1 §3 K4；31 號 v0.1.1（c1c1f9ab…）；05 號 RR-010；13 號 C7
- 撰寫時間：2026-10-27 17:02:20（UTC+8，取自 `date`）

## 1. 結論
**K4 成立，K4 關閉。** gosec v2.30.0 在 CI 上實際擋下了 HIGH 等級的 G404 假樣本，阻擋門檻有效。CR-005 的 K1～K5 已全部完成，只剩 K6：升 Go 1.27.3 的 CR 最遲 2026-11-20 提出，對應 RR-012。**RR-010 結案。**

## 2. 背景事實
| 項目 | 內容 |
|---|---|
| CR-005 合併 | PR [#29](https://github.com/444112029012/short-url/pull/29) 於 2026-10-27 16:37:12（UTC+8，取自 git commit 時間）合併，merge commit `db380a40c3330bf23498ae7fda2f123cf7dd4cc4`，父節點為 `36f4717c`＋`1545143c`。合併後 main push run 40474778561 四項皆 success（K3） |
| K4 期限 | 合併後 7 日內，最遲 2026-11-03；實際於 2026-10-27 執行，**在期限內** |
| 執行核准 | 據維運部轉述，使用者 16:58 在維運部聊天核准執行。這是轉述，安全部沒有見到原話 |
| 金絲雀 commit | `e9db9e4951a8363999a7be90670c1b27ba817e5b`，parent 為 main `db380a40`，只新增 `internal/k4canary/canary.go`（sha256 `eab0217e78a62e20df70882b1232bd87ef944f0441438b48b31ba45be71976c6`）。跟 evidence/cr-005/11 號樣本相比，只有第 1、2 行（註解、package 名稱）不同，第 1 行沒有 `#nosec`。安全部推送前已用本機 gosec v2.30.0 預檢：G404 HIGH、Issues 1、exit 1（本機是 Go 1.27.3，不算證據） |

## 3. 查證（安全部自行執行）
| 項目 | 結果 |
|---|---|
| run 40477332321（公開 API） | workflow_dispatch，head_branch `canary/k4-gosec-v2300-red`，head_sha `e9db9e49…`，conclusion failure，建立時間 2026-10-27 16:58:49（UTC+8；API 值 08:58:49Z） |
| jobs | sast-gosec（122098695670）failure；secrets-gitleaks、sca-govulncheck、unit-test 都是 success |
| sast-gosec 步驟（公開 API） | setup-go、Install gosec、**Verify gosec module hash** 都是 success；只有 **Run gosec** failure |
| log：版本與完整性（20 號檔） | 第 140 行 `go version go1.27.2 linux/amd64`；第 256 行 `mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`；第 257 行 `gosec h1 OK` |
| log：偵測（20 號檔） | 第 375 行 `internal/k4canary/canary.go:12 - G404 (CWE-338) … (Confidence: MEDIUM, Severity: HIGH)`；第 387 行 `Issues : 1`；第 389 行 `##[error]Process completed with exit code 1.` |
| log：checkout | 第 93～107 行取出的是 `e9db9e49…`，對應 `origin/canary/k4-gosec-v2300-red` |
| log 時間戳與 API | 各步驟時間一致：開始 08:58:54Z，Verify 08:59:48Z，Run gosec 09:00:07Z 結束 |
| 分支已刪 | 安全部在 /tmp/v23 執行 `git ls-remote origin 'refs/heads/canary/*'`，輸出 0 行；branches API 回 404。維運部回報的刪除時間是 17:00:31（UTC+8），這個時間是轉述 |
| 沒有合併 | 沒有開 PR；main 仍是 `db380a40` |

## 4. 證據（維運部存在本機 `docs/09-cicd-env/evidence/cr-005/`，還沒進 repo）
| 檔案 | sha256 |
|---|---|
| 19-k4-run-40477332321.json | `43bf70797031816cc54c9c1cb3162a32851af215abdf75cd16d3649caa6572c2` |
| 20-k4-run-40477332321-sast-gosec.log（405 行） | `b780d4a8d6d53c288f2bed5c321633099deb56deeb70e94d88b0b5f61e673004` |
| 21-k4-branch-deleted-ls-remote.txt（空檔） | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| SHA256SUMS | `d204af5aa356ade1efc4211091327a2aa13882f1170422e943cae66889327d4e`（`sha256sum -c` 全部 OK） |

job log 要有 admin 權限才能下載，所以這份 log 是維運部下載的。安全部已核對 log 內容、時間戳都和公開 API 一致。這些證據之後會由維運部用另一個文件 PR 送進 repo，推送前須經使用者核准。

## 5. 連帶處置
- 05 號 v0.1.13：RR-010 結案。
- 03 號 v0.3.9：§2 SAST 改為發行版 v2.30.0，暫行例外結束。
- 13 號 v0.1.3：C7 結束；改為每月檢視工具版本，依 impact/安全部.md v0.2.1 的 S6 裁定。每月檢視至少涵蓋：gosec 的發行版與安全公告、Go 的修補版（含 RR-012）、govulncheck 與 gitleaks 的新版本。
- 如果 CR-005 回滾，會回到 pseudo-version，RR-010 重新開啟，C7 也恢復每週檢查。

## 修訂紀錄
| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 安全部 | 初版：K4 關閉、RR-010 結案 |
