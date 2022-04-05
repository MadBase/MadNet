// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/interfaces/IAToken.sol";
import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt ATokenMinter
/// @custom:deploy-type deployUpgradeable
contract ATokenMinter is ImmutableAToken {
    constructor() ImmutableFactory(msg.sender) ImmutableAToken() {}

    function mint(address to, uint256 amount) public onlyFactory {
        IAToken(_aTokenAddress()).externalMint(to, amount);
    }
}
