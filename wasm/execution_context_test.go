package wasm

import (
	"io/ioutil"
	"testing"
)

func TestInitialize(t *testing.T) {
	wasmBytes, _ := ioutil.ReadFile("eosio.token.wasm")
	context := ExecutionContext{}
	context.Initialize()
	context.Exec(wasmBytes)
}
