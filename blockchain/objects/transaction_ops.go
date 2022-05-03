package objects

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type TxOpts struct {
	TxHashes     []common.Hash
	Nonce        *big.Int
	GasFeeCap    *big.Int
	GasTipCap    *big.Int
	MinedInBlock uint64
}

func (t *TxOpts) GetHexTxsHashes() string {
	var hashes strings.Builder
	for _, txHash := range t.TxHashes {
		hashes.WriteString(txHash.Hex())
		hashes.WriteString(" ")
	}
	return hashes.String()
}
