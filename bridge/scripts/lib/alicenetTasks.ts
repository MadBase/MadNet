import toml from "@iarna/toml";
import { BigNumber, ContractTransaction, Transaction } from "ethers";
import fs from "fs";
import { ethers } from "hardhat";
import { task, types } from "hardhat/config";
import { ALICENET_FACTORY, DEFAULT_CONFIG_OUTPUT_DIR } from "./constants";
import { readDeploymentArgs } from "./deployment/deploymentConfigUtil";
import { getATokenAddress, getATokenMinterAddress, getBTokenAddress, getDefaultFactoryAddress } from "./deployment/factoryStateUtil";

function delay(milliseconds: number) {
  return new Promise((resolve) => setTimeout(resolve, milliseconds));
}

export async function getTokenIdFromTx(ethers: any, tx: ContractTransaction) {
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
  .addFlag(
    "migrateFromLegacy",
    "Flag indicating if the script should try to migrate legacy tokens to ATokens. " +
      "Only necessary if you only have legacy tokens and want to register validators with them!"
  )
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
    const lockTime = 1;
    const validatorAddresses: string[] = taskArgs.addresses;
    const stakingTokenIds: BigNumber[] = [];

    const aToken = await hre.ethers.getContractAt(
      "AToken",
      await factory.lookup(hre.ethers.utils.formatBytes32String("AToken"))
    );
    console.log(`AToken Address: ${aToken.address}`);

    const publicStaking = await hre.ethers.getContractAt(
      "PublicStaking",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("PublicStaking")
      )
    );
    console.log(`publicStaking Address: ${publicStaking.address}`);
    const validatorPool = await hre.ethers.getContractAt(
      "ValidatorPool",
      await factory.lookup(
        hre.ethers.utils.formatBytes32String("ValidatorPool")
      )
    );
    console.log(`validatorPool Address: ${validatorPool.address}`);
    console.log(await validatorPool.getMaxNumValidators());
    const stakeAmountATokenWei = await validatorPool.getStakeAmount();

    console.log(
      `Minimum amount ATokenWei to stake: ${stakeAmountATokenWei.toBigInt()}`
    );

    // Make sure that admin is the named account at position 0
    const [admin] = await hre.ethers.getSigners();
    console.log(`Admin address: ${admin.address}`);

    const totalStakeAmt = stakeAmountATokenWei.mul(validatorAddresses.length);
    if (taskArgs.migrateFromLegacy) {
      const legacyToken = await hre.ethers.getContractAt(
        "ERC20",
        await aToken.getLegacyTokenAddress()
      );

      const input = aToken.interface.encodeFunctionData("allowMigration");
      await factory.connect(admin).callAny(aToken.address, 0, input);

      console.log(`Legacy Token Address: ${legacyToken.address}`);
      console.log(
        `Legacy Token Balance: ${(
          await legacyToken.balanceOf(admin.address)
        ).toBigInt()}`
      );
      await (
        await legacyToken.connect(admin).approve(aToken.address, totalStakeAmt)
      ).wait();
      await (await aToken.connect(admin).migrate(totalStakeAmt)).wait();
    }
    // approve tokens
    let tx = await aToken
      .connect(admin)
      .approve(
        publicStaking.address,
        stakeAmountATokenWei.mul(validatorAddresses.length)
      );
    await tx.wait();
    console.log(
      `Approved allowance to validatorPool of: ${stakeAmountATokenWei
        .mul(validatorAddresses.length)
        .toNumber()} ATokenWei`
    );

    console.log("Starting the registration process...");
    // mint PublicStaking positions to validators
    for (let i = 0; i < validatorAddresses.length; i++) {
      let tx = await publicStaking
        .connect(admin)
        .mintTo(factory.address, stakeAmountATokenWei, lockTime);
      await tx.wait();
      const tokenId = BigNumber.from(await getTokenIdFromTx(hre.ethers, tx));
      console.log(`Minted PublicStaking.tokenID ${tokenId}`);
      stakingTokenIds.push(tokenId);
      const iface = new hre.ethers.utils.Interface([
        "function approve(address,uint256)",
      ]);
      const input = iface.encodeFunctionData("approve", [
        validatorPool.address,
        tokenId,
      ]);
      tx = await factory
        .connect(admin)
        .callAny(publicStaking.address, 0, input);

      await tx.wait();
      console.log(`Approved tokenID:${tokenId} to ValidatorPool`);
    }

    console.log(
      `registering ${validatorAddresses.length} validators with ValidatorPool...`
    );
    // add validators to the ValidatorPool
    // await validatorPool.registerValidators(validatorAddresses, stakingTokenIds)
    const iface = new hre.ethers.utils.Interface([
      "function registerValidators(address[],uint256[])",
    ]);
    const input = iface.encodeFunctionData("registerValidators", [
      validatorAddresses,
      stakingTokenIds,
    ]);
    tx = await factory.connect(admin).callAny(validatorPool.address, 0, input);
    await tx.wait();
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


  task(
    "transferEth",
    "Forcing consensus to stop on block number defined by --input"
  )
    .addParam("receiver", "address of the account to fund")
    .addParam("amount", "amount of eth to transfer")
    .setAction(async (taskArgs, hre) => {
      let accounts = await hre.ethers.getSigners();
      let ownerBal = await hre.ethers.provider.getBalance(accounts[0].address)
      let wei = BigNumber.from(parseInt(taskArgs.amount, 16)).mul(BigNumber.from("10").pow(BigInt(18)))
      let amount = wei
      const target = taskArgs.receiver
      console.log(`previous owner balance: ${ownerBal.toString()}`);
      let receiverBal = await hre.ethers.provider.getBalance(target)
      console.log(`previous receiver balance: ${receiverBal.toString()}`)
      let txRequest = await accounts[0].populateTransaction({from:accounts[0].address, value: amount, to: target});
      let txResponse = await accounts[0].sendTransaction(txRequest);
      receiverBal = await hre.ethers.provider.getBalance(target)
      console.log(`new receiver balance: ${receiverBal}`);
      let ownerBal2 = await hre.ethers.provider.getBalance(accounts[0].address)
      console.log(`new owner balance: ${ownerBal.sub(ownerBal2).toString()}`);
    });

  task(
    "mintATokenTo",
    "mints A token to an address"
  )
    .addParam("factoryAddress", "address of the factory deploying the contract")
    .addParam("amount", "amount to mint")
    .addParam("to", "address of the recipient")
    .setAction(async (taskArgs, hre) => {
      const network = hre.network.name;
      let aTokenMinterBase = await hre.ethers.getContractFactory("ATokenMinter")
      const factory = await hre.ethers.getContractAt(
        "AliceNetFactory",
        taskArgs.factoryAddress
      );
      const  aTokenMinterAddr = await factory.lookup(hre.ethers.utils.formatBytes32String("ATokenMinter"));  
      const aToken = await hre.ethers.getContractAt(
        "AToken",
        await factory.lookup(hre.ethers.utils.formatBytes32String("AToken"))
      );
      let bal1 = await aToken.callStatic.balanceOf(taskArgs.to)
      let calldata = aTokenMinterBase.interface.encodeFunctionData("mint", [taskArgs.to, taskArgs.amount]);
      //use the factory to call the A token minter
      let txResponse = await factory.callAny(aTokenMinterAddr, 0, calldata)
      await txResponse.wait();
      let bal2 = await aToken.balanceOf(taskArgs.to);
      console.log(`Minted ${bal2.sub(bal1).toString()} to account ${taskArgs.to}`)
    });

  task(
    "getATokenBalance",
    "gets AToken balance of account"
  )
    .addParam("factoryAddress", "address of the factory deploying the contract")
    .addParam("account", "address of account to get balance of")
    .setAction(async (taskArgs, hre) => {
      const factory = await hre.ethers.getContractAt(
        "AliceNetFactory",
        taskArgs.factoryAddress
      );
      const aToken = await hre.ethers.getContractAt(
        "AToken",
        await factory.lookup(hre.ethers.utils.formatBytes32String("AToken"))
      );
      let bal = await aToken.callStatic.balanceOf(taskArgs.account)
      console.log(bal)
      return bal    
    });

    task(
      "mintBTokenTo",
      "mints B token to an address"
    )
      .addParam("factoryAddress", "address of the factory deploying the contract")
      .addParam("amount", "amount to mint")
      .addParam("numWei", "amount of eth to use")
      .addParam("to", "address of the recipient")
      .setAction(async (taskArgs, hre) => {
        const network = hre.network.name;
        if (
          taskArgs.factoryAddress === undefined ||
          taskArgs.factoryAddress === ""
        ) {
          throw new Error("Expected a factory address to be passed!");
        }
        const factory = await hre.ethers.getContractAt(
          "AliceNetFactory",
          taskArgs.factoryAddress
        );
        const bToken = await hre.ethers.getContractAt(
          "BToken",
          await factory.lookup(hre.ethers.utils.formatBytes32String("BToken"))
        );
        let bal1 = await bToken.callStatic.balanceOf(taskArgs.to)
        let txResponse = await bToken.mintTo(taskArgs.to, taskArgs.amount, {value: taskArgs.numWei})
        await txResponse.wait();
        let bal2 = await bToken.callStatic.balanceOf(taskArgs.to);
        console.log(`Minted ${bal2.sub(bal1).toString()} BToken to account ${taskArgs.to}`)
      });

  task(
    "getBTokenBalance",
    "gets BToken balance of account"
  )
    .addParam("factoryAddress", "address of the factory deploying the contract")
    .addParam("account", "address of account to get balance of")
    .setAction(async (taskArgs, hre) => {
      const factory = await hre.ethers.getContractAt(
        "AliceNetFactory",
        taskArgs.factoryAddress
      );  
      const bToken = await hre.ethers.getContractAt(
        "BToken",
        await factory.lookup(hre.ethers.utils.formatBytes32String("BToken"))
      );
      let bal = await bToken.callStatic.balanceOf(taskArgs.account)
      console.log(bal)
      return bal    
    });

  task(
    "getEthBalance",
    "gets AToken balance of account"
  )
    .addParam("account", "address of account to get balance of")
    .setAction(async (taskArgs, hre) => {
      const bal = await hre.ethers.provider.getBalance(taskArgs.account)
      console.log(bal)
      return bal    
    });

  task(
    "ethToBToken",
    "gets AToken balance of account"
  )
    .addParam("factoryAddress", "address of the factory deploying the contract")
    
    .setAction(async (taskArgs, hre) => {
      const factory = await hre.ethers.getContractAt(
        "AliceNetFactory",
        taskArgs.factoryAddress
      );  
      const bToken = await hre.ethers.getContractAt(
        "BToken",
        await factory.lookup(hre.ethers.utils.formatBytes32String("BToken"))
      );
      for(let i = 0; i < 100; i++){
        let poolBal = await bToken.getPoolBalance();
        let eth = await bToken.ethToBTokens(poolBal , i);
        console.log(i, eth.toNumber())
      }
    });

function getPriceForBToken(numBtoken: number){

} 