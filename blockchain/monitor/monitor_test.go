package monitor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	aobjs "github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/blockchain"
	"github.com/MadBase/MadNet/blockchain/etest"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/monitor"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/blockchain/tasks"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/test/mocks"

	mockrequire "github.com/derision-test/go-mockgen/testutil/require"

	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEthereum(t *testing.T, mineInterval time.Duration) interfaces.Ethereum {

	n := 4
	privKeys := etest.SetupPrivateKeys(n)
	eth, err := blockchain.NewEthereumSimulator(
		privKeys,
		1,
		time.Second*2,
		time.Second*5,
		0,
		big.NewInt(math.MaxInt64))
	assert.Nil(t, err, "Failed to build Ethereum endpoint...")
	assert.True(t, eth.IsEthereumAccessible(), "Web3 endpoint is not available.")
	defer eth.Close()

	// c := eth.Contracts()

	go func() {
		for {
			time.Sleep(mineInterval)
			eth.Commit()
		}
	}()

	// Unlock deploy account and make sure it has a balance
	acct := eth.GetDefaultAccount()
	err = eth.UnlockAccount(acct)
	assert.Nil(t, err, "Failed to unlock deploy account")

	// _, _, err = c.DeployContracts(context.TODO(), acct)
	// assert.Nil(t, err, "Failed to deploy contracts...")
	panic("needs deployment")

	return eth
}

//
//
//
func createSharedKey(addr common.Address) [4]*big.Int {

	b := addr.Bytes()

	return [4]*big.Int{
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b)}
}

func createValidator(addrHex string, idx uint8) objects.Validator {
	addr := common.HexToAddress(addrHex)
	return objects.Validator{
		Account:   addr,
		Index:     idx,
		SharedKey: createSharedKey(addr),
	}
}

func populateMonitor(state *objects.MonitorState, addr0 common.Address, EPOCH uint32) {
	state.EthDKG.Account = accounts.Account{
		Address: addr0,
		URL: accounts.URL{
			Scheme: "keystore",
			Path:   ""}}
	state.EthDKG.Index = 1
	state.EthDKG.SecretValue = big.NewInt(512)
	meAsAParticipant := &objects.Participant{
		Address: state.EthDKG.Account.Address,
		Index:   state.EthDKG.Index,
	}
	state.EthDKG.Participants[addr0] = meAsAParticipant
	state.EthDKG.Participants[addr0].GPKj = [4]*big.Int{
		big.NewInt(44), big.NewInt(33), big.NewInt(22), big.NewInt(11)}
	state.EthDKG.Participants[addr0].Commitments = make([][2]*big.Int, 3)
	state.EthDKG.Participants[addr0].Commitments[0][0] = big.NewInt(5)
	state.EthDKG.Participants[addr0].Commitments[0][1] = big.NewInt(2)

	state.ValidatorSets[EPOCH] = objects.ValidatorSet{
		ValidatorCount:        4,
		NotBeforeMadNetHeight: 321,
		GroupKey:              [4]*big.Int{big.NewInt(3), big.NewInt(2), big.NewInt(1), big.NewInt(5)}}

	state.Validators[EPOCH] = []objects.Validator{
		createValidator("0x546F99F244b7B58B855330AE0E2BC1b30b41302F", 1),
		createValidator("0x9AC1c9afBAec85278679fF75Ef109217f26b1417", 2),
		createValidator("0x26D3D8Ab74D62C26f1ACc220dA1646411c9880Ac", 3),
		createValidator("0x615695C4a4D6a60830e5fca4901FbA099DF26271", 4)}

}

//
// Mock implementation of interfaces.Task
//
type mockTask struct {
	DoneCalled bool
	State      *objects.DkgState
}

func (mt *mockTask) DoDone(logger *logrus.Entry) {
	mt.DoneCalled = true
}

func (mt *mockTask) DoRetry(context.Context, *logrus.Entry, interfaces.Ethereum) error {
	return nil
}

func (mt *mockTask) DoWork(context.Context, *logrus.Entry, interfaces.Ethereum) error {
	return nil
}

func (mt *mockTask) Initialize(context.Context, *logrus.Entry, interfaces.Ethereum, interface{}) error {
	return nil
}

func (mt *mockTask) ShouldRetry(context.Context, *logrus.Entry, interfaces.Ethereum) bool {
	return false
}

//
// Mock implementation of interfaces.DepositHandler
//
type mockDepositHandler struct {
}

func (dh *mockDepositHandler) Add(*badger.Txn, uint32, []byte, *big.Int, *aobjs.Owner) error {
	return nil
}

