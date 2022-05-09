// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "contracts/utils/ImmutableAuth.sol";
import "contracts/utils/MerkleTree.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import "hardhat/console.sol";
import {BridgePoolErrorCodes} from "contracts/libraries/errorCodes/BridgePoolErrorCodes.sol";
import "contracts/libraries/parsers/MerkleProofParserLibrary.sol";
import "contracts/libraries/MerkleProofLibrary.sol";

/// @custom:salt BridgePool
/// @custom:deploy-type deployStatic
contract BridgePool is Initializable, ImmutableFactory, MerkleTree {
    address internal _erc20TokenContract;
    address internal _bTokenContract;
    struct UTXO {
        uint32 chainID;
        address owner;
        uint256 value;
        uint256 fee;
        bytes32 txHash;
    }
    event Deposit(address from, uint256 value);

    constructor(address erc20TokenContract_, address bTokenContract_)
        public
        ImmutableFactory(msg.sender)
    {
        _erc20TokenContract = erc20TokenContract_;
        _bTokenContract = bTokenContract_;
    }

    function initialize(address erc20TokenContract_, address bTokenContract_)
        public
        onlyFactory
        initializer
    {
        _erc20TokenContract = erc20TokenContract_;
        _bTokenContract = bTokenContract_;
    }

    function deposit(
        uint8 accountType_,
        address aliceNetAddress_,
        uint256 erc20Amount_,
        uint256 bTokenAmount_,
        address ethAddress_
    ) public onlyFactory {
        ERC20(_erc20TokenContract).transferFrom(ethAddress_, address(this), erc20Amount_);
        ERC20(_bTokenContract).transferFrom(ethAddress_, address(this), bTokenAmount_);
        BToken(_bTokenContract).burnTo(address(this), bTokenAmount_, 0);
        emit Deposit(aliceNetAddress_, erc20Amount_);
    }

    function withdraw(
        bytes32 merkleRoot,
        bytes32 merkleKeyHash,
        bytes32 merkleValueHash,
        bytes memory encodedBurnedUTXO,
        bytes32[] memory auditPath,
        address receiver
    ) public onlyFactory {
        UTXO memory burnedUTXO = abi.decode(encodedBurnedUTXO, (UTXO));
        uint16 trieHeight = 256;
        bytes1 diff = bytes2(trieHeight - uint16(auditPath.length))[1];
        bytes32 leafHash = keccak256(abi.encodePacked(merkleKeyHash, merkleValueHash, diff));
        require(
            burnedUTXO.owner == receiver,
            string(
                abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER)
            )
        );
        require(
            merkleRoot == _calculateRootFromMerkleProof(auditPath, 0, merkleKeyHash, leafHash),
            string(abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_PROOF_OF_BURN_NOT_VERIFIED))
        );
        ERC20(_erc20TokenContract).approve(address(this), burnedUTXO.value);
        ERC20(_erc20TokenContract).transferFrom(address(this), receiver, burnedUTXO.value);
    }

    using MerkleProofParserLibrary for bytes;
    using MerkleProofLibrary for MerkleProofParserLibrary.MerkleProof;

    function withdrawWithBinaryProof(
        bytes memory encodedMerkleProof,
        bytes memory encodedBurnedUTXO,
        bytes32 stateRoot,
        address receiver
    ) public onlyFactory {
        MerkleProofParserLibrary.MerkleProof memory merkleProof = encodedMerkleProof
            .extractMerkleProof();
        UTXO memory burnedUTXO = abi.decode(encodedBurnedUTXO, (UTXO));
        bytes32 leafHash = merkleProof.computeLeafHash2();
        require(
            burnedUTXO.owner == receiver,
            string(
                abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER)
            )
        );
        require(
            merkleProof.verifyInclusion(stateRoot),
            string(abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_PROOF_OF_BURN_NOT_VERIFIED))
        );

        ERC20(_erc20TokenContract).approve(address(this), burnedUTXO.value);
        ERC20(_erc20TokenContract).transferFrom(address(this), receiver, burnedUTXO.value);
    }

    receive() external payable {}
}
