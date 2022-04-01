package process

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/logger"
	"github.com/heron-sense/gadk/rpc"
	"time"
)

func (i *island) processPack(rpcConn *rpc.IslandSession, routine *subroutineProfile, pk rpc.FlowPack) fsc.FlowStateCode {
	beginTime := time.Now()
	reply, fsCode := routine.routine.Handle(pk)
	if !fsCode.Finished() {
		logger.Vital("handle msg not finished")
		return fsCode
	}

	replyTime := time.Now()
	remainingTime, available := pk.CalRemainingTime(uint64(replyTime.UnixNano() / 1e6))
	if !available {
		logger.Vital("flow expire")
		return fsc.FlowExpireCancelled
	}

	serialized, fsCode := pk.GenReply(routine.replyDirective, uint64(beginTime.UnixNano()/1e6),
		remainingTime, uint32(fsCode), reply, nil)
	if !fsCode.Finished() {
		logger.Error("gen reply err: %v", fsCode)
		return fsCode
	}

	fsCode = rpcConn.ReplyFlow(serialized)
	if !fsCode.Finished() {
		logger.Error("send msg error:code", fsCode)
		return fsCode
	}

	return fsc.FlowFinished
}
