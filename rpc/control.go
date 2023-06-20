package rpc

import fsc "github.com/heron-sense/gadk/flow-state-code"

func EnableDesisting() {
	//一个开关，默认关闭。
	//打开这个开关，内存中的全局变量
}

func DesistWithResultIfMatch(proj string, headerFeature PackMeta, rsp []byte, fsCode fsc.FlowStateCode) {
	//放在一个context里
}
