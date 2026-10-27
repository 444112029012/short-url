---
文件：G4-QA-R2 處置紀錄（方案 B：經 PR #27、PR #28 進 main）
版本：v0.2
狀態：審查中
負責角色：品保部
最後更新：2026-10-27（Asia/Taipei，UTC+8）
依據：`19-G4-overall-qa-conclusion-v0.1.md` §6（sha256 `86606bd3cfb059c52dc1ac9984161d9c8cc879eeeebec9c46795f4d4960b0d70`）；`12-gate-G4-qa-precheck-v0.1.md` §8 G4-QA-R2；本文 v0.1（sha256 `dcc63af8f886f7dca13a8130a17967525786c29519b1f33eaf1c75cd73118dfb`）§0、§3、§4 所列關閉條件
比對基準：main `36f4717caab8a957e60f0a67ef5b01400f69ef52`（PR #28 merge commit）；另對最新 main `db380a40c3330bf23498ae7fda2f123cf7dd4cc4`（PR #29）複核
專案代號：SHORTURL
---

# G4-QA-R2 處置紀錄

## 0. 結論：**G4-QA-R2 結案（關閉）**

本文 v0.1 列出的關閉條件，品保已用 GitHub 唯讀方式逐項查證，**全部符合**：

1. v0.1 §3 列的 16 個缺漏檔（7 份舊版文件＋9 份證據檔）已隨 PR #28 進 main，main 與本機 sha256 **16／16 一致**。
2. 13 號（v0.1.1，sha256 `8f692c33…`）main 與本機一致。
3. 20 號（本文 v0.1，sha256 `dcc63af8…`）main 與本機一致（自我參照，說明見 §1）。
4. v0.1 §4 列的 main 上 6 處簡體字（04×2、05、06、10、README）在 `36f4717c` 與最新 main `db380a40` 皆為 **0 處**。
5. 另外複核 v0.1 §2 原本一致的 19 檔，至今仍一致。07-testing 全目錄 **37／37 檔一致，0 不一致**，main 上沒有本機以外的多餘檔，本機也沒有 main 缺少的檔。

依此，G4-QA-R2（07-testing 文件更新進 main）**結案**；DEF-004 在 main 上也已全數修正。

> 「PR #28 已於 2026-10-27 12:01（UTC+8）合併，merge commit `36f4717c`」原為 agent 轉述；品保已自行用 GitHub 唯讀查證，結果見 §2，與轉述相符。
> 「使用者已簽核 G4 整體關卡」：仍屬 agent 轉述，品保未直接查證。本結案不以此為前提。

## 1. 自我參照說明（README 與本文）

`README.md` 與本文（20 號）都在 PR #28 的變更清單內，而這次又會修改這兩個檔案，屬於自我參照。處理原則如下：

- 本文的比對與結案判定，**以本次修改前的內容為準**：README 為 v0.4（sha256 `25c1807b8aa8c47a71dd8d635fa07d3331910a5a27b6c68eb2a8a1c40edd030c`），本文為 v0.1（sha256 `dcc63af8f886f7dca13a8130a17967525786c29519b1f33eaf1c75cd73118dfb`）。這兩個值在 main `36f4717c` 與本機修改前完全相同。
- 修改後的新版（本文 v0.2、README v0.5）尚未進 main，會進下一個純文件 PR。合併前 main 與本機在這兩檔上的 sha256 不同屬預期，不列為差異，也不影響本次結案。
- README v0.5 的索引列記有本文 v0.2 的 sha256。本文無法反向記入 README v0.5 的 sha256（兩檔互相引用會形成循環），因此本文只記 README 修改前的 sha256 作為比對基準；README v0.5 的 sha256 另見交接紀錄。

## 2. PR #28 查證（gh／GitHub API 唯讀，2026-10-27 17:04～17:08 UTC+8）

