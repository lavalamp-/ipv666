 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

test:
	$(GOTEST) ./...

test-verbose:
	$(GOTEST) -v ./...

get-packr:
	$(GOGET) -u github.com/gobuffalo/packr/v2/packr2

build:
	$(GOBUILD) -v -o bin/ipv666 ipv666/main.go

build-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o bin/ipv666 ipv666/main.go
