# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=drone-release
VERSION=0.1.0

all: deps test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOMOD) download
docker-build:
	    docker build -f ./Dockerfile -t yerhy/drone-release:$(VERSION) .
docker-push:
        docker push yerhy/drone-release:$(VERSION); yerhy/drone-release:latest