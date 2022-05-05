#!/bin/sh

set -x

CURRENT_WD=$PWD
BRIDGE_DIR=./bridge

cd $BRIDGE_DIR || exit

npx hardhat node --show-stack-traces &
HARDHAT_NODE_PID="$!"

#trap "echo 'Intercepted SIGTERM hardhat.sh - $$ - $HARDHAT_NODE_PID' && kill -9 $HARDHAT_NODE_PID" SIGTERM SIGINT SIGKILL EXIT
echo "Intercepted SIGTERM main.sh - $$ - $HARDHAT_NODE_PID"
trap "trap - SIGTERM && kill -- $HARDHAT_NODE_PID" SIGTERM SIGINT SIGKILL EXIT

wait

cd "$CURRENT_WD" || exit

