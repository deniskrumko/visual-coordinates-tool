GOBIN ?= $$(go env GOPATH)/bin

run:
	go run main.go serve --config config/config.toml

run-noconfig:
	go run main.go serve

up:
	docker compose up --build

down:
	docker compose down --remove-orphans

tidy:
	go mod tidy

download-dependencies:
	go mod download

build: download-dependencies
	go build out/bin/main

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

# Uber Nilaway: https://github.com/uber-go/nilaway
nilaway:
	nilaway ./...

.PHONY: tests
tests:
	gotestsum ./...

check-coverage:
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

install-tools:
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go install gotest.tools/gotestsum@latest
