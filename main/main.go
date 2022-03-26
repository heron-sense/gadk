package main

import (
	"fmt"
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

func main() {
	NewClient(WithStatefulRoute(true), WithTimeout(32)).DoSomething()
	fmt.Println(incr())
}
