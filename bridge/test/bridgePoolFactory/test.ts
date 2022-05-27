import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { ethers } from "hardhat";
import { expect } from "../chai-setup";
import {
  Fixture,
  getContractAddressFromDeployedProxyEvent,
  getFixture,
} from "../setup";

let fixture: Fixture;
describe("BridgePool Contract Factory", () => {
  let firstOwner: SignerWithAddress;
  let fixture: Fixture;

  beforeEach(async () => {
    fixture = await getFixture(false, false, false);
    [firstOwner] = await ethers.getSigners();
  });
  describe("Testing Access control", () => {
    it("should not deploy new BridgePool with BridgePoolFactory not being delegator", async () => {
      await expect(
        fixture.bridgePoolFactory.deployNewPool(
          fixture.aToken.address,
          fixture.bToken.address
        )
      ).to.be.revertedWith("900");
    });

    it("should not deploy new BridgePool with BridgePoolFactory trying to access factory logic directly", async () => {
      await expect(
        fixture.bridgePoolFactory.deployViaFactoryLogic(
          fixture.aToken.address,
          fixture.bToken.address
        )
      ).to.be.revertedWith("2018");
    });

    it("should deploy new BridgePool with BridgePoolFactory being delegator", async () => {
      await fixture.factory.setDelegator(fixture.bridgePoolFactory.address);
      const deployNewPoolTransaction =
        await fixture.bridgePoolFactory.deployNewPool(
          fixture.aToken.address,
          fixture.bToken.address
        );
      const bridgePoolAddress = await getContractAddressFromDeployedProxyEvent(
        deployNewPoolTransaction
      );
      const bridgePool = (await ethers.getContractFactory("BridgePool")).attach(
        bridgePoolAddress
      );
      await expect(
        bridgePool.deposit(1, firstOwner.address, 1, 1)
      ).to.be.revertedWith("ERC20: insufficient allowance");
    });

    it("should not deploy two BridgePools with same ERC20 contract", async () => {
      await fixture.factory.setDelegator(fixture.bridgePoolFactory.address);
      const deployNewPoolTransaction =
        await fixture.bridgePoolFactory.deployNewPool(
          fixture.aToken.address,
          fixture.bToken.address
        );
      const bridgePoolAddress = await getContractAddressFromDeployedProxyEvent(
        deployNewPoolTransaction
      );
      await expect(
        fixture.bridgePoolFactory.deployNewPool(
          fixture.aToken.address,
          fixture.bToken.address
        )
      ).to.be.revertedWith("901");
    });
  });
});
