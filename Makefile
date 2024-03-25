TEST?=$$(go list ./... | grep -v 'vendor')
###### chang variables below according to your own modules ###
NAMESPACE=kusionstack
NAME=kawesome
VERSION=v0.1.0
BINARY=../bin/kusion-module-${NAME}_${VERSION}

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
GOARCH_LOCAL := amd64
else
GOARCH_LOCAL := $(LOCAL_ARCH)
endif
export GOOS_LOCAL := $(shell uname|tr 'A-Z' 'a-z')
export OS_ARCH ?= $(GOARCH_LOCAL)

default: install

build-darwin:
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY} ./${NAME}

install: build-darwin
# copy module binary to $KUSION_HOME. e.g. ~/.kusion/modules/kusionstack/kawesome/v0.1.0/darwin/arm64/kusion-module-kawesome_v0.1.0
	mkdir -p ${KUSION_HOME}/modules/${NAMESPACE}/${NAME}/${VERSION}/${GOOS_LOCAL}/${OS_ARCH}
	cp ${BINARY} ${KUSION_HOME}/modules/${NAMESPACE}/${NAME}/${VERSION}/${GOOS_LOCAL}/${OS_ARCH}

test:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 5m
