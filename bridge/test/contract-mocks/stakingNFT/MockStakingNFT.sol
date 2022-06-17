// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/libraries/StakingNFT/StakingNFT.sol";
import "contracts/libraries/StakingNFT/StakingNFTStorage.sol";

contract MockStakingNFT is StakingNFT {
    function mintNFTMock(address to_, uint256 amount_) public returns (uint256) {
        return StakingNFT._mintNFT(to_, amount_);
    }

    function burnMock(
        address from_,
        address to_,
        uint256 tokenID_
    ) public returns (uint256 payoutEth, uint256 payoutToken) {
        return StakingNFT._burn(from_, to_, tokenID_);
    }

    function collectMock(
        uint256 shares_,
        Accumulator memory state_,
        Position memory p_,
        uint256 positionAccumulatorValue_
    )
        public
        returns (
            Accumulator memory,
            Position memory,
            uint256,
            uint256
        )
    {
        return StakingNFT._collect(shares_, state_, p_, positionAccumulatorValue_);
    }

    function depositMock(
        uint256 shares_,
        uint256 delta_,
        Accumulator memory state_
    ) public returns (Accumulator memory) {
        return StakingNFT._deposit(shares_, delta_, state_);
    }

    function depositPure(
        uint256 shares_,
        uint256 delta_,
        Accumulator memory state_
    ) public pure returns (Accumulator memory) {
        return StakingNFT._deposit(shares_, delta_, state_);
    }

    function slushSkimMock(
        uint256 shares_,
        uint256 accumulator_,
        uint256 slush_
    ) public returns (uint256, uint256) {
        return StakingNFT._slushSkim(shares_, accumulator_, slush_);
    }

    function slushSkimPure(
        uint256 shares_,
        uint256 accumulator_,
        uint256 slush_
    ) public pure returns (uint256, uint256) {
        return StakingNFT._slushSkim(shares_, accumulator_, slush_);
    }
}
