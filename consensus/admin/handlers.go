package admin

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/dynamics"
	"github.com/MadBase/MadNet/errorz"
	"github.com/MadBase/MadNet/interfaces"
	"github.com/MadBase/MadNet/ipc"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// Todo: Retry logic on snapshot submission; this will cause deadlock
//		 if snapshots are taken too close together

// Handlers Is a private bus for internal service use.
// At this time the only reason to use this bus is to
// enable blockchain events to be fed into the system
// and for accusations to be fed out of the system and
// into the Ethereum blockchain.
type Handlers struct {
	sync.RWMutex

	// constructor assigned fields
	ctx         context.Context
	cancelFunc  func()
	logger      *logrus.Logger
	database    *db.Database
	chainID     uint32
	appHandler  interfaces.Application
	ethPubk     []byte
	secret      []byte
	ethAcct     []byte
	ReceiveLock chan interfaces.Lockable
	storage     dynamics.StorageGetter
	ipcServer   *ipc.Server
	ethHeights  chan uint32
	madHeaders  chan *objs.BlockHeader
	osListener  chan *objs.OwnState

	// other fields
	closeOnce sync.Once
	isInit    bool
	isSync    bool
}

// Init creates all fields and binds external services
func (ah *Handlers) Init(chainID uint32, database *db.Database, secret []byte, appHandler interfaces.Application, ethPubk []byte, storage dynamics.StorageGetter, ipcs *ipc.Server) {
	ctx := context.Background()
	subCtx, cancelFunc := context.WithCancel(ctx)
	ah.ctx = subCtx
	ah.cancelFunc = cancelFunc
	ah.logger = logging.GetLogger(constants.LoggerConsensus)
	ah.database = database
	ah.chainID = chainID
	ah.appHandler = appHandler
	ah.ethPubk = ethPubk
	ah.secret = utils.CopySlice(secret)
	ah.ethAcct = crypto.GetAccount(ethPubk)
	ah.ReceiveLock = make(chan interfaces.Lockable)
	ah.storage = storage
	ah.ipcServer = ipcs
	ah.ethHeights = make(chan uint32, 4)
	ah.madHeaders = make(chan *objs.BlockHeader, 4)
}

// Close shuts down all workers
func (ah *Handlers) Close() {
	ah.closeOnce.Do(func() {
		ah.cancelFunc()
	})
}

func (ah *Handlers) getLock() (interfaces.Lockable, bool) {
	select {
	case lock := <-ah.ReceiveLock:
		return lock, true
	case <-ah.ctx.Done():
		return nil, false
	}
}

// AddValidatorSet adds a validator set to the db
// This function also creates the first block and initializes
// the genesis state when the first validator set is written
func (ah *Handlers) AddValidatorSet(v *objs.ValidatorSet) error {
	mutex, ok := ah.getLock()
	if !ok {
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	return ah.database.Update(func(txn *badger.Txn) error {
		// Checking if we can exit earlier (mainly when reconstructing the chain
		// from ethereum data)
		{
			height := uint32(1)
			if v.NotBefore >= 1 {
				height = v.NotBefore
			}

			vSet, err := ah.database.GetValidatorSet(txn, height)
			if err != nil {
				if err != badger.ErrKeyNotFound {
					utils.DebugTrace(ah.logger, err)
					return err
				}
				// do nothing
			}
			bhHeight := height - 1
			if v.NotBefore == 0 {
				bhHeight = 1
			}
			bh, err := ah.database.GetCommittedBlockHeader(txn, bhHeight)
			if err != nil {
				if err != badger.ErrKeyNotFound {
					utils.DebugTrace(ah.logger, err)
					return err
				}
			}
			// If we have a committed blocker header, and the current validator
			// set in memory is equal to the validator set that we are
			// receiving, we are good and we don't need to execute the steps
			// below
			if bh != nil && vSet != nil && bytes.Equal(v.GroupKey, vSet.GroupKey) {
				return nil
			}
		}
		// Adding new validators in case of epoch boundary
		if v.NotBefore%constants.EpochLength == 0 {
			return ah.epochBoundaryValidator(txn, v)
		}
		// reset case (we received from ethereum an event with group key fields
		// all zeros).
		if bytes.Equal(v.GroupKey, make([]byte, len(v.GroupKey))) {
			return ah.database.SetValidatorSet(txn, v)
		}
		// Setting a new set of validator outside the epoch boundaries and after
		// a the reset case above
		return ah.AddValidatorSetEdgecase(txn, v)
	})
}

// AddValidatorSetEdgecase adds a validator set to the db if we have the
// expected block at the height 'v.NotBefore-1' (e.g syncing from the ethereum
// data). Otherwise, it will mark the change to happen in the future once we
// have the required block
func (ah *Handlers) AddValidatorSetEdgecase(txn *badger.Txn, v *objs.ValidatorSet) error {
	bh, err := ah.database.GetCommittedBlockHeader(txn, v.NotBefore-1)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			utils.DebugTrace(ah.logger, err)
			return err
		}
		return ah.database.SetValidatorSetPostApplication(txn, v, v.NotBefore)
	}
	rcert, err := bh.GetRCert()
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	isValidator, err := ah.initValidatorsRoundState(txn, v, rcert)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	// If we are not a validator we need to start our Round State
	if !isValidator {
		err = ah.initOwnRoundState(txn, v, rcert)
		if err != nil {
			utils.DebugTrace(ah.logger, err)
			return err
		}
	}
	err = ah.database.SetValidatorSet(txn, v)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	return nil
}

