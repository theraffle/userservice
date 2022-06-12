# Current  Version
VERSION ?= v0.0.1-alpha
REGISTRY ?= changjjjjjjjj

# Image URL to use all building/pushing image targets
IMG_USER_SERVICE ?= $(REGISTRY)/raffle-user-service:$(VERSION)

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build . -f Dockerfile -t ${IMG_USER_SERVICE}

# Push the docker image
.PHONY: docker-push
docker-push:
	docker push ${IMG_USER_SERVICE}

# Test code lint
test-lint:
	golint ./...
