package admin

import (
	"context"
	"fmt"
	"testing"

	aobjs "github.com/MadBase/MadNet/application/objs"
	trie "github.com/MadBase/MadNet/badgerTrie"
	"github.com/MadBase/MadNet/config"
	"github.com/MadBase/MadNet/consensus/appmock"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	mncrypto "github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/dynamics"
	"github.com/MadBase/MadNet/interfaces"
	"github.com/MadBase/MadNet/ipc"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ahTestProxy struct {
	logger    *logrus.Logger
	db        *db.Database
	ah        *Handlers
	secretKey []byte
}

var _ appmock.Application = &ahTestProxy{}

const (
	notImpl = "not implemented"
)

// SetNextValidValue is defined on the interface object
func (p *ahTestProxy) SetNextValidValue(vv *objs.Proposal) {
	panic(notImpl)
}

// ApplyState is defined on the interface object
func (p *ahTestProxy) ApplyState(txn *badger.Txn, chainID, height uint32, txs []interfaces.Transaction) ([]byte, error) {
	fmt.Printf("ahTestProxy.ApplyState()\n")
	//err := p.SetTxCacheItem() AddTxs(txn, 1, []interfaces.Transaction{tx})
	//assert.Nil(t, err)
	return nil, nil
}

//GetValidProposal is defined on the interface object
func (p *ahTestProxy) GetValidProposal(txn *badger.Txn, chainID, height, maxBytes uint32) ([]interfaces.Transaction, []byte, error) {
	return nil, nil, nil
}

// PendingTxAdd is defined on the interface object
func (p *ahTestProxy) PendingTxAdd(txn *badger.Txn, chainID, height uint32, txs []interfaces.Transaction) error {
	return nil
}

//IsValid is defined on the interface object
func (p *ahTestProxy) IsValid(txn *badger.Txn, chainID uint32, height uint32, stateHash []byte, _ []interfaces.Transaction) (bool, error) {
	return false, nil
}

// MinedTxGet is defined on the interface object
func (p *ahTestProxy) MinedTxGet(*badger.Txn, [][]byte) ([]interfaces.Transaction, [][]byte, error) {
	return nil, nil, nil
}

// PendingTxGet is defined on the interface object
func (p *ahTestProxy) PendingTxGet(txn *badger.Txn, height uint32, txhashes [][]byte) ([]interfaces.Transaction, [][]byte, error) {
	return nil, nil, nil
}

//PendingTxContains is defined on the interface object
func (p *ahTestProxy) PendingTxContains(txn *badger.Txn, height uint32, txHashes [][]byte) ([][]byte, error) {
	return nil, nil
}

// UnmarshalTx is defined on the interface object
func (p *ahTestProxy) UnmarshalTx(v []byte) (interfaces.Transaction, error) {
	tx := &aobjs.Tx{}
	err := tx.UnmarshalBinary(v)
	if err != nil {
		utils.DebugTrace(p.logger, err)
		return nil, err
	}
	return tx, nil
}

// StoreSnapShotNode is defined on the interface object
func (p *ahTestProxy) StoreSnapShotNode(txn *badger.Txn, batch []byte, root []byte, layer int) ([][]byte, int, []trie.LeafNode, error) {
	panic(notImpl)
}

// GetSnapShotNode is defined on the interface object
func (p *ahTestProxy) GetSnapShotNode(txn *badger.Txn, height uint32, key []byte) ([]byte, error) {
	panic(notImpl)
}

// StoreSnapShotStateData is defined on the interface object
func (p *ahTestProxy) StoreSnapShotStateData(txn *badger.Txn, key []byte, value []byte, data []byte) error {
	panic(notImpl)
}

// GetSnapShotStateData is defined on the interface object
func (p *ahTestProxy) GetSnapShotStateData(txn *badger.Txn, key []byte) ([]byte, error) {
	panic(notImpl)
}

// FinalizeSnapShotRoot is defined on the interface object
func (p *ahTestProxy) FinalizeSnapShotRoot(txn *badger.Txn, root []byte, height uint32) error {
	panic(notImpl)
}

// BeginSnapShotSync is defined on the interface object
func (p *ahTestProxy) BeginSnapShotSync(txn *badger.Txn) error {
	panic(notImpl)
}

// FinalizeSync is defined on the interface object
func (p *ahTestProxy) FinalizeSync(txn *badger.Txn) error {
	panic(notImpl)
}

// MockTransaction is defined on the interface object
type MockTransaction struct {
	V []byte
}

// TxHash is defined on the interface object
func (m *MockTransaction) TxHash() ([]byte, error) {
	return mncrypto.Hasher(m.V), nil
}

//MarshalBinary is defined on the interface object
func (m *MockTransaction) MarshalBinary() ([]byte, error) {
	return m.V, nil
}

//XXXIsTx is defined on the interface object
func (m *MockTransaction) XXXIsTx() {}

// setupAHTests
func setupAHTests(t *testing.T) (testProxy *ahTestProxy, closeFn func()) {
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

	var chainID uint32 = 1337
	ctx := context.Background()
	nodeCtx, cf := context.WithCancel(ctx)
	deferables = append(deferables, cf)

	// Initialize consensus db: stores all state the consensus mechanism requires to work
	rawConsensusDb, err := utils.OpenBadger(nodeCtx.Done(), "", true)
	assert.Nil(t, err)
	var closeDB func() = func() {
		err := rawConsensusDb.Close()
		if err != nil {
			panic(fmt.Errorf("error closing rawConsensusDb: %v", err))
		}
	}
	deferables = append(deferables, closeDB)

	db := &db.Database{}
	db.Init(rawConsensusDb)

	secretKey := mncrypto.Hasher([]byte("someSuperFancySecretThatWillBeHashed"))

	testProxy = &ahTestProxy{
		logger:    logger,
		db:        db,
		secretKey: secretKey,
	}

	ethPubKey := []byte("b904C0A2d203Ceb2B518055f116064666C028240")
	storage := &dynamics.Storage{}

	ipcServer := ipc.NewServer(config.Configuration.Firewalld.SocketFile)
	deferables = append(deferables, ipcServer.Close)

	testProxy.ah = &Handlers{}
	testProxy.ah.Init(chainID, db, secretKey, testProxy, ethPubKey, storage, ipcServer)
	deferables = append(deferables, testProxy.ah.Close)

	return
}

func TestAdminHandlerSetup(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	err := ahTestProxy.db.View(func(txn *badger.Txn) error {
		txHashes := make([][]byte, 0)
		txHashes = append(txHashes, []byte("aaa"))
		pendingTxs, err := ahTestProxy.PendingTxContains(txn, 1, txHashes)
		assert.Nil(t, err)
		assert.Empty(t, pendingTxs)

		return err
	})
	assert.Nil(t, err)
}

func TestAdminHandler_GetKey(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	key, err := ahTestProxy.ah.GetKey([]byte("aaa"))
	assert.Nil(t, err)
	assert.NotNil(t, key)
	// todo: this was unexpected just by looking at the API docs
	// todo: should this be renamed from GetKey to GetSecretKeyHash ?
	assert.Equal(t, key, ahTestProxy.secretKey)

}

func TestAdminHandler_DontGetPrivK(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	key, err := ahTestProxy.ah.GetPrivK([]byte("a random name 123"))
	assert.NotNil(t, err)
	assert.Nil(t, key)
}

func TestAdminHandler_AddPrivateKey_CurveSecp256k1(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	err := ahTestProxy.ah.AddPrivateKey([]byte("a random private key 123"), constants.CurveSecp256k1)
	assert.Nil(t, err)
}

func TestAdminHandler_AddPrivateKey_CurveBN256Eth(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	err := ahTestProxy.ah.AddPrivateKey([]byte("a random private key 123"), constants.CurveBN256Eth)
	assert.Nil(t, err)
}
