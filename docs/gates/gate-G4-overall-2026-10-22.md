# G4（實作／PR 合併）整體關卡審查紀錄

- 日期：2026-10-22（台北時間 UTC+8）｜審查角色：審查部（只審不改；未在 GitHub 留言、未改 repo）
- 送審物：`docs/gates/gate-G4-overall-submission-2026-10-22-v0.1.md`，自算 sha256 `1ee9355c27ff2489f3e5e2288ad7900153ca13aa74b751da276c6adee5782669`（＝宣稱）
- 依據：stage-gate-checklist §G4；部門規範「七、審查部」A–G
- 審查基準：main `f646beeb25eec7a20b358df0af88450c76a85af5`

## 一、結論：**退回＋必改**（G4-R1～R5）

實作面（程式、CI、各 PR 綁定 head）全部核實相符，無程式碼必改。退回原因在**前置專業結論缺漏**、**G4 證據未入庫（可追溯性）**、**兩項使用者決定未經可查核紀錄**。必改都不需要變動 main 程式碼。結論請使用者拍板。

## 二、獨立核實結果

| 項 | 宣稱 | 審查部實查 | 判定 |
|---|---|---|---|
| main HEAD | f646bee… | `git ls-remote`：HEAD＝refs/heads/main＝`f646beeb25eec7a20b358df0af88450c76a85af5` | ✅ |
| 開著的 PR | 無 | `list_pull_requests state=open` → `[]` | ✅ |
| main push run 39737480926 | 四綠 | event=push、head_sha=`f646beeb…`、conclusion=success；jobs：sast-gosec／secrets-gitleaks／unit-test／sca-govulncheck 皆 success（govulncheck 內 Probe／Dependency review 為 PR 限定步驟，push 時 skipped，屬設計）；完成 15:01:49（UTC+8） | ✅ |
| 送審材料 sha256 | 1ee9355c… | 相同 | ✅ |
| §4 列 34 檔＋CR-004 本文 sha256 | 如表 | 35 檔逐一重算，全數相符 | ✅ |
| ① CR-003 本文 | c55ad3ca… | `docs/99-changelog/CR-003-branch-docs-exception-v0.1.md` 與 `/workspace/snap/CR-003-c55ad3ca.md` 皆＝`c55ad3ca7c97497e07c84cfce2503bc74063bdff9de6a7692b3a6604d44f57b2` | ✅（送審材料標「未核實」已可解除） |
| ② #24 merge record | 36cb275c… | ＝`36cb275c972d769799de5f2dc73398373003a87401a601318f08b26c76fd0700`；merge 5084d6a parents＝2a02819（main）＋0379b64c（head），屬實 | 檔案 ✅ |
| ② t257s1017 使用者原話 | 「是，用 admin bypass 合併」14:20 | 審查部可用工具中**無 ReadTranscript**（cursor namespace 查無），無法讀 Grok Bot 67cef426-… 原始紀錄 | **未核實** |
| ③ ENG-020 企劃原件 | 缺 | 確認缺 | 開放（G5 前） |

### 合併 commit 第二 parent ＝「可合併」綁定 head

| PR | 綁定 head（審查紀錄） | merge commit | 第二 parent | 判定 |
|---|---|---|---|---|
| #7 | fc9a7a7（PR7 recheck-01） | 24b8da0 | fc9a7a7f… | ✅ |
| #8 | 0e47abe（PR8 recheck-01） | e9f4396 | 0e47abe1… | ✅ |
| #9 | 6c7104c（PR9） | 63b373d | 6c7104c6… | ✅ |
| #10 | 7c82929（PR10 recheck-01） | f8c4011 | 7c829291… | ✅ |
| #11 | b623461（PR11 recheck-01） | 3d3a9cd | b623461f… | ✅ |
| #12 | c6af336（PR12 recheck-01） | 28fc4f7 | c6af3368… | ✅ |
| #13 | 6674a59（PR13） | 2e54f7e | 6674a59d… | ✅ |
| #15 | 412cefa（PR15 recheck-01） | c341be1 | 412cefaf… | ✅ |
| #16 | 95e7561（PR16） | ed62039 | 95e7561b… | ✅ |
| #17 | e6d9ab1（PR17 recheck-02） | 433140c | e6d9ab12… | ✅ |
| #18 | f547a6e（PR18 recheck-01） | 643bf8a | f547a6eb… | ✅ |
| #19 | 4804370（PR19 recheck-02） | 141fcda | 48043702… | ✅ |
| #20 | 44b099c（PR20 recheck-01） | abbec9e | 44b099c2… | ✅ |
| #21 | 3f634f1（PR21） | db51cd5 | 3f634f1e… | ✅ |
| #22 | 2bd0e81（PR22） | 363f169 | 2bd0e811… | ✅ |
| #23 | c3155b3（PR23 recheck-01） | 2a02819 | c3155b3b… | ✅ |
| #24 | 0379b64（PR24 recheck-02） | 5084d6a | 0379b64c… | ✅ |
| #25 | e54d900（PR25 recheck-01） | d0749f5 | e54d9003… | ✅ |
| #26 | 49e29d8（PR26 recheck-01） | f646bee | 49e29d80… | ✅ |

