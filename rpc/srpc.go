package rpc

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"time"
)

func (session *IslandSession) Dispatch(fCtx *FlowContext, method string, location []byte, msg []byte) ([]byte, fsc.FlowStateCode) {
	deadline := fCtx.InitiateTime.Add(time.Duration(fCtx.RemainingTime) * time.Millisecond)
	directive := GenDirective(method, location)

	req, fsCode := NewPack(fCtx, directive, 0, msg, nil)
	if !fsCode.Finished() {
		if logError != nil {
			logError("failed to create flow context: fsCode=%d", fsCode)
		}
		return nil, fsCode
	}

	serialized, fsCode := req.Serialize()
	if fsCode := session.SendPack(serialized); !fsCode.Finished() {
		if logError != nil {
			logError("send err:%s", fsCode)
		}
		return nil, fsCode
	}

	logVital("Pack[FlowTracingID=%s,Digest=%s] was Sent Successfully", req.GetFlowTracingId(), req.Digest)

	replyPack, fsCode := session.RecvPack(deadline)
	if !fsCode.Finished() {
		if logError != nil {
			logError("Rcv err:%s", fsCode)
		}
		return nil, fsCode
	}

	finishTime := time.Now()
	/**
	 * 发送时间 =
	 */
	sndTime := time.UnixMilli(int64(replyPack.GetInitiatedTime())).Sub(fCtx.InitiateTime).Milliseconds()
	procTime := fCtx.RemainingTime - replyPack.GetRemainingTime() - uint16(sndTime)
	rtt := uint16(0)
	if uint16(finishTime.Sub(fCtx.InitiateTime).Milliseconds()) > procTime {
		rtt = uint16(finishTime.Sub(fCtx.InitiateTime).Milliseconds()) - procTime
	}
	if logVital != nil {
		logVital("peer process time[%d], rtt[%d]", procTime, rtt)
	}

	return replyPack.GetData(), fsCode
}
