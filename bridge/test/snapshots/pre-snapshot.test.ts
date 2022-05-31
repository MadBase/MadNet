import { Signer } from "ethers";
import { ethers } from "hardhat";
import { expect } from "../chai-setup";
import {
  Fixture,
  getFixture,
  getTokenIdFromTx,
  getValidatorEthAccount,
  mineBlocks,
} from "../setup";
import { validatorsSnapshotsG1 } from "../sharedConstants/4-validators-snapshots-100-Group1";

describe("Snapshots: Tests Snapshots methods", () => {
  let fixture: Fixture;
  let adminSigner: Signer;
  let randomSigner: Signer;
  const stakeAmount = 20000;
  const stakeAmountATokenWei = ethers.utils.parseUnits(
    stakeAmount.toString(),
    18
  );
  const lockTime = 1;
  let validators: any[];
  let stakingTokenIds: any[];

  beforeEach(async function () {
    validators = [];
    stakingTokenIds = [];
    fixture = await getFixture(true, false, undefined, true);
    const [admin, , , , , randomUser] = fixture.namedSigners;
    adminSigner = await getValidatorEthAccount(admin.address);
    randomSigner = await getValidatorEthAccount(randomUser.address);
    for (const validator of validatorsSnapshotsG1) {
      validators.push(validator.address);
    }

    await fixture.aToken.approve(
      fixture.validatorPool.address,
      stakeAmountATokenWei.mul(validators.length)
    );
    await fixture.aToken.approve(
      fixture.publicStaking.address,
      stakeAmountATokenWei.mul(validators.length)
    );

    for (const validator of validatorsSnapshotsG1) {
      const tx = await fixture.publicStaking
        .connect(adminSigner)
        .mintTo(validator.address, stakeAmountATokenWei, lockTime);
      const tokenId = getTokenIdFromTx(tx);
      stakingTokenIds.push(tokenId);
      await fixture.publicStaking
        .connect(await getValidatorEthAccount(validator))
        .setApprovalForAll(fixture.validatorPool.address, true);
    }

    await fixture.validatorPool
      .connect(adminSigner)
      .registerValidators(validators, stakingTokenIds);
  });

  it("Does not allow snapshot if sender is not validator", async function () {
    const junkData =
      "0x0000000000000000000000000000000000000000000000000000006d6168616d";
    await expect(
      fixture.snapshots.connect(randomSigner).snapshot(junkData, junkData)
    ).to.be.revertedWith("400");
  });

  it("Does not allow snapshot consensus is not running", async function () {
    const junkData =
      "0x0000000000000000000000000000000000000000000000000000006d6168616d";
    const validValidator = await getValidatorEthAccount(
      validatorsSnapshotsG1[0]
    );
    await mineBlocks(
      (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
    );
    await expect(
      fixture.snapshots.connect(validValidator).snapshot(junkData, junkData)
    ).to.be.revertedWith(`401`);
  });
});
