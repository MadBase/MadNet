// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;
import "contracts/AliceNetFactory.sol";
import "contracts/libraries/factory/AliceNetFactoryBase.sol";
import "contracts/BridgePool.sol";
import "contracts/Proxy.sol";
import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt BridgePoolFactory
/// @custom:deploy-type deployStatic
contract BridgePoolFactory is ImmutableFactory, ImmutableBridgePoolFactory {
    constructor() ImmutableFactory(msg.sender) {}

    //TODO: Set proper accessControl?
    function deployNewPool(address erc20Contract_, address token2) public {
        bytes memory callData = abi.encodeWithSelector(
            this.deployViaFactoryLogic.selector,
            erc20Contract_,
            token2
        );
        address impAddress = Proxy(payable(address(this))).getImplementationAddress();
        AliceNetFactory(_factoryAddress()).delegateCallAny(impAddress, callData);
    }

    function deployViaFactoryLogic(address erc20Contract_, address token2)
        public
        onlyBridgePoolFactory
        returns (address)
    {
        bytes memory initializers = abi.encode(erc20Contract_, token2);
        bytes memory deployCode = bytes.concat(type(BridgePool).creationCode, initializers);
        AliceNetFactory(_factoryAddress()).deployTemplate(deployCode);
        bytes32 salt = getSaltFromERC20Address(erc20Contract_);
        //TODO: set proper value
        uint256 value = 0;
        address contractAddress = AliceNetFactory(_factoryAddress()).deployCreate2(
            value,
            salt,
            deployCode
        );
        address proxyAddress = AliceNetFactory(_factoryAddress()).deployProxy(salt);
        bytes memory initCallData = abi.encodeWithSelector(
            BridgePool.initialize.selector,
            erc20Contract_,
            token2
        );
        AliceNetFactory(_factoryAddress()).upgradeProxy(salt, contractAddress, initCallData);
        return proxyAddress;
    }

    /**
     * @dev getSaltFromAddress calculates salt for a BridgePool contract based on ERC20 contract's address
     * @param erc20Contract_ address of ERC20 contract of BridgePool
     * @return calculated salt
     */
    function getSaltFromERC20Address(address erc20Contract_) internal pure returns (bytes32) {
        return
            keccak256(bytes.concat(keccak256(abi.encodePacked(erc20Contract_)), keccak256("ERC")));
    }
}
