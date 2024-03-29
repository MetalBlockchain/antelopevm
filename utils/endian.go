package utils

func EndianReverseU32(x uint32) uint32 {
	return ((x >> 0x18) & 0xFF) |
		((x>>0x10)&0xFF)<<0x08 |
		((x>>0x08)&0xFF)<<0x10 |
		(x&0xFF)<<0x18
}
