package account

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Account) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Name":
			err = z.Name.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "VmType":
			z.VmType, err = dc.ReadUint8()
			if err != nil {
				err = msgp.WrapError(err, "VmType")
				return
			}
		case "VmVersion":
			z.VmVersion, err = dc.ReadUint8()
			if err != nil {
				err = msgp.WrapError(err, "VmVersion")
				return
			}
		case "Privileged":
			z.Privileged, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Privileged")
				return
			}
		case "LastCodeUpdate":
			err = z.LastCodeUpdate.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "LastCodeUpdate")
				return
			}
		case "CodeVersion":
			err = z.CodeVersion.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "CodeVersion")
				return
			}
		case "CreationDate":
			err = z.CreationDate.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "CreationDate")
				return
			}
		case "Code":
			err = z.Code.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Code")
				return
			}
		case "Abi":
			err = z.Abi.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Abi")
				return
			}
		case "AbiVersion":
			err = z.AbiVersion.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "AbiVersion")
				return
			}
		case "RecvSequence":
			z.RecvSequence, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "RecvSequence")
				return
			}
		case "AuthSequence":
			z.AuthSequence, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "AuthSequence")
				return
			}
		case "CodeSequence":
			z.CodeSequence, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "CodeSequence")
				return
			}
		case "AbiSequence":
			z.AbiSequence, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "AbiSequence")
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
func (z *Account) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 15
	// write "ID"
	err = en.Append(0x8f, 0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = z.ID.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "Name"
	err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = z.Name.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	// write "VmType"
	err = en.Append(0xa6, 0x56, 0x6d, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint8(z.VmType)
	if err != nil {
		err = msgp.WrapError(err, "VmType")
		return
	}
	// write "VmVersion"
	err = en.Append(0xa9, 0x56, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint8(z.VmVersion)
	if err != nil {
		err = msgp.WrapError(err, "VmVersion")
		return
	}
	// write "Privileged"
	err = en.Append(0xaa, 0x50, 0x72, 0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Privileged)
	if err != nil {
		err = msgp.WrapError(err, "Privileged")
		return
	}
	// write "LastCodeUpdate"
	err = en.Append(0xae, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65)
	if err != nil {
		return
	}
	err = z.LastCodeUpdate.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "LastCodeUpdate")
		return
	}
	// write "CodeVersion"
	err = en.Append(0xab, 0x43, 0x6f, 0x64, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = z.CodeVersion.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "CodeVersion")
		return
	}
	// write "CreationDate"
	err = en.Append(0xac, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65)
	if err != nil {
		return
	}
	err = z.CreationDate.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "CreationDate")
		return
	}
	// write "Code"
	err = en.Append(0xa4, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return
	}
	err = z.Code.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	// write "Abi"
	err = en.Append(0xa3, 0x41, 0x62, 0x69)
	if err != nil {
		return
	}
	err = z.Abi.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Abi")
		return
	}
	// write "AbiVersion"
	err = en.Append(0xaa, 0x41, 0x62, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = z.AbiVersion.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "AbiVersion")
		return
	}
	// write "RecvSequence"
	err = en.Append(0xac, 0x52, 0x65, 0x63, 0x76, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.RecvSequence)
	if err != nil {
		err = msgp.WrapError(err, "RecvSequence")
		return
	}
	// write "AuthSequence"
	err = en.Append(0xac, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.AuthSequence)
	if err != nil {
		err = msgp.WrapError(err, "AuthSequence")
		return
	}
	// write "CodeSequence"
	err = en.Append(0xac, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.CodeSequence)
	if err != nil {
		err = msgp.WrapError(err, "CodeSequence")
		return
	}
	// write "AbiSequence"
	err = en.Append(0xab, 0x41, 0x62, 0x69, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.AbiSequence)
	if err != nil {
		err = msgp.WrapError(err, "AbiSequence")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Account) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 15
	// string "ID"
	o = append(o, 0x8f, 0xa2, 0x49, 0x44)
	o, err = z.ID.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// string "Name"
	o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o, err = z.Name.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	// string "VmType"
	o = append(o, 0xa6, 0x56, 0x6d, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendUint8(o, z.VmType)
	// string "VmVersion"
	o = append(o, 0xa9, 0x56, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint8(o, z.VmVersion)
	// string "Privileged"
	o = append(o, 0xaa, 0x50, 0x72, 0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x64)
	o = msgp.AppendBool(o, z.Privileged)
	// string "LastCodeUpdate"
	o = append(o, 0xae, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65)
	o, err = z.LastCodeUpdate.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "LastCodeUpdate")
		return
	}
	// string "CodeVersion"
	o = append(o, 0xab, 0x43, 0x6f, 0x64, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o, err = z.CodeVersion.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "CodeVersion")
		return
	}
	// string "CreationDate"
	o = append(o, 0xac, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65)
	o, err = z.CreationDate.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "CreationDate")
		return
	}
	// string "Code"
	o = append(o, 0xa4, 0x43, 0x6f, 0x64, 0x65)
	o, err = z.Code.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	// string "Abi"
	o = append(o, 0xa3, 0x41, 0x62, 0x69)
	o, err = z.Abi.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Abi")
		return
	}
	// string "AbiVersion"
	o = append(o, 0xaa, 0x41, 0x62, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o, err = z.AbiVersion.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "AbiVersion")
		return
	}
	// string "RecvSequence"
	o = append(o, 0xac, 0x52, 0x65, 0x63, 0x76, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.RecvSequence)
	// string "AuthSequence"
	o = append(o, 0xac, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.AuthSequence)
	// string "CodeSequence"
	o = append(o, 0xac, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.CodeSequence)
	// string "AbiSequence"
	o = append(o, 0xab, 0x41, 0x62, 0x69, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.AbiSequence)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Account) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Name":
			bts, err = z.Name.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "VmType":
			z.VmType, bts, err = msgp.ReadUint8Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "VmType")
				return
			}
		case "VmVersion":
			z.VmVersion, bts, err = msgp.ReadUint8Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "VmVersion")
				return
			}
		case "Privileged":
			z.Privileged, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Privileged")
				return
			}
		case "LastCodeUpdate":
			bts, err = z.LastCodeUpdate.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastCodeUpdate")
				return
			}
		case "CodeVersion":
			bts, err = z.CodeVersion.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "CodeVersion")
				return
			}
		case "CreationDate":
			bts, err = z.CreationDate.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "CreationDate")
				return
			}
		case "Code":
			bts, err = z.Code.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Code")
				return
			}
		case "Abi":
			bts, err = z.Abi.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Abi")
				return
			}
		case "AbiVersion":
			bts, err = z.AbiVersion.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "AbiVersion")
				return
			}
		case "RecvSequence":
			z.RecvSequence, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RecvSequence")
				return
			}
		case "AuthSequence":
			z.AuthSequence, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AuthSequence")
				return
			}
		case "CodeSequence":
			z.CodeSequence, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CodeSequence")
				return
			}
		case "AbiSequence":
			z.AbiSequence, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AbiSequence")
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
func (z *Account) Msgsize() (s int) {
	s = 1 + 3 + z.ID.Msgsize() + 5 + z.Name.Msgsize() + 7 + msgp.Uint8Size + 10 + msgp.Uint8Size + 11 + msgp.BoolSize + 15 + z.LastCodeUpdate.Msgsize() + 12 + z.CodeVersion.Msgsize() + 13 + z.CreationDate.Msgsize() + 5 + z.Code.Msgsize() + 4 + z.Abi.Msgsize() + 11 + z.AbiVersion.Msgsize() + 13 + msgp.Uint64Size + 13 + msgp.Uint64Size + 13 + msgp.Uint64Size + 12 + msgp.Uint64Size
	return
}