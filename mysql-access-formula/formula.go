package maf

import (
	"database/sql"
	"fmt"
)

type selectAllCols interface {
	SelectAllCols() *MysqlAccessFormula
}

type onlyCols interface {
	OnlyCols(cols ...Column) *MysqlAccessFormula
}

type excludeCols interface {
	ExcludeCols(cols ...Column) *MysqlAccessFormula
}

type withColumnInWhereClause interface {
	WithColumnInWhereClause(cols ...Column) *MysqlAccessFormula
}

type MysqlAccessFormula struct {
	db         *sql.DB
	datasource string
	Definition TableDefinition
	selectAllCols
	excludeCols
	onlyCols
	withColumnInWhereClause
}

func (f *MysqlAccessFormula) WithColumnInWhereClause(cols ...Column) *MysqlAccessFormula {
	return f
}

func (f *MysqlAccessFormula) OnlyCols(cols ...Column) *MysqlAccessFormula {
	return f
}

func (f *MysqlAccessFormula) SelectAllCols() *MysqlAccessFormula {
	return f
}

func (f *MysqlAccessFormula) ExcludeCols(cols ...Column) *MysqlAccessFormula {
	return f
}

func CreateDbAccessFormula(dsCfg DataSourceCfg, cols ...Column) (*MysqlAccessFormula, error) {
	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dsCfg.User, dsCfg.Passwd, dsCfg.Host, dsCfg.Schema)
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)
	formula := &MysqlAccessFormula{
		db: db,
		Definition: TableDefinition{
			cols: cols,
		},
	}
	return formula, nil
}
