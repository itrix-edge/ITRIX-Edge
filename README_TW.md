![img](docs/images/ITRIX-Edge-logo-small.png) 
ITRIX-Edge: 軟體定義邊緣
================================

ITRIX-Edge旨為透過容器化與kubernetes(K8S)經驗，使得其上的邊緣運算可簡易的連接到各大主要雲端供應商。
此項目的主要重點是優化部署和維護，通過雲端供應商的內部連結，讓一般用戶可以輕鬆的將邊緣運算連接到雲端。

# 總覽

當前項目分為硬體與軟體兩部分。在硬體部分，主機針對節點做機器重啟與硬體相關監控。而軟體部分使用容器技術和K8S排程管理系統，
將系統服務和客戶應用系統放在K8S裡運作，提供應用管理以及雲端整合工作。
更多相關訊息請參考 [Overview document](doc/Overview.md).

# 貢獻

請參考 [CONTRIBUTING.md](CONTRIBUTING.md) 並了解如何使用.

如果您是硬體供應商，請撿查硬體是否相容，並驗證軟體的相容性，
我們會將您的硬體設備列入相容性列表裡，並歡迎提供硬體給我們進行測試。

# 社群

我們對外透明化項目排程的會議，只要對此項目有興趣的可以通過Zoom軟體加入會議一起討論。

## 定期會議
- 每隔週三 下午14:15~15:15(台灣時間Taiwan Standard Time; TST)
- Microsoft Team線上會議: https://teams.microsoft.com/l/meetup-join/19%3ameeting_Nzc1ODAxNjEtMjBiMy00NjY2LTk4NzUtYTE3ZjhlNjZmYzll%40thread.v2/0?context=%7b%22Tid%22%3a%2273ffe322-3a8c-4d1c-936d-665676f559aa%22%2c%22Oid%22%3a%2235847b3d-47f0-4419-9476-427efc9d7281%22%7d

## 定期會議紀錄
- 請參考 [Google doc](https://docs.google.com/document/d/1wQb8q7dXOevTFSIFiWSf9xacT_8qqiqOgxSLDL-Gn3E)
