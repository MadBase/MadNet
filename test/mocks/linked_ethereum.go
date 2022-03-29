package mocks

import (
	"math/big"

	"github.com/MadBase/MadNet/blockchain/interfaces"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

type EthMockset struct {
	Ethereum   *MockEthereum
	GethClient *MockGethClient
	Queue      *MockTxnQueue
	Contracts  *MockContracts

	ETHDKG           *MockIETHDKG
	Governance       *MockIGovernance
	MadByte          *MockIMadByte
	MadToken         *MockIMadToken
	PublicStaking    *MockIPublicStaking
	Snapshots        *MockISnapshots
	ValidatorPool    *MockIValidatorPool
	ValidatorStaking *MockIValidatorStaking
}

var _ interfaces.Ethereum = (*MockEthereum)(nil)

func NewMockLinkedEthereum() (*MockEthereum, *EthMockset) {
	eth := NewMockEthereum()
	eth.GetCurrentHeightFunc.SetDefaultReturn(1024, nil)
	eth.GetTransactionOptsFunc.SetDefaultReturn(&bind.TransactOpts{}, nil)
	eth.SignTxFunc.SetDefaultHook(func(a common.Address, tx *types.DynamicFeeTx) (*types.Transaction, error) {
		return types.NewTx(tx), nil
	})

	geth := NewMockLinkedGethClient()
	eth.GetGethClientFunc.SetDefaultReturn(geth)

	queue := NewMockLinkedQueue()
	eth.QueueFunc.SetDefaultReturn(queue)

	contracts := NewMockContracts()
	eth.ContractsFunc.SetDefaultReturn(contracts)

	ethdkg := NewMockIETHDKG()
	contracts.EthdkgFunc.SetDefaultReturn(ethdkg)

	governance := NewMockIGovernance()
	contracts.GovernanceFunc.SetDefaultReturn(governance)

	madbyte := NewMockIMadByte()
	contracts.MadByteFunc.SetDefaultReturn(madbyte)

	madtoken := NewMockIMadToken()
	contracts.MadTokenFunc.SetDefaultReturn(madtoken)

	publicstaking := NewMockIPublicStaking()
	contracts.PublicStakingFunc.SetDefaultReturn(publicstaking)

	snapshots := NewMockLinkedSnapshots()
	contracts.SnapshotsFunc.SetDefaultReturn(snapshots)

	validatorpool := NewMockIValidatorPool()
	contracts.ValidatorPoolFunc.SetDefaultReturn(validatorpool)

	validatorstaking := NewMockIValidatorStaking()
	contracts.ValidatorStakingFunc.SetDefaultReturn(validatorstaking)

	return eth, &EthMockset{
		Ethereum:   eth,
		GethClient: geth,
		Queue:      queue,
		Contracts:  contracts,

		ETHDKG:           ethdkg,
		Governance:       governance,
		MadByte:          madbyte,
		MadToken:         madtoken,
		PublicStaking:    publicstaking,
		Snapshots:        snapshots,
		ValidatorPool:    validatorpool,
		ValidatorStaking: validatorstaking,
	}
}

func NewMockLinkedSnapshots() *MockISnapshots {
	m := NewMockISnapshots()
	m.SnapshotFunc.SetDefaultHook(func(*bind.TransactOpts, []byte, []byte) (*types.Transaction, error) { return NewMockSnapshotTx(), nil })
	return m
}

func NewMockLinkedQueue() *MockTxnQueue {
	queue := NewMockTxnQueue()
	queue.WaitTransactionFunc.SetDefaultReturn(&types.Receipt{Status: 1}, nil)
	return queue
}

func NewMockLinkedGethClient() *MockGethClient {
	geth := NewMockGethClient()
	geth.SuggestGasTipCapFunc.SetDefaultReturn(big.NewInt(15000), nil)
	geth.SuggestGasPriceFunc.SetDefaultReturn(big.NewInt(1000), nil)
	return geth
}
