# 審查紀錄：PR #31（CR-006，Go 1.27.3）新 head 複審（recheck-02）

- 審查部（獨立關卡，只審不改）；2026-11-02 16:34～16:45（UTC+8）
- 對象：head `474756381bd170beb3a5d723b63f1b8ba48b2add`；base main `164b62c956a9ffcb1b206b7fd2b21cd47bd93661`；run 41403811744
- 依據：change-management、部門規範「七、審查部」、stage-gate-checklist G4；前輪 review-PR31-2026-11-02.md（d2b76da8）、review-docs-cr006-followups-2026-11-02.md；CR-006 v0.4（818f48eb，重算相符）
- 前置結論：安全部 ☑ 36 號 v0.1.1（2936baec…，§8 綁新 head）；品保部 ☑ qa-evidence-pr31-47475638/results.md（SUMS 1af3cf98…，10(b)(c) PASS 綁新 head）
- 全程未改專案檔、未留言、未推送、未 rerun；sha256 全部自算

## 結論
**審查內容通過；必改：無（PR31-R 清單為空）。品保 10(b)(c) 已完成並綁新 head，品保前提滿足。**
**可請使用者核准審查部留「可合併」（綁 head 47475638…2add、base 164b62c9…3661、run 41403811744）。**
附帶前提（不需改檔，但須使用者明示）：CR v0.4 §9.1 S1「單一 commit」與 QA-CR006-14 原文，因 D-10 一般 merge 同步而多一個 merge commit；品保 results.md 第 36、43 行仍寫「請拍板」（非品保自行判定可接受），故核准句須一併寫明接受此偏離（見 PR31-S1）。合併由使用者本人執行，最遲 2026-11-20；head 一變即失效。

## 1. 實測核對
| 項目 | 實測 | 宣稱 |
|---|---|---|
| head | ls-remote refs/pull/31/head＝refs/heads/ci/CR-006-go-1.27.3＝474756381bd170beb3a5d723b63f1b8ba48b2add；API head 同；open、未合併、2 commits、2 files、0 留言、0 review、mergeable_state blocked（登入視角；安全部未登入記為 clean） | 相符 |
| base | main＝API base＝merge-base＝164b62c956a9ffcb1b206b7fd2b21cd47bd93661 | 相符 |
| #32 | merged=true，2026-11-02 14:47:43（UTC+8），merged_by 444112029012，merge_commit 164b62c9（parents ca09137a＋fb270d90） | 相符 |
| main 文件（CR6-S3） | 35 號 1000dd9f151f1c33f1c7a2e7b5a4480e093ea1eee1c136a2ff495311b583f38e、34 號 ed210033290835d2693e1a57c221514d1a8986e133796d1a29e68379a782e072、36 號 b7444593cb2e8f5cafe5f0f1836c96dec15585eee1570e7facafb23bb8a5f3ec、README 5dc17e080da695d892f72bfe440f7ffc4ea6383ddc1729cbdc97db0288e4d94c、decision-log 9bf6523b…a5e5 | **CR6-S3 滿足** |
| diff(base,head) | 只有 .github/workflows/ci-devsecops.yml（4/4）、go.mod（1/1） | 相符 |
| 兩檔 sha256（head） | workflow 06833a314722d4bc7b436533a6268b2be29d5c56d028eeb1bfd20f9151def245；go.mod 02842cf5c5b0a74916d55ffa85fcb05db9afc800e4609e8b1c5c35e045f28ebf＝D-20 | 相符 |
| merge commit | parents 826d5d34＋164b62c9；`diff(826d5d34,head)` 與 `diff(ca09137a,164b62c9)` 逐位元組相同（docs 5 檔）；`--remerge-diff` 0 bytes（無衝突解決內容） | 無夾帶 |
| run 41403811744 | pull_request、head_sha＝新 head、success、attempt 1、建立 15:03:25（UTC+8）；jobs secrets-gitleaks 125141230565、sast-gosec 125141230822、sca-govulncheck 125141230836、unit-test 125141230942 皆 success；head check-runs 四項 success；該 head 只有此一 run | 相符 |
| 自取 log（gh api --allow-escape-sequences） | 三個 Go job：go1.27.3、GOROOT 1.27.3、GOVERSION go1.27.3；1.27.2 字樣 0；switching/downloading 0；gosec h1 OK（274 行）、57/7485/0/0；govulncheck「No vulnerabilities found.」、GO-2026-6809 0；L2 key `setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d42…a025`（三 job 命中快取，只記錄）；ENG-011 `--- PASS` | 相符 |
| L7（自取 artifact go-test-json 13122343667） | 181 pass／0 fail／0 skip、TestCT_ pass 14；與品保 art-pr 逐位元組相同 | 成立 |
| 副本比對 | 四份 log 與維運 evidence-pr31-47475638/job-*.log、品保 api-job-*.log 逐位元組相同（品保 job-*.log 為 gh run view 格式，本就不同） | 相符 |

## 2. 證據 sha256 與 -c
| 證據 | sha256（自算） | -c |
|---|---|---|
| 維運 evidence-pr31-47475638/SHA256SUMS | 45cff365e74ec3d6ba1919362c3ecbfa698132b2f42f316b6d3da0f5aa377a6b | 6/6 OK |
| 品保 qa-evidence-pr31-47475638/SHA256SUMS（新版） | 1af3cf98cc9e852cf5f8fc2e82b3a729259380762532b4da515377b0dfe61dab（取代 78c55146…037d；舊版 26 列漏 30～35 號檔） | 32/32 OK，目錄所有檔皆列入 |
| 安全部 M5 cr-006-m5/474756381bd1…/SHA256SUMS | 43d33bb47bed52571dff61b56776787ff4e838a8634da79565aa533122a659bd | 14/14 OK |
| 36 號 v0.1.1（本機 projects/short-url） | 2936baec4ff6d20474a4b31a30e9820e8a3a1177de17544518dd20179d58d598；檔頭與 §8 綁 474756381bd1…2add、run 41403811744 | — |
三份證據與 36 號 v0.1.1 都綁同一 head／run。

