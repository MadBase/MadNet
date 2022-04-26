// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/ATokenBurner.sol";
import "contracts/ATokenMinter.sol";
import "hardhat/console.sol";

contract BridgePool is Initializable {
    struct Deposit {
        address account;
        uint256 eths;
        bool proofOfBurn;
    }
    mapping(uint256 => Deposit) internal _deposits;
    address internal _aTokenContract;
    address internal _bTokenContract;
    uint256 internal _depositID;

    function initialize(address aTokenContract_, address bTokenContract_) public initializer {
        _aTokenContract = aTokenContract_;
        _bTokenContract = bTokenContract_;
    }

    function deposit(address to_, uint256 amount_) public returns (uint256) {
        return _deposit(to_, amount_);
    }

    function withdraw(uint256 depositID_) public {
        _withdraw(depositID_);
    }

    function confirmProofOfBurn(uint256 depositID_) public {
        //TODO: set access control.
        _deposits[depositID_].proofOfBurn = true;
    }

    function distribute(uint256 depositID_) public {
        _distribute(depositID_);
    }

    function _deposit(address to_, uint256 amount_) internal returns (uint256) {
        ERC20(_aTokenContract).transferFrom(msg.sender, address(this), amount_);
        ATokenBurner(_aTokenContract).burn(to_, 0);
        // console.log("Eths for Madbytes", eths);
        // require(eths > 0, "BridgePool: Could not burn tokens for deposit");
        uint256 depositID = _createDeposit(to_, 0, false);
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
        //{value: _deposits[depositID_].eths}
        ATokenMinter(_aTokenContract).mint(_deposits[depositID_].account, 0);
        delete _deposits[depositID_];
    }

    receive() external payable {}
}
