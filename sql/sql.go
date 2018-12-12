package sql

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	goqu "gopkg.in/doug-martin/goqu.v5"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/mysql"
)

// QueryOp ..
type QueryOp string

//
const (
	EQ      QueryOp = "eq"
	NEQ     QueryOp = "neq"
	BETWEEN QueryOp = "between"
	IN      QueryOp = "in"
	NOTIN   QueryOp = "notIn"
	GT      QueryOp = "gt"
	GTE     QueryOp = "gte"
	LT      QueryOp = "lt"
	LTE     QueryOp = "lte"
	LIKE    QueryOp = "like"
	IS      QueryOp = "is"
)

// DataBase ...
type DataBase struct {
	*gorm.DB
	Goqu *goqu.Database
}

// ArgsFilter ...
type ArgsFilter struct {
	ex        goqu.Ex
	filterMap map[string]interface{}
	exOr      goqu.ExOr
}

// QueryFilter ...
func QueryFilter(filterMap map[string]interface{}) *ArgsFilter {
	return &ArgsFilter{
		ex:        make(goqu.Ex),
		filterMap: filterMap,
		exOr:      make(goqu.ExOr),
	}
}

// Ex ...
func (f *ArgsFilter) Ex() map[string]interface{} {
	return f.ex
}

// Where ...
func (f *ArgsFilter) Where(field string, op QueryOp,
	filterField string) *ArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.ex[field] = opEx
		}

	}
	return f
}

// Or ...
func (f *ArgsFilter) Or(field string, op QueryOp, filterField string) *ArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.exOr[field] = opEx
		}
	}
	return f
}

func (f *ArgsFilter) execOp(op QueryOp, value interface{}) goqu.Op {
	var opEx goqu.Op
	typ := reflect.TypeOf(value)
	switch op {
	case EQ, NEQ, GT, LT, LTE, GTE, IS:
		opEx = goqu.Op{string(op): value}
	case BETWEEN:
		if typ == nil {
			return opEx
		}
		kind := typ.Kind()
		v := reflect.ValueOf(value)
		if kind == reflect.Slice && v.Len() == 2 {
			opEx = goqu.Op{
				"between": goqu.RangeVal{
					Start: v.Index(0).Interface(),
					End:   v.Index(1).Interface(),
				},
			}
		}
	case IN, NOTIN:
		if typ == nil {
			return opEx
		}
		kind := typ.Kind()
		v := reflect.ValueOf(value)
		if kind == reflect.Slice && v.Len() > 0 {
			opEx = goqu.Op{string(op): value}
		}

	case LIKE:
		opEx = goqu.Op{string(LIKE): "%" + fmt.Sprintf("%v", value) + "%"}
	default:
	}

	return opEx
}

// End ...
func (f *ArgsFilter) End() []goqu.Expression {
	var ex []goqu.Expression
	if len(f.ex) > 0 {
		ex = append(ex, f.ex)
	}
	if len(f.exOr) > 0 {
		ex = append(ex, f.exOr)
	}
	return ex
}

// PageQuery ...
func (db *DataBase) PageQuery(query *goqu.Dataset, scaner *gorm.DB, pageIndex int64,
	pageSize int64, outRows interface{}, selectEx ...interface{}) (int64, error) {
	var selectQuery = query
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

	sql, args, err := selectQuery.ToSql()
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
func (db *DataBase) QueryCount(query *goqu.Dataset, selectEx ...interface{}) (int64, error) {
	var (
		selectQuery = query
		count       int64
	)

	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := db.Goqu.From(selectQuery.As("query_count")).Select(goqu.COUNT(goqu.L("*"))).Prepared(true).ToSql()
	if err != nil {
		return 0, err
	}

	result := db.Raw(sql, args...).Count(&count)
	if result.Error != nil {
		return 0, err
	}
	return count, nil
}

// Query ...
func (db *DataBase) Query(query *goqu.Dataset, scaner *gorm.DB,
	outRows interface{}, selectEx ...interface{}) error {

	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
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

// QueryScan ...
func (db *DataBase) QueryScan(query *goqu.Dataset,
	outRows interface{}, selectEx ...interface{}) error {
	selectQuery := query
	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
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

// QueryFirst ..
func (db *DataBase) QueryFirst(query *goqu.Dataset, scaner *gorm.DB,
	outRows interface{}, selectEx ...interface{}) error {

	selectQuery := query

	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return err
	}

	// use gorm to scan rows
	err = scaner.Raw(sql, args...).Limit(1).Find(outRows).Error
	if err != nil {
		return err
	}

	return nil
}

// QueryRows ..
func (db *DataBase) QueryRows(sqlBuilder *goqu.Dataset, scaner *gorm.DB,
	selectEx ...interface{}) (*sql.Rows, error) {
	selectQuery := sqlBuilder
	if selectEx != nil {
		selectQuery = selectQuery.Select(selectEx...)
	}
	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return nil, err
	}
	return scaner.Raw(sql, args...).Rows()
}

// DebugSQL ...
func DebugSQL(sqlBuilder *goqu.Dataset) string {
	sql, args, err := sqlBuilder.ToSql()
	return fmt.Sprint("Sql:", sql, args, err)
}

// StringOutRows ...
func StringOutRows(rows []string, pri ...string) (outs []interface{}) {
	var goquout []goqu.AliasedExpression
	for _, row := range rows {
		if len(pri) != 0 {
			goquout = append(goquout, goqu.I(pri[0]+"."+row).As(row))
		} else {
			goquout = append(goquout, goqu.I(row).As(row))
		}
	}

	for _, row := range goquout {
		outs = append(outs, row)
	}

	return
}
