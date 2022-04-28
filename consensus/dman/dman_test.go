package dman

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	aobjs "github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/application/objs/uint256"
	trie "github.com/MadBase/MadNet/badgerTrie"
	"github.com/MadBase/MadNet/consensus/appmock"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/interfaces"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type dmanTestProxy struct {
	sync.Mutex
	callIndex     int
	expectedCalls []testingProxyCall
	returns       [][]interface{}
	skipCallCheck bool
	db            *db.Database
	logger        *logrus.Logger
	testTx        *aobjs.Tx
}

// assert struct `dmanTestProxy` implements `reqBusView` , `appmock.Application`, `databaseView` interfaces
var _ reqBusView = &dmanTestProxy{}
var _ appmock.Application = &dmanTestProxy{}
var _ databaseView = &dmanTestProxy{}

//var _ typeProxyIface = &dmanTestProxy{}

// implementation of reqBusView interface

func (p *dmanTestProxy) RequestP2PGetPendingTx(ctx context.Context, txHashes [][]byte, opts ...grpc.CallOption) ([][]byte, error) {
	defer func() {
		p.callIndex++
	}()
	fmt.Println("RequestP2PGetPendingTx()")
	// cType := pendingTxCall
	p.Lock()
	defer p.Unlock()
	// if p.callIndex == len(p.expectedCalls) {
	// 	panic(fmt.Sprintf("got unexpected call of type %s : expected calls %v", cType, p.expectedCalls))
	// }
	// if p.expectedCalls[p.callIndex] != cType {
	// 	panic(fmt.Sprintf("got unexpected call of type %s at index %v : expected calls %v", cType, p.callIndex, p.expectedCalls))
	// }
	// if ctx == nil {
	// 	panic(fmt.Sprintf("ctx was nil in test mock object of call type %s", cType))
	// }

	//ret := [][]byte{make([]byte, constants.HashLen)}

	bin, err := p.testTx.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("RequestP2PGetPendingTx() error tx.MarshalBinary(): %v", err))
	}

	hash, err := p.testTx.TxHash()
	if err != nil {
		panic(fmt.Errorf("RequestP2PGetPendingTx() error tx.TxHash(): %v", err))
	}
	fmt.Printf("RequestP2PGetPendingTx returning %v", hash)

	ret := [][]byte{bin}

	// returnTuple := p.returns[p.callIndex]
	// tx := returnTuple[0].([][]byte)
	// err, ok := returnTuple[1].(error)
	// if ok {
	// 	return tx, err
	// }
	return ret, nil
}

func (p *dmanTestProxy) RequestP2PGetMinedTxs(ctx context.Context, txHashes [][]byte, opts ...grpc.CallOption) ([][]byte, error) {
	defer func() {
		p.callIndex++
	}()
	fmt.Println("RequestP2PGetMinedTxs()")
	// cType := minedTxCall
	p.Lock()
	defer p.Unlock()
	// if p.callIndex == len(p.expectedCalls) {
	// 	panic(fmt.Sprintf("got unexpected call of type %s : expected calls %v", cType, p.expectedCalls))
	// }
	// if p.expectedCalls[p.callIndex] != cType {
	// 	panic(fmt.Sprintf("got unexpected call of type %s at index %v : expected calls %v", cType, p.callIndex, p.expectedCalls))
	// }
	// if ctx == nil {
	// 	panic(fmt.Sprintf("ctx was nil in test mock object of call type %s", cType))
	// }

	//ret := [][]byte{make([]byte, constants.HashLen)}

	bin, err := p.testTx.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("RequestP2PGetMinedTxs() error tx.MarshalBinary(): %v", err))
	}

	hash, err := p.testTx.TxHash()
	if err != nil {
		panic(fmt.Errorf("RequestP2PGetMinedTxs() error tx.TxHash(): %v", err))
	}
	fmt.Printf("RequestP2PGetMinedTxs returning %v", hash)

	ret := [][]byte{bin}

	// returnTuple := p.returns[p.callIndex]
	// tx := returnTuple[0].([][]byte)
	// err, ok := returnTuple[1].(error)
	// if ok {
	// 	return tx, err
	// }
	return ret, nil
}

