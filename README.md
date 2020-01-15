# Gsigo Web socketio and cmd Framework

Gsigo是一个用Go (Golang)编写的web、socketio、command框架。

gsigo 主要基于下面的包进行了封装, 保留了原有包的用法

https://github.com/gin-gonic/gin

https://github.com/googollee/go-socket.io

https://github.com/sirupsen/logrus

https://github.com/jinzhu/gorm

https://github.com/gomodule/redigo/redis

# 目录

- [安装](#安装)
- [快速开始](#快速开始)

## 安装

1. 首先需要安装 [Go](https://golang.org/) (**version 1.10+**), 可以使用下面的命令进行安装 Gsigo.

```sh
$ go get github.com/whf-sky/gsigo
```

2. 导入你的代码

```go
import "github.com/whf-sky/gsigo"
```


## 快速开始

```sh
# 假设文件 main.go 中有如下代码：
$ cat main.go
```

```go
package main

import (
	"github.com/whf-sky/gsigo"
	"net/http"
)

type IndexController struct {
	gsigo.Controller
}

func (this *IndexController) Get() {
	this.Ctx.String(http.StatusOK, "test")
}

func main()  {
	gsigo.GET("/", &IndexController{})
	gsigo.Run()
}
```


