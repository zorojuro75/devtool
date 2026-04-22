# Build metadata
VERSION := $(shell git describe --tags --always --dirty 2>nul || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>nul || echo unknown)
DATE    := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'" 2>nul || echo unknown)

LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(DATE) -X main.commit=$(COMMIT)"
BINARY  := devtool
DIST    := dist

.PHONY: build test lint clean release install

build:
	go build $(LDFLAGS) -o bin\$(BINARY).exe .

install:
	go install $(LDFLAGS) .

test:
	go test -v -count=1 ./...

lint:
	go vet ./...

clean:
	if exist bin rd /s /q bin
	if exist dist rd /s /q dist
	if exist myapi rd /s /q myapi

release:
	@if not exist $(DIST) mkdir $(DIST)
	SET GOOS=linux&& SET GOARCH=amd64&& go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-amd64 .
	SET GOOS=linux&& SET GOARCH=arm64&& go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64 .
	SET GOOS=darwin&& SET GOARCH=amd64&& go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-amd64 .
	SET GOOS=darwin&& SET GOARCH=arm64&& go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-arm64 .
	SET GOOS=windows&& SET GOARCH=amd64&& go build $(LDFLAGS) -o $(DIST)/$(BINARY)-windows-amd64.exe .