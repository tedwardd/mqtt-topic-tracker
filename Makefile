BINARY_NAME=mqtt-topic-tracker

.PHONY: list deps build docker podman vendor lint clean
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

# IMPORTANT: The line above must be indented by (at least one) 
#            *actual TAB character* - *spaces* do *not* work.

deps:
	go mod download

build:
	$(MAKE) deps
	go build -o build/${BINARY_NAME} ./cmd/mqtt-topic-tracker

docker:
	docker build -t ${BINARY_NAME} .

podman:
	podman build -t ${BINARY_NAME} .

vendor:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run --enable-all

clean:
	go clean
	rm -rf build
