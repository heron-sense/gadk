package maf

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"testing"
)

var (
	ColumnUserID   = DeclAsU64("user_id", true, 1000, math.MaxUint64)
	ColumnUserName = DeclAsStr("user_name", true, 4, 32)

	BaseColumns       = []Column{ColumnUserID, ColumnUserName}
	userAccessFormula *MysqlAccessFormula
)

func TestCreateDbAccessFormula(t *testing.T) {
	userAccessFormula, err := CreateDbAccessFormula(DataSourceCfg{
		Passwd: "Heron-sense.com#8611",
		User:   "account",
		Schema: "heron_store",
		Host:   "192.168.161.128",
	}, BaseColumns...)
	if err != nil {
		t.Fatal("init datasource err", err)
	}
	matched, err := userAccessFormula.OnlyCols(BaseColumns...).
		WithColumnInWhereClause(ColumnUserID).Exec().Find(context.TODO(), func(elem interface{}) error {
		return nil
	})
	if err != nil {
		t.Errorf("query err[%s]", err)
		return
	}
	if matched == 0 {
		t.Errorf("none rows found")
	}
}
