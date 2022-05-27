import { BigNumber, ContractTransaction, Signer } from "ethers";
import { ethers } from "hardhat";
import {
  AliceNetFactory,
  AToken,
  PublicStaking,
  ValidatorPool,
  ValidatorPoolMock,
  ValidatorStaking,
} from "../../typechain-types";
import { ValidatorRawData } from "../ethdkg/setup";
import {
  factoryCallAny,
  Fixture,
  getTokenIdFromTx,
  getValidatorEthAccount,
} from "../setup";

interface Contract {
  PublicStaking: bigint;
  ValNFT: bigint;
  ATK: bigint;
  ETH: bigint;
  Addr: string;
}

interface Admin {
  PublicStaking: bigint;
  ValNFT: bigint;
  ATK: bigint;
  Addr: string;
}

interface Validator {
  NFT: bigint;
  ATK: bigint;
  Addr: string;
  Reg: boolean;
  ExQ: boolean;
  Acc: boolean;
  Idx: number;
}
interface State {
  Admin: Admin;
  PublicStaking: Contract;
  ValidatorStaking: Contract;
  ValidatorPool: Contract;
  Factory: Contract;
  validators: Array<Validator>;
}

export const commitSnapshots = async (
  fixture: Fixture,
  numSnapshots: number
) => {
  for (let i = 0; i < numSnapshots; i++) {
    await fixture.snapshots.snapshot("0x00", "0x00");
  }
};

export const getCurrentStateWFixture = async (
  fixture: Fixture,
  _validators: string[]
): Promise<State> => {
  return getCurrentState(
    fixture.factory,
    fixture.publicStaking,
    fixture.validatorStaking,
    fixture.aToken,
    fixture.validatorPool,
    _validators
  );
};

export const getCurrentState = async (
  factory: AliceNetFactory,
  publicStaking: PublicStaking,
  validatorStaking: ValidatorStaking,
  aToken: AToken,
  validatorPool: ValidatorPool | ValidatorPoolMock,
  _validators: string[]
): Promise<State> => {
  // System state
  const state: State = {
    Admin: {
      PublicStaking: BigInt(0),
      ValNFT: BigInt(0),
      ATK: BigInt(0),
      Addr: "0x0",
    },
    PublicStaking: {
      PublicStaking: BigInt(0),
      ValNFT: BigInt(0),
      ATK: BigInt(0),
      ETH: BigInt(0),
      Addr: "0x0",
    },
    ValidatorStaking: {
      PublicStaking: BigInt(0),
      ValNFT: BigInt(0),
      ATK: BigInt(0),
      ETH: BigInt(0),
      Addr: "0x0",
    },
    ValidatorPool: {
      PublicStaking: BigInt(0),
      ValNFT: BigInt(0),
      ATK: BigInt(0),
      ETH: BigInt(0),
      Addr: "0x0",
    },
    Factory: {
      PublicStaking: BigInt(0),
      ValNFT: BigInt(0),
      ATK: BigInt(0),
      ETH: BigInt(0),
      Addr: "0x0",
    },
    validators: [],
  };
  const [adminSigner] = await ethers.getSigners();
  // Get state for admin
  state.Admin.PublicStaking = (
    await publicStaking.balanceOf(adminSigner.address)
  ).toBigInt();
  state.Admin.ValNFT = (
    await validatorStaking.balanceOf(adminSigner.address)
  ).toBigInt();
  state.Admin.ATK = (await aToken.balanceOf(adminSigner.address)).toBigInt();
  state.Admin.Addr = adminSigner.address;

  // Get state for validators
  for (let i = 0; i < _validators.length; i++) {
    const validator: Validator = {
      Idx: i,
      NFT: (await publicStaking.balanceOf(_validators[i])).toBigInt(),
      ATK: (await aToken.balanceOf(_validators[i])).toBigInt(),
      Addr: _validators[i],
      Reg: await validatorPool.isValidator(_validators[i]),
      ExQ: await validatorPool.isInExitingQueue(_validators[i]),
      Acc: await validatorPool.isAccusable(_validators[i]),
    };
    state.validators.push(validator);
  }
  // Contract data
  const contractData = [
    {
      contractState: state.PublicStaking,
      contractAddress: publicStaking.address,
    },
    {
      contractState: state.ValidatorStaking,
      contractAddress: validatorStaking.address,
    },
    {
      contractState: state.ValidatorPool,
      contractAddress: validatorPool.address,
    },
    {
      contractState: state.Factory,
      contractAddress: factory.address,
    },
  ];
  // Get state for contracts
  for (let i = 0; i < contractData.length; i++) {
    contractData[i].contractState.PublicStaking = (
      await publicStaking.balanceOf(contractData[i].contractAddress)
    ).toBigInt();
    contractData[i].contractState.ValNFT = (
      await validatorStaking.balanceOf(contractData[i].contractAddress)
    ).toBigInt();
    contractData[i].contractState.ATK = (
      await aToken.balanceOf(contractData[i].contractAddress)
    ).toBigInt();
    contractData[i].contractState.ETH = (
      await ethers.provider.getBalance(contractData[i].contractAddress)
    ).toBigInt();
    contractData[i].contractState.Addr = contractData[i].contractAddress;
  }
  return state;
};

export const showState = async (title: string, state: State): Promise<void> => {
  if (process.env.npm_config_detailed === "true") {
    // execute "npm --detailed=true  run test" to see this output
    console.log(title);
    console.log(state);
  }
};

export const createValidatorsWFixture = async (
  fixture: Fixture,
  _validatorsSnapshots: ValidatorRawData[]
): Promise<string[]> => {
  return await createValidators(
    fixture.factory,
    fixture.aToken,
    fixture.validatorPool,
    fixture.publicStaking,
    fixture.validatorStaking,
    _validatorsSnapshots
  );
};

export const createValidators = async (
  factory: AliceNetFactory,
  aToken: AToken,
  validatorPool: ValidatorPool | ValidatorPoolMock,
  publicStaking: PublicStaking,
  validatorStaking: ValidatorStaking,
  _validatorsSnapshots: ValidatorRawData[]
): Promise<string[]> => {
  const validators: string[] = [];
  const stakeAmountATokenWei = await validatorPool.getStakeAmount();
  const [adminSigner] = await ethers.getSigners();
  // Approve ValidatorPool to withdraw ATK tokens of validators
  await aToken.approve(
    validatorPool.address,
    stakeAmountATokenWei.mul(_validatorsSnapshots.length)
  );
  for (let i = 0; i < _validatorsSnapshots.length; i++) {
    const validator = _validatorsSnapshots[i];
    await getValidatorEthAccount(validator);
    validators.push(validator.address);
    // Send ATK tokens to each validator
    await aToken.transfer(validator.address, stakeAmountATokenWei);
  }
  await aToken
    .connect(adminSigner)
    .approve(
      publicStaking.address,
      stakeAmountATokenWei.mul(_validatorsSnapshots.length)
    );
  await showState(
    "After creating:",
    await getCurrentState(
      factory,
      publicStaking,
      validatorStaking,
      aToken,
      validatorPool,
      validators
    )
  );
  return validators;
};

export const stakeValidatorsWFixture = async (
  fixture: Fixture,
  validators: string[]
): Promise<BigNumber[]> => {
  await showState(
    "After staking:",
    await getCurrentStateWFixture(fixture, validators)
  );
  return await stakeValidators(
    fixture.factory,
    fixture.aToken,
    fixture.validatorPool as ValidatorPool,
    fixture.publicStaking,
    fixture.validatorStaking,
    validators
  );
};

export const stakeValidators = async (
  factory: AliceNetFactory,
  aToken: AToken,
  validatorPool: ValidatorPool | ValidatorPoolMock,
  publicStaking: PublicStaking,
  validatorStaking: ValidatorStaking,

  validators: string[]
): Promise<BigNumber[]> => {
  const stakingTokenIds: BigNumber[] = [];
  const [adminSigner] = await ethers.getSigners();
  const stakeAmountATokenWei = await validatorPool.getStakeAmount();
  const lockTime = 1;
  for (let i = 0; i < validators.length; i++) {
    // Stake all ATK tokens
    const tx = await publicStaking
      .connect(adminSigner)
      .mintTo(factory.address, stakeAmountATokenWei, lockTime);
    // Get the proof of staking (NFT's tokenID)
    const tokenID = await getTokenIdFromTx(tx);
    stakingTokenIds.push(tokenID);
    factoryCallAny(factory, publicStaking, "approve", [
      validatorPool.address,
      tokenID,
    ]);
  }
  await showState(
    "After staking:",
    await getCurrentState(
      factory,
      publicStaking,
      validatorStaking,
      aToken,
      validatorPool,
      validators
    )
  );
  return stakingTokenIds;
};

export const claimPosition = async (
  fixture: Fixture,
  validator: ValidatorRawData
): Promise<BigNumber> => {
  const claimTx = (await fixture.validatorPool
    .connect(await getValidatorEthAccount(validator))
    .claimExitingNFTPosition()) as ContractTransaction;
  const receipt = await ethers.provider.getTransactionReceipt(claimTx.hash);
  return BigNumber.from(receipt.logs[0].topics[3]);
};

export const getPublicStakingFromMinorSlashEvent = async (
  tx: ContractTransaction
): Promise<bigint> => {
  const receipt = await ethers.provider.getTransactionReceipt(tx.hash);
  const intrface = new ethers.utils.Interface([
    "event ValidatorMinorSlashed(address indexed account, uint256 publicStaking)",
  ]);
  const data = receipt.logs[receipt.logs.length - 1].data;
  const topics = receipt.logs[receipt.logs.length - 1].topics;
  const event = intrface.decodeEventLog("ValidatorMinorSlashed", data, topics);
  return event.publicStaking;
};

/**
 * Mint a publicStaking and burn it to the ValidatorPool contract. Besides a contract self destructing
 * itself, this is a method to send eth accidentally to the validatorPool contract
 * @param fixture
 * @param etherAmount
 * @param aTokenAmount
 * @param adminSigner
 */
export const burnStakeTo = async (
  fixture: Fixture,
  etherAmount: BigNumber,
  aTokenAmount: BigNumber,
  adminSigner: Signer
) => {
  await fixture.aToken
    .connect(adminSigner)
    .approve(fixture.publicStaking.address, aTokenAmount);
  const tx = await fixture.publicStaking
    .connect(adminSigner)
    .mint(aTokenAmount);
  const tokenID = await getTokenIdFromTx(tx);
  await fixture.publicStaking.depositEth(42, {
    value: etherAmount,
  });
  await fixture.publicStaking
    .connect(adminSigner)
    .burnTo(fixture.validatorPool.address, tokenID);
};

/**
 * Mint a publicStaking
 * @param fixture
 * @param etherAmount
 * @param aTokenAmount
 * @param adminSigner
 */
export const mintPublicStaking = async (
  fixture: Fixture,
  aTokenAmount: BigNumber,
  adminSigner: Signer
) => {
  await fixture.aToken
    .connect(adminSigner)
    .approve(fixture.publicStaking.address, aTokenAmount);
  const tx = await fixture.publicStaking
    .connect(adminSigner)
    .mint(aTokenAmount);
  return await getTokenIdFromTx(tx);
};