func (p *dmanTestProxy) RequestP2PGetBlockHeaders(ctx context.Context, blockNums []uint32, opts ...grpc.CallOption) ([]*objs.BlockHeader, error) {
	defer func() {
		p.callIndex++
	}()
	// cType := blockHeaderCall
	fmt.Println("RequestP2PGetBlockHeaders()")
	p.Lock()
	defer p.Unlock()
	// if p.callIndex == len(p.expectedCalls) {
	// 	panic(fmt.Sprintf("got unexpected call of type %s : expected calls %v, callIndex %v", cType, p.expectedCalls, p.callIndex))
	// }
	// if p.expectedCalls[p.callIndex] != cType {
	// 	panic(fmt.Sprintf("got unexpcted call of type %s at index %v : expected calls %v, callIndex %v", cType, p.callIndex, p.expectedCalls, p.callIndex))
	// }
	// if ctx == nil {
	// 	panic(fmt.Sprintf("ctx was nil in test mock object of call type %s", cType))
	// }

	//bh := makeGoodBlock(t)

	// returnTuple := p.returns[p.callIndex]
	// bh := returnTuple[0].([]*objs.BlockHeader)
	// err, ok := returnTuple[1].(error)
	// if ok {
	// 	return bh, err
	// }
	return nil, errors.New("could not request block header from P2P")
}

// implementation of databaseView interface

func (p *dmanTestProxy) SetTxCacheItem(txn *badger.Txn, height uint32, txHash []byte, tx []byte) error {
	fmt.Printf("SetTxCacheItem mocked. height: %v, txHash: %x\n", height, txHash)
	return p.db.SetTxCacheItem(txn, height, txHash, tx)
}

func (p *dmanTestProxy) GetTxCacheItem(txn *badger.Txn, height uint32, txHash []byte) ([]byte, error) {
	fmt.Printf("GetTxCacheItem mocked. height: %v, txHash: %x\n", height, txHash)
	return p.db.GetTxCacheItem(txn, height, txHash)
}

func (p *dmanTestProxy) SetCommittedBlockHeader(txn *badger.Txn, v *objs.BlockHeader) error {
	return p.db.SetCommittedBlockHeader(txn, v)
}

func (p *dmanTestProxy) TxCacheDropBefore(txn *badger.Txn, beforeHeight uint32, maxKeys int) error {
	return p.db.TxCacheDropBefore(txn, beforeHeight, maxKeys)
}

// implementation of appmock.Application interface

// //MockApplication is the the receiver for TxHandler interface
// type MockApplication struct {
// 	logger     *logrus.Logger
// 	validValue *objs.Proposal
// 	MissingTxn bool
// }

const (
	notImpl = "not implemented"
)

// SetNextValidValue is defined on the interface object
func (p *dmanTestProxy) SetNextValidValue(vv *objs.Proposal) {
	panic(notImpl)
}

// ApplyState is defined on the interface object
func (p *dmanTestProxy) ApplyState(txn *badger.Txn, chainID, height uint32, txs []interfaces.Transaction) ([]byte, error) {
	fmt.Printf("dmanTestProxy.ApplyState()\n")
	//err := p.SetTxCacheItem() AddTxs(txn, 1, []interfaces.Transaction{tx})
	//assert.Nil(t, err)
	return nil, nil
}

//GetValidProposal is defined on the interface object
func (p *dmanTestProxy) GetValidProposal(txn *badger.Txn, chainID, height, maxBytes uint32) ([]interfaces.Transaction, []byte, error) {
	return nil, nil, nil
}

// PendingTxAdd is defined on the interface object
func (p *dmanTestProxy) PendingTxAdd(txn *badger.Txn, chainID, height uint32, txs []interfaces.Transaction) error {
	return nil
}

//IsValid is defined on the interface object
func (p *dmanTestProxy) IsValid(txn *badger.Txn, chainID uint32, height uint32, stateHash []byte, _ []interfaces.Transaction) (bool, error) {
	return false, nil
}

