---
文件：CI 流水線設計與分支保護
版本：v0.1
狀態：審查中
負責角色：維運部
最後更新：2026-10-07
對應需求：NFR-004、SI-05、SEC-009、SEC-010、SEC-011
專案代號：SHORTURL
---

# CI 流水線設計與分支保護（v0.1）

> **本文件為設計稿**。實際 `.github/workflows` 與 GitHub 分支保護規則，於 **repo 就緒後由維運啟用**。  
> G3 審查可標「**文件就緒／實作待 repo**」（若章程允許）；啟用後依 §7 檢查清單再驗。  
> 阻擋語意對齊安全部 [`../06-security/03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md)，**不得放寬**。

---

## 1. 適用範圍

| 項目 | 說明 |
|---|---|
| 技術堆疊 | Go + chi + SQLite（D-04）；短碼 base62 長度 8、CSPRNG |
| 觸發 | 對 `main` 的 Pull Request；建議另於 `push` 至 `main` 再跑（合併後確認） |
| 預算 | 零現金；僅使用免費／開源工具與 GitHub 免費額度 |
| 不適用 | 本機未推送實驗分支（仍建議開發者自跑同等指令） |

---

## 2. 工具鎖定建議（零預算）

| 掃描 | 對應 SEC | 選定工具 | 授權／成本 | 備註 |
|---|---|---|---|---|
| 機密 | SEC-011 | **gitleaks**（CLI 或 gitleaks-action） | 開源／個人與公開 repo 可用 | **另建議**啟用 GitHub **secret scanning**＋**push protection**（若 repo 可見性允許） |
| SAST | SEC-009 | **gosec** | 開源 | Go 專用；High／Critical 擋 |
| SAST（選） | SEC-009 | **CodeQL**（`languages: go`） | GitHub 免費（公開／符合條件之私有） | Job 名固定為 `sast-codeql`；非 required 亦可先加 |
| SCA | SEC-010 | **govulncheck** | 官方 Go 工具、免費 | 模組漏洞；見 §3.3 實作策略 |
| SCA（PR 輔助） | SEC-010 | **dependency-review-action** | GitHub 免費（Dependency graph） | 建議 `fail-on-severity: critical`；與 govulncheck 併用於 PR |
| 單元／靜態 | — | **`go test`／`go vet`** | 標準工具鏈 | 必跑；失敗擋合併 |
| SBOM | SEC-010 輔助 | **syft** 產 SPDX（建議）或至少確保 **`go.sum` 存在** | 開源 | 鎖定檔必存在；SBOM 產物可上傳 artifact |
| 容器映像 | — | **Trivy image** | 開源 | **有 Dockerfile 時啟用** |
| IaC／組態 | — | Checkov 或 Trivy misconfig | 開源 | **暫不適用**（尚無 Terraform）；**有 Dockerfile 時**可掃 misconfig |

> Action 版本字串於範例檔僅供參考；**實作前須釘定 commit SHA**（見 [`examples/ci-devsecops.yml.example`](./examples/ci-devsecops.yml.example)）。

---

## 3. Job 名稱（固定；供分支保護 required checks）

下列名稱 **不得更名**（已與安全部暫定對齊）：

| Job 名稱 | 必選 | 內容摘要 | 失敗語意 |
|---|---|---|---|
| `secrets-gitleaks` | ✅ | gitleaks 掃工作區／歷史政策依實作選定（至少掃 PR diff／全 repo） | 疑似真實機密 → **exit ≠ 0，擋合併** |
| `sast-gosec` | ✅ | `gosec ./...`；依報告嚴重度判定 | **High／Critical** → 擋 |
| `sca-govulncheck` | ✅ | govulncheck；可同 job 或同 workflow 含 dependency-review | 見 §3.3 |
| `unit-test` | ✅ | `go vet ./...` ＋ `go test ./...`（建議含 race 視 CI 資源） | 測試或 vet 失敗 → 擋 |
| `sast-codeql` | 選 | CodeQL analyze（go） | High／Critical 等級發現依 CodeQL／分支規則擋（若列為 required） |

### 3.1 機密（`secrets-gitleaks`）— SEC-011

- 偵測疑似**真實**機密（雲端金鑰、私鑰、密碼、權杖）→ **一律擋**。
- 明顯測試用假值僅得經 **allowlist**（見 [`03-secrets-management-v0.1.md`](./03-secrets-management-v0.1.md)）且須安全部審核。
- **禁止**無紀錄關閉整條規則；單一路徑例外須安全部＋書面紀錄（對齊門檻文件 §2／§5）。
- 建議並行：GitHub secret scanning ＋ push protection。

### 3.2 SAST（`sast-gosec`）— SEC-009

- **High／Critical** → 流水線失敗並阻擋合併。
- Medium／Low：建立缺陷單追蹤，不強制擋（G5 前應清或接受）。
- 例外：須書面誤判理由、替代控制、核准人（安全部；高風險另呈使用者）。**不得**為「先合併」關閉規則。

### 3.3 SCA（`sca-govulncheck`）— SEC-010

對齊安全部：

| 條件 | 行為 |
|---|---|
| 已知漏洞 **Critical** | **擋** |
| **High** 且有可用修補版本，自修補發布起算 **超過 14 日**仍未升級 | **擋** |
| Medium／Low | 追蹤；不強制擋 |
| 鎖檔／SBOM | **`go.sum` 必須存在**；建議產 SPDX |

**實作策略（擇一或併用，語意不得放寬）：**

1. **主路徑（建議）**：`govulncheck` 發現 **HIGH／CRITICAL**，且漏洞條目標示有 fix／已修補版本 → **fail**（涵蓋 Critical 一律擋，以及 High 有修補之阻擋；14 日期限以「合併當下仍未升級至含修補版本」實作，CI 不放行未升級之 High+fix）。
2. **輔助**：PR 上 `dependency-review-action`，`fail-on-severity: critical`（Critical 擋）。
3. **替代輔助（可選）**：Trivy fs，severity `CRITICAL,HIGH`，並啟用 **ignore-unfixed**（僅擋有修補者）——用於補強 High 逾修補可用仍未升之情境；若採此路徑須在啟用時寫入 workflow 註解並請安全部知悉。

修補時限政策（Critical 7 日／High 14 日等）見安全部門檻 §3；CI 阻擋為合併閘，不等同免除追蹤時限。

### 3.4 單元測試（`unit-test`）

- 必跑：`go vet ./...`、`go test ./...`。
- 建議：模組下載使用 `go.sum` 驗證；CI 快取 Go modules。
- 覆蓋率門檻非本文件強制（品保另訂）；本 job 以「通過／失敗」為合併閘。

### 3.5 SBOM／容器／IaC

| 項目 | G3 設計 |
|---|---|
| SBOM | 建議於建置 job 以 syft 產 SPDX 並上傳 artifact；**至少**版控含 `go.mod`／`go.sum` |
| Dockerfile | **有 Dockerfile 時**新增 job（例如 `image-trivy`）掃映像；Critical／High 策略對齊 SCA 精神並與安全部確認後列入 required |
| IaC | 尚無 Terraform → **暫不適用**；有 Dockerfile 時可加 Checkov 或 Trivy misconfig（非 G3 必列 required） |

---

## 4. 流水線階段（邏輯順序）

```text
checkout
  → secrets-gitleaks          （可並行）
  → sast-gosec                （可並行）
  → sca-govulncheck           （可並行；含或不含 dependency-review）
  → unit-test                 （可並行；需 Go toolchain）
  → （選）sast-codeql
  → （建議）SBOM 產製
  → （有 Dockerfile 時）映像建置＋Trivy
```

任一 **required** job 失敗 → PR 不可合併。

---

## 5. 分支保護（`main`）

| 規則 | 設定 |
|---|---|
| 直接 push | **禁止**（僅能經 PR） |
| PR 核准 | 至少 **1** 人核准 |
| Required status checks | **`secrets-gitleaks`**、**`sast-gosec`**、**`sca-govulncheck`**、**`unit-test`**；（選）`sast-codeql` |
| Force push | **禁止** |
| 管理員 | **同樣適用**（Include administrators） |
| 刪除分支 | 建議限制；合併後可刪 feature 分支 |
| 線性歷史 | 建議啟用（依團隊習慣） |

### 5.1 CODEOWNERS（建議路徑）

於 repo 根目錄 `.github/CODEOWNERS`（擁有者帳號於實作時填入）：

```text
# 範例；實作時替換為實際 GitHub 使用者或團隊
*                              @shorturl-maintainers
/.github/                      @shorturl-maintainers
/docs/06-security/             @shorturl-security
/docs/09-cicd-env/             @shorturl-ops
/docs/08-engineering/          @shorturl-dev
```

> 實際擁有者字串待 repo／團隊帳號就緒後由維運填入；本段為路徑建議。

---

## 6. 與安全部「與維運對齊待辦」對照（逐項）

來源：[`../06-security/03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md) §4。

