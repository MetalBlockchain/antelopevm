package resource

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ResourceUsage) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ID":
			err = z.ID.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Owner":
			err = z.Owner.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Owner")
				return
			}
		case "NetUsage":
			err = z.NetUsage.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "NetUsage")
				return
			}
		case "CpuUsage":
			err = z.CpuUsage.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "CpuUsage")
				return
			}
		case "RamUsage":
			z.RamUsage, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "RamUsage")
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
func (z *ResourceUsage) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "ID"
	err = en.Append(0x85, 0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = z.ID.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "Owner"
	err = en.Append(0xa5, 0x4f, 0x77, 0x6e, 0x65, 0x72)
	if err != nil {
		return
	}
	err = z.Owner.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Owner")
		return
	}
	// write "NetUsage"
	err = en.Append(0xa8, 0x4e, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return
	}
	err = z.NetUsage.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "NetUsage")
		return
	}
	// write "CpuUsage"
	err = en.Append(0xa8, 0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return
	}
	err = z.CpuUsage.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "CpuUsage")
		return
	}
	// write "RamUsage"
	err = en.Append(0xa8, 0x52, 0x61, 0x6d, 0x55, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.RamUsage)
	if err != nil {
		err = msgp.WrapError(err, "RamUsage")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ResourceUsage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "ID"
	o = append(o, 0x85, 0xa2, 0x49, 0x44)
	o, err = z.ID.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// string "Owner"
	o = append(o, 0xa5, 0x4f, 0x77, 0x6e, 0x65, 0x72)
	o, err = z.Owner.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Owner")
		return
	}
	// string "NetUsage"
	o = append(o, 0xa8, 0x4e, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65)
	o, err = z.NetUsage.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "NetUsage")
		return
	}
	// string "CpuUsage"
	o = append(o, 0xa8, 0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65)
	o, err = z.CpuUsage.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "CpuUsage")
		return
	}
	// string "RamUsage"
	o = append(o, 0xa8, 0x52, 0x61, 0x6d, 0x55, 0x73, 0x61, 0x67, 0x65)
	o = msgp.AppendUint64(o, z.RamUsage)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ResourceUsage) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ID":
			bts, err = z.ID.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Owner":
			bts, err = z.Owner.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Owner")
				return
			}
		case "NetUsage":
			bts, err = z.NetUsage.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "NetUsage")
				return
			}
		case "CpuUsage":
			bts, err = z.CpuUsage.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "CpuUsage")
				return
			}
		case "RamUsage":
			z.RamUsage, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RamUsage")
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
func (z *ResourceUsage) Msgsize() (s int) {
	s = 1 + 3 + z.ID.Msgsize() + 6 + z.Owner.Msgsize() + 9 + z.NetUsage.Msgsize() + 9 + z.CpuUsage.Msgsize() + 9 + msgp.Uint64Size
	return
}