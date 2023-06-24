package logger

import (
	"fmt"
	"sync/atomic"
	"time"
)

func Fatal(format string, argv ...interface{}) {
	rec := genRecord("FATAL", format, argv...)
	writeAndStat(rec, &stats.Written.Fatal, &stats.Redirected.Fatal, &stats.Failed.Fatal)
}

func writeAndStat(rec *record, succStat *Stats, redirected *Stats, failed *Stats) {
	if config != nil {
		select {
		case logChan <- rec:
			atomic.AddInt64(&succStat.Bytes, int64(len(rec.serialized)))
			atomic.AddInt64(&succStat.Times, 1)
			return
		case <-time.After(20 * time.Millisecond):
			nBytes, _ := fmt.Print(rec.serialized)
			if nBytes == len(rec.serialized) {
				atomic.AddInt64(&redirected.Bytes, int64(len(rec.serialized)))
				atomic.AddInt64(&redirected.Times, 1)
			} else {
				atomic.AddInt64(&failed.Bytes, int64(len(rec.serialized)-nBytes))
				atomic.AddInt64(&failed.Times, 1)
			}
			break
		}
	}
}

func Vital(format string, argv ...interface{}) {
	rec := genRecord("VITAL", format, argv...)
	writeAndStat(rec, &stats.Written.Vital, &stats.Redirected.Vital, &stats.Failed.Vital)
}

func Debug(format string, argv ...interface{}) {
	rec := genRecord("DEBUG", format, argv...)
	writeAndStat(rec, &stats.Written.Debug, &stats.Redirected.Debug, &stats.Failed.Debug)
}

func Error(format string, argv ...interface{}) {
	rec := genRecord("ERROR", format, argv...)
	writeAndStat(rec, &stats.Written.Error, &stats.Redirected.Error, &stats.Failed.Error)
}

func Alert(format string, argv ...interface{}) {
	rec := genRecord("ALERT", format, argv...)
	writeAndStat(rec, &stats.Written.Alert, &stats.Redirected.Alert, &stats.Failed.Alert)
}
