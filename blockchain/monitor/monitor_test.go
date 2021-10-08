package monitor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	aobjs "github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/blockchain"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/monitor"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/blockchain/tasks"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/logging"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func setupEthereum(t *testing.T, mineInterval time.Duration) interfaces.Ethereum {

	eth, err := blockchain.NewEthereumSimulator(
		"../../assets/test/keys",
		"../../assets/test/passcodes.txt",
		3,
		2*time.Second,
		5*time.Second,
		0,
		big.NewInt(9223372036854775807),
		"0x26D3D8Ab74D62C26f1ACc220dA1646411c9880Ac",
		"0x546F99F244b7B58B855330AE0E2BC1b30b41302F")

	assert.Nil(t, err, "Failed to build Ethereum endpoint...")
	assert.True(t, eth.IsEthereumAccessible(), "Web3 endpoint is not available.")

	c := eth.Contracts()

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

	deployAcct, _ := eth.GetAccount(common.HexToAddress("0x546F99F244b7B58B855330AE0E2BC1b30b41302F"))
	err = eth.UnlockAccount(deployAcct)
	assert.Nil(t, err, "Failed to unlock deploy account")

	bal, err := eth.GetBalance(acct.Address)
	assert.Nil(t, err, "Can't check balance for %v.", deployAcct.Address.Hex())
	t.Logf(" deploy account (%v) balance: %v", deployAcct.Address.Hex(), bal)

	// Unlock testing account and make sure it has a balance
	err = eth.UnlockAccount(acct)
	assert.Nil(t, err, "Failed to unlock default account")

	bal, err = eth.GetBalance(acct.Address)
	assert.Nil(t, err, "Can't check balance for %v.", acct.Address.Hex())
	t.Logf("default account (%v) balance: %v", acct.Address.Hex(), bal)

	// Transfer some eth
	testingEth := big.NewInt(200000)
	t.Logf("Transfering %v from %v to %v", testingEth, deployAcct.Address.Hex(), acct.Address.Hex())
	_, err = eth.TransferEther(deployAcct.Address, acct.Address, testingEth)
	assert.Nil(t, err, "Failed to transfer ether to default account")

	_, _, err = c.DeployContracts(context.TODO(), acct)
	assert.Nil(t, err, "Failed to deploy contracts...")

	return eth
}

//
//
//
func createSharedKey(t *testing.T, addr common.Address) [4]*big.Int {

	b := addr.Bytes()

	return [4]*big.Int{
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b),
		(&big.Int{}).SetBytes(b)}
}

func createValidator(t *testing.T, addrHex string, idx uint8) objects.Validator {
	addr := common.HexToAddress(addrHex)
	return objects.Validator{
		Account:   addr,
		Index:     idx,
		SharedKey: createSharedKey(t, addr),
	}
}

