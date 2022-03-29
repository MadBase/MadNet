package dkgtasks

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/MadBase/MadNet/blockchain/dkg"
	"github.com/MadBase/MadNet/blockchain/dkg/math"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

// MPKSubmissionTask stores the data required to submit the mpk
type MPKSubmissionTask struct {
	Start          uint64
	End            uint64
	State          *objects.DkgState
	Success        bool
	StartBlockHash common.Hash
}

// asserting that MPKSubmissionTask struct implements interface interfaces.Task
var _ interfaces.Task = &MPKSubmissionTask{}

// NewMPKSubmissionTask creates a new task
func NewMPKSubmissionTask(state *objects.DkgState, start uint64, end uint64) *MPKSubmissionTask {
	return &MPKSubmissionTask{
		Start:   start,
		End:     end,
		State:   state,
		Success: false,
	}
}

// Initialize prepares for work to be done in MPKSubmission phase.
// Here we load all key shares and construct the master public key
// to submit in DoWork.
func (t *MPKSubmissionTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum, state interface{}) error {
	dkgData, ok := state.(objects.ETHDKGTaskData)
	if !ok {
		return objects.ErrCanNotContinue
	}

	t.State = dkgData.State

	t.State.Lock()
	defer t.State.Unlock()

	logger.Info("MPKSubmissionTask Initialize()...")

	if t.State.Phase != objects.MPKSubmission {
		return fmt.Errorf("%w because it's not in MPKSubmission phase", objects.ErrCanNotContinue)
	}

	// setup leader election
	block, err := eth.GetGethClient().BlockByNumber(ctx, big.NewInt(int64(t.Start)))
	if err != nil {
		return fmt.Errorf("MPKSubmissionTask.Initialize(): error getting block by number: %v", err)
	}

	logger.Infof("block hash: %v\n", block.Hash())
	t.StartBlockHash.SetBytes(block.Hash().Bytes())

	// prepare MPK
	g1KeyShares := make([][2]*big.Int, t.State.NumberOfValidators)
	g2KeyShares := make([][4]*big.Int, t.State.NumberOfValidators)

	var participantsList = t.State.GetSortedParticipants()
	validMPK := true
	for idx, participant := range participantsList {
		// Bringing these in from state but could directly query contract
		g1KeyShares[idx] = t.State.Participants[participant.Address].KeyShareG1s
		g2KeyShares[idx] = t.State.Participants[participant.Address].KeyShareG2s

		logger.Debugf("INIT idx:%v pidx:%v address:%v g1:%v g2:%v", idx, participant.Index, participant.Address.Hex(), g1KeyShares[idx], g2KeyShares[idx])

		for i := range g1KeyShares[idx] {
			if g1KeyShares[idx][i] == nil {
				logger.Errorf("Missing g1Keyshare[%v][%v] for %v.", idx, i, participant.Address.Hex())
				validMPK = false
			}
		}

		for i := range g2KeyShares[idx] {
			if g2KeyShares[idx][i] == nil {
				logger.Errorf("Missing g2Keyshare[%v][%v] for %v.", idx, i, participant.Address.Hex())
				validMPK = false
			}
		}
	}

	logger.Infof("# Participants: %v\n", len(t.State.Participants))

	mpk, err := math.GenerateMasterPublicKey(g1KeyShares, g2KeyShares)
	if err != nil && validMPK {
		return dkg.LogReturnErrorf(logger, "Failed to generate master public key:%v", err)
	}

	if !validMPK {
		mpk = [4]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	}

	// Master public key is all we generate here so save it
	t.State.MasterPublicKey = mpk

	return nil
}

// DoWork is the first attempt at submitting the mpk
func (t *MPKSubmissionTask) DoWork(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

// DoRetry is all subsequent attempts at submitting the mpk
func (t *MPKSubmissionTask) DoRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	return t.doTask(ctx, logger, eth)
}

func (t *MPKSubmissionTask) doTask(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	t.State.Lock()
	defer t.State.Unlock()

	logger.Info("MPKSubmissionTask doTask()")

	if !t.shouldSubmitMPK(ctx, eth, logger) {
		t.Success = true
		return nil
	}

	// submit if I'm a leader for this task
	if !t.AmILeading(ctx, eth, logger) {
		return errors.New("not leading MPK submission yet")
	}

	// Setup
	txnOpts, err := eth.GetTransactionOpts(ctx, t.State.Account)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "getting txn opts failed: %v", err)
	}

	// Submit MPK
	logger.Infof("submitting master public key:%v", t.State.MasterPublicKey)
	txn, err := eth.Contracts().Ethdkg().SubmitMasterPublicKey(txnOpts, t.State.MasterPublicKey)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "submitting master public key failed: %v", err)
	}

	//TODO: add retry logic, add timeout

	// Waiting for receipt
	receipt, err := eth.Queue().QueueAndWait(ctx, txn)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "waiting for receipt failed: %v", err)
	}
	if receipt == nil {
		return dkg.LogReturnErrorf(logger, "missing receipt")
	}

	// Check receipt to confirm we were successful
	if receipt.Status != uint64(1) {
		return dkg.LogReturnErrorf(logger, "master public key (%v) indicates failure: %v", receipt.Status, receipt.Logs)
	}
	t.Success = true

	return nil
}

// ShouldRetry checks if it makes sense to try again
// Predicates:
// -- we haven't passed the last block
func (t *MPKSubmissionTask) ShouldRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) bool {
	t.State.Lock()
	defer t.State.Unlock()

	logger.Info("MPKSubmissionTask ShouldRetry()")

	generalRetry := GeneralTaskShouldRetry(ctx, logger, eth, t.Start, t.End)
	if !generalRetry {
		return false
	}

	if t.State.Phase != objects.MPKSubmission {
		return false
	}

	return t.shouldSubmitMPK(ctx, eth, logger)
}

// DoDone creates a log entry saying task is complete
func (t *MPKSubmissionTask) DoDone(logger *logrus.Entry) {
	t.State.Lock()
	defer t.State.Unlock()

	logger.WithField("Success", t.Success).Infof("MPKSubmissionTask done")
}

func (t *MPKSubmissionTask) shouldSubmitMPK(ctx context.Context, eth interfaces.Ethereum, logger *logrus.Entry) bool {
	if t.State.MasterPublicKey[0].Cmp(big.NewInt(0)) == 0 &&
		t.State.MasterPublicKey[1].Cmp(big.NewInt(0)) == 0 &&
		t.State.MasterPublicKey[2].Cmp(big.NewInt(0)) == 0 &&
		t.State.MasterPublicKey[3].Cmp(big.NewInt(0)) == 0 {
		return false
	}

	isMPKSet, err := eth.Contracts().Ethdkg().IsMasterPublicKeySet(eth.GetCallOpts(ctx, t.State.Account))
	if err == nil && isMPKSet {
		return false
	}

	return !isMPKSet
}

func (t *MPKSubmissionTask) AmILeading(ctx context.Context, eth interfaces.Ethereum, logger *logrus.Entry) bool {
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
