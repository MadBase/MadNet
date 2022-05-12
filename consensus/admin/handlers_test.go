package admin

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"testing"
	"time"

	aobjs "github.com/MadBase/MadNet/application/objs"
	utxo "github.com/MadBase/MadNet/application/utxohandler"
	trie "github.com/MadBase/MadNet/badgerTrie"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/config"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/crypto/bn256"
	"github.com/MadBase/MadNet/dynamics"
	"github.com/MadBase/MadNet/interfaces"
	"github.com/MadBase/MadNet/ipc"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ahTestProxy struct {
	logger         *logrus.Logger
	db             *db.Database
	ah             *Handlers
	secretKey      []byte
	rawConsensusDb *badger.DB
	utxoHandler    *utxo.UTXOHandler
}

var _ interfaces.Application = &ahTestProxy{}

const (
	notImpl = "not implemented"
)

// SetNextValidValue is defined on the interface object
func (p *ahTestProxy) SetNextValidValue(vv *objs.Proposal) {
	panic(notImpl)
}

// ApplyState is defined on the interface object
func (p *ahTestProxy) ApplyState(txn *badger.Txn, chainID, height uint32, txs []interfaces.Transaction) ([]byte, error) {
	return p.utxoHandler.ApplyState(txn, aobjs.TxVec{}, 1)
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
	return crypto.Hasher(m.V), nil
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

	hndlr := utxo.NewUTXOHandler(rawConsensusDb)
	err = hndlr.Init(1)
	if err != nil {
		t.Fatal(err)
	}

	secretKey := crypto.Hasher([]byte("someSuperFancySecretThatWillBeHashed"))

	testProxy = &ahTestProxy{
		logger:         logger,
		db:             db,
		secretKey:      secretKey,
		rawConsensusDb: rawConsensusDb,
		utxoHandler:    hndlr,
	}

	ethPubKey := []byte("b904C0A2d203Ceb2B518055f116064666C028240")
	storage := &dynamics.Storage{}

	ipcServer := ipc.NewServer(config.Configuration.Firewalld.SocketFile)
	deferables = append(deferables, ipcServer.Close)

	testProxy.ah = &Handlers{}
	testProxy.ah.Init(chainID, db, secretKey, testProxy, ethPubKey, storage, ipcServer)
	deferables = append(deferables, testProxy.ah.Close)

	assert.False(t, testProxy.ah.IsInitialized())

	// start goroutine to emulate Synchronizer.adminInteruptLoop()
	closeCh := make(chan struct{})
	deferables = append(deferables, func() { close(closeCh) })
	synchronizer := &synchronizerMock{}

	go func(closeChan chan struct{}) {
		defer func() { fmt.Println("Stopping AdminInterupt loop") }()
		fmt.Println("Starting AdminInterupt loop")
		for {
			select {
			case testProxy.ah.ReceiveLock <- synchronizer:
				continue
			case <-closeChan:
				return
			}
		}
	}(closeCh)

	// start goroutine to emulate Handler.InitializationMonitor()
	closeChInitializationMonitor := make(chan struct{})
	deferables = append(deferables, func() { close(closeChInitializationMonitor) })
	go testProxy.ah.InitializationMonitor(closeChInitializationMonitor)

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

	err := ahTestProxy.ah.AddPrivateKey([]byte("a random 32byte private key 1234"), constants.CurveSecp256k1)
	assert.Nil(t, err)
}

func TestAdminHandler_AddPrivateKey_CurveBN256Eth(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	err := ahTestProxy.ah.AddPrivateKey([]byte("a random private key 123"), constants.CurveBN256Eth)
	assert.Nil(t, err)
}

func TestAdminHandler_Sinchronization(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	assert.False(t, ahTestProxy.ah.IsSynchronized())

	ahTestProxy.ah.SetSynchronized(true)
	assert.True(t, ahTestProxy.ah.IsSynchronized())

	ahTestProxy.ah.SetSynchronized(false)
	assert.False(t, ahTestProxy.ah.IsSynchronized())
}

func TestAdminHandler_AddValidatorSet(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()

	vs := generateValidatorSet(t)

	err := ahTestProxy.ah.AddValidatorSet(vs)
	assert.Nil(t, err)
}

