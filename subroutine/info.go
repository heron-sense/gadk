package subroutine

import (
	"context"
	"github.com/heron-sense/gadk/logger"
)

//go:generate msgp
type GetInfo struct {
	Group string `msg:"group"`
}

func (z *GetInfo) Register() {
	//TODO implement me
	panic("implement me")
}

//go:generate msgp
type GetInfoReply struct {
	Group logger.StatGroup `msg:"group_stats"`
}

func (z *GetInfo) Handle(ctx context.Context) (Reply, error) {
	return &GetInfoReply{Group: logger.GetStats(z.Group)}, nil
}
