---
文件編號：EVD-G3-RULESETS-20261007
文件：Rulesets API 快照（D-08 路徑 A）
版本：v0.3
狀態：草稿待審查（合併後改定稿）
負責：維運部
日期：2026-10-07
專案代號：SHORTURL
對應：ENG-020；NFR-004（SI-05）；SEC-009、SEC-010、SEC-011；D-08、RR-006；`docs/09-cicd-env/06-rulesets-setup-guide-v0.1.md`
---

# Rulesets API 快照（2026-10-07）

## 目的

本文件是 **D-08 合併路徑 A** 的落地證據，用來補 G3 證據，並取代舊的 classic branch protection 證據。內容只記錄唯讀 API 回應與對照結果，沒有建立、修改或刪除任何 ruleset 或 branch protection。

## 呼叫清單

時區：`Date` 回應標頭為 GMT（下表 UTC），台北時間為 UTC+8。權杖未寫入本文件或原始 JSON。

| # | 方法與路徑 | HTTP | UTC | 台北時間 | Request ID |
|---|---|---|---|---|---|
| 1 | `GET /repos/444112029012/short-url/rulesets?includes_parents=true` | 200 | 2026-10-07 05:24:34 | 2026-10-07 13:24:34 | `0D18:387339:3F63C4:D00071:6AC5D792` |
| 2 | `GET /repos/444112029012/short-url/rulesets/24629805` | 200 | 2026-10-07 05:24:42 | 2026-10-07 13:24:42 | `DBC4:3A132F:4A15B8:F3A931:6AC5D79A` |
| 3 | `GET /repos/444112029012/short-url/rulesets/24629854` | 200 | 2026-10-07 05:24:42 | 2026-10-07 13:24:42 | `52B8:281AEA:46A157:E8DB42:6AC5D799` |
| 4 | `GET /repos/444112029012/short-url/rules/branches/main` | 200 | 2026-10-07 05:24:42 | 2026-10-07 13:24:42 | `2742:1D8314:3EE163:CEC2F1:6AC5D79A` |
| 5 | `GET /repos/444112029012/short-url/branches/main/protection` | 403 | 2026-10-07 05:24:42 | 2026-10-07 13:24:42 | `219B:199E9:430550:DBC47F:6AC5D79A` |
| 6 | `GET /repos/444112029012/short-url/branches/main` | 200 | 2026-10-07 05:24:43 | 2026-10-07 13:24:43 | `2C16:56F8E:3E73B8:CE8B67:6AC5D79A` |

呼叫 5 的回應本文為 `{"message":"Resource not accessible by integration",...}`。回應標頭 `X-Accepted-Github-Permissions: administration=read`。同一整合權杖對 ruleset 明細的接受權限為 `metadata=read`。

## 對照結果

| 檢查項 | 結果 | 依據 |
|---|---|---|
| 規則一 `main-required-checks`：`enforcement=active` | PASS | id `24629805` |
| 規則一：`target=branch` | PASS | `"target": "branch"` |
| 規則一：條件含 main（`refs/heads/main` 或 `~DEFAULT_BRANCH`） | PASS | `conditions.ref_name.include` 為 `["~DEFAULT_BRANCH"]` |
| 規則一：`bypass_actors` 為空 | PASS（人工確認：使用者截圖） | [Bypass list is empty](rulesets-2026-10-07/screenshots/rule1-main-required-checks-bypass-empty.png) |
| 規則一：`required_status_checks` 恰為四個 context | PASS | `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` |
| 規則一：`strict_required_status_checks_policy` | 已記錄 | `false` |
| 規則一：各 check 的 `integration_id` | 已記錄 | 四項皆為 `15368` |
| 規則一：含 `deletion` 與 `non_fast_forward` | PASS | rules 陣列含此二型別 |
| 規則二 `main-pr-review`：`enforcement=active` | PASS | id `24629854` |
| 規則二：目標為 main | PASS | `include` 為 `["~DEFAULT_BRANCH"]` |
| 規則二：`pull_request.required_approving_review_count=1` | PASS | 見下方參數 |
| 規則二：`bypass_actors` 恰為 Repository admin（`bypass_mode=pull_request`），且無其他項目 | PASS（人工確認：使用者截圖；ruleset 名稱未入鏡，依排除法歸屬，見下節） | [僅 Repository admin，Allow for pull requests only](rulesets-2026-10-07/screenshots/rule2-main-pr-review-bypass.png)；`actor_type=RepositoryRole`、`actor_id=5` 為 Repository admin 角色的 API 表示，屬推定（API 未回 `bypass_actors`，畫面只顯示 “Repository admin (Roles)”） |
| 沒有其他指向 main 的 active ruleset | PASS | `includes_parents=true` 的清單只有上述兩筆，皆為 active |
| 舊 ruleset 名稱 `main` 不存在 | PASS | 清單中無此名稱；舊 ruleset `GET /rulesets/24629586` 回 404（審查部約 13:39 重驗，見 `docs/gates/gate-G4-PR5-2026-10-07.md` §三 #14） |
| classic branch protection 已不存在 | PASS（人工確認：使用者截圖） | [Classic branch protections have not been configured](rulesets-2026-10-07/screenshots/branches-no-classic-protection.png) |

