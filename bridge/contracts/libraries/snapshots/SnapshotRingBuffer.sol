// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11; 

import "contracts/libraries/parsers/BClaimsParserLibrary.sol"; 
import "contracts/interfaces/ISnapshots.sol";

struct Epoch {
    uint32 _value;
}

struct SnapshotBuffer {
    Snapshot[6] _array;
}


library RingBuffer {

    function get(SnapshotBuffer memory self_, uint32 epoch_)
        internal
        pure
        returns (Snapshot memory)
    {
        return self_._array[_indexFor(self_, epoch_)];
    }
    //writes the new snapshot to the buffer
    function set(SnapshotBuffer storage self_, function (uint32) returns(uint32) epochFor_, Snapshot memory new_)
        internal
        returns(uint32)
    {
        //get the epoch corresponding to the blocknumber 
        uint32 epoch = epochFor_(new_.blockClaims.height);
        //gets the snapshot that was at that location
        Snapshot storage old = self_._array[_indexFor(self_, epoch)];
        //checks if the new snapshot height is greater than the previous 
        require(new_.blockClaims.height > old.blockClaims.height, "invalid new blockheight");
        unsafeSet(self_, new_, epoch);
        return epoch;
    }

    function unsafeSet(SnapshotBuffer storage self_, Snapshot memory new_, uint32 epoch_) internal {
        self_._array[_indexFor(self_, epoch_)] = new_;
    }
    /**
    * @dev calculates the congruent value for current epoch in respect to the array length 
    * for index to be replaced with most recent epoch
    * @param epoch_ epoch_ number associated with the snapshot 
    */
    function _indexFor(SnapshotBuffer memory self_, uint32 epoch_) internal pure returns(uint256) {
        //TODO determine if this is necessary 
        // require(epoch_ > 0);
        return epoch_ % self_._array.length;
    }
}

library EpochLib {
    function set(Epoch storage self_, uint32 value_) internal {
        self_._value = value_;
    }

    function get(Epoch memory self_) internal view returns(uint32) {
        return _max(1, self_._value);
    }
    //TODO determine if useful
    function isZero(Epoch storage self_) internal view returns(bool) {
        return self_._value == 0;
    }

    function _max(uint32 a, uint32 b) internal pure returns(uint32) {
        if (a > b ) {
            return a;
        }
        return b;
    }
}

abstract contract SnapshotRingBuffer {
    using RingBuffer for SnapshotBuffer;
    using EpochLib for Epoch;

    // Must be defined in storage contract
    function _getEpochFromHeight(uint32) internal virtual view returns(uint32);

    // Must be defined in storage contract
    function _snapshots() internal virtual view returns(SnapshotBuffer storage);

    // Must be defined in storage contract
    function _epoch() internal virtual view returns(Epoch storage);
    
    function _setEpoch(uint32 epoch_) internal {
        _epoch().set(epoch_);
    }
    function _getEpoch() internal view returns (uint32) {
        return _epoch().get();
    }
    /**
    * @notice Assigns the snapshot to correct index and updates __epoch
    * @param snapshot_ to be stored
    * @return epoch of the passed snapshot
    */
    function _setSnapshot(Snapshot memory snapshot_) internal returns(uint32) {
        uint32 epoch = _snapshots().set(_getEpochFromHeight, snapshot_);
        _epoch().set(epoch);
        return epoch;
    }

    /**
    * @notice Returns the snapshot for the passed epoch and safety flag
    * @param epoch_ of the snapshot
    * @return ok if the struct is valid and the snapshot struct itself
    */
    function _getSnapshot(uint32 epoch_) internal view returns (bool ok, Snapshot memory snapshot) {
        //get the pointer to the specified epoch snapshot
        Snapshot memory temp = _snapshots().get(epoch_);
        if (_getEpochFromHeight(temp.blockClaims.height) == epoch_) {
            ok = true;
            snapshot = temp;
        }
        return (ok, snapshot);
    }

    /**
    * @notice Returns the snapshot for the passed height and safety flag
    * @param height_ of the snapshot
    * @return ok if the struct is valid and the snapshot struct itself
    */
    function _getSnapshotForHeight(uint32 height_) internal view returns (bool ok, Snapshot memory snapshot) {
        return _getSnapshot(_getEpochFromHeight(height_));
    }

    /**
    * @return ok if the struct is valid and the snapshot struct itself
    */
    function _getLastSnapshot() internal view returns (bool ok, Snapshot memory snapshot) {
        return _getSnapshot(_epoch().get());
    }
}

