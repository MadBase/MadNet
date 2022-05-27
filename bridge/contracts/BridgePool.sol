// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "contracts/utils/ImmutableAuth.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "contracts/BToken.sol";
import "hardhat/console.sol";
import {BridgePoolErrorCodes} from "contracts/libraries/errorCodes/BridgePoolErrorCodes.sol";
import "contracts/libraries/parsers/MerkleProofParserLibrary.sol";
import "contracts/libraries/MerkleProofLibrary.sol";
import "contracts/Snapshots.sol";
import "contracts/libraries/parsers/BClaimsParserLibrary.sol";
import "contracts/utils/ERC20SafeTransfer.sol";

/// @custom:salt BridgePool
/// @custom:deploy-type deployStatic
contract BridgePool is Initializable, ImmutableSnapshots, ERC20SafeTransfer {
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

    //TODO: remove and use DepositNotifier
    event Deposited(
        uint256 nonce,
        address ercContract,
        address owner,
        uint256 number,
        uint256 networkId
    );

    constructor(address erc20TokenContract_, address bTokenContract_) ImmutableFactory(msg.sender) {
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
        uint256 bTokenAmount_
    ) public {
        require(
            ERC20(_erc20TokenContract).transferFrom(msg.sender, address(this), erc20Amount_),
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
        BToken(_bTokenContract).burnTo(address(this), bTokenAmount_, 0);
        //TODO: remove and use DepositNotifier
        emit Deposited(accountType_, _erc20TokenContract, aliceNetAddress_, erc20Amount_, 0);
        // Uncomment upon merging of PR-126
        // DepositNotifier(_depositNotifierAddress()).doEmit(
        //     _saltForBridgePool(),
        //     _erc20TokenContract,
        //     erc20Amount_,i
        //     aliceNetAddress_
        // );
    }

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
        _safeTransferERC20(IERC20Transferable(_erc20TokenContract), msg.sender, burnedUTXO.value);
    }

    receive() external payable {}
}
