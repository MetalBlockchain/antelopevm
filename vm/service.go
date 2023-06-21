package vm

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	json "encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MetalBlockchain/antelopevm/abi"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/MetalBlockchain/antelopevm/vm/service/chain_api_plugin"
	"github.com/MetalBlockchain/antelopevm/vm/service/history_api_plugin"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/inconshreveable/log15"
)

var (
	errAccountNotFound = fmt.Errorf("account not found")
)

// Service is the API service for this VM
type Service struct{ vm *VM }

type RequestHandler struct {
	handler func(http.ResponseWriter, *http.Request) ([]byte, error)
}

func NewRequestHandler(handler func(http.ResponseWriter, *http.Request) ([]byte, error)) *RequestHandler {
	return &RequestHandler{
		handler: handler,
	}
}

func (req *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	log15.Info("got request", "path", r.RequestURI)

	if data, err := req.handler(w, r); err != nil {
		log15.Error("failed to serve request", "error", err)

		json.NewEncoder(w).Encode(service.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	}
}

func (s *Service) GetInfo(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	version, _ := s.vm.Version(context.Background())
	lastAcceptedId, _ := s.vm.LastAccepted(context.Background())
	lastAccepted, _ := s.vm.getBlock(core.BlockHash(lastAcceptedId))
	info := chain_api_plugin.NewChainInfoResponse(version, lastAccepted, s.vm.controller.ChainId)

	return json.Marshal(info)
}

func (s *Service) GetKeyAccounts(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	response := history_api_plugin.GetKeyAccountsResponse{
		AccountNames: []name.AccountName{name.StringToName("joe"), name.StringToName("eosio")},
	}
	return json.Marshal(response)
}

func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetBlockRequest
	json.NewDecoder(r.Body).Decode(&body)

	session := s.vm.state.CreateSession(false)
	defer session.Discard()

	if val, err := strconv.ParseUint(string(body.BlockNumOrId), 10, 64); err == nil {
		block, err := session.FindBlockByIndex(val)

		if err != nil {
			w.WriteHeader(404)
			return nil, fmt.Errorf("could not parse block num")
		}

		return json.Marshal(service.NewGetBlockResponse(block))
	}

	blockHash, err := hex.DecodeString(string(body.BlockNumOrId))

	if err != nil {
		w.WriteHeader(400)
		return nil, errAccountNotFound
	}

	blockID, err := ids.ToID(blockHash)

	if err != nil {
		w.WriteHeader(400)
		return nil, errAccountNotFound
	}

	block, err := session.FindBlockByHash(core.BlockHash(blockID))

	if err != nil {
		w.WriteHeader(404)
		return nil, errAccountNotFound
	}

	return json.Marshal(service.NewGetBlockResponse(block))
}

func (s *Service) GetBlockInfo(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetBlockInfoRequest
	json.NewDecoder(r.Body).Decode(&body)

	session := s.vm.state.CreateSession(false)
	defer session.Discard()

	block, err := session.FindBlockByIndex(body.BlockNum)

	if err != nil {
		w.WriteHeader(404)
		return nil, errAccountNotFound
	}

	return json.Marshal(service.NewGetBlockInfoResponse(block))
}

func (s *Service) PushTransaction(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var trx core.PackedTransaction

	if err := json.NewDecoder(r.Body).Decode(&trx); err != nil {
		return nil, err
	}

	if ok := s.vm.mempool.Add(&trx); !ok {
		return nil, fmt.Errorf("could not submit trx")
	}

	w.WriteHeader(202)

	return json.Marshal(chain_api_plugin.PushTransactionResults{
		TransactionId: trx.Id.String(),
		Processed: core.TransactionTrace{
			Hash: trx.Id,
			Receipt: core.TransactionReceipt{
				TransactionReceiptHeader: core.TransactionReceiptHeader{
					Status: core.TransactionStatusExecuted,
				},
			},
			ActionTraces: make([]core.ActionTrace, 0),
		},
	})
}

func (s *Service) GetRequiredKeys(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.RequiredKeysRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	authorizationManager := s.vm.controller.GetAuthorizationManager(session)
	keySet := ecc.NewPublicKeySetFromArray(body.AvailableKeys)
	data, err := authorizationManager.GetRequiredKeys(body.Transaction, keySet)

	if err != nil {
		return nil, err
	}

	return json.Marshal(service.RequiredKeysResponse{RequiredKeys: data})
}

