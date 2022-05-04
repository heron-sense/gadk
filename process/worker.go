package process

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/logger"
	"github.com/heron-sense/gadk/rpc"
	"github.com/heron-sense/gadk/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

func (i *island) processPack(rpcConn *rpc.IslandSession, routine *subroutineProfile, pk rpc.FlowPack) fsc.FlowStateCode {
	spanCtx, _ := i.Tracing.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(pk.GetExtension()))
	span := i.Tracing.Tracer.StartSpan(pk.GetDirective(), ext.RPCServerOption(spanCtx))
	defer span.Finish()
	span.SetTag(tracing.TagTracingID, pk.GetFlowTracingId())
	span.SetTag(tracing.TagIslandTerm, "unknown")
	span.SetTag(tracing.TagIslandName, i.name)

	beginTime := time.Now()
	reply, fsCode := routine.routine.Handle(pk)
	if !fsCode.Finished() {
		span.SetTag(tracing.TagFlowState, fsCode)
		logger.Vital("handle msg not finished")
		return fsCode
	}

	replyTime := time.Now()
	remainingTime, available := pk.CalRemainingTime(uint64(replyTime.UnixNano() / 1e6))
	if !available {
		logger.Vital("flow expire")
		return fsc.FlowExpireCancelled
	}

	span.SetTag(tracing.TagRemainingTime, remainingTime)

	serialized, fsCode := pk.GenReply(routine.replyDirective, uint64(beginTime.UnixNano()/1e6),
		remainingTime, uint32(fsCode), reply, nil)
	if !fsCode.Finished() {
		span.SetTag(tracing.TagFlowState, fsCode)
		logger.Error("gen reply err: %v", fsCode)
		return fsCode
	}

	fsCode = rpcConn.ReplyFlow(serialized)
	if !fsCode.Finished() {
		span.SetTag(tracing.TagFlowState, fsCode)
		logger.Error("send msg error:code", fsCode)
		return fsCode
	}

	span.LogFields(
		log.String("event", "string-format"),
	)

	span.SetTag(tracing.TagFlowState, 0)
	return fsc.FlowFinished
}
