# PR #30 secrets-gitleaks 命中裁定（安全部）

- 文件編號：06-security/33
- 版本：v0.1（Draft）
- 負責：安全部
- 關聯：PR [#30](https://github.com/444112029012/short-url/pull/30)（head `be3d1bde`，base `db380a40`）；run 40478991708；03 號 v0.3.9 §4、§12
- 撰寫時間：2026-10-27 17:15:34（UTC+8，取自 `date`）

## 1. 判定：誤報，不是機密
| 命中 | 內容 | 查證 |
|---|---|---|
| evidence/cr-005/03-version-m-gosec-pseudo707f227.txt:48；15-diff-pseudo-vs-v2300.txt:15 | `dep google.golang.org/api v0.301.0 h1:rTvI8szs…（全值見本機原檔）` | 與 sum.golang.org/lookup/google.golang.org/api@v0.301.0 的 h1 逐字相同 |
| evidence/cr-005/03-version-m-gosec-v2300.txt:48；15-diff-pseudo-vs-v2300.txt:19 | `dep google.golang.org/api v0.304.0 h1:2R1dYSiI…（全值見本機原檔）` | 與 sum.golang.org/lookup/google.golang.org/api@v0.304.0 的 h1 逐字相同 |

這些值是公開透明日誌裡的模組內容雜湊，不是憑證。generic-api-key 規則會被觸發，是因為同一行有 `api` 關鍵字（模組路徑 `google.golang.org/api`），後面又接著一段高熵字串。其他 h1 行沒有 `api` 關鍵字，所以沒被命中。

## 2. 裁定：只把命中的 3 個檔案留在本機，不放寬門檻，也不遮蔽證據
- **不採 C（.gitleaksignore）**：這等於放寬 secrets 門檻，依 03 §12 必須走 CR，並經安全部核准。PR 掃描的 fingerprint 綁定 commit，rebase 之後就會失效，而且之後每次遇到同類輸出都要再加。
- **不採 B（遮蔽 h1）**：遮蔽後就不是原始證據，跟審查部「收原始書面紀錄」的要求衝突；而且那些 h1 正是用來證明依賴版本的內容。
- **採 A 的縮小版**：只把命中的 3 個檔案從 #30 拿掉，其餘 00～21 照收。這 3 個檔案留在本機原樣保存，sha256 記在本檔 §3，維持可追溯。repo 內要收一份新的 SHA256SUMS，只涵蓋實際收進 repo 的檔案；這份清單由維運部重出（是維運部的檔案，安全部不改）。本機那份涵蓋 00～21 的舊版 SHA256SUMS 保留不動。
- 以後的證據如果含 `go version -m` 原文，同樣先在本機跑 gitleaks 再收。

## 3. 本機保存、不進 repo 的檔案
| 檔案 | sha256 |
|---|---|
| 03-version-m-gosec-pseudo707f227.txt | `197f2ec8832b8a6723ab4e2b06002fbaec6749947224a689e7cd55e9b2ef2b34` |
| 03-version-m-gosec-v2300.txt | `a60760a4c2bf1ea8b05caeb886968b59c70bd0ad22651f9bdca16e1e93586126` |
| 15-diff-pseudo-vs-v2300.txt | `f27236cece2657d23e627409c4e5890c215d3b239c8e2e9f8305faae78f56061` |

這 3 個檔案只用來佐證「依賴只升版、沒有新增」（K5）。gosec 本身的版本和 h1，在 18、20 號 CI log 以及 06 號重現檔裡都有，所以這 3 個檔案不進 repo，不影響 K1～K4 的結論。

## 4. 合併前條件
- 新 head 上四項 check 都要通過，secrets-gitleaks 要轉綠，不能動 CI 設定或 .gitleaksignore。
- diff 只能在 docs/。
- 安全部會對新 head 做一次只看同步範圍的複查。

> 註：本檔的 h1 只寫前 8 碼，避免本檔自己被 generic-api-key 命中。

## 修訂紀錄
| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 安全部 | 初版 |
