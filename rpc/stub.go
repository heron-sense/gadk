package rpc

import (
	"context"
	fsc "github.com/heron-sense/gadk/flow-state-code"
)

type Stub struct {
	BeforeRead      func(ctx context.Context) fsc.FlowStateCode
	AfterRead       func(ctx context.Context) fsc.FlowStateCode
	BeforeWrite     func(ctx context.Context) fsc.FlowStateCode
	AfterWrite      func(ctx context.Context) fsc.FlowStateCode
	OnTimeoutClose  func(ctx context.Context) fsc.FlowStateCode
	BeforePeerClose func(ctx context.Context) fsc.FlowStateCode
}

//分包

//序列化

//反序列化
