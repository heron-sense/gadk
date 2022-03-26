package rpc

func GenDirective(method string, location []byte) string {
	directive := make([]byte, 0, len(method)+len(location)+1)
	directive = append(directive, method...)
	directive = append(directive, ':')
	directive = append(directive, location...)
	return string(directive)
}
