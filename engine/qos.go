package engine

import (
	lrfCall "github.com/heron-sense/gadk/rpc"
	"sync"
	"time"
)

type subroutineProfile struct {
	replyDirective []byte
	mutex          sync.RWMutex
	durationList   [35]uint64
	avgDuration    uint64
	appendPos      uint16
}

func tooLate(tm time.Time, profile *subroutineProfile, pack *lrfCall.PackMeta) bool {
	profile.mutex.RLock()
	defer profile.mutex.RUnlock()
	if uint64(tm.UnixNano()/1e6)-pack.InitiatedTime >= uint64(pack.RemainingTime) {
		return true
	}
	return false
}

func perf(begin, end time.Time, profile *subroutineProfile) {
	if profile == nil {
		return
	}

	profile.mutex.Lock()
	profile.avgDuration -= profile.durationList[profile.appendPos]
	profile.durationList[profile.appendPos] = uint64(end.Sub(begin).Milliseconds())
	profile.avgDuration += profile.durationList[profile.appendPos]
	profile.appendPos++
	if int(profile.appendPos) > len(profile.durationList) {
		profile.appendPos = 0
	}
	profile.mutex.Unlock()
}
