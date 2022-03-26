package rpc

import (
	"context"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"io"
	"time"
)

type FlowContext struct {
	ctx context.Context
	FlowPack
	BeginAcceptTime time.Time
	AcceptTime      time.Time
	conn            *IslandSession
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
	remainingTime, exceed := ctx.GetRemainingTime(uint64(time.Now().UnixNano() / 1e6))
	if exceed {
		remainingTime = 0
	}

	reply := &_pack{
		Header: Header{
			FlowTracingId:   [FlowTracingIdLength]byte{},
			TrackSequence:   pk.GetTrackSequence(),
			InitiatedTime:   uint64(ctx.AcceptTime.UnixNano() / 1e6),
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
	reply.InitiatedTime = uint64(ctx.AcceptTime.UnixNano() / 1e6)
	return reply, fsc.FlowFinished
}

func InitJaeger(service string) (int, io.Closer) {
	return 0, nil
}

func NewFlowContext(directive string, flowTracingId [FlowTracingIdLength]byte, stateCode uint32, data []byte, extension []byte) (*FlowContext, fsc.FlowStateCode) {
	length := uint32(len(data))
	fCtx := &FlowContext{
		ctx: nil,
		FlowPack: &_pack{
			Header: Header{
				FlowTracingId:  flowTracingId,
				TrackSequence:  0,
				InitiatedTime:  0,
				RemainingTime:  0,
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
			flowTracingId: string(flowTracingId[:]),
			Directive:     directive,
			Data:          data,
			Extension:     extension,
		},
		BeginAcceptTime: time.Time{},
		AcceptTime:      time.Time{},
		conn:            nil,
	}

	return fCtx, fsc.FlowFinished
}
