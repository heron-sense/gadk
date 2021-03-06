package rpc

import (
	"time"

	fsc "github.com/heron-sense/gadk/flow-state-code"
)

func (session *IslandSession) Dispatch(fCtx *FlowContext, method string, location []byte, msg []byte, ext []byte) ([]byte, fsc.FlowStateCode) {
	deadline := fCtx.InitiateTime.Add(time.Duration(fCtx.RemainingTime) * time.Millisecond)
	directive := GenDirective(method, location)

	req, fsCode := NewPack(fCtx, directive, 0, msg, ext)
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
	 * ειζΆι΄ =
	 */
	rtt := uint16(finishTime.Sub(fCtx.InitiateTime).Milliseconds())
	procTime := fCtx.RemainingTime - replyPack.GetRemainingTime() - rtt

	if logVital != nil {
		logVital("peer process time[%d], rtt[%d],data[%s]", procTime, rtt, string(replyPack.GetData()))
	}

	return replyPack.GetData(), fsCode
}
