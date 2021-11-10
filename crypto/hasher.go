package crypto

import (
	eth "github.com/ethereum/go-ethereum/crypto"
)

// HashFunc256 is our universal hash function with 32-byte digest.
// Changing this function will change *all* locations.
func Hasher(msg ...[]byte) []byte {
	return eth.Keccak256(msg...)
}
