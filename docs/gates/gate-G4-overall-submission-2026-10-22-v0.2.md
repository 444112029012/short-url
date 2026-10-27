# SHORTURL G4 整體關卡送審材料 v0.2

- 日期：2026-10-22（時間皆為台北時間 UTC+8）
- Repo：444112029012/short-url；main HEAD `f646beeb25eec7a20b358df0af88450c76a85af5`（gh 核實）
- 依據：`/home/box/agent-data/workflows/stage-gate-checklist/SKILL.md` §G4 PR 合併
- 狀態：v0.2，待審查部複審（v0.1 保留於同目錄，sha256 1ee9355c27ff2489f3e5e2288ad7900153ca13aa74b751da276c6adee5782669）；隨 R3 文件 PR 入庫

## 1. G4 過關條件對照

| # | 條件 | 自評 | 證據 |
|---|---|---|---|
| 1 | PR 標示對應 REQ／SEC 編號 | 符合（標題帶任務編號） | §2 PR 標題：ENG-011／007／021／020／013、DR-G4-01 |
| 2 | 公開方法簽名與已簽核規格一致（差異附 CR） | 本輪 PR 無公開方法變更（#21–#24 為 docs、#25 CI、#26 測試）；流程差異走 CR-003／CR-004 | `gates/review-CR-003-2026-10-22.md`、`gates/review-CR-004-2026-10-22.md` |
| 3 | 單元測試涵蓋正常 I/O、邊界、規格例外且通過 | 符合；#26 契約測試 14 PASS | `07-testing/evidence/pr26-ct-run.txt`（14 筆 PASS，`ok internal/httpapi`）、`07-testing/13-G4-QA-R1-acceptance-PR26-v0.1.md` |
| 4 | SAST／SCA／機密掃描無超門檻發現 | 符合 | §3 run 39737480926 四項 success |
| 5 | 安全審查清單完成、阻擋項已修正 | 符合；開放項僅低／資訊級 | §4、§5 |
| 6 | 介面變更同步 OpenAPI 與文件 | #24 README 431 措辭對齊 API 說明；#26 補契約測試 | `06-security/20-*`、`06-security/27-*` |
| 7 | 符合 commit 規範與 DoD | 符合；#24 以 admin bypass 合併，已於 2026-10-22 16:35（台北時間）由使用者追認（原話見 §9 R4） | `07-testing/08-definition-of-done-v0.1.md`、`09-cicd-env/evidence/cr-004/pr24-merge-record-v0.2.md` |

前置：CI 全綠（§3）；安全審查無未解決阻擋項（§5）。除本 PR 外，目前無開啟中的 PR（`gh pr list --state open` 核實）。

## 2. 本輪已合併 PR 對照表（gh 核實）

| PR | 標題 | 合併時間（台北） | merge commit | head | G4 審查紀錄 |
|---|---|---|---|---|---|
| #21 | docs: G4 第二批安全審查文件與維運文件更新 (ENG-011) | 10-22 08:57 | db51cd53d775c398dc2fbfcd9892f2052ba9bf36 | 3f634f1e | `gates/gate-G4-PR21-2026-10-19.md` |
| #22 | docs: G4-PR20 S16 裁示狀態列與 SCR-041 本機同步 (ENG-007, DR-G4-01) | 10-22 09:08 | 363f1692e56228b39c70c4dfe1da58976a3c9871 | 2bd0e811 | `gates/gate-G4-PR22-2026-10-22.md` |
| #23 | docs: 安全審查 08 v0.1.5 與 06-security README v0.22 (ENG-011) | 10-22 12:09 | 2a0281934aec6aa74d3f1ccd9c6b2635d698d4e5 | c3155b3b | `gates/gate-G4-PR23-2026-10-22.md`、`-recheck-01.md` |
| #24 | docs: README 431 門檻與例外回應措辭對齊 API 說明 (ENG-021) | 10-22 13:59 | 5084d6a95a0ac1d1a3884e990089fa742822c3b1 | 0379b64c | `gates/gate-G4-PR24-2026-10-22.md`、`-recheck-01/02.md` |
| #25 | fix: CI 強化 SCR-032／019／021／004 (ENG-020) | 10-22 14:29 | d0749f57038aa0f56e4fa26b2c2b3c97ceb71265 | e54d9003 | `gates/gate-G4-PR25-2026-10-22.md`、`-recheck-01.md` |
| #26 | test: 補 OpenAPI 契約測試 CT-* (ENG-013) | 10-22 14:59 | f646beeb25eec7a20b358df0af88450c76a85af5 | 49e29d80 | `gates/gate-G4-PR26-2026-10-22.md`、`-recheck-01.md` |

