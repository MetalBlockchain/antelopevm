package vm

import (
	"context"
	"encoding/hex"
	json "encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/inconshreveable/log15"
)

// Service is the API service for this VM
type Service struct{ vm *VM }

type RequestHandler struct {
	handler func(http.ResponseWriter, *http.Request) error
}

func NewRequestHandler(handler func(http.ResponseWriter, *http.Request) error) *RequestHandler {
	return &RequestHandler{
		handler: handler,
	}
}

func (req *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")

	if err := req.handler(w, r); err != nil {
		json.NewEncoder(w).Encode(service.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
}

func (s *Service) GetInfo(w http.ResponseWriter, r *http.Request) error {
	version, _ := s.vm.Version(context.Background())
	lastAcceptedId, _ := s.vm.LastAccepted(context.Background())
	lastAccepted, _ := s.vm.getBlock(lastAcceptedId)
	info := service.NewChainInfoResponse(version, lastAccepted, s.vm.controller.ChainId)

	json.NewEncoder(w).Encode(info)

	return nil
}

func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) error {
	var body service.GetBlockRequest
	json.NewDecoder(r.Body).Decode(&body)

	log15.Info("got request", "data", body)

	if val, err := strconv.ParseUint(string(body.BlockNumOrId), 10, 64); err == nil {
		block, err := s.vm.state.GetBlockByIndex(val)

		if err != nil {
			w.WriteHeader(404)
			return nil
		}

		json.NewEncoder(w).Encode(service.NewGetBlockResponse(block))
		return nil
	}

	blockHash, err := hex.DecodeString(string(body.BlockNumOrId))

	if err != nil {
		w.WriteHeader(400)
		return nil
	}

	blockID, err := ids.ToID(blockHash)

	if err != nil {
		w.WriteHeader(400)
		return nil
	}

	block, err := s.vm.state.GetBlock(blockID)

	if err != nil {
		w.WriteHeader(404)
		return nil
	}

	json.NewEncoder(w).Encode(service.NewGetBlockResponse(block))

	return nil
}

func (s *Service) GetBlockInfo(w http.ResponseWriter, r *http.Request) error {
	var body service.GetBlockInfoRequest
	json.NewDecoder(r.Body).Decode(&body)

	block, err := s.vm.state.GetBlockByIndex(body.BlockNum)

	if err != nil {
		w.WriteHeader(404)
		return nil
	}

	json.NewEncoder(w).Encode(service.NewGetBlockInfoResponse(block))
	return nil
}

func (s *Service) PushTransaction(w http.ResponseWriter, r *http.Request) error {
	log15.Info("PushTransaction")
	var trx core.PackedTransaction

	if err := json.NewDecoder(r.Body).Decode(&trx); err != nil {
		return err
	}

	if ok := s.vm.mempool.Add(&trx); !ok {
		return fmt.Errorf("could not submit trx")
	}

	return nil
}

func (s *Service) GetRequiredKeys(w http.ResponseWriter, r *http.Request) error {
	var body service.RequiredKeysRequest
	json.NewDecoder(r.Body).Decode(&body)
	authorizationManager := s.vm.controller.GetAuthorizationManager(s.vm.state)
	data, err := authorizationManager.GetRequiredKeys(body.Transaction, body.AvailableKeys)

	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(service.RequiredKeysResponse{RequiredKeys: data})
	return nil
}

func (s *Service) GetTransaction(w http.ResponseWriter, r *http.Request) error {
	var body service.GetTransactionRequest
	json.NewDecoder(r.Body).Decode(&body)

	hash := core.TransactionIdType(*crypto.NewSha256String(body.TransactionId))
	trx, err := s.vm.state.GetTransaction(hash)

	log15.Info("requested", "id", hash.String(), "req", body.TransactionId)

	if err != nil {
		w.WriteHeader(400)
		return nil
	}

	lastAcceptedId, _ := s.vm.LastAccepted(context.Background())
	lastAccepted, _ := s.vm.getBlock(lastAcceptedId)
	signedTrx, _ := trx.Receipt.Transaction.GetSignedTransaction()
	response := &service.GetTransactionResponse{
		BlockNum:              trx.BlockNum,
		BlockTime:             trx.BlockTime.String(),
		HeadBlockNum:          uint32(lastAccepted.Index),
		Id:                    trx.Id,
		Irreversible:          true,
		LastIrreversibleBlock: uint32(lastAccepted.Index),
		TransactionNum:        0,
		Traces:                trx.ActionTraces,
		MetaData: service.TransactionMetaData{
			Receipt:     service.TransactionReceipt(trx.Receipt),
			Transaction: *signedTrx,
		},
	}

	json.NewEncoder(w).Encode(response)
	return nil
}
