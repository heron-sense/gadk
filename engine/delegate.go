package engine

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/rpc"
)

type Handle struct {
	TraceID   [rpc.FlowTracingIdLength]byte
	Directive []byte
	Host      string
	Proj      string
	UrlPath   string
	Method    string
	Data      []byte //优先使用http-body填充, 对于DELETE/GET如果无body则使用query-string填充
	fsCode    fsc.FlowStateCode
}

//dispose
