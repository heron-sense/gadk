package logger

import (
	"time"
)

const (
	recvPackTimeout = 100
)

var (
	StartID [4]uint8
)

func init() {
	tmpRand := time.Now().UnixNano()
	for offset := uint32(0); offset < 32; offset += 8 {
		StartID[offset/8] = uint8(tmpRand >> (24 - offset))
	}
}

type LogFunc func(format string, argv ...interface{})

func ConfigTrack() {
}
