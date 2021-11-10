/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package trie

import (
	"bytes"
	"runtime"

	"time"

	"fmt"
	"testing"
)

func prefixFn() []byte {
	return []byte("aa")
}

func TestSmtEmptyTrie(t *testing.T) {
	smt := NewSMT()
	if !bytes.Equal([]byte{}, smt.Root) {
		t.Fatal("empty trie root hash not correct")
	}
}

func TestParseBatch(t *testing.T) {
	smt := NewSMT()
	for i := 0; i < 100; i++ {
		for j := 0; j < 2000; j++ {
			temp := i % 32
			var data []byte
			if i > 32 {
				for k := 0; k < (i / 32); k++ {
					data = append(data, GetFreshData(1, 32)[0]...)
				}
				data = append(data, GetFreshData(1, temp)[0]...)
			} else {
				data = GetFreshData(1, i)[0]
			}
			smt.parseBatch(data)
		}
	}
}

func TestSmtUpdateAndGet(t *testing.T) {
	smt := NewSMT()
	// Add data to empty trie
	keys := GetFreshData(10, 32)
	values := GetFreshData(10, 32)
	ch := make(chan mresult, 1)
	defer close(ch)
	smt.update(smt.Root, keys, values, nil, 0, smt.TrieHeight, ch)
	res := <-ch
	root := res.update

	// Check all keys have been stored
	for i, key := range keys {
		value, _ := smt.get(root, key, nil, 0, smt.TrieHeight)
		if !bytes.Equal(values[i], value) {
			t.Fatal("value not updated")
		}
	}

	// Append to the trie
	newKeys := GetFreshData(5, 32)
	newValues := GetFreshData(5, 32)
	ch = make(chan mresult, 1)
	defer close(ch)
	smt.update(root, newKeys, newValues, nil, 0, smt.TrieHeight, ch)
	res = <-ch
	newRoot := res.update
	if bytes.Equal(root, newRoot) {
		t.Fatal("trie not updated")
	}
	for i, newKey := range newKeys {
		newValue, _ := smt.get(newRoot, newKey, nil, 0, smt.TrieHeight)
		if !bytes.Equal(newValues[i], newValue) {
			t.Fatal("failed to get value")
		}
	}
	// Check old keys are still stored
	for i, key := range keys {
		value, _ := smt.get(newRoot, key, nil, 0, smt.TrieHeight)
		if !bytes.Equal(values[i], value) {
			t.Fatal("failed to get value")
		}
	}
}

func TestSmtPublicUpdateAndGet(t *testing.T) {
	smt := NewSMT()
	// Add data to empty trie
	keys := GetFreshData(5, 32)
	values := GetFreshData(5, 32)
	root, err := smt.Update(keys, values)
	if err != nil {
		t.Fatal(err)
	}

	// Check all keys have been stored
	for i, key := range keys {
		value, _ := smt.Get(key)
		if !bytes.Equal(values[i], value) {
			t.Fatal("trie not updated")
		}
	}
	if !bytes.Equal(root, smt.Root) {
		t.Fatal("Root not stored")
	}

	newValues := GetFreshData(5, 32)
	_, err = smt.Update(keys, newValues)
	if err == nil {
		t.Fatal("multiple updates don't cause an error")
	}
}

