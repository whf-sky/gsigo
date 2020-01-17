# Gsigo Web socketio and cmd Framework

###### Gsigo是一个用Go (Golang)编写的web、socketio、command框架。

###### gsigo 主要基于下面的包进行了封装, 保留了原有包的用法

https://github.com/gin-gonic/gin

https://github.com/googollee/go-socket.io

https://github.com/sirupsen/logrus

https://github.com/jinzhu/gorm

# 目录

- [安装](#安装)
- [快速开始](#快速开始)
- [配置文件](#配置文件)
    - [应用配置文件](#应用配置文件)
    - [REDIS配置文件](#Redis配置文件)
    - [数据库配置文件](#数据库配置文件)
- [路由规则](#路由规则)
    - [WEB路由规则](#WEB路由规则)
    - [SOCKETIO路由规则](#SOCKETIO路由规则)  
    - [CMD路由规则](#CMD路由规则)  
- [WEB应用](#WEB应用)  
- [SOCKETIO应用](#SOCKETIO应用)
- [CMD应用](#CMD应用)  
- [数据库](#数据库)
    - [CURD](#CURD)
        - [事务](#事务)
        - [主从强制切换](#主从强制切换)
        - [Create](#Create)
        - [Delete](#Delete)
        - [Update](#Update)
        - [Query](#Query)
    - [MODEL](#MODEL)
- [REDIS](#REDIS)
- [日志](#日志)
- [环境变量](#环境变量)  

## 安装

###### 1. 首先需要安装 [Go](https://golang.org/) (**version 1.10+**), 可以使用下面的命令进行安装 Gsigo.

```sh
$ go get github.com/whf-sky/gsigo
```

###### 2. 导入你的代码

```go
import "github.com/whf-sky/gsigo"
```

如使用go mod包依赖管理工具,请参考下面命令

###### Windows 下开启 GO111MODULE 的命令为：
```sh
$ set GO111MODULE=on
```

###### MacOS 或者 Linux 下开启 GO111MODULE 的命令为：
```sh
$ export GO111MODULE=on
```

###### Windows 下设置 GOPROXY 的命令为：
```sh
$ go env -w GOPROXY=https://goproxy.cn,direct
```

###### MacOS 或 Linux 下设置 GOPROXY 的命令为：
```sh
$ export GOPROXY=https://goproxy.cn
```



## 快速开始

###### 假设文件 main.go 中有如下代码：

```sh
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
	gsigo.Run()
}

func init() {
	gsigo.GET("/", &IndexController{})
}
```
## 配置文件

###### gsigo 默认是不加载配置文件的，配置文件格式.ini文件

### 应用配置文件

配置文件的使用

```go
package main

func main()  {
    gsigo.Run("./config/.app.ini")
}
```

不加载配置文件的默认参数：

```ini
app.name = "gsigo"
app.debug = true
app.host = 0.0.0.0
app.port = "8080"
app.mode = 'default'

socket.ping_timeout = 60
socket.ping_interval = 20

log.hook = "default"
log.formatter = "text"
```


不同级别的配置：

###### 当使用环境变量时，当前环境变量会替换调公共环境变量信息，环境变量需自定义

```go
app.name = "gsigo"
app.debug = true
app.host = 0.0.0.0
app.port = "8080"
app.mode = 'default'

socket.ping_timeout = 60
socket.ping_interval = 20

log.hook = "stdout"
log.formatter = "text"
log.params.priority = "LOG_LOCAL0"

[production]
app.name = "test"

[develop]


[testing]
```

App 配置

- app.name

###### 应用名称，默认值`gsigo`。

###### 配置文件中设置

```ini
app.name = "gsigo"
````

###### 代码中调用

```go
gsigo.Config.APP.Name
````

- app.debug

###### 应用debug，默认值`true`。

###### 配置文件中设置

```ini
app.debug = true
````

###### 代码中调用

```go
gsigo.Config.APP.Debug
````


- app.host

###### 应用HOST，默认值`0.0.0.0`

###### 配置文件中设置

```ini
app.host = 0.0.0.0
````

###### 代码中调用

```go
gsigo.Config.APP.Host
````

- app.port

###### 应用PORT，默认值 `8080`。

###### 配置文件中设置

```ini
app.port = "8080"
````

###### 代码中调用

```go
gsigo.Config.APP.Port
````

- app.mode

`default` `gin` `cmd`

###### 应用模式，默认值 `default`(默认：gin+socketio)。


###### 配置文件中设置

```ini
app.mode = 'default'
````

###### 代码中调用

```go
gsigo.Config.APP.Mode
````

SOCKETIO配置

- socket.ping_timeout

###### ping 超时时间，默认值 `60`。

###### 配置文件中设置

```ini
socket.ping_timeout = 60
````

###### 代码中调用

```go
gsigo.Config.Socket.PingTimeout
````


- **socket.ping_interval**

###### ping 时间间隔，默认值 `20`。

###### 配置文件中设置

```ini
socket.ping_interval = 20
````

###### 代码中调用

```go
gsigo.Config.Socket.PingInterval
````

日志配置

- log.hook


`default` `syslog`

###### 日志钩子，默认值 `default`，可自定义钩子。

###### 配置文件中设置

```ini
log.hook = "stdout"
````

###### 代码中调用

```go
gsigo.Config.Log.Hook
````

- log.formatter

`text` `json`

###### 日志输出格式，默认值 `text`。

###### 配置文件中设置

```ini
log.formatter = "text"
````

###### 代码中调用

```go
gsigo.Config.Log.Formatter
````


- log.params

`text` `json`

###### 日志需要的参数，无默认值.

###### 配置文件中设置,syslog例子

```ini
log.params.priority = "LOG_LOCAL0"
log.params.tag = ""
log.params.network = ""
log.params.addr = ""
````

###### 代码中调用

```go
gsigo.Config.Log.params["priority"]
````

### REDIS配置文件

###### 存放路径 项目目录/config/`gsigo.ENV`/redis.ini

```ini
;分组
[redis]

;链接地址
address = 127.0.0.1:6379

;redis密码
password =

;redis库
select = 0

;保持链接时间，单位小时
keep_alive = 10

;连接池，开启链接数量
max_idle = 10

;主
master.address = 127.0.0.1:6379
master.max_idle = 10

;从
slave.max_idle = 10
slave.address[] = 127.0.0.1:6379
slave.address[] = 127.0.0.1:6379
slave.address[] = 127.0.0.1:6379
```

### 数据库配置文件

###### 存放路径 项目目录/config/`gsigo.ENV`/database.ini

```ini
;分组
[english]
;数据库驱动
driver = mysql

;数据库dsn
dsn = root:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local

;打开到数据库的最大连接数。
max_open =  20

;空闲连接池中的最大连接数
max_idle = 10

;可重用连接的最大时间
max_lifetime = 1

;主库
master.dsn = root:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local
master.max_open =  20
master.max_idle = 10

;从库
slave.max_open =  20
slave.max_idle = 10
slave.dsn[] = root:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local
slave.dsn[] = root:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local
slave.dsn[] = root:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local

```

## 路由规则

### WEB路由规则

> [参考gin](https://github.com/gin-gonic/gin)

##### 分组

```go
gsigo.Group(relativePath string, controller ...ControllerInterface) *router
```

###### 示例

```go

package routers

import (
	"github.com/whf-sky/gsigo"
	"test/controllers/web/index"
)

func init()  {
	rootGin := gsigo.Group("/root/")
	{
		rootGin.GET("/", &index.IndexController{})
	}
}
```

##### 使用中间件

```go
gsigo.Use(controller ControllerInterface)
```

##### 静态文件路由规则

```go
gsigo.Static(relativePath string, filePath string)
```

##### POST

```go
gsigo.POST(relativePath string, controller ControllerInterface)
```
##### GET

```go
gsigo.GET(relativePath string, controller ControllerInterface)
```

##### DELETE

```go
gsigo.DELETE(relativePath string, controller ControllerInterface)
```

##### PATCH

```go
gsigo.PATCH(relativePath string, controller ControllerInterface)
```
##### PUT

```go
gsigo.PUT(relativePath string, controller ControllerInterface)
```

##### OPTIONS

```go
gsigo.OPTIONS(relativePath string, controller ControllerInterface)
```

##### HEAD

```go
gsigo.HEAD(relativePath string, controller ControllerInterface)
```

##### Any

`GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE`

```go
gsigo.Any(relativePath string, controller ControllerInterface)
```

### SOCKETIO路由规则

###### 示例

```go
package routers

import (
	"github.com/whf-sky/gsigo"
	"test/controllers/sio/chat"
	"test/controllers/sio/root"
)

func init()  {
	rootRouter := gsigo.Nsp("/")
	{
		rootRouter.OnConnect(&root.ConnectEvent{})
		rootRouter.OnDisconnect(&root.DisconnectEvent{})
		rootRouter.OnError(&root.ErrorEvent{})
		rootRouter.OnEvent("notice", &root.NoticeEvent{})
		rootRouter.OnEvent("bye", &root.ByeEvent{})
	}

	chatRouter := gsigo.Nsp("/chat")
	{
		//如需要ack需要按照如下设置，否则不设置
		chatRouter.OnEvent("msg", &chat.MsgEvent{gsigo.Event{Ack: true},})
	}

}
```

##### Nsp 命名空间相当于WEB组

```go
gsigo.Nsp(nsp string, event ...EventInterface) *router
```

##### OnConnect

```go
gsigo.OnConnect(event EventInterface)
```

##### OnEvent

```go
gsigo.OnEvent(eventName string, event EventInterface)
```


##### OnError

```go
gsigo.OnError(event EventInterface)
```

##### OnDisconnect

```go
gsigo.OnDisconnect(event EventInterface)
```

### CMD路由规则

###### 示例

```go

package routers

import (
	"github.com/whf-sky/gsigo"
	"test/controllers/cmd"
)

func init()  {
	gsigo.CmdRouter(&cmd.TestCmd{})
}

```

```go
gsigo.CmdRouter(cmd CmdInterface) *router
```

## WEB应用

> [参考gin](https://github.com/gin-gonic/gin)

###### 遵循RESTFUL设计风格和开发方式

##### 示例

```go
package index

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
```

##### 结构体定义必须内嵌`gsigo.Controller`结构体 

##### 可定义的Action，参考gin

- `Get()`

- `Post()`

- `Delete()`

- `Put()`

- `Head()`

- `Patch()`

- `Options()` 

- `Any()`  包含请求 GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.

- `Group()` 组

- `Use()` 中间件

- `Prepare() ` 上述方法执行前执行 

- `Finish()` 在执行上述方法之后执行 


##### 可用属性

###### Ctx 用法参考 gin.Context

```go
type Controller struct {
	Ctx  *gin.Context
}
```

##### 可调用方法

###### 获取组名

```go
func (c *Controller) GetGroup() string
```

###### 获取控制器名称

```go
func (c *Controller) GetController() string
```

###### 获取操作名称

```go
func (c *Controller) GetAction() string
```

## SOCKETIO应用

> [参考socket.io](https://github.com/googollee/go-socket.io)

##### 示例

```go
package chat

import (
	"fmt"
	"github.com/whf-sky/gsigo"
)

type MsgEvent struct {
	gsigo.Event
}

func (this *MsgEvent) Execute() {
	this.Conn.Emit("reply", "have "+this.GetMessage())
}
```
##### 结构体定义必须内嵌`gsigo.Event`结构体

##### 可定义的方法

- `Execute` 执行方法

- `Prepare()` 在执行 `Execute` 前执行

- `Finish()` 在执行 `Execute` 后执行

##### 可用属性

###### Ctx 用法参考 gin.Context

```go
type Event struct {
	//是否发送ack
	Ack bool
	//socket 链接用法参考 socketio.Conn
	Conn socketio.Conn
	//事件类型 connect/event/errordisconnect
	EventType string
}
```



##### 可调用方法

###### 绑定业务用户

```go
func (e *Event) SetUser(uid string)
```

###### 获取业务用户

```go
func (e *Event) GetUser() string
```

###### 根据用户获取所有的链接编号,map[Conn.ID()]无意义占位符

```go
func (e *Event) GetCidsByUser(uid string) map[string]int
```

###### 是否ACK消息

```go
func (e *Event) IsAck() bool
```

###### 获取消息

```go
func (e *Event) GetMessage() string
```

###### 设置ACK消息

```go
func (e *Event) SetAckMsg(msg string)
```


###### 获取ACK消息

```go
func (e *Event) GetAckMsg() string
```

###### 获取命名空间

```go
func (e *Event) GetNamespace() string
```

###### 设置错误消息

```go
func (e *Event) SetError(text string)
```

###### 获取错误消息

```go
func (e *Event) GetError() error
```

## CMD应用

##### 示例

```go
package cmd

import (
	"fmt"
	"github.com/whf-sky/gsigo"
)

type TestCmd struct {
	gsigo.Cmd
}

func (this * TestCmd)  Execute(){
	for {
	    fmt.Println("test")	
	}  
}

```
##### 可定义的方法

- `Execute` 执行方法

## 数据库

> [参考gorm](https://gorm.io/docs/index.html)

###### 代码实现了读写分离操作

### CURD

###### 与gorm不同之处增加回调函数

[MODEL详细文档](https://gorm.io/docs/models.html)

######  实例DB

```go
NewDB(gname ...string) *DB 
```

###### 使用配置的组，如不使用`NewDB`需自己实例化使用此方法

```go
func (d *DB) Using(gname ...string) *DB
```

#### 事务

[gorm 事务 文档](https://gorm.io/docs/transactions.html)

###### 回调函数中使用事务

```go
func (d *DB) Transaction (fc func(tx *transaction) error) (err error) 
```

###### 开启事务

```go
func (d *DB) Begin() *transaction 
```

#### 主从强制切换

###### 强制切换到主库

```go
func (d *DB) Master() *DB 
```

###### 强制切换到从库

```go
func (d *DB) Slave() *DB 
```

#### Create

[gorm Create 文档](https://gorm.io/docs/create.html)


###### 示例

```go
type User struct {
  orm.Model
  Name string
  Age  sql.NullInt64 `gorm:"default:18"`
}

user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

db := NewDB("user")

db.Create(&user)

```

###### 插入数据

```go
func (d *DB) Create(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 

//Create别名
func (d *DB) Insert(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

#### Delete

[gorm Delete 文档](https://gorm.io/docs/delete.html)

###### 示例

```go
type User struct {
  orm.Model
  Name string
  Age  sql.NullInt64 `gorm:"default:18"`
}

user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

db := NewDB("user")

// Delete an existing record
db.Delete(&email)
//// DELETE from emails where id=10;

// Add extra SQL option for deleting SQL
db.Delete(&email,func(db *gorm.DB) *gorm.DB {
    db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)")
})
//// DELETE from emails where id=10 OPTION (OPTIMIZE FOR UNKNOWN);

```

###### 删除数据

```go
func (d *DB) Delete(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

#### Update

[gorm Update 文档](https://gorm.io/docs/update.html)

###### 示例

```go
type User struct {
  orm.Model
  Name string
  Age  sql.NullInt64 `gorm:"default:18"`
}

db := NewDB("user")
// Update single attribute if it is changed
db.Update([]string{"name", "hello"}, func(db *gorm.DB) *gorm.DB {
	db.Model(&user)
})
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update single attribute with combined conditions
db.Update([]string{"name", "hello"}, func(db *gorm.DB) *gorm.DB {
	db.Model(&user).Where("active = ?", true)
})
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

// Update multiple attributes with `map`, will only update those changed fields
db.Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false}, func(db *gorm.DB) *gorm.DB {
	db.Model(&user)
})
//// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update multiple attributes with `struct`, will only update those changed & non blank fields
db.Update(User{Name: "hello", Age: 18}, func(db *gorm.DB) *gorm.DB {
	db.Model(&user)
})
//// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// WARNING when update with struct, GORM will only update those fields that with non blank value
// For below Update, nothing will be updated as "", 0, false are blank values of their types
db.Update(User{Name: "", Age: 0, Actived: false}, func(db *gorm.DB) *gorm.DB {
	db.Model(&user)
})
```

###### 改变单个字段

```go
func (d *DB) Update(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 修改多个字段

```go
func (d *DB) Updates(values interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 修改列数据

```go
func (d *DB) UpdateColumn(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

#### Query

[gorm Query 文档](https://gorm.io/docs/query.html)


###### 示例

```go
type User struct {
  orm.Model
  Name string
  Age  sql.NullInt64 `gorm:"default:18"`
}

user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

db := NewDB("user")

// Get first matched record
db.First(&user, func(db *gorm.DB) *gorm.DB {
    db.Where("name = ?", "jinzhu")
})
//// SELECT * FROM users WHERE name = 'jinzhu' limit 1;

// Get all matched records
db.Find(&users, func(db *gorm.DB) *gorm.DB {
   db.Where("name = ?", "jinzhu")
})
//// SELECT * FROM users WHERE name = 'jinzhu';

// <>
db.Find(&users, func(db *gorm.DB) *gorm.DB {
  db.Where("name <> ?", "jinzhu")
})
//// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Find(&users, func(db *gorm.DB) *gorm.DB {
  db.Where("name IN (?)", []string{"jinzhu", "jinzhu 2"})
})
//// SELECT * FROM users WHERE name in ('jinzhu','jinzhu 2');

// LIKE
db.Find(&users, func(db *gorm.DB) *gorm.DB {
  db.Where("name LIKE ?", "%jin%")
})
//// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Find(&users, func(db *gorm.DB) *gorm.DB {
   db.Where("name = ? AND age >= ?", "jinzhu", "22")
})
//// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

// Time
db.Find(&users, func(db *gorm.DB) *gorm.DB {
  db.Where("updated_at > ?", lastWeek)
})
//// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Find(&users, func(db *gorm.DB) *gorm.DB {
 db.Where("created_at BETWEEN ? AND ?", lastWeek, today)
})
//// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';

```

###### 查询第一条数据，按主键正序排序
```go
func (d *DB) First(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取一条记录，没有指定的顺序

```go
func (d *DB) Take(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取最后一条数据，按照主键倒叙排序
```go
func (d *DB) Last(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取多条数据
```go
func (d *DB) Find(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取第一个匹配的记录，或者在给定条件下初始化一个新记录(只适用于结构，映射条件)
```go
func (d *DB) FirstOrInit(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取第一个匹配的记录，或者在给定的条件下创建一个新的记录(只适用于struct, map条件)
```go
func (d *DB) FirstOrCreate(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 获取一个模型有多少条记录
```go
func (d *DB) Count(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

###### 从模型中查询单个列作为映射，如果您想要查询多个列，则应该使用Scan
```go
func (d *DB) Pluck(column string, value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```
#### 原生sql

[gorm raw 文档](https://gorm.io/docs/query.html)

###### 将结果扫描到另一个结构中。

```go
func (d *DB) Scan(dest interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB
```

###### 运行原始SQL,单条查询，它不能与其他方法链接
```go
func (d *DB) Row(funcs ...func(db *gorm.DB) *gorm.DB ) *sql.Row
```

###### 运行原始SQL,多条查询，它不能与其他方法链接
```go
func (d *DB) Rows(funcs ...func(db *gorm.DB) *gorm.DB ) (*sql.Rows, error)
```

###### 运行原始SQL,将结果扫描到另一个结构中。
```go
func (d *DB) Raw(sql string, values ...interface{}) *gorm.DB
```

###### 运行原始SQL,执行影响操作的SQL
```go
func (d *DB) Exec(sql string, values ...interface{}) *gorm.DB
```

### MODEL

[gorm Models 详细文档](https://gorm.io/docs/models.html)

```go
type Test struct {
    Id int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
}

func (w *Words) BeforeCreate(scope *gorm.Scope) error {
    scope.SetColumn("CreateTime", time.Now())
    return nil
}
```

## REDIS

> [参考redis](https://github.com/gomodule/redigo/redis)

###### 代码实现了读写分离操作

######  实例redis

```go
func NewRedis(gname ...string) *Redis
```

###### 使用的redis配置的组，如不使用`redis.NewRedis`需自己实例化使用此方法

```go
func (d *Redis) Using(gname ...string) *Redis 
```

###### 强制切换到主库

```go
func (r *Redis) Master() *Redis
```

###### 强制切换到从库

```go
func (r *Redis) Slave() *Redis{
```

###### 执行命令

```go
func (r *Redis) Do(cmd string, args ...interface{}) (reply interface{}, err error)
```

###### 将命令写入客户机的输出缓冲区

```go
func (r *Redis) Send(commandName string, args ...interface{}) error
```

###### 将输出缓冲区刷新到Redis服务器。

```go
func (r *Redis) Flush() error
```

###### 接收来自Redis服务器的单个回复

```go
func (r *Redis) Receive() (reply interface{}, err error) 
```


###### 发布订阅

```go
func (r *Redis) PubSub() redis.PubSubConn
```

###### 返回一个新的脚本对象

```go
func (r *Redis) Script(keyCount int, src string) *script
```

## 日志

> [参考logrus](https://github.com/sirupsen/logrus)

## 环境变量

###### 环境变量的使用示例

```sh
$ go run main.go -env=develop
````

```sh
$ export gsigo_env=develop
````