package dkgtasks

import (
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/ethereum/go-ethereum/common"
)

type ExecutionData struct {
	Start          uint64
	End            uint64
	State          *objects.DkgState
	Success        bool
	StartBlockHash common.Hash
	TxOpts         *objects.TxOpts
}

func (d *ExecutionData) Clear() {
	d.TxOpts = &objects.TxOpts{
		TxHashes: make([]common.Hash, 0),
	}
}

func NewExecutionData(state *objects.DkgState, start uint64, end uint64) *ExecutionData {
	return &ExecutionData{
		State:   state,
		Start:   start,
		End:     end,
		Success: false,
		TxOpts:  &objects.TxOpts{TxHashes: make([]common.Hash, 0)},
	}
}
