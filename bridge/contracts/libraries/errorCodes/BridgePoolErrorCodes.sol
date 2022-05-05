// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library BridgePoolErrorCodes {
    // BridgePool error codes
    bytes32 public constant BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER = "2400"; //"BridgePool: Deposit can only be requested for the owner in burned UTXO"
    bytes32 public constant BRIDGEPOOL_PROOF_OF_BURN_NOT_VERIFIED = "2401"; //"BridgePool: Proof of burn in aliceNet could not be verified"
}
