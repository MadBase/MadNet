// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.13;

import "contracts/utils/ImmutableAuth.sol";
import "contracts/utils/DeterministicAddress.sol";

abstract contract RoleBasedAccessControl is ImmutableFactory {
    string public constant SALT_ERC20 = "ERC20";
    string public constant SALT_ERC721 = "ERC721";
    string public constant SALT_ERC1155 = "ERC1155";

    uint256 public immutable networkID;

    constructor(uint256 networkID_) ImmutableFactory(msg.sender) {
        networkID = networkID_;
    }

    modifier onlyERC20(address erc20Address_) {
        address expectedAddr = getMetamorphicContractAddress(
            keccak256(abi.encodePacked(SALT_ERC20, networkID, erc20Address_)),
            _factoryAddress()
        );

        require(expectedAddr == msg.sender, "Not ERC20");
        _;
    }
    modifier onlyERC721(address erc20Address_) {
        address expectedAddr = getMetamorphicContractAddress(
            keccak256(abi.encodePacked(SALT_ERC721, networkID, erc20Address_)),
            _factoryAddress()
        );

        require(expectedAddr == msg.sender, "Not ERC721");
        _;
    }
    modifier onlyERC1155(address erc20Address_) {
        address expectedAddr = getMetamorphicContractAddress(
            keccak256(abi.encodePacked(SALT_ERC1155, networkID, erc20Address_)),
            _factoryAddress()
        );

        require(expectedAddr == msg.sender, "Not ERC721");
        _;
    }
}
