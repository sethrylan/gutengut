# Go parameters
# GOCMD=go
# GOBUILD=$(GOCMD) build
# GOCLEAN=$(GOCMD) clean
# GOTEST=$(GOCMD) test
# GOGET=$(GOCMD) get
BINARY_NAME=snippet
BINARY_UNIX=$(BINARY_NAME)_unix
ARTIFACT_NAME=$(BINARY_NAME).zip
FUNCTION_NAME=guten-snippet

all: test build

build: vet lint
	go build -o ./bin/$(BINARY_NAME) -v

deploy: build-lambda
	# aws cloudformation package --template-file template.yml --s3-bucket ${S3_BUCKET} --output-template-file packaged.yml
	aws lambda update-function-code --zip-file=fileb://$(ARTIFACT_NAME) --function-name=$(FUNCTION_NAME)

fmt:
	go fmt *.go
fmtchk:
	diff -u <(echo -n) <(go fmt -d ./)
lint:
	$(GOPATH)/bin/golint
vet:
	go vet -v ./...
test:
	go test -v ./...
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	./bin/$(BINARY_NAME)
deps:
	go get github.com/golang/lint/golint
	go get -t ./...

# Cross compilation
build-lambda: vet lint
	GOOS=linux go build -o $(BINARY_NAME)
	zip $(ARTIFACT_NAME) $(BINARY_NAME)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