註：#21 的 G4 審查檔日期為 10-19（檔名 `gate-G4-PR21-2026-10-19.md`），合併於 10-22。

## 3. CI 證據

- main push run **39737480926**：event=push，headSha=`f646beeb…`（＝#26 merge commit），conclusion=success；sast-gosec／secrets-gitleaks／unit-test／sca-govulncheck 全 success（gh 核實）。
- 本機全量執行紀錄：`07-testing/evidence/g4-full-run-2026-10-22.txt`
- #26 契約測試：`07-testing/evidence/pr26-ct-run.txt`；另有 `pr26-cover.txt`、`pr26-race.txt`、`pr26-vet.txt`、`pr26-mutation.txt`、`pr26-gh-ci.txt`
- Ruleset 24629854 現況（gh 核實）：`dismiss_stale_reviews_on_push=true`、`require_last_push_approval=true`、`required_approving_review_count=1`

## 4. 各部門原件清單（sha256 為 2026-10-22 本機重算）

### 安全
| 檔案 | sha256 |
|---|---|
| docs/06-security/20-pr-024-scr-046-readme-431-security-confirmation-v0.1.md | b10707d0c28eb2631e85f65047332adeab7541e0ba12ace18466ca79b6c719e9 |
| docs/06-security/21-pr-025-ci-hardening-security-review-v0.1.md | 08fe6783e0c76f7a2fffa8be090e09fef7b1076a16b0f22f2bf065f636b56365 |
| docs/06-security/22-cr3-r8-stale-approval-residual-risk-v0.1.md | ca0d6106e697bd4e1415be8c464016baecd3875ca65147c6e038409c38f3f316 |
| docs/06-security/25-kin-openapi-test-dependency-sca-opinion-v0.1.md | 0975ec992f6efc0e58c3288aad63093d5f4666afbf869987e610e767a99e31e5 |
| docs/06-security/26-cr-004-ruleset-stale-approval-security-opinion-v0.1.md | 8f4f3c62684b6683e8d2b54c449b87a6e44e1d26195829b6cc36e99dcfe39a8b |
| docs/06-security/27-pr-026-contract-tests-security-confirmation-v0.1.md | 56b20589e3346b8906f9854605e6c7bff71977edf12f4ce36d355414f28acf6d |
| docs/06-security/05-residual-risk-register-v0.1.md | 540d72321c4721cc8ad69f6b0aeec88414e40579395a60554f13fd41365c91d0 |

### 品保
| 檔案 | sha256 |
|---|---|
| docs/07-testing/12-gate-G4-qa-precheck-v0.1.md | 85831dc8c12f6b26a1899e55b5f534f7fa11de2e6c1cce1f501420ccdb29f09b |
| docs/07-testing/13-G4-QA-R1-acceptance-PR26-v0.1.md | 2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c |
| docs/07-testing/17-CR-004-signoff-2-2-qa-confirmation-v0.1.md | e34b76f793b4f5e874d00880c602687d7f0d8b331116ac7ecc46758674c66587 |
| docs/07-testing/18-CR-004-AC6-9-qa-verdict-v0.1.md（內文版本 v0.4，判首測整體通過） | 543763dad54eb6dd1bbd3b68afcdfa4c80e0b375cfea912fb24a02ef02d61288 |
| docs/07-testing/evidence/g4-full-run-2026-10-22.txt | 8dd97911aa38c2c4248ad69e72975bdd261f3396011d56788a9f2756e4eeb5d8 |
| docs/07-testing/evidence/pr26-ct-run.txt | 3209871dd317775c79574f3ecc5750375dedf8553bf4cf7d2f007066431fc1bd |

