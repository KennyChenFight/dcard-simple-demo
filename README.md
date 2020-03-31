# dcard-simple-demo

此專案是針對Dcard公司的申請作業，題目描述：

Dcard 每天午夜都有大量使用者湧入抽卡，為了不讓伺服器過載，請設計一個 middleware：

- 限制每小時來自同一個 IP 的請求數量不得超過 1000
- 在 response headers 中加入剩餘的請求數量 (X-RateLimit-Remaining) 以及 rate limit 歸零的時間 (X-RateLimit-Reset)
- 如果超過限制的話就回傳 429 (Too Many Requests)
- 可以使用各種資料庫達成

因此本專案除了做了該middleware，並模擬了小型Dcard系統的專案，例如有抽卡API及創建user API，因此背後需要有PostgreSQL、Redis兩種資料庫的支持。

## 專案講解

### 使用環境及工具

採用Golang語言的Gin框架開發RESTful API，使用PostgreSQL作為資料庫，以Redis做為快取資料庫，負責IP次數檢查的middleware。

### 如何運行該專案(使用docker-compose)

可利用本專案的docker-compose.yaml會一次啟動Backend、PostgreSQL、Redis，方便直接運行測試。請確保主機有docker環境，如果是Linux環境則需要另外安裝docker-compose套件。而如果是Windows、Mac則只需要安裝Docker Desktop即可。

#### Clone 專案

```bash
# 透過 git clone 專案到主機任意路徑下
```

#### 運行專案

````bash
# 在本專案的根目錄下執行以下指令即可
# -d 代表背景運行(Optional)
docker-compose up -d
````

總共起了以下五個Docker Container

+ postgres

  裡面預設會自動放入一些測試資料，可參考[測試資料](schema/testing_data.md)

+ postgres-client

  這個是採用pgAdmin作為client端，開啟**localhost:5432**，並需要輸入帳號密碼: **dcard@example.com**、**dcard**，登入後連接postgres即可，注意hostname是**postgres**，至於連接的帳號密碼可參考[創建db說明](schema/create_db.sql)

+ redis

+ redis-client

  這個是採用phpRedisAdmin作為client端，開啟**localhost:8081**，登入帳密為**admin**、**admin**，登入進去就會直接連接Redis，在左邊即可看到當前的所有key List，點擊key即可看到儲存的資料內容。

+ dcard-simple-demo

  這個Golang的Backend Service。

  可打開URL：localhost:8080/api-docs，利用Swagger框架打造RESTful API，可透過該框架直接測試API。

### 如何架設並運行此專案(針對Ubuntu)

#### Go 1.13+

建議去官網下載Golang 1.13+的編譯器，而不在Ubuntu直接安裝，因為Ubuntu的Golang是舊版的。

安裝指令如下：

```bash
# 下載 Golang 壓縮檔
wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
# 解壓縮
sudo tar -xvf go1.13.linux-amd64.tar.gz
# 放置到/usr/local
sudo mv go /usr/local
# 設定Golang環境變數：GOROOT、GOPATH、
nano ~/.profile
# 加入以下內容
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
# 生效設定
source ~/.profile
```

Note：Golang開發慣例是在GOPATH路徑下，會有三個資料夾需要創建，分別是，**bin**、**pkg**、**src**，建議可以先創建好。

#### Clone專案

```bash
# 透過 git clone 專案到GOPATH/src 路徑下
```

#### PostgreSQL  10.12+

可參考[schema目錄的安裝教學](schema/readme.md)

並且會教你如何一步一步地建立資料庫內所需的物件。

#### Redis  3.0.5+

