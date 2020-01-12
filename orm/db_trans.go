package orm

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

// Transaction
type transaction struct {
	tx *gorm.DB
}

// Create insert the value into database
func (t *transaction) Insert(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Create(value)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
// WARNING If model has DeletedAt field, GORM will only set field DeletedAt's value to current time
func (t *transaction) Delete(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Delete(value)
}

// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
// WARNING when update with struct, GORM will not update fields that with zero value
func (t *transaction) Update(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Update(attrs...)
}

// Updates update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (t *transaction) Updates(values interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Updates(values)
}

// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (t *transaction) UpdateColumn(attrs []interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.UpdateColumn(attrs...)
}

// First find first record that match given conditions, order by primary key
func (t *transaction) First(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.First(out)
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (t *transaction) Take(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Take(out)
}

// Last find last record that match given conditions, order by primary key
func (t *transaction) Last(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Last(out)
}

// Find find records that match given conditions
func (t *transaction) Find(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Find(out)
}

// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorinit
func (t *transaction) FirstOrInit(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.FirstOrInit(out)
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (t *transaction) FirstOrCreate(out interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.FirstOrCreate(out)
}

// Count get how many records for a model
func (t *transaction) Count(value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Count(value)
}

// Row return `*sql.Row` with given conditions
func (t *transaction) Row(funcs ...func(db *gorm.DB) *gorm.DB ) *sql.Row {
	return t.tx.Row()
}

// Rows return `*sql.Rows` with given conditions
func (t *transaction) Rows(funcs ...func(db *gorm.DB) *gorm.DB ) (*sql.Rows, error) {
	return t.tx.Rows()
}

//var ages []int64
//db.Find(&users).Pluck("age", &ages)
func (t *transaction) Pluck(column string, value interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Pluck(column, value)
}

// Scan scan value to a struct
func (t *transaction) Scan(dest interface{}, funcs ...func(db *gorm.DB) *gorm.DB ) *gorm.DB {
	return t.tx.Scan(dest)
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func (t *transaction) Raw(sql string, values ...interface{}) *gorm.DB {
	return t.tx.Raw(sql)
}

// Exec execute raw sql
func (t *transaction) Exec(sql string, values ...interface{}) *gorm.DB {
	return t.tx.Exec(sql, values...)
}