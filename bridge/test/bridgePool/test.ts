import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber, Contract } from "ethers";
import { defaultAbiCoder } from "ethers/lib/utils";
import { ethers } from "hardhat";
import { expect } from "../chai-setup";
import {
  callFunctionAndGetReturnValues,
  factoryCallAnyFixture,
  Fixture,
  getFixture,
} from "../setup";
import { getState, init, showState, state } from "./setup";
import keccak256 = require("keccak256");

describe("Testing BridgePool methods", async () => {
  let immutableAuthErrorCodesContract: Contract;
  let bridgePoolErrorCodesContract: Contract;
  let admin: SignerWithAddress;
  let user: SignerWithAddress;
  let user2: SignerWithAddress;
  let fixture: Fixture;
  const bTokenFeeInETH = 10;
  const totalErc20Amount = BigNumber.from(20000).toBigInt();
  const erc20Amount = BigNumber.from(100).toBigInt();
  let bTokenFee: BigNumber;
  let ethIn: BigNumber;
  let ethsReceived: BigNumber;
  let erc20AmountWei: BigNumber;
  let expectedState: state;
  // Mock a merkle proof for a burned UTXO on alicenet
  // Merkle tree hashed key
  let keyHash = ethers.utils.keccak256(Buffer.from("A"));
  // Merkle tree value (alicenet burned UTXO)
  let burnedUTXO = {
    chainId: 0,
    owner: "0x9AC1c9afBAec85278679fF75Ef109217f26b1417",
    value: 100,
    fee: 1,
    txHash:
      "0x0000000000000000000000000000000000000000000000000000000000000000",
  };
  // Merkle tree hashed value
  let valueHash = ethers.utils.keccak256(
    Buffer.from(JSON.stringify(burnedUTXO))
  );
  // Encode burned UTXO
  let encodedBurnedUTXO = ethers.utils.defaultAbiCoder.encode(
    [
      "tuple(uint256 chainId, address owner, uint256 value, uint256 fee, bytes32 txHash)",
    ],
    [burnedUTXO]
  );
  // These values can be obtained from accusation_builder_test.go
  let merkleProof = // capnproto
    "0x010005cda80a6c60e1215c1882b25b4744bd9d95c1218a2fd17827ab809c68196fd9bf0000000000000000000000000000000000000000000000000000000000000000af469f3b9864a5132323df8bdd9cbd59ea728cd7525b65252133a5a02f1566ee00010003a8793650a7050ac58cf53ea792426b97212251673788bf0b4045d0bb5bdc3843aafb9eb5ced6edc2826e734abad6235c8cf638c812247fd38f04e7080d431933b9c6d6f24756341fde3e8055dd3a83743a94dddc122ab3f32a3db0c4749ff57bad";
  let stateRoot = // stateRoot
    "0x0d66a8a0babec3d38b67b5239c1683f15a57e087f3825fac3d70fd6a243ed30b";

  beforeEach(async function () {
    const BridgePoolErrorCodesContract = await ethers.getContractFactory(
      "BridgePoolErrorCodes"
    );
    bridgePoolErrorCodesContract = await BridgePoolErrorCodesContract.deploy();
    await bridgePoolErrorCodesContract.deployed();
    const ImmutableAuthErrorCodesContract = await ethers.getContractFactory(
      "ImmutableAuthErrorCodes"
    );
    immutableAuthErrorCodesContract =
      await ImmutableAuthErrorCodesContract.deploy();
    await immutableAuthErrorCodesContract.deployed();
    fixture = await getFixture(false, true, false);
    let signers = await ethers.getSigners();
    [admin, user, user2] = signers;
    await init(fixture);
    // let expectedState = await getState(contractAddresses, userAddresses);
    // await factoryCallAnyFixture(fixture, "aToken", "setAdmin", [admin.address]);
    ethIn = ethers.utils.parseEther(bTokenFeeInETH.toString());
    erc20AmountWei = ethers.utils.parseUnits(erc20Amount.toString());
    // mint and approve some ERC20 tokens to deposit
    await factoryCallAnyFixture(fixture, "aTokenMinter", "mint", [
      user.address,
      totalErc20Amount,
    ]);
    await fixture.aToken
      .connect(user)
      .approve(fixture.bridgePool.address, totalErc20Amount);
    // mint and approve some bTokens to deposit (and burn)
    [bTokenFee] = await callFunctionAndGetReturnValues(
      fixture.bToken,
      "mintTo",
      admin,
      [user.address, 0],
      ethIn
    );
    await fixture.bToken
      .connect(user)
      .approve(fixture.bridgePool.address, BigNumber.from(bTokenFee));
    // Calculate eths to be received by burning bTokens
    ethsReceived = await fixture.bToken.bTokensToEth(
      await fixture.bToken.getPoolBalance(),
      await fixture.bToken.totalSupply(),
      bTokenFee
    );
    showState("Initial", await getState(fixture));
  });

  describe("Testing access control", async () => {
    it("Should not make a deposit if not via factory contract", async () => {
      await expect(
        fixture.bridgePool.deposit(
          1,
          user.address,
          erc20Amount,
          bTokenFee,
          user.address
        )
      ).to.be.revertedWith(
        immutableAuthErrorCodesContract.IMMUTEABLEAUTH_ONLY_FACTORY()
      );
    });
    it("Should not make a withdraw if not via factory contract", async () => {
      await expect(
        fixture.bridgePool.withdraw(
          merkleProof,
          encodedBurnedUTXO,
          stateRoot,
          user.address
        )
      ).to.be.revertedWith(
        immutableAuthErrorCodesContract.IMMUTEABLEAUTH_ONLY_FACTORY()
      );
    });
  });

  describe("Testing business logic", async () => {
    it("Should make a deposit with parameters and emit correspondent event", async () => {
      expectedState = await getState(fixture);
      expectedState.Balances.aToken.user -= erc20Amount;
      expectedState.Balances.aToken.bridgePool += erc20Amount;
      expectedState.Balances.bToken.user -= bTokenFee.toBigInt();
      expectedState.Balances.bToken.totalSupply -= bTokenFee.toBigInt();
      expectedState.Balances.eth.bridgePool += ethsReceived.toBigInt();
      let receipt = await factoryCallAnyFixture(
        fixture,
        "bridgePool",
        "deposit",
        [1, user.address, erc20Amount, bTokenFee, user.address]
      );
      let encodedEventEmitted = receipt.events?.at(-1)?.data as string;
      let nonce = 1;
      let networkId = 0;
      let eventExpected = [
        BigNumber.from(nonce),
        fixture.aToken.address,
        user.address,
        BigNumber.from(erc20Amount),
        BigNumber.from(networkId),
      ];
      let eventEmitted = defaultAbiCoder.decode(
        ["uint256", "address", "address", "uint256", "uint256"],
        encodedEventEmitted
      );
      expect(eventExpected).to.be.deep.equal(eventEmitted);
      showState("After Deposit", await getState(fixture));
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should make a withdraw for amount specified on burned UTXO with verified proof", async () => {
      // Make first a deposit to withdraw afterwards
      await factoryCallAnyFixture(fixture, "bridgePool", "deposit", [
        1,
        user.address,
        erc20Amount,
        bTokenFee,
        user.address,
      ]);
      showState("After Deposit", await getState(fixture));
      expectedState = await getState(fixture);
      expectedState.Balances.aToken.user += erc20Amount;
      expectedState.Balances.aToken.bridgePool -= erc20Amount;
      await factoryCallAnyFixture(fixture, "bridgePool", "withdraw", [
        merkleProof,
        encodedBurnedUTXO,
        stateRoot,
        user.address,
      ]);
      showState("After withdraw", await getState(fixture));
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should not make a withdraw for amount specified on burned UTXO with wrong merkle proof", async () => {
      let wrongMerkleProof =
        "0x016665cda80a6c60e1215c1882b25b4744bd9d95c1218a2fd17827ab809c68196fd9bf0000000000000000000000000000000000000000000000000000000000000000af469f3b9864a5132323df8bdd9cbd59ea728cd7525b65252133a5a02f1566ee00010003a8793650a7050ac58cf53ea792426b97212251673788bf0b4045d0bb5bdc3843aafb9eb5ced6edc2826e734abad6235c8cf638c812247fd38f04e7080d431933b9c6d6f24756341fde3e8055dd3a83743a94dddc122ab3f32a3db0c4749ff57bad";
      await expect(
        factoryCallAnyFixture(fixture, "bridgePool", "withdraw", [
          wrongMerkleProof,
          encodedBurnedUTXO,
          stateRoot,
          user.address,
        ])
      ).to.be.revertedWith(
        bridgePoolErrorCodesContract.BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER()
      );
    });

    it("Should not make a withdraw for amount specified on burned UTXO with wrong root", async () => {
      let wrongStateRoot =
        "0xFF66a8a0babec3d38b67b5239c1683f15a57e087f3825fac3d70fd6a243ed30b";
      await expect(
        factoryCallAnyFixture(fixture, "bridgePool", "withdraw", [
          merkleProof,
          encodedBurnedUTXO,
          wrongStateRoot,
          user.address,
        ])
      ).to.be.revertedWith(
        bridgePoolErrorCodesContract.BRIDGEPOOL_PROOF_OF_BURN_NOT_VERIFIED()
      );
    });

    it("Should not make a withdraw for an address that is not the owner", async () => {
      await expect(
        factoryCallAnyFixture(fixture, "bridgePool", "withdraw", [
          merkleProof,
          encodedBurnedUTXO,
          stateRoot,
          user2.address,
        ])
      ).to.be.revertedWith(
        bridgePoolErrorCodesContract.BRIDGEPOOL_RECEIVER_NOT_PROOF_OF_BURN_OWNER()
      );
    });

    it("Should make a withdraw for amount specified on burned UTXO with verified binary proof", async () => {
      // Make first a deposit to withdraw afterwards
      await factoryCallAnyFixture(fixture, "bridgePool", "deposit", [
        1,
        user.address,
        erc20Amount,
        bTokenFee,
        user.address,
      ]);
      showState("After Deposit", await getState(fixture));
      expectedState = await getState(fixture);
      expectedState.Balances.aToken.user += erc20Amount;
      expectedState.Balances.aToken.bridgePool -= erc20Amount;
      await factoryCallAnyFixture(fixture, "bridgePool", "withdraw", [
        merkleProof,
        encodedBurnedUTXO,
        stateRoot,
        user.address,
      ]);
      showState("After withdraw", await getState(fixture));
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });
  });
});