//
// Actual tests
//
func TestMonitorPersist(t *testing.T) {
	rawDb, err := utils.OpenBadger(context.Background().Done(), "", true)
	assert.Nil(t, err)

	database := &db.Database{}
	database.Init(rawDb)

	eth := mocks.NewMockEthereum()
	mon, err := monitor.NewMonitor(database, database, mocks.NewMockAdminHandler(), &mockDepositHandler{}, eth, 1*time.Second, time.Minute, 1)
	assert.Nil(t, err)

	addr0 := common.HexToAddress("0x546F99F244b7B58B855330AE0E2BC1b30b41302F")
	EPOCH := uint32(1)
	populateMonitor(mon.State, addr0, EPOCH)
	raw, err := json.Marshal(mon)
	assert.Nil(t, err)
	t.Logf("Raw: %v", string(raw))

	mon.PersistState()

	//
	newMon, err := monitor.NewMonitor(database, database, mocks.NewMockAdminHandler(), &mockDepositHandler{}, eth, 1*time.Second, time.Minute, 1)
	assert.Nil(t, err)

	newMon.LoadState()

	newRaw, err := json.Marshal(mon)
	assert.Nil(t, err)
	t.Logf("NewRaw: %v", string(newRaw))
}

func TestBidirectionalMarshaling(t *testing.T) {

	// setup
	adminHandler := mocks.NewMockAdminHandler()
	depositHandler := &mockDepositHandler{}
	eth := mocks.NewMockEthereum()
	logger := logging.GetLogger("test")

	addr0 := common.HexToAddress("0x546F99F244b7B58B855330AE0E2BC1b30b41302F")

	EPOCH := uint32(1)

	// Setup monitor state
	mon, err := monitor.NewMonitor(&db.Database{}, &db.Database{}, adminHandler, depositHandler, eth, 2*time.Second, time.Minute, 1)
	assert.Nil(t, err)
	populateMonitor(mon.State, addr0, EPOCH)

	// Schedule some tasks
	_, err = mon.State.Schedule.Schedule(1, 2, &mockTask{})
	assert.Nil(t, err)

	_, err = mon.State.Schedule.Schedule(3, 4, &mockTask{})
	assert.Nil(t, err)

	_, err = mon.State.Schedule.Schedule(5, 6, &mockTask{})
	assert.Nil(t, err)

	_, err = mon.State.Schedule.Schedule(7, 8, &mockTask{})
	assert.Nil(t, err)

	// Marshal
	mon.TypeRegistry.RegisterInstanceType(&mockTask{})
	raw, err := json.Marshal(mon)
	assert.Nil(t, err)
	t.Logf("RawData:%v", string(raw))

	// Unmarshal
	newMon, err := monitor.NewMonitor(&db.Database{}, &db.Database{}, adminHandler, depositHandler, eth, 2*time.Second, time.Minute, 1)
	assert.Nil(t, err)

	newMon.TypeRegistry.RegisterInstanceType(&mockTask{})
	err = json.Unmarshal(raw, newMon)
	assert.Nil(t, err)

	// Compare raw data for mon and newMon
	newRaw, err := json.Marshal(newMon)
	assert.Nil(t, err)
	assert.Equal(t, len(raw), len(newRaw))
	t.Logf("Len(RawData): %v", len(raw))

	// Do comparisons
	validator0 := createValidator("0x546F99F244b7B58B855330AE0E2BC1b30b41302F", 1)

	assert.Equal(t, 0, validator0.SharedKey[0].Cmp(newMon.State.Validators[EPOCH][0].SharedKey[0]))
	assert.Equal(t, 0, big.NewInt(44).Cmp(newMon.State.EthDKG.Participants[addr0].GPKj[0]))

	// Compare the schedules
	_, err = newMon.State.Schedule.Find(9)
	assert.Equal(t, objects.ErrNothingScheduled, err)

	//
	taskID, err := newMon.State.Schedule.Find(1)
	assert.Nil(t, err)

	task, err := newMon.State.Schedule.Retrieve(taskID)
	assert.Nil(t, err)

	//
	taskID2, err := newMon.State.Schedule.Find(3)
	assert.Nil(t, err)

	task2, err := newMon.State.Schedule.Retrieve(taskID2)
	assert.Nil(t, err)

	taskStruct := task.(*mockTask)
	assert.False(t, taskStruct.DoneCalled)

	taskStruct2 := task2.(*mockTask)
	assert.False(t, taskStruct2.DoneCalled)

	t.Logf("State:%p State2:%p", taskStruct.State, taskStruct2.State)
	assert.Equal(t, taskStruct.State, taskStruct2.State)

	wg := &sync.WaitGroup{}
	tasks.StartTask(logger.WithField("Task", "Mocked"), wg, eth, task, nil)
	wg.Wait()

	assert.True(t, taskStruct.DoneCalled)
}

