---
文件：ADR-003 短碼編碼（長度／字元集／CSPRNG／熵）
版本：v0.1
狀態：審查中
負責角色：設計部
最後更新：2026-10-07
對應需求：REQ-006, REQ-010, SEC-006, SEC-013
專案代號：SHORTURL
---

# ADR-003：短碼編碼

| 欄位 | 內容 |
|---|---|
| 狀態 | **已採納**（D-04，2026-10-07，方案 A：base62／長度8／CSPRNG） |
| 日期 | 2026-10-07 |

## 背景

SEC-006 要求 CSPRNG 與文件化熵；REQ-006 要求可機器驗證之字元集／長度；OpenAPI path pattern 須一致。

## 選項

### 方案 A（建議）：base62（A-Za-z0-9）、長度 8、CSPRNG

- 字元集大小 62；空間 62^8 ≈ 2.18×10^14
- 熵 ≈ 8 × log2(62) ≈ **47.6 bit**
- 小流量演示下暴力枚舉不具實務價值；URL 友善、無填充字元

### 方案 B：base64url、長度 8

- 字元含 `-` `_`；熵略高但路徑／文件逃逸需注意；與「僅字母數字」建議不一致

### 方案 C：base62、長度 6

- 熵 ≈ 35.7 bit；空間較小，演示雖可但安全裕度較方案 A 低

## 熵估算（方案 A）

```text
alphabet = 62
length   = 8
space    = 62^8 ≈ 2.183×10^14
entropy  ≈ log2(62^8) ≈ 47.63 bits
```

假設攻擊者每秒嘗試 1e6 次（遠高於 SEC-007 限流），窮舉期望時間仍極長；加上速率限制，枚舉在 MVP 假設下不具實務價值。

## 建議傾向

**已採納方案 A**（D-04）：Go `crypto/rand`（CSPRNG；**禁止** math/rand 等弱來源）＋ base62 ＋ 長度 8。碰撞時重試生成（REQ-010）。

## 後果

- OpenAPI：`^[A-Za-z0-9]{8}$`
- UrlValidator.validateShortCode、ShortCodeGenerator 對齊
- 變更長度／字元集須 CR 並同步 OpenAPI／SEC-013 測試

## 對應需求

REQ-006、REQ-010；SEC-006、SEC-013。
| v0.1+D04 | 2026-10-07 | 設計部核對 | D-04 採納方案 A；對齊 Go crypto/rand |
