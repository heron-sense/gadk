package maf

import (
	"context"
	"fmt"
	"log"
)

type FindHandle struct {
	f *MysqlAccessFormula
}

func (h *FindHandle) ApplyCols(cols ...Column) *FindHandle {
	return h
}

func (h *FindHandle) Get(ctx context.Context) (bool, error) {
	return false, fmt.Errorf("not explemented")
}

func (h *FindHandle) CountAndFetch(ctx context.Context) (int, []interface{}, error) {
	return 0, nil, fmt.Errorf("not explemented")
}

func (h *FindHandle) Fetch(ctx context.Context) ([]interface{}, error) {
	stmt, err := h.f.db.PrepareContext(ctx, "SELECT id,name FROM heron_store.tbl_user WHERE id >= ?")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(100)

	cols, err := rows.ColumnTypes()
	for _, col := range cols {
		fmt.Printf("%+v\n", col)
	}
	for rows.Next() {
		var id int
		var name []byte

		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Printf("what err:%s\n", err)
		}
		fmt.Printf("id=%d,name=%s\n", id, string(name))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}