func (s *Service) GetTransaction(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetTransactionRequest
	json.NewDecoder(r.Body).Decode(&body)

	hash := core.TransactionIdType(*crypto.NewSha256String(body.TransactionId))
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	trx, err := session.FindTransactionByHash(hash)

	if err != nil {
		w.WriteHeader(400)
		return nil, errAccountNotFound
	}

	lastAcceptedId, _ := s.vm.LastAccepted(context.Background())
	lastAccepted, _ := s.vm.getBlock(core.BlockHash(lastAcceptedId))
	signedTrx, _ := trx.Receipt.Transaction.GetSignedTransaction()
	response := &service.GetTransactionResponse{
		BlockNum:              trx.BlockNum,
		BlockTime:             trx.BlockTime.String(),
		HeadBlockNum:          uint32(lastAccepted.Header.Index),
		Id:                    trx.Hash,
		Irreversible:          true,
		LastIrreversibleBlock: uint32(lastAccepted.Header.Index),
		TransactionNum:        0,
		Traces:                trx.ActionTraces,
		MetaData: service.TransactionMetaData{
			Receipt:     service.TransactionReceipt(trx.Receipt),
			Transaction: *signedTrx,
		},
	}

	for index, trace := range response.Traces {
		if acc, err := session.FindAccountByName(trace.Action.Account); err == nil {
			if len(acc.Abi) > 0 {
				if abi, err := abi.NewABI(acc.Abi); err == nil {
					if data, err := abi.DecodeAction(trace.Action.Name, trace.Action.Data); err == nil {
						parsedData := map[string]interface{}{}

						if err := json.Unmarshal(data, &parsedData); err == nil {
							response.Traces[index].Action.ParsedData = parsedData
						} else {
							log15.Error("err 2", "e", err)
						}
					} else {
						log15.Error("err", "e", err)
					}
				}
			}
		}
	}

	return json.Marshal(response)
}

func (s *Service) GetCodeHash(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetCodeHashRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

	if err != nil {
		return nil, errAccountNotFound
	}

	return json.Marshal(service.GetCodeHashResponse{AccountName: body.AccountName, CodeHash: acc.CodeVersion})
}

func (s *Service) GetAccount(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetAccountRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

	if err != nil {
		return nil, errAccountNotFound
	}

	permissions := make([]*core.Permission, 0)
	iterator := session.FindPermissionsByOwner(acc.Name)
	defer iterator.Close()

	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		if item, err := iterator.Item(); err == nil {
			permissions = append(permissions, item)
		} else {
			return nil, err
		}
	}

	if err != nil {
		return nil, errAccountNotFound
	}

	response := service.GetAccountResponse{
		AccountName: body.AccountName,
		CpuLimit: service.Limit{
			Available:           16346040,
			CurrentUsed:         0,
			LastUsageUpdateTime: core.Now().String(),
			Max:                 16346421,
			Used:                381,
		},
		CpuWeight:         500000,
		Created:           acc.CreationDate.String(),
		CoreLiquidBalance: "1000.0000 SYS",
		HeadBlockNum:      0,
		HeadBlockTime:     core.Now().String(),
		LastCodeUpdate:    acc.LastCodeUpdate.String(),
		NetLimit: service.Limit{
			Available:           88094630,
			CurrentUsed:         0,
			LastUsageUpdateTime: core.Now().String(),
			Max:                 88094878,
			Used:                248,
		},
		NetWeight:   500000,
		Permissions: make([]service.Permission, 0),
		Privileged:  acc.Privileged,
		RamQuota:    525686,
		RamUsage:    5544,
		TotalResources: service.Resources{
			CpuWeight: "50.0000 SYS",
			NetWeight: "50.0000 SYS",
			Owner:     body.AccountName,
			RamBytes:  524286,
		},
	}

	for _, permission := range permissions {
		parent := ""

		for _, p := range permissions {
			if p.ID != permission.ID && permission.Parent == p.ID {
				parent = p.Name.String()
			}
		}

		response.Permissions = append(response.Permissions, service.Permission{
			Parent: parent,
			Name:   permission.Name.String(),
			Auth:   permission.Auth,
		})
	}

	return json.Marshal(response)
}

func (s *Service) GetAbi(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetAbiRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

	if err != nil {
		return nil, errAccountNotFound
	}

	if len(acc.Abi) == 0 {
		response := service.GetAbiResponse{
			AccountName: acc.Name.String(),
		}

		return json.Marshal(response)
	}

	abi, err := abi.NewABI(acc.Abi)

	if err != nil {
		return nil, fmt.Errorf("could not decode ABI")
	}

	response := service.GetAbiResponse{
		AccountName: acc.Name.String(),
		Abi:         *abi,
	}

	return json.Marshal(response)
}

