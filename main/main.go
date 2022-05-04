package main

import (
	"fmt"
	"net/http"
)

type Option func(cli *_client)

func WithTimeout(timeout uint32) Option {
	return func(cli *_client) {
		cli.Timeout = timeout
	}
}

func WithStatefulRoute(stateful bool) Option {
	return func(cli *_client) {
		cli.Stateful = stateful
	}
}

type Client interface {
	DoSomething()
}

type _client struct {
	Timeout  uint32
	Stateful bool
}

func (cli *_client) DoSomething() {
}

func NewClient(opts ...Option) Client {
	cli := &_client{}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func incr() (v int) {
	defer func(v int) {
		err := recover()
		fmt.Printf("raw=%d\n", err)
		v++
	}(v)

	panic(v)
}

func GenDirective(method string, location string) []byte {
	directive := make([]byte, 0, len(method)+len(location)+1)
	directive = append(directive, method...)
	directive = append(directive, ':')
	directive = append(directive, location...)
	return directive
}

func ParseExtension() http.Header {
	ext := []byte("a=b&c=d&e=")
	hdr := http.Header{}

	var kPos, vPos int
	for idx := 0; idx < len(ext); idx++ {
		switch chr := ext[idx]; chr {
		case '&':
			if vPos > 0 {
				key := string(ext[kPos : vPos-1])
				val := string(ext[vPos:idx])
				fmt.Printf("%s=%s\n", key, val)
				hdr.Set(key, val)
			}
			kPos = idx + 1
			vPos = 0
		case '=':
			vPos = idx + 1
			continue
		}
	}
	if vPos > 0 {
		key := string(ext[kPos : vPos-1])
		val := string(ext[vPos:])
		fmt.Printf("%s=%s\n", key, val)
		hdr.Set(key, val)
	}

	return hdr
}

func main() {
	fmt.Printf("%s\n", string([]byte("abcd")[1:3]))
	ParseExtension()
	fmt.Printf("finished")
	return
	NewClient(WithStatefulRoute(true), WithTimeout(32)).DoSomething()
	fmt.Println(incr())
}
