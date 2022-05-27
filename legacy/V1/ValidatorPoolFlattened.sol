// Sources flattened with hardhat v2.9.1 https://hardhat.org

// File contracts/interfaces/IStakingNFT.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

interface IStakingNFT {
    function skimExcessEth(address to_) external returns (uint256 excess);

    function skimExcessToken(address to_) external returns (uint256 excess);

    function depositToken(uint8 magic_, uint256 amount_) external;

    function depositEth(uint8 magic_) external payable;

    function lockPosition(
        address caller_,
        uint256 tokenID_,
        uint256 lockDuration_
    ) external returns (uint256);

    function lockOwnPosition(uint256 tokenID_, uint256 lockDuration_)
        external
        returns (uint256);

    function lockWithdraw(uint256 tokenID_, uint256 lockDuration_)
        external
        returns (uint256);

    function mint(uint256 amount_) external returns (uint256 tokenID);

    function mintTo(
        address to_,
        uint256 amount_,
        uint256 lockDuration_
    ) external returns (uint256 tokenID);

    function burn(uint256 tokenID_)
        external
        returns (uint256 payoutEth, uint256 payoutAToken);

    function burnTo(address to_, uint256 tokenID_)
        external
        returns (uint256 payoutEth, uint256 payoutAToken);

    function collectEth(uint256 tokenID_) external returns (uint256 payout);

    function collectToken(uint256 tokenID_) external returns (uint256 payout);

    function collectEthTo(address to_, uint256 tokenID_)
        external
        returns (uint256 payout);

    function collectTokenTo(address to_, uint256 tokenID_)
        external
        returns (uint256 payout);

    function getPosition(uint256 tokenID_)
        external
        view
        returns (
            uint256 shares,
            uint256 freeAfter,
            uint256 withdrawFreeAfter,
            uint256 accumulatorEth,
            uint256 accumulatorToken
        );

    function getAccumulatorScaleFactor() external view returns (uint256);

    function getTotalShares() external view returns (uint256);

    function getTotalReserveEth() external view returns (uint256);

    function getTotalReserveAToken() external view returns (uint256);

    function estimateEthCollection(uint256 tokenID_)
        external
        view
        returns (uint256 payout);

    function estimateTokenCollection(uint256 tokenID_)
        external
        view
        returns (uint256 payout);

    function estimateExcessToken() external view returns (uint256 excess);

    function estimateExcessEth() external view returns (uint256 excess);

    function getEthAccumulator()
        external
        view
        returns (uint256 accumulator, uint256 slush);

    function getTokenAccumulator()
        external
        view
        returns (uint256 accumulator, uint256 slush);
}

// File contracts/interfaces/IERC20Transferable.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

interface IERC20Transferable {
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);

    function transfer(address recipient, uint256 amount)
        external
        returns (bool);

    function approve(address spender, uint256 amount) external returns (bool);

    function balanceOf(address account) external view returns (uint256);
}

// File contracts/interfaces/IERC721Transferable.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

interface IERC721Transferable {
    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId
    ) external;
}

// File contracts/libraries/errorCodes/BClaimsParserLibraryErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library BClaimsParserLibraryErrorCodes {
    // BClaimsParserLibrary error codes
    bytes32 public constant BCLAIMSPARSERLIB_SIZE_THRESHOLD_EXCEEDED = "1100"; //"BClaimsParserLibrary: The size of the data section should be 1 or 2 words!"
    bytes32 public constant BCLAIMSPARSERLIB_DATA_OFFSET_OVERFLOW = "1101"; //"BClaimsParserLibrary: Invalid parsing. Overflow on the dataOffset parameter"
    bytes32 public constant BCLAIMSPARSERLIB_NOT_ENOUGH_BYTES = "1102"; //"BClaimsParserLibrary: Invalid parsing. Not enough bytes to extract BClaims"
    bytes32 public constant BCLAIMSPARSERLIB_CHAINID_ZERO = "1103"; //"BClaimsParserLibrary: Invalid parsing. The chainId should be greater than 0!"
    bytes32 public constant BCLAIMSPARSERLIB_HEIGHT_ZERO = "1104"; //"BClaimsParserLibrary: Invalid parsing. The height should be greater than 0!"
}

// File contracts/libraries/errorCodes/BaseParserLibraryErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library BaseParserLibraryErrorCodes {
    // BaseParserLibrary error codes
    bytes32 public constant BASEPARSERLIB_OFFSET_PARAMETER_OVERFLOW = "1000"; // "BaseParserLibrary: An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_OFFSET_OUT_OF_BOUNDS = "1001"; // "BaseParserLibrary: Trying to read an offset out of boundaries in the src binary!"
    bytes32 public constant BASEPARSERLIB_LE_UINT16_OFFSET_PARAMETER_OVERFLOW =
        "1002"; //  "BaseParserLibrary: Error extracting uin16! An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_LE_UINT16_OFFSET_OUT_OF_BOUNDS =
        "1003"; //  "BaseParserLibrary: UINT16 ERROR! Trying to read an offset out of boundaries!"
    bytes32 public constant BASEPARSERLIB_BE_UINT16_OFFSET_PARAMETER_OVERFLOW =
        "1004"; // "BaseParserLibrary: UINT16 ERROR! An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_BE_UINT16_OFFSET_OUT_OF_BOUNDS =
        "1005"; // "BaseParserLibrary: UINT16 ERROR! Trying to read an offset out of boundaries!"
    bytes32 public constant BASEPARSERLIB_BOOL_OFFSET_PARAMETER_OVERFLOW =
        "1006"; // "BaseParserLibrary: BOOL ERROR: OVERFLOW!"
    bytes32 public constant BASEPARSERLIB_BOOL_OFFSET_OUT_OF_BOUNDS = "1007"; //  "BaseParserLibrary: BOOL ERROR: OFFSET OUT OF BOUNDARIES!"
    bytes32 public constant BASEPARSERLIB_LE_UINT256_OFFSET_PARAMETER_OVERFLOW =
        "1008"; //  "BaseParserLibrary: Error extracting uin16! An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_LE_UINT256_OFFSET_OUT_OF_BOUNDS =
        "1009"; //  "BaseParserLibrary: UINT16 ERROR! Trying to read an offset out of boundaries!"
    bytes32 public constant BASEPARSERLIB_BE_UINT256_OFFSET_PARAMETER_OVERFLOW =
        "1010"; // "BaseParserLibrary: UINT16 ERROR! An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_BE_UINT256_OFFSET_OUT_OF_BOUNDS =
        "1011"; // "BaseParserLibrary: UINT16 ERROR! Trying to read an offset out of boundaries!"
    bytes32 public constant BASEPARSERLIB_BYTES_OFFSET_PARAMETER_OVERFLOW =
        "1012"; // "BaseParserLibrary: An overflow happened with the offset or the howManyBytes parameter!"
    bytes32 public constant BASEPARSERLIB_BYTES_OFFSET_OUT_OF_BOUNDS = "1013"; //   "BaseParserLibrary: Not enough bytes to extract in the src binary"
    bytes32 public constant BASEPARSERLIB_BYTES32_OFFSET_PARAMETER_OVERFLOW =
        "1014"; // "BaseParserLibrary: An overflow happened with the offset parameter!"
    bytes32 public constant BASEPARSERLIB_BYTES32_OFFSET_OUT_OF_BOUNDS = "1015"; //   "BaseParserLibrary: not enough bytes to extract"
}