func (s *Service) GetRawAbi(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetRawAbiRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

	if err != nil {
		return nil, errAccountNotFound
	}

	if len(acc.Abi) == 0 {
		response := service.GetRawAbiResponse{
			AccountName: acc.Name.String(),
		}

		return json.Marshal(response)
	}

	rawAbi := base64.StdEncoding.EncodeToString(acc.Abi)
	response := service.GetRawAbiResponse{
		AccountName: acc.Name.String(),
		CodeHash:    acc.CodeVersion,
		AbiHash:     *crypto.NewSha256String("bf13acab1b4bc2676ef6f0afcf1765ab5db3ffa1ac18453a628e1e65fe26e045"),
		Abi:         rawAbi,
	}

	return json.Marshal(response)
}

func (s *Service) GetCurrencyBalance(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetCurrencyBalanceRequest
	json.NewDecoder(r.Body).Decode(&body)
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	results := make([]core.Asset, 0)

	table, err := session.FindTableByCodeScopeTable(body.Code, body.Account, name.StringToName("accounts"))

	if err != nil {
		return json.Marshal(results)
	}

	iterator := session.FindKeyValuesByScope(table.ID)
	defer iterator.Close()

	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item, err := iterator.Item()

		if err == nil {
			asset := core.Asset{}

			if err := rlp.DecodeBytes(item.Value, &asset); err == nil {

				if len(body.Symbol) == 0 || body.Symbol == asset.Symbol.Symbol {
					results = append(results, asset)

					if len(body.Symbol) > 0 {
						break
					}
				}
			}
		}
	}

	return json.Marshal(results)
}

func (s *Service) GetActions(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetActionsRequest
	json.NewDecoder(r.Body).Decode(&body)

	return json.Marshal(
		service.GetActionsResponse{
			Actions:                  []string{},
			HeadBlockNum:             0,
			LastIrreversibleBlockNum: 0,
		},
	)
}

func (s *Service) GetTableRows(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetTableRowsRequest
	json.NewDecoder(r.Body).Decode(&body)
	response := &service.GetTableRowsResponse{
		More: false,
		Rows: make([]service.TableRow, 0),
	}
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	acc, err := session.FindAccountByName(body.Code)

	if err != nil {
		return nil, fmt.Errorf("account with name %s does not exist", body.Code)
	}

	abi, err := abi.NewABI(acc.Abi)

	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI")
	}

	tableDef := abi.TableForName(body.Table)

	if tableDef == nil {
		return nil, fmt.Errorf("table %s is not specified in the ABI", body.Table)
	}

	table, err := session.FindTableByCodeScopeTable(body.Code, body.Scope, body.Table)

	if err != nil {
		return json.Marshal(response)
	}

	keyValues := make([]*core.KeyValue, 0)
	keyValueIterator := session.FindKeyValuesByScope(table.ID)
	defer keyValueIterator.Close()

	for keyValueIterator.Rewind(); keyValueIterator.Valid(); keyValueIterator.Next() {
		if item, err := keyValueIterator.Item(); err == nil {
			keyValues = append(keyValues, item)
		}
	}

	for _, keyValue := range keyValues {
		structDef := abi.StructForName(tableDef.Type)

		if data, err := abi.DecodeStruct(structDef.Name, keyValue.Value); err == nil {
			parsedData := map[string]interface{}{}

			if err := json.Unmarshal(data, &parsedData); err == nil {
				response.Rows = append(response.Rows, service.TableRow{
					Data:  parsedData,
					Payer: keyValue.Payer.String(),
				})
			} else {
				log15.Error("failed to encode json", "err", err)
			}
		} else {
			log15.Error("failed to decode struct", "err", err)
		}
	}

	return json.Marshal(response)
}

func (s *Service) GetCurrencyStats(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var body service.GetCurrencyStatsRequest
	json.NewDecoder(r.Body).Decode(&body)
	response := service.GetCurrencyStatsResponse{}
	session := s.vm.state.CreateSession(false)
	defer session.Discard()
	symbol, err := core.StringToSymbol(0, strings.ToUpper(body.Symbol))

	if err != nil {
		return nil, fmt.Errorf("invalid symbol")
	}

	scope := symbol >> 8
	table, err := session.FindTableByCodeScopeTable(body.Code, name.Name(scope), name.StringToName("stat"))

	if err != nil {
		return nil, fmt.Errorf("could not find currency stats")
	}

	iterator := session.FindKeyValuesByScope(table.ID)
	defer iterator.Close()

	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		obj, err := iterator.Item()

		if err != nil {
			continue
		}

		ds := rlp.NewDecoder(obj.Value)
		result := service.CurrencyStats{}

		if err := ds.Decode(&result.Supply); err != nil {
			continue
		}

		if err := ds.Decode(&result.MaxSupply); err != nil {
			continue
		}

		if err := ds.Decode(&result.Issuer); err != nil {
			continue
		}

		response[result.Supply.Symbol.Name()] = result
	}

	return json.Marshal(response)
}
