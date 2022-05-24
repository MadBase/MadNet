import { assert, expect } from "chai";
import { ethers } from "hardhat";
import { Accusations } from "../../typechain-types";
import { Fixture, getFixture } from "../setup";
import {
  addValidators,
  generateSigAndPClaims0,
  generateSigAndPClaims1,
  generateSigAndPClaimsDifferentChainId,
  generateSigAndPClaimsDifferentHeight,
  generateSigAndPClaimsDifferentRound,
} from "./accusations-test-helpers";

describe("StakingPositionDescriptor: Tests StakingPositionDescriptor methods", async () => {
  let fixture: Fixture;

  let accusation: Accusations;

  beforeEach(async function () {
    fixture = await getFixture(true, true);

    accusation = fixture.accusations;
  });

  describe.only("recoverSigner:", async () => {
    it("returns signer when valid", async function () {
      const sig =
        "0x" +
        "cba766e2ba024aad86db556635cec9f104e76644b235f77759ff80bfefc990c5774d2d5ff3069a5099e4f9fadc9b08ab20472e2ef432fba94498d93c10cc584b00";
      const prefix = ethers.utils.toUtf8Bytes("");
      const message =
        "0x" +
        "54686520717569636b2062726f776e20666f782064696420736f6d657468696e67";
      const expectedAddress = "0x38e959391dD8598aE80d5d6D114a7822A09d313A";

      const who = await accusation.recoverSigner(sig, prefix, message);

      assert.equal(expectedAddress, who);
    });

    it("returns signer with pclaims data", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const prefix = ethers.utils.toUtf8Bytes("Proposal");
      const expectedAddress = "0x38e959391dD8598aE80d5d6D114a7822A09d313A";

      const who = await accusation.recoverSigner(sig0, prefix, pClaims0);

      assert.equal(expectedAddress, who);
    });
  });
  describe("recoverMadNetSigner:", async () => {
    it("returns signer when valid", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } = generateSigAndPClaims1();
      const expectedAddress = "0x38e959391dD8598aE80d5d6D114a7822A09d313A";

      const signerAccount0 = await accusation.recoverMadNetSigner(
        sig0,
        pClaims0
      );
      const signerAccount1 = await accusation.recoverMadNetSigner(
        sig1,
        pClaims1
      );

      assert.equal(expectedAddress, signerAccount0);
      assert.equal(expectedAddress, signerAccount1);
    });
  });
  describe("AccuseMultipleProposal:", async () => {
    it("returns signer when valid", async function () {
      // const madID = generateMadID(987654321);

      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } = generateSigAndPClaims1();

      // const signerAccount0 = await accusation.recoverMadNetSigner(sig0, pClaims0);
      const signerAccount0 = "0x38e959391dD8598aE80d5d6D114a7822A09d313A";

      await addValidators(fixture.validatorPool, [signerAccount0]);

      const signer = await accusation.AccuseMultipleProposal(
        sig0,
        pClaims0,
        sig1,
        pClaims1
      );

      assert.equal(signer, signerAccount0);
    });

    it("reverts when signer is not valid", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } = generateSigAndPClaims1();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims1)
      ).to.be.revertedWith(
        "Accusations: the signer of these proposals is not a valid validator!"
      );
    });

    it("reverts when duplicate data for pClaims0 and pClaims1", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig0, pClaims0)
      ).to.be.revertedWith("Accusations: the PClaims are equal!");
    });

    it("reverts when proposals have different signature", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1 } = generateSigAndPClaims1();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims0)
      ).to.be.revertedWith(
        "Accusations: the signers of the proposals should be the same"
      );
    });

    it("reverts when proposals have different block height", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } =
        generateSigAndPClaimsDifferentHeight();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims1)
      ).to.be.revertedWith(
        "Accusations: the block heights between the proposals are different!"
      );
    });

    it("reverts when proposals have different round", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } =
        generateSigAndPClaimsDifferentRound();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims1)
      ).to.be.revertedWith(
        "Accusations: the round between the proposals are different!"
      );
    });

    it("reverts when proposals have different chain id", async function () {
      const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
      const { sig: sig1, pClaims: pClaims1 } =
        generateSigAndPClaimsDifferentChainId();

      await expect(
        accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims1)
      ).to.be.revertedWith(
        "Accusations: the chainId between the proposals are different!"
      );
    });

    // TODO find a way to set the chain id to 0
    // it("reverts when proposals are for different chain id than current chain", async function () {
    //   const { sig: sig0, pClaims: pClaims0 } = generateSigAndPClaims0();
    //   const { sig: sig1, pClaims: pClaims1 } = generateSigAndPClaims1();
    //   participants.setChainId(0);
    //   await expect(
    //     accusation.AccuseMultipleProposal(sig0, pClaims0, sig1, pClaims1)
    //   ).to.be.revertedWith(
    //     "Accusations: the chainId between the proposals are different!"
    //   );
    // });
  });
});
