// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/interfaces/IValidatorPool.sol";
import "contracts/interfaces/IETHDKG.sol";
import "contracts/utils/ImmutableAuth.sol";
import "contracts/libraries/snapshots/SnapshotRingBuffer.sol";



abstract contract SnapshotsStorage is ImmutableETHDKG, ImmutableValidatorPool, SnapshotRingBuffer {
    using RingBuffer for SnapshotBuffer;
    using EpochLib for Epoch;
    uint256 internal immutable _epochLength;

    uint256 internal immutable _chainId;

    // uint32 internal _epoch;

    // Number of ethereum blocks that we should wait between snapshots. Mainly used to prevent the
    // submission of snapshots in short amount of time by validators that could be potentially being
    // malicious
    uint32 internal _minimumIntervalBetweenSnapshots;

    // after how many eth blocks of not having a snapshot will we start allowing more validators to
    // make it
    uint32 internal _snapshotDesperationDelay;

    // how quickly more validators will be allowed to make a snapshot, once
    // _snapshotDesperationDelay has passed
    uint32 internal _snapshotDesperationFactor;

    // mapping(uint256 => Snapshot) internal _snapshots;
    //epoch counter wrapped in a struct
    Epoch internal __epoch;
    //new snapshot ring buffer
    SnapshotBuffer internal __snapshots;

    constructor(uint256 chainId_, uint256 epochLength_)
        ImmutableFactory(msg.sender)
        ImmutableETHDKG()
        ImmutableValidatorPool()
    {
        _chainId = chainId_;
        _epochLength = epochLength_;
    }
    
    function _getEpochFromHeight(uint32 height_) internal override view returns(uint32){
        if (height_ <= _epochLength) {
            return 1;
        }
        if (height_ % _epochLength == 0) {
            //TODO make sure this does what i think it does
            uint32 epoch = uint32(height_ / _epochLength);
            return uint32(height_ / _epochLength);
        }
        return uint32((height_ / _epochLength) + 1);
    }
    
    function _snapshots() internal override view returns(SnapshotBuffer storage){
        return __snapshots;
    }
    
    function _epoch() internal override view returns(Epoch storage){
        return __epoch;
    }
}