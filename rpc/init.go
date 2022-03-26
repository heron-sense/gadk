package rpc

import (
	"time"
)

const (
	recvPackTimeout = 100
)

var (
	logFatal LogFunc
	logVital LogFunc
	logAlert LogFunc
	logError LogFunc
	logDebug LogFunc
	StartID  [4]uint8
)

func init() {
	tmpRand := time.Now().UnixNano()
	for offset := uint32(0); offset < 32; offset += 8 {
		StartID[offset/8] = uint8(tmpRand >> (24 - offset))
	}
}

type LogFunc func(format string, argv ...interface{})

func Init(fatal LogFunc, vital LogFunc, alert LogFunc, error LogFunc, debug LogFunc) {
	logFatal = fatal
	logVital = vital
	logAlert = alert
	logError = error
	logDebug = debug
}

func ConfigTrack() {
}

type Strategy uint16

const (
	StrategyMin  = 0
	StrategyHash = 11
)
