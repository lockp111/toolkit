package quorm

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	// mysql drive
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// DataBase ...
type DataBase struct {
	*gorm.DB
	Goqu *goqu.Database
}

// Connect ...
func (db *DataBase) Connect(debug bool, dialect, dburl string) error {
	gdb, err := gorm.Open(dialect, dburl)
	if err != nil {
		return err
	}

	gdb.SingularTable(true)
	if debug {
		gdb.LogMode(true)
	}

	db.DB = gdb
	db.Goqu = goqu.New(dialect, gdb.DB())
	return nil
}

// PageQuery ...
func (db *DataBase) PageQuery(scaner *gorm.DB, query *goqu.SelectDataset, pageIndex, pageSize int64,
	outRows interface{}) (int64, error) {

	count, err := db.QueryCount(query)
	if err != nil {
		return 0, err
	}

	query = query.
		Offset(uint((pageIndex - 1) * pageSize)).
		Limit(uint(pageSize))

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return 0, err
	}

	// use gorm to scan rows
	result := scaner.Raw(sql, args...).Find(outRows)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// QueryCount ...
func (db *DataBase) QueryCount(query *goqu.SelectDataset) (int64, error) {
	var (
		count int64
	)

	sql, args, err := db.Goqu.From(
		query.As("query_count"),
	).Select(
		goqu.COUNT(goqu.Star()),
	).Prepared(true).ToSQL()
	if err != nil {
		return 0, err
	}

	result := db.Raw(sql, args...).Count(&count)
	if result.Error != nil {
		return 0, err
	}
	return count, nil
}

// QueryScan ...
func (db *DataBase) QueryScan(query *goqu.SelectDataset, outRows interface{}) error {
	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	return db.Raw(sql, args...).Scan(outRows).Error
}

// QueryRows ...
func (db *DataBase) QueryRows(query *goqu.SelectDataset, fn func(*sql.Rows) error) error {
	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := fn(rows); err != nil {
			return err
		}
	}

	return nil
}

// Transaction ...
func (db *DataBase) Transaction(f func(*gorm.DB) error) error {

	gdb := db.Begin()

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("critical error in db transaction: %v", err)
		}
	}()

	err := f(gdb)
	if err != nil {
		log.Errorf("db transaction failed: %v", err)
		gdb.Rollback()
		return err
	}

	err = gdb.Commit().Error
	if err != nil {
		log.Errorf("db transaction commit failed: %v", err)
		return err
	}

	return nil
}

// QueryPluck ...
func (db *DataBase) QueryPluck(query *goqu.SelectDataset, column string, outRows interface{}) error {
	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	return db.Raw(sql, args...).Pluck(column, outRows).Error
}
