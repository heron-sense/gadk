package lf


type LogicalForm interface {
	ForeachElem(selector string, proc func(elem interface{}) error) error
	AsInt64(field string, elem *int64) error
	AsUint64(field string, elem *uint64) error
	AsString(field string, elem *string) error
}

func NewLf() LogicalForm {
	return &_logicalForm{}
}