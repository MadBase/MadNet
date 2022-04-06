package tasks_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	mockrequire "github.com/derision-test/go-mockgen/testutil/require"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/blockchain/tasks"
	"github.com/MadBase/MadNet/test/mocks"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/stretchr/testify/assert"
)

func TestStartTask_initializeTask_HappyPath(t *testing.T) {
	eth := mocks.NewMockEthereum()
	task := mocks.NewMockTask()

	wg := sync.WaitGroup{}
	tasks.StartTask(mocks.NewMockLogger().WithField("", nil), &wg, eth, task, nil)
	wg.Wait()

	mockrequire.Called(t, task.DoWorkFunc)
	mockrequire.NotCalled(t, task.DoRetryFunc)
	mockrequire.Called(t, task.DoDoneFunc)
}

func TestStartTask_initializeTask_Error(t *testing.T) {
	eth := mocks.NewMockEthereum()
	task := mocks.NewMockTask()
	task.InitializeFunc.SetDefaultReturn(errors.New("initialize error"))

	wg := sync.WaitGroup{}
	tasks.StartTask(mocks.NewMockLogger().WithField("", nil), &wg, eth, task, nil)
	wg.Wait()

	mockrequire.NotCalled(t, task.DoWorkFunc)
	mockrequire.NotCalled(t, task.DoRetryFunc)
	mockrequire.Called(t, task.DoDoneFunc)

}

func TestStartTask_executeTask_ErrorRetry(t *testing.T) {
	eth := mocks.NewMockEthereum()
	eth.RetryCountFunc.SetDefaultReturn(10)

	task := mocks.NewMockTask()
	task.ShouldRetryFunc.SetDefaultReturn(true)
	task.DoWorkFunc.SetDefaultReturn(errors.New("DoWork_error"))
	task.DoRetryFunc.SetDefaultReturn(errors.New(tasks.NonceToLowError))

	wg := sync.WaitGroup{}
	tasks.StartTask(mocks.NewMockLogger().WithField("Task", 0), &wg, eth, task, nil)
	wg.Wait()

	mockrequire.Called(t, task.DoWorkFunc)
	mockrequire.CalledN(t, task.DoRetryFunc, 10)
	mockrequire.Called(t, task.DoDoneFunc)
}

// Happy path with mined tx present after finality delay
func TestStartTask_handleExecutedTask_FinalityDelay1(t *testing.T) {
	state := objects.NewDkgState(accounts.Account{})
	task := mocks.NewMockTaskWithExecutionData(1, 100)
	task.ExecutionData.TxOpts.TxHashes = append(task.ExecutionData.TxOpts.TxHashes, common.BigToHash(big.NewInt(123871239)))

	wg := sync.WaitGroup{}

	eth := mocks.NewMockEthereum()
	eth.GethClientMock.TransactionByHashFunc.SetDefaultReturn(&types.Transaction{}, false, nil)
	eth.GethClientMock.TransactionReceiptFunc.SetDefaultReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(1)}, nil)

	tasks.StartTask(mocks.NewMockLogger().WithField("Task", 0), &wg, eth, task, state)

	wg.Wait()

	mockrequire.Called(t, task.DoWorkFunc)
	mockrequire.Called(t, task.DoDoneFunc)
	assert.Len(t, task.ExecutionData.TxOpts.TxHashes, 1)
	assert.Equal(t, uint64(1), task.ExecutionData.TxOpts.MinedInBlock)
}