// File contracts/libraries/parsers/BaseParserLibrary.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library BaseParserLibrary {
    // Size of a word, in bytes.
    uint256 internal constant _WORD_SIZE = 32;
    // Size of the header of a 'bytes' array.
    uint256 internal constant _BYTES_HEADER_SIZE = 32;

    /// @notice Extracts a uint32 from a little endian bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return val a uint32
    /// @dev ~559 gas
    function extractUInt32(bytes memory src, uint256 offset)
        internal
        pure
        returns (uint32 val)
    {
        require(
            offset + 4 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= offset + 4,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_OFFSET_OUT_OF_BOUNDS
                )
            )
        );

        assembly {
            val := shr(sub(256, 32), mload(add(add(src, 0x20), offset)))
            val := or(
                or(
                    or(
                        shr(24, and(val, 0xff000000)),
                        shr(8, and(val, 0x00ff0000))
                    ),
                    shl(8, and(val, 0x0000ff00))
                ),
                shl(24, and(val, 0x000000ff))
            )
        }
    }

    /// @notice Extracts a uint16 from a little endian bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return val a uint16
    /// @dev ~204 gas
    function extractUInt16(bytes memory src, uint256 offset)
        internal
        pure
        returns (uint16 val)
    {
        require(
            offset + 2 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_LE_UINT16_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= offset + 2,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_LE_UINT16_OFFSET_OUT_OF_BOUNDS
                )
            )
        );

        assembly {
            val := shr(sub(256, 16), mload(add(add(src, 0x20), offset)))
            val := or(shr(8, and(val, 0xff00)), shl(8, and(val, 0x00ff)))
        }
    }

    /// @notice Extracts a uint16 from a big endian bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return val a uint16
    /// @dev ~204 gas
    function extractUInt16FromBigEndian(bytes memory src, uint256 offset)
        internal
        pure
        returns (uint16 val)
    {
        require(
            offset + 2 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BE_UINT16_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= offset + 2,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BE_UINT16_OFFSET_OUT_OF_BOUNDS
                )
            )
        );

        assembly {
            val := and(
                shr(sub(256, 16), mload(add(add(src, 0x20), offset))),
                0xffff
            )
        }
    }

    /// @notice Extracts a bool from a bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return a bool
    /// @dev ~204 gas
    function extractBool(bytes memory src, uint256 offset)
        internal
        pure
        returns (bool)
    {
        require(
            offset + 1 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BOOL_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= offset + 1,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BOOL_OFFSET_OUT_OF_BOUNDS
                )
            )
        );
        uint256 val;
        assembly {
            val := shr(sub(256, 8), mload(add(add(src, 0x20), offset)))
            val := and(val, 0x01)
        }
        return val == 1;
    }

    /// @notice Extracts a uint256 from a little endian bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return val a uint256
    /// @dev ~5155 gas
    function extractUInt256(bytes memory src, uint256 offset)
        internal
        pure
        returns (uint256 val)
    {
        require(
            offset + 31 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_LE_UINT256_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length > offset + 31,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_LE_UINT256_OFFSET_OUT_OF_BOUNDS
                )
            )
        );

        assembly {
            val := mload(add(add(src, 0x20), offset))
        }
    }

    /// @notice Extracts a uint256 from a big endian bytes array.
    /// @param src the binary data
    /// @param offset place inside `src` to start reading data from
    /// @return val a uint256
    /// @dev ~1400 gas
    function extractUInt256FromBigEndian(bytes memory src, uint256 offset)
        internal
        pure
        returns (uint256 val)
    {
        require(
            offset + 31 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BE_UINT256_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length > offset + 31,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BE_UINT256_OFFSET_OUT_OF_BOUNDS
                )
            )
        );

        uint256 srcDataPointer;
        uint32 val0 = 0;
        uint32 val1 = 0;
        uint32 val2 = 0;
        uint32 val3 = 0;
        uint32 val4 = 0;
        uint32 val5 = 0;
        uint32 val6 = 0;
        uint32 val7 = 0;

        assembly {
            srcDataPointer := mload(add(add(src, 0x20), offset))
            val0 := and(srcDataPointer, 0xffffffff)
            val1 := and(shr(32, srcDataPointer), 0xffffffff)
            val2 := and(shr(64, srcDataPointer), 0xffffffff)
            val3 := and(shr(96, srcDataPointer), 0xffffffff)
            val4 := and(shr(128, srcDataPointer), 0xffffffff)
            val5 := and(shr(160, srcDataPointer), 0xffffffff)
            val6 := and(shr(192, srcDataPointer), 0xffffffff)
            val7 := and(shr(224, srcDataPointer), 0xffffffff)

            val0 := or(
                or(
                    or(
                        shr(24, and(val0, 0xff000000)),
                        shr(8, and(val0, 0x00ff0000))
                    ),
                    shl(8, and(val0, 0x0000ff00))
                ),
                shl(24, and(val0, 0x000000ff))
            )
            val1 := or(
                or(
                    or(
                        shr(24, and(val1, 0xff000000)),
                        shr(8, and(val1, 0x00ff0000))
                    ),
                    shl(8, and(val1, 0x0000ff00))
                ),
                shl(24, and(val1, 0x000000ff))
            )
            val2 := or(
                or(
                    or(
                        shr(24, and(val2, 0xff000000)),
                        shr(8, and(val2, 0x00ff0000))
                    ),
                    shl(8, and(val2, 0x0000ff00))
                ),
                shl(24, and(val2, 0x000000ff))
            )
            val3 := or(
                or(
                    or(
                        shr(24, and(val3, 0xff000000)),
                        shr(8, and(val3, 0x00ff0000))
                    ),
                    shl(8, and(val3, 0x0000ff00))
                ),
                shl(24, and(val3, 0x000000ff))
            )
            val4 := or(
                or(
                    or(
                        shr(24, and(val4, 0xff000000)),
                        shr(8, and(val4, 0x00ff0000))
                    ),
                    shl(8, and(val4, 0x0000ff00))
                ),
                shl(24, and(val4, 0x000000ff))
            )
            val5 := or(
                or(
                    or(
                        shr(24, and(val5, 0xff000000)),
                        shr(8, and(val5, 0x00ff0000))
                    ),
                    shl(8, and(val5, 0x0000ff00))
                ),
                shl(24, and(val5, 0x000000ff))
            )
            val6 := or(
                or(
                    or(
                        shr(24, and(val6, 0xff000000)),
                        shr(8, and(val6, 0x00ff0000))
                    ),
                    shl(8, and(val6, 0x0000ff00))
                ),
                shl(24, and(val6, 0x000000ff))
            )
            val7 := or(
                or(
                    or(
                        shr(24, and(val7, 0xff000000)),
                        shr(8, and(val7, 0x00ff0000))
                    ),
                    shl(8, and(val7, 0x0000ff00))
                ),
                shl(24, and(val7, 0x000000ff))
            )

            val := or(
                or(
                    or(
                        or(
                            or(
                                or(
                                    or(shl(224, val0), shl(192, val1)),
                                    shl(160, val2)
                                ),
                                shl(128, val3)
                            ),
                            shl(96, val4)
                        ),
                        shl(64, val5)
                    ),
                    shl(32, val6)
                ),
                val7
            )
        }
    }

    /// @notice Reverts a bytes array. Can be used to convert an array from little endian to big endian and vice-versa.
    /// @param orig the binary data
    /// @return reversed the reverted bytes array
    /// @dev ~13832 gas
    function reverse(bytes memory orig)
        internal
        pure
        returns (bytes memory reversed)
    {
        reversed = new bytes(orig.length);
        for (uint256 idx = 0; idx < orig.length; idx++) {
            reversed[orig.length - idx - 1] = orig[idx];
        }
    }

    /// @notice Copy 'len' bytes from memory address 'src', to address 'dest'. This function does not check the or destination, it only copies the bytes.
    /// @param src the pointer to the source
    /// @param dest the pointer to the destination
    /// @param len the len of data to be copied
    function copy(
        uint256 src,
        uint256 dest,
        uint256 len
    ) internal pure {
        // Copy word-length chunks while possible
        for (; len >= _WORD_SIZE; len -= _WORD_SIZE) {
            assembly {
                mstore(dest, mload(src))
            }
            dest += _WORD_SIZE;
            src += _WORD_SIZE;
        }
        // Returning earlier if there's no leftover bytes to copy
        if (len == 0) {
            return;
        }
        // Copy remaining bytes
        uint256 mask = 256**(_WORD_SIZE - len) - 1;
        assembly {
            let srcpart := and(mload(src), not(mask))
            let destpart := and(mload(dest), mask)
            mstore(dest, or(destpart, srcpart))
        }
    }

    /// @notice Returns a memory pointer to the data portion of the provided bytes array.
    /// @param bts the bytes array to get a pointer from
    /// @return addr the pointer to the `bts` bytes array
    function dataPtr(bytes memory bts) internal pure returns (uint256 addr) {
        assembly {
            addr := add(bts, _BYTES_HEADER_SIZE)
        }
    }

    /// @notice Extracts a bytes array with length `howManyBytes` from `src`'s `offset` forward.
    /// @param src the bytes array to extract from
    /// @param offset where to start extracting from
    /// @param howManyBytes how many bytes we want to extract from `src`
    /// @return out the extracted bytes array
    /// @dev Extracting the 32-64th bytes out of a 64 bytes array takes ~7828 gas.
    function extractBytes(
        bytes memory src,
        uint256 offset,
        uint256 howManyBytes
    ) internal pure returns (bytes memory out) {
        require(
            offset + howManyBytes >= offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BYTES_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= offset + howManyBytes,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BYTES_OFFSET_OUT_OF_BOUNDS
                )
            )
        );
        out = new bytes(howManyBytes);
        uint256 start;

        assembly {
            start := add(add(src, offset), _BYTES_HEADER_SIZE)
        }

        copy(start, dataPtr(out), howManyBytes);
    }

    /// @notice Extracts a bytes32 extracted from `src`'s `offset` forward.
    /// @param src the source bytes array to extract from
    /// @param offset where to start extracting from
    /// @return out the bytes32 data extracted from `src`
    /// @dev ~439 gas
    function extractBytes32(bytes memory src, uint256 offset)
        internal
        pure
        returns (bytes32 out)
    {
        require(
            offset + 32 > offset,
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BYTES32_OFFSET_PARAMETER_OVERFLOW
                )
            )
        );
        require(
            src.length >= (offset + 32),
            string(
                abi.encodePacked(
                    BaseParserLibraryErrorCodes
                        .BASEPARSERLIB_BYTES32_OFFSET_OUT_OF_BOUNDS
                )
            )
        );
        assembly {
            out := mload(add(add(src, _BYTES_HEADER_SIZE), offset))
        }
    }
}

// File contracts/libraries/parsers/BClaimsParserLibrary.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

/// @title Library to parse the BClaims structure from a blob of capnproto data
library BClaimsParserLibrary {
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
        uint16 dataSectionSize = BaseParserLibrary.extractUInt16(
            src,
            dataOffset
        );
        require(
            dataSectionSize > 0 && dataSectionSize <= 2,
            string(
                abi.encodePacked(
                    BClaimsParserLibraryErrorCodes
                        .BCLAIMSPARSERLIB_SIZE_THRESHOLD_EXCEEDED
                )
            )
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
    function extractBClaims(bytes memory src)
        internal
        pure
        returns (BClaims memory bClaims)
    {
        return
            extractInnerBClaims(
                src,
                _CAPNPROTO_HEADER_SIZE,
                getPointerOffsetAdjustment(src, 4)
            );
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
            string(
                abi.encodePacked(
                    BClaimsParserLibraryErrorCodes
                        .BCLAIMSPARSERLIB_DATA_OFFSET_OVERFLOW
                )
            )
        );
        require(
            src.length >= dataOffset + _BCLAIMS_SIZE - pointerOffsetAdjustment,
            string(
                abi.encodePacked(
                    BClaimsParserLibraryErrorCodes
                        .BCLAIMSPARSERLIB_NOT_ENOUGH_BYTES
                )
            )
        );

        if (pointerOffsetAdjustment == 0) {
            bClaims.txCount = BaseParserLibrary.extractUInt32(
                src,
                dataOffset + 8
            );
        } else {
            // In case the txCount is 0, the value is not included in the binary
            // blob by capnproto.
            bClaims.txCount = 0;
        }

        bClaims.chainId = BaseParserLibrary.extractUInt32(src, dataOffset);
        require(
            bClaims.chainId > 0,
            string(
                abi.encodePacked(
                    BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_CHAINID_ZERO
                )
            )
        );
        bClaims.height = BaseParserLibrary.extractUInt32(src, dataOffset + 4);
        require(
            bClaims.height > 0,
            string(
                abi.encodePacked(
                    BClaimsParserLibraryErrorCodes.BCLAIMSPARSERLIB_HEIGHT_ZERO
                )
            )
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

// File contracts/interfaces/ISnapshots.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

struct Snapshot {
    uint256 committedAt;
    BClaimsParserLibrary.BClaims blockClaims;
}

interface ISnapshots {
    event SnapshotTaken(
        uint256 chainId,
        uint256 indexed epoch,
        uint256 height,
        address indexed validator,
        bool isSafeToProceedConsensus,
        bytes signatureRaw
    );

    function setSnapshotDesperationDelay(uint32 desperationDelay_) external;

    function setSnapshotDesperationFactor(uint32 desperationFactor_) external;

    function setMinimumIntervalBetweenSnapshots(
        uint32 minimumIntervalBetweenSnapshots_
    ) external;

    function snapshot(bytes calldata signatureGroup_, bytes calldata bClaims_)
        external
        returns (bool);

    function migrateSnapshots(
        bytes[] memory groupSignature_,
        bytes[] memory bClaims_
    ) external returns (bool);

    function getSnapshotDesperationDelay() external view returns (uint256);

    function getSnapshotDesperationFactor() external view returns (uint256);

    function getMinimumIntervalBetweenSnapshots()
        external
        view
        returns (uint256);

    function getChainId() external view returns (uint256);

    function getEpoch() external view returns (uint256);

    function getEpochLength() external view returns (uint256);

    function getChainIdFromSnapshot(uint256 epoch_)
        external
        view
        returns (uint256);

    function getChainIdFromLatestSnapshot() external view returns (uint256);

    function getBlockClaimsFromSnapshot(uint256 epoch_)
        external
        view
        returns (BClaimsParserLibrary.BClaims memory);

    function getBlockClaimsFromLatestSnapshot()
        external
        view
        returns (BClaimsParserLibrary.BClaims memory);

    function getCommittedHeightFromSnapshot(uint256 epoch_)
        external
        view
        returns (uint256);

    function getCommittedHeightFromLatestSnapshot()
        external
        view
        returns (uint256);

    function getAliceNetHeightFromSnapshot(uint256 epoch_)
        external
        view
        returns (uint256);

    function getAliceNetHeightFromLatestSnapshot()
        external
        view
        returns (uint256);

    function getSnapshot(uint256 epoch_)
        external
        view
        returns (Snapshot memory);

    function getLatestSnapshot() external view returns (Snapshot memory);

    function getEpochFromHeight(uint256 height) external view returns (uint256);

    function mayValidatorSnapshot(
        uint256 numValidators,
        uint256 myIdx,
        uint256 blocksSinceDesperation,
        bytes32 blsig,
        uint256 desperationFactor
    ) external pure returns (bool);
}

// File contracts/libraries/errorCodes/CustomEnumerableMapsErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library CustomEnumerableMapsErrorCodes {
    // CustomEnumerableMaps error codes
    bytes32 public constant CUSTOMENUMMAP_KEY_NOT_IN_MAP = "1900"; //"Error: Key not in the map!"
}

// File contracts/utils/CustomEnumerableMaps.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

struct ValidatorData {
    address _address;
    uint256 _tokenID;
}

struct ExitingValidatorData {
    uint128 _tokenID;
    uint128 _freeAfter;
}

struct ValidatorDataMap {
    ValidatorData[] _values;
    mapping(address => uint256) _indexes;
}

library CustomEnumerableMaps {
    /**
     * @dev Add a value to a map. O(1).
     *
     * Returns true if the value was added to the map, that is if it was not
     * already present.
     */
    function add(ValidatorDataMap storage map, ValidatorData memory value)
        internal
        returns (bool)
    {
        if (!contains(map, value._address)) {
            map._values.push(value);
            // The value is stored at length-1, but we add 1 to all indexes
            // and use 0 as a sentinel value
            map._indexes[value._address] = map._values.length;
            return true;
        } else {
            return false;
        }
    }

    /**
     * @dev Removes a value from a map using its address. O(1).
     *
     * Returns true if the value was removed from the map, that is if it was
     * present.
     */
    function remove(ValidatorDataMap storage map, address key)
        internal
        returns (bool)
    {
        // We read and store the key's index to prevent multiple reads from the same storage slot
        uint256 valueIndex = map._indexes[key];

        if (valueIndex != 0) {
            // Equivalent to contains(map, key)
            // To delete an element from the _values array in O(1), we swap the element to delete with the last one in
            // the array, and then remove the last element (sometimes called as 'swap and pop').
            // This modifies the order of the array, as noted in {at}.

            uint256 toDeleteIndex = valueIndex - 1;
            uint256 lastIndex = map._values.length - 1;

            if (lastIndex != toDeleteIndex) {
                ValidatorData memory lastValue = map._values[lastIndex];

                // Move the last value to the index where the value to delete is
                map._values[toDeleteIndex] = lastValue;
                // Update the index for the moved value
                map._indexes[lastValue._address] = valueIndex; // Replace lastValue's index to valueIndex
            }

            // Delete the slot where the moved key was stored
            map._values.pop();

            // Delete the index for the deleted slot
            delete map._indexes[key];

            return true;
        } else {
            return false;
        }
    }

    /**
     * @dev Returns true if the key is in the map. O(1).
     */
    function contains(ValidatorDataMap storage map, address key)
        internal
        view
        returns (bool)
    {
        return map._indexes[key] != 0;
    }

    /**
     * @dev Returns the number of values in the map. O(1).
     */
    function length(ValidatorDataMap storage map)
        internal
        view
        returns (uint256)
    {
        return map._values.length;
    }

    /**
     * @dev Returns the value stored at position `index` in the map. O(1).
     *
     * Note that there are no guarantees on the ordering of values inside the
     * array, and it may change when more values are added or removed.
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function at(ValidatorDataMap storage map, uint256 index)
        internal
        view
        returns (ValidatorData memory)
    {
        return map._values[index];
    }

    /**
     * @dev Tries to returns the value associated with `key`.  O(1).
     * Does not revert if `key` is not in the map.
     *
     */
    function tryGet(ValidatorDataMap storage map, address key)
        internal
        view
        returns (bool, ValidatorData memory)
    {
        uint256 index = map._indexes[key];
        if (index == 0) {
            return (false, ValidatorData(address(0), 0));
        } else {
            return (true, map._values[index - 1]);
        }
    }

    /**
     * @dev Returns the value associated with `key`.  O(1).
     *
     * Requirements:
     *
     * - `key` must be in the map.
     */
    function get(ValidatorDataMap storage map, address key)
        internal
        view
        returns (ValidatorData memory)
    {
        (bool success, ValidatorData memory value) = tryGet(map, key);
        require(
            success,
            string(
                abi.encodePacked(
                    CustomEnumerableMapsErrorCodes.CUSTOMENUMMAP_KEY_NOT_IN_MAP
                )
            )
        );
        return value;
    }

    /**
     * @dev Return the entire map in an array
     *
     * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
     * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
     * this function has an unbounded cost, and using it as part of a state-changing function may render the function
     * uncallable if the map grows to a point where copying to memory consumes too much gas to fit in a block.
     */
    function values(ValidatorDataMap storage map)
        internal
        view
        returns (ValidatorData[] memory)
    {
        return map._values;
    }

    /**
     * @dev Return the address of every entry in the entire map in an array
     *
     * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
     * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
     * this function has an unbounded cost, and using it as part of a state-changing function may render the function
     * uncallable if the map grows to a point where copying to memory consumes too much gas to fit in a block.
     */
    function addressValues(ValidatorDataMap storage map)
        internal
        view
        returns (address[] memory)
    {
        ValidatorData[] memory _values = values(map);
        address[] memory addresses = new address[](_values.length);
        for (uint256 i = 0; i < _values.length; i++) {
            addresses[i] = _values[i]._address;
        }
        return addresses;
    }
}

// File contracts/interfaces/IValidatorPool.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

interface IValidatorPool {
    event ValidatorJoined(
        address indexed account,
        uint256 validatorStakingTokenID
    );
    event ValidatorLeft(address indexed account, uint256 publicStakingTokenID);
    event ValidatorMinorSlashed(
        address indexed account,
        uint256 publicStakingTokenID
    );
    event ValidatorMajorSlashed(address indexed account);
    event MaintenanceScheduled();

    function setStakeAmount(uint256 stakeAmount_) external;

    function setMaxNumValidators(uint256 maxNumValidators_) external;

    function setDisputerReward(uint256 disputerReward_) external;

    function setLocation(string calldata ip) external;

    function scheduleMaintenance() external;

    function initializeETHDKG() external;

    function completeETHDKG() external;

    function pauseConsensus() external;

    function pauseConsensusOnArbitraryHeight(uint256 aliceNetHeight) external;

    function registerValidators(
        address[] calldata validators,
        uint256[] calldata publicStakingTokenIDs
    ) external;

    function unregisterValidators(address[] calldata validators) external;

    function unregisterAllValidators() external;

    function collectProfits()
        external
        returns (uint256 payoutEth, uint256 payoutToken);

    function claimExitingNFTPosition() external returns (uint256);

    function majorSlash(address dishonestValidator_, address disputer_)
        external;

    function minorSlash(address dishonestValidator_, address disputer_)
        external;

    function getValidatorsCount() external view returns (uint256);

    function getValidatorsAddresses() external view returns (address[] memory);

    function getValidator(uint256 index) external view returns (address);

    function getValidatorData(uint256 index)
        external
        view
        returns (ValidatorData memory);

    function getLocation(address validator)
        external
        view
        returns (string memory);

    function getLocations(address[] calldata validators_)
        external
        view
        returns (string[] memory);

    function getStakeAmount() external view returns (uint256);

    function getMaxNumValidators() external view returns (uint256);

    function getDisputerReward() external view returns (uint256);

    function tryGetTokenID(address account_)
        external
        view
        returns (
            bool,
            address,
            uint256
        );

    function isValidator(address participant) external view returns (bool);

    function isInExitingQueue(address participant) external view returns (bool);

    function isAccusable(address participant) external view returns (bool);

    function isMaintenanceScheduled() external view returns (bool);

    function isConsensusRunning() external view returns (bool);
}

// File contracts/utils/DeterministicAddress.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

abstract contract DeterministicAddress {
    function getMetamorphicContractAddress(bytes32 _salt, address _factory)
        public
        pure
        returns (address)
    {
        // byte code for metamorphic contract
        // 6020363636335afa1536363636515af43d36363e3d36f3
        bytes32 metamorphicContractBytecodeHash_ = 0x1c0bf703a3415cada9785e89e9d70314c3111ae7d8e04f33bb42eb1d264088be;
        return
            address(
                uint160(
                    uint256(
                        keccak256(
                            abi.encodePacked(
                                hex"ff",
                                _factory,
                                _salt,
                                metamorphicContractBytecodeHash_
                            )
                        )
                    )
                )
            );
    }
}

// File contracts/libraries/errorCodes/ImmutableAuthErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library ImmutableAuthErrorCodes {
    // ImmutableAuth error codes
    bytes32 public constant IMMUTEABLEAUTH_ONLY_FACTORY = "2000"; //"onlyFactory"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKEN = "2001"; //"onlyAToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_FOUNDATION = "2002"; //"onlyFoundation"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_GOVERNANCE = "2003"; // "onlyGovernance"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_LIQUIDITYPROVIDERSTAKING =
        "2004"; // "onlyLiquidityProviderStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_BTOKEN = "2005"; // "onlyBToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_MADTOKEN = "2006"; // "onlyMadToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_PUBLICSTAKING = "2007"; // "onlyPublicStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_SNAPSHOTS = "2008"; // "onlySnapshots"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_STAKINGPOSITIONDESCRIPTOR =
        "2009"; // "onlyStakingPositionDescriptor"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_VALIDATORPOOL = "2010"; // "onlyValidatorPool"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_VALIDATORSTAKING = "2011"; // "onlyValidatorStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKENBURNER = "2012"; // "onlyATokenBurner"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKENMINTER = "2013"; // "onlyATokenMinter"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKGACCUSATIONS = "2014"; // "onlyETHDKGAccusations"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKGPHASES = "2015"; // "onlyETHDKGPhases"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKG = "2016"; // "onlyETHDKG"
}

// File contracts/utils/ImmutableAuth.sol

// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;

    modifier onlyFactory() {
        require(
            msg.sender == _factory,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FACTORY
                )
            )
        );
        _;
    }

    constructor(address factory_) {
        _factory = factory_;
    }

    function _factoryAddress() internal view returns (address) {
        return _factory;
    }
}

