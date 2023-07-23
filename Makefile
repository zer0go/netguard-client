default: help

.PHONY: build

version=`git describe --tags || echo "0.1.0"`

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' -e 's/:.*#/ #/'

install: ## Install the binary
	go install
	go install honnef.co/go/tools/cmd/staticcheck@latest

build: ## Build the application
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/ngclient main.go

build-all: ## Build application for supported architectures
	@echo "version: ${version}"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-x86_64 main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-aarch64 main.go

run: ## Run the application (eg: make run arg=show)
	@go run main.go ${arg}

lint: ## Check lint errors
	staticcheck ./...