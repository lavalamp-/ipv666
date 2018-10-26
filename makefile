 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build-scan:
	$(GOBUILD) -v -o build/666scan cmd/666scan/main.go

build-scan-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o build/666scan cmd/666scan/main.go

build-blgen:
	$(GOBUILD) -v -o build/666blgen cmd/666blgen/main.go

build-blgen-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o build/666blgen cmd/666blgen/main.go

build-clean:
	$(GOBUILD) -v -o build/666clean cmd/666clean/main.go

build-clean-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o build/666clean cmd/666clean/main.go

build-all:
	make build-scan
	make build-blgen
	make build-clean

build-all-linux:
	make build-scan-linux
	make build-blgen-linux
	make build-clean-linux