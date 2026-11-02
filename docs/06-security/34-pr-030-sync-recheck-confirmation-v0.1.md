# PR #30 同步範圍複查確認書

| 欄位 | 內容 |
|---|---|
| 文件編號 | 06-security/34 |
| 版本 | v0.1 |
| 負責 | 安全部 |
| 綁定 head | 51dc255c4911a5dc3f809376e7cbc97c190b5705 |
| base | db380a40c3330bf23498ae7fda2f123cf7dd4cc4 |
| 複查時間 | 2026-10-28 09:15:27 (UTC+8，取自 `date`) |
| 依據 | 06-security/33 §合併條件 |

## 1. 結論
**通過。** 本確認書只對 head 51dc255c 有效，head 一有變動就失效（D-09）。

推送核准來源：Grok Bot 轉述使用者在 10/28 09:13 回覆 widget t328s1268，原文為「同意強制推送 docs/cr-005-followups，以 51dc255c 取代 be3d1bde」。這句話是由 agent 轉述，本文件只記錄為轉述，不視為安全部親自收到的核准。審查部可以用 ReadTranscript 查證原文。

## 2. 查核項目（公開 clone /tmp/v23 與未登入 API）
| # | 項目 | 結果 |
|---|---|---|
| 1 | PR head / base（pulls API） | 51dc255c / db380a40，state open，mergeable_state clean |
| 2 | commit 的 parent | db380a40（只有一個 parent） |
| 3 | diff 範圍 | 共 82 個檔案，全部在 docs/ 底下（6 個中文檔名被 git 加了引號，路徑仍是 docs/99-changelog/CR-005/impact/） |
| 4 | CI 設定與例外清單 | 沒有改動 .github/ 或 .gitleaksignore |
| 5 | 33 號裁定要排除的 3 個檔案 | 03-version-m-gosec-pseudo707f227.txt、03-version-m-gosec-v2300.txt、15-diff-pseudo-vs-v2300.txt 都不在此 head |
| 6 | 06-security 各檔 sha256 | 31 c1c1f9ab…、32 74bac45e…、33 0b7ef6cf…、03 v0.3 624f08d7…、05 df8636d8…、13 32f7f8c2…、README 84fd26aa…，全部與安全部的定稿相符 |
| 7 | impact/安全部.md | 5ea337a3…，與 v0.2.1 相符 |
| 8 | evidence/cr-005 的雜湊清單 | SHA256SUMS.repo 為 932da82b…，共 32 列，`sha256sum -c` 全部 OK。原本的 SHA256SUMS（d204af5a…）仍留在 repo，內容沒有改動。目錄裡有 34 個檔案，也就是 32 個證據檔加上這 2 份清單 |
| 9 | CI run 40613719732 | 事件為 pull_request，head_sha 是 51dc255c，結論 success。secrets-gitleaks、sca-govulncheck、unit-test、sast-gosec 四項都是 success，最後一項在 09:14:54（UTC+8）完成 |

## 3. 說明
- 雜湊清單的檔案歸維運部，安全部只做驗證，沒有修改。
- 本確認書與 README v0.25 不在 #30 裡，等下次文件同步時再一起收進 repo，避免改動這次已確認的 head。
- 是否合併由使用者決定。

## 修訂紀錄
| 版本 | 日期 | 負責 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-28 | 安全部 | 初版；小版修訂不改檔名（檔名保留主版號） |
