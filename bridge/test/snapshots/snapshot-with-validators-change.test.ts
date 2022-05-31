import { ethers } from "hardhat";
import { Snapshots } from "../../typechain-types";
import { expect } from "../chai-setup";
import { completeETHDKGRound } from "../ethdkg/setup";
import {
  factoryCallAnyFixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
  SNAPSHOT_BUFFER_LENGTH,
} from "../setup";
import {
  validatorsSnapshotsG1,
  validSnapshot7168,
  validSnapshot8192,
} from "../sharedConstants/4-validators-snapshots-100-Group1";
import {
  validatorsSnapshotsG2,
  validSnapshot8192 as validSnapShot8192G2,
} from "../sharedConstants/4-validators-snapshots-100-Group2";
import {
  createValidatorsWFixture,
  stakeValidatorsWFixture,
} from "../validatorPool/setup";

describe("Snapshots: With successful ETHDKG round completed and validatorPool", () => {
  it("Successfully performs snapshot then change the validators and perform another snapshot", async function () {
    let expectedChainId = 1;
    let expectedEpoch = SNAPSHOT_BUFFER_LENGTH + 1;
    let expectedHeight = validSnapshot7168.height as number;
    let expectedSafeToProceedConsensus = false;
    const fixture = await getFixture(undefined, undefined, undefined, true);
    const snapshots = fixture.snapshots as Snapshots;
    let validators: Array<string> = [];
    for (let validator of validatorsSnapshotsG1) {
      validators.push(validator.address);
    }
    await factoryCallAnyFixture(
      fixture,
      "validatorPool",
      "scheduleMaintenance"
    );
    await mineBlocks(
      (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
    );
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .snapshot(validSnapshot7168.GroupSignature, validSnapshot7168.BClaims)
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshotsG1[0].address),
        expectedSafeToProceedConsensus,
        validSnapshot7168.GroupSignature
      );
    await factoryCallAnyFixture(
      fixture,
      "validatorPool",
      "unregisterValidators",
      [validators]
    );

    // registering the new validators
    const newValidators = await createValidatorsWFixture(
      fixture,
      validatorsSnapshotsG2
    );
    const newStakingTokenIds = await stakeValidatorsWFixture(
      fixture,
      newValidators
    );
    await factoryCallAnyFixture(
      fixture,
      "validatorPool",
      "registerValidators",
      [newValidators, newStakingTokenIds]
    );
    await factoryCallAnyFixture(fixture, "validatorPool", "initializeETHDKG");
    await completeETHDKGRound(
      validatorsSnapshotsG2,
      {
        ethdkg: fixture.ethdkg,
        validatorPool: fixture.validatorPool,
      },
      expectedEpoch,
      expectedHeight,
      (await snapshots.getCommittedHeightFromLatestSnapshot()).toNumber()
    );
    await mineBlocks(
      (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
    );
    expectedChainId = 1;
    expectedEpoch = (validSnapshot8192.height as number) / 1024;
    expectedHeight = validSnapshot8192.height as number;
    expectedSafeToProceedConsensus = true;
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG2[0]))
        .snapshot(
          validSnapShot8192G2.GroupSignature,
          validSnapShot8192G2.BClaims
        )
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshotsG2[0].address),
        expectedSafeToProceedConsensus,
        validSnapShot8192G2.GroupSignature
      );
  });
});
