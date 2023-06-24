package facilities

import (
	"github.com/heron-sense/gadk/subroutine"
	"sync"
	"sync/atomic"
)

var (
	isEnable         int32
	endpointMapMutex sync.RWMutex
	endpointMap      map[string]subroutine.Subroutine
)

func getResource(serviceName, methodName, configKey string) string {
	return configKey + "." + serviceName + "." + methodName
}

func IsEnable() bool {
	return atomic.LoadInt32(&isEnable) == 1
}

func GetEndpoint(serviceName, methodName, configKey string) (subroutine.Subroutine, bool) {
	endpointMapMutex.RLock()
	defer endpointMapMutex.RUnlock()
	ep, isGot := endpointMap[getResource(serviceName, methodName, configKey)]
	return ep, isGot
}

func ActivateForTest() {
	endpointMapMutex.Lock()
	defer endpointMapMutex.Unlock()

	endpointMap = make(map[string]subroutine.Subroutine)
	atomic.StoreInt32(&isEnable, 1)
}

func DeactivateForTest() {
	endpointMapMutex.Lock()
	defer endpointMapMutex.Unlock()

	endpointMap = make(map[string]subroutine.Subroutine)
	atomic.StoreInt32(&isEnable, 0)
}

func RegisterForTest(serviceName, methodName, configKey string, endpoint subroutine.Subroutine) {
	if !IsEnable() {
		return
	}
	endpointMapMutex.Lock()
	defer endpointMapMutex.Unlock()
	endpointMap[getResource(serviceName, methodName, configKey)] = endpoint
}
