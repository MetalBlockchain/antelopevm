package api

type Context interface {
	GetApplyContext() ApplyContext
	ReadMemory(start uint32, length uint32) []byte
	WriteMemory(start uint32, data []byte)
}
