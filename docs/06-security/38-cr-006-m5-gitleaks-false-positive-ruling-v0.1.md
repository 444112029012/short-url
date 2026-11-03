---
文件：CR-006 M5 證據 gitleaks 命中裁定（安全部）
版本：v0.1
狀態：已裁定（2026-11-03 11:44 UTC+8 原裁定，本檔為獨立書面，回應審查部 PR33-S3）
負責角色：安全部
專案代號：SHORTURL
比照：33-pr-030-gitleaks-false-positive-ruling-v0.1.md（PR #30）
撰寫時間：開始於 2026-11-03 15:51:18（UTC+8，取自 `date`），完成時間以檔案 mtime 為準
---

# CR-006 M5 證據 gitleaks 命中裁定

## 1. 命中
docs/cr-006-closeout 本機 commit 3e19aa3f 上，gitleaks 掃全部歷史時回報 4 筆 `generic-api-key`：
| 檔案（兩個 M5 目錄各一份） | 行 | 內容 | 判定 |
|---|---|---|---|
| evidence/cr-006-m5/<head>/01-m5-1-head.txt | 17、19 | `(ii) API=<commit SHA>`（474756381bd170beb3a5d723b63f1b8ba48b2add、826d5d34d676e11f52e204e155ea271c7d811581） | 誤判：公開 repo 的 commit SHA |
| evidence/cr-006-m5/<head>/04-m5-4-gosec-buildinfo.txt | 49 | `google.golang.org/api v0.304.0 h1:2R1dYSiI…` | 誤判：公開模組校驗和 |
安全部 2026-11-03 11:44 以 `sed` 讀原檔核對內容（審查部 PR #33 紀錄 §4 另行核對相符）。

## 2. 裁定：命中的 4 個檔案只留本機，不放寬門檻，不遮蔽證據
- 不改 `.gitleaksignore`、不改 CI 設定（放寬 secrets 門檻須走 CR，理由同 33 號 §2）。
- 不遮蔽、不改寫（須保留原始書面證據，理由同 33 號 §2）。
- 原檔位置：`/workspace/projects/short-url/docs/06-security/evidence/cr-006-m5/` 兩目錄。
- SHA256SUMS 不改：47475638… 目錄 43d33bb47bed52571dff61b56776787ff4e838a8634da79565aa533122a659bd；826d5d34… 目錄 33c96bae0c0dc14d41e47bd3f0aff20f3c25c364d33edb0a8264be7ac19c2ef9。
- repo 內驗證必須用 `sha256sum -c --ignore-missing SHA256SUMS`；缺的 2 檔以本機原檔完整 `sha256sum -c` 為準。
- 須以重做 commit 方式排除，使整段分支歷史 gitleaks 0 筆（不可在後面加刪除 commit）。

## 3. 結果
PR #33 head 87943c91：4 檔不在檔案樹；gitleaks 全歷史 0 筆、secrets-gitleaks success（審查部 PR #33 紀錄 §3、§4 實測）。

## 4. 修訂紀錄
| 版本 | 日期 | 負責 | 說明 |
|---|---|---|---|
| v0.1 | 2026-11-03 | 安全部 | 初版：把 11:44 的訊息裁定與 36 號 v0.1.2 注意段落獨立成檔（PR33-S3） |
