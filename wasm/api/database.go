package api

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/utils"
)

func GetDatabaseFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	/**
	 * interface for primary index
	 */
	functions["db_store_i64"] = storeI64(context)
	functions["db_update_i64"] = updateI64(context)
	functions["db_remove_i64"] = removeI64(context)
	functions["db_get_i64"] = getI64(context)
	functions["db_next_i64"] = nextI64(context)
	functions["db_previous_i64"] = previousI64(context)
	functions["db_find_i64"] = findI64(context)
	functions["db_lowerbound_i64"] = lowerboundI64(context)
	functions["db_upperbound_i64"] = upperboundI64(context)
	functions["db_end_i64"] = endI64(context)

	/**
	 * interface for uint64_t secondary
	 */
	functions["db_idx64_store"] = storeIdx64(context)
	functions["db_idx64_update"] = updateIdx64(context)
	functions["db_idx64_remove"] = removeIdx64(context)
	functions["db_idx64_find_secondary"] = findIdx64Secondary(context)
	functions["db_idx64_find_primary"] = findIdx64Primary(context)
	functions["db_idx64_lowerbound"] = lowerboundIdx64(context)
	functions["db_idx64_upperbound"] = upperboundIdx64(context)
	functions["db_idx64_end"] = endIdx64(context)
	functions["db_idx64_next"] = nextIdx64(context)
	functions["db_idx64_previous"] = previousIdx64(context)

	/**
	 * interface for uint128_t secondary
	 */
	functions["db_idx128_store"] = storeIdx64(context)
	functions["db_idx128_update"] = updateIdx64(context)
	functions["db_idx128_remove"] = removeIdx64(context)
	functions["db_idx128_find_secondary"] = findIdx64Secondary(context)
	functions["db_idx128_find_primary"] = findIdx64Primary(context)
	functions["db_idx128_lowerbound"] = lowerboundIdx64(context)
	functions["db_idx128_upperbound"] = upperboundIdx64(context)
	functions["db_idx128_end"] = endIdx64(context)
	functions["db_idx128_next"] = nextIdx64(context)
	functions["db_idx128_previous"] = previousIdx64(context)

	/**
	 * interface for 256-bit interger secondary
	 */
	functions["db_idx256_store"] = storeIdx64(context)
	functions["db_idx256_update"] = updateIdx64(context)
	functions["db_idx256_remove"] = removeIdx64(context)
	functions["db_idx256_find_secondary"] = findIdx64Secondary(context)
	functions["db_idx256_find_primary"] = findIdx64Primary(context)
	functions["db_idx256_lowerbound"] = lowerboundIdx64(context)
	functions["db_idx256_upperbound"] = upperboundIdx64(context)
	functions["db_idx256_end"] = endIdx64(context)
	functions["db_idx256_next"] = nextIdx64(context)
	functions["db_idx256_previous"] = previousIdx64(context)

	/**
	 * interface for double secondary
	 */
	functions["db_idx_double_store"] = storeIdx64(context)
	functions["db_idx_double_update"] = updateIdx64(context)
	functions["db_idx_double_remove"] = removeIdx64(context)
	functions["db_idx_double_find_secondary"] = findIdx64Secondary(context)
	functions["db_idx_double_find_primary"] = findIdx64Primary(context)
	functions["db_idx_double_lowerbound"] = lowerboundIdx64(context)
	functions["db_idx_double_upperbound"] = upperboundIdx64(context)
	functions["db_idx_double_end"] = endIdx64(context)
	functions["db_idx_double_next"] = nextIdx64(context)
	functions["db_idx_double_previous"] = previousIdx64(context)

	/**
	 * interface for long double secondary
	 */
	functions["db_idx_long_double_store"] = storeIdx64(context)
	functions["db_idx_long_double_update"] = updateIdx64(context)
	functions["db_idx_long_double_remove"] = removeIdx64(context)
	functions["db_idx_long_double_find_secondary"] = findIdx64Secondary(context)
	functions["db_idx_long_double_find_primary"] = findIdx64Primary(context)
	functions["db_idx_long_double_lowerbound"] = lowerboundIdx64(context)
	functions["db_idx_long_double_upperbound"] = upperboundIdx64(context)
	functions["db_idx_long_double_end"] = endIdx64(context)
	functions["db_idx_long_double_next"] = nextIdx64(context)
	functions["db_idx_long_double_previous"] = previousIdx64(context)

	return functions
}

func storeI64(context Context) func(core.ScopeName, core.TableName, core.AccountName, uint64, uint32, uint32) int32 {
	return func(scope core.ScopeName, table core.TableName, payer core.AccountName, id uint64, buffer uint32, bufferSize uint32) int32 {
		data := context.ReadMemory(buffer, bufferSize)
		code := context.GetApplyContext().GetReceiver()
		iterator, err := context.GetApplyContext().StoreI64(code, scope, table, payer, id, data)

		if err != nil {
			panic("failed to store i64 object: " + err.Error())
		}

		return int32(iterator)
	}
}

