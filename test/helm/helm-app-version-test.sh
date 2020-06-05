#!/bin/bash

# Script to test that all helm charts' Chart.yaml and values.yaml are updated to reflect the latest release
# Exit with error when values.yaml is not updated with the current release version

set -euo pipefail

REPOPATH="$( cd "$(dirname "$0")"; cd ../../; pwd -P )"
MAKEFILEPATH=$REPOPATH/Makefile
RELEASE_VERSION=$(make -s -f $MAKEFILEPATH LATEST_TAG)
HELM_CHART_PATH=$REPOPATH/helm

for c in $HELM_CHART_PATH/*; do
    chart_name=$(echo $c | tr '/' '\n' | tail -1)
    if [ -d $c ]; then

        # verify AEMM version in Chart/yaml
        # note: appVersion in Chart.yaml is informational only
        app_version=$(sed -n 's/appVersion: //p' $c/Chart.yaml)
        if [[ $app_version != $RELEASE_VERSION ]]; then
            echo "Please update the appVersion in $chart_name helm chart's Chart.yaml and increment the chart version, if not already done. Expected version $RELEASE_VERSION, but got $app_version"
        fi

        # verify AEMM version in values.yaml in order to test the latest release
        values_image_tag=$(sed -n 's/[^\s]tag: //p' $c/values.yaml | tr -d '""' | tr -d ' ')
        if [[ $values_image_tag != $RELEASE_VERSION ]]; then
            echo "Please update the release version in $chart_name helm chart's values.yaml, image.tag field and increment the chart version, if not already done. Expected version $RELEASE_VERSION, but got $values_image_tag"
            exit 1
        fi
    fi
done
