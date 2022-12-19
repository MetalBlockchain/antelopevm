package vm

import (
	"context"
	"encoding/hex"
	json2 "encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/chain/types"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/metalgo/ids"
)

var (
	errBadData               = errors.New("data must be base 58 repr. of 32 bytes")
	errNoSuchBlock           = errors.New("couldn't get block from database. Does it exist?")
	errCannotGetLastAccepted = errors.New("problem getting last accepted")
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
		json2.NewEncoder(w).Encode(service.ErrorResponse{
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

	json2.NewEncoder(w).Encode(info)

	return nil
}

func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) error {
	var body service.GetBlockRequest
	json2.NewDecoder(r.Body).Decode(&body)

	if val, err := strconv.ParseUint(body.BlockNumOrId, 10, 64); err == nil {
		block, err := s.vm.state.GetBlockByIndex(val)

		if err != nil {
			w.WriteHeader(404)
			return nil
		}

		json2.NewEncoder(w).Encode(service.NewGetBlockResponse(block))
		return nil
	}

	blockHash, err := hex.DecodeString(body.BlockNumOrId)

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

	json2.NewEncoder(w).Encode(service.NewGetBlockResponse(block))

	return nil
}

func (s *Service) GetBlockInfo(w http.ResponseWriter, r *http.Request) error {
	var body service.GetBlockInfoRequest
	json2.NewDecoder(r.Body).Decode(&body)

	block, err := s.vm.state.GetBlockByIndex(body.BlockNum)

	if err != nil {
		w.WriteHeader(404)
		return nil
	}

	json2.NewEncoder(w).Encode(service.NewGetBlockInfoResponse(block))
	return nil
}

func (s *Service) PushTransaction(w http.ResponseWriter, r *http.Request) error {
	var trx types.PackedTransaction

	if err := json2.NewDecoder(r.Body).Decode(&trx); err != nil {
		return err
	}

	s.vm.mempool.Add(&trx)

	return nil
}

func (s *Service) GetRequiredKeys(w http.ResponseWriter, r *http.Request) error {
	var body service.RequiredKeysRequest
	json2.NewDecoder(r.Body).Decode(&body)

	data, err := s.vm.controller.Authorization.GetRequiredKeys(body.Transaction, body.AvailableKeys)

	if err != nil {
		return err
	}

	json2.NewEncoder(w).Encode(service.RequiredKeysResponse{RequiredKeys: data})
	return nil
}
