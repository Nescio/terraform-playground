.DEFAULT_GOAL := build
NAME          := petstore
OS            := $(shell go env GOOS)
ARCH          := $(shell go env GOARCH)
PLUGIN_PATH   := ${HOME}/.terraform.d/plugins/hacknights.club/${NAME}-provider/${NAME}/1.0.0/${OS}_${ARCH}
PLUGIN_NAME   := terraform-provider-${NAME}
DIST_PATH     := dist/${OS}_${ARCH}
GO_PACKAGES   := $(shell go list ./... | grep -v /vendor/)
GO_FILES      := $(shell find . -type f -name '*.go')


.PHONY: all
all: clean test build

.PHONY: test
test: test-all

.PHONY: test-all
test-all:
	@TF_LOG=TRACE PETSTORE_ADDRESS=http://localhost:8080 TF_ACC=1 go test -v -race $(GO_PACKAGES)

${DIST_PATH}/${PLUGIN_NAME}: ${GO_FILES}
	mkdir -p $(DIST_PATH); \
	go build -o $(DIST_PATH)/${PLUGIN_NAME}

.PHONY: build
build: ${DIST_PATH}/${PLUGIN_NAME}

.PHONY: install
install: build
	mkdir -p $(PLUGIN_PATH); \
	rm -rf $(PLUGIN_PATH)/${PLUGIN_NAME}; \
	install -m 0755 $(DIST_PATH)/${PLUGIN_NAME} $(PLUGIN_PATH)/${PLUGIN_NAME} 

.PHONY: clean
clean:
	rm -rf ${DIST_PATH}/*

.PHONY: update
update:
	go get