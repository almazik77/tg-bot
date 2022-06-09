GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=./bin/adapter
DOCKER_IMAGE_NAME=docker.technokratos.com/hackathon/back

all: build docker run_server
docker: docker_build docker_publish

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/app
	echo "binary build"

docker_build:
	docker build . -t $(DOCKER_IMAGE_NAME):$(CI_BUILD_REF_NAME)-$(CI_COMMIT_SHORT_SHA)

docker_publish:
	docker push $(DOCKER_IMAGE_NAME):$(CI_BUILD_REF_NAME)-$(CI_COMMIT_SHORT_SHA)

run_server:
	 LOG_LEVEL=info