## 3. 逐項
- **L1～L7**：全成立。維運只列 L1、L3～L6，原因是 Grok 16:27 給維運的指示只寫「重做 L1/L3/L4/L5」（維運聊天紀錄）；L2 由品保 QA-CR006-2 判 PASS、審查部自取 log 確認；L7 依 CR 對照表由品保（QA-CR006-7/8）負責，審查部自取 artifact 確認。維運取證缺 L2 擷取 → PR31-S2。
- **QA-CR006-1～15**：1～9、10(a)(b)(c)、11、12 PASS；14 PASS 但標「請拍板」（PR31-S1）；13、15 及合併後部分待合併。
- **35 號 S1～S7**：S1 的「兩檔、同一 commit、blob 相符」成立；S2 對新 head 成立；S3～S7 屬合併後。
- **D-10**：以一般 merge 同步、未 force（826d5d34 仍為祖先）。

## 4. 使用者原話（ReadTranscript）
| 對話 | id | 時間（UTC+8） | 原文 |
|---|---|---|---|
| Grok 67cef426 | t365s1390（widget） | 14:45 | 「同意審查部在 PR #32 留「可合併」」 |
| Grok 67cef426 | t368u | 14:48 | 「合併了」（#32 由使用者帳號 444112029012 於 14:47:43 合併） |
| 維運 d7e93bb8 | t108u | 14:55 | 「核准以一般 merge 同步 PR #31：推送 ci/CR-006-go-1.27.3 新 head 98c6f679（merge 826d5d34 與 main 164b62c9，不 force）」——98c6f679 未推送（本機遺失） |
| 維運 d7e93bb8 | t109u | 15:03 | 「核准以一般 merge 同步 PR #31：推送 ci/CR-006-go-1.27.3 新 head 47475638（merge 826d5d34 與 main 164b62c9，不 force）」→ 推送後 run 15:03:25 建立 |
註：#32 沒有「核准合併」的 widget，但合併由使用者本人帳號執行並自述（t368u），符合「合併由使用者本人執行」。

## 5. 必改清單
無。

## 6. S 項（不阻擋）
| 編號 | 條件 | 問題 | 建議內容 | 負責角色 |
|---|---|---|---|---|
| PR31-S1 | 一致性（CR v0.4 §9.1 S1、QA-CR006-14） | PR 為 2 commit（826d5d34＋merge 47475638），與「單一 commit」原文不符；品保 results.md 仍寫「請拍板」，不是品保自行接受；t109u 只核准推送 | 「可合併」核准句一併寫明接受此偏離；CR 下一版把判準改為「淨差異 2 檔＋merge 無額外內容」 | 總協調（呈報）；維運部；品保部 |
| PR31-S2 | 完整性（§9.4 ⑦） | 維運證據只列 L1、L3～L6，缺 L2 快取 key 擷取（源自 Grok 指示漏列） | 維運補 L2 擷取並更新 SUMS（不影響 head） | 維運部；總協調 |
| PR31-S3 | 一致性（版本慣例） | 36 號 v0.1.1 就地覆寫 `-v0.1.md`：同一路徑 main（#32 已入庫）為 b7444593…，本機為 2936baec…，兩者不一致，且本機無 v0.1 原文 | 改用新檔名（如 `-v0.1.1.md` 或 37 號），以文件 PR 入庫；勿以同名覆蓋 repo 的 b7444593 | 安全部 |
| PR31-S4 | 證據格式 | 安全部 E/01-m5-1-head.txt 未記 ls-remote／curl exit code | 補記 exit code | 安全部 |
| PR31-S5 | SCR-060～063 | 仍開放（未登入 log 403 等）；本輪審查部已自取 log 補第二來源 | 照安全部追蹤 | 安全部 |
| PR31-S6 | 證據入庫 | 維運／品保／安全部新 head 證據與 06-security README（本機 7edd4194）未入 repo | 合併後 §9.4 ⑭ 文件 PR 一併入庫 | 總協調 |
| PR31-S7 | 合併方式 | ruleset 需 approving review＋last push approval，PR blocked | 使用者本人 bypass 合併並核對 head | 使用者 |
| 沿用 | CR6-S1（README 第 14 行仍 1.27.2）、S2、S4、S6～S9；PR31-S2（D-20 的 7c4e4e2b 未向使用者展示）；QA 13、15 合併後 | 照原定 | 見前次紀錄 |

## 7. 與 Grok 宣稱不符
1. 「品保 SUMS 78c55146」已作廢：該版 26 列、漏 30～35 號檔、10(b)(c) 待定；現行 1af3cf98。
2. 「品保判定 14 可接受／S12 不阻擋」：品保書面（results.md 36、43 行）實為「請使用者拍板」。
3. 「維運 L1、L3～L6 PASS」正確，但 L2、L7 未列，是因 Grok 指示只要求 L1/L3/L4/L5；L2/L7 實際成立（品保＋審查部自取）。
4. 安全部 M5-1 記 mergeable_state clean（未登入），登入 API 為 blocked；不影響判定。
其餘（head、base、#32 合併、run、四綠、45cff365、43d33bb4、2936baec、兩檔 sha）相符。

— 審查部　2026-11-02（UTC+8）
