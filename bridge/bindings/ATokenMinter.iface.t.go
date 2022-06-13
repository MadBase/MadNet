// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// IATokenMinterTransactor ...
type IATokenMinterTransactor interface {
	// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
	//
	// Solidity: function mint(address to, uint256 amount) returns()
	Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)
}
