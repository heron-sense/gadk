package fsc

type FlowStateCode int32

const (
	FlowSubroutineUndefined FlowStateCode = 20001
	FlowAllSessionOccupied  FlowStateCode = 20002
	FlowNewSessionFailed    FlowStateCode = 20003
	FlowSendNotFinished     FlowStateCode = 20004
	FlowServerUnreachable   FlowStateCode = 20005
	FlowBadRequest          FlowStateCode = 20006
	FlowPermissionDenied    FlowStateCode = 20007
	FlowNoReplyReceived     FlowStateCode = 20008
	FlowRecvNotFinished     FlowStateCode = 20009
	FlowEncodeFailed        FlowStateCode = 20010
	FlowDecodeFailed        FlowStateCode = 20011
	FlowDecodeIncomplete    FlowStateCode = 20012
	FlowBadReply            FlowStateCode = 20013
	FlowCriticalSetUpFailed FlowStateCode = 20014
	FlowDirectiveLength     FlowStateCode = 20014
	FlowDataLength          FlowStateCode = 20015
	FlowRecvErrorOccurred   FlowStateCode = 20016
	FlowAssertFailed        FlowStateCode = 20018
	FlowExpireCancelled     FlowStateCode = 20019
	FlowFinished            FlowStateCode = 0
)

func (flow FlowStateCode) Finished() bool {
	return flow == FlowFinished
}
