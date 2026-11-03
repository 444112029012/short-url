# PR33-S2：L2 取證補記（維運部，2026-11-03，UTC+8）

判準（CR-006 v0.4 §5.1 L2）：快取 key 以 `setup-go-Linux-x64-ubuntu24-go-1.27.3-` 開頭，雜湊段＝`cf6b0d429736099b597565b998fd2f9b82c99c65d9cffc4d8c9270df02c0a025`；命中與否只記錄，不作判定（QA-M2）。
來源：既有本機日誌（未重新下載、未 rerun），原檔未修改，逐行引用如下（時間為日誌原始 UTC）。

## run 41403811744（PR #31 head 47475638）— 來源 /workspace/rebuild/cr-006/evidence-pr31-47475638/
- sast-gosec（job 125141230822）第 148 行：`Cache hit for: setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d42…c0a025`
- sca-govulncheck（job 125141230836）第 149 行：同上 key，Cache hit
- unit-test（job 125141230942）第 149 行：同上 key，Cache hit
- 判定：三個 job 的前綴與雜湊段皆符合 → L2 符合（命中，僅記錄）

## run 41559091290（main 1b53fcde）— 來源 /workspace/rebuild/cr-006/evidence-main-1b53fcde/
- 三個 job 還原時皆為 `Cache is not found`（sast-gosec 第 132 行、sca-govulncheck 第 133 行、unit-test 第 133 行）
- sca-govulncheck 第 214 行：`Cache saved with the key: setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d42…c0a025`
- sast-gosec 第 381 行、unit-test 第 253 行：`Failed to save: Unable to reserve cache with key setup-go-Linux-x64-ubuntu24-go-1.27.3-cf6b0d42…c0a025, another job …`（同 key 由另一 job 寫入）
- 判定：key 前綴與雜湊段符合 → L2 符合（未命中，僅記錄；main 不使用 PR ref 快取，與 CR §R1 預期一致）
