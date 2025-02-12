VERSION ?= $(shell git describe --tags --always --dirty)
LATEST_RELEASE_TAG=$(shell git describe --tags --abbrev=0)
PREVIOUS_RELEASE_TAG=$(shell git describe --abbrev=0 --tags `git rev-list --tags --skip=1  --max-count=1`)
LATEST_COMMIT_HASH=$(shell git rev-parse HEAD)
LATEST_COMMIT_CHART_VERSION=$(shell git --no-pager show ${LATEST_COMMIT_HASH}:helm/amazon-ec2-metadata-mock/Chart.yaml | grep 'version:' | cut -d' ' -f2 | tr -d '[:space:]')
REPO_FULL_NAME=aws/amazon-ec2-metadata-mock
ECR_REGISTRY ?= public.ecr.aws/aws-ec2
ECR_REPO ?= ${ECR_REGISTRY}/amazon-ec2-metadata-mock
ECR_REPO_CHART ?= amazon-ec2-metadata-mock
IMG ?= amazon/amazon-ec2-metadata-mock
IMG_TAG ?= ${VERSION}
IMG_W_TAG = ${IMG}:${IMG_TAG}
GOOS ?= linux
GOARCH ?= amd64
GOPROXY ?= "https://proxy.golang.org,direct"
SUPPORTED_PLATFORMS_LINUX ?= "linux/amd64,linux/arm64,linux/arm,darwin/amd64,darwin/arm64"
SUPPORTED_PLATFORMS_WINDOWS ?= "windows/amd64"
MAKEFILE_PATH = $(dir $(realpath -s $(firstword $(MAKEFILE_LIST))))
BUILD_DIR_PATH = ${MAKEFILE_PATH}/build
BIN_DIR = ${MAKEFILE_PATH}/bin
BINARY_NAME ?= ec2-metadata-mock
THIRD_PARTY_LICENSES = ${MAKEFILE_PATH}/THIRD_PARTY_LICENSES.md
GOLICENSES = ${BIN_DIR}/go-licenses
AMAZON_ECR_CREDENTIAL_HELPER_VERSION = 0.7.1

$(shell mkdir -p ${BUILD_DIR_PATH} && touch ${BUILD_DIR_PATH}/_go.mod)

$(GOLICENSES):
	GOBIN="$(BIN_DIR)" go install github.com/google/go-licenses@v1.6.0

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*$$' $(MAKEFILE_LIST) | sort

version:
	@echo ${VERSION}

chart-version:
	@echo ${LATEST_COMMIT_CHART_VERSION}

latest-release-tag:
	@echo ${LATEST_RELEASE_TAG}

previous-release-tag:
	@echo ${PREVIOUS_RELEASE_TAG}

repo-full-name:
	@echo ${REPO_FULL_NAME}

binary-name:
	@echo ${BINARY_NAME}

image:
	@echo ${IMG_W_TAG}

clean:
	rm -rf ${BUILD_DIR_PATH}

compile:
	@echo ${MAKEFILE_PATH}
	go build -a -tags aemm${GOOS} -o ${BUILD_DIR_PATH}/${BINARY_NAME} ./cmd/...

validate-json:
	${MAKEFILE_PATH}/scripts/validators/json-validator

build: validate-json compile

unit-test:
	go test -bench=. ${MAKEFILE_PATH}/... -v -coverprofile=coverage.out -covermode=atomic -outputdir=${BUILD_DIR_PATH}

e2e-test: build
	${MAKEFILE_PATH}/test/e2e/run-tests

helm-lint-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -l

helm-install-e2e-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -i

helm-mock-ip-count-test:
	${MAKEFILE_PATH}/test/helm/chart-test.sh -m

shellcheck:
	${MAKEFILE_PATH}/test/shellcheck/run-shellcheck

.PHONY: third-party-licenses
third-party-licenses: $(GOLICENSES)
	@$(GOLICENSES) report \
		--include_tests \
		--template "${MAKEFILE_PATH}/templates/third-party-licenses.tmpl" \
		"${MAKEFILE_PATH}/..." > "${THIRD_PARTY_LICENSES}"

.PHONY: license-test
license-test: $(GOLICENSES)
	@$(GOLICENSES) check \
		--allowed_licenses="Apache-2.0,BSD-2-Clause,BSD-3-Clause,BSD-4-Clause,ISC,MIT,MPL-2.0" \
		--include_tests \
		"${MAKEFILE_PATH}/..." \
		&& echo "✅ Passed" || echo "❌ Failed"

spellcheck:
	${MAKEFILE_PATH}/test/readme-test/run-readme-spellcheck

test: spellcheck shellcheck unit-test e2e-test helm-install-e2e-test license-test

build-binaries:
	${MAKEFILE_PATH}/scripts/build-binaries -d -p ${SUPPORTED_PLATFORMS_LINUX},${SUPPORTED_PLATFORMS_WINDOWS} -v ${VERSION}

