package api

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/MetalBlockchain/antelopevm/utils"
)

func init() {
	/**
	 * interface for primary index
	 */
	Functions["db_store_i64"] = storeI64
	Functions["db_update_i64"] = updateI64
	Functions["db_remove_i64"] = removeI64
	Functions["db_get_i64"] = getI64
	Functions["db_next_i64"] = nextI64
	Functions["db_previous_i64"] = previousI64
	Functions["db_find_i64"] = findI64
	Functions["db_lowerbound_i64"] = lowerboundI64
	Functions["db_upperbound_i64"] = upperboundI64
	Functions["db_end_i64"] = endI64

	/**
	 * interface for uint64_t secondary
	 */
	Functions["db_idx64_store"] = storeIdx64
	Functions["db_idx64_update"] = updateIdx64
	Functions["db_idx64_remove"] = removeIdx64
	Functions["db_idx64_find_secondary"] = findIdx64Secondary
	Functions["db_idx64_find_primary"] = findIdx64Primary
	Functions["db_idx64_lowerbound"] = lowerboundIdx64
	Functions["db_idx64_upperbound"] = upperboundIdx64
	Functions["db_idx64_end"] = endIdx64
	Functions["db_idx64_next"] = nextIdx64
	Functions["db_idx64_previous"] = previousIdx64

	/**
	 * interface for uint128_t secondary
	 */
	Functions["db_idx128_store"] = storeIdx128
	Functions["db_idx128_update"] = updateIdx128
	Functions["db_idx128_remove"] = removeIdx128
	Functions["db_idx128_find_secondary"] = findIdx128Secondary
	Functions["db_idx128_find_primary"] = findIdx128Primary
	Functions["db_idx128_lowerbound"] = lowerboundIdx128
	Functions["db_idx128_upperbound"] = upperboundIdx128
	Functions["db_idx128_end"] = endIdx128
	Functions["db_idx128_next"] = nextIdx128
	Functions["db_idx128_previous"] = previousIdx128

	/**
	 * interface for 256-bit interger secondary
	 */
	Functions["db_idx256_store"] = storeIdx256
	Functions["db_idx256_update"] = updateIdx256
	Functions["db_idx256_remove"] = removeIdx256
	Functions["db_idx256_find_secondary"] = findIdx256Secondary
	Functions["db_idx256_find_primary"] = findIdx256Primary
	Functions["db_idx256_lowerbound"] = lowerboundIdx256
	Functions["db_idx256_upperbound"] = upperboundIdx256
	Functions["db_idx256_end"] = endIdx256
	Functions["db_idx256_next"] = nextIdx256
	Functions["db_idx256_previous"] = previousIdx256

	/**
	 * interface for double secondary
	 */
	Functions["db_idx_double_store"] = storeIdxDouble
	Functions["db_idx_double_update"] = updateIdxDouble
	Functions["db_idx_double_remove"] = removeIdxDouble
	Functions["db_idx_double_find_secondary"] = findIdxDoubleSecondary
	Functions["db_idx_double_find_primary"] = findIdxDoublePrimary
	Functions["db_idx_double_lowerbound"] = lowerboundIdxDouble
	Functions["db_idx_double_upperbound"] = upperboundIdxDouble
	Functions["db_idx_double_end"] = endIdxDouble
	Functions["db_idx_double_next"] = nextIdxDouble
	Functions["db_idx_double_previous"] = previousIdxDouble

	/**
	 * interface for long double secondary
	 */
	Functions["db_idx_long_double_store"] = storeIdx64
	Functions["db_idx_long_double_update"] = updateIdx64
	Functions["db_idx_long_double_remove"] = removeIdx64
	Functions["db_idx_long_double_find_secondary"] = findIdx64Secondary
	Functions["db_idx_long_double_find_primary"] = findIdx64Primary
	Functions["db_idx_long_double_lowerbound"] = lowerboundIdx64
	Functions["db_idx_long_double_upperbound"] = upperboundIdx64
	Functions["db_idx_long_double_end"] = endIdx64
	Functions["db_idx_long_double_next"] = nextIdx64
	Functions["db_idx_long_double_previous"] = previousIdx64
}

func storeI64(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, buffer uint32, bufferSize uint32) int32 {
		data := context.ReadMemory(buffer, bufferSize)
		code := context.GetApplyContext().GetReceiver()
		iterator, err := context.GetApplyContext().StoreI64(code, scope, table, payer, id, data)

		if err != nil {
			panic("failed to store i64 object: " + err.Error())
		}

		return int32(iterator)
	}
}

