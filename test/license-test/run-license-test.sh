#!/bin/bash
set -euo pipefail

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
BUILD_BIN="$SCRIPTPATH/../../build"

BINARY_NAME="amazon-ec2-metadata-mock"
LICENSE_TEST_TAG="aemm-license-test"

make -s -f $SCRIPTPATH/../../Makefile build
docker build --build-arg=GOPROXY=direct -t $LICENSE_TEST_TAG $SCRIPTPATH/
docker run -it -e GITHUB_TOKEN --rm -v $SCRIPTPATH/:/test -v $BUILD_BIN/:/aemm-bin $LICENSE_TEST_TAG golicense /test/license-config.hcl /aemm-bin/$BINARY_NAME
