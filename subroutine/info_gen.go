package subroutine

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *GetInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "group":
			z.Group, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Group")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z GetInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "group"
	err = en.Append(0x81, 0xa5, 0x67, 0x72, 0x6f, 0x75, 0x70)
	if err != nil {
		return
	}
	err = en.WriteString(z.Group)
	if err != nil {
		err = msgp.WrapError(err, "Group")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z GetInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "group"
	o = append(o, 0x81, 0xa5, 0x67, 0x72, 0x6f, 0x75, 0x70)
	o = msgp.AppendString(o, z.Group)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GetInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "group":
			z.Group, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Group")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z GetInfo) Msgsize() (s int) {
	s = 1 + 6 + msgp.StringPrefixSize + len(z.Group)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *GetInfoReply) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "group_stats":
			err = z.Group.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Group")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *GetInfoReply) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "group_stats"
	err = en.Append(0x81, 0xab, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73)
	if err != nil {
		return
	}
	err = z.Group.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Group")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GetInfoReply) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "group_stats"
	o = append(o, 0x81, 0xab, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73)
	o, err = z.Group.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Group")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GetInfoReply) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "group_stats":
			bts, err = z.Group.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Group")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *GetInfoReply) Msgsize() (s int) {
	s = 1 + 12 + z.Group.Msgsize()
	return
}
