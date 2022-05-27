import { expect } from "chai";
import { contract } from "hardhat";
import {
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
} from "../setup";
import {
  signedData,
  validatorsSnapshotsG1,
} from "../sharedConstants/4-validators-snapshots-100-Group1";

contract("SnapshotRingBuffer", async () => {
  let fixture: Fixture;
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
      const signedSnapshots = signedData;
      const numSnaps = epochs + 6;
      let snapshots = fixture.snapshots.connect(
        await getValidatorEthAccount(validatorsSnapshotsG1[0])
      );
      //take 6 snapshots
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

    it("attempts to get a snapshot that is not in the buffer", async () => {});
  });
});
