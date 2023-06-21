package config

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *WasmConfig) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "MaxMutableGlobalBytes":
			z.MaxMutableGlobalBytes, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxMutableGlobalBytes")
				return
			}
		case "MaxTableElements":
			z.MaxTableElements, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxTableElements")
				return
			}
		case "MaxSectionElements":
			z.MaxSectionElements, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxSectionElements")
				return
			}
		case "MaxLinearMemoryInit":
			z.MaxLinearMemoryInit, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxLinearMemoryInit")
				return
			}
		case "MaxFuncLocalBytes":
			z.MaxFuncLocalBytes, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxFuncLocalBytes")
				return
			}
		case "MaxNestedStructures":
			z.MaxNestedStructures, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxNestedStructures")
				return
			}
		case "MaxSymbolBytes":
			z.MaxSymbolBytes, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxSymbolBytes")
				return
			}
		case "MaxModuleBytes":
			z.MaxModuleBytes, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxModuleBytes")
				return
			}
		case "MaxCodeBytes":
			z.MaxCodeBytes, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxCodeBytes")
				return
			}
		case "MaxPages":
			z.MaxPages, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxPages")
				return
			}
		case "MaxCallDepth":
			z.MaxCallDepth, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "MaxCallDepth")
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
func (z *WasmConfig) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 11
	// write "MaxMutableGlobalBytes"
	err = en.Append(0x8b, 0xb5, 0x4d, 0x61, 0x78, 0x4d, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxMutableGlobalBytes)
	if err != nil {
		err = msgp.WrapError(err, "MaxMutableGlobalBytes")
		return
	}
	// write "MaxTableElements"
	err = en.Append(0xb0, 0x4d, 0x61, 0x78, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxTableElements)
	if err != nil {
		err = msgp.WrapError(err, "MaxTableElements")
		return
	}
	// write "MaxSectionElements"
	err = en.Append(0xb2, 0x4d, 0x61, 0x78, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxSectionElements)
	if err != nil {
		err = msgp.WrapError(err, "MaxSectionElements")
		return
	}
	// write "MaxLinearMemoryInit"
	err = en.Append(0xb3, 0x4d, 0x61, 0x78, 0x4c, 0x69, 0x6e, 0x65, 0x61, 0x72, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x69, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxLinearMemoryInit)
	if err != nil {
		err = msgp.WrapError(err, "MaxLinearMemoryInit")
		return
	}
	// write "MaxFuncLocalBytes"
	err = en.Append(0xb1, 0x4d, 0x61, 0x78, 0x46, 0x75, 0x6e, 0x63, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxFuncLocalBytes)
	if err != nil {
		err = msgp.WrapError(err, "MaxFuncLocalBytes")
		return
	}
	// write "MaxNestedStructures"
	err = en.Append(0xb3, 0x4d, 0x61, 0x78, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxNestedStructures)
	if err != nil {
		err = msgp.WrapError(err, "MaxNestedStructures")
		return
	}
	// write "MaxSymbolBytes"
	err = en.Append(0xae, 0x4d, 0x61, 0x78, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxSymbolBytes)
	if err != nil {
		err = msgp.WrapError(err, "MaxSymbolBytes")
		return
	}
	// write "MaxModuleBytes"
	err = en.Append(0xae, 0x4d, 0x61, 0x78, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxModuleBytes)
	if err != nil {
		err = msgp.WrapError(err, "MaxModuleBytes")
		return
	}
	// write "MaxCodeBytes"
	err = en.Append(0xac, 0x4d, 0x61, 0x78, 0x43, 0x6f, 0x64, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxCodeBytes)
	if err != nil {
		err = msgp.WrapError(err, "MaxCodeBytes")
		return
	}
	// write "MaxPages"
	err = en.Append(0xa8, 0x4d, 0x61, 0x78, 0x50, 0x61, 0x67, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxPages)
	if err != nil {
		err = msgp.WrapError(err, "MaxPages")
		return
	}
	// write "MaxCallDepth"
	err = en.Append(0xac, 0x4d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x44, 0x65, 0x70, 0x74, 0x68)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MaxCallDepth)
	if err != nil {
		err = msgp.WrapError(err, "MaxCallDepth")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *WasmConfig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 11
	// string "MaxMutableGlobalBytes"
	o = append(o, 0x8b, 0xb5, 0x4d, 0x61, 0x78, 0x4d, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxMutableGlobalBytes)
	// string "MaxTableElements"
	o = append(o, 0xb0, 0x4d, 0x61, 0x78, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendUint32(o, z.MaxTableElements)
	// string "MaxSectionElements"
	o = append(o, 0xb2, 0x4d, 0x61, 0x78, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendUint32(o, z.MaxSectionElements)
	// string "MaxLinearMemoryInit"
	o = append(o, 0xb3, 0x4d, 0x61, 0x78, 0x4c, 0x69, 0x6e, 0x65, 0x61, 0x72, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x69, 0x74)
	o = msgp.AppendUint32(o, z.MaxLinearMemoryInit)
	// string "MaxFuncLocalBytes"
	o = append(o, 0xb1, 0x4d, 0x61, 0x78, 0x46, 0x75, 0x6e, 0x63, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxFuncLocalBytes)
	// string "MaxNestedStructures"
	o = append(o, 0xb3, 0x4d, 0x61, 0x78, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxNestedStructures)
	// string "MaxSymbolBytes"
	o = append(o, 0xae, 0x4d, 0x61, 0x78, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxSymbolBytes)
	// string "MaxModuleBytes"
	o = append(o, 0xae, 0x4d, 0x61, 0x78, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxModuleBytes)
	// string "MaxCodeBytes"
	o = append(o, 0xac, 0x4d, 0x61, 0x78, 0x43, 0x6f, 0x64, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxCodeBytes)
	// string "MaxPages"
	o = append(o, 0xa8, 0x4d, 0x61, 0x78, 0x50, 0x61, 0x67, 0x65, 0x73)
	o = msgp.AppendUint32(o, z.MaxPages)
	// string "MaxCallDepth"
	o = append(o, 0xac, 0x4d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x44, 0x65, 0x70, 0x74, 0x68)
	o = msgp.AppendUint32(o, z.MaxCallDepth)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *WasmConfig) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "MaxMutableGlobalBytes":
			z.MaxMutableGlobalBytes, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxMutableGlobalBytes")
				return
			}
		case "MaxTableElements":
			z.MaxTableElements, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxTableElements")
				return
			}
		case "MaxSectionElements":
			z.MaxSectionElements, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxSectionElements")
				return
			}
		case "MaxLinearMemoryInit":
			z.MaxLinearMemoryInit, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxLinearMemoryInit")
				return
			}
		case "MaxFuncLocalBytes":
			z.MaxFuncLocalBytes, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxFuncLocalBytes")
				return
			}
		case "MaxNestedStructures":
			z.MaxNestedStructures, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxNestedStructures")
				return
			}
		case "MaxSymbolBytes":
			z.MaxSymbolBytes, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxSymbolBytes")
				return
			}
		case "MaxModuleBytes":
			z.MaxModuleBytes, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxModuleBytes")
				return
			}
		case "MaxCodeBytes":
			z.MaxCodeBytes, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxCodeBytes")
				return
			}
		case "MaxPages":
			z.MaxPages, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxPages")
				return
			}
		case "MaxCallDepth":
			z.MaxCallDepth, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxCallDepth")
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
func (z *WasmConfig) Msgsize() (s int) {
	s = 1 + 22 + msgp.Uint32Size + 17 + msgp.Uint32Size + 19 + msgp.Uint32Size + 20 + msgp.Uint32Size + 18 + msgp.Uint32Size + 20 + msgp.Uint32Size + 15 + msgp.Uint32Size + 15 + msgp.Uint32Size + 13 + msgp.Uint32Size + 9 + msgp.Uint32Size + 13 + msgp.Uint32Size
	return
}