| 項 | agent 轉述 | 品保實查 | 判定 |
|---|---|---|---|
| 狀態 | 已合併 | `gh pr view 28`：state MERGED；`gh api repos/444112029012/short-url/pulls/28`：merged=true | ✅ |
| 合併時間 | 2026-10-27 12:01（UTC+8） | mergedAt `2026-10-27T04:01:01Z`（＝2026-10-27 12:01:01 UTC+8）；merge commit committer 時間 `2026-10-27T12:01:01+08:00` | ✅ |
| merge commit | `36f4717c` | merge_commit_sha `36f4717caab8a957e60f0a67ef5b01400f69ef52` | ✅ |
| 合併方式 | — | 一般合併（merge commit）：該 commit 有兩個 parent，`2425129b`（main，PR #27 merge commit）＋`f1c96bc1`（PR head `f1c96bc1c0ea7638f666127481bec9497e3f1220`），訊息為「Merge pull request #28 from 444112029012/docs/g4-followups」。GitHub API 不直接回傳合併方式，此為依 parent 數與訊息推定；非 squash、非 rebase | ✅（推定） |
| 合併者 | — | merged_by `444112029012` | 記錄 |
| 分支／標題 | — | `docs/g4-followups` → `main`；「docs: G4 整體關卡收尾（S2～S6、QA-R2 檔案入庫、D-18）」；2 個 commit（`7f64b150`、`f1c96bc1`）；24 檔 | — |
| checks | — | PR 四項（secrets-gitleaks、sast-gosec、sca-govulncheck、unit-test）皆 SUCCESS（2026-10-27 09:05～09:06 UTC+8）；main push run `40447379645`（ci-devsecops）success（2026-10-27 12:01:04 UTC+8 觸發） | ✅ |
| 審查狀態 | — | reviewDecision REVIEW_REQUIRED，reviews 為空（無核准，推定以 admin bypass 合併；品保未查到 bypass 紀錄） | 記錄 |
| 含 07-testing 檔 | — | 18 檔：04、05、06、07、09、10、13、20（新增）、README，以及 evidence 9 檔（新增） | ✅ |

最新 main：`db380a40c3330bf23498ae7fda2f123cf7dd4cc4`（PR #29，2026-10-27 16:37:12 UTC+8），parent 為 `36f4717c`＋`1545143c`；與 `36f4717c` 相比只改 `.github/workflows/ci-devsecops.yml`，`docs/` 完全相同，故以下比對兩個 commit 結果一致。

## 3. 逐檔比對（sha256，2026-10-27 17:05 UTC+8 重算）

方法：`git clone` 取 origin/main，以 `git show <commit>:<path> | sha256sum` 計算 main 值，以 `sha256sum` 計算本機值（README、20 號為本次修改前的值）。`/workspace/dev/short-url` 停在舊分支，不作為 main 依據。

### 3.1 v0.1 §3 原列缺漏檔（16 檔）

| # | 檔案（`docs/07-testing/`） | 版本 | main `36f4717c` sha256 | 本機 sha256 | 判定 |
|---|---|---|---|---|---|
| 1 | 04-design-testability-review-v0.1.md（DEF-004） | v0.1.1 | 12002fa436755e189f46769318d65c87e44fef73298460dd31def4df13a5b3ca | 12002fa436755e189f46769318d65c87e44fef73298460dd31def4df13a5b3ca | ✅ 一致 |
| 2 | 05-gate-G2-qa-precheck-v0.1.md（DEF-004） | v0.1.1 | b5d018d19a7804edb6d572ac9ea602f2ca9c230446c1b1178ce3b698c09c5b1d | b5d018d19a7804edb6d572ac9ea602f2ca9c230446c1b1178ce3b698c09c5b1d | ✅ 一致 |
| 3 | 06-test-strategy-v0.1.md（DEF-004） | v0.1 | 409187e193ada518970957f0fc24c613fb2c5906c4ec667ef679dcf60d8e557d | 409187e193ada518970957f0fc24c613fb2c5906c4ec667ef679dcf60d8e557d | ✅ 一致 |
| 4 | 07-test-cases-v0.1.md | v0.2 | ef8ed11d6e8324b3d52a5cb6d597b12b13c6e0896730dd7ef458672978f098e5 | ef8ed11d6e8324b3d52a5cb6d597b12b13c6e0896730dd7ef458672978f098e5 | ✅ 一致 |
| 5 | 09-dev-task-acceptance-map-v0.1.md | v0.3 | f7ca9727ab9c594fbe5c01efe3ed18c2542115d4d7656e2d71fca2ca38568792 | f7ca9727ab9c594fbe5c01efe3ed18c2542115d4d7656e2d71fca2ca38568792 | ✅ 一致 |
| 6 | 10-gate-G3-qa-precheck-v0.1.md（DEF-004） | v0.1 | d229e23e5ef87d41790288ff65959869d9182767ab3bf993c80561b7df0cb344 | d229e23e5ef87d41790288ff65959869d9182767ab3bf993c80561b7df0cb344 | ✅ 一致 |
| 7 | README.md（修改前） | v0.4 | 25c1807b8aa8c47a71dd8d635fa07d3331910a5a27b6c68eb2a8a1c40edd030c | 25c1807b8aa8c47a71dd8d635fa07d3331910a5a27b6c68eb2a8a1c40edd030c | ✅ 一致 |
| 8 | evidence/g4-full-run-2026-10-22.txt | — | 8dd97911aa38c2c4248ad69e72975bdd261f3396011d56788a9f2756e4eeb5d8 | 8dd97911aa38c2c4248ad69e72975bdd261f3396011d56788a9f2756e4eeb5d8 | ✅ 一致 |
| 9 | evidence/s6-eng011-run-2026-10-22.txt | — | c381c361bbbe2d851625d126ac951d9168b1098d907d3c94826d2b336413d4a4 | c381c361bbbe2d851625d126ac951d9168b1098d907d3c94826d2b336413d4a4 | ✅ 一致 |
| 10 | evidence/pr26-cover.txt | — | eb004f4c1541a95f003ee5e07d721692d356c606c5e6980850b60caa333f1ac1 | eb004f4c1541a95f003ee5e07d721692d356c606c5e6980850b60caa333f1ac1 | ✅ 一致 |
| 11 | evidence/pr26-env.txt | — | 3231de7959ade2625e9113e67fd4b6df1f84283f22ff93ac0a0315a3adb7ec6c | 3231de7959ade2625e9113e67fd4b6df1f84283f22ff93ac0a0315a3adb7ec6c | ✅ 一致 |
| 12 | evidence/pr26-gh-ci.txt | — | 592a4386cb87e3e91f067d15416c576c573b86c4714b29275530b838ecef8ac5 | 592a4386cb87e3e91f067d15416c576c573b86c4714b29275530b838ecef8ac5 | ✅ 一致 |
| 13 | evidence/pr26-gh-pr.json | — | 712f3a2a7a371cf4ee84e50decabfb7f4a7b75cccd220c59b80e9f5472e0af4c | 712f3a2a7a371cf4ee84e50decabfb7f4a7b75cccd220c59b80e9f5472e0af4c | ✅ 一致 |
| 14 | evidence/pr26-mutation.txt | — | 7b80f35803c72a54ed05761126b13176ce1c9d39b4ee3e3a10e269312e6f091e | 7b80f35803c72a54ed05761126b13176ce1c9d39b4ee3e3a10e269312e6f091e | ✅ 一致 |
| 15 | evidence/pr26-race.txt | — | 5a18e0414f56dcc71f1f7d094462392879f2a1a093b820f3898c15e0d2376cb3 | 5a18e0414f56dcc71f1f7d094462392879f2a1a093b820f3898c15e0d2376cb3 | ✅ 一致 |
| 16 | evidence/pr26-vet.txt | — | 19eaf43821a7660ec323a87c8457bf74823beb296c39f5e01aa8a683aa50f061 | 19eaf43821a7660ec323a87c8457bf74823beb296c39f5e01aa8a683aa50f061 | ✅ 一致 |

