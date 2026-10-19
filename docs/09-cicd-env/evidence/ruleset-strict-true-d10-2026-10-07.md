---
文件編號：EVD-D10-STRICT-20261007
文件：規則一 main-required-checks 啟用 strict（D-10）生效取證
版本：v0.1
狀態：完成
負責：維運部
日期：2026-10-07
專案代號：SHORTURL
對應：ENG-020；NFR-004（SI-05）；SEC-009、SEC-010、SEC-011；D-08、D-10、RR-006；安全部 `docs/06-security/03-ci-security-gates-v0.3.md` §7.1；SCR-011、SCR-012
---

# 規則一 strict=true 生效取證（D-10）

## 1. 範圍與時點
- 對象：repo `444112029012/short-url`，ruleset **24629805 `main-required-checks`**（target `main`）。
- 設定者：使用者本人於 GitHub 網頁 UI 勾選「Require branches to be up to date before merging」並儲存（約 14:33，UTC+8）。維運部**未修改** Rulesets，只做唯讀 API 讀取。
- 設定前狀態：14:32:14（UTC+8）唯讀讀取時 `strict_required_status_checks_policy=false`，`updated_at` 為 2026-10-07 13:19:23（UTC+8）。
- 時間一律換算為 UTC+8（台北）；括號內附 API 原始 UTC 值。

## 2. 讀取來源
| # | 方式 | 讀取時間（UTC+8） | 原始檔 |
|---|---|---|---|
| A | 未驗證身分的公開 REST API `GET /repos/444112029012/short-url/rulesets/24629805`（HTTP 200） | 14:34:10（06:34:10Z，取自回應 `date` 標頭） | [`rulesets-2026-10-07/ruleset-24629805-d10-public-read.json`](./rulesets-2026-10-07/ruleset-24629805-d10-public-read.json) |
| B | 雲端代理的 GitHub token，同一 API（唯讀；回應標頭 accepted permission 為 `metadata=read`，非 admin） | 14:36:20（06:36:20Z，取自回應 `Date` 標頭） | [`rulesets-2026-10-07/ruleset-24629805-d10-auth-read.json`](./rulesets-2026-10-07/ruleset-24629805-d10-auth-read.json) |
| C | 使用者於 GitHub 網頁 UI 儲存後的畫面（由總協調轉述）：顯示「Ruleset updated」，Bypass list 為空 | 14:33 | 使用者截圖（未收入 repo） |

> `bypass_actors` 只回給具 repo admin 權限的呼叫者。A（公開）與 B（`metadata=read`）都**沒有回傳**此欄位，代表權限不足、讀不到，**不代表清單為空，也不代表非空**；B 另回 `current_user_can_bypass: "never"`（只說明該 token 本身不能繞過）。因此 bypass 清單以使用者 UI 畫面（C）為準；API 層的 bypass 佐證留待有 admin 讀取權限時補取（見 §4）。

## 3. 檢核結果
| # | 檢核項目 | 預期 | 實際（來源） | 結果 |
|---|---|---|---|---|
| 1 | `enforcement` | `active` | `active`（A） | PASS |
| 2 | `rules[required_status_checks].parameters.strict_required_status_checks_policy` | `true` | `true`（A） | PASS |
| 3 | required checks 名稱（四項，不改名） | `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` | 四項名稱完全一致，皆 `integration_id: 15368`（GitHub Actions）（A） | PASS |
| 4 | 其他規則未變 | `deletion`、`non_fast_forward` 仍在 | 兩項仍在（A） | PASS |
| 5 | `do_not_enforce_on_create` | `false` | `false`（A） | PASS |
| 6 | Bypass list | 空 | UI 顯示為空（C）；API 讀不到此欄位（A、B 皆未回傳） | PASS（依 UI 畫面；API 未能佐證） |
| 7 | `updated_at` 落在使用者儲存時點 | 約 14:33（UTC+8） | 14:33:28（2026-10-07T06:33:28.486Z）（A、B 一致） | PASS |

## 4. 影響（維運說明）
- strict 生效後，PR 的四項 required checks 必須是在「已含最新 main」的 head 上跑出的綠燈才可合併。任何其他 PR 先合併後，其餘 PR 須先 update branch、重取四綠，再請審查部依新 head 確認（D-10；合併順序見 `decision-log` D-10）。
- 本次只改規則一；規則二 24629854 `main-pr-review` 未在本次範圍內讀取與變更。
- 待補：以具 admin 讀取權限的方式（例如使用者本人在已登入的 GitHub 網頁開啟 ruleset JSON，或 admin token 唯讀呼叫）取得 `bypass_actors: []` 的 API 佐證，補入本檔 v0.2。

## 5. 修訂紀錄
| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | 初版：D-10 strict=true 生效取證（公開讀取、metadata=read 讀取、使用者 UI 畫面）；bypass 以 UI 為準，API 佐證待補 |
