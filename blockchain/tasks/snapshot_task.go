package tasks

import (
	"context"
	"fmt"
	"sync"

	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// SnapshotTask pushes a snapshot to Ethereum
type SnapshotTask struct {
	sync.RWMutex
	acct        accounts.Account
	blockHeader *objs.BlockHeader
	rawBclaims  []byte
	consDb      *db.Database
	tx          *types.Transaction
	txState     objs.SnapshotTxState
	madEpoch    uint32
	ethHeight   uint32
}

func NewSnapshotTask(account accounts.Account, db *db.Database, bh *objs.BlockHeader) *SnapshotTask {
	return &SnapshotTask{acct: account, consDb: db, blockHeader: bh}
}

func NewSnapshotTaskFromCached(account accounts.Account, db *db.Database, ss *objs.CachedSnapshotTx) *SnapshotTask {
	return &SnapshotTask{acct: account, consDb: db, tx: ss.Tx, txState: ss.State}
}

func (t *SnapshotTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum, _ interface{}) error {
	if t.txState < objs.SnapshotTxCreated || t.tx == nil {
		rawBClaims, err := t.blockHeader.BClaims.MarshalBinary()
		if err != nil {
			return fmt.Errorf("unable to marshal block header: %v", err)
		}

		t.Lock()
		defer t.Unlock()

		t.rawBclaims = rawBClaims
	}
	return nil
}
func (t *SnapshotTask) DoWork(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	t.DoTask(ctx, logger, eth)
	return nil // swallow error
}

func (t *SnapshotTask) DoRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	t.DoTask(ctx, logger, eth)
	return nil // swallow error
}

func (t *SnapshotTask) DoTask(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	t.Lock()
	defer t.Unlock()

	if t.txState < objs.SnapshotTxCreated || t.tx == nil {
		err := t.createSnapshotTx(ctx, logger, eth)
		if err != nil {
			return err
		}
	}

	if t.txState < objs.SnapshotTxSubmitted {
		err := t.submitSnapshotTx(ctx, logger, eth)
		if err != nil {
			return err
		}
	}

	err := t.awaitSnapshotTx(ctx, logger, eth)
	return err
}

func (t *SnapshotTask) ShouldRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) bool {

	t.RLock()
	defer t.RUnlock()

	opts := eth.GetCallOpts(ctx, t.acct)

	epoch, err := eth.Contracts().Snapshots().GetEpoch(opts)
	if err != nil {
		logger.Errorf("Failed to determine current epoch: %v", err)
		return true
	}

	height, err := eth.Contracts().Snapshots().GetMadnetHeightFromSnapshot(opts, epoch)
	if err != nil {
		logger.Errorf("Failed to determine height: %v", err)
		return true
	}

	// This means the block height we want to snapshot is older than (or same as) what's already been snapshotted
	if t.blockHeader.BClaims.Height != 0 && t.blockHeader.BClaims.Height < uint32(height.Uint64()) {
		return false
	}

	return true
}

func (*SnapshotTask) DoDone(logger *logrus.Entry) {
}

func (t *SnapshotTask) createSnapshotTx(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	if t.tx != nil || t.txState >= objs.SnapshotTxCreated {
		return fmt.Errorf("snapshot already created")
	}

	bal, err := eth.GetBalance(t.acct.Address)
	if err != nil {
		return err
	}

	txnOpts, err := eth.GetTransactionOpts(ctx, t.acct)
	if err != nil {
		return fmt.Errorf("failed to generate transaction options: %v", err)
	}
	logger.Debugf("Snapshot tx options {acct:%v, gas:%v from:%v, balance:%v, sigGroup:%v, bclaims:%v}", t.acct, txnOpts.GasLimit, txnOpts.From, bal, t.blockHeader.SigGroup, t.rawBclaims)
	tx, err := eth.Contracts().Snapshots().Snapshot(txnOpts, t.blockHeader.SigGroup, t.rawBclaims)
	if err != nil {
		return err
	}

	epoch := t.blockHeader.BClaims.Height / constants.EpochLength
	err = t.consDb.Update(func(txn *badger.Txn) error {
		return t.consDb.SetSnapshotTx(txn, &objs.CachedSnapshotTx{
			State:    objs.SnapshotTxCreated,
			Tx:       tx,
			MadEpoch: epoch,
		})
	})
	if err != nil {
		return fmt.Errorf("failed to persist transaction: %v", err)
	}

	t.tx = tx
	t.txState = objs.SnapshotTxCreated
	t.madEpoch = epoch
	return nil
}

func (t *SnapshotTask) submitSnapshotTx(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	if t.tx == nil || t.txState < objs.SnapshotTxCreated {
		return fmt.Errorf("no transaction created")
	}

	h, err := eth.GetCurrentHeight(ctx)
	if err != nil {
		return err
	}

	err = eth.Queue().QueueTransactionSync(ctx, t.tx)
	if err != nil {
		return fmt.Errorf("failed to queue transaction: %v", err)
	}

	err = t.consDb.Update(func(txn *badger.Txn) error {
		return t.consDb.SetSnapshotTx(txn, &objs.CachedSnapshotTx{
			State:     objs.SnapshotTxSubmitted,
			Tx:        t.tx,
			MadEpoch:  t.madEpoch,
			EthHeight: uint32(h),
		})
	})
	if err != nil {
		return fmt.Errorf("failed to persist transaction being sent: %v", err)
	}

	t.txState = objs.SnapshotTxSubmitted
	t.ethHeight = uint32(h)
	return nil
}

func (t *SnapshotTask) awaitSnapshotTx(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	if t.tx == nil || t.txState < objs.SnapshotTxSubmitted {
		return fmt.Errorf("no transaction sent")
	}

	rcpt, err := eth.Queue().WaitTransaction(ctx, t.tx)
	if err != nil {
		return fmt.Errorf("snapshot failed to retreive receipt: %v", err)
	}
	if rcpt.Status != 1 {
		return fmt.Errorf("snapshot receipt status %v", rcpt.Status)
	}

	err = t.consDb.Update(func(txn *badger.Txn) error {
		return t.consDb.SetSnapshotTx(txn, &objs.CachedSnapshotTx{
			State:     objs.SnapshotTxVerified,
			Tx:        t.tx,
			MadEpoch:  t.madEpoch,
			EthHeight: t.ethHeight,
		})
	})

	if err != nil {
		return fmt.Errorf("failed to persist transaction being succeeded: %v", err)
	}
	return nil
}
