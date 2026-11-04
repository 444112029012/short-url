# cr-005 證據目錄：SHA256SUMS 與 SHA256SUMS.repo 說明

| 項目 | 內容 |
|---|---|
| 撰寫部門 | 維運部 |
| 日期 | 2026-11-03（UTC+8） |
| 基準 commit | e65eeaa3a1b3d482af77fd73045fd2ae2c370f14（main） |
| 對應審查項 | PR33-S9、PR34-R2 |
| 版本 | v1.0 |
| 狀態 | 定稿（已核對） |

## 1. 狀況
本目錄的 `SHA256SUMS` 共 35 筆，是維運部當初在本機產生證據時的原始紀錄。其中 3 個檔案沒有入庫，在 repo 裡找不到：

| 檔名 | SHA256SUMS 所列 sha256 |
|---|---|
| 03-version-m-gosec-pseudo707f227.txt | 197f2ec8832b8a6723ab4e2b06002fbaec6749947224a689e7cd55e9b2ef2b34 |
| 03-version-m-gosec-v2300.txt | a60760a4c2bf1ea8b05caeb886968b59c70bd0ad22651f9bdca16e1e93586126 |
| 15-diff-pseudo-vs-v2300.txt | f27236cece2657d23e627409c4e5890c215d3b239c8e2e9f8305faae78f56061 |

查證方式：在 e65eeaa3 執行 `git log --all -- docs/09-cicd-env/evidence/cr-005/<檔名>`，3 個檔案都沒有任何 commit；在任何路徑下也查不到同名檔案。也就是說，這 3 個檔案從來沒有提交進 repo。

### 1.1 依據
這 3 個檔案不入庫，是依安全部 33 號裁定：gitleaks generic-api-key 規則把檔內公開的模組 h1 雜湊誤判為機密（誤報），因此原檔只留在本機原樣保存。本機原檔可用 `SHA256SUMS` 驗證。

| 文件 | sha256 | 引用行 | 原文摘錄 |
|---|---|---|---|
| `docs/06-security/33-pr-030-gitleaks-false-positive-ruling-v0.1.md` | 0b7ef6cf1e899cd441d37d2f140787b467f63c68d54670dc46e1007dee566bef | 第 9、17、20 行 | 「## 1. 判定：誤報，不是機密」；「## 2. 裁定：只把命中的 3 個檔案留在本機，不放寬門檻，也不遮蔽證據」；「這 3 個檔案留在本機原樣保存，sha256 記在本檔 §3，維持可追溯。」 |
| `docs/06-security/34-pr-030-sync-recheck-confirmation-v0.1.md` | ed210033290835d2693e1a57c221514d1a8986e133796d1a29e68379a782e072 | 第 25、28 行 | 「33 號裁定要排除的 3 個檔案 \| 03-version-m-gosec-pseudo707f227.txt、03-version-m-gosec-v2300.txt、15-diff-pseudo-vs-v2300.txt 都不在此 head」；「SHA256SUMS.repo 為 932da82b…，共 32 列，`sha256sum -c` 全部 OK。」 |

上表 3 個檔案的 sha256 與 33 號 §3（第 26～28 行）逐字相同。

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

## 修訂紀錄
| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-11-03 | 維運部 | 初版（草稿，PR33-S9） |
| v1.0 | 2026-11-03 | 維運部 | PR34-R2：補 §1.1 依據（33 號、34 號），狀態改為定稿 |
