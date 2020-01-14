module github.com/whf-sky/gsigo

go 1.12

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20191210023423-ac6580df4449

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20191217033636-bbbf87ae2631

replace golang.org/x/mod => github.com/golang/mod v0.1.0

replace golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543

replace golang.org/x/net => github.com/golang/net v0.0.0-20191209160850-c0dbc17a3553

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191206172530-e9b2fee46413

replace golang.org/x/text => github.com/golang/text v0.3.2

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/googollee/go-engine.io v1.4.2
	github.com/googollee/go-socket.io v1.4.2
	github.com/sirupsen/logrus v1.4.2
	github.com/tabalt/gracehttp v1.3.0
	gopkg.in/yaml.v2 v2.2.7
)
