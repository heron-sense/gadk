package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	pool sync.Pool
)

type record struct {
	time       time.Time
	serialized string
}

type Stats struct {
	Bytes int64 `json:"bytes"`
	Times int64 `json:"times"`
}

//go:generate msgp
type StatGroup struct {
	Fatal Stats `json:"fatal"`
	Error Stats `json:"error"`
	Alert Stats `json:"alert"`
	Vital Stats `json:"vital"`
	Debug Stats `json:"debug"`
}

type logStats struct {
	Written    StatGroup
	Failed     StatGroup
	Redirected StatGroup
	IoError    Stats  `json:"io_error"`
	VolumeID   uint64 `json:"next_id"`
}

func GetStats(group string) StatGroup {
	if group == "redirected" {
		return stats.Redirected
	} else if group == "failed" {
		return stats.Failed
	} else {
		return stats.Written
	}
}

func genRecord(level string, format string, argv ...interface{}) *record {
	node := pool.Get()
	var r *record
	if node != nil {
		r = node.(*record)
	} else {
		r = new(record)
	}
	r.time = time.Now()

	argList := make([]interface{}, 0, len(argv)+6)
	argList = append(argList, r.time.Format("2006-01-02 15:04:05.000000"))
	if pc, file, line, ok := runtime.Caller(2); ok {
		//2023-06-24 09:20:26.524492-DEBUG(parsereq,engine.go:182): new pack:%!s(MISSING)
		format = "%s[%s:%s:%d]-%s:" + format + "\n"
		file = strings.Split(file, "heron-sense/")[1]
		function := filepath.Ext(runtime.FuncForPC(pc).Name())

		argList = append(argList, file, function, line, level)
	} else {
		argList = append(argList, level)
		format = "%s-%s:" + format + "\n"
	}

	argList = append(argList, argv...)
	r.serialized = fmt.Sprintf(format, argList...)

	return r
}
