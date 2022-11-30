BIN := "./bin"
DOCKER_IMG="calendar:develop"
DOCKER_IMG_SENDER="calendar_sender:develop"
DOCKER_IMG_SCHEDULER="calendar_scheduler:develop"
DOCKER_IMG_INT_TESTS="calend_integration_test:develop"

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

build-img_main:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/mainDockerfile .

build-img_sender:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/senderDockerfile .

build-img_scheduler:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/schedulerDockerfile .

build-img_int_tests:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_INT_TESTS) \
		-f integration_tests/Dockerfile .

build-compose: build-img_main build-img_scheduler build-img_sender build-img_int_tests
	docker-compose build    

run_compose: build-compose
	docker-compose up -d postgres_db 
	docker-compose up -d mainSevice 
	docker-compose up -d scheduler
	docker-compose up -d scheduler
	docker-compose up -d sender

up: run_compose

down: build-compose down

integration_tests: run_compose
	docker-compose up integraton_test


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

.PHONY: build run build-img run-img version test lint build-img_main build-img_sender build-img_scheduler