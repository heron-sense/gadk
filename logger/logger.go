package logger

import (
	"context"
	"encoding/json"
	"fmt"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"math"
	"os"
	"unsafe"
)

type CfgInfo struct {
	FileName string `toml:"file_name" json:"file_name"`
	FileSize uint64 `toml:"file_size" json:"file_size"`
	MaxFiles uint8  `toml:"max_files" json:"max_files"`
}

var (
	config            *CfgInfo
	stats             logStats
	currentVolumeSize int
	logVolume         *os.File
	logChan           chan *record
)

func Stat() (string, fsc.FlowStateCode) {
	serialized, err := json.Marshal(stats)
	if err != nil {
		return "", fsc.FlowEncodeFailed
	}
	return string(serialized), fsc.FlowFinished
}

func createFile(name string) (*os.File, fsc.FlowStateCode) {
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	file, err := os.OpenFile(name, flags, os.FileMode(0644))
	if err != nil {
		return nil, fsc.FlowCriticalSetUpFailed
	}
	return file, fsc.FlowFinished
}

func Init(ctx context.Context, cfg *CfgInfo) fsc.FlowStateCode {
	if cfg.FileSize < uint64(math.MaxUint16) || cfg.FileSize >= uint64(math.MaxUint32) {
		line := genRecord("ERROR", "invalid file-size[%d] specified, suggested range:[%d,%d)", cfg.FileSize, math.MaxUint16, math.MaxUint32)
		fmt.Printf("err: %s", line.serialized)
		return fsc.FlowPermissionDenied
	}
	config = cfg
	logChan = make(chan *record, 4096)

	go func() {
		stopped := false
		for {
			select {
			case _, exit := <-ctx.Done():
				if exit {
					stopped = true
				}
			case rec, ok := <-logChan:
				if !ok {
					if stopped {
						break
					}
					continue
				}

				doWrite(rec)
				pool.Put(rec)
			}
		}
	}()
	return 0
}

func doWrite(rec *record) fsc.FlowStateCode {
	if logVolume == nil || currentVolumeSize >= int(config.FileSize) {
		stats.VolumeID++
		name := fmt.Sprintf("%s.%d.txt", config.FileName, stats.VolumeID)
		file, fsCode := createFile(name)
		if !fsCode.Finished() {
			line := genRecord("ERROR", "failed to create new updated.sql file[%s]:%s", name, fsCode)
			fmt.Printf("%s", line.serialized)
			os.Exit(0)
		}

		if logVolume != nil {
			err := logVolume.Close()
			line := genRecord("ERROR", "failed to close updated.sql file[%d]:%s", stats.VolumeID-1, err)
			fmt.Printf("%s", line.serialized)
			os.Exit(0)
		}

		logVolume = file
	}

	strMeta := (*[2]uintptr)(unsafe.Pointer(&rec.serialized))
	bufMeta := [3]uintptr{strMeta[0], strMeta[1], strMeta[1]}
	switch nWritten, err := logVolume.Write(*(*[]byte)(unsafe.Pointer(&bufMeta))); {
	case err != nil:
		stats.WriteErrorTimes++
	case nWritten != len(rec.serialized):
		stats.WritePartialTimes++
	default:
		currentVolumeSize += len(rec.serialized)
	}
	return fsc.FlowFinished
}
