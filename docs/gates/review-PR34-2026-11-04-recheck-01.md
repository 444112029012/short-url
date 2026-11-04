# 審查紀錄：PR #34（docs/pr33-followups）複審 recheck-01

- 審查部（獨立關卡，只審不改）；2026-11-04 09:29～09:40（UTC+8）
- 對象：head `9ec7b01d5dd968e72edc45c73c52613fd01726b1`（parent 3bf1d02c091ec62eb315c83344d4d7b3d08d022e），base main `e65eeaa3a1b3d482af77fd73045fd2ae2c370f14`，run 41608874768
- 前輪：review-PR34-2026-11-03.md（efb3e95c…）
- 全程未改專案檔、未留言、未推送、未 rerun；sha256 自算；暫存 worktree 用畢移除

## 結論：**通過（可合併）**。必改：無。PR34-R1、R2 關閉。

## 1. 綁定值與 CI
| 項目 | 實測 |
|---|---|
| head | ls-remote refs/pull/34/head＝PR API＝9ec7b01d…26b1；commit 2026-11-03 17:08:28（UTC+8），parent 3bf1d02c；OPEN、0 留言、0 review |
| base | ls-remote main＝e65eeaa3a1b3d482af77fd73045fd2ae2c370f14（暫停期間未變）；e65eeaa3..head＝2 commits |
| run | 41608874768（pull_request，head_sha＝9ec7b01d，success，attempt 1，建立 2026-11-03 17:10:34 UTC+8）；secrets-gitleaks 125825449531、sca-govulncheck 125825449031、unit-test 125825449564、sast-gosec 125825449508 全 success（17:12:36 完成） |
| 推送核准 | Grok 67cef426 widget **t405s1543**（2026-11-03 17:10，題目「要推送 9ec7b01d 到 PR #34 嗎？」）：「同意推送 9ec7b01d 到 PR #34」。run 建立於 17:10:34，在核准之後（核准時間只精確到分鐘，同一分鐘內，先後與 Grok 紀錄一致） |
| diff 3bf1d02c..head | 4 檔，全在 docs/：cr-005/README-SHA256SUMS.md、decision-log.md、gates/CR-006-review-records.SHA256SUMS、gates/review-PR34-2026-11-03.md（新增） |
| 凍結檔 | 35 號 1000dd9f、34 號 ed210033、CR 本文 818f48eb、36 號快照 2936baec 未變 |
| SHA256SUMS | 全 repo 以 `--ignore-missing` 全 OK（一般 `-c` 只缺 M5 01／04 與 cr-005 的 3 檔，皆為既有）；cr-005 SHA256SUMS.repo 32/32 OK；gates SUMS 8/8 OK，新增的 review-PR34-2026-11-03.md＝efb3e95c…，與 /workspace 原檔相同 |

## 2. 必改項複審
- **PR34-R1 關閉**：decision-log 在 D-24 之前**追加**一列「D-23 更正」，改指 `docs/gates/review-PR33-2026-11-03-recheck-01.md` §7（該檔在 repo 內）。D-23 原列逐字未動，只追加。
- **PR34-R2 關閉**：README-SHA256SUMS.md 新增 §1.1 依據。33 號 0b7ef6cf（第 9、17、20 行）、34 號 ed210033（第 25、28 行）的摘錄都與原文相符，寫明 gitleaks generic-api-key 誤判公開 h1、3 檔只留本機、本機用 SHA256SUMS 驗。狀態改為「定稿（已核對）」，版本 v1.0，有修訂列。
- **PR34-S1**：D-24 已補註 t120u（2026-11-03 16:31），但註記寫「此則為審查部轉述，審查部無法以 ReadTranscript 讀到原文」。依使用者本輪說明，t120u 是使用者在審查部聊天親口說的，審查部已親自確認，所以這段註記已經過時 → S 項。另外，D-24 是在原列儲存格內追加補記，比照 D-20 先例，可以接受。

## 3. 使用者原話（ReadTranscript，Grok 67cef426）
| id | 時間（UTC+8） | 原文 |
|---|---|---|
| t405s1543 | 2026-11-03 17:10 | 「同意推送 9ec7b01d 到 PR #34」 |
| t404s1540 | 2026-11-03 17:08 | 「21 號跟 36、37 號一樣：舊版存快照，同一檔名放最新版，之後的文件也照這個通則」 |
| t120u（審查部聊天） | 2026-11-03 16:31 | 「合併了」（使用者本輪告知審查部已親自確認，本輪不再查） |

## 4. S 項（不阻擋）
| 編號 | 內容 | 負責 |
|---|---|---|
| PR34-S3 | D-24 的 t120u 註記仍寫「審查部轉述、讀不到原文」，與審查部已親自確認不符；下個文件 PR 追加更正 | 總協調 |
| PR34-S4 | 使用者 t404s1540 把快照做法延伸為所有文件的通則（解了 PR34-S2），但 decision-log 還沒記 | 總協調 |
| 沿用 | CR6-S4、S7、S9；S6（G5 前）；SCR-053（11-30）；PR33-S8 | 見結案紀錄 §7 |

## 5. 與 Grok 宣稱不符
1. Grok 說 S1 已補 t120u，屬實；但註記內容是「審查部轉述、讀不到」，與目前事實（審查部已確認）不符（S3）。
2. Grok 說四綠，屬實；run 在 2026-11-03 17:10:34 建立，不是今天。
其餘都經實測相符。

— 審查部　2026-11-04（UTC+8）
