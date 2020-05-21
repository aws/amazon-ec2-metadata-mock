#!/bin/bash
set -euo pipefail

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
BUILD_BIN="$SCRIPTPATH/../../build/bin"

BINARY_NAME="ec2-metadata-mock-linux-amd64"
LICENSE_TEST_TAG="aemm-license-test"

SUPPORTED_PLATFORMS="linux/amd64" make -s -f $SCRIPTPATH/../../Makefile build-binaries
docker build --build-arg=GOPROXY=direct -t $LICENSE_TEST_TAG $SCRIPTPATH/
docker run -it -e GITHUB_TOKEN --rm -v $SCRIPTPATH/:/test -v $BUILD_BIN/:/aemm-bin $LICENSE_TEST_TAG golicense /test/license-config.hcl /aemm-bin/$BINARY_NAME
