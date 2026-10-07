---
文件：09-cicd-env 章節索引
版本：v0.1.5
狀態：審查中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05、SEC-007、SEC-009～SEC-012、SEC-017
專案代號：SHORTURL
---

# 09 CI／環境／機密｜章節索引

> G3 維運必交物。手冊路徑依 `project-handbook-template` 之 **`docs/09-cicd-env/`**（非 08-ops）。  
> 資安阻擋語意以安全部 [`../06-security/03-ci-security-gates-v0.3.md`](../06-security/03-ci-security-gates-v0.3.md) 為準（現行，D-07；v0.2 已由 v0.3 取代），**不得放寬**。

## 本章文件

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-ci-pipeline-v0.2.md](./01-ci-pipeline-v0.2.md) | v0.2 | 核准（D-07、D-08） | **現行** CI 流水線設計：Go 1.27.x／gosec v2.29.0／govulncheck v1.8.0（CR-001）、dependency review fail-closed、合併路徑 A（Rulesets） |
| [01-ci-pipeline-v0.1.md](./01-ci-pipeline-v0.1.md) | v0.1 | 已由 v0.2 取代 | 歷史版本：CI 流水線設計、工具鎖定、分支保護、啟用檢查清單 |
| [02-environments-v0.1.md](./02-environments-v0.1.md) | v0.1 | 審查中 | local／staging／production 區隔、限流預設、SQLite、託管選項 |
| [03-secrets-management-v0.1.md](./03-secrets-management-v0.1.md) | v0.1 | 審查中 | 機密不進版控、存放、輪替、allowlist、gitignore |
| [04-g3r1-landing-checklist-v0.2.md](./04-g3r1-landing-checklist-v0.2.md) | v0.2 | 完成 | **現行** G3-R1：workflow／分支保護落地檢查與證據（classic → Rulesets，D-08） |
| [04-g3r1-landing-checklist-v0.1.md](./04-g3r1-landing-checklist-v0.1.md) | v0.1 | 已由 v0.2 取代 | 歷史版本：G3-R1 落地檢查 |
| [evidence/g3r1-landing-evidence-2026-10-07.md](./evidence/g3r1-landing-evidence-2026-10-07.md) | v0.2 | 完成 | G3-R1 落地證據（Actions 四 checks、classic 保護時點與改用 Rulesets 說明） |
| [05-pr2-ci-change-review-v0.1.md](./05-pr2-ci-change-review-v0.1.md) | v0.1 | 會簽完成 | PR #2 CI 變更的維運意見（同意 (a)(b)） |
| [06-rulesets-setup-guide-v0.1.md](./06-rulesets-setup-guide-v0.1.md) | v0.1 | 已完成 | main 分支 Rulesets 設定步驟（D-08 方案 A） |
| [evidence/g3-addendum-pr2-toolchain-2026-10-07.md](./evidence/g3-addendum-pr2-toolchain-2026-10-07.md) | v0.2 | 草稿待審查 | G3 證據補註：PR #2 工具鏈升版的影響、dependency review 恢復 fail-closed |
| [evidence/rulesets-snapshot-2026-10-07.md](./evidence/rulesets-snapshot-2026-10-07.md) | v0.3 | 草稿待審查 | Rulesets API 快照與截圖證據（13 項 PASS＋2 項已記錄） |
| [evidence/ruleset-strict-true-d10-2026-10-07.md](./evidence/ruleset-strict-true-d10-2026-10-07.md) | v0.1 | 完成 | 規則一 strict=true（D-10）生效取證 |
| [examples/ci-devsecops.yml.example](./examples/ci-devsecops.yml.example) | 範例 | 參考 | 範例 workflow（action 版本為範例；實作前釘 SHA） |

## 與安全部對齊摘要（已選定，勿改名）

