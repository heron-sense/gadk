package rpc

import (
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"math"
	"net"
	"sync"
	"time"
)

var (
	sessionPool sync.Map
)

func GetSession(remainedTime time.Duration, address string) (*IslandSession, fsc.FlowStateCode) {
	mgr, ok := sessionPool.Load(address)
	if !ok {
		mgr, _ = sessionPool.LoadOrStore(address, NewSessionPool(math.MaxUint16))
	}

	sPool, ok := mgr.(*SessionPool)
	if !ok {
		return nil, fsc.FlowAssertFailed
	}

	return sPool.GetSession(remainedTime, address)
}

func Connect(remainedTime time.Duration, address string) (*IslandSession, fsc.FlowStateCode) {
	conn, err := net.DialTimeout("tcp", address, remainedTime)
	if err != nil {
		logError("dial err:%s", err)
		return nil, fsc.FlowNewSessionFailed
	}
	return NewNetConn(conn), fsc.FlowFinished
}

func NewNetConn(conn net.Conn) *IslandSession {
	return &IslandSession{
		createTime:   uint64(time.Now().UnixNano() / 1e3),
		conn:         conn,
		sendBytes:    0,
		rcvBytes:     0,
		sndTimes:     0,
		rcvTimes:     0,
		nextSequence: 1,
	}
}

func NewSessionPool(capacity uint16) *SessionPool {
	return &SessionPool{
		freeConn: make([]*IslandSession, 0, capacity),
		openerCh: make(chan struct{}),
		maxOpen:  capacity,
		maxIdle:  capacity,
		numOpen:  0,
		mu:       sync.Mutex{},
	}
}

type SessionPool struct {
	numOpen     uint16     // number of opened and pending open connections
	numBusy     uint16     // number of busy connections
	maxOpen     uint16     // 0 means unlimited
	maxIdle     uint16     // 0 means do not allow idle
	mu          sync.Mutex // protects following fields
	freeConn    []*IslandSession
	nextRequest uint64 // Next key to use in connRequests.
	openerCh    chan struct{}
}

func (pl *SessionPool) GetSession(remainedTime time.Duration, addr string) (*IslandSession, fsc.FlowStateCode) {
	pl.mu.Lock()
	var session *IslandSession
	if len(pl.freeConn) > 0 {
		session = pl.freeConn[0]
		pl.freeConn = pl.freeConn[1:]
		pl.mu.Unlock()
		return session, fsc.FlowFinished
	}

	if uint16(len(pl.freeConn)) < pl.maxOpen {
		pl.mu.Unlock()
		return Connect(remainedTime, addr)
	}

	return nil, fsc.FlowFinished
}

func (session *IslandSession) Dispose() {
	mgr, _ := sessionPool.Load(session.addr)
	sPool, ok := mgr.(*SessionPool)
	if !ok {
		return
	}

	sPool.mu.Lock()
	if uint16(len(sPool.freeConn)) < sPool.maxOpen {
		sPool.freeConn = append(sPool.freeConn, session)
		sPool.mu.Unlock()
		return
	}

	sPool.mu.Unlock()
	session.Close()
}
