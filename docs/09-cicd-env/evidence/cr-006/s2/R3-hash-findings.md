# PR33-R3：CR-006 v0.4 原件雜湊（維運部，2026-11-03，UTC+8）

檔案：CR-006-ci-go-1.27.3-v0.4.md（由 /workspace/rebuild/cr-006/ 原封複製，cmp 相同，未修改）
- sha256sum：818f48eb8fdbcc38bdecee2f0bbe3d5de0e0e1329088892dff261971a0409734 → **以 818f48eb 開頭，與紀錄「818f48eb…734」相符**
- git hash-object（git blob SHA-1）：781766af128ddbf8bc20241c8e55e459064f74ba → **不是 818f48eb**

結論：818f48eb 是檔案的 sha256，不是 git blob id。審查紀錄 PR33-R3 寫「blob 必須等於 818f48eb」，用詞與實際不符；入庫後 git blob 會是 781766af…，sha256 會是 818f48eb…。建議審查部把條件改為「sha256＝818f48eb…734（git blob 781766af…）」。檔案不做任何修改。
