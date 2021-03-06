package orm

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

const (
	ModeDefault = iota
	ModeMaster
	ModeSlave
)

func NewDB() *DB  {
	return &DB{
		trans:		false,
		validate:	validator.New(),
	}
}

type DB struct {
	db 			*gorm.DB
	groups 		map[string]*group
	group 		*group
	mode 		int
	curSlave 	int
	trans 		bool
	validate   *validator.Validate
}

//设置组信息
func (d *DB) SetGroups(groups map[string]*group) *DB {
	d.groups = groups
	return d
}

//设置gorm.db
func (d *DB) SetDb(db *gorm.DB) *DB {
	d.db = db
	return d
}

//open database group
func (d *DB) Using(gname string) *DB {
	if gname == "" {
		panic("No gname information")
	}
	group, ok := d.groups[gname];
	if !ok {
		panic("database group not exist")
	}
	d.group = group
	return d
}

//curd for  master database connect
func (d *DB) Master() *DB {
	d.mode = ModeMaster
	return d
}

//curd for slave database connect
func (d *DB) Slave() *DB {
	d.mode = ModeSlave
	return d
}

// Transaction start a transaction as a block,
// return error will rollback, otherwise to commit.
func (d *DB) Transaction (fc func(tx *DB) error) (err error) {
	panicked := true
	tx := d.Begin()
	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if panicked || err != nil {
			tx.Rollback()
		}
	}()
	err = fc(tx)
	if err == nil {
		err = tx.Commit().Error
	}
	panicked = false
	return
}

// Begin begins a transaction
func (d *DB) Begin() *DB {
	d.Db(ModeMaster)
	c := d.clone()
	c.db = c.db.Begin()
	return c
}

// Commit commit a transaction
func (d *DB) Commit() *gorm.DB {
	return d.Db(ModeMaster).Commit()
}

// Rollback rollback a transaction
func (d *DB) Rollback() *gorm.DB {
	return d.Db(ModeMaster).Rollback()
}

type gormFunc func(db *gorm.DB) *gorm.DB

// Create insert the value into database
func (d *DB) Create(value interface{}, funcs ...gormFunc ) (*gorm.DB, error) {
	err := d.validate.Struct(value)
	if err != nil {
		return nil, err
	}
	return d.Db(ModeMaster, funcs...).Create(value), nil
}

// Create insert the value into database
func (d *DB) Insert(value interface{}, funcs ...gormFunc ) (*gorm.DB, error) {
	return d.Create(value, funcs...)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
// WARNING If model has DeletedAt field, GORM will only set field DeletedAt's value to current time
func (d *DB) Delete(value interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Delete(value)
}

// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
// WARNING when update with struct, GORM will not update fields that with zero value
func (d *DB) Update(model interface{}, attrs []interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Model(model).Update(attrs...)
}

// Updates update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (d *DB) Updates(model interface{}, values interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Model(model).Updates(values)
}

// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (d *DB) UpdateColumn(model interface{}, attrs []interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Model(model).UpdateColumn(attrs...)
}

// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (d *DB) UpdateColumns(model interface{}, values interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Model(model).UpdateColumns(values)
}

// First find first record that match given conditions, order by primary key
func (d *DB) First(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).First(out)
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (d *DB) Take(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Take(out)
}

// Last find last record that match given conditions, order by primary key
func (d *DB) Last(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Last(out)
}

// Find find records that match given conditions
func (d *DB) Find(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Find(out)
}

// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorinit
func (d *DB) FirstOrInit(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).FirstOrInit(out)
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (d *DB) FirstOrCreate(out interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).FirstOrCreate(out)
}

// Count get how many records for a model
func (d *DB) Count(model interface{}, value interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Model(model).Count(value)
}

// Row return `*sql.Row` with given conditions
func (d *DB) Row(funcs ...gormFunc ) *sql.Row {
	return d.Db(ModeSlave, funcs...).Row()
}

// Rows return `*sql.Rows` with given conditions
func (d *DB) Rows(funcs ...gormFunc ) (*sql.Rows, error) {
	return d.Db(ModeSlave, funcs...).Rows()
}

//var ages []int64
//db.Find(&users).Pluck("age", &ages)
func (d *DB) Pluck(model interface{}, column string, value interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Model(model).Pluck(column, value)
}

// Scan scan value to a struct
func (d *DB) Scan(dest interface{}, funcs ...gormFunc ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Scan(dest)
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func (d *DB) Raw(sql string, values ...interface{}) *gorm.DB {
	return d.Db(ModeSlave).Raw(sql, values...)
}

// Exec execute raw sql
func (d *DB) Exec(sql string, values ...interface{}) *gorm.DB {
	return d.Db(ModeMaster).Exec(sql, values...)
}

// get a opened database specified by its database driver name and a
// driver-specific data source name for database group, usually consisting of at least a
// database name and connection information.
func (d *DB) Db(mode int, funcs ...gormFunc) *gorm.DB {
	//当是事务的时候直接返回
	if d.trans {
		return d.db
	}
	//当db被直接设置时使用
	if d.db != nil {
		return d.db
	}
	//还原mode为默认值
	defer func() {
		d.db = nil
		d.mode = ModeDefault
	}()
	//当DB mode 不等于默认值时对BD mode进行修改
	if d.mode != ModeDefault {
		mode = d.mode
	}
	//当数据连接没有配置主库的时候使用默认数据库连接
	if d.group.Master == nil {
		d.db = d.group.Db
		return d.callback(funcs)
	}
	//当存在主库配置并且mode值为主库时使用主数据库连接
	if mode == ModeMaster {
		d.db = d.group.Master
		return d.callback(funcs)
	}
	//使用从库的数据连接
	sCnt := len(d.group.Slave)
	if d.curSlave > sCnt-1 {
		d.curSlave = 0
	}
	d.db = d.group.Slave[d.curSlave]
	d.curSlave++

	return d.callback(funcs)
}

//clone a DB
func  (d *DB)  clone() *DB {
	return &DB{
		db:       	d.db,
		group:   	d.group,
		mode:     	d.mode,
		curSlave: 	d.curSlave,
		trans:		true,
	}
}

//数据库操作的回调函数
func (d *DB) callback(funcs []gormFunc) *gorm.DB {
	if len(funcs) > 0 {
		for _, fun := range funcs  {
			d.db = fun(d.db)
		}
	}
	return d.db
}