`current_user_can_bypass` 在兩筆明細都是 `"never"`。這是呼叫端整合身分自己能否繞過，不是 `bypass_actors` 清單。

v0.1 有三列 CANNOT VERIFY。原因是 API 權杖沒有 `administration=read`：ruleset 明細（HTTP 200，接受權限 `metadata=read`）省略 `bypass_actors`；`GET /branches/main/protection` 為 403。這三項改由儲存庫擁有者於 2026-10-07 提供的設定畫面截圖人工確認，維運部已檢視，結果改記為 PASS（人工確認：使用者截圖）。

## 規則二截圖的歸屬（排除法）

`rule2-main-pr-review-bypass.png` 上緣裁切，**ruleset 名稱不在畫面內**，畫面也沒有網址列，因此無法只靠這張圖看出是哪一條 ruleset。本截圖**不重拍**，改以下列排除法歸屬為規則二 `main-pr-review`（id `24629854`）：

1. **沒有其他 ruleset**：API 清單（呼叫 1，[`list.json`](rulesets-2026-10-07/list.json)）只回兩筆：`24629805`（`main-required-checks`）與 `24629854`（`main-pr-review`）。
2. **規則一另有截圖且名稱可見**：[`rule1-main-required-checks-bypass-empty.png`](rulesets-2026-10-07/screenshots/rule1-main-required-checks-bypass-empty.png) 標頭與 Ruleset Name 欄位皆顯示 `main-required-checks`，Bypass list 為空；而本圖的 Bypass list 有一列 Repository admin，故不可能是規則一。
3. **舊 ruleset `main` 已不存在**：清單無此名稱，`GET /rulesets/24629586` 回 404（審查部重驗）。
4. **ruleset id 與 API 一致**：`list.json` 中 id 與名稱一一對應（`24629805`↔`main-required-checks`、`24629854`↔`main-pr-review`）；`/rules/branches/main`（[`rules-branches-main.json`](rulesets-2026-10-07/rules-branches-main.json)）的 6 條生效規則只來自這兩個 id，其中 `pull_request` 規則屬 `24629854`。本圖可見的設定（Target「Default」、勾選 Restrict deletions、有 bypass）也與 `24629854` 明細相符。

剩下唯一可能即為 `main-pr-review`。

> **後續截圖慣例**：之後的證據截圖一律保留瀏覽器**網址列**（可看出 repo 與 ruleset id），並盡量讓 ruleset 名稱或頁面標頭入鏡。

## 結論

13 項 PASS＋2 項已記錄（符合指南），無需修正。「已記錄」兩項為規則一 `strict_required_status_checks_policy=false` 與各 check `integration_id=15368`，數值皆符合設定指南。

規則二同時含 `deletion` 與 `non_fast_forward`，與規則一重複，無害：規則一的 Bypass list 為空，四項 required check、禁刪除與禁 force push 沒有繞過名單。`require_extra_approval_for_unattributed_changes=true` 是 GitHub 目前的預設值。

