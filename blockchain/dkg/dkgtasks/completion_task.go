package dkgtasks

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/MadBase/MadNet/blockchain/dkg"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

// CompletionTask contains required state for safely performing a registration
type CompletionTask struct {
	Start          uint64
	End            uint64
	State          *objects.DkgState
	Success        bool
	StartBlockHash common.Hash
}

// asserting that CompletionTask struct implements interface interfaces.Task
var _ interfaces.Task = &CompletionTask{}

// NewCompletionTask creates a background task that attempts to call Complete on ethdkg
func NewCompletionTask(state *objects.DkgState, start uint64, end uint64) *CompletionTask {
	return &CompletionTask{
		Start:   start,
		End:     end,
		State:   state,
		Success: false,
	}
}

// Initialize prepares for work to be done in the Completion phase
func (t *CompletionTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum, state interface{}) error {

	logger.Info("CompletionTask Initialize()...")

	dkgData, ok := state.(objects.ETHDKGTaskData)
	if !ok {
		return objects.ErrCanNotContinue
	}

	t.State = dkgData.State

	t.State.Lock()
	defer t.State.Unlock()

	if t.State.Phase != objects.DisputeGPKJSubmission {
		return fmt.Errorf("%w because it's not in DisputeGPKJSubmission phase", objects.ErrCanNotContinue)
	}

	// setup leader election
	block, err := eth.GetGethClient().BlockByNumber(ctx, big.NewInt(int64(t.Start)))
	if err != nil {
		return fmt.Errorf("CompletionTask.Initialize(): error getting block by number: %v", err)
	}

	logger.Infof("block hash: %v\n", block.Hash())
	t.StartBlockHash.SetBytes(block.Hash().Bytes())

	return nil
}

// DoWork is the first attempt
func (t *CompletionTask) DoWork(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

// DoRetry is all subsequent attempts
func (t *CompletionTask) DoRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

func (t *CompletionTask) doTask(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {

	t.State.Lock()
	defer t.State.Unlock()

	logger.Info("CompletionTask doTask()")

	if t.isTaskCompleted(ctx, eth) {
		t.Success = true
		return nil
	}

	// submit if I'm a leader for this task
	if !t.AmILeading(ctx, eth, logger) {
		return errors.New("not leading Completion yet")
	}

	// Setup
	c := eth.Contracts()
	txnOpts, err := eth.GetTransactionOpts(ctx, t.State.Account)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "getting txn opts failed: %v", err)
	}

	// Register
	txn, err := c.Ethdkg().Complete(txnOpts)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "completion failed: %v", err)
	}

	logger.Info("CompletionTask sent completed call")

	//TODO: add retry logic, add timeout

	// Waiting for receipt
	receipt, err := eth.Queue().QueueAndWait(ctx, txn)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "waiting for receipt failed: %v", err)
	}
	if receipt == nil {
		return dkg.LogReturnErrorf(logger, "missing completion receipt")
	}

	// Check receipt to confirm we were successful
	if receipt.Status != uint64(1) {
		return dkg.LogReturnErrorf(logger, "completion status (%v) indicates failure: %v", receipt.Status, receipt.Logs)
	}

	logger.Info("CompletionTask complete!")

	t.Success = true

	return nil
}

// ShouldRetry checks if it makes sense to try again
// Predicates:
// -- we haven't passed the last block
// -- the registration open hasn't moved, i.e. ETHDKG has not restarted
func (t *CompletionTask) ShouldRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) bool {

	t.State.Lock()
	defer t.State.Unlock()

	logger.Info("CompletionTask ShouldRetry()")

	generalRetry := GeneralTaskShouldRetry(ctx, logger, eth, t.Start, t.End)
	if !generalRetry {
		return false
	}

	if t.isTaskCompleted(ctx, eth) {
		logger.WithFields(logrus.Fields{
			"t.State.Phase":      t.State.Phase,
			"t.State.PhaseStart": t.State.PhaseStart,
		}).Info("CompletionTask ShouldRetry - will not retry because it's done")
		return false
	}

	logger.Info("CompletionTask ShouldRetry() will retry")

	return true
}

// DoDone creates a log entry saying task is complete
func (t *CompletionTask) DoDone(logger *logrus.Entry) {
	t.State.Lock()
	defer t.State.Unlock()

	logger.WithField("Success", t.Success).Infof("CompletionTask done")
}

func (t *CompletionTask) isTaskCompleted(ctx context.Context, eth interfaces.Ethereum) bool {
	c := eth.Contracts()
	phase, err := c.Ethdkg().GetETHDKGPhase(eth.GetCallOpts(ctx, t.State.Account))
	if err != nil {
		return false
	}

	return phase == uint8(objects.Completion)
}

func (t *CompletionTask) AmILeading(ctx context.Context, eth interfaces.Ethereum, logger *logrus.Entry) bool {
	// check if I'm a leader for this task
	currentHeight, err := eth.GetCurrentHeight(ctx)
	if err != nil {
		return false
	}

	blocksSinceDesperation := int(currentHeight) - int(t.Start) - constants.ETHDKGDesperationDelay
	amILeading := dkg.AmILeading(t.State.NumberOfValidators, t.State.Index-1, blocksSinceDesperation, t.StartBlockHash.Bytes(), logger)

	logger.WithFields(logrus.Fields{
		"currentHeight":                    currentHeight,
		"t.Start":                          t.Start,
		"constants.ETHDKGDesperationDelay": constants.ETHDKGDesperationDelay,
		"blocksSinceDesperation":           blocksSinceDesperation,
		"amILeading":                       amILeading,
	}).Infof("dkg.AmILeading")

	return amILeading
}
