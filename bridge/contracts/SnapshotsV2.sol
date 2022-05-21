// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/interfaces/ISnapshots.sol";
import "contracts/interfaces/IValidatorPool.sol";
import "contracts/interfaces/IETHDKG.sol";
import "contracts/libraries/parsers/RCertParserLibrary.sol";
import "contracts/libraries/parsers/BClaimsParserLibrary.sol";
import "contracts/libraries/math/CryptoLibrary.sol";
import "contracts/libraries/snapshots/SnapshotsStorageV2.sol";
import "contracts/utils/DeterministicAddress.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {SnapshotsErrorCodes} from "contracts/libraries/errorCodes/SnapshotsErrorCodes.sol";
import "hardhat/console.sol";

/// @custom:salt Snapshots
/// @custom:deploy-type deployUpgradeable
contract SnapshotsV2 is Initializable, SnapshotsStorageV2, ISnapshots {
    constructor(uint256 chainID_, uint256 epochLength_) SnapshotsStorageV2(chainID_, epochLength_) {}

    function initialize(uint32 desperationDelay_, uint32 desperationFactor_)
        public
        onlyFactory
        initializer
    {   
        // considering that in optimum conditions 1 Sidechain block is at every 3 seconds and 1 block at
        // ethereum is approx at 13 seconds
        _minimumIntervalBetweenSnapshots = uint32(_epochLength / 4);
        _snapshotDesperationDelay = desperationDelay_;
        _snapshotDesperationFactor = desperationFactor_;
    }
    
    function integrate() public onlyFactory {
        require(getEpoch() == 0, "already integrated");
        if(_epoch > 3){
            _setSnapshot(_snapshots[_epoch-2]);
            _setSnapshot(_snapshots[_epoch-1]);
            _setSnapshot(_snapshots[_epoch]);
        }
    }

    function setSnapshotDesperationDelay(uint32 desperationDelay_) public onlyFactory {
        _snapshotDesperationDelay = desperationDelay_;
    }

    function setSnapshotDesperationFactor(uint32 desperationFactor_) public onlyFactory {
        _snapshotDesperationFactor = desperationFactor_;
    }

    function setMinimumIntervalBetweenSnapshots(uint32 minimumIntervalBetweenSnapshots_)
        public
        onlyFactory
    {
        _minimumIntervalBetweenSnapshots = minimumIntervalBetweenSnapshots_;
    }

    /// @notice Saves next snapshot
    /// @param groupSignature_ The group signature used to sign the snapshots' block claims
    /// @param bClaims_ The claims being made about given block
    /// @return Flag whether we should kick off another round of key generation
    function snapshot(bytes calldata groupSignature_, bytes calldata bClaims_)
        public
        returns (bool)
    {
        //require the sender to be a validator
        require(
            IValidatorPool(_validatorPoolAddress()).isValidator(msg.sender),
            string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_ONLY_VALIDATORS_ALLOWED))
        );
        //require that consensus is running
        require(
            IValidatorPool(_validatorPoolAddress()).isConsensusRunning(),
            string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_CONSENSUS_RUNNING))
        );
        uint32 epoch = __epoch._value;
        // console.log("epoch: %s, height: %s", epoch, blockClaims.height);
        //get the last snapshot 
        (, Snapshot memory lastSnapshot) = _getLastSnapshot();
        //TODO determine if we need to check 
        // require that the current block number is greater than 
        //blocknumber at last snapshot plus min interval between snaps
        require(
            block.number >= lastSnapshot.committedAt + _minimumIntervalBetweenSnapshots,
            string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_MIN_BLOCKS_INTERVAL_NOT_PASSED))
        );
        epoch = epoch + 1;
        // // TODO: BRING BACK AFTER GOLANG LOGIC IS DEBUGGED AND MERGED
        // {
        //     // Check if sender is the elected validator allowed to make the snapshot
        //     (bool success, uint256 validatorIndex) = IETHDKG(_ethdkgAddress())
        //         .tryGetParticipantIndex(msg.sender);
        //     require(success, "Snapshots: Caller didn't participate in the last ethdkg round!");

        //     uint256 ethBlocksSinceLastSnapshot = block.number - _snapshots[epoch - 1].committedAt;

        //     uint256 blocksSinceDesperation = ethBlocksSinceLastSnapshot >= _snapshotDesperationDelay
        //         ? ethBlocksSinceLastSnapshot - _snapshotDesperationDelay
        //         : 0;

        //     require(
        //         _mayValidatorSnapshot(
        //             IValidatorPool(_validatorPoolAddress()).getValidatorsCount(),
        //             validatorIndex - 1,
        //             blocksSinceDesperation,
        //             keccak256(bClaims_),
        //             uint256(_snapshotDesperationFactor)
        //         ),
        //         "Snapshots: Validator not elected to do snapshot!"
        //     );
        // }
        //this block verifies the signature of the claim matches the one on ethdkg
        {   //extract the masterPublicKey and signature from the group signature
            (uint256[4] memory masterPublicKey, uint256[2] memory signature) = RCertParserLibrary
                .extractSigGroup(groupSignature_, 0);
            //require the masterPubKey from group sig to match the one on ethdkg            
            require(
                keccak256(abi.encodePacked(masterPublicKey)) ==
                    IETHDKG(_ethdkgAddress()).getMasterPublicKeyHash(),
                string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_WRONG_MASTER_PUBLIC_KEY))
            );
            //verify the claim with the signature and master pub key
            require(
                CryptoLibrary.verifySignatureASM(
                    abi.encodePacked(keccak256(bClaims_)),
                    signature,
                    masterPublicKey
                ),
                string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_SIGNATURE_VERIFICATION_FAILED))
            );
        }

        BClaimsParserLibrary.BClaims memory blockClaims = BClaimsParserLibrary.extractBClaims(
            bClaims_
        );
        //require the new claim to have
        // console.log("epoch: %s, height: %s", epoch, blockClaims.height);
        require(
            epoch * _epochLength == blockClaims.height,
            string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_INCORRECT_BLOCK_HEIGHT))
        );

        require(
            blockClaims.chainId == _chainId,
            string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_INCORRECT_CHAIN_ID))
        );

        bool isSafeToProceedConsensus = true;
        if (IValidatorPool(_validatorPoolAddress()).isMaintenanceScheduled()) {
            isSafeToProceedConsensus = false;
            IValidatorPool(_validatorPoolAddress()).pauseConsensus();
        }
        //add the snapshot to the buffer and returns the new epoch
        _setSnapshot(Snapshot(block.number, blockClaims));
        //set the new epoch
        //_epochRegister.set(epoch);

        emit SnapshotTaken(
            _chainId,
            epoch,
            blockClaims.height,
            msg.sender,
            isSafeToProceedConsensus,
            groupSignature_
        );
        return isSafeToProceedConsensus;
    }

    /// @notice Saves next snapshot
    /// @param groupSignature_ The group signature used to sign the snapshots' block claims
    /// @param bClaims_ The claims being made about given block
    /// @return Flag whether we should kick off another round of key generation
    function migrateSnapshots(bytes[] memory groupSignature_, bytes[] memory bClaims_)
        public
        onlyFactory
        returns (bool)
    {
        {
            require(
                getEpoch() == 0,
                string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_MIGRATION_NOT_ALLOWED))
            );
            require(
                groupSignature_.length == bClaims_.length && groupSignature_.length >= 1,
                string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_MIGRATION_INPUT_DATA_MISMATCH))
            );
        }

        uint256 epoch;
        for (uint256 i = 0; i < bClaims_.length; i++) {
            BClaimsParserLibrary.BClaims memory blockClaims = BClaimsParserLibrary.extractBClaims(
                bClaims_[i]
            );
            require(
                blockClaims.height % _epochLength == 0,
                string(abi.encodePacked(SnapshotsErrorCodes.SNAPSHOT_INCORRECT_BLOCK_HEIGHT))
            );
            epoch = getEpochFromHeight(blockClaims.height);
            _setSnapshot(Snapshot(block.number, blockClaims));
            emit SnapshotTaken(
                _chainId,
                epoch,
                blockClaims.height,
                msg.sender,
                true,
                groupSignature_[i]
            );
        }
        _setEpoch(uint32(epoch));
        return true;
    }

    function getSnapshotDesperationFactor() public view returns (uint256) {
        return _snapshotDesperationFactor;
    }

    function getSnapshotDesperationDelay() public view returns (uint256) {
        return _snapshotDesperationDelay;
    }

    function getMinimumIntervalBetweenSnapshots() public view returns (uint256) {
        return _minimumIntervalBetweenSnapshots;
    }

    function getChainId() public view returns (uint256) {
        return _chainId;
    }

    function getEpoch() public view returns (uint256) {
        return __epoch._value;
    }

    function getEpochLength() public view returns (uint256) {
        return _epochLength;
    }

    function getChainIdFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return getSnapshot(epoch_).blockClaims.chainId;
    }

    function getChainIdFromLatestSnapshot() public view returns (uint256) {
        return getLatestSnapshot().blockClaims.chainId;
    }

    function getBlockClaimsFromSnapshot(uint256 epoch_)
        public
        view
        returns (BClaimsParserLibrary.BClaims memory)
    {
        return getSnapshot(epoch_).blockClaims;
    }

    function getBlockClaimsFromLatestSnapshot()
        public
        view
        returns (BClaimsParserLibrary.BClaims memory)
    {
        return getLatestSnapshot().blockClaims;
    }

    function getCommittedHeightFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return getSnapshot(epoch_).committedAt;
    }

    function getCommittedHeightFromLatestSnapshot() public view returns (uint256) {
        // return __snapshots.get(__epoch._value).committedAt;
        return getLatestSnapshot().committedAt;
    }

    function getAliceNetHeightFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return getSnapshot(epoch_).blockClaims.height;
    }

    function getAliceNetHeightFromLatestSnapshot() public view returns (uint256) {
        // return __snapshots.get(__epoch._value).blockClaims.height;
        return getLatestSnapshot().blockClaims.height;
    }
 
    function getSnapshot(uint256 epoch_) public view returns (Snapshot memory) {
        return __snapshots._array[epoch_%__snapshots._array.length];
    }

    function getLatestSnapshot() public view returns (Snapshot memory) {
        return __snapshots._array[getEpoch() % __snapshots._array.length];
    }
    function getEpochFromHeight(uint256 height) public view returns (uint256) {
        if (height <= _epochLength) {
            return 1;
        }
        if (height % _epochLength == 0) {
            return height / _epochLength;
        }
        return (height / _epochLength) + 1;
    }
    function mayValidatorSnapshot(
        uint256 numValidators,
        uint256 myIdx,
        uint256 blocksSinceDesperation,
        bytes32 blsig,
        uint256 desperationFactor
    ) public pure returns (bool) {
        return
            _mayValidatorSnapshot(
                numValidators,
                myIdx,
                blocksSinceDesperation,
                blsig,
                desperationFactor
            );
    }

    function _mayValidatorSnapshot(
        uint256 numValidators,
        uint256 myIdx,
        uint256 blocksSinceDesperation,
        bytes32 blsig,
        uint256 desperationFactor
    ) internal pure returns (bool) {
        uint256 numValidatorsAllowed = 1;

        uint256 desperation = 0;
        while (desperation < blocksSinceDesperation && numValidatorsAllowed <= numValidators / 3) {
            desperation += desperationFactor / numValidatorsAllowed;
            numValidatorsAllowed++;
        }

        uint256 rand = uint256(blsig);
        uint256 start = (rand % numValidators);
        uint256 end = (start + numValidatorsAllowed) % numValidators;

        if (end > start) {
            return myIdx >= start && myIdx < end;
        } else {
            return myIdx >= start || myIdx < end;
        }
    }
}
