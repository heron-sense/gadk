package process

import (
	"github.com/heron-sense/gadk/logger"
	"github.com/heron-sense/gadk/tracing"
	"github.com/opentracing/opentracing-go"
	"io"
)

func NewIsland(name string, addr string, subroutineMap map[string]Subroutine) *island {
	tracer, closer := tracing.InitJaeger(name)
	reactor := &island{
		Tracing: struct {
			Tracer opentracing.Tracer
			io.Closer
		}{Tracer: tracer, Closer: closer},
		SubroutineMap: make(map[string]*subroutineProfile),
		addr:          addr,
		name:          name,
	}
	opentracing.SetGlobalTracer(tracer)

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
