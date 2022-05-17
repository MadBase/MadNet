// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;
import "contracts/utils/DeterministicAddress.sol";
import "contracts/Proxy.sol";
import "contracts/AliceNetFactory.sol";

/// @custom:salt BridgePoolFactory
contract BridgePoolFactory is AliceNetFactory {
    /**
     * @dev The constructor encodes the proxy deploy byte code with the _UNIVERSAL_DEPLOY_CODE at the
     * head and the factory address at the tail, and deploys the proxy byte code using create OpCode.
     * The result of this deployment will be a contract with the proxy contract deployment bytecode with
     * its constructor at the head, runtime code in the body and constructor args at the tail. The
     * constructor then sets proxyTemplate_ state var to the deployed proxy template address the deploy
     * account will be set as the first owner of the factory.
     * @param selfAddr_ is the factory contracts
     * address (address of itself)
     */
    constructor(address selfAddr_) AliceNetFactory(selfAddr_) {}
}
