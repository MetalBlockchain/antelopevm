package wasm

import (
	"context"
	"fmt"
	"time"

	"github.com/MetalBlockchain/antelopevm/math"
	wasmApi "github.com/MetalBlockchain/antelopevm/wasm/api"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var _ wasmApi.Context = &ExecutionContext{}

type ExecutionContext struct {
	context               context.Context
	cancelFunc            context.CancelFunc
	engine                wazero.Runtime
	module                api.Module
	controller            wasmApi.Controller
	applyContext          wasmApi.ApplyContext
	authorizationManager  wasmApi.AuthorizationManager
	resourceLimitsManager wasmApi.ResourceLimitsManager
	idx64                 wasmApi.MultiIndex[uint64]
	idx128                wasmApi.MultiIndex[math.Uint128]
	idx256                wasmApi.MultiIndex[math.Uint256]
	idxDouble             wasmApi.MultiIndex[float64]
	idxLongDouble         wasmApi.MultiIndex[math.Float128]
}

func NewWasmExecutionContext(context context.Context,
	controller wasmApi.Controller,
	applyContext wasmApi.ApplyContext,
	authorizationManager wasmApi.AuthorizationManager,
	resourceLimitsManager wasmApi.ResourceLimitsManager,
	idx64 wasmApi.MultiIndex[uint64],
	idx128 wasmApi.MultiIndex[math.Uint128],
	idx256 wasmApi.MultiIndex[math.Uint256],
	idxDouble wasmApi.MultiIndex[float64],
	idxLongDouble wasmApi.MultiIndex[math.Float128],
) *ExecutionContext {
	return &ExecutionContext{
		context:               context,
		controller:            controller,
		applyContext:          applyContext,
		authorizationManager:  authorizationManager,
		resourceLimitsManager: resourceLimitsManager,
		idx64:                 idx64,
		idx128:                idx128,
		idx256:                idx256,
		idxDouble:             idxDouble,
		idxLongDouble:         idxLongDouble,
	}
}

func (c *ExecutionContext) Exec(wasmCode []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), 32*time.Nanosecond)
	runtime := wazero.NewRuntime(c.context)
	// This closes everything this runtime created
	defer runtime.Close(ctx)
	builder := c.engine.NewHostModuleBuilder("env")

	for name, f := range wasmApi.GetActionFunctions(c) {
		builder.NewFunctionBuilder().WithFunc(f).Export(name)
	}

	if _, err := builder.Instantiate(ctx); err != nil {
		return err
	}

	module, err := runtime.Instantiate(ctx, wasmCode)
	if err != nil {
		return err
	}

	// All Leap contracts export the apply function as the main entrypoint
	applyFunc := module.ExportedFunction("apply")
	if applyFunc == nil {
		return fmt.Errorf("failed to find apply function")
	}

	receiver := c.applyContext.GetReceiver()
	code := c.applyContext.GetAction().Account
	actionName := c.applyContext.GetAction().Name

	// Run the apply function with the given data
	if _, resultErr := applyFunc.Call(c.context, uint64(receiver), uint64(code), uint64(actionName)); resultErr != nil {
		return fmt.Errorf("execution failed: %s", resultErr)
	}

	return nil
}

// This function will read an array of bytes from the WASM memory, it panics on purpose when the read is out of range to kill the WASM execution environment
func (c *ExecutionContext) ReadMemory(start uint32, length uint32) []byte {
	if data, ok := c.module.Memory().Read(start, length); !ok {
		panic("memory read out of range")
	} else {
		return data
	}
}

// This function will write an array of bytes to the WASM memory, it panics on purpose when the write is out of range to kill the WASM execution environment
func (c *ExecutionContext) WriteMemory(start uint32, data []byte) {
	if ok := c.module.Memory().Write(start, data); !ok {
		panic("memory write out of range")
	}
}

func (c *ExecutionContext) GetMemorySize() uint32 {
	return c.module.Memory().Size()
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

func (c *ExecutionContext) GetResourceLimitsManager() wasmApi.ResourceLimitsManager {
	return c.resourceLimitsManager
}

func (c *ExecutionContext) GetIdx64() wasmApi.MultiIndex[uint64] {
	return c.idx64
}

func (c *ExecutionContext) GetIdx128() wasmApi.MultiIndex[math.Uint128] {
	return c.idx128
}

func (c *ExecutionContext) GetIdx256() wasmApi.MultiIndex[math.Uint256] {
	return c.idx256
}

func (c *ExecutionContext) GetIdxDouble() wasmApi.MultiIndex[float64] {
	return c.idxDouble
}

func (c *ExecutionContext) GetIdxLongDouble() wasmApi.MultiIndex[math.Float128] {
	return c.idxLongDouble
}

// Shutdown kills the running WASM context
func (c *ExecutionContext) Shutdown() {
	c.cancelFunc()
}
