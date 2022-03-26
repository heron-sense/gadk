package mysql_delegate

import (
	"database/sql"
	"fmt"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	lrf "github.com/heron-sense/gadk/logic-form"
)

type HSession struct {
	db *sql.DB
}

func (s HSession) SelectByCond(fetchList []*lrf.LogicalForm, clauseList []*Clause, parser func([]interface{}) (interface{}, fsc.FlowStateCode)) ([]interface{}, fsc.FlowStateCode) {
	res, err := s.db.Query("statement", clauseList)
	if err != nil {
		return nil, fsc.FlowRecvErrorOccurred
	}

	fieldList := make([]interface{}, 0, len(fetchList))

	list := make([]interface{}, 0)
	for res.Next() {
		err := res.Scan(fieldList...)
		if err != nil {
			return nil, fsc.FlowRecvErrorOccurred
		}

		obj, fsCode := parser(fieldList)
		if !fsCode.Finished() {
			return nil, fsCode
		}

		list = append(list, obj)
	}

	return nil, 0
}

func (s HSession) Insert(obj HObject) (int64, fsc.FlowStateCode) {
	return 0, 0
}

func (s HSession) Delete(obj HObject) (int64, fsc.FlowStateCode) {
	return 0, 0
}

func (s HSession) UpdateByCond(obj HObject, cond HObject) (int64, fsc.FlowStateCode) {
	res, err := s.db.Exec("statement")
	if err != nil {
		fmt.Printf("err=%+v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, 343243
	}
	return affected, fsc.FlowFinished
}

type HObject struct {
}

func (o *HObject) VerifyObject() bool {
	return false
}
