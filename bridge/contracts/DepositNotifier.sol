// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt DepositNotifier
/// @custom:deploy-type deployUpgradeable
contract DepositNotifier is ImmutableFactory {
    uint256 internal _nonce;
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

    function doEmit(bytes32 salt, address ercContract, uint256 number, address owner) onlyAllowed(salt) public {
        uint256 n = _nonce + 1;
        emit Deposited(n, ercContract, owner, number, _networkId);
        _nonce = n;
    }

    modifier onlyAllowed(bytes32 salt) {
        address expected = getMetamorphicContractAddress(keccak256(abi.encodePacked(salt, "ERC20")), _factoryAddress());
        require(msg.sender == expected, "not allowed");
        _;
    }
}

/// @custom:salt BridgeUSDT
/// @custom:auth-salt ERC20
/// @custom:deploy-type deployUpgradeable
contract BridgeUSDT {
    // final salt = salt + authSalt == BridgeUSDTERC20
}