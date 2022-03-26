package mysql_delegate

import (
	"bytes"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	lsf "github.com/heron-sense/gadk/logic-form"
)

type Stream interface {
	Serialize(buf bytes.Buffer) (uint32, fsc.FlowStateCode)
	SetStr(field lsf.LogicalForm, str string)
	PutInt(field lsf.LogicalForm, val int)
	PutFloat(field lsf.LogicalForm, val float64)
	Get(field lsf.LogicalForm)
}
