// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/MadByte.sol";
import "contracts/EventEmitter.sol";
import "hardhat/console.sol";

contract BridgePool is Initializable {
    mapping(address => uint256) internal _depositsETH;
    mapping(address => uint256) internal _stateProofs;
    address internal _eventEmitter;
    address internal _madByte;

    function initialize(address madByte_, address eventEmitter_) public initializer {
        _eventEmitter = eventEmitter_;
        _madByte = madByte_;
    }

    function deposit(uint256 amountMB_) public returns (uint256) {
        return _deposit(msg.sender, amountMB_);
    }

    function requestWithdrawal(uint256 depositID_) public {
        _requestWithdrawal(msg.sender, depositID_);
    }

    function confirmProofOfBurn(address account_, uint256 amountMB_) public {
        //TODO: set access control. Probably onlyETHDKG
        _stateProofs[account_] += amountMB_;
    }

    function withdraw(uint256 amountMB_) public {
        _withdraw(msg.sender, amountMB_);
    }

    function _deposit(address account_, uint256 amountMB_) internal returns (uint256) {
        require(amountMB_ > 0, "BridgePool: Deposit amountMB must be greater that 0");
        ERC20(_madByte).transferFrom(account_, address(this), amountMB_);
        uint256 amountETH = MadByte(_madByte).burn(amountMB_, 0);
        _depositsETH[account_] += amountETH;
        require(amountETH > 0, "BridgePool: Could not burn tokens for deposit");
        _emitDepositEvent(account_, amountMB_);
        return 0;
    }

    function _requestWithdrawal(address account_, uint256 amountMB_) internal {
        _emitWithdrawalRequestEvent(account_, amountMB_);
    }

    function _withdraw(address account_, uint256 amountMB_) internal {
        require(amountMB_ > 0, "BridgePool: Withdrawal amountMB must be greater that 0");
        require(
            amountMB_ <= _stateProofs[account_],
            "BridgePool: Withdrawal amountMB greater than available balance"
        );
        _stateProofs[account_] -= amount;
        uint256 amountETH = _depositsETH[account_];
        uint256 madBytes = MadByte(_madByte).mintTo{value: amountETH}(account_, 0);
        _emitDistributeEvent(account_, madBytes);
        delete _depositsETH[account_];
    }

    function _emitDepositEvent(address account_, uint256 amountMB_) internal {
        bytes memory encodedEvent = abi.encode("BridgePool", "deposit", account_, amountMB_);
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    function _emitWithdrawalRequestEvent(address account_, uint256 amountMB_) internal {
        bytes memory encodedEvent = abi.encode(
            "BridgePool",
            "requestWithdrawal",
            account_,
            amountMB_
        );
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    function _emitDistributeEvent(address account_, uint256 amountMB_) internal {
        bytes memory encodedEvent = abi.encode("BridgePool", "withdrawal", account_, amountMB_);
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    receive() external payable {}
}
