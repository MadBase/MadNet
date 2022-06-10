import { ethers } from "hardhat";
import { CallAny, DepositNotifier } from "../../typechain-types";
import {
  deployAliceNetFactory,
  deployUpgradeableWithFactory,
  expect,
  preFixtureSetup,
} from "../setup";

describe("depositNotifier", () => {
  it("admin not allowed to call doEmit", async () => {
    await preFixtureSetup();
    const [admin] = await ethers.getSigners();
    const factory = await deployAliceNetFactory(admin);

    const contract = (await deployUpgradeableWithFactory(
      factory,
      "DepositNotifier",
      "salt1",
      undefined,
      [1]
    )) as DepositNotifier;

    await expect(
      contract
        .connect(admin)
        .doEmit(
          ethers.utils.formatBytes32String("salt1"),
          "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
          1337,
          "0xcd3B766CCDd6AE721141F452C550Ca635964ce71"
        )
        .then((resp) => resp.wait())
    ).to.be.rejectedWith("not allowed");
  });

  it("doEmit succeeds when factory-deployed contract is calls it with its own salt", async () => {
    await preFixtureSetup();
    const [admin] = await ethers.getSigners();
    const factory = await deployAliceNetFactory(admin);

    const dnContract = (await deployUpgradeableWithFactory(
      factory,
      "DepositNotifier",
      "salt1",
      undefined,
      [777]
    )) as DepositNotifier;
    const auxillaryContract = (await deployUpgradeableWithFactory(
      factory,
      "CallAny",
      "salt2"
    )) as CallAny;

    const encodedArgs = dnContract.interface.encodeFunctionData("doEmit", [
      ethers.utils.formatBytes32String("salt2"),
      "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      1337,
      "0xcd3B766CCDd6AE721141F452C550Ca635964ce71",
    ]);
    const receipt = await auxillaryContract
      .connect(admin)
      .callAny(dnContract.address, 0, encodedArgs)
      .then((resp) => resp.wait());

    const encodedArgs2 = dnContract.interface.encodeFunctionData("doEmit", [
      ethers.utils.formatBytes32String("salt2"),
      "0xcf0a769a379d2aae3cf64b38abbafbe7653a6418",
      42,
      "0x5e0d927f4bfe097e62f6b15a11aec9843f47ba90",
    ]);
    const receipt2 = await auxillaryContract
      .connect(admin)
      .callAny(dnContract.address, 0, encodedArgs2)
      .then((resp) => resp.wait());

    if (receipt.events === undefined) {
      throw Error("receipt 1 has no events");
    } else {
      const parsed = dnContract.interface.parseLog(receipt.events[0]);
      expect(parsed.name).to.eq("Deposited");
      expect(parsed.args.nonce.toBigInt()).to.eq(1n);
      expect(parsed.args.networkId.toBigInt()).to.eq(777n);
      expect(parsed.args.ercContract).to.eq(
        "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
      );
      expect(parsed.args.number.toBigInt()).to.eq(1337n);
      expect(parsed.args.owner).to.eq(
        "0xcd3B766CCDd6AE721141F452C550Ca635964ce71"
      );
    }

    if (receipt2.events === undefined) {
      throw Error("receipt 2 has no events");
    } else {
      const parsed = dnContract.interface.parseLog(receipt2.events[0]);
      expect(parsed.name).to.eq("Deposited");
      expect(parsed.args.nonce.toBigInt()).to.eq(2n);
      expect(parsed.args.networkId.toBigInt()).to.eq(777n);
      expect(parsed.args.ercContract).to.eq(
        "0xcF0a769A379d2aaE3Cf64b38AbBafBe7653A6418"
      );
      expect(parsed.args.number.toBigInt()).to.eq(42n);
      expect(parsed.args.owner).to.eq(
        "0x5e0d927F4Bfe097E62F6B15a11aEC9843F47Ba90"
      );
    }
  });

  it("doEmit fails when factory-deployed contract is calls it with the wrong salt", async () => {
    await preFixtureSetup();
    const [admin] = await ethers.getSigners();
    const factory = await deployAliceNetFactory(admin);

    const dnContract = (await deployUpgradeableWithFactory(
      factory,
      "DepositNotifier",
      "salt1",
      undefined,
      [1]
    )) as DepositNotifier;
    const auxillaryContract = (await deployUpgradeableWithFactory(
      factory,
      "CallAny",
      "salt2"
    )) as CallAny;

    const encodedArgs = dnContract.interface.encodeFunctionData("doEmit", [
      ethers.utils.formatBytes32String("salt2XXX"),
      "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      1337,
      "0xcd3B766CCDd6AE721141F452C550Ca635964ce71",
    ]);

    await expect(
      auxillaryContract
        .connect(admin)
        .callAny(dnContract.address, 0, encodedArgs)
        .then((resp) => resp.wait())
    ).to.be.rejectedWith("not allowed");
  });
});
