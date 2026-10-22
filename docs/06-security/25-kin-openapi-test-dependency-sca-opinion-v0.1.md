---
文件：第三方元件 SCA 書面意見：github.com/getkin/kin-openapi（僅測試用，G4-QA-R1 契約測試）
版本：v0.1
狀態：**不同意**（評等：低到中；使用者政策只接受低風險）
負責角色：安全部
作者：安全部
日期：2026-10-22
最後更新：2026-10-22
對應：品保必改 G4-QA-R1；研發部引入申請（經總協調 2026-10-22 09:07 轉達）；審查部 PR #26 S1；書面結論已於 2026-10-22 09:09 送研發部，並副本給總協調
專案代號：SHORTURL
時區：Asia/Taipei（UTC+8）
---

> 以下是 2026-10-22 評估當時的工作紀錄全文，原樣收錄；結論與處置見文末「安全部結論」。

- 對象：short-url main db51cd5，go 1.27.2（工具鏈 go1.27.3）
- 評估日期：2026-10-22（Asia/Taipei）；實驗目錄 /tmp/kinsca.4fsP（拋棄式）

## 結論：不同意（評等：低到中）

## 1. 版本與維護
- 最新穩定版 v0.151.0，proxy.golang.org 時間 2026-10-15T10:00:58Z（台北時間 2026-10-15 18:00），commit a07ec19。
- 近 90 天 commit ≥100（API 單頁上限）；最後 push 2026-10-22 00:46（台北時間）；未封存。
- 未關閉 issue+PR：84（GitHub open_issues_count 含 PR）。
- 貢獻者集中：fenollp 326、reuvenharrison 93、jban332 36，其餘個位數到 23；實際活躍維護者約 1–2 人。

## 2. 授權（全部為寬鬆授權）
kin-openapi MIT；regexp2 MIT；oasdiff/yaml MIT；oasdiff/yaml3、gopkg.in/yaml.v3 MIT+Apache-2.0；go-openapi/jsonpointer、swag/jsonname、testify/v2 Apache-2.0；santhosh-tekuri/jsonschema/v6 Apache-2.0；gorilla/mux BSD-3；golang.org/x/text BSD-3；stretchr/testify MIT；go-spew ISC；go-difflib BSD-3；go-internal BSD-3；check.v1 BSD-2；creack/pty、kr/pretty、kr/text MIT（GitHub API）；pkg/diff BSD-3（GitHub API）。無 copyleft。

## 3. 依賴圖實驗
- go get kin-openapi@v0.151.0 + go mod tidy，新增 22 個 module（見回報清單）。
- 最小 _test.go import openapi3 / openapi3filter，載入 docs/04-api/openapi.yaml 並 Validate：PASS。
- `go list -deps ./cmd/... | grep -c kin` → 0；`go mod why` 只經 internal/contracttest(.test)。
- govulncheck v1.9.0：不加 -test「No vulnerabilities found.」；-test「No vulnerabilities found.」
- `go mod verify` → all modules verified。
- go.sum：v0.151.0 h1:p8NIstL53h620DnlO9R4pR4QRkEzJ3jbw9flb4kc1hg= ；/go.mod h1:1+BHDzstro+P5CKtPy1X4PfofnFgmRe6uvMy9+r9fKY=

## 4. 安全公告歷史
GitHub repo 公告 15 筆，其中 14 筆在 2026 年：critical 1（GHSA-r277-6w6q-xmqw 認證繞過，<0.144.0）、high 7、medium 6；最新 GHSA-g3f9-hcw7-r9w2（2026-10-14 發布，影響 0.95.0–0.149.0，修補 0.150.0，未在 OSV）。
- 正則相關 DoS：GHSA-h4r5-mxv5-8mjc（deepObject 重複編譯正則的資源耗用，<0.146.0）；未查到 ReDoS（正則回溯過深）專屬公告，但依賴 regexp2（回溯型引擎）用於 pattern，理論面存在。
- 外部 $ref SSRF / 讀本機檔：未查到 CVE/GHSA。原始碼：`IsExternalRefsAllowed` 預設 false（實測 ext=false）；但註解明載自訂 ReadFromURIFunc 會繞過此旗標。
- 大部分漏洞在 openapi3filter（執行期驗證），僅測試使用時暴露面低。

## 5. 評等：低到中
依據：測試圖隔離已實證、授權乾淨、目前版本無已知漏洞（技術面偏低）；但 2026 年公告密度極高（14 筆、含 critical/high，且有「修補不完整」案例 GHSA-vrcc），最新修補僅一週，維護集中於 1–2 人。依政策「低到中」→ 不同意。

## 若日後重審可附的條件（本次不適用）
釘 v0.151.0 及上列 go.sum hash；只限 _test.go；IsExternalRefsAllowed=false 且不設自訂 ReadFromURIFunc；CI 檢查 `go list -deps ./cmd/...` 無 kin-openapi；不得 import 到非測試套件。
替代：用既有工具做規格靜態檢查（例如 CI 中獨立執行的 linter），不進 go.mod。

## 無法查證
維護者正式名單/權限（僅有 contributor 統計）；GHSA-g3f9 的 CVE 與 Go vulndb 收錄；ReDoS 是否可實際觸發（未做 fuzz）；commit 數超過 100 的精確值；簽章/發布流程（無 sigstore 驗證）。

## 安全部結論
- **不同意引入。** 技術面偏低：只在測試依賴圖內，govulncheck 0，授權寬鬆。但 2026 年共 14 筆 GHSA（critical 1、high 8），曾有修補不完整，最新一筆在 10/14 才發布，實際維護者約 1–2 人，因此評為低到中。依使用者「只接受低風險，低到中一律排除」的政策，予以排除。
- 處置：研發部改用標準函式庫手寫 6 個契約測試；不得加入任何 YAML／JSON Schema 第三方套件。
- 若日後要重審，必須重新評估，並由使用者重新拍板；草稿中的條件僅供參考，不構成核准。

## 修訂紀錄
| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| v0.1 | 2026-10-22 | 安全部 | 由 09:09 書面結論的工作紀錄正式收錄為原件（回應審查部 PR #26 S1） |