// AddSnapshot stores a snapshot to the database
func (ah *Handlers) AddSnapshot(bh *objs.BlockHeader, safeToProceedConsensus bool) error {
	ah.logger.Debugf("inside adminHandler.AddSnapshot")
	mutex, ok := ah.getLock()
	if !ok {
		return errors.New("could not get adminHandler lock")
	}
	mutex.Lock()
	defer mutex.Unlock()
	err := ah.database.Update(func(txn *badger.Txn) error {
		safeToProceed, err := ah.database.GetSafeToProceed(txn, bh.BClaims.Height+1)
		if err != nil {
			utils.DebugTrace(ah.logger, err)
			return err
		}
		if !safeToProceed {
			ah.logger.Debugf("Did validators change in the previous epoch:%v Setting is safe to proceed for height %d to: %v", !safeToProceedConsensus, bh.BClaims.Height+1, safeToProceedConsensus)
			// set that it's safe to proceed to the next block
			if err := ah.database.SetSafeToProceed(txn, bh.BClaims.Height+1, safeToProceedConsensus); err != nil {
				utils.DebugTrace(ah.logger, err)
				return err
			}
		}
		if bh.BClaims.Height > 1 {
			err = ah.database.SetSnapshotBlockHeader(txn, bh)
			if err != nil {
				utils.DebugTrace(ah.logger, err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	ah.logger.Debugf("successfully saved state on adminHandler.AddSnapshot")
	return nil
}

// UpdateDynamicStorage updates dynamic storage values.
func (ah *Handlers) UpdateDynamicStorage(txn *badger.Txn, key, value string, epoch uint32) error {
	mutex, ok := ah.getLock()
	if !ok {
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()

	update, err := dynamics.NewUpdate(key, value, epoch)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	err = ah.storage.UpdateStorage(txn, update)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	return nil
}

// IsInitialized returns if the database has been initialized yet
func (ah *Handlers) IsInitialized() bool {
	ah.RLock()
	defer ah.RUnlock()
	return ah.isInit
}

// IsSynchronized returns if the ethereum BC has been/is synchronized
func (ah *Handlers) IsSynchronized() bool {
	ah.RLock()
	defer ah.RUnlock()
	return ah.isSync
}

// SetSynchronized allows the BC monitor to set the sync state for ethereum
func (ah *Handlers) SetSynchronized(v bool) {
	ah.Lock()
	defer ah.Unlock()
	ah.isSync = v
}

// RegisterSnapshotCallbacks allows callbacks to be registered that are used for the snapshotting logic
func (ah *Handlers) RegisterSnapshotCallbacks(
	newSS func(bh *objs.BlockHeader) error,
	resumeSS func(bh *objs.CachedSnapshotTx) error,
	getSSEthHeight func(epoch *big.Int) (*big.Int, error),
	refreshTx func(context.Context, *types.Transaction, int) (*types.Transaction, error),
) {
	ah.database.SubscribeBroadcastBlockHeader(ah.ctx, func(b []byte) error {
		bh := &objs.BlockHeader{}
		err := bh.UnmarshalBinary(b)

		if err != nil {
			ah.logger.Errorf("Could not unmarshal received mad block header: %v", err)
			return nil // swallow error (else subscription is lost)
		}

		for { // push mad block to channel, remove oldest item if full
			select {

			case ah.madHeaders <- bh:
				return nil

			default:

				select {
				case <-ah.madHeaders:
				default:
				}

			}
		}
	})

	go func() { // continously check if snapshots should be made each time we see a new mad or eth block
		var myAddr []byte
		var bh *objs.BlockHeader
		var eh uint32
		var lastSS struct {
			ethHeight uint32
			madEpoch  uint32
		}

		ah.logger.Infof("Starting snapshot check loop...")
		defer ah.logger.Infof("Ended snapshot check loop")

		err := ah.database.View(func(txn *badger.Txn) (err error) {
			os, err := ah.database.GetOwnState(txn)
			if err != nil {
				return err
			}
			bh = os.SyncToBH
			myAddr = os.VAddr
			return nil
		})
		if err == badger.ErrKeyNotFound {
			ah.Lock()
			ah.osListener = make(chan *objs.OwnState, 4)
			ah.Unlock()

			os := <-ah.osListener
			bh = os.SyncToBH
			myAddr = os.VAddr

			ah.Lock()
			ah.osListener = nil
			ah.Unlock()
		} else if err != nil {
			ah.logger.Errorf("Failed to initialize snapshot check loop. Could not read state from db: %v", err)
			return
		}

		for {
			done := ah.ctx.Done()
			for { // receive block headers and eth heights until we have at least one of both
				select {
				case bh = <-ah.madHeaders:
				case eh = <-ah.ethHeights:
				case <-done:
					return
				}
				ah.logger.Info(eh, bh)

				if eh != 0 && bh != nil && bh.BClaims != nil {
					break
				}
			}

			for loop := true; loop; { // receive all remaining block headers and eth heights
				select {
				case bh = <-ah.madHeaders:
				case eh = <-ah.ethHeights:
				case <-done:
					return
				default:
					loop = false
				}
			}

			if bh.BClaims.Height%constants.EpochLength != 0 { // only check if madnet is stuck at epoch boundary
				continue
			}

			ah.logger.Debugf("Snapshot check started for height %v, epoch %v...", bh.BClaims.Height, bh.BClaims.Height/constants.EpochLength)

			var vs *objs.ValidatorSet
			var tx *objs.CachedSnapshotTx

			err := ah.database.View(func(txn *badger.Txn) (err error) {
				tx, err = ah.database.GetSnapshotTx(txn)
				if err != nil && err != badger.ErrKeyNotFound {
					return err
				}
				vs, err = ah.database.GetValidatorSet(txn, bh.BClaims.Height+1)
				return err
			})
			if err != nil {
				ah.logger.Errorf("Snapshot check failed. Could not read from db: %v", err)
				continue
			}

			index := -1
			for i := 0; i < len(vs.Validators); i++ {
				val := vs.Validators[i]
				if bytes.Equal(val.VAddr, myAddr) {
					index = i
					break
				}
			}
			if index == -1 {
				ah.logger.Debugf("Snapshot not needed. Not a validator")
				continue
			}

			currentEpoch := bh.BClaims.Height / constants.EpochLength
			lastEpoch := currentEpoch - 1
			if lastSS.madEpoch != lastEpoch {
				ssEh, err := getSSEthHeight(big.NewInt(int64(lastEpoch)))
				if err != nil {
					ah.logger.Errorf("Snapshot check failed. Could not retrieve ethheight of last snapshot: %v", err)
					continue
				}
				lastSS.ethHeight = uint32(ssEh.Int64())
				lastSS.madEpoch = lastEpoch
			}

			ethBlocksSinceDesperation := int(eh) - int(lastSS.ethHeight) - constants.SnapshotDesperationDelay

			blockhash, err := bh.BClaims.BlockHash()
			if err != nil {
				ah.logger.Errorf("Snapshot check failed. Could not marshal blockhash")
				continue
			}

			if !mayValidatorSnapshot(len(vs.Validators), index, ethBlocksSinceDesperation, blockhash) {
				ah.logger.Debugf("Snapshot not needed. Not my turn")
				continue
			}

			ah.logger.Infof("Snapshot needed. Starting snapshot submission logic...")

			switch {

			// send new transaction if no relevant cached transaction exists
			case tx == nil || tx.MadEpoch != currentEpoch:
				ah.logger.Infof("Persisting snapshot for height %v...", bh.BClaims.Height)
				err := newSS(bh)
				if err != nil {
					ah.logger.Errorf("Could not persist new snapshot: %v", err)
					continue
				}
				ah.logger.Infof("Successfully persisted snapshot for height %v...", bh.BClaims.Height)

			// if cached transaction was created but not sent, send
			case tx.State == objs.SnapshotTxCreated:
				ah.logger.Infof("Submitting previously unsent snapshot tx for height %v...", bh.BClaims.Height)
				err := resumeSS(tx)
				if err != nil {
					ah.logger.Errorf("Could not resubmit previously created snapshot tx: %v", err)
					continue
				}
				ah.logger.Infof("Successfully persisted previously unsent snapshot tx for height %v...", bh.BClaims.Height)

			// if cached transaction was sent but not verified for a while
			case tx.State == objs.SnapshotTxSubmitted:
				newTx, err := refreshTx(ah.ctx, tx.Tx, int(eh)-int(tx.EthHeight))
				if newTx == nil {
					ah.logger.Debugf("No refresh needed")
					continue
				}
				if err != nil {
					ah.logger.Errorf("Failed to refresh transaction: %v", err)
					continue
				}
				ah.logger.Infof("Resubmitting stale snapshot tx for height %v...", bh.BClaims.Height)

				err = resumeSS(&objs.CachedSnapshotTx{
					Tx:       newTx,
					State:    objs.SnapshotTxCreated,
					MadEpoch: currentEpoch,
				})
				if err != nil {
					ah.logger.Errorf("Could not resubmit stale snapshot tx snapshot: %v", err)
					continue
				}
				ah.logger.Infof("Successfully persisted previously stale snapshot tx for height %v...", bh.BClaims.Height)
			default:
			}
		}
	}()
}

// AddPrivateKey stores a private key from an EthDKG run into an encrypted
// keystore in the DB
func (ah *Handlers) AddPrivateKey(pk []byte, curveSpec constants.CurveSpec) error {
	mutex, ok := ah.getLock()
	if !ok {
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	// ah.logger.Error("!!! OPEN AddPrivateKey TXN")
	// defer func() { ah.logger.Error("!!! CLOSE AddPrivateKey TXN") }()
	err := ah.database.Update(func(txn *badger.Txn) error {
		switch curveSpec {
		case constants.CurveSecp256k1:
			privk := utils.CopySlice(pk)
			// secp key
			signer := crypto.Secp256k1Signer{}
			err := signer.SetPrivk(privk)
			if err != nil {
				return err
			}
			pubkey, err := signer.Pubkey()
			if err != nil {
				return err
			}
			name := crypto.GetAccount(pubkey)
			ec := &objs.EncryptedStore{
				Name:      name,
				ClearText: privk,
				Kid:       constants.AdminHandlerKid(),
			}
			err = ec.Encrypt(ah)
			if err != nil {
				return err
			}
			return ah.database.SetEncryptedStore(txn, ec)
		case constants.CurveBN256Eth:
			privk := utils.CopySlice(pk)
			// bn key
			signer := crypto.BNGroupSigner{}
			err := signer.SetPrivk(privk)
			if err != nil {
				return err
			}
			pubkey, err := signer.PubkeyShare()
			if err != nil {
				return err
			}
			ec := &objs.EncryptedStore{
				Name:      pubkey,
				ClearText: privk,
				Kid:       constants.AdminHandlerKid(),
			}
			err = ec.Encrypt(ah)
			if err != nil {
				return err
			}
			return ah.database.SetEncryptedStore(txn, ec)
		default:
			panic("not an allowed curve type")
		}
	})
	if err != nil {
		panic(err)
	}
	return nil
}

// updateEthHeight
func (ah *Handlers) UpdateEthHeight(ethHeight uint32) {
	for { // push ethHeight to channel, remove oldest item if full
		select {
		case ah.ethHeights <- ethHeight:
			return
		default:
			select {
			case <-ah.ethHeights:
			default:
			}
		}
	}
}

// GetPrivK returns an decrypted private key from an EthDKG run to the caller
func (ah *Handlers) GetPrivK(name []byte) ([]byte, error) {
	var privk []byte
	err := ah.database.View(func(txn *badger.Txn) error {
		ec, err := ah.database.GetEncryptedStore(txn, name)
		if err != nil {
			return err
		}
		err = ec.Decrypt(ah)
		if err != nil {
			return err
		}
		privk = utils.CopySlice(ec.ClearText)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return privk, nil
}

// GetKey allows the admin handler to act as a key resolver for decrypting
// stored private keys
func (ah *Handlers) GetKey(kid []byte) ([]byte, error) {
	out := make([]byte, len(ah.secret))
	copy(out[:], ah.secret)
	return out, nil
}

// InitializationMonitor polls the database for the existence of a snapshot
// It sets IsInitialized when one is found and returns
func (ah *Handlers) InitializationMonitor(closeChan <-chan struct{}) {
	ah.logger.Debug("InitializationMonitor loop starting")
	defer func() {
		ah.logger.Debug("InitializationMonitor loop stopping")
	}()

	for {
		ok, err := func() (bool, error) {
			select {
			case <-closeChan:
				return false, errorz.ErrClosing
			case <-ah.ctx.Done():
				return false, nil
			case <-time.After(2 * time.Second):
				err := ah.database.View(func(txn *badger.Txn) error {
					_, err := ah.database.GetLastSnapshot(txn)
					return err
				})
				if err != nil {
					return false, nil
				}
				ah.Lock()
				ah.isInit = true
				ah.Unlock()
				return true, nil
			}
		}()
		if err != nil {
			return
		}
		if ok {
			break
		}
	}
}

func mayValidatorSnapshot(numValidators int, myIdx int, blocksSinceDesperation int, blockhash []byte) bool {
	var numValidatorsAllowed int = 1
	for i := int(blocksSinceDesperation); i >= 0; {
		i -= constants.SnapshotDesperationFactor / numValidatorsAllowed
		numValidatorsAllowed++

		if numValidatorsAllowed > numValidators/3 {
			break
		}
	}

	// use the random nature of blockhash to deterministically define the range of validators that are allowed to make a snapshot
	rand := (&big.Int{}).SetBytes(blockhash)

	start := int((&big.Int{}).Mod(rand, big.NewInt(int64(numValidators))).Int64())
	end := (start + numValidatorsAllowed) % numValidators

	if end > start {
		return myIdx >= start && myIdx < end
	} else {
		return myIdx >= start || myIdx < end
	}
}

func (ah *Handlers) epochBoundaryValidator(txn *badger.Txn, v *objs.ValidatorSet) error {
	bh, err := ah.database.GetSnapshotByHeight(txn, v.NotBefore)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			utils.DebugTrace(ah.logger, err)
			return err
		}
	}
	if bh == nil || v.NotBefore == 0 {
		bh, err = ah.initDB(txn, v)
		if err != nil {
			utils.DebugTrace(ah.logger, err)
			return err
		}
	}
	rcert, err := bh.GetRCert()
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	isValidator, err := ah.initValidatorsRoundState(txn, v, rcert)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	// If we are not a validator we need to start our Round State
	if !isValidator {
		err = ah.initOwnRoundState(txn, v, rcert)
		if err != nil {
			utils.DebugTrace(ah.logger, err)
			return err
		}
	}

	// fix zero epoch event in chain
	switch v.NotBefore {
	case 0:
		v.NotBefore = 1
	default:
		v.NotBefore = rcert.RClaims.Height
	}

	err = ah.database.SetValidatorSet(txn, v)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	err = ah.database.SetSafeToProceed(txn, v.NotBefore, true)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	return nil
}

// Re-Initializes our own Round State object
func (ah *Handlers) initOwnRoundState(txn *badger.Txn, v *objs.ValidatorSet, rcert *objs.RCert) error {
	rs, err := ah.database.GetCurrentRoundState(txn, ah.ethAcct)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	if (rs == nil) || (!bytes.Equal(rs.GroupKey, v.GroupKey) && v.NotBefore >= rcert.RClaims.Height) {
		rs = &objs.RoundState{
			VAddr:      ah.ethAcct,
			GroupKey:   v.GroupKey,
			GroupShare: make([]byte, constants.CurveBN256EthPubkeyLen),
			GroupIdx:   0,
			RCert:      rcert,
		}
	}
	err = ah.database.SetCurrentRoundState(txn, rs)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return err
	}
	return nil
}

// Re-Initializes all the validators Round State objects
func (ah *Handlers) initValidatorsRoundState(txn *badger.Txn, v *objs.ValidatorSet, rcert *objs.RCert) (bool, error) {
	isValidator := false
	for i := 0; i < len(v.Validators); i++ {
		val := v.Validators[i]
		rs, err := ah.database.GetCurrentRoundState(txn, val.VAddr)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				utils.DebugTrace(ah.logger, err)
				return false, err
			}
		}
		rcertTemp := rcert
		if rs != nil && rs.RCert.RClaims.Height > rcert.RClaims.Height {
			rcertTemp = rs.RCert
		}
		rs = &objs.RoundState{
			VAddr:      utils.CopySlice(val.VAddr),
			GroupKey:   utils.CopySlice(v.GroupKey),
			GroupShare: utils.CopySlice(val.GroupShare),
			GroupIdx:   uint8(i),
			RCert:      rcertTemp,
		}
		err = ah.database.SetCurrentRoundState(txn, rs)
		if err != nil {
			utils.DebugTrace(ah.logger, err)
			return false, err
		}
		if bytes.Equal(rs.VAddr, ah.ethAcct) {
			isValidator = true
		}
	}
	return isValidator, nil
}

// Init the validators DB and objects
func (ah *Handlers) initDB(txn *badger.Txn, v *objs.ValidatorSet) (*objs.BlockHeader, error) {
	stateRoot, err := ah.appHandler.ApplyState(txn, ah.chainID, 1, nil)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	vlst := [][]byte{}
	for i := 0; i < len(v.Validators); i++ {
		val := v.Validators[i]
		vlst = append(vlst, crypto.Hasher(val.VAddr))
	}
	prevBlock, err := objs.MakeTxRoot(vlst)
	if err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	bh := &objs.BlockHeader{
		BClaims: &objs.BClaims{
			ChainID:    ah.chainID,
			Height:     1,
			PrevBlock:  prevBlock,
			StateRoot:  stateRoot,
			HeaderRoot: make([]byte, constants.HashLen),
			TxRoot:     txRoot,
		},
		SigGroup: make([]byte, constants.CurveBN256EthSigLen),
		TxHshLst: [][]byte{},
	}
	if err := ah.database.SetSnapshotBlockHeader(txn, bh); err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	if err := ah.database.SetCommittedBlockHeader(txn, bh); err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	ownState := &objs.OwnState{
		VAddr:             ah.ethAcct,
		SyncToBH:          bh,
		MaxBHSeen:         bh,
		CanonicalSnapShot: bh,
		PendingSnapShot:   bh,
	}
	if err := ah.database.SetOwnState(txn, ownState); err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	ownValidatingState := new(objs.OwnValidatingState)
	ownValidatingState.SetRoundStarted()
	if err := ah.database.SetOwnValidatingState(txn, ownValidatingState); err != nil {
		utils.DebugTrace(ah.logger, err)
		return nil, err
	}
	return bh, nil
}
