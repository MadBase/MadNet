// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BridgePoolDepositNotifierMetaData contains all meta data concerning the BridgePoolDepositNotifier contract.
var BridgePoolDepositNotifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"networkId_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"ercContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"networkId\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"ercContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"doEmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_salt\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"getMetamorphicContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// BridgePoolDepositNotifierABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgePoolDepositNotifierMetaData.ABI instead.
var BridgePoolDepositNotifierABI = BridgePoolDepositNotifierMetaData.ABI

// BridgePoolDepositNotifier is an auto generated Go binding around an Ethereum contract.
type BridgePoolDepositNotifier struct {
	BridgePoolDepositNotifierCaller     // Read-only binding to the contract
	BridgePoolDepositNotifierTransactor // Write-only binding to the contract
	BridgePoolDepositNotifierFilterer   // Log filterer for contract events
}

// BridgePoolDepositNotifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgePoolDepositNotifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgePoolDepositNotifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgePoolDepositNotifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgePoolDepositNotifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgePoolDepositNotifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgePoolDepositNotifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgePoolDepositNotifierSession struct {
	Contract     *BridgePoolDepositNotifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// BridgePoolDepositNotifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgePoolDepositNotifierCallerSession struct {
	Contract *BridgePoolDepositNotifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// BridgePoolDepositNotifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgePoolDepositNotifierTransactorSession struct {
	Contract     *BridgePoolDepositNotifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// BridgePoolDepositNotifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgePoolDepositNotifierRaw struct {
	Contract *BridgePoolDepositNotifier // Generic contract binding to access the raw methods on
}

// BridgePoolDepositNotifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgePoolDepositNotifierCallerRaw struct {
	Contract *BridgePoolDepositNotifierCaller // Generic read-only contract binding to access the raw methods on
}

// BridgePoolDepositNotifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgePoolDepositNotifierTransactorRaw struct {
	Contract *BridgePoolDepositNotifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgePoolDepositNotifier creates a new instance of BridgePoolDepositNotifier, bound to a specific deployed contract.
func NewBridgePoolDepositNotifier(address common.Address, backend bind.ContractBackend) (*BridgePoolDepositNotifier, error) {
	contract, err := bindBridgePoolDepositNotifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgePoolDepositNotifier{BridgePoolDepositNotifierCaller: BridgePoolDepositNotifierCaller{contract: contract}, BridgePoolDepositNotifierTransactor: BridgePoolDepositNotifierTransactor{contract: contract}, BridgePoolDepositNotifierFilterer: BridgePoolDepositNotifierFilterer{contract: contract}}, nil
}

// NewBridgePoolDepositNotifierCaller creates a new read-only instance of BridgePoolDepositNotifier, bound to a specific deployed contract.
func NewBridgePoolDepositNotifierCaller(address common.Address, caller bind.ContractCaller) (*BridgePoolDepositNotifierCaller, error) {
	contract, err := bindBridgePoolDepositNotifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgePoolDepositNotifierCaller{contract: contract}, nil
}

// NewBridgePoolDepositNotifierTransactor creates a new write-only instance of BridgePoolDepositNotifier, bound to a specific deployed contract.
func NewBridgePoolDepositNotifierTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgePoolDepositNotifierTransactor, error) {
	contract, err := bindBridgePoolDepositNotifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgePoolDepositNotifierTransactor{contract: contract}, nil
}

// NewBridgePoolDepositNotifierFilterer creates a new log filterer instance of BridgePoolDepositNotifier, bound to a specific deployed contract.
func NewBridgePoolDepositNotifierFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgePoolDepositNotifierFilterer, error) {
	contract, err := bindBridgePoolDepositNotifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgePoolDepositNotifierFilterer{contract: contract}, nil
}

// bindBridgePoolDepositNotifier binds a generic wrapper to an already deployed contract.
func bindBridgePoolDepositNotifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgePoolDepositNotifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgePoolDepositNotifier.Contract.BridgePoolDepositNotifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.BridgePoolDepositNotifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.BridgePoolDepositNotifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgePoolDepositNotifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.contract.Transact(opts, method, params...)
}

// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
//
// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierCaller) GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error) {
	var out []interface{}
	err := _BridgePoolDepositNotifier.contract.Call(opts, &out, "getMetamorphicContractAddress", _salt, _factory)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
