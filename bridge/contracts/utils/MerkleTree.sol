// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract MerkleTree {
    // returns the merkle root by hashing the merkle proof items
    function _calculateRootFromMerkleProof(
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
                        _calculateRootFromMerkleProof(
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
                    _calculateRootFromMerkleProof(
                        auditPath,
                        merkleKeyIndex + 1,
                        merkleKeyHash,
                        merkleLeafHash
                    ),
                    auditPath[auditPath.length - merkleKeyIndex - 1]
                )
            );
    }

    function _bitIsSet(bytes32 bits, uint256 merkleKeyIndex) internal pure returns (bool) {
        return (bits[merkleKeyIndex / 8] & (bytes1(0x01) << (7 - (merkleKeyIndex % 8))) != 0);
    }
}
