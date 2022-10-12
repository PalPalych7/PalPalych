BIN := "./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_calendar:
	go build -v -o $(BIN)/calendar -ldflags "$(LDFLAGS)" ./cmd/calendar
run: build_calendar
	$(BIN)/calendar -config configs/calendar_config.toml

build_scheduler:
	go build -v -o $(BIN)/scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
run_scheduler: build_scheduler
	$(BIN)/scheduler -config ./configs/scheduler_config.toml

build_sender:
	go build -v -o $(BIN)/sender -ldflags "$(LDFLAGS)" ./cmd/sender
run_sender: build_sender
	$(BIN)/sender -config ./configs/sender_config.toml

build: build_calendar build_scheduler build_sender

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-rabMQ:
	docker run -d --name rb -p 15672:15672 -p 5672:5672 rabbitmq:3-management

version: build
	$(BIN) version

test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	CGO_ENABLED=0 golangci-lint run ./...

grpc_generate:
	protoc --proto_path=api --go_out=pb/ --go-grpc_out=pb/ api/*.proto

.PHONY: build run build-img run-img version test lint