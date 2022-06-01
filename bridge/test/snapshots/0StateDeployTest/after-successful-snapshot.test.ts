import { BigNumber } from "ethers";
import { contract } from "hardhat";
import { expect } from "../../chai-setup";
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
contract("Snapshots 0state", async () => {
  describe("Snapshots: With successful snapshot completed", () => {
    let fixture: Fixture;
    let snapshotNumber: BigNumber;
    const initialEpoch = 6;
    const nextEpochVal = 7;
    beforeEach(async function () {
      fixture = await getFixture(
        true,
        false,
        undefined,
        true,
        undefined,
        undefined,
        true
      );
      snapshotNumber = await fixture.snapshots.getEpoch();
    });

    it("Should succeed doing a valid snapshot for next epoch", async function () {
      const validValidator = await getValidatorEthAccount(
        validatorsSnapshotsG1[0]
      );
      expect(await fixture.snapshots.getEpoch()).to.be.equal(initialEpoch);
      await mineBlocks(
        (
          await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
        ).toBigInt()
      );
      await fixture.snapshots
        .connect(validValidator)
        .snapshot(
          signedData1[initialEpoch].GroupSignature,
          signedData1[initialEpoch].BClaims
        );
      expect((await fixture.snapshots.getEpoch()).toNumber()).to.be.equal(
        nextEpochVal
      );
    });

    it("Should not allow committing a snapshot for next epoch before time", async function () {
      const validValidator = await getValidatorEthAccount(
        validatorsSnapshotsG1[0]
      );
      expect(await fixture.snapshots.getEpoch()).to.be.equal(initialEpoch);
      await expect(
        fixture.snapshots
          .connect(validValidator)
          .snapshot(
            signedData1[initialEpoch].GroupSignature,
            signedData1[initialEpoch].BClaims
          )
      ).to.be.revertedWith("402");
      expect((await fixture.snapshots.getEpoch()).toNumber()).to.be.equal(
        initialEpoch
      );
    });

    it("Does not allow snapshot with data from previous snapshot", async function () {
      const validValidator = await getValidatorEthAccount(
        validatorsSnapshotsG1[0]
      );
      await mineBlocks(
        (
          await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
        ).toBigInt()
      );
      await expect(
        fixture.snapshots
          .connect(validValidator)
          .snapshot(
            signedData1[initialEpoch - 1].GroupSignature,
            signedData1[initialEpoch - 1].BClaims
          )
      ).to.be.revertedWith("406");
    });

    it("Does not allow snapshot if ETHDKG round is Running", async function () {
      await fixture.validatorPool.scheduleMaintenance();
      await mineBlocks(
        (
          await fixture.snapshots.getMinimumIntervalBetweenSnapshots()
        ).toBigInt()
      );
      await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .snapshot(
          signedData1[initialEpoch].GroupSignature,
          signedData1[initialEpoch].BClaims
        );
      await fixture.validatorPool.initializeETHDKG();
      const junkData =
        "0x0000000000000000000000000000000000000000000000000000006d6168616d";
      const validValidator = await getValidatorEthAccount(
        validatorsSnapshotsG1[0]
      );
      await expect(
        fixture.snapshots.connect(validValidator).snapshot(junkData, junkData)
      ).to.be.revertedWith(`401`);
    });

    it("getLatestSnapshot returns correct snapshot data", async function () {
      const expectedChainId = BigNumber.from(1);
      const expectedHeight = BigNumber.from(initialEpoch * 1024);
      const expectedTxCount = BigNumber.from(0);
      const expectedPrevBlock = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedTxRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedStateRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedHeaderRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );

      const snapshotData = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getLatestSnapshot();

      const blockClaims = snapshotData.blockClaims;
      await expect(blockClaims.chainId).to.be.equal(expectedChainId);
      await expect(blockClaims.height).to.be.equal(expectedHeight);
      await expect(blockClaims.txCount).to.be.equal(expectedTxCount);
      await expect(blockClaims.prevBlock).to.be.equal(expectedPrevBlock);
      await expect(blockClaims.txRoot).to.be.equal(expectedTxRoot);
      await expect(blockClaims.stateRoot).to.be.equal(expectedStateRoot);
      await expect(blockClaims.headerRoot).to.be.equal(expectedHeaderRoot);
    });

    it("getBlockClaimsFromSnapshot returns correct data", async function () {
      const epoch = initialEpoch - 1;
      const expectedChainId = BigNumber.from(1);
      const expectedHeight = BigNumber.from(epoch * 1024);
      const expectedTxCount = BigNumber.from(0);
      const expectedPrevBlock = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedTxRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedStateRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedHeaderRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );

      const blockClaims = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getBlockClaimsFromSnapshot(epoch);

      await expect(blockClaims.chainId).to.be.equal(expectedChainId);
      await expect(blockClaims.height).to.be.equal(expectedHeight);
      await expect(blockClaims.txCount).to.be.equal(expectedTxCount);
      await expect(blockClaims.prevBlock).to.be.equal(expectedPrevBlock);
      await expect(blockClaims.txRoot).to.be.equal(expectedTxRoot);
      await expect(blockClaims.stateRoot).to.be.equal(expectedStateRoot);
      await expect(blockClaims.headerRoot).to.be.equal(expectedHeaderRoot);
    });

    it("getBlockClaimsFromLatestSnapshot returns correct data", async function () {
      const expectedChainId = BigNumber.from(1);
      const expectedHeight = BigNumber.from(initialEpoch * 1024);
      const expectedTxCount = BigNumber.from(0);
      const expectedPrevBlock = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedTxRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedStateRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );
      const expectedHeaderRoot = BigNumber.from(
        "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
      );

      const blockClaims = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getBlockClaimsFromLatestSnapshot();

      await expect(blockClaims.chainId).to.be.equal(expectedChainId);
      await expect(blockClaims.height).to.be.equal(expectedHeight);
      await expect(blockClaims.txCount).to.be.equal(expectedTxCount);
      await expect(blockClaims.prevBlock).to.be.equal(expectedPrevBlock);
      await expect(blockClaims.txRoot).to.be.equal(expectedTxRoot);
      await expect(blockClaims.stateRoot).to.be.equal(expectedStateRoot);
      await expect(blockClaims.headerRoot).to.be.equal(expectedHeaderRoot);
    });

    it("getAliceNetHeightFromSnapshot returns correct data", async function () {
      const epoch = initialEpoch - 1;
      const expectedHeight = BigNumber.from(epoch * 1024);
      const height = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getAliceNetHeightFromSnapshot(epoch);

      expect(height).to.be.equal(expectedHeight);
    });

    it("getAliceNetHeightFromLatestSnapshot returns correct data", async function () {
      const expectedHeight = BigNumber.from(initialEpoch * 1024);
      const height = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getAliceNetHeightFromLatestSnapshot();

      await expect(height).to.be.equal(expectedHeight);
    });

    it("getChainIdFromSnapshot returns correct chain id", async function () {
      const expectedChainId = 1;
      const chainId = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getChainIdFromSnapshot(snapshotNumber);

      await expect(chainId).to.be.equal(expectedChainId);
    });

    it("getChainIdFromLatestSnapshot returns correct chain id", async function () {
      const expectedChainId = 1;
      const chainId = await fixture.snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshotsG1[0]))
        .getChainIdFromLatestSnapshot();
      await expect(chainId).to.be.equal(expectedChainId);
    });
  });
});
