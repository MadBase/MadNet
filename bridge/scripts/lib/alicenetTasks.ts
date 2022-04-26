import toml from "@iarna/toml";
import fs from "fs";
import { task, types } from "hardhat/config";
import { DEFAULT_CONFIG_OUTPUT_DIR } from "./constants";
import { readDeploymentArgs } from "./deployment/deploymentConfigUtil";

function delay(milliseconds: number) {
  return new Promise((resolve) => setTimeout(resolve, milliseconds));
}

export async function getTokenIdFromTx(ethers: any, tx: any) {
  const abi = [
    "event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)",
  ];
  const iface = new ethers.utils.Interface(abi);
  const receipt = await tx.wait();
  const log = iface.parseLog(receipt.logs[2]);
  return log.args[2];
}

async function waitBlocks(waitingBlocks: number, hre: any) {
  let constBlock = await hre.ethers.provider.getBlockNumber();
  const expectedBlock = constBlock + waitingBlocks;
  console.log(
    `Current block: ${constBlock} Waiting for ${waitingBlocks} blocks to be mined!`
  );
  while (constBlock < expectedBlock) {
    constBlock = await hre.ethers.provider.getBlockNumber();
    console.log(`Current block: ${constBlock}`);
    await delay(10000);
  }
}

task(
  "deployLegacyTokenAndUpdateDeploymentArgs",
  "Computes factory address and to the deploymentArgs file"
)
  .addOptionalParam(
    "deploymentArgsTemplatePath",
    "path of the deploymentArgsTemplate file",
    DEFAULT_CONFIG_OUTPUT_DIR + "/deploymentArgsTemplate"
  )
  .addOptionalParam(
    "outputFolder",
    "path of the output folder where new deploymentArgsTemplate file will be saved",
    "../scripts/generated"
  )
  .setAction(async (taskArgs, hre) => {
    if (!fs.existsSync(taskArgs.deploymentArgsTemplatePath)) {
      throw new Error(
        `Error: Could not find deployment Args file expected at ${taskArgs.deploymentArgsTemplatePath}`
      );
    }
    if (!fs.existsSync(taskArgs.outputFolder)) {
      throw new Error(
        `Error: Output folder  ${taskArgs.outputFolder} doesn't exist!`
      );
    }
    console.log(
      `Loading deploymentArgs from: ${taskArgs.deploymentArgsTemplatePath}`
    );

    const deploymentConfig: any = await readDeploymentArgs(
      taskArgs.deploymentArgsTemplatePath
    );

    const expectedContract = "contracts/AToken.sol:AToken";
    const expectedField = "legacyToken_";
    if (deploymentConfig.constructor[expectedContract] === undefined) {
      throw new Error(
        `Couldn't find ${expectedField} in the constructor area for` +
          ` ${expectedContract} inside the ${taskArgs.deploymentArgsTemplatePath}`
      );
    }

    // Make sure that admin is the named account at position 0
    const [admin] = await hre.ethers.getSigners();
    console.log(`Admin address: ${admin.address}`);

    const legacyToken = await (
      await hre.ethers.getContractFactory("LegacyToken")
    )
      .connect(admin)
      .deploy();

    await (await legacyToken.connect(admin).initialize()).wait();
    console.log(
      `Minted ${await legacyToken.balanceOf(admin.address)} tokens for user: ${
        admin.address
      }`
    );

    console.log(`Deployed legacy token at: ${legacyToken.address}`);
    deploymentConfig.constructor[expectedContract][0] = {
      legacyToken_: legacyToken.address,
    };

    const data = toml.stringify(deploymentConfig);
    fs.writeFileSync(taskArgs.outputFolder + "/deploymentArgsTemplate", data);
  });

