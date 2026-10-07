---
文件：G3-R1 落地生效證據
版本：v0.2
狀態：完成
負責角色：維運部
最後更新：2026-10-07
專案代號：SHORTURL
對應：G3-R1（workflow 落地＋main 分支保護）；ENG-020；NFR-004（SI-05）；SEC-009、SEC-010、SEC-011
---

# G3-R1 落地生效證據

> **時點說明**：本證據記錄 2026-10-07 12:15（Asia/Taipei）時點的 **classic branch protection**，屬歷史時點。之後已依 D-08 改為 **Rulesets**（ruleset 24629805 `main-required-checks`、24629854 `main-pr-review`；classic 已刪除），現況證據見 [`rulesets-snapshot-2026-10-07.md`](./rulesets-snapshot-2026-10-07.md)。

## 1. Workflow 已入庫

| 項目 | 值 |
|---|---|
| Repo | https://github.com/444112029012/short-url |
| 路徑 | `.github/workflows/ci-devsecops.yml` |
| Seed commit | [`d09cceb`](https://github.com/444112029012/short-url/commit/d09cceb4d2284f2dbcb99e6f1a5833cd9baecd66) |
| 四 required job 名稱 | `secrets-gitleaks`／`sast-gosec`／`sca-govulncheck`／`unit-test` |
| 備註 | 取證當時 CI Go toolchain 為 `1.24.x`（避開 1.22 stdlib `GO-2025-3750`）；掃描門檻未放寬。ENG-007／PR #2 合併後改為 Go `1.27.x`（gosec v2.29.0、govulncheck v1.8.0；D-07），影響說明見 [`g3-addendum-pr2-toolchain-2026-10-07.md`](./g3-addendum-pr2-toolchain-2026-10-07.md) |

Workflow 連結：https://github.com/444112029012/short-url/blob/main/.github/workflows/ci-devsecops.yml

## 2. Actions 綠燈（四 job 皆 success）

| 項目 | 值 |
|---|---|
| Run | https://github.com/444112029012/short-url/actions/runs/37569477572 |
| Event | push → `main` |
| Conclusion | **success** |
| secrets-gitleaks | success |
| sast-gosec | success |
| sca-govulncheck | success |
| unit-test | success |

## 3. main 分支保護生效

取證時點：2026-10-07 12:15（Asia/Taipei），classic branch protection。

取證來源：cursor-github `list_branches`（protected=true）輸出的 `main` 分支欄位，與公開分支 API `GET https://api.github.com/repos/444112029012/short-url/branches/main` 的同名欄位相同（審查部於 G3 複審以 WebFetch 該公開 API 獨立確認相同值）。維運部整理的**摘要**（非原始 API 回應）存於 [`branches-main-protection-summary-2026-10-07.json`](./branches-main-protection-summary-2026-10-07.json)。

摘錄（摘要）：

```json
{
  "name": "main",
  "protected": true,
  "protection": {
    "enabled": true,
    "required_status_checks": {
      "enforcement_level": "everyone",
      "contexts": [
        "secrets-gitleaks",
        "sast-gosec",
        "sca-govulncheck",
        "unit-test"
      ]
    }
  }
}
```

| 檢查項 | 結果 |
|---|---|
| `protected` | true |
| Required checks 四名完全一致 | 通過 |
| `enforcement_level` | `everyone`（管理員亦適用） |
| 使用者確認 | 須 PR、至少 1 核准、禁 force push／刪分支（Settings 已設；完整 protection JSON 需 Administration token，Cursor App 無此權限故以公開 API＋使用者確認補齊） |

## 4. 結論

G3-R1 審查退回項已落地並可複審：workflow 入庫、四 check 綠燈、main 分支保護（含四 required checks）生效。

## 5. 後續變更（2026-10-07）

- **工具鏈**：本證據的綠燈 run 基於 Go 1.24.x、gosec v2.21.4、govulncheck v1.1.3（main 當時版本）。ENG-007／PR #2 依 D-07 改為 Go 1.27.x、gosec v2.29.0、govulncheck v1.8.0；四個 required job 名稱與 High／Critical 阻擋語意不變。詳見 [`g3-addendum-pr2-toolchain-2026-10-07.md`](./g3-addendum-pr2-toolchain-2026-10-07.md)。
- **合併路徑**：依 D-08 改為 Rulesets（方案 A），classic protection 已刪除；現況見 [`rulesets-snapshot-2026-10-07.md`](./rulesets-snapshot-2026-10-07.md)。

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3-R1 落地生效證據（12:15 時點 classic branch protection；run 37569477572） |
| v0.2 | 2026-10-07 | 維運部 | 依 G4-PR1 審查建議：S1 加時點說明（classic → 依 D-08 改 Rulesets，指向 `rulesets-snapshot-2026-10-07.md`）；S2 刪除重複 JSON，只保留並標明摘要檔 `branches-main-protection-summary-2026-10-07.json`，來源描述與 JSON `source` 欄一致；S4 新增修訂紀錄；S5 註明當時 toolchain 1.24.x、ENG-007／PR #2 後為 1.27.x，連結 `g3-addendum-pr2-toolchain-2026-10-07.md`；本機先前附加、未推送的「補註」節與 g3-addendum 重複，併入第 5 節後刪除；對應補 ENG-020、NFR-004、SEC-009～011 |
