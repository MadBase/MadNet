package localrpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/MadBase/MadNet/application"
	"github.com/MadBase/MadNet/application/deposit"
	"github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/application/objs/uint256"
	"github.com/MadBase/MadNet/application/utxohandler"
	"github.com/MadBase/MadNet/config"
	"github.com/MadBase/MadNet/consensus"
	"github.com/MadBase/MadNet/consensus/admin"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/dman"
	"github.com/MadBase/MadNet/consensus/evidence"
	"github.com/MadBase/MadNet/consensus/gossip"
	"github.com/MadBase/MadNet/consensus/lstate"
	"github.com/MadBase/MadNet/consensus/request"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	mncrypto "github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/ipc"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/peering"
	"github.com/MadBase/MadNet/proto"
	"github.com/MadBase/MadNet/utils"
	mnutils "github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"
)

var address string = "localhost:8884"
var timeout time.Duration = time.Second * 10
var utxoIDs [][]byte
var account []byte
var owner *objs.ValueStoreOwner
var signer *crypto.Secp256k1Signer
var txHash []byte
var pubKey []byte
var tx *objs.Tx
var geth = exec.Cmd{}
var err error
var badgerDB *badger.DB
var srpc *Handlers
var lrpc *Client
var stateDB *db.Database
var ctx context.Context

type mutexUint32 struct {
	sync.RWMutex
	value uint32
}

func TestMain(m *testing.M) {

	file, _ := os.Open("validator.toml")
	bs, _ := ioutil.ReadAll(file)
	reader := bytes.NewReader(bs)
	viper.SetConfigType("toml")
	viper.ReadConfig(reader)
	viper.Unmarshal(&config.Configuration)
	chainID = uint32(config.Configuration.Chain.ID)

	signer = &crypto.Secp256k1Signer{}
	signer.SetPrivk(crypto.Hasher([]byte("secret")))
	pubKey, err = signer.Pubkey()

	// dbFile := "../scripts/generated/stateDBs/validator1"
	dbFile := "testDB"

	ctx = context.Background()
	app := &application.Application{}
	// storage := &dynamics.Storage{}
	appDepositHandler := &deposit.Handler{} // watches ETH blockchain about deposits
	stateDB = &db.Database{}
	consTxPool := &evidence.Pool{}
	consGossipHandlers := &gossip.Handlers{}
	// consensus p2p comm
	consReqClient := &request.Client{}
	// consReqHandler := &request.Handler{}

	// core of consensus algorithm: where outside stake relies, how gossip ends up, how state modifications occur
	consLSEngine := &lstate.Engine{}
	consLSHandler := &lstate.Handlers{}
	consDlManager := &dman.DMan{}
	consGossipClient := &gossip.Client{}
	consAdminHandlers := &admin.Handlers{}

	consSync := &consensus.Synchronizer{}

	rawStateDB, _ := mnutils.OpenBadger(ctx.Done(), dbFile, false)
	rawTxPoolDB, _ := mnutils.OpenBadger(ctx.Done(), "", true)

	stateDB.Init(rawStateDB)
	consTxPool.Init(stateDB)
	storage := makeMockStorageGetter()
	ipcServer := ipc.NewServer(config.Configuration.Firewalld.SocketFile)

	appDepositHandler.Init()
	// storage.Init(stateDB, nil)
	app.Init(stateDB, rawTxPoolDB, appDepositHandler, storage)
	pubKey, _ := signer.Pubkey()
	peerManager := initPeerManager(consGossipHandlers, nil)
	/* 	consReqClient.Init(peerManager.P2PClient(), storage)
	   	consReqHandler.Init(stateDB, app, storage)
	   	consDlManager.Init(stateDB, app, consReqClient)
	   	consLSHandler.Init(stateDB, consDlManager) */

	consGossipHandlers.Init(1337, stateDB, peerManager.P2PClient(), app, consLSHandler, storage)
	consGossipClient.Init(stateDB, peerManager.P2PClient(), app, storage)
	consAdminHandlers.Init(1337, stateDB, mncrypto.Hasher([]byte(config.Configuration.Validator.SymmetricKey)), app, pubKey, storage, ipcServer)
	consLSEngine.Init(stateDB, consDlManager, app, signer, consAdminHandlers, pubKey, consReqClient, storage)
	consSync.Init(stateDB, nil, nil, consGossipClient, consGossipHandlers, consTxPool, consLSEngine, app, consAdminHandlers, peerManager, storage)
	consSync.Start()

	srpc = &Handlers{
		ctx:         ctx,
		cancelCtx:   nil,
		database:    stateDB,
		sstore:      nil,
		AppHandler:  app,
		GossipBus:   consGossipHandlers,
		Storage:     storage,
		logger:      nil,
		ethAcct:     nil,
		EthPubk:     nil,
		safeHandler: func() bool { return true },
		safecount:   1,
	}
	srpc.Init(srpc.database, srpc.AppHandler, srpc.GossipBus, srpc.EthPubk, srpc.safeHandler, srpc.Storage)
	go srpc.Start()
	defer srpc.Stop()

	fmt.Println("ADDD", config.Configuration.Transport.LocalStateListeningAddress)

	lrpc = &Client{
		Mutex:       sync.Mutex{},
		closeChan:   nil,
		closeOnce:   sync.Once{},
		Address:     config.Configuration.Transport.LocalStateListeningAddress,
		TimeOut:     timeout,
		conn:        nil,
		client:      nil,
		wg:          sync.WaitGroup{},
		isConnected: false,
	}
	go lrpc.Connect(ctx)
	defer lrpc.Close()

	localStateServer := initLocalStateServer(srpc)
	go localStateServer.Serve()
	defer localStateServer.Close()

	time.Sleep(1 * time.Second)

	utxoIDs, account, txHash = insertTestUTXO()
	fmt.Println("utxoIDs:" + hex.EncodeToString(utxoIDs[0]))
	fmt.Println("account:" + hex.EncodeToString(account))
	fmt.Println("txHash:" + hex.EncodeToString(txHash))

	/* 	utxoIDs, _ = hex.DecodeString("cfa168a40031808356aa7d26eed4e4d3fef46ce02e0969be2c21a80c3fabd7a8")
	   	account, _ = hex.DecodeString("38e959391dd8598ae80d5d6d114a7822a09d313a")
	   	txHash, _ = hex.DecodeString("63f98d384fbb0fad03b3456b93c7605b4c2296358fdae857eb0d1a8ef479713a")
	*/
	/*
		paginationToken := objs.PaginationToken{
			LastPaginatedType: objs.LastPaginatedDeposit,
			TotalValue:        uint256.Zero(),
			LastKey:           utxoIDs[0],
		}

		binaryPaginationToken, err := paginationToken.MarshalBinary()
		if err != nil {
			fmt.Errorf("could not create paginationToken %v\n", err)
		}
		fmt.Println("binaryPaginationToken", hex.EncodeToString(binaryPaginationToken))
	*/

	//Start tests after validator is running
	exitVal := m.Run()

	os.Exit(exitVal)
}

