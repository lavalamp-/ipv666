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

get-deps:
	$(GOGET) -d ./...
	$(GOGET) -d github.com/stretchr/testify/assert

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

build-alias:
	$(GOBUILD) -v -o build/666alias cmd/666alias/main.go

build-alias-linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o build/666alias cmd/666alias/main.go

build-all:
	make build-scan
	make build-blgen
	make build-clean
	make build-alias

build-all-linux:
	make build-scan-linux
	make build-blgen-linux
	make build-clean-linux
	make build-alias-linux