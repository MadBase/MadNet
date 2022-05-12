import { BigNumber } from "ethers";
import { ethers } from "hardhat";
import { expect } from "../chai-setup";
import { Fixture, getFixture } from "../setup";
import {
  getMockBlockClaimsForStateRoot,
  getState,
  init,
  showState,
  state,
  testData,
} from "./setup";
let fixture: Fixture;
let expectedState: state;
const encodedBurnedUTXO = ethers.utils.defaultAbiCoder.encode(
  [
    "tuple(uint256 chainId, address owner, uint256 value, uint256 fee, bytes32 txHash)",
  ],
  [testData.burnedUTXO]
);

describe("Testing BridgePool methods", async () => {
  beforeEach(async function () {
    fixture = await getFixture(false, true, false);
    await init(fixture);
  });

  describe("Testing business logic", async () => {
    it("Should make a deposit with parameters and emit correspondent event", async () => {
      expectedState = await getState(fixture);
      expectedState.Balances.aToken.user -= testData.erc20Amount;
      expectedState.Balances.aToken.bridgePool += testData.erc20Amount;
      expectedState.Balances.bToken.user -= testData.bTokenAmount;
      expectedState.Balances.bToken.totalSupply -= testData.bTokenAmount;
      expectedState.Balances.eth.bridgePool += testData.ethsReceived.toBigInt();
      const nonce = 1;
      const networkId = 0;
      await expect(
        fixture.bridgePool
          .connect(testData.user)
          .deposit(
            1,
            testData.user.address,
            testData.erc20Amount,
            testData.bTokenAmount
          )
      )
        // TODO: change to  .to.emit(fixture.depositNotifier, "Deposited") upon merging of PR-126
        .to.emit(fixture.bridgePool, "Deposited")
        .withArgs(
          BigNumber.from(nonce),
          fixture.aToken.address,
          testData.user.address,
          BigNumber.from(testData.erc20Amount),
          BigNumber.from(networkId)
        );
      showState("After Deposit", await getState(fixture));
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should make a withdraw for amount specified on burned UTXO upon proof verification", async () => {
      // Make first a deposit to withdraw afterwards
      await fixture.bridgePool
        .connect(testData.user)
        .deposit(
          1,
          testData.user.address,
          testData.erc20Amount,
          testData.bTokenAmount
        );
      showState("After Deposit", await getState(fixture));
      expectedState = await getState(fixture);
      expectedState.Balances.aToken.user += testData.erc20Amount;
      expectedState.Balances.aToken.bridgePool -= testData.erc20Amount;
      await fixture.bridgePool
        .connect(testData.user)
        .withdraw(testData.merkleProof, encodedBurnedUTXO);
      showState("After withdraw", await getState(fixture));
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should not make a withdraw for amount specified on burned UTXO with not verified merkle proof", async () => {
      const wrongMerkleProof =
        "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000";
      expectedState = await getState(fixture);
      const reason = ethers.utils.parseBytes32String(
        await testData.bridgePoolErrorCodesContract.BRIDGEPOOL_COULD_NOT_VERIFY_PROOF_OF_BURN()
      );
      await expect(
        fixture.bridgePool
          .connect(testData.user)
          .withdraw(wrongMerkleProof, encodedBurnedUTXO)
      ).to.be.revertedWith(reason);
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should not make a withdraw for amount specified on burned UTXO with wrong root", async () => {
      const wrongStateRoot =
        "0x0000000000000000000000000000000000000000000000000000000000000000";
      const encodedMockBlockClaims =
        getMockBlockClaimsForStateRoot(wrongStateRoot);
      fixture.snapshots.snapshot(Buffer.from("0x0"), encodedMockBlockClaims);
      expectedState = await getState(fixture);
      const reason = ethers.utils.parseBytes32String(
        await testData.bridgePoolErrorCodesContract.BRIDGEPOOL_COULD_NOT_VERIFY_PROOF_OF_BURN()
      );
      await expect(
        fixture.bridgePool
          .connect(testData.user)
          .withdraw(testData.merkleProof, encodedBurnedUTXO)
      ).to.be.revertedWith(reason);
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });

    it("Should not make a withdraw to an address that is not the owner in burned UTXO", async () => {
      const reason = ethers.utils.parseBytes32String(
        await testData.bridgePoolErrorCodesContract.BRIDGEPOOL_RECEIVER_IS_NOT_OWNER_ON_PROOF_OF_BURN_UTXO()
      );
      expectedState = await getState(fixture);
      await expect(
        fixture.bridgePool
          .connect(testData.user2)
          .withdraw(testData.merkleProof, encodedBurnedUTXO)
      ).to.be.revertedWith(reason);
      expect(await getState(fixture)).to.be.deep.equal(expectedState);
    });
  });
});
