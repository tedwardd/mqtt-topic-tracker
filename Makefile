LOGGER=mqtt-topic-tracker
FRONTEND=mqtt-topic-frontend

.PHONY: list deps build docker podman vendor lint clean
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

# IMPORTANT: The line above must be indented by (at least one) 
#            *actual TAB character* - *spaces* do *not* work.

deps:
	go mod download

logger:
	$(MAKE) deps
	go build -o build/${LOGGER} ./cmd/mqtt-topic-tracker

frontend:
	$(MAKE) deps
	go build -o build/${FRONTEND} ./cmd/frontend

build:
	$(MAKE) deps
	go build -o build/${LOGGER} ./cmd/mqtt-topic-tracker
	go build -o build/${FRONTEND} ./cmd/frontend

docker:
	docker build -t ${LOGGER} docker/logger
	docker build -t ${FRONTEND} docker/frontend

podman:
	podman build -t ${LOGGER} docker/logger
	podman build -t ${FRONTEND} docker/frontend

vendor:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run --enable-all

clean:
	go clean
	rm -rf build
