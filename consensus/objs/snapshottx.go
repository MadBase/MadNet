package objs

import "github.com/ethereum/go-ethereum/core/types"

type SnapshotTxState byte

const (
	SnapshotTxCreated SnapshotTxState = iota + 1
	SnapshotTxSubmitted
	SnapshotTxVerified
)

type CachedSnapshotTx struct {
	State     SnapshotTxState
	Tx        *types.Transaction
	MadEpoch  uint32
	EthHeight uint32
}
