package quorm

import (
	"github.com/doug-martin/goqu/v8"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	// mysql drive
	_ "github.com/doug-martin/goqu/v8/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Quorm ...
type Quorm struct {
	*gorm.DB
	Goqu *goqu.Database
}

// Connect ...
func (db *Quorm) Connect(debug bool, dialect, dburl string) error {
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
func (db *Quorm) PageQuery(scaner *gorm.DB, query *goqu.SelectDataset, pageIndex, pageSize int64,
	outRows interface{}, selectEx ...interface{}) (int64, error) {

	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	count, err := db.QueryCount(query, selectEx...)
	if err != nil {
		return 0, err
	}

	selectQuery = query.
		Offset(uint((pageIndex - 1) * pageSize)).
		Limit(uint(pageSize))

	sql, args, err := selectQuery.Prepared(true).ToSQL()
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
func (db *Quorm) QueryCount(query *goqu.SelectDataset, selectEx ...interface{},
) (int64, error) {
	var (
		selectQuery = query
		count       int64
	)

	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := db.Goqu.From(
		selectQuery.As("query_count"),
	).Select(
		goqu.COUNT(goqu.L("*")),
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
func (db *Quorm) QueryScan(query *goqu.SelectDataset, outRows interface{},
	selectEx ...interface{}) error {

	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = db.Raw(sql, args...).Scan(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

// Transaction ...
func (db *Quorm) Transaction(f func(*gorm.DB) error) error {

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
