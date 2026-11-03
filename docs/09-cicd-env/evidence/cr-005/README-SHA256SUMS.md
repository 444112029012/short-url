# cr-005 證據目錄：SHA256SUMS 與 SHA256SUMS.repo 說明

| 項目 | 內容 |
|---|---|
| 撰寫部門 | 維運部 |
| 日期 | 2026-11-03（UTC+8） |
| 基準 commit | e65eeaa3a1b3d482af77fd73045fd2ae2c370f14（main） |
| 對應審查項 | PR33-S9 |
| 狀態 | 草稿 |

## 1. 狀況
本目錄的 `SHA256SUMS` 共 35 筆，是維運部當初在本機產生證據時的原始紀錄。其中 3 個檔案沒有入庫，在 repo 裡找不到：

| 檔名 | SHA256SUMS 所列 sha256 |
|---|---|
| 03-version-m-gosec-pseudo707f227.txt | 197f2ec8832b8a6723ab4e2b06002fbaec6749947224a689e7cd55e9b2ef2b34 |
| 03-version-m-gosec-v2300.txt | a60760a4c2bf1ea8b05caeb886968b59c70bd0ad22651f9bdca16e1e93586126 |
| 15-diff-pseudo-vs-v2300.txt | f27236cece2657d23e627409c4e5890c215d3b239c8e2e9f8305faae78f56061 |

查證方式：在 e65eeaa3 執行 `git log --all -- docs/09-cicd-env/evidence/cr-005/<檔名>`，3 個檔案都沒有任何 commit；在任何路徑下也查不到同名檔案。也就是說，這 3 個檔案從來沒有提交進 repo。

## 2. 處理原則
- `SHA256SUMS` 保留原狀，不修改，作為原始本機紀錄。
- `SHA256SUMS.repo`（32 筆）只列出實際入庫的檔案，內容是 `SHA256SUMS` 扣掉上面 3 筆，其餘雜湊值都相同。
- 驗證 repo 內的副本時，一律使用 `SHA256SUMS.repo`。

## 3. 驗證指令與預期結果
```
cd docs/09-cicd-env/evidence/cr-005 && sha256sum -c SHA256SUMS.repo
```
預期結果：32 行都是 `<檔名>: OK`，沒有 FAILED，也沒有 missing，結束碼為 0。第一行是 `00-go1.27.2-download.txt: OK`，最後一行是 `21-k4-branch-deleted-ls-remote.txt: OK`。維運部 2026-11-03 已在 e65eeaa3 的副本上實測，結果和上面相同。

如果改用 `sha256sum -c SHA256SUMS`，上面 3 個檔案會出現 `No such file or directory`／`FAILED open or read`，這是預期中的結果，不代表證據遭到竄改。
