---
文件：G3-R1 落地檢查清單（workflow／分支保護）
版本：v0.2
狀態：完成
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05；SEC-009～SEC-011；ENG-020
專案代號：SHORTURL
---

# G3-R1 落地檢查清單

> 審查退回 G3-R1：須於**實作 repo** 落地 workflow 並啟用 `main` 分支保護 required checks，附生效證據。  
> **狀態：已完成**（證據見 `evidence/g3r1-landing-evidence-2026-10-07.md`）  
> **時點說明**：本清單記錄的是 2026-10-07 12:15（Asia/Taipei）時點的 **classic branch protection**；之後已依 D-08 改為 **Rulesets**（ruleset 24629805 `main-required-checks`、24629854 `main-pr-review`；classic 已刪除），現況證據見 `evidence/rulesets-snapshot-2026-10-07.md`。

## 前置

- [x] 使用者建立 GitHub 空庫（`444112029012/short-url`）
- [x] 將 repo URL 交給總協調／維運
- [x] Cursor GitHub 連線可存取該 repo

## 落地步驟

- [x] 推入落地包（docs／workflow／scaffold）→ commit `d09cceb`
- [x] Actions 四 check 綠燈 → [run 37569477572](https://github.com/444112029012/short-url/actions/runs/37569477572)
- [x] main 分支保護（四 required checks；`protected: true`，`enforcement_level: everyone`）——12:15 時點為 classic protection，已依 D-08 改為 Rulesets
- [x] secret scanning＋push protection（repo 已 enabled）

## 證據交付物

| 項目 | 證據 |
|---|---|
| Workflow 已入庫 | https://github.com/444112029012/short-url/blob/main/.github/workflows/ci-devsecops.yml ＋ [d09cceb](https://github.com/444112029012/short-url/commit/d09cceb4d2284f2dbcb99e6f1a5833cd9baecd66) |
| 四 job 通過 | https://github.com/444112029012/short-url/actions/runs/37569477572 |
| 分支保護生效（12:15 時點，classic） | `GET .../branches/main` → `protected: true`；詳見 `evidence/g3r1-landing-evidence-2026-10-07.md` 與摘要 `evidence/branches-main-protection-summary-2026-10-07.json` |
| 分支保護現況（D-08 後，Rulesets） | `evidence/rulesets-snapshot-2026-10-07.md` |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | 備妥落地檔與檢查清單；待空庫 URL |
| v0.2 | 2026-10-07 | 維運部 | 落地完成並取證；G3-R1 可送複審 |
| v0.2 | 2026-10-07 | 維運部 | 依 G4-PR1 審查建議：S3 另存 `-v0.2.md` 使檔名與 frontmatter 一致（`-v0.1.md` 已在 main，保留不動）；S1 註明 12:15 時點為 classic protection、已依 D-08 改 Rulesets 並指向 `rulesets-snapshot-2026-10-07.md`；S2 證據 JSON 只引用摘要檔；對應需求補 ENG-020 |
