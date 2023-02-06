package wasm

import (
	"bytes"
	"strings"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	log "github.com/inconshreveable/log15"
)

func GetCryptoFunctions(context *ExecutionContext) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["assert_recover_key"] = assertRecoverKey(context)
	functions["recover_key"] = recoverKey(context)
	functions["assert_sha256"] = assertSha256(context)
	functions["assert_sha1"] = assertSha1(context)
	functions["assert_sha512"] = assertSha512(context)
	functions["assert_ripemd160"] = assertRipemd160(context)
	functions["sha1"] = sha1(context)
	functions["sha256"] = sha256(context)
	functions["sha512"] = sha512(context)
	functions["ripemd160"] = ripemd160(context)

	return functions
}

func assertRecoverKey(context *ExecutionContext) func(uint32, uint32, uint32, uint32, uint32) {
	return func(digest uint32, signature uint32, signatureLength uint32, publicKey uint32, publicKeyLength uint32) {
		log.Info("assert_recover_key", "digest", digest, "signature", signature, "signatureLength", signatureLength, "publicKey", publicKey, "publicKeyLength", publicKeyLength)

		digestBytes, ok := context.module.Memory().Read(context.context, digest, 32)

		if !ok {
			panic("could not read digest bytes")
		}

		signatureBytes, ok := context.module.Memory().Read(context.context, signature, signatureLength)

		if !ok {
			panic("could not read signature bytes")
		}

		publicKeyBytes, ok := context.module.Memory().Read(context.context, publicKey, publicKeyLength)

		if !ok {
			panic("could not read public key bytes")
		}

		sig := ecc.NewSigNil()
		pub := ecc.NewPublicKeyNil()
		rlp.DecodeBytes(signatureBytes, sig)
		rlp.DecodeBytes(publicKeyBytes, pub)

		check, err := sig.PublicKey(digestBytes)

		if err != nil {
			panic("could not form public key from digest bytes")
		}

		if strings.Compare(check.String(), pub.String()) != 0 {
			panic("expected different key than recovered key")
		}
	}
}

func recoverKey(context *ExecutionContext) func(uint32, uint32, uint32, uint32, uint32) int32 {
	return func(digest uint32, signature uint32, signatureLength uint32, publicKey uint32, publicKeyLength uint32) int32 {
		log.Info("recover_key", "digest", digest, "signature", signature, "signatureLength", signatureLength, "publicKey", publicKey, "publicKeyLength", publicKeyLength)

		digestBytes, ok := context.module.Memory().Read(context.context, digest, 32)

		if !ok {
			panic("could not read digest bytes")
		}

		signatureBytes, ok := context.module.Memory().Read(context.context, signature, signatureLength)

		if !ok {
			panic("could not read signature bytes")
		}

		sig := ecc.NewSigNil()
		rlp.DecodeBytes(signatureBytes, sig)
		check, err := sig.PublicKey(digestBytes)

		if err != nil {
			panic("could not form public key from digest bytes")
		}

		encoded, err := rlp.EncodeToBytes(check)

		if err != nil {
			panic("could not encode public key to rlp")
		}

		bufferSize := len(encoded)

		if bufferSize > int(publicKeyLength) {
			bufferSize = int(publicKeyLength)
		}

		context.module.Memory().Write(context.context, publicKey, encoded[0:bufferSize])

		return int32(bufferSize)
	}
}

func assertSha256(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("assert_sha256", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		hashBytes, ok := context.module.Memory().Read(context.context, hash, 32)

		if !ok {
			panic("could not read hash bytes")
		}

		s := crypto.NewSha256()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha256 hash mismatch")
		}
	}
}

func assertSha1(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("assert_sha1", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		hashBytes, ok := context.module.Memory().Read(context.context, hash, 20)

		if !ok {
			panic("could not read hash bytes")
		}

		s := crypto.NewSha1()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha1 hash mismatch")
		}
	}
}

func assertSha512(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("assert_sha512", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		hashBytes, ok := context.module.Memory().Read(context.context, hash, 64)

		if !ok {
			panic("could not read hash bytes")
		}

		s := crypto.NewSha512()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha512 hash mismatch")
		}
	}
}

func assertRipemd160(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("assert_ripemd160", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		hashBytes, ok := context.module.Memory().Read(context.context, hash, 20)

		if !ok {
			panic("could not read hash bytes")
		}

		s := crypto.NewRipemd160()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("ripemd160 hash mismatch")
		}
	}
}

func sha1(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("sha1", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		s := crypto.NewSha1()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if ok := context.module.Memory().Write(context.context, hash, calculatedHash[0:20]); !ok {
			panic("could not write sha1 hash to memory")
		}
	}
}

func sha256(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("sha256", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		s := crypto.NewSha256()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if ok := context.module.Memory().Write(context.context, hash, calculatedHash[0:32]); !ok {
			panic("could not write sha256 hash to memory")
		}
	}
}

func sha512(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("sha512", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		s := crypto.NewSha512()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if ok := context.module.Memory().Write(context.context, hash, calculatedHash[0:64]); !ok {
			panic("could not write sha512 hash to memory")
		}
	}
}

func ripemd160(context *ExecutionContext) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		log.Info("ripemd160", "data", data, "dataLength", dataLength, "hash", hash)

		dataBytes, ok := context.module.Memory().Read(context.context, data, dataLength)

		if !ok {
			panic("could not read data bytes")
		}

		s := crypto.NewRipemd160()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if ok := context.module.Memory().Write(context.context, hash, calculatedHash[0:20]); !ok {
			panic("could not write ripemd160 hash to memory")
		}
	}
}
