package process

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/rpc"
	"sync"
)

type Subroutine interface {
	Handle(pack rpc.FlowPack) ([]byte, fsc.FlowStateCode)
}

type subroutineProfile struct {
	routine        Subroutine
	replyDirective []byte
	mutex          sync.RWMutex
	durationList   [35]uint64
	avgDuration    uint64
	appendPos      uint16
}
