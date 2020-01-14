# Gsigo Web socketio and cmd Framework

Gsigo是一个用Go (Golang)编写的web、socketio、command框架。

gsigo 主要基于下面的包进行了封装,大部分保留了原有包的用法

https://github.com/gin-gonic/gin

https://github.com/googollee/go-socket.io

https://github.com/sirupsen/logrus

https://github.com/jinzhu/gorm

https://github.com/gomodule/redigo/redis

# 目录

- [安装](#安装)
- [快速开始](#快速开始)

## 安装

1. 首先需要安装 [Go](https://golang.org/) (**version 1.10+ is required**), 可以使用下面的命令进行安装 Gsigo.

```sh
$ go get -u github.com/whf-sky/gsigo
```

2. 导入你的代码

```go
import "github.com/whf-sky/gsigo"
```


## 快速开始

假设项目目录如下

- test
  - controllers
    - cmd
        - test.go
    - gin
        - index.go
    - sio
      - chat
        - message.go
      - root
        - connect.go
        - disconnect.go
        - error.go
        - notice .go
  - config
    - production
        - system.yml
  - routers
    - sio.go
    - gin.go
    - cmd.go
  - static
    - index.html
  - sio.go
  - gin.go
  - cmd.go
 
```sh
# 假设文件 test/main.go 中有如下代码：
$ cat main.go
```

```go

package main
 
 import (
 	"github.com/whf-sky/gsigo"
 	_ "test/routers"
 )
 
 func main() {
 	gsigo.Run()
 }
```

```sh
# 假设文件 test/main.go 中有如下代码：
$ cat main.go
```

```
# run example.go and visit 0.0.0.0:8080/ping (for windows "localhost:8080/ping") on browser
$ go run example.go
```


