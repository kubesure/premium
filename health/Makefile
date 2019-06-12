GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
DOCKER=docker
DBUILD=$(DOCKER) build
DTAG= $(DOCKER) tag
DPUSH= $(DOCKER) push

BINARY_NAME=premium
BINARY_VERSION=v0.2
BINARY_UNIX=$(BINARY_NAME)
TAG_LOCAL = premiumcalc:$(BINARY_VERSION)
TAG_HUB = bikertales/premiumcalc:$(BINARY_VERSION)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

.PHONY: install
install:
	$(GOBUILD) -o $(BINARY_UNIX) -v ./...

.PHONY: run
run:
	$(GORUN) health.go

.PHONY: dbuild
dbuild:
	$(DBUILD) . -t $(TAG_LOCAL)

.PHONY: dtag
dtag: dbuild
	$(DTAG) $(TAG_LOCAL) $(TAG_HUB)

.PHONY: dpush
dpush: dtag
	$(DPUSH) $(TAG_HUB)
