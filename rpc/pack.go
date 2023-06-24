package rpc

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	lib_enhanced "github.com/heron-sense/gadk/extension"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/logger"
)

type _pack struct {
	PackMeta
	DstAddr       string
	SrcAddr       string
	flowTracingId string
	Directive     string
	Digest        string
	Data          []byte
	Extension     []byte
}

func (pk *_pack) GetFlowTracingId() string {
	return pk.flowTracingId
}

func (pk *_pack) GetTrackSequence() uint32 {
	return pk.TrackSequence
}

func (pk *_pack) GetDstAddr() string {
	return pk.DstAddr
}

func (pk *_pack) GetSrcAddr() string {
	return pk.SrcAddr
}

func (pk *_pack) GenTrack(span uint8) uint32 {
	return GenTrackID(pk.TrackSequence, span)
}

func (pk *_pack) GetSha1Padding() []byte {
	return pk.PackMeta.PackSignature[:]
}

func (pk *_pack) GetData() []byte {
	return pk.Data
}
func (pk *_pack) GetTrack() uint32 {
	return pk.TrackSequence
}

func (pk *_pack) GetExtensionLen() uint32 {
	return uint32(pk.ExtensionNotes)
}

func (pk *_pack) GetDirectiveLen() uint32 {
	return uint32(pk.PackMeta.DirectiveNotes)
}

func (pk *_pack) GetDataLen() uint32 {
	var length uint32
	if pk != nil {
		for pos := 0; pos < len(pk.PackMeta.DataLength); pos++ {
			length = length<<8 | uint32(pk.PackMeta.DataLength[pos])
		}
	}
	return length
}

func (pk *_pack) GetDirective() string {
	if pk != nil {
		return pk.Directive
	}
	return ""
}

func (pk *_pack) GetRemainingTime() uint16 {
	return pk.RemainingTime
}

func (pk *_pack) CalRemainingTime(nowMs uint64) (uint16, bool) {
	if nowMs >= pk.InitiatedTime {
		elapseMs := nowMs - pk.InitiatedTime
		if uint64(pk.RemainingTime) > elapseMs {
			return pk.RemainingTime - uint16(elapseMs), true
		}
		return 0, false
	} else {
		return pk.RemainingTime, false
	}
}

func (pk *_pack) GetInitiatedTime() uint64 {
	return pk.PackMeta.InitiatedTime
}

func GenTrackID(trackID uint32, spanID uint8) uint32 {
	const strategyMask = uint32(3) << 30
	const trackMask = 0x3FFF_FFFF

	return (trackID & strategyMask) | ((trackID & trackMask) << 5) | uint32(spanID)
}

func (pk *_pack) Serialize() ([]byte, fsc.FlowStateCode) {
	totalLength := int(PackHeaderLength) + len(pk.Directive) + len(pk.Data) + len(pk.Extension)
	buf := bytes.NewBuffer(make([]byte, 0, totalLength))
	if err := binary.Write(buf, binary.BigEndian, pk.PackMeta); err != nil {
		logger.Error("write err:%s", err)
		return nil, fsc.FlowEncodeFailed
	}

	buf.Write([]byte(pk.Directive))
	buf.Write(pk.Data)
	buf.Write(pk.Extension)

	serialized := buf.Bytes()
	sum := sha1.Sum(serialized)
	sign := lib_enhanced.He32ofRaw(sum[:])

	copy(serialized[PackHeaderLength-PackSignatureLength:], sign)
	pk.Digest = string(sign)
	return serialized, fsc.FlowFinished
}

func (pk *_pack) GenReply(directive []byte, initiatedTime uint64, remainingTime uint16, stateCode uint32, data []byte, extension []byte) ([]byte, fsc.FlowStateCode) {
	dataLen := len(data)
	if (dataLen) > 0xFFFFFF {
		return nil, fsc.FlowDataLength
	}

	replyHeader := &PackMeta{
		FlowTracingId:   pk.FlowTracingId,
		TrackSequence:   pk.TrackSequence,
		InitiatedTime:   initiatedTime,
		RemainingTime:   remainingTime,
		DirectiveNotes:  uint16(len(directive)),
		DataRepFormat:   0,
		DataLength:      [3]uint8{uint8(dataLen >> 16), uint8(dataLen >> 8), uint8(dataLen)},
		ExtensionNotes:  uint16(len(extension)),
		RoutingStrategy: 0,
		Reserved:        0,
		StateCode:       stateCode,
		PackSignature:   [PackSignatureLength]uint8{},
	}

	totalLength := PackHeaderLength + len(directive) + len(data) + len(extension)
	buf := bytes.NewBuffer(make([]byte, 0, totalLength))
	if err := binary.Write(buf, binary.BigEndian, replyHeader); err != nil {
		logger.Error("write err:%s", err)
		return nil, fsc.FlowEncodeFailed
	}

	buf.Write(directive)
	buf.Write(data)
	buf.Write(extension)

	serialized := buf.Bytes()
	sum := sha1.Sum(serialized)
	sign := lib_enhanced.He32ofRaw(sum[:])
	copy(serialized[PackHeaderLength-PackSignatureLength:], sign)

	return serialized, fsc.FlowFinished
}
