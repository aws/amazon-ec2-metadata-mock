MAKEFILE_PATH = $(dir $(realpath -s $(firstword $(MAKEFILE_LIST))))
BUILD_DIR_PATH = ${MAKEFILE_PATH}/build
BINARY_NAME=amazon-ec2-metadata-mock
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
	go build -a -ldflags '-X "${DEFAULT_VALUES_VAR}=${ENCODED_METADATA_DEFAULTS}"' -o ${BUILD_DIR_PATH}/${BINARY_NAME} ${MAKEFILE_PATH}/cmd/amazon-ec2-metadata-mock.go

build: create-build-dir validate-json compile

unit-test: create-build-dir
	go test ${MAKEFILE_PATH}/... -v -coverprofile=coverage.txt -covermode=atomic -outputdir=${BUILD_DIR_PATH}

e2e-test: build
	${MAKEFILE_PATH}/test/e2e/run-tests

validate-json:
	${MAKEFILE_PATH}/test/json-validator

license-test:
	${MAKEFILE_PATH}/test/license-test/run-license-test.sh

go-report-card-test:
	${MAKEFILE_PATH}/test/go-report-card-test/run-report-card-test.sh

test: unit-test e2e-test license-test go-report-card-test