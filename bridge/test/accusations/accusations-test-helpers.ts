import { expect } from "chai";
import { ValidatorPool, ValidatorPoolMock } from "../../typechain-types";

export const generateSigAndPClaims0 = () => {
  const pClaims =
    "0x" +
    "0000000000000200" +
    "04000000" +
    "02000400" +
    "58000000" +
    "00000200" +
    // BClaim
    "01000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "0d000000" +
    "02010000" +
    "19000000" +
    "02010000" +
    "25000000" +
    "02010000" +
    "31000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "de8b68a6643fa528a513f99a1ea30379927197a097ca86d9108e4c29d684b1ec" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    // Rcert
    "04000000" +
    "02000100" +
    "1d000000" +
    "02060000" +
    // RClaim
    "01000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "01000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    // SigGroup
    "258aa89365a642358d92db67a13cb25d73e6eedf0d25100d8d91566882fac54b" +
    "1ccedfb0425434b54999a88cd7d993e05411955955c0cfec9dd33066605bd4a6" +
    "0f6bbfbab37349aaa762c23281b5749932c514f3b8723cf9bb05f9841a7f2d0e" +
    "0f75e42fd6c8e9f0edadac3dcfb7416c2d4b2470f4210f2afa93138615b1deb1" +
    "1ff56a9538b079e16dd77a8ef81318497b195ad81b8cd1c5ea5d48b0c160f599" +
    "12387b5ab69538ef4cda0f7a879982f9b4943291b1e6d998abefe7bb4ebb6993";
  const sig =
    "0x" +
    "05d08b3bfd0fcb21e00a1468a9013fb023aa5eb86d714600dd69675ef9acce8c" +
    "3247fe575e3d16a3e32d1e0ea10a30474744e7aab3166daea7c591776c1e942500";
  return { sig, pClaims };
};

export const generateSigAndPClaims1 = () => {
  // Second PClaims, this time containg 2 transactions instead of 1.
  const pClaims =
    "0x" +
    "0000000000000200" +
    "04000000" +
    "02000400" +
    "58000000" +
    "00000200" +
    "01000000" +
    "02000000" +
    "02000000" + // txCount is different
    "00000000" +
    "0d000000" +
    "02010000" +
    "19000000" +
    "02010000" +
    "25000000" +
    "02010000" +
    "31000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "8db49c526748abbf9eabf4b49e9edd6d91ca3c970791b027e815b628c148afd0" + // txRoot is different
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "04000000" +
    "02000100" +
    "1d000000" +
    "02060000" +
    "01000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "01000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "258aa89365a642358d92db67a13cb25d73e6eedf0d25100d8d91566882fac54b" +
    "1ccedfb0425434b54999a88cd7d993e05411955955c0cfec9dd33066605bd4a6" +
    "0f6bbfbab37349aaa762c23281b5749932c514f3b8723cf9bb05f9841a7f2d0e" +
    "0f75e42fd6c8e9f0edadac3dcfb7416c2d4b2470f4210f2afa93138615b1deb1" +
    "1ff56a9538b079e16dd77a8ef81318497b195ad81b8cd1c5ea5d48b0c160f599" +
    "12387b5ab69538ef4cda0f7a879982f9b4943291b1e6d998abefe7bb4ebb6993";

  const sig =
    "0x" +
    "534ebe41176f66ffaaa3dd387e097b998b51576ed0a8fbc3f9c8a1b14699adb2" +
    "0dc69ad5e32d78ccd4a964af065cfa76c23b6009ae877c1d748494776abfae6f00";
  return { sig, pClaims };
};

export const generateSigAndPClaimsDifferentHeight = () => {
  const pClaims =
    "0x" +
    "0000000000000200" +
    "04000000" +
    "02000400" +
    "58000000" +
    "00000200" +
    "01000000" +
    "03000000" + // different height
    "01000000" +
    "00000000" +
    "0d000000" +
    "02010000" +
    "19000000" +
    "02010000" +
    "25000000" +
    "02010000" +
    "31000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "de8b68a6643fa528a513f99a1ea30379927197a097ca86d9108e4c29d684b1ec" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "04000000" +
    "02000100" +
    "1d000000" +
    "02060000" +
    "01000000" +
    "03000000" + // different height
    "01000000" +
    "00000000" +
    "01000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "258aa89365a642358d92db67a13cb25d73e6eedf0d25100d8d91566882fac54b" +
    "1ccedfb0425434b54999a88cd7d993e05411955955c0cfec9dd33066605bd4a6" +
    "0f6bbfbab37349aaa762c23281b5749932c514f3b8723cf9bb05f9841a7f2d0e" +
    "0f75e42fd6c8e9f0edadac3dcfb7416c2d4b2470f4210f2afa93138615b1deb1" +
    "1ff56a9538b079e16dd77a8ef81318497b195ad81b8cd1c5ea5d48b0c160f599" +
    "12387b5ab69538ef4cda0f7a879982f9b4943291b1e6d998abefe7bb4ebb6993";
  const sig =
    "0x" +
    "b9c42c4e41a2df9040f061756dcd6ca47ec2042cbd9f22519309a553be6a1dcc" +
    "69ab71d5d8b002b9d8814fb1fee17e071fe37cc052cced7f8a57354e9bb8956a01";
  return { sig, pClaims };
};

