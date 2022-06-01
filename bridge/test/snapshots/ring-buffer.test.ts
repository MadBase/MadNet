import { expect } from "chai";
import { contract, ethers } from "hardhat";
import { Snapshots } from "../../typechain-types";
import { completeETHDKGRound } from "../ethdkg/setup";
import {
  deployLogicAndUpgradeWithFactory,
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
  SNAPSHOT_BUFFER_LENGTH,
} from "../setup";
import {
  invalidSnapshot1024ChainID2,
  invalidSnapshot500,
  invalidSnapshotIncorrectSig1024,
  signedData1,
  validatorsSnapshotsG1,
} from "../sharedConstants/4-validators-snapshots-100-Group1";

contract("SnapshotRingBuffer", async () => {
  let fixture: Fixture;
  let snapshots: Snapshots;
  describe("Snapshot without migration", async () => {
    beforeEach(async () => {
      fixture = await getFixture(true, false, undefined, false);
      fixture.snapshots = (await deployLogicAndUpgradeWithFactory(
        fixture.factory,
        "Snapshots",
        fixture.snapshots.address,
        undefined,
        [],
        [1, 1024]
      )) as Snapshots;
      snapshots = fixture.snapshots;
      // const validators = await createValidatorsWFixture(fixture, validatorsSnapshotsG1)
      // const stakingTokenIds = await stakeValidatorsWFixture(fixture, validators);

      // await factoryCallAnyFixture(fixture, "validatorPool", "registerValidators", [
      //   validators,
      //   stakingTokenIds,
      // ]);
      // await factoryCallAnyFixture(fixture, "validatorPool", "initializeETHDKG")
      await completeETHDKGRound(validatorsSnapshotsG1, {
        ethdkg: fixture.ethdkg,
        validatorPool: fixture.validatorPool,
      });

      await mineBlocks(
        (
          await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
        ).toBigInt()
      );
      // upgraded contract has no prior snapshots
    });

    it("Reverts when snapshot data contains invalid height", async function () {
      await expect(
        snapshots
          .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
          .snapshot(
            invalidSnapshot500.GroupSignature,
            invalidSnapshot500.BClaims
          )
      ).to.be.revertedWith("406");
    });

    it("Reverts when snapshot data contains invalid chain id", async function () {
      await expect(
        snapshots
          .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
          .snapshot(
            invalidSnapshot1024ChainID2.GroupSignature,
            invalidSnapshot1024ChainID2.BClaims
          )
      ).to.be.revertedWith("407");
    });

    // todo wrong public key failure happens first with this data
    it("Reverts when snapshot data contains incorrect signature", async function () {
      await expect(
        snapshots
          .connect(
            await getValidatorEthAccount(
              validatorsSnapshotsG1[
                invalidSnapshotIncorrectSig1024.validatorIndex
              ]
            )
          )
          .snapshot(
            signedData1[SNAPSHOT_BUFFER_LENGTH].GroupSignature,
            invalidSnapshotIncorrectSig1024.BClaims
          )
      ).to.be.revertedWith("405");
    });

    it("Reverts when snapshot data contains incorrect public key", async function () {
      await expect(
        snapshots
          .connect(
            await getValidatorEthAccount(
              validatorsSnapshotsG1[
                invalidSnapshotIncorrectSig1024.validatorIndex
              ]
            )
          )
          .snapshot(
            invalidSnapshotIncorrectSig1024.GroupSignature,
            invalidSnapshotIncorrectSig1024.BClaims
          )
      ).to.be.revertedWith("404");
    });

    it("Successfully performs snapshot", async function () {
      const expectedChainId = 1;
      const expectedEpoch = 1;
      const expectedHeight = expectedEpoch * 1024;
      const expectedSafeToProceedConsensus = true;
      await expect(
        snapshots
          .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
          .snapshot(signedData1[0].GroupSignature, signedData1[0].BClaims)
      )
        .to.emit(snapshots, `SnapshotTaken`)
        .withArgs(
          expectedChainId,
          expectedEpoch,
          expectedHeight,
          ethers.utils.getAddress(validatorsSnapshotsG1[0].address),
          expectedSafeToProceedConsensus,
          signedData1[0].GroupSignature
        );
    });
  });

  describe("Snapshot upgrade integration", async () => {
    beforeEach(async () => {
      fixture = await getFixture(true, false, undefined, true);
    });

    it("verifies epoch value and snapshot migration onto ring buffer", async () => {
      const epochs = (await fixture.snapshots.getEpoch()).toNumber();
      for (let i = epochs - 5; i <= epochs; i++) {
        const snap = await fixture.snapshots.getSnapshot(i);
        expect(snap.blockClaims.height).to.equal(i * 1024);
      }
      expect((await fixture.snapshots.getEpoch()).toNumber()).to.equal(6);
    });

    it("adds 6 new snapshots to the snapshot buffer", async () => {
      let epochs = (await fixture.snapshots.getEpoch()).toNumber();
      const signedSnapshots = signedData1;
      const numSnaps = epochs + 6;
      const snapshots = fixture.snapshots.connect(
        await getValidatorEthAccount(validatorsSnapshotsG1[0])
      );
      // take 6 snapshots
      for (let i = epochs + 1; i <= numSnaps; i++) {
        await mineBlocks(
          (await snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
        );
        const contractTx = await snapshots.snapshot(
          signedSnapshots[i - 1].GroupSignature,
          signedSnapshots[i - 1].BClaims,
          { gasLimit: 30000000 }
        );
        await contractTx.wait();
      }
      epochs = (await snapshots.getEpoch()).toNumber();
      const lastSnapshot = await snapshots.getLatestSnapshot();
      expect(lastSnapshot.blockClaims.height).to.equal(epochs * 1024);
    });
  });
});
