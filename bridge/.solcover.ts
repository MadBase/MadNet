module.exports = {
  skipFiles: ["libraries", "interfaces", "utils"],
  measureFunctionCoverage: false,
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
