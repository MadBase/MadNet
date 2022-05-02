package tasks

import (
	"context"
	"errors"
	dangerousRand "math/rand"
	"sync"
	"time"

	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/sirupsen/logrus"
)

// SnapshotTask pushes a snapshot to Ethereum
type SnapshotTask struct {
	sync.RWMutex
	acct        accounts.Account
	BlockHeader *objs.BlockHeader
	rawBclaims  []byte
	rawSigGroup []byte
}

// asserting that SnapshotTask struct implements interface interfaces.Task
var _ interfaces.Task = &SnapshotTask{}

func NewSnapshotTask(account accounts.Account) *SnapshotTask {
	return &SnapshotTask{
		acct: account,
	}
}

func (t *SnapshotTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum, _ interface{}) error {

	if t.BlockHeader == nil {
		return errors.New("BlockHeader must be assigned before initializing")
	}

	rawBClaims, err := t.BlockHeader.BClaims.MarshalBinary()
	if err != nil {
		logger.Errorf("Unable to marshal block header: %v", err)
		return err
	}

	t.Lock()
	defer t.Unlock()

	t.rawBclaims = rawBClaims
	t.rawSigGroup = t.BlockHeader.SigGroup

	return nil
}
func (t *SnapshotTask) DoWork(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

func (t *SnapshotTask) DoRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

func (t *SnapshotTask) doTask(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {

	t.RLock()
	defer t.RUnlock()

	for {
		dangerousRand.Seed(time.Now().UnixNano())
		n := dangerousRand.Intn(60) // n will be between 0 and 60
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(n) * time.Second):
		}
		if !t.ShouldRetry(ctx, logger, eth) {
			logger.Debug("Should not retry snapshots!")
			// no need for snapshot
			return nil
		}

		// do the actual snapshot
		err := func() error {
			txnOpts, err := eth.GetTransactionOpts(ctx, t.acct)
			if err != nil {
				logger.Debugf("Failed to generate transaction options: %v", err)
				return nil
			}

			txn, err := eth.Contracts().Snapshots().Snapshot(txnOpts, t.rawSigGroup, t.rawBclaims)
			if err != nil {
				logger.Debugf("Snapshot failed: %v", err)
				return nil
			} else {
				rcpt, err := eth.Queue().QueueAndWait(ctx, txn)
				if err != nil {
					logger.Debugf("Snapshot failed to retrive receipt: %v", err)
					return nil
				}

				if rcpt.Status != 1 {
					logger.Debugf("Snapshot receipt status != 1")
					return context.DeadlineExceeded
				}

				logger.Info("Snapshot succeeded")
			}

			return nil
		}()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				select {
				case <-ctx.Done():
					return err
				default:
				}
				logger.Debugf("Retrying snapshot after failed tx")
				continue
			}
		}

		logger.Debugf("Waiting for finality delay")
		// check/wait for finality delay
		err = func() error {
			subctx, cf := context.WithTimeout(ctx, 5*time.Second)
			defer cf()
			initialHeight, err := eth.GetCurrentHeight(subctx)
			if err != nil {
				logger.Debugf("Error to get current eth height")
				return err
			}

			currentHeight := initialHeight
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 5):
				}

				err := func() error {
					subctx, cf := context.WithTimeout(ctx, 5*time.Second)
					defer cf()
					testHeight, err := eth.GetCurrentHeight(subctx)
					if err != nil {
						logger.Debugf("Error to get test eth height")
						return err
					}

					if testHeight > initialHeight+eth.GetFinalityDelay() {
						logger.Debugf("testHeight: %v", testHeight)
						logger.Debugf("initialHeight: %v", initialHeight)
						logger.Debugf("getFinalityDelay: %v", eth.GetFinalityDelay())
						return nil
					}

					if testHeight > currentHeight {
						if !t.ShouldRetry(ctx, logger, eth) {
							// no need for snapshot
							currentHeight = testHeight
							if currentHeight >= initialHeight+eth.GetFinalityDelay() {
								logger.Debugf("Finished waiting finality delay!")
								// todo: figure how to get the doTask() func to return nil
								return nil
							}
						}
					}
					logger.Debugf("Finished waiting finality delay 2!")
					return nil
				}()

				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						select {
						case <-ctx.Done():
							return err
						default:
						}
						logger.Debugf("Retrying snapshot after finality delay")
						continue
					}
				}
			}
		}()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				select {
				case <-ctx.Done():
					return err
				default:
				}
				logger.Debugf("Retrying snapshot after everything")
				continue
			}
		}

	}
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

	height, err := eth.Contracts().Snapshots().GetAliceNetHeightFromSnapshot(opts, epoch)
	if err != nil {
		logger.Errorf("Failed to determine height: %v", err)
		return true
	}

	// This means the block height we want to snapshot is older than (or same as) what's already been snapshotted
	if t.BlockHeader.BClaims.Height != 0 && t.BlockHeader.BClaims.Height < uint32(height.Uint64()) {
		return false
	}

	return true
}

func (*SnapshotTask) DoDone(logger *logrus.Entry) {
}

func (*SnapshotTask) GetExecutionData() interface{} {
	return nil
}