| 掃描 | Job 名稱（required checks） | 工具（零預算） | 阻擋門檻 |
|---|---|---|---|
| 機密 | `secrets-gitleaks` | gitleaks（＋建議 GitHub secret scanning／push protection） | 疑似真實機密 **一律擋**（SEC-011） |
| SAST | `sast-gosec` | gosec（選：CodeQL language=go → job `sast-codeql`） | **High／Critical 擋**（SEC-009） |
| SCA | `sca-govulncheck` | govulncheck（PR 可＋dependency-review-action） | **Critical 擋**；**High 有修補逾 14 日擋**（SEC-010） |
| 單元測試 | `unit-test` | `go test`／`go vet` | 失敗擋合併 |

一致性簽核：**通過**（2026-10-07）。維運不修改 `docs/06-security/`；實作層（workflow／分支保護／擋測）另排程。

## 文件狀態 vs 實作狀態

| 項目 | G3 現況 |
|---|---|
| 維運設計文件 | **就緒／審查中** |
| `.github/workflows` 與分支保護 | **已落地**：workflow 已在 main；分支保護採 Rulesets 方案 A（D-08），證據見 [Rulesets 快照](./evidence/rulesets-snapshot-2026-10-07.md)（G3 原紀錄：待 repo 就緒後由維運啟用） |
| 正式環境 HTTPS／託管選定 | **待 G5 前選定**（見 02） |

## 交接單：維運 G3 → 總協調／安全／研發／審查

- **交出角色 → 接收角色**：維運部 → 總協調（並請安全部一致性簽核；研發知悉環境變數與限流；審查部引用本索引與三份設計稿）
- **日期**：2026-10-07
- **交付物清單**：本索引、01 CI 流水線、02 環境、03 機密管理、範例 workflow（可選）
- **版本**：上表皆 v0.1；狀態：審查中
- **對應需求編號**：NFR-004、SI-05、SEC-007、SEC-009～SEC-012、SEC-017
- **已知問題／限制**：
  - 實際 CI／分支保護尚未啟用（文件就緒／實作待 repo）
  - 零預算託管與 production HTTPS 待 G5 前選定
  - MVP 可能尚無雲端 API 金鑰；執行期機密以「有則集中保管、無則文件標明」為準
- **需要下游注意**：
  - 安全部：工具／job 名／門檻已與 [`03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md) 對齊；**一致性簽核通過**（2026-10-07；G3 時點紀錄，現行門檻為 [`03-ci-security-gates-v0.3.md`](../06-security/03-ci-security-gates-v0.3.md)）；實作落地（workflow＋分支保護＋擋測）仍待 repo 就緒
  - 研發：`.env.example`、`.gitignore`、`go.sum` 鎖檔；短碼 CSPRNG／base62 長度 8（D-04）；限流讀環境變數
  - 審查：G3 可標「文件就緒／實作待 repo」若章程允許；啟用後再抽查 required checks
- **待決問題**：無（工具與 job 名已與安全部暫定對齊）；託管選型非本閘阻擋項

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3 初版：CI／環境／機密設計與交接單 |
| v0.1.1 | 2026-10-07 | 維運部 | 一致性簽核：門檻引用改指 `03-ci-security-gates-v0.2.md`；job／阻擋語意未放寬 |
| v0.1.2 | 2026-10-07 | 維運部 | 引用同步：開頭 CI 資安門檻引用由 `03-ci-security-gates-v0.2.md` 改指現行 `03-ci-security-gates-v0.3.md`（v0.2 已由 v0.3 取代，D-07）；G3 交接單之 v0.2 簽核紀錄保留並註記現行為 v0.3；frontmatter 版本對齊修訂紀錄 |
| v0.1.3 | 2026-10-07 | 維運部 | 索引補文件 PR B 新增檔：01 v0.2（現行；v0.1 標為已取代）、05、06、evidence 下 G3 補註與 Rulesets 快照；G3 現況表「workflow 與分支保護」更新為已落地（Rulesets 方案 A），保留原紀錄 |
| v0.1.4 | 2026-10-07 | 維運部 | 索引補文件 PR A 新增檔：04 v0.2（現行；v0.1 標為已取代）、evidence 下 G3-R1 落地證據 |
| v0.1.5 | 2026-10-07 | 維運部 | 索引補 D-10 strict=true 生效取證 |
