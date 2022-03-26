package logger

import (
	"sync/atomic"
)

func Fatal(format string, argv ...interface{}) {
	rec := genRecord("FATAL", format, argv...)
	select {
	case logChan <- rec:
	default:
		atomic.AddInt64(&stats.Fatal.Bytes, int64(len(rec.serialized)))
		atomic.AddInt64(&stats.Fatal.Lines, 1)
	}
}

func Vital(format string, argv ...interface{}) {
	rec := genRecord("VITAL", format, argv...)
	select {
	case logChan <- rec:
	default:
		atomic.AddInt64(&stats.Vital.Bytes, int64(len(rec.serialized)))
		atomic.AddInt64(&stats.Vital.Lines, 1)
	}
}

func Debug(format string, argv ...interface{}) {
	rec := genRecord("DEBUG", format, argv...)
	select {
	case logChan <- rec:
	default:
		atomic.AddInt64(&stats.Debug.Bytes, int64(len(rec.serialized)))
		atomic.AddInt64(&stats.Debug.Lines, 1)
	}
}

func Error(format string, argv ...interface{}) {
	rec := genRecord("ERROR", format, argv...)
	select {
	case logChan <- rec:
	default:
		atomic.AddInt64(&stats.Error.Bytes, int64(len(rec.serialized)))
		atomic.AddInt64(&stats.Error.Lines, 1)
	}
}

func Alert(format string, argv ...interface{}) {
	rec := genRecord("ALERT", format, argv...)
	select {
	case logChan <- rec:
	default:
		atomic.AddInt64(&stats.Alert.Bytes, int64(len(rec.serialized)))
		atomic.AddInt64(&stats.Alert.Lines, 1)
	}
}
