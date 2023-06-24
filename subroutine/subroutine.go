package subroutine

import (
	"context"
	"github.com/tinylib/msgp/msgp"
)

type Reply interface {
	Msgsize() (s int)
	MarshalMsg([]byte) ([]byte, error)
}

type Subroutine interface {
	msgp.Unmarshaler
	Register()
	Handle(ctx context.Context) (rsp Reply, err error)
}
