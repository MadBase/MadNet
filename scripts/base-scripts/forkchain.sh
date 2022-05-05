#!/bin/sh


CURRENT_WD=$PWD
BRIDGE_DIR=./bridge

cd $BRIDGE_DIR

npx hardhat node --fork https://testnet.eth.mnexplore.com/ --fork-block-number  --show-stack-traces &


wait

cd $CURRENT_WD
