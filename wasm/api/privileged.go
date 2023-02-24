package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	log "github.com/inconshreveable/log15"
)

func GetPrivilegedFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["is_feature_active"] = isFeatureActive(context)
	functions["activate_feature"] = activateFeature(context)
	functions["preactivate_feature"] = preactivateFeature(context)
	functions["set_resource_limits"] = setResourceLimits(context)
	functions["get_resource_limits"] = getResourceLimits(context)
	functions["get_wasm_parameters_packed"] = getWasmParametersPacked(context)
	functions["set_wasm_parameters_packed"] = setWasmParametersPacked(context)
	functions["set_proposed_producers"] = setProposedProducers(context)
	functions["set_proposed_producers_ex"] = setProposedProducersEx(context)
	functions["get_blockchain_parameters_packed"] = getBlockchainParametersPacked(context)
	functions["set_blockchain_parameters_packed"] = setBlockchainParametersPacked(context)
	functions["get_parameters_packed"] = getParametersPacked(context)
	functions["set_parameters_packed"] = setParametersPacked(context)
	functions["is_privileged"] = isPrivileged(context)
	functions["set_privileged"] = setPrivileged(context)

	return functions
}

func isFeatureActive(context Context) func(core.Name) int32 {
	return func(featureName core.Name) int32 {
		log.Info("is_feature_active", "featureName", featureName)

		return 0
	}
}

func activateFeature(context Context) func(core.Name) {
	return func(featureName core.Name) {
		log.Info("activate_feature", "featureName", featureName)

		panic("not supported")
	}
}

func preactivateFeature(context Context) func(uint32) {
	return func(ptr uint32) {
		log.Info("preactivate_feature", "ptr", ptr)

		panic("not supported")
	}
}

func setResourceLimits(context Context) func(core.AccountName, int64, int64, int64) {
	return func(account core.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) {
		log.Info("set_resource_limits", "account", account, "ramBytes", ramBytes, "netWeight", netWeight, "cpuWeight", cpuWeight)

		panic("not supported")
	}
}

func getResourceLimits(context Context) func(core.AccountName, uint32, uint32, uint32) {
	return func(account core.AccountName, ramPtr uint32, netPtr uint32, cpuPtr uint32) {
		log.Info("get_resource_limits", "account", account, "ramPtr", ramPtr, "netPtr", netPtr, "cpuPtr", cpuPtr)

		panic("not supported")
	}
}

func getWasmParametersPacked(context Context) func(uint32, uint32, uint32) uint32 {
	return func(ptr uint32, length uint32, maxVersion uint32) uint32 {
		log.Info("get_wasm_parameters_packed", "ptr", ptr, "length", length, "maxVersion", maxVersion)

		panic("not supported")
	}
}

func setWasmParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("set_wasm_parameters_packed", "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func setProposedProducers(context Context) func(uint32, uint32) int64 {
	return func(ptr uint32, length uint32) int64 {
		log.Info("set_proposed_producers", "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func setProposedProducersEx(context Context) func(uint64, uint32, uint32) int64 {
	return func(format uint64, ptr uint32, length uint32) int64 {
		log.Info("set_proposed_producers_ex", "format", format, "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func getBlockchainParametersPacked(context Context) func(uint32, uint32) uint32 {
	return func(ptr uint32, length uint32) uint32 {
		log.Info("get_blockchain_parameters_packed", "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func setBlockchainParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("set_blockchain_parameters_packed", "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func getParametersPacked(context Context) func(uint32, uint32, uint32, uint32) uint32 {
	return func(idsPtr uint32, idsLength uint32, parametersPtr uint32, parametersLength uint32) uint32 {
		log.Info("get_parameters_packed", "idsPtr", idsPtr, "idsLength", idsLength)

		panic("not supported")
	}
}

func setParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		log.Info("set_parameters_packed", "ptr", ptr, "length", length)

		panic("not supported")
	}
}

func isPrivileged(context Context) func(core.AccountName) int32 {
	return func(account core.AccountName) int32 {
		log.Info("is_privileged", "account", account)

		panic("not supported")
	}
}

func setPrivileged(context Context) func(core.AccountName, int32) {
	return func(account core.AccountName, isPrivileged int32) {
		log.Info("set_privileged", "account", account, "isPrivileged", isPrivileged)

		panic("not supported")
	}
}

func setProposedProducersCommon(context Context, validateKeys bool) {
	panic("not supported")
}
