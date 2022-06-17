package objects

import (
	"context"
	"errors"
	"sync"

	"github.com/MadBase/MadNet/blockchain/ethereum"
	"github.com/MadBase/MadNet/blockchain/executor/interfaces"
	"github.com/MadBase/MadNet/blockchain/transaction"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type InternalTransaction struct {
	mu          sync.RWMutex
	transaction *types.Transaction
}

type Task struct {
	isInitialized       bool
	id                  string
	name                string
	start               uint64
	end                 uint64
	allowMultiExecution bool
	ctx                 context.Context
	cancelFunc          context.CancelFunc
	database            *db.Database
	logger              *logrus.Entry
	client              ethereum.Network
	taskResponseChan    interfaces.ITaskResponseChan
	subscribedTxn       *InternalTransaction
	subscribeOptions    *transaction.SubscribeOptions
}

func NewTask(name string, start uint64, end uint64, allowMultiExecution bool, subscribeOptions *transaction.SubscribeOptions) *Task {
	ctx, cf := context.WithCancel(context.Background())

	return &Task{
		name:                name,
		start:               start,
		end:                 end,
		allowMultiExecution: allowMultiExecution,
		subscribeOptions:    subscribeOptions,
		ctx:                 ctx,
		cancelFunc:          cf,
	}
}

// Initialize default implementation for the ITask interface
func (t *Task) Initialize(ctx context.Context, cancelFunc context.CancelFunc, database *db.Database, logger *logrus.Entry, eth ethereum.Network, id string, taskResponseChan interfaces.ITaskResponseChan) error {
	if !t.isInitialized {
		t.id = id
		t.ctx = ctx
		t.cancelFunc = cancelFunc
		t.database = database
		t.logger = logger
		t.client = eth
		t.taskResponseChan = taskResponseChan
		t.isInitialized = true
		return nil
	}
	return errors.New("trying to initialize task twice!")
}

// GetId default implementation for the ITask interface
func (t *Task) GetId() string {
	return t.id
}

// GetStart default implementation for the ITask interface
func (t *Task) GetStart() uint64 {
	return t.start
}

// GetEnd default implementation for the ITask interface
func (t *Task) GetEnd() uint64 {
	return t.end
}

// GetName default implementation for the ITask interface
func (t *Task) GetName() string {
	return t.name
}

// GetAllowMultiExecution default implementation for the ITask interface
func (t *Task) GetAllowMultiExecution() bool {
	return t.allowMultiExecution
}

func (t *Task) SetSubscribedTx(txn *types.Transaction) {
	t.subscribedTxn.mu.Lock()
	defer t.subscribedTxn.mu.Unlock()
	t.subscribedTxn.transaction = txn
}

func (t *Task) GetSubscribedTx() *types.Transaction {
	t.subscribedTxn.mu.RLock()
	defer t.subscribedTxn.mu.RUnlock()
	return t.subscribedTxn.transaction
}

func (t *Task) GetSubscribeOptions() *transaction.SubscribeOptions {
	return t.subscribeOptions
}

// GetCtx default implementation for the ITask interface
func (t *Task) GetCtx() context.Context {
	return t.ctx
}

// GetEth default implementation for the ITask interface
func (t *Task) GetClient() ethereum.Network {
	return t.client
}

// GetLogger default implementation for the ITask interface
func (t *Task) GetLogger() *logrus.Entry {
	return t.logger
}

// Close default implementation for the ITask interface
func (t *Task) Close() {
	if t.cancelFunc != nil {
		t.cancelFunc()
	}
}

// Finish default implementation for the ITask interface
func (t *Task) Finish(err error) {
	if err != nil {
		t.logger.WithError(err).Errorf("Id: %s, name: %s task is done", t.id, t.name)
	} else {
		t.logger.Infof("Id: %s, name: %s task is done", t.id, t.name)
	}

	t.taskResponseChan.Add(interfaces.TaskResponse{Id: t.id, Err: err})
}

func (t *Task) GetDB() *db.Database {
	return t.database
}