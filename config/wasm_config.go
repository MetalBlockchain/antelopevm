package config

import "fmt"

//go:generate msgp
type WasmConfig struct {
	MaxMutableGlobalBytes uint32
	MaxTableElements      uint32
	MaxSectionElements    uint32
	MaxLinearMemoryInit   uint32
	MaxFuncLocalBytes     uint32
	MaxNestedStructures   uint32
	MaxSymbolBytes        uint32
	MaxModuleBytes        uint32
	MaxCodeBytes          uint32
	MaxPages              uint32
	MaxCallDepth          uint32
}

func DefaultInitialWasmConfiguration() WasmConfig {
	return WasmConfig{
		MaxMutableGlobalBytes: DefaultMaxWasmMutableGlobalBytes,
		MaxTableElements:      DefaultMaxWasmTableElements,
		MaxSectionElements:    DefaultMaxWasmSectionElements,
		MaxLinearMemoryInit:   DefaultMaxWasmLinearMemoryInit,
		MaxFuncLocalBytes:     DefaultMaxWasmFuncLocalBytes,
		MaxNestedStructures:   DefaultMaxWasmNestedStructures,
		MaxSymbolBytes:        DefaultMaxWasmSymbolBytes,
		MaxModuleBytes:        DefaultMaxWasmModuleBytes,
		MaxCodeBytes:          DefaultMaxWasmCodeBytes,
		MaxPages:              DefaultMaxWasmPages,
		MaxCallDepth:          DefaultMaxWasmCallDepth,
	}
}

func (w WasmConfig) Validate() error {
	if w.MaxSectionElements < 4 {
		return fmt.Errorf("max_section_elements cannot be less than 4")
	}

	if w.MaxFuncLocalBytes < 8 {
		return fmt.Errorf("max_func_local_bytes cannot be less than 8")
	}

	if w.MaxNestedStructures < 1 {
		return fmt.Errorf("max_nested_structures cannot be less than 1")
	}

	if w.MaxSymbolBytes < 32 {
		return fmt.Errorf("max_symbol_bytes cannot be less than 32")
	}

	if w.MaxModuleBytes < 256 {
		return fmt.Errorf("max_module_bytes cannot be less than 256")
	}

	if w.MaxCodeBytes < 32 {
		return fmt.Errorf("max_code_bytes cannot be less than 32")
	}

	if w.MaxPages < 1 {
		return fmt.Errorf("max_pages cannot be less than 1")
	}

	if w.MaxCallDepth < 2 {
		return fmt.Errorf("max_call_depth cannot be less than 2")
	}

	return nil
}