### 3.2 13 號與 20 號

| 檔案 | 版本 | main `36f4717c` sha256 | 本機 sha256 | 判定 |
|---|---|---|---|---|
| 13-G4-QA-R1-acceptance-PR26-v0.1.md | v0.1.1 | 8f692c33ba5a88ee191b1320f84c5dd9fd1a27981a7c04b2595f1bf3376b4c43 | 8f692c33ba5a88ee191b1320f84c5dd9fd1a27981a7c04b2595f1bf3376b4c43 | ✅ 一致 |
| 20-G4-QA-R2-closure-note-v0.1.md（本文，修改前） | v0.1 | dcc63af8f886f7dca13a8130a17967525786c29519b1f33eaf1c75cd73118dfb | dcc63af8f886f7dca13a8130a17967525786c29519b1f33eaf1c75cd73118dfb | ✅ 一致 |

### 3.3 v0.1 §2 原一致檔複核（19 檔，仍一致）

| 檔案（`docs/07-testing/`） | sha256（main `36f4717c`＝本機） | 判定 |
|---|---|---|
| 01-requirements-testability-review-v0.1.md | 120ddbf65e649aa196c892dcb88b5be6956adf30a7a59ee9a0ad66cf786724b3 | ✅ |
| 02-traceability-matrix-v0.1.md（v0.5） | 8a005aefcf24a32fa769c803373e88a3bb411aa71652ff14426feb18907b9932 | ✅ |
| 03-gate-G1-qa-precheck-v0.1.md | dbabd14abd278b16d87d554868d5c017e0ace08f4b0904d97da99aded24987c1 | ✅ |
| 08-definition-of-done-v0.1.md | be865c1cd31e6ea7b25e0bf84f09088182bc86acf6efaaf5e2ebd198fd33dd06 | ✅ |
| 11-s6-eng011-disposition-v0.1.md | 02e2973d4686d2945a31873d8e9acb701974d208209d1188e1aa96c1511217b4 | ✅ |
| 12-gate-G4-qa-precheck-v0.1.md | 85831dc8c12f6b26a1899e55b5f534f7fa11de2e6c1cce1f501420ccdb29f09b | ✅ |
| 14-CR-003-qa-confirmation-v0.1.md | 3fb5d45727a337d9d9901ed7a5399922457d166b135cba10845b408f68b30ae9 | ✅ |
| 15-CR-004-R3-qa-impact-row-v0.1.md | 4c3684fbfb92ed4463a24cdd4a6d0e14349ce4418bfb7315b440a49f1617ade2 | ✅ |
| 16-CR-004-R4-qa-opinion-v0.1.md | ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c | ✅ |
| 17-CR-004-signoff-2-2-qa-confirmation-v0.1.md | e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587 | ✅ |
| 18-CR-004-AC6-9-qa-verdict-v0.1.md（v0.4） | 543763dad54eb6dd1bbd3b68afcdfa4c80e0b375cfea912fb24a02ef02d61288 | ✅ |
| 19-G4-overall-qa-conclusion-v0.1.md | 86606bd3cfb059c52dc1ac9984161d9c8cc879eeeebec9c46795f4d4960b0d70 | ✅ |
| evidence/g4-overall-f646-cover.txt | 2d4a41869fcc9f106babebcc8eb4bb95f2953bb0d4d5efbe0c0c31397c269c29 | ✅ |
| evidence/g4-overall-f646-ct.txt | e363d42c357983b5f34a117d508e459345af12b641ef9f020e5b4cd2c486e990 | ✅ |
| evidence/g4-overall-f646-env.txt | 73138c31bfb06fb5d5699641b93b029d364cd2a65824f0d2d431a583e9eec537 | ✅ |
| evidence/g4-overall-f646-race.txt | 42f375378befa50f51b9255f994c4bde7e54a00b04804f0d78d86abdca21ad86 | ✅ |
| evidence/g4-overall-f646-test.json | 71353d9a8d74dea36ca6fb9916b07cfd93bfba0f6f2dafcd50e0867a5bcb89a7 | ✅ |
| evidence/g4-overall-f646-vet.txt | 00a5f1d995f1697b6f134bcca3e00d503345a5241662ca72a16f8b0ac5768294 | ✅ |
| evidence/pr26-ct-run.txt | 3209871dd317775c79574f3ecc5750375dedf8553bf4cf7d2f007066431fc1bd | ✅ |

