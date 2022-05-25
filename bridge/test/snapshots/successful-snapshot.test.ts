import { ethers } from "hardhat";
import { Snapshots } from "../../typechain-types";
import { expect } from "../chai-setup";
import {
  signedData,
  validatorsSnapshots as validatorsSnapshots1,
} from "../math/assets/4-validators-1000-snapshots";
import {
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
  SNAPSHOT_BUFFER_LENGTH,
} from "../setup";
import {
  invalidSnapshot500,
  invalidSnapshotChainID2,
  invalidSnapshotIncorrectSig,
} from "./assets/4-validators-snapshots-1";
describe("Snapshots: With successful ETHDKG round completed", () => {
  let fixture: Fixture;
  let snapshots: Snapshots;
  beforeEach(async function () {
    fixture = await getFixture(true, false, undefined, true);
    await mineBlocks(
      (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
    );
    snapshots = fixture.snapshots;
  });

  it("Reverts when validator not elected to do snapshot", async function () {
    const junkData =
      "0x0000000000000000000000000000000000000000000000000000006d6168616d";
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshots1[0]))
        .snapshot(junkData, junkData)
    ).to.be.revertedWith("1401");
  });

  it("Reverts when snapshot data contains invalid height", async function () {
    await expect(
      snapshots
        .connect(
          await getValidatorEthAccount(
            validatorsSnapshots1[invalidSnapshot500.validatorIndex]
          )
        )
        .snapshot(invalidSnapshot500.GroupSignature, invalidSnapshot500.BClaims)
    ).to.be.revertedWith("406");
  });

  it("Reverts when snapshot data contains invalid chain id", async function () {
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshots1[0]))
        .snapshot(
          invalidSnapshotChainID2.GroupSignature,
          invalidSnapshotChainID2.BClaims
        )
    ).to.be.revertedWith("407");
  });

  // todo wrong public key failure happens first with this data
  it("Reverts when snapshot data contains incorrect signature", async function () {
    await expect(
      snapshots
        .connect(
          await getValidatorEthAccount(
            validatorsSnapshots1[invalidSnapshotIncorrectSig.validatorIndex]
          )
        )
        .snapshot(
          signedData[SNAPSHOT_BUFFER_LENGTH].GroupSignature,
          invalidSnapshotIncorrectSig.BClaims
        )
    ).to.be.revertedWith("405");
  });

  it("Reverts when snapshot data contains incorrect public key", async function () {
    await expect(
      snapshots
        .connect(
          await getValidatorEthAccount(
            validatorsSnapshots1[invalidSnapshotIncorrectSig.validatorIndex]
          )
        )
        .snapshot(
          invalidSnapshotIncorrectSig.GroupSignature,
          invalidSnapshotIncorrectSig.BClaims
        )
    ).to.be.revertedWith("404");
  });

  it("Successfully performs snapshot", async function () {
    const expectedChainId = 1;
    const expectedEpoch = SNAPSHOT_BUFFER_LENGTH + 1;
    const expectedHeight = expectedEpoch * 1024;
    const expectedSafeToProceedConsensus = true;
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshots1[0]))
        .snapshot(
          signedData[SNAPSHOT_BUFFER_LENGTH].GroupSignature,
          signedData[SNAPSHOT_BUFFER_LENGTH].BClaims
        )
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshots1[0].address),
        expectedSafeToProceedConsensus,
        signedData[SNAPSHOT_BUFFER_LENGTH].GroupSignature
      );
  });

  /*
  FYI this scenario is not possible to cover due to the fact that no validators can be registered but not participate in the ETHDKG round.

  it('Does not allow snapshot caller did not participate in the last ETHDKG round', async function () {
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshots1[0]))
        .snapshot(validSnapshot1024.GroupSignature, validSnapshot1024.BClaims)
    ).to.be.revertedWith(
      `Snapshots: Caller didn't participate in the last ethdkg round!`
    )
  }) */
});
