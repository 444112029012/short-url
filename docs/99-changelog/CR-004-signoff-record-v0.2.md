# CR-004 簽核彙整（v0.2）

- 檔案版本：v0.2（2026-10-22，依審查部 CR-004 複審 #01 必改 CR4-RR01～RR05 改版）

- 對象：docs/99-changelog/CR-004-ruleset-stale-approval-v0.3.md，sha256 eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4（定稿，本文凍結）
- 原則：簽核與勘誤一律記在本檔，不回寫 CR 本文（避免簽核作廢，比照 CR-003）。
- 彙整人：Grok Bot（總協調），2026-10-22

## 1. 原件

| 部門 | 性質 | 原件 | sha256 | 綁定 | 結論 |
|---|---|---|---|---|---|
| 企劃部 | 影響列 | docs/01-requirements/06-cr-004-scope-schedule-impact-v0.1.md | eb6934d4049f40a335d8d189af3b54f6a6b75de15a8d59f2b05dd2e6f3cc26b2 | eae83378 | 範圍不適用，里程碑不變 |
| 設計部 | 影響列 | docs/02-architecture/12-impact-CR-004-design-v1.0.md | 746866809c0ae638f5824cb8ff2110d77c6b41d3d0fe69ba81b4940eb608b42b | eae83378 | 不適用 |
| 研發部 | 影響列 | docs/08-engineering/06-cr004-impact-rd-v0.1.md | cc72797909f13a304d41808a4dc0131d585711a5defbbf8b6430d823c8195b41 | eae83378 | 不改程式，工作量 0 |
| 安全部 | 意見書（R1） | docs/06-security/26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md（v0.1.1） | 8f4f3c62684b6683e8d2b54c449b87a6e44e1d26195829b6cc36e99dcfe39a8b | eae83378 | 附條件同意；RR-011 維持中、SCR-042 不關閉 |
| 品保部 | 影響列（R3） | docs/07-testing/15-CR-004-R3-qa-impact-row-v0.1.md | 4c3684fbfb92ed4463a24cdd4a6d0e14349ce4418bfb7315b440a49f1617ade2 | eae83378 | 無 REQ／NFR／SEC／TC 變更，產品回歸範圍不變；屬流程控制驗證，依 16 號 AC 驗收；品保估 0.5 人日 |
| 品保部 | 書面意見（R4） | docs/07-testing/16-CR-004-R4-qa-opinion-v0.1.md | ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c | eae83378 | 有條件同意，AC-CR004-1～9 |
| 審查部 | 影響列 | docs/gates/impact-CR-004-review-dept-v0.1.md | d53d581db6aafe673ed9db6ac6505452fbb11f95b6b37e132ad4861671623965 | eae83378 | 不適用 |

| 審查部 | 複審結論 | docs/gates/review-CR-004-2026-10-22.md（複審 #01） | （審查部自存） | eae83378 | 通過，CR 本文不改，R1～R7 關閉；惟須彙整檔改版並經複核 |

## 2. 本文外的決定與勘誤（以本節為準）

1. 安全部 26 號 sha：CR §0／§9 引用的是 v0.1 舊 sha 869bc215…7d53；現行原件為 v0.1.1 8f4f3c62…9a8b（上表），以上表為準。出處：品保 16 號 §6-2。
2. 執行判定：CR §4 步驟 3 的預期輸出字串**不作為判定依據**（gh 實際輸出鍵序為字母序）。判定一律以品保 16 號 AC-CR004-2～4 為準：24629854 的 dismiss_stale_reviews_on_push 與 require_last_push_approval 為 true；其餘欄位以 `jq -S` 比對不變（updated_at 除外）；24629805 的 .rules 不變。執行紀錄必須引用 16 號（ad44316f…495c）。出處：品保 16 號 §1 條件 (4)、§3、§6-1。待品保確認本項寫法。
3. 執行時機（R5）：使用者 2026-10-22 10:01（UTC+8）在維運部聊天回答 widget t67s322，原話「Switch it on after #23 merges and test it on #24」，對應 CR §3.1 方案 A：#23 合併後啟用，由 #24 做首次實測（程序見品保 16 號 §4）。出處：維運部聊天（審查部複審 #01 已核對）。
4. ENG-022 編號衝突（CR §11 S5）：**待確認，不在 CR-004 範圍**。目前的處理方向（#25 改掛既有 ENG-020、分支名保留記為偏差）只見於部門訊息，尚無書面原件；另案由企劃部出原件後處理，於 #25 同步複審時核對。

## 3. 核准條件（執行時必須全部滿足）

- 安全部 26 號 §4 條件 1～5（與 CR v0.3 §8 第 190～194 行逐字一致，安全部已確認）；RR-011 維持中、SCR-042 不關閉，CR 不宣稱修補 RR-011。
- 品保 16 號 §1 條件 (1)～(4)：依 AC-CR004-1～9 執行並留存證據；#24 實測紀錄明載「未實測撤銷行為」「未實測最後一次 push 須他人核准」；D-09／SCR-042 人工核 head 照舊，不得放寬；比對改用 `jq -S`（即本檔 §2-2）。

## 4. 執行前提

1. CR-004 經使用者簽核（下表）。
2. PR #23 已合併。
3. PR #24 尚未同步 main（AC-CR004-1）。
4. 由使用者本人修改 ruleset 24629854；審查部通過不等於可以執行。

## 5. 使用者簽核

| 項目 | 結論 | 時間（UTC+8） | 出處 |
|---|---|---|---|
| CR-004 v0.3 核准 | 核准 | 2026-10-22 10:09 | 使用者在 Grok Bot 聊天回答 widget t242s965：「核准 CR-004」 |
| 實際執行 ruleset 變更 | 待執行 | | |
