VERSION = $(shell git describe --tags --always --dirty)
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

create-build-dir:
	mkdir -p ${BUILD_DIR_PATH}

clean:
	rm -rf ${BUILD_DIR_PATH}

fmt:
	goimports -w ./

compile:
	@echo ${MAKEFILE_PATH}
	go build -a -tags aemm${GOOS} -ldflags '-X "${DEFAULT_VALUES_VAR}=${ENCODED_METADATA_DEFAULTS}"' -o ${BUILD_DIR_PATH}/${BINARY_NAME} ${MAKEFILE_PATH}/cmd/amazon-ec2-metadata-mock.go

build: create-build-dir validate-json compile

build-and-test: build test

build-binaries:
	${MAKEFILE_PATH}/scripts/build-binaries -d -p ${SUPPORTED_PLATFORMS} -v ${VERSION}

docker-build:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${GOOS}/${GOARCH} -r ${IMG} -v ${VERSION}

docker-run:
	docker run ${IMG_W_TAG}

docker-push:
	@echo ${DOCKERHUB_TOKEN} | docker login -u ${DOCKERHUB_USERNAME} --password-stdin
	docker push ${IMG_W_TAG}

build-docker-images:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${SUPPORTED_PLATFORMS} -r ${IMG} -v ${VERSION}

push-docker-images:
	@echo ${DOCKERHUB_TOKEN} | docker login -u ${DOCKERHUB_USERNAME} --password-stdin
	${MAKEFILE_PATH}/scripts/push-docker-images -p ${SUPPORTED_PLATFORMS} -r ${IMG} -v ${VERSION} -m

version:
	@echo ${VERSION}

image:
	@echo ${IMG_W_TAG}

upload-resources-to-github:
	${MAKEFILE_PATH}/scripts/upload-resources-to-github

generate-k8s-yaml:
	${MAKEFILE_PATH}/scripts/generate-k8s-yaml

sync-readme-to-dockerhub:
	${MAKEFILE_PATH}/scripts/sync-readme-to-dockerhub

unit-test: create-build-dir
	go test ${MAKEFILE_PATH}/... -v -coverprofile=coverage.txt -covermode=atomic -outputdir=${BUILD_DIR_PATH}

e2e-test: build
	${MAKEFILE_PATH}/test/e2e/run-tests

validate-json:
	${MAKEFILE_PATH}/test/json-validator

validate-readme:
	${MAKEFILE_PATH}/test/readme-validator

helm-lint-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -l

helm-e2e-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh

helm-app-version-test:
	${MAKEFILE_PATH}/test/helm/helm-app-version-test.sh

helm-tests:
	helm-e2e-test helm-app-version-test

gen-helm-chart-archives:
	${MAKEFILE_PATH}/scripts/generate-helm-chart-archives

license-test:
	${MAKEFILE_PATH}/test/license-test/run-license-test.sh

go-report-card-test:
	${MAKEFILE_PATH}/test/go-report-card-test/run-report-card-test.sh

test: unit-test e2e-test helm-e2e-test helm-app-version-test license-test go-report-card-test

release: create-build-dir build-binaries build-docker-images push-docker-images generate-k8s-yaml gen-helm-chart-archives upload-resources-to-github
