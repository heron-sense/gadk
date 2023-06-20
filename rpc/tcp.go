package rpc

import (
	"bytes"
	"crypto/sha1"
	"github.com/heron-sense/gadk/extension"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"io"
	"net"
	"strings"
	"time"
	"unsafe"
)

type IslandSession struct {
	addr         string
	createTime   uint64
	conn         net.Conn
	sendBytes    uint64
	rcvBytes     uint64
	sndTimes     uint32
	rcvTimes     uint32
	nextSequence uint64
}

func (session *IslandSession) hgiSend(buf *bytes.Buffer, directive []byte, data []byte) fsc.FlowStateCode {
	if _, err := buf.Write(directive); err != nil {
		if logError != nil {
			logError("write directive to buf[%s] err:%s", err)
		}
		return fsc.FlowEncodeFailed
	}

	if _, err := buf.Write(data); err != nil {
		if logError != nil {
			logError("write data to buf[%s] err:%s", err)
		}
		return fsc.FlowEncodeFailed
	}

	sign := sha1.Sum(buf.Bytes())
	_, err := buf.Write(extension.EncodeHgiPadding(sign))
	if err != nil {
		if logError != nil {
			logError("padding sign[%s] err:%s", err)
		}
		return fsc.FlowEncodeFailed
	}

	stream := buf.Bytes()
	nSent, err := session.conn.Write(stream)
	if err != nil {
		if logError != nil {
			logError("write session[%s] err:%s", session.conn.RemoteAddr(), err)
		}
		return fsc.FlowSendNotFinished
	}

	if nSent != len(stream) {
		if logError != nil {
			logError("sent[%d] of len[%d] to %s", nSent, len(stream), session.conn.RemoteAddr())
		}
		return fsc.FlowSendNotFinished
	}

	logError("%v, bytes sent:%d", stream, nSent)
	return fsc.FlowFinished
}

func (session *IslandSession) Close() {
	session.conn.Close()
}

func (session *IslandSession) SendPack(pack []byte) fsc.FlowStateCode {
	nSent, err := session.conn.Write(pack)
	if err != nil {
		if logError != nil {
			logError("write session[%s] err:%s", session.conn.RemoteAddr(), err)
		}
		return fsc.FlowSendNotFinished
	}

	if nSent != len(pack) {
		if logError != nil {
			logError("sent[%d] of len[%d] to %s", nSent, len(pack), session.conn.RemoteAddr())
		}
		return fsc.FlowSendNotFinished
	}

	return fsc.FlowFinished
}

func (session *IslandSession) ReplyFlow(reply []byte) fsc.FlowStateCode {
	deadline := time.Now().Add(100 * time.Millisecond)
	err := session.conn.SetWriteDeadline(deadline)
	if err != nil {
		if logError != nil {
			logError("set deadline err:%s", session.conn.RemoteAddr(), deadline, err)
		}
		return fsc.FlowCriticalSetUpFailed
	}

	sendBytes, err := session.conn.Write(reply)
	if err != nil {
		if logError != nil {
			logError("write session[%s] replyDeadline[%s] err:%s", session.conn.RemoteAddr(), deadline, err)
		}
		return fsc.FlowSendNotFinished
	}

	if sendBytes != len(reply) {
		return fsc.FlowSendNotFinished
	}

	session.sndTimes++
	session.sendBytes += uint64(sendBytes)
	return fsc.FlowFinished
}

func (session *IslandSession) RecvPack(deadline time.Time) (FlowPack, fsc.FlowStateCode) {
	err := session.conn.SetReadDeadline(deadline)
	if err != nil {
		if logFatal != nil {
			logFatal("write session[%s] err:%s", session.conn.RemoteAddr(), err)
		}
		return nil, fsc.FlowCriticalSetUpFailed
	}

	if logDebug != nil {
		logDebug("will read[%d] read deadline[%s]", PackHeaderLength, deadline)
	}

	var buf [PackHeaderLength]byte
	_, err = io.ReadFull(session.conn, buf[:])
	if err != nil {
		if err == io.EOF {
			return nil, fsc.FlowFinished
		}

		if logError != nil {
			logError("read data from[%s] err:%s", session.conn.RemoteAddr(), err)
		}

		if strings.Contains(err.Error(), "timeout") {
			return nil, fsc.FlowRecvNotFinished
		}
		return nil, fsc.FlowRecvErrorOccurred
	}

	pk, fCode := ParseMeta(buf)
	if !fCode.Finished() {
		if logError != nil {
			logError("parse header err:%s", fCode)
		}
		return nil, fsc.FlowRecvNotFinished
	}

	pk.DstAddr = session.conn.LocalAddr().String()
	pk.SrcAddr = session.conn.RemoteAddr().String()

	if logDebug != nil {
		logDebug("received new pack[FlowTracingId=%s,SrcAddr=%s,DstAddr=%s] directiveNodes=%v dataLength=[%v] extensionNode=[%v]",
			pk.flowTracingId, pk.SrcAddr, pk.DstAddr,
			pk.DirectiveNotes, pk.DataLength, pk.ExtensionNotes)
	}

	err = session.conn.SetReadDeadline(deadline)
	if err != nil {
		if logFatal != nil {
			logFatal("write session[%s] err:%s", session.conn.RemoteAddr(), err)
		}
		return nil, fsc.FlowCriticalSetUpFailed
	}

	readLen := pk.GetDirectiveLen() + pk.GetDataLen() + pk.GetExtensionLen()
	data := make([]byte, readLen, readLen)
	_, err = io.ReadFull(session.conn, data)
	if err != nil {
		if logError != nil {
			logError("read err:%s", err)
		}

		if strings.Contains(err.Error(), "timeout") {
			return nil, fsc.FlowRecvNotFinished
		}
		return nil, fsc.FlowRecvErrorOccurred
	}

	directive := data[:pk.GetDirectiveLen()]
	pk.Directive = *(*string)(unsafe.Pointer(&directive))
	pk.Data = data[pk.GetDirectiveLen() : pk.GetDirectiveLen()+pk.GetDataLen()]
	pk.Extension = data[pk.GetDirectiveLen()+pk.GetDataLen():]
	if logDebug != nil {
		logDebug("receive data finished[FlowTracingId=%s,directive=%s,extension=%s]", pk.flowTracingId, pk.Directive, pk.Extension)
	}

	return pk, fsc.FlowFinished
}
