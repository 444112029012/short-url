# CR-005 PR #29 K1／K2 書面確認（安全部）

- 文件編號：06-security/31
- 版本：v0.1.1（Draft；v0.1 sha256 f36f45b4…）
- 負責：安全部
- 關聯：CR-005 v0.3（bfaa6d1d…）、D-19；30 號 v0.1.1（e4579d31…）§3 K1～K6；impact/安全部.md v0.2.1（5ea337a3…）
- 綁定：PR [#29](https://github.com/444112029012/short-url/pull/29) head `1545143ca94a3a92f6877a30afa1762de994eedf`，base main `36f4717caab8a957e60f0a67ef5b01400f69ef52`
- 撰寫時間：2026-10-27 13:17:14（UTC+8，取自 `date`）

## 1. 結論
**K1 成立；K2 成立。** 本確認只對上述 head 有效；head 一變動就失效，須重新確認。K3（D-09 核 head、合併後 main 全綠）、K4（合併後 7 日內 G404 canary）、K6（Go 1.27.3 CR 最遲 2026-11-20）仍待完成。合併由使用者決定。

## 2. K1 查證（安全部在 /tmp/v23 以公開 clone 自行執行）
| 項目 | 結果 |
|---|---|
| head 父節點 | 只有一個父節點 `36f4717c`，等於目前 origin/main |
| 與 origin/main 的 diff | 只有 `.github/workflows/ci-devsecops.yml`，增 3、刪 3 |
| 改後檔案 sha256 | `a67a0be4f0fba76f030cc7394202db91d9de6635b2b89c17a67e783211dece8f`，與 `/workspace/rebuild/cr-005/ci-devsecops.proposed.yml` 相同 |
| 整份 diff sha256 | `6c13fa7a2ccf6008a9e8ba653c28f93972b6a1e9f12df8c3eb87dfd4650ba7c4`，與維運部回報相同 |
| 3 行內容 | 第 8 行註解（CR-005 換回 v2.30.0）；`go install …/gosec@v2.30.0`；`expected='h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0='` |
| 門檻未變 | 仍是 `"$(go env GOPATH)/bin/gosec" -fmt=text -no-fail=false -severity=high -tests ./...`（絕對路徑，SCR-019；-tests，SCR-032）；setup-go 仍是 `1.27.2`、`check-latest: false` |
| main 的 workflow | 2425129 到 36f4717c 之間 `.github/` 沒有變動，main 檔案 sha256 前綴 `da73b404` |

## 3. K2 查證（GitHub 未登入 API）
- run 40454105428：event pull_request，head_sha `1545143c…`，completed／success，建立時間 2026-10-27 13:13:41（UTC+8；API 值為 05:13:41Z）。
- head 的 check-runs 共四項，全部 success：sast-gosec（job 122023300733）、sca-govulncheck、secrets-gitleaks、unit-test。
- sast-gosec 各步驟都是 success：setup-go → Install gosec → **Verify gosec module hash** → Run gosec。h1 驗證步驟寫死 expected＝`h1:67aF6…`，不相符就 exit 1；這一步成功，代表執行檔的 gosec 模組 h1 和 v2.30.0 的鎖定值相同。
- Annotations 只有 Node.js 20 淘汰警告和 ubuntu-latest 遷移通知，跟資安無關。
- **限制**：未登入的帳號下載 job log 會回 403（Must have admin rights）。所以 log 裡的四行原文（`go1.27.2`、`mod … v2.30.0 h1:67aF6…`、`gosec h1 OK`、`Issues : 0`）是維運部轉述的，安全部沒有直接讀到 log。安全部是從 workflow 內容加上各步驟都成功，推論 h1 和 Go 版本成立。

## 4. 補件（不阻擋本確認，合併前完成）
- 比照 C6 的作法（`evidence/cr-002-c6/run-40430102886.log`），請維運部把 run 40454105428 的 sast-gosec 原始 log 和 run JSON 存到 `docs/09-cicd-env/evidence/cr-005/`，並更新 SHA256SUMS。安全部收到後會核對上述四行原文，再在本檔小版補記。

## 5. 補件核對（v0.1.1，2026-10-27 13:18:30 UTC+8，取自 `date`）
維運部已把原始 log 存在本機 `docs/09-cicd-env/evidence/cr-005/`，**還沒進 repo**；合併後會用另一個文件 PR 送進去，推送前要先取得使用者核准。這些檔案**不可**加進 #29，否則 head 會改變、K1 也會失效。
| 檔案 | sha256 | 核對結果 |
|---|---|---|
| 17-pr29-run-40454105428.json | 6c66a1be61fc54fb8a325ebafe5093e0d6c055ea07ce6564826c1d62f57438dc | headSha 1545143c、pull_request、success；四個 job 都是 success，job id 與 API 相符 |
| 18-pr29-run-40454105428-sast-gosec.log | 40dd0eb81b6421f1628796c8e63fb8a0b0baa8fd46c5ba35fc748edc65d5a9be | 共 417 行。第 155 行 `go version go1.27.2 linux/amd64`；第 271 行 `mod github.com/securego/gosec/v2 v2.30.0 h1:67aF6LHa7SAL6QEc4es857GrXTT50gUqmrl0XPhtDx0=`；第 272 行 `gosec h1 OK`；第 392 行 `Issues : 0`。第 120 行 checkout 為 `Merge 1545143c… into 36f4717c…`；第 277 行指令含 `-severity=high -tests`、絕對路徑；沒有 `::error` 輸出 |
| SHA256SUMS | 1e3a0dca41a87faf307ac3e2cc517a4bb4d7c2cec925bda48cfbea5f08115dc1 | `sha256sum -c` 全部 OK |

log 的時間戳和公開 API 的步驟時間一致：開始 05:13:44Z，h1 驗證 05:14:57Z，結束 05:15:30Z。log 是維運部用 admin 權限下載的，安全部無法另外取得原始來源比對；就內容和時間戳來看，沒有不一致的地方。**§3 的限制已解除**，四行原文已由安全部直接核對。K1、K2 結論不變，仍只綁 head 1545143c。

## 修訂紀錄
| 版本 | 日期 | 作者 | 說明 |
|---|---|---|---|
| v0.1 | 2026-10-27 | 安全部 | 初版：綁 head 1545143c，K1、K2 成立 |
| v0.1.1 | 2026-10-27 | 安全部 | 加入 §5：核對維運部存下的原始 log 和 run JSON，解除 §3 的限制；結論不變；小版修訂不改檔名 |
