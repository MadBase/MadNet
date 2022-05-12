#!/bin/sh


CURRENT_WD=$PWD
BRIDGE_DIR=./bridge

cd $BRIDGE_DIR

npx hardhat create-local-test-genesis-node


cd $CURRENT_WD
