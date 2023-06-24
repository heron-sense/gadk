package rpc

func GenDirective(method string, location []byte) string {
	directive := make([]byte, 0, len(method)+len(location)+1)
	directive = append(directive, method...)
	directive = append(directive, ':')
	directive = append(directive, location...)
	return string(directive)
}

//protocol :// hostname[:port] / path / [;parameters][?query]#fragment [编辑本段]格式说明： URL的组成

func GetRoute() {
	//首先从ctx中找。
	//从动态配置找
	//从静态配置中找
	//ctx中可以控制过程，和查找范围
}
