package producer

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ProducerKey) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ProducerName":
			err = z.ProducerName.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "ProducerName")
				return
			}
		case "BlockSigningKey":
			err = z.BlockSigningKey.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "BlockSigningKey")
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
func (z *ProducerKey) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ProducerName"
	err = en.Append(0x82, 0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = z.ProducerName.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ProducerName")
		return
	}
	// write "BlockSigningKey"
	err = en.Append(0xaf, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = z.BlockSigningKey.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "BlockSigningKey")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ProducerKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ProducerName"
	o = append(o, 0x82, 0xac, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65)
	o, err = z.ProducerName.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ProducerName")
		return
	}
	// string "BlockSigningKey"
	o = append(o, 0xaf, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79)
	o, err = z.BlockSigningKey.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "BlockSigningKey")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ProducerKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ProducerName":
			bts, err = z.ProducerName.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ProducerName")
				return
			}
		case "BlockSigningKey":
			bts, err = z.BlockSigningKey.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlockSigningKey")
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
func (z *ProducerKey) Msgsize() (s int) {
	s = 1 + 13 + z.ProducerName.Msgsize() + 16 + z.BlockSigningKey.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SharedBlockSigningAuthority) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Threshold":
			z.Threshold, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "Threshold")
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
func (z SharedBlockSigningAuthority) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Threshold"
	err = en.Append(0x81, 0xa9, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.Threshold)
	if err != nil {
		err = msgp.WrapError(err, "Threshold")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z SharedBlockSigningAuthority) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Threshold"
	o = append(o, 0x81, 0xa9, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64)
	o = msgp.AppendUint32(o, z.Threshold)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SharedBlockSigningAuthority) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Threshold":
			z.Threshold, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Threshold")
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
func (z SharedBlockSigningAuthority) Msgsize() (s int) {
	s = 1 + 10 + msgp.Uint32Size
	return
}
