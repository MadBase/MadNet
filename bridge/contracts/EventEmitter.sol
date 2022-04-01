// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import "hardhat/console.sol";

contract EventEmitter is Initializable {
    event BridgePoolDepositReceived(uint256 depositID_, address to_, uint256 amount_);
    event Generic(string title, address to_, uint256 amount_);
    event GenericEvent(bytes calldata_);

    fallback() external {
        // emit Delegated(msg.data);
        _emitEvent();
        console.log("fallback");
        console.logBytes(msg.data);
    }

    function _emitEvent() internal {
        console.log("emitEvent");
        emit Generic("ppp", address(this), 0);
    }

    function emitGenericEvent(bytes memory encodedEvent) public {
        emit GenericEvent(encodedEvent);
    }

    function emitBridgePoolDepositReceived(
        uint256 depositID_,
        address to_,
        uint256 amount_
    ) public {
        console.logBytes(msg.data);

        emit BridgePoolDepositReceived(depositID_, to_, amount_);
    }

    function emitGeneric(
        string calldata title,
        address to_,
        uint256 amount_
    ) public {
        console.logBytes(msg.data);

        emit Generic(title, to_, amount_);
    }

    function initialize() public initializer {}
}
