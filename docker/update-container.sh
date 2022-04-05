#!/bin/bash

# This script rebuilds an image and if the image has changed, recreates its container
# The image and container are linked by having the same tag
# The script also ensures the container is stopped, so it can always be started right after

set -euo pipefail
export MSYS_NO_PATHCONV=1;

YELLOW='\033[0;33m';
NOCOL='\033[0m';

DOCKERFILE=$1;
TAG=$2;
DOCKER_CREATE_ARGS=${3-};
CMD=${4-};

DOCKER_BUILD_ARGS="-f $DOCKERFILE -t $TAG";

PASS_PERMVARS=${PASS_PERMVARS-};
if [ "$PASS_PERMVARS" != "" ]; then
	source ./docker/set-permvars.sh;
	DOCKER_BUILD_ARGS="$DOCKER_BUILD_ARGS --build-arg BUILDER_UID --build-arg BUILDER_GIDS";
fi;

echo -e "$YELLOW# Building image... $NOCOL";
DOCKER_BUILDKIT=1 docker build . $DOCKER_BUILD_ARGS;

EXISTING_CONTAINER_IMAGE=$(docker ps -a --filter name=$TAG --format {{.Image}});
if [ "$EXISTING_CONTAINER_IMAGE" = "$TAG" ]; then
	IS_RUNNING=$(docker ps --filter name=$TAG --format {{.Image}});
	if [ "$IS_RUNNING" != "" ]; then
		echo -e "$YELLOW# Stopping existing container... $NOCOL";
		docker stop -t 0 $TAG;
	fi;
	
	echo -e "$YELLOW# Container ready! $NOCOL";
else
	if [ "$EXISTING_CONTAINER_IMAGE" != "" ]; then
		echo -e "$YELLOW# Removing old image and container... $NOCOL";
		docker container rm -vf $TAG;
		docker image rm -f $EXISTING_CONTAINER_IMAGE
		docker builder prune -f
	fi;
	
	echo -e "$YELLOW# Creating new container... $NOCOL";
	docker create --name $TAG -it $DOCKER_CREATE_ARGS $TAG $CMD;
	echo -e "$YELLOW# Container created! $NOCOL";
fi;