// MinedTxGet is defined on the interface object
func (p *dmanTestProxy) MinedTxGet(*badger.Txn, [][]byte) ([]interfaces.Transaction, [][]byte, error) {
	return nil, nil, nil
}

// PendingTxGet is defined on the interface object
func (p *dmanTestProxy) PendingTxGet(txn *badger.Txn, height uint32, txhashes [][]byte) ([]interfaces.Transaction, [][]byte, error) {
	return nil, nil, nil
}

//PendingTxContains is defined on the interface object
func (p *dmanTestProxy) PendingTxContains(txn *badger.Txn, height uint32, txHashes [][]byte) ([][]byte, error) {
	return nil, nil
}

// UnmarshalTx is defined on the interface object
func (p *dmanTestProxy) UnmarshalTx(v []byte) (interfaces.Transaction, error) {
	tx := &aobjs.Tx{}
	err := tx.UnmarshalBinary(v)
	if err != nil {
		utils.DebugTrace(p.logger, err)
		return nil, err
	}
	return tx, nil
}

// StoreSnapShotNode is defined on the interface object
func (p *dmanTestProxy) StoreSnapShotNode(txn *badger.Txn, batch []byte, root []byte, layer int) ([][]byte, int, []trie.LeafNode, error) {
	panic(notImpl)
}

// GetSnapShotNode is defined on the interface object
func (p *dmanTestProxy) GetSnapShotNode(txn *badger.Txn, height uint32, key []byte) ([]byte, error) {
	panic(notImpl)
}

// StoreSnapShotStateData is defined on the interface object
func (p *dmanTestProxy) StoreSnapShotStateData(txn *badger.Txn, key []byte, value []byte, data []byte) error {
	panic(notImpl)
}

// GetSnapShotStateData is defined on the interface object
func (p *dmanTestProxy) GetSnapShotStateData(txn *badger.Txn, key []byte) ([]byte, error) {
	panic(notImpl)
}

// FinalizeSnapShotRoot is defined on the interface object
func (p *dmanTestProxy) FinalizeSnapShotRoot(txn *badger.Txn, root []byte, height uint32) error {
	panic(notImpl)
}

// BeginSnapShotSync is defined on the interface object
func (p *dmanTestProxy) BeginSnapShotSync(txn *badger.Txn) error {
	panic(notImpl)
}

// FinalizeSync is defined on the interface object
func (p *dmanTestProxy) FinalizeSync(txn *badger.Txn) error {
	panic(notImpl)
}

// MockTransaction is defined on the interface object
type MockTransaction struct {
	V []byte
}

// TxHash is defined on the interface object
func (m *MockTransaction) TxHash() ([]byte, error) {
	return crypto.Hasher(m.V), nil
}

//MarshalBinary is defined on the interface object
func (m *MockTransaction) MarshalBinary() ([]byte, error) {
	return m.V, nil
}

//XXXIsTx is defined on the interface object
func (m *MockTransaction) XXXIsTx() {}

// test setup

func initDatabase(ctx context.Context, path string, inMemory bool) *badger.DB {
	db, err := utils.OpenBadger(ctx.Done(), path, inMemory)
	if err != nil {
		panic(err)
	}
	return db
}

func setupDmanTests(t *testing.T) (testProxy *dmanTestProxy, dman *DMan, ownerSigner aobjs.Signer, closeFn func()) {
	logger := logging.GetLogger("Test")
	deferables := make([]func(), 0)

	closeFn = func() {
		// iterate in reverse order because deferables behave like a stack:
		// the last added deferable should be the first executed
		totalDeferables := len(deferables)
		for i := totalDeferables - 1; i >= 0; i-- {
			deferables[i]()
		}
	}

	ctx := context.Background()
	nodeCtx, cf := context.WithCancel(ctx)
	deferables = append(deferables, cf)

	// Initialize consensus db: stores all state the consensus mechanism requires to work
	rawConsensusDb := initDatabase(nodeCtx, "", true)
	var closeDB func() = func() {
		err := rawConsensusDb.Close()
		if err != nil {
			panic(fmt.Errorf("error closing rawConsensusDb: %v", err))
		}
	}
	deferables = append(deferables, closeDB)

	db := &db.Database{}
	db.Init(rawConsensusDb)

	testProxy = &dmanTestProxy{
		db:     db,
		logger: logger,
	}

	ra := &RootActor{}
	err := ra.Init(logger, testProxy)
	if err != nil {
		closeFn()
		t.Fatal(err)
		return
	}
	ra.Start()
	deferables = append(deferables, ra.Close)

	dman = &DMan{
		ra,
		testProxy,
		testProxy,
		nil,
		logger,
	}
	dman.Init(testProxy, testProxy, testProxy)
	deferables = append(deferables, dman.Close)

	dman.Start()

	ownerSigner = testingOwner()

	return
}