詳細安裝教學可參考我的Blog文章：[Redis安裝教學](https://kennychen-blog.herokuapp.com/2020/02/22/Redis-高流量系統不能或缺的資料庫！安裝教學！/)

簡單指令如下：

```bash
# 下載Redis原始碼
wget http://download.redis.io/releases/redis-3.0.5.tar.gz
# 解壓縮
tar xzf redis-3.0.5.tar.gz
# build原始碼
cd redis-3.0.5
make
# 啟動Redis Server
./src/redis-server
```

Note：由於本專案，預設連接Redis時需要密碼，但因為安裝完Redis，Redis預設是沒有設定密碼的，因此需要多加設定密碼的步驟，如果不想要設定密碼，則需要去main.go刪除**連接Redis的Password選項**。

Redis設定密碼：

修改配置文件：**/redis-3.0.5/redis.conf**

```bash
# Warning: since Redis is pretty fast an outside user can try up to
# 150k passwords per second against a good box. This means that you should
# use a very strong password otherwise it will be very easy to break.
#
# requirepass foobared
```

找到被註解掉的**requirepass foobared**

因此把註解移除，以及把**foobared**換成你想要設定的密碼即可～

本專案預設的Redis密碼為**root**。

接著儲存配置檔案，需要將redis-server重啟：

```bash
./redis-cli shutdown
./redis-server
```

#### 檢查環境變數是否一致

當前面Golang、PostgreSQL、Redis環境都建置好之後，在本專案下有一個.env檔案，裡面定義了本專案需要的環境變數。

基本上，如果是按照前面教學一步一步操作，則裡面的環境變數的值不需要做任何更改，直接用本專案的預設值即可。若有改變，請自行修正。

環境變數介紹：

+ JWT 設定

  本專案採用JWT作為認證機制，因此需要**SECERT_KEY**、**TOKEN_LIFETIME**。

+ PostgreSQL 設定

  需要定義**DB_HOST**、**DB_PORT**、**DB_NAME**、**DB_USERNAME**、**DB_PASSWORD**、**DB_SSL_MODE**。

+ Redis 設定

  需要定義**REDIS_ENDPOINT**、**REDIS_PASSWORD**、**REDIS_POOL_SIZE**。

### RESTful API 文件參考

當前面必要的環境都架設好之後，並成功運行後

接著，可打開URL：localhost:8080/api-docs

這邊利用的Swagger框架，作為靜態檔案直接放置在本專案上，並定義好本專案所寫的RESTful API，方便測試用，就不需要額外使用其他如Postman等工具來測試。

### 針對題目的Middleware講解

#### 為何選擇Redis

+ Redis讀寫快速，且專門存放Hot Data，就算這類Hot Data遺失了對於系統的損失是輕微的，不會影響整個系統的運作。而且對應該middleware的功能，Redis只需存放IP連線次數及到期時間，這類資料不需要也沒必要放在系統的主RDBMS。
+ 不應該使用PostgreSQL這類的RDBMS去做該middleware，系統的主RDBMS應該是要負責重要的API端口使用，否則由於Dcard的使用人數過多，導致高併發流量，若重要的API前面都加了該middleware去做檢查，會導致在middleware這層，RDBMS效能就會被大量拖垮。

所以，考量高併發的要求，使用Redis這類資料庫在前面middleware層作緩解是必要的方法，主RDBMS理應是Backend的最後一層且不應輕易被使用到。

#### 如何實現IP次數檢查

在Redis利用Hash結構存放資料：

+ Hash Name

  將API URL + API Method + IP位址當作Hash Name，因為考慮到其他API可能也會用該middleware，以URL+Method為前綴可以區分不同API的訪問流量控管。

+ count

  裡面存放名為count的key name，值存放的是訪問API次數，最高數值只會到10，符合題目要求，而每當過了一個小時後，會將數值歸0重新計算。

+ reset

  裡面存放名為reset的key name，值存放的是到期時間，格式為timestamp，而當到期時間小於現在時間時，會將到期時間改為現今時間+一小時在重新存入。

就我對題目的理解，每次回傳要給剩餘次數及到期時間，因此透過Hash結構存會比較好操作。

而在程式實作上有兩種方法：

+ Golang利用Redis Client Library操作

  簡單來說，就是透過Golang呼叫Redis的GET、SET等操作，但是會有以下缺點：

  由於每次使用者訪問都需要透過讀取Redis數據，才能判斷要存入怎樣的次數或是到期時間，再加上Redis沒辦法將多個操作都視為Atomic操作，就算將這些操作都包在Redis Transaction結構裡面，GET操作在Transaction是沒辦法得到值的，造就程式判斷錯誤，進而造成Race Condition，會使得該使用者可能會得到相同的剩餘次數。也就是說使用者可以在一個小時內拜訪不只1000次，無法完全阻擋住。

+ 利用Lua腳本來對Redis進行操作

  Redis官方有說：推薦使用Lua腳本對Redis進行操作，而非使用Transaction方式。並且Redis對於Lua腳本是原生支持原子性，因此不會有第一種方法的問題，加上在Lua腳本寫會有更大的彈性。

本專案採用第二種方法，可以嚴格的讓使用者一小時內確實只能拜訪一千次，而不會造成Race Condition。但其實兩種方法各有優缺點：對於第一種方法而言，如果系統並沒有高流量，即使使用者會多拜訪個幾次，對於系統負擔造成不是很大，就可以採用第一種方法，因為第二種方法雖然不會Race Condition，但是需要利用Lua語言來寫腳本，會是額外的負擔。但我覺得對於極度高併發的系統，使用第二種方法會是較好的方法。

#### 解釋Lua腳本

```lua
local key = KEYS[1]
local now = tonumber(ARGV[1])
local ipLimit = tonumber(ARGV[2])
local period = tonumber(ARGV[3])
local userInfo = redis.call('HGETALL', ip)
local reset = tonumber(userInfo[4])
local result = {}
if #userInfo == 0 or reset < now then
    reset = now + period
    redis.call('HMSET', key, "count", 1, "reset", reset)
    result[1] = ipLimit - 1
    result[2] = reset
    return result
end

local count = tonumber(userInfo[2])
if count < ipLimit then
    local newCount = redis.call('HINCRBY', key, "count", 1)	
    result[1] = ipLimit - newCount
    result[2] = reset
    return result
else
    result[1] = -1
    result[2] = reset
    return result
end
```

**key**、**now**、**ipLimit**、**period**，為外部參數，由Golang程式丟進去，如此可以客製化，而不是寫死在Lua腳本，特別注意**key**是由**API URL + API Method + IP位址**所組成的。

講解該兩種if判斷含意：

+ 檢查userInfo是否為0 或是 到期時間 < 現今時間
  1. userInfo為0代表第一次的訪問，Redis沒有該使用者之前的資料
  2. 到期時間 < 現今時間，代表Redis內該使用者的資料要重設訪問次數以及到期時間
+ 檢查count < ipLimit
  1. 如果count < ipLimit，代表使用者還可以繼續訪問
  2. 反之，代表使用者訪問次數已經超量

Lua腳本執行完會回傳一個陣列，長度為二，第一格資料是使用者還可訪問的次數(如果回傳-1代表超量訪問)，第二格資料是到期時間。如此以來，計算訪問次數及到期時間都不是在Golang程式進行計算，更防止Race Condition的可能性。

### 總結

本專案是將**/v1/pairs/**(POST 抽卡配對使用者)、**/v1/pairs/**(GET 取得配對對象)前面都加了IP Limit檢查middleware，因為我覺得這兩隻API是高流量的，因此在運行的時候會看到Redis資料裡面會有這兩隻API所對應的使用者拜訪次數。

