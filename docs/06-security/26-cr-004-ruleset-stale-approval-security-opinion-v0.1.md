---
文件：CR-004（ruleset 24629854 開啟 dismiss_stale_reviews_on_push、require_last_push_approval）安全評估
版本：v0.1.1
狀態：**v0.1.1：對 CR-004 v0.3 再確認，附條件同意**（見 §6）；v0.1：**附條件同意**；RR-011 修補後**不降級**，仍維持中，理由見 §3
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22（v0.1.1）
對應：CR-004 v0.3 /workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.3.md（sha256 eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4）；v0.1 審查對象為維運部 CR-004 草稿 /workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.1.md（sha256 7a5aaaaba96b56d58057e946628274d003648892da442b9b558ff854004d7bea）；[22](./22-cr3-r8-stale-approval-residual-risk-v0.1.md)（RR-011，中）；SCR-042；D-09
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# CR-004 安全評估

## 1. 核對結果（2026-10-22 09:54，本機檔案）
| 檔案 | sha256 | 結果 |
|---|---|---|
| ruleset-24629854-before-2026-10-22.json | af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12 | 與維運部提供的值相符；兩個旗標都是 false |
| change-payload.json | c6b261ac4aee12e4c1f99f1b65d3ab27f7daaff12931743b8e0854a60676e02e | 與 rollback 相比只差兩個旗標（false→true），其餘欄位與快照一致，bypass_actors 沒有被改動 |
| rollback-payload.json | dfcae65fd02fd144ca84b84ff01c4e453f57c46461270aa5722bb527a1aa8683 | 等同快照的可寫欄位 |

## 2. 兩個指定問題
1. **使用者用 admin bypass 合併仍可行嗎？** 依快照，`bypass_actors` 為 RepositoryRole 5（Admin）、`bypass_mode: pull_request`，`current_user_can_bypass: pull_requests_only`。兩個新旗標屬於同一條 pull_request 規則，bypass 時會整條略過，所以使用者應仍可透過 PR 以 bypass 合併。這是依據 GitHub 規則集的行為推定，尚未實測；請在變更後的第一個 PR 實際確認。
2. **人工核 head（D-09、SCR-042）是否保留？** **必須保留。** 理由見 §3。

## 3. 關鍵發現：改完後的保護範圍有限
- 維運部說 #23～#26 目前都沒有任何核准。從 #17 起，各 PR 的實際合併路徑很可能都是 admin bypass，因為單一使用者帳號無法核准自己開的 PR；但這點尚未逐一查證。
- 用 bypass 合併時，「至少 1 人核准」、`dismiss_stale_reviews_on_push`、`require_last_push_approval` 都不會生效。因此 CR-004 只在「確實有第二人核准、而且不使用 bypass」時，才能強制核准綁定最終 push。照目前的作業方式，RR-011 的缺口基本上仍然存在。
- 結論：CR-004 本身風險低、可回滾，同意實施；但 **RR-011 修補後仍評為中，SCR-042 不關閉**。22 號檔 §3 第 3 點「修補後可降為低」以本檔為準，在此更正。

## 4. 同意條件
1. 由使用者本人執行變更；變更後以 GET 取得 ruleset，留存快照與 sha256，並確認只有兩個旗標改變、bypass_actors 不變。
2. 變更後第一個 PR 實測：bypass 合併可行；同步 push 後，既有核准會被撤銷（如果當時有核准的話）。
3. D-09 人工核 head 照舊。每次用 bypass 合併時，合併紀錄要寫明 bypass、被合併的 head＝安全部確認的 SHA，並附 `git diff <確認 SHA> <merge>` 的結果。
4. 不得在 CR-004 裡順帶改動 bypass_actors、required checks（ruleset 24629805）或任何其他欄位。
5. 若要真正降低 RR-011，可另案評估以下兩種做法，並由使用者拍板：(a) 設第二個人工核准帳號、移除或收窄 bypass；(b) 以 CI 檢查合併前的 head 是否等於安全部確認的 SHA。

## 5. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版，評估 CR-004 草稿 7a5aaaab… |

## 6. v0.1.1 再確認：CR-004 v0.3（回應審查部 CR4-R1）
2026-10-22 10:00 在本機重算 sha256，結果如下：
| 檔案 | sha256 | 結果 |
|---|---|---|
| CR-004-ruleset-stale-approval-v0.3.md | eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4 | 與總協調提供的值相符 |
| change-payload.json | c6b261ac4aee12e4c1f99f1b65d3ab27f7daaff12931743b8e0854a60676e02e | 與 v0.1 審查時相同 |
| rollback-payload.json | dfcae65fd02fd144ca84b84ff01c4e453f57c46461270aa5722bb527a1aa8683 | 與 v0.1 審查時相同 |
| ruleset-24629854-before-2026-10-22.json | af1c9e4a958468c332275f38fdd62da611ce616d3106e42c637f34588646ee12 | 與 v0.1 審查時相同 |

- **五項條件：** v0.3 §8 第 190～194 行與本檔 §4 條件 1～5 逐字比對一致，沒有改寫或重新編號。各條件的落實位置見 v0.3 第 197 行，安全部同意。
- **RR-011 說法：** v0.3 明寫「本 CR 不修補 RR-011」，RR-011 維持中、SCR-042 不關閉、D-09 保留，符合本檔 §3。
- **補充，不影響結論：**
  1. v0.3 §3.1 說明 #24 上無法實測「核准被撤銷」，並要求在紀錄中明載「未實測撤銷行為」，安全部接受。條件 2 的這半句在單帳號下不適用，須等到有第二位審查者時再實測。
  2. v0.3 §4 步驟 2 的 PUT 警語（以當下最新 GET 產生完整 payload；GET 與快照不一致就停止）安全部同意，並列為執行時的必要檢查。
  3. v0.3 第 45 行提到「D-12 允許以使用者帳號留審查結論」。這是使用者的決策，不在本 CR 範圍；安全部本身仍不在 PR 上留言或核准。
- **結論：** 對 CR-004 v0.3（sha256 eae83378…7ed4）及上列三個附件**附條件同意**，條件即 §4 條件 1～5。CR-004 若再修改，或附件 sha256 有變，本確認即失效。

## 7. 修訂紀錄（續）
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1.1 | 2026-10-22 | 安全部 | 對 CR-004 v0.3 再確認並綁定 sha256；小版修訂不改檔名（檔名保留主版號） |
