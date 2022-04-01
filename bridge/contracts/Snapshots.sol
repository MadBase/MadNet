// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/interfaces/ISnapshots.sol";
import "contracts/interfaces/IValidatorPool.sol";
import "contracts/interfaces/IETHDKG.sol";
import "contracts/libraries/parsers/RCertParserLibrary.sol";
import "contracts/libraries/parsers/BClaimsParserLibrary.sol";
import "contracts/libraries/math/CryptoLibrary.sol";
import "contracts/libraries/snapshots/SnapshotsStorage.sol";
import "contracts/utils/DeterministicAddress.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {SnapshotsErrorCodes} from "contracts/libraries/errorCodes/SnapshotsErrorCodes.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @custom:salt Snapshots
/// @custom:deploy-type deployUpgradeable
contract Snapshots is Initializable, SnapshotsStorage, ISnapshots {
    using Strings for uint16;

    constructor(uint256 chainID_, uint256 epochLength_) SnapshotsStorage(chainID_, epochLength_) {}

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
        require(
            IValidatorPool(_validatorPoolAddress()).isValidator(msg.sender),
            SnapshotsErrorCodes.SNAPSHOT_ONLY_VALIDATORS_ALLOWED.toString()
        );
        require(
            IValidatorPool(_validatorPoolAddress()).isConsensusRunning(),
            SnapshotsErrorCodes.SNAPSHOT_CONSENSUS_RUNNING.toString()
        );

        require(
            block.number >= _snapshots[_epoch].committedAt + _minimumIntervalBetweenSnapshots,
            SnapshotsErrorCodes.SNAPSHOT_MIN_BLOCKS_INTERVAL_NOT_PASSED.toString()
        );

        (bool success, uint256 validatorIndex) = IETHDKG(_ethdkgAddress()).tryGetParticipantIndex(
            msg.sender
        );
        //todo:remove this, dummy operation only to silence linter
        validatorIndex;
        require(success, SnapshotsErrorCodes.SNAPSHOT_CALLER_NOT_ETHDKG_PARTICIPANT.toString());

        uint32 epoch = _epoch + 1;
        // uint256 ethBlocksSinceLastSnapshot = block.number - _snapshots[epoch - 1].committedAt;

        // TODO: BRING BACK AFTER GOLANG LOGIC IS DEBUGGED AND MERGED
        /*
        uint256 blocksSinceDesperation = ethBlocksSinceLastSnapshot >= _snapshotDesperationDelay
            ? ethBlocksSinceLastSnapshot - _snapshotDesperationDelay
            : 0;
        */

        // Check if sender is the elected validator allowed to make the snapshot
        // TODO: BRING BACK AFTER GOLANG LOGIC IS DEBUGGED AND MERGED
        /*
        require(
            _mayValidatorSnapshot(
                IValidatorPool(_validatorPoolAddress()).getValidatorsCount(),
                validatorIndex - 1,
                blocksSinceDesperation,
                keccak256(bClaims_),
                uint256(_snapshotDesperationFactor)
            ),
             "Snapshots: Validator not elected to do snapshot!"
        );
        */

        {
            (uint256[4] memory masterPublicKey, uint256[2] memory signature) = RCertParserLibrary
                .extractSigGroup(groupSignature_, 0);

            require(
                keccak256(abi.encodePacked(masterPublicKey)) ==
                    keccak256(abi.encodePacked(IETHDKG(_ethdkgAddress()).getMasterPublicKey())),
                SnapshotsErrorCodes.SNAPSHOT_WRONG_MASTER_PUBLIC_KEY.toString()
            );

            require(
                CryptoLibrary.Verify(
                    abi.encodePacked(keccak256(bClaims_)),
                    signature,
                    masterPublicKey
                ),
                SnapshotsErrorCodes.SNAPSHOT_SIGNATURE_VERIFICATION_FAILED.toString()
            );
        }

        BClaimsParserLibrary.BClaims memory blockClaims = BClaimsParserLibrary.extractBClaims(
            bClaims_
        );

        require(
            epoch * _epochLength == blockClaims.height,
            SnapshotsErrorCodes.SNAPSHOT_INCORRECT_BLOCK_HEIGHT.toString()
        );

        require(
            blockClaims.chainId == _chainId,
            SnapshotsErrorCodes.SNAPSHOT_INCORRECT_CHAIN_ID.toString()
        );

        bool isSafeToProceedConsensus = true;
        if (IValidatorPool(_validatorPoolAddress()).isMaintenanceScheduled()) {
            isSafeToProceedConsensus = false;
            IValidatorPool(_validatorPoolAddress()).pauseConsensus();
        }

        _snapshots[epoch] = Snapshot(block.number, blockClaims);
        _epoch = epoch;

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
        return _epoch;
    }

    function getEpochLength() public view returns (uint256) {
        return _epochLength;
    }

    function getChainIdFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return _snapshots[epoch_].blockClaims.chainId;
    }

    function getChainIdFromLatestSnapshot() public view returns (uint256) {
        return _snapshots[_epoch].blockClaims.chainId;
    }

    function getBlockClaimsFromSnapshot(uint256 epoch_)
        public
        view
        returns (BClaimsParserLibrary.BClaims memory)
    {
        return _snapshots[epoch_].blockClaims;
    }

    function getBlockClaimsFromLatestSnapshot()
        public
        view
        returns (BClaimsParserLibrary.BClaims memory)
    {
        return _snapshots[_epoch].blockClaims;
    }

    function getCommittedHeightFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return _snapshots[epoch_].committedAt;
    }

    function getCommittedHeightFromLatestSnapshot() public view returns (uint256) {
        return _snapshots[_epoch].committedAt;
    }

    function getAliceNetHeightFromSnapshot(uint256 epoch_) public view returns (uint256) {
        return _snapshots[epoch_].blockClaims.height;
    }

    function getAliceNetHeightFromLatestSnapshot() public view returns (uint256) {
        return _snapshots[_epoch].blockClaims.height;
    }

    function getSnapshot(uint256 epoch_) public view returns (Snapshot memory) {
        return _snapshots[epoch_];
    }

    function getLatestSnapshot() public view returns (Snapshot memory) {
        return _snapshots[_epoch];
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
