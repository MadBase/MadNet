// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// IBridgePoolDepositNotifierCaller ...
type IBridgePoolDepositNotifierCaller interface {
	// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
	//
	// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
	GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error)
}