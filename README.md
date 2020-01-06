# Chapter 1. Our First API in Go
## 配置Unbuntu环境
### 配置Go
* 解压tar包
* 配置Go
```
> sudo chown -R root:root ./go
> sudo mv go /usr/local
> vim ~/.profile
    export GOPATH=$HOME/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
> source ~/.profile
```
### 配置MySQL
* 安装mysql-server
```sh
> sudo apt-get install mysql-server
> sudo systemctl status mysql
```
* 创建一个用户ben
```sh
> sudo mysql
# 新建一个用户，他可以来在任何地址
mysql> CREATE USER 'ben'@'%' IDENTIFIED BY 'user_password';
# 查看新建用户
mysql> select host, user from mysql.user;
# 给这个用户开所有数据库的权限
mysql> GRANT ALL PRIVILEGES ON *.* TO 'ben'@'%';
# 查看权限
mysql> SHOW GRANTS FOR 'ben'@'%';
# 刷新
mysql> FLUSH PRIVILEGES;
```
* 用新用户测试mysql
```sh
# 登录
> mysql -u ben -p
# 新建一个hello的数据库
mysql> CREATE DATABASE hello;
# 显示所有已存在的数据库
mysql> show databases;
# 使用新建的hello数据库
mysql> USE hello;
# 查看新建数据库中的表格
mysql> SHOW TABLES;
```
### 配置Redis
* 安装
```sh
# sudo apt-get install redis-server
```
* 配置
```
# 开启权限
> sudo vim /etc/redis/redis.con
    #修改 `supervised no` 为 `supervised systemd`
> sudo systemctl restart redis.service
```
* 测试
```sh
> sudo systemctl status redis
> redis-cli
redis> ping
    output> PONG
redis> set test "It's working!"
    outupt> OK
redis> get test
    output> "It's working!"
redis> exit
```
* 使能Redis可以从任何地方访问，模式只能从localhost访问
```sh
> sudo vim /etc/redis/redis.conf
    找到`bind 127.0.0.1 ::1`，取消注释
# 查看redis端口
> sudo netstat -lnp | grep redis
```
* 其他配置，例如设置密码，参考：[link](https:#www.digitalocean.com/community/tutorials/how-to-install-and-secure-redis-on-ubuntu-18-04)

## 写入数据到MySQL
```sh
> mysql -u ben -p
# 新建数据库
mysql> create database social_network;
# 选择数据库
mysql> use social_network;
# 新建表格
# PRIMARY KEY关键字用于定义列为主键, 唯一的标识该表的每条信息
# UNIQUE 标识该属性的值是唯一的
mysql> CREATE TABLE users (
    user_id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_nickname VARCHAR(32) NOT NULL,
    user_first VARCHAR(32) NOT NULL,
    user_last VARCHAR(32) NOT NULL,
    user_email VARCHAR(128) NOT NULL,
    PRIMARY KEY (user_id),
    UNIQUE INDEX user_nickname (user_nickname)
);
# 查看表格
mysql> show tables;
# 查看表格中的数据
mysql> select * from users;
```

## Hello World via API
* [http-hello-to-client](./code/hello/hello.go)
    - open [http://localhost:8080/api](http://localhost:8080/api)
    - 服务器向浏览器发送Hello字符串
* [http-hello-to-server](./code/hello/hello-to-server.go)
    - 在浏览器输入<br>
    > http://localhost:8080/send?word=hello
    - 在服务器CMD窗口中显示hello字样
    - 工作过程：<br>
        - 浏览器通过GET包传来URL: http://localhost:8080/send?word=hello
        - 被Go语言的HTTP包解析出RawQuery: `word=hello`
        - 再通过`http.Request.FormValue("word")`得到RawQuery里面以`word`为key的字符串`hello`
* 往Http包的Header中写入信息：[hello-set-header.go](./code/hello/hello-set-header.go)
    - 输入下面命令测试：<br>
    > curl --head http://localhost:8080
    - 返回输入数据：
    ```
    HTTP/1.1 200 OK
    Pragma: no-cache
    Date: Sun, 05 Jan 2020 01:33:14 GMT
    Content-Length: 5
    Content-Type: text/plain; charset=utf-8
    ```
    其中`Pragma: no-cache`就是通过`w.Header().Set("Pragma", "no-cache")`配置进去的。

## Building first route
* Multiplexer
    - refers to taking URLs or URL paterns and translating them into internal functions
    - mapping from a request to a function
        - /api/user -> func apiUser
        - /api/message -> func apiMessage
        - /api/status -> func apiStatus
* [router-gorilla](./code/router-gorilla/hello.go)
    - `gorilla/mux`是一个GO语言的URL Router and Dispatcher
    - open [http://localhost:8080/api/123](http://localhost:8080/api/123)

## Getting and Setting data via HTTP
* 连接MySQL: [mysql-connect](./code/mysql/connect.go)
    - `root`账户可能会没有权限写入数据库
* 与MySQL进行数据交互：[mysql-get-set-data](./code/mysql/get-set-data.go)
    - 写数据到数据库<br>
    > http://localhost:8080/api/user/create?user=nkozyra&first=Nathan&last=Kozyra&email=nathan@nathankozyra.com
    - 从数据库读数据<br>
    > http://localhost:8080/api/user/read/1
    - 创建用户的过程<br>
        - 浏览器通过GET包（正常应该用Post包发送，这里为了方便，将数据直接从URI中发送过来了）将数据通过`http.Request`发送到服务器
        - `http.Request`根据Key值，解析出对应的Value
        - 在`CreateUser`回调函数中，将数据组装好后，写入数据库
        - 打印数据后返回
    - 读取用户的过程<br>
        - 浏览器发送GET包到服务器，并带上`ID`参数
        - 在`GetUser`回调函数中，从数据库中找到此`ID`对应的数据
        - 将数据库返回的数据编码成JSON格式后，通过`http.ResponseWriter`返回给浏览器

# Chapter 2. RESTful Services in Go
## 配置MySQL数据库
* 将`user_email`加上`UNIQUE INDEX`属性
```sql
-- 准备工作
> show databases;
> use social_network;
> show tables;
> select * from users;
-- 删除重复的行
> DELETE FROM users WHERE user_id=3;
-- 修改属性
> ALTER TABLE users ADD UNIQUE INDEX user_email (user_email);
```
* 添加新表，用于表示用户关系
```sql
-- 添加users_relationships表
CREATE TABLE users_relationships (
    users_relationship_id INT(13) NOT NULL,
    from_user_id INT(10) NOT NULL,
    to_user_id INT(10) UNSIGNED NOT NULL,
    users_relationship_type VARCHAR(10) NOT NULL,
    users_relationship_timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (users_relationship_id),
    INDEX from_user_id (from_user_id),
    INDEX to_user_id (to_user_id),
    INDEX from_user_id_to_user_id (from_user_id, to_user_id),
    INDEX from_user_id_user_id_users_relationship_type (from_user_id, to_user_id, users_relationship_type)
);
-- 查看新建表的Index
SHOW INDEX FROM users_relationships;
```
## 编码格式
### JSON
* [json-hello](./code/json/hello.go)
### XML
* [xml-hello](./code/xml/hello.go)
### YAML
* [yaml-hello](./code/yaml/hello.go)

## Comparing the HTTP actions and methods
### POST and PUT
* Post - Createing data
* PUT - Updating data
* 例如对同样的URI：/api/users/1234
    - Post request to `/api/users/1234` will tell our web service that we will be creating a new user resource based on the data within.
    - Put request to `/api/users/1234` will tell our web service that we're accepting data that will update or overwrite the user resource data for our user with the ID 1234.
### PUT and PATCH
* Put - 更新整个资源
* Patch - 更新部分资源
### CRUD and REST
* CREATE -- POST
* RETRIEVE(READ) -- GET
* UPDATE -- PUT/PATCH
* DELETE -- DELETE
### Endpoints 设计
|    Endpoint   | Method  | Purpose                                           |
|:-------------:|---------|---------------------------------------------------|
| /api          | OPTIONS | To outline the available actions within the API   |
| /api/users    | GET     | To return users with optional filtering paramters |
| /api/users    | POST    | To create a user                                  |
| /api/user/123 | PUT     | To update a user with the ID 123                  |
| /api/user/123 | DELETE  | To delete a user with the ID 123                  |

* [mysql-get-set-data-at-same-endpoint](./code/mysql/get-set-data-at-same-endpoint.go)
    - 利用Gorilla的`Methods()`方法去区分这个包是GET还是POST

# Chapter 3. Routing and Bootstrapping
## Writing custom routers in Go
* 对URI进行简单的正则表达式检查：[router-regex-raw](./code/router-regex/raw.go)
    - 浏览器输入：[http://localhost:8080/testing1](http://localhost:8080/testing13)，不满足
    - 浏览器输入：[http://localhost:8080/testing1](http://localhost:8080/testing1234)，满足

## WebSockets
* [How to Use Websockets in Golang](https://yalantis.com/blog/how-to-build-websockets-in-go/)
* Network socket
    - Datagram sockets (SOCK_DGRAM): UDP, connectionless sockets
    - Stream sockets (SOCK_STREAM): TCP, SCTP, DCCP, connection-oriented sockets
    - Raw sockets (raw IP sockets): available in routeres and other networking equipment
* WebSockets
    - Don't require clients to send a request in order to get a response
    - One handshake between a browser and server for establishing a connection that will remain active throughout its lifetime
    - a client request looks like:<br>
    ```
    GET /chat HTTP/1.1
    Host: server.example.com
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
    Sec-WebSocket-Protocol: chat, superchat
    Sec-WebSocket-Version: 13
    Origin: http://example.com
    ```
    - server response:<br>
    ```
    HTTP/1.1 101 Switching Protocols
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=
    Sec-WebSocket-Protocol: chat
    ```
* WebSocket是对RESTful API的一种补充，应对某些特殊的S/C交互。WebSocket其实是一个遵守HTTP协议的TCP的长连接，提供的接口和普通`socket`类似（`read/write`函数）。
* WebSocket Hello
    - [websocket-hello-server](./code/websocket/hello-server.go)
    - [websocket-hello-client](./code/websocket/hello-client.html)
    - 启动服务器后，在浏览器输入：[http://localhost:8080/wsclient](http://localhost:8080/wsclient)
    - 由于常规情况下，浏览器无法向服务器发送带websocket格式(`Upgrade: websocket`)的包，所以client要利用JS和WebSocket服务器通信。
    - URI:`/wsclient`会打开`hello-client.html`网页，此网页会连接服务器的URI:`/ws`，然后进行WebSocket通信。
* WebSocket Echo Length
    - [websocket-echo-length-server](./code/websocket/echo-length-server.go)
    - [websocket-echo-length-client](./code/websocket/echo-length-client.html)
    - 启动服务器后，在浏览器中输入：[http://localhost:12345/websocket](http://localhost:12345/websocket)

# Chapter 5. Templates and Options in Go
## HTTPS
* [https-hello](./code/https/hello.go)
    - 利用Go工具生成钥匙：`cert.go`和`key.pem`<br>
    > go run generate_cert.go --host localhost --ca true