//
// Mock implementation of interfaces.Task
//
type mockTask struct {
	DoneCalled bool
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

func (mt *mockTask) Initialize(context.Context, *logrus.Entry, interfaces.Ethereum) error {
	return nil
}

func (mt *mockTask) ShouldRetry(context.Context, *logrus.Entry, interfaces.Ethereum) bool {
	return false
}

//
// Mock implementation of interfaces.AdminHandler
//
type mockAdminHandler struct {
}

func (ah *mockAdminHandler) AddPrivateKey([]byte, constants.CurveSpec) error {
	return nil
}

func (ah *mockAdminHandler) AddSnapshot(*objs.BlockHeader, bool) error {
	return nil
}
func (ah *mockAdminHandler) AddValidatorSet(*objs.ValidatorSet) error {
	return nil
}

func (ah *mockAdminHandler) RegisterSnapshotCallback(func(*objs.BlockHeader) error) {

}

func (ah *mockAdminHandler) SetSynchronized(v bool) {

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
// Mock implementation of interfaces.Ethereum
//
type mockEthereum struct {
}

func (eth *mockEthereum) ChainID() *big.Int {
	return nil
}

func (eth *mockEthereum) Close() error {
	return nil
}

func (eth *mockEthereum) Commit() {

}

func (eth *mockEthereum) IsEthereumAccessible() bool {
	return false
}

func (eth *mockEthereum) GetCallOpts(context.Context, accounts.Account) *bind.CallOpts {
	return nil
}

func (eth *mockEthereum) GetTransactionOpts(context.Context, accounts.Account) (*bind.TransactOpts, error) {
	return nil, nil
}

func (eth *mockEthereum) LoadAccounts(string) {}

func (eth *mockEthereum) LoadPasscodes(string) error {
	return nil
}

func (eth *mockEthereum) UnlockAccount(accounts.Account) error {
	return nil
}

func (eth *mockEthereum) TransferEther(common.Address, common.Address, *big.Int) (*types.Transaction, error) {
	return nil, nil
}

func (eth *mockEthereum) GetAccount(addr common.Address) (accounts.Account, error) {
	return accounts.Account{Address: addr}, nil
}
func (eth *mockEthereum) GetAccountKeys(addr common.Address) (*keystore.Key, error) {
	return nil, nil
}
func (eth *mockEthereum) GetBalance(common.Address) (*big.Int, error) {
	return nil, nil
}
func (eth *mockEthereum) GetGethClient() interfaces.GethClient {
	return nil
}

func (eth *mockEthereum) GetCoinbaseAddress() common.Address {
	return eth.GetDefaultAccount().Address
}

func (eth *mockEthereum) GetCurrentHeight(context.Context) (uint64, error) {
	return 0, nil
}

func (eth *mockEthereum) GetDefaultAccount() accounts.Account {
	return accounts.Account{}
}
func (eth *mockEthereum) GetEndpoint() string {
	return "na"
}
func (eth *mockEthereum) GetEvents(ctx context.Context, firstBlock uint64, lastBlock uint64, addresses []common.Address) ([]types.Log, error) {
	return nil, nil
}
func (eth *mockEthereum) GetFinalizedHeight(context.Context) (uint64, error) {
	return 0, nil
}
func (eth *mockEthereum) GetPeerCount(context.Context) (uint64, error) {
	return 0, nil
}
func (eth *mockEthereum) GetSnapshot() ([]byte, error) {
	return nil, nil
}
func (eth *mockEthereum) GetSyncProgress() (bool, *ethereum.SyncProgress, error) {
	return false, nil, nil
}
func (eth *mockEthereum) GetTimeoutContext() (context.Context, context.CancelFunc) {
	return nil, nil
}
func (eth *mockEthereum) GetValidators(context.Context) ([]common.Address, error) {
	return nil, nil
}

func (eth *mockEthereum) KnownSelectors() interfaces.SelectorMap {
	return nil
}

func (eth *mockEthereum) Queue() interfaces.TxnQueue {
	return nil
}

func (eth *mockEthereum) RetryCount() int {
	return 0
}
func (eth *mockEthereum) RetryDelay() time.Duration {
	return time.Second
}

func (eth *mockEthereum) Timeout() time.Duration {
	return time.Second
}

func (eth *mockEthereum) Contracts() interfaces.Contracts {
	return nil
}

//
// Actual tests
//
func TestMonitor(t *testing.T) {
	eth := setupEthereum(t, time.Second)
	c := eth.Contracts()

	acct := eth.GetDefaultAccount()

	t.Logf("eth scheme:%v url:%v", acct.URL.Scheme, acct.URL.Path)

	txnOpts, err := eth.GetTransactionOpts(context.TODO(), eth.GetDefaultAccount())
	assert.Nil(t, err, "Failed to build txnOpts endpoint... %v", err)

	_, err = c.Ethdkg().InitializeState(txnOpts)
	assert.Nil(t, err, "Failed to Initialize state... %v", err)

	eth.Commit()
}

func TestWTF(t *testing.T) {
	type foo struct {
		sync.Mutex
		count int
	}

	f := foo{}
	for {
		f.Lock()
		f.count++
		fmt.Printf("Count:%v\n", f.count)
		time.Sleep(1 * time.Second)

		if f.count > 10 {
			f.Unlock()
			break
		}
		f.Unlock()
	}
}

func TestBidirectionalMarshaling(t *testing.T) {

	// setup
	adminHandler := &mockAdminHandler{}
	depositHandler := &mockDepositHandler{}
	eth := &mockEthereum{}
	logger := logging.GetLogger("test")

	addr0 := common.HexToAddress("0x546F99F244b7B58B855330AE0E2BC1b30b41302F")

	EPOCH := uint32(1)

	// Setup monitor state
	mon, err := monitor.NewMonitor(&db.Database{}, adminHandler, depositHandler, eth, 2*time.Second, time.Minute, 1)
	assert.Nil(t, err)

	mon.State.EthDKG.Account = accounts.Account{
		Address: addr0,
		URL: accounts.URL{
			Scheme: "keystore",
			Path:   "/home/agdean/Projects/MadNet/assets/test/keys/UTC--2020-03-24T13-41-44.886736400Z--26d3d8ab74d62c26f1acc220da1646411c9880ac"}}
	mon.State.EthDKG.Index = 1
	mon.State.EthDKG.SecretValue = big.NewInt(512)
	mon.State.EthDKG.GroupPublicKeys[addr0] = [4]*big.Int{
		big.NewInt(44), big.NewInt(33), big.NewInt(22), big.NewInt(11)}
	mon.State.EthDKG.Commitments[addr0] = make([][2]*big.Int, 3)
	mon.State.EthDKG.Commitments[addr0][0][0] = big.NewInt(5)
	mon.State.EthDKG.Commitments[addr0][0][1] = big.NewInt(2)

	mon.State.ValidatorSets[EPOCH] = objects.ValidatorSet{
		ValidatorCount:        4,
		NotBeforeMadNetHeight: 321,
		GroupKey:              [4]*big.Int{big.NewInt(3), big.NewInt(2), big.NewInt(1), big.NewInt(5)}}

	mon.State.Validators[EPOCH] = []objects.Validator{
		createValidator(t, "0x546F99F244b7B58B855330AE0E2BC1b30b41302F", 1),
		createValidator(t, "0x9AC1c9afBAec85278679fF75Ef109217f26b1417", 2),
		createValidator(t, "0x26D3D8Ab74D62C26f1ACc220dA1646411c9880Ac", 3),
		createValidator(t, "0x615695C4a4D6a60830e5fca4901FbA099DF26271", 4)}

	mon.State.Schedule.Schedule(5, 10, &mockTask{})

	// Marshal
	mon.TypeRegistry.RegisterInstanceType(&mockTask{})
	raw, err := json.Marshal(mon)
	assert.Nil(t, err)
	t.Logf("RawData:%v", string(raw))

	// Unmarshal
	newMon, err := monitor.NewMonitor(&db.Database{}, adminHandler, depositHandler, eth, 2*time.Second, time.Minute, 1)
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
	validator0 := createValidator(t, "0x546F99F244b7B58B855330AE0E2BC1b30b41302F", 1)

	assert.Equal(t, 0, validator0.SharedKey[0].Cmp(newMon.State.Validators[EPOCH][0].SharedKey[0]))
	assert.Equal(t, 0, big.NewInt(44).Cmp(newMon.State.EthDKG.GroupPublicKeys[addr0][0]))

	// Compare the schedules
	_, err = newMon.State.Schedule.Find(2)
	assert.Equal(t, objects.ErrNothingScheduled, err)

	taskID, err := newMon.State.Schedule.Find(6)
	assert.Nil(t, err)

	task, err := newMon.State.Schedule.Retrieve(taskID)
	assert.Nil(t, err)

	taskStruct := task.(*mockTask)
	assert.False(t, taskStruct.DoneCalled)

	wg := &sync.WaitGroup{}
	tasks.StartTask(logger.WithField("Task", "Mocked"), wg, eth, task)

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