### 研發
| 檔案 | sha256 |
|---|---|
| docs/08-engineering/04-g4-implementation-status-v0.1.md | fbd351d2a8a235590855c5b819f52b358e69f1628d315b099fd9c22ab6375123 |
| docs/08-engineering/05-cr003-approval-rd-v0.1.md | 6d3506375774f4f0da265042a542e5d87623a489d9bae191dd6e5cf26cbbfce2 |
| docs/08-engineering/06-cr004-impact-rd-v0.1.md | cc72797909f13a304d41808a4dc0131d585711a5defbbf8b6430d823c8195b41 |

### 維運（CR-004 證據，目錄 docs/09-cicd-env/evidence/cr-004/）
| 檔案 | sha256 |
|---|---|
| ac1-5-record-v0.2.md | 040fced08f6e09e32060cb7ffdf717913b31b8ed4e8bdb7ae9c2da7ea9905b55 |
| pr24-first-test-record-v0.1.md | 3cbe705e5bbc1c8e5c00a5f959dc9c416bb0b3f096931e99da564be7a6ebf99b |
| pr24-merge-record-v0.2.md | 36cb275c972d769799de5f2dc73398373003a87401a601318f08b26c76fd0700 |
| ruleset-24629854-at-pr24-test.json | 1404dab6a8767b6e6235d7c1f9ccadf8dd309cc97606f15d53c2a3952b7ccefa |
| /workspace/rebuild/cr-004/CR-004-ruleset-stale-approval-v0.3.md（CR 本文，repo 外） | eae83378ea0234b936714803fcf9a93515ef72fd9aa36b5c0891c7401ded7ed4 |

### 審查（docs/gates/）
| 檔案 | sha256 |
|---|---|
| gate-G4-PR21-2026-10-19.md | b5d4d491541178c874538cc0a7272599314ede7eb09bcb955c4eb68d377f8e1a |
| gate-G4-PR22-2026-10-22.md | eaea33ce54d80052bf5bb8fa84e220b1493aa5d59af1b0ff55681235d47ccf8c |
| gate-G4-PR23-2026-10-22.md | 8864c87757b0841718c5dd7a59af82dcafe4bbe275453bbac5242cba046c0d2d |
| gate-G4-PR23-2026-10-22-recheck-01.md | e2a5925e65fd1d1b0517da705ebb86d40607c31044d669bf81f6015a3107aab9 |
| gate-G4-PR24-2026-10-22.md | 21e7c40d91037730ef4f77eb128597eb18b76c3ceedc877514be9b95f88809c9 |
| gate-G4-PR24-2026-10-22-recheck-01.md | 13067c45ac74cf7ef3302a39510174289dbe6898b67e9873d709023f46569ee8 |
| gate-G4-PR24-2026-10-22-recheck-02.md | 9a00e133b1ad44e518054f13c6caed36f78f272c730ba86d7db0ec1ab6e4649c |
| gate-G4-PR25-2026-10-22.md | 660c0ffb3fc45242afdaa404cfb7ebf462dab3fc73cce759cf64db83ab14055b |
| gate-G4-PR25-2026-10-22-recheck-01.md | f83d1ed8ef1d5e672a084c76b8da452005eff5c79e0d5f4c5d38ccc729b2bc19 |
| gate-G4-PR26-2026-10-22.md | 2efc4eeb84b4537c9d3b844026821d91584ef5bf1792cb5ec27bef43d542d7b9 |
| gate-G4-PR26-2026-10-22-recheck-01.md | 1460233ed1baf7914da8b44a44641fc7584b9799df27b0dfeb1c571076ce703d |
| review-CR-003-2026-10-22.md | e3618d13f810c77d5048732536a86e33da0fa21667e3f021cc181cd89a170b91 |
| review-CR-004-2026-10-22.md | 325b9c6269ce6fdeb48ee3e4bb0e41ed7919776be9343bab43205576c0a43f15 |
| impact-CR-004-review-dept-v0.1.md | d53d581db6aafe673ed9db6ac6505452fbb11f95b6b37e132ad4861671623965 |