func TestAdminHandler_RegisterSnapshotCallback(t *testing.T) {
	ahTestProxy, closeFn := setupAHTests(t)
	defer closeFn()
	var didRunCallback bool = false

	ahTestProxy.ah.RegisterSnapshotCallback(func(bh *objs.BlockHeader) error {
		didRunCallback = true
		ahTestProxy.logger.Printf("SnapshotCallback is being called!")
		return nil
	})

	// add validator set
	vs := generateValidatorSet(t)
	err := ahTestProxy.ah.AddValidatorSet(vs)
	assert.Nil(t, err, err)

	bhs := ahTestProxy.makeGoodBlock(t, 32)

	// add own state with custom validator address
	err = ahTestProxy.db.Update(func(txn *badger.Txn) error {
		addrBytes, err := hex.DecodeString("9AC1c9afBAec85278679fF75Ef109217f26b1417")
		assert.Nil(t, err)
		ownState := &objs.OwnState{
			VAddr:             addrBytes,
			SyncToBH:          bhs[0],
			MaxBHSeen:         bhs[0],
			CanonicalSnapShot: bhs[0],
			PendingSnapShot:   bhs[0],
		}
		return ahTestProxy.db.SetOwnState(txn, ownState)
	})
	assert.Nil(t, err)

	for i := 0; i < 32; i++ {
		binary, err := bhs[i].MarshalBinary()
		assert.Nil(t, err)

		bht := &objs.BlockHeader{}
		err = bht.UnmarshalBinary(binary)
		assert.Nil(t, err)

		if i == 0 {
			err = ahTestProxy.db.Update(func(txn *badger.Txn) error {
				err := ahTestProxy.db.SetSnapshotBlockHeader(txn, bhs[i])
				assert.Nil(t, err)

				return err
			})
			assert.Nil(t, err)
		} else {
			err = ahTestProxy.db.Update(func(txn *badger.Txn) error {
				err := ahTestProxy.db.SetBroadcastBlockHeader(txn, bhs[i])
				assert.Nil(t, err)

				_, err = ahTestProxy.ah.database.GetBroadcastBlockHeader(txn)
				assert.Nil(t, err)

				return err
			})
			assert.Nil(t, err)
		}

		assert.Nil(t, err)
	}

	<-time.After(3 * time.Second)
	assert.True(t, didRunCallback)
	assert.True(t, ahTestProxy.ah.IsInitialized())
}

func (ah *ahTestProxy) makeGoodBlock(t *testing.T, nBlocks int) []*objs.BlockHeader {
	bclaimsList, txHashListList, err := ah.generateChain(nBlocks)
	if err != nil {
		t.Fatal(err)
	}

	bhs := make([]*objs.BlockHeader, 0)

	gk := crypto.BNGroupSigner{}
	err = gk.SetPrivk(crypto.Hasher([]byte("secret")))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < nBlocks; i++ {
		bclaims := bclaimsList[i]
		bhsh, err := bclaims.BlockHash()
		if err != nil {
			t.Fatal(err)
		}

		sig, err := gk.Sign(bhsh)
		if err != nil {
			t.Fatal(err)
		}
		bh := &objs.BlockHeader{
			BClaims:  bclaims,
			SigGroup: sig,
			TxHshLst: txHashListList[i],
		}
		bhs = append(bhs, bh)
	}

	return bhs
}

func (ah *ahTestProxy) generateChain(nBlocks int) ([]*objs.BClaims, [][][]byte, error) {
	chain := []*objs.BClaims{}
	txHashes := [][][]byte{}

	for i := 0; i < nBlocks; i++ {

		txhash := crypto.Hasher([]byte(strconv.Itoa(i + 1)))
		txHshLst := [][]byte{txhash}
		txRoot, err := objs.MakeTxRoot(txHshLst)
		if err != nil {
			return nil, nil, err
		}
		txHashes = append(txHashes, txHshLst)

		prevBlockHash := crypto.Hasher([]byte("foo"))

		if i > 0 {
			prevBlockHash, err = chain[i-1].BlockHash()
			if err != nil {
				panic(fmt.Errorf("could not create prevBlockHash: %v\n", err))
			}
		}

		var stateRoot []byte
		err = ah.db.Update(func(txn *badger.Txn) error {
			stateRoot, err = ah.ApplyState(txn, 1, uint32(i)+1, nil)
			return err
		})
		if err != nil || stateRoot == nil {
			panic(err)
		}

		bclaims := &objs.BClaims{
			ChainID:    1,
			Height:     uint32(i) + 1,
			TxCount:    1,
			PrevBlock:  prevBlockHash,
			TxRoot:     txRoot,
			StateRoot:  stateRoot,
			HeaderRoot: crypto.Hasher([]byte("header root")),
		}
		chain = append(chain, bclaims)
	}

	return chain, txHashes, nil
}

