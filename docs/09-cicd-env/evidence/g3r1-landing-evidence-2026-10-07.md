---
文件：G3-R1 落地生效證據
版本：v0.1
狀態：完成
負責角色：維運部
最後更新：2026-10-07
專案代號：SHORTURL
對應：G3-R1（workflow 落地＋main 分支保護）
---

# G3-R1 落地生效證據

## 1. Workflow 已入庫

| 項目 | 值 |
|---|---|
| Repo | https://github.com/444112029012/short-url |
| 路徑 | `.github/workflows/ci-devsecops.yml` |
| Seed commit | [`d09cceb`](https://github.com/444112029012/short-url/commit/d09cceb4d2284f2dbcb99e6f1a5833cd9baecd66) |
| 四 required job 名稱 | `secrets-gitleaks`／`sast-gosec`／`sca-govulncheck`／`unit-test` |
| 備註 | CI Go toolchain 為 `1.24.x`（避開 1.22 stdlib `GO-2025-3750`）；掃描門檻未放寬 |

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

驗收來源（公開分支 API，2026-10-07）：  
`GET https://api.github.com/repos/444112029012/short-url/branches/main`

摘錄：

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
