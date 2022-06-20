// Code generated by go-mockgen 1.1.4; DO NOT EDIT.

package mocks

import (
	"context"
	"sync"

	transaction "github.com/MadBase/MadNet/layer1/transaction"
	types "github.com/ethereum/go-ethereum/core/types"
)

// MockWatcher is a mock implementation of the Watcher interface (from the
// package github.com/MadBase/MadNet/layer1/transaction) used for unit
// testing.
type MockWatcher struct {
	// CloseFunc is an instance of a mock function object controlling the
	// behavior of the method Close.
	CloseFunc *WatcherCloseFunc
	// StartFunc is an instance of a mock function object controlling the
	// behavior of the method Start.
	StartFunc *WatcherStartFunc
	// SubscribeFunc is an instance of a mock function object controlling
	// the behavior of the method Subscribe.
	SubscribeFunc *WatcherSubscribeFunc
	// SubscribeAndWaitFunc is an instance of a mock function object
	// controlling the behavior of the method SubscribeAndWait.
	SubscribeAndWaitFunc *WatcherSubscribeAndWaitFunc
	// WaitFunc is an instance of a mock function object controlling the
	// behavior of the method Wait.
	WaitFunc *WatcherWaitFunc
}

// NewMockWatcher creates a new mock of the Watcher interface. All methods
// return zero values for all results, unless overwritten.
func NewMockWatcher() *MockWatcher {
	return &MockWatcher{
		CloseFunc: &WatcherCloseFunc{
			defaultHook: func() {
				return
			},
		},
		StartFunc: &WatcherStartFunc{
			defaultHook: func() error {
				return nil
			},
		},
		SubscribeFunc: &WatcherSubscribeFunc{
			defaultHook: func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
				return nil, nil
			},
		},
		SubscribeAndWaitFunc: &WatcherSubscribeAndWaitFunc{
			defaultHook: func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error) {
				return nil, nil
			},
		},
		WaitFunc: &WatcherWaitFunc{
			defaultHook: func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error) {
				return nil, nil
			},
		},
	}
}

// NewStrictMockWatcher creates a new mock of the Watcher interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockWatcher() *MockWatcher {
	return &MockWatcher{
		CloseFunc: &WatcherCloseFunc{
			defaultHook: func() {
				panic("unexpected invocation of MockWatcher.Close")
			},
		},
		StartFunc: &WatcherStartFunc{
			defaultHook: func() error {
				panic("unexpected invocation of MockWatcher.Start")
			},
		},
		SubscribeFunc: &WatcherSubscribeFunc{
			defaultHook: func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
				panic("unexpected invocation of MockWatcher.Subscribe")
			},
		},
		SubscribeAndWaitFunc: &WatcherSubscribeAndWaitFunc{
			defaultHook: func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error) {
				panic("unexpected invocation of MockWatcher.SubscribeAndWait")
			},
		},
		WaitFunc: &WatcherWaitFunc{
			defaultHook: func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error) {
				panic("unexpected invocation of MockWatcher.Wait")
			},
		},
	}
}

// NewMockWatcherFrom creates a new mock of the MockWatcher interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockWatcherFrom(i transaction.Watcher) *MockWatcher {
	return &MockWatcher{
		CloseFunc: &WatcherCloseFunc{
			defaultHook: i.Close,
		},
		StartFunc: &WatcherStartFunc{
			defaultHook: i.Start,
		},
		SubscribeFunc: &WatcherSubscribeFunc{
			defaultHook: i.Subscribe,
		},
		SubscribeAndWaitFunc: &WatcherSubscribeAndWaitFunc{
			defaultHook: i.SubscribeAndWait,
		},
		WaitFunc: &WatcherWaitFunc{
			defaultHook: i.Wait,
		},
	}
}

// WatcherCloseFunc describes the behavior when the Close method of the
// parent MockWatcher instance is invoked.
type WatcherCloseFunc struct {
	defaultHook func()
	hooks       []func()
	history     []WatcherCloseFuncCall
	mutex       sync.Mutex
}

// Close delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockWatcher) Close() {
	m.CloseFunc.nextHook()()
	m.CloseFunc.appendCall(WatcherCloseFuncCall{})
	return
}

