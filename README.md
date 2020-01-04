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
        - 被Go语言的HTTP包解析出RequestURI: `/send?word=hello`
        - 再通过`http.Request.FormValue("word")`得到RequestURI里面以`word`为key的字符串`hello`

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
    - 写数据的过程
