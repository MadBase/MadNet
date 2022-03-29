package interfaces

import (
	"context"
	"math/big"

	aobjs "github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/core/types"
)

type AdminHandler interface {
	AddPrivateKey([]byte, constants.CurveSpec) error
	AddSnapshot(header *objs.BlockHeader, safeToProceedConsensus bool) error
	AddValidatorSet(*objs.ValidatorSet) error
	RegisterSnapshotCallbacks(
		func(*objs.BlockHeader) error,
		func(*objs.CachedSnapshotTx) error,
		func(*big.Int) (*big.Int, error),
		func(context.Context, *types.Transaction, int) (*types.Transaction, error),
	)
	SetSynchronized(v bool)
	UpdateEthHeight(ethHeight uint32)
}

type DepositHandler interface {
	Add(*badger.Txn, uint32, []byte, *big.Int, *aobjs.Owner) error
}
