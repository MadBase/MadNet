// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import {
    BClaimsParserLibraryErrorCodes
} from "contracts/libraries/errorCodes/BClaimsParserLibraryErrorCodes.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

import "./BaseParserLibrary.sol";

/// @title Library to parse the BClaims structure from a blob of capnproto data
library BClaimsParserLibrary {
    using Strings for uint16;
    struct BClaims {
        uint32 chainId;
        uint32 height;
        uint32 txCount;
        bytes32 prevBlock;
        bytes32 txRoot;
        bytes32 stateRoot;
        bytes32 headerRoot;
    }

    /** @dev size in bytes of a BCLAIMS cap'npro structure without the cap'n
      proto header bytes*/
    uint256 internal constant _BCLAIMS_SIZE = 176;
    /** @dev Number of bytes of a capnproto header, the data starts after the
      header */
    uint256 internal constant _CAPNPROTO_HEADER_SIZE = 8;

    /**
    @notice This function computes the offset adjustment in the pointer section
    of the capnproto blob of data. In case the txCount is 0, the value is not
    included in the binary blob by capnproto. Therefore, we need to deduce 8
    bytes from the pointer's offset.
    */
    /// @param src Binary data containing a BClaims serialized struct
    /// @param dataOffset Blob of binary data with a capnproto serialization
    /// @return pointerOffsetAdjustment the pointer offset adjustment in the blob data
    /// @dev Execution cost: 499 gas
    function getPointerOffsetAdjustment(bytes memory src, uint256 dataOffset)
        internal
        pure
        returns (uint16 pointerOffsetAdjustment)
    {
        // Size in capnproto words (16 bytes) of the data section
        uint16 dataSectionSize = BaseParserLibrary.extractUInt16(src, dataOffset);
        require(
            dataSectionSize > 0 && dataSectionSize <= 2,
            BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_SIZE_THRESHOLD_EXCEEDED.toString()
        );
        // In case the txCount is 0, the value is not included in the binary
        // blob by capnproto. Therefore, we need to deduce 8 bytes from the
        // pointer's offset.
        if (dataSectionSize == 1) {
            pointerOffsetAdjustment = 8;
        } else {
            pointerOffsetAdjustment = 0;
        }
    }

    /**
    @notice This function is for deserializing data directly from capnproto
            BClaims. It will skip the first 8 bytes (capnproto headers) and
            deserialize the BClaims Data. This function also computes the right
            PointerOffset adjustment (see the documentation on
            `getPointerOffsetAdjustment(bytes, uint256)` for more details). If
            BClaims is being extracted from inside of other structure (E.g
            PClaims capnproto) use the `extractInnerBClaims(bytes, uint,
            uint16)` instead.
    */
    /// @param src Binary data containing a BClaims serialized struct with Capn Proto headers
    /// @return bClaims the BClaims struct
    /// @dev Execution cost: 2484 gas
    function extractBClaims(bytes memory src) internal pure returns (BClaims memory bClaims) {
        return extractInnerBClaims(src, _CAPNPROTO_HEADER_SIZE, getPointerOffsetAdjustment(src, 4));
    }

    /**
    @notice This function is for deserializing the BClaims struct from an defined
            location inside a binary blob. E.G Extract BClaims from inside of
            other structure (E.g PClaims capnproto) or skipping the capnproto
            headers.
    */
    /// @param src Binary data containing a BClaims serialized struct without Capn proto headers
    /// @param dataOffset offset to start reading the BClaims data from inside src
    /// @param pointerOffsetAdjustment Pointer's offset that will be deduced from the pointers location, in case txCount is missing in the binary
    /// @return bClaims the BClaims struct
    /// @dev Execution cost: 2126 gas
    function extractInnerBClaims(
        bytes memory src,
        uint256 dataOffset,
        uint16 pointerOffsetAdjustment
    ) internal pure returns (BClaims memory bClaims) {
        require(
            dataOffset + _BCLAIMS_SIZE - pointerOffsetAdjustment > dataOffset,
            BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_DATA_OFFSET_OVERFLOW.toString()
        );
        require(
            src.length >= dataOffset + _BCLAIMS_SIZE - pointerOffsetAdjustment,
            BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_NOT_ENOUGH_BYTES.toString()
        );

        if (pointerOffsetAdjustment == 0) {
            bClaims.txCount = BaseParserLibrary.extractUInt32(src, dataOffset + 8);
        } else {
            // In case the txCount is 0, the value is not included in the binary
            // blob by capnproto.
            bClaims.txCount = 0;
        }

        bClaims.chainId = BaseParserLibrary.extractUInt32(src, dataOffset);
        require(
            bClaims.chainId > 0,
            BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_CHAINID_ZERO.toString()
        );
        bClaims.height = BaseParserLibrary.extractUInt32(src, dataOffset + 4);
        require(
            bClaims.height > 0,
            BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_HEIGHT_ZERO.toString()
        );
        bClaims.prevBlock = BaseParserLibrary.extractBytes32(
            src,
            dataOffset + 48 - pointerOffsetAdjustment
        );
        bClaims.txRoot = BaseParserLibrary.extractBytes32(
            src,
            dataOffset + 80 - pointerOffsetAdjustment
        );
        bClaims.stateRoot = BaseParserLibrary.extractBytes32(
            src,
            dataOffset + 112 - pointerOffsetAdjustment
        );
        bClaims.headerRoot = BaseParserLibrary.extractBytes32(
            src,
            dataOffset + 144 - pointerOffsetAdjustment
        );
    }
}