// Tx was mined, but it's not present after finality delay
func TestStartTask_handleExecutedTask_FinalityDelay2(t *testing.T) {
	minedInBlock := 9

	state := objects.NewDkgState(accounts.Account{})
	task := mocks.NewMockTaskWithExecutionData(1, 100)
	task.ExecutionData.TxOpts.TxHashes = append(task.ExecutionData.TxOpts.TxHashes, common.BigToHash(big.NewInt(123871239)))

	wg := sync.WaitGroup{}

	eth := mocks.NewMockEthereum()
	eth.GethClientMock.TransactionByHashFunc.SetDefaultReturn(&types.Transaction{}, false, nil)
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(2)}, nil)
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{}, errors.New("error getting receipt"))
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(int64(minedInBlock))}, nil)
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(int64(minedInBlock))}, nil)

	tasks.StartTask(mocks.NewMockLogger().WithField("Task", 0), &wg, eth, task, state)

	wg.Wait()

	mockrequire.Called(t, task.DoWorkFunc)
	mockrequire.Called(t, task.DoDoneFunc)
	assert.Len(t, task.ExecutionData.TxOpts.TxHashes, 1)
	assert.Equal(t, task.ExecutionData.TxOpts.MinedInBlock, uint64(minedInBlock))
}

// Tx was mined after a retry because of a failed receipt
func TestStartTask_handleExecutedTask_RetrySameFee(t *testing.T) {
	minedInBlock := 7

	state := objects.NewDkgState(accounts.Account{})
	task := mocks.NewMockTaskWithExecutionData(1, 100)
	task.ExecutionData.TxOpts.TxHashes = append(task.ExecutionData.TxOpts.TxHashes, common.BigToHash(big.NewInt(123871239)))
	task.ShouldRetryFunc.SetDefaultReturn(true)

	wg := sync.WaitGroup{}

	eth := mocks.NewMockEthereum()
	eth.GethClientMock.TransactionByHashFunc.SetDefaultHook(func(context.Context, common.Hash) (*types.Transaction, bool, error) {
		time.Sleep(20 * time.Millisecond)
		return &types.Transaction{}, false, nil
	})
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: 0}, nil)
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(int64(minedInBlock))}, nil)
	eth.GethClientMock.TransactionReceiptFunc.PushReturn(&types.Receipt{Status: uint64(1), BlockNumber: big.NewInt(int64(minedInBlock))}, nil)

	tasks.StartTask(mocks.NewMockLogger().WithField("Task", 0), &wg, eth, task, state)

	wg.Wait()

	mockrequire.Called(t, task.DoWorkFunc)
	mockrequire.Called(t, task.DoRetryFunc)
	mockrequire.Called(t, task.DoDoneFunc)
	assert.Len(t, task.ExecutionData.TxOpts.TxHashes, 1)
	assert.Equal(t, task.ExecutionData.TxOpts.MinedInBlock, uint64(minedInBlock))
}

