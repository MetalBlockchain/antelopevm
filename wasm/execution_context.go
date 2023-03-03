package wasm

import (
	"context"
	"fmt"
	"time"

	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var _ wasmApi.Context = &ExecutionContext{}

type ExecutionContext struct {
	context              context.Context
	cancelFunc           context.CancelFunc
	engine               wazero.Runtime
	module               api.Module
	controller           wasmApi.Controller
	applyContext         wasmApi.ApplyContext
	authorizationManager wasmApi.AuthorizationManager
	multiIndex           wasmApi.MultiIndex
}

func NewWasmExecutionContext(context context.Context, controller wasmApi.Controller, applyContext wasmApi.ApplyContext, authorizationManager wasmApi.AuthorizationManager, multiIndex wasmApi.MultiIndex) *ExecutionContext {
	return &ExecutionContext{
		context:              context,
		controller:           controller,
		applyContext:         applyContext,
		authorizationManager: authorizationManager,
		multiIndex:           multiIndex,
	}
}

func (c *ExecutionContext) Initialize() error {
	c.context, c.cancelFunc = context.WithTimeout(context.Background(), 32*time.Nanosecond)
	c.engine = wazero.NewRuntime(c.context)

	if _, err := c.engine.NewHostModuleBuilder("env").
		ExportFunctions(wasmApi.GetActionFunctions(c)).
		ExportFunctions(wasmApi.GetAccountFunctions(c)).
		ExportFunctions(wasmApi.GetContextFreeSystemFunctions(c)).
		ExportFunctions(wasmApi.GetContextFreeTransactionFunctions(c)).
		ExportFunctions(wasmApi.GetMemoryFunctions(c)).
		ExportFunctions(wasmApi.GetConsoleFunctions(c)).
		ExportFunctions(wasmApi.GetContextFreeFunctions(c)).
		ExportFunctions(GetMathFunctions(c)).
		ExportFunctions(GetSystemFunctions(c)).
		ExportFunctions(wasmApi.GetDatabaseFunctions(c)).
		ExportFunctions(wasmApi.GetCryptoFunctions(c)).
		ExportFunctions(wasmApi.GetPermissionFunctions(c)).
		ExportFunctions(wasmApi.GetPrivilegedFunctions(c)).
		ExportFunctions(wasmApi.GetProducerFunctions(c)).
		ExportFunctions(wasmApi.GetTransactionFunctions(c)).
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
		Instantiate(c.context, c.engine); err != nil {
		return err
	}

	return nil
}

func (c *ExecutionContext) Exec(wasmCode []byte) error {
	module, err := c.engine.InstantiateModuleFromBinary(c.context, wasmCode)

	if err != nil {
		return fmt.Errorf("failed to instantiate the module: %s", err)
	}

	defer module.Close(c.context)

	c.module = module
	exportedFunc := module.ExportedFunction("apply")

	if exportedFunc == nil {
		return fmt.Errorf("failed to find apply function")
	}

	receiver := c.applyContext.GetReceiver()
	code := c.applyContext.GetAction().Account
	actionName := c.applyContext.GetAction().Name
	_, resultErr := exportedFunc.Call(c.context, uint64(receiver), uint64(code), uint64(actionName))

	if resultErr != nil {
		return fmt.Errorf("execution failed: %s", resultErr)
	}

	return nil
}

// This function will read an array of bytes from the WASM memory, it panics on purpose when the read is out of range to kill the WASM execution environment
func (c *ExecutionContext) ReadMemory(start uint32, length uint32) []byte {
	if data, ok := c.module.Memory().Read(c.context, start, length); !ok {
		panic("memory read out of range")
	} else {
		return data
	}
}

// This function will write an array of bytes to the WASM memory, it panics on purpose when the write is out of range to kill the WASM execution environment
func (c *ExecutionContext) WriteMemory(start uint32, data []byte) {
	if ok := c.module.Memory().Write(c.context, start, data); !ok {
		panic("memory write out of range")
	}
}

func (c *ExecutionContext) GetMemorySize() uint32 {
	return c.module.Memory().Size(c.context)
}

func (c *ExecutionContext) GetController() wasmApi.Controller {
	return c.controller
}

func (c *ExecutionContext) GetApplyContext() wasmApi.ApplyContext {
	return c.applyContext
}

func (c *ExecutionContext) GetAuthorizationManager() wasmApi.AuthorizationManager {
	return c.authorizationManager
}

func (c *ExecutionContext) GetMultiIndex() wasmApi.MultiIndex {
	return c.multiIndex
}

// Shutdown kills the running WASM context
func (c *ExecutionContext) Shutdown() {
	c.cancelFunc()
}
