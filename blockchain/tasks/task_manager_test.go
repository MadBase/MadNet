package tasks_test

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	"github.com/MadBase/MadNet/blockchain/dkg/dkgtasks"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/blockchain/tasks"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/test/mocks"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//
// Mock implementation of interfaces.Task
//
type mockState struct {
	sync.Mutex
	count int
}

type mockTask struct {
	DoneCalled bool
	State      *mockState
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

func (mt *mockTask) Initialize(ctx context.Context, logger *logrus.Entry, eth interfaces.Ethereum, state interface{}) error {
	dkgState := state.(*mockState)

	mt.State = dkgState

	mt.State.Lock()
	defer mt.State.Unlock()

	mt.State.count += 1

	return nil
}

func (mt *mockTask) ShouldRetry(context.Context, *logrus.Entry, interfaces.Ethereum) bool {
	return false
}

func TestFoo(t *testing.T) {
	var s map[string]string

	raw, err := json.Marshal(s)
	assert.Nil(t, err)

	t.Logf("Raw data:%v", string(raw))
}

func TestType(t *testing.T) {
	state := objects.NewDkgState(accounts.Account{})
	ct := dkgtasks.NewCompletionTask(state, 1, 10)

	var task interfaces.Task = ct
	raw, err := json.Marshal(task)
	assert.Nil(t, err)
	assert.Greater(t, len(raw), 1)

	tipe := reflect.TypeOf(task)
	t.Logf("type0:%v", tipe.String())

	if tipe.Kind() == reflect.Ptr {
		tipe = tipe.Elem()
	}

	typeName := tipe.String()
	t.Logf("type1:%v", typeName)

}

func TestSharedState(t *testing.T) {
	logger := logging.GetLogger("test")

	state := &mockState{}

	task0 := &mockTask{}
	task1 := &mockTask{}

	wg := sync.WaitGroup{}

	eth := mocks.NewMockEthereum()

	tasks.StartTask(logger.WithField("Task", 0), &wg, eth, task0, state)
	tasks.StartTask(logger.WithField("Task", 1), &wg, eth, task1, state)

	wg.Wait()

	assert.Equal(t, 2, state.count)
}

func TestIsAdminClient(t *testing.T) {
	adminInterface := reflect.TypeOf((*interfaces.AdminClient)(nil)).Elem()

	task := &dkgtasks.GPKjSubmissionTask{}
	isAdminClient := reflect.TypeOf(task).Implements(adminInterface)

	assert.True(t, isAdminClient)
}
