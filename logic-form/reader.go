package lf

import (
	"bytes"
)

func From(buf bytes.Buffer) (*_logicalForm, error) {
	return nil,nil
}

type _logicalForm struct {
	list [8]*_octTreeNode
}


func (f *_logicalForm) ForeachElem(selector string, proc func(elem interface{}) error) error {
	return nil
}

func (f *_logicalForm) AsInt64(selector string, elem *int64) error {
	return nil
}

func (f *_logicalForm) AsUint64(selector string, elem *uint64) error {
	return nil
}

func (f *_logicalForm) AsString(selector string, elem *string) error {
	return nil
}