func Test_DManProxy(t *testing.T) {
	var p *dmanTestProxy = &dmanTestProxy{}
	var dman *DMan = &DMan{}
	dman.Init(p, p, p)
	dman.Close()
}

func Test_GetNonExistent(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)

	hash, err := tx.TxHash()
	assert.Nil(t, err)

	err = testProxy.db.View(func(txn *badger.Txn) error {
		tx2Binary, err := dman.database.GetTxCacheItem(txn, 1, hash)
		assert.NotNil(t, err)
		assert.Nil(t, tx2Binary)

		return nil
	})

	assert.Nil(t, err)

}

func Test_SetAndGet(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)

	hash, err := tx.TxHash()
	assert.Nil(t, err)
	binary, err := tx.MarshalBinary()
	assert.Nil(t, err)

	txsToGet := make([][]byte, 0)

	// test
	err = testProxy.db.Update(func(txn *badger.Txn) error {
		err = dman.database.SetTxCacheItem(txn, 1, hash, binary)
		assert.Nil(t, err)
		txsToGet = append(txsToGet, hash)

		tx2Binary, err := dman.database.GetTxCacheItem(txn, 1, hash)
		assert.Nil(t, err)
		assert.NotNil(t, tx2Binary)

		tx2 := &aobjs.Tx{}
		err = tx2.UnmarshalBinary(tx2Binary)
		assert.Nil(t, err)

		binary2, err := tx2.MarshalBinary()
		assert.Nil(t, err)
		assert.Equal(t, binary, binary2)
		return err
	})

	assert.Nil(t, err)

	err = testProxy.db.View(func(txn *badger.Txn) error {
		// read through dman.GetTxs
		txs, missing, err := dman.GetTxs(txn, 1, 1, txsToGet)
		assert.Nil(t, err)
		assert.Len(t, txs, 1)
		assert.Len(t, missing, 0)

		binary2, err := txs[0].MarshalBinary()
		assert.Nil(t, err)
		assert.Equal(t, binary, binary2)

		// read through dman.database.GetTxCacheItem
		tx2Binary, err := dman.database.GetTxCacheItem(txn, 1, hash)
		assert.Nil(t, err)
		assert.NotNil(t, tx2Binary)

		tx2 := &aobjs.Tx{}
		err = tx2.UnmarshalBinary(tx2Binary)
		assert.Nil(t, err)

		binary2, err = tx2.MarshalBinary()
		assert.Nil(t, err)
		assert.Equal(t, binary, binary2)
		assert.Equal(t, binary, tx2Binary)
		assert.Equal(t, binary2, tx2Binary)

		return err
	})

	assert.Nil(t, err)

}

func Test_FlushCacheToDisk(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)

	hash, err := tx.TxHash()
	assert.Nil(t, err)
	binary, err := tx.MarshalBinary()
	assert.Nil(t, err)

	// test
	err = testProxy.db.Update(func(txn *badger.Txn) error {
		err = dman.database.SetTxCacheItem(txn, 1, hash, binary)
		assert.Nil(t, err)

		err := dman.FlushCacheToDisk(txn, 1)
		assert.Nil(t, err)

		assert.False(t, dman.downloadActor.bhc.Contains(1))

		return err
	})

	assert.Nil(t, err)
}

