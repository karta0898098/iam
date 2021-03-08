IAM
========

IAM 是一個簡單示範 identity & access token 管理伺服器。主要提供
user 認證、登入、註冊以及角色管理等。目前只為了展示作者的 golang 風格而開發的。

## Require Environment
啟動 Server 所需要的服務
* MySQL 
* Redis 

P.S 專案底下的 [deployments/environment](https://github.com/karta0898098/iam/tree/master/deployments/environment "link") 提供了所需要環境的 [docker-compose.dev.yml](https://github.com/karta0898098/iam/blob/master/deployments/environment/docker-compose.dev.yml "link") 方便快速啟動環境。</br>
也可以參考 Makefile 來啟動所需要環境
```
    make iam.dev.env
```

## Config 
啟動本服務時會預設載入 ./deployments/config/[app.toml](https://github.com/karta0898098/iam/blob/master/deployments/config/app.toml "link") 的設定檔案。請依真實環境配置

## 專案說明
此專案使用三層式設計，並參考了DDD的設計，將服務拆分成聚合層、商業邏輯層以及存儲資料層來方便維護以及開發。</br>

* Identity Module 提供 User 的身份別，註冊登入等相關邏輯。
* Access Module 提供 User 具備何種角色，以及發放 Token 的管理。

**NOTE: 目前尚未開發完成**

## License
IAM source code is available under an MIT [License](https://github.com/karta0898098/iam/blob/master/LICENSE "link").




