package orm

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

const (
	ModeDefault = iota
	ModeMaster
	ModeSlave
)

func NewDB(gname ...string) *DB  {
	return (&DB{}).Using(gname...)
}

type DB struct {
	db 			*gorm.DB
	groups 		map[string]*group
	gnames 		[]string
	mode 		int
	curSlave 	int
}

//open database group
func (d *DB) Using(gname ...string) *DB {
	d.gnames = gname
	d.groups = using(gname...)
	return d
}

// Transaction start a transaction as a block,
// return error will rollback, otherwise to commit.
func (d *DB) Transaction (fc func(tx *transaction) error) (err error) {
	return d.Db(ModeMaster).Transaction(func(tx *gorm.DB) error {
		return fc(&transaction{tx:tx})
	})
}

// Begin begins a transaction
func (d *DB) Begin() *transaction {
	tx := d.Db(ModeMaster).Begin()
	return &transaction{tx:tx}
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

// Create insert the value into database
func (d *DB) Create(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Create(value)
}

// Create insert the value into database
func (d *DB) Insert(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Create(value, funcs...)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
// WARNING If model has DeletedAt field, GORM will only set field DeletedAt's value to current time
func (d *DB) Delete(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Delete(value)
}

// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
// WARNING when update with struct, GORM will not update fields that with zero value
func (d *DB) Update(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Update(attrs...)
}

// Updates update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (d *DB) Updates(values interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).Updates(values)
}

// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (d *DB) UpdateColumn(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeMaster, funcs...).UpdateColumn(attrs...)
}

// First find first record that match given conditions, order by primary key
func (d *DB) First(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).First(out)
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (d *DB) Take(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Take(out)
}

// Last find last record that match given conditions, order by primary key
func (d *DB) Last(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Last(out)
}

// Find find records that match given conditions
func (d *DB) Find(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Find(out)
}

// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorinit
func (d *DB) FirstOrInit(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).FirstOrInit(out)
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (d *DB) FirstOrCreate(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).FirstOrCreate(out)
}

// Count get how many records for a model
func (d *DB) Count(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Count(value)
}

// Row return `*sql.Row` with given conditions
func (d *DB) Row(funcs ...func(db *gorm.DB) *gorm.DB ) *sql.Row {
	return d.Db(ModeSlave, funcs...).Row()
}

// Rows return `*sql.Rows` with given conditions
func (d *DB) Rows(funcs ...func(db *gorm.DB) *gorm.DB ) (*sql.Rows, error) {
	return d.Db(ModeSlave, funcs...).Rows()
}

//var ages []int64
//db.Find(&users).Pluck("age", &ages)
func (d *DB) Pluck(column string, value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Pluck(column, value)
}

// Scan scan value to a struct
func (d *DB) Scan(dest interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return d.Db(ModeSlave, funcs...).Scan(dest)
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func (d *DB) Raw(sql string, values ...interface{}) *gorm.DB {
	return d.Db(ModeSlave).Raw(sql)
}

// Exec execute raw sql
func (d *DB) Exec(sql string, values ...interface{}) *gorm.DB {
	return d.Db(ModeMaster).Exec(sql, values...)
}

// get a opened database specified by its database driver name and a
// driver-specific data source name for database group, usually consisting of at least a
// database name and connection information.
func (d *DB) Db(mode int, funcs ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	defer func() {
		d.mode = ModeDefault
	}()
	//set mode
	if d.mode != ModeDefault {
		mode = d.mode
	}

	//run functions
	runFuncs := func() *gorm.DB {
		if len(funcs) > 0 {
			for _, fun := range funcs  {
				d.db = fun(d.db)
			}
		}
		return d.db
	}

	//get group
	group, ok := d.groups[d.gnames[0]]
	if !ok {
		panic("this dao without using group")
	}

	//get gorm Db without set read write database
	if group.Master == nil {
		d.db = group.Db
		return runFuncs()
	}
	//get master gorm Db
	if mode == ModeMaster {
		d.db = group.Master
		return runFuncs()
	}

	//get slave  gorm Db
	sCnt := len(group.Slave)
	if d.curSlave > sCnt-1 {
		d.curSlave = 0
	}
	d.db = group.Slave[d.curSlave]
	d.curSlave++

	return runFuncs()
}