func Test_CleanCache(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)

	hash, err := tx.TxHash()
	assert.Nil(t, err)
	binary, err := tx.MarshalBinary()
	assert.Nil(t, err)

	// test
	err = testProxy.db.Update(func(txn *badger.Txn) error {
		err = dman.database.SetTxCacheItem(txn, 1, hash, binary)
		assert.Nil(t, err)

		err := dman.CleanCache(txn, 1)
		assert.Nil(t, err)

		assert.False(t, dman.downloadActor.bhc.Contains(1))

		return err
	})

	assert.Nil(t, err)
}

func Test_AddTxs(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)
	var txs []interfaces.Transaction = []interfaces.Transaction{tx}

	hash, err := tx.TxHash()
	assert.Nil(t, err)
	binary, err := tx.MarshalBinary()
	assert.Nil(t, err)

	txsToGet := make([][]byte, 0)
	txsToGet = append(txsToGet, hash)

	// add Txs
	err = testProxy.db.Update(func(txn *badger.Txn) error {
		err := dman.AddTxs(txn, 1, txs)
		assert.Nil(t, err)

		return err
	})

	assert.Nil(t, err)

	err = testProxy.db.View(func(txn *badger.Txn) error {
		// read through dman.GetTxs
		txs, missing, err := dman.GetTxs(txn, 1, 1, txsToGet)
		assert.Nil(t, err)
		assert.Len(t, txs, 1)
		assert.Len(t, missing, 0)

		binary2, err := txs[0].MarshalBinary()
		assert.Nil(t, err)
		assert.Equal(t, binary, binary2)

		// read through dman.database.GetTxCacheItem
		tx2Binary, err := dman.database.GetTxCacheItem(txn, 1, hash)
		assert.Nil(t, err)
		assert.NotNil(t, tx2Binary)

		tx2 := &aobjs.Tx{}
		err = tx2.UnmarshalBinary(tx2Binary)
		assert.Nil(t, err)

		binary2, err = tx2.MarshalBinary()
		assert.Nil(t, err)
		assert.Equal(t, binary, binary2)
		assert.Equal(t, binary, tx2Binary)
		assert.Equal(t, binary2, tx2Binary)

		return err
	})

	assert.Nil(t, err)
}

func Test_DownloadTxs(t *testing.T) {
	testProxy, dman, ownerSigner, closeFn := setupDmanTests(t)
	defer closeFn()

	assert.False(t, dman.downloadActor.bhc.Contains(1))

	/*consumedUTXOs*/
	_, tx := makeTxInitial(ownerSigner)
	var txs []interfaces.Transaction = []interfaces.Transaction{tx}

	hash, err := tx.TxHash()
	assert.Nil(t, err)
	testProxy.testTx = tx

	txsToGet := make([][]byte, 0)
	txsToGet = append(txsToGet, hash)

	assert.False(t, dman.downloadActor.txc.Contains(hash))

	// add Txs
	err = testProxy.db.Update(func(txn *badger.Txn) error {
		err := dman.AddTxs(txn, 1, txs)
		assert.Nil(t, err)

		return err
	})

	assert.Nil(t, err)

	assert.False(t, dman.downloadActor.txc.Contains(hash))
	assert.False(t, dman.downloadActor.bhc.Contains(1))

	// download the Txs
	dman.DownloadTxs(1, 1, txsToGet)

	// wait some time for actors to download Txs
	<-time.After(5 * time.Second)

	t.Logf("expecting hash: %v", hash)

	assert.True(t, dman.downloadActor.txc.Contains(hash))
	assert.False(t, dman.downloadActor.bhc.Contains(1))
}

