package wasm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MetalBlockchain/antelopevm/utils"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type ExecutionContext struct {
	context context.Context
	engine  wazero.Runtime
	store   *wasmer.Store
	imports *wasmer.ImportObject
	ram     *wasmer.Memory
	module  api.Module
}

func (c *ExecutionContext) Initialize() {
	c.context = context.Background()
	c.engine = wazero.NewRuntime(c.context)
	c.imports = wasmer.NewImportObject()
	c.engine.NewHostModuleBuilder("env").
		ExportFunctions(GetAccountFunctions(c)).
		ExportFunctions(GetMemoryFunctions(c)).
		ExportFunctions(GetConsoleFunctions(c)).
		ExportFunctions(GetActionFunctions(c)).
		ExportFunctions(GetMathFunctions(c)).
		ExportFunction("eosio_assert", EosIoAssert(c)).
		ExportFunction("db_find_i64", FindI64(c)).
		ExportFunction("abort", Abort(c)).
		ExportFunction("db_update_i64", UpdateI64(c)).
		ExportFunction("db_store_i64", StoreI64(c)).
		ExportFunction("db_next_i64", NextI64(c)).
		ExportFunction("__extendsftf2", Extendsftf2(c)).
		ExportFunction("__floatsitf", Floatsitf(c)).
		ExportFunction("__multf3", Multf3(c)).
		ExportFunction("__floatunsitf", Floatunsitf(c)).
		ExportFunction("__divtf3", Divtf3(c)).
		ExportFunction("__addtf3", Addtf3(c)).
		ExportFunction("__extenddftf2", Extenddftf2(c)).
		ExportFunction("__eqtf2", Eqtf2(c)).
		ExportFunction("__letf2", Letf2(c)).
		ExportFunction("__netf2", Netf2(c)).
		ExportFunction("__subtf3", Subtf3(c)).
		ExportFunction("__trunctfdf2", Trunctfdf2(c)).
		ExportFunction("__getf2", Getf2(c)).
		ExportFunction("__trunctfsf2", Trunctfsf2(c)).
		ExportFunction("__unordtf2", Unordtf2(c)).
		ExportFunction("__fixunstfsi", Fixunstfsi(c)).
		ExportFunction("__fixtfsi", Fixtfsi(c)).
		ExportFunction("eosio_assert_code", AssertCode(c)).
		ExportFunction("db_get_i64", GetI64(c)).
		ExportFunction("db_remove_i64", RemoveI64(c)).
		Instantiate(c.context, c.engine)
	/* map[string]wasmer.IntoExtern{
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
	) */
}

func (c *ExecutionContext) Exec(code []byte) {
	module, err := c.engine.InstantiateModuleFromBinary(c.context, code)

	if err != nil {
		panic(fmt.Sprintln("Failed to instantiate the module: ", err))
	}

	c.module = module
	exportedFunc := module.ExportedFunction("apply")

	if exportedFunc == nil {
		panic(fmt.Sprintln("Failed to find apply function"))
	}

	receiver, _ := utils.StringToName("glenn")
	actionName, _ := utils.StringToName("transfer")

	start := time.Now()
	result, resultErr := exportedFunc.Call(c.context, receiver, receiver, actionName)
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	if resultErr != nil {
		panic(fmt.Sprintln("Execution failed: ", resultErr))
	}

	fmt.Printf("%v", result)
}