func TestSmtDelete(t *testing.T) {
	smt := NewSMT()
	// Add data to empty trie
	keys := GetFreshData(10, 32)
	values := GetFreshData(10, 32)
	ch := make(chan mresult, 1)
	defer close(ch)
	smt.update(smt.Root, keys, values, nil, 0, smt.TrieHeight, ch)
	res := <-ch
	root := res.update
	value, _ := smt.get(root, keys[0], nil, 0, smt.TrieHeight)
	if !bytes.Equal(values[0], value) {
		t.Fatal("trie not updated")
	}

	// Delete from trie
	// To delete a key, just set it's value to Default leaf hash.
	ch = make(chan mresult, 1)
	defer close(ch)
	smt.update(root, keys[0:1], [][]byte{DefaultLeaf}, nil, 0, smt.TrieHeight, ch)
	// smt.Commit(txn)
	res = <-ch
	newRoot := res.update
	newValue, _ := smt.get(newRoot, keys[0], nil, 0, smt.TrieHeight)
	if len(newValue) != 0 {
		t.Fatal("Failed to delete from trie")
	}
	// Remove deleted key from keys and check root with a clean trie.
	smt2 := NewSMT()
	ch = make(chan mresult, 1)
	defer close(ch)
	smt2.update(smt2.Root, keys[1:], values[1:], nil, 0, smt.TrieHeight, ch)
	// smt2.Commit(txn)
	res = <-ch
	cleanRoot := res.update
	if !bytes.Equal(newRoot, cleanRoot) {
		t.Fatal("roots mismatch")
	}

	//Empty the trie
	var newValues [][]byte
	for i := 0; i < 10; i++ {
		newValues = append(newValues, DefaultLeaf)
	}
	ch = make(chan mresult, 1)
	defer close(ch)
	smt.update(root, keys, newValues, nil, 0, smt.TrieHeight, ch)
	res = <-ch
	root = res.update
	if len(root) != 0 {
		t.Fatal("empty trie root hash not correct")
	}
	// Test deleting an already empty key
	smt = NewSMT()
	keys = GetFreshData(2, 32)
	values = GetFreshData(2, 32)
	root, _ = smt.Update(keys, values)
	key0 := make([]byte, 32, 32)
	key1 := make([]byte, 32, 32)
	smt.Update([][]byte{key0, key1}, [][]byte{DefaultLeaf, DefaultLeaf})
	if !bytes.Equal(root, smt.Root) {
		t.Fatal("deleting a default key shouldnt' modify the tree")
	}
}

