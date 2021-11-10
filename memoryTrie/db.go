package trie

func (db *cacheDB) getNodeDB(key []byte) (val []byte, err error) {
	var node Hash
	copy(node[:], key)
	value := db.store[node]
	if value == nil {
		return nil, nil
	}
	buf := make([]byte, len(value))
	copy(buf[:], value)
	return buf, nil
}

func (db *cacheDB) setNodeDB(key, value []byte) error {
	var node Hash
	copy(node[:], key)
	db.store[node] = value
	return nil
}

func (db *cacheDB) deleteNodeDB(key []byte) error {
	var node Hash
	copy(node[:], key)
	delete(db.store, node)
	return nil
}
