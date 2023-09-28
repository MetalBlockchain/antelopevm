package utils

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type Pair[T any, S any] struct {
	First  T
	Second S
}

func MakeCanonicalLeft(val crypto.Sha256) crypto.Sha256 {
	canonicalLeft := val
	canonicalLeft.Hash[0] &= 0xFFFFFFFFFFFFFF7F
	return canonicalLeft
}

func MakeCanonicalRight(val crypto.Sha256) crypto.Sha256 {
	canonicalRight := val
	canonicalRight.Hash[0] |= 0x0000000000000080
	return canonicalRight
}

func MakeCanonicalPair(l crypto.Sha256, r crypto.Sha256) Pair[crypto.Sha256, crypto.Sha256] {
	return Pair[crypto.Sha256, crypto.Sha256]{
		First:  MakeCanonicalLeft(l),
		Second: MakeCanonicalRight(r),
	}
}

func Merkle(hashes []crypto.Sha256) crypto.Sha256 {
	if len(hashes) == 0 {
		return crypto.NewSha256Nil()
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		for i := 0; i < len(hashes)/2; i++ {
			hashes[i] = *crypto.Hash256(MakeCanonicalPair(hashes[2*i], hashes[(2*i)+1]))
		}

		hashes = hashes[:len(hashes)/2]
	}

	return hashes[0]
}