**未核實**：CR-003 本文原件（本機未找到 `/workspace/rebuild/cr-003`）；ENG-020 企劃書面原件（待辦，尚不存在）。

## 5. 開放風險與發現項

| 編號 | 等級 | 狀態 | 依據 |
|---|---|---|---|
| RR-011（stale approval 殘餘風險） | 中 | 維持 | `06-security/22-*`、`26-*` |
| SCR-042（人工核 head） | — | 未關；單帳號無法實測撤銷，安全意見要求保留人工核 head | `06-security/26-*` §2 |
| SCR-050 | 低 | 開放 | `06-security/21-*`／`27-*`（檔內提及） |
| SCR-051、SCR-052 | 資訊 | 開放 | 同上 |
| SCR-004／019／021 | — | 已關（#25） | `06-security/21-*` |
| SCR-032 | — | 已關，附 G5 前自核條件 | `06-security/21-*` |
| SCR-046 | — | 隨 #24 合併關閉 | `06-security/20-*` |
| 檔 25 kin-openapi（測試相依 SCA） | — | 申請結案（待確認） | `06-security/25-*` |

各 SCR 的開關狀態沿用任務已知事實，本彙整只核實檔案存在與 sha256，未逐條重讀判定文字。

## 6. 未結待辦

1. DEF-005、DEF-006：G5 前處理
2. ENG-020 企劃書面原件：待補
3. ENG-022：重做
4. C6 canary：期限 2026-10-26
5. 各 PR（#21–#26）與 CR-003／004 的非阻擋 S 項
6. 文件 PR：把 `06-security/20`、`09-cicd-env/evidence/cr-004/`、`gates/`、CR-003／004 本文與簽核紀錄等收進 repo（目前僅在本機）
7. SCR-032 的 G5 前自核

## 7. 需使用者拍板事項

1. **G4 整體是否過關**（或附條件過關，條件列 §6）。
2. **追認 #24 admin bypass 合併**（14:20 已口頭確認；是否需要書面簽核入檔）。
3. **RR-011 中風險接受**與 SCR-042 在單帳號下無法實測撤銷的處置（保留人工核 head）。
4. **檔 25 kin-openapi 結案**是否核准。
5. **文件 PR** 何時開、範圍（含 repo 外的 CR-004 本文）。
6. ENG-022 重做與 C6 canary（10/26）排程。

## 8. 核實摘要

- gh 核實相符：#21–#26 皆 MERGED、merge commit 與 head 與已知事實一致；無開啟中 PR；run 39737480926 四項 success 且 headSha＝main HEAD；ruleset 24629854 兩參數皆為 true。
- 本機核實：§4 所列檔案皆存在；#26 契約測試 14 筆 PASS；品保 18 號內文版本 v0.4。
- 未核實：CR-003 本文、ENG-020 企劃原件。（#24 使用者追認已於 16:35 取得，見 §9 R4）
- 07-testing/13 號檔 sha256 於 v0.2 重算為 2627c2da9ecb2d38ce3ffe51833fa18461ff8cae84cf789b5134e7878e5f4d6c（v0.1 §4 已為此值，未出現舊值 e688bd6d…）。

## 9. 退回單 G4-R4／R5 決定紀錄（v0.2 新增）

### R4（2026-10-22 16:35 台北時間，使用者原話，逐字引用）

- 「追認 #24 以 admin bypass 合併」
- 「接受 RR-011 維持中風險」
- 「SCR-042 維持不關」
- 「同意由你開文件 PR 收證據」

已追記於 `99-changelog/decision-log.md`（D-17）。

### R5（D-15）：已確認

- 使用者 2026-10-22 16:34 於審查部聊天（widget t92s384）回覆，原話：「算，以 G3 通過為準補標已簽核」。
- 依此，`08-engineering/01-task-breakdown-v0.1.md` 與 `08-engineering/02-branch-commit-conventions-v0.1.md` 檔頭狀態補標「已簽核（G3 通過，gate-G3-2026-10-07-recheck-01；D-15）」，版本號不變。
- 已追記於 `99-changelog/decision-log.md`（D-15 確認列）與 `99-changelog/CR-003-signoff-record-v0.1.md`（CR-003 R10 結案）。
