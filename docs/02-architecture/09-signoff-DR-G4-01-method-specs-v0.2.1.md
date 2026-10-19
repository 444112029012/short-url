# 設計部書面簽核：方法規格 v0.2.1（DR-G4-01／PR #17）

文件編號：09-signoff-DR-G4-01-method-specs-v0.2.1
版本：v1.0
日期：2026-10-19（Asia/Taipei）
簽核者：設計部（方法規格表負責角色）
狀態：已簽核

## 一、綁定標的

| 項目 | 值 |
|---|---|
| PR | https://github.com/444112029012/short-url/pull/17 |
| 綁定 head（H） | `5b5017b2b9109b9e547e9273787fd6c790ae92d5` |
| H 的 parents | `347a02759552d8f8a1e18484eaa3b70fb5fe55e1`、`ed620394b1c9986cc4ac243eb1543c8b677d92cb`（main，含 #16） |
| H 上 `docs/02-architecture/06-method-specs-v0.2.md` 之 sha256 | `c89e336e84dc524f74961dab9874bde1d8669cec7373e798f5edef5a35aea687` |
| H 的 CI run | 39203891218（ci-devsecops，head_sha＝H，conclusion＝success；unit-test、sast-gosec、sca-govulncheck、secrets-gitleaks 四項 success） |

## 二、必改項關閉狀態

| 項目 | 內容 | 狀態 | 首次確認 head |
|---|---|---|---|
| 1 | §13.4.1 補回裁定全文並改為已生效寫法；外框 `{"url":""}`；刪「現行 8192」 | 已關閉 | 3b1aa5b |
| 2 | §13.2／§13.4 行號改為函式／常數名引用，內容正確 | 已關閉 | 3b1aa5b |
| 3 | §13.4 新增標頭上限列：MaxHeaderBytes＝16384，超過由 net/http 回 431（text/plain，不經 ErrorMapper，非錯誤信封，不屬四類錯誤） | 已關閉 | 3b1aa5b |
| 4 | metadata 版本 v0.2.1、狀態更新；修訂紀錄單列 v0.2.1，作者設計部（研發部草擬） | 已關閉 | 3b1aa5b |
| 5 | C4(a) 測試：`https://example.com/`＋2028 個 U+1F600 surrogate pair 跳脫，201、本文 >8192 且 ≤32768、存入 2048 rune | 已關閉 | 3b1aa5b |
| R2（SCR-040） | §13.2 RateLimitGuard／ObservabilityHooks 改為現行實作（PR #7、#11），刪 stub／no-op | 已關閉 | 347a027 |
| R3 | §13.3／§13.5 改為名稱引用，§13 無殘留行號 | 已關閉 | 347a027 |

H 與 347a027 之規格檔 sha256 相同（設計部於 347a027 實算並複核通過）；H 與 main 之差異僅 PR #17 原 5 檔。

## 三、聲明

1. 方法規格 v0.2.1 以 PR #17 於 H 之版本為準；設計部本機版不再另行提交（SCR-026、SCR-041），合併後本機檔同步為 repo 版。
2. 本簽核僅在下列條件下有效：本檔以唯一 parent 為 H 之 commit 加入，且該 commit 僅新增本檔。H 之後若有其他 commit 或 main 前進，須自 H 重新複核與重簽。
3. 非阻擋待辦（不影響本簽核）：metadata／§13 前言之「對照依據 PR #2 head」、02-architecture/README 索引版本，由設計部合併後文件 PR 處理；安全部意見書 10-dr-g4-01 由安全部補入 repo；若合併日非 2026-10-19，文件 PR 校正 §13.4.1 生效日。

設計部 2026-10-19
