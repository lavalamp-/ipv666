 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build-scan:
	$(GOBUILD) -v -o 666scan cmd/666scan/main.go

build-scan-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o 666scan cmd/666scan/main.go

build-blgen:
	$(GOBUILD) -v -o 666blgen cmd/666blgen/main.go

build-blgen-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o 666blgen cmd/666blgen/main.go