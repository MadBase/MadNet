import { contract, ethers } from "hardhat";
import { 
    Snapshots,
    AliceNetFactory,
    AToken,
    ATokenBurner,
    ATokenMinter,
    BToken,
    ETHDKG,
    Foundation,
    LegacyToken,
    LiquidityProviderStaking,
    PublicStaking,
    SnapshotsMock,
    StakingPositionDescriptor,
    ValidatorPool,
    ValidatorPoolMock,
    ValidatorStaking,
    SnapshotsV2, } from "../../typechain-types";
import {
    deployFactoryAndBaseTokens,
    deployLogicAndUpgradeWithFactory,
    deployUpgradeableWithFactory,
    Fixture,
    getFixture,
    getValidatorEthAccount,
    mineBlocks,
    posFixtureSetup,
    preFixtureSetup,
} from "../setup";
import {
    validatorsSnapshots,
    signedData
} from "../math/assets/4-validators-1000-snapshots";
import { completeETHDKGRound } from "../ethdkg/setup";
import { BigNumber, ContractTransaction } from "ethers";
import { SnapshotStructOutput } from "../../typechain-types/SnapshotsMock";
import { expect } from "chai";


contract("SnapshotRingBuffer", async () => {
    let fixture: Fixture;
    let snapshots: Snapshots;
    let numEpochs: number;
    let initialTestSnapshots: Array<SnapshotStructOutput> = [];
    let signedSnapshots
    describe("Snapshot upgrade integration",async () => {
      
      beforeEach(async () => {
        numEpochs = 0;
        fixture = await getFixture(true, false);
        await completeETHDKGRound(validatorsSnapshots,
          {
            ethdkg: fixture.ethdkg,
            validatorPool: fixture.validatorPool,
          }
        );
        signedSnapshots = signedData;
        await mineBlocks(
          (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
        );
        snapshots = fixture.snapshots.connect(await getValidatorEthAccount(validatorsSnapshots[0])) as Snapshots;          
        //take 6 snapshots
        for(let i = 0; i < 6; i++){
          await mineBlocks(
            (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
            );
            const contractTx = await snapshots.snapshot(signedSnapshots[i].GroupSignature, signedSnapshots[i].BClaims, {gasLimit:30000000})
            await contractTx.wait();
          numEpochs++;
        }
        //upgrade the snapshot contract
        fixture.snapshotsV2 = await deployLogicAndUpgradeWithFactory(fixture.factory, "SnapshotsV2", snapshots.address, undefined, [],[1, 1024]) as SnapshotsV2; 
      });
      
      it("verifies epoch value and snapshot migration onto ring buffer", async () => {
        let snapshotsV2 = 
          fixture.snapshotsV2 === undefined ? 
          await ethers.getContractAt("SnapshotsV2", fixture.snapshots.address) 
          : fixture.snapshotsV2;
        for(let i = numEpochs; i >= numEpochs - 2; i--){
          const snap = await snapshotsV2.getSnapshot(i); 
          expect(snap.blockClaims.height).to.equal(i*1024)
        }
        expect(await snapshots.getEpoch()).to.equal(numEpochs)
      })

      it("adds 6 new snapshots to the snapshot buffer", async () => {
        let snapshotsV2 = 
          fixture.snapshotsV2 === undefined ? 
          await ethers.getContractAt("SnapshotsV2", fixture.snapshots.address) 
          : fixture.snapshotsV2;
        const signedSnapshots = signedData
        const numSnaps = numEpochs + 6;
        snapshotsV2 = snapshotsV2.connect(await getValidatorEthAccount(validatorsSnapshots[0]));          
        //take 6 snapshots
        for(let i = numEpochs; i < numSnaps; i++){
          console.log(i)
          await mineBlocks(
            (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
          );
          const contractTx = await snapshotsV2.snapshot(signedSnapshots[i].GroupSignature, signedSnapshots[i].BClaims, {gasLimit:30000000})
          await contractTx.wait();
          numEpochs++;
        }
        console.log(await snapshotsV2.getEpoch())
        const lastSnapshot = await snapshotsV2.getLatestSnapshot()

      });

      it("attempts to get a snapshot that is not in the buffer", async () => {

      });
    })


});