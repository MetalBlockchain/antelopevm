package api

import (
	"bytes"
	"strings"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

func GetCryptoFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["assert_recover_key"] = assertRecoverKey(context)
	functions["recover_key"] = recoverKey(context)
	functions["assert_sha256"] = assertSha256(context)
	functions["assert_sha1"] = assertSha1(context)
	functions["assert_sha512"] = assertSha512(context)
	functions["assert_ripemd160"] = assertRipemd160(context)
	functions["sha256"] = sha256(context)
	functions["sha1"] = sha1(context)
	functions["sha512"] = sha512(context)
	functions["ripemd160"] = ripemd160(context)

	return functions
}

func assertRecoverKey(context Context) func(uint32, uint32, uint32, uint32, uint32) {
	return func(digest uint32, signature uint32, signatureLength uint32, publicKey uint32, publicKeyLength uint32) {
		digestBytes := context.ReadMemory(digest, 32)
		signatureBytes := context.ReadMemory(signature, signatureLength)
		publicKeyBytes := context.ReadMemory(publicKey, publicKeyLength)
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

func recoverKey(context Context) func(uint32, uint32, uint32, uint32, uint32) int32 {
	return func(digest uint32, signature uint32, signatureLength uint32, publicKey uint32, publicKeyLength uint32) int32 {
		digestBytes := context.ReadMemory(digest, 32)
		signatureBytes := context.ReadMemory(signature, signatureLength)

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

		context.WriteMemory(publicKey, encoded[0:bufferSize])

		return int32(bufferSize)
	}
}

func assertSha256(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		hashBytes := context.ReadMemory(hash, 32)
		s := crypto.NewSha256()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha256 hash mismatch")
		}
	}
}

func assertSha1(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		hashBytes := context.ReadMemory(hash, 20)
		s := crypto.NewSha1()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha1 hash mismatch")
		}
	}
}

func assertSha512(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		hashBytes := context.ReadMemory(hash, 64)
		s := crypto.NewSha512()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("sha512 hash mismatch")
		}
	}
}

func assertRipemd160(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		hashBytes := context.ReadMemory(hash, 20)
		s := crypto.NewRipemd160()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		if !bytes.Equal(calculatedHash, hashBytes) {
			panic("ripemd160 hash mismatch")
		}
	}
}

func sha1(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha1()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:20])
	}
}

func sha256(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha256()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:32])
	}
}

func sha512(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha512()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:64])
	}
}

func ripemd160(context Context) func(uint32, uint32, uint32) {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewRipemd160()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:20])
	}
}
