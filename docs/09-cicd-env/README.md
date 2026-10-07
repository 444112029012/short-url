---
文件：09-cicd-env 章節索引
版本：v0.1
狀態：審查中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05、SEC-007、SEC-009～SEC-012、SEC-017
專案代號：SHORTURL
---

# 09 CI／環境／機密｜章節索引

> G3 維運必交物。手冊路徑依 `project-handbook-template` 之 **`docs/09-cicd-env/`**（非 08-ops）。  
> 資安阻擋語意以安全部 [`../06-security/03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md) 為準，**不得放寬**。

## 本章文件

| 文件 | 版本 | 狀態 | 說明 |
|---|---|---|---|
| [01-ci-pipeline-v0.1.md](./01-ci-pipeline-v0.1.md) | v0.1 | 審查中 | CI 流水線設計、工具鎖定、分支保護、啟用檢查清單 |
| [02-environments-v0.1.md](./02-environments-v0.1.md) | v0.1 | 審查中 | local／staging／production 區隔、限流預設、SQLite、託管選項 |
| [03-secrets-management-v0.1.md](./03-secrets-management-v0.1.md) | v0.1 | 審查中 | 機密不進版控、存放、輪替、allowlist、gitignore |
| [04-g3r1-landing-checklist-v0.1.md](./04-g3r1-landing-checklist-v0.1.md) | v0.1 | 進行中 | G3-R1：workflow／分支保護落地檢查與證據 |
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
| `.github/workflows` 與分支保護 | **待 repo 就緒後由維運啟用**（見 01 啟用檢查清單） |
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
  - 安全部：工具／job 名／門檻已與 [`03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md) 對齊；**一致性簽核通過**（2026-10-07）；實作落地（workflow＋分支保護＋擋測）仍待 repo 就緒
  - 研發：`.env.example`、`.gitignore`、`go.sum` 鎖檔；短碼 CSPRNG／base62 長度 8（D-04）；限流讀環境變數
  - 審查：G3 可標「文件就緒／實作待 repo」若章程允許；啟用後再抽查 required checks
- **待決問題**：無（工具與 job 名已與安全部暫定對齊）；託管選型非本閘阻擋項

---

## 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3 初版：CI／環境／機密設計與交接單 |
| v0.1.1 | 2026-10-07 | 維運部 | 一致性簽核：門檻引用改指 `03-ci-security-gates-v0.2.md`；job／阻擋語意未放寬 |
