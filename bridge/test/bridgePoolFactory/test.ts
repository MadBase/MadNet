import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Contract } from "ethers";
import { ethers } from "hardhat";
import { AliceNetFactory } from "../../typechain-types";
import { expect } from "../chai-setup";
import {
  deployFactory,
  deployStaticWithFactory,
  Fixture,
  getFixture,
} from "../setup";
import {} from "./setup";

let fixture: Fixture;
describe("BridgePool Contract Factory", () => {
  let firstOwner: SignerWithAddress;
  let bridgePoolFactory: AliceNetFactory;
  let bridgePool: Contract;
  let fixture: Fixture;

  beforeEach(async () => {
    fixture = await getFixture(false, false, false);
    [firstOwner] = await ethers.getSigners();
    bridgePoolFactory = (await deployFactory(
      "BridgePoolFactory",
      firstOwner
    )) as AliceNetFactory;
    console.log("BridgePoolFactory deployed to:", bridgePoolFactory.address);
  });

  it("should deploy BridgePool contract", async () => {
    bridgePool = await deployStaticWithFactory(
      bridgePoolFactory,
      "BridgePool",
      undefined,
      [fixture.aToken.address, fixture.bToken.address],
      [fixture.aToken.address, fixture.bToken.address]
    );
    await expect(
      bridgePool.deposit(1, firstOwner.getAddress(), 1, 1)
    ).to.be.revertedWith("ERC20: insufficient allowance");
  });
});
