DHCPRELAY_IMAGE ?= dhcprelay:latest

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: go-imports
go-imports:
	go run golang.org/x/tools/cmd/goimports -w .

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: build-binary
build-binary:
	@CGO_ENABLED=0 go build -o bin/dhcprelay cmd/main.go

.PHONY: build-image
build-image:
	@docker build -t ${DHCPRELAY_IMAGE} .