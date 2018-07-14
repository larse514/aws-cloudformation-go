# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOCOVERAGE=$(GOTEST) -cover -coverprofile=coverage.out
GOCOVERAGEOUT=coverage.out
GOGET=$(GOCMD) get

all: clean dependencies coverage

package:
	echo "unimplemnted"
test: 
	$(GOTEST) -v ./...
coverage: 
	$(GOCOVERAGE)
clean: 
	$(GOCLEAN)
	rm -f $(GOCOVERAGEOUT)
dependencies: 
	@go get github.com/aws/aws-sdk-go/service/cloudformation
