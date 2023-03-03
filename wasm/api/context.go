package api

type Context interface {
	GetController() Controller
	GetApplyContext() ApplyContext
	GetMultiIndex() MultiIndex
	GetAuthorizationManager() AuthorizationManager
	ReadMemory(start uint32, length uint32) []byte
	WriteMemory(start uint32, data []byte)
	GetMemorySize() uint32
}
