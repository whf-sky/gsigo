package orm

import (
	"github.com/jinzhu/gorm"
)

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
	db, err = gorm.Open(driver, dsn)
	return
}