generate-k8s-yaml:
	${MAKEFILE_PATH}/scripts/generate-k8s-yaml

gen-helm-chart-archives:
	${MAKEFILE_PATH}/scripts/generate-helm-chart-archives

upload-resources-to-github:
	${MAKEFILE_PATH}/scripts/upload-resources-to-github

build-release-assets: build-binaries generate-k8s-yaml gen-helm-chart-archives

build-docker-images-linux:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${SUPPORTED_PLATFORMS_LINUX} -r ${IMG} -v ${VERSION}

build-docker-images-windows:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${SUPPORTED_PLATFORMS_WINDOWS} -r ${IMG} -v ${VERSION}

push-docker-images-linux:
	${MAKEFILE_PATH}/scripts/retag-docker-images -p ${SUPPORTED_PLATFORMS_LINUX} -v ${VERSION} -o ${IMG} -n ${ECR_REPO}
	@ECR_REGISTRY=${ECR_REGISTRY} ${MAKEFILE_PATH}/scripts/ecr-public-login
	${MAKEFILE_PATH}/scripts/push-docker-images -p ${SUPPORTED_PLATFORMS_LINUX} -r ${ECR_REPO} -v ${VERSION} -m

push-docker-images-windows:
	${MAKEFILE_PATH}/scripts/retag-docker-images -p ${SUPPORTED_PLATFORMS_WINDOWS} -v ${VERSION} -o ${IMG} -n ${ECR_REPO}
	bash ${MAKEFILE_PATH}/scripts/install-amazon-ecr-credential-helper $(AMAZON_ECR_CREDENTIAL_HELPER_VERSION)
	${MAKEFILE_PATH}/scripts/push-docker-images -p ${SUPPORTED_PLATFORMS_WINDOWS} -r ${ECR_REPO} -v ${VERSION} -m

push-helm-chart:
	@ECR_REGISTRY=${ECR_REGISTRY} ${MAKEFILE_PATH}/scripts/helm-login
	${MAKEFILE_PATH}/scripts/push-helm-chart -r ${ECR_REPO_CHART} -v ${LATEST_COMMIT_CHART_VERSION} -h ${ECR_REGISTRY}

sync-readme-to-ecr-public:
	@ECR_REGISTRY=${ECR_REGISTRY} ${MAKEFILE_PATH}/scripts/ecr-public-login
	${MAKEFILE_PATH}/scripts/sync-readme-to-ecr-public

sync-catalog-information-for-helm-chart:
	@ECR_REGISTRY=${ECR_REGISTRY} ${MAKEFILE_PATH}/scripts/helm-login
	${MAKEFILE_PATH}/scripts/sync-catalog-information-for-helm-chart

homebrew-sync-dry-run:
	${MAKEFILE_PATH}/scripts/sync-to-aws-homebrew-tap -d -b ${BINARY_NAME} -r ${REPO_FULL_NAME} -p ${SUPPORTED_PLATFORMS_LINUX} -v ${LATEST_RELEASE_TAG}

homebrew-sync:
	${MAKEFILE_PATH}/scripts/sync-to-aws-homebrew-tap -b ${BINARY_NAME} -r ${REPO_FULL_NAME} -p ${SUPPORTED_PLATFORMS_LINUX}

validate-release-version:
	${MAKEFILE_PATH}/scripts/validators/release-version-validator

release-github: build-release-assets upload-resources-to-github

release-docker-linux: build-docker-images-linux push-docker-images-linux sync-readme-to-ecr-public

release-docker-windows: build-docker-images-windows push-docker-images-windows

release: release-github release-docker-linux release-docker-windows

# Targets intended for local use
fmt:
	goimports -w ./ && gofmt -s -w ./

build-and-test: build test

docker-build:
	${MAKEFILE_PATH}/scripts/build-docker-images -d -p ${GOOS}/${GOARCH} -r ${IMG} -v ${VERSION}

docker-run:
	docker run ${IMG_W_TAG}

## Targets intended to be run in preparation for a new release
create-local-release-tag-major:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -m

create-local-release-tag-minor:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -i

create-local-release-tag-patch:
	${MAKEFILE_PATH}/scripts/create-local-tag-for-release -p

create-release-pr:
	${MAKEFILE_PATH}/scripts/prepare-for-release

create-release-pr-draft:
	${MAKEFILE_PATH}/scripts/prepare-for-release -d

release-prep-major: create-local-release-tag-major create-release-pr

release-prep-minor: create-local-release-tag-minor create-release-pr

release-prep-patch: create-local-release-tag-patch create-release-pr

release-prep-custom: # Run make NEW_VERSION=v1.2.3 release-prep-custom to prep for a custom release version
ifdef NEW_VERSION
	$(shell echo "${MAKEFILE_PATH}/scripts/create-local-tag-for-release -v $(NEW_VERSION) && echo && make create-release-pr")
endif

