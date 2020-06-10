#!/bin/bash

# Script to test helm charts for the following:
## YAML validity
## chart deployability
## chart version increment when chart is changed
## chart installation on Kind cluster
## run helm tests that make http requests to AEMM service running on the test cluster

# Prerequisites:
## Docker

set -euo pipefail

# KIND / Kubernetes
readonly K8s_1_18="v1.18.2"
readonly K8s_1_17="v1.17.5"
readonly K8s_1_16="v1.16.9"
readonly K8s_1_15="v1.15.11"
readonly K8s_1_14="v1.14.10"
KIND_IMAGE="$K8s_1_18"
readonly KIND_VERSION="v0.8.1"
readonly CLUSTER_NAME="kind-ct"
readonly REPO_PATH="$( cd "$(dirname "$0")"; cd ../../ ; pwd -P )"
readonly CLUSTER_CONFIG="$REPO_PATH/test/helm/kind-config.yaml"
readonly TMP_DIR="$REPO_PATH/build/tmp-$CLUSTER_NAME"
readonly KUBECONFIG_TMP_PATH="$TMP_DIR/kubeconfig"
readonly KIND_EXEC_ARGS="--context kind-$CLUSTER_NAME --kubeconfig $KUBECONFIG_TMP_PATH"

# Helm/chart-testing
CT_TAG="v3.0.0-rc.1"
readonly CT_CONFIG="test/helm/ct.yaml"
AEMM_DOCKER_IMG="amazon/amazon-ec2-metadata-mock:$(make latest-tag)"

# chart-testing container
readonly CT_CONTAINER_NAME="ct"
readonly CT_EXEC="docker exec --interactive $CT_CONTAINER_NAME"

readonly HELP=$(cat << 'EOM'
Test Helm charts by linting and installing helm charts in a provisioned environment. Only changed charts are tested.
Provisioned environment includes a Kind Kubernetes cluster and Docker container with helm/chart-testing image.

Usage:
  chart-test [options]

Examples:
  chart-test -l                       Run the test for linting only for default Kubernetes version
  chart-test -k 1.17                  Run the test for linting and installation for Kubernetes version 1.17
  chart-test -r -c v3.0.0-rc.1        Run the test for linting and installation for default Kubernetes version, and reuse previously provisioned test environment

Options:
  -k     Kubernetes version / kindest node image tag to use for the test (default: 1.18) (options: 1.14, 1.15, 1.16, 1.17)
  -c     chart-testing image tag to use for the test
  -l     test charts for linting only (helm lint, version checking, YAML validation, maintainer validation)
  -p     preserve the provisioned environment after test runs
  -r     reuse kind cluster and docker chart-testing container previously provisioned by this tool
  -d     debug, enables set -x, printing primary commands before executing
  -h     help message
EOM
)

LINT_ONLY=false
DEBUG=false
PRESERVE=false
REUSE_ENV=false

export TERM="xterm"
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
MAGENTA=$(tput setaf 5)
RESET_FMT=$(tput sgr 0)
BOLD=$(tput bold)

setup_ct_container() {
    c_echo "Provisioning and running chart-testing container named $CT_CONTAINER_NAME..."
    docker run --rm --interactive --detach --network host --name $CT_CONTAINER_NAME \
        --volume "$REPO_PATH/$CT_CONFIG:/etc/ct/ct.yaml" \
        --volume "$REPO_PATH:/workdir" \
        --workdir /workdir \
        "quay.io/helmpack/chart-testing:$CT_TAG"
    echo
}

install_kind() {
    c_echo "Installing kind..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/$KIND_VERSION/kind-$(uname)-amd64
    chmod +x ./kind
    mv ./kind $TMP_DIR/kind
    export PATH=$TMP_DIR:$PATH
}

create_kind_cluster() {
    c_echo "Creating kind Kubernetes cluster with kubeconfig in $KUBECONFIG_TMP_PATH"
    kind create cluster --name $CLUSTER_NAME --config $CLUSTER_CONFIG --image "kindest/node:$KIND_IMAGE" --kubeconfig $KUBECONFIG_TMP_PATH --wait 60s

    build_and_load_docker_image

    c_echo "Copying kubeconfig to container..."
    $CT_EXEC mkdir -p /root/.kube
    docker cp $KUBECONFIG_TMP_PATH ct:/root/.kube/config

    c_echo "Cluster ready!\n"
}

# build and load Docker image for AEMM to be able to install the latest
build_and_load_docker_image() {

    ## pre-pull golang builder to prevent intermittent timeouts from dockerhub
    timeout 120 docker pull golang:1.14-- || :

    c_echo "🥑 Building the amazon-ec2-metadata-mock docker image"
    docker build -t $AEMM_DOCKER_IMG "$REPO_PATH" 
    c_echo "👍 Built the amazon-ec2-metadata-mock docker image"
    c_echo "$AEMM_DOCKER_IMG" > $TMP_DIR/aemm-docker-img

    c_echo "🥑 Loading AEMM image into the cluster"
    kind load docker-image --name $CLUSTER_NAME --nodes=$CLUSTER_NAME-worker,$CLUSTER_NAME-control-plane $AEMM_DOCKER_IMG
    c_echo "👍 Loaded AEMM image into the cluster"
}

