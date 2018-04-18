# .PHONY: build doc fmt lint run test vet

GOPATH := ${PWD}
export GOPATH

# default: build

# all: test build

# build: vet
# 	go build -v -o ./bin/guten-snippet ./src/main

# doc:
# 	godoc -http=:6060 -index

# # http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
# fmt:
# 	go fmt main

# # https://github.com/golang/lint
# lint:
# 	./bin/golint main

# run: build
# 	./bin/guten-snippet

# test:
# 	go test -cover -v ./test/...

# vet:
# 	go vet main

# deps:
# 	go get github.com/golang/lint/golint
# 	go get github.com/haya14busa/goverage
# 	go get github.com/aws/aws-lambda-go/lambda

# zip:
# 	tar czvf guten-snippet.tar.gz --exclude=".DS_Store" Makefile readme.md ./src ./test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=snippet
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: vet
	$(GOBUILD) -o $(BINARY_NAME) -v
# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt main
# https://github.com/golang/lint
lint:
	./bin/golint main
vet:
	go vet main
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/golang/lint/golint
	$(GOGET) github.com/haya14busa/goverage
	$(GOGET) github.com/aws/aws-lambda-go/lambda

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
