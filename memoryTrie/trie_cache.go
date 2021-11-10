/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package trie

import (
	"sync"
)

type cacheDB struct {
	// lock for cacheDB
	lock sync.RWMutex
	// updatedNodes that will be flushed to disk
	updatedNodes map[Hash][][]byte
	// updatedMux is a lock for updatedNodes
	updatedMux sync.RWMutex

	store map[Hash][]byte
}

// commit stores the updated nodes to disk.
func (db *cacheDB) commit(s *SMT) error {
	for key, batch := range db.updatedNodes {
		var node []byte
		err := db.setNodeDB(append(node, key[:]...), db.serializeBatch(batch))
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *cacheDB) serializeBatch(batch [][]byte) []byte {
	serialized := make([]byte, 4)
	if batch[0][0] == 1 {
		// the batch node is a shortcut
		bitSet(serialized, 31)
	}
	for i := 1; i < 31; i++ {
		if len(batch[i]) != 0 {
			bitSet(serialized, i-1)
			serialized = append(serialized, batch[i]...)
		}
	}
	return serialized
}
