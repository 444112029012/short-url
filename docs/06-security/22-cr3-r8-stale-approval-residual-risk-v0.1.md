---
文件：CR3-R8 殘餘風險評等：ruleset 不強制「核准綁定 SHA」
版本：v0.1
狀態：評等完成：**中**；建議修補（開啟兩個旗標），修補前以人工核 head 補償
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22
對應：審查部 CR3-R8；SCR-042（[17](./17-pr-017-dr-g4-01-body-limit-security-confirmation-v0.1.md)）；D-09；CR-003；殘餘風險登錄 [05](./05-residual-risk-register-v0.1.md)（預定新增 RR-011）
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# CR3-R8 殘餘風險評等

## 1. 風險描述
main-pr-review ruleset（id 24629854）中，`dismiss_stale_reviews_on_push=false` 和 `require_last_push_approval=false`。因此 PR 核准之後再 push 新 commit，原核准仍然有效，GitHub 不會強制「合併的 head＝核准的 head」。這項保證目前只靠合併前人工核對 head（D-09、SCR-042）。
（設定值引自審查部 CR3-R8 與 SCR-042 既有紀錄；本次因 box 上 gh 未登入，安全部未重新讀取 ruleset。）

## 2. 評等（STRIDE：竄改 T、否認 R）
| 項目 | 評估 |
|---|---|
| 可能性 | **中**：多個代理共用同一帳號推送，strict=true 本來就常需要同步 push；人工核 head 只是流程，沒有工具強制，漏核一次就成立 |
| 影響 | **中**：未經審查的變更可能進入 main，繞過安全部與審查部的關卡結論。但四項必要檢查（gitleaks、gosec、govulncheck、unit-test）仍會對最終 head 執行，所以已知類型的機密外洩、漏洞依賴與靜態問題仍會被擋下；剩下的缺口是邏輯與授權類變更 |
| 補償控制 | 合併前人工核 head（D-09），以及安全部每次的「只同步」範圍核對。截至 #20 已全數執行，未發現不符 |
| **殘餘評等** | **中**（未達高，依現行規則不必交使用者拍板，但不建議長期接受） |

## 3. 建議
1. **修補（建議採行，成本低）**：維運部把 ruleset 24629854 的 `require_last_push_approval` 設為 true，並評估同時設 `dismiss_stale_reviews_on_push=true`。代價：每次同步 main 後要重新核准一次。屬於環境與規則變更，須走 change-management，並取得使用者同意後才能變更。
2. **修補前**：維持 D-09 人工核 head。合併紀錄必須寫明「合併的 head＝核准 SHA」，以及 `git diff <核准 SHA> <merge>` 的檢查結果。
3. 修補完成後，複評可降為低，SCR-042 關閉。

## 4. 結論
殘餘風險評等為**中**，登錄為 RR-011（05 號檔下一版寫入）。若不修補而要長期接受，須由使用者書面接受。

## 5. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版，回應 CR3-R8 |
