SHELL:=/bin/bash

export DOCKER_IMAGE:=$(or $(DOCKER_IMAGE), hellofront)
export DOCKER_VERSION:=$(or $(DOCKER_VERSION), latest)

.PHONY:
fmt: ## format Go code
	go fmt .

.PHONY: fmt-check
fmt-check: ## check if the code is formatted properly
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then echo "Formatting errors found! Please fix them by running: make fmt" ; exit 1; fi

.PHONY: run
run: build ## run the application
	./build/hellofront

.PHONY: test
test: ## run unit tests
	go test

.PHONY: build
build: ## build the executable
	go build -o build/hellofront

.PHONY: docker-build
docker-build: ## build the Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_VERSION) .

.PHONY: docker-push
docker-push: ## push the Docker image
	docker push $(DOCKER_IMAGE):$(DOCKER_VERSION)

.PHONY: docker-release
docker-release: docker-build docker-push ## build and publish the Docker image

.PHONY: load-test
load-test: ## run load tests usiong k6
	k6 run -e SVC_URL=$(or $(SVC_URL), http://localhost) tests/simple.js

.PHONY: clean
clean: ## remove artifacts
	rm -rf ./build