// test updating and deleting at the same time
func TestTrieUpdateAndDelete(t *testing.T) {
	smt := NewSMT()
	keys := GetFreshData(2, 32)
	values := GetFreshData(2, 32)
	_, err := smt.Update(keys, values)
	if err != nil {
		t.Fatal(err)
	}

	vv, err := smt.Get(keys[0])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(vv, values[0]) {
		t.Fatal("key not inserted after commit")
	}

	newvalues := [][]byte{}
	newvalues = append(newvalues, DefaultLeaf)
	newkeys := [][]byte{}
	newkeys = append(newkeys, keys[0])
	root2, err := smt.Update(newkeys, newvalues)
	if err != nil {
		t.Fatal(err)
	}

	vvv, err := smt.Get(keys[0])
	if err != nil {
		t.Fatal(err)
	}
	if vvv != nil {
		t.Fatal("key not deleted")
	}

	var node Hash
	copy(node[:], root2)

	keys = GetFreshData(2, 32)
	values = GetFreshData(2, 32)
	_, err = smt.Update(keys, values)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifySubtree(t *testing.T) {
	smt := NewSMT()
	keys := GetFreshData(7, 32)
	values := GetFreshData(7, 32)
	smt.Update(keys, values)
	var node Hash
	copy(node[:], smt.Root)
	batch, err := smt.loadBatch(smt.Root)
	if err != nil {
		t.Fatal(err)
	}
	_, res := smt.verifyBatch(batch, 0, 4, 252, smt.Root, false)
	if res != true {
		t.Fatal("the sub tree verification did not succeed")
	}
}

func TestTrieMerkleProof(t *testing.T) {
	smt := NewSMT()
	keys := GetFreshData(10, 32)
	values := GetFreshData(10, 32)
	smt.Update(keys, values)
	for i, key := range keys {
		ap, _, k, v, _ := smt.MerkleProof(key)
		if !smt.VerifyInclusion(ap, key, values[i]) {
			t.Fatalf("failed to verify inclusion proof")
		}
		if !bytes.Equal(key, k) && !bytes.Equal(values[i], v) {
			t.Fatalf("merkle proof didnt return the correct key-value pair")
		}
	}
	emptyKey := Hasher([]byte("non-member"))
	ap, included, proofKey, proofValue, _ := smt.MerkleProof(emptyKey)
	if included {
		t.Fatalf("failed to verify non inclusion proof")
	}
	if !smt.VerifyNonInclusion(ap, emptyKey, proofValue, proofKey) {
		t.Fatalf("failed to verify non inclusion proof")
	}
}

func TestTrieMerkleProofCompressed(t *testing.T) {
	smt := NewSMT()
	// Add data to empty trie
	keys := GetFreshData(10, 32)
	values := GetFreshData(10, 32)
	smt.Update(keys, values)

	for i, key := range keys {
		bitmap, ap, length, _, k, v, _ := smt.MerkleProofCompressed(key)
		if !smt.VerifyInclusionC(bitmap, key, values[i], ap, length) {
			t.Fatalf("failed to verify inclusion proof")
		}
		if !bytes.Equal(key, k) && !bytes.Equal(values[i], v) {
			t.Fatalf("merkle proof didnt return the correct key-value pair")
		}
	}
	emptyKey := Hasher([]byte("non-member"))
	bitmap, ap, length, included, proofKey, proofValue, _ := smt.MerkleProofCompressed(emptyKey)
	if included {
		t.Fatalf("failed to verify non inclusion proof")
	}
	if !smt.VerifyNonInclusionC(ap, length, bitmap, emptyKey, proofValue, proofKey) {
		t.Fatalf("failed to verify non inclusion proof")
	}
}

func TestSmtCommit(t *testing.T) {
	smt := NewSMT()
	keys := GetFreshData(32, 32)
	values := GetFreshData(32, 32)
	smt.Update(keys, values)
	for i := range keys {
		value, _ := smt.Get(keys[i])
		if !bytes.Equal(values[i], value) {
			t.Fatal("failed to get value in committed db")
		}
	}

	// test loading a shortcut batch
	smt = NewSMT()
	keys = GetFreshData(1, 32)
	values = GetFreshData(1, 32)
	smt.Update(keys, values)
	value, _ := smt.Get(keys[0])
	if !bytes.Equal(values[0], value) {
		t.Fatal("failed to get value in committed db")
	}
}

/*
 func TestDoubleUpdate(t *testing.T) {
	 dir, err := ioutil.TempDir("", "badger-test")
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer func() {
		 if err := os.RemoveAll(dir); err != nil {
			 t.Fatal(err)
		 }
	 }()
	 opts := badger.DefaultOptions(dir)
	 db, err := badger.Open(opts)
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer db.Close()
	 err = db.Update(func(outerTxn *badger.Txn) error {
		 err := outerTxn.Set([]byte("b"), []byte("outer"))
		 if err != nil {
			 return err
		 }
		 err = outerTxn.Set([]byte("a"), []byte("outer"))
		 if err != nil {
			 return err
		 }
		 _, err = outerTxn.Get([]byte("b"))
		 if err != nil {
			 return err
		 }
		 return db.Update(func(innerTxn *badger.Txn) error {
			 err := innerTxn.Set([]byte("c"), []byte("inner"))
			 if err != nil {
				 return err
			 }
			 return nil
		 })
	 })
	 if err != nil {
		 t.Error(err)
	 }
	 err = db.View(func(txn *badger.Txn) error {
		 item, err := txn.Get([]byte("a"))
		 if err != nil {
			 return err
		 }
		 value, err := item.ValueCopy(nil)
		 if err != nil {
			 return err
		 }
		 if !bytes.Equal(value, []byte("outer")) {
			 t.Error("a not equal to inner")
		 }
		 return nil
	 })
	 if err != nil {
		 t.Error(err)
	 }
	 err = db.View(func(txn *badger.Txn) error {
		 item, err := txn.Get([]byte("c"))
		 if err != nil {
			 return err
		 }
		 value, err := item.ValueCopy(nil)
		 if err != nil {
			 return err
		 }
		 if !bytes.Equal(value, []byte("inner")) {
			 t.Error("a not equal to inner")
		 }
		 return nil
	 })
	 if err != nil {
		 t.Error(err)
	 }
 }

*/

func TestSmtRaisesError(t *testing.T) {
	smt := NewSMT()
	// Add data to empty trie
	keys := GetFreshData(10, 32)
	values := GetFreshData(10, 32)
	smt.Update(keys, values)
	smt.db.updatedNodes = make(map[Hash][][]byte)
	smt.loadDefaultHashes()
	// Check errors are raised is a keys is not in cache nor db
	for _, key := range keys {
		_, err := smt.Get(key)
		if err == nil {
			t.Fatal("Error not created if database doesnt have a node")
		}
	}
	//_, _, err := smt.MerkleProofCompressed( keys[0])
	//if err == nil {
	//	t.Fatal("Error not created if database doesnt have a node")
	//}
	_, err := smt.Update(keys, values)
	if err == nil {
		t.Fatal("Error not created if database doesnt have a node")
	}
}

/*
 func TestDiscard(t *testing.T) {
	 dir, err := ioutil.TempDir("", "badger-test")
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer func() {
		 if err := os.RemoveAll(dir); err != nil {
			 t.Fatal(err)
		 }
	 }()
	 opts := badger.DefaultOptions(dir)
	 db, err := badger.Open(opts)
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer db.Close()
	 var rootTest []byte
	 smt := NewSMT()
	 keys := GetFreshData(20, 32)
	 fn1 := func(txn *badger.Txn) error {
		 // Add data to empty trie
		 values := GetFreshData(20, 32)
		 root, _ := smt.Update(keys, values)
		 rootTest = root
		 smt.Commit(txn)
		 return nil
	 }
	 fn2 := func(txn *badger.Txn) error {
		 keys = GetFreshData(20, 32)
		 values := GetFreshData(20, 32)
		 smt.Update(keys, values)
		 smt.Discard(txn)
		 return nil
	 }
	 fn3 := func(txn *badger.Txn) error {
		 keys = GetFreshData(20, 32)
		 values := GetFreshData(20, 32)
		 root, err := smt.Update(keys, values)
		 if err != nil {
			 t.Error(err)
		 }
		 rootTest = root
		 smt.Commit(txn)
		 return nil
	 }
	 fn4 := func(txn *badger.Txn) error {
		 keys := GetFreshData(20, 32)
		 values := GetFreshData(20, 32)
		 root, err := smt.Update(keys, values)
		 if err != nil {
			 t.Error(err)
		 }
		 rootTest = root
		 smt.Discard(txn)
		 return nil
	 }
	 err = db.Update(fn1)
	 if err != nil {
		 t.Error(err)
	 }
	 err = db.Update(fn2)
	 if err != nil {
		 t.Error(err)
	 }
	 if !bytes.Equal(smt.Root, rootTest) {
		 t.Fatal("Trie not rolled back")
	 }
	 if len(smt.db.updatedNodes) != 0 {
		 t.Fatal("Trie not rolled back")
	 }
	 err = db.Update(fn3)
	 if err != nil {
		 t.Error(err)
	 }
	 if !bytes.Equal(smt.Root, rootTest) {
		 t.Fatal("Trie not rolled back")
	 }
	 if len(smt.db.updatedNodes) != 0 {
		 t.Fatal("Trie not rolled back")
	 }
	 err = db.Update(fn4)
	 if err != nil {
		 t.Error(err)
	 }
 }
*/

func TestBigDelete(t *testing.T) {
	smt := NewSMT()

	for i := 0; i < 50; i++ {
		keys := GetFreshData(12, 32)
		values := GetFreshData(12, 32)
		_, err := smt.Update(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		deletes := make([][]byte, 2)
		for i := 0; i < 2; i++ {
			deletes[i] = DefaultLeaf
		}
		_, err = smt.Update(keys[:2], deletes)
		if err != nil {
			t.Fatal(err)
		}

		for j := 0; j < 2; j++ {
			v, err := smt.Get(keys[j])
			if err != nil {
				t.Fatal(err)
			}
			if len(v) > 0 {
				t.Fatal("deleted key still present")
			}
		}
		for j := 2; j < 12; j++ {
			v, err := smt.Get(keys[j])
			if err != nil {
				t.Fatal(err)
			}
			if len(v) < 32 {
				t.Fatal("non-deleted key not present")
			}
		}
	}

}

func benchmark10MAccounts10Ktps(smt *SMT, b *testing.B) {
	fmt.Println("\nLoading b.N x 1000 accounts")
	newkeys := GetFreshData(1000, 32)
	newvalues := GetFreshData(1000, 32)
	for index := 0; index < b.N; index++ {
		newvalues = GetFreshData(1000, 32)
		newkeys = GetFreshData(1000, 32)
		var end time.Time

		smt.Update(newkeys, newvalues)
		end = time.Now()

		start := time.Now()

		end2 := time.Now()

		for i, key := range newkeys {
			val, _ := smt.Get(key)
			if !bytes.Equal(val, newvalues[i]) {
				b.Fatal("new key not included")
			}
		}

		end3 := time.Now()
		elapsed := end.Sub(start)
		elapsed2 := end2.Sub(end)
		elapsed3 := end3.Sub(end2)
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Println(index, " : update time : ", elapsed, "commit time : ", elapsed2,
			"\n1000 Get time : ", elapsed3,
			"\nRAM : ", m.Sys/1024/1024, " MiB")
	}
}

/*
 //go test -run=BenchmarkSMT -bench=. -benchmem -test.benchtime=20s
 func BenchmarkSMT(b *testing.B) {
	 dir, err := ioutil.TempDir("", "badger-test")
	 if err != nil {
		 b.Fatal(err)
	 }
	 defer func() {
		 if err := os.RemoveAll(dir); err != nil {
			 b.Fatal(err)
		 }
	 }()
	 opts := badger.DefaultOptions(dir)
	 db, err := badger.Open(opts)
	 if err != nil {
		 b.Fatal(err)
	 }
	 defer db.Close()
	 smt := NewSMT()
	 benchmark10MAccounts10Ktps(db, smt, b)
 }

*/

/*

 // if bit31 set => batch1 and batch2 are leaf nodes
 func TestPrintDB(t *testing.T) {
	 dir, err := ioutil.TempDir("", "badger-test")
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer func() {
		 if err := os.RemoveAll(dir); err != nil {
			 t.Fatal(err)
		 }
	 }()
	 opts := badger.DefaultOptions(dir)
	 db, err := badger.Open(opts)
	 if err != nil {
		 t.Fatal(err)
	 }
	 defer db.Close()
	 smt := NewSMT()
	 keys := GetFreshData(10, 32)
	 //keys[0] = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	 //keys[1] = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	 values := GetFreshData(10, 32)
	 //values[0] = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}
	 //values[1] = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}
	 for i := 0; i < len(keys); i++ {
		 t.Logf("\nk::%x  \nv::%x\n", keys[i], values[i])
	 }
	 var root []byte
	 fn := func(txn *badger.Txn) error {
		 // Add data to empty trie
		 root, _ = smt.Update( keys, values)
		 smt.Commit(txn)
		 return nil
	 }
	 err = db.Update(fn)
	 if err != nil {
		 t.Error(err)
	 }
	 err = db.View(func(txn *badger.Txn) error {
		 opts := badger.DefaultIteratorOptions
		 opts.PrefetchSize = 10
		 prefix := []byte{}
		 prefix = append(prefix, smt.db.prefixFunc()...)
		 prefix = append(prefix, prefixNode()...)
		 opts.Prefix = prefix
		 it := txn.NewIterator(opts)
		 defer it.Close()
		 j := 0
		 t.Logf("\nROOT: %x\n", root)
		 for it.Rewind(); it.Valid(); it.Next() {
			 item := it.Item()
			 k := item.Key()
			 j++
			 err := item.Value(func(v []byte) error {
				 t.Logf("\nknum %d ::: key=%x, bitflag=%08b\n", j, k[len(smt.db.prefixFunc())+1:], v[:4])
				 if len(v) > 0 {
					 batch := smt.parseBatch(v)
					 wasKey := false
					 shortcut := false
					 for i, b := range batch {
						 if i == 0 {
							 if bytes.Equal(b, []byte{1}) {
								 shortcut = true
							 }
						 } else {
							 isLeaf := false
							 isKey := false
							 for k := 0; k < len(values); k++ {
								 if len(b) > 30 {
									 if bytes.Equal(values[k], b[:32]) {
										 if !wasKey {
											 t.Fatal("was key fail")
										 }
										 isLeaf = true
										 wasKey = false
									 }
									 if bytes.Equal(keys[k], b[:32]) {
										 isKey = true
										 wasKey = true
									 }
								 }
							 }
							 t.Logf("SC: %t  knum %d iK:%t iL:%t value %d =%x\n", shortcut, j, isKey, isLeaf, i, b)
						 }
					 }
				 }
				 return nil
			 })
			 if err != nil {
				 return err
			 }
		 }
		 return nil
	 })
	 if err != nil {
		 t.Error(err)
	 }
 }
*/
