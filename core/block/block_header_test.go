package block

import (
	"testing"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCalculateId(t *testing.T) {
	time, err := core.FromIsoString("2018-06-09T11:56:30.000")
	if err != nil {
		panic(err)
	}
	microSinceEpoch := time.TimeSinceEpoch()
	msecSinceEpoch := microSinceEpoch.Count() / 1000
	slot := (msecSinceEpoch - 946684800000) / 500
	block := BlockHeader{
		Timestamp:             uint32(slot),
		Producer:              name.StringToName("eosio"),
		Confirmed:             0,
		Previous:              *crypto.NewSha256String("00000001405147477ab2f5f51cda427b638191c66d2c59aa392d5c2c98076cb0"),
		TransactionMerkleRoot: *crypto.NewSha256String("0000000000000000000000000000000000000000000000000000000000000000"),
		ActionMerkleRoot:      *crypto.NewSha256String("e0244db4c02d68ae64dec160310e247bb04e5cb599afb7c14710fbf3f4576c0e"),
		ScheduleVersion:       0,
		Producers:             nil,
		Extensions:            nil,
	}
	hash := block.CalculateId()
	assert.Equal(t, hash.String(), "0000000267f3e2284b482f3afc2e724be1d6cbc1804532ec62d4e7af47c30693")
}