//
// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierSession) GetMetamorphicContractAddress(_salt [32]byte, _factory common.Address) (common.Address, error) {
	return _BridgePoolDepositNotifier.Contract.GetMetamorphicContractAddress(&_BridgePoolDepositNotifier.CallOpts, _salt, _factory)
}

// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
//
// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierCallerSession) GetMetamorphicContractAddress(_salt [32]byte, _factory common.Address) (common.Address, error) {
	return _BridgePoolDepositNotifier.Contract.GetMetamorphicContractAddress(&_BridgePoolDepositNotifier.CallOpts, _salt, _factory)
}

// DoEmit is a paid mutator transaction binding the contract method 0xc58b0a4f.
//
// Solidity: function doEmit(bytes32 salt, address ercContract, uint256 number, address owner) returns()
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierTransactor) DoEmit(opts *bind.TransactOpts, salt [32]byte, ercContract common.Address, number *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.contract.Transact(opts, "doEmit", salt, ercContract, number, owner)
}

// DoEmit is a paid mutator transaction binding the contract method 0xc58b0a4f.
//
// Solidity: function doEmit(bytes32 salt, address ercContract, uint256 number, address owner) returns()
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierSession) DoEmit(salt [32]byte, ercContract common.Address, number *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.DoEmit(&_BridgePoolDepositNotifier.TransactOpts, salt, ercContract, number, owner)
}

// DoEmit is a paid mutator transaction binding the contract method 0xc58b0a4f.
//
// Solidity: function doEmit(bytes32 salt, address ercContract, uint256 number, address owner) returns()
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierTransactorSession) DoEmit(salt [32]byte, ercContract common.Address, number *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BridgePoolDepositNotifier.Contract.DoEmit(&_BridgePoolDepositNotifier.TransactOpts, salt, ercContract, number, owner)
}

// BridgePoolDepositNotifierDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the BridgePoolDepositNotifier contract.
type BridgePoolDepositNotifierDepositedIterator struct {
	Event *BridgePoolDepositNotifierDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgePoolDepositNotifierDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgePoolDepositNotifierDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgePoolDepositNotifierDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgePoolDepositNotifierDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgePoolDepositNotifierDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgePoolDepositNotifierDeposited represents a Deposited event raised by the BridgePoolDepositNotifier contract.
type BridgePoolDepositNotifierDeposited struct {
	Nonce       *big.Int
	ErcContract common.Address
	Owner       common.Address
	Number      *big.Int
	NetworkId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0xefd7bb4f992963087670aa168590baee24d227ff5c18ae3790ef7ec22bde6274.
//
// Solidity: event Deposited(uint256 nonce, address ercContract, address owner, uint256 number, uint256 networkId)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierFilterer) FilterDeposited(opts *bind.FilterOpts) (*BridgePoolDepositNotifierDepositedIterator, error) {

	logs, sub, err := _BridgePoolDepositNotifier.contract.FilterLogs(opts, "Deposited")
	if err != nil {
		return nil, err
	}
	return &BridgePoolDepositNotifierDepositedIterator{contract: _BridgePoolDepositNotifier.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0xefd7bb4f992963087670aa168590baee24d227ff5c18ae3790ef7ec22bde6274.
//
// Solidity: event Deposited(uint256 nonce, address ercContract, address owner, uint256 number, uint256 networkId)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *BridgePoolDepositNotifierDeposited) (event.Subscription, error) {

	logs, sub, err := _BridgePoolDepositNotifier.contract.WatchLogs(opts, "Deposited")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgePoolDepositNotifierDeposited)
				if err := _BridgePoolDepositNotifier.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0xefd7bb4f992963087670aa168590baee24d227ff5c18ae3790ef7ec22bde6274.
//
// Solidity: event Deposited(uint256 nonce, address ercContract, address owner, uint256 number, uint256 networkId)
func (_BridgePoolDepositNotifier *BridgePoolDepositNotifierFilterer) ParseDeposited(log types.Log) (*BridgePoolDepositNotifierDeposited, error) {
	event := new(BridgePoolDepositNotifierDeposited)
	if err := _BridgePoolDepositNotifier.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
