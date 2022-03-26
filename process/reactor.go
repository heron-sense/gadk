package process

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/logger"
	"github.com/heron-sense/gadk/rpc"
	"math"
	"net"
	"syscall"
	"time"
)

var (
	SubroutineNotFound = []byte(`SUBROUTINE.NOT-FOUND`)
)

type island struct {
	SubroutineMap map[string]*subroutineProfile
	addr          string
	AcceptErr     uint64
	AcceptSuc     uint64
}

func (i *island) ProcessUnknown(rpcConn *rpc.IslandSession, pk *rpc.FlowContext, acceptTime time.Time) fsc.FlowStateCode {
	defer rpcConn.Close()

	remainingTime, expired := pk.GetRemainingTime(uint64(acceptTime.UnixNano() / 1e6))
	if expired {
		remainingTime = 0
	}
	replyData, fsCode := pk.GenReply(SubroutineNotFound, uint64(acceptTime.UnixNano()/1e6),
		remainingTime,
		uint32(fsc.FlowSubroutineUndefined), nil, nil)
	if !fsCode.Finished() {
		logger.Alert("failed to gen reply for pack[%+v], which specified unknown subroutine[%s]",
			pk.FlowPack, pk.GetDirective())
		return fsCode
	}

	fsCode = rpcConn.ReplyFlow(replyData)
	if !fsCode.Finished() {
		logger.Alert("failed reply pack[%+v], which specified unknown subroutine[%s]",
			pk.FlowPack, pk.GetDirective())
		return fsCode
	}

	logger.Alert("reply pack[%+v], which specified unknown subroutine[%s]",
		pk.FlowPack, pk.GetDirective())
	return fsc.FlowFinished
}

func (i *island) ProcessConn(conn net.Conn, timeout time.Duration) fsc.FlowStateCode {
	logger.Vital("new conn[%s] from remote[%s] accepted",
		conn.LocalAddr().Network(),
		conn.RemoteAddr().Network())
	rpcConn := rpc.NewNetConn(conn)
	for {
		startTime := time.Now()
		pk, fsCode := rpcConn.RecvPack(timeout)
		if !fsCode.Finished() {
			logger.Error("recv over conn err:%d", fsCode)
			break
		}
		if pk == nil {
			logger.Error("conn closed")
			break
		}
		pk.BeginAcceptTime = startTime
		pk.AcceptTime = time.Now()

		routine, exist := i.SubroutineMap[pk.GetDirective()]
		if !exist {
			i.ProcessUnknown(rpcConn, pk, startTime)
			break
		}

		track := pk.GetFlowTracingId()
		logger.Debug("new pack:%s", track)
		fsCode = i.processPack(rpcConn, routine, pk)
		if !fsCode.Finished() {
			return fsCode
		}
	}

	return fsc.FlowFinished
}

func (i *island) ListenThenStart() fsc.FlowStateCode {
	listenCtx, err := net.Listen("tcp", i.addr)
	if err != nil {
		logger.Error("listen failed,err:%s", err)
		syscall.Exit(int(fsc.FlowNewSessionFailed))
	}

	logger.Vital("init finished, waiting for requests")
	for {
		con, err := listenCtx.Accept()
		if err != nil {
			logger.Fatal("accept failed", err)
			i.AcceptErr++
			if i.AcceptSuc >= 32 &&
				float64(i.AcceptErr) >= math.Log2(float64(i.AcceptSuc)) {
				break
			} else {
				continue
			}
		}

		i.AcceptSuc++
		go i.ProcessConn(con, 65535*time.Millisecond)
	}
	err = listenCtx.Close()
	if err != nil {
		logger.Vital("close listen ctx err:%s", err)
	}

	return fsc.FlowFinished
}