func generateValidatorSet(t *testing.T) *objs.ValidatorSet {
	gpkj1, ok := big.NewInt(0).SetString("14395602319113363333690669395961581081803242358678131578916981232954633806960", 10)
	assert.True(t, ok)
	gpkj2, ok := big.NewInt(0).SetString("300089735810954642595088127891607498572672898349379085034409445552605516765", 10)
	assert.True(t, ok)
	gpkj3, ok := big.NewInt(0).SetString("17169409825226096532229555694191340178889298261881998623204757401596570351688", 10)
	assert.True(t, ok)
	gpkj4, ok := big.NewInt(0).SetString("19780380227412019371988923760536598779715024137904246485146692590642474692882", 10)
	assert.True(t, ok)

	v1 := objects.Validator{
		Account: common.HexToAddress("0x9AC1c9afBAec85278679fF75Ef109217f26b1417"),
		Index:   1,
		SharedKey: [4]*big.Int{
			gpkj1,
			gpkj2,
			gpkj3,
			gpkj4,
		},
	}

	gpkj1, ok = big.NewInt(0).SetString("21154017404198718862920160130737623556546602199694661996869957208062851500379", 10)
	assert.True(t, ok)
	gpkj2, ok = big.NewInt(0).SetString("19389833000731437962153734187923001234830293448701992540723746507685386979412", 10)
	assert.True(t, ok)
	gpkj3, ok = big.NewInt(0).SetString("21289029302611008572663530729853170393569891172031986702208364730022339833735", 10)
	assert.True(t, ok)
	gpkj4, ok = big.NewInt(0).SetString("15926764275937493411567546154328577890519582979565228998979506880914326856186", 10)
	assert.True(t, ok)

	v2 := objects.Validator{
		Account: common.HexToAddress("0x615695C4a4D6a60830e5fca4901FbA099DF26271"),
		Index:   2,
		SharedKey: [4]*big.Int{
			gpkj1,
			gpkj2,
			gpkj3,
			gpkj4,
		},
	}

	gpkj1, ok = big.NewInt(0).SetString("15079629603150363557558188402860791995814736941924946256968815481986722866449", 10)
	assert.True(t, ok)
	gpkj2, ok = big.NewInt(0).SetString("11164680325282976674805760467491699367894125557056167854003650409966070344792", 10)
	assert.True(t, ok)
	gpkj3, ok = big.NewInt(0).SetString("18616624374737795490811424594534628399519274885945803292205658067710235197668", 10)
	assert.True(t, ok)
	gpkj4, ok = big.NewInt(0).SetString("4331613963825409904165282575933135091483251249365224295595121580000486079984", 10)
	assert.True(t, ok)

	v3 := objects.Validator{
		Account: common.HexToAddress("0x63a6627b79813A7A43829490C4cE409254f64177"),
		Index:   3,
		SharedKey: [4]*big.Int{
			gpkj1,
			gpkj2,
			gpkj3,
			gpkj4,
		},
	}

	gpkj1, ok = big.NewInt(0).SetString("10875965504600753744265546216544158224793678652818595873355677460529088515116", 10)
	assert.True(t, ok)
	gpkj2, ok = big.NewInt(0).SetString("7912658035712558991777053184829906144303269569825235765302768068512975453162", 10)
	assert.True(t, ok)
	gpkj3, ok = big.NewInt(0).SetString("11324169944454120842956077363729540506362078469024985744551121054724657909930", 10)
	assert.True(t, ok)
	gpkj4, ok = big.NewInt(0).SetString("11005450895245397587287710270721947847266013997080161834700568409163476112947", 10)
	assert.True(t, ok)

	v4 := objects.Validator{
		Account: common.HexToAddress("0x16564cF3e880d9F5d09909F51b922941EbBbC24d"),
		Index:   4,
		SharedKey: [4]*big.Int{
			gpkj1,
			gpkj2,
			gpkj3,
			gpkj4,
		},
	}

	validators := []objects.Validator{v1, v2, v3, v4}
	ptrGroupKey := [4]*big.Int{
		v1.SharedKey[0],
		v1.SharedKey[1],
		v1.SharedKey[2],
		v1.SharedKey[3],
	}
	groupKey, err := bn256.MarshalG2Big(ptrGroupKey)
	assert.Nil(t, err)
	vs := &objs.ValidatorSet{
		GroupKey:   groupKey,
		Validators: make([]*objs.Validator, len(validators)),
		NotBefore:  0,
	}

	for _, validator := range validators {
		v := &objs.Validator{
			VAddr:      validator.Account.Bytes(),
			GroupShare: groupKey,
		}
		vs.Validators[validator.Index-1] = v
	}

	return vs
}

type synchronizerMock struct {
	sync.Mutex
}

// assert Synchronizer struct implements interfaces.Lockable interface
var _ interfaces.Lockable = &synchronizerMock{}