/*
// Tx reached replacement timeout, tx mined after retry with replacement
func TestStartTask_handleExecutedTask_RetryReplacingFee(t *testing.T) {
	logger := logging.GetLogger("test")
	minedInBlock := 10

	state := objects.NewDkgState(accounts.Account{})
	dkgTaskMock := mocks.NewMockTask()
	dkgTaskMock.On("Initialize", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	dkgTaskMock.On("DoWork", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	dkgTaskMock.On("ShouldRetry", mock.Anything, mock.Anything, mock.Anything).Return(true)
	dkgTaskMock.On("DoRetry", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	wg := sync.WaitGroup{}

	geth := mocks.NewMockGethClient()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, true, nil).Once()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, true, nil).Once()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, false, nil)
	receiptOk := &types.Receipt{
		Status:      uint64(1),
		BlockNumber: big.NewInt(int64(minedInBlock)),
	}
	geth.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receiptOk, nil)

	eth := mocks.NewMockEthereum()
	ethMock.On("GetGethClient").Return(geth)
	ethMock.On("GetFinalityDelay").Return(2)
	ethMock.On("RetryCount").Return(3)
	ethMock.On("RetryDelay").Return(1 * time.Millisecond)
	ethMock.On("GetTxCheckFrequency").Return(3 * time.Second)
	ethMock.On("GetTxTimeoutForReplacement").Return(6 * time.Second)
	ethMock.On("GetTxFeePercentageToIncrease").Return(43)
	ethMock.On("GetTxMaxFeeThresholdInGwei").Return(uint64(1000000))

	ethMock.On("GetCurrentHeight", mock.Anything).Return(1, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(2, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(3, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(4, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(5, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(6, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(7, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(8, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(9, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(10, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(11, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(12, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(13, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(14, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(15, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(16, nil).Once()

	tasks.StartTask(logger.WithField("Task", 0), &wg, ethMock, dkgTaskMock, state, nil)

	expectedGasFeeCap := big.NewInt(203569)
	expectedGasTipCap := big.NewInt(52)

	wg.Wait()

	assert.False(t, dkgTaskMock.Success)
	assert.NotEqual(t, 0, len(dkgTaskMock.TxOpts.TxHashes))
	assert.Equal(t, uint64(minedInBlock), dkgTaskMock.TxOpts.MinedInBlock)
	assert.Equal(t, expectedGasFeeCap, dkgTaskMock.TxOpts.GasFeeCap)
	assert.Equal(t, expectedGasTipCap, dkgTaskMock.TxOpts.GasTipCap)
}

// Tx reached replacement timeout, tx mined after retry with replacement
func TestStartTask_handleExecutedTask_RetryReplacingFeeExceedingThreshold(t *testing.T) {
	logger := logging.GetLogger("test")
	minedInBlock := 10

	state := objects.NewDkgState(accounts.Account{})
	dkgTaskMock := mocks.NewMockTask()
	dkgTaskMock.On("Initialize", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	dkgTaskMock.On("DoWork", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	dkgTaskMock.On("ShouldRetry", mock.Anything, mock.Anything, mock.Anything).Return(true)
	dkgTaskMock.On("DoRetry", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	wg := sync.WaitGroup{}

	geth := mocks.NewMockGethClient()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, true, nil).Once()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, true, nil).Once()
	geth.On("TransactionByHash", mock.Anything, mock.Anything).Return(&types.Transaction{}, false, nil)
	receiptOk := &types.Receipt{
		Status:      uint64(1),
		BlockNumber: big.NewInt(int64(minedInBlock)),
	}
	geth.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receiptOk, nil)

	eth := mocks.NewMockEthereum()
	ethMock.On("GetGethClient").Return(geth)
	ethMock.On("GetFinalityDelay").Return(2)
	ethMock.On("RetryCount").Return(3)
	ethMock.On("RetryDelay").Return(1 * time.Millisecond)
	ethMock.On("GetTxCheckFrequency").Return(3 * time.Second)
	ethMock.On("GetTxTimeoutForReplacement").Return(6 * time.Second)
	ethMock.On("GetTxFeePercentageToIncrease").Return(143)
	ethMock.On("GetTxMaxFeeThresholdInGwei").Return(uint64(200000))

	ethMock.On("GetCurrentHeight", mock.Anything).Return(1, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(2, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(3, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(4, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(5, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(6, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(7, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(8, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(9, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(10, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(11, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(12, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(13, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(14, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(15, nil).Once()
	ethMock.On("GetCurrentHeight", mock.Anything).Return(16, nil).Once()

	tasks.StartTask(logger.WithField("Task", 0), &wg, ethMock, dkgTaskMock, state, nil)

	expectedGasFeeCap := big.NewInt(200000)
	expectedGasTipCap := big.NewInt(89)

	wg.Wait()

	assert.False(t, dkgTaskMock.Success)
	assert.NotEqual(t, 0, len(dkgTaskMock.TxOpts.TxHashes))
	assert.Equal(t, uint64(minedInBlock), dkgTaskMock.TxOpts.MinedInBlock)
	assert.Equal(t, expectedGasFeeCap, dkgTaskMock.TxOpts.GasFeeCap)
	assert.Equal(t, expectedGasTipCap, dkgTaskMock.TxOpts.GasTipCap)
}
*/
