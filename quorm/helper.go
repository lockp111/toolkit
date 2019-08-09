package quorm

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
	"github.com/jinzhu/gorm"
)

var db = new(DataBase)

// NewDB ...
func NewDB() *DataBase {
	return &DataBase{}
}

// Open ...
func Open(debug bool, dialect, dburl string) error {
	return db.Connect(debug, dialect, dburl)
}

// DB ...
func DB() *DataBase {
	return db
}

// Gorm ...
func Gorm() *gorm.DB {
	return db.DB
}

// Goqu ...
func Goqu() *goqu.Database {
	return db.Goqu
}

// DebugSQL ...
func DebugSQL(sqlBuilder *goqu.SelectDataset) string {
	sql, args, err := sqlBuilder.Prepared(true).ToSQL()
	return fmt.Sprint("SQL:", sql, args, err)
}

// StringOutRows ...
func StringOutRows(rows []string, pri ...string) []exp.AliasedExpression {
	var goquout = make([]exp.AliasedExpression, 0, len(rows))
	for _, row := range rows {
		if len(pri) != 0 {
			goquout = append(goquout, goqu.I(pri[0]+"."+row).As(row))
		} else {
			goquout = append(goquout, goqu.I(row).As(row))
		}
	}

	return goquout
}

// Now ..
func Now() *time.Time {
	t := time.Now()
	return &t
}

// RecordCount ...
func RecordCount(tx *gorm.DB) (int64, error) {
	c := &struct {
		Size int64
	}{}

	err := tx.Select("count(*) as `size`").Scan(c).Error
	return c.Size, err
}

// Transaction ...
func Transaction(f func(*gorm.DB) error) error {
	return db.Transaction(f)
}

// Query ...
func Query(scaner *gorm.DB, query *goqu.SelectDataset,
	outRows interface{}, selectEx ...interface{}) error {

	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = scaner.Raw(sql, args...).Find(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

// QueryFirst ..
func QueryFirst(scaner *gorm.DB, query *goqu.SelectDataset,
	outRow interface{}, selectEx ...interface{}) error {

	selectQuery := query
	if selectEx != nil {
		query = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = scaner.Raw(sql, args...).Limit(1).Find(outRow).Error
	if err != nil {
		return err
	}

	return nil
}

// QueryRows ..
func QueryRows(scaner *gorm.DB, query *goqu.SelectDataset,
	selectEx ...interface{}) (*sql.Rows, error) {

	selectQuery := query
	if selectEx != nil {
		selectQuery = selectQuery.Select(selectEx...)
	}
	sql, args, err := selectQuery.Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scaner.Raw(sql, args...).Rows()
}

// Exec ...
func Exec(tx *gorm.DB, update *goqu.UpdateDataset) (rowsAffected int64, err error) {
	sql, args, err := update.Prepared(true).ToSQL()
	if err != nil {
		return 0, err
	}

	ret := tx.Exec(sql, args)
	return ret.RowsAffected, ret.Error
}

// QueryCount ...
func QueryCount(query *goqu.SelectDataset, selectEx ...interface{},
) (int64, error) {
	return db.QueryCount(query, selectEx...)
}

// PageQuery ...
func PageQuery(scaner *gorm.DB, query *goqu.SelectDataset, pageIndex int64,
	pageSize int64, outRows interface{}, selectEx ...interface{}) (int64, error) {
	return db.PageQuery(scaner, query, pageIndex, pageSize, outRows, selectEx...)
}

// QueryScan ...
func QueryScan(query *goqu.SelectDataset, outRows interface{},
	selectEx ...interface{}) error {
	return db.QueryScan(query, outRows, selectEx...)
}
