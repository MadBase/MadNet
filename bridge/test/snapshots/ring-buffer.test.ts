import { expect } from "chai";
import { contract } from "hardhat";
import { Fixture, getFixture } from "../setup";

contract("SnapshotRingBuffer", async () => {
  let fixture: Fixture;
  describe("Snapshot upgrade integration", async () => {
    beforeEach(async () => {
      fixture = await getFixture(true, false, undefined, true);
    });

    it.only("verifies epoch value and snapshot migration onto ring buffer", async () => {
      const epochs = (await fixture.snapshots.getEpoch()).toNumber();
      console.log(await fixture.snapshots.getSnapshot(1));
      console.log(epochs);
      for (let i = epochs - 5; i <= epochs; i++) {
        console.log(i);
        const snap = await fixture.snapshots.getSnapshot(i);
        expect(snap.blockClaims.height).to.equal(i * 1024);
      }
      expect((await fixture.snapshots.getEpoch()).toNumber()).to.equal(6);
    });

    // it("adds 6 new snapshots to the snapshot buffer", async () => {
    //   if (fixture.snapshotsV2 === undefined) {
    //     throw new Error("failed to upgrade snapshotsV2");
    //   }
    //   const signedSnapshots = signedData;
    //   const numSnaps = numEpochs + 6;
    //   fixture.snapshotsV2 = fixture.snapshotsV2.connect(
    //     await getValidatorEthAccount(validatorsSnapshots[0])
    //   );
    //   //take 6 snapshots
    //   for (let i = numEpochs; i < numSnaps; i++) {
    //     await mineBlocks(
    //       (
    //         await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
    //       ).toBigInt()
    //     );
    //     const contractTx = await fixture.snapshotsV2.snapshot(
    //       signedSnapshots[i].GroupSignature,
    //       signedSnapshots[i].BClaims,
    //       { gasLimit: 30000000 }
    //     );
    //     await contractTx.wait();
    //     numEpochs++;
    //   }
    //   const lastSnapshot = await fixture.snapshotsV2.getLatestSnapshot();
    // });

    it("attempts to get a snapshot that is not in the buffer", async () => {});
  });
});
