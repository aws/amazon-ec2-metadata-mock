VERSION ?= $(shell git describe --tags --always --dirty)
LATEST_RELEASE_TAG=$(shell git tag | tail -1)
PREVIOUS_RELEASE_TAG=$(shell git tag | tail -2 | head -1)
REPO_FULL_NAME=aws/amazon-ec2-metadata-mock
IMG ?= amazon/amazon-ec2-metadata-mock
IMG_TAG ?= ${VERSION}
IMG_W_TAG = ${IMG}:${IMG_TAG}
DOCKERHUB_USERNAME ?= ""
DOCKERHUB_TOKEN ?= ""
GOOS ?= linux
GOARCH ?= amd64
GOPROXY ?= "https://proxy.golang.org,direct"
SUPPORTED_PLATFORMS ?= "linux/amd64,linux/arm64,linux/arm,darwin/amd64,windows/amd64"
MAKEFILE_PATH = $(dir $(realpath -s $(firstword $(MAKEFILE_LIST))))
BUILD_DIR_PATH = ${MAKEFILE_PATH}/build
BINARY_NAME=ec2-metadata-mock
METADATA_DEFAULTS_FILE=${MAKEFILE_PATH}/pkg/config/defaults/aemm-metadata-default-values.json
ENCODED_METADATA_DEFAULTS=$(shell cat ${METADATA_DEFAULTS_FILE} | base64 | tr -d \\n)
DEFAULT_VALUES_VAR=github.com/aws/amazon-ec2-metadata-mock/pkg/config/defaults.encodedDefaultValues
ROOT_VERSION_VAR=github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/root.version

$(shell mkdir -p ${BUILD_DIR_PATH} && touch ${BUILD_DIR_PATH}/_go.mod)

help:
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | sort

version:
	@echo ${VERSION}

latest-release-tag:
	@echo ${LATEST_RELEASE_TAG}

previous-release-tag:
	@echo ${PREVIOUS_RELEASE_TAG}

repo-full-name:
	@echo ${REPO_FULL_NAME}

image:
	@echo ${IMG_W_TAG}

clean:
	rm -rf ${BUILD_DIR_PATH}

compile:
	@echo ${MAKEFILE_PATH}
	go build -a -tags aemm${GOOS} -ldflags '-X "${DEFAULT_VALUES_VAR}=${ENCODED_METADATA_DEFAULTS}" -X "${ROOT_VERSION_VAR}=${VERSION}"' -o ${BUILD_DIR_PATH}/${BINARY_NAME} ${MAKEFILE_PATH}/cmd/amazon-ec2-metadata-mock.go

validate-json:
	${MAKEFILE_PATH}/scripts/validators/json-validator

build: validate-json compile

unit-test:
	go test -bench=. ${MAKEFILE_PATH}/... -v -coverprofile=coverage.out -covermode=atomic -outputdir=${BUILD_DIR_PATH}

validate-readme:
	${MAKEFILE_PATH}/scripts/validators/readme-validator

e2e-test: build
	${MAKEFILE_PATH}/test/e2e/run-tests

helm-lint-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -l

helm-install-e2e-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -i

license-test:
	${MAKEFILE_PATH}/test/license-test/run-license-test.sh

go-report-card-test:
	${MAKEFILE_PATH}/test/go-report-card-test/run-report-card-test.sh

test: unit-test e2e-test helm-install-e2e-test license-test go-report-card-test

build-binaries:
	${MAKEFILE_PATH}/scripts/build-binaries -d -p ${SUPPORTED_PLATFORMS} -v ${VERSION}

generate-k8s-yaml:
	${MAKEFILE_PATH}/scripts/generate-k8s-yaml

gen-helm-chart-archives:
	${MAKEFILE_PATH}/scripts/generate-helm-chart-archives

upload-resources-to-github:
	${MAKEFILE_PATH}/scripts/upload-resources-to-github

build-release-assets: build-binaries generate-k8s-yaml gen-helm-chart-archives

build-docker-images:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${SUPPORTED_PLATFORMS} -r ${IMG} -v ${VERSION}

push-docker-images:
	@echo ${DOCKERHUB_TOKEN} | docker login -u ${DOCKERHUB_USERNAME} --password-stdin
	${MAKEFILE_PATH}/scripts/push-docker-images -p ${SUPPORTED_PLATFORMS} -r ${IMG} -v ${VERSION} -m

sync-readme-to-dockerhub:
	${MAKEFILE_PATH}/scripts/sync-readme-to-dockerhub

validate-release-version:
	${MAKEFILE_PATH}/scripts/validators/release-version-validator

release-github: build-release-assets upload-resources-to-github

release-docker: build-docker-images push-docker-images sync-readme-to-dockerhub

release: release-github release-docker

# Targets intended for local use 
fmt:
	goimports -w ./ && gofmt -s -w ./

build-and-test: build test

docker-build:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${GOOS}/${GOARCH} -r ${IMG} -v ${VERSION}

docker-run:
	docker run ${IMG_W_TAG}

docker-push:
	@echo ${DOCKERHUB_TOKEN} | docker login -u ${DOCKERHUB_USERNAME} --password-stdin
	docker push ${IMG_W_TAG}

## Targets intended to be run in preparation for a new release
create-local-tag-for-major-release:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -m

create-local-tag-for-minor-release:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -i

create-local-tag-for-patch-release:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -p

create-release-prep-pr:
	${MAKEFILE_PATH}/scripts/prepare-for-release

release-prep-major: create-local-tag-for-major-release create-release-prep-pr

release-prep-minor: create-local-tag-for-minor-release create-release-prep-pr

release-prep-patch: create-local-tag-for-patch-release create-release-prep-pr

release-prep-custom: # Run make NEW_VERSION=1.2.3 release-prep-custom to prep for a custom release version
ifdef NEW_VERSION
	$(shell echo "${MAKEFILE_PATH}/scripts/create-local-tag-for-release -v $(NEW_VERSION) && echo && make create-release-prep-pr")
endif
	