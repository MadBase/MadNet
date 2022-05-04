// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import "hardhat/console.sol";

/// @custom:salt BridgePool
/// @custom:deploy-type deployStatic
contract BridgePool is Initializable {
    address internal _erc20TokenContract;
    address internal _bTokenContract;
    uint256 internal _depositID;

    function initialize(address erc20TokenContract_, address bTokenContract_) public initializer {
        _erc20TokenContract = erc20TokenContract_;
        _bTokenContract = bTokenContract_;
    }

    constructor(address erc20TokenContract_, address bTokenContract_) public initializer {
        _erc20TokenContract = erc20TokenContract_;
        _bTokenContract = bTokenContract_;
    }

    struct UTXO {
        uint32 chainID;
        address owner;
        uint256 value;
        uint256 fee;
        bytes32 txHash;
    }
    event Deposit(address from, uint256 value);

    function deposit(
        uint8 accountType_,
        address aliceNetAddress_,
        uint256 erc20Amount_,
        uint256 bTokenAmount_,
        address ethAddress_
    ) public {
        ERC20(_erc20TokenContract).transferFrom(ethAddress_, address(this), erc20Amount_);
        ERC20(_bTokenContract).transferFrom(ethAddress_, address(this), bTokenAmount_);
        uint256 eths = BToken(_bTokenContract).burnTo(address(this), bTokenAmount_, 0);
        emit Deposit(aliceNetAddress_, erc20Amount_);
    }

    function withdraw(
        bytes32 merkleRoot,
        bytes32 merkleKeyHash,
        bytes32 merkleValueHash,
        bytes memory encodedBurnedUTXO,
        bytes32[] memory auditPath,
        address to
    ) public {
        UTXO memory burnedUTXO = abi.decode(encodedBurnedUTXO, (UTXO));
        uint16 trieHeight = 256;
        bytes1 diff = bytes2(trieHeight - uint16(auditPath.length))[1];
        bytes32 leafHash = keccak256(abi.encodePacked(merkleKeyHash, merkleValueHash, diff));
        require(
            burnedUTXO.owner == to,
            "BridgePool: deposit can only be requested for the owner in burned UTXO"
        );
        require(
            merkleRoot == verifyInclusion(auditPath, 0, merkleKeyHash, leafHash),
            "BridgePool: Proof of burn in aliceNet could not be verified"
        );
        ERC20(_erc20TokenContract).approve(address(this), burnedUTXO.value);
        ERC20(_erc20TokenContract).transferFrom(address(this), to, burnedUTXO.value);
    }

    receive() external payable {}

    // verifyInclusion returns the merkle root by hashing the merkle proof items
    function verifyInclusion(
        bytes32[] memory auditPath,
        uint16 merkleKeyIndex,
        bytes32 merkleKeyHash,
        bytes32 merkleLeafHash
    ) public returns (bytes32) {
        if (merkleKeyIndex == auditPath.length) {
            return merkleLeafHash;
        }
        if (_bitIsSet(merkleKeyHash, merkleKeyIndex)) {
            return
                keccak256(
                    abi.encodePacked(
                        auditPath[auditPath.length - merkleKeyIndex - 1],
                        verifyInclusion(
                            auditPath,
                            merkleKeyIndex + 1,
                            merkleKeyHash,
                            merkleLeafHash
                        )
                    )
                );
        }
        return
            keccak256(
                abi.encodePacked(
                    verifyInclusion(auditPath, merkleKeyIndex + 1, merkleKeyHash, merkleLeafHash),
                    auditPath[auditPath.length - merkleKeyIndex - 1]
                )
            );
    }

    function _bitIsSet(bytes32 bits, uint256 merkleKeyIndex) internal returns (bool) {
        return (bits[merkleKeyIndex / 8] & (bytes1(0x01) << (7 - (merkleKeyIndex % 8))) != 0);
    }
}