func TestWrapDoNotContinue(t *testing.T) {
	genErr := objects.ErrCanNotContinue
	specErr := errors.New("neutrinos")

	niceErr := errors.Wrapf(genErr, "Caused by %v", specErr)
	assert.True(t, errors.Is(niceErr, genErr))

	t.Logf("NiceErr: %v", niceErr)

	nice2Err := fmt.Errorf("%w because %v", genErr, specErr)
	assert.True(t, errors.Is(nice2Err, genErr))

	t.Logf("Nice2Err: %v", nice2Err)
}

func GetSnapshotCallbacks(t *testing.T) (*mocks.EthMockset, mocks.AdminHandlerRegisterSnapshotCallbacksFuncCall, func() *objs.CachedSnapshotTx) {
	adminHandler := mocks.NewMockAdminHandler()
	depositHandler := &mockDepositHandler{}
	eth, mockset := mocks.NewMockLinkedEthereum()

	cdb := mocks.NewMockDb()
	_, err := monitor.NewMonitor(cdb, mocks.NewMockDb(), adminHandler, depositHandler, eth, 2*time.Second, time.Minute, 1)
	require.Nil(t, err)

	mockrequire.Called(t, adminHandler.RegisterSnapshotCallbacksFunc)
	callbacks := adminHandler.RegisterSnapshotCallbacksFunc.History()[0]

	return mockset, callbacks, func() *objs.CachedSnapshotTx {
		var tx *objs.CachedSnapshotTx
		cdb.View(func(txn *badger.Txn) (err error) {
			tx, _ = cdb.GetSnapshotTx(txn)
			return nil
		})
		return tx
	}
}

func TestNewSS(t *testing.T) {
	mockset, callbacks, getDbTx := GetSnapshotCallbacks(t)
	newSS := callbacks.Arg0

	blockHeader := &objs.BlockHeader{}
	blockHeader.UnmarshalBinary([]byte{0, 0, 0, 0, 0, 0, 3, 0, 8, 0, 0, 0, 1, 0, 4, 0, 89, 0, 0, 0, 2, 6, 0, 0, 181, 0, 0, 0, 2, 0, 0, 0, 42, 0, 0, 0, 8, 3, 0, 0, 13, 0, 0, 0, 2, 1, 0, 0, 25, 0, 0, 0, 2, 1, 0, 0, 37, 0, 0, 0, 2, 1, 0, 0, 49, 0, 0, 0, 2, 1, 0, 0, 125, 56, 56, 255, 62, 64, 136, 59, 115, 108, 129, 228, 18, 133, 160, 220, 127, 56, 179, 7, 55, 215, 39, 111, 187, 195, 120, 118, 22, 203, 242, 201, 197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 206, 149, 115, 249, 102, 3, 185, 58, 36, 197, 107, 9, 40, 75, 179, 75, 161, 110, 58, 11, 133, 64, 12, 26, 144, 171, 90, 190, 173, 128, 60, 206, 31, 127, 172, 22, 115, 42, 90, 205, 18, 11, 164, 170, 60, 91, 40, 65, 170, 211, 47, 172, 44, 247, 93, 58, 91, 198, 254, 95, 56, 0, 12, 145, 20, 223, 81, 247, 203, 176, 58, 150, 81, 152, 211, 30, 159, 117, 124, 215, 247, 243, 135, 30, 155, 76, 194, 237, 96, 57, 174, 171, 197, 56, 239, 144, 25, 233, 194, 32, 143, 161, 28, 164, 130, 145, 219, 25, 218, 129, 134, 165, 202, 50, 142, 130, 94, 240, 142, 111, 239, 190, 137, 174, 189, 13, 194, 74, 42, 34, 151, 20, 115, 75, 95, 78, 250, 216, 5, 12, 36, 204, 133, 118, 173, 38, 28, 11, 16, 64, 11, 204, 37, 233, 110, 62, 217, 0, 185, 87, 16, 199, 104, 254, 97, 159, 166, 106, 239, 71, 255, 158, 63, 80, 171, 211, 175, 43, 93, 114, 134, 3, 0, 211, 177, 136, 111, 26, 22, 108, 118, 199, 29, 106, 82, 61, 246, 187, 142, 226, 14, 167, 116, 171, 194, 244, 157, 203, 217, 127, 150, 130, 49, 129, 224, 242, 60, 229, 35, 70, 107, 245, 20, 122})

	require.Nil(t, getDbTx())
	err := newSS(blockHeader)
	require.Nil(t, err)

	mockrequire.Called(t, mockset.Queue.QueueTransactionSyncFunc)
	require.Equal(t, objs.SnapshotTxVerified, getDbTx().State)
}

