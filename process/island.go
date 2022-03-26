package process

import (
	"github.com/heron-sense/gadk/logger"
)

func NewIsland(title string, addr string, subroutineMap map[string]Subroutine) *island {

	reactor := &island{
		SubroutineMap: make(map[string]*subroutineProfile),
		addr:          addr,
	}

	for id, subroutine := range subroutineMap {
		if subroutine == nil {
			logger.Error("invalid subroutine[%d]", id)
			return nil
		}
		reactor.SubroutineMap[id] = &subroutineProfile{
			routine: subroutine,
		}
	}
	return reactor
}
