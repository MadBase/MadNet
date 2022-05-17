#!/bin/sh

set -x

CURRENT_WD=$PWD
BRIDGE_DIR=./bridge

cd $BRIDGE_DIR

npx hardhat setHardhatIntervalMining --network dev --interval 1000
npx hardhat --network dev --show-stack-traces registerValidators --factory-address "$@"
npx hardhat setHardhatIntervalMining --network dev --enable-auto-mine

cd $CURRENT_WD