func generateFullChain(length int) ([]*objs.BClaims, [][][]byte, error) {
	chain := []*objs.BClaims{}
	txHashes := [][][]byte{}
	txhash := crypto.Hasher([]byte(strconv.Itoa(1)))
	txHshLst := [][]byte{txhash}
	txRoot, err := objs.MakeTxRoot(txHshLst)
	if err != nil {
		return nil, nil, err
	}
	txHashes = append(txHashes, txHshLst)
	bclaims := &objs.BClaims{
		ChainID:    1,
		Height:     1,
		TxCount:    1,
		PrevBlock:  crypto.Hasher([]byte("foo")),
		TxRoot:     txRoot,
		StateRoot:  crypto.Hasher([]byte("")),
		HeaderRoot: crypto.Hasher([]byte("")),
	}
	chain = append(chain, bclaims)
	for i := 1; i < length; i++ {
		bhsh, err := chain[i-1].BlockHash()
		if err != nil {
			return nil, nil, err
		}
		txhash := crypto.Hasher([]byte(strconv.Itoa(i)))
		txHshLst := [][]byte{txhash}
		txRoot, err := objs.MakeTxRoot(txHshLst)
		if err != nil {
			return nil, nil, err
		}
		txHashes = append(txHashes, txHshLst)
		bclaims := &objs.BClaims{
			ChainID:    1,
			Height:     uint32(len(chain) + 1),
			TxCount:    1,
			PrevBlock:  bhsh,
			TxRoot:     txRoot,
			StateRoot:  chain[i-1].StateRoot,
			HeaderRoot: chain[i-1].HeaderRoot,
		}
		chain = append(chain, bclaims)
	}
	return chain, txHashes, nil
}

func testingOwner() aobjs.Signer {
	signer := &crypto.Secp256k1Signer{}
	err := signer.SetPrivk(crypto.Hasher([]byte("secret")))
	if err != nil {
		panic(err)
	}
	return signer
}

func accountFromSigner(s aobjs.Signer) []byte {
	pubk, err := s.Pubkey()
	if err != nil {
		panic(err)
	}
	return crypto.GetAccount(pubk)
}

func makeVS(ownerSigner aobjs.Signer) *aobjs.TXOut {
	cid := uint32(2)
	//val := uint32(1)
	val := uint256.One()

	ownerAcct := accountFromSigner(ownerSigner)
	owner := &aobjs.ValueStoreOwner{}
	owner.New(ownerAcct, constants.CurveSecp256k1)

	vsp := &aobjs.VSPreImage{
		ChainID: cid,
		Value:   val,
		Owner:   owner,
		Fee:     uint256.Zero(),
	}
	vs := &aobjs.ValueStore{
		VSPreImage: vsp,
		TxHash:     make([]byte, constants.HashLen),
	}
	utxInputs := &aobjs.TXOut{}
	err := utxInputs.NewValueStore(vs)
	if err != nil {
		panic(err)
	}
	return utxInputs
}

func makeVSTXIn(ownerSigner aobjs.Signer, txHash []byte) (*aobjs.TXOut, *aobjs.TXIn) {
	vs := makeVS(ownerSigner)
	vss, err := vs.ValueStore()
	if err != nil {
		panic(err)
	}
	if txHash == nil {
		txHash = make([]byte, constants.HashLen)
		rand.Read(txHash)
	}
	vss.TxHash = txHash

	txIn, err := vss.MakeTxIn()
	if err != nil {
		panic(err)
	}
	return vs, txIn
}

func makeTxInitial(ownerSigner aobjs.Signer) (aobjs.Vout, *aobjs.Tx) {
	consumedUTXOs := aobjs.Vout{}
	txInputs := []*aobjs.TXIn{}
	for i := 0; i < 2; i++ {
		utxo, txin := makeVSTXIn(ownerSigner, nil)
		consumedUTXOs = append(consumedUTXOs, utxo)
		txInputs = append(txInputs, txin)
	}
	generatedUTXOs := aobjs.Vout{}
	for i := 0; i < 2; i++ {
		generatedUTXOs = append(generatedUTXOs, makeVS(ownerSigner))
	}
	err := generatedUTXOs.SetTxOutIdx()
	if err != nil {
		panic(err)
	}
	txfee := uint256.Zero()
	tx := &aobjs.Tx{
		Vin:  txInputs,
		Vout: generatedUTXOs,
		Fee:  txfee,
	}
	err = tx.SetTxHash()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 2; i++ {
		vs, err := consumedUTXOs[i].ValueStore()
		if err != nil {
			panic(err)
		}
		err = vs.Sign(tx.Vin[i], ownerSigner)
		if err != nil {
			panic(err)
		}
	}
	return consumedUTXOs, tx
}