func updateI64(context Context) func(uint32, core.AccountName, uint32, uint32) {
	return func(iterator uint32, payer core.AccountName, buffer uint32, bufferSize uint32) {
		data := context.ReadMemory(buffer, bufferSize)

		if err := context.GetApplyContext().UpdateI64(int(iterator), payer, data, int(bufferSize)); err != nil {
			panic("failed to update i64 object: " + err.Error())
		}
	}
}

func removeI64(context Context) func(int32) {
	return func(iterator int32) {
		if err := context.GetApplyContext().RemoveI64(int(iterator)); err != nil {
			panic("failed to remove i64")
		}
	}
}

func getI64(context Context) func(int32, uint32, uint32) int32 {
	return func(iterator int32, buffer uint32, bufferSize uint32) int32 {
		bytes := make([]byte, bufferSize)
		size, err := context.GetApplyContext().GetI64(int(iterator), bytes, int(bufferSize))

		if err != nil {
			panic("failed to get i64")
		} else if bufferSize == 0 {
			return int32(size)
		}

		context.WriteMemory(buffer, bytes[0:size])

		return int32(size)
	}
}

func nextI64(context Context) func(int32, uint32) int32 {
	return func(iterator int32, ptr uint32) int32 {
		var primaryKey uint64

		if result, err := context.GetApplyContext().NextI64(int(iterator), &primaryKey); err == nil {
			context.WriteMemory(ptr, utils.Uint64ToLittleEndian(primaryKey))

			return int32(result)
		} else {
			panic(err)
		}
	}
}

func previousI64(context Context) func(int32, uint32) int32 {
	return func(iterator int32, ptr uint32) int32 {
		var primaryKey uint64

		if result, err := context.GetApplyContext().PreviousI64(int(iterator), &primaryKey); err == nil {
			context.WriteMemory(ptr, utils.Uint64ToLittleEndian(primaryKey))

			return int32(result)
		} else {
			panic(err)
		}
	}
}

func findI64(context Context) func(core.Name, core.ScopeName, core.TableName, uint64) int32 {
	return func(code core.Name, scope core.ScopeName, table core.TableName, id uint64) int32 {
		iterator := context.GetApplyContext().FindI64(code, scope, table, id)

		return int32(iterator)
	}
}

func lowerboundI64(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint64) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, id uint64) int32 {
		if res, err := context.GetApplyContext().LowerboundI64(code, scope, table, id); err == nil {
			return int32(res)
		} else {
			panic(err)
		}
	}
}

func upperboundI64(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint64) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, id uint64) int32 {
		if res, err := context.GetApplyContext().UpperboundI64(code, scope, table, id); err == nil {
			return int32(res)
		} else {
			panic(err)
		}
	}
}

func endI64(context Context) func(core.AccountName, core.ScopeName, core.TableName) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName) int32 {
		res, err := context.GetApplyContext().EndI64(code, scope, table)

		if err != nil {
			panic(err)
		}

		return int32(res)
	}
}

func storeIdx64(context Context) func(core.ScopeName, core.TableName, core.AccountName, uint64, uint32) int32 {
	return func(scope core.ScopeName, table core.TableName, payer core.AccountName, id uint64, ptr uint32) int32 {
		panic("not implemented")
	}
}

func updateIdx64(context Context) func(int32, core.AccountName, uint32) {
	return func(iterator int32, payer core.AccountName, ptr uint32) {
		panic("not implemented")
	}
}

func removeIdx64(context Context) func(int32) {
	return func(iterator int32) {
		panic("not implemented")
	}
}

func findIdx64Secondary(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint32, uint32, uint32, uint32) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, ptrSecondary, ptrSecondaryLength, ptrPrimary, ptrPrimaryLength uint32) int32 {
		panic("not implemented")
	}
}

func findIdx64Primary(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint32, uint64) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, ptrSecondary uint32, primary uint64) int32 {
		panic("not implemented")
	}
}

func lowerboundIdx64(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint32, uint32) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		panic("not implemented")
	}
}

func upperboundIdx64(context Context) func(core.AccountName, core.ScopeName, core.TableName, uint32, uint32) int32 {
	return func(code core.AccountName, scope core.ScopeName, table core.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		panic("not implemented")
	}
}

func endIdx64(context Context) func(core.AccountName, core.ScopeName, core.TableName) int32 {
	return func(code, scope, table core.Name) int32 {
		panic("not implemented")
	}
}

func nextIdx64(context Context) func(int32, uint32) int32 {
	return func(iterator int32, ptr uint32) int32 {
		panic("not implemented")
	}
}

func previousIdx64(context Context) func(int32, uint32) int32 {
	return func(iterator int32, ptr uint32) int32 {
		panic("not implemented")
	}
}
