# PR33-S2（recheck PR31-S4）：01 號 exit code 補記（維運部，2026-11-03，UTC+8）

01 號檔屬安全部（docs/06-security/evidence/cr-006-m5/…/01-m5-1-head.txt），依指示不得修改，也依安全部裁定不入庫；本檔只做旁註，不改原檔、不改其 SHA256SUMS。

- 826d5d34 目錄的 01 號：`git ls-remote` 與 `curl` 後各有一行 `exit 0`；clone／fetch／checkout 均為 `exit 0`。已有 exit code。
- 47475638 目錄的 01 號：`git ls-remote` 與 `curl … | jq` 兩個指令**沒有記錄 exit code**；clone／fetch／checkout 為 `exit 0`。
  兩指令的輸出完整（ls-remote 回傳 `474756381bd1…` refs/pull/31/head，API 回傳 head 相同），與（i）（ii）（iii）三值一致，可間接佐證指令成功，但**當時的 exit code 無法事後補得**，維運部不補寫數值。
- 建議：由安全部在下一版裁定或結案紀錄註明「47475638 之 01 號 ls-remote／curl exit code 未記錄，以輸出內容與審查部自取 log 為準」。
