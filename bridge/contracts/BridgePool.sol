// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "contracts/utils/ImmutableAuth.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import "hardhat/console.sol";
import {BridgePoolErrorCodes} from "contracts/libraries/errorCodes/BridgePoolErrorCodes.sol";
import "contracts/libraries/parsers/MerkleProofParserLibrary.sol";
import "contracts/libraries/MerkleProofLibrary.sol";
import "contracts/DepositNotifier.sol";

/// @custom:salt BridgePool
/// @custom:deploy-type deployStatic
contract BridgePool is
    Initializable,
    ImmutableFactory,
    ImmutableBridgePool,
    ImmutableDepositNotifier
{
    address internal immutable _erc20TokenContract;
    address internal immutable _bTokenContract;
    using MerkleProofParserLibrary for bytes;
    using MerkleProofLibrary for MerkleProofParserLibrary.MerkleProof;

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
    {}

    function deposit(
        uint8 accountType_,
        address aliceNetAddress_,
        uint256 erc20Amount_,
        uint256 bTokenAmount_,
        address from_
    ) public onlyFactory {
        ERC20(_erc20TokenContract).transferFrom(from_, address(this), erc20Amount_);
        ERC20(_bTokenContract).transferFrom(from_, address(this), bTokenAmount_);
        BToken(_bTokenContract).burnTo(address(this), bTokenAmount_, 0);
        DepositNotifier(_depositNotifierAddress()).doEmit(
            _saltForBridgePool(),
            _erc20TokenContract,
            erc20Amount_,
            from_
        );
        // emit Deposit(aliceNetAddress_, erc20Amount_);
    }

    function withdraw(
        bytes memory encodedMerkleProof,
        bytes memory encodedBurnedUTXO,
        bytes32 stateRoot,
        address receiver
    ) public onlyFactory {
        MerkleProofParserLibrary.MerkleProof memory merkleProof = encodedMerkleProof
            .extractMerkleProof();
        UTXO memory burnedUTXO = abi.decode(encodedBurnedUTXO, (UTXO));
        bytes32 leafHash = merkleProof.computeLeafHash();
        require(
            burnedUTXO.owner == receiver,
            string(
                abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER)
            )
        );
        require(
            merkleProof.checkProof(stateRoot, leafHash),
            string(abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_PROOF_OF_BURN_NOT_VERIFIED))
        );
        ERC20(_erc20TokenContract).approve(address(this), burnedUTXO.value);
        ERC20(_erc20TokenContract).transferFrom(address(this), receiver, burnedUTXO.value);
    }

    receive() external payable {}
}
