---
文件編號：EVD-G3-RULESETS-20261007
文件：Rulesets API 快照（D-08 路徑 A）
版本：v0.2
狀態：草稿待審查
負責：維運部
日期：2026-10-07
專案代號：SHORTURL
對應：D-08、RR-006、`docs/09-cicd-env/06-rulesets-setup-guide-v0.1.md`
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
| 規則二：`bypass_actors` 恰為 Repository admin（`actor_type=RepositoryRole`、`actor_id=5`、`bypass_mode=pull_request`），且無其他項目 | PASS（人工確認：使用者截圖） | [僅 Repository admin，Allow for pull requests only](rulesets-2026-10-07/screenshots/rule2-main-pr-review-bypass.png) |
| 沒有其他指向 main 的 active ruleset | PASS | `includes_parents=true` 的清單只有上述兩筆，皆為 active |
| 舊 ruleset 名稱 `main` 不存在 | PASS | 清單中無此名稱 |
| classic branch protection 已不存在 | PASS（人工確認：使用者截圖） | [Classic branch protections have not been configured](rulesets-2026-10-07/screenshots/branches-no-classic-protection.png) |

`current_user_can_bypass` 在兩筆明細都是 `"never"`。這是呼叫端整合身分自己能否繞過，不是 `bypass_actors` 清單。

v0.1 有三列 CANNOT VERIFY。原因是 API 權杖沒有 `administration=read`：ruleset 明細（HTTP 200，接受權限 `metadata=read`）省略 `bypass_actors`；`GET /branches/main/protection` 為 403。這三項改由儲存庫擁有者於 2026-10-07 提供的設定畫面截圖人工確認，維運部已檢視，結果改記為 PASS（人工確認：使用者截圖）。

## 結論

全數 PASS，無需修正。

規則二同時含 `deletion` 與 `non_fast_forward`，與規則一重複，無害：規則一的 Bypass list 為空，四項 required check、禁刪除與禁 force push 沒有繞過名單。`require_extra_approval_for_unattributed_changes=true` 是 GitHub 目前的預設值。

## PR #5 CI

本證據首版（head `dd0a5c0c6475f7eeee4e55e260ca63b234261692`）的 `ci-devsecops` pull_request 跑次為綠燈：https://github.com/444112029012/short-url/actions/runs/37576407599 。四個 job（`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`）皆為 pass。

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

擁有者於 2026-10-07 提供。檔案在 [`rulesets-2026-10-07/screenshots/`](rulesets-2026-10-07/screenshots/)。

| 檔案 | 說明 |
|---|---|
| [rule1-main-required-checks-bypass-empty.png](rulesets-2026-10-07/screenshots/rule1-main-required-checks-bypass-empty.png) | 規則一 `main-required-checks` 為 Active，Bypass list 顯示 “Bypass list is empty”。 |
| [rule1-main-required-checks-status-checks.png](rulesets-2026-10-07/screenshots/rule1-main-required-checks-status-checks.png) | 規則一要求四個 GitHub Actions 狀態檢查（`sca-govulncheck`、`sast-gosec`、`secrets-gitleaks`、`unit-test`），不要求 up to date，並勾選 Block force pushes。 |
| [rule2-main-pr-review-bypass.png](rulesets-2026-10-07/screenshots/rule2-main-pr-review-bypass.png) | 規則二 `main-pr-review` 為 Active；Bypass list 只有 “Repository admin (Roles)”，模式為 “Allow for pull requests only”；目標 Default；勾選 Restrict deletions。 |
| [branches-no-classic-protection.png](rulesets-2026-10-07/screenshots/branches-no-classic-protection.png) | Settings > Branches 顯示 “Classic branch protections have not been configured”。 |
| [rulesets-list-before-old-main-deleted.png](rulesets-2026-10-07/screenshots/rulesets-list-before-old-main-deleted.png) | 刪除舊 ruleset `main` 之前的清單（僅供上下文）。API [`list.json`](rulesets-2026-10-07/list.json) 確認該名稱現已不在清單中。 |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | API 快照；三項因權杖無 `administration=read` 記為 CANNOT VERIFY |
| v0.2 | 2026-10-07 | 維運部 | 以擁有者截圖關閉三項 CANNOT VERIFY；結論全數 PASS，無需修正 |
