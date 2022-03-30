import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber } from "ethers";
import { ethers } from "hardhat";
import { BridgePool } from "../../typechain-types";
import { expect } from "../chai-setup";
import {
  callFunctionAndGetReturnValues,
  factoryCallAny,
  Fixture,
  getFixture,
} from "../setup";
import { getState, init, showState } from "./setup";

describe("Testing BridgePool methods", async () => {
  let admin: SignerWithAddress;
  let user: SignerWithAddress;
  let user2: SignerWithAddress;
  let fixture: Fixture;
  let eth = 10;
  let mad = 10;
  let ethIn: BigNumber;
  let madBytes: BigNumber;
  let madDeposit: BigNumber;
  let minMadBytes = 0;
  let bridgePool: BridgePool;

  beforeEach(async function () {
    fixture = await getFixture();
    // const EventEmitter = await ethers.getContractFactory("EventEmitter");
    // let eventEmitter = (await EventEmitter.deploy()) as EventEmitter;
    // await eventEmitter.deployed();
    // const BridgePool = await ethers.getContractFactory("BridgePool");
    // bridgePool = (await BridgePool.deploy(eventEmitter.address)) as BridgePool;
    // await bridgePool.deployed();

    let signers = await ethers.getSigners();
    [admin, user, user2] = signers;
    await init(fixture);
    // let expectedState = await getState(contractAddresses, userAddresses);
    await factoryCallAny(fixture, "madByte", "setAdmin", [admin.address]);
    showState("Initial", await getState(fixture));
    ethIn = ethers.utils.parseEther(eth.toString());
    madDeposit = ethers.utils.parseUnits(mad.toString());
    [madBytes] = await callFunctionAndGetReturnValues(
      fixture.madByte,
      "mint",
      admin,
      [minMadBytes],
      ethIn
    );
    await fixture.madByte.approve(fixture.bridgePool.address, madBytes);
    showState("After Mint", await getState(fixture));
  });

  it("Should make a deposit transfering and burning tokens and emit correspondent event", async () => {
    // Expect a deposit event from BridgePool Contract
    let encodedEventArgs = ethers.utils.defaultAbiCoder.encode(
      ["string", "string", "uint256", "address", "uint256"],
      ["BridgePool", "deposit", BigNumber.from(1), admin.address, madBytes]
    );
    showState("Before Deposit", await getState(fixture));
    await expect(fixture.bridgePool.deposit(admin.address, madBytes))
      .to.emit(fixture.eventEmitter, "GenericEvent")
      .withArgs(encodedEventArgs);
    showState("After Deposit", await getState(fixture));
  });

  it("Should prepare a withdrawal and emit correspondent event", async () => {
    const [depositID] = await callFunctionAndGetReturnValues(
      fixture.bridgePool,
      "deposit",
      admin,
      [admin.address, madBytes]
    );
    // Expect a withdrawal event from BridgePool Contract
    let encodedEventArgs = ethers.utils.defaultAbiCoder.encode(
      ["string", "string", "uint256"],
      ["BridgePool", "withdrawal", depositID]
    );
    showState("Before Withdrawal", await getState(fixture));
    await expect(fixture.bridgePool.withdraw(depositID))
      .to.emit(fixture.eventEmitter, "GenericEvent")
      .withArgs(encodedEventArgs);
    showState("After Withdrawal", await getState(fixture));
  });

  it.only("Should withraw eth if proof of burn was confirmed", async () => {
    const [depositID] = await callFunctionAndGetReturnValues(
      fixture.bridgePool,
      "deposit",
      admin,
      [admin.address, madBytes]
    );
    showState("After Deposit", await getState(fixture));
    // Prepare the withdrawal
    await fixture.bridgePool.withdraw(depositID);
    showState("After Withdraw", await getState(fixture));
    // Simulate successful validation of a proof of burn inside of the state proof
    await fixture.bridgePool.confirmProofOfBurn(depositID);
    // let eths = await fixture.madByte.madByteToEth(
    //   await fixture.madByte.getPoolBalance(),
    //   await fixture.madByte.totalSupply(),
    //   madBytes
    // );
    // Expect a distribute event from BridgePool Contract
    let encodedEventArgs = ethers.utils.defaultAbiCoder.encode(
      ["string", "string", "uint256", "address"],
      ["BridgePool", "distribute", depositID, admin.address]
    );
    await expect(fixture.bridgePool.distribute(depositID))
      .to.emit(fixture.eventEmitter, "GenericEvent")
      .withArgs(encodedEventArgs);
    showState("After Distribute", await getState(fixture));
  });
});
