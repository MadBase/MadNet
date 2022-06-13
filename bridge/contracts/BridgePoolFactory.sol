// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;
import "contracts/AliceNetFactory.sol";
import "contracts/libraries/factory/AliceNetFactoryBase.sol";
import "contracts/BridgePool.sol";
import "contracts/Proxy.sol";
import "contracts/utils/ImmutableAuth.sol";

/// @custom:salt BridgePoolFactory
/// @custom:deploy-type deployUpgradeable
contract BridgePoolFactory is
    ImmutableFactory,
    ImmutableBridgePoolFactory,
    ImmutableBridgePoolDepositNotifier
{
    uint256 internal immutable _networkId;
    string internal constant _BRIDGE_POOL_TAG = "ERC";

    constructor(uint256 networkId_) ImmutableFactory(msg.sender) {
        _networkId = networkId_;
    }

    /**
     * @notice deployNewPool delegates call to this contract's method "deployViaFactoryLogic" through alicenet factory
     * @param erc20Contract_ address of ERC20 token contract
     * @param token address of bridge token contract
     */
    function deployNewPool(address erc20Contract_, address token) public {
        bytes memory callData = abi.encodeWithSelector(
            this.deployViaFactoryLogic.selector,
            erc20Contract_,
            token
        );
        address impAddress = Proxy(payable(address(this))).getImplementationAddress();
        AliceNetFactory(_factoryAddress()).delegateCallAny(impAddress, callData);
    }

    /**
     * @notice deployViaFactoryLogic deploys a BridgePool contract between ERC20 token and token
     * @param erc20Contract_ address of ERC20 token contract
     * @param token address of bridge token contract
     * @return contract address
     */
    function deployViaFactoryLogic(address erc20Contract_, address token)
        public
        onlyBridgePoolFactory
        returns (address)
    {
        bytes memory initializers = abi.encode(erc20Contract_, token);
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
            token
        );
        AliceNetFactory(_factoryAddress()).upgradeProxy(salt, contractAddress, initCallData);
        return proxyAddress;
    }

    /**
     * @notice getSaltFromAddress calculates salt for a BridgePool contract based on ERC20 contract's address
     * @param erc20Contract_ address of ERC20 contract of BridgePool
     * @return calculated salt
     */
    function getSaltFromERC20Address(address erc20Contract_)
        public
        view
        returns (
            //onlyBridgePoolDepositNotifier
            bytes32
        )
    {
        return
            keccak256(
                bytes.concat(
                    keccak256(abi.encodePacked(erc20Contract_)),
                    keccak256(abi.encodePacked(_BRIDGE_POOL_TAG)),
                    keccak256(abi.encodePacked(_networkId))
                )
            );
    }
}
