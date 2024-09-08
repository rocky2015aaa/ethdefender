# Define variables for reusability
DETECTOR_IMAGE_NAME = rocky2015aaa/ethdefender_detector
PREVENTER_IMAGE_NAME = rocky2015aaa/ethdefender_preventer
REPORTER_IMAGE_NAME = rocky2015aaa/ethdefender_reporter
DETECTOR_CONTAINER_NAME = ethdefender_detector
PREVENTER_CONTAINER_NAME = ethdefender_preventer
REPORTER_CONTAINER_NAME = ethdefender_reporter
PORT = 8080
VERSION := test
BUILD := test
DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')
CURRENT_DIR := $(shell pwd)
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
TESTER_NAME := ethtester

# Check if Docker image exists
detector_image_exists = $(shell docker images -q $(DETECTOR_IMAGE_NAME):latest)
preventer_image_exists = $(shell docker images -q $(PREVENTER_IMAGE_NAME):latest)
reporter_image_exists = $(shell docker images -q $(REPORTER_IMAGE_NAME):latest)

# Check if Docker container exists
detector_container_exists = $(shell docker ps -aq -f name=$(DETECTOR_CONTAINER_NAME))
preventer_container_exists = $(shell docker ps -aq -f name=$(PREVENTER_CONTAINER_NAME))
reporter_container_exists = $(shell docker ps -aq -f name=$(REPORTER_CONTAINER_NAME))

# Target to build the Docker image if it doesn't already exist
build: go-build-test
ifeq ($(detector_image_exists),)
	@echo "Building Docker image: $(DETECTOR_IMAGE_NAME):latest"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		-t $(DETECTOR_IMAGE_NAME):latest \
		-f $(CURRENT_DIR)/deployment/Dockerfile_detector \
    	$(CURRENT_DIR)
else
	@echo "Docker image $(DETECTOR_IMAGE_NAME):latest already exists."
endif

ifeq ($(preventer_image_exists),)
	@echo "Building Docker image: $(PREVENTER_IMAGE_NAME):latest"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		-t $(PREVENTER_IMAGE_NAME):latest \
		-f $(CURRENT_DIR)/deployment/Dockerfile_preventer \
    	$(CURRENT_DIR)
else
	@echo "Docker image $(PREVENTER_IMAGE_NAME):latest already exists."
endif

ifeq ($(reporter_image_exists),)
	@echo "Building Docker image: $(REPORTER_IMAGE_NAME):latest"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		-t $(REPORTER_IMAGE_NAME):latest \
		-f $(CURRENT_DIR)/deployment/Dockerfile_reporter \
    	$(CURRENT_DIR)
else
	@echo "Docker image $(REPORTER_IMAGE_NAME):latest already exists."
endif

up:
	docker-compose up -d

# Target to build and run only if it’s the first time (initial setup)
setup: build up

down:
	docker-compose down

clean-image:
	if [ -f $(GOBIN)/$(TESTER_NAME) ]; then rm $(GOBIN)/$(TESTER_NAME); fi

ifeq ($(detector_image_exists),)
	@echo "Docker image $(DETECTOR_IMAGE_NAME):latest does not exist."
else
	@echo "Removing Docker image: $(DETECTOR_IMAGE_NAME):latest"
	docker rmi $(DETECTOR_IMAGE_NAME):latest
endif

ifeq ($(preventer_image_exists),)
	@echo "Docker image $(PREVENTER_IMAGE_NAME):latest does not exist."
else
	@echo "Removing Docker image: $(PREVENTER_IMAGE_NAME):latest"
	docker rmi $(PREVENTER_IMAGE_NAME):latest
endif

ifeq ($(reporter_image_exists),)
	@echo "Docker image $(REPORTER_IMAGE_NAME):latest does not exist."
else
	@echo "Removing Docker image: $(REPORTER_IMAGE_NAME):latest"
	docker rmi $(REPORTER_IMAGE_NAME):latest
endif

# Target to stop and remove container and image
clean-all: down clean-image

# Target to rebuild and rerun everything with conditional checks
rebuild: clean-all go-build-test
	@if [ -z "$(shell docker images -q $(DETECTOR_IMAGE_NAME):latest)" ] || \
		[ -z "$(shell docker images -q $(PREVENTER_IMAGE_NAME):latest)" ] || \
		[ -z "$(shell docker images -q $(REPORTER_IMAGE_NAME):latest)" ]; then \
		make build; \
	fi
	@if [ -z "$(shell docker ps -aq -f name=$(DETECTOR_CONTAINER_NAME))" ] || \
		[ -z "$(shell docker ps -aq -f name=$(PREVENTER_CONTAINER_NAME))" ] || \
		[ -z "$(shell docker ps -aq -f name=$(REPORTER_CONTAINER_NAME))" ]; then \
		make up; \
	fi


go-build-test:
	@echo "  >  Building tx test binary..."
	GOBIN=$(GOBIN) go build -o $(GOBIN)/$(TESTER_NAME) ./test

# .PHONY prevents targets from being mistaken for files
.PHONY: build up down clean-image clean-all rebuild go-build-test