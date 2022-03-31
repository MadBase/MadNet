// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/MadByte.sol";
import "contracts/EventEmitter.sol";
import "hardhat/console.sol";

contract BridgePool1 is Initializable {
    struct Deposit {
        address account;
        uint256 eths;
        bool proofOfBurn;
    }
    mapping(uint256 => Deposit) internal _deposits;
    address internal _eventEmitter;
    address internal _madByte;
    uint256 internal _depositID;

    function initialize(address madByte_, address eventEmitter_) public initializer {
        _eventEmitter = eventEmitter_;
        _madByte = madByte_;
    }

    function deposit(address to_, uint256 amount_) public returns (uint256) {
        return _deposit(to_, amount_);
    }

    function withdraw(uint256 depositID_) public {
        _withdraw(depositID_);
    }

    function confirmProofOfBurn(uint256 depositID_) public {
        //TODO: set access control. Probably onlyETHDKG
        _deposits[depositID_].proofOfBurn = true;
    }

    function distribute(uint256 depositID_) public {
        _distribute(depositID_);
    }

    function _deposit(address to_, uint256 amount_) internal returns (uint256) {
        ERC20(_madByte).transferFrom(msg.sender, address(this), amount_);
        uint256 eths = MadByte(_madByte).burn(amount_, 0);
        console.log("Eths for Madbytes", eths);
        require(eths > 0, "BridgePool: Could not burn tokens for deposit");
        uint256 depositID = _createDeposit(to_, eths, false);
        _emitDepositEvent(depositID, to_, amount_);
        return depositID;
    }

    function _createDeposit(
        address to_,
        uint256 amount_,
        bool proofOfBurn_
    ) internal returns (uint256) {
        uint256 depositID = _depositID + 1;
        _deposits[depositID] = Deposit(to_, amount_, proofOfBurn_);
        _depositID = depositID;
        return depositID;
    }

    function _withdraw(uint256 depositID_) internal {
        require(
            _deposits[depositID_].account == msg.sender,
            "BridgePool: Only owner of deposit can withdraw"
        );
        _emitWithdrawalRequestEvent(depositID_);
    }

    function _distribute(uint256 depositID_) internal {
        require(
            _deposits[depositID_].account != address(0),
            "BridgePool: Deposit with provided Id does not exist"
        );
        require(
            _deposits[depositID_].account == msg.sender,
            "BridgePool: Only owner of deposit can withdraw"
        );
        require(
            _deposits[depositID_].proofOfBurn == true,
            "BridgePool: No proof of Burn confirmed for this deposit. Can't distribute yet."
        );
        MadByte(_madByte).mintTo{value: _deposits[depositID_].eths}(
            _deposits[depositID_].account,
            0
        );
        _emitDistributeEvent(depositID_, _deposits[depositID_].account, _deposits[depositID_].eths);
        delete _deposits[depositID_];
    }

    function _emitDepositEvent(
        uint256 depositID_,
        address to_,
        uint256 amount_
    ) internal {
        bytes memory encodedEvent = abi.encode("BridgePool", "deposit", depositID_, to_, amount_);
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    function _emitWithdrawalRequestEvent(uint256 depositID_) internal {
        bytes memory encodedEvent = abi.encode("BridgePool", "withdrawal", depositID_);
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    function _emitDistributeEvent(
        uint256 depositID_,
        address to_,
        uint256 amount_
    ) internal {
        bytes memory encodedEvent = abi.encode(
            "BridgePool",
            "distribute",
            depositID_,
            to_,
            amount_
        );
        EventEmitter(_eventEmitter).emitGenericEvent(encodedEvent);
    }

    receive() external payable {}
}
