package logger

import (
	"fmt"
	"path"
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

type Dropped struct {
	Bytes int64 `json:"bytes"`
	Lines int64 `json:"lines"`
}

type logStats struct {
	Fatal Dropped `json:"dropped_fatal"`
	Error Dropped `json:"dropped_error"`
	Alert Dropped `json:"dropped_alert"`
	Vital Dropped `json:"dropped_vital"`
	Debug Dropped `json:"dropped_debug"`

	Latency struct {
	} `json:"latency"`
	WritePartialTimes uint64 `json:"write_partial_times"`
	WriteErrorTimes   uint64 `json:"write_error_times"`
	VolumeID          uint64 `json:"next_id"`
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

	argList := []interface{}{r.time.Format("Jan.02 15:04:05.000000"), level}
	if pc, file, line, ok := runtime.Caller(2); ok {
		format = "%s-%s(%s,%s:%d): " + format + "\n"
		file = path.Base(file)
		function := filepath.Ext(runtime.FuncForPC(pc).Name())
		function = strings.TrimLeft(function, ".")
		argList = append(argList, function, file, line)
		argList = append(argList, argv...)
	} else {
		format = "%s %s " + format + "\n"
		argList = append(argList, argv...)
	}

	r.serialized = fmt.Sprintf(format, argList...)

	return r
}