func updateI64(context Context) interface{} {
	return func(iterator uint32, payer name.AccountName, buffer uint32, bufferSize uint32) {
		data := context.ReadMemory(buffer, bufferSize)

		if err := context.GetApplyContext().UpdateI64(int(iterator), payer, data, int(bufferSize)); err != nil {
			panic("failed to update i64 object: " + err.Error())
		}
	}
}

func removeI64(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetApplyContext().RemoveI64(int(iterator)); err != nil {
			panic("failed to remove i64")
		}
	}
}

func getI64(context Context) interface{} {
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

func nextI64(context Context) interface{} {
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

func previousI64(context Context) interface{} {
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

func findI64(context Context) interface{} {
	return func(code name.Name, scope name.ScopeName, table name.TableName, id uint64) int32 {
		iterator := context.GetApplyContext().FindI64(code, scope, table, id)

		return int32(iterator)
	}
}

func lowerboundI64(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) int32 {
		if res, err := context.GetApplyContext().LowerboundI64(code, scope, table, id); err == nil {
			return int32(res)
		} else {
			panic(err)
		}
	}
}

func upperboundI64(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, id uint64) int32 {
		if res, err := context.GetApplyContext().UpperboundI64(code, scope, table, id); err == nil {
			return int32(res)
		} else {
			panic(err)
		}
	}
}

func endI64(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName) int32 {
		res, err := context.GetApplyContext().EndI64(code, scope, table)

		if err != nil {
			panic(err)
		}

		return int32(res)
	}
}

// Idx64
func storeIdx64(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, ptr uint32) int32 {
		secondaryKey := readUint64(context, ptr)

		if iterator, err := context.GetIdx64().Store(scope, table, payer, id, secondaryKey); err != nil {
			panic(err)
		} else {
			return int32(iterator)
		}
	}
}

func updateIdx64(context Context) interface{} {
	return func(iterator int32, payer name.AccountName, ptr uint32) {
		secondaryKey := readUint64(context, ptr)

		if err := context.GetIdx64().Update(int(iterator), payer, secondaryKey); err != nil {
			panic(err)
		}
	}
}

func removeIdx64(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetIdx64().Remove(int(iterator)); err != nil {
			panic(err)
		}
	}
}

