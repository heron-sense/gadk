package lf

import (
	"bytes"
)

/**
name=s32:heron.sense.zeal
balance=f10:7.12345678
total=i1:5
stat=d355:{orders=d14:{to_pay=i2:10},create_time=i10:1623445564}
list=l<e.order_id=$1>33:[d9:{$1=i1:1},d9:{$1=i1:2}]
order_list=l<e=i>8:[5,7,44]
*/

func (f *_logicalForm) Serialize(track []byte, buf *bytes.Buffer) (uint32, error) {
	serialized := buf.Len()

	for ptr:=f.list[0];ptr!=nil;ptr = ptr{
		buf.WriteByte('=')

		switch ptr.Symbol {
		case DataTypeLst:
			length:=2
			str:=make([]byte,0,length)
			buf.Write(str)
		case DataTypeRaw:
			length:=2
			str:=make([]byte,0,length)
			buf.Write(str)
		case DataTypeFlt:
			length:=2
			str:=make([]byte,0,length)
			buf.Write(str)
		case DataTypeBool:
 		default:
		}

	}

	return uint32(buf.Len()-serialized), nil
}

