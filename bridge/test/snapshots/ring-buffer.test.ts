import { expect } from "chai";
import { contract } from "hardhat";
import { SnapshotStructOutput } from "../../typechain-types/SnapshotsMock";
import {
  signedData,
  validatorsSnapshots,
} from "../math/assets/4-validators-1000-snapshots";
import {
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
} from "../setup";

contract("SnapshotRingBuffer", async () => {
  let fixture: Fixture;
  let numEpochs: number;
  let initialTestSnapshots: Array<SnapshotStructOutput> = [];
  let signedSnapshots;
  describe("Snapshot upgrade integration", async () => {
    beforeEach(async () => {
      fixture = await getFixture(true, false, undefined, true);
    });

    it.only("verifies epoch value and snapshot migration onto ring buffer", async () => {
      if (fixture.snapshotsV2 === undefined) {
        throw new Error("failed to upgrade snapshotsV2");
      }
      for (let i = numEpochs; i >= numEpochs - 2; i--) {
        const snap = await fixture.snapshotsV2.getSnapshot(i);
        expect(snap.blockClaims.height).to.equal(i * 1024);
      }
      expect(await fixture.snapshotsV2.getEpoch()).to.equal(6);
    });

    it("adds 6 new snapshots to the snapshot buffer", async () => {
      if (fixture.snapshotsV2 === undefined) {
        throw new Error("failed to upgrade snapshotsV2");
      }
      const signedSnapshots = signedData;
      const numSnaps = numEpochs + 6;
      fixture.snapshotsV2 = fixture.snapshotsV2.connect(
        await getValidatorEthAccount(validatorsSnapshots[0])
      );
      //take 6 snapshots
      for (let i = numEpochs; i < numSnaps; i++) {
        await mineBlocks(
          (
            await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
          ).toBigInt()
        );
        const contractTx = await fixture.snapshotsV2.snapshot(
          signedSnapshots[i].GroupSignature,
          signedSnapshots[i].BClaims,
          { gasLimit: 30000000 }
        );
        await contractTx.wait();
        numEpochs++;
      }
      const lastSnapshot = await fixture.snapshotsV2.getLatestSnapshot();
    });

    it("attempts to get a snapshot that is not in the buffer", async () => {});
  });
});