## 本證據所在 PR 的 CI

CI 證據以**本證據所在 PR 頁面上最新一次** `ci-devsecops` run 為準（須對應 PR 目前 head，四個 required job `secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test` 皆 success）。本文件不再引用特定 run 編號，以免 head 變動後引用失效。

## 關鍵欄位（verbatim）

### `main-required-checks`（id 24629805）

```json
{
  "id": 24629805,
  "name": "main-required-checks",
  "enforcement": "active",
  "target": "branch",
  "conditions": {
    "ref_name": {
      "exclude": [],
      "include": ["~DEFAULT_BRANCH"]
    }
  },
  "rules": [
    {"type": "deletion"},
    {"type": "non_fast_forward"},
    {
      "type": "required_status_checks",
      "parameters": {
        "strict_required_status_checks_policy": false,
        "do_not_enforce_on_create": false,
        "required_status_checks": [
          {"context": "sca-govulncheck", "integration_id": 15368},
          {"context": "sast-gosec", "integration_id": 15368},
          {"context": "secrets-gitleaks", "integration_id": 15368},
          {"context": "unit-test", "integration_id": 15368}
        ]
      }
    }
  ]
}
```

回應中沒有 `bypass_actors`。

### `main-pr-review`（id 24629854）

```json
{
  "id": 24629854,
  "name": "main-pr-review",
  "enforcement": "active",
  "target": "branch",
  "conditions": {
    "ref_name": {
      "exclude": [],
      "include": ["~DEFAULT_BRANCH"]
    }
  },
  "rules": [
    {"type": "deletion"},
    {"type": "non_fast_forward"},
    {
      "type": "pull_request",
      "parameters": {
        "required_approving_review_count": 1,
        "dismiss_stale_reviews_on_push": false,
        "required_reviewers": [],
        "require_code_owner_review": false,
        "require_last_push_approval": false,
        "required_review_thread_resolution": false,
        "require_extra_approval_for_unattributed_changes": true,
        "allowed_merge_methods": ["merge", "squash", "rebase"]
      }
    }
  ]
}
```

回應中沒有 `bypass_actors`。

### `GET /branches/main` 與保護狀態相關欄位

```json
{
  "name": "main",
  "protected": true,
  "protection": {
    "enabled": false,
    "required_status_checks": {
      "enforcement_level": "off",
      "contexts": [],
      "checks": []
    }
  }
}
```

`protected: true` 只表示此分支受規則約束（ruleset 也會讓此欄位為 true）。`protection.enabled: false` 與「classic protection 未啟用」一致。專用端點仍是 403，classic protection 已刪除是由擁有者截圖確認，不是由 404 回應確認。

## 與設定指南的差異

對照 `06-rulesets-setup-guide-v0.1`（D-08 方案 A）：

- 指南要求規則一 Bypass list 空白、規則二僅 Repository admin 且模式為 For pull requests only。API 沒有回傳 `bypass_actors`；v0.2 改以擁有者截圖確認，與指南相符。
- 指南要求規則二一併勾選 Restrict deletions 與 Block force pushes。快照的 `deletion`、`non_fast_forward` 與此相符。這兩條與規則一重複，無害，因為規則一沒有 bypass。
- 指南要求不要勾「Require branches to be up to date before merging」與「Do not require status checks on creation」。快照為 `strict_required_status_checks_policy: false`、`do_not_enforce_on_create: false`；狀態檢查截圖亦顯示不要求 up to date。
- 四個 check 的 `integration_id` 皆為 `15368`（GitHub Actions）。指南要求來源選 GitHub Actions。名稱集合與指南一致；API 陣列順序是 `sca-govulncheck`、`sast-gosec`、`secrets-gitleaks`、`unit-test`。
- 規則二 `pull_request` 的 `require_extra_approval_for_unattributed_changes=true` 是 GitHub 預設。指南寫「其他子選項維持預設」。核准數為 1。
- 指南預期刪除 classic branch protection 後，`GET /branches/main/protection` 應為 404。本次 API 為 403（權限 `administration=read`）。v0.2 改以 Branches 設定截圖確認「Classic branch protections have not been configured」。

