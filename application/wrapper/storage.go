package wrapper

import (
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/MadBase/MadNet/application/objs/uint256"
)

func bigFromBase10(s string) *uint256.Uint256 {
	n, _ := new(big.Int).SetString(s, 10)
	nn, _ := new(uint256.Uint256).FromBigInt(n)
	return nn
}

const maxBytes = 3000000

type feeResponse struct {
	MinTxFee      string `json:"MinTxFee,omitempty"`
	ValueStoreFee string `json:"ValueStoreFee,omitempty"`
	DataStoreFee  string `json:"DataStoreFee,omitempty"`
	AtomicSwapFee string `json:"AtomicSwapFee,omitempty"`
}

// Storage wraps the dynamics.StorageGetter interface to make
// it easier to interact within application logic
type Storage struct {
	overRide feeResponse
}

// NewStorage creates a new storage struct which wraps
// the StorageGetter interface
func NewStorage(updated string) (*Storage, error) {
	f := new(feeResponse)
	err := json.Unmarshal([]byte(updated), f)
	if err != nil {
		return nil, err
	}
	storage := &Storage{*f}
	return storage, nil
}

// GetMaxBytes returns MaxBytes
func (s *Storage) GetMaxBytes() uint32 {
	return maxBytes
}

// GetAtomicSwapFee returns the fee for AtomicSwap
func (s *Storage) GetAtomicSwapFee() (*uint256.Uint256, error) {
	if len(s.overRide.AtomicSwapFee) > 0 {
		u, err := hex.DecodeString(s.overRide.AtomicSwapFee)
		if err != nil {
			return nil, err
		}
		return uint256.Uint256FromBytes(u)
	}
	return bigFromBase10("2"), nil
}

// GetDataStoreEpochFee returns the per-epoch fee of DataStore
func (s *Storage) GetDataStoreEpochFee() (*uint256.Uint256, error) {
	if len(s.overRide.DataStoreFee) > 0 {
		u, err := hex.DecodeString(s.overRide.DataStoreFee)
		if err != nil {
			return nil, err
		}
		return uint256.Uint256FromBytes(u)
	}
	return bigFromBase10("3"), nil
}

// GetValueStoreFee returns the fee of ValueStore
func (s *Storage) GetValueStoreFee() (*uint256.Uint256, error) {
	if len(s.overRide.ValueStoreFee) > 0 {
		u, err := hex.DecodeString(s.overRide.ValueStoreFee)
		if err != nil {
			return nil, err
		}
		return uint256.Uint256FromBytes(u)
	}

	return bigFromBase10("1"), nil
}

// GetMinTxFee returns the minimum TxFee
func (s *Storage) GetMinTxFee() (*uint256.Uint256, error) {
	if len(s.overRide.MinTxFee) > 0 {
		u, err := hex.DecodeString(s.overRide.MinTxFee)
		if err != nil {
			return nil, err
		}
		return uint256.Uint256FromBytes(u)
	}

	return bigFromBase10("4"), nil
}
