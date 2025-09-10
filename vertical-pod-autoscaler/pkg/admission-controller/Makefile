all: build

TAG?=dev
REGISTRY?=gcr.io/k8s-staging-autoscaling
FLAGS=
TEST_ENVVAR=LD_FLAGS=-s GO111MODULE=on
ENVVAR=CGO_ENABLED=0 $(TEST_ENVVAR)
GOOS?=linux
COMPONENT=admission-controller
FULL_COMPONENT=vpa-${COMPONENT}

BUILD_LD_FLAGS?=-s
ALL_ARCHITECTURES?=amd64 arm arm64 ppc64le s390x

# Buildx configuration
BUILDX_BUILDER?=multiarch-builder
PLATFORMS=$(subst $(space),$(comma),$(foreach arch,$(ALL_ARCHITECTURES),linux/$(arch)))
comma:=,
space:= $(empty) $(empty)

# Build targets
.PHONY: build
build: clean
	$(ENVVAR) GOOS=$(GOOS) go build ./...
	$(ENVVAR) GOOS=$(GOOS) go build -o ${COMPONENT}

.PHONY: build-binary
build-binary: clean
	$(ENVVAR) GOOS=$(GOOS) go build -o ${COMPONENT}

.PHONY: test-unit
test-unit: clean build
	$(TEST_ENVVAR) go test --test.short -race ./... $(FLAGS)

# Docker buildx setup
.PHONY: setup-buildx
setup-buildx:
	@if ! docker buildx ls | grep -q $(BUILDX_BUILDER); then \
		echo "Creating buildx builder $(BUILDX_BUILDER)..."; \
		docker buildx create --name $(BUILDX_BUILDER) --driver docker-container --use; \
		docker buildx inspect --bootstrap; \
	else \
		echo "Using existing buildx builder $(BUILDX_BUILDER)"; \
		docker buildx use $(BUILDX_BUILDER); \
	fi

.PHONY: cleanup-buildx
cleanup-buildx:
	@if docker buildx ls | grep -q $(BUILDX_BUILDER); then \
		echo "Removing buildx builder $(BUILDX_BUILDER)..."; \
		docker buildx rm $(BUILDX_BUILDER); \
	fi

# Multi-arch docker build
.PHONY: docker-build-multi-arch
docker-build-multi-arch: check-vars setup-buildx
	@echo "Building for platforms: $(PLATFORMS)"
	docker buildx build \
		--platform=$(PLATFORMS) \
		--build-arg BUILD_LD_FLAGS="${BUILD_LD_FLAGS}" \
		-t $(REGISTRY)/${FULL_COMPONENT}:${TAG} \
		-f ./Dockerfile \
		../../

# Multi-arch docker build and push in one step (most efficient)
.PHONY: docker-build-push-multi-arch
docker-build-push-multi-arch: check-vars setup-buildx
	@echo "Building and pushing for platforms: $(PLATFORMS)"
	docker buildx build \
		--platform=$(PLATFORMS) \
		--build-arg BUILD_LD_FLAGS="${BUILD_LD_FLAGS}" \
		--push \
		-t $(REGISTRY)/${FULL_COMPONENT}:${TAG} \
		-f ./Dockerfile \
		../../

# Build and load single architecture (for local testing)
.PHONY: docker-build-local
docker-build-local: check-vars setup-buildx
	@echo "Building for local architecture only"
	docker buildx build \
		--load \
		--build-arg BUILD_LD_FLAGS="${BUILD_LD_FLAGS}" \
		-t $(REGISTRY)/${FULL_COMPONENT}:${TAG} \
		-f ./Dockerfile \
		../../

# Legacy single-arch builds (if needed)
.PHONY: docker-build-single-arch
docker-build-single-arch: check-vars
	$(if $(ARCH),,$(error ARCH is undefined))
	@echo "Building for single architecture: $(ARCH)"
	docker buildx build \
		--platform=linux/$(ARCH) \
		--load \
		--build-arg BUILD_LD_FLAGS="${BUILD_LD_FLAGS}" \
		-t ${REGISTRY}/${FULL_COMPONENT}-$(ARCH):${TAG} \
		-f ./Dockerfile \
		../../

# Push pre-built images (if using docker-build-multi-arch without --push)
.PHONY: docker-push-multi-arch
docker-push-multi-arch: check-vars
	docker buildx imagetools create \
		--tag $(REGISTRY)/${FULL_COMPONENT}:${TAG} \
		$(REGISTRY)/${FULL_COMPONENT}:${TAG}

# Release process
.PHONY: release
release: docker-build-push-multi-arch
	@echo "Full multi-arch release ${FULL_COMPONENT}:${TAG} completed"

# Development release (build without push)
.PHONY: release-dev
release-dev: docker-build-multi-arch
	@echo "Multi-arch build ${FULL_COMPONENT}:${TAG} completed (not pushed)"

# Utility targets
.PHONY: check-vars
check-vars:
ifndef REGISTRY
	$(error REGISTRY is undefined)
endif
ifndef TAG
	$(error TAG is undefined)
endif

.PHONY: show-git-info
show-git-info:
	@echo '=============== local git status ==============='
	@git status
	@echo '=============== last commit ==============='
	@git log -1
	@echo '=============== bulding from the above ==============='

.PHONY: clean
clean:
	rm -f ${COMPONENT}

.PHONY: format
format:
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -d {} + | tee /dev/stderr)" || \
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -w {} + | tee /dev/stderr)"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build                      - Build Go binaries"
	@echo "  test-unit                  - Run unit tests"
	@echo "  docker-build               - Build Docker image using configured REGISTRY"
	@echo "  docker-build-local         - Build Docker image for localhost registry"
	@echo "  docker-build-multi-arch    - Build Docker images for all architectures"
	@echo "  docker-build-push-multi-arch - Build and push multi-arch images in one step"
	@echo "  docker-push                - Push image to registry (handles --insecure for localhost)"
	@echo "  docker-push-local          - Push image to localhost registry"
	@echo "  release                    - Full release (build and push multi-arch)"
	@echo "  release-dev                - Development build (multi-arch, no push)"
	@echo "  setup-buildx               - Setup Docker buildx builder"
	@echo "  cleanup-buildx             - Remove Docker buildx builder"
	@echo ""
	@echo "Variables:"
	@echo "  REGISTRY=$(REGISTRY)"
	@echo "  TAG=$(TAG)"
	@echo "  PLATFORMS=$(PLATFORMS)"
	@echo "  BUILDX_BUILDER=$(BUILDX_BUILDER)"