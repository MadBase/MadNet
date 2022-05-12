import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber } from "ethers/lib/ethers";
import { defaultAbiCoder } from "ethers/lib/utils";
import { ethers } from "hardhat";
import {
  BridgePoolErrorCodes,
  ImmutableAuthErrorCodes,
} from "../../typechain-types";
import {
  callFunctionAndGetReturnValues,
  factoryCallAnyFixture,
  Fixture,
} from "../setup";

let admin: SignerWithAddress;
let user: SignerWithAddress;
let user2: SignerWithAddress;

export interface state {
  Balances: {
    aToken: {
      address: string;
      admin: bigint;
      user: bigint;
      bridgePool: bigint;
      totalSupply: bigint;
    };
    bToken: {
      address: string;
      admin: bigint;
      user: bigint;
      bridgePool: bigint;
      totalSupply: bigint;
    };
    eth: {
      address: string;
      // We leave user balance as number to round values and avoid consumed gas comparison
      admin: number;
      user: number;
      bridgePool: bigint;
      aToken: bigint;
    };
  };
}

export async function getState(fixture: Fixture) {
  [admin, user, user2] = await ethers.getSigners();
  let state: state = {
    Balances: {
      aToken: {
        address: fixture.aToken.address.slice(-4),
        admin: (await fixture.aToken.balanceOf(admin.address)).toBigInt(),
        user: (await fixture.aToken.balanceOf(user.address)).toBigInt(),
        bridgePool: (
          await fixture.aToken.balanceOf(fixture.bridgePool.address)
        ).toBigInt(),
        totalSupply: (await fixture.aToken.totalSupply()).toBigInt(),
      },
      bToken: {
        address: fixture.bToken.address.slice(-4),
        admin: (await fixture.bToken.balanceOf(admin.address)).toBigInt(),
        user: (await fixture.bToken.balanceOf(user.address)).toBigInt(),
        bridgePool: (
          await fixture.bToken.balanceOf(fixture.bridgePool.address)
        ).toBigInt(),
        totalSupply: (await fixture.bToken.totalSupply()).toBigInt(),
      },
      eth: {
        address: "0000",
        admin: format(await ethers.provider.getBalance(admin.address)),
        user: format(await ethers.provider.getBalance(user.address)),
        bridgePool: (
          await ethers.provider.getBalance(fixture.bridgePool.address)
        ).toBigInt(),
        aToken: (
          await ethers.provider.getBalance(fixture.aToken.address)
        ).toBigInt(),
      },
    },
  };
  return state;
}

export function showState(title: string, state: state) {
  if (process.env.npm_config_detailed == "true") {
    // execute "npm --detailed=true test" to see this output
    console.log(title, state);
  }
}

export function format(number: BigNumber) {
  return parseFloat((+ethers.utils.formatEther(number)).toFixed(0));
}

export function formatBigInt(number: BigNumber) {
  return BigInt(parseFloat((+ethers.utils.formatEther(number)).toFixed(0)));
}

export var testData = {
  immutableAuthErrorCodesContract: {} as ImmutableAuthErrorCodes,
  bridgePoolErrorCodesContract: {} as BridgePoolErrorCodes,
  admin: {} as SignerWithAddress,
  user: {} as SignerWithAddress,
  user2: {} as SignerWithAddress,
  ethIn: BigNumber.from(0),
  ethsReceived: BigNumber.from(0),
  bTokenFeeInETH: 10,
  totalErc20Amount: BigNumber.from(20000).toBigInt(),
  erc20Amount: BigNumber.from(100).toBigInt(),
  bTokenAmount: BigNumber.from(100).toBigInt(),
  // The following merkle proof and stateRoot values can be obtained from accusation_builder_test.go execution
  merkleProof:
    "0x010005cda80a6c60e1215c1882b25b4744bd9d95c1218a2fd17827ab809c68196fd9bf0000000000000000000000000000000000000000000000000000000000000000af469f3b9864a5132323df8bdd9cbd59ea728cd7525b65252133a5a02f1566ee00010003a8793650a7050ac58cf53ea792426b97212251673788bf0b4045d0bb5bdc3843aafb9eb5ced6edc2826e734abad6235c8cf638c812247fd38f04e7080d431933b9c6d6f24756341fde3e8055dd3a83743a94dddc122ab3f32a3db0c4749ff57bad", // capnproto
  stateRoot:
    "0x0d66a8a0babec3d38b67b5239c1683f15a57e087f3825fac3d70fd6a243ed30b", // stateRoot
  // Mock a merkle proof for a burned UTXO on alicenet
  burnedUTXO: {
    chainId: 0,
    owner: "0x9AC1c9afBAec85278679fF75Ef109217f26b1417",
    value: 100,
    fee: 1,
    txHash:
      "0x0000000000000000000000000000000000000000000000000000000000000000",
  },
};

export async function init(fixture: Fixture) {
  let signers = await ethers.getSigners();
  [testData.admin, testData.user, testData.user2] = signers;
  const BridgePoolErrorCodesContract = await ethers.getContractFactory(
    "BridgePoolErrorCodes"
  );
  testData.bridgePoolErrorCodesContract =
    await BridgePoolErrorCodesContract.deploy();
  await testData.bridgePoolErrorCodesContract.deployed();
  const ImmutableAuthErrorCodesContract = await ethers.getContractFactory(
    "ImmutableAuthErrorCodes"
  );
  testData.immutableAuthErrorCodesContract =
    await ImmutableAuthErrorCodesContract.deploy();
  await testData.immutableAuthErrorCodesContract.deployed();
  testData.ethIn = ethers.utils.parseEther(testData.bTokenFeeInETH.toString());
  // mint and approve some ERC20 tokens to deposit
  await factoryCallAnyFixture(fixture, "aTokenMinter", "mint", [
    testData.user.address,
    testData.totalErc20Amount,
  ]);
  await fixture.aToken
    .connect(testData.user)
    .approve(fixture.bridgePool.address, testData.totalErc20Amount);
  // mint and approve some bTokens to deposit (and burn)
  await callFunctionAndGetReturnValues(
    fixture.bToken,
    "mintTo",
    testData.admin,
    [testData.user.address, 0],
    testData.ethIn
  );
  await fixture.bToken
    .connect(testData.user)
    .approve(fixture.bridgePool.address, BigNumber.from(testData.bTokenAmount));
  // Calculate eths to be received by burning bTokens
  testData.ethsReceived = await fixture.bToken.bTokensToEth(
    await fixture.bToken.getPoolBalance(),
    await fixture.bToken.totalSupply(),
    testData.bTokenAmount
  );
  let encodedMockBlockClaims = getMockBlockClaimsForStateRoot(
    testData.stateRoot
  );
  // Take a mock snapshot
  fixture.snapshots.snapshot(Buffer.from("0x0"), encodedMockBlockClaims);
  showState("Initial", await getState(fixture));
}

export function getMockBlockClaimsForStateRoot(stateRoot: string) {
  let encodedMockBlockClaims = defaultAbiCoder.encode(
    ["uint32", "uint32", "uint32", "bytes32", "bytes32", "bytes32", "bytes32"],
    [
      0,
      0,
      0,
      "0x0000000000000000000000000000000000000000000000000000000000000000",
      "0x0000000000000000000000000000000000000000000000000000000000000000",
      stateRoot,
      "0x0000000000000000000000000000000000000000000000000000000000000000",
    ]
  );
  return encodedMockBlockClaims;
}
