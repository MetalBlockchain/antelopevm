package api

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/core/producer"
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

func isFeatureActive(context Context) func(name.Name) int32 {
	return func(featureName name.Name) int32 {
		checkPrivileged(context)

		return 0
	}
}

func activateFeature(context Context) func(name.Name) {
	return func(featureName name.Name) {
		checkPrivileged(context)

		panic("unsupported hardfork detected")
	}
}

func preactivateFeature(context Context) func(uint32) {
	return func(ptr uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setResourceLimits(context Context) func(name.AccountName, int64, int64, int64) {
	return func(account name.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) {
		checkPrivileged(context)

		eosAssert(ramBytes >= -1, "invalid value for ram resource limit expected [-1,INT64_MAX]")
		eosAssert(netWeight >= -1, "invalid value for net resource limit expected [-1,INT64_MAX]")
		eosAssert(cpuWeight >= -1, "invalid value for cpu resource limit expected [-1,INT64_MAX]")

		decreasedLimit, err := context.GetResourceLimitsManager().SetAccountLimits(account, ramBytes, netWeight, cpuWeight)

		if err != nil {
			panic(err)
		}

		if decreasedLimit {
			// TODO: context.trx_context.validate_ram_usage.insert( account );
		}
	}
}

func getResourceLimits(context Context) func(name.AccountName, uint32, uint32, uint32) {
	return func(account name.AccountName, ramBytesPtr uint32, netWeightPtr uint32, cpuWeightPtr uint32) {
		checkPrivileged(context)

		var ramBytes, netWeight, cpuWeight int64

		if err := context.GetResourceLimitsManager().GetAccountLimits(account, &ramBytes, &netWeight, &cpuWeight); err != nil {
			panic(err)
		}

		setUint64(context, ramBytesPtr, uint64(ramBytes))
		setUint64(context, netWeightPtr, uint64(netWeight))
		setUint64(context, cpuWeightPtr, uint64(cpuWeight))
	}
}

func getWasmParametersPacked(context Context) func(uint32, uint32, uint32) uint32 {
	return func(ptr uint32, length uint32, maxVersion uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setWasmParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setProposedProducers(context Context) func(uint32, uint32) int64 {
	return func(ptr uint32, length uint32) int64 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setProposedProducersEx(context Context) func(uint64, uint32, uint32) int64 {
	return func(format uint64, ptr uint32, length uint32) int64 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func getBlockchainParametersPacked(context Context) func(uint32, uint32) uint32 {
	return func(ptr uint32, length uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setBlockchainParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func getParametersPacked(context Context) func(uint32, uint32, uint32, uint32) uint32 {
	return func(idsPtr uint32, idsLength uint32, parametersPtr uint32, parametersLength uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setParametersPacked(context Context) func(uint32, uint32) {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func isPrivileged(context Context) func(name.AccountName) int32 {
	return func(account name.AccountName) int32 {
		checkPrivileged(context)

		privileged, err := context.GetApplyContext().IsPrivileged(account)

		if err != nil {
			panic(err)
		}

		if privileged {
			return 1
		}

		return 0
	}
}

func setPrivileged(context Context) func(name.AccountName, int32) {
	return func(account name.AccountName, isPrivileged int32) {
		checkPrivileged(context)

		if err := context.GetApplyContext().SetPrivileged(account, isPrivileged != 0); err != nil {
			panic(err)
		}
	}
}

func setProposedProducersCommon(context Context, producers []producer.ProducerKey, validateKeys bool) {
	panic("not implemented")
}

func checkPrivileged(context Context) {
	if !context.GetApplyContext().IsContextPrivileged() {
		panic(context.GetApplyContext().GetReceiver().String() + " does not have permission to call this API")
	}
}