| 待辦 | 維運回應 | 狀態 |
|---|---|---|
| 選定具體 SAST／SCA／secrets 工具與版本鎖定 | SAST=`gosec`（選 CodeQL）；SCA=`govulncheck`（＋dependency-review）；secrets=`gitleaks`。版本於啟用時釘 SHA／release tag 並記入 workflow | **文件已選定**；版本釘選待啟用 |
| 流水線 job 名稱與失敗退出碼寫入維運 CI 文件 | Job：`secrets-gitleaks`／`sast-gosec`／`sca-govulncheck`／`unit-test`；（選）`sast-codeql`。失敗＝非零退出碼／GitHub check failure | **已寫入本文** |
| 分支保護：上述 job 列為 required checks | 見 §5 | **設計已寫**；repo 就緒後啟用 |
| 例外／allowlist 存放位置與審查流程 | 見 [`03-secrets-management-v0.1.md`](./03-secrets-management-v0.1.md)；誤判／例外**不得**無紀錄關閉規則，程序指向安全部門檻文件 | **已寫入** |
| 與門檻表做一次一致性簽核 | 門檻 v0.2 與本文件 job／阻擋語意一致；維運回覆「一致性簽核通過」（2026-10-07） | **通過（書面）** |

---

## 7. 啟用檢查清單（repo 就緒後維運執行）

