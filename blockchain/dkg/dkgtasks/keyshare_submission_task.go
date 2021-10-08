package dkgtasks

import (
	"context"
	"fmt"

	"github.com/MadBase/MadNet/blockchain/dkg"
	"github.com/MadBase/MadNet/blockchain/dkg/math"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/sirupsen/logrus"
)

// DisputeTask stores the data required to dispute shares
type KeyshareSubmissionTask struct {
	OriginalRegistrationEnd uint64
	State                   *objects.DkgState
	Success                 bool
}

// NewKeyshareSubmissionTask creates a new task
func NewKeyshareSubmissionTask(state *objects.DkgState) *KeyshareSubmissionTask {
	return &KeyshareSubmissionTask{
		OriginalRegistrationEnd: state.RegistrationEnd, // If these quit being equal, this task should be abandoned
		State:                   state,
	}
}

// This is not exported and does not lock so can only be called from within task. Return value indicates whether task has been initialized.
func (t *KeyshareSubmissionTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {

	t.State.Lock()
	defer t.State.Unlock()

	if !t.State.Dispute {
		return fmt.Errorf("%w because dispute phase not successful", objects.ErrCanNotContinue)
	}

	// Generate the key shares
	g1KeyShare, g1Proof, g2KeyShare, err := math.GenerateKeyShare(t.State.SecretValue)
	if err != nil {
		return err
	}

	// t.State.KeyShareG1s[state.Account.Address]
	me := t.State.Account.Address

	logger.Infof("generating key shares for %v from %v", me.Hex(), t.State.SecretValue.String())

	t.State.KeyShareG1s[me] = g1KeyShare
	t.State.KeyShareG1CorrectnessProofs[me] = g1Proof
	t.State.KeyShareG2s[me] = g2KeyShare

	return nil
}

// DoWork is the first attempt at registering with ethdkg
func (t *KeyshareSubmissionTask) DoWork(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	logger.Info("DoWork() ...")
	return t.doTask(ctx, logger, eth)
}

// DoRetry is all subsequent attempts at registering with ethdkg
func (t *KeyshareSubmissionTask) DoRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	logger.Info("DoRetry() ...")
	return t.doTask(ctx, logger, eth)
}

func (t *KeyshareSubmissionTask) doTask(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) error {
	t.State.Lock()
	defer t.State.Unlock()

	// Setup
	me := t.State.Account

	txnOpts, err := eth.GetTransactionOpts(ctx, me)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "getting txn opts failed: %v", err)
	}

	// Submit Keyshares
	logger.Infof("submitting key shares: %v %p %p %p",
		me.Address,
		t.State.KeyShareG1s[me.Address],
		t.State.KeyShareG1CorrectnessProofs[me.Address],
		t.State.KeyShareG2s[me.Address])
	txn, err := eth.Contracts().Ethdkg().SubmitKeyShare(txnOpts, me.Address,
		t.State.KeyShareG1s[me.Address],
		t.State.KeyShareG1CorrectnessProofs[me.Address],
		t.State.KeyShareG2s[me.Address])
	if err != nil {
		return dkg.LogReturnErrorf(logger, "submitting keyshare failed: %v", err)
	}

	// Waiting for receipt
	receipt, err := eth.Queue().QueueAndWait(ctx, txn)
	if err != nil {
		return dkg.LogReturnErrorf(logger, "waiting for receipt failed: %v", err)
	}
	if receipt == nil {
		return dkg.LogReturnErrorf(logger, "missing submit keyshare receipt")
	}

	// Check receipt to confirm we were successful
	if receipt.Status != uint64(1) {
		dkg.LogReturnErrorf(logger, "submit keyshare status (%v) indicates failure: %v", receipt.Status, receipt.Logs)
	}

	t.Success = true

	return nil
}

// ShouldRetry checks if it makes sense to try again
// Predicates:
// -- we haven't passed the last block
// -- the registration open hasn't moved, i.e. ETHDKG has not restarted
func (t *KeyshareSubmissionTask) ShouldRetry(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum) bool {
	t.State.Lock()
	defer t.State.Unlock()

	state := t.State

	// This wraps the retry logic for the general case
	generalRetry := GeneralTaskShouldRetry(ctx, state.Account, logger, eth,
		state.TransportPublicKey, t.OriginalRegistrationEnd, state.KeyShareSubmissionEnd)

	// If it's generally good to retry, let's try to be more specific
	if generalRetry {
		me := state.Account
		callOpts := eth.GetCallOpts(ctx, me)
		status, err := CheckKeyShare(ctx, eth.Contracts().Ethdkg(), logger, callOpts, me.Address, state.KeyShareG1s[me.Address])
		if err != nil {
			return true
		}
		logger.Infof("Key shared status: %v", status)

		// If we have already shared a key, there is no reason to retry. Regardless of whether it's right or wrong.
		if status == KeyShared || status == BadKeyShared {
			return false
		}
	}

	return generalRetry
}

// DoDone creates a log entry saying task is complete
func (t *KeyshareSubmissionTask) DoDone(logger *logrus.Entry) {
	logger.Infof("done")

	t.State.Lock()
	defer t.State.Unlock()

	t.State.KeyShareSubmission = t.Success
}
