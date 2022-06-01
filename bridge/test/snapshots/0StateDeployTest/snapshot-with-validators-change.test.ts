import { ethers } from "hardhat";
import { Snapshots } from "../../../typechain-types";
import { expect } from "../../chai-setup";
import { completeETHDKGRound } from "../../ethdkg/setup";
import {
  factoryCallAnyFixture,
  getFixture,
  getValidatorEthAccount,
  mineBlocks,
} from "../../setup";
import {
  validatorsSnapshotsG1,
  validSnapshot1024,
  validSnapshot2048,
} from "../../sharedConstants/4-validators-snapshots-100-Group1";
import {
  validatorsSnapshotsG2,
  validSnapshot2048G2,
} from "../../sharedConstants/4-validators-snapshots-100-Group2";
import {
  createValidatorsWFixture,
  stakeValidatorsWFixture,
} from "../../validatorPool/setup";

describe("Snapshots 0state: With successful ETHDKG round completed and validatorPool", () => {
  it("Successfully performs snapshot then change the validators and perform another snapshot", async function () {
    let expectedChainId = 1;
    let expectedEpoch = 1;
    let expectedHeight = validSnapshot1024.height as number;
    let expectedSafeToProceedConsensus = false;
    const fixture = await getFixture(undefined, undefined, undefined, true);

    const snapshots = fixture.snapshots as Snapshots;
    const validators: Array<string> = [];
    for (const validator of validatorsSnapshotsG1) {
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
        .snapshot(validSnapshot1024.GroupSignature, validSnapshot1024.BClaims)
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshotsG1[0].address),
        expectedSafeToProceedConsensus,
        validSnapshot1024.GroupSignature
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
    expectedEpoch = (validSnapshot2048.height as number) / 1024;
    expectedHeight = validSnapshot2048.height as number;
    expectedSafeToProceedConsensus = true;
    await expect(
      snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG2[0]))
        .snapshot(
          validSnapshot2048G2.GroupSignature,
          validSnapshot2048G2.BClaims
        )
    )
      .to.emit(snapshots, `SnapshotTaken`)
      .withArgs(
        expectedChainId,
        expectedEpoch,
        expectedHeight,
        ethers.utils.getAddress(validatorsSnapshotsG2[0].address),
        expectedSafeToProceedConsensus,
        validSnapshot2048G2.GroupSignature
      );
  });
});
