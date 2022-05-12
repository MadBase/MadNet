#!/bin/sh


CURRENT_WD=$PWD
BRIDGE_DIR=./bridge

cd $BRIDGE_DIR

#make the genesis node 
#aliceNetBlocknum = $(npx hardhat create-local-test-genesis-node)

testnetBlocknum=$(npx hardhat get-latest-blockheight)

npx hardhat node --fork https://eth-ropsten.alchemyapi.io/v2/4ynIWs9XdY4lQdv0xthFfqvV4qmumPTB --fork-block-number $testnetBlocknum --show-stack-traces &
sleep 10
#turn on impersonating
npx hardhat enable-hardhat-impersonate --account 0x137425E39a2A981ed83Fe490dedE1aB139840B87 --network dev
#mine 9000 blocks
npx hardhat mine-num-blocks --num-blocks 9000 --network dev
#pause validator at
#npx hardhat pause-consensus-at-height  --height $aliceNetBlocknum --network dev
#evict validators 
#npx hardhat unregister-all-validators --network dev
#npx hardhat start-local-genesis-node --network dev
wait
cd $CURRENT_WD
