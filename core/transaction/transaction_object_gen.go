package transaction

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *TransactionObject) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Expiration":
			err = z.Expiration.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Expiration")
				return
			}
		case "TrxId":
			err = z.TrxId.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "TrxId")
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
func (z *TransactionObject) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "ID"
	err = en.Append(0x83, 0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = z.ID.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "Expiration"
	err = en.Append(0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = z.Expiration.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Expiration")
		return
	}
	// write "TrxId"
	err = en.Append(0xa5, 0x54, 0x72, 0x78, 0x49, 0x64)
	if err != nil {
		return
	}
	err = z.TrxId.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "TrxId")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TransactionObject) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "ID"
	o = append(o, 0x83, 0xa2, 0x49, 0x44)
	o, err = z.ID.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// string "Expiration"
	o = append(o, 0xaa, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o, err = z.Expiration.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Expiration")
		return
	}
	// string "TrxId"
	o = append(o, 0xa5, 0x54, 0x72, 0x78, 0x49, 0x64)
	o, err = z.TrxId.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "TrxId")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TransactionObject) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Expiration":
			bts, err = z.Expiration.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Expiration")
				return
			}
		case "TrxId":
			bts, err = z.TrxId.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "TrxId")
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
func (z *TransactionObject) Msgsize() (s int) {
	s = 1 + 3 + z.ID.Msgsize() + 11 + z.Expiration.Msgsize() + 6 + z.TrxId.Msgsize()
	return
}
