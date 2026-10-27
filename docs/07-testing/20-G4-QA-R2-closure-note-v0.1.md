---
文件：G4-QA-R2 處置紀錄（方案 B：經 PR #27 進 main）
版本：v0.1
狀態：審查中
負責角色：品保部
最後更新：2026-10-27（Asia/Taipei，UTC+8）
依據：`19-G4-overall-qa-conclusion-v0.1.md` §6（sha256 `86606bd3cfb059c52dc1ac9984161d9c8cc879eeeebec9c46795f4d4960b0d70`）；`12-gate-G4-qa-precheck-v0.1.md` §8 G4-QA-R2
比對基準：main `2425129b2384d6a8b2c158691e8e4d3076369164`（PR #27 merge commit）
專案代號：SHORTURL
---

# G4-QA-R2 處置紀錄

## 0. 結論：**有條件（未關閉）**

採 19 號 §6 **方案 B**：07-testing 文件經 PR #27 進 main。PR #27 收錄的 07-testing 檔（含 19 號、矩陣 v0.5）與本機 sha256 **全部一致**；但 07-testing 目錄另有 **7 個文件在 main 上仍為舊版、9 個證據檔未進 main**，含 G4-QA-R2 原列之案例 v0.2、任務對照 v0.3、README v0.4，以及 DEF-004 的簡體字修正。依規則有不一致不得寫成關閉。

關閉條件：以純文件 PR 補入 §3 所列 16 檔（內容須與 §3 本機 sha256 一致），合併後品保重新比對即關閉。開 PR／推送須使用者本人同意；本文不構成核准。

> 「使用者已簽核 G4 整體關卡」：**agent 轉述，品保未直接查證。** 本處置不以此為前提。

## 1. PR #27 查證（gh 唯讀，2026-10-27 08:58～09:02 UTC+8）

| 項 | 回報 | 品保實查 | 判定 |
|---|---|---|---|
| 狀態 | 已合併 | `gh pr view 27`：MERGED，mergedAt 2026-10-27T00:56:42Z（＝08:56:42 UTC+8），mergedBy 444112029012 | ✅ |
| merge commit | 2425129b | `2425129b2384d6a8b2c158691e8e4d3076369164`；parents f646beeb（main）＋a4276fba（head）；`gh api …/commits/main` 同值 | ✅ |
| 分支／標題 | — | `docs/g4-overall-evidence`；「docs: 收錄 G4 整體關卡證據（R3）」；71 檔 | — |
| checks | — | PR 四項 success；main push run 40429569968 success | ✅ |
| 審查狀態 | — | reviewDecision REVIEW_REQUIRED（無核准，推定以 admin bypass 合併；品保未查到 bypass 紀錄） | 記錄 |
| 含 07-testing 檔 | 含 19 號、矩陣 v0.5 | 含 02、11～19、evidence/g4-overall-f646-*（6 檔）、evidence/pr26-ct-run.txt，共 17 檔 | ✅（範圍見 §3） |

## 2. 一致檔（main `2425129b` 與本機 sha256 相同，2026-10-27 09:00 重算）

| 檔案（`docs/07-testing/`） | sha256 |
|---|---|
| 01-requirements-testability-review-v0.1.md | 120ddbf65e649aa196c892dcb88b5be6956adf30a7a59ee9a0ad66cf786724b3 |
| 02-traceability-matrix-v0.1.md（v0.5） | 8a005aefcf24a32fa769c803373e88a3bb411aa71652ff14426feb18907b9932 |
| 03-gate-G1-qa-precheck-v0.1.md | dbabd14abd278b16d87d554868d5c017e0ace08f4b0904d97da99aded24987c1 |
| 08-definition-of-done-v0.1.md | be865c1cd31e6ea7b25e0bf84f09088182bc86acf6efaaf5e2ebd198fd33dd06 |
| 11-s6-eng011-disposition-v0.1.md | 02e2973d4686d2945a31873d8e9acb701974d208209d1188e1aa96c1511217b4 |
| 12-gate-G4-qa-precheck-v0.1.md | 85831dc8c12f6b26a1899e55b5f534f7fa11de2e6c1cce1f501420ccdb29f09b |
| 13-G4-QA-R1-acceptance-PR26-v0.1.md（main＝v0.1 附錄版） | 2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c |
| 14-CR-003-qa-confirmation-v0.1.md | 3fb5d45727a337d9d9901ed7a5399922457d166b135cba10845b408f68b30ae9 |
| 15-CR-004-R3-qa-impact-row-v0.1.md | 4c3684fbfb92ed4463a24cdd4a6d0e14349ce4418bfb7315b440a49f1617ade2 |
| 16-CR-004-R4-qa-opinion-v0.1.md | ad44316f22ae5cf8c5b2266141dc2cc6493e6e249e1dcc753ecb9d9ce737495c |
| 17-CR-004-signoff-2-2-qa-confirmation-v0.1.md | e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587 |
| 18-CR-004-AC6-9-qa-verdict-v0.1.md（v0.4） | 543763dad54eb6dd1bbd3b68afcdfa4c80e0b375cfea912fb24a02ef02d61288 |
| 19-G4-overall-qa-conclusion-v0.1.md | 86606bd3cfb059c52dc1ac9984161d9c8cc879eeeebec9c46795f4d4960b0d70 |
| evidence/g4-overall-f646-cover.txt | 2d4a41869fcc9f106babebcc8eb4bb95f2953bb0d4d5efbe0c0c31397c269c29 |
| evidence/g4-overall-f646-ct.txt | e363d42c357983b5f34a117d508e459345af12b641ef9f020e5b4cd2c486e990 |
| evidence/g4-overall-f646-env.txt | 73138c31bfb06fb5d5699641b93b029d364cd2a65824f0d2d431a583e9eec537 |
| evidence/g4-overall-f646-race.txt | 42f375378befa50f51b9255f994c4bde7e54a00b04804f0d78d86abdca21ad86 |
| evidence/g4-overall-f646-test.json | 71353d9a8d74dea36ca6fb9916b07cfd93bfba0f6f2dafcd50e0867a5bcb89a7 |
| evidence/g4-overall-f646-vet.txt | 00a5f1d995f1697b6f134bcca3e00d503345a5241662ca72a16f8b0ac5768294 |
| evidence/pr26-ct-run.txt | 3209871dd317775c79574f3ecc5750375dedf8553bf4cf7d2f007066431fc1bd |

註：13 號本機已於 2026-10-27 升 v0.1.1（S2，只加修訂紀錄一列），故本機現值與 main 不同；屬比對之後的預期改版，不列為差異。20 號（本文）為新檔，不在比對範圍。

## 3. 差異（**阻擋關閉**）

### 3.1 main 為舊版（PR #27 未收錄）

| 檔案 | main 版本／sha256 | 本機版本／sha256 | 差異行數 |
|---|---|---|---|
| 04-design-testability-review-v0.1.md | v0.1.1／c45c7bab56a8b7a8f18e15870873aca5eaee69480d92fae494d7b7577d6bdb4e | v0.1.1／12002fa436755e189f46769318d65c87e44fef73298460dd31def4df13a5b3ca | 4（DEF-004 修正 2 處） |
| 05-gate-G2-qa-precheck-v0.1.md | v0.1.1／a539f8308c7e2e3e01be64f3f1ab226da7c31ed0f7355b66d36036b0d400faf3 | v0.1.1／b5d018d19a7804edb6d572ac9ea602f2ca9c230446c1b1178ce3b698c09c5b1d | 2（DEF-004） |
| 06-test-strategy-v0.1.md | v0.1／12173f9213dc6377ed9638c5854ebefd669a1a62ed7ebc10118e2e7783857151 | v0.1／409187e193ada518970957f0fc24c613fb2c5906c4ec667ef679dcf60d8e557d | 2（DEF-004） |
| 07-test-cases-v0.1.md | v0.1／61628cc81d9e69d9faa03344e7e2b6301cd30a0e9917d6fffb97d172f0c7df9a | **v0.2**／ef8ed11d6e8324b3d52a5cb6d597b12b13c6e0896730dd7ef458672978f098e5 | 35 |
| 09-dev-task-acceptance-map-v0.1.md | v0.2／905141696dff1909014978796eb3db9870f825a27a67b8aa0bb60d42dac931ee | **v0.3**／f7ca9727ab9c594fbe5c01efe3ed18c2542115d4d7656e2d71fca2ca38568792 | 28 |
| 10-gate-G3-qa-precheck-v0.1.md | v0.1／ad8435b99a102014bb5952f0e14c852a65ceb2e7b7b7c71038a0607b636ecacd | v0.1／d229e23e5ef87d41790288ff65959869d9182767ab3bf993c80561b7df0cb344 | 2（DEF-004） |
| README.md | v0.3／f7fbcec807a901aa66063b537a1942dd090de671392bef54d3493fca50902a5a | **v0.4**／25c1807b8aa8c47a71dd8d635fa07d3331910a5a27b6c68eb2a8a1c40edd030c | 26 |