## RR-006 補償控制

管理員經規則二 bypass 合併時，適用 RR-006 補償控制：PR 留言須有審查部書面「可合併」結論、貼上四項 check 綠燈的 Actions run 連結，並由使用者本人按合併。

## 原始 JSON

- [`rulesets-2026-10-07/list.json`](rulesets-2026-10-07/list.json)
- [`rulesets-2026-10-07/ruleset-24629805.json`](rulesets-2026-10-07/ruleset-24629805.json)
- [`rulesets-2026-10-07/ruleset-24629854.json`](rulesets-2026-10-07/ruleset-24629854.json)
- [`rulesets-2026-10-07/rules-branches-main.json`](rulesets-2026-10-07/rules-branches-main.json)
- [`rulesets-2026-10-07/branch-main.json`](rulesets-2026-10-07/branch-main.json)
- [`rulesets-2026-10-07/branch-main-protection.json`](rulesets-2026-10-07/branch-main-protection.json)（HTTP 403 的回應本文）

## 截圖

擁有者於 2026-10-07 提供。檔案在 [`rulesets-2026-10-07/screenshots/`](rulesets-2026-10-07/screenshots/)。5 張皆為 GitHub Settings 頁裁切畫面，無網址列、無 repo 名稱、無時間戳；兩條 ruleset 的 `updated_at` 與 `created_at` 相差不到 0.1 秒（建立後未再修改），故截圖代表建立後現況。

| 檔案 | 說明 |
|---|---|
| [rule1-main-required-checks-bypass-empty.png](rulesets-2026-10-07/screenshots/rule1-main-required-checks-bypass-empty.png) | 規則一 `main-required-checks` 為 Active，Bypass list 顯示 “Bypass list is empty”。 |
| [rule1-main-required-checks-status-checks.png](rulesets-2026-10-07/screenshots/rule1-main-required-checks-status-checks.png) | 規則一要求四個 GitHub Actions 狀態檢查（`sca-govulncheck`、`sast-gosec`、`secrets-gitleaks`、`unit-test`），不要求 up to date，並勾選 Block force pushes。 |
| [rule2-main-pr-review-bypass.png](rulesets-2026-10-07/screenshots/rule2-main-pr-review-bypass.png) | 某 ruleset 為 Active；Bypass list 只有 “Repository admin (Roles)”，模式為 “Allow for pull requests only”；目標 Default；勾選 Restrict deletions。**ruleset 名稱未入鏡**，依排除法歸屬規則二 `main-pr-review`（見「規則二截圖的歸屬（排除法）」）。 |
| [branches-no-classic-protection.png](rulesets-2026-10-07/screenshots/branches-no-classic-protection.png) | Settings > Branches 顯示 “Classic branch protections have not been configured”。 |
| [rulesets-list-before-old-main-deleted.png](rulesets-2026-10-07/screenshots/rulesets-list-before-old-main-deleted.png) | 刪除舊 ruleset `main` 之前的清單（僅供上下文）。API [`list.json`](rulesets-2026-10-07/list.json) 確認該名稱現已不在清單中。 |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | API 快照；三項因權杖無 `administration=read` 記為 CANNOT VERIFY |
| v0.2 | 2026-10-07 | 維運部 | 以擁有者截圖關閉三項 CANNOT VERIFY；結論全數 PASS，無需修正 |
| v0.3 | 2026-10-07 | 維運部 | 依 G4-PR5 審查建議：S1 規則二截圖名稱未入鏡，改寫為排除法歸屬（不重拍），並訂後續截圖保留網址列；S3 CI 節改為以 PR 頁面最新 run 為準，移除首版 head 的舊 run 引用；S4 frontmatter「對應」補 ENG-020、NFR-004（SI-05）、SEC-009～SEC-011（狀態維持草稿待審查，合併後改定稿）；S5 結論改「13 項 PASS＋2 項已記錄（符合指南）」；舊 ruleset 404 補審查部重驗出處 |
