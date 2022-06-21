package state

import (
	"encoding/json"
	"fmt"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/layer1/executor/tasks"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/accounts"
)

// asserting that SnapshotState struct implements interface tasks.ITaskState
var _ tasks.TaskState = &SnapshotState{}

type SnapshotState struct {
	Account     accounts.Account
	RawBClaims  []byte
	RawSigGroup []byte
	BlockHeader *objs.BlockHeader
}

func (state *SnapshotState) PersistState(txn *badger.Txn) error {
	rawData, err := json.Marshal(state)
	if err != nil {
		return err
	}

	key := dbprefix.PrefixSnapshotState()
	if err = utils.SetValue(txn, key, rawData); err != nil {
		return err
	}

	return nil
}

func (state *SnapshotState) LoadState(txn *badger.Txn) error {
	key := dbprefix.PrefixSnapshotState()
	rawData, err := utils.GetValue(txn, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawData, state)
	if err != nil {
		return err
	}

	return nil

}

func GetSnapshotState(monDB *db.Database) (*SnapshotState, error) {
	snapshotState := &SnapshotState{}
	err := monDB.View(func(txn *badger.Txn) error {
		return snapshotState.LoadState(txn)
	})
	if err != nil {
		return nil, err
	}
	return snapshotState, nil
}

func SaveSnapshotState(monDB *db.Database, snapshotState *SnapshotState) error {
	err := monDB.Update(func(txn *badger.Txn) error {
		return snapshotState.PersistState(txn)
	})
	if err != nil {
		return err
	}
	if err = monDB.Sync(); err != nil {
		return fmt.Errorf("Failed to set sync of snapshotState: %v", err)
	}
	return nil
}