task(
  "deployStateMigrationContract",
  "Deploy state migration contract and run migrations"
)
  .addParam(
    "factoryAddress",
    "the default factory address from factoryState will be used if not set"
  )
  .addOptionalParam("migrationAddress", "the address of the migration contract")
  .addFlag(
    "skipFirstTransaction",
    "The task executes 2 tx to execute the migrations." +
      " Use this flag if you want to skip the first tx where we mint the NFT."
  )
  .setAction(async (taskArgs, hre) => {
    if (
      taskArgs.factoryAddress === undefined ||
      taskArgs.factoryAddress === ""
    ) {
      throw new Error("Expected a factory address to be passed!");
    }
    // Make sure that admin is the named account at position 0
    const [admin] = await hre.ethers.getSigners();
    console.log(`Admin address: ${admin.address}`);

    const factory = await hre.ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );

    let stateMigration;
    if (
      taskArgs.migrationAddress === undefined ||
      taskArgs.migrationAddress === ""
    ) {
      console.log("Deploying migration contract!");
      stateMigration = await (
        await hre.ethers.getContractFactory("StateMigration")
      )
        .connect(admin)
        .deploy(taskArgs.factoryAddress);

      await waitBlocks(6, hre);

      console.log("Deployed migration contract at " + stateMigration.address);
    } else {
      stateMigration = await hre.ethers.getContractAt(
        "StateMigration",
        taskArgs.migrationAddress
      );
      console.log(
        "Using migration contract deployed at " + stateMigration.address
      );
    }

    if (
      taskArgs.skipFirstTransaction === undefined ||
      taskArgs.skipFirstTransaction === false
    ) {
      console.log("Calling the contract first time to mint and stake NFTs!");
      await (
        await factory.delegateCallAny(
          stateMigration.address,
          stateMigration.interface.encodeFunctionData("doMigrationStep")
        )
      ).wait();

      await waitBlocks(3, hre);
    }
    console.log(
      "Calling the contract second time to register and migrate state!"
    );
    await (
      await factory.delegateCallAny(
        stateMigration.address,
        stateMigration.interface.encodeFunctionData("doMigrationStep")
      )
    ).wait();

    await waitBlocks(3, hre);
  });

task("registerValidators", "registers validators")
  .addParam("factoryAddress", "address of the factory deploying the contract")
  .addVariadicPositionalParam(
    "addresses",
    "validators' addresses",
    undefined,
    types.string,
    false
  )
  .setAction(async (taskArgs, hre) => {
    console.log("registerValidators", taskArgs.addresses);
    const factory = await hre.ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );

    // checking factory address
    factory
      .lookup(hre.ethers.utils.formatBytes32String("AToken"))
      .catch((error) => {
        console.log(
          `Invalid factory-address ${taskArgs.factoryAddress}!\n${error}`
        );
      });
    const validatorAddresses: string[] = taskArgs.addresses;
    console.log(validatorAddresses);
    // Make sure that admin is the named account at position 0
    const [admin] = await hre.ethers.getSigners();
    console.log(`Admin address: ${admin.address}`);

    const registrationContract = await (
      await hre.ethers.getContractFactory("RegisterValidators")
    )
      .connect(admin)
      .deploy(taskArgs.factoryAddress);
    await registrationContract.deployTransaction.wait(3);

    const validatorPool = await hre.ethers.getContractAt(
      "ValidatorPool",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("ValidatorPool")
      )
    );
    console.log(`validatorPool Address: ${validatorPool.address}`);
    console.log("Staking validators");
    await (
      await factory.delegateCallAny(
        registrationContract.address,
        registrationContract.interface.encodeFunctionData("stakeValidators", [
          validatorAddresses.length,
        ])
      )
    ).wait(3);

    console.log("Registering validators");
    await (
      await factory.delegateCallAny(
        registrationContract.address,
        registrationContract.interface.encodeFunctionData(
          "registerValidators",
          [validatorAddresses]
        )
      )
    ).wait(3);

    console.log("done");
  });

task("ethdkgInput", "calculate the initializeETHDKG selector").setAction(
  async (taskArgs, hre) => {
    const { ethers } = hre;
    const iface = new ethers.utils.Interface(["function initializeETHDKG()"]);
    const input = iface.encodeFunctionData("initializeETHDKG");
    console.log("input", input);
  }
);

