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
        - [主从强制切换](#主从强制切换)
        - [事务](#事务)
        - [Create](#Create)
        - [Delete](#Delete)
        - [Update](#Update)
        - [Query](#Query)
    - [MODEL](#MODEL)
- [REDIS](#REDIS)
- [日志](#日志)
- [工具包](#工具包)
    - [验证器](#验证器)
    - [Curl](#Curl)
    - [ApiResult](#ApiResult)
    - [一致性hash](#一致性hash)
    - [登录](#登录)
    - [接口签名](#接口签名)
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
    gsigo.Run("./config/app.ini")
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

##### 示例

```go

package routers

import (
	"github.com/whf-sky/gsigo"
	"test/controllers/cmd"
)

func init()  {
	gsigo.Cmd("test", &cmd.TestCmdController{})
}

```

##### 路由

###### Cmd

```go
gsigo.func Cmd(requestUri string, cmd CmdInterface)
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

##### 执行

```go
go run cmd.go  -request_uri=requestUri
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


#### 主从强制切换

###### 强制切换到主库

```go
func (d *DB) Master() *DB 
```

###### 强制切换到从库

```go
func (d *DB) Slave() *DB 
```


#### 事务

[gorm 事务 文档](https://gorm.io/docs/transactions.html)

##### 回调函数中使用事务

```go
func (d *DB) Transaction (fc func(tx *DB) error) (err error)
```

###### 示例

```go
db.Transaction(func(tx *DB) error {
    // do some database operations in the transaction (use 'tx' from this point, not 'db')
    if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
      // return any error will rollback
      return err
    }

    if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
      return err
    }

    // return nil will commit
    return nil
  })
```


##### 开启事务

```go
func (d *DB) Begin() *DB
```

##### 提交事务

```go
func (d *DB) Commit() *gorm.DB
```

##### 回滚事务

```go
func (d *DB) Rollback() *gorm.DB
```

##### 事务示例

```go
func CreateAnimals(db *DB) error {
  // Note the use of tx as the database handle once you are within a transaction
  tx := db.Begin()
  defer func() {
    if r := recover(); r != nil {
      tx.Rollback()
    }
  }()

  if err := tx.Error; err != nil {
    return err
  }

  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
     tx.Rollback()
     return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
     tx.Rollback()
     return err
  }

  return tx.Commit().Error
}
```

#### Create

[gorm Create 文档](https://gorm.io/docs/create.html)

##### 方法

```go
func (d *DB) Create(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) (*gorm.DB, error) 
```

##### 别名

```go
func (d *DB) Insert(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB )  (*gorm.DB, error) 
```

[error 参见文档](https://github.com/go-playground/validator)


##### 示例

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

#### Delete

[gorm Delete 文档](https://gorm.io/docs/delete.html)

##### 方法

```go
func (d *DB) Delete(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB
```

##### 示例

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

#### Update

[gorm Update 文档](https://gorm.io/docs/update.html)


##### 改变单个字段

```go
func (d *DB) Update(model interface{}, attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 修改多个字段

```go
func (d *DB) Updates(model interface{}, values interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 修改列数据

```go
func (d *DB) UpdateColumn(model interface{}, attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 修改多列数据

```go
func (d *DB) UpdateColumns(model interface{}, values interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB```

##### 示例

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

#### Query

[gorm Query 文档](https://gorm.io/docs/query.html)

##### 查询第一条数据，按主键正序排序
```go
func (d *DB) First(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取一条记录，没有指定的顺序

```go
func (d *DB) Take(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取最后一条数据，按照主键倒叙排序
```go
func (d *DB) Last(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取多条数据
```go
func (d *DB) Find(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取第一个匹配的记录，或者在给定条件下初始化一个新记录(只适用于结构，映射条件)
```go
func (d *DB) FirstOrInit(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取第一个匹配的记录，或者在给定的条件下创建一个新的记录(只适用于struct, map条件)
```go
func (d *DB) FirstOrCreate(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 获取一个模型有多少条记录
```go
func (d *DB) Count(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 从模型中查询单个列作为映射，如果您想要查询多个列，则应该使用Scan
```go
func (d *DB) Pluck(column string, value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB 
```

##### 示例

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
#### 原生sql

[gorm raw 文档](https://gorm.io/docs/query.html)

##### 将结果扫描到另一个结构中。

```go
func (d *DB) Scan(dest interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB
```

##### 运行原始SQL,单条查询，它不能与其他方法链接
```go
func (d *DB) Row(funcs ...func(db *gorm.DB) *gorm.DB ) *sql.Row
```

##### 运行原始SQL,多条查询，它不能与其他方法链接
```go
func (d *DB) Rows(funcs ...func(db *gorm.DB) *gorm.DB ) (*sql.Rows, error)
```

##### 运行原始SQL,将结果扫描到另一个结构中。
```go
func (d *DB) Raw(sql string, values ...interface{}) *gorm.DB
```

##### 运行原始SQL,执行影响操作的SQL
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

##### 代码实现了读写分离操作

#####  实例redis

```go
func NewRedis(gname ...string) *Redis
```

##### 使用的redis配置的组，如不使用`redis.NewRedis`需自己实例化使用此方法

```go
func (d *Redis) Using(gname ...string) *Redis 
```

##### 强制切换到主库

```go
func (r *Redis) Master() *Redis
```

##### 强制切换到从库

```go
func (r *Redis) Slave() *Redis{
```

##### 执行命令

```go
func (r *Redis) Do(cmd string, args ...interface{}) (reply interface{}, err error)
```

##### 将命令写入客户机的输出缓冲区

```go
func (r *Redis) Send(commandName string, args ...interface{}) error
```

##### 将输出缓冲区刷新到Redis服务器。

```go
func (r *Redis) Flush() error
```

##### 接收来自Redis服务器的单个回复

```go
func (r *Redis) Receive() (reply interface{}, err error) 
```


##### 发布订阅

```go
func (r *Redis) PubSub() redis.PubSubConn
```

##### 返回一个新的脚本对象

```go
func (r *Redis) Script(keyCount int, src string) *script
```

## 日志

> [参考logrus](https://github.com/sirupsen/logrus)

```go
gsigo.Log.
```

##### 示例

```go
gsigo.Log.Trace("this is test!")
```

## 工具包

### 验证器

[参见文档](https://github.com/go-playground/validator)

### Curl

#### 示例

```go
package main

import "github.com/whf-sky/gsigo/utils"


func main()  {
    utils.NewCurl("https://www.baidu.com").Get()
}

```

#### 方法

##### 实例化

###### NewCurl

```go
func NewCurl(url string) *Curl
```

##### 选项 

###### Url

```go
func (c *Curl) Url(url string) *Curl
```

###### Query

```go
func  (c *Curl) Query(query map[string]string) *Curl 
```

###### SetQuery

```go
func (c *Curl) SetQuery(name string, value string) *Curl 
```

###### Body

```go
func (c *Curl) Body(body io.Reader) *Curl 
```

###### Cookie

```go
func (c *Curl) Cookie(cookie map[string]string) *Curl
```

###### SetCookie

```go
func (c *Curl) SetCookie(name string, value string) *Curl
```

###### Header

```go
func (c *Curl) Header(header map[string]string) *Curl 
```

###### SetHeader

```go
func (c *Curl) SetHeader(name string, value string) *Curl 
```

###### SetPostForm

```go
func (c *Curl) SetPostForm() *Curl
```

###### GetBody

```go
func (c *Curl) GetBody(method string) ([]byte, error) 
```

###### GetHeader

```go
func (c *Curl) GetHeader() http.Header
```

##### 请求方法

###### Get

```go
func (c *Curl) Get() ([]byte, error)
```

###### Post

```go
func (c *Curl) Post() ([]byte, error)
```
###### Options

```go
func (c *Curl) Options() ([]byte, error)
```

###### Head

```go
func (c *Curl) Head() ([]byte, error)
```

###### Put

```go
func (c *Curl) Put() ([]byte, error)
```

###### Delete

```go
func (c *Curl) Delete() ([]byte, error)
```

###### Patch

```go
func (c *Curl) Patch() ([]byte, error)
```

###### Connect

```go
func (c *Curl) Connect() ([]byte, error)
```


### ApiResult

#### 示例

```go
package main

import "github.com/whf-sky/gsigo/utils"


func main()  {
    utils.NewApiResult().SetError("类型错误", 100)
}

```

#### 方法

###### 实例化

```go
func NewApiResult() *ApiResult
```

###### 设置状态码

```go
func (a *ApiResult) SetCode(code int) map[string]interface{}
```

###### 设置成功信息

```go
func (a *ApiResult) SetSuccess(data interface{}) map[string]interface{}
```

###### 设置错误信息

```go
func (a *ApiResult) SetError(msg string, code int) map[string]interface{} 
```

### 一致性hash

#### 示例

```go
package main

import (
	"fmt"
	"github.com/whf-sky/gsigo/utils"
)

func main(){

	hash := utils.NewConsistentHashing(10)
	//添加节点
	hash.AddNode("test1")
	hash.AddNode("test2")

	//获取落点
	fmt.Println(hash.GetLocation("我在哪"))
	fmt.Println(hash.GetLocation("我是谁"))
    //删除落点
	hash.DeleteNode("test1")
}
```


#### 方法

###### 实例化

```go
func NewConsistentHashing(virtualNodeNum int) *ConsistentHashing
```

###### SetVirtualNodeNum

```go
//设置虚拟节点数量
//num 节点数
func (c *ConsistentHashing) SetVirtualNodeNum (num int) 
```


###### AddNode

```go
//添加节点
//node 节点
func (c *ConsistentHashing) AddNode (node string) 
```

###### GetLocation

```go
//寻找字符串所在位置
//str 字符串
func (c *ConsistentHashing) GetLocation (str string) string
```

###### DeleteNode

```go
//删除一个节点
//node 节点
func (c *ConsistentHashing) DeleteNode (node string) 
```


### 登录

#### 示例

```go
package index

import (
	"github.com/whf-sky/gsigo"
	"github.com/whf-sky/gsigo/utils"
)

type IndexController struct {
	gsigo.Controller
}

func (this *IndexController) Get() {
	login := utils.NewLogin(this.Ctx)
	login.SetExpire(0)
	login.SetName("login")
	login.SetVerifyExpire(false)
	login.SetValue("1")

	val, err := login.GetValue()
	if err != nil {
		this.Ctx.String(200, "%s" , err.Error())
		return
	}
	this.Ctx.String(200,  "%s", val)
}
```

#### 方法

###### 实例化

```go
func NewLogin(ctx *gin.Context) *Login
```

###### SetCtx

```go
func (l *Login) SetCtx(ctx *gin.Context) *Login
```


###### SetVerifyExpire

```go
//设置cookie 的有效期。
//verify true:验证/false:不验证
//默认验证
func (l *Login) SetVerifyExpire(verify bool) *Login
```

###### SetExpire

```go
//设置cookie 的有时间。
//expire 秒
func (l *Login) SetExpire(maxAge int) *Login
```

###### SetPath

```go
//设置cookie路径。
func (l *Login) SetPath(path string) *Login
```

###### SetDomain

```go
//设置cookie域名
func (l *Login) SetDomain(domain string) *Login
```

###### Delete

```go
//删除cookie
func (l *Login) Delete() *Login
```
###### SetValue

```go
//设置cookie值
//使用此方法会最终生成cookie，放在最后调用
func (l *Login) SetValue(value string) *Login
```
###### GetValue

```go
//获取cookie
func  (l *Login) GetValue() (string, error)
```

### 接口签名

#### 验证示例

```go
package index

import (
	"github.com/whf-sky/gsigo"
	"github.com/whf-sky/gsigo/utils"
)

type IndexController struct {
	gsigo.Controller
}

func (this *IndexController) Get() {
	sign := utils.NewSign()
	sign.SetSecret("!@!@#%^&&*(())hhjkk")
	sign.SetApiKey("saas-1")
	sign.SetCtx(this.Ctx)
	sign.SetOpenExpired(true)
	sign.SetExpired(1000)
	err := sign.Verify("我是谁？")
	if err != nil {
		this.Ctx.String(200, "%s" , err.Error())
		return
	}
	this.Ctx.String(200, "%s" , "验证成功")
}
```
#### 签名生成示例

```go
package index

import (
	"github.com/whf-sky/gsigo"
	"github.com/whf-sky/gsigo/utils"
)

type IndexController struct {
	gsigo.Controller
}

func (this *IndexController) Get() {
    sign := utils.NewSign()
    sign.SetSecret("!@!@#%^&&*(())hhjkk")
    apikey := "saas-1"
    nonce := fmt.Sprintf("%f", rand.Float64())
    timestamp := fmt.Sprintf("%d",time.Now().Unix())
    signStr := sign.GenerateSignature(apikey, nonce, timestamp, "我是谁？")

    fmt.Println("HTTP_API_KEY:",apikey)
    fmt.Println("HTTP_API_NONCE:",nonce)
    fmt.Println("HTTP_API_TIMESTAMP:",timestamp)
    fmt.Println("HTTP_API_SIGNATURE:",signStr)
}
```
#### 头信息

- HTTP_API_KEY

###### 随机字符串

- HTTP_API_NONCE

###### 时间戳

- HTTP_API_TIMESTAMP

###### 签名

- HTTP_API_SIGNATURE

#### 方法

###### 实例化

```go
func NewSign() *Sign
```

###### SetCtx

```go
//设置 gin Context
func (s *Sign) SetCtx(ctx *gin.Context) *Sign
```


###### SetSecret

```go
//设置秘钥
func (s *Sign) SetSecret(secret string) *Sign
```

###### SetApiKey

```go
//设置应用key
func (s *Sign) SetApiKey(apiKey string) *Sign
```

###### SetExpired

```go
//设置签名有效期
//expired 秒
func (s *Sign) SetExpired(expired int) *Sign
```

###### SetOpenExpired

```go
//设置签名有效期
func (s *Sign) SetOpenExpired(open bool) *Sign
```

###### Verify

```go
//验证签名
func (s *Sign) Verify(str string) error
```

###### GenerateSignature

```go
//生成签名
func (s *Sign) GenerateSignature(apikey, nonce, timestamp, str string) string
```

## 环境变量

##### 环境变量的使用示例


```sh
$ export GSIGO_ENV=develop
````

或

```sh
$ go run main.go -env=develop
````