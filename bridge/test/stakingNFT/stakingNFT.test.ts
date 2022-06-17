import { expect } from "chai";
import { BigNumber } from "ethers";
import { contract, ethers } from "hardhat";
import { MockStakingNFT } from "../../typechain-types";

contract("StakingNFT", async (accounts) => {
  let stakingNFT: MockStakingNFT;
  const owner = accounts[0];

  before(async () => {
    const stakingNFTBase = await ethers.getContractFactory("MockStakingNFT");
    stakingNFT = await stakingNFTBase.deploy();
  });

  describe("tripCB", async () => {});

  describe("skimExcessEth", async () => {
    it("attempts to skim excess eth as owner, should fail", async () => {});

    it("successfully skims excess eth", async () => {});
  });

  describe("slushSkim", async () => {
    it("slushSkimPure Test 1; normal", async () => {
      const shares = 1000000;
      const accumulator = 1234567890;
      const slush = 5678901234;
      const deltaAccum = Math.floor(slush / shares);
      const slushExp = slush - deltaAccum * shares;
      const accumulatorExp = accumulator + deltaAccum;
      let accumulatorRet;
      let slushRet;
      const returned = await stakingNFT.slushSkimPure(
        shares,
        accumulator,
        slush
      );
      accumulatorRet = returned[0];
      slushRet = returned[1];
      expect(slushRet.toNumber()).to.eq(slushExp);
      expect(accumulatorRet.toNumber()).to.eq(accumulatorExp);
    });

    it("slushSkimPure Test 2; wrap around", async () => {
      const shares = BigNumber.from("1000");
      // 2**168 - 1
      const accumulator = BigNumber.from(
        "374144419156711147060143317175368453031918731001855"
      );
      const slush = BigNumber.from("1001");
      const deltaAccum = slush.div(shares);
      const slushExp = BigNumber.from("1");
      const accumulatorExp = BigNumber.from("0");
      let accumulatorRet;
      let slushRet;
      const returned = await stakingNFT.slushSkimPure(
        shares,
        accumulator,
        slush
      );
      accumulatorRet = returned[0];
      slushRet = returned[1];
      // console.log(returned);
      expect(slushRet).to.eq(slushExp);
      expect(accumulatorRet).to.eq(accumulatorExp);
    });
    // Need to know how to work with BigInt

    // Curly braces on the fly
    it("slushSkim gas Test", async () => {
      // Base Tx gas cost: 21K gas
      // Thus, subtract 21K off of all stated gas costs.
      const shares = 1111111111111111;
      const accumulator = 123456789012345;
      const slush = 567890123456789;
      const txResponse = await stakingNFT.slushSkimMock(
        shares,
        accumulator,
        slush
      );
      const receipt = await txResponse.wait();
      expect(receipt.status).to.eq(1);
    });
  });

  describe("deposit", async () => {
    it("depositPure Test 1; normal (deposit == 0; shares == 0)", async () => {
      const shares = BigNumber.from("0");
      const delta = BigNumber.from("0");
      const stateAccum = BigNumber.from("5678901234");
      const stateSlush = BigNumber.from("1234");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      const stateSlushExp = stateSlush.add(delta.mul(accumScaleFactor));
      const stateExp = { accumulator: stateAccum, slush: stateSlushExp };
      const returned = await stakingNFT.depositPure(shares, delta, state);
      const stateRet = returned;
      expect(stateRet.accumulator).to.eq(stateExp.accumulator);
      expect(stateRet.slush).to.eq(stateExp.slush);
    });

    it("depositPure Test 2; normal (deposit > 0; shares == 0)", async () => {
      const shares = BigNumber.from("0");
      const delta = BigNumber.from("1234567890");
      const stateAccum = BigNumber.from("5678901234");
      const stateSlush = BigNumber.from("1234");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      const stateSlushExp = stateSlush.add(delta.mul(accumScaleFactor));
      const stateExp = { accumulator: stateAccum, slush: stateSlushExp };
      const returned = await stakingNFT.depositPure(shares, delta, state);
      const stateRet = returned;
      expect(stateRet.accumulator).to.eq(stateExp.accumulator);
      expect(stateRet.slush).to.eq(stateExp.slush);
    });

    it("depositPure Test 3; normal (deposit > 0; shares > 0)", async () => {
      // Set initial values
      const shares = BigNumber.from("25519");
      const delta = BigNumber.from("1234567890");
      const stateAccumInitial = BigNumber.from("5678901234");
      const stateSlushInitial = BigNumber.from("1234");
      const state = {
        accumulator: stateAccumInitial,
        slush: stateSlushInitial,
      };
      // Compute return values
      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      const stateSlushUpdate = stateSlushInitial.add(
        delta.mul(accumScaleFactor)
      );
      let slushSkimRet = await stakingNFT.slushSkimPure(
        shares,
        stateAccumInitial,
        stateSlushUpdate
      );
      let stateAccumExp = slushSkimRet[0];
      let stateSlushExp = slushSkimRet[1];
      const stateExp = { accumulator: stateAccumExp, slush: stateSlushExp };
      // Call deposit
      const returned = await stakingNFT.depositPure(shares, delta, state);
      const stateRet = returned;
      expect(stateRet.accumulator).to.eq(stateExp.accumulator);
      expect(stateRet.slush).to.eq(stateExp.slush);
    });

    it("depositPure Test 4; fail: (failed require [slush overflow])", async () => {
      // Set initial values
      const shares = BigNumber.from("0");
      const delta = BigNumber.from("0");
      const stateAccumInitial = BigNumber.from("0");
      // 2**167
      const stateSlushInitial = BigNumber.from(
        "187072209578355573530071658587684226515959365500928"
      );
      const state = {
        accumulator: stateAccumInitial,
        slush: stateSlushInitial,
      };
      // Call deposit
      const tx = stakingNFT.depositPure(shares, delta, state);
      // Error Code 608 for Slush Too Large
      await expect(tx).to.be.revertedWith("608");
    });
  });

  describe("collect", async () => {});
  describe("burn", async () => {});
});
