// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "contracts/utils/ImmutableAuth.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import {BridgePoolErrorCodes} from "contracts/libraries/errorCodes/BridgePoolErrorCodes.sol";
import "contracts/libraries/parsers/MerkleProofParserLibrary.sol";
import "contracts/libraries/MerkleProofLibrary.sol";
import "contracts/Snapshots.sol";
import "contracts/libraries/parsers/BClaimsParserLibrary.sol";
import "contracts/utils/ERC20SafeTransfer.sol";
import "contracts/BridgePoolDepositNotifier.sol";
import "contracts/BridgePoolFactory.sol";
import "contracts/Foundation.sol";
import "contracts/utils/MagicEthTransfer.sol";
import "contracts/Proxy.sol";

import "hardhat/console.sol";

/// @custom:salt BridgePool
/// @custom:deploy-type deployStatic
contract BridgePool is
    Initializable,
    ImmutableSnapshots,
    ERC20SafeTransfer,
    ImmutableBridgePool,
    ImmutableBridgePoolDepositNotifier,
    ImmutableBridgePoolFactory,
    ImmutableBToken,
    MagicEthTransfer
{
    using MerkleProofParserLibrary for bytes;
    using MerkleProofLibrary for MerkleProofParserLibrary.MerkleProof;

    address internal immutable _ercTokenContract;
    address payable internal immutable _bTokenContract;

    struct UTXO {
        uint32 chainID;
        address owner;
        uint256 value;
        uint256 fee;
        bytes32 txHash;
    }

    constructor(address erc20TokenContract_, address payable bTokenContract_)
        ImmutableFactory(msg.sender)
    {
        _ercTokenContract = erc20TokenContract_;
        _bTokenContract = bTokenContract_;
    }

    function initialize(address erc20TokenContract_, address bTokenContract_)
        public
        onlyFactory
        initializer
    {}

    /// @notice Transfer tokens from sender and emit a "Deposited" event for minting correspondent tokens in sidechain
    /// @param accountType_ The type of account
    /// @param aliceNetAddress_ The address on the sidechain where to mint the tokens
    /// @param ercAmount_ The amount of ERC tokens to deposit
    /// @param bTokenAmount_ The fee for deposit in bTokens
    function deposit(
        uint8 accountType_,
        address aliceNetAddress_,
        uint256 ercAmount_,
        uint256 bTokenAmount_
    ) public {
        require(
            ERC20(_ercTokenContract).transferFrom(msg.sender, address(this), ercAmount_),
            string(
                abi.encodePacked(
                    BridgePoolErrorCodes.BRIDGEPOOL_COULD_NOT_TRANSFER_DEPOSIT_AMOUNT_FROM_SENDER
                )
            )
        );
        require(
            ERC20(_bTokenContract).transferFrom(msg.sender, address(this), bTokenAmount_),
            string(
                abi.encodePacked(
                    BridgePoolErrorCodes.BRIDGEPOOL_COULD_NOT_TRANSFER_DEPOSIT_FEE_FROM_SENDER
                )
            )
        );
        console.log("BToken ETH Balance before burn", address(_bTokenContract).balance);

        //        require(BToken(_bTokenContract).destroyTokens(bTokenAmount_), "unable to burn fees");

        uint256 ethFromBurning = BToken(_bTokenContract).burnTo(address(this), bTokenAmount_, 0);
        address addr;
        uint256 value;
        console.log("BridgePool ETH Balance before selfDestruct", address(this).balance);
        console.log("BToken ETH Balance before selfDestruct", address(_bTokenContract).balance);
        bytes memory deployCode = hex"6020595937595180ff";
        bytes memory encodedDeployCode = abi.encodePacked(_bTokenAddress(), deployCode);
        console.logBytes(encodedDeployCode);
        assembly {
            value := ethFromBurning
            //mstore(0x00, 0x000000000000000000000000000000000000000000006020595937595180fffe)
            mstore(0x00, 0x3e4b55722e041c0fc55a90b6eba50ea5d3f7c02a00006020595937595180ff)
            addr := create(value, 0x00, 0x20)
            if iszero(addr) {
                returndatacopy(0x00, 0x00, returndatasize())
                revert(0x00, returndatasize())
            }
        }
        console.log("Selfdestruct deployed at", addr);
        console.log("BridgePool ETH Balance after selfDestruct", address(this).balance);

        console.log("BToken ETH Balance after selfDesctruct", address(_bTokenContract).balance);
        BridgePoolDepositNotifier(_bridgePoolDepositNotifierAddress()).doEmit(
            BridgePoolFactory(_bridgePoolFactoryAddress()).getSaltFromERC20Address(
                _ercTokenContract
            ),
            _ercTokenContract,
            ercAmount_,
            msg.sender
        );
    }

    /// @notice Transfer funds to sender upon a verificable proof of burn in sidechain
    /// @param encodedMerkleProof The merkle proof
    /// @param encodedBurnedUTXO The burned UTXO
    function withdraw(bytes memory encodedMerkleProof, bytes memory encodedBurnedUTXO) public {
        BClaimsParserLibrary.BClaims memory bClaims = Snapshots(_snapshotsAddress())
            .getBlockClaimsFromLatestSnapshot();
        MerkleProofParserLibrary.MerkleProof memory merkleProof = encodedMerkleProof
            .extractMerkleProof();
        UTXO memory burnedUTXO = abi.decode(encodedBurnedUTXO, (UTXO));
        require(
            burnedUTXO.owner == msg.sender,
            string(
                abi.encodePacked(
                    BridgePoolErrorCodes.BRIDGEPOOL_RECEIVER_IS_NOT_OWNER_ON_PROOF_OF_BURN_UTXO
                )
            )
        );
        require(
            merkleProof.checkProof(bClaims.stateRoot, merkleProof.computeLeafHash()),
            string(abi.encodePacked(BridgePoolErrorCodes.BRIDGEPOOL_COULD_NOT_VERIFY_PROOF_OF_BURN))
        );
        _safeTransferERC20(IERC20Transferable(_ercTokenContract), msg.sender, burnedUTXO.value);
    }

    receive() external payable {}
}
