GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
SSM_CONNECT_VERSION := "0.2.3"

GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?=$(shell arch)

.PHONY: all build install

all: build install

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build OS ARCH
build: guard-SSM_CONNECT_VERSION mod-tidy clean
	@echo "================================================="
	@echo "Building ssm-connect"
	@echo "=================================================\n"

	@if [ ! -d "${GOOS}" ]; then \
		mkdir "${GOOS}"; \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${GOOS}/ssm-connect"
	sleep 2
	tar -C "${GOOS}" -czvf "ssmconnect_${SSM_CONNECT_VERSION}_${GOOS}_${GOARCH}.tgz" ssm-connect; \

.PHONY: clean
clean:
	@echo "================================================="
	@echo "Cleaning ssm-connect"
	@echo "=================================================\n"
	@for OS in darwin linux; do \
		if [ -f $${OS}/ssm-connect ]; then \
			rm -f $${OS}/ssm-connect; \
		fi; \
	done

.PHONY: clean-all
clean-all: clean
	@echo "================================================="
	@echo "Cleaning tarballs"
	@echo "=================================================\n"
	@rm -f *.tgz 2>/dev/null

.PHONY: install
install:
	@echo "================================================="
	@echo "Installing ssm-connect in ${GOPATH}/bin"
	@echo "=================================================\n"

	go install -race

#
# General targets
#
guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
