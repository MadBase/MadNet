import { ethers } from "hardhat";
import { Snapshots } from "../../typechain-types";
import { expect } from "../chai-setup";
import {
  Fixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
  SNAPSHOT_BUFFER_LENGTH,
} from "../setup";
import {
  invalidSnapshot7168ChainID2,
  invalidSnapshot7668,
  invalidSnapshotIncorrectSig7168,
  signedData1,
  validatorsSnapshotsG1,
  validSnapshot7168,
} from "../sharedConstants/4-validators-snapshots-100-Group1";

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
  // this doesnt make sense
  xit("Reverts when validator not elected to do snapshot", async function () {
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .snapshot(validSnapshot7168.GroupSignature, validSnapshot7168.BClaims)
    ).to.be.revertedWith("1401");
  });

  it("Reverts when snapshot data contains invalid height", async function () {
    await expect(
      snapshots
        .connect(
          await getValidatorEthAccount(
            validatorsSnapshotsG1[invalidSnapshot7668.validatorIndex]
          )
        )
        .snapshot(
          invalidSnapshot7668.GroupSignature,
          invalidSnapshot7668.BClaims
        )
    ).to.be.revertedWith("406");
  });

  it("Reverts when snapshot data contains invalid chain id", async function () {
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .snapshot(
          invalidSnapshot7168ChainID2.GroupSignature,
          invalidSnapshot7168ChainID2.BClaims
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
              invalidSnapshotIncorrectSig7168.validatorIndex
            ]
          )
        )
        .snapshot(
          signedData1[SNAPSHOT_BUFFER_LENGTH].GroupSignature,
          invalidSnapshotIncorrectSig7168.BClaims
        )
    ).to.be.revertedWith("405");
  });

  it("Reverts when snapshot data contains incorrect public key", async function () {
    await expect(
      snapshots
        .connect(
          await getValidatorEthAccount(
            validatorsSnapshotsG1[
              invalidSnapshotIncorrectSig7168.validatorIndex
            ]
          )
        )
        .snapshot(
          invalidSnapshotIncorrectSig7168.GroupSignature,
          invalidSnapshotIncorrectSig7168.BClaims
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
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .snapshot(
          signedData1[SNAPSHOT_BUFFER_LENGTH].GroupSignature,
          signedData1[SNAPSHOT_BUFFER_LENGTH].BClaims
        )
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshotsG1[0].address),
        expectedSafeToProceedConsensus,
        signedData1[SNAPSHOT_BUFFER_LENGTH].GroupSignature
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
