/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package trie

// The Package Trie implements a sparse merkle trie.

import (
	"bytes"
)

// Get fetches the value of a key by going down the current trie root.
func (s *SMT) Get(key []byte) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.get(s.Root, key, nil, 0, s.TrieHeight)
}

// get fetches the value of a key in a given trie
// defined by root
func (s *SMT) get(root []byte, key []byte, batch [][]byte, iBatch, height int) ([]byte, error) {
	if len(root) == 0 {
		return nil, nil
	}
	if height == 0 {
		return root[:HashLength], nil
	}
	// Fetch the children of the node
	batch, iBatch, lnode, rnode, isShortcut, err := s.loadChildren(root, height, iBatch, batch)
	if err != nil {
		return nil, err
	}
	if isShortcut {
		if bytes.Equal(lnode[:HashLength], key) {
			return rnode[:HashLength], nil
		}
		return nil, nil
	}
	if bitIsSet(key, s.TrieHeight-height) {
		return s.get(rnode, key, batch, 2*iBatch+2, height-1)
	}
	return s.get(lnode, key, batch, 2*iBatch+1, height-1)
}