export const generateSigAndPClaimsDifferentRound = () => {
  const pClaims =
    "0x" +
    "0000000000000200" +
    "04000000" +
    "02000400" +
    "58000000" +
    "00000200" +
    "01000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "0d000000" +
    "02010000" +
    "19000000" +
    "02010000" +
    "25000000" +
    "02010000" +
    "31000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "de8b68a6643fa528a513f99a1ea30379927197a097ca86d9108e4c29d684b1ec" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "04000000" +
    "02000100" +
    "1d000000" +
    "02060000" +
    "01000000" +
    "02000000" +
    "02000000" + // different round
    "00000000" +
    "01000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "258aa89365a642358d92db67a13cb25d73e6eedf0d25100d8d91566882fac54b" +
    "1ccedfb0425434b54999a88cd7d993e05411955955c0cfec9dd33066605bd4a6" +
    "0f6bbfbab37349aaa762c23281b5749932c514f3b8723cf9bb05f9841a7f2d0e" +
    "0f75e42fd6c8e9f0edadac3dcfb7416c2d4b2470f4210f2afa93138615b1deb1" +
    "1ff56a9538b079e16dd77a8ef81318497b195ad81b8cd1c5ea5d48b0c160f599" +
    "12387b5ab69538ef4cda0f7a879982f9b4943291b1e6d998abefe7bb4ebb6993";
  const sig =
    "0x" +
    "2f8df0cad1b95c978d0feb132f777459c2efc8204ec79bd6d17d31e4d50611c8" +
    "7c5a37d4f4a6838c9427b53749fda98970227dbd7364713abd24c2a14f5a36a601";
  return { sig, pClaims };
};

export const generateSigAndPClaimsDifferentChainId = () => {
  const pClaims =
    "0x" +
    "0000000000000200" +
    "04000000" +
    "02000400" +
    "58000000" +
    "00000200" +
    "0b000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "0d000000" +
    "02010000" +
    "19000000" +
    "02010000" +
    "25000000" +
    "02010000" +
    "31000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "de8b68a6643fa528a513f99a1ea30379927197a097ca86d9108e4c29d684b1ec" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" +
    "04000000" +
    "02000100" +
    "1d000000" + // different chainId
    "02060000" +
    "0b000000" +
    "02000000" +
    "01000000" +
    "00000000" +
    "01000000" +
    "02010000" +
    "f75f3eb17cd8136aeb15cca22b01ad5b45c795cb78787e74e55e088a7aa5fa16" +
    "258aa89365a642358d92db67a13cb25d73e6eedf0d25100d8d91566882fac54b" +
    "1ccedfb0425434b54999a88cd7d993e05411955955c0cfec9dd33066605bd4a6" +
    "0f6bbfbab37349aaa762c23281b5749932c514f3b8723cf9bb05f9841a7f2d0e" +
    "0f75e42fd6c8e9f0edadac3dcfb7416c2d4b2470f4210f2afa93138615b1deb1" +
    "1ff56a9538b079e16dd77a8ef81318497b195ad81b8cd1c5ea5d48b0c160f599" +
    "12387b5ab69538ef4cda0f7a879982f9b4943291b1e6d998abefe7bb4ebb6993";

  const sig =
    "0x" +
    "4343b877b15485ba84c8907cf982957f138466e2907572235d0ffe848b57c031" +
    "577e854335594ef6c99efe01e18f3fa40093b4741a3fc67080da7c1ae8111b7a01";

  return { sig, pClaims };
};

// Helper functions to create validators
export const generateMadID = (id: Number) => {
  return [id, id];
};

export const addValidators = async (
  validatorPool: ValidatorPoolMock | ValidatorPool,
  validators: string[]
) => {
  if ((<ValidatorPoolMock>validatorPool).isMock) {
    for (const validator of validators) {
      expect(await validatorPool.isValidator(validator)).to.equal(false);
      await validatorPool.registerValidators([validator], [0]);
      expect(await validatorPool.isValidator(validator)).to.equal(true);
    }
  }
};
