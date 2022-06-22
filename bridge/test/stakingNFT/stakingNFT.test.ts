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

  //describe("tripCB", async () => {});

  /*
  describe("skimExcessEth", async () => {
    it("attempts to skim excess eth as owner, should fail", async () => {});

    it("successfully skims excess eth", async () => {});
  });
  */

  //describe("burn", async () => {});

  describe("collect", async () => {
    it("collectPure Test 1; normal (all zero)", async () => {
      // Test Setup
      const shares = BigNumber.from("1000");
      const stateAccum = BigNumber.from("0");
      const stateSlush = BigNumber.from("0");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const positionShares = BigNumber.from("1000");
      const positionFreeAfter = BigNumber.from("1");
      const positionWithdrawFreeAfter = BigNumber.from("1");
      const positionAccumulatorEth = BigNumber.from("0");
      const positionAccumulatorToken = BigNumber.from("0");
      const position = {
        shares: positionShares,
        freeAfter: positionFreeAfter,
        withdrawFreeAfter: positionWithdrawFreeAfter,
        accumulatorEth: positionAccumulatorEth,
        accumulatorToken: positionAccumulatorToken,
      };
      let positionAccumulatorValue = BigNumber.from("0");

      // Expected Values
      const expPayout = BigNumber.from("0");
      const expUpdatedPositionAccumulatorValue = positionAccumulatorValue;
      const expStateAccumulator = BigNumber.from("0");
      const expStateSlush = BigNumber.from("0");
      const expPositionShares = positionShares;

      // Run collect
      const returned = await stakingNFT.collectPure(
        shares,
        state,
        position,
        positionAccumulatorValue
      );
      let retUpdatedState = returned[0];
      let retUpdatedPosition = returned[1];
      let retUpdatedPositionAccumulatorValue = returned[2];
      let retPayout = returned[3];

      // Verify returned values
      expect(retPayout).to.eq(expPayout);
      expect(retUpdatedPositionAccumulatorValue).to.eq(
        expUpdatedPositionAccumulatorValue
      );
      expect(retUpdatedState.accumulator).to.eq(expStateAccumulator);
      expect(retUpdatedState.slush).to.eq(expStateSlush);
      expect(retUpdatedPosition.shares).to.eq(expPositionShares);
      expect(retUpdatedPositionAccumulatorValue).to.eq(
        expUpdatedPositionAccumulatorValue
      );
      expect(retPayout).to.eq(expPayout);
    });

    it("collectPure Test 2; normal (state nonzero, no slush, no wraparound)", async () => {
      // Test Setup
      const shares = BigNumber.from("3");
      const stateAccum = BigNumber.from("10000000000000000000");
      const stateSlush = BigNumber.from("2");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const positionShares = BigNumber.from("1");
      const positionFreeAfter = BigNumber.from("1");
      const positionWithdrawFreeAfter = BigNumber.from("1");
      const positionAccumulatorEth = BigNumber.from("0");
      const positionAccumulatorToken = BigNumber.from("0");
      const position = {
        shares: positionShares,
        freeAfter: positionFreeAfter,
        withdrawFreeAfter: positionWithdrawFreeAfter,
        accumulatorEth: positionAccumulatorEth,
        accumulatorToken: positionAccumulatorToken,
      };
      let positionAccumulatorValue = BigNumber.from("0");

      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      // Compute Expected Values
      console.log("accumScaleFactor:", accumScaleFactor);
      let expStateSlush = stateSlush;
      const accumDelta = stateAccum.sub(positionAccumulatorValue);
      console.log("accumDelta:", accumDelta);
      let tmp = accumDelta.mul(positionShares);
      if (shares == positionShares) {
        tmp = tmp.add(stateSlush);
        expStateSlush = BigNumber.from("0");
      }
      console.log("tmp:", tmp);
      const expPayout = tmp.div(accumScaleFactor);
      console.log("expPayout:", expPayout);
      const payoutRem = tmp.sub(expPayout.mul(accumScaleFactor));
      console.log("payoutRem:", payoutRem);
      const expPositionAccumulatorValue = stateAccum;
      console.log("expPositionAccum:", expPositionAccumulatorValue);
      const expStateAccumulator = stateAccum;
      console.log("expStateAccum:", expStateAccumulator);
      expStateSlush = expStateSlush.add(payoutRem);
      console.log("expStateSlush:", expStateSlush);
      const expPositionShares = positionShares;
      console.log("expPositionShares:", expPositionShares);

      // Run collect
      const returned = await stakingNFT.collectPure(
        shares,
        state,
        position,
        positionAccumulatorValue
      );
      let retState = returned[0];
      let retPosition = returned[1];
      let retPositionAccumulatorValue = returned[2];
      let retPayout = returned[3];

      // Verify returned values
      expect(retPositionAccumulatorValue).to.eq(expPositionAccumulatorValue);
      expect(retState.accumulator).to.eq(expStateAccumulator);
      expect(retState.slush).to.eq(expStateSlush);
      expect(retPosition.shares).to.eq(expPositionShares);
      console.log("retPayout:", retPayout);
      console.log("expPayout:", expPayout);
      expect(retPayout).to.eq(expPayout);
    });
  });

  describe("slushSkim", async () => {
    it("slushSkimPure Test 1; normal", async () => {
      const shares = BigNumber.from("1000000");
      const accumulator = BigNumber.from("1234567890");
      const slush = BigNumber.from("5678901234");
      const deltaAccum = slush.div(shares);
      const slushExp = slush.sub(deltaAccum.mul(shares));
      const accumulatorExp = accumulator.add(deltaAccum);
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
      const shares = BigNumber.from("1111111111111111");
      const accumulator = BigNumber.from("123456789012345");
      const slush = BigNumber.from("567890123456789");
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
    it("depositPure Test 1; normal (deposit == 0)", async () => {
      const delta = BigNumber.from("0");
      const stateAccum = BigNumber.from("5678901234");
      const stateSlush = BigNumber.from("1234");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      const stateSlushExp = stateSlush.add(delta.mul(accumScaleFactor));
      const stateExp = { accumulator: stateAccum, slush: stateSlushExp };
      const returned = await stakingNFT.depositPure(delta, state);
      const stateRet = returned;
      expect(stateRet.accumulator).to.eq(stateExp.accumulator);
      expect(stateRet.slush).to.eq(stateExp.slush);
    });

    it("depositPure Test 2; normal (deposit > 0)", async () => {
      const delta = BigNumber.from("1234567890");
      const stateAccum = BigNumber.from("5678901234");
      const stateSlush = BigNumber.from("1234");
      const state = { accumulator: stateAccum, slush: stateSlush };
      const accumScaleFactor = await stakingNFT.getAccumulatorScaleFactor();
      const stateSlushExp = stateSlush.add(delta.mul(accumScaleFactor));
      const stateExp = { accumulator: stateAccum, slush: stateSlushExp };
      const returned = await stakingNFT.depositPure(delta, state);
      const stateRet = returned;
      expect(stateRet.accumulator).to.eq(stateExp.accumulator);
      expect(stateRet.slush).to.eq(stateExp.slush);
    });

    it("depositPure Test 3; fail: (failed require [slush overflow])", async () => {
      // Set initial values
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
      const tx = stakingNFT.depositPure(delta, state);
      // Error Code 608 for Slush Too Large; must have slush < 2**167
      await expect(tx).to.be.revertedWith("608");
    });
  });
});
