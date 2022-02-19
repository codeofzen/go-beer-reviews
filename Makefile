GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet

CMD_FILES=./cmd/beer-reviews/main.go

BINARY_NAME=beer-review-service
VERSION?=0.0.0

SERVICE_PORT?=3000
DOCKER_REGISTRY?= #if set it should finished by /


.phony: all test build 

build: # build the project
	mkdir -p out/bin
	GO111MODULE=on $(GOCMD) build -mod vendor -o out/bin/$(BINARY_NAME) $(CMD_FILES)	


test: # run all unit and integration tests
	$(GOTEST) -v -race ./...