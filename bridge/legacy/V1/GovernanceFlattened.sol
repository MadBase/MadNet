// Sources flattened with hardhat v2.9.1 https://hardhat.org

// File contracts/interfaces/IGovernor.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

interface IGovernor {
    event ValueUpdated(
        uint256 indexed epoch,
        uint256 indexed key,
        bytes32 indexed value,
        address who
    );

    function updateValue(
        uint256 epoch,
        uint256 key,
        bytes32 value
    ) external;
}

// File contracts/libraries/errorCodes/GovernanceErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library GovernanceErrorCodes {
    // Governance error codes
    bytes32 public constant GOVERNANCE_ONLY_FACTORY_ALLOWED = "200"; //"Governance: Only factory allowed!"
}

// File contracts/Governance.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

/// @custom:salt Governance
/// @custom:deploy-type deployUpgradeable
contract Governance is IGovernor {
    // dummy contract
    address internal immutable _factory;

    constructor() {
        _factory = msg.sender;
    }

    function updateValue(
        uint256 epoch,
        uint256 key,
        bytes32 value
    ) external {
        require(
            msg.sender == _factory,
            string(abi.encodePacked(GovernanceErrorCodes.GOVERNANCE_ONLY_FACTORY_ALLOWED))
        );
        emit ValueUpdated(epoch, key, value, msg.sender);
    }
}
