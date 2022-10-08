package wasm

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/utils"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type ExecutionContext struct {
	engine  *wasmer.Engine
	store   *wasmer.Store
	imports *wasmer.ImportObject
	ram     *wasmer.Memory
}

func (context *ExecutionContext) Initialize() {
	context.engine = wasmer.NewEngine()
	context.store = wasmer.NewStore(context.engine)
	context.imports = wasmer.NewImportObject()
	limits, _ := wasmer.NewLimits(1, 4)
	context.ram = wasmer.NewMemory(context.store, wasmer.NewMemoryType(limits))
	context.imports.Register(
		"env",
		map[string]wasmer.IntoExtern{
			"require_auth":      RequireAuth(context),
			"eosio_assert":      EosIoAssert(context),
			"db_find_i64":       FindI64(context),
			"current_receiver":  CurrentReceiver(context),
			"abort":             Abort(context),
			"memset":            MemSet(context),
			"memcpy":            MemCopy(context),
			"db_update_i64":     UpdateI64(context),
			"db_store_i64":      StoreI64(context),
			"is_account":        IsAccount(context),
			"require_recipient": RequireRecipient(context),
			"has_auth":          HasAuth(context),
			"db_next_i64":       NextI64(context),
			"action_data_size":  ActionDataSize(context),
			"read_action_data":  ReacActionData(context),
			"memmove":           MemMove(context),
			"__extendsftf2":     Extendsftf2(context),
			"__floatsitf":       Floatsitf(context),
			"__multf3":          Multf3(context),
			"__floatunsitf":     Floatunsitf(context),
			"__divtf3":          Divtf3(context),
			"__addtf3":          Addtf3(context),
			"__extenddftf2":     Extenddftf2(context),
			"__eqtf2":           Eqtf2(context),
			"__letf2":           Letf2(context),
			"__netf2":           Netf2(context),
			"__subtf3":          Subtf3(context),
			"__trunctfdf2":      Trunctfdf2(context),
			"__getf2":           Getf2(context),
			"__trunctfsf2":      Trunctfsf2(context),
			"prints_l":          PrintsL(context),
			"__unordtf2":        Unordtf2(context),
			"__fixunstfsi":      Fixunstfsi(context),
			"__fixtfsi":         Fixtfsi(context),
			"eosio_assert_code": AssertCode(context),
			"db_get_i64":        GetI64(context),
			"db_remove_i64":     RemoveI64(context),
		},
	)
}

func (context *ExecutionContext) Exec(code []byte) {
	module, _ := wasmer.NewModule(context.store, code)
	instance, err := wasmer.NewInstance(module, context.imports)

	if err != nil {
		panic(fmt.Sprintln("Failed to instantiate the module: ", err))
	}

	exportedFunc, exportedFuncErr := instance.Exports.GetFunction("apply")

	if exportedFuncErr != nil {
		panic(fmt.Sprintln("Failed to instantiate the function: ", exportedFuncErr))
	}

	receiver := &utils.Name{Value: "glenn"}
	actionName := &utils.Name{Value: "transfer"}
	result, resultErr := exportedFunc(receiver.ToInt64(), receiver.ToInt64(), actionName.ToInt64())

	if resultErr != nil {
		panic(fmt.Sprintln("Execution failed: ", resultErr))
	}

	fmt.Printf("%v", result)
}