影響：main 上矩陣 v0.5 標頭與 §7 引用「案例內容版 v0.2」，main 上 07 號仍為 v0.1，追溯在 main 內不自洽。

### 3.2 證據檔未進 main

| 檔案（`evidence/`） | 本機 sha256 | 被引用處 |
|---|---|---|
| g4-full-run-2026-10-22.txt | 8dd97911aa38c2c4248ad69e72975bdd261f3396011d56788a9f2756e4eeb5d8 | 矩陣 §7、12 號 §2 |
| s6-eng011-run-2026-10-22.txt | c381c361bbbe2d851625d126ac951d9168b1098d907d3c94826d2b336413d4a4 | 12 號 §2、11 號 |
| pr26-cover.txt | eb004f4c1541a95f003ee5e07d721692d356c606c5e6980850b60caa333f1ac1 | 13 號 §3 |
| pr26-env.txt | 3231de7959ade2625e9113e67fd4b6df1f84283f22ff93ac0a0315a3adb7ec6c | 13 號 |
| pr26-gh-ci.txt | 592a4386cb87e3e91f067d15416c576c573b86c4714b29275530b838ecef8ac5 | 13 號 §1 |
| pr26-gh-pr.json | 712f3a2a7a371cf4ee84e50decabfb7f4a7b75cccd220c59b80e9f5472e0af4c | 13 號 §1 |
| pr26-mutation.txt | 7b80f35803c72a54ed05761126b13176ce1c9d39b4ee3e3a10e269312e6f091e | 13 號 §4 |
| pr26-race.txt | 5a18e0414f56dcc71f1f7d094462392879f2a1a093b820f3898c15e0d2376cb3 | 13 號 §3 |
| pr26-vet.txt | 19eaf43821a7660ec323a87c8457bf74823beb296c39f5e01aa8a683aa50f061 | 13 號 §3 |

main 上 07-testing 沒有本機以外的多餘檔。

## 4. DEF-004（簡體字修正）是否進 main：**否（部分）**

- 指令：在 main `2425129b` 工作樹執行 `rg` 搜尋 DEF-004 所列兩個簡體字（字樣見 12 號 §7 DEF-004 重現步驟）。
- main 結果：仍有 **6 處**：04 號第 126、141 行；05 號第 33 行；06 號第 114 行；10 號第 130 行；README 第 39 行。只有 02 號（隨 PR #27 矩陣 v0.5）已修正。另 12 號第 99、111 行為缺陷描述中引用字樣，屬記錄用途，不計。
- 本機：上述 6 檔已無簡體字（排除 12 號引用）。
- 判定：DEF-004 維持「本機已修正，**main 未全數修正**」；隨 §3.1 補 PR 一併處理。

## 5. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 品保部 | PR #27（merge 2425129b）查證；採 19 號 §6 方案 B；20 檔一致、7 檔舊版、9 證據檔缺；G4-QA-R2 有條件（未關閉）；DEF-004 main 未全數修正 |
