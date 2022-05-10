// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

import "contracts/utils/ImmutableAuth.sol";
import "hardhat/console.sol";

/// @custom:salt DepositNotifier
/// @custom:deploy-type deployUpgradeable
contract DepositNotifier is ImmutableFactory {
    uint256 internal _nonce = 0;
    uint256 internal immutable _networkId;

    event Deposited(
        uint256 nonce,
        address ercContract,
        address owner,
        uint256 number, // If fungible, this is the amount. If non-fungible, this is the id
        uint256 networkId
    );

    constructor(uint256 id) ImmutableFactory(msg.sender) {
        _networkId = id;
    }

    function initialize(uint256 id) public onlyFactory {}

    function doEmit(
        bytes32 salt,
        address ercContract,
        uint256 number,
        address owner
    ) public onlyFactoryChildren(salt) {
        uint256 n = _nonce + 1;
        emit Deposited(n, ercContract, owner, number, _networkId);
        _nonce = n;
    }
}
