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
    ValidatorStaking, } from "../../typechain-types";
import {
    deployFactoryAndBaseTokens,
    Fixture,
    getValidatorEthAccount,
    mineBlocks,
    preFixtureSetup,
} from "../setup";
import {
    validatorsSnapshots,
    validSnapshot1024,
    validSnapshot2048,
} from "./assets/4-validators-snapshots-1";
import { completeETHDKGRound } from "../ethdkg/setup";
import { BigNumber } from "ethers";


contract("SnapshotRingBuffer", async () => {
    let fixture: Fixture;
    let snapshots: Snapshots;
    let snapshotNumber: BigNumber;
    beforeEach(async () => {
        fixture = await getV1Fixture(true, false);

        await completeETHDKGRound(validatorsSnapshots,
            {
                ethdkg: fixture.ethdkg,
                validatorPool: fixture.validatorPool,
            });
        await mineBlocks(
        (await fixture.snapshots.getMinimumIntervalBetweenSnapshots()).toBigInt()
        );
        snapshots = fixture.snapshots as Snapshots;
        await snapshots
        .connect(await getValidatorEthAccount(validatorsSnapshots[0]))
        .snapshot(validSnapshot1024.GroupSignature, validSnapshot1024.BClaims);
        snapshotNumber = BigNumber.from(1);
    });
});

export const getV1Fixture = async (
    mockValidatorPool?: boolean,
    mockSnapshots?: boolean,
    mockETHDKG?: boolean
  ): Promise<Fixture> => {
    await preFixtureSetup();
    const namedSigners = await ethers.getSigners();
    const [admin] = namedSigners;
    // Deploy the factory and base token
    const { factory, aToken, bToken, legacyToken, publicStaking } =
      await deployFactoryAndBaseTokens(admin);
  
    // ValidatorStaking is not considered a base token since is only used by validators
    const validatorStaking = (await deployUpgradeableWithFactory(
      factory,
      "ValidatorStaking",
      "ValidatorStaking",
      []
    )) as ValidatorStaking;
  
    // LiquidityProviderStaking
    const liquidityProviderStaking = (await deployUpgradeableWithFactory(
      factory,
      "LiquidityProviderStaking",
      "LiquidityProviderStaking",
      []
    )) as LiquidityProviderStaking;
  
    // Foundation
    const foundation = (await deployUpgradeableWithFactory(
      factory,
      "Foundation",
      undefined
    )) as Foundation;
  
    let validatorPool;
    if (typeof mockValidatorPool !== "undefined" && mockValidatorPool) {
      // ValidatorPoolMock
      validatorPool = (await deployUpgradeableWithFactory(
        factory,
        "ValidatorPoolMock",
        "ValidatorPool"
      )) as ValidatorPoolMock;
    } else {
      // ValidatorPool
      validatorPool = (await deployUpgradeableWithFactory(
        factory,
        "ValidatorPool",
        "ValidatorPool",
        [
          ethers.utils.parseUnits("20000", 18),
          10,
          ethers.utils.parseUnits("3", 18),
        ]
      )) as ValidatorPool;
    }
  
    // ETHDKG Accusations
    await deployUpgradeableWithFactory(factory, "ETHDKGAccusations");
  
    // StakingPositionDescriptor
    const stakingPositionDescriptor = (await deployUpgradeableWithFactory(
      factory,
      "StakingPositionDescriptor"
    )) as StakingPositionDescriptor;
  
    // ETHDKG Phases
    await deployUpgradeableWithFactory(factory, "ETHDKGPhases");
  
    // ETHDKG
    let ethdkg;
    if (typeof mockETHDKG !== "undefined" && mockETHDKG) {
      // ValidatorPoolMock
      ethdkg = (await deployUpgradeableWithFactory(
        factory,
        "ETHDKGMock",
        "ETHDKG",
        [BigNumber.from(40), BigNumber.from(6)]
      )) as ETHDKG;
    } else {
      // ValidatorPool
      ethdkg = (await deployUpgradeableWithFactory(factory, "ETHDKG", "ETHDKG", [
        BigNumber.from(40),
        BigNumber.from(6),
      ])) as ETHDKG;
    }