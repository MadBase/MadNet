module.exports = {
  configureYulOptimizer: true,
  solcOptimizer: {
    enabled: true,
    runs: 20000,
    solcOptimizerDetails: {
      peephole: false,
      inliner: true,
      jumpdestRemover: true,
      orderLiterals: true, // <-- TRUE! Stack too deep when false
      deduplicate: true,
      cse: false,
      constantOptimizer: false,
      yul: true,
    },
  },
};
