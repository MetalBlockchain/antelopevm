package api

import (
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/producer"
)

func init() {
	Functions["is_feature_active"] = isFeatureActive
	Functions["activate_feature"] = activateFeature
	Functions["preactivate_feature"] = preactivateFeature
	Functions["set_resource_limits"] = setResourceLimits
	Functions["get_resource_limits"] = getResourceLimits
	Functions["get_wasm_parameters_packed"] = getWasmParametersPacked
	Functions["set_wasm_parameters_packed"] = setWasmParametersPacked
	Functions["set_proposed_producers"] = setProposedProducers
	Functions["set_proposed_producers_ex"] = setProposedProducersEx
	Functions["get_blockchain_parameters_packed"] = getBlockchainParametersPacked
	Functions["set_blockchain_parameters_packed"] = setBlockchainParametersPacked
	Functions["get_parameters_packed"] = getParametersPacked
	Functions["set_parameters_packed"] = setParametersPacked
	Functions["is_privileged"] = isPrivileged
	Functions["set_privileged"] = setPrivileged
}

func isFeatureActive(context Context) interface{} {
	return func(featureName name.Name) int32 {
		checkPrivileged(context)

		return 0
	}
}

func activateFeature(context Context) interface{} {
	return func(featureName name.Name) {
		checkPrivileged(context)

		panic("unsupported hardfork detected")
	}
}

func preactivateFeature(context Context) interface{} {
	return func(ptr uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setResourceLimits(context Context) interface{} {
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

func getResourceLimits(context Context) interface{} {
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

func getWasmParametersPacked(context Context) interface{} {
	return func(ptr uint32, length uint32, maxVersion uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setWasmParametersPacked(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setProposedProducers(context Context) interface{} {
	return func(ptr uint32, length uint32) int64 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setProposedProducersEx(context Context) interface{} {
	return func(format uint64, ptr uint32, length uint32) int64 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func getBlockchainParametersPacked(context Context) interface{} {
	return func(ptr uint32, length uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setBlockchainParametersPacked(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func getParametersPacked(context Context) interface{} {
	return func(idsPtr uint32, idsLength uint32, parametersPtr uint32, parametersLength uint32) uint32 {
		checkPrivileged(context)

		panic("not supported")
	}
}

func setParametersPacked(context Context) interface{} {
	return func(ptr uint32, length uint32) {
		checkPrivileged(context)

		panic("not supported")
	}
}

func isPrivileged(context Context) interface{} {
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

func setPrivileged(context Context) interface{} {
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