handle_errors_and_cleanup() {
    # cleanup
    echo
    MSG_PREFIX="-> "
    if [ $PRESERVE == true ]; then
        c_echo "The test environment is preserved. Reuse with the '-r' option.\n${MSG_PREFIX}List Docker container: docker ps --filter \"name=^ct$\""

        if [[ $LINT_ONLY == true ]]; then
            c_echo "Cleanup commands:\n  * docker kill ct > /dev/null 2>&1" $MSG_PREFIX
        else
            c_echo "List kind cluster: kind get clusters" $MSG_PREFIX
            c_echo "Cluster config can be found in $KUBECONFIG_TMP_PATH" $MSG_PREFIX
            c_echo "Cleanup commands:\n  * docker kill ct > /dev/null 2>&1\n  * kind delete cluster --name $CLUSTER_NAME\n  * rm -r $TMP_DIR" $MSG_PREFIX
            c_echo "Kubectl commands:\n  * kubectl get pods $KIND_EXEC_ARGS" $MSG_PREFIX
        fi
    else
        c_echo "Cleaning up resources..."
        c_echo "Deleting ct container..." $MSG_PREFIX
        docker kill ct > /dev/null 2>&1

         if [[ $LINT_ONLY == false ]]; then
            c_echo "Deleting kind cluster $CLUSTER_NAME..." $MSG_PREFIX
            kind delete cluster --name $CLUSTER_NAME --quiet

            c_echo "Deleting tmp dir '$TMP_DIR'" $MSG_PREFIX
            rm -r $TMP_DIR
        fi
    fi

    # error handling
    if [ $1 != "0" ]; then
        FAILED_COMMAND=${@:2}
        echo -e "\n❌ ${RED}One or more tests failed / error occurred while running command ${BOLD}${FAILED_COMMAND}${RESET_FMT} ❌"
        exit 1
    fi

    echo -e "\n✅✅ ${GREEN}All tests passed and cleaned up${RESET_FMT} ✅✅\n"
}

test_charts() {
    if [ $DEBUG == true ]; then
        set -x
    fi

    # provision test env
    if [[ $REUSE_ENV == false ]]; then
        # setup env to run chart-testing commands
        setup_ct_container

        if [[ $LINT_ONLY == false ]]; then
            # setup env for chart installation
            mkdir -p $TMP_DIR
            install_kind
            create_kind_cluster
        fi
    fi

    c_echo "Linting and validating helm charts"
    git remote add upstream https://github.com/aws/amazon-ec2-metadata-mock.git &> /dev/null || true
    git fetch upstream
    $CT_EXEC ct lint

    [[ $? == 0 ]] && echo -e "✅ ${GREEN}All charts linted successfully${RESET_FMT}"
    c_echo "------------------------------------------------------------------------------------------------------------------------"

    if [[ $LINT_ONLY == false ]]; then
        c_echo "Installing helm charts and running tests...\n"

        git remote add upstream https://github.com/aws/amazon-ec2-metadata-mock.git &> /dev/null || true
        git fetch upstream

        if [[ $DEBUG == true ]]; then
            $CT_EXEC ct install --debug
        else
            $CT_EXEC ct install
        fi
        [[ $? == 0 ]] && echo -e "✅ ${GREEN}All charts installed and tested successfully${RESET_FMT}"
    fi

    if [ $DEBUG == true ]; then
        set +x
    fi
}

# $1=message to echo; [$2]=indication of sub-echo
c_echo() {
    DEFAULT_PREFIX="🥑"
    PREFIX="${2:-$DEFAULT_PREFIX}"
    echo -e "${MAGENTA}${PREFIX} ${1}${RESET_FMT}"
}

process_args() {
    while getopts "hdlprk:c:" opt; do
        case ${opt} in
            h )
              echo -e "$HELP" 1>&2
              exit 0
              ;;
            d )
              DEBUG=true
              ;;
            l )
              LINT_ONLY=true
              ;;
            p )
              PRESERVE=true
              ;;
            r )
              REUSE_ENV=true
              ;;
            k )
              OPTARG="K8s_$(echo $OPTARG | sed 's/\./\_/g')"
              KIND_IMAGE="${!OPTARG}"
              ;;
            c )
              CT_TAG=$OPTARG
              ;;
            \? )
              echo "$HELP" 1>&2
              exit 0
              ;;
        esac
    done
    shift $((OPTIND -1))
}

main() {
    process_args $@

    trap 'handle_errors_and_cleanup $? $BASH_COMMAND' EXIT

    c_echo "Testing Helm charts in a newly provisioned test environment"
    if [[ $LINT_ONLY == true ]]; then
        c_echo "Using:\n${BOLD}  * helm/chart-testing version=$CT_TAG\n  * lint only=$LINT_ONLY\n  * preserve test env=$PRESERVE\n  * reuse=$REUSE_ENV\n  * debug=$DEBUG\n${RESET_FMT}"
    else
        c_echo "Using:\n${BOLD}  * kind version=$KIND_VERSION\n  * Kubernetes version=$KIND_IMAGE\n  * helm/chart-testing version=$CT_TAG\n  * lint only=$LINT_ONLY\n  * preserve test env=$PRESERVE\n  * reuse=$REUSE_ENV\n  * debug=$DEBUG\n${RESET_FMT}"
    fi

    test_charts
}

main $@
