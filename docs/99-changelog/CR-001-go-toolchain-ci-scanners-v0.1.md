---
文件：變更請求 CR-001
版本：v0.1
日期：2026-10-07
狀態：Pending Approval（安全／維運已有條件同意方向；待正式裁示定稿＋使用者拍板）
提出：研發部
專案代號：SHORTURL
關聯：PR https://github.com/444112029012/short-url/pull/2；SEC-009／SEC-010；D-04 工具鏈約束旁支
時區：Asia/Taipei（UTC+8）
---

# CR-001｜Go 工具鏈與 CI 掃描器升版／dependency-review 暫行 fail-open

## 1. 摘要

為使 SCA／SAST 在受支援 Go 線上可編譯並通過漏洞掃描，將 Go 自 **1.24.x** 升至 **1.27.x**，並同步釘版升級 `govulncheck`／`gosec`。另因倉庫 Dependency graph 未啟用，對 `dependency-review-action` 採探測後略過（fail-open），**略過必須發出 GitHub Actions `::warning::`**，且 **`sca-govulncheck` 維持 required**；暫行最遲 **G5 送審前**失效並改回 fail-closed。

## 2. 變更內容

| 項目 | 變更前 | 變更後 | 理由 |
|---|---|---|---|
| Go toolchain | 1.24.x（維運避 GO-2025-3750） | **1.27.x**（`go.mod` 釘 1.27.1；CI `setup-go` 1.27.x） | 1.24.13 為 1.24 最後版；現行 vuln DB 仍標可達標準庫問題，修補僅在受支援版本 |
| gosec | v2.21.4 | **v2.29.0** | 舊版無法於 Go 1.25+ 編譯；`-severity=high`／job 名不變 |
| govulncheck | v1.1.3 | **v1.8.0** | 同上；非零退出仍擋 |
| dependency-review | PR 必跑 action | 探測 compare API＝200 才跑；否則 **`::warning::` 略過** | 依賴圖未開時 action 直接失敗；主擋仍為 govulncheck |
| Required checks | 四項 job 名 | **不變**（含 `sca-govulncheck`） | 分支保護無需改 |

### 略過警告文案（必須可見於 Actions Annotations）

```bash
echo "::warning::dependency-review skipped: Dependency graph disabled; fail-open until enabled; expires before G5"
```

## 3. 影響與風險

- **正面**：govulncheck 可在支援線回報無已知漏洞；避免長期紅燈或被迫放寬門檻。
- **風險**：dependency-review fail-open 期間，若依賴圖仍未開，PR 依賴差分審查弱於預期 → 以 govulncheck required＋G5 前失效補償。
- **範圍外**：不放寬 High／Critical 阻擋規則；不改四個 job 名稱。

## 4. 驗收條件

1. PR #2（或後續同變更分支）workflow 在依賴圖未開時：`sca-govulncheck` **success**，且 Annotations 可見上述 `::warning::`。
2. 依賴圖開啟後：dependency-review 實際執行；Critical 仍擋；可移除暫行略過。
3. `unit-test`／`sast-gosec`／`secrets-gitleaks`／`sca-govulncheck` 仍為 required。
4. 安全／維運正式裁示定稿後，使用者拍板；決策記入 `decision-log.md`（建議 D-0x）。

## 5. 回滾

若裁示否決：還原 `go.mod`／workflow 至 1.24.x＋原掃描器釘版，並另開 CR 處理 1.24 EOL／漏洞庫衝突（可能需接受紅燈或改門檻——不在本 CR 範圍）。

## 6. 依據與意見

- 維運意見：`docs/09-cicd-env/05-pr2-ci-change-review-v0.1.md`（有條件同意）
- 總協調指示：略過須 `::warning::`；完整裁示定稿前不合併 PR #2
- 門檻文件：`docs/06-security/03-ci-security-gates-*.md`、`docs/09-cicd-env/01-ci-pipeline-*.md`

## 7. 簽核

| 角色 | 狀態 | 備註 |
|---|---|---|
| 研發部（提出） | 已提出 | 本文件 v0.1 |
| 安全部 | 待正式定稿 | 方向已透過總協調／維運文件轉達 |
| 維運部 | 有條件同意 | 見 05-pr2-ci-change-review |
| 審查部 | 待審 | |
| 使用者（王昱殊） | 待拍板 | 合併前 |

