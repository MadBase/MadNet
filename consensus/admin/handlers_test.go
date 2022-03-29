package admin

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/test/mocks"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logging.GetLogger(constants.LoggerConsensus).SetLevel(logrus.ErrorLevel)

	code := m.Run()
	os.Exit(code)
}

func initialize(VAddr byte, validatorAddrs []byte, initialHeight int) (*Handlers, *db.Database) {
	ah := &Handlers{}
	db := mocks.NewMockDb()
	db.Update(func(txn *badger.Txn) error {
		vs := make([]*objs.Validator, len(validatorAddrs))
		for i, addr := range validatorAddrs {
			vs[i] = &objs.Validator{VAddr: []byte{addr}}
		}
		err := db.SetValidatorSet(txn, &objs.ValidatorSet{Validators: vs, NotBefore: 1})
		if err != nil {
			panic(err)
		}

		if initialHeight >= 0 {
			os := mocks.NewMockOwnState()
			os.VAddr = []byte{VAddr}
			os.SyncToBH.BClaims.Height = uint32(initialHeight)
			err = db.SetOwnState(txn, os)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	ah.Init(1, db, []byte{}, nil, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, mocks.NewMockStorage(), nil)

	return ah, db
}

// TestFirstSnapshot tests submission of the very first snapshot
// In this case, the ethereum height of the previous snapshot is not considered
func TestFirstSnapshot(t *testing.T) {
	var myTests = []struct {
		name           string
		addr           byte
		addrs          []byte
		madHeight      int
		ethHeight      int
		shouldSnapshot bool
	}{
		// Snapshots should only be taken at an epoch boundary
		{name: "EpochNotReached", shouldSnapshot: false,
			addr: 4, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: int(constants.EpochLength) - 1,
			ethHeight: 1e6},
		{name: "EpochReached", shouldSnapshot: true,
			addr: 4, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: int(constants.EpochLength),
			ethHeight: 1e6},
	}

	for _, test := range myTests {
		t.Run(test.name, func(t *testing.T) {

			ah, db := initialize(test.addr, test.addrs, 1)

			called := false
			ah.RegisterSnapshotCallbacks(
				func(bh *objs.BlockHeader) error { called = true; return nil },
				nil,
				nil,
				nil,
			)

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			ah.UpdateEthHeight(uint32(test.ethHeight))

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			db.Update(func(txn *badger.Txn) error {
				bh := mocks.NewMockBlockHeader()
				bh.BClaims.Height = uint32(test.madHeight)
				return db.SetBroadcastBlockHeader(txn, bh)
			})

			time.Sleep(time.Millisecond)
			if test.shouldSnapshot {
				require.Truef(t, called, "expected snapshot to have been taken")
			} else {
				require.Falsef(t, called, "expected no snapshot to have been taken")
			}
		})
	}
}

// TestFirstSnapshot tests submission a snapshot later than the first one
// In this case, the ethereum height of the previous snapshot defines desperation, e.g.: how badly we want a snapshot, and therefore how many validators may make one
func TestSubsequentSnapshot(t *testing.T) {
	var myTests = []struct {
		name            string
		addr            byte
		addrs           []byte
		madHeight       int
		ethHeight       int
		prevSSEthHeight int
		shouldSnapshot  bool
	}{
		// Snapshots should only be taken at an epoch boundary
		{name: "EpochNotReached", shouldSnapshot: false,
			addr: 12, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3*int(constants.EpochLength) - 1,
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6},
		{name: "EpochReached", shouldSnapshot: true,
			addr: 12, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6},

		// Range of validators allowed to snapshot should only start expanding after DesperationDelay is reached
		{name: "DesperationDelayNotReached", shouldSnapshot: false,
			addr: 0, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay + 1},
		{name: "DesperationDelayReached", shouldSnapshot: true,
			addr: 0, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay},

		// Rate of expansion of the range validators allowed to snapshot should follow SnapshotDesperationFactor
		{name: "DesperationFactorNotReached", shouldSnapshot: false,
			addr: 1, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay - constants.SnapshotDesperationFactor + 1},
		{name: "DesperationFactorReached", shouldSnapshot: true,
			addr: 1, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay - constants.SnapshotDesperationFactor},
		{name: "DesperationFactorNotReachedAgain", shouldSnapshot: false,
			addr: 2, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay - constants.SnapshotDesperationFactor - constants.SnapshotDesperationFactor/2 + 1},
		{name: "DesperationFactorReachedAgain", shouldSnapshot: true,
			addr: 2, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 1e6 - constants.SnapshotDesperationDelay - constants.SnapshotDesperationFactor - constants.SnapshotDesperationFactor/2},

		// Range of validators allowed to snapshot should never be more than 1/3, regardless of desperation
		{name: "DesperationLimit", shouldSnapshot: false,
			addr: 4, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6, prevSSEthHeight: 0},
	}

	for _, test := range myTests {
		t.Run(test.name, func(t *testing.T) {

			ah, db := initialize(test.addr, test.addrs, 1)

			called := false
			ah.RegisterSnapshotCallbacks(
				func(bh *objs.BlockHeader) error { called = true; return nil },
				nil,
				func(epoch *big.Int) (*big.Int, error) { return big.NewInt(int64(test.prevSSEthHeight)), nil },
				nil,
			)

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			ah.UpdateEthHeight(uint32(test.ethHeight))

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			db.Update(func(txn *badger.Txn) error {
				bh := mocks.NewMockBlockHeader()
				bh.BClaims.Height = uint32(test.madHeight)
				return db.SetBroadcastBlockHeader(txn, bh)
			})

			time.Sleep(time.Millisecond)
			if test.shouldSnapshot {
				require.Truef(t, called, "expected snapshot to have been taken")
			} else {
				require.Falsef(t, called, "expected no snapshot to have been taken")
			}
		})
	}
}

// TestResumeSnapshot asserts that previously used snapshot txs are properly resumed if the tx was not successful
// Only previous snapshot txs of the current epoch are resumed
func TestResumeSnapshot(t *testing.T) {
	var myTests = []struct {
		name            string
		addr            byte
		addrs           []byte
		madHeight       int
		ethHeight       int
		prevSSEthHeight int
		cachedTx        *objs.CachedSnapshotTx
		refreshedTx     *types.Transaction
		shouldResume    bool
		shouldRefresh   bool
	}{
		// Snapshots that were only created, but not submitted to the eth network, should just be submitted in their existing state
		{name: "CreatedTxFromPrevEpoch", shouldResume: false,
			addr: 12, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6,
			cachedTx: &objs.CachedSnapshotTx{State: objs.SnapshotTxCreated, Tx: mocks.NewMockSnapshotTx(), MadEpoch: 2}},
		{name: "CreatedTx", shouldResume: true,
			addr: 8, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6,
			cachedTx: &objs.CachedSnapshotTx{State: objs.SnapshotTxCreated, Tx: mocks.NewMockSnapshotTx(), MadEpoch: 3}},

		// Snapshots that were submitted to the eth network, but not verified, should be resubmitted to the eth network
		// The injected refreshTx function dictates the transaction should be refreshed or left untouched
		{name: "SubmittedTxFromPrevEpoch", shouldResume: false,
			addr: 8, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6,
			cachedTx:    &objs.CachedSnapshotTx{State: objs.SnapshotTxSubmitted, Tx: mocks.NewMockSnapshotTx(), MadEpoch: 2},
			refreshedTx: mocks.NewMockSnapshotTx2()},
		{name: "SubmittedTxAged", shouldResume: true, shouldRefresh: true,
			addr: 8, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6,
			cachedTx:    &objs.CachedSnapshotTx{State: objs.SnapshotTxSubmitted, Tx: mocks.NewMockSnapshotTx(), MadEpoch: 3},
			refreshedTx: mocks.NewMockSnapshotTx2()},
		{name: "SubmittedTxRecent", shouldResume: false,
			addr: 8, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: 3 * int(constants.EpochLength),
			ethHeight: 1e6 + constants.SnapshotDesperationDelay - 1, prevSSEthHeight: 1e6,
			cachedTx: &objs.CachedSnapshotTx{State: objs.SnapshotTxSubmitted, Tx: mocks.NewMockSnapshotTx(), MadEpoch: 3}},
	}

	for _, test := range myTests {
		t.Run(test.name, func(t *testing.T) {

			ah, db := initialize(test.addr, test.addrs, test.madHeight)
			db.Update(func(txn *badger.Txn) error {
				err := db.SetSnapshotTx(txn, test.cachedTx)
				if err != nil {
					panic(err)
				}
				return nil
			})

			var called *objs.CachedSnapshotTx
			ah.RegisterSnapshotCallbacks(
				func(bh *objs.BlockHeader) error { return nil },
				func(bh *objs.CachedSnapshotTx) error { called = bh; return nil },
				func(epoch *big.Int) (*big.Int, error) { return big.NewInt(int64(test.prevSSEthHeight)), nil },
				func(c context.Context, tx *types.Transaction, a int) (*types.Transaction, error) {
					return test.refreshedTx, nil
				},
			)

			time.Sleep(time.Millisecond)
			require.Nil(t, called, "expected no snapshot to be taken yet")

			ah.UpdateEthHeight(uint32(test.ethHeight))

			time.Sleep(time.Millisecond)
			if test.shouldResume {
				require.NotNilf(t, called, "expected snapshot to have been resumed")
				if test.shouldRefresh {
					require.Equalf(t, mocks.NewMockSnapshotTx2().Hash().Hex(), called.Tx.Hash().Hex(), "expected snapshot to have been refreshed")
				} else {
					require.Equalf(t, test.cachedTx.Tx.Hash().Hex(), called.Tx.Hash().Hex(), "expected cached snapshot to have been used")
				}
			} else {
				require.Nil(t, called, "expected no snapshot to have been resumed")
			}
		})
	}
}

// TestSnapshotDelayedInitialization ensures that the snapshot logic properly starts up,
// even when there is no Ownstate to be read from the db at the beginning
func TestSnapshotDelayedInitialization(t *testing.T) {
	var myTests = []struct {
		name           string
		addr           byte
		addrs          []byte
		madHeight      int
		ethHeight      int
		shouldSnapshot bool
	}{
		{name: "EpochNotReached", shouldSnapshot: false,
			addr: 3, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: int(constants.EpochLength) - 1,
			ethHeight: constants.SnapshotDesperationDelay - 1},

		{name: "EpochReached", shouldSnapshot: true,
			addr: 3, addrs: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			madHeight: int(constants.EpochLength),
			ethHeight: constants.SnapshotDesperationDelay - 1},
	}

	for _, test := range myTests {
		t.Run(test.name, func(t *testing.T) {

			ah, _ := initialize(0, test.addrs, -1)

			called := false
			ah.RegisterSnapshotCallbacks(
				func(bh *objs.BlockHeader) error { called = true; return nil },
				nil,
				nil,
				nil,
			)

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			ah.UpdateEthHeight(uint32(test.ethHeight))

			time.Sleep(time.Millisecond)
			require.Falsef(t, called, "expected no snapshot to be taken yet")

			os := mocks.NewMockOwnState()
			os.VAddr = []byte{test.addr}
			os.SyncToBH.BClaims.Height = uint32(test.madHeight)
			ah.osListener <- os

			time.Sleep(time.Millisecond)
			if test.shouldSnapshot {
				require.Truef(t, called, "expected snapshot to have been taken")
			} else {
				require.Falsef(t, called, "expected no snapshot to have been taken")
			}
		})
	}
}
