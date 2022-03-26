package lf

import (
	"fmt"
	"testing"
	"time"
)

func TestEscape(t *testing.T)  {
	//expr:="/0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	expr:="heron-sense.com"
	mm:=make(map[string]int)
	mm[expr]=0
	//rr:=make([]byte,0,len(expr)*2)
	begin:=time.Now()
	for n:=0;n<65536*100;n++{
		//_,_=track(expr,rr)
		_,_=mm[expr]
	}
	fmt.Printf("%s", time.Since(begin))
}