task("virtualMintDeposit", "Virtually creates a deposit on the side chain")
  .addParam(
    "factoryAddress",
    "the default factory address from factoryState will be used if not set",
    undefined,
    types.string
  )
  .addParam(
    "depositOwnerAddress",
    "the address of the account that will have ownership over the newly created deposit",
    undefined,
    types.string
  )
  .addParam(
    "depositAmount",
    "Amount of BTokens to be deposited",
    undefined,
    types.int
  )
  .addParam(
    "accountType",
    "For ethereum based address use number: 1  For BN curve addresses user number: 2",
    1,
    types.int
  )
  .setAction(async (taskArgs, hre) => {
    const { ethers } = hre;
    const iface = new ethers.utils.Interface([
      "function virtualMintDeposit(uint8 accountType_,address to_,uint256 amount_)",
    ]);
    const input = iface.encodeFunctionData("virtualMintDeposit", [
      taskArgs.accountType,
      taskArgs.depositOwnerAddress,
      taskArgs.depositAmount,
    ]);
    const [admin] = await ethers.getSigners();
    const adminSigner = await ethers.getSigner(admin.address);
    const factory = await ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );
    const bToken = await ethers.getContractAt(
      "BToken",
      await factory.lookup(hre.ethers.utils.formatBytes32String("BToken"))
    );
    const tx = await factory
      .connect(adminSigner)
      .callAny(bToken.address, 0, input);
    await tx.wait();
    const receipt = await ethers.provider.getTransactionReceipt(tx.hash);
    console.log(receipt);
    const intrface = new ethers.utils.Interface([
      "event DepositReceived(uint256 indexed depositID, uint8 indexed accountType, address indexed depositor, uint256 amount)",
    ]);
    const data = receipt.logs[0].data;
    const topics = receipt.logs[0].topics;
    const event = intrface.decodeEventLog("DepositReceived", data, topics);
    console.log(event);
  });

task("scheduleMaintenance", "Calls schedule Maintenance")
  .addParam(
    "factoryAddress",
    "the default factory address from factoryState will be used if not set"
  )
  .setAction(async (taskArgs, hre) => {
    const { ethers } = hre;
    const iface = new ethers.utils.Interface([
      "function scheduleMaintenance()",
    ]);
    const input = iface.encodeFunctionData("scheduleMaintenance", []);
    console.log("input", input);
    const [admin] = await ethers.getSigners();
    const adminSigner = await ethers.getSigner(admin.address);
    const factory = await ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );
    const validatorPool = await hre.ethers.getContractAt(
      "ValidatorPool",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("ValidatorPool")
      )
    );
    await (
      await factory
        .connect(adminSigner)
        .callAny(validatorPool.address, 0, input)
    ).wait();
  });

task(
  "pauseEthdkgArbitraryHeight",
  "Forcing consensus to stop on block number defined by --input"
)
  .addParam("alicenetHeight", "The block number after the latest block mined")
  .addParam(
    "factoryAddress",
    "the default factory address from factoryState will be used if not set"
  )
  .setAction(async (taskArgs, hre) => {
    const { ethers } = hre;
    const iface = new ethers.utils.Interface([
      "function pauseConsensusOnArbitraryHeight(uint256)",
    ]);
    const input = iface.encodeFunctionData("pauseConsensusOnArbitraryHeight", [
      taskArgs.alicenetHeight,
    ]);
    const [admin] = await ethers.getSigners();
    const adminSigner = await ethers.getSigner(admin.address);
    const factory = await ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );
    const validatorPool = await hre.ethers.getContractAt(
      "ValidatorPool",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("ValidatorPool")
      )
    );
    await (
      await factory
        .connect(adminSigner)
        .callAny(validatorPool.address, 0, input)
    ).wait();
  });

task("initializeEthdkg", "Start the ethdkg process")
  .addParam(
    "factoryAddress",
    "the default factory address from factoryState will be used if not set"
  )
  .setAction(async (taskArgs, hre) => {
    const { ethers } = hre;

    const [admin] = await ethers.getSigners();
    const adminSigner = await ethers.getSigner(admin.address);
    const factory = await ethers.getContractAt(
      "AliceNetFactory",
      taskArgs.factoryAddress
    );
    const validatorPool = await hre.ethers.getContractAt(
      "ValidatorPool",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("ValidatorPool")
      )
    );

    console.log("Initializing ETHDKG");
    await (
      await factory
        .connect(adminSigner)
        .callAny(
          validatorPool.address,
          0,
          validatorPool.interface.encodeFunctionData("initializeETHDKG")
        )
    ).wait(3);
    console.log("Done");
  });