// SetDefaultHook sets function that is called when the Close method of the
// parent MockWatcher instance is invoked and the hook queue is empty.
func (f *WatcherCloseFunc) SetDefaultHook(hook func()) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Close method of the parent MockWatcher instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *WatcherCloseFunc) PushHook(hook func()) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *WatcherCloseFunc) SetDefaultReturn() {
	f.SetDefaultHook(func() {
		return
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *WatcherCloseFunc) PushReturn() {
	f.PushHook(func() {
		return
	})
}

func (f *WatcherCloseFunc) nextHook() func() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *WatcherCloseFunc) appendCall(r0 WatcherCloseFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of WatcherCloseFuncCall objects describing the
// invocations of this function.
func (f *WatcherCloseFunc) History() []WatcherCloseFuncCall {
	f.mutex.Lock()
	history := make([]WatcherCloseFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// WatcherCloseFuncCall is an object that describes an invocation of method
// Close on an instance of MockWatcher.
type WatcherCloseFuncCall struct{}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c WatcherCloseFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c WatcherCloseFuncCall) Results() []interface{} {
	return []interface{}{}
}

// WatcherStartFunc describes the behavior when the Start method of the
// parent MockWatcher instance is invoked.
type WatcherStartFunc struct {
	defaultHook func() error
	hooks       []func() error
	history     []WatcherStartFuncCall
	mutex       sync.Mutex
}

// Start delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockWatcher) Start() error {
	r0 := m.StartFunc.nextHook()()
	m.StartFunc.appendCall(WatcherStartFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Start method of the
// parent MockWatcher instance is invoked and the hook queue is empty.
func (f *WatcherStartFunc) SetDefaultHook(hook func() error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Start method of the parent MockWatcher instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *WatcherStartFunc) PushHook(hook func() error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *WatcherStartFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func() error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *WatcherStartFunc) PushReturn(r0 error) {
	f.PushHook(func() error {
		return r0
	})
}

func (f *WatcherStartFunc) nextHook() func() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *WatcherStartFunc) appendCall(r0 WatcherStartFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of WatcherStartFuncCall objects describing the
// invocations of this function.
func (f *WatcherStartFunc) History() []WatcherStartFuncCall {
	f.mutex.Lock()
	history := make([]WatcherStartFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// WatcherStartFuncCall is an object that describes an invocation of method
// Start on an instance of MockWatcher.
type WatcherStartFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c WatcherStartFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c WatcherStartFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// WatcherSubscribeFunc describes the behavior when the Subscribe method of
// the parent MockWatcher instance is invoked.
type WatcherSubscribeFunc struct {
	defaultHook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error)
	hooks       []func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error)
	history     []WatcherSubscribeFuncCall
	mutex       sync.Mutex
}

// Subscribe delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockWatcher) Subscribe(v0 context.Context, v1 *types.Transaction, v2 *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
	r0, r1 := m.SubscribeFunc.nextHook()(v0, v1, v2)
	m.SubscribeFunc.appendCall(WatcherSubscribeFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Subscribe method of
// the parent MockWatcher instance is invoked and the hook queue is empty.
func (f *WatcherSubscribeFunc) SetDefaultHook(hook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Subscribe method of the parent MockWatcher instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *WatcherSubscribeFunc) PushHook(hook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *WatcherSubscribeFunc) SetDefaultReturn(r0 transaction.ReceiptResponse, r1 error) {
	f.SetDefaultHook(func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *WatcherSubscribeFunc) PushReturn(r0 transaction.ReceiptResponse, r1 error) {
	f.PushHook(func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
		return r0, r1
	})
}

func (f *WatcherSubscribeFunc) nextHook() func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (transaction.ReceiptResponse, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *WatcherSubscribeFunc) appendCall(r0 WatcherSubscribeFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of WatcherSubscribeFuncCall objects describing
// the invocations of this function.
func (f *WatcherSubscribeFunc) History() []WatcherSubscribeFuncCall {
	f.mutex.Lock()
	history := make([]WatcherSubscribeFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// WatcherSubscribeFuncCall is an object that describes an invocation of
// method Subscribe on an instance of MockWatcher.
type WatcherSubscribeFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *types.Transaction
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 *transaction.SubscribeOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 transaction.ReceiptResponse
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c WatcherSubscribeFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c WatcherSubscribeFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// WatcherSubscribeAndWaitFunc describes the behavior when the
// SubscribeAndWait method of the parent MockWatcher instance is invoked.
type WatcherSubscribeAndWaitFunc struct {
	defaultHook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error)
	hooks       []func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error)
	history     []WatcherSubscribeAndWaitFuncCall
	mutex       sync.Mutex
}

// SubscribeAndWait delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockWatcher) SubscribeAndWait(v0 context.Context, v1 *types.Transaction, v2 *transaction.SubscribeOptions) (*types.Receipt, error) {
	r0, r1 := m.SubscribeAndWaitFunc.nextHook()(v0, v1, v2)
	m.SubscribeAndWaitFunc.appendCall(WatcherSubscribeAndWaitFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the SubscribeAndWait
// method of the parent MockWatcher instance is invoked and the hook queue
// is empty.
func (f *WatcherSubscribeAndWaitFunc) SetDefaultHook(hook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// SubscribeAndWait method of the parent MockWatcher instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *WatcherSubscribeAndWaitFunc) PushHook(hook func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *WatcherSubscribeAndWaitFunc) SetDefaultReturn(r0 *types.Receipt, r1 error) {
	f.SetDefaultHook(func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *WatcherSubscribeAndWaitFunc) PushReturn(r0 *types.Receipt, r1 error) {
	f.PushHook(func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error) {
		return r0, r1
	})
}

func (f *WatcherSubscribeAndWaitFunc) nextHook() func(context.Context, *types.Transaction, *transaction.SubscribeOptions) (*types.Receipt, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *WatcherSubscribeAndWaitFunc) appendCall(r0 WatcherSubscribeAndWaitFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of WatcherSubscribeAndWaitFuncCall objects
// describing the invocations of this function.
func (f *WatcherSubscribeAndWaitFunc) History() []WatcherSubscribeAndWaitFuncCall {
	f.mutex.Lock()
	history := make([]WatcherSubscribeAndWaitFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// WatcherSubscribeAndWaitFuncCall is an object that describes an invocation
// of method SubscribeAndWait on an instance of MockWatcher.
type WatcherSubscribeAndWaitFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *types.Transaction
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 *transaction.SubscribeOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *types.Receipt
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c WatcherSubscribeAndWaitFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c WatcherSubscribeAndWaitFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// WatcherWaitFunc describes the behavior when the Wait method of the parent
// MockWatcher instance is invoked.
type WatcherWaitFunc struct {
	defaultHook func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error)
	hooks       []func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error)
	history     []WatcherWaitFuncCall
	mutex       sync.Mutex
}

// Wait delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockWatcher) Wait(v0 context.Context, v1 transaction.ReceiptResponse) (*types.Receipt, error) {
	r0, r1 := m.WaitFunc.nextHook()(v0, v1)
	m.WaitFunc.appendCall(WatcherWaitFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Wait method of the
// parent MockWatcher instance is invoked and the hook queue is empty.
func (f *WatcherWaitFunc) SetDefaultHook(hook func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Wait method of the parent MockWatcher instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *WatcherWaitFunc) PushHook(hook func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *WatcherWaitFunc) SetDefaultReturn(r0 *types.Receipt, r1 error) {
	f.SetDefaultHook(func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *WatcherWaitFunc) PushReturn(r0 *types.Receipt, r1 error) {
	f.PushHook(func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error) {
		return r0, r1
	})
}

func (f *WatcherWaitFunc) nextHook() func(context.Context, transaction.ReceiptResponse) (*types.Receipt, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *WatcherWaitFunc) appendCall(r0 WatcherWaitFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of WatcherWaitFuncCall objects describing the
// invocations of this function.
func (f *WatcherWaitFunc) History() []WatcherWaitFuncCall {
	f.mutex.Lock()
	history := make([]WatcherWaitFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// WatcherWaitFuncCall is an object that describes an invocation of method
// Wait on an instance of MockWatcher.
type WatcherWaitFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 transaction.ReceiptResponse
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *types.Receipt
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c WatcherWaitFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c WatcherWaitFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
