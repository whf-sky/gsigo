package orm

import (
	"github.com/jinzhu/gorm"
)

//启动所有的组的数据库连接
func OpenGroup(configRead func(out interface{}) error) (map[string]*group, error) {
	//配置信息
	var configs map[string]dbGroupsCnf
	//获取配置文件配置信息
	err := configRead(&configs)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, nil
	}
	//数据库连接组
	groups := map[string]*group{}
	//启动组的数据连接
	for gname, config := range configs {
		groups[gname] = newGroup().openGroup(config)
	}
	return groups, nil
}

// Open initialize a new db connection, need to import driver first, e.g:
//
//     import _ "github.com/go-sql-driver/mysql"
//	   import "github.com/whf-sky/gsigo/orm"
//     func main() {
//       db, err := orm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
//     }
// gsigo orm has wrapped some drivers, for easier to remember driver's import path, so you could import the mysql driver with
//    import _ "github.com/jinzhu/gorm/dialects/mysql"
//    // import _ "github.com/jinzhu/gorm/dialects/postgres"
//    // import _ "github.com/jinzhu/gorm/dialects/sqlite"
//    // import _ "github.com/jinzhu/gorm/dialects/mssql"
func Open(driver string, dsn string) (db *gorm.DB, err error){
	// Open initialize a new db connection
	return gorm.Open(driver, dsn)
}