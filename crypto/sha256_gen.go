package crypto

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Sha256) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Hash":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Hash")
				return
			}
			if zb0002 != uint32(4) {
				err = msgp.ArrayError{Wanted: uint32(4), Got: zb0002}
				return
			}
			for za0001 := range z.Hash {
				z.Hash[za0001], err = dc.ReadUint64()
				if err != nil {
					err = msgp.WrapError(err, "Hash", za0001)
					return
				}
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
func (z *Sha256) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Hash"
	err = en.Append(0x81, 0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(4))
	if err != nil {
		err = msgp.WrapError(err, "Hash")
		return
	}
	for za0001 := range z.Hash {
		err = en.WriteUint64(z.Hash[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Hash", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Sha256) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Hash"
	o = append(o, 0x81, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o = msgp.AppendArrayHeader(o, uint32(4))
	for za0001 := range z.Hash {
		o = msgp.AppendUint64(o, z.Hash[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Sha256) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Hash":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Hash")
				return
			}
			if zb0002 != uint32(4) {
				err = msgp.ArrayError{Wanted: uint32(4), Got: zb0002}
				return
			}
			for za0001 := range z.Hash {
				z.Hash[za0001], bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Hash", za0001)
					return
				}
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
func (z *Sha256) Msgsize() (s int) {
	s = 1 + 5 + msgp.ArrayHeaderSize + (4 * (msgp.Uint64Size))
	return
}