---
文件：G3-R1 落地檢查清單（workflow／分支保護）
版本：v0.1
狀態：進行中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05；SEC-009～SEC-011
專案代號：SHORTURL
---

# G3-R1 落地檢查清單

> 審查退回 G3-R1：須於**實作 repo** 落地 workflow 並啟用 `main` 分支保護 required checks，附生效證據。  
> 本機已備妥：`.github/workflows/ci-devsecops.yml`、`.gitignore`、`.env.example`、最小 Go scaffold（供 CI 可跑）。

## 前置（阻塞）

- [ ] 使用者建立 GitHub **空庫**（建議名：`short-url`；帳號 `444112029012`）
- [ ] 將 repo URL 交給總協調／維運
- [ ] Cursor GitHub 連線可存取該 repo（若工具回 not found，走 access 授權）

## 落地步驟（repo URL 到齊後）

1. 推入（或 cloud agent 提交）本專案樹：`docs/`、`.github/workflows/ci-devsecops.yml`、`.gitignore`、`.env.example`、`go.mod`／scaffold  
2. 確認 Actions 於 `main`／PR 跑出四個 check：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`  
3. 設定分支保護（Ruleset 或 Branch protection）：
   - 禁直推 `main`、禁 force push
   - 至少 1 核准
   - Required status checks＝上列四個（名稱完全一致）
   - 管理員亦適用（或繞過須留紀錄）
4. （建議）啟用 secret scanning＋push protection  
5. 擋測證據（擇一附檔／截圖／run URL）：
   - 假機密 commit → `secrets-gitleaks` 紅
   - 或引用一次綠燈 run＋分支保護設定頁截圖／API 輸出

## 證據交付物（關單用）

| 項目 | 證據形式 |
|---|---|
| Workflow 已入庫 | 路徑 `.github/workflows/ci-devsecops.yml`＋commit／PR URL |
| 四 job 曾通過或曾阻擋 | Actions run URL |
| 分支保護生效 | Settings 截圖或 `gh api` rulesets／protection JSON |

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3-R1：備妥落地檔與檢查清單；待空庫 URL |