- [ ] 建立 `.github/workflows/`（可依 [`examples/ci-devsecops.yml.example`](./examples/ci-devsecops.yml.example) 改寫）
- [ ] 所有第三方 action **釘定 commit SHA**（勿只留浮動 major tag）
- [ ] 確認四個 required job 名稱完全一致：`secrets-gitleaks`、`sast-gosec`、`sca-govulncheck`、`unit-test`
- [ ] `main` 分支保護：禁直接 push、至少 1 核准、required checks、禁 force push、管理員適用
- [ ] （建議）啟用 secret scanning ＋ push protection
- [ ] （建議）新增 `.github/CODEOWNERS`
- [ ] 版控含 `go.mod`／`go.sum`；`.env` 已於 `.gitignore`
- [ ] 以含假值之測試 PR 驗證：gitleaks 對「真實樣式」會失敗；gosec／govulncheck／test 綠燈路徑可合併
- [x] 與安全部門檻 v0.2 一致性簽核通過（書面；2026-10-07）
- [ ] 抽查：故意引入 High SAST 或 Critical SCA 時 PR 被擋

---

## 8. 例外／誤判

- **不得**無紀錄關閉掃描規則或整條 job。
- 單一發現之誤判或路徑例外：書面理由＋替代控制＋安全部核准（高風險另呈使用者）；指向 [`../06-security/03-ci-security-gates-v0.2.md`](../06-security/03-ci-security-gates-v0.2.md) §2「例外核准」與 §5「禁止事項」。
- 機密掃描原則不准關閉整條規則（SEC-011）。

---

## 9. 修訂紀錄

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-07 | 維運部 | G3 設計：工具與 job 名鎖定、門檻對齊、分支保護與啟用清單 |
| v0.1.1 | 2026-10-07 | 維運部 | 一致性簽核：門檻引用改指 `03-ci-security-gates-v0.2.md`；job／阻擋語意未放寬 |
