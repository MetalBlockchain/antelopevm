package block_test

import (
	"fmt"
	"testing"

	"github.com/MetalBlockchain/antelopevm/chain/block"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/chain/time"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCalculateId(t *testing.T) {
	time, err := time.FromIsoString("2020-04-22T17:00:00.000")
	if err != nil {
		panic(err)
	}
	fmt.Printf("block.NewBlockTimeStampFromTimePoint(time): %v\n", block.NewBlockTimeStampFromTimePoint(time))
	block := block.BlockHeader{
		Timestamp:             block.NewBlockTimeStampFromTimePoint(time),
		Producer:              name.StringToName(""),
		Confirmed:             1,
		Previous:              *crypto.NewSha256String("0000000000000000000000000000000000000000000000000000000000000000"),
		TransactionMerkleRoot: *crypto.NewSha256String("0000000000000000000000000000000000000000000000000000000000000000"),
		ActionMerkleRoot:      *crypto.NewSha256String("384da888112027f0321850a169f737c33e53b388aad48b5adace4bab97f437e0"),
		ScheduleVersion:       0,
		NewProducers:          nil,
		Extensions:            nil,
	}
	hash := block.CalculateId()
	fmt.Printf("block.BlockNum(): %v\n", block.BlockNum())
	assert.Equal(t, hash.String(), "000000018421bd47ce23d4c47706e0bb98604157afedc67d56d05c82d5aa10c5")
}
