package rpc

import (
	"bytes"
	"encoding/binary"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"unsafe"
)

const (
	FlowTracingIdLength   = 32
	TrackSequenceLength   = 4
	InitiatedTimeLength   = 8
	RemainingTimeLength   = 2
	DirectiveNotesLength  = 2
	DataRepFormatLength   = 1
	DataLengthLength      = 3
	ExtensionNotesLength  = 2
	RoutingStrategyLength = 1
	ReservedLength        = 1
	StateCodeLength       = 4
	PackSignatureLength   = 32

	_hakDataMaxLength = 0xFF_FFFF
	PackHeaderLength  = FlowTracingIdLength +
		TrackSequenceLength + InitiatedTimeLength +
		RemainingTimeLength + DirectiveNotesLength +
		DataRepFormatLength + DataLengthLength +
		ExtensionNotesLength + RoutingStrategyLength +
		ReservedLength + StateCodeLength + PackSignatureLength
)

type Header struct {
	FlowTracingId   [FlowTracingIdLength]uint8 //7 bytes run-id + 13-bytes timestamp + 4 bytes rand
	TrackSequence   uint32                     //2bits(reserved) + 5bits + 5bits + 5bits + 5bits + 5bits + 5bits
	InitiatedTime   uint64                     //unix timestamp, ms
	RemainingTime   uint16                     //measured in ms
	DirectiveNotes  uint16                     //the MOST significant 5 bits for Method=POST/DELETE/PATCH/PUT/GET, and the remaining for Location Length
	DataRepFormat   uint8                      //Data Representation Format:=Url/Proto/Json/Raw
	DataLength      [3]uint8                   //big-endian
	ExtensionNotes  uint16                     //the MOST significant 5 bits for Extension
	RoutingStrategy byte                       //'H':=hash, 'Q':=最快响应, 'R':=round-robin
	Reserved        uint8
	StateCode       uint32
	PackSignature   [PackSignatureLength]uint8
}

type FlowPack interface {
	GetFlowTracingId() string
	GetTrackSequence() uint32
	CalRemainingTime(nowMs uint64) (uint16, bool)
	GetRemainingTime() uint16
	GetInitiatedTime() uint64
	GetData() []byte
	GetDirective() string
	GetDstAddr() string
	GetSrcAddr() string
	GetSha1Padding() []byte
	Serialize() ([]byte, fsc.FlowStateCode)
	GenReply(directive []byte, initiatedTime uint64, remainingTime uint16, stateCode uint32, data []byte, extension []byte) ([]byte, fsc.FlowStateCode)
}

func GenTrackHeader(traceID [FlowTracingIdLength]byte, strategy Strategy) *Header {
	hdr := &Header{
		FlowTracingId: traceID,
		TrackSequence: uint32(strategy),
	}
	hdr.TrackSequence <<= 30
	return hdr
}

func ParseHeader(buf [PackHeaderLength]byte) (*_pack, fsc.FlowStateCode) {
	pk := &_pack{}
	err := binary.Read(bytes.NewReader(buf[:]), binary.BigEndian, &pk.Header)
	if err != nil {
		if logError != nil {
			logError("err:%s", err)
		}
		return nil, fsc.FlowDecodeFailed
	}

	bufMeta := pk.FlowTracingId[:]
	pk.flowTracingId = *(*string)(unsafe.Pointer(&bufMeta))
	return pk, fsc.FlowFinished
}
