/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package trie

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"sort"

	eth "github.com/ethereum/go-ethereum/crypto"
)

var (
	// DefaultLeaf is the value that may be passed to Update in order to delete
	// a key from the database.
	DefaultLeaf = Hasher([]byte{0})
)

const (
	// HashLength is the number of bytes in the hash function being used
	// in the trie
	HashLength   = 32
	maxPastTries = 1024
)

// Hash is used to convert a hash into a byte array
type Hash [HashLength]byte

// GetFreshData -
func GetFreshData(size, length int) [][]byte {
	var data [][]byte
	for i := 0; i < size; i++ {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			panic(err)
		}
		data = append(data, Hasher(key)[:length])
	}
	sort.Sort(DataArray(data))
	return data
}

func convNilToBytes(byteArray []byte) []byte {
	if byteArray == nil {
		return []byte{}
	}
	return byteArray
}

func bitIsSet(bits []byte, i int) bool {
	return bits[i/8]&(1<<uint(7-i%8)) != 0
}

func bitSet(bits []byte, i int) {
	bits[i/8] |= 1 << uint(7-i%8)
}

// Hasher is a default hasher to use
// this hasher utilizes Ethereum keccak
func Hasher(data ...[]byte) []byte {
	return eth.Keccak256(data...)
}

// DataArray is for sorting
type DataArray [][]byte

func (d DataArray) Len() int {
	return len(d)
}
func (d DataArray) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d DataArray) Less(i, j int) bool {
	return bytes.Compare(d[i], d[j]) == -1
}

func marshalUint32(v uint32) []byte {
	vv := make([]byte, 4)
	binary.BigEndian.PutUint32(vv, v)
	return vv[:]
}

func unmarshalUint32(v []byte) (uint32, error) {
	if len(v) != 4 {
		return 0, errors.New("invalid")
	}
	vv := binary.BigEndian.Uint32(v)
	return vv, nil
}