### 3.4 小計

- 原列結案條件檔：16＋2＝**18 檔，一致 18，不一致 0**。
- 07-testing 全目錄：main `36f4717c` 37 檔，本機 37 檔，**一致 37，不一致 0**，兩邊都沒有多出的檔。
- 最新 main `db380a40`：37 檔 sha256 與 `36f4717c` 完全相同。

## 4. DEF-004（簡體字修正）是否進 main：**是（全數修正）**

- 指令：`git grep -n -E '<DEF-004 所列兩個簡體字>' <commit> -- docs/07-testing`（字樣見 12 號 §7 DEF-004 重現步驟；本文不重抄字樣，以免本文本身含簡體字）。
- PR #27 merge `2425129b`（對照）：8 處，即 v0.1 §4 列的 6 處（04 號第 126、141 行；05 號第 33 行；06 號第 114 行；10 號第 130 行；README 第 39 行）＋12 號第 99、111 行。
- PR #28 merge `36f4717c`：**只剩 12 號第 99、111 行**，屬缺陷描述中引用字樣、記錄用途，依 v0.1 §4 不計。原列 6 處 **0 處**。
- 最新 main `db380a40`：同 `36f4717c`。
- 本機：同上，只有 12 號兩行引用。
- 判定：DEF-004 **main 已全數修正**（12 號狀態欄的「待 R2 進 main」已滿足；12 號本身是否改寫狀態欄，留待 12 號下次改版處理，不影響本結案）。

## 5. 不能直接查證或屬推定的項目

| 項 | 說明 |
|---|---|
| 合併方式 | GitHub API 不直接提供；依 merge commit 有兩個 parent 與預設訊息推定為一般合併（merge commit） |
| admin bypass | reviewDecision REVIEW_REQUIRED 且無任何 review 卻已合併，推定以管理者權限略過；品保未查 ruleset／稽核紀錄（本次只能唯讀，不查 ruleset 設定變更） |
| 使用者簽核 G4 整體關卡 | agent 轉述，品保未直接查證；本結案不以此為前提 |
| 本文 v0.2、README v0.5 | 本次修改後的新版尚未進 main，需下一個純文件 PR；開 PR／推送須使用者本人同意，本文不構成核准 |

## 6. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 品保部 | PR #27（merge 2425129b）查證；採 19 號 §6 方案 B；20 檔一致、7 檔舊版、9 證據檔缺；G4-QA-R2 有條件（未關閉）；DEF-004 main 未全數修正 |
| v0.2 | 2026-10-27 | 品保部 | PR #28（merge 36f4717c，2026-10-27 12:01:01 UTC+8，一般合併）唯讀查證；v0.1 原列 16 缺漏檔＋13 號＋20 號共 18 檔 sha256 全一致，07-testing 全目錄 37／37 一致；main 原列 6 處簡體字已為 0；**G4-QA-R2 結案**；DEF-004 main 全數修正；新增自我參照說明（以修改前的 README v0.4、本文 v0.1 為比對基準） |
