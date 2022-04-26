import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber } from "ethers";
import { ethers } from "hardhat";
import { BridgePool } from "../../typechain-types";
import { expect } from "../chai-setup";
import { Fixture, getFixture } from "../setup";
import { getState, init, showState } from "./setup";

describe("Testing BridgePool methods", async () => {
  let admin: SignerWithAddress;
  let user: SignerWithAddress;
  let user2: SignerWithAddress;
  let fixture: Fixture;
  let eth = 10;
  let aAmount = 10;
  let ethIn: BigNumber;
  let aDeposit: BigNumber;
  let bridgePool: BridgePool;

  beforeEach(async function () {
    fixture = await getFixture();
    let signers = await ethers.getSigners();
    [admin, user, user2] = signers;
    await init(fixture);
    // let expectedState = await getState(contractAddresses, userAddresses);
    showState("Initial", await getState(fixture));
    // await factoryCallAnyFixture(fixture, "aToken", "setAdmin", [admin.address]);
    ethIn = ethers.utils.parseEther(eth.toString());
    aDeposit = ethers.utils.parseUnits(aAmount.toString());
  });

  it.only("Should make a deposit and emit generic event", async () => {
    await expect(fixture.bridgePool.deposit(admin.address, ethIn.toString()))
      .to.emit(fixture.eventEmitter, "Generic")
      .withArgs("BridgePool:deposit", admin.address, ethIn);
  });

  it("Should make a deposit and emit specific event", async () => {
    await expect(fixture.bridgePool.deposit(admin.address, ethIn.toString()))
      .to.emit(fixture.eventEmitter, "BridgePoolDepositReceived")
      .withArgs(0, admin.address, ethIn);
  });

  it("Should make a deposit and emit delegate event", async () => {
    await expect(fixture.bridgePool.deposit(admin.address, ethIn.toString()))
      .to.emit(fixture.eventEmitter, "Delegated")
      .withArgs(
        "0xf8a9ea73000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000546f99f244b7b58b855330ae0e2bc1b30b41302f0000000000000000000000000000000000000000000000008ac7230489e80000000000000000000000000000000000000000000000000000000000000000000a427269646765506f6f6c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000076465706f73697400000000000000000000000000000000000000000000000000"
      );
  });
});