func initPeerManager(consGossipHandlers *gossip.Handlers, consReqHandler *request.Handler) *peering.PeerManager {
	p2pDispatch := proto.NewP2PDispatch()

	peerManager, err := peering.NewPeerManager(
		proto.NewGeneratedP2PServer(p2pDispatch),
		uint32(config.Configuration.Chain.ID),
		config.Configuration.Transport.PeerLimitMin,
		config.Configuration.Transport.PeerLimitMax,
		config.Configuration.Transport.FirewallMode,
		config.Configuration.Transport.FirewallHost,
		config.Configuration.Transport.P2PListeningAddress,
		config.Configuration.Transport.PrivateKey,
		config.Configuration.Transport.UPnP)
	if err != nil {
		panic(err)
	}
	p2pDispatch.RegisterP2PGetPeers(peerManager)
	p2pDispatch.RegisterP2PGossipTransaction(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipProposal(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipPreVote(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipPreVoteNil(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipPreCommit(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipPreCommitNil(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipNextRound(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipNextHeight(consGossipHandlers)
	p2pDispatch.RegisterP2PGossipBlockHeader(consGossipHandlers)
	p2pDispatch.RegisterP2PGetBlockHeaders(consReqHandler)
	p2pDispatch.RegisterP2PGetMinedTxs(consReqHandler)
	p2pDispatch.RegisterP2PGetPendingTxs(consReqHandler)
	p2pDispatch.RegisterP2PGetSnapShotNode(consReqHandler)
	p2pDispatch.RegisterP2PGetSnapShotStateData(consReqHandler)
	p2pDispatch.RegisterP2PGetSnapShotHdrNode(consReqHandler)

	return peerManager
}

func initLocalStateServer(localStateHandler *Handlers) *Handler {
	localStateDispatch := proto.NewLocalStateDispatch()
	localStateServer, err := NewStateServerHandler(
		logging.GetLogger(constants.LoggerTransport),
		config.Configuration.Transport.LocalStateListeningAddress,
		proto.NewGeneratedLocalStateServer(localStateDispatch),
	)
	if err != nil {
		panic(err)
	}
	localStateDispatch.RegisterLocalStateGetBlockNumber(localStateHandler)
	localStateDispatch.RegisterLocalStateGetEpochNumber(localStateHandler)
	localStateDispatch.RegisterLocalStateGetBlockHeader(localStateHandler)
	localStateDispatch.RegisterLocalStateGetChainID(localStateHandler)
	localStateDispatch.RegisterLocalStateSendTransaction(localStateHandler)
	localStateDispatch.RegisterLocalStateGetValueForOwner(localStateHandler)
	localStateDispatch.RegisterLocalStateGetUTXO(localStateHandler)
	localStateDispatch.RegisterLocalStateGetTransactionStatus(localStateHandler)
	localStateDispatch.RegisterLocalStateGetMinedTransaction(localStateHandler)
	localStateDispatch.RegisterLocalStateGetPendingTransaction(localStateHandler)
	localStateDispatch.RegisterLocalStateGetRoundStateForValidator(localStateHandler)
	localStateDispatch.RegisterLocalStateGetValidatorSet(localStateHandler)
	localStateDispatch.RegisterLocalStateIterateNameSpace(localStateHandler)
	localStateDispatch.RegisterLocalStateGetData(localStateHandler)
	localStateDispatch.RegisterLocalStateGetTxBlockNumber(localStateHandler)
	localStateDispatch.RegisterLocalStateGetFees(localStateHandler)

	return localStateServer
}

func insertTestUTXO() ([][]byte, []byte, []byte) {
	if err != nil {
		fmt.Printf("could not create signer %v \n", err)
	}
	accountAddress := crypto.GetAccount(pubKey)
	owner := &objs.ValueStoreOwner{
		SVA:       objs.ValueStoreSVA,
		CurveSpec: constants.CurveSecp256k1,
		Account:   accountAddress,
	}
	hndlr := utxohandler.NewUTXOHandler(stateDB.DB())
	err = hndlr.Init(1)
	if err != nil {
		fmt.Printf("could not create utxo handler %v \n", err)
	}
	value, _ := new(uint256.Uint256).FromUint64(6)
	vs := &objs.ValueStore{
		VSPreImage: &objs.VSPreImage{
			TXOutIdx: constants.MaxUint32,
			Value:    value,
			ChainID:  1337,
			Owner:    owner,
		},
		TxHash: utils.ForceSliceToLength([]byte(strconv.Itoa(1)), constants.HashLen),
	}
	utxoDep := &objs.TXOut{}
	err = utxoDep.NewValueStore(vs)
	if err != nil {
		fmt.Printf("Could not create ValueStore %v \n", err)
	}
	tx, txHash := makeTxs(signer, vs)
	utxoIDs, err := tx.GeneratedUTXOID()
	if err != nil {
		fmt.Printf("Could not get utxoIds %v \n", err)
	}
	err = stateDB.Update(func(txn *badger.Txn) error {
		_, err := hndlr.ApplyState(txn, []*objs.Tx{tx}, 2)
		if err != nil {
			fmt.Printf("Could not validate %v \n", err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not update DB %v \n", err)
	}
	return utxoIDs, accountAddress, txHash
}

func makeTxs(s objs.Signer, v *objs.ValueStore) (*objs.Tx, []byte) {
	txIn, err := v.MakeTxIn()
	value, err := v.Value()
	chainID, err := txIn.ChainID()
	pubkey, err := s.Pubkey()
	tx = &objs.Tx{}
	tx.Vin = []*objs.TXIn{txIn}
	newValueStore := &objs.ValueStore{
		VSPreImage: &objs.VSPreImage{
			ChainID: chainID,
			Value:   value,
			Owner: &objs.ValueStoreOwner{
				SVA:       objs.ValueStoreSVA,
				CurveSpec: constants.CurveSecp256k1,
				Account:   crypto.GetAccount(pubkey)},
			TXOutIdx: 0,
			Fee:      new(uint256.Uint256).SetZero(),
		},
		TxHash: make([]byte, constants.HashLen),
	}
	newUTXO := &objs.TXOut{}
	err = newUTXO.NewValueStore(newValueStore)
	tx.Vout = append(tx.Vout, newUTXO)
	tx.Fee = uint256.Zero()
	err = tx.SetTxHash()
	err = v.Sign(tx.Vin[0], s)
	if err != nil {
		fmt.Printf("Could not create Txs %v \n", err)
	}
	return tx, tx.Vin[0].TXInLinker.TxHash
}