abstract contract ImmutableAToken is ImmutableFactory {
    address private immutable _aToken;

    modifier onlyAToken() {
        require(
            msg.sender == _aToken,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKEN
                )
            )
        );
        _;
    }

    constructor() {
        _aToken = getMetamorphicContractAddress(
            0x41546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenAddress() internal view returns (address) {
        return _aToken;
    }

    function _saltForAToken() internal pure returns (bytes32) {
        return
            0x41546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenBurner is ImmutableFactory {
    address private immutable _aTokenBurner;

    modifier onlyATokenBurner() {
        require(
            msg.sender == _aTokenBurner,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENBURNER
                )
            )
        );
        _;
    }

    constructor() {
        _aTokenBurner = getMetamorphicContractAddress(
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenBurnerAddress() internal view returns (address) {
        return _aTokenBurner;
    }

    function _saltForATokenBurner() internal pure returns (bytes32) {
        return
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenMinter is ImmutableFactory {
    address private immutable _aTokenMinter;

    modifier onlyATokenMinter() {
        require(
            msg.sender == _aTokenMinter,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENMINTER
                )
            )
        );
        _;
    }

    constructor() {
        _aTokenMinter = getMetamorphicContractAddress(
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenMinterAddress() internal view returns (address) {
        return _aTokenMinter;
    }

    function _saltForATokenMinter() internal pure returns (bytes32) {
        return
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableBToken is ImmutableFactory {
    address private immutable _bToken;

    modifier onlyBToken() {
        require(
            msg.sender == _bToken,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_BTOKEN
                )
            )
        );
        _;
    }

    constructor() {
        _bToken = getMetamorphicContractAddress(
            0x42546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _bTokenAddress() internal view returns (address) {
        return _bToken;
    }

    function _saltForBToken() internal pure returns (bytes32) {
        return
            0x42546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableFoundation is ImmutableFactory {
    address private immutable _foundation;

    modifier onlyFoundation() {
        require(
            msg.sender == _foundation,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FOUNDATION
                )
            )
        );
        _;
    }

    constructor() {
        _foundation = getMetamorphicContractAddress(
            0x466f756e646174696f6e00000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _foundationAddress() internal view returns (address) {
        return _foundation;
    }

    function _saltForFoundation() internal pure returns (bytes32) {
        return
            0x466f756e646174696f6e00000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableGovernance is ImmutableFactory {
    address private immutable _governance;

    modifier onlyGovernance() {
        require(
            msg.sender == _governance,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_GOVERNANCE
                )
            )
        );
        _;
    }

    constructor() {
        _governance = getMetamorphicContractAddress(
            0x476f7665726e616e636500000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _governanceAddress() internal view returns (address) {
        return _governance;
    }

    function _saltForGovernance() internal pure returns (bytes32) {
        return
            0x476f7665726e616e636500000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableLiquidityProviderStaking is ImmutableFactory {
    address private immutable _liquidityProviderStaking;

    modifier onlyLiquidityProviderStaking() {
        require(
            msg.sender == _liquidityProviderStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_LIQUIDITYPROVIDERSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _liquidityProviderStaking = getMetamorphicContractAddress(
            0x4c697175696469747950726f76696465725374616b696e670000000000000000,
            _factoryAddress()
        );
    }

    function _liquidityProviderStakingAddress()
        internal
        view
        returns (address)
    {
        return _liquidityProviderStaking;
    }

    function _saltForLiquidityProviderStaking()
        internal
        pure
        returns (bytes32)
    {
        return
            0x4c697175696469747950726f76696465725374616b696e670000000000000000;
    }
}

abstract contract ImmutablePublicStaking is ImmutableFactory {
    address private immutable _publicStaking;

    modifier onlyPublicStaking() {
        require(
            msg.sender == _publicStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_PUBLICSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _publicStaking = getMetamorphicContractAddress(
            0x5075626c69635374616b696e6700000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _publicStakingAddress() internal view returns (address) {
        return _publicStaking;
    }

    function _saltForPublicStaking() internal pure returns (bytes32) {
        return
            0x5075626c69635374616b696e6700000000000000000000000000000000000000;
    }
}

abstract contract ImmutableSnapshots is ImmutableFactory {
    address private immutable _snapshots;

    modifier onlySnapshots() {
        require(
            msg.sender == _snapshots,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_SNAPSHOTS
                )
            )
        );
        _;
    }

    constructor() {
        _snapshots = getMetamorphicContractAddress(
            0x536e617073686f74730000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _snapshotsAddress() internal view returns (address) {
        return _snapshots;
    }

    function _saltForSnapshots() internal pure returns (bytes32) {
        return
            0x536e617073686f74730000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakingPositionDescriptor is ImmutableFactory {
    address private immutable _stakingPositionDescriptor;

    modifier onlyStakingPositionDescriptor() {
        require(
            msg.sender == _stakingPositionDescriptor,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_STAKINGPOSITIONDESCRIPTOR
                )
            )
        );
        _;
    }

    constructor() {
        _stakingPositionDescriptor = getMetamorphicContractAddress(
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000,
            _factoryAddress()
        );
    }

    function _stakingPositionDescriptorAddress()
        internal
        view
        returns (address)
    {
        return _stakingPositionDescriptor;
    }

    function _saltForStakingPositionDescriptor()
        internal
        pure
        returns (bytes32)
    {
        return
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000;
    }
}

abstract contract ImmutableValidatorPool is ImmutableFactory {
    address private immutable _validatorPool;

    modifier onlyValidatorPool() {
        require(
            msg.sender == _validatorPool,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORPOOL
                )
            )
        );
        _;
    }

    constructor() {
        _validatorPool = getMetamorphicContractAddress(
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorPoolAddress() internal view returns (address) {
        return _validatorPool;
    }

    function _saltForValidatorPool() internal pure returns (bytes32) {
        return
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000;
    }
}

abstract contract ImmutableValidatorStaking is ImmutableFactory {
    address private immutable _validatorStaking;

    modifier onlyValidatorStaking() {
        require(
            msg.sender == _validatorStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _validatorStaking = getMetamorphicContractAddress(
            0x56616c696461746f725374616b696e6700000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorStakingAddress() internal view returns (address) {
        return _validatorStaking;
    }

    function _saltForValidatorStaking() internal pure returns (bytes32) {
        return
            0x56616c696461746f725374616b696e6700000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGAccusations is ImmutableFactory {
    address private immutable _ethdkgAccusations;

    modifier onlyETHDKGAccusations() {
        require(
            msg.sender == _ethdkgAccusations,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_ETHDKGACCUSATIONS
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkgAccusations = getMetamorphicContractAddress(
            0x455448444b4741636375736174696f6e73000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAccusationsAddress() internal view returns (address) {
        return _ethdkgAccusations;
    }

    function _saltForETHDKGAccusations() internal pure returns (bytes32) {
        return
            0x455448444b4741636375736174696f6e73000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGPhases is ImmutableFactory {
    address private immutable _ethdkgPhases;

    modifier onlyETHDKGPhases() {
        require(
            msg.sender == _ethdkgPhases,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKGPHASES
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkgPhases = getMetamorphicContractAddress(
            0x455448444b475068617365730000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgPhasesAddress() internal view returns (address) {
        return _ethdkgPhases;
    }

    function _saltForETHDKGPhases() internal pure returns (bytes32) {
        return
            0x455448444b475068617365730000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKG is ImmutableFactory {
    address private immutable _ethdkg;

    modifier onlyETHDKG() {
        require(
            msg.sender == _ethdkg,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKG
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkg = getMetamorphicContractAddress(
            0x455448444b470000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAddress() internal view returns (address) {
        return _ethdkg;
    }

    function _saltForETHDKG() internal pure returns (bytes32) {
        return
            0x455448444b470000000000000000000000000000000000000000000000000000;
    }
}

// File @openzeppelin/contracts-upgradeable/utils/AddressUpgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (utils/Address.sol)

pragma solidity ^0.8.1;

/**
 * @dev Collection of functions related to the address type
 */
library AddressUpgradeable {
    /**
     * @dev Returns true if `account` is a contract.
     *
     * [IMPORTANT]
     * ====
     * It is unsafe to assume that an address for which this function returns
     * false is an externally-owned account (EOA) and not a contract.
     *
     * Among others, `isContract` will return false for the following
     * types of addresses:
     *
     *  - an externally-owned account
     *  - a contract in construction
     *  - an address where a contract will be created
     *  - an address where a contract lived, but was destroyed
     * ====
     *
     * [IMPORTANT]
     * ====
     * You shouldn't rely on `isContract` to protect against flash loan attacks!
     *
     * Preventing calls from contracts is highly discouraged. It breaks composability, breaks support for smart wallets
     * like Gnosis Safe, and does not provide security since it can be circumvented by calling from a contract
     * constructor.
     * ====
     */
    function isContract(address account) internal view returns (bool) {
        // This method relies on extcodesize/address.code.length, which returns 0
        // for contracts in construction, since the code is only stored at the end
        // of the constructor execution.

        return account.code.length > 0;
    }

    /**
     * @dev Replacement for Solidity's `transfer`: sends `amount` wei to
     * `recipient`, forwarding all available gas and reverting on errors.
     *
     * https://eips.ethereum.org/EIPS/eip-1884[EIP1884] increases the gas cost
     * of certain opcodes, possibly making contracts go over the 2300 gas limit
     * imposed by `transfer`, making them unable to receive funds via
     * `transfer`. {sendValue} removes this limitation.
     *
     * https://diligence.consensys.net/posts/2019/09/stop-using-soliditys-transfer-now/[Learn more].
     *
     * IMPORTANT: because control is transferred to `recipient`, care must be
     * taken to not create reentrancy vulnerabilities. Consider using
     * {ReentrancyGuard} or the
     * https://solidity.readthedocs.io/en/v0.5.11/security-considerations.html#use-the-checks-effects-interactions-pattern[checks-effects-interactions pattern].
     */
    function sendValue(address payable recipient, uint256 amount) internal {
        require(
            address(this).balance >= amount,
            "Address: insufficient balance"
        );

        (bool success, ) = recipient.call{value: amount}("");
        require(
            success,
            "Address: unable to send value, recipient may have reverted"
        );
    }

    /**
     * @dev Performs a Solidity function call using a low level `call`. A
     * plain `call` is an unsafe replacement for a function call: use this
     * function instead.
     *
     * If `target` reverts with a revert reason, it is bubbled up by this
     * function (like regular Solidity function calls).
     *
     * Returns the raw returned data. To convert to the expected return value,
     * use https://solidity.readthedocs.io/en/latest/units-and-global-variables.html?highlight=abi.decode#abi-encoding-and-decoding-functions[`abi.decode`].
     *
     * Requirements:
     *
     * - `target` must be a contract.
     * - calling `target` with `data` must not revert.
     *
     * _Available since v3.1._
     */
    function functionCall(address target, bytes memory data)
        internal
        returns (bytes memory)
    {
        return functionCall(target, data, "Address: low-level call failed");
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`], but with
     * `errorMessage` as a fallback revert reason when `target` reverts.
     *
     * _Available since v3.1._
     */
    function functionCall(
        address target,
        bytes memory data,
        string memory errorMessage
    ) internal returns (bytes memory) {
        return functionCallWithValue(target, data, 0, errorMessage);
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
     * but also transferring `value` wei to `target`.
     *
     * Requirements:
     *
     * - the calling contract must have an ETH balance of at least `value`.
     * - the called Solidity function must be `payable`.
     *
     * _Available since v3.1._
     */
    function functionCallWithValue(
        address target,
        bytes memory data,
        uint256 value
    ) internal returns (bytes memory) {
        return
            functionCallWithValue(
                target,
                data,
                value,
                "Address: low-level call with value failed"
            );
    }

    /**
     * @dev Same as {xref-Address-functionCallWithValue-address-bytes-uint256-}[`functionCallWithValue`], but
     * with `errorMessage` as a fallback revert reason when `target` reverts.
     *
     * _Available since v3.1._
     */
    function functionCallWithValue(
        address target,
        bytes memory data,
        uint256 value,
        string memory errorMessage
    ) internal returns (bytes memory) {
        require(
            address(this).balance >= value,
            "Address: insufficient balance for call"
        );
        require(isContract(target), "Address: call to non-contract");

        (bool success, bytes memory returndata) = target.call{value: value}(
            data
        );
        return verifyCallResult(success, returndata, errorMessage);
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
     * but performing a static call.
     *
     * _Available since v3.3._
     */
    function functionStaticCall(address target, bytes memory data)
        internal
        view
        returns (bytes memory)
    {
        return
            functionStaticCall(
                target,
                data,
                "Address: low-level static call failed"
            );
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-string-}[`functionCall`],
     * but performing a static call.
     *
     * _Available since v3.3._
     */
    function functionStaticCall(
        address target,
        bytes memory data,
        string memory errorMessage
    ) internal view returns (bytes memory) {
        require(isContract(target), "Address: static call to non-contract");

        (bool success, bytes memory returndata) = target.staticcall(data);
        return verifyCallResult(success, returndata, errorMessage);
    }

    /**
     * @dev Tool to verifies that a low level call was successful, and revert if it wasn't, either by bubbling the
     * revert reason using the provided one.
     *
     * _Available since v4.3._
     */
    function verifyCallResult(
        bool success,
        bytes memory returndata,
        string memory errorMessage
    ) internal pure returns (bytes memory) {
        if (success) {
            return returndata;
        } else {
            // Look for revert reason and bubble it up if present
            if (returndata.length > 0) {
                // The easiest way to bubble the revert reason is using memory via assembly

                assembly {
                    let returndata_size := mload(returndata)
                    revert(add(32, returndata), returndata_size)
                }
            } else {
                revert(errorMessage);
            }
        }
    }
}

// File @openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (proxy/utils/Initializable.sol)

pragma solidity ^0.8.0;

/**
 * @dev This is a base contract to aid in writing upgradeable contracts, or any kind of contract that will be deployed
 * behind a proxy. Since proxied contracts do not make use of a constructor, it's common to move constructor logic to an
 * external initializer function, usually called `initialize`. It then becomes necessary to protect this initializer
 * function so it can only be called once. The {initializer} modifier provided by this contract will have this effect.
 *
 * TIP: To avoid leaving the proxy in an uninitialized state, the initializer function should be called as early as
 * possible by providing the encoded function call as the `_data` argument to {ERC1967Proxy-constructor}.
 *
 * CAUTION: When used with inheritance, manual care must be taken to not invoke a parent initializer twice, or to ensure
 * that all initializers are idempotent. This is not verified automatically as constructors are by Solidity.
 *
 * [CAUTION]
 * ====
 * Avoid leaving a contract uninitialized.
 *
 * An uninitialized contract can be taken over by an attacker. This applies to both a proxy and its implementation
 * contract, which may impact the proxy. To initialize the implementation contract, you can either invoke the
 * initializer manually, or you can include a constructor to automatically mark it as initialized when it is deployed:
 *
 * [.hljs-theme-light.nopadding]
 * ```
 * /// @custom:oz-upgrades-unsafe-allow constructor
 * constructor() initializer {}
 * ```
 * ====
 */
abstract contract Initializable {
    /**
     * @dev Indicates that the contract has been initialized.
     */
    bool private _initialized;

    /**
     * @dev Indicates that the contract is in the process of being initialized.
     */
    bool private _initializing;

    /**
     * @dev Modifier to protect an initializer function from being invoked twice.
     */
    modifier initializer() {
        // If the contract is initializing we ignore whether _initialized is set in order to support multiple
        // inheritance patterns, but we only do this in the context of a constructor, because in other contexts the
        // contract may have been reentered.
        require(
            _initializing ? _isConstructor() : !_initialized,
            "Initializable: contract is already initialized"
        );

        bool isTopLevelCall = !_initializing;
        if (isTopLevelCall) {
            _initializing = true;
            _initialized = true;
        }

        _;

        if (isTopLevelCall) {
            _initializing = false;
        }
    }

    /**
     * @dev Modifier to protect an initialization function so that it can only be invoked by functions with the
     * {initializer} modifier, directly or indirectly.
     */
    modifier onlyInitializing() {
        require(_initializing, "Initializable: contract is not initializing");
        _;
    }

    function _isConstructor() private view returns (bool) {
        return !AddressUpgradeable.isContract(address(this));
    }
}

// File contracts/libraries/ethdkg/ETHDKGStorage.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

enum Phase {
    RegistrationOpen,
    ShareDistribution,
    DisputeShareDistribution,
    KeyShareSubmission,
    MPKSubmission,
    GPKJSubmission,
    DisputeGPKJSubmission,
    Completion
}

// State of key generation
struct Participant {
    uint256[2] publicKey;
    uint64 nonce;
    uint64 index;
    Phase phase;
    bytes32 distributedSharesHash;
    uint256[2] commitmentsFirstCoefficient;
    uint256[2] keyShares;
    uint256[4] gpkj;
}

abstract contract ETHDKGStorage is
    Initializable,
    ImmutableFactory,
    ImmutableSnapshots,
    ImmutableValidatorPool
{
    // ISnapshots internal immutable _snapshots;
    // IValidatorPool internal immutable _validatorPool;
    //address internal immutable _factory;
    uint256 internal constant _MIN_VALIDATORS = 4;

    uint64 internal _nonce;
    uint64 internal _phaseStartBlock;
    Phase internal _ethdkgPhase;
    uint32 internal _numParticipants;
    uint16 internal _badParticipants;
    uint16 internal _phaseLength;
    uint16 internal _confirmationLength;

    // AliceNet height used to start the new validator set in arbitrary height points if the AliceNet
    // Consensus is halted
    uint256 internal _customAliceNetHeight;

    address internal _admin;

    uint256[4] internal _masterPublicKey;
    uint256[2] internal _mpkG1;
    bytes32 internal _masterPublicKeyHash;

    mapping(address => Participant) internal _participants;

    constructor()
        ImmutableFactory(msg.sender)
        ImmutableSnapshots()
        ImmutableValidatorPool()
    {}
}

// File contracts/interfaces/IETHDKG.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

interface IETHDKG {
    function setPhaseLength(uint16 phaseLength_) external;

    function setConfirmationLength(uint16 confirmationLength_) external;

    function setCustomAliceNetHeight(uint256 aliceNetHeight) external;

    function initializeETHDKG() external;

    function register(uint256[2] memory publicKey) external;

    function distributeShares(
        uint256[] memory encryptedShares,
        uint256[2][] memory commitments
    ) external;

    function submitKeyShare(
        uint256[2] memory keyShareG1,
        uint256[2] memory keyShareG1CorrectnessProof,
        uint256[4] memory keyShareG2
    ) external;

    function submitMasterPublicKey(uint256[4] memory masterPublicKey_) external;

    function submitGPKJ(uint256[4] memory gpkj) external;

    function complete() external;

    function migrateValidators(
        address[] memory validatorsAccounts_,
        uint256[] memory validatorIndexes_,
        uint256[4][] memory validatorShares_,
        uint8 validatorCount_,
        uint256 epoch_,
        uint256 sideChainHeight_,
        uint256 ethHeight_,
        uint256[4] memory masterPublicKey_
    ) external;

    function accuseParticipantNotRegistered(address[] memory dishonestAddresses)
        external;

    function accuseParticipantDidNotDistributeShares(
        address[] memory dishonestAddresses
    ) external;

    function accuseParticipantDistributedBadShares(
        address dishonestAddress,
        uint256[] memory encryptedShares,
        uint256[2][] memory commitments,
        uint256[2] memory sharedKey,
        uint256[2] memory sharedKeyCorrectnessProof
    ) external;

    function accuseParticipantDidNotSubmitKeyShares(
        address[] memory dishonestAddresses
    ) external;

    function accuseParticipantDidNotSubmitGPKJ(
        address[] memory dishonestAddresses
    ) external;

    function accuseParticipantSubmittedBadGPKJ(
        address[] memory validators,
        bytes32[] memory encryptedSharesHash,
        uint256[2][][] memory commitments,
        address dishonestAddress
    ) external;

    function isETHDKGRunning() external view returns (bool);

    function isMasterPublicKeySet() external view returns (bool);

    function getNonce() external view returns (uint256);

    function getPhaseStartBlock() external view returns (uint256);

    function getPhaseLength() external view returns (uint256);

    function getConfirmationLength() external view returns (uint256);

    function getETHDKGPhase() external view returns (Phase);

    function getNumParticipants() external view returns (uint256);

    function getBadParticipants() external view returns (uint256);

    function getMinValidators() external view returns (uint256);

    function getParticipantInternalState(address participant)
        external
        view
        returns (Participant memory);

    function getMasterPublicKey() external view returns (uint256[4] memory);

    function getMasterPublicKeyHash() external view returns (bytes32);

    function tryGetParticipantIndex(address participant)
        external
        view
        returns (bool, uint256);
}

// File contracts/utils/EthSafeTransfer.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract EthSafeTransfer {
    /// @notice _safeTransferEth performs a transfer of Eth using the call
    /// method / this function is resistant to breaking gas price changes and /
    /// performs call in a safe manner by reverting on failure. / this function
    /// will return without performing a call or reverting, / if amount_ is zero
    function _safeTransferEth(address to_, uint256 amount_) internal {
        if (amount_ == 0) {
            return;
        }
        require(
            to_ != address(0),
            "EthSafeTransfer: cannot transfer ETH to address 0x0"
        );
        address payable caller = payable(to_);
        (bool success, ) = caller.call{value: amount_}("");
        require(success, "EthSafeTransfer: Transfer failed.");
    }
}

// File contracts/utils/ERC20SafeTransfer.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract ERC20SafeTransfer {
    // _safeTransferFromERC20 performs a transferFrom call against an erc20 contract in a safe manner
    // by reverting on failure
    // this function will return without performing a call or reverting
    // if amount_ is zero
    function _safeTransferFromERC20(
        IERC20Transferable contract_,
        address sender_,
        uint256 amount_
    ) internal {
        if (amount_ == 0) {
            return;
        }
        require(
            address(contract_) != address(0x0),
            "ERC20SafeTransfer: Cannot call methods on contract address 0x0."
        );
        bool success = contract_.transferFrom(sender_, address(this), amount_);
        require(success, "ERC20SafeTransfer: Transfer failed.");
    }

    // _safeTransferERC20 performs a transfer call against an erc20 contract in a safe manner
    // by reverting on failure
    // this function will return without performing a call or reverting
    // if amount_ is zero
    function _safeTransferERC20(
        IERC20Transferable contract_,
        address to_,
        uint256 amount_
    ) internal {
        if (amount_ == 0) {
            return;
        }
        require(
            address(contract_) != address(0x0),
            "ERC20SafeTransfer: Cannot call methods on contract address 0x0."
        );
        bool success = contract_.transfer(to_, amount_);
        require(success, "ERC20SafeTransfer: Transfer failed.");
    }
}

// File contracts/libraries/errorCodes/MagicValueErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library MagicValueErrorCodes {
    // MagicValue error codes
    bytes32 public constant MAGICVALUE_BAD_MAGIC = "2200"; //"BAD MAGIC"
}

// File contracts/utils/MagicValue.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract MagicValue {
    // _MAGIC_VALUE is a constant that may be used to prevent
    // a user from calling a dangerous method without significant
    // effort or ( hopefully ) reading the code to understand the risk
    uint8 internal constant _MAGIC_VALUE = 42;

    modifier checkMagic(uint8 magic_) {
        require(
            magic_ == _getMagic(),
            string(abi.encodePacked(MagicValueErrorCodes.MAGICVALUE_BAD_MAGIC))
        );
        _;
    }

    // _getMagic returns the magic constant
    function _getMagic() internal pure returns (uint8) {
        return _MAGIC_VALUE;
    }
}

// File contracts/libraries/validatorPool/ValidatorPoolStorage.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

abstract contract ValidatorPoolStorage is
    ImmutableFactory,
    ImmutableSnapshots,
    ImmutableETHDKG,
    ImmutablePublicStaking,
    ImmutableValidatorStaking,
    ImmutableAToken
{
    // POSITION_LOCK_PERIOD describes the maximum interval a PublicStaking Position may be locked after
    // being given back to validator exiting the pool
    uint256 public constant POSITION_LOCK_PERIOD = 172800;
    // Interval in AliceNet Epochs that a validator exiting the pool should before claiming is
    // PublicStaking position
    uint256 public constant CLAIM_PERIOD = 3;

    // Maximum number the ethereum blocks allowed without a validator committing a snapshot
    uint256 public constant MAX_INTERVAL_WITHOUT_SNAPSHOTS = 8192;

    // Minimum amount to stake
    uint256 internal _stakeAmount;
    // Max number of validators allowed in the pool
    uint256 internal _maxNumValidators;
    // Value in WEIs to be discounted of dishonest validator in case of slashing event. This value
    // is usually sent back to the disputer
    uint256 internal _disputerReward;

    // Boolean flag to be read by the snapshot contract in order to decide if the validator set
    // needs to be changed or not (i.e if a validator is going to be removed or added).
    bool internal _isMaintenanceScheduled;
    // Boolean flag to keep track if the consensus is running in the side chain or not. Validators
    // can only join or leave the pool in case this value is false.
    bool internal _isConsensusRunning;

    // The internal iterable mapping that tracks all ACTIVE validators in the Pool
    ValidatorDataMap internal _validators;

    // Mapping that keeps track of the validators leaving the Pool. Validators assets are hold by
    // `CLAIM_PERIOD` epochs before the user being able to claim the assets back in the form a new
    // PublicStaking position.
    mapping(address => ExitingValidatorData) internal _exitingValidatorsData;

    // Mapping to keep track of the active validators IPs.
    mapping(address => string) internal _ipLocations;

    constructor()
        ImmutableFactory(msg.sender)
        ImmutableSnapshots()
        ImmutableETHDKG()
        ImmutablePublicStaking()
        ImmutableValidatorStaking()
        ImmutableAToken()
    {}
}

// File @openzeppelin/contracts/token/ERC721/IERC721Receiver.sol@v4.5.0

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (token/ERC721/IERC721Receiver.sol)

pragma solidity ^0.8.0;

/**
 * @title ERC721 token receiver interface
 * @dev Interface for any contract that wants to support safeTransfers
 * from ERC721 asset contracts.
 */
interface IERC721Receiver {
    /**
     * @dev Whenever an {IERC721} `tokenId` token is transferred to this contract via {IERC721-safeTransferFrom}
     * by `operator` from `from`, this function is called.
     *
     * It must return its Solidity selector to confirm the token transfer.
     * If any other value is returned or the interface is not implemented by the recipient, the transfer will be reverted.
     *
     * The selector can be obtained in Solidity with `IERC721.onERC721Received.selector`.
     */
    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4);
}

// File @openzeppelin/contracts/token/ERC721/utils/ERC721Holder.sol@v4.5.0

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (token/ERC721/utils/ERC721Holder.sol)

pragma solidity ^0.8.0;

/**
 * @dev Implementation of the {IERC721Receiver} interface.
 *
 * Accepts all token transfers.
 * Make sure the contract is able to use its token with {IERC721-safeTransferFrom}, {IERC721-approve} or {IERC721-setApprovalForAll}.
 */
contract ERC721Holder is IERC721Receiver {
    /**
     * @dev See {IERC721Receiver-onERC721Received}.
     *
     * Always returns `IERC721Receiver.onERC721Received.selector`.
     */
    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onERC721Received.selector;
    }
}

// File @openzeppelin/contracts/utils/introspection/IERC165.sol@v4.5.0

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/introspection/IERC165.sol)

pragma solidity ^0.8.0;

/**
 * @dev Interface of the ERC165 standard, as defined in the
 * https://eips.ethereum.org/EIPS/eip-165[EIP].
 *
 * Implementers can declare support of contract interfaces, which can then be
 * queried by others ({ERC165Checker}).
 *
 * For an implementation, see {ERC165}.
 */
interface IERC165 {
    /**
     * @dev Returns true if this contract implements the interface defined by
     * `interfaceId`. See the corresponding
     * https://eips.ethereum.org/EIPS/eip-165#how-interfaces-are-identified[EIP section]
     * to learn more about how these ids are created.
     *
     * This function call must use less than 30 000 gas.
     */
    function supportsInterface(bytes4 interfaceId) external view returns (bool);
}

// File @openzeppelin/contracts/token/ERC721/IERC721.sol@v4.5.0

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (token/ERC721/IERC721.sol)

pragma solidity ^0.8.0;

/**
 * @dev Required interface of an ERC721 compliant contract.
 */
interface IERC721 is IERC165 {
    /**
     * @dev Emitted when `tokenId` token is transferred from `from` to `to`.
     */
    event Transfer(
        address indexed from,
        address indexed to,
        uint256 indexed tokenId
    );

    /**
     * @dev Emitted when `owner` enables `approved` to manage the `tokenId` token.
     */
    event Approval(
        address indexed owner,
        address indexed approved,
        uint256 indexed tokenId
    );

    /**
     * @dev Emitted when `owner` enables or disables (`approved`) `operator` to manage all of its assets.
     */
    event ApprovalForAll(
        address indexed owner,
        address indexed operator,
        bool approved
    );

    /**
     * @dev Returns the number of tokens in ``owner``'s account.
     */
    function balanceOf(address owner) external view returns (uint256 balance);

    /**
     * @dev Returns the owner of the `tokenId` token.
     *
     * Requirements:
     *
     * - `tokenId` must exist.
     */
    function ownerOf(uint256 tokenId) external view returns (address owner);

    /**
     * @dev Safely transfers `tokenId` token from `from` to `to`, checking first that contract recipients
     * are aware of the ERC721 protocol to prevent tokens from being forever locked.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must exist and be owned by `from`.
     * - If the caller is not `from`, it must be have been allowed to move this token by either {approve} or {setApprovalForAll}.
     * - If `to` refers to a smart contract, it must implement {IERC721Receiver-onERC721Received}, which is called upon a safe transfer.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId
    ) external;

    /**
     * @dev Transfers `tokenId` token from `from` to `to`.
     *
     * WARNING: Usage of this method is discouraged, use {safeTransferFrom} whenever possible.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must be owned by `from`.
     * - If the caller is not `from`, it must be approved to move this token by either {approve} or {setApprovalForAll}.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 tokenId
    ) external;

    /**
     * @dev Gives permission to `to` to transfer `tokenId` token to another account.
     * The approval is cleared when the token is transferred.
     *
     * Only a single account can be approved at a time, so approving the zero address clears previous approvals.
     *
     * Requirements:
     *
     * - The caller must own the token or be an approved operator.
     * - `tokenId` must exist.
     *
     * Emits an {Approval} event.
     */
    function approve(address to, uint256 tokenId) external;

    /**
     * @dev Returns the account approved for `tokenId` token.
     *
     * Requirements:
     *
     * - `tokenId` must exist.
     */
    function getApproved(uint256 tokenId)
        external
        view
        returns (address operator);

    /**
     * @dev Approve or remove `operator` as an operator for the caller.
     * Operators can call {transferFrom} or {safeTransferFrom} for any token owned by the caller.
     *
     * Requirements:
     *
     * - The `operator` cannot be the caller.
     *
     * Emits an {ApprovalForAll} event.
     */
    function setApprovalForAll(address operator, bool _approved) external;

    /**
     * @dev Returns if the `operator` is allowed to manage all of the assets of `owner`.
     *
     * See {setApprovalForAll}
     */
    function isApprovedForAll(address owner, address operator)
        external
        view
        returns (bool);

    /**
     * @dev Safely transfers `tokenId` token from `from` to `to`.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must exist and be owned by `from`.
     * - If the caller is not `from`, it must be approved to move this token by either {approve} or {setApprovalForAll}.
     * - If `to` refers to a smart contract, it must implement {IERC721Receiver-onERC721Received}, which is called upon a safe transfer.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId,
        bytes calldata data
    ) external;
}

// File contracts/libraries/errorCodes/ValidatorPoolErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library ValidatorPoolErrorCodes {
    // ValidatorPool error codes
    bytes32 public constant VALIDATORPOOL_CALLER_NOT_VALIDATOR = "800"; //"ValidatorPool: Only validators allowed!"
    bytes32 public constant VALIDATORPOOL_CONSENSUS_RUNNING = "801"; //"ValidatorPool: Error AliceNet Consensus should be halted!"
    bytes32 public constant VALIDATORPOOL_ETHDKG_ROUND_RUNNING = "802"; //"ValidatorPool: There's an ETHDKG round running!"
    bytes32 public constant VALIDATORPOOL_ONLY_CONTRACTS_ALLOWED = "803"; //"Only NFT contracts allowed to send ethereum!"
    bytes32 public constant VALIDATORPOOL_MIN_BLOCK_INTERVAL_NOT_MET = "804"; //"ValidatorPool: Condition not met to stop consensus!"
    bytes32 public constant VALIDATORPOOL_MAX_VALIDATORS_MET = "805"; //"ValidatorPool: There are not enough free spots for all new validators!"
    bytes32
        public constant VALIDATORPOOL_REGISTRATION_PARAMETER_LENGTH_MISMATCH =
        "806"; //"ValidatorPool: Both input array should have same length!"
    bytes32 public constant VALIDATORPOOL_FACTORY_SHOULD_OWN_POSITION = "807"; //"ValidatorPool: The factory should be the owner of the StakeNFT position!"
    bytes32 public constant VALIDATORPOOL_VALIDATORS_GREATER_THAN_AVAILABLE =
        "808"; //"ValidatorPool: There are not enough validators to be removed!"
    bytes32
        public constant VALIDATORPOOL_PROFITS_ONLY_CLAIMABLE_DURING_CONSENSUS =
        "809"; //"ValidatorPool: Profits can only be claimable when consensus is running!"
    bytes32 public constant VALIDATORPOOL_TOKEN_BALANCE_CHANGED = "810"; //"ValidatorPool: Invalid transaction, token balance of the contract changed!"
    bytes32 public constant VALIDATORPOOL_ETH_BALANCE_CHANGED = "811"; //"ValidatorPool: Invalid transaction, eth balance of the contract changed!"
    bytes32 public constant VALIDATORPOOL_SENDER_NOT_IN_EXITING_QUEUE = "812"; //"ValidatorPool: Address not in the exitingQueue!"
    bytes32 public constant VALIDATORPOOL_WAITING_PERIOD_NOT_MET = "813"; //"ValidatorPool: The waiting period is not over yet!"
    bytes32 public constant VALIDATORPOOL_DISHONEST_VALIDATOR_NOT_ACCUSABLE =
        "814"; //"ValidatorPool: DishonestValidator should be a validator or be in the exiting line!"
    bytes32 public constant VALIDATORPOOL_INVALID_INDEX = "815"; //"Index out boundaries!"
    bytes32 public constant VALIDATORPOOL_ADDRESS_ALREADY_VALIDATOR = "816"; // "ValidatorPool: Address is already a validator or it is in the exiting line!"
    bytes32 public constant VALIDATORPOOL_ADDRESS_NOT_VALIDATOR = "817"; // "ValidatorPool: Address is not a validator_!"
    bytes32 public constant VALIDATORPOOL_MINIMUM_STAKE_NOT_MET = "818"; // "ValidatorStakeNFT: Error, the Stake position doesn't have enough funds!"
    bytes32 public constant VALIDATORPOOL_PAYOUT_TOO_LOW = "819"; // "ValidatorPool: Miner shares greater then the total payout in tokens!"
    bytes32 public constant VALIDATORPOOL_ADDRESS_NOT_ACCUSABLE = "820"; // "ValidatorPool: Address is not accusable!"
    bytes32 public constant VALIDATORPOOL_INSUFFICIENT_FUNDS_IN_STAKE_POSITION =
        "821"; // "ValidatorPool: Error, the Stake position doesn't have enough funds!"
}

// File contracts/ValidatorPool.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

/// @custom:salt ValidatorPool
/// @custom:deploy-type deployUpgradeable
contract ValidatorPool is
    Initializable,
    ValidatorPoolStorage,
    IValidatorPool,
    MagicValue,
    EthSafeTransfer,
    ERC20SafeTransfer,
    ERC721Holder
{
    using CustomEnumerableMaps for ValidatorDataMap;

    modifier onlyValidator() {
        require(
            _isValidator(msg.sender),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_CALLER_NOT_VALIDATOR
                )
            )
        );
        _;
    }

    modifier assertNotConsensusRunning() {
        require(
            !_isConsensusRunning,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_CONSENSUS_RUNNING
                )
            )
        );
        _;
    }

    modifier assertNotETHDKGRunning() {
        require(
            !IETHDKG(_ethdkgAddress()).isETHDKGRunning(),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETHDKG_ROUND_RUNNING
                )
            )
        );
        _;
    }

    constructor() ValidatorPoolStorage() {}

    receive() external payable {
        require(
            msg.sender == _validatorStakingAddress() ||
                msg.sender == _publicStakingAddress(),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ONLY_CONTRACTS_ALLOWED
                )
            )
        );
    }

    function initialize(
        uint256 stakeAmount_,
        uint256 maxNumValidators_,
        uint256 disputerReward_
    ) public onlyFactory initializer {
        _stakeAmount = stakeAmount_;
        _maxNumValidators = maxNumValidators_;
        _disputerReward = disputerReward_;
    }

    function setStakeAmount(uint256 stakeAmount_) public onlyFactory {
        _stakeAmount = stakeAmount_;
    }

    function setMaxNumValidators(uint256 maxNumValidators_) public onlyFactory {
        _maxNumValidators = maxNumValidators_;
    }

    function setDisputerReward(uint256 disputerReward_) public onlyFactory {
        _disputerReward = disputerReward_;
    }

    function setLocation(string calldata ip_) public onlyValidator {
        _ipLocations[msg.sender] = ip_;
    }

    function scheduleMaintenance() public onlyFactory {
        _isMaintenanceScheduled = true;
        emit MaintenanceScheduled();
    }

    function initializeETHDKG()
        public
        onlyFactory
        assertNotETHDKGRunning
        assertNotConsensusRunning
    {
        IETHDKG(_ethdkgAddress()).initializeETHDKG();
    }

    function completeETHDKG() public onlyETHDKG {
        _isMaintenanceScheduled = false;
        _isConsensusRunning = true;
    }

    // todo: check async in AliceNet
    function pauseConsensus() public onlySnapshots {
        _isConsensusRunning = false;
    }

    function pauseConsensusOnArbitraryHeight(uint256 aliceNetHeight_)
        public
        onlyFactory
    {
        require(
            block.number >
                ISnapshots(_snapshotsAddress())
                    .getCommittedHeightFromLatestSnapshot() +
                    MAX_INTERVAL_WITHOUT_SNAPSHOTS,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_MIN_BLOCK_INTERVAL_NOT_MET
                )
            )
        );
        _isConsensusRunning = false;
        IETHDKG(_ethdkgAddress()).setCustomAliceNetHeight(aliceNetHeight_);
    }

    function registerValidators(
        address[] memory validators_,
        uint256[] memory stakerTokenIDs_
    ) public onlyFactory assertNotETHDKGRunning assertNotConsensusRunning {
        require(
            validators_.length + _validators.length() <= _maxNumValidators,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_MAX_VALIDATORS_MET
                )
            )
        );
        require(
            validators_.length == stakerTokenIDs_.length,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_REGISTRATION_PARAMETER_LENGTH_MISMATCH
                )
            )
        );

        for (uint256 i = 0; i < validators_.length; i++) {
            require(
                msg.sender ==
                    IERC721(_publicStakingAddress()).ownerOf(
                        stakerTokenIDs_[i]
                    ),
                string(
                    abi.encodePacked(
                        ValidatorPoolErrorCodes
                            .VALIDATORPOOL_FACTORY_SHOULD_OWN_POSITION
                    )
                )
            );
            _registerValidator(validators_[i], stakerTokenIDs_[i]);
        }
    }

    function unregisterValidators(address[] memory validators_)
        public
        onlyFactory
        assertNotETHDKGRunning
        assertNotConsensusRunning
    {
        require(
            validators_.length <= _validators.length(),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_VALIDATORS_GREATER_THAN_AVAILABLE
                )
            )
        );
        for (uint256 i = 0; i < validators_.length; i++) {
            _unregisterValidator(validators_[i]);
        }
    }

    function unregisterAllValidators()
        public
        onlyFactory
        assertNotETHDKGRunning
        assertNotConsensusRunning
    {
        while (_validators.length() > 0) {
            address validator = _validators
                .at(_validators.length() - 1)
                ._address;
            _unregisterValidator(validator);
        }
    }

    function collectProfits()
        public
        onlyValidator
        returns (uint256 payoutEth, uint256 payoutToken)
    {
        require(
            _isConsensusRunning,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_PROFITS_ONLY_CLAIMABLE_DURING_CONSENSUS
                )
            )
        );

        uint256 balanceBeforeToken = IERC20Transferable(_aTokenAddress())
            .balanceOf(address(this));
        uint256 balanceBeforeEth = address(this).balance;

        uint256 validatorTokenID = _validators.get(msg.sender)._tokenID;
        payoutEth = IStakingNFT(_validatorStakingAddress()).collectEthTo(
            msg.sender,
            validatorTokenID
        );
        payoutToken = IStakingNFT(_validatorStakingAddress()).collectTokenTo(
            msg.sender,
            validatorTokenID
        );

        require(
            balanceBeforeToken ==
                IERC20Transferable(_aTokenAddress()).balanceOf(address(this)),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_TOKEN_BALANCE_CHANGED
                )
            )
        );
        require(
            balanceBeforeEth == address(this).balance,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETH_BALANCE_CHANGED
                )
            )
        );

        return (payoutEth, payoutToken);
    }

    function claimExitingNFTPosition() public returns (uint256) {
        ExitingValidatorData memory data = _exitingValidatorsData[msg.sender];
        require(
            data._freeAfter > 0,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_SENDER_NOT_IN_EXITING_QUEUE
                )
            )
        );
        require(
            ISnapshots(_snapshotsAddress()).getEpoch() > data._freeAfter,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_WAITING_PERIOD_NOT_MET
                )
            )
        );

        _removeExitingQueueData(msg.sender);

        IStakingNFT(_publicStakingAddress()).lockOwnPosition(
            data._tokenID,
            POSITION_LOCK_PERIOD
        );

        IERC721Transferable(_publicStakingAddress()).safeTransferFrom(
            address(this),
            msg.sender,
            data._tokenID
        );

        return data._tokenID;
    }

    function majorSlash(address dishonestValidator_, address disputer_)
        public
        onlyETHDKG
    {
        require(
            _isAccusable(dishonestValidator_),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_DISHONEST_VALIDATOR_NOT_ACCUSABLE
                )
            )
        );
        uint256 balanceBeforeToken = IERC20Transferable(_aTokenAddress())
            .balanceOf(address(this));
        uint256 balanceBeforeEth = address(this).balance;

        (uint256 minerShares, uint256 payoutEth, uint256 payoutToken) = _slash(
            dishonestValidator_
        );
        // deciding which state to clean based if the accusable person was a active validator or was
        // in the exiting line
        if (isValidator(dishonestValidator_)) {
            _removeValidatorData(dishonestValidator_);
        } else {
            _removeExitingQueueData(dishonestValidator_);
        }
        // redistribute the dishonest staking equally with the other validators

        IERC20Transferable(_aTokenAddress()).approve(
            _validatorStakingAddress(),
            minerShares
        );
        IStakingNFT(_validatorStakingAddress()).depositToken(
            _getMagic(),
            minerShares
        );
        // transfer to the disputer any profit that the dishonestValidator had when his
        // position was burned + the disputerReward
        _transferEthAndTokens(disputer_, payoutEth, payoutToken);

        require(
            balanceBeforeToken ==
                IERC20Transferable(_aTokenAddress()).balanceOf(address(this)),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_TOKEN_BALANCE_CHANGED
                )
            )
        );
        require(
            balanceBeforeEth == address(this).balance,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETH_BALANCE_CHANGED
                )
            )
        );

        emit ValidatorMajorSlashed(dishonestValidator_);
    }

    function minorSlash(address dishonestValidator_, address disputer_)
        public
        onlyETHDKG
    {
        require(
            _isAccusable(dishonestValidator_),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_DISHONEST_VALIDATOR_NOT_ACCUSABLE
                )
            )
        );
        uint256 balanceBeforeToken = IERC20Transferable(_aTokenAddress())
            .balanceOf(address(this));
        uint256 balanceBeforeEth = address(this).balance;

        (uint256 minerShares, uint256 payoutEth, uint256 payoutToken) = _slash(
            dishonestValidator_
        );
        uint256 stakeTokenID;
        // In case there's not enough shares to create a new PublicStaking position, state is just
        // cleaned and the rest of the funds is sent to the disputer
        if (minerShares > 0) {
            stakeTokenID = _mintPublicStakingPosition(minerShares);
            _moveToExitingQueue(dishonestValidator_, stakeTokenID);
        } else {
            if (isValidator(dishonestValidator_)) {
                _removeValidatorData(dishonestValidator_);
            } else {
                _removeExitingQueueData(dishonestValidator_);
            }
        }
        _transferEthAndTokens(disputer_, payoutEth, payoutToken);

        require(
            balanceBeforeToken ==
                IERC20Transferable(_aTokenAddress()).balanceOf(address(this)),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_TOKEN_BALANCE_CHANGED
                )
            )
        );
        require(
            balanceBeforeEth == address(this).balance,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETH_BALANCE_CHANGED
                )
            )
        );

        emit ValidatorMinorSlashed(dishonestValidator_, stakeTokenID);
    }

    /// skimExcessEth will allow the Admin role to refund any Eth sent to this contract in error by a
    /// user. This function should only be necessary if a user somehow manages to accidentally
    /// selfDestruct a contract with this contract as the recipient or use the PublicStaking burnTo with the
    /// address of this contract.
    function skimExcessEth(address to_)
        public
        onlyFactory
        returns (uint256 excess)
    {
        // This contract shouldn't held any eth balance.
        // todo: revisit this when we have the dutch auction
        excess = address(this).balance;
        _safeTransferEth(to_, excess);
        return excess;
    }

    /// skimExcessToken will allow the Admin role to refund any AToken sent to this contract in error
    /// by a user.
    function skimExcessToken(address to_)
        public
        onlyFactory
        returns (uint256 excess)
    {
        // This contract shouldn't held any token balance.
        IERC20Transferable aToken = IERC20Transferable(_aTokenAddress());
        excess = aToken.balanceOf(address(this));
        _safeTransferERC20(aToken, to_, excess);
        return excess;
    }

    function getStakeAmount() public view returns (uint256) {
        return _stakeAmount;
    }

    function getMaxNumValidators() public view returns (uint256) {
        return _maxNumValidators;
    }

    function getDisputerReward() public view returns (uint256) {
        return _disputerReward;
    }

    function getValidatorsCount() public view returns (uint256) {
        return _validators.length();
    }

    function getValidatorsAddresses() public view returns (address[] memory) {
        return _validators.addressValues();
    }

    function getValidator(uint256 index_) public view returns (address) {
        require(
            index_ < _validators.length(),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_INVALID_INDEX
                )
            )
        );
        return _validators.at(index_)._address;
    }

    function getValidatorData(uint256 index_)
        public
        view
        returns (ValidatorData memory)
    {
        require(
            index_ < _validators.length(),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_INVALID_INDEX
                )
            )
        );
        return _validators.at(index_);
    }

    function getLocation(address validator_)
        public
        view
        returns (string memory)
    {
        return _ipLocations[validator_];
    }

    function getLocations(address[] calldata validators_)
        public
        view
        returns (string[] memory)
    {
        string[] memory ret = new string[](validators_.length);
        for (uint256 i = 0; i < validators_.length; i++) {
            ret[i] = _ipLocations[validators_[i]];
        }
        return ret;
    }

    /// @notice Try to get the NFT tokenID for an account.
    /// @param account_ address of the account to try to retrieve the tokenID
    /// @return tuple (bool, address, uint256). Return true if the value was found, false if not.
    /// Returns the address of the NFT contract and the tokenID. In case the value was not found, tokenID
    /// and address are 0.
    function tryGetTokenID(address account_)
        public
        view
        returns (
            bool,
            address,
            uint256
        )
    {
        if (_isValidator(account_)) {
            return (
                true,
                _validatorStakingAddress(),
                _validators.get(account_)._tokenID
            );
        } else if (_isInExitingQueue(account_)) {
            return (
                true,
                _publicStakingAddress(),
                _exitingValidatorsData[account_]._tokenID
            );
        } else {
            return (false, address(0), 0);
        }
    }

    function isValidator(address account_) public view returns (bool) {
        return _isValidator(account_);
    }

    function isInExitingQueue(address account_) public view returns (bool) {
        return _isInExitingQueue(account_);
    }

    function isAccusable(address account_) public view returns (bool) {
        return _isAccusable(account_);
    }

    function isMaintenanceScheduled() public view returns (bool) {
        return _isMaintenanceScheduled;
    }

    function isConsensusRunning() public view returns (bool) {
        return _isConsensusRunning;
    }

    function _transferEthAndTokens(
        address to_,
        uint256 payoutEth_,
        uint256 payoutToken_
    ) internal {
        _safeTransferERC20(
            IERC20Transferable(_aTokenAddress()),
            to_,
            payoutToken_
        );
        _safeTransferEth(to_, payoutEth_);
    }

    function _registerValidator(address validator_, uint256 stakerTokenID_)
        internal
        returns (
            uint256 validatorTokenID,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        require(
            _validators.length() <= _maxNumValidators,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_MAX_VALIDATORS_MET
                )
            )
        );
        require(
            !_isAccusable(validator_),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_ADDRESS_ALREADY_VALIDATOR
                )
            )
        );

        uint256 balanceBeforeToken = IERC20Transferable(_aTokenAddress())
            .balanceOf(address(this));
        uint256 balanceBeforeEth = address(this).balance;
        (
            validatorTokenID,
            payoutEth,
            payoutToken
        ) = _swapPublicStakingForValidatorStaking(msg.sender, stakerTokenID_);

        _validators.add(ValidatorData(validator_, validatorTokenID));
        // transfer back any profit that was available for the PublicStaking position by the time that we
        // burned it
        _transferEthAndTokens(validator_, payoutEth, payoutToken);
        require(
            balanceBeforeToken ==
                IERC20Transferable(_aTokenAddress()).balanceOf(address(this)),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_TOKEN_BALANCE_CHANGED
                )
            )
        );
        require(
            balanceBeforeEth == address(this).balance,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETH_BALANCE_CHANGED
                )
            )
        );

        emit ValidatorJoined(validator_, validatorTokenID);
    }

    function _unregisterValidator(address validator_)
        internal
        returns (
            uint256 stakeTokenID,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        require(
            _isValidator(validator_),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ADDRESS_NOT_VALIDATOR
                )
            )
        );

        uint256 balanceBeforeToken = IERC20Transferable(_aTokenAddress())
            .balanceOf(address(this));
        uint256 balanceBeforeEth = address(this).balance;
        (
            stakeTokenID,
            payoutEth,
            payoutToken
        ) = _swapValidatorStakingForPublicStaking(validator_);

        _moveToExitingQueue(validator_, stakeTokenID);

        // transfer back any profit that was available for the PublicStaking position by the time that we
        // burned it
        _transferEthAndTokens(validator_, payoutEth, payoutToken);
        require(
            balanceBeforeToken ==
                IERC20Transferable(_aTokenAddress()).balanceOf(address(this)),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_TOKEN_BALANCE_CHANGED
                )
            )
        );
        require(
            balanceBeforeEth == address(this).balance,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ETH_BALANCE_CHANGED
                )
            )
        );

        emit ValidatorLeft(validator_, stakeTokenID);
    }

    function _swapPublicStakingForValidatorStaking(
        address to_,
        uint256 stakerTokenID_
    )
        internal
        returns (
            uint256 validatorTokenID,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        (uint256 stakeShares, , , , ) = IStakingNFT(_publicStakingAddress())
            .getPosition(stakerTokenID_);
        uint256 stakeAmount = _stakeAmount;
        require(
            stakeShares >= stakeAmount,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes
                        .VALIDATORPOOL_INSUFFICIENT_FUNDS_IN_STAKE_POSITION
                )
            )
        );
        IERC721Transferable(_publicStakingAddress()).safeTransferFrom(
            to_,
            address(this),
            stakerTokenID_
        );
        (payoutEth, payoutToken) = IStakingNFT(_publicStakingAddress()).burn(
            stakerTokenID_
        );

        // Subtracting the shares from PublicStaking profit. The shares will be used to mint the new
        // ValidatorPosition
        //payoutToken should always have the minerShares in it!
        require(
            payoutToken >= stakeShares,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_PAYOUT_TOO_LOW
                )
            )
        );
        payoutToken -= stakeAmount;

        validatorTokenID = _mintValidatorStakingPosition(stakeAmount);

        return (validatorTokenID, payoutEth, payoutToken);
    }

    function _swapValidatorStakingForPublicStaking(address validator_)
        internal
        returns (
            uint256,
            uint256,
            uint256
        )
    {
        (
            uint256 minerShares,
            uint256 payoutEth,
            uint256 payoutToken
        ) = _burnValidatorStakingPosition(validator_);
        //payoutToken should always have the minerShares in it!
        require(
            payoutToken >= minerShares,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_PAYOUT_TOO_LOW
                )
            )
        );
        payoutToken -= minerShares;

        uint256 stakeTokenID = _mintPublicStakingPosition(minerShares);

        return (stakeTokenID, payoutEth, payoutToken);
    }

    function _mintValidatorStakingPosition(uint256 minerShares_)
        internal
        returns (uint256 validatorTokenID)
    {
        // We should approve the validatorStaking to transferFrom the tokens of this contract
        IERC20Transferable(_aTokenAddress()).approve(
            _validatorStakingAddress(),
            minerShares_
        );
        validatorTokenID = IStakingNFT(_validatorStakingAddress()).mint(
            minerShares_
        );
    }

    function _mintPublicStakingPosition(uint256 minerShares_)
        internal
        returns (uint256 stakeTokenID)
    {
        // We should approve the PublicStaking to transferFrom the tokens of this contract
        IERC20Transferable(_aTokenAddress()).approve(
            _publicStakingAddress(),
            minerShares_
        );
        stakeTokenID = IStakingNFT(_publicStakingAddress()).mint(minerShares_);
    }

    function _burnValidatorStakingPosition(address validator_)
        internal
        returns (
            uint256 minerShares,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        uint256 validatorTokenID = _validators.get(validator_)._tokenID;
        (minerShares, payoutEth, payoutToken) = _burnNFTPosition(
            validatorTokenID,
            _validatorStakingAddress()
        );
    }

    function _burnExitingPublicStakingPosition(address validator_)
        internal
        returns (
            uint256 minerShares,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        uint256 stakerTokenID = _exitingValidatorsData[validator_]._tokenID;
        (minerShares, payoutEth, payoutToken) = _burnNFTPosition(
            stakerTokenID,
            _publicStakingAddress()
        );
    }

    function _burnNFTPosition(uint256 tokenID_, address stakeContractAddress_)
        internal
        returns (
            uint256 minerShares,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        IStakingNFT stakeContract = IStakingNFT(stakeContractAddress_);
        (minerShares, , , , ) = stakeContract.getPosition(tokenID_);
        (payoutEth, payoutToken) = stakeContract.burn(tokenID_);
    }

    function _slash(address dishonestValidator_)
        internal
        returns (
            uint256 minerShares,
            uint256 payoutEth,
            uint256 payoutToken
        )
    {
        require(
            _isAccusable(dishonestValidator_),
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_ADDRESS_NOT_ACCUSABLE
                )
            )
        );
        // If the user accused is a valid validator, we should burn is validatorStaking position,
        // otherwise we burn the user's PublicStaking in the exiting line
        if (_isValidator(dishonestValidator_)) {
            (
                minerShares,
                payoutEth,
                payoutToken
            ) = _burnValidatorStakingPosition(dishonestValidator_);
        } else {
            (
                minerShares,
                payoutEth,
                payoutToken
            ) = _burnExitingPublicStakingPosition(dishonestValidator_);
        }
        uint256 disputerReward = _disputerReward;
        if (minerShares >= disputerReward) {
            minerShares -= disputerReward;
        } else {
            // In case there's not enough shares to cover the _disputerReward, minerShares is set to
            // 0 and the rest of the payout Token is sent to disputer
            minerShares = 0;
        }
        //payoutToken should always have the minerShares in it!
        require(
            payoutToken >= minerShares,
            string(
                abi.encodePacked(
                    ValidatorPoolErrorCodes.VALIDATORPOOL_PAYOUT_TOO_LOW
                )
            )
        );
        payoutToken -= minerShares;
    }

    function _moveToExitingQueue(address validator_, uint256 stakeTokenID_)
        internal
    {
        if (_isValidator(validator_)) {
            _removeValidatorData(validator_);
        }
        _exitingValidatorsData[validator_] = ExitingValidatorData(
            uint128(stakeTokenID_),
            uint128(ISnapshots(_snapshotsAddress()).getEpoch() + CLAIM_PERIOD)
        );
    }

    function _removeValidatorData(address validator_) internal {
        _validators.remove(validator_);
        delete _ipLocations[validator_];
    }

    function _removeExitingQueueData(address validator_) internal {
        delete _exitingValidatorsData[validator_];
    }

    function _isValidator(address account_) internal view returns (bool) {
        return _validators.contains(account_);
    }

    function _isInExitingQueue(address account_) internal view returns (bool) {
        return _exitingValidatorsData[account_]._freeAfter > 0;
    }

    function _isAccusable(address account_) internal view returns (bool) {
        return _isValidator(account_) || _isInExitingQueue(account_);
    }
}