func TestResumeSSCreated(t *testing.T) {
	mockset, callbacks, getDbTx := GetSnapshotCallbacks(t)
	resumeSS := callbacks.Arg1

	sstx := mocks.NewMockSnapshotTx()

	require.Nil(t, getDbTx())
	err := resumeSS(&objs.CachedSnapshotTx{State: objs.SnapshotTxCreated, Tx: sstx})
	require.Nil(t, err)

	mockrequire.CalledWith(t, mockset.Queue.QueueTransactionSyncFunc, mockrequire.Values(mockrequire.Skip, sstx))
	require.Equal(t, objs.SnapshotTxVerified, getDbTx().State)
}

func TestResumeSSSubmitted(t *testing.T) {
	mockset, callbacks, getDbTx := GetSnapshotCallbacks(t)
	resumeSS := callbacks.Arg1

	sstx := mocks.NewMockSnapshotTx()

	require.Nil(t, getDbTx())
	err := resumeSS(&objs.CachedSnapshotTx{State: objs.SnapshotTxSubmitted, Tx: sstx})
	require.Nil(t, err)

	mockrequire.NotCalled(t, mockset.Queue.QueueTransactionSyncFunc)
	mockrequire.CalledWith(t, mockset.Queue.WaitTransactionFunc, mockrequire.Values(mockrequire.Skip, sstx))
	require.Equal(t, objs.SnapshotTxVerified, getDbTx().State)
}

func TestResumeSSDbWrite(t *testing.T) {
	mockset, callbacks, _ := GetSnapshotCallbacks(t)
	resumeSS := callbacks.Arg1

	sstx := mocks.NewMockSnapshotTx()
	err := resumeSS(&objs.CachedSnapshotTx{State: objs.SnapshotTxSubmitted, Tx: sstx})
	require.Nil(t, err)

	mockrequire.NotCalled(t, mockset.Queue.QueueTransactionSyncFunc)
	mockrequire.CalledWith(t, mockset.Queue.WaitTransactionFunc, mockrequire.Values(mockrequire.Skip, sstx))
}

func TestGetSSEthHeight(t *testing.T) {
	mockset, callbacks, _ := GetSnapshotCallbacks(t)
	getSSEthHeight := callbacks.Arg2

	_, err := getSSEthHeight(big.NewInt(12))
	require.Nil(t, err)

	mockrequire.Called(t, mockset.Snapshots.GetCommittedHeightFromSnapshotFunc)
	require.Equal(t, mockset.Snapshots.GetCommittedHeightFromSnapshotFunc.History()[0].Arg1, big.NewInt(12))
}

func TestRefreshTxRecent(t *testing.T) {
	_, callbacks, _ := GetSnapshotCallbacks(t)
	refreshTx := callbacks.Arg3

	tx1 := mocks.NewMockSnapshotTx()
	tx2, err := refreshTx(context.Background(), tx1, 1)
	require.Nil(t, err)
	require.Nil(t, tx2)
}

func TestRefreshTxMined(t *testing.T) {
	mockset, callbacks, _ := GetSnapshotCallbacks(t)
	refreshTx := callbacks.Arg3

	tx1 := mocks.NewMockSnapshotTx()
	mockset.GethClient.TransactionByHashFunc.SetDefaultReturn(tx1, false, nil)
	tx2, err := refreshTx(context.Background(), tx1, 64)
	require.Nil(t, err)
	require.Nil(t, tx2)
}

func TestRefreshTxPending(t *testing.T) {
	mockset, callbacks, _ := GetSnapshotCallbacks(t)
	refreshTx := callbacks.Arg3

	tx1 := mocks.NewMockSnapshotTx()
	mockset.GethClient.TransactionByHashFunc.SetDefaultReturn(tx1, true, nil)
	tx2, err := refreshTx(context.Background(), tx1, 64)
	require.Nil(t, err)

	require.Greater(t, tx2.GasFeeCap().Int64(), 2*tx1.GasFeeCap().Int64())
	require.Greater(t, tx2.GasTipCap().Int64(), tx1.GasTipCap().Int64()+tx1.GasTipCap().Int64()/10)
}

func TestRefreshTxFailed(t *testing.T) {
	mockset, callbacks, _ := GetSnapshotCallbacks(t)
	refreshTx := callbacks.Arg3

	tx1 := mocks.NewMockSnapshotTx()
	mockset.GethClient.TransactionByHashFunc.SetDefaultReturn(tx1, false, fmt.Errorf("failed"))
	tx2, err := refreshTx(context.Background(), tx1, 64)
	require.Nil(t, tx2)
	require.NotNil(t, err)
}
