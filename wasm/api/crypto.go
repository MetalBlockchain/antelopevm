package api

import (
	"bytes"
	"strings"

	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

func init() {
	Functions["assert_recover_key"] = assertRecoverKey
	Functions["recover_key"] = recoverKey
	Functions["assert_sha256"] = assertSha256
	Functions["assert_sha1"] = assertSha1
	Functions["assert_sha512"] = assertSha512
	Functions["assert_ripemd160"] = assertRipemd160
	Functions["sha256"] = sha256
	Functions["sha1"] = sha1
	Functions["sha512"] = sha512
	Functions["ripemd160"] = ripemd160
}

func assertRecoverKey(context Context) interface{} {
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

func recoverKey(context Context) interface{} {
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

func assertSha256(context Context) interface{} {
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

func assertSha1(context Context) interface{} {
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

func assertSha512(context Context) interface{} {
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

func assertRipemd160(context Context) interface{} {
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

func sha1(context Context) interface{} {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha1()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:20])
	}
}

func sha256(context Context) interface{} {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha256()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:32])
	}
}

func sha512(context Context) interface{} {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewSha512()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:64])
	}
}

func ripemd160(context Context) interface{} {
	return func(data uint32, dataLength uint32, hash uint32) {
		dataBytes := context.ReadMemory(data, dataLength)
		s := crypto.NewRipemd160()
		s.Write(dataBytes)
		calculatedHash := s.Sum(nil)

		context.WriteMemory(hash, calculatedHash[0:20])
	}
}
