# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOCOVERAGE=$(GOTEST) -cover -coverprofile=coverage.out
GOGET=$(GOCMD) get
CORE_BINARY_NAME=main
BINARY_NAME=main
SAM_OUTPUT=sam_output.yml
SAM_FILE=sam.yml

all: clean dependencies test build package

build: 
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CORE_BINARY_NAME) main.go
package:
	zip main.zip $(CORE_BINARY_NAME)
	aws cloudformation package --template-file $(SAM_FILE) --output-template-file $(SAM_OUTPUT) --s3-bucket $(DEPLOYMENT_BUCKET) 
test: 
	$(GOTEST) -v ./...
coverage: 
	$(GOCOVERAGE)
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME).zip
	rm -f $(SAM_OUTPUT)
deploy:
	aws cloudformation deploy --template-file $(SAM_OUTPUT) --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --parameter-overrides Environment=$(ENV) TableBaseName=$(TABLE_NAME)

dependencies: 
	@go get github.com/aws/aws-sdk-go/service/cloudformation
