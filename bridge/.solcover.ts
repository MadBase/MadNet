import {preFixtureSetup} from "./test/setup";
module.exports = {
  configureYulOptimizer: true,
  solcOptimizerDetails: {
    peephole: true,
    inliner: true,
    jumpdestRemover: true,
    orderLiterals: false, // <-- TRUE! Stack too deep when false
    deduplicate: false,
    cse: true,
    constantOptimizer: true,
    yul: true,
  },
  
};
