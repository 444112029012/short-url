---
文件：CR-003（分支命名：純文件 PR 例外）安全意見原件——回應審查部 CR3-R11
版本：v0.1
狀態：**同意，附三項前提**；僅適用於 CR-003 sha256 `c55ad3ca7c97497e07c84cfce2503bc74063bdff9de6a7692b3a6604d44f57b2` 版
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22
對應：docs/99-changelog/CR-003-branch-docs-exception-v0.1.md（§2.2、§3 資安列、AC-022-6、§7）；[22](./22-cr3-r8-stale-approval-residual-risk-v0.1.md)（CR3-R8，評等中，sha256 ca0d6106e697bd4e1415be8c464016baecd3875ca65147c6e038409c38f3f316）；SCR-042；D-09
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

# CR-003 安全意見

## 1. 核對對象
安全部於 2026-10-22 09:21 在本機計算 `docs/99-changelog/CR-003-branch-docs-exception-v0.1.md` 的 sha256，得到 `c55ad3ca7c97497e07c84cfce2503bc74063bdff9de6a7692b3a6604d44f57b2`，與總協調指定版本相同。已逐條核對該版 §3 資安列、§2.2 第 51 行、AC-022-6 與 §7 安全部列，內容與以下三項前提一致。

## 2. 結論：同意，附三項前提
資安面無影響：不新增、不調整 SEC，威脅模型不變。前提如下：
1. **例外只涵蓋分支命名。** 純文件 PR 仍須跑四項 required checks（gitleaks、gosec、govulncheck、unit-test），審查與核准規則不變（strict、至少 1 人核准）。核准與 head 的綁定只靠合併前人工核對 head（D-09、SCR-042），因為 Ruleset 24629854 的 `dismiss_stale_reviews_on_push=false`、`require_last_push_approval=false`。這項殘餘風險評等為中，見 22 號檔，預定登錄為 RR-011。不得以此為由在 workflow 加 paths-ignore。
2. **「只改 docs/」以 PR 實際變更檔判定。** .github/、go.mod／go.sum、任何 *.go 檔與根目錄設定檔都算 docs/ 以外。審查時須附 `git diff --name-only` 的輸出作為證據。
3. **docs/ 裡會被程式、測試或 CI 讀取的檔案照常走資安審查。** 目前是 docs/04-api/openapi.yaml，以及契約測試讀取的檔案；不論 PR 是否使用例外分支名都一樣。研發部已決定把這項寫進 02 規範 v0.2 §2.2（CR-003 第 51 行）。

另外，研發部的附帶條件是範圍擴大到 docs/ 以外時須經安全部評估，安全部同意。

## 3. 效力
- 本意見只適用於上述 sha256 的 CR-003 版本。CR-003 若再修改，且改到 §2.2、§3 資安列、AC-022-6 或 §7 安全部列，須重新取得安全部確認。
- 22 號檔建議的修補（開啟 require_last_push_approval）須依變更管理處理，並先取得使用者同意；不在本 CR 範圍內。

## 4. 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 初版，回應 CR3-R11，綁定 CR-003 c55ad3ca… |
