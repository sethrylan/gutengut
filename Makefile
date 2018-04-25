# Go parameters
# GOCMD=go
# GOBUILD=$(GOCMD) build
# GOCLEAN=$(GOCMD) clean
# GOTEST=$(GOCMD) test
# GOGET=$(GOCMD) get
BINARY_NAME=guten-snippet
BINARY_UNIX=$(BINARY_NAME)_unix
ARTIFACT_NAME=$(BINARY_NAME).zip
FUNCTION_NAME=guten-snippet
AWS_BUCKET_NAME=gutengut
AWS_STACK_NAME=gutengut
AWS_REGION=us-east-1

all: test build

build: vet lint
	go build -o ./dist/$(BINARY_NAME) -v
# 	@for dir in `ls handler`; do \
# 		GOOS=linux go build -o dist/handler/$$dir github.com/sbstjn/go-lambda-example/handler/$$dir; \
# 	done

clean:
	go clean
	@rm -f package.yml
	@rm -f guten-snippet.zip
	@rm -rf dist

###### AWS Deployment ######
configure:
	aws s3api create-bucket \
		--bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION)

package: build-lambda
	@aws cloudformation package \
		--template-file template.yml \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--output-template-file package.yml

deploy:
	# aws lambda update-function-code --zip-file=fileb://$(ARTIFACT_NAME) --function-name=$(FUNCTION_NAME)
	@aws cloudformation deploy \
		--template-file package.yml \
		--region $(AWS_REGION) \
		--capabilities CAPABILITY_IAM \
		--stack-name $(AWS_STACK_NAME)

describe:
		@aws cloudformation describe-stacks \
			--region $(AWS_REGION) \
			--stack-name $(AWS_STACK_NAME) \

###### Go tools ######
fmt:
	go fmt *.go
fmtchk:
	diff -u <(echo -n) <(go fmt -d ./)
lint:
	$(GOPATH)/bin/golint
vet:
	go vet -v ./...
test:
	go test -v ./... --cover
deps:
	go get github.com/golang/lint/golint
	go get -t ./...

# Cross compilation
build-lambda: clean vet lint
	GOOS=linux go build -o dist/$(BINARY_NAME)
	zip $(ARTIFACT_NAME) dist/*
build-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_UNIX) -v
docker-build: clean
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
