package api

import "github.com/MetalBlockchain/antelopevm/math"

type Context interface {
	GetController() Controller
	GetApplyContext() ApplyContext
	GetIdx64() MultiIndex[uint64]
	GetIdx128() MultiIndex[math.Uint128]
	GetIdx256() MultiIndex[math.Uint256]
	GetIdxDouble() MultiIndex[float64]
	GetIdxLongDouble() MultiIndex[math.Float128]
	GetAuthorizationManager() AuthorizationManager
	GetResourceLimitsManager() ResourceLimitsManager
	ReadMemory(start uint32, length uint32) []byte
	WriteMemory(start uint32, data []byte)
	GetMemorySize() uint32
}
