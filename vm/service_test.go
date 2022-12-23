package vm

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/stretchr/testify/assert"
)

func TestInitialKey(t *testing.T) {
	requestBody, _ := os.ReadFile("./test.json")
	body := &service.RequiredKeysRequest{}
	json.Unmarshal(requestBody, body)
	assert.Equal(t, body.AvailableKeys[0].String(), "a")
}
