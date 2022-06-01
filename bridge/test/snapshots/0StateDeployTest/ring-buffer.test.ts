import { expect } from "chai";
import { contract } from "hardhat";
import { Snapshots } from "../../../typechain-types";
import {
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
} from "../../setup";
import {
  signedData1,
  validatorsSnapshotsG1,
} from "../../sharedConstants/4-validators-snapshots-100-Group1";

contract("SnapshotRingBuffer 0state", async () => {
  const epochLength = 1024;
  let fixture: Fixture;
  let snapshots: Snapshots;
  describe("Snapshot upgrade integration", async () => {
    beforeEach(async () => {
      //deploys the new snapshot contract with buffer and zero state
      fixture = await getFixture(
        true,
        false,
        undefined,
        true,
        undefined,
        undefined,
        true
      );
    });

    it("checks if 6 snapshots has been done and that epoch number matches snapshots", async () => {
      const epochs = (await fixture.snapshots.getEpoch()).toNumber();
      for (let i = epochs - 5; i <= epochs; i++) {
        const snap = await fixture.snapshots.getSnapshot(i);
        expect(snap.blockClaims.height).to.equal(i * epochLength);
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
      expect(lastSnapshot.blockClaims.height).to.equal(epochs * epochLength);
    });
  });
});
