package rpc

import (
	"bytes"
	"encoding/binary"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"time"
)

func (fc *FlowContext) SpanJsonCall(directive []byte, spanID uint8, msg []byte) ([]byte, fsc.FlowStateCode) {
	pk, ok := fc.FlowPack.(*_pack)
	if !ok {
		return nil, fsc.FlowAssertFailed
	}
	startTime := time.Now()
	nowMs := uint64(startTime.Nanosecond() / 1e6)
	remainingMs, enough := pk.GetRemainingTime(nowMs)
	if !enough {
		return nil, fsc.FlowExpireCancelled
	}
	reqHeader := &Header{
		FlowTracingId: pk.FlowTracingId,
		TrackSequence: pk.GenTrack(spanID),
		InitiatedTime: nowMs,
		RemainingTime: remainingMs,
	}

	buf := bytes.NewBuffer(make([]byte, 0, int(PackHeaderLength)+len(directive)+len(msg)))
	err := binary.Write(buf, binary.BigEndian, reqHeader)
	if err != nil {
		if logError != nil {
			logError("write buf[%s] err:%s", err)
		}
		return nil, fsc.FlowEncodeFailed
	}

	if fsCode := fc.conn.hgiSend(buf, directive, msg); !fsCode.Finished() {
		return nil, fsCode
	}

	replyPack, fsCode := fc.conn.RecvPack(time.Duration(remainingMs) * time.Millisecond)
	if !fsCode.Finished() {
		return nil, fsCode
	}

	finishTime := time.Now()
	rttTime := uint64(finishTime.UnixNano()/1e6) - replyPack.GetInitiatedTime()
	procTime := finishTime.Sub(startTime) - 2*time.Duration(rttTime)
	if logVital != nil {
		logVital("peer process time[%dms], rtt[%dms]", procTime, rttTime)
	}

	return replyPack.GetData(), fsCode
}

func (session *IslandSession) Dispatch(method string, location []byte, traceID [FlowTracingIdLength]byte, remainingTime uint16, msg []byte) ([]byte, fsc.FlowStateCode) {
	startTime := time.Now()
	directive := GenDirective(method, location)
	fCtx, fsCode := NewFlowContext(directive, traceID, 0, msg, nil)
	if !fsCode.Finished() {
		if logError != nil {
			logError("failed to create flow context: fsCode=%d", fsCode)
		}
		return nil, fsCode
	}

	serialized, fsCode := fCtx.FlowPack.Serialize()
	if fsCode := session.SendPack(serialized); !fsCode.Finished() {
		if logError != nil {
			logError("send err:%s", fsCode)
		}
		return nil, fsCode
	}

	logVital("send :%+v", serialized)

	replyPack, fsCode := session.RecvPack(time.Duration(remainingTime) * time.Millisecond)
	if !fsCode.Finished() {
		if logError != nil {
			logError("recv err:%s", fsCode)
		}
		return nil, fsCode
	}

	finishTime := time.Now()
	rttTime := uint64(finishTime.UnixNano()/1e6) - replyPack.GetInitiatedTime()
	procTime := finishTime.Sub(startTime) - 2*time.Duration(rttTime)
	if logVital != nil {
		logVital("peer process time[%dms], rtt[%dms]", procTime, rttTime)
	}

	return replyPack.GetData(), fsCode
}