（#6 為事後補審，PR6 recheck-01 已處置；#14 未合併已關閉。）

## 三、前置專業結論

- **安全部**：只有逐 PR 確認（08、11、12、14、16～21、23、27 號）與 CR／殘餘風險意見（22、24～26），**沒有針對 G4 整體（main f646bee）的書面結論**（「安全審查無未解決阻擋項」由送審方自評，非安全部原件）→ G4-R1。
- **品保部**：12 號為 G4 前置，只盤點 #13～#20、基準 db51cd5，結論「有條件可送審」，明載「R1、R2 關閉前不出具通過」。R1（DEF-001）已關；**G4-QA-R2（矩陣 v0.4 等文件進 main）仍未關**（main 上 `02-traceability-matrix` 仍為 v0.3）。13 號只驗收 #26。無品保對 f646bee 的 G4 整體結論 → G4-R2。

依規範 B 第 1 條，前置不到位應直接退回；因實作面已逐 PR 審過，以下仍留四層勾選供複審縮小範圍。

## 四、四層檢查

### ① 完整性
- [x] 各 PR 描述、CI、逐 PR 審查紀錄齊（#7～#26）
- [ ] 安全部 G4 整體結論 → G4-R1
- [ ] 品保部 G4 整體結論 → G4-R2
- [ ] ENG-020 企劃原件（G5 前，S 項）

### ② 一致性
- [x] 送審 §2／§3／§4 數據與實查相符
- [ ] 送審 §1 第 7 列寫「使用者 14:20 親口確認」，§8 又列為「未核實」；merge record 亦註「維運部未直接見證」→ 自相矛盾，見 G4-R4
- [ ] CR-003 本文以「依 D-15 視為已簽核」書寫，但 review-CR-003 結論為「通過，惟 R10 待使用者確認後生效」，D-15 未在審查部確認 → 見 G4-R5
- [x] 送審 §4 「CR-003 本文未核實」：本次已核實 c55ad3ca（材料下一版應更新）

### ③ 可追溯性
- [x] 任務→PR：PR 標題帶 ENG-xxx／DR 編號，逐 PR 紀錄已核 REQ／SEC
- [x] 公開簽名→G2 方法規格：main 上 internal/ 51 個匯出函式／方法，40 個字面出現在 06-method-specs v0.2＋09 簽核 v0.2.1；其餘 11 個（NewRedirectService、NewURLValidator、NewRouter、NewRouterWithOptions、NewErrorMapper、NewHTTPAPIAdapter、Load、ParseTrustedCIDRs、Status、InvalidURL、WriteHeader）多為建構子／組裝／介面實作，9 個已在逐 PR 審查或架構文件處理；NewRedirectService、NewURLValidator 無任何文件引用 → S2
- [x] 單測覆蓋 I/O 與例外：品保實跑 167 PASS／0 FAIL／0 SKIP、覆蓋 80.9%；CT 14 PASS；殘餘 DEF-005（CT 寫死常數）、DEF-006（redirect／stats 500 無 CT）G5 前
- [ ] **G4 證據未入庫**：main 上 `docs/gates/` 只有 8 檔（G0～G3），全部 G4 審查紀錄、安全 20～27 號、品保 11～18 號與矩陣 v0.4、`09-cicd-env/evidence/cr-004/`、CR-003／CR-004 本文與簽核紀錄皆只在本機快照 → G4-R3

### ④ 目標達成驗證（G4：PR 變更真對應 REQ、公開簽名與規格一致、測試與安全門檻通過）
- [x] PR 對應 REQ／任務：逐 PR 已核
- [x] 簽名一致：無未附 CR 的簽名差異（DR-G4-01 已簽 v0.2.1）
- [x] 測試／安全門檻：main push run 四綠，綁定 head 全部吻合
- [ ] 「安全審查無未解決阻擋項」缺安全部原件證明（G4-R1）；RR-011（中）屬殘餘風險接受，須使用者決定（G4-R4 一併）

