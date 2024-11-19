LOGGER=mqtt-topic-tracker
FRONTEND=mqtt-topic-frontend
GIT_ROOT=$(shell git rev-parse --show-toplevel)

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
	cp ${GIT_ROOT}/docker/logger/Dockerfile ${GIT_ROOT}/Dockerfile
	docker build -t ${LOGGER} .
	rm ${GIT_ROOT}/Dockerfile
	cp ${GIT_ROOT}/docker/frontend/Dockerfile ${GIT_ROOT}/Dockerfile
	docker build -t ${FRONTEND} .
	rm ${GIT_ROOT}/Dockerfile

podman:
	cp ${GIT_ROOT}/docker/logger/Dockerfile ${GIT_ROOT}/Dockerfile
	podman build -t ${LOGGER} .
	rm ${GIT_ROOT}/Dockerfile
	cp ${GIT_ROOT}/docker/frontend/Dockerfile ${GIT_ROOT}/Dockerfile
	podman build -t ${FRONTEND} .
	rm ${GIT_ROOT}/Dockerfile

vendor:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run --enable-all

clean:
	go clean
	rm -rf build
