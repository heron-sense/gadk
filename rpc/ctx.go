package rpc

import (
	"context"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"io"
	"time"
)

type FlowContext struct {
	Ctx           context.Context
	InitiateTime  time.Time
	FlowTracingId [FlowTracingIdLength]byte
	RemainingTime uint16
}

func (fc FlowContext) Deadline() (time.Time, bool) {
	return time.Now(), true
}

func (fc FlowContext) Done() <-chan struct{} {
	return nil
}

func (fc FlowContext) Err() error {
	return nil
}

func (fc FlowContext) Value(key interface{}) interface{} {
	return nil
}

func (ctx *FlowContext) Reply(pk FlowPack, data []byte) (FlowPack, fsc.FlowStateCode) {
	remainingTime, exceed := pk.CalRemainingTime(uint64(time.Now().UnixNano() / 1e6))
	if exceed {
		remainingTime = 0
	}

	reply := &_pack{
		PackMeta: PackMeta{
			FlowTracingId:   [FlowTracingIdLength]byte{},
			TrackSequence:   pk.GetTrackSequence(),
			InitiatedTime:   uint64(ctx.InitiateTime.UnixNano() / 1e6),
			RemainingTime:   remainingTime,
			DirectiveNotes:  0,
			DataRepFormat:   0,
			DataLength:      [3]uint8{},
			ExtensionNotes:  0,
			RoutingStrategy: 0,
			Reserved:        0,
			StateCode:       0,
			PackSignature:   [PackSignatureLength]uint8{},
		},
		Data: data,
	}
	reply.InitiatedTime = uint64(ctx.InitiateTime.UnixNano() / 1e6)
	return reply, fsc.FlowFinished
}

func InitJaeger(service string) (int, io.Closer) {
	return 0, nil
}

func NewPack(fCtx *FlowContext, directive string, stateCode uint32, data []byte, extension []byte) (*_pack, fsc.FlowStateCode) {
	length := uint32(len(data))
	nowMs := fCtx.InitiateTime.UnixNano() / 1e6
	pack := &_pack{
		PackMeta: PackMeta{
			FlowTracingId:  fCtx.FlowTracingId,
			TrackSequence:  0,
			InitiatedTime:  uint64(nowMs),
			RemainingTime:  fCtx.RemainingTime,
			DirectiveNotes: uint16(len(directive)),
			DataRepFormat:  0,
			DataLength: [3]uint8{
				uint8(length >> 16),
				uint8(length >> 8),
				uint8(length),
			},
			ExtensionNotes:  uint16(len(extension)),
			RoutingStrategy: 0,
			Reserved:        0,
			StateCode:       stateCode,
			PackSignature:   [PackSignatureLength]uint8{},
		},
		DstAddr:       "",
		SrcAddr:       "",
		flowTracingId: string(fCtx.FlowTracingId[:]),
		Directive:     directive,
		Data:          data,
		Extension:     extension,
	}

	return pack, fsc.FlowFinished
}