func findIdx64Secondary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint64(context, ptrSecondary)
		iterator := context.GetIdx64().FindSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func findIdx64Primary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary uint32, primary uint64) int32 {
		var secondaryKey uint64
		iterator := context.GetIdx64().FindPrimary(code, scope, table, &secondaryKey, primary)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func lowerboundIdx64(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint64(context, ptrSecondary)
		iterator := context.GetIdx64().LowerboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func upperboundIdx64(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint64(context, ptrSecondary)
		iterator := context.GetIdx64().UpperboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func endIdx64(context Context) interface{} {
	return func(code, scope, table name.Name) int32 {
		iterator := context.GetIdx64().EndSecondary(code, scope, table)

		return int32(iterator)
	}
}

func nextIdx64(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx64().NextSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func previousIdx64(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx64().PreviousSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

// Idx128
func storeIdx128(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, ptr uint32) int32 {
		secondaryKey := readUint128(context, ptr)

		if iterator, err := context.GetIdx128().Store(scope, table, payer, id, secondaryKey); err != nil {
			panic(err)
		} else {
			return int32(iterator)
		}
	}
}

func updateIdx128(context Context) interface{} {
	return func(iterator int32, payer name.AccountName, ptr uint32) {
		secondaryKey := readUint128(context, ptr)

		if err := context.GetIdx128().Update(int(iterator), payer, secondaryKey); err != nil {
			panic(err)
		}
	}
}

func removeIdx128(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetIdx128().Remove(int(iterator)); err != nil {
			panic(err)
		}
	}
}

func findIdx128Secondary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint128(context, ptrSecondary)
		iterator := context.GetIdx128().FindSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func findIdx128Primary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary uint32, primary uint64) int32 {
		var secondaryKey math.Uint128
		iterator := context.GetIdx128().FindPrimary(code, scope, table, &secondaryKey, primary)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func lowerboundIdx128(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint128(context, ptrSecondary)
		iterator := context.GetIdx128().LowerboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func upperboundIdx128(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint128(context, ptrSecondary)
		iterator := context.GetIdx128().UpperboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func endIdx128(context Context) interface{} {
	return func(code, scope, table name.Name) int32 {
		iterator := context.GetIdx128().EndSecondary(code, scope, table)

		return int32(iterator)
	}
}

func nextIdx128(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx128().NextSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func previousIdx128(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx128().PreviousSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

// Idx256
func storeIdx256(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, ptr uint32) int32 {
		secondaryKey := readUint256(context, ptr)

		if iterator, err := context.GetIdx256().Store(scope, table, payer, id, secondaryKey); err != nil {
			panic(err)
		} else {
			return int32(iterator)
		}
	}
}

func updateIdx256(context Context) interface{} {
	return func(iterator int32, payer name.AccountName, ptr uint32) {
		secondaryKey := readUint256(context, ptr)

		if err := context.GetIdx256().Update(int(iterator), payer, secondaryKey); err != nil {
			panic(err)
		}
	}
}

func removeIdx256(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetIdx256().Remove(int(iterator)); err != nil {
			panic(err)
		}
	}
}

func findIdx256Secondary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint256(context, ptrSecondary)
		iterator := context.GetIdx256().FindSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func findIdx256Primary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary uint32, primary uint64) int32 {
		var secondaryKey math.Uint256
		iterator := context.GetIdx256().FindPrimary(code, scope, table, &secondaryKey, primary)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint256(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func lowerboundIdx256(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint256(context, ptrSecondary)
		iterator := context.GetIdx256().LowerboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint256(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func upperboundIdx256(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readUint256(context, ptrSecondary)
		iterator := context.GetIdx256().UpperboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setUint256(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func endIdx256(context Context) interface{} {
	return func(code, scope, table name.Name) int32 {
		iterator := context.GetIdx256().EndSecondary(code, scope, table)

		return int32(iterator)
	}
}

func nextIdx256(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx256().NextSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func previousIdx256(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdx256().PreviousSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

// IdxDouble
func storeIdxDouble(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, ptr uint32) int32 {
		secondaryKey := readFloat64(context, ptr)

		if iterator, err := context.GetIdxDouble().Store(scope, table, payer, id, secondaryKey); err != nil {
			panic(err)
		} else {
			return int32(iterator)
		}
	}
}

func updateIdxDouble(context Context) interface{} {
	return func(iterator int32, payer name.AccountName, ptr uint32) {
		secondaryKey := readFloat64(context, ptr)

		if err := context.GetIdxDouble().Update(int(iterator), payer, secondaryKey); err != nil {
			panic(err)
		}
	}
}

func removeIdxDouble(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetIdxDouble().Remove(int(iterator)); err != nil {
			panic(err)
		}
	}
}

func findIdxDoubleSecondary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat64(context, ptrSecondary)
		iterator := context.GetIdxDouble().FindSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func findIdxDoublePrimary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary uint32, primary uint64) int32 {
		var secondaryKey float64
		iterator := context.GetIdxDouble().FindPrimary(code, scope, table, &secondaryKey, primary)

		if iterator <= -1 {
			return int32(iterator)
		}

		setFloat64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func lowerboundIdxDouble(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat64(context, ptrSecondary)
		iterator := context.GetIdxDouble().LowerboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setFloat64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func upperboundIdxDouble(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat64(context, ptrSecondary)
		iterator := context.GetIdxDouble().UpperboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setFloat64(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func endIdxDouble(context Context) interface{} {
	return func(code, scope, table name.Name) int32 {
		iterator := context.GetIdxDouble().EndSecondary(code, scope, table)

		return int32(iterator)
	}
}

func nextIdxDouble(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdxDouble().NextSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func previousIdxDouble(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdxDouble().PreviousSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

// IdxLongDouble
func storeIdxLongDouble(context Context) interface{} {
	return func(scope name.ScopeName, table name.TableName, payer name.AccountName, id uint64, ptr uint32) int32 {
		secondaryKey := readFloat128(context, ptr)

		if iterator, err := context.GetIdxLongDouble().Store(scope, table, payer, id, secondaryKey); err != nil {
			panic(err)
		} else {
			return int32(iterator)
		}
	}
}

func updateIdxLongDouble(context Context) interface{} {
	return func(iterator int32, payer name.AccountName, ptr uint32) {
		secondaryKey := readFloat128(context, ptr)

		if err := context.GetIdxLongDouble().Update(int(iterator), payer, secondaryKey); err != nil {
			panic(err)
		}
	}
}

func removeIdxLongDouble(context Context) interface{} {
	return func(iterator int32) {
		if err := context.GetIdxLongDouble().Remove(int(iterator)); err != nil {
			panic(err)
		}
	}
}

func findIdxLongDoubleSecondary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat128(context, ptrSecondary)
		iterator := context.GetIdxLongDouble().FindSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func findIdxLongDoublePrimary(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary uint32, primary uint64) int32 {
		var secondaryKey math.Float128
		iterator := context.GetIdxLongDouble().FindPrimary(code, scope, table, &secondaryKey, primary)

		if iterator <= -1 {
			return int32(iterator)
		}

		setFloat128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func lowerboundIdxLongDouble(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat128(context, ptrSecondary)
		iterator := context.GetIdxLongDouble().LowerboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setFloat128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func upperboundIdxLongDouble(context Context) interface{} {
	return func(code name.AccountName, scope name.ScopeName, table name.TableName, ptrSecondary, ptrPrimary uint32) int32 {
		var primaryKey uint64
		secondaryKey := readFloat128(context, ptrSecondary)
		iterator := context.GetIdxLongDouble().UpperboundSecondary(code, scope, table, &secondaryKey, &primaryKey)

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)
		setFloat128(context, ptrSecondary, secondaryKey)

		return int32(iterator)
	}
}

func endIdxLongDouble(context Context) interface{} {
	return func(code, scope, table name.Name) int32 {
		iterator := context.GetIdxLongDouble().EndSecondary(code, scope, table)

		return int32(iterator)
	}
}

func nextIdxLongDouble(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdxLongDouble().NextSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func previousIdxLongDouble(context Context) interface{} {
	return func(itr int32, ptrPrimary uint32) int32 {
		var primaryKey uint64
		iterator, err := context.GetIdxLongDouble().PreviousSecondary(int(itr), &primaryKey)

		if err != nil {
			panic(err)
		}

		if iterator <= -1 {
			return int32(iterator)
		}

		setUint64(context, ptrPrimary, primaryKey)

		return int32(iterator)
	}
}

func readUint64(context Context, ptr uint32) uint64 {
	var ret uint64
	value := context.ReadMemory(ptr, 8)
	if err := rlp.DecodeBytes(value, &ret); err != nil {
		panic(err)
	}
	return ret
}

func setUint64(context Context, ptr uint32, value uint64) {
	bytes, err := rlp.EncodeToBytes(value)

	if err != nil {
		panic(err)
	}

	context.WriteMemory(ptr, bytes)
}

func readUint128(context Context, ptr uint32) math.Uint128 {
	var ret math.Uint128
	value := context.ReadMemory(ptr, 16)
	if err := rlp.DecodeBytes(value, &ret); err != nil {
		panic(err)
	}
	return ret
}

func setUint128(context Context, ptr uint32, value math.Uint128) {
	bytes, err := rlp.EncodeToBytes(value)

	if err != nil {
		panic(err)
	}

	context.WriteMemory(ptr, bytes)
}

func readUint256(context Context, ptr uint32) math.Uint256 {
	var ret math.Uint256
	value := context.ReadMemory(ptr, 32)
	if err := rlp.DecodeBytes(value, &ret); err != nil {
		panic(err)
	}
	return ret
}

func setUint256(context Context, ptr uint32, value math.Uint256) {
	bytes, err := rlp.EncodeToBytes(value)

	if err != nil {
		panic(err)
	}

	context.WriteMemory(ptr, bytes)
}

func readFloat64(context Context, ptr uint32) float64 {
	var ret float64
	value := context.ReadMemory(ptr, 8)
	if err := rlp.DecodeBytes(value, &ret); err != nil {
		panic(err)
	}
	return ret
}

func setFloat64(context Context, ptr uint32, value float64) {
	bytes, err := rlp.EncodeToBytes(value)

	if err != nil {
		panic(err)
	}

	context.WriteMemory(ptr, bytes)
}

func readFloat128(context Context, ptr uint32) math.Float128 {
	var ret math.Float128
	value := context.ReadMemory(ptr, 16)
	if err := rlp.DecodeBytes(value, &ret); err != nil {
		panic(err)
	}
	return ret
}

func setFloat128(context Context, ptr uint32, value math.Float128) {
	bytes, err := rlp.EncodeToBytes(value)

	if err != nil {
		panic(err)
	}

	context.WriteMemory(ptr, bytes)
}