### G4 過關條件勾選
1. [x] PR 標示 REQ／SEC　2. [x] 簽名一致（差異附 CR）　3. [x] 單測涵蓋且通過（DEF-005／006 G5 前）　4. [x] SAST／SCA／機密掃描無超門檻　5. [ ] 安全審查清單完成——缺整體結論（R1）　6. [x] 介面變更同步 OpenAPI（#20、#24）；文件入庫見 R3　7. [ ] commit 規範與 DoD——#24 bypass 合併未有可查核授權（R4）；CR-003 例外未生效（R5）；DoD D2 待 R3

## 五、D-15 是否阻擋 G4

**阻擋，但只需使用者一個決定即可解除。** 理由：CR-003 是 docs 分支命名例外，#21～#24 等文件 PR 的條件 7（commit／分支規範）依賴它成立；review-CR-003 已明定「R10 待使用者確認後生效」，D-15 的證據（09:16「這個事情你幫我決定」）在總協調聊天，審查部無法查核。使用者在本次拍板時一併確認 D-15（01／02 視為已簽核基線，或改走升 v1.0）即可關閉 G4-R5，不需重做任何 PR。

## 六、必改清單

| 編號 | 條件 | 問題 | 必改內容 | 負責角色 |
|---|---|---|---|---|
| G4-R1 | 前置；條件 5 | 無安全部對 G4 整體（main f646bee）的書面結論 | 出具 G4 安全結論原件：綁定 main `f646beeb…`、run 39737480926，逐條列 SCR-042、SCR-050／051／052、RR-011、檔 25 狀態，明寫「無未解決阻擋項」或列出阻擋項 | 安全部 |
| G4-R2 | 前置；③可追溯 | 品保僅有 #13～#20 前置（有條件）＋#26 驗收，G4-QA-R2 未關 | 出具品保 G4 整體結論，涵蓋 #13～#26、基準 f646bee（main 重跑結果），並關閉或明確處置 G4-QA-R2 | 品保部 |
| G4-R3 | ③可追溯；條件 6；DoD D2 | G4 全部審查紀錄與證據只在本機，main 不可查 | 開純文件 PR，把 `docs/gates/` G4 紀錄、`06-security/20～27`、`07-testing/11～18`＋矩陣 v0.4＋evidence、`09-cicd-env/evidence/cr-004/`、`99-changelog/` CR-003／004 本文與簽核紀錄、decision-log 收入 repo；送審部 PR 審查（內容須與本紀錄所列 sha256 一致） | 總協調（各部提供原件） |
| G4-R4 | 條件 7；②一致性 | #24 admin bypass 的使用者授權（t257s1017）審查部無法查核；送審 §1 與 §8 說法矛盾；RR-011（中）尚未經使用者接受 | 使用者在審查部聊天或書面親自確認：(a) 追認 #24 admin bypass 合併；(b) 接受 RR-011 與 SCR-042 保留人工核 head。總協調將原話入檔並修正送審 §1 措辭為「待追認」 | 使用者（總協調記錄） |
| G4-R5 | 條件 7；CR-003 R10 | D-15 未經使用者在審查部確認，CR-003 未生效 | 使用者確認 D-15（01／02 視為已簽核）或另擇升 v1.0；總協調更新 decision-log 與 CR-003 簽核紀錄 | 使用者（總協調記錄） |

複審只看 R1～R5 與受影響文件；main 若無新 commit，不需重核 CI 與綁定 head。

## 七、S 項（非必改，不阻擋）

- S1：ENG-020 企劃書面原件 G5 前補（企劃部）
- S2：方法規格表補建構子 NewRedirectService、NewURLValidator 等（或註明建構子不列）（設計部，G5 前）
- S3：DEF-005、DEF-006 G5 前關閉；DEF-002、DEF-003 排期（研發部／品保部）
- S4：SCR-050（低）G5 前處置；SCR-051／052 資訊級追蹤；SCR-032 G5 前自核（安全部）
- S5：SCR-042／D-09 人工核 head 持續執行直到 ruleset 可實測（審查部／維運部）
- S6：檔 25 kin-openapi 結案、ENG-022 重做、C6 canary（期限 2026-10-26）排程由使用者決定（總協調）
- S7：送審材料下一版更新「CR-003 本文已核實 c55ad3ca」；CR-004 本文目前在 `/workspace/rebuild/cr-004/`（eae83378…），與 `docs/99-changelog/` 副本一併以 R3 入庫
- S8：各 PR 遺留 S 項（含 PR21～PR26 紀錄所列）由總協調登入退回項追蹤表

## 八、交使用者拍板

1. 是否同意本次「退回＋必改 R1～R5」。
2. R4：是否追認 #24 admin bypass 合併、接受 RR-011（中）與 SCR-042 保留人工核 head。
3. R5：是否確認 D-15。
4. R3 文件 PR 的開立時間。

本紀錄只新增本檔；未修改任何專案檔、未推送、未在 GitHub 留